# Policy Hook Reference

## initialize

```rust
pub fn initialize(e: &Env, admin: Address)
```

Sets the admin address. Called once at deployment.

## set_jurisdiction

```rust
pub fn set_jurisdiction(e: &Env, country: Symbol, allowed: bool, operator: Address)
```

Sets jurisdiction allowlist entry. Requires admin auth.

## set_daily_cap

```rust
pub fn set_daily_cap(e: &Env, cap: i128, operator: Address)
```

Sets daily spending cap per address. Requires admin auth.

## set_travel_rule_threshold

```rust
pub fn set_travel_rule_threshold(e: &Env, threshold: i128, operator: Address)
```

Sets threshold that triggers Travel Rule flag. Requires admin auth.

## check_policy

```rust
pub fn check_policy(e: &Env, from: Address, to: Address, amount: i128) -> PolicyResult
```

Evaluates policy and returns Approve / Reject / Flag.

## get_jurisdiction

```rust
pub fn get_jurisdiction(e: &Env, country: Symbol) -> bool
```

Returns allowlist status for a country.

## get_daily_cap

```rust
pub fn get_daily_cap(e: &Env) -> i128
```

Returns current daily cap.
