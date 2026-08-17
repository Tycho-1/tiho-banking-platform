# OTLP traces to the platform Collector (Kind)

Turns tracing **on** for **frontend** and **product-catalog**, pointing at the Kind OpenTelemetry Collector over **OTLP HTTP**.

`disable-gcp-telemetry` sets `ENABLE_TRACING=false` (no Cloud Trace on Kind). This component runs **after** that and sets `ENABLE_TRACING=true` plus `OTEL_EXPORTER_OTLP_ENDPOINT`.

**Not on `gke-dev`.** That overlay does not set `OTEL_EXPORTER_OTLP_ENDPOINT`. Frontend then uses **Cloud Trace** (when `ENABLE_TRACING=true`). Product-catalog has no Cloud Trace exporter — Go tracing stays off without an OTLP endpoint.

## Prerequisites (platform)

- OpenTelemetry Collector in `observability` (OTLP HTTP `:4318`)
- A trace backend the Collector exports to (Tempo or Jaeger) + Grafana datasource

Default endpoint:

```text
http://otel-collector.observability.svc:4318
```

Change the patches if your Collector DNS differs.

## Enabled on

- [`kind-local`](../../overlays/kind-local/) (inherited by ambient / CNPG / ESO overlays)
- [`kind-local-gateway-api`](../../overlays/kind-local-gateway-api/)
