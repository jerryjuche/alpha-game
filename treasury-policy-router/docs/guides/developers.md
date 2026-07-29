# Developer Guide

## Setup

1. Clone the repo.
2. Install Rust 1.84+ and wasm32v1-none target.
3. Install Node 20 and pnpm 9.
4. Run `pnpm install` in `app/`.
5. Deploy PolicyHook to Futurenet using `bash scripts/deploy.sh`.

## SDK

```typescript
import { PolicyHookClient } from '@tpr/sdk';

const client = new PolicyHookClient({
  rpcUrl: process.env.STELLAR_RPC_URL,
  contractId: process.env.POLICY_HOOK_CONTRACT_ID,
});

const result = await client.checkPolicy(from, to, amount);
console.log(result.approved, result.flag, result.reason);
```

## API

```bash
curl "http://localhost:3000/reports/policy-summary?start=2024-01-01&end=2024-01-31"
```
