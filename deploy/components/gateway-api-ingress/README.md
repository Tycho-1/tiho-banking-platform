# gateway-api-ingress

Kubernetes **[Gateway API](https://gateway-api.sigs.k8s.io/docs/introduction/)** north-south entry for Bank of Anthos:

| Resource | Role |
|----------|------|
| `Gateway` `banking-platform` | Listener :80, `gatewayClassName: istio` |
| `HTTPRoute` `frontend` | `/` → Service `frontend:80` |

**Implementation:** Istio’s `GatewayClass` `istio` (must exist on the cluster). This is the portable *API*; Istio is one of many possible controllers.

Used by overlay `kind-local-gateway-api`.
