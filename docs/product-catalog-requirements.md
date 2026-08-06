# Product Catalog — Go service requirements

Requirements summary for a new **Go** microservice that extends Tiho Banking Platform with a retail **products / offerings** page (savings, investments, and related bank products).

**Status:** implemented (v1) — Go service + frontend `/products` page. **Not committed yet.**

**Versions (this change):** frontend `v0.6.11` · product-catalog `v0.1.0`

**Verified so far:** `go test ./...` · local `curl` with JWT · `kubectl kustomize` on overlays. **Not yet:** browser UI · full Kind deploy with new images.

**Origin:** discussed in agent session (Jul 2026); see repo git history when committed.

---

## 1. Goals

| Goal | Detail |
|------|--------|
| **Portfolio value** | Show polyglot microservices: Python + Java + **Go** |
| **Domain extension** | Move from “checking + payments” to “what the bank offers” |
| **Low risk** | Read-only catalog in v1 — no ledger or money movement |
| **Platform fit** | Same deploy, JWT, Istio, and CI patterns as existing services |

---

## 2. Service identity

| Item | Value |
|------|-------|
| **Working name** | `product-catalog` (alt: `bank-offerings`) |
| **Language** | Go (1.26 — `go.mod`, Dockerfile, CI aligned) |
| **Source path** | `src/products/` |
| **K8s Service name** | `product-catalog` |
| **Default port** | `8080` |
| **Frontend env var** | `PRODUCTS_API_ADDR` → `product-catalog:8080` |

---

## 3. Functional requirements

### 3.1 Product types (v1)

| Type | Banking terms | Example fields |
|------|---------------|----------------|
| **Savings** | High-yield savings | name, APY, minimum balance, description |
| **Term deposits / CDs** | Fixed-term savings | term (e.g. 12m), rate, early-withdrawal note |
| **Investments (demo)** | Simple fund list | fund name, risk level, disclaimer |
| **Loans** | — | Out of scope for v1; reserved for later |

### 3.2 User-facing behaviour

- **Authenticated users only** — same login as the rest of the app (`testuser` / `bankofanthos`).
- New frontend page at **`/products`**, nav label **“Save & invest”** (locked name for v1).
- Page shows a card grid of bank offerings with optional **type tabs/filters** (All / Savings / CDs / Investments).
- Data is **informational only** — no purchase, signup, or money transfer from this page in v1.
- UI must include a clear **demo disclaimer** (same spirit as upstream: *not a real bank, not financial advice*).

### 3.3 Web UI navigation (v1)

Today’s frontend has a single authenticated view (`/home`) and nav with **Sign out** only. v1 adds:

| UI element | Value |
|------------|-------|
| **Nav link label** | Save & invest |
| **Route** | `/products` |
| **Companion link** | Home → `/home` (existing checking page) |
| **Visibility** | Only when logged in (valid JWT cookie) |
| **On-page filters** | All · Savings · CDs · Investments |

Example nav after login:

```text
[ Bank name ]              Home | Save & invest |  Test User ▼
```

| Tab / filter | API call |
|--------------|----------|
| All (default) | `GET /api/products` |
| Savings | `GET /api/products?type=savings` |
| CDs | `GET /api/products?type=cd` |
| Investments | `GET /api/products?type=investment` |

### 3.4 API behaviour (v1)

| Endpoint | Auth | Description |
|----------|------|-------------|
| `GET /health` | None | Liveness/readiness for Kubernetes |
| `GET /ready` | None | Readiness probe |
| `GET /version` | None | Build/version string (optional, matches BoA pattern) |
| `GET /api/products` | **JWT required** | List all products |
| `GET /api/products?type=savings` | **JWT required** | Filter by type (`savings`, `investment`, `cd`) |

Unauthenticated or invalid token → **401 Unauthorized** on product endpoints.

**Response shape (example):**

```json
{
  "products": [
    {
      "id": "savings-hysa",
      "type": "savings",
      "name": "High-Yield Savings",
      "apy": "4.25",
      "minimumBalance": "0",
      "description": "No monthly fees. FDIC-style demo product."
    }
  ]
}
```

### 3.5 Demo data

- Ship **5–10 static products** in v1 (embedded in Go or loaded from a JSON file at startup).
- Postgres-backed catalog is **optional** for v1; prefer static data to ship faster.

---

## 4. Non-functional requirements

| Area | Requirement |
|------|-------------|
| **Performance** | Low memory footprint; suitable for Kind |
| **Availability** | Stateless; no database dependency in v1 |
| **Observability** | Structured logs to stdout; optional `/metrics` later |
| **Security** | JWT required on product API; public key from `jwt-key` secret; no secrets in code |
| **Licensing** | New files are project-owned; upstream BoA attribution unchanged |

---

## 5. Authentication

Bank of Anthos uses **app-level JWT** (not service-to-service login). The Go service follows the same model as `contacts` and `balancereader`.

```text
User → frontend → userservice (login) → JWT in cookie
User opens /products → frontend verifies cookie → renders page or redirects to login
frontend → product-catalog with Authorization: Bearer <jwt>
product-catalog → verify RS256 with jwt-key public key → return products
```

### v1 — authenticated end-to-end

| Layer | Auth required? | Behaviour |
|-------|----------------|-----------|
| **Frontend page** (`/products`) | **Yes** | Same `verify_token()` pattern as `/home`; redirect to login if missing/invalid |
| **Go API** (`/api/products`) | **Yes** | Require `Authorization: Bearer <jwt>`; **401** without valid token |
| **Health probes** (`/health`, `/ready`) | No | Kubernetes probes only |

Direct `curl` to the Go service without a token must **not** return product data.

### JWT implementation

- Mount **`jwt-key`** secret; read public key from `PUB_KEY_PATH` or env `JWT_PUBLIC_KEY_PATH` (same as other services).
- Verify RS256 with `github.com/golang-jwt/jwt/v5`.
- Valid claims available for logging or future use: `user`, `acct`, `name`.
- **Do not** call `userservice` on every request — only verify the token the frontend already holds.
- **Do not** integrate Keycloak — JWT from `userservice` is sufficient.

### Later (post-v1)

| Phase | Scope |
|-------|-------|
| **v2** | `GET /api/products/recommended` — personalized list using `acct` claim |
| **v3** | Optional Istio `AuthorizationPolicy` as an extra mesh layer (app-level JWT remains primary) |

---

## 6. Frontend integration

The **frontend is the only existing microservice that changes** in v1. All other current services (userservice, contacts, ledger*, loadgenerator) stay unchanged.

### 6.1 Config

Add to `deploy/base/config.yaml` → `service-api-config`:

```yaml
PRODUCTS_API_ADDR: "product-catalog:8080"
```

Wire in `src/frontend/frontend.py` the same way as `CONTACTS_API_ADDR` and `USERSERVICE_API_ADDR`.

### 6.2 Frontend code changes

| File | Change |
|------|--------|
| `src/frontend/frontend.py` | New `@app.route("/products")` with JWT gate; helper to call product-catalog with Bearer token |
| `src/frontend/templates/products.html` | **New** — product cards, type tabs, demo disclaimer; pass same context as `/home` (`name`, etc.) for nav |
| `src/frontend/templates/shared/navigation.html` | Add **Home** and **Save & invest** links when `name` is set (logged in) |
| `src/frontend/static/scripts/products.js` | **Optional** — client-side tab switching / fetch |
| `deploy/base/config.yaml` | Add `PRODUCTS_API_ADDR` |
| `.github/workflows/frontend.yml` | Path filter already covers `src/frontend/**` — no structural change required |

### 6.3 Page behaviour

1. User logs in → lands on `/home`.
2. Clicks **Save & invest** in nav → `GET /products`.
3. Frontend verifies cookie JWT (redirect to login if invalid).
4. Frontend calls `http://${PRODUCTS_API_ADDR}/api/products` with `Authorization: Bearer <token>`.
5. Renders cards; tab click adds `?type=…` (server-side reload or client-side fetch).

### 6.4 Istio / ingress

- No special ingress rules — traffic stays **frontend → product-catalog** inside the mesh.
- Same namespace as other app services (`banking-platform` or overlay namespace).

---

## 7. Changes by service (v1)

Summary of what to touch vs leave alone.

| Service / area | Change? | What to do |
|----------------|---------|------------|
| **`product-catalog` (new)** | **New** | Go API, Dockerfile, tests, K8s Deployment + Service, CI workflow |
| **frontend** | **Yes — small** | Route, template, nav links, `PRODUCTS_API_ADDR`, Bearer call to Go API |
| **userservice** | No | Unchanged — still issues JWT at login |
| **contacts** | No | Unchanged |
| **balancereader** | No | Unchanged |
| **ledgerwriter** | No | Unchanged |
| **transactionhistory** | No | Unchanged |
| **accounts-db / ledger-db** | No | Unchanged — no schema changes |
| **CNPG clusters** | No | Unchanged — product catalog has no DB in v1 |
| **loadgenerator** | No | Unchanged (optional: ignore) |
| **`deploy/base/`** | **Yes** | Add `product-catalog.yaml`; update `config.yaml`, `kustomization.yaml` |
| **`deploy/overlays/kind-local/`** | **Yes** | Image pin for `product-catalog` (inherited by ambient overlays) |
| **`deploy/overlays/gke-dev/`** | **Yes** | Image pin for `product-catalog` |
| **`deploy/components/disable-gcp-telemetry/`** | **Yes** | Add `patch-product-catalog.yaml` for Kind/non-GKE overlays |
| **Istio / Gateway** | No | No new ingress; in-cluster traffic only |
| **Platform repo** (`tiho-local-kind-cluster`) | No | No change required for v1 |

**Rule of thumb:** new Go service + frontend glue only; no changes to ledger, accounts, or auth flow.

---

## 8. Deployment (Kustomize)

### 8.1 New files

```text
src/products/
  main.go
  go.mod
  go.sum
  Dockerfile
  README.md

deploy/base/product-catalog.yaml    # Deployment + Service (+ jwt-key volume mount)

.github/workflows/product-catalog.yml
```

`product-catalog` Deployment must mount **`jwt-key`** (public key) the same way Java/Python backends do.

### 8.2 Base kustomization

Add `product-catalog.yaml` to `deploy/base/kustomization.yaml` `resources`.

### 8.3 Overlays

Each overlay that runs the full app must include the new Deployment and image pin, same as `contacts` / `userservice`:

- `deploy/overlays/kind-local/` — image pin + inherits on ambient overlays
- `deploy/overlays/gke-dev/`
- `deploy/components/use-ghcr-images/` — frontend `v0.6.11`, product-catalog `v0.1.0`

### 8.4 Container image

There is **no upstream Google image** for this service (see §15.3). Build your own image and pin it in overlay `images:` blocks.

- Base manifest: set `image:` in `deploy/base/product-catalog.yaml` (name matches Kustomize `images[].name`).
- **Image pins:** add to `deploy/overlays/kind-local/kustomization.yaml` and `deploy/overlays/gke-dev/kustomization.yaml`. Nested overlays (`kind-local-ambient`, `kind-local-ambient-cnpg`) inherit from `kind-local` — no separate `images:` block there.
- CI push: `ghcr.io/${{ github.repository }}/product-catalog` when `ENABLE_IMAGE_PUSH=true` (same pattern as `contacts.yml`).

---

## 9. CI / CD

New workflow: `.github/workflows/product-catalog.yml`

| Trigger | Path filter |
|---------|-------------|
| `push` / `pull_request` on `main` | `src/products/**`, workflow file |

Frontend changes are covered by existing `.github/workflows/frontend.yml` path filters (`src/frontend/**`).

| Step | Tool |
|------|------|
| Unit tests | `go test ./...` |
| Lint (optional) | `go vet`, `staticcheck` |
| Build | `docker build` from `src/products/Dockerfile` |
| Push (when enabled) | GHCR, same `ENABLE_IMAGE_PUSH` pattern as Python/Java workflows |

---

## 10. Local development

### 10.1 Go on laptop, rest in Kind

```text
Kind cluster:  userservice, frontend, ledger*, DBs
Laptop:        go run ./src/products  (localhost:8080)
```

```bash
# JWT public key from cluster (required for product API)
kubectl get secret jwt-key -n banking-platform \
  -o jsonpath='{.data.jwtRS256.key\.pub}' | base64 -d > /tmp/jwt.pub

export JWT_PUBLIC_KEY_PATH=/tmp/jwt.pub
go run ./src/products
```

### 10.2 Frontend → local Go

Temporarily point frontend at the laptop:

- `PRODUCTS_API_ADDR=host.docker.internal:8080` (Docker Desktop), or
- port-forward / extraPortMapping as needed for Kind.

### 10.3 Manual API test

```bash
# Login via UI, copy token from cookie, then:
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/products
curl -H "Authorization: Bearer <token>" "http://localhost:8080/api/products?type=savings"

# Without token — must fail:
curl -i http://localhost:8080/api/products   # expect 401
```

---

## 11. Architecture placement

```text
Browser
  └── frontend (Python)
        ├── userservice ──► accounts-db
        ├── contacts ─────► accounts-db
        ├── balancereader ──► ledger-db
        ├── transactionhistory ──► ledger-db
        ├── ledgerwriter ──► ledger-db
        └── product-catalog (Go, new)   ← read-only, no DB in v1
```

**In scope:** product listing, frontend page, deploy, CI.

**Out of scope (v1):**

- Real investment trading or account opening
- Ledger writes or balance changes
- New PostgreSQL database (unless added in a later version)
- Keycloak / external OIDC
- Calls from Go service to other backends (unless added in v2+)

---

## 12. Alternatives considered

| Service idea | Why not first choice |
|--------------|----------------------|
| **notifications** | No visible UI page |
| **rates** | Narrower story than full product shelf |
| **savings-accounts** | Needs DB + ledger integration — higher complexity |
| **statements** | Medium complexity, less visible |
| **fraud-scoring** | Medium–high complexity, no user-facing page |

**Chosen:** `product-catalog` — best match for Go + banking terminology + portfolio UI.

---

## 13. Acceptance criteria (v1)

### Code & deploy (implemented)

- [x] `src/products/` builds and passes `go test ./...`
- [x] `GET /api/products` returns demo savings / CD / investment entries **with valid JWT** (local curl)
- [x] `GET /api/products` returns **401** without JWT (unit test + local curl)
- [x] `GET /api/products?type=…` filters correctly (unit test)
- [x] `/health` and `/ready` return 200 (no auth)
- [x] Dockerfile present
- [x] `deploy/base/product-catalog.yaml` in base; mounts `jwt-key`
- [x] Frontend **Save & invest** nav + **Home** when logged in
- [x] `/products` route with JWT gate (redirect to login if invalid)
- [x] `/products` template with cards, type tabs, demo disclaimer
- [x] Frontend passes Bearer token to product-catalog
- [x] `PRODUCTS_API_ADDR` in `service-api-config`
- [x] `.github/workflows/product-catalog.yml` with path filters
- [x] `kubectl kustomize deploy/overlays/kind-local-ambient-cnpg` (and ESO overlay) builds
- [x] `disable-gcp-telemetry` patch for `product-catalog`
- [x] **No changes** to userservice, contacts, or ledger services
- [x] frontend `release.yaml` bumped to **0.6.11**; `use-ghcr-images` updated

### Still to verify before calling v1 “done”

- [ ] Dockerfile produces a runnable image (CI or local `docker build`)
- [ ] `/products` in browser (frontend local or Kind)
- [ ] Kind deploy: login → **Save & invest** → products visible
- [ ] GHCR push: `frontend:v0.6.11` + `product-catalog:v0.1.0` after commit + push

---

## 14. Future work (post-v1)

- `GET /api/products/recommended` — personalized list using `acct` claim from JWT
- Postgres table + CNPG for editable product catalog
- Admin API or GitOps-managed ConfigMap for product data
- Loan and card product types
- OpenTelemetry traces (align with platform observability)
- Flux GitOps wiring in `tiho-local-kind-cluster`

---

## 15. Implementation references (for agents)

Use this section when implementing v1. Read the requirements above first; use these repo paths as concrete patterns.

### 15.1 Copy patterns from

| What | File to mirror |
|------|----------------|
| K8s Deployment + Service + JWT volume | [`deploy/base/contacts.yaml`](../deploy/base/contacts.yaml) |
| CI workflow (test + Docker build + optional GHCR push) | [`.github/workflows/contacts.yml`](../.github/workflows/contacts.yml) |
| Frontend auth gate (redirect if no token) | [`src/frontend/frontend.py`](../src/frontend/frontend.py) — `@app.route("/home")` |
| Frontend backend URI wiring | same file — `CONTACTS_URI` / `app.config[...]` block (~line 661) |
| Overlay image pins | [`deploy/overlays/kind-local/kustomization.yaml`](../deploy/overlays/kind-local/kustomization.yaml) `images:` |
| JWT secret (Kind overlays) | [`deploy/components/demo-jwt/`](../deploy/components/demo-jwt/) — included via `kind-local` |

### 15.2 JWT mount (must match other backends)

From `contacts` Deployment:

- Secret: **`jwt-key`**
- Mount path: **`/tmp/.ssh/publickey`** (key `jwtRS256.key.pub` → file `publickey`)
- ConfigMap `environment-config`: **`PUB_KEY_PATH: "/tmp/.ssh/publickey"`**
- Go service reads public key from `PUB_KEY_PATH` (or `JWT_PUBLIC_KEY_PATH` for local dev)

### 15.3 Container image (no upstream image exists)

Unlike frontend/contacts, **Google does not publish a `product-catalog` image**. For v1:

1. Add `src/products/Dockerfile` and build locally or via CI.
2. Add to each overlay `images:` block, e.g.:

```yaml
  - name: ghcr.io/<owner>/tiho-banking-platform/product-catalog
    newTag: v0.1.0
```

3. In `deploy/base/product-catalog.yaml`, set `image:` to that name (tag overridden by Kustomize).
4. Until `ENABLE_IMAGE_PUSH` is on, build/push manually or use `imagePullPolicy: Never` + `kind load docker-image` for local Kind testing.

### 15.4 Telemetry component

Kind overlays use [`deploy/components/disable-gcp-telemetry`](../deploy/components/disable-gcp-telemetry/). Add **`patch-product-catalog.yaml`** there (same `ENABLE_TRACING: "false"`, `ENABLE_METRICS: "false"`) and register it in that component’s `kustomization.yaml`.

### 15.5 Primary test overlay

Implement and verify against:

```bash
kubectl apply -k deploy/overlays/kind-local-ambient-cnpg
```

Then: login → **Save & invest** → products visible; unauthenticated `/products` redirects to login.

Port-forward ingress (ambient Kind):

```bash
kubectl port-forward -n istio-ingress svc/istio-ingress 8080:80
```

Open **http://localhost:8080/** — demo user `testuser` / `bankofanthos`.

### 15.6 Agent handoff prompt

Paste this when starting a new agent:

```text
Implement product-catalog v1 per docs/product-catalog-requirements.md.

Rules:
- Mirror deploy/base/contacts.yaml for K8s + JWT; mirror .github/workflows/contacts.yml for CI.
- Mirror /home auth pattern in frontend.py for /products.
- Only change: new src/products/, frontend (route/template/nav/config), deploy/, CI, disable-gcp-telemetry patch.
- Do NOT modify userservice, contacts, ledger*, or databases.
- JWT required on /api/products (401 without token); page auth required.
- Nav: Home | Save & invest; route /products; tabs All/Savings/CDs/Investments.
- No upstream image — build ghcr.io/.../product-catalog and pin in overlay images:.
- Definition of done: acceptance criteria in §13 of the requirements doc.
- Test overlay: deploy/overlays/kind-local-ambient-cnpg
```

---

## 16. Related docs

| Document | Contents |
|----------|----------|
| [architecture.md](architecture.md) | Service map including product-catalog |
| [TODO.md](TODO.md) | Platform follow-ups |
| [deploy/README.md](../deploy/README.md) | Kustomize overlays and deploy commands |
| [src/products/README.md](../src/products/README.md) | Run/test the Go service locally |
