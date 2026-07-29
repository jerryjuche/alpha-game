# API Overview

Base URL: `http://localhost:3000`

All routes return JSON. Errors use standard HTTP status codes.

## Endpoints

### GET /reports/policy-summary

Returns summary of policy checks.

**Query params:**
- `start` — ISO date
- `end` — ISO date
- `contractId` — optional PolicyHook contract ID

### GET /reports/:contractId/flags

Returns flagged transactions.

### POST /webhooks/policy-flag

Registers a webhook URL to receive flagged checks.
