# Kind local overlay

| Piece | Location |
|-------|----------|
| Namespace `banking-platform` | `namespace.yaml` |
| Demo JWT secret | `deploy/components/demo-jwt/` |
| Istio Gateway + VirtualService | `istio/gateway-helm.yaml` (default) |
| Image tags | `kustomization.yaml` → `images:` |
| GCP telemetry off | component `disable-gcp-telemetry` |
| Frontend ClusterIP | component `frontend-clusterip` |

**Mesh:** namespace label `istio-injection: enabled` (sidecar, default). For **ambient** Istio, use overlay [`kind-local-ambient`](../kind-local-ambient/) instead.
