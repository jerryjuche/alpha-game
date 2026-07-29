# Development Setup

## Prerequisites

- Rust 1.84+
- Node.js 20 LTS
- pnpm 9.x
- Docker and Docker Compose
- `stellar` CLI
- `gh` CLI

## Contract

```bash
cd contract
cargo build --target wasm32v1-none
bash scripts/setup.sh
bash scripts/deploy.sh
cargo test
```

## App

```bash
cd app
cp .env.example .env.local
pnpm install
docker compose up -d
pnpm dev
```

## Generate Issues

```bash
bash scripts/generate-issues.sh
```
