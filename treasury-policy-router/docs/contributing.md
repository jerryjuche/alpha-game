# Contributing

## Prerequisites

- Rust 1.84+
- Node 20 LTS
- pnpm 9.x
- Docker

## Local dev

```bash
# Contract
cd contract && cargo build

# App
cd app && pnpm install && pnpm dev
```

## Commit rules

- One commit per logical unit
- Push immediately
- Follow Conventional Commits

## Labels

- `enhancement` — new feature
- `bug` — broken behavior
- `documentation` — docs only
- `complexity:100/150/200` — Wave points estimate
