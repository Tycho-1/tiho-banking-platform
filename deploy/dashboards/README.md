# Banking platform metrics dashboard

`banking-platform-metrics.json` is a Grafana dashboard for the metrics this app's microservices expose. It is not a cluster or Kubernetes dashboard. Panels query the series added with the observability work in this repo.

Import it by hand: **Dashboards → New → Import → Upload dashboard JSON**, then choose the Prometheus datasource that scrapes the Kind ServiceMonitors. Nothing in the Kustomize tree loads this file.

TODO: wrap this JSON in a ConfigMap that the Kind Grafana sidecar watches, and include that manifest from the Kind overlay so Flux loads the dashboard.

## Services

| Service | Scrape path | What the panels use |
|---------|-------------|---------------------|
| `frontend` | `/metrics` | `frontend_http_*` and `frontend_backend_*` from `src/frontend/metrics.py` |
| `contacts` | `/metrics` | `contacts_http_*` from `src/accounts/contacts/metrics.py` |
| `userservice` | `/metrics` | `userservice_http_*` from `src/accounts/userservice/metrics.py` |
| `product-catalog` | `/metrics` | `product_catalog_*` from `src/products/metrics.go` |
| `balancereader`, `transactionhistory`, `ledgerwriter` | `/actuator/prometheus` | Spring `http_server_requests_seconds_*`, `jvm_memory_used_bytes`, `process_cpu_usage`, and the Guava `cache_*` meters |

The **job** variable matches either `ledgerwriter` or `namespace/ledgerwriter`, and the same pattern for the other services above.

Traces stay in Jaeger. Pod logs are in `banking-platform-logs.json`.

## Logs dashboard

`banking-platform-logs.json` shows pod logs for the same seven app services. Import it the same way and pick the Loki datasource (`Loki` or `webstore-logs`).

The queries use `namespace="banking-platform"` and `service_name`. Promtail already writes those labels, so the panels can show lines before the OTel log images are deployed. OTel export is a second copy of the application lines and is not what these panels count.
