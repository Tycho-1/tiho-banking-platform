# GitHub Actions — Tiho Banking Platform

Path-filtered CI per microservice under `src/`. **Build and test run without registry credentials.**

## Versioning (per-service semver)

Each service has a small `release.yaml` next to its code:

```text
src/frontend/release.yaml
src/accounts/userservice/release.yaml
…
```

Example:

```yaml
# Service release metadata — CI reads version for GHCR image tags (vX.Y.Z).
version: "0.6.10"
```

Publish jobs read **`$SERVICE_PATH/release.yaml`** → `version` and tag:

```text
ghcr.io/<owner>/<repo>/<service>:v0.6.10   # from that service’s release.yaml
ghcr.io/<owner>/<repo>/<service>:<git-sha> # immutable build id
```

No root version file and no `:latest` tag. Services can diverge (frontend `0.6.10`, userservice `0.6.9`). Add more keys to `release.yaml` later if needed — keep **one** file, don’t add sibling `VERSION`/`OWNER` files.

**Release one service**

1. Bump `version` in `src/<service>/release.yaml`.
2. Commit with your code change; push to `main` (`ENABLE_IMAGE_PUSH=true`).
3. Pin that service in GitOps (`use-ghcr-images` or overlay `images:`) to `v0.6.10`.

Upstream BoA images remain `v0.6.9` until you opt into GHCR for a given service.

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
| `deploy-validate.yml` | `deploy/` | — | `kubectl kustomize` all five overlays |

## Image push (optional — off by default)

Push jobs run only when **all** of:

- `github.ref == refs/heads/main`
- event is `push` (not PR)
- repository variable **`ENABLE_IMAGE_PUSH`** = `true`

### Enable GHCR push

1. Repo → **Settings → Secrets and variables → Actions → Variables**
2. Add variable: `ENABLE_IMAGE_PUSH` = `true`
3. Prefer **Settings → Actions → General → Workflow permissions** = Read (publish jobs already set `packages: write`)

Images push to (owner/repo lowercased for GHCR):

```text
ghcr.io/tycho-1/tiho-banking-platform/<service>:v<release.yaml-version>
ghcr.io/tycho-1/tiho-banking-platform/<service>:<git-sha>
```

Workflows set `IMAGE` with `tr '[:upper:]' '[:lower:]'` because `github.repository` keeps the GitHub username casing (`Tycho-1`) and GHCR rejects uppercase.
See [deploy/components/use-ghcr-images/README.md](../../deploy/components/use-ghcr-images/README.md). Pin `newTag` per service to match that service’s `release.yaml`.

### GCP Artifact Registry (later)

Replace the publish job login with `google-github-actions/auth` + `gcloud auth configure-docker` and set `IMAGE_REGISTRY` variable (see workflow comments).
