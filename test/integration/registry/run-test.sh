#!/bin/bash
# SPDX-License-Identifier: GPL-3.0
set -e

cd "$(dirname "$0")"

trap 'kill $(jobs -p) &>/dev/null' EXIT

anvil >/dev/null &
ipfs daemon >/dev/null &
sleep 2

PROVIDER_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
COMPETITOR_KEY=0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
PUBLISHER_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a

forge create --root ../../../contracts Registry --private-key $PROVIDER_KEY
forge create --root ../../../contracts MockToken --private-key $PROVIDER_KEY

REGISTRY_ADDRESS=0x5FbDB2315678afecb367f032d93F642f64180aa3
TOKEN_ADDRESS=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512

REGION="bul-east-1"

echo -e "Registering 2 providers, 2 pricing tables:\n"

go run ../../../cmd/tpodserver registry register \
  --config config.yaml \
  --ethereum-key "$PROVIDER_KEY" \
  --registry-contract "$REGISTRY_ADDRESS" \
  --token-contract "$TOKEN_ADDRESS" \

go run ../../../cmd/tpodserver registry register \
  --config competitor.yaml \
  --ethereum-key "$COMPETITOR_KEY" \
  --registry-contract "$REGISTRY_ADDRESS" \
  --token-contract "$TOKEN_ADDRESS" \

echo -e "\nPrint tables:\n"
go run ../../../cmd/trustedpods registry get --config config.yaml --ethereum-key "$PUBLISHER_KEY" --registry-contract "$REGISTRY_ADDRESS" --token-contract "$TOKEN_ADDRESS"

echo -e "\nUnsubscribe provider1 from the first table, subscribe to the second, and return providers in bul-east-1 region (skip tables with no subscribers)\n"

go run ../../../cmd/tpodserver registry unsubscribe 1 --ethereum-key "$PROVIDER_KEY" --registry-contract "$REGISTRY_ADDRESS" --config config.yaml
go run ../../../cmd/tpodserver registry subscribe 2 --ethereum-key "$PROVIDER_KEY" --registry-contract "$REGISTRY_ADDRESS" --config config.yaml
go run ../../../cmd/trustedpods registry get --config config.yaml --ethereum-key "$PUBLISHER_KEY" --registry-contract "$REGISTRY_ADDRESS" --token-contract "$TOKEN_ADDRESS" --region "$REGION"

