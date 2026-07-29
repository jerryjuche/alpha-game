# Service Topology

```
User → Frontend (Vercel) → API (Render) ←→ PostgreSQL (Render, same region)
                                    ↓
                            Soroban RPC (public/Validation Cloud)
                                    ↓
                         PolicyHook + multisig + DAO contracts
```

## Service linkage

- Frontend calls API via `NEXT_PUBLIC_API_URL` (never `localhost` in prod)
- API calls indexer via internal Render service link or same-process call
- Indexer calls Soroban RPC directly; no API hop
- Database connection string uses Render internal hostname when co-located

## Failure trace

If frontend hits `localhost:3000` in prod, root cause is missing `NEXT_PUBLIC_API_URL` env var. Fix at Vercel env UI, not in code.
