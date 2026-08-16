# product-catalog

Read-only Go microservice for bank product offerings (savings, CDs, investments).

## Endpoints

| Path | Auth | Description |
|------|------|-------------|
| `GET /health` | No | Liveness |
| `GET /ready` | No | Readiness |
| `GET /version` | No | Build version |
| `GET /metrics` | No | Prometheus metrics (for scraping) |
| `GET /api/products` | JWT | List products |
| `GET /api/products?type=savings\|cd\|investment` | JWT | Filter by type |

## Local run

```bash
cd src/products
go test ./...
go vet ./...

# Demo JWT public key from repo (same as kind-local demo-jwt component)
grep 'jwtRS256.key.pub:' ../../deploy/components/demo-jwt/jwt-secret.yaml \
  | awk '{print $2}' | base64 -d > /tmp/jwt.pub

export PUB_KEY_PATH=/tmp/jwt.pub
go run .
```

If your cluster is up (non-ESO overlay with demo-jwt):

```bash
kubectl get secret jwt-key -n banking-platform \
  -o jsonpath='{.data.jwtRS256.key\.pub}' | base64 -d > /tmp/jwt.pub
```

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/products
curl http://localhost:8080/metrics
```

## Prometheus metrics (`/metrics`)

This service exposes [Prometheus](https://prometheus.io/) metrics on **`GET /metrics`** (no JWT — scrapers must reach the pod/network, same pattern as `/health`).

### How it works (mental model)

```text
Browser / user     →  /api/products     (business API, JWT required)
Prometheus server  →  /metrics          (telemetry, scrape every N seconds)
Grafana            →  queries Prometheus  (dashboards, alerts)
```

Prometheus **pulls** metrics from your app; the app does not push to Grafana.

### What we expose

| Metric | Type | Labels | Meaning |
|--------|------|--------|---------|
| `product_catalog_http_requests_total` | Counter | `method`, `path`, `status` | Request count per route (excludes `/metrics` scrapes) |
| `product_catalog_http_request_duration_seconds` | Histogram | `method`, `path` | Latency per route |
| `product_catalog_products_requests_total` | Counter | `type` | Successful product API calls (`all`, `savings`, `cd`, …) |
| `product_catalog_catalog_products_loaded` | Gauge | — | Products loaded from JSON at startup |
| `go_*`, `process_*` | — | — | Go runtime + process metrics (via collectors) |

Implementation: `metrics.go` — `prometheus/client_golang`, HTTP middleware wraps all routes.

### Local check

After `go run .`:

```bash
curl -s http://localhost:8080/health
curl -s http://localhost:8080/metrics | grep product_catalog_http
```

Example line after a few requests:

```text
product_catalog_http_requests_total{method="GET",path="/health",status="200"} 3
```

### Kubernetes / Kind (next step)

The endpoint alone is not enough — your **platform** Prometheus must be configured to **scrape** `product-catalog:8080/metrics` (Pod annotation or `ServiceMonitor`). Until then, metrics exist on the pod but won't appear in Grafana.

## Product data

Edit `data/products.json` to add or update demo products. Rebuild the image after changes.

## CI

`.github/workflows/product-catalog.yml` runs `go test`, `go vet`, and a Docker build on
changes under `src/products/**`. Publishing to GHCR requires the repository variable
`ENABLE_IMAGE_PUSH=true`; the image tag is read from `release.yaml`.
