# Kind local overlay

| Piece | Location |
|-------|----------|
| Namespace `banking-platform` | `namespace.yaml` |
| Demo JWT secret | `deploy/components/demo-jwt/` |
| Istio Gateway + VirtualService | `istio/gateway-helm.yaml` (default) |
| **Container images** | component `use-ghcr-images` → `ghcr.io/tycho-1/tiho-banking-platform/<service>` (tags in [use-ghcr-images](../../components/use-ghcr-images/kustomization.yaml)) |
| GCP telemetry off | component `disable-gcp-telemetry` |
| Frontend ClusterIP | component `frontend-clusterip` |
| Prometheus scrape | component `prometheus-servicemonitors` (product-catalog `/metrics`) |

**Mesh:** namespace label `istio-injection: enabled` (sidecar, default). For **ambient** Istio, use overlay [`kind-local-ambient`](../kind-local-ambient/) instead.

## Upstream BoA images instead of GHCR

This overlay **defaults to self-built GHCR images** from this repo’s CI. To use **original** Google Artifact Registry images (no GHCR pull):

1. In `kustomization.yaml`, **remove** the `use-ghcr-images` line from `components:`.
2. Keep (or use) the `images:` block below — it pins upstream tags on `us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/*`.

```bash
kubectl kustomize deploy/overlays/kind-local   # should show us-central1-docker.pkg.dev/... URLs
```

Requires pull access to upstream BoA images only (no GHCR packages).
