---
description: Add OpenTelemetry OTLP tracing for one service on kind-local via use-otel-otlp patches. Use after or with Prom work.
---

# Observability — add OpenTelemetry

Add **OTLP traces and logs** for **one service** that does not have them yet. The seven app services already export both on Kind. Use component `deploy/components/use-otel-otlp/`.

## Before coding

Confirm:

- **Service name**. Do not redo frontend, contacts, userservice, balancereader, transactionhistory, ledgerwriter, or product-catalog.
- **Reference:** `balancereader` (Java traces + `OtlpLogConfig`), `frontend` (Python `_setup_tracing` / `_setup_otlp_logs`), or `product-catalog` (`tracing.go` / `logs.go`)
- Overlay includes `use-otel-otlp` (default on `kind-local`)

Collector endpoint (Kind): `http://otel-collector.observability.svc:4318`

## Steps

1. Read reference patch + service Deployment in `deploy/base/`.

2. Create or extend `deploy/components/use-otel-otlp/patch-<service>.yaml`:
   - `ENABLE_TRACING=true`
   - `OTEL_SERVICE_NAME`
   - `OTEL_EXPORTER_OTLP_ENDPOINT` / protocol
   - Java: `MANAGEMENT_OTLP_TRACING_ENDPOINT`, disable GCP trace if needed
   - Disable Cloud Trace where `disable-gcp-telemetry` applies

3. Add patch to `deploy/components/use-otel-otlp/kustomization.yaml`.

4. Application code: traces and logs, matching the reference. Python logs turn on only when tracing is on and the endpoint is set. Java and Go logs turn on from the endpoint alone. Stdout stays.

5. Run `bash scripts/agent/observability-verify.sh`.

6. **User must verify in Jaeger, Grafana, and Loki** after deploy. Look for `OTLP logs enabled` on the new service. A diff is not proof.

## Do not

- Enable on `gke-dev` without user approval (Cloud Trace default there)
- Skip deploy + trace UI verification
