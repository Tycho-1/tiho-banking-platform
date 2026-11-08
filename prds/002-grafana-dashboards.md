# PRD 002 — Grafana dashboards from Git (Kind)

**Status:** draft

The three dashboard JSON files exist and work when imported by hand. Nothing in Git applies them, so an agent cannot repeat the import and a cluster rebuild loses them.

## Problem

`deploy/dashboards/` holds `banking-platform-metrics.json`, `banking-platform-traces.json`, and `banking-platform-logs.json`. Grafana loads them only after **Dashboards → Import**. The Kind sidecar does not watch this repo.

The sidecar in `tiho-local-kind-cluster` (`kind-cluster/observability/kube-prometheus-stack-values.yaml`) has `grafana.sidecar.dashboards.searchNamespace: observability`. It loads ConfigMaps labeled `grafana_dashboard: "1"` in that namespace. The pattern already used by the platform is `kind-cluster/observability/dashboards/kustomization.yaml` (`configMapGenerator`, `namespace: observability`).

`deploy/overlays/kind-local` sets `namespace: banking-platform`. A ConfigMap added there stays in `banking-platform` and the sidecar never sees it.

## Scope

- The three existing JSON files. Same panels, same dashboard uids (`banking-platform-metrics`, `banking-platform-traces`, `banking-platform-logs`).
- Kind only. Provision them the same way as the other Kind dashboards: ConfigMap in `observability`, label `grafana_dashboard: "1"`.
- Datasource UIDs must be real before the sidecar loads the file. Hand import fills `${DS_PROMETHEUS}` and `${DS_LOKI}` through Grafana's `__inputs` block. The sidecar does not. Use the Kind datasources that already point at this Prometheus and this Loki: `webstore-metrics` and `webstore-logs`.

## Out of scope

- New panels, alerts, or a fourth dashboard
- `gke-dev` and Cloud Ops
- Changing scrape, traces, or log export
- Putting the ConfigMap in the `banking-platform` namespace

## Acceptance

- [ ] The three JSON files are inputs to a `configMapGenerator` in the Kind observability dashboards kustomization (`namespace: observability`, label `grafana_dashboard: "1"`)
- [ ] `__inputs` and `${DS_PROMETHEUS}` / `${DS_LOKI}` are replaced with `webstore-metrics` and `webstore-logs`
- [ ] `kind-local` does not gain a ConfigMap in `banking-platform` for this
- [ ] `gke-dev` is unchanged
- [ ] After apply, `kubectl --context kind-tiho-test get cm -n observability -l grafana_dashboard=1` lists the new ConfigMaps
- [ ] Grafana shows the three dashboards without a manual import, and the existing panels still return data
- [ ] `deploy/dashboards/README.md` describes the Git path and no longer says the only load path is a hand import

## Tasks

- [ ] Copy or generate the three dashboards into `tiho-local-kind-cluster/kind-cluster/observability/dashboards/` and register them in that `kustomization.yaml` (`disableNameSuffixHash: true`, same as the existing entries)
- [ ] Rewrite datasource references for sidecar provisioning
- [ ] Apply with the platform task (`task apply_grafana_dashboards`) or server-side apply if the JSON is too large for a client-side apply annotation
- [ ] Confirm the three uids in Grafana. The provisioned object replaces the hand-imported one with the same uid
