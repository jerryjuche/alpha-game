# Post-Approval Issues

## Contract issues

1. feat(contract): add Travel Rule amount splitting
   Depends on app/api: POST /reports/travel-rule-split
   
2. feat(contract): add set_admin with auth guard
   Depends on none

3. test(contract): add property tests for daily cap overflow
   Depends on none

4. feat(contract): add compliance report export event
   Depends on app/indexer: parse new event type

## App issues

5. feat(sdk): add PolicyHookTypeScript client
   Depends on contract: deployed to testnet

6. feat(indexer): cursor-based pagination for contract events
   Depends on none

7. feat(api): add POST /reports/policy-summary
   Depends on indexer: event data available

8. feat(web): add policy dashboard with jurisdiction map
   Depends on api: endpoints ready

9. chore: add Docker compose for local dev
   Depends on none

10. feat(api): add webhook delivery for policy flags
    Depends on indexer: flag events available

11. docs: add operator guide with worked numbers
    Depends on none

12. docs: add API reference with cURL examples
    Depends on api: endpoints implemented
