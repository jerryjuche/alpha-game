#[cfg(test)]
mod tests {
    use soroban_sdk::{Address, BytesN, Env, testutils::*};
    use soroban_sdk_test::TestEnv;
    use crate::{PrizeEscrow, PrizeEscrowClient};

    #[test]
    fn test_pool_lifecycle() {
        let env = TestEnv::new();

        let admin = Address::generate(&env);
        let sponsor = Address::generate(&env);
        let player = Address::generate(&env);
        let asset = Address::generate(&env);

        let contract_id = env.register_contract(None, PrizeEscrow);
        let client = PrizeEscrowClient::new(&env, &contract_id);

        client.initialize(&admin);

        let game_id = BytesN::new(&env, &[1u8; 32]);
        let pool_id = client.create_pool(
            &sponsor,
            &asset,
            &1000,
            &(env.ledger().timestamp() + 86400),
            &game_id,
        );

        let pool = client.pool_info(&pool_id);
        assert_eq!(pool.total_amount, 1000);
        assert_eq!(pool.status, 0);

        client.finalize_pool(&pool_id);

        let proof = BytesN::new(&env, &[2u8; 32]);
        client.payout(&pool_id, &player, &500, &proof);

        let claim = client.player_claim(&pool_id, &player);
        assert_eq!(claim.amount, 500);
        assert!(claim.claimed);
    }
}
