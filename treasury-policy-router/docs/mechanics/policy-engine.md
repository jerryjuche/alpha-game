# Policy Engine

The PolicyHook contract evaluates three rules on every check:

1. Daily cap: sum of outgoing amounts per sender per UTC day.
2. Jurisdiction allowlist: receiver country must be in allowlist.
3. Travel Rule: transactions above threshold emit a flag.

All amounts use i128. Basis points are represented as integer values (1e7 = 100%).

## Worked example

- Daily cap: 1_000_000_0000 (1000 XLM in stroops, assuming 7 decimals)
- Travel Rule threshold: 100_000_000 (10 XLM)
- Amount: 150_000_000
- Result: approved=true, flag=true, reason=travel_rule_flag
