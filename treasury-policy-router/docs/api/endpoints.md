# API Endpoints Reference

## GET /reports/policy-summary

Response:
```json
{
  "total": 42,
  "approved": 38,
  "rejected": 2,
  "flagged": 2
}
```

## GET /reports/:contractId/flags

Response:
```json
[
  {
    "from": "G...",
    "to": "G...",
    "amount": "150000000",
    "reason": "travel_rule_flag"
  }
]
```

## POST /webhooks/policy-flag

Body:
```json
{
  "url": "https://example.com/hook",
  "secret": "supersecret"
}
```

Headers sent:
```
X-TPR-Signature: sha256=...
```
