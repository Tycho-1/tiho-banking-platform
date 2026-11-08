#!/usr/bin/env bash
set -ueo pipefail
input=$(cat)
command=$(echo "$input" | jq -r '.tool_input.command // empty')
if [[ -z "$command" ]]; then exit 0; fi
if echo "$command" | grep -Eq 'git push.*(origin|upstream)?[[:space:]]*(main|master)\b'; then
  jq -n '{
    hookSpecificOutput: {
      hookEventName: "PreToolUse",
      permissionDecision: "deny",
      permissionDecisionReason: "Push to main blocked. Use a feature branch or disable hook for solo work."
    }
  }'
  exit 0
fi
exit 0
