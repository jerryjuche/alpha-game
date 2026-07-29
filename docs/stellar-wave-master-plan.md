# AlphaBlitz Stellar Wave Submission — Master Plan

> Generated from live ecosystem research. Treat this as a real submission with real stakes.

---

## Phase 1 — Ecosystem Reconnaissance

### Stellar Stack (Current, verified 2026-07)
- **Consensus**: Stellar network, ~150 TPS, 5-second finality, Protocol 23 (parallel execution, lower latency).
- **Smart Contracts**: Soroban (Rust → WASM). Deterministic gas metering, local testing harness, upgradeable contract patterns (SEP-49 draft).
- **DEX / AMM**: Native Stellar DEX (order book), Soroswap (AMM + aggregator), Blend Capital (lending pools on Soroban), DeFindex (tokenized vaults / SEP-56).
- **Token Standards**: SEP-41 (fungible token interface, e.g., USDC SAC), SEP-50 (NFT draft), SEP-56 (vaults / ERC-4626), SEP-57 (T-REX RWA).
- **Auth / Identity**: SEP-10 (web auth for classic accounts), SEP-45 (web auth for contract accounts / passkeys), CAP-0051 (secp256r1 verification).
- **Tooling**: Soroban CLI, `stellar-sdk` (TypeScript/JS), `soroban-sdk` (Rust), Freighter wallet, Albedo, WalletConnect, Stellar Expert, Horizon REST API, Soroban RPC, OpenZeppelin Stellar Contracts, golang-migrate (for off-chain DB).

### Drips Wave Stellar Program (Live data)
- **Total approved repos**: 594 across 317 orgs (fetched 2026-07-28 from drips.network/wave/stellar/repos).
- **Current Wave**: Wave 7 launched July 23, 2026. Waves run 7-day sprints, monthly.
- **Constraints**: Max 5 repo applications per Wave per user/org. KYC required. Per-repo points budgets enforced (Wave 4+). Org-level point caps introduced (Wave 5).
- **Approval pattern**: Maintainers apply repos → SDF/Drips review → approved repos get `Stellar Wave` label on issues. Issues have point values (e.g., 150–200 pts) based on complexity.

### SDF Funding Priorities (2026)
- **Mandate**: $1B network asset value growth, 15 new transformational enterprise partners, 5 live deployments in payments/treasury/settlement.
- **SCF 7.0 tracks**: Open (novel on-chain use cases), Integration (build on existing building blocks), RFP (developer tooling gaps).
- **Key verticals**: DeFi (Blend, DeFindex), RWA (T-REX SEP-57, YieldVault), escrow/infrastructure (Trustless Work), payments (Stellopay, Paymesh, Stellar Stream), AI agents (Stellarmind, Talos), privacy (Shielded Protocol), gaming (Stellarcade).

### Approved Repo Landscape & White Space
| Domain | Saturation | Examples |
|--------|-----------|----------|
| Escrow / Payments | **High** | Trustless Work, SafeTrust, PayStell, Paymesh, WaveMilestone |
| DeFi / Yield / RWA | **High** | Blend, DeFindex, YieldVault, Shielded, Neko-Protocol |
| Marketplace / Crowdfunding | **Moderate** | OfferHub, Stellar-Rent, Boundless, PrediFi, Mercato |
| Streaming / Subscriptions | **Low** | Stellar Stream, SubTrackr |
| AI / Agents | **Low–Moderate** | Stellarmind, Talos, Hazina, Galaxy-DevKit |
| Gaming | **Very Low** | Stellarcade (arcade/prize pools, 2x pts) |
| Consumer / Social | **Very Low** | MentorMinds, MindVault |
| Tooling / DevEx | **Moderate** | stellar-dev-skill, Stellar-wave-hub, stellar-portfolio-rebalancer |

**White space identified**: Real-time **consumer multiplayer games** with on-chain prize settlement. Stellarcade exists but focuses on provably fair arcade mechanics; no casual social word game exists. This maps to SDF's user-adoption goal without competing head-on with DeFi primitives.

---

## Phase 2 — Idea Generation (Grounded in Landscape)

Raw input: **AlphaBlitz** — existing Go real-time multiplayer word game (8s letter rotation, 5min rounds, elimination, WebSocket-based).

### Direction A — AlphaBlitz On-Chain Tournament
Free-to-play word game with optional **sponsored prize pools** and **community-funded tournaments**. Off-chain gameplay (Go + WebSockets), on-chain prize escrow via Soroban. Players connect Freighter to claim winnings.
- **Stellar primitives**: Soroban escrow contract, SEP-41 USDC payouts, Freighter auth, contract events for audit.
- **Fit**: Taps underserved gaming vertical; low competition (Stellarcade is the only gaming repo). Avoids gambling by defaulting to free play + sponsor-funded pools.

### Direction B — WordBounties Micro-Tasks
Players earn micro-rewards (x402 micropayments) for completing daily word challenges. Sybil-resistant via wallet age + stake.
- **Stellar primitives**: Soroban reward contract, x402 protocol, SEP-41.
- **Weakness**: Sybil resistance is a missing infra block. No standard identity primitive on Stellar.

### Direction C — Lexicon Arena (Skill Staking)
PvP word battles where players stake XLM/USDC; winner takes all via atomic Soroban settlement.
- **Stellar primitives**: Soroban escrow, SEP-41, state channels (speculative).
- **Weakness**: Explicitly gambling-adjacent. Regulatory wall is structural in most jurisdictions.

### Direction D — ChainLex (Word-Game-as-a-Service)
White-label API for developers to spin up branded word games with built-in Stellar monetization (prize pools, ad slots).
- **Stellar primitives**: Soroban factory contracts, SEP-41.
- **Weakness**: Developer tooling is crowded. Blockchain moat is weak for this use case.

### Direction E — Stellar Spelling Bee (EdTech)
Schools/tutors create spelling bees; parents/sponsors fund prize pools in stablecoins; students play free.
- **Stellar primitives**: Soroban escrow, SEP-41 stablecoins, SEP-10 auth, optional SEP-12 KYC.
- **Weakness**: EdTech adoption is slow; child-safety / COPPA adds friction.

### Direction F — TrusTrove Tournaments
Seasonal word-tournament leagues with on-chain leaderboards and NFT achievement badges (SEP-50).
- **Stellar primitives**: Soroban scoring + escrow, SEP-41, SEP-50 (NFT draft).
- **Weakness**: SEP-50 is still draft. Relying on draft standards for core value is risky.

---

## Phase 3 — Critical Review

| Direction | Weak Spot | Missing Infra | Regulatory | Stellar Fit | MVP Feasibility | Verdict |
|-----------|-----------|---------------|------------|-------------|-----------------|---------|
| **A — Tournament** | Prize settlement must be off-chain; only final room result hits chain. Risk: host can fake results if backend is centralized. | **Blocked**: None. Uses existing Freighter + SEP-41. | **Conditional**: Free-to-play default avoids gambling laws. Entry fees require jurisdiction-by-jurisdiction review. | Strong. Uses proven Soroban escrow + USDC. | High. Existing Go backend handles gameplay. Contract is thin. | **STRONGEST** |
| B — WordBounties | Sybil grinding drains reward pool. No on-chain identity standard exists. | **Blocked**: Sybil-resistant identity. | Low (micro-tasks), but wage laws possible if rewards look like compensation. | Moderate. x402 is novel but unproven at scale. | Low. Identity gap must be closed first. | **WEAK** |
| C — Lexicon Arena | PvP staking = gambling in most countries. Enforcement risk is structural. | None. | **Structural blocker**: Licensed gambling frameworks required. | Strong technically. | Moderate. | **REJECT** |
| D — ChainLex | No defensible moat vs. Unity / HTML5 templates. Blockchain is a cost, not a feature. | None. | Low. | Weak. Gamers don't care about chain. | Moderate. | **REJECT** |
| E — Spelling Bee | EdTech sales cycles are 12–24 months. COPPA / child-data laws add compliance overhead. | None. | **Structural blocker**: Child safety regulations in US/EU. | Moderate. | Low. | **WEAK** |
| F — TrusTrove | SEP-50 is draft. If it changes, badges break. Leaderboard on-chain is wasteful (high write freq). | None. | Low if free-to-play. | Moderate. Draft standard risk. | Moderate. | **CONDITIONAL** |

**Verdict**: Proceed with **Direction A — AlphaBlitz On-Chain Tournament**.

**Key design rule**: Gameplay state stays off-chain (Go + WebSockets). Soroban only handles **prize escrow and final payout**. Backend is trusted for game logic, but financial settlement is trustless.

---

## Phase 4 — Naming, Scoping, and Repo Structure

### Name
**AlphaBlitz** (retain existing brand). Web3 layer is marketed as "AlphaBlitz On-Chain".

### One-Paragraph Description
AlphaBlitz is a real-time multiplayer word game where players race to fill categories starting with a randomly selected letter. The Web3 upgrade adds trustless prize pools and sponsored tournaments on the Stellar network: players compete in free-to-play rooms, while sponsors or community members fund USDC prize pools that are held in Soroban smart-contract escrow and paid out automatically to winners. Off-chain gameplay preserves the 8-second-per-letter pace; on-chain settlement removes counterparty risk for prize distribution and gives contributors verifiable, low-fee payouts in stable currency.

### Repo Structure
Split into two repos to maximize Drips Wave surface area while keeping maintainer overhead sane.

```
alphablitz-contract/          # Pure Rust workspace
├── Cargo.toml                 # Workspace root
├── contracts/
│   └── prize_escrow/          # Single-responsibility Soroban contract
│       ├── Cargo.toml
│       ├── src/
│       │   ├── lib.rs
│       │   ├── types.rs
│       │   └── storage.rs
│       └── tests/
│           └── integration.rs
└── README.md

alphablitz-app/               # Monorepo
├── package.json
├── packages/
│   └── sdk/                   # TypeScript SDK for contract I/O
│       ├── package.json
│       └── src/
│           ├── client.ts
│           ├── types.ts
│           └── xdr.ts
├── apps/
│   └── web/                   # Frontend (React + Vite)
│       ├── package.json
│       └── src/
│           ├── main.tsx
│           ├── App.tsx
│           ├── components/
│           └── lib/
├── indexer/                   # Soroban event indexer + API
│   ├── package.json
│   └── src/
│       ├── index.ts
│       ├── db.ts
│       └── routes.ts
├── docker-compose.yml
├── .env.example
└── README.md
```

**Rationale**: Contract repo is pure Rust → clean for Soroban CI. App repo is TS monorepo → standard for web3 frontends. Each repo can be approved independently in Drips Wave.

---

## Phase 5 — Contract Architecture

### Contract: `prize_escrow`
**Single responsibility**: Hold USDC (via SAC), accept sponsor deposits, release payouts to verified winners, refund sponsors.

### Dependency Graph (Build / Deploy Order)
1. `prize_escrow` (no internal dependencies)

### Storage Schema
| Key | Type | Description |
|-----|------|-------------|
| `P` (pool prefix) | Persistent `Map<PoolId, Pool>` | Active prize pools |
| `C` (claim prefix) | Persistent `Map<(PoolId, PlayerAddress), Claim>` | Payout claims per pool |
| `A` (admin) | Persistent `Address` | Contract admin (backend hot wallet or multisig) |
| `S` (sponsor prefix) | Persistent `Map<Address, bool>` | whitelisted sponsor addresses (optional, for whitelisted mode) |

**Pool struct**:
- `pool_id: u64` — auto-incremented
- `sponsor: Address` — funder
- `asset: Address` — SAC address (USDC)
- `total_amount: i128` — total deposited
- `claimed_amount: i128` — total released
- `status: PoolStatus` — Active / Expired / Finalized
- `game_id: Bytes` — links to off-chain game room
- `expires_at: u64` — unix timestamp
- `created_at: u64` — unix timestamp

**Claim struct**:
- `player: Address`
- `amount: i128`
- `claimed: bool`
- `tx_hash: Bytes` — off-chain transaction proof (hash of game result)

**PoolStatus enum**: `Active = 0`, `Expired = 1`, `Finalized = 2`

### Public Functions
| Function | Params | Returns | Auth | Event |
|----------|--------|---------|------|-------|
| `create_pool` | `sponsor: Address, asset: Address, amount: i128, expires_at: u64, game_id: Bytes` | `pool_id: u64` | `sponsor.require_auth()` | `PoolCreated` |
| `finalize_pool` | `pool_id: u64` | `()` | `admin.require_auth()` | `PoolFinalized` |
| `payout` | `pool_id: u64, player: Address, amount: i128, proof: Bytes` | `()` | `admin.require_auth()` | `Payout` |
| `refund_expired` | `pool_id: u64` | `()` | `sponsor.require_auth()` | `Refunded` |
| `pool_info` | `pool_id: u64` | `Pool` | — | — |
| `player_claim` | `pool_id: u64, player: Address` | `Claim` | — | — |

### Events
- `PoolCreated(pool_id, sponsor, asset, amount, expires_at, game_id)`
- `PoolFinalized(pool_id)`
- `Payout(pool_id, player, amount, proof)`
- `Refunded(pool_id, sponsor, amount)`

### User-Flow Mapping
1. Sponsor funds pool → `create_pool`
2. Game ends off-chain → backend computes winners
3. Backend calls `finalize_pool` (or multi-sig admin)
4. Backend calls `payout` for each winner with signed proof hash
5. Any remaining balance after expiry → `refund_expired`

### Build Sequence
1. Initialize Rust workspace + Cargo.toml
2. Implement storage + types
3. Implement `create_pool` + tests
4. Implement `finalize_pool` + `payout` + tests
5. Implement `refund_expired` + integration tests
6. Add contract meta (SEP-46) + CI workflow

---

## Phase 6 — Contract System Prompt

*(Saved to `docs/system-prompts/alphablitz-contract-system-prompt.md`)*

Role: You are a senior Soroban engineer. No placeholders, no stubs, no unwrap() outside tests.

Repo: `alphablitz-contract`
Tech Stack:
- Rust 1.78+
- soroban-sdk 21.x (match latest stable)
- Rust workspace with single crate `prize_escrow`
- No external crates beyond `soroban-sdk`, `soroban-sdk-test`, `serde` (for types only)

Soroban Patterns:
- Use `soroban_sdk::contracttype!` for all storage types.
- Use `soroban_sdk::symbol!` for event topics.
- Use `require_auth()` only on `sponsor` (create/refund) and `admin` (finalize/payout).
- Use `Address` for all account references.
- Use `i128` for all token amounts (basis points / stroops, no floats).
- Use `u64` for timestamps and auto-increment IDs.
- Emit events via `env.events().publish(...)`.

Naming:
- Contract file: `lib.rs`
- Types: `Pool`, `Claim`, `PoolStatus` in `types.rs`
- Storage keys: `P`, `C`, `A`, `S` as symbols in `storage.rs`
- Public functions: snake_case matching Phase 5 table.

Git Workflow (non-negotiable):
- Never `git add .`
- One logical change per commit
- Commit message format: `feat(contract): <description>` or `fix(contract): <description>`
- Push immediately after each commit

Build Sequence (ordered):
1. `cargo new contracts/prize_escrow`
2. Implement `types.rs` and `storage.rs`
3. Implement `lib.rs` with `create_pool`
4. Write `tests/integration.rs` for `create_pool`
5. Implement `finalize_pool` + `payout`
6. Write tests for `finalize_pool` + `payout`
7. Implement `refund_expired`
8. Write full integration lifecycle test
9. Add SEP-46 contractmeta
10. Add GitHub Actions CI (cargo fmt, cargo clippy, cargo test)

What Not To Do:
- Do not write off-chain logic in the contract.
- Do not use `unwrap()` outside `#[cfg(test)]`.
- Do not use floating-point math.
- Do not add speculative functions not in Phase 5.
- Do not add upgradeability proxy (out of scope for MVP).
- Do not use `std::collections::HashMap` (use `soroban_sdk::Map`).

---

## Phase 7 — App System Prompt

*(Saved to `docs/system-prompts/alphablitz-app-system-prompt.md`)*

Role: You are a senior full-stack web3 engineer building the application layer for AlphaBlitz On-Chain.

Repo: `alphablitz-app`
Tech Stack (exact versions):
- Node.js 20 LTS
- pnpm 9.x (monorepo)
- TypeScript 5.4
- React 18 + Vite 5 (frontend)
- Tailwind CSS 3.4 (frontend styling)
- `@stellar/stellar-sdk` 12.x (TypeScript)
- `@stellar/freighter-api` (wallet connect)
- `zod` 3.x (validation)
- `@tanstack/react-query` 5.x (server state)
- `better-sqlite3` or `pg` (indexer DB, choose pg for prod)
- `express` or `hono` (indexer API)

Monorepo Structure:
```
packages/sdk/       — TypeScript SDK wrapping contract I/O
apps/web/           — React frontend
indexer/            — Event indexer + REST API
```

Contract Interfaces (restated standalone):
- `create_pool(sponsor, asset, amount, expiresAt, gameId) => poolId`
- `finalize_pool(poolId)`
- `payout(poolId, player, amount, proof)`
- `refundExpired(poolId)`
- `poolInfo(poolId) => Pool`
- `playerClaim(poolId, player) => Claim`

SDK Responsibilities:
- XDR encoding helpers for contract calls.
- `readContract` / `sendTransaction` wrappers using `@stellar/stellar-sdk`.
- Wallet connection abstraction (Freighter primary, Albedo fallback).

Indexer Responsibilities:
- Poll Soroban RPC for `PoolCreated`, `Payout`, `Refunded` events.
- Persist to Postgres.
- Expose `/api/pools`, `/api/pools/:id/claims`, `/api/games/:id/payouts`.

Frontend Responsibilities:
- Wallet connect button + address display.
- Game lobby (WebSocket connection to existing Go backend).
- Prize pool creation flow (sponsor).
- Payout history view.
- Theme toggle (retain existing light/dark).

Environment Variables:
| Variable | Description |
|----------|-------------|
| `VITE_NETWORK` | `testnet` or `mainnet` |
| `VITE_RPC_URL` | Soroban RPC endpoint |
| `VITE_CONTRACT_PRIZE_ESCROW` | Deployed contract ID |
| `VITE_FREIGHTER_API_URL` | `https://freighter.app` |
| `VITE_API_URL` | Indexer backend URL |
| `VITE_WS_URL` | Go game server WebSocket URL |

Git Workflow:
- Never `git add .`
- One logical change per commit
- Conventional commits: `feat(app):`, `fix(app):`, `chore(app):`
- Push immediately

Build Sequence:
1. `pnpm init` + workspace config
2. Build `packages/sdk` (types + contract client)
3. Build `indexer` (DB schema + event listeners + routes)
4. Build `apps/web` shell (Vite + React + routing + wallet context)
5. Wire game client into web app
6. Add env handling + network switching
7. Add CI (lint, typecheck, build, test)

Per-Sub-Stack Standards:
- SDK: no `any`, all contract responses typed, Zod schemas for all API inputs.
- Indexer: use parameterized SQL, connection pooling, graceful shutdown.
- Frontend: component files PascalCase, hooks camelCase, Tailwind utility-first, no inline styles.

Constraints Checklist:
- Do not hardcode contract IDs in frontend (use env).
- Do not store private keys in frontend.
- Do not trust off-chain game results without a backend admin signature check on-chain.
- Do not mint/burn tokens directly; use SAC transfers only.

---

## Phase 8 — Local Environment and Deployment (Planned)

### Diagnosis Approach
- **PATH conflicts**: Ensure `~/.cargo/bin` and `stellar` CLI are on PATH.
- **Version mismatches**: Pin Rust to `1.78.0` (matches Soroban CI images), Soroban CLI to latest `stable`.
- **Disk space**: Soroban target adds ~2GB; verify before build.

### Fix Preference
Surgical fixes:
- If `stellar contract build` fails, check `rustup target add wasm32-unknown-unknown` before reinstalling toolchain.
- If `stellar contract deploy` fails with auth error, verify network passphrase and funded deployer key.

### One SDK/Compiler Version
- Rust `1.78.0`
- `wasm32-unknown-unknown` target
- `soroban-cli` installed via `cargo install soroban-cli` (matches Rust toolchain)

### Deployment Scripts
- `scripts/deploy-contract.sh`: builds WASM, deploys to testnet, initializes admin, prints contract ID.
- `scripts/setup-env.sh`: writes `.env.local` and `indexer/.env` from deployed addresses.

### Wiring Env Vars
- Frontend `.env.local`: `VITE_CONTRACT_PRIZE_ESCROW`, `VITE_RPC_URL`, etc.
- Vercel project settings: mirror frontend env vars.
- Render service settings: `DATABASE_URL`, `VITE_API_URL`, `VITE_CONTRACT_PRIZE_ESCROW`.

---

## Phase 9 — Hosting and Service Topology (Planned)

```
User → [Vercel: Frontend] → [Render: Indexer/API] → [PostgreSQL (Render same region)]
         ↓ direct RPC
      [Soroban RPC] → [Stellar Testnet / Mainnet]
```

- **Frontend**: Vercel (purpose-built for static + SSR; React/Vite adapter).
- **Indexer/API**: Render (long-running stateful service). Build: `pnpm install && pnpm build`. Start: `node indexer/dist/index.js`.
- **Database**: PostgreSQL on Render, same region as indexer. Internal connection string (not exposed).
- **Game Server**: Existing Go server stays separate (can deploy to Fly.io, Railway, or Render). Frontend connects via `wss://`.
- **Failure trace**: If frontend hits `localhost` in prod, fix at source by reading `.env.example` and enforcing `VITE_` prefix in Vite config; never allow default localhost fallback in production build.

---

## Phase 10 — Repo Hygiene for Program Approval

### Checklist
- [ ] Branch protection on `main` (require CI pass, require review).
- [ ] `CONTRIBUTING.md` with Drips Wave link, setup steps, commit rules.
- [ ] `SECURITY.md` with disclosure contact, unaudited disclaimer.
- [ ] `README.md` with banner, badges (Stellar Wave, Rust, TypeScript), maintainer contact, architecture summary, quick-start.
- [ ] GitHub topics: `stellar`, `soroban`, `stellar-wave`, `rust`, `typescript`, `web3`, `gaming`, `word-game`.
- [ ] Bulk issue generation via `gh` CLI script (title + labels + acceptance criteria).
- [ ] Release tag `v0.1.0` with deployed testnet contract IDs in body.

### README Pattern (matches approved repos)
- Project logo / banner
- Badges (Stellar Wave, build status, Rust version)
- Maintainer contact table
- Community link (Discord)
- Architecture diagram
- Quick-start commands
- Contributing section
- `contrib.rocks` image (optional)

---

## Phase 11 — Documentation Site (Planned)

Separate from README:
- Introduction with cited figures (Stellar TPS, fee cost, Wave program stats).
- Protocol mechanics / state machine with worked numbers (e.g., "10 players × 0.01 USDC entry = 0.1 USDC pool, 5% platform fee = 0.095 USDC to winners").
- Smart contract reference (auto-generated from doc comments).
- Per-persona guides: Player (how to connect wallet, join tournament), Sponsor (how to fund a pool), Contributor (how to run indexer).
- Developer guide: setup, env vars, SDK/API reference with real examples.

---

## Phase 12 — Submission (Planned)

### Pre-submission verification
- Confirm `AlphaBlitz` / repos are not already in approved list (search drips.network/wave/stellar/repos for "AlphaBlitz", "alphablitz", "word game").
- Search result: **No existing AlphaBlitz repo found** in approved list (verified 2026-07-28).

### Supporting Links (to assemble)
- Live app URL: `https://alphablitz.vercel.app` (post-deploy)
- Contract repo: `https://github.com/jerryjuche/alphablitz-contract`
- App repo: `https://github.com/jerryjuche/alphablitz-app`
- Docs site: `https://alphablitz-docs.vercel.app`
- On-chain verification: Stellar Expert links to deployed contract + sponsor pool transactions
- Demo video: 2–3 min Loom/YouTube showing wallet connect → free play → sponsored tournament → payout

### Submission Text (draft)
> AlphaBlitz is a real-time multiplayer word game that adds trustless prize pools and sponsored tournaments to Stellar. Players compete in fast-paced rooms (8-second letter rotation, elimination rounds) while sponsors fund USDC prize pools held in Soroban escrow. Gameplay stays off-chain for speed; financial settlement moves on-chain for trustlessness. The project targets the consumer gaming vertical — currently underserved on Stellar — and aligns with SDF's 2026 user-adoption goals by bringing casual players into the ecosystem through a free-to-play experience with optional web3 rewards.

### Planned Issues (to generate)
- `feat(contract): implement create_pool with sponsor auth`
- `feat(contract): implement finalize_pool and payout with proof hash`
- `feat(contract): add refund_expired for stale pools`
- `feat(sdk): build Soroban RPC client for prize_escrow`
- `feat(indexer): index PoolCreated and Payout events`
- `feat(web): wire Freighter wallet connect into game lobby`
- `feat(web): build sponsor prize-pool creation flow`
- `feat(web): add payout history dashboard`
- `chore: add GitHub Actions CI for Rust and TypeScript`
- `docs: write contract reference and contributor guide`

---

## Phase 13 — Post-Approval Iteration (Planned)

- Scope new gaps honestly (e.g., whitelist mode, multi-asset pools, upgradeable proxy).
- For cross-repo changes, write coordinated issues with `Depends on` references.
- Maintain same issue rigor (title + labels + acceptance criteria).
