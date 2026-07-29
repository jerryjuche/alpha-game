#!/bin/bash
set -eux
stellar contract build --package policy-hook
stellar contract deploy --wasm target/wasm32-unknown-unknown/release/policy_hook.wasm --source admin --network futurenet
echo "Deployment complete. Save the contract ID shown above."
