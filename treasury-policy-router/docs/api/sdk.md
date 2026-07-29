# SDK Reference

## PolicyHookClient

```typescript
class PolicyHookClient {
  constructor(config: { rpcUrl: string; contractId: string })
  checkPolicy(from: string, to: string, amount: string): Promise<PolicyResult>
  getJurisdiction(country: string): Promise<boolean>
  getDailyCap(): Promise<string>
}
```

## MultisigClient

```typescript
class MultisigClient {
  constructor(config: { rpcUrl: string; contractId: string })
  getTransactions(cursor?: string): Promise<TransactionPage>
}
```

## GovernanceClient

```typescript
class GovernanceClient {
  constructor(config: { rpcUrl: string; contractId: string })
  getProposals(cursor?: string): Promise<ProposalPage>
}
```
