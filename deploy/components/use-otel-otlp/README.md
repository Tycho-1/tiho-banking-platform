# OTLP traces to the platform Collector (Kubernetes)

Turns tracing **on** for **frontend**, **product-catalog**, and **balancereader**, pointing at the platform OpenTelemetry Collector over **OTLP HTTP**.

`disable-gcp-telemetry` sets `ENABLE_TRACING=false` (no legacy GCP trace push on local overlays). This component runs **after** that and sets `ENABLE_TRACING=true` plus `OTEL_EXPORTER_OTLP_ENDPOINT`.

**Not on `gke-dev` yet.** That overlay does not set `OTEL_EXPORTER_OTLP_ENDPOINT`. Frontend still falls back to **Cloud Trace** when `ENABLE_TRACING=true`. Product-catalog and balancereader need an OTLP endpoint to export (balancereader uses Micrometer OTel only — see [balancereader README](../../../src/ledger/balancereader/README.md#tracing)).

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
