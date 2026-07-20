# GKE dev overlay

Closest to the **upstream Bank of Anthos** quickstart: same `base/` manifests with **no** Kind-specific components.

## vs `kind-local`

| | **gke-dev** | **kind-local** |
|---|-------------|----------------|
| **Upstream default** | Yes — LoadBalancer frontend, GCP telemetry **on** | No — adapted for local Kind |
| **Components** | None | `disable-gcp-telemetry`, `strip-gke-workload-identity`, `frontend-clusterip` |
| **Ingress** | Frontend `LoadBalancer` (like upstream README) | Istio Gateway + `ClusterIP` frontend |
| **Namespace labels** | Plain namespace | `istio-injection: enabled` |
| **JWT** | `deploy/components/demo-jwt/` (shared) | Same |

No contradiction: **`base`** stays upstream-shaped; each **overlay** picks what that cluster needs.

## Prerequisites

- GKE cluster + `kubectl` credentials
- Pull access to `us-central1-docker.pkg.dev/bank-of-anthos-ci/...`
- For **GCP telemetry** (`ENABLE_METRICS` / `ENABLE_TRACING` = true in base): either
  - **Workload Identity** + Cloud Operations configured (see [upstream environments.md](https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/main/docs/environments.md)), or
  - Temporarily add component `../../components/disable-gcp-telemetry` to `kustomization.yaml` until Cloud Ops is set up (Java pods otherwise look for GCP credentials).

Optional: patch `ServiceAccount` `bank-of-anthos` Workload Identity annotation to **your** Google service account (base still ships the upstream demo annotation).

## Install

```bash
kubectl apply -k deploy/overlays/gke-dev
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
kubectl get svc frontend -n banking-platform   # EXTERNAL-IP when LB is ready
```

Demo login: **`testuser`** / **`bankofanthos`**

## Optional: Istio on GKE

Upstream extras live under `bank-of-anthos/extras/istio/`. Add Gateway + VirtualService to this overlay (or a `components/` block) when using Anthos Service Mesh / Istio on GKE — same idea as `kind-local/istio/`, with GKE ingress/LB layout.

## GitOps

```yaml
spec:
  path: ./deploy/overlays/gke-dev
```
