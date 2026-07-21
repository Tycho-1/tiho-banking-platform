# Use GHCR images (self-built)

Optional Kustomize component. Remaps upstream BoA image names to this repo’s **GHCR** images (built by CI).

**Enabled by default** on overlay **`kind-local`**. Other overlays still use upstream Google Artifact Registry unless you add this component.

## What it does

Rewrites each service image to this repo’s GHCR path. **Tags are per-service** (from each `src/.../release.yaml`):

```text
frontend        → ghcr.io/tycho-1/tiho-banking-platform/frontend:v0.6.10
userservice     → …/userservice:v0.6.9
…
```

Only remap services you have actually published to GHCR when adding this component to **other** overlays. On **`kind-local`**, all nine app images are expected on GHCR after CI.

## Default on `kind-local`

Already listed in [overlays/kind-local/kustomization.yaml](../../overlays/kind-local/kustomization.yaml). To switch back to **upstream** BoA images, remove this component from that overlay — see [overlays/kind-local/README.md](../../overlays/kind-local/README.md).

## Per-service versioning

Each microservice owns its semver in `release.yaml`:

```yaml
# src/frontend/release.yaml
version: "0.6.10"
```

CI reads `$SERVICE_PATH/release.yaml` and pushes `v$(version)` (+ git SHA). Bumping only frontend does **not** change other services’ tags.

Keep `newTag` in this file in sync with each service’s `release.yaml` when you promote.

## Enable on other overlays

1. Bump `version` in the service’s `release.yaml`, enable `ENABLE_IMAGE_PUSH=true`, push to `main` under that service path.
2. Set matching `newTag` here for that service.
3. Add the component to the overlay:

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
