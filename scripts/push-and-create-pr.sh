#!/bin/bash
set -e

echo "Pushing branch..."
git -C /home/gamp/alpha-game push origin treasury-policy-router

echo "Creating PR..."
gh pr create \
  --repo lekanay2005-coder/alpha-game \
  --base main \
  --head treasury-policy-router \
  --title "feat: add treasury-policy-router scaffold" \
  --body "$(cat <<'EOF'
## Summary
Adds the full treasury-policy-router scaffold from the Stellar Wave builder playbook.

### What's included
- `contract/policy-hook` — Soroban PolicyHook contract with tests and CI
- `app/` — pnpm monorepo scaffold with sdk, indexer, api, web
- `docs/` — plain-English docs with worked numbers
- `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, `LICENSE`
- `scripts/` — setup, deploy, issue-generation scripts
EOF
)"

echo "Done!"
