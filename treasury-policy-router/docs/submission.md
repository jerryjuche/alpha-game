# Submission Package

## Duplicate check
Confirmed via drips.network/wave/stellar/repos and github.com/stellar-wave that no approved repo implements treasury policy monitoring.

## Supporting links
- Live app URL: https://treasury-policy-router.vercel.app
- Repo 1: https://github.com/yourname/treasury-policy-router-contract
- Repo 2: https://github.com/yourname/treasury-policy-router
- Docs site: https://treasury-policy-router.vercel.app/docs
- On-chain verification: https://stellar.expert/explorer/public/contract/<CONTRACT_ID>
- Demo video: https://youtube.com/watch?v=...

## Repo relationship description
treasury-policy-router-contract is a lightweight Soroban hook contract. treasury-policy-router is the app layer (indexer, REST API, dashboard) that watches events from policy-hook, multisig vaults, and DAO governors. The contract is optional; the app can index transactions even without it, but the hook enables real-time on-chain enforcement.

## Planned issues description
Generated 12 issues across the two repos: 4 contract issues (config, policy logic, events, tests), 6 app issues (SDK, indexer, API, reports, CI, Docker), 2 docs issues (operator guide, API reference).

## Project description
A surveillance camera and bookkeeper for blockchain treasuries. It watches Soroban multisig vault and DAO governor transactions in real time, applies configurable treasury policies, and generates compliance-ready audit reports without requiring changes to existing approved contracts.
