# Kind local — ambient + CNPG + ESO

Same stack as [`kind-local-ambient-cnpg`](../kind-local-ambient-cnpg/), plus Vault-backed secrets:

| Component | Role |
|-----------|------|
| [`use-eso-vault`](../../components/use-eso-vault/) | App Secrets (`accounts-db-secrets`, `ledger-db-secrets`, `jwt-key`, `demo-data-secrets`) |
| [`use-eso-vault-cnpg`](../../components/use-eso-vault-cnpg/) | `cnpg-connection-secrets` + CNPG bootstrap creds in `cnpg-banking` |

`SecretStore` + auth: **platform** (both `banking-platform` and `cnpg-banking`).

## Prerequisites

- Everything for [`kind-local-ambient-cnpg`](../kind-local-ambient-cnpg/README.md)
- Platform: ESO, Vault, `SecretStore` `vault-banking-platform` in both namespaces
- Vault paths — phase 1 + CNPG (see component READMEs); `db/accounts` + `db/ledger` must use **CNPG hostnames**

## Install

```bash
kubectl apply -k deploy/overlays/kind-local-ambient-cnpg-eso

kubectl get secretstore,externalsecret -n banking-platform
kubectl get secretstore,externalsecret -n cnpg-banking
kubectl wait --for=condition=Ready externalsecret --all -n banking-platform --timeout=120s
kubectl wait --for=condition=Ready externalsecret --all -n cnpg-banking --timeout=120s
```

Then wait for CNPG clusters / app Deployments as in the parent overlay README.

## Notes

- ConfigMap `cnpg-connection-targets` (JDBC URL without password) stays in git via the parent overlay.
- Do not mix with `kind-local-ambient-cnpg` (plain Secrets) in the same namespaces at once.
- This overlay does **not** set Kustomize `namespace:` — CNPG is in `cnpg-banking`, app in `banking-platform`; ESO CRs set namespaces explicitly.
