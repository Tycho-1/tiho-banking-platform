---
description: Close an observability PRD — verify acceptance checklist and set status done.
---

# Observability — PRD done

## Steps

1. Read the active PRD path from user (`prds/NNN-*.md`).

2. Run `bash scripts/agent/observability-verify.sh`.

3. Walk that PRD's **Acceptance** checkboxes. Mark `[x]` only with evidence for that item. Grafana, Jaeger, or Loki counts only when the PRD asks for it. Do not close 001 or 003 again. 002 is done only when the sidecar loads the three dashboards without a hand import.

4. Set `**Status:** done` at top.

5. Summarize remaining follow-ups from `observability-review` if any (new PRDs or TODO.md).
