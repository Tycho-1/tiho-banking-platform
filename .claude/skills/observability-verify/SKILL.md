---
description: Run observability pre-commit gates and review diff. Use before commit or after prom/otel edits.
---

# Observability — verify

## Steps

1. Run:

   ```bash
   bash scripts/agent/observability-verify.sh
   ```

2. Review `git diff`. Confirm:
   - Only intended service + manifests touched
   - ServiceMonitor selector matches Service labels
   - Overlay `newTag` bumped if app image changed
   - No secrets in diff

3. Report **pass** or **fail** with file list and any scrape/OTel risks.

4. The script renders `kind-local`, requires the seven ServiceMonitors and seven OTLP patches, and fails if `gke-dev` includes those components.

5. Remind the user: after merge and deploy, Grafana, Jaeger, and Loki are the proof. The script does not check them.
