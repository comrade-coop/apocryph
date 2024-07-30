#!/bin/bash

docker rm --force anvil || true

# (NOTE: Unfortunately, we cannot use a port other than 8545, or otherwise the eth-rpc service will break)
docker run -d -p 8545:8545 --restart=always --name=anvil \
  ghcr.io/foundry-rs/foundry:nightly-619f3c56302b5a665164002cb98263cd9812e4d5 \
  -- 'anvil --host 0.0.0.0 --state /anvil-state.json' 2>/dev/null || {
    docker exec anvil ash -c 'kill 1 && rm -f /anvil-state.json' # Reset anvil state
}
sleep 5

# deploy the contracts
DEPLOYER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(0\)/ {print $2; exit}') # anvil.accounts[0]
( cd ../../../contracts; forge script script/Deploy.s.sol --private-key "$DEPLOYER_KEY" --rpc-url http://localhost:8545 --broadcast)
