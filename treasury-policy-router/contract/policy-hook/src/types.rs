use soroban_sdk::{contracttype, Address, Map, Symbol};

#[contracttype]
pub struct PolicyResult {
    pub approved: bool,
    pub flag: bool,
    pub reason: Symbol,
}

#[contracttype]
pub struct PolicyConfig {
    pub jurisdiction_allowlist: Map<Symbol, bool>,
    pub daily_cap: i128,
    pub travel_rule_threshold: i128,
}

#[contracttype]
pub enum DataKey {
    Config,
    Admin,
    DailySpend(Address, u64),
}
