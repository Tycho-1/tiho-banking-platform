# TODO — Tiho Banking Platform

## Done so far

- [x] **Kustomize deploy layout** — `deploy/base` + reusable `components` + Kind/GKE overlays.
- [x] **Kind overlays** — sidecar (`kind-local`), Gateway API ingress (`kind-local-gateway-api`), Istio ambient (`kind-local-ambient`).
- [x] **CNPG** — `kind-local-ambient-cnpg` (clusters + seed Jobs for accounts + ledger demo data); apps pointed at CNPG via `use-cnpg-databases`.
- [x] **GHCR images** — `use-ghcr-images` default on Kind overlays; per-service `release.yaml` + CI workflows (push gated by `ENABLE_IMAGE_PUSH`).
- [x] **External Secrets + Vault** — components `use-eso-vault` / `use-eso-vault-cnpg`; overlays `kind-local-eso` and `kind-local-ambient-cnpg-eso`. App ships ExternalSecrets only; platform owns `SecretStore` (`vault-banking-platform`) + Vault KV paths.
- [x] **Deploy CI** — `deploy-validate` kustomize-builds all overlays (including ESO); `checkout@v5` / `setup-kubectl@v5`.
- [x] **Public app repo** — banking app on GitHub with `src/`, `deploy/`, workflows, docs.
- [ ] **Product catalog v1** — Go `product-catalog` + frontend `/products` (local API tested; Kind/browser pending) — [product-catalog-requirements.md](product-catalog-requirements.md)

## README & docs

- [ ] **Replace README screenshots** with your own captures (Kind + Istio ambient, CNPG overlay, Gateway API UI). Keep upstream images in `docs/upstream.md` or `docs/img/upstream/` with attribution, or note “screenshots from local deploy”.

## Deploy & platform

- [ ] Re-test **`kind-local-ambient-cnpg`** / **`-eso`** (seed Jobs + waits) when cluster has headroom (scale down `loadgenerator` if OOM).
- [ ] **Harden Vault auth** — move SecretStores off long-lived `vault-token` to Kubernetes auth (or AppRole); platform-side.
- [ ] **Flux** from a platform GitOps repo → chosen overlay (e.g. `kind-local-ambient-cnpg-eso` or `kind-local-gateway-api`).
- [ ] **Gateway API without Istio** — install a non-Istio controller on Kind (e.g. Envoy Gateway); set `gatewayClassName` accordingly; optional: drop `istio-injection` on that overlay.
- [ ] **Do not** make Helm the primary deploy path — Kustomize stays source of truth. Optional thin Helm wrapper only if a real consumer asks for `helm install`.

## Observability (**hard / high value**)

Upstream BoA assumes GCP Cloud Ops (`ENABLE_METRICS` / `ENABLE_TRACING`). On Kind you currently disable GCP telemetry. Replacing that properly is non-trivial (app wiring + platform stack + dashboards).

- [ ] **OpenTelemetry** — export traces/metrics/logs from services to an OTEL Collector (replace or complement GCP exporters); document env/flags per language (Python / Java / frontend).
- [ ] **Prometheus stack (ref)** — platform installs Prometheus (+ Grafana, Alertmanager as needed); app exposes scrape targets / ServiceMonitors (or OTEL→Prometheus bridge). Prefer **platform-owned** stack; app repo only adds scrape/instrumentation hooks + optional overlay notes.
- [ ] Optional: **Istio/ambient telemetry** into the same backend (mesh metrics + app OTEL) so one Grafana story covers both.
- [ ] Docs: Kind prerequisites + how to open Grafana / find BoA dashboards; keep Cloud Ops path for `gke-dev` if useful.

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
