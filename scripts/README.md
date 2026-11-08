# Scripts

Helper scripts for this repo. Application code stays in `src/`. Deploy manifests stay in `deploy/`.

## Agent

| Script | Purpose |
|--------|---------|
| [`agent/observability-verify.sh`](agent/observability-verify.sh) | Check the Kind observability wiring before a commit |

Run it from the repo root:

```bash
bash scripts/agent/observability-verify.sh
```

It passes only when all of these are true:

- `deploy/overlays/kind-local` renders (`kustomize`, or `kubectl kustomize` if `kustomize` is not installed)
- Each of the seven app services has a ServiceMonitor and a `use-otel-otlp` patch: `frontend`, `contacts`, `userservice`, `balancereader`, `transactionhistory`, `ledgerwriter`, `product-catalog`
- `deploy/overlays/gke-dev` does not include `prometheus-servicemonitors` or `use-otel-otlp`

Exit code `0` prints `ALL GATES PASSED`. Anything missing exits `1`.

The script does not query Grafana, Jaeger, or Loki. Those are checked after Flux is running the new image.

Cursor and Claude run this script from their stop hook. The skill `observability-verify` uses the same command. Map: [AGENTS.md](../AGENTS.md). Detail: [docs/agent-observability.md](../docs/agent-observability.md).
