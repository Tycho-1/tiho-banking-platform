# Kind local — ESO + Vault

Same stack as [`kind-local`](../kind-local/) (sidecar, Istio Gateway/VS, GHCR images), but **demo Secrets are not taken from git**. Component [`use-eso-vault`](../../components/use-eso-vault/) deletes them and syncs the same Secret names from Vault.

| Piece | Location |
|-------|----------|
| Base overlay | Reuses `../kind-local` |
| ExternalSecrets | `deploy/components/use-eso-vault/` |
| SecretStore + Vault auth | **Platform** — `vault-banking-platform` in `banking-platform` |

## Prerequisites

- Everything for [`kind-local`](../kind-local/README.md) (Kind, Istio sidecar, ingress, GHCR pulls)
- **ESO** + Vault + **`SecretStore` `vault-banking-platform`** (and auth) from the **platform** repo
- Vault phase-1 paths populated — see [use-eso-vault README](../../components/use-eso-vault/README.md)

## Install

```bash
kubectl apply -k deploy/overlays/kind-local-eso

kubectl get secretstore,externalsecret -n banking-platform
kubectl wait --for=condition=Ready externalsecret --all -n banking-platform --timeout=120s

kubectl wait --for=condition=available deployment --all -n banking-platform --timeout=300s
```

Demo login (unchanged seed): `testuser` / `bankofanthos`.

## vs `kind-local`

| | `kind-local` | **`kind-local-eso`** |
|--|--------------|----------------------|
| Secrets | Plain YAML in git | Vault → ESO → K8s Secrets |
| JWT | `demo-jwt` component | Same Secret name `jwt-key`, from Vault |
| One-shot without Vault | Yes | No — needs platform ESO + Vault |

Do not mix both overlays in the same namespace at once.

## CNPG

For ambient + CNPG + ESO, use [`kind-local-ambient-cnpg-eso`](../kind-local-ambient-cnpg-eso/).
