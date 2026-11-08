---
description: After metrics/traces work in Grafana or Jaeger — review coverage, gaps, and TODOs. Use when deploy proof is already OK.
---

# Observability — review (post-deploy)

Run **after** user confirms Grafana metrics and/or Jaeger traces look OK.

## Steps

1. Ask which services were in scope and what user saw in Grafana/Jaeger.

2. **Metrics review**
   - The seven app services already have ServiceMonitors. Do not list them as missing.
   - Are default JVM/HTTP metrics enough, or is a business counter still missing (ledgerwriter accepted / rejected / duplicate)?
   - PromQL smoke queries from `deploy/components/prometheus-servicemonitors/README.md`

3. **Traces and logs review**
   - Cross-service propagation (frontend → backend)?
   - Python and Java OTLP logs can carry the trace id. Go `log.Printf` does not.
   - Logs dashboard counts the Promtail stream (`namespace="banking-platform"`), not the OTel copy.

4. **Open PRD**
   - [002](../../../prds/002-grafana-dashboards.md) is the draft for loading dashboards through the sidecar. Do not treat a hand import as that PRD being done.
   - 001 and 003 are done. Do not reopen them.

5. Output:
   - **Coverage summary** (good / partial / missing)
   - **TODO list** — follow-ups as new PRs, not extra scope in the current change
   - Suggest `observability-prd-create` only for a new slice (alerts, ledger result counter, Go trace ids, an EKS overlay)

6. Append TODOs to `docs/TODO.md` or the active PRD if the user agrees.

## Do not

- Mark observability "complete" without user confirmation they checked UI
- One-shot perfect observability — prefer iterative PRs
