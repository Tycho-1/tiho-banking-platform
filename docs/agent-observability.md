# Agent workflow — observability

How skills, hooks, and PRDs apply now that the Kind baseline is on `main`.

Entry point: [AGENTS.md](../AGENTS.md)

---

## Baseline

These seven services scrape metrics, export OTLP traces, and export OTLP logs when `OTEL_EXPORTER_OTLP_ENDPOINT` is set (Kind):

| Service | Metrics | Traces and logs |
|---------|---------|-----------------|
| `frontend`, `contacts`, `userservice` | `metrics.py`, `GET /metrics` | `_setup_tracing`, `_setup_otlp_logs` |
| `balancereader`, `transactionhistory`, `ledgerwriter` | Actuator `GET /actuator/prometheus` | Spring OTLP traces, `OtlpLogConfig` + Log4j2 |
| `product-catalog` | `metrics.go`, `GET /metrics` | `tracing.go`, `logs.go` |

Kind overlay `deploy/overlays/kind-local` includes `prometheus-servicemonitors` and `use-otel-otlp`. Inherited ambient and CNPG overlays pick both up. `gke-dev` includes neither.

Python OTLP logs also require `ENABLE_TRACING=true`. Java and Go logs turn on from the endpoint alone. Python and Java log records can carry the trace id. Go `log.Printf` cannot.

Dashboards in `deploy/dashboards/` are imported by hand: metrics, trace rates, pod logs. Jaeger is the waterfall.

---

## Skills

| Skill | When |
|-------|------|
| **observability-add-prom** | A service that still has no ServiceMonitor |
| **observability-add-otel** | A service that still has no Kind OTLP patch |
| **observability-verify** | Before commit |
| **observability-review** | After you have looked at Grafana, Jaeger, or Loki |
| **observability-prd-create** | Multi-step initiative |
| **observability-prd-done** | Close a PRD |

### Example session (next slice, not a missing service)

```
User: Add a ledgerwriter result counter (accepted, rejected, duplicate) on kind-local.

Agent: copy the Java metrics path, bump the image tag, observability-verify
User: deploy, check Grafana
Agent: observability-review — file follow-ups, do not add alerts in the same change
```

---

## Hooks

| Hook | Behavior |
|------|----------|
| Block push to `main` | Feature branch + PR. Cursor: `.cursor/hooks.json`. Claude: `.claude/settings.json` |
| Stop | Runs `scripts/agent/observability-verify.sh` |

`.cursor/` and `.claude/` are both committed. Keep the skill text the same in both trees. Hook scripts differ because each tool reads different JSON.

---

## Your workflow mapped

| Step | Agent support |
|------|----------------|
| Define the slice | You still define service, signal, and overlay |
| Look at code | Skills point at frontend, balancereader, or product-catalog |
| Commit | `observability-verify` first |
| CI, Flux, image | Unchanged. Remind to confirm the running tag |
| Grafana / Jaeger / Loki | `observability-review` only after you confirm the UI |

---

## Reference services

| Copy from | For |
|-----------|-----|
| `balancereader` | Java Prom, OTLP traces, OTLP logs |
| `frontend` | Python Prom, traces, logs, outbound backend metrics |
| `product-catalog` | Go `/metrics`, traces, logs |
| `contacts` | Python OTLP patch shape without frontend’s backend counters |
