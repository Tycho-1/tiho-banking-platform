# Architecture

Retail banking demo: HTTP microservices, JWT auth, two PostgreSQL databases.

## Service map

```
Browser
  └── frontend (Python)
        ├── userservice ──► accounts-db
        ├── contacts ─────► accounts-db
        ├── balancereader ──► ledger-db
        ├── transactionhistory ──► ledger-db
        ├── ledgerwriter ──► ledger-db
        └── product-catalog (Go)   ← read-only products API, no DB (v1)
```

| Service | Role |
|---------|------|
| **frontend** | Web UI — login, signup, home, deposit, transfer, **Save & invest** (`/products`) |
| **product-catalog** | Read-only bank products API (Go) — savings, CDs, investments demo |
| **userservice** | Accounts, bcrypt passwords, **issues RS256 JWT** |
| **contacts** | Per-user contact list for payment forms |
| **balancereader** | Cached balances from ledger |
| **transactionhistory** | Cached transaction history |
| **ledgerwriter** | Validates and writes transactions |
| **accounts-db** | User accounts (embedded StatefulSet **or** CNPG) |
| **ledger-db** | Ledger (embedded StatefulSet **or** CNPG) |
| **loadgenerator** | Locust synthetic traffic (optional on Kind) |

## Authentication

1. User posts credentials to **frontend** → **userservice** `GET /login`.
2. **userservice** returns JWT signed with `jwt-key` secret (claims: `user`, `acct`, `name`).
3. **frontend** stores token in cookie; backends receive `Authorization: Bearer …`.
4. Java services verify JWT with public key and enforce `acct` matches requested account.

Non-secret config: ConfigMaps in `deploy/base/config.yaml`. Passwords / JWT material: Kubernetes **Secrets** (demo values in git; ESO later).

Demo user: `testuser` / `bankofanthos` (when `USE_DEMO_DATA=True`).

## Money flow

- **Deposit / transfer:** frontend → ledgerwriter → ledger-db.
- **Balance:** frontend → balancereader (in-memory cache updated from ledger stream).
- **History:** frontend → transactionhistory.

## Deployment (this repo)

| Path | Role |
|------|------|
| `deploy/base/` | Shared Deployments, Services, ConfigMaps, Secrets |
| `deploy/components/` | Portable patches (JWT, telemetry off, Gateway API, CNPG, …) |
| `deploy/overlays/` | Per-target Kind / GKE compose |

**Kind overlays** (pick one at a time in a namespace):

| Overlay | Ingress | Mesh / data |
|---------|---------|-------------|
| `kind-local` | Classic Istio Gateway + VirtualService → cluster ingress LB | Sidecar |
| `kind-local-gateway-api` | Kubernetes Gateway API (`Gateway` + `HTTPRoute`) | Sidecar |
| `kind-local-ambient` | Same as `kind-local` (inherits VS) | Ambient + waypoint |
| `kind-local-ambient-cnpg` | Same as ambient | Ambient + **CloudNativePG** (no embedded DB StatefulSets) |

**GKE:** `gke-dev` — frontend LoadBalancer, upstream-shaped (GCP telemetry on).

Full install/access: [deploy/README.md](../deploy/README.md).

GitOps (Flux) from a platform repo → chosen overlay path — planned.

## Mesh and ingress

Istio is not in application code. Overlays set namespace mesh labels and either:

- **Classic Istio** `networking.istio.io` Gateway + VirtualService, or  
- **Gateway API** `gateway.networking.k8s.io` Gateway + HTTPRoute  

The cluster must provide the dataplane (Istio sidecar/ambient and/or a GatewayClass controller). That bootstrap is outside this app repo.
