---
description: Add Prometheus metrics and ServiceMonitor for one banking microservice on kind-local. Use when user asks to scrape a new service in Grafana.
---

# Observability — add Prometheus

Add **metrics + ServiceMonitor** for **one service** on **`kind-local`** overlays (inherits to ambient/CNPG variants).

## Before coding

Confirm with user:

- **Service name**. The seven app services already have ServiceMonitors. This skill is for a service that still does not.
- **Stack:** Java Spring → `/actuator/prometheus`; Go → `/metrics`; Python → `/metrics` + `prometheus_client`
- **Reference to copy:** `transactionhistory` (Java), `product-catalog` (Go), or `frontend` (Python)

## Steps

1. Read reference service **source** + existing ServiceMonitor in `deploy/components/prometheus-servicemonitors/`.

2. **Application metrics**
   - Java: Micrometer/Actuator (often already present — verify endpoint)
   - Go: expose `/metrics`
   - Python: follow `src/frontend/metrics.py` pattern

3. **ServiceMonitor** — create `servicemonitor-<service>.yaml`:
   - Labels: `release: kube-prometheus-stack` (match cluster Prometheus)
   - `selector.matchLabels.app` must match Service in `deploy/base/`
   - Correct `path` and port `http`

4. Register in `deploy/components/prometheus-servicemonitors/kustomization.yaml`.

5. If the image changed: bump `src/<service>/release.yaml`, `newTag` in `deploy/components/use-ghcr-images/kustomization.yaml`, and the overlay pins (`kind-local` and `gke-dev`). `gke-dev` gets the tag only. Do not add ServiceMonitors there.

6. Run:

   ```bash
   bash scripts/agent/observability-verify.sh
   ```

7. Tell user **post-deploy checks:**
   - CI green
   - Flux synced → new image
   - Grafana: query metrics (see `deploy/components/prometheus-servicemonitors/README.md`)

## Do not

- Add ServiceMonitor to `gke-dev` without explicit user request
- Assume scrape works without matching Service labels (frontend needed `app: frontend`)
