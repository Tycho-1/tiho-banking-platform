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

"""Prometheus metrics for the contacts Flask app."""

import time

from flask import Response, g, request
from prometheus_client import CONTENT_TYPE_LATEST, Counter, Histogram, generate_latest

# Inbound: contacts_http_requests_total, contacts_http_request_duration_seconds
_HTTP_REQUESTS = Counter(
    "contacts_http_requests_total",
    "Total HTTP requests handled by the contacts service.",
    ["method", "endpoint", "status"],
)
_HTTP_DURATION = Histogram(
    "contacts_http_request_duration_seconds",
    "HTTP request latency in seconds.",
    ["method", "endpoint"],
    buckets=(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0),
)

_EXCLUDED_PATHS = frozenset({"/metrics", "/ready", "/version"})


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
