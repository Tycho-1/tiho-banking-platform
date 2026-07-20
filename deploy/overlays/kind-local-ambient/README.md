# Kind local — Istio ambient

Same app stack as [`kind-local`](../kind-local/) (upstream images, portable components, Istio ingress gateway), but joins an **ambient** mesh instead of sidecar injection.

| Piece | Location |
|-------|----------|
| Base overlay | Reuses `../kind-local` |
| Ambient labels + waypoint | `deploy/components/istio-ambient-mesh/` |

## Prerequisites

- Kind cluster with **Istio ambient** only (sidecar injection off; ambient dataplane + waypoint)
- **Gateway API CRDs** (needed for the ambient waypoint `Gateway`)
- **Istio CNI** enabled for healthy ztunnel on Kind (do not skip CNI on ambient installs)

Do **not** mix sidecar and ambient dataplanes on one cluster.

## Install

```bash
kubectl apply -k deploy/overlays/kind-local-ambient
kubectl wait --for=condition=Available deployment --all -n banking-platform --timeout=300s
```

Access and demo login are the same as **kind-local** — see [deploy/README.md](../../README.md#access).

## GitOps

```yaml
spec:
  path: ./deploy/overlays/kind-local-ambient
```
