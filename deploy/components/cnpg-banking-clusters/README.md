# cnpg-banking-clusters

Deploys two dedicated CNPG `Cluster` resources for Bank of Anthos:

- `banking-accounts` — schema + demo users (from upstream `initdb/`)
- `banking-ledger` — transactions schema

Demo credentials in `credentials-secrets.yaml` (plain Kubernetes Secrets). Swap for External Secrets later.

Omit this component when reusing an existing CNPG cluster; point apps via overlay `cnpg-connection-targets.yaml` instead.
