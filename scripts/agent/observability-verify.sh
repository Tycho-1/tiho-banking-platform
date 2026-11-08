#!/usr/bin/env bash
# Agent verify — kind-local renders, the seven app signals are wired, gke-dev is not.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OVERLAY="${ROOT}/deploy/overlays/kind-local"
PROM_DIR="${ROOT}/deploy/components/prometheus-servicemonitors"
OTLP_DIR="${ROOT}/deploy/components/use-otel-otlp"
GKE="${ROOT}/deploy/overlays/gke-dev/kustomization.yaml"

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

cd "$ROOT"

echo "== render deploy/overlays/kind-local =="
if command -v kustomize >/dev/null 2>&1; then
  kustomize build "$OVERLAY" >/dev/null
elif command -v kubectl >/dev/null 2>&1; then
  kubectl kustomize "$OVERLAY" >/dev/null
else
  fail "neither kustomize nor kubectl is in PATH"
fi
echo "render: OK"

echo "== ServiceMonitors =="
monitors=(frontend contacts userservice balancereader transactionhistory ledgerwriter product-catalog)
for name in "${monitors[@]}"; do
  f="${PROM_DIR}/servicemonitor-${name}.yaml"
  [[ -f "$f" ]] || fail "missing ${f}"
  echo "  servicemonitor-${name}.yaml"
done

echo "== use-otel-otlp patches =="
for name in "${monitors[@]}"; do
  f="${OTLP_DIR}/patch-${name}.yaml"
  [[ -f "$f" ]] || fail "missing ${f}"
  echo "  patch-${name}.yaml"
done

echo "== gke-dev stays on Cloud Ops =="
if grep -E 'prometheus-servicemonitors|use-otel-otlp' "$GKE" >/dev/null; then
  fail "gke-dev kustomization includes a Kind observability component"
fi
echo "gke-dev: OK"

echo "ALL GATES PASSED"
