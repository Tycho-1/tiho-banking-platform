---
description: Create a PRD for a multi-step observability rollout (Prom, OTel, or both). Does not implement code.
---

# Observability — PRD create

Create `prds/NNN-<slug>.md` for an observability initiative spanning one or more services.

## Steps

1. Get a title if missing. 001 (Kind signals) and 003 (product-catalog) are done. 002 (Grafana dashboards) is the open draft. Do not create another PRD for those.

2. Next ID: highest `prds/*.md` number + 1. That is **004** while 001–003 exist.

3. Write file:

```markdown
# PRD NNN — {title}

**Status:** draft

## Problem

{why observability needed for these services}

## Scope

- Services: ...
- Overlay: kind-local (default)
- Prom / OTel / both

## Out of scope

- gke-dev Cloud Ops migration
- ...

## Acceptance

- [ ] {proof that matches this slice — Grafana, Jaeger, Loki, a page, or a kustomize build. Not a copy of another PRD}

## Tasks

- [ ] ...
```

4. Tell user: implement with `observability-add-prom` / `observability-add-otel`, close with `observability-prd-done`.
