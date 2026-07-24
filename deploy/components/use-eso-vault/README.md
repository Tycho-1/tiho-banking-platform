# Use ESO + Vault (phase 1)

Kustomize **component**. Removes git-backed demo Secrets and replaces them with `ExternalSecret` CRs that recreate the **same** Kubernetes Secret names and keys from Vault.

App Deployments are unchanged — they still `secretKeyRef` / mount Secrets.

## Ownership

| Piece | Where it lives |
|-------|----------------|
| ESO + Vault | **Platform** |
| `SecretStore` **`vault-banking-platform`** (+ auth, e.g. `vault-token`) | **Platform** (in `banking-platform`) |
| **`ExternalSecret` CRs** (this component) | **App** |

This repo does **not** ship `SecretStore` manifests. ExternalSecrets reference:

```yaml
secretStoreRef:
  name: vault-banking-platform
  kind: SecretStore
```

### Platform example (not applied from this repo)

```yaml
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: vault-banking-platform
  namespace: banking-platform
spec:
  provider:
    vault:
      server: http://vault-active.vault.svc:8200
      path: platform-kv
      version: v2
      auth:
        tokenSecretRef:
          name: vault-token
          key: token
```

Prefer Vault Kubernetes auth / AppRole instead of a long-lived token when you harden the platform.

## Prerequisites

| Piece | Expected |
|-------|----------|
| ESO installed | CRDs `ExternalSecret`, `SecretStore` |
| Namespace | **`banking-platform`** |
| Store + auth | Provided by platform (see above) |
| Vault secrets | Paths below under mount `platform-kv` |

## Vault paths (`platform-kv`)

KV v2 mount **`platform-kv`**. Folder layout for this component (phase 1 — STS / `kind-local-eso`):

```text
platform-kv/
└── tiho-banking-platform/
    ├── db/
    │   ├── accounts          # → K8s Secret accounts-db-secrets
    │   └── ledger            # → K8s Secret ledger-db-secrets
    ├── auth/
    │   └── jwt               # → K8s Secret jwt-key
    └── demo/
        └── login             # → K8s Secret demo-data-secrets
```

| Vault path | K8s Secret | Keys |
|------------|------------|------|
| `tiho-banking-platform/db/accounts` | `accounts-db-secrets` | `ACCOUNTS_DB_URI`, `POSTGRES_PASSWORD` |
| `tiho-banking-platform/db/ledger` | `ledger-db-secrets` | `POSTGRES_PASSWORD`, `SPRING_DATASOURCE_PASSWORD` |
| `tiho-banking-platform/auth/jwt` | `jwt-key` | `jwtRS256.key`, `jwtRS256.key.pub` (full PEM with BEGIN/END) |
| `tiho-banking-platform/demo/login` | `demo-data-secrets` | `DEMO_LOGIN_PASSWORD` |

Login password must stay `bankofanthos` unless you change the seed hash. After rotating `jwt-key`, restart workloads that mount it.

For CNPG extras under `tiho-banking-platform/cnpg/`, see [`use-eso-vault-cnpg`](../use-eso-vault-cnpg/).

## Opt in

```yaml
components:
  - ../../components/use-eso-vault
```

When enabled: plain Secrets from `deploy/base/` and `demo-jwt` are **deleted** from the build; ESO owns the same names.

## Verify

```bash
kubectl get secretstore,externalsecret -n banking-platform
kubectl get secret accounts-db-secrets ledger-db-secrets jwt-key demo-data-secrets -n banking-platform
```

`SecretStore` should be **Ready** (platform); each `ExternalSecret` **SecretSynced**.

## CNPG

See [`use-eso-vault-cnpg`](../use-eso-vault-cnpg/).
