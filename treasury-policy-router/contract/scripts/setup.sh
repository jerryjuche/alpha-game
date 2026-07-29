#!/bin/bash
set -e
stellar network add --global futurenet --rpc-url https://rpc-futurenet.stellar.org --network-passphrase "Test SDF Future Network ; October 2022"
stellar keys generate --global admin --fund --network futurenet
cp .env.example .env.local
