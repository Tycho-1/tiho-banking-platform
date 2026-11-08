# Tiho Banking Platform — agent map

Cursor / Claude Code entry point.

**Kind baseline (done on `main`):** the seven app services have Prometheus scrape, OTLP traces, and OTLP logs. Services: `frontend`, `contacts`, `userservice`, `balancereader`, `transactionhistory`, `ledgerwriter`, `product-catalog`. `loadgenerator` and the databases are out of this set. `gke-dev` stays on GCP Cloud Ops.

A new task is one slice on top of that baseline (an alert, a business counter, a missing span, Go trace ids on log lines, dashboard GitOps). It is not “add Prom to ledgerwriter.”

---

## Commands

| Command | Purpose |
|---------|---------|
| `bash scripts/agent/observability-verify.sh` | Kind overlay builds, seven monitors and OTLP patches exist, `gke-dev` does not include them |
| `kubectl kustomize deploy/overlays/kind-local` | Render the Kind overlay |

After a deploy (human): confirm Flux synced the image, then check Grafana and Jaeger. A diff is not runtime proof.

---

## Observability layout

| Area | Path |
|------|------|
| ServiceMonitors (Kind) | `deploy/components/prometheus-servicemonitors/` |
| OTLP env patches (Kind) | `deploy/components/use-otel-otlp/` |
| Kind overlay (both on) | `deploy/overlays/kind-local/` |
| Dashboards (import by hand) | `deploy/dashboards/` |
| **Not on `gke-dev`** | GCP Cloud Ops — no ServiceMonitor or OTLP component |

### Reference services (copy these)

| Service | Prom path | Stack |
|---------|-----------|--------|
| product-catalog | `/metrics` | Go — `metrics.go`, `tracing.go`, `logs.go` |
| frontend | `/metrics` | Python — `src/frontend/metrics.py`, `_setup_tracing` / `_setup_otlp_logs` |
| balancereader | `/actuator/prometheus` | Java — Actuator, Spring OTLP traces, `OtlpLogConfig` + Log4j2 |

Contacts and userservice follow frontend. Transactionhistory and ledgerwriter follow balancereader. **ServiceMonitor** `spec.selector.matchLabels.app` must match the Service.

Collector (Kind): `http://otel-collector.observability.svc:4318`. Java traces use that URL plus `/v1/traces`.

---

## Agent skills

| Skill | When |
|-------|------|
| `observability-add-prom` | A service that does not yet have a ServiceMonitor |
| `observability-add-otel` | OTLP traces and logs for a service that does not yet have a Kind patch |
| `observability-verify` | Before commit, after Prom or OTel edits |
| `observability-review` | After Grafana/Jaeger/Loki look right — gaps beyond the baseline |
| `observability-prd-create` | A multi-step initiative |
| `observability-prd-done` | Close a PRD when acceptance is met |

Guide: [docs/agent-observability.md](docs/agent-observability.md)

---

## Workflow

1. **Define** the slice: which signal, which service, which overlay (`kind-local` unless the user says otherwise).
2. **Implement** by copying the reference service in that language.
3. **`observability-verify`** before commit.
4. Commit → CI → Flux → confirm the deployed image.
5. Grafana / Jaeger / Loki — runtime proof.
6. **`observability-review`** for follow-ups. Do not widen the current change.

---

## Hooks

- Block push to `main` (feature branch + PR). Cursor: `.cursor/hooks.json`. Claude: `.claude/settings.json`.
- Stop hook runs `scripts/agent/observability-verify.sh`.

`.cursor/` is the Cursor copy (rules, skills, hooks). `.claude/` is the Claude copy. Both are committed. Skill text should stay the same in both. Hook scripts differ because each tool reads different JSON.

---

## PRDs

`prds/*.md`. [001](prds/001-microservices-prom-otel.md) is the closed Kind baseline (Prometheus scrape, OTLP traces, OTLP logs for the seven app services). There is not one PRD per service. [002](prds/002-grafana-dashboards.md) is the open draft: load the three dashboards through the Kind Grafana sidecar. [003](prds/003-product-catalog.md) is done: the product-catalog service, the Save & invest page, and the Kind deploy. Other follow-ups are not files yet: alerts, a ledger result counter, Go log trace ids.

---

## Still open after the baseline

- Alerts on scrape down, HTTP 5xx, and p95
- Ledgerwriter counter for accepted / rejected / duplicate
- Go `log.Printf` lines do not carry a trace id (Python and Java OTLP logs can)
- Dashboard JSON is imported by hand. Git load is [PRD 002](prds/002-grafana-dashboards.md)
