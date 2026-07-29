use soroban_sdk::{
    contract, contractimpl, contractmeta, token::Client as TokenClient, Address, BytesN, Env, Symbol,
    Vec,
};
use crate::{storage::*, types::*};

mod storage;
mod types;

contractmeta!(key = "sep", val = "41");

#[contract]
pub struct PrizeEscrow;

#[contractimpl]
impl PrizeEscrow {
    pub fn initialize(env: Env, admin: Address) {
        if admin_exists(&env) {
            panic!("already initialized");
        }
        set_admin(&env, &admin);
        set_counter(&env, &0);
    }

    pub fn create_pool(
        env: Env,
        sponsor: Address,
        asset: Address,
        amount: i128,
        expires_at: u64,
        game_id: BytesN<32>,
    ) -> u64 {
        sponsor.require_auth();
        if amount <= 0 {
            panic!("invalid amount");
        }
        let pool_id = get_counter(&env) + 1;
        let now = env.ledger().timestamp();
        let token = TokenClient::new(&env, &asset);
        let contract_addr = env.current_contract_address();
        token.transfer(&sponsor, &contract_addr, &amount);

        let pool = Pool {
            pool_id,
            sponsor: sponsor.clone(),
            asset: asset.clone(),
            total_amount: amount,
            claimed_amount: 0,
            status: PoolStatus::Active,
            game_id,
            expires_at,
            created_at: now,
        };
        set_pool(&env, pool_id, &pool);
        set_counter(&env, &pool_id);

        env.events().publish(
            (symbol!("pool_created"),),
            (pool_id, sponsor, asset, amount, expires_at, game_id),
        );
        pool_id
    }

    pub fn finalize_pool(env: Env, pool_id: u64) {
        let admin = get_admin(&env);
        admin.require_auth();
        let mut pool = get_pool(&env, pool_id).expect("pool not found");
        if pool.status == PoolStatus::Finalized {
            panic!("already finalized");
        }
        if pool.status == PoolStatus::Expired {
            panic!("pool expired");
        }
        pool.status = PoolStatus::Finalized;
        set_pool(&env, pool_id, &pool);
        env.events().publish((symbol!("pool_finalized"),), pool_id);
    }

    pub fn payout(env: Env, pool_id: u64, player: Address, amount: i128, proof: BytesN<32>) {
        let admin = get_admin(&env);
        admin.require_auth();
        if amount <= 0 {
            panic!("invalid amount");
        }
        let mut pool = get_pool(&env, pool_id).expect("pool not found");
        if pool.status != PoolStatus::Finalized {
            panic!("pool not finalized");
        }
        let remaining = pool.total_amount - pool.claimed_amount;
        if amount > remaining {
            panic!("insufficient balance");
        }
        let token = TokenClient::new(&env, &pool.asset);
        let contract_addr = env.current_contract_address();
        token.transfer(&contract_addr, &player, &amount);

        pool.claimed_amount += amount;
        set_pool(&env, pool_id, &pool);

        let claim = Claim {
            player: player.clone(),
            amount,
            claimed: true,
            proof,
        };
        set_claim(&env, pool_id, &player, &claim);

        env.events()
            .publish((symbol!("payout"),), (pool_id, player, amount, proof));
    }

    pub fn refund_expired(env: Env, pool_id: u64) {
        let mut pool = get_pool(&env, pool_id).expect("pool not found");
        pool.sponsor.require_auth();
        let now = env.ledger().timestamp();
        if now < pool.expires_at {
            panic!("pool not expired");
        }
        if pool.status == PoolStatus::Expired {
            panic!("already expired");
        }
        let remaining = pool.total_amount - pool.claimed_amount;
        if remaining > 0 {
            let token = TokenClient::new(&env, &pool.asset);
            let sponsor_addr = pool.sponsor.clone();
            token.transfer(&env.current_contract_address(), &sponsor_addr, &remaining);
        }
        pool.status = PoolStatus::Expired;
        set_pool(&env, pool_id, &pool);
        env.events()
            .publish((symbol!("refunded"),), (pool_id, pool.sponsor, remaining));
    }

    pub fn pool_info(env: Env, pool_id: u64) -> Pool {
        get_pool(&env, pool_id).expect("pool not found")
    }

    pub fn player_claim(env: Env, pool_id: u64, player: Address) -> Claim {
        get_claim(&env, pool_id, &player).unwrap_or(Claim {
            player: player.clone(),
            amount: 0,
            claimed: false,
            proof: BytesN::new(&env, &[0u8; 32]),
        })
    }
}

fn admin_exists(e: &Env) -> bool {
    e.storage().persistent().has(&KEY_ADMIN)
}
