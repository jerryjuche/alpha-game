use soroban_sdk::{contracttype, symbol, Address, BytesN, Env, Vec};

#[contracttype]
#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Pool {
    pub pool_id: u64,
    pub sponsor: Address,
    pub asset: Address,
    pub total_amount: i128,
    pub claimed_amount: i128,
    pub status: PoolStatus,
    pub game_id: BytesN<32>,
    pub expires_at: u64,
    pub created_at: u64,
}

#[contracttype]
#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Claim {
    pub player: Address,
    pub amount: i128,
    pub claimed: bool,
    pub proof: BytesN<32>,
}

#[contracttype]
#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub enum PoolStatus {
    Active = 0,
    Expired = 1,
    Finalized = 2,
}

#[contracttype]
#[derive(Clone, Debug, Eq, PartialEq)]
pub struct Error;

impl Error {
    pub const POOL_NOT_FOUND: u32 = 1;
    pub const POOL_NOT_FINALIZED: u32 = 2;
    pub const INSUFFICIENT_BALANCE: u32 = 3;
    pub const POOL_NOT_EXPIRED: u32 = 4;
    pub const ALREADY_FINALIZED: u32 = 5;
    pub const ALREADY_EXPIRED: u32 = 6;
    pub const INVALID_AMOUNT: u32 = 7;
    pub const UNAUTHORIZED: u32 = 8;
}
