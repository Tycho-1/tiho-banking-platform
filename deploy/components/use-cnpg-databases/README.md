# use-cnpg-databases

Removes embedded `accounts-db` / `ledger-db` StatefulSets (and their Services). Apps keep using ConfigMaps + Secrets for DB config; the overlay supplies CNPG endpoints via Kustomize **replacements**.

**Requires** in the overlay:

- ConfigMap `cnpg-connection-targets` (non-secret URLs, e.g. `SPRING_DATASOURCE_URL`)
- Secret `cnpg-connection-secrets` (URIs/passwords)
- `replacements:` in `kustomization.yaml` into `accounts-db-secrets` / `ledger-db-secrets` / `ledger-db-config`

See `overlays/kind-local-ambient-cnpg/`.

Does **not** deploy CNPG clusters — pair with `cnpg-banking-clusters` or your own Cluster CRs.
