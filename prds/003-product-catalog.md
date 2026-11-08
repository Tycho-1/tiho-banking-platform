# PRD 003 — product-catalog

**Status:** done — Go service, frontend page, image, and Kind deploy are on `main`.

Spec: [docs/product-catalog-requirements.md](../docs/product-catalog-requirements.md). Prometheus scrape and OTLP for this service were closed later in [PRD 001](001-microservices-prom-otel.md).

## Problem

The app can show balances and payments. It cannot show what the bank offers. Add a read-only catalog (savings, CDs, investments) and a page to browse it. No ledger writes and no new database.

## Scope

- New Go service `product-catalog` in `src/products/`, port `8080`.
- `GET /api/products` and `GET /api/products?type=savings|cd|investment`. JWT required. `/health`, `/ready`, and `/version` stay open.
- Static JSON catalog, about 5–10 products. Loaded at startup.
- Frontend page `/products`, nav label **Save & invest**, only when the user is logged in. Frontend calls the service with the bearer token. Config key `PRODUCTS_API_ADDR`.
- Deploy like the other app services: base Deployment and Service, `jwt-key` public key mounted, image pin on `kind-local` and `gke-dev`, GHCR via `use-ghcr-images`, CI workflow from `release.yaml`.
- Kind telemetry flags off through `disable-gcp-telemetry`, same as the other services. Metrics and OTLP are a later slice, once the service exists.

## Out of scope

- Buying, signup, or moving money
- Loans, recommendations, Postgres, CNPG schema changes
- Changes to userservice, contacts, or the ledger services
- A new ingress. The browser talks to the frontend. The frontend talks to product-catalog inside the cluster.
- Keycloak or an Istio AuthorizationPolicy. App JWT stays the check.
- `gke-dev` Cloud Ops, and a dedicated Grafana board for this service alone

## Acceptance

- [x] `GET /api/products` without a token returns 401
- [x] A valid JWT returns the static list. `?type=savings` returns only that type
- [x] Logged-in user can open **Save & invest** and see the cards. Logged-out user is sent to login
- [x] `go test ./...` passes in `src/products`
- [x] `kubectl kustomize` succeeds for `kind-local` and `gke-dev`
- [x] CI publishes `ghcr.io/tycho-1/tiho-banking-platform/product-catalog` from `src/products/release.yaml`
- [x] Flux on Kind runs that image and the page loads in the cluster

## Tasks

- [x] Service: `main.go`, JWT verify, catalog JSON, Dockerfile, `src/products/README.md`
- [x] Tests for auth failure, list, and type filter
- [x] `deploy/base/product-catalog.yaml` and `PRODUCTS_API_ADDR` in the frontend config
- [x] Frontend route, template, and nav link
- [x] Image pin in `use-ghcr-images`, `kind-local`, and `gke-dev`. Kind inherit overlays pick up the `kind-local` pin
- [x] `.github/workflows/product-catalog.yml` and `src/products/release.yaml`
- [x] `disable-gcp-telemetry` patch so Kind does not push GCP telemetry from this Deployment
- [x] Proved with tests, a kustomize build, and the page on Kind after the image deployed

## Follow-ups

Closed in PRD 001: `/metrics`, OTLP traces, and OTLP logs on Kind.
