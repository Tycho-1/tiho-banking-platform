# Deploy — Tiho Banking Platform

Kubernetes manifests for all microservices, structured for **Kustomize** and **GitOps** (Flux / Argo CD).

## Layout

```
deploy/
├── base/                    # Portable app manifests (all services)
├── components/              # Reusable Kustomize components (patch sets)
│   ├── demo-jwt/            # Shared demo JWT (kind-local + gke-dev)
│   ├── disable-gcp-telemetry/
│   ├── strip-gke-workload-identity/
│   ├── frontend-clusterip/
│   ├── istio-ambient-mesh/  # Ambient labels + waypoint (kind-local-ambient)
│   ├── gateway-api-ingress/ # Gateway API Gateway + HTTPRoute (kind-local-gateway-api)
│   ├── use-cnpg-databases/  # Point apps at CNPG; drop embedded StatefulSets
│   ├── cnpg-banking-clusters/  # Optional dedicated CNPG Cluster CRs
│   ├── use-ghcr-images/     # Remap base images → GHCR (enabled on kind-local)
│   ├── use-eso-vault/       # ExternalSecret → Vault (opt-in; replaces demo Secrets)
│   └── use-eso-vault-cnpg/  # ExternalSecret → Vault for CNPG connection + bootstrap creds
└── overlays/
    ├── kind-local/          # Local Kind — Istio sidecar + Istio Gateway/VS
    ├── kind-local-eso/      # kind-local + External Secrets (Vault)
    ├── kind-local-gateway-api/  # Kind — sidecar + Kubernetes Gateway API ingress
    ├── kind-local-ambient/  # Same as kind-local, Istio ambient + waypoint
    ├── kind-local-ambient-cnpg/  # Ambient + CloudNativePG
    ├── kind-local-ambient-cnpg-eso/  # ambient + CNPG + ESO/Vault
    └── gke-dev/             # GKE dev — upstream-shaped (LoadBalancer, GCP telemetry on)
```

| Layer | Responsibility |
|-------|----------------|
| **base** | Deployments, Services, StatefulSets, ConfigMaps + Secrets — close to upstream `kubernetes-manifests/` |
| **components** | Reusable patches — Kind overlays typically use JWT, telemetry-off, WI strip, ClusterIP; Gateway API / ambient / CNPG as needed |
| **overlays** | Per-target: namespace, secrets, ingress model, image tags |

## Overlays at a glance

| Overlay | Target | Intent |
|---------|--------|--------|
| **`kind-local`** | Kind + Istio sidecar + ingress LB | Local mesh — **Istio sidecar** + classic Istio Gateway/VS |
| **`kind-local-eso`** | Same as `kind-local` + ESO + Vault | Secrets from Vault — [README](overlays/kind-local-eso/README.md) |
| **`kind-local-gateway-api`** | Kind + Gateway API CRDs + a `GatewayClass` | Sidecar + **[Gateway API](https://gateway-api.sigs.k8s.io/docs/introduction/)** ingress (`Gateway`/`HTTPRoute`) — [README](overlays/kind-local-gateway-api/README.md) |
| **`kind-local-ambient`** | Kind + Istio ambient | Same stack, **Istio ambient** + waypoint — see [kind-local-ambient/README.md](overlays/kind-local-ambient/README.md) |
| **`kind-local-ambient-cnpg`** | Kind + ambient + CNPG operator | Ambient + **CloudNativePG** instead of embedded Postgres — [kind-local-ambient-cnpg/README.md](overlays/kind-local-ambient-cnpg/README.md) |
| **`kind-local-ambient-cnpg-eso`** | Same + ESO + Vault | CNPG + secrets from Vault — [README](overlays/kind-local-ambient-cnpg-eso/README.md) |
| **`gke-dev`** | GKE development | **Upstream default** — LoadBalancer frontend, GCP telemetry **on** (see [gke-dev/README.md](overlays/gke-dev/README.md)) |

Same **`base/`** for all — no contradiction. **`gke-dev`** is what Bank of Anthos was written for; Kind overlays adapt it for a local cluster you provide.

## Prerequisites

Overlays assume a **Kubernetes cluster you create** (typically Kind). This app repo does **not** install the cluster or mesh/operators. Document prerequisites per overlay; a public Kind/platform reference can be linked later ([docs/TODO.md](../docs/TODO.md)).

**Kind (sidecar + classic Istio ingress):** Kind + **Istio sidecar** + an ingress gateway Service (e.g. Helm `istio-ingress` + MetalLB). Overlay **`kind-local`** pulls **GHCR** images by default (`use-ghcr-images`); packages must be public or use an `imagePullSecret`.

**Kind (ESO + Vault):** Same as `kind-local`, but secrets from Vault — overlay **`kind-local-eso`**. Requires platform **ESO + Vault + `SecretStore` `vault-banking-platform`** (and auth) — [README](overlays/kind-local-eso/README.md).

**Kind (Gateway API ingress):** Kind + **Gateway API CRDs** + a `GatewayClass` (today often `istio` if Istio is installed). Overlay **`kind-local-gateway-api`** — [README](overlays/kind-local-gateway-api/README.md). Uses Kubernetes `Gateway` + `HTTPRoute` instead of Istio `VirtualService`. Later, point `gatewayClassName` at a non-Istio controller if you install one on the cluster.

**Kind (ambient):** Kind with **Istio ambient** only, Gateway API CRDs, Istio CNI for ztunnel — overlay **`kind-local-ambient`**, not `kind-local`.

**Kind (ambient + CNPG):** CNPG operator installed; use **`kind-local-ambient-cnpg`**. Edit `cnpg-connection-targets.yaml` to point at your clusters.

**Kind (ambient + CNPG + ESO):** Overlay **`kind-local-ambient-cnpg-eso`** — same as ambient-cnpg with ExternalSecrets; platform provides SecretStores in `banking-platform` and `cnpg-banking` — [README](overlays/kind-local-ambient-cnpg-eso/README.md).

**GKE:** GKE cluster + `kubectl`; for full GCP telemetry, Workload Identity / Cloud Ops (or add `disable-gcp-telemetry` component until configured).

**Later:** `eks-dev`, `aks-dev` — same pattern as `gke-dev` / `kind-local`.

## Container images

| Overlay | Default image source |
|---------|---------------------|
| **`kind-local`**, **`kind-local-eso`**, **`kind-local-gateway-api`** | **GHCR** — `ghcr.io/tycho-1/tiho-banking-platform/<service>` (component `use-ghcr-images` in overlay) |
| **`kind-local-ambient`**, **`kind-local-ambient-cnpg`**, **`kind-local-ambient-cnpg-eso`** | **GHCR** — inherits `use-ghcr-images` from [`kind-local`](overlays/kind-local/kustomization.yaml) |
| **`gke-dev`** | **Upstream** BoA on Google Artifact Registry |

Tags for GHCR come from each service’s **`src/<service>/release.yaml`** (e.g. frontend `v0.6.10`, others `v0.6.9`). CI publishes on push to `main` when **`ENABLE_IMAGE_PUSH=true`** — see [`.github/workflows/README.md`](../.github/workflows/README.md) and [components/use-ghcr-images/README.md](components/use-ghcr-images/README.md).

### Use upstream BoA images on Kind overlays

Remove **`use-ghcr-images`** from `components:` in [overlays/kind-local/kustomization.yaml](overlays/kind-local/kustomization.yaml) and/or [overlays/kind-local-gateway-api/kustomization.yaml](overlays/kind-local-gateway-api/kustomization.yaml). The overlay `images:` block then pins upstream tags:

```text
us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/<service>:v0.6.9
```

No GHCR packages required — good for a fresh clone without CI.

### Enable GHCR on other overlays

Add the component to the overlay you deploy:

```yaml
components:
  - ../../components/use-ghcr-images
```

Pin tags in [components/use-ghcr-images/kustomization.yaml](components/use-ghcr-images/kustomization.yaml) to match each `release.yaml`.

## Secrets (demo in git, correct K8s types)

Passwords and JWT material use **Secret** objects (not ConfigMaps). On **`kind-local`** / most overlays, values remain in git so `kubectl apply -k` stays one-shot — same ease as upstream BoA, better typing for RBAC.

For **Vault-backed secrets**, use **`kind-local-eso`** or **`kind-local-ambient-cnpg-eso`**. This repo ships **ExternalSecrets** only; the platform provides `SecretStore` `vault-banking-platform` (+ auth). See [`use-eso-vault`](components/use-eso-vault/) / [`use-eso-vault-cnpg`](components/use-eso-vault-cnpg/).

Expected KV v2 layout under mount **`platform-kv`**:

```text
platform-kv/
└── tiho-banking-platform/
    ├── db/
    │   ├── accounts          # → accounts-db-secrets          (use-eso-vault)
    │   └── ledger            # → ledger-db-secrets            (use-eso-vault)
    ├── auth/
    │   └── jwt               # → jwt-key                      (use-eso-vault)
    ├── demo/
    │   └── login             # → demo-data-secrets            (use-eso-vault)
    └── cnpg/                 # CNPG overlays only (use-eso-vault-cnpg)
        ├── accounts-creds    # → banking-accounts-credentials
        ├── ledger-creds      # → banking-ledger-credentials
        └── connection        # → cnpg-connection-secrets
```

| Secret | Keys (examples) |
|--------|-----------------|
| `jwt-key` | RSA keypair (`demo-jwt` on non-ESO overlays; Vault on `kind-local-eso`) |
| `accounts-db-secrets` | `ACCOUNTS_DB_URI`, `POSTGRES_PASSWORD` (`demo-password-change-me`) |
| `ledger-db-secrets` | `POSTGRES_PASSWORD`, `SPRING_DATASOURCE_PASSWORD` (`demo-password-change-me`) |
| `demo-data-secrets` | `DEMO_LOGIN_PASSWORD` (`bankofanthos` — matches seeded bcrypt hash) |

ConfigMaps keep non-secret config (`POSTGRES_DB`, hosts/JDBC URL, `USE_DEMO_DATA`, usernames).

## Install (manual)

**Kind (sidecar)** — overlay `kind-local`:

```bash
# Preview
kubectl kustomize deploy/overlays/kind-local

# Apply
kubectl apply -k deploy/overlays/kind-local

# Wait
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
kubectl get pods -n banking-platform
```

**Kind (sidecar + ESO)** — overlay `kind-local-eso` (see [overlays/kind-local-eso/README.md](overlays/kind-local-eso/README.md)):

```bash
kubectl apply -k deploy/overlays/kind-local-eso
kubectl wait --for=condition=Ready externalsecret --all -n banking-platform --timeout=120s
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
```

**Kind (Gateway API ingress)** — overlay `kind-local-gateway-api` ([README](overlays/kind-local-gateway-api/README.md)):

```bash
kubectl apply -k deploy/overlays/kind-local-gateway-api
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
kubectl wait --for=condition=Programmed gateway/banking-platform -n banking-platform --timeout=120s
kubectl get gateway,httproute,svc -n banking-platform
```

Access via the **Gateway** LoadBalancer EXTERNAL-IP (Istio-managed), not `istio-ingress` — see overlay README.

**Kind (ambient)** — overlay `kind-local-ambient` (your cluster today):

```bash
kubectl apply -k deploy/overlays/kind-local-ambient
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
```

**Kind (ambient + CNPG)** — overlay `kind-local-ambient-cnpg`:

```bash
kubectl apply -k deploy/overlays/kind-local-ambient-cnpg
kubectl wait --for=condition=Ready cluster/banking-accounts cluster/banking-ledger -n cnpg-banking --timeout=300s

# Demo users + demo ledger transactions (seed Jobs; idempotent)
kubectl wait --for=condition=complete job/banking-accounts-seed -n cnpg-banking --timeout=120s
kubectl wait --for=condition=complete job/banking-ledger-seed -n cnpg-banking --timeout=120s

kubectl rollout restart deployment -n banking-platform
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
```

Specify CNPG endpoints in **`deploy/overlays/kind-local-ambient-cnpg/cnpg-connection-targets.yaml`** (`<cluster>-rw.<namespace>.svc.cluster.local`).

**GKE** — overlay `gke-dev` (upstream-shaped — see [overlays/gke-dev/README.md](overlays/gke-dev/README.md)):

```bash
kubectl apply -k deploy/overlays/gke-dev
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
kubectl get svc frontend -n banking-platform
```

## Access

Demo login (all overlays): **`testuser`** / **`bankofanthos`**

### Kind — Gateway API (`kind-local-gateway-api`)

Uses **[Kubernetes Gateway API](https://gateway-api.sigs.k8s.io/docs/introduction/)** (`Gateway` + `HTTPRoute`), not the classic Istio `VirtualService` path and **not** the shared `istio-ingress` Service.

Istio (as `GatewayClass` `istio`) creates a LoadBalancer Service in `banking-platform` (e.g. `banking-platform-istio`). MetalLB (or similar) assigns an EXTERNAL-IP — often a different address than the shared `istio-ingress` Service.

```bash
kubectl get gateway,httproute,svc -n banking-platform
kubectl get gateway banking-platform -n banking-platform -o jsonpath='{.status.addresses[0].value}{"\n"}'
```

Open **`http://<GATEWAY-EXTERNAL-IP>/`**.

```bash
# Example check
curl -sS -o /dev/null -w "%{http_code}\n" http://<GATEWAY-EXTERNAL-IP>/
```

Full detail: [overlays/kind-local-gateway-api/README.md](overlays/kind-local-gateway-api/README.md).

Do not mix this overlay with `kind-local` in the same namespace (different ingress models).

### Kind — classic Istio ingress (`kind-local`, `kind-local-ambient`, `kind-local-ambient-cnpg`)

These overlays route traffic through **Istio**, not the frontend Service directly. The app creates a classic Istio `Gateway` + `VirtualService` in `banking-platform` that sends `/` to the `frontend` Service. External traffic enters via an **Istio ingress gateway** installed on the cluster (not by this repo) — typically a LoadBalancer Service such as `istio-ingress` in namespace `istio-ingress`.

Default layout: Helm `istio/gateway` in namespace **`istio-ingress`**, service **`istio-ingress`**, listener port **80**.

**Option A — LoadBalancer (recommended on Kind + MetalLB)**

Find the ingress IP (wait until `EXTERNAL-IP` is assigned):

```bash
kubectl get svc -n istio-ingress istio-ingress
# or just the IP:
kubectl get svc -n istio-ingress istio-ingress -o jsonpath='{.status.loadBalancer.ingress[0].ip}{"\n"}'
```

Open **`http://<EXTERNAL-IP>/`** (port 80).

**Option B — port-forward (no MetalLB / browser on another machine)**

Forward local port 8080 to ingress port 80:

```bash
kubectl port-forward -n istio-ingress svc/istio-ingress 8080:80
```

Open **`http://localhost:8080/`**. Leave the command running while you use the app.

**Verify routing is applied**

```bash
kubectl get gateway.networking.istio.io,virtualservice.networking.istio.io -n banking-platform
# Expect: banking-platform-gateway, banking-platform-frontend
```

**Wait for app readiness before opening the UI**

```bash
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
```

For **`kind-local-ambient-cnpg`**, also wait for CNPG clusters (and optionally the seed Job) — see [kind-local-ambient-cnpg/README.md](overlays/kind-local-ambient-cnpg/README.md).

### GKE (`gke-dev`)

Upstream style: the **frontend** Service is a **LoadBalancer** (no Istio ingress in this overlay).

```bash
kubectl get svc frontend -n banking-platform
```

Open **`http://<EXTERNAL-IP>/`** (port 80).

### Classic istioctl ingress (Kind only, optional)

If you swapped `gateway-helm.yaml` for `gateway-istioctl.yaml` in `kind-local/kustomization.yaml`, use **`istio-ingressgateway`** in **`istio-system`**, port **8080** — see [Istio gateway variants](#istio-gateway-variants-kind-local-only) below.

## Istio gateway variants (kind-local only)

| File | When to use |
|------|-------------|
| `overlays/kind-local/istio/gateway-helm.yaml` | **Default** — `istio-ingress` namespace, selector `istio: ingress`, port 80 |
| `overlays/kind-local/istio/gateway-istioctl.yaml` | Classic `istio-ingressgateway` in `istio-system`, port 8080 |

To switch: in `overlays/kind-local/kustomization.yaml`, replace `istio/gateway-helm.yaml` with `istio/gateway-istioctl.yaml`.

## GitOps (Flux example)

Point a `Kustomization` at this repo path:

```yaml
spec:
  path: ./deploy/overlays/kind-local
  sourceRef:
    kind: GitRepository
    name: tiho-banking-platform
```

To roll **one service**, change only its `images[].newTag` in `kind-local/kustomization.yaml` (or a small patch file) — Flux reconciles the diff.

## Optional: disable load generator

Add to `kind-local/kustomization.yaml`:

```yaml
patches:
  - target:
      kind: Deployment
      name: loadgenerator
    patch: |
      - op: replace
        path: /spec/replicas
        value: 0
```

## Future overlays

When you add cloud clusters, use **consistent names** — for learning and portfolio work, **`dev` on each cloud** is enough before you introduce staging/prod:

```text
deploy/overlays/kind-local/              # local Kind — Istio sidecar + Gateway/VS
deploy/overlays/kind-local-eso/          # local Kind — sidecar + ESO/Vault
deploy/overlays/kind-local-gateway-api/  # local Kind — Gateway API ingress (present now)
deploy/overlays/kind-local-ambient/      # local Kind — Istio ambient
deploy/overlays/kind-local-ambient-cnpg/ # local Kind — ambient + CNPG
deploy/overlays/kind-local-ambient-cnpg-eso/ # ambient + CNPG + ESO/Vault
deploy/overlays/gke-dev/                 # GKE development (upstream-shaped) — present now
deploy/overlays/eks-dev/                 # AWS EKS development — planned
deploy/overlays/aks-dev/                 # Azure AKS development — planned
```

Later, if you need promotion lanes: `eks-staging`, `eks-prod`, etc. — same cloud, different overlay per **environment**, not mixed names across clouds.

Each overlay reuses `base/` + `components/` and adds only what that **target** needs (ingress, secrets operator, DB endpoints, image tags).

## Intent (all environments)

| Layer | Purpose |
|-------|---------|
| **`base/`** | Shared manifests, kept close to upstream BoA `kubernetes-manifests/` for sync |
| **`components/`** | Optional portable patches — used by **kind-local** (not required on **gke-dev**) |
| **`overlays/*`** | Per-target only: namespace, ingress (Istio VS **or** Gateway API), image tags, cloud-specific config |

**This file** is the main deploy reference for every environment. Cluster bootstrap (Kind, Istio, CNPG operator, MetalLB, …) lives outside this repo.
