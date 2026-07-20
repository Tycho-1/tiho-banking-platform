# TODO — Tiho Banking Platform

## README & docs

- [ ] **Replace README screenshots** with your own captures (Kind + Istio ambient, CNPG overlay, Gateway API UI). Keep upstream images in `docs/upstream.md` or `docs/img/upstream/` with attribution, or note “screenshots from local deploy”.

## Deploy & platform

- [ ] Re-test **`kind-local-ambient-cnpg`** (seed Job + waits) when cluster has headroom (scale down `loadgenerator` if OOM).
- [ ] **External Secrets Operator** for CNPG bootstrap + JWT (replace plain Secrets in git with Vault sync).
- [ ] **Flux** from a platform GitOps repo → chosen overlay (e.g. `kind-local-ambient-cnpg` or `kind-local-gateway-api`)
- [ ] **Gateway API without Istio** — install a non-Istio controller on Kind (e.g. Envoy Gateway); set `gatewayClassName` accordingly; optional: drop `istio-injection` on that overlay.

## Public Kind / platform reference (**important**)

This app repo documents **prerequisites** only; cluster bootstrap is out of tree and currently **private**. For a public portfolio, employers need somewhere to see how the Kind platform is built.

- [ ] **Publish a public platform reference repo** (new name OK — sanitized Kind bootstrap: Istio and/or Gateway API controller, MetalLB, CNPG operator, optional Vault/observability).
- [ ] Link it from root `README.md` + `deploy/README.md` (replace the “cluster you create / private bootstrap” wording).
- [ ] Keep personal/experimental Kind work private if needed; the **public** repo is the interview-facing reference.

Until then: banking app can still be public; docs stay prerequisite-based (no broken private GitHub links).

## CI / images

- [ ] Enable **`ENABLE_IMAGE_PUSH`** + GHCR when ready; update `deploy/overlays/*/kustomization.yaml` `images[].name` to your registry.
- [ ] Optional: auto-bump `newTag` in overlay after image push.

## New public repo (this banking app)

- [ ] Create/publish GitHub repo; push `src/`, `deploy/`, `.github/workflows/`, `docs/img/`, `LICENSE`, `NOTICE`.
- [ ] Revisit **Jib vs Dockerfile/buildpacks** for Java services if you want a different portfolio story (not required).
