mod types;

use soroban_sdk::{contracttype, Address, Env, Symbol};
use crate::types::{Claim, Pool};

pub const PREFIX_POOL: Symbol = symbol!("P");
pub const PREFIX_CLAIM: Symbol = symbol!("C");
pub const KEY_ADMIN: Symbol = symbol!("A");
pub const KEY_COUNTER: Symbol = symbol!("N");

pub fn get_admin(e: &Env) -> Address {
    e.storage().persistent().get(&KEY_ADMIN).unwrap_or_else(|| panic!("admin not set"))
}

pub fn set_admin(e: &Env, admin: &Address) {
    e.storage().persistent().set(&KEY_ADMIN, admin);
}

pub fn get_counter(e: &Env) -> u64 {
    e.storage().persistent().get(&KEY_COUNTER).unwrap_or(0)
}

pub fn set_counter(e: &Env, counter: &u64) {
    e.storage().persistent().set(&KEY_COUNTER, counter);
}

pub fn pool_key(e: &Env, pool_id: u64) -> (Symbol, u64) {
    (PREFIX_POOL, pool_id)
}

pub fn get_pool(e: &Env, pool_id: u64) -> Option<Pool> {
    e.storage().persistent().get(&pool_key(e, pool_id))
}

pub fn set_pool(e: &Env, pool_id: u64, pool: &Pool) {
    e.storage().persistent().set(&pool_key(e, pool_id), pool);
}

pub fn claim_key(e: &Env, pool_id: u64, player: &Address) -> (Symbol, u64, Address) {
    (PREFIX_CLAIM, pool_id, *player)
}

pub fn get_claim(e: &Env, pool_id: u64, player: &Address) -> Option<Claim> {
    e.storage().persistent().get(&claim_key(e, pool_id, player))
}

pub fn set_claim(e: &Env, pool_id: u64, player: &Address, claim: &Claim) {
    e.storage().persistent().set(&claim_key(e, pool_id, player), claim);
}
