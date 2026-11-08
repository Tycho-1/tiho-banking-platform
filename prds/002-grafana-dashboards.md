# PRD 002 — Grafana dashboards from Git (Kind)

**Status:** draft

The three dashboard JSON files exist and work when imported by hand. Nothing in Git applies them, so an agent cannot repeat the import and a cluster rebuild loses them.

## Problem

`deploy/dashboards/` holds `banking-platform-metrics.json`, `banking-platform-traces.json`, and `banking-platform-logs.json`. Grafana loads them only after **Dashboards → Import**. Flux does not apply these files today.

GitOps path: Flux already reconciles `flux-fleet` `infrastructure/observability/dashboards/`. That kustomization builds a ConfigMap in namespace `observability` with label `grafana_dashboard: "1"`. Grafana does not read Git. A container in the Grafana pod (the chart's dashboard sidecar, not the Istio sidecar) watches that namespace and loads the JSON. Flux applies the ConfigMap. The container picks it up.

`deploy/overlays/kind-local` sets `namespace: banking-platform`. A ConfigMap added there stays in `banking-platform` and Grafana never sees it. Do not use `task apply_grafana_dashboards` for this. That task is a one-off apply in the Kind lab repo, not this cluster's GitOps.

## Scope

- The three existing JSON files. Same panels, same dashboard uids (`banking-platform-metrics`, `banking-platform-traces`, `banking-platform-logs`).
- Kind only. Flux applies them from `flux-fleet` `infrastructure/observability/dashboards/`: ConfigMap in `observability`, label `grafana_dashboard: "1"`.
- Datasource UIDs must be real in the JSON. Hand import fills `${DS_PROMETHEUS}` and `${DS_LOKI}` through Grafana's `__inputs` block. A ConfigMap load does not. Use the Kind datasources that already point at this Prometheus and this Loki: `webstore-metrics` and `webstore-logs`.

## Out of scope

- New panels, alerts, or a fourth dashboard
- `gke-dev` and Cloud Ops
- Changing scrape, traces, or log export
- Putting the ConfigMap in the `banking-platform` namespace

## Acceptance

- [ ] The three JSON files are inputs to the Flux dashboards kustomization `flux-fleet/infrastructure/observability/dashboards/` (`namespace: observability`, label `grafana_dashboard: "1"`)
- [ ] `__inputs` and `${DS_PROMETHEUS}` / `${DS_LOKI}` are replaced with `webstore-metrics` and `webstore-logs`
- [ ] `kind-local` does not gain a ConfigMap in `banking-platform` for this
- [ ] `gke-dev` is unchanged
- [ ] After apply, `kubectl --context kind-tiho-test get cm -n observability -l grafana_dashboard=1` lists the new ConfigMaps
- [ ] Grafana shows the three dashboards without a manual import, and the existing panels still return data
- [ ] `deploy/dashboards/README.md` describes the Git path and no longer says the only load path is a hand import

## Tasks

- [ ] Add the three dashboards to `flux-fleet/infrastructure/observability/dashboards/` and register them on the `configMapGenerator` (`disableNameSuffixHash: true`)
- [ ] Rewrite datasource references so the JSON does not depend on the import dialog
- [ ] Let Flux reconcile. Do not `kubectl apply` or `task apply_grafana_dashboards` as the way to keep them
- [ ] Confirm the three uids in Grafana. The Flux-loaded object replaces the hand-imported one with the same uid
