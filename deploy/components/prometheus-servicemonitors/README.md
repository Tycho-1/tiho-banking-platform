# Prometheus ServiceMonitors (Kubernetes)

Tells **Prometheus Operator** to scrape app metrics. Kubernetes pattern: ServiceMonitor → Service → pod. Grafana queries Prometheus.

**`gke-dev` does not include this component** (GCP Cloud Ops). Include it on any overlay whose cluster has Prometheus Operator + Grafana.

**Not PodMonitor.** These services already have a Service with named port `http`. `ServiceMonitor` follows that Service’s Endpoints. Use `PodMonitor` when there is no Service (CNPG controller is an example).

| Service | Path |
|---------|------|
| product-catalog | `/metrics` |
| balancereader | `/actuator/prometheus` (Spring Actuator) |

## Prerequisites

- Prometheus Operator CRDs (`ServiceMonitor`)
- Prometheus that selects monitors labeled `release: kube-prometheus-stack` (change the label if your Prometheus Helm release name differs)

## Included from (this repo)

- [`kind-local`](../../overlays/kind-local/) (inherited by ambient / CNPG / ESO overlays)
- [`kind-local-gateway-api`](../../overlays/kind-local-gateway-api/)

## Grafana (after scrape)

```promql
product_catalog_catalog_products_loaded
{__name__=~"product_catalog_.*"}
sum by (path, status) (rate(product_catalog_http_requests_total[5m]))

# Java / Spring (balancereader) — JVM + HTTP; Guava cache as cache_*
jvm_memory_used_bytes{job="balancereader"}
http_server_requests_seconds_count{job="balancereader"}
```
