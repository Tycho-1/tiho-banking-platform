# TODO — Tiho Banking Platform

## Done so far

- [x] **Kustomize deploy layout** — `deploy/base` + reusable `components` + Kind/GKE overlays.
- [x] **Kind overlays** — sidecar (`kind-local`), Gateway API ingress (`kind-local-gateway-api`), Istio ambient (`kind-local-ambient`).
- [x] **CNPG** — `kind-local-ambient-cnpg` (clusters + seed Jobs for accounts + ledger demo data); apps pointed at CNPG via `use-cnpg-databases`.
- [x] **GHCR images** — `use-ghcr-images` default on Kind overlays; per-service `release.yaml` + CI workflows (push gated by `ENABLE_IMAGE_PUSH`).
- [x] **External Secrets + Vault** — components `use-eso-vault` / `use-eso-vault-cnpg`; overlays `kind-local-eso` and `kind-local-ambient-cnpg-eso`. App ships ExternalSecrets only; platform owns `SecretStore` (`vault-banking-platform`) + Vault KV paths.
- [x] **Deploy CI** — `deploy-validate` kustomize-builds all overlays (including ESO); `checkout@v5` / `setup-kubectl@v5`.
- [x] **Public app repo** — banking app on GitHub with `src/`, `deploy/`, workflows, docs.
- [x] **Product catalog v1** — Go `product-catalog` + frontend `/products` on Kind. PRD: [prds/003-product-catalog.md](../prds/003-product-catalog.md). Spec: [product-catalog-requirements.md](product-catalog-requirements.md)

## README & docs

- [ ] **Replace README screenshots** with your own captures (Kind + Istio ambient, CNPG overlay, Gateway API UI). Keep upstream images in `docs/upstream.md` or `docs/img/upstream/` with attribution, or note “screenshots from local deploy”.

## Deploy & platform

- [ ] Re-test **`kind-local-ambient-cnpg`** / **`-eso`** (seed Jobs + waits) when cluster has headroom (scale down `loadgenerator` if OOM).
- [ ] **Harden Vault auth** — move SecretStores off long-lived `vault-token` to Kubernetes auth (or AppRole); platform-side.
- [ ] **Flux** from a platform GitOps repo → chosen overlay (e.g. `kind-local-ambient-cnpg-eso` or `kind-local-gateway-api`).
- [ ] **Gateway API without Istio** — install a non-Istio controller on Kind (e.g. Envoy Gateway); set `gatewayClassName` accordingly; optional: drop `istio-injection` on that overlay.
- [ ] **Do not** make Helm the primary deploy path — Kustomize stays source of truth. Optional thin Helm wrapper only if a real consumer asks for `helm install`.

## Observability

Upstream BoA assumes GCP Cloud Ops. Kind turns that off (`disable-gcp-telemetry`) and uses Prometheus scrape plus OTLP. `gke-dev` stays on Cloud Ops.

- [x] **Seven app services on Kind** — Prometheus ServiceMonitors, OTLP traces, and OTLP logs (`frontend`, `contacts`, `userservice`, `balancereader`, `transactionhistory`, `ledgerwriter`, `product-catalog`). Metrics stay a scrape, not an OTLP metrics export. Agent map: [agent-observability.md](agent-observability.md).
- [x] **Dashboards** — hand-imported JSON in `deploy/dashboards/` (metrics, trace rates, pod logs). Jaeger remains the waterfall.
- [ ] **Alerts** — scrape down, HTTP 5xx, p95. No `PrometheusRule` yet.
- [ ] **Ledger business counter** — accepted / rejected / duplicate. HTTP status hides that.
- [ ] **Go log trace ids** — Python and Java OTLP logs can carry the trace id. `product-catalog` `log.Printf` does not.
- [ ] **Dashboard GitOps** — [PRD 002](../prds/002-grafana-dashboards.md). Flux applies the JSON from `flux-fleet` into namespace `observability`, not from `kind-local`.
- [ ] Optional: **Istio/ambient telemetry** in the same Grafana story.

## Public Kind / platform reference (**important**)

This app repo documents **prerequisites** only; cluster bootstrap is out of tree and currently **private**. For a public portfolio, employers need somewhere to see how the Kind platform is built.

- [ ] **Publish a public platform reference repo** (new name OK — sanitized Kind bootstrap: Istio and/or Gateway API controller, MetalLB, CNPG operator, Vault/ESO, optional observability).
- [ ] Link it from root `README.md` + `deploy/README.md` (replace the “cluster you create / private bootstrap” wording).
- [ ] Keep personal/experimental Kind work private if needed; the **public** repo is the interview-facing reference.

Until then: banking app can still be public; docs stay prerequisite-based (no broken private GitHub links).

## CI / images

- [ ] Confirm **`ENABLE_IMAGE_PUSH=true`** on the GitHub repo when GHCR packages should rebuild on `main` (workflows already wired).
- [ ] Optional: auto-bump `newTag` in overlay / `release.yaml` after image push.

## Nice-to-have

- [ ] Revisit **Jib vs Dockerfile/buildpacks** for Java services if you want a different portfolio story (not required).
- [ ] Later overlays: `eks-dev`, `aks-dev` — same pattern as `gke-dev` / Kind.
