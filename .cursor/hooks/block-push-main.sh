#!/usr/bin/env bash
set -euo pipefail
input=$(cat)
command=$(echo "$input" | jq -r '.command // empty')
if [[ -z "$command" ]]; then echo '{"permission":"allow"}'; exit 0; fi
if echo "$command" | grep -Eq 'git push.*(origin|upstream)?[[:space:]]*(main|master)\b'; then
  jq -n '{
    "permission": "deny",
    "user_message": "Push to main blocked. Use a feature branch or disable this hook in .cursor/hooks.json for solo work.",
    "agent_message": "Do not push to main. Use a feature branch and open a PR."
  }'
  exit 0
fi
echo '{"permission":"allow"}'
exit 0
