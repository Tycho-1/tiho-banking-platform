# Upstream — Bank of Anthos

## Attribution

Application code and initial manifests are derived from:

**[GoogleCloudPlatform/bank-of-anthos](https://github.com/GoogleCloudPlatform/bank-of-anthos)** — Apache License 2.0

Formal attribution: see **[NOTICE](../NOTICE)** in the repo root (required for Apache-2.0 redistribution).

## Version pin

| Item | Value |
|------|--------|
| Upstream tag | `v0.6.9` (adjust when you sync) |
| Container registry | `us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/` |
| Last manual sync | 2026-07-06 |
| README screenshots | `docs/img/{login,transactions,architecture}.png` — copied from upstream `docs/img/` (Apache-2.0) |

## What to sync from upstream

```bash
git remote add upstream https://github.com/GoogleCloudPlatform/bank-of-anthos.git
git fetch upstream
git diff HEAD upstream/main -- src/
```

| Path | Sync strategy |
|------|----------------|
| `src/` | Merge or cherry-pick on new upstream **tags** |
| `deploy/base/` | Diff against `upstream/kubernetes-manifests/`; port changes manually |
| `deploy/overlays/` | **Ours** — do not overwrite from upstream |

## What intentionally diverges from upstream

- Central `deploy/` instead of `kubernetes-manifests/` + `src/**/k8s/`
- No per-service `cloudbuild.yaml` / `skaffold.yaml` in `src/`
- Kind / GitOps / Istio overlays (non-GKE telemetry disabled)
- No `iac/`, Cloud Deploy, or Workload Identity from upstream
