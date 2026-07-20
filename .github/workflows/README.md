# GitHub Actions — Tiho Banking Platform

Path-filtered CI per microservice under `src/`. **Build and test run without registry credentials.**

## Workflows

| Workflow | Path | Test | Build |
|----------|------|------|-------|
| `frontend.yml` | `src/frontend/` | — | Docker |
| `userservice.yml` | `src/accounts/userservice/` | pytest (uv) | Docker |
| `contacts.yml` | `src/accounts/contacts/` | pytest (uv) | Docker |
| `loadgenerator.yml` | `src/loadgenerator/` | — | Docker |
| `accounts-db.yml` | `src/accounts/accounts-db/` | — | Docker |
| `ledger-db.yml` | `src/ledger/ledger-db/` | — | Docker |
| `balancereader.yml` | `src/ledger/balancereader/` | Maven test (Java 17) | Jib → Docker |
| `ledgerwriter.yml` | `src/ledger/ledgerwriter/` | Maven test (Java 17) | Jib → Docker |
| `transactionhistory.yml` | `src/ledger/transactionhistory/` | Maven test (Java 17) | Jib → Docker |
| `deploy-validate.yml` | `deploy/` | — | `kubectl kustomize` all five overlays (`kind-local`, `kind-local-gateway-api`, `kind-local-ambient`, `kind-local-ambient-cnpg`, `gke-dev`) |

## Image push (optional — off by default)

Push jobs run only when **all** of:

- `github.ref == refs/heads/main`
- event is `push` (not PR)
- repository variable **`ENABLE_IMAGE_PUSH`** = `true`

### Enable GHCR push

1. Repo → **Settings → Secrets and variables → Actions → Variables**
2. Add variable: `ENABLE_IMAGE_PUSH` = `true`
3. Repo → **Settings → Actions → General → Workflow permissions** → read/write for packages

Images push to:

```text
ghcr.io/<owner>/tiho-banking-platform/<service>:<git-sha>
ghcr.io/<owner>/tiho-banking-platform/<service>:latest
```

Update `deploy/overlays/*/kustomization.yaml` (or add component **`deploy/components/use-ghcr-images`**) when switching from upstream BoA images — see [deploy/components/use-ghcr-images/README.md](../../deploy/components/use-ghcr-images/README.md).

### GCP Artifact Registry (later)

Replace the publish job login with `google-github-actions/auth` + `gcloud auth configure-docker` and set `IMAGE_REGISTRY` variable (see workflow comments).
