#!/usr/bin/env bash
set -euo pipefail
ROOT="${CLAUDE_PROJECT_DIR:-$(cd "$(dirname "$0")/../.." && pwd)}"
cd "$ROOT"
if ! bash scripts/agent/observability-verify.sh 2>&1; then
  echo "Stop hook: observability-verify failed." >&2
  exit 2
fi
exit 0
