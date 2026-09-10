# Copyright 2026 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Prometheus metrics for the frontend Flask app."""

import time

import requests
from flask import Response, g, request
from prometheus_client import CONTENT_TYPE_LATEST, Counter, Histogram, generate_latest
from requests.exceptions import RequestException

# Inbound: frontend_http_requests_total, frontend_http_request_duration_seconds
_HTTP_REQUESTS = Counter(
    "frontend_http_requests_total",
    "Total HTTP requests handled by the frontend.",
    ["method", "endpoint", "status"],
)
_HTTP_DURATION = Histogram(
    "frontend_http_request_duration_seconds",
    "HTTP request latency in seconds.",
    ["method", "endpoint"],
    buckets=(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0),
)

# Outbound BFF calls: frontend_backend_requests_total, frontend_backend_request_duration_seconds
_BACKEND_REQUESTS = Counter(
    "frontend_backend_requests_total",
    "Total outbound HTTP requests from the frontend to backend services.",
    ["service", "status"],
)
_BACKEND_DURATION = Histogram(
    "frontend_backend_request_duration_seconds",
    "Outbound backend HTTP request latency in seconds.",
    ["service"],
    buckets=(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0),
)

# ApiCall display_name → trace/Jaeger-style service name
_BACKEND_SERVICE_NAMES = {
    "balance": "balancereader",
    "transaction_list": "transactionhistory",
    "contacts": "contacts",
}

_EXCLUDED_PATHS = frozenset({"/metrics", "/ready", "/version"})


def resolve_backend_service(display_name):
    """Map ApiCall display names to backend service labels."""
    return _BACKEND_SERVICE_NAMES.get(display_name, display_name)


def record_backend_call(service, status, duration_seconds):
    """Record one outbound backend HTTP call."""
    status_label = str(status)
    _BACKEND_REQUESTS.labels(service, status_label).inc()
    _BACKEND_DURATION.labels(service).observe(duration_seconds)


def backend_get(service, **kwargs):
    """requests.get wrapper that records backend metrics."""
    return _backend_request(service, requests.get, **kwargs)


def backend_post(service, **kwargs):
    """requests.post wrapper that records backend metrics."""
    return _backend_request(service, requests.post, **kwargs)


def _backend_request(service, method, **kwargs):
    start = time.perf_counter()
    try:
        response = method(**kwargs)
        record_backend_call(
            service, response.status_code, time.perf_counter() - start
        )
        return response
    except RequestException:
        record_backend_call(service, "error", time.perf_counter() - start)
        raise


def _endpoint_label():
    """Stable route name for labels; avoids high cardinality from path params."""
    return request.endpoint or request.path


def init_metrics(app):
    """Register /metrics and request counters on the Flask app."""

    @app.before_request
    def _metrics_start_timer():
        if request.path in _EXCLUDED_PATHS:
            return
        g._metrics_start_time = time.perf_counter()

    @app.after_request
    def _metrics_record(response):
        if request.path in _EXCLUDED_PATHS:
            return response
        start = getattr(g, "_metrics_start_time", None)
        if start is None:
            return response
        endpoint = _endpoint_label()
        status = str(response.status_code)
        _HTTP_REQUESTS.labels(request.method, endpoint, status).inc()
        _HTTP_DURATION.labels(request.method, endpoint).observe(
            time.perf_counter() - start
        )
        return response

    @app.route("/metrics", methods=["GET"])
    def metrics():
        return Response(generate_latest(), mimetype=CONTENT_TYPE_LATEST)
