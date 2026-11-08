# PRDs

A PRD is for a multi-step initiative: a new service, a new overlay, or a signal that crosses code and Kustomize. A one-file tweak does not need one. There is not one PRD per microservice.

| Skill | Action |
|-------|--------|
| `observability-prd-create` | New file here. Next number is the highest `NNN` plus one |
| `observability-prd-done` | Close when acceptance is met |

## Files

| PRD | Status |
|-----|--------|
| [001-microservices-prom-otel.md](001-microservices-prom-otel.md) | done — Kind Prometheus scrape, OTLP traces, and OTLP logs for the seven app services |
| [002-grafana-dashboards.md](002-grafana-dashboards.md) | draft — load the three existing dashboards through the Kind Grafana sidecar |
| [003-product-catalog.md](003-product-catalog.md) | done — Go catalog service, Save & invest page, image, and Kind deploy |

Later slices (not files yet): alerts, a ledger result counter, Go trace ids on logs. The next PRD number is **004**.
