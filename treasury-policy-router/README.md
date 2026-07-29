# Treasury Policy Router

[![Stellar](https://img.shields.io/badge/Stellar-Soroban-blue)](https://stellar.org/soroban)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Drips Wave](https://img.shields.io/badge/Drips-Wave_Program-blue)](https://drips.network/wave/stellar)

> A surveillance camera and bookkeeper for blockchain treasuries. Watches Soroban multisig vault and DAO governor transactions in real time, applies configurable treasury policies, and generates compliance-ready audit reports.

## Maintainer

| Name | GitHub | Email | Discord |
|------|--------|-------|---------|
| Your Name | [@yourname](https://github.com/yourname) | you@example.com | yourname#0000 |

## Community

- [Discussions](https://github.com/yourname/treasury-policy-router/discussions)
- [Discord](https://discord.gg/yourinvite)

## Architecture

treasury-policy-router has two parts:
1. `policy-hook` — a lightweight Soroban contract that exposes `check_policy` for real-time treasury policy enforcement.
2. `treasury-policy-router` — the app layer (indexer, REST API, dashboard) that watches events from `policy-hook`, multisig vaults, and DAO governors.

## Quick Start

### Contract

```bash
cd contract
cargo build --target wasm32v1-none
bash scripts/setup.sh
bash scripts/deploy.sh
```

### App

```bash
cd app
cp .env.example .env.local
pnpm install
docker compose up -d
pnpm dev
```

## Documentation

See [docs/](docs/) for full documentation.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE)
