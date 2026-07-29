# Operator Guide

This guide is for treasury operators running multisig vaults or DAO governors.

## Configuring policies

1. Deploy PolicyHook.
2. Call `set_daily_cap` with your desired cap.
3. Call `set_travel_rule_threshold` with your threshold.
4. Call `set_jurisdiction` for any jurisdictions you want to block.

## Reading reports

Use the dashboard at `apps/web` or the `/reports/policy-summary` endpoint.

## Worked example

- Cap: 1_000_000_0000 stroops/day
- Threshold: 100_000_000 stroops
- Sender sends 150_000_000 stroops to blocked jurisdiction (IR)
- Result: rejected, reason=jurisdiction_blocked
