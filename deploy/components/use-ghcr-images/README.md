# Use GHCR images (self-built)

Optional Kustomize component. Default overlays keep **upstream** Bank of Anthos images so demos work without CI push.

## What it does

Rewrites each service image to this repo’s GHCR path. **Tags are per-service** (from each `src/.../release.yaml`):

```text
frontend        → ghcr.io/tycho-1/tiho-banking-platform/frontend:v0.6.10
userservice     → …/userservice:v0.6.9
…
```

Only opt in for services you have actually published to GHCR. Until then, leave this component off (or override a single image in the overlay).

## Per-service versioning

Each microservice owns its semver in `release.yaml`:

```yaml
# src/frontend/release.yaml
version: "0.6.10"
```

CI reads `$SERVICE_PATH/release.yaml` and pushes `v$(version)` (+ git SHA). Bumping only frontend does **not** change other services’ tags.

Keep `newTag` in this file in sync with each service’s `release.yaml` when you promote.

## Opt in

1. Bump `version` in the service’s `release.yaml`, enable `ENABLE_IMAGE_PUSH=true`, push to `main` under that service path.
2. Set matching `newTag` here for that service.
3. Add the component (or pin only one image in the overlay):

```yaml
components:
  - ../../components/use-ghcr-images
```

Frontend-only smoke test without remapping everything:

```yaml
images:
  - name: us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/frontend
    newName: ghcr.io/tycho-1/tiho-banking-platform/frontend
    newTag: v0.6.10
```

## CNPG note

On `kind-local-ambient-cnpg`, embedded `accounts-db` / `ledger-db` Deployments are removed — remapping those two image names is harmless leftover. App services still use GHCR.

## If the owner/repo differs

Change `newName` prefixes in `kustomization.yaml` to match your GHCR path.
