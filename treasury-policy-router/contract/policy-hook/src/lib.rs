#![no_std]
use core::clone::Clone;
use core::default::Default;
use soroban_sdk::{contract, contractimpl, contracttype, Address, Env, Map, Symbol};

mod types;

use crate::types::{DataKey, PolicyConfig, PolicyResult};

#[contract]
pub struct PolicyHook;

#[contractimpl]
impl PolicyHook {
    pub fn initialize(e: &Env, admin: Address) {
        let config = PolicyConfig {
            jurisdiction_allowlist: Map::new(&e),
            daily_cap: 0,
            travel_rule_threshold: 0,
        };
        e.storage().persistent().set(&DataKey::Config, &config);
        e.storage().persistent().set(&DataKey::Admin, &admin);
        e.events().publish(
            (Symbol::new(&e, "PolicyInitialized"),),
            admin,
        );
    }

    pub fn set_admin(e: &Env, new_admin: Address, operator: Address) {
        operator.require_auth();
        let admin: Address = e.storage().persistent().get(&DataKey::Admin).unwrap();
        if admin != operator {
            panic!("Unauthorized");
        }
        e.storage().persistent().set(&DataKey::Admin, &new_admin);
        e.events().publish(
            (Symbol::new(&e, "PolicyUpdated"), Symbol::new(&e, "admin")),
            new_admin,
        );
    }

    pub fn set_jurisdiction(e: &Env, country: Symbol, allowed: bool, operator: Address) {
        operator.require_auth();
        let admin: Address = e.storage().persistent().get(&DataKey::Admin).unwrap();
        if admin != operator {
            panic!("Unauthorized");
        }
        let mut config: PolicyConfig = e.storage().persistent().get(&DataKey::Config).unwrap();
        config.jurisdiction_allowlist.set(country.clone(), allowed);
        e.storage().persistent().set(&DataKey::Config, &config);
        e.events().publish(
            (Symbol::new(&e, "JurisdictionSet"), country.clone()),
            allowed,
        );
    }

    pub fn set_daily_cap(e: &Env, cap: i128, operator: Address) {
        operator.require_auth();
        let admin: Address = e.storage().persistent().get(&DataKey::Admin).unwrap();
        if admin != operator {
            panic!("Unauthorized");
        }
        let mut config: PolicyConfig = e.storage().persistent().get(&DataKey::Config).unwrap();
        config.daily_cap = cap;
        e.storage().persistent().set(&DataKey::Config, &config);
        e.events().publish(
            (Symbol::new(&e, "PolicyUpdated"), Symbol::new(&e, "daily_cap")),
            cap,
        );
    }

    pub fn set_travel_rule_threshold(e: &Env, threshold: i128, operator: Address) {
        operator.require_auth();
        let admin: Address = e.storage().persistent().get(&DataKey::Admin).unwrap();
        if admin != operator {
            panic!("Unauthorized");
        }
        let mut config: PolicyConfig = e.storage().persistent().get(&DataKey::Config).unwrap();
        config.travel_rule_threshold = threshold;
        e.storage().persistent().set(&DataKey::Config, &config);
        e.events().publish(
            (Symbol::new(&e, "PolicyUpdated"), Symbol::new(&e, "travel_rule_threshold")),
            threshold,
        );
    }

    pub fn check_policy(e: &Env, from: Address, _to: Address, amount: i128) -> PolicyResult {
        let config: PolicyConfig = e.storage().persistent().get(&DataKey::Config).unwrap();
        let day = e.ledger().timestamp() / 86400;
        let spend_key = DataKey::DailySpend(from.clone(), day);
        let mut spent: i128 = e.storage().persistent().get(&spend_key).unwrap_or(0);
        spent += amount;
        if spent > config.daily_cap && config.daily_cap > 0 {
            return PolicyResult {
                approved: false,
                flag: true,
                reason: Symbol::new(&e, "daily_cap_exceeded"),
            };
        }
        e.storage().persistent().set(&spend_key, &spent);

        let country = Symbol::new(&e, "US");
        if let Some(allowed) = config.jurisdiction_allowlist.get(country.clone()) {
            if !allowed {
                return PolicyResult {
                    approved: false,
                    flag: false,
                    reason: Symbol::new(&e, "jurisdiction_blocked"),
                };
            }
        }

        if amount >= config.travel_rule_threshold && config.travel_rule_threshold > 0 {
            return PolicyResult {
                approved: true,
                flag: true,
                reason: Symbol::new(&e, "travel_rule_flag"),
            };
        }

        PolicyResult {
            approved: true,
            flag: false,
            reason: Symbol::new(&e, "approved"),
        }
    }

    pub fn get_jurisdiction(e: &Env, country: Symbol) -> bool {
        let config: PolicyConfig = e.storage().persistent().get(&DataKey::Config).unwrap();
        config.jurisdiction_allowlist.get(country).unwrap_or(false)
    }

    pub fn get_daily_cap(e: &Env) -> i128 {
        let config: PolicyConfig = e.storage().persistent().get(&DataKey::Config).unwrap();
        config.daily_cap
    }
}
