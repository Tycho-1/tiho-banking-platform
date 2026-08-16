# Prometheus ServiceMonitors (Kind)

Tells the **platform** Prometheus (kube-prometheus-stack) to scrape app `/metrics` endpoints.

**Not PodMonitor.** Product-catalog already has a Service with named port `http`. `ServiceMonitor` follows that Service’s Endpoints. Use `PodMonitor` when there is no Service (CNPG controller is an example).

## Prerequisites

- Prometheus Operator CRDs (`ServiceMonitor`)
- Prometheus that selects monitors labeled `release: kube-prometheus-stack` (Helm default on this Kind cluster)

**`gke-dev` does not include this component** — that overlay keeps GCP Cloud Ops.

## Enabled on

- [`kind-local`](../../overlays/kind-local/) (inherited by ambient / CNPG / ESO overlays)
- [`kind-local-gateway-api`](../../overlays/kind-local-gateway-api/)

## Grafana (after scrape)

```promql
product_catalog_catalog_products_loaded
{__name__=~"product_catalog_.*"}
sum by (path, status) (rate(product_catalog_http_requests_total[5m]))
```
