#!/bin/bash
set -e
REPO_CONTRACT="yourname/treasury-policy-router-contract"
REPO_APP="yourname/treasury-policy-router"

create_issue() {
  local repo=$1
  local title=$2
  local labels=$3
  local body=$4
  gh issue create --repo "$repo" --title "$title" --label "$labels" --body "$body"
}

echo "Creating contract issues..."
create_issue "$REPO_CONTRACT" \
  "feat(contract): add set_admin with auth guard" \
  "enhancement,medium" \
  "Summary: allow admin transfer with require_auth. Acceptance Criteria: - [ ] function exists - [ ] emits PolicyUpdated - [ ] old admin is invalidated Tech Stack: Rust, Soroban SDK"

create_issue "$REPO_CONTRACT" \
  "feat(contract): add Travel Rule amount splitting" \
  "enhancement,high" \
  "Summary: split transactions above threshold. Acceptance Criteria: - [ ] returns split amounts - [ ] emits two PolicyChecked events Tech Stack: Rust, Soroban SDK"

create_issue "$REPO_CONTRACT" \
  "test(contract): add property tests for daily cap overflow" \
  "documentation,medium" \
  "Summary: ensure daily cap cannot overflow i128. Acceptance Criteria: - [ ] proptest passes for 1_000 iterations - [ ] edge cases documented Tech Stack: Rust, proptest"

create_issue "$REPO_CONTRACT" \
  "feat(contract): add compliance report export event" \
  "enhancement,medium" \
  "Summary: emit structured report event. Acceptance Criteria: - [ ] event includes from, to, amount, result - [ ] indexer parses without custom logic Tech Stack: Rust, Soroban SDK"

echo "Creating app issues..."
create_issue "$REPO_APP" \
  "feat(sdk): add PolicyHookTypeScript client" \
  "enhancement,medium" \
  "Summary: typed wrapper for contract calls. Acceptance Criteria: - [ ] check_policy method - [ ] retry on 5xx - [ ] typed return Tech Stack: TypeScript, @stellar/stellar-sdk"

create_issue "$REPO_APP" \
  "feat(indexer): cursor-based pagination for contract events" \
  "enhancement,medium" \
  "Summary: paginate getTransactions with cursor. Acceptance Criteria: - [ ] cursor state stored in Redis - [ ] backoff on RPC error - [ ] idempotent processing Tech Stack: TypeScript, PostgreSQL, Redis"

create_issue "$REPO_APP" \
  "feat(api): add POST /reports/policy-summary" \
  "enhancement,high" \
  "Summary: generate JSON report of policy checks. Acceptance Criteria: - [ ] date range filter - [ ] CSV export - [ ] paginated list Tech Stack: TypeScript, Hono, Zod"

create_issue "$REPO_APP" \
  "feat(web): add policy dashboard with jurisdiction map" \
  "enhancement,high" \
  "Summary: UI for viewing policy checks. Acceptance Criteria: - [ ] filter by date, result - [ ] chart for flags - [ ] mobile layout Tech Stack: React, Tailwind, Recharts"

create_issue "$REPO_APP" \
  "chore: add Docker compose for local dev" \
  "enhancement,low" \
  "Summary: one-command local env. Acceptance Criteria: - [ ] postgres + redis - [ ] npm run dev starts all - [ ] documented in README Tech Stack: Docker, docker-compose"

create_issue "$REPO_APP" \
  "feat(api): add webhook delivery for policy flags" \
  "enhancement,medium" \
  "Summary: POST flagged checks to webhook URL. Acceptance Criteria: - [ ] retry with backoff - [ ] signature header - [ ] configurable per org Tech Stack: TypeScript, Hono"

create_issue "$REPO_APP" \
  "docs: add operator guide with worked numbers" \
  "documentation,medium" \
  "Summary: plain-English guide + examples. Acceptance Criteria: - [ ] daily cap example - [ ] Travel Rule math - [ ] screenshots Tech Stack: Markdown"

create_issue "$REPO_APP" \
  "docs: add API reference with cURL examples" \
  "documentation,medium" \
  "Summary: every endpoint documented. Acceptance Criteria: - [ ] request/response JSON - [ ] error codes - [ ] auth flow Tech Stack: Markdown"
