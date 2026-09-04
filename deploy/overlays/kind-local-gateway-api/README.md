# Kind local — Gateway API ingress

Bank of Anthos on Kind using **[Kubernetes Gateway API](https://gateway-api.sigs.k8s.io/docs/introduction/)** for north-south traffic instead of classic Istio `Gateway` + `VirtualService`.

| Piece | Location |
|-------|----------|
| Namespace `banking-platform` | `namespace.yaml` (`istio-injection: enabled`) |
| Demo JWT | `deploy/components/demo-jwt/` |
| **Gateway + HTTPRoute** | `deploy/components/gateway-api-ingress/` |
| **Container images** | component `use-ghcr-images` → `ghcr.io/tycho-1/tiho-banking-platform/<service>` (same as [`kind-local`](../kind-local/)) |
| GCP telemetry off | `disable-gcp-telemetry` |
| Frontend ClusterIP | `frontend-clusterip` |
| Prometheus scrape | `prometheus-servicemonitors` (product-catalog `/metrics`, balancereader `/actuator/prometheus`) |
| OTLP traces | `use-otel-otlp` (frontend, product-catalog, balancereader → Collector) |

## Why this overlay

- **Portable API:** `Gateway` / `HTTPRoute` / `GatewayClass` — how many companies expose apps on Kubernetes ([Gateway API intro](https://gateway-api.sigs.k8s.io/docs/introduction/)).
- **Implementation here:** `gatewayClassName: istio` (requires a `GatewayClass` named `istio` on the cluster — typical when Istio is installed with Gateway API support). Istio runs the data plane; you write standard K8s resources.
- **Not “no Istio”:** mesh sidecars still inject. For mesh-off + port-forward only, that would be a separate thin overlay later.
- **Not classic Ingress:** we intentionally skip `networking.k8s.io/Ingress`.

Compare:

| Overlay | Ingress model |
|---------|----------------|
| `kind-local` | Istio `networking.istio.io` Gateway + VirtualService → Helm `istio-ingress` |
| **`kind-local-gateway-api`** | Gateway API → Istio-managed Gateway Service (new LB) |
| `kind-local-ambient` | Ambient + waypoint (+ classic VS ingress from parent) |

## Prerequisites

- Kind cluster with **Istio** and **Gateway API CRDs** (and a `GatewayClass`, e.g. `istio`)
- `GatewayClass` named `istio` present: `kubectl get gatewayclass`
- MetalLB (or similar) so the Gateway `Service` gets an EXTERNAL-IP

## Install

```bash
kubectl apply -k deploy/overlays/kind-local-gateway-api

kubectl wait --for=condition=available deployment --all -n banking-platform --timeout=300s
kubectl wait --for=condition=Programmed gateway/banking-platform -n banking-platform --timeout=120s
```

## Access

Istio creates a LoadBalancer Service for the Gateway (often named like the Gateway). Find the IP:

```bash
kubectl get gateway,httproute,svc -n banking-platform
# Look for EXTERNAL-IP on the gateway Service (MetalLB)
curl -sS -o /dev/null -w "%{http_code}\n" http://<EXTERNAL-IP>/
```

Open `http://<EXTERNAL-IP>/` in a browser. Demo login: `testuser` / `bankofanthos`.

## Upstream BoA images instead of GHCR

Remove `use-ghcr-images` from `components:` in `kustomization.yaml`. The overlay `images:` block then pins upstream tags (same pattern as [`kind-local`](../kind-local/README.md)).

Optional hostname later: set `spec.listeners[].hostname` / `HTTPRoute` `hostnames` and point DNS at the LB IP (TLS via cert-manager when ready).

## Delete

```bash
kubectl delete -k deploy/overlays/kind-local-gateway-api
```

Do not mix with `kind-local` in the same namespace (both expect to own ingress differently). Pick one overlay at a time.
