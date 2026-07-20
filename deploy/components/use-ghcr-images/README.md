# Use GHCR images (self-built)

Optional Kustomize component. Default overlays keep **upstream** Bank of Anthos images so demos work without CI push.

## What it does

Rewrites each service image:

```text
us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/<service>
  → ghcr.io/tycho-1/tiho-banking-platform/<service>:latest
```

Names match `.github/workflows/*` publish tags (`ghcr.io/<owner>/<repo>/<service>`), lowercased for GHCR.

## Opt in

1. Enable CI push: repo variable `ENABLE_IMAGE_PUSH=true`, push to `main`.
2. Add the component to an overlay, e.g. `deploy/overlays/kind-local/kustomization.yaml`:

```yaml
components:
  - ../../components/demo-jwt
  # ...
  - ../../components/use-ghcr-images
```

3. Apply as usual:

```bash
kubectl apply -k deploy/overlays/kind-local
```

## Promote a specific build

Edit `newTag` in this component (or add overlay `images:` keyed by the **GHCR** name):

```yaml
images:
  - name: ghcr.io/tycho-1/tiho-banking-platform/frontend
    newTag: abcdef0123...   # git sha from CI
```

## CNPG note

On `kind-local-ambient-cnpg`, embedded `accounts-db` / `ledger-db` Deployments are removed — remapping those two image names is harmless leftover. App services still use GHCR.

## If the owner/repo differs

Change `newName` prefixes in `kustomization.yaml` to match your GHCR path.
