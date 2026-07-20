# Kind ambient + CNPG

Bank of Anthos with **CloudNativePG** instead of embedded `accounts-db` / `ledger-db` StatefulSets.

| Piece | Location |
|-------|----------|
| App overlay | Reuses `../kind-local-ambient` |
| CNPG clusters (optional) | `deploy/components/cnpg-banking-clusters/` |
| App → DB wiring | `deploy/components/use-cnpg-databases/` |
| **Your cluster endpoints** | **`cnpg-connection-targets.yaml`** (this overlay) |

## How you specify the CNPG cluster

Edit **`cnpg-connection-targets.yaml`** in this overlay (ConfigMap for non-secret URLs + Secret for passwords/URIs). That file is the single source of truth for where apps connect.

CNPG exposes a read-write Service per cluster:

```text
<cluster-name>-rw.<namespace>.svc.cluster.local:5432
```

**Bundled defaults** (when `cnpg-banking-clusters` component is enabled):

| Database | CNPG Cluster | RW host |
|----------|--------------|---------|
| Accounts | `banking-accounts` in `cnpg-banking` | `banking-accounts-rw.cnpg-banking.svc.cluster.local` |
| Ledger | `banking-ledger` in `cnpg-banking` | `banking-ledger-rw.cnpg-banking.svc.cluster.local` |

**Using an existing CNPG cluster** (any namespace — not only the bundled `cnpg-banking` clusters):

1. Remove `cnpg-banking-clusters` from `kustomization.yaml` `components:`.
2. Create databases/users on that cluster (manual SQL or CNPG `Database` CR).
3. Update `cnpg-connection-targets.yaml` hosts and credentials.
4. Keep `use-cnpg-databases` + `replacements` as-is.

Example — edit **both** the ConfigMap URL and the Secret URI/passwords:

```yaml
# ConfigMap cnpg-connection-targets
data:
  SPRING_DATASOURCE_URL: jdbc:postgresql://my-cluster-rw.my-ns.svc.cluster.local:5432/postgresdb

# Secret cnpg-connection-secrets
stringData:
  ACCOUNTS_DB_URI: postgresql://accounts-admin:demo-password-change-me@my-cluster-rw.my-ns.svc.cluster.local:5432/accounts-db
  POSTGRES_PASSWORD: demo-password-change-me
  SPRING_DATASOURCE_PASSWORD: demo-password-change-me
  POSTGRES_PASSWORD_LEDGER: demo-password-change-me
```

## Prerequisites

- CNPG operator installed on the cluster
- Istio **ambient** on the Kind cluster
- If switching from embedded DBs: delete old `accounts-db` / `ledger-db` StatefulSets or use a fresh namespace

## Install

```bash
kubectl apply -k deploy/overlays/kind-local-ambient-cnpg

kubectl wait --for=condition=Ready cluster/banking-accounts cluster/banking-ledger -n cnpg-banking --timeout=300s

# Demo users — seed Job (also covers bootstrap race / missed postInit SQL)
kubectl wait --for=condition=complete job/banking-accounts-seed -n cnpg-banking --timeout=120s

# Demo ledger transactions (embedded ledger-db loads these via USE_DEMO_DATA; CNPG needs this Job)
kubectl wait --for=condition=complete job/banking-ledger-seed -n cnpg-banking --timeout=120s

kubectl rollout restart deployment -n banking-platform
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
```

### Why test data might be missing after one `kubectl apply -k`

CNPG **`postInitApplicationSQLRefs` runs only once** at first cluster bootstrap (empty PVC). It does **not** re-run on later applies.

Common causes:

1. **Cluster already existed** from an earlier attempt (bootstrap skipped).
2. **Apply race** — Cluster CR reconciled before ConfigMap `banking-accounts-init-sql` existed; bootstrap may skip or fail init SQL silently.
3. **Bootstrap error** — check `kubectl describe cluster banking-accounts -n cnpg-banking`.

The **`banking-accounts-seed` Job** re-applies `schema.sql` + `testdata.sql` idempotently (`ON CONFLICT DO NOTHING` on inserts). Wait for it after every deploy.

Verify users:

```bash
kubectl exec -n cnpg-banking banking-accounts-1 -- \
  psql -U accounts-admin -d accounts-db -c 'SELECT username FROM users;'
```

Re-run seed only:

```bash
kubectl delete job banking-accounts-seed -n cnpg-banking --ignore-not-found
kubectl apply -k deploy/overlays/kind-local-ambient-cnpg
kubectl wait --for=condition=complete job/banking-accounts-seed -n cnpg-banking --timeout=120s
```

Full re-bootstrap (destructive — deletes DB data):

```bash
kubectl delete cluster banking-accounts banking-ledger -n cnpg-banking
kubectl delete pvc -n cnpg-banking -l cnpg.io/cluster
kubectl apply -k deploy/overlays/kind-local-ambient-cnpg
```

## Secrets

Demo values stay **in git** for easy `kubectl apply -k` (same as upstream intent), but use the **correct K8s types**:

| Object | Kind | Contents |
|--------|------|----------|
| `jwt-key` | Secret | JWT RSA keypair (`demo-jwt` component) |
| `accounts-db-secrets` | Secret | `ACCOUNTS_DB_URI`, `POSTGRES_PASSWORD` |
| `ledger-db-secrets` | Secret | `POSTGRES_PASSWORD`, `SPRING_DATASOURCE_PASSWORD` |
| `demo-data-secrets` | Secret | `DEMO_LOGIN_PASSWORD` (`bankofanthos`) |
| `accounts-db-config` / `ledger-db-config` / `demo-data-config` | ConfigMap | Non-secret: DB names, users, JDBC URL host, `USE_DEMO_DATA`, username |

**CNPG overlay:** `cnpg-connection-targets` (ConfigMap hosts/URLs) + `cnpg-connection-secrets` (passwords/URIs) replace into the app Secrets via Kustomize `replacements`. CNPG bootstrap still uses `credentials-secrets.yaml`.

**Later:** External Secrets Operator — same Secret names/keys, values from Vault (not git).


## Ledger demo transactions

Embedded `ledger-db` loads demo transaction history when `USE_DEMO_DATA=True` (see `src/ledger/ledger-db/initdb/1_create_transactions.sh`). CNPG bootstrap loads **schema only**; the **`banking-ledger-seed` Job** creates the same payroll deposits and inter-user payments for `testuser` / `alice` / `bob` / `eve`.

```bash
kubectl wait --for=condition=complete job/banking-ledger-seed -n cnpg-banking --timeout=120s
```

Re-run only if demo payroll deposits are missing:

```bash
kubectl delete job banking-ledger-seed -n cnpg-banking --ignore-not-found
kubectl apply -k deploy/overlays/kind-local-ambient-cnpg
kubectl wait --for=condition=complete job/banking-ledger-seed -n cnpg-banking --timeout=120s
```

Accounts demo users (`testuser` / `bankofanthos`) are seeded via init SQL + `banking-accounts-seed`.

**Table ownership:** CNPG runs init SQL as `postgres`; schema files end with `ALTER TABLE ... OWNER TO` the app user (`admin` / `accounts-admin`). If you bootstrapped before that fix, run the manual SQL in the troubleshooting section below.

## Troubleshooting

### `permission denied for table transactions`

CNPG ran init SQL as superuser `postgres`, so tables were owned by `postgres` while apps connect as `admin`. Fix on a **running** cluster:

```bash
# Ledger
kubectl exec -n cnpg-banking banking-ledger-1 -- psql -U postgres -d postgresdb -c \
  'ALTER TABLE transactions OWNER TO admin;'

# Accounts (if login/contacts fail similarly)
kubectl exec -n cnpg-banking banking-accounts-1 -- psql -U postgres -d accounts-db -c \
  'ALTER TABLE users OWNER TO "accounts-admin"; ALTER TABLE contacts OWNER TO "accounts-admin";'

kubectl rollout restart deployment -n banking-platform
```

Init SQL in `deploy/components/cnpg-banking-clusters/init-sql/` now includes these `ALTER TABLE ... OWNER` statements for **new** clusters. Existing clusters need the SQL above (or delete PVCs and re-bootstrap).

### Demo users (`testuser`) missing

CNPG bootstrap SQL runs **once**; a single `kubectl apply -k` can also race (Cluster before ConfigMap). The **`banking-accounts-seed` Job** fixes this — re-run:

```bash
kubectl delete job banking-accounts-seed -n cnpg-banking --ignore-not-found
kubectl apply -k deploy/overlays/kind-local-ambient-cnpg
kubectl wait --for=condition=Ready cluster/banking-accounts -n cnpg-banking --timeout=300s
kubectl wait --for=condition=complete job/banking-accounts-seed -n cnpg-banking --timeout=120s
kubectl exec -n cnpg-banking banking-accounts-1 -- \
  psql -U accounts-admin -d accounts-db -c 'SELECT username FROM users;'
```

### Seed Job: `secret "banking-accounts-superuser" not found`

Older Job manifests referenced CNPG’s operator-generated superuser secret (exists only after the Cluster is Ready). Current Job uses **`banking-accounts-credentials`** (from `credentials-secrets.yaml`) and `pg_isready` without auth in the init container. Delete the failed Job and re-apply the overlay.
```

## GitOps

```yaml
spec:
  path: ./deploy/overlays/kind-local-ambient-cnpg
```
