# Use ESO + Vault for CNPG secrets (phase 2)

Kustomize **component**. Use **with** [`use-eso-vault`](../use-eso-vault/) on a CNPG overlay.

Removes git-backed CNPG Secrets and replaces them with `ExternalSecret` CRs.

## Ownership

Same split as phase 1: **platform** provides `SecretStore` **`vault-banking-platform`** (and auth) in **`banking-platform`** and **`cnpg-banking`**. This component only ships `ExternalSecret`s.

## What phase 1 vs phase 2 covers

| Secrets | Component | Used by |
|---------|-----------|---------|
| `accounts-db-secrets`, `ledger-db-secrets`, `jwt-key`, `demo-data-secrets` | **`use-eso-vault`** | App (+ STS or CNPG via Vault values) |
| `cnpg-connection-secrets` | **this component** | Overlay connection material (`banking-platform`) |
| `banking-accounts-credentials`, `banking-ledger-credentials` | **this component** | CNPG `Cluster` bootstrap + seed Jobs (`cnpg-banking`) |

## Vault paths (`platform-kv`)

KV v2 mount **`platform-kv`**. Extra folder layout for this component (phase 2 — CNPG / `kind-local-ambient-cnpg-eso`):

```text
platform-kv/
└── tiho-banking-platform/
    ├── db/ … auth/ … demo/   # phase 1 — still required (see use-eso-vault)
    └── cnpg/
        ├── accounts-creds    # → banking-accounts-credentials (ns cnpg-banking)
        ├── ledger-creds      # → banking-ledger-credentials (ns cnpg-banking)
        └── connection        # → cnpg-connection-secrets (ns banking-platform)
```

| Vault path | K8s Secret | Namespace | Keys |
|------------|------------|-----------|------|
| `tiho-banking-platform/cnpg/connection` | `cnpg-connection-secrets` | `banking-platform` | `ACCOUNTS_DB_URI`, `POSTGRES_PASSWORD`, `SPRING_DATASOURCE_PASSWORD`, `POSTGRES_PASSWORD_LEDGER` |
| `tiho-banking-platform/cnpg/accounts-creds` | `banking-accounts-credentials` | `cnpg-banking` | `username`, `password` |
| `tiho-banking-platform/cnpg/ledger-creds` | `banking-ledger-credentials` | `cnpg-banking` | `username`, `password` |

Also set **CNPG hostnames** in phase-1 paths `db/accounts` and `db/ledger` (app Secrets). Do not rely on kustomize Secret replacements when ESO owns those Secrets.

## Opt in

```yaml
components:
  - ../../components/use-eso-vault
  - ../../components/use-eso-vault-cnpg
```

ConfigMap `cnpg-connection-targets` (non-secret JDBC URL) stays in the overlay.

## Verify

```bash
kubectl get secretstore,externalsecret -n banking-platform
kubectl get secretstore,externalsecret -n cnpg-banking
```
