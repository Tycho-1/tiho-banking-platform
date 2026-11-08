# PRD 001 — Prometheus and OpenTelemetry for the banking microservices (Kind)

**Status:** done — on `main` for the seven app services.

The application code ships in every image. This PRD is the Kind wiring: Kustomize on `kind-local` points that code at Prometheus and the OpenTelemetry Collector. A GKE backend (Cloud Ops, or the same Collector on GKE) is a new PRD.

Metrics are a Prometheus scrape. Traces and logs go to the collector over OTLP. This is not an OTLP metrics export.

## Problem

Upstream Bank of Anthos sends telemetry to GCP. Kind needs the same application signals in Prometheus, Jaeger, and Loki. `gke-dev` stays on Cloud Ops.

## Scope

Services: `frontend`, `contacts`, `userservice`, `balancereader`, `transactionhistory`, `ledgerwriter`, `product-catalog`.

Overlay: `kind-local` (inherited by the ambient, CNPG, and gateway overlays).

| Signal | How |
|--------|-----|
| Metrics | ServiceMonitor. Python and Go: `/metrics`. Java: `/actuator/prometheus` |
| Traces | `use-otel-otlp` → `http://otel-collector.observability.svc:4318`. Java also sets `MANAGEMENT_OTLP_TRACING_ENDPOINT` with `/v1/traces` |
| Logs | OTLP on that same endpoint. Stdout stays, so Promtail still fills the logs dashboard |

## Application vs overlay

The seven services always expose their metrics endpoint, and they export OTLP traces and logs only when the endpoint environment variables are set. On Kind, `use-otel-otlp` sets those variables and `prometheus-servicemonitors` scrapes the endpoint. On `gke-dev` those components are absent, so the same image keeps stdout logs and the existing Cloud Trace path.

## Out of scope

- `gke-dev` and Cloud Ops — same binaries, different overlay. A GKE PRD would choose the backend and the Kustomize component. It would not re-instrument the seven services.
- `loadgenerator`, `accounts-db`, `ledger-db`
- Alerts, a ledger result counter, Go trace ids on log lines, dashboard GitOps

## Acceptance

- [x] One ServiceMonitor per service, registered in `prometheus-servicemonitors`
- [x] One OTLP patch per service, registered in `use-otel-otlp`
- [x] App code exports the signal: Python `_setup_tracing` / `_setup_otlp_logs`, Java Actuator + `OtlpLogConfig`, Go `metrics.go` / `tracing.go` / `logs.go`
- [x] `gke-dev` does not include those two components
- [x] Dashboards exist for metrics, trace rates, and pod logs (`deploy/dashboards/`)
- [x] Kind runtime was checked per language while the images rolled out (Grafana, Jaeger, Loki `OTLP logs enabled`)

## Follow-ups

These are new work, not open checkboxes on this PRD:

- Alerts on scrape down, HTTP 5xx, and p95
- Ledgerwriter counter for accepted, rejected, and duplicate
- Trace id on product-catalog log lines
- Load the dashboard JSON from Flux
