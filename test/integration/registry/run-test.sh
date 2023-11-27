#!/bin/sh
cd "$(dirname "$0")"

anvil > anvil_output.txt &
ipfs daemon >/dev/null 2>&1 &
sleep 2

PROVIDER_KEY=$(awk '/Private Keys/ {flag=1; next} flag && /^\(0\)/ {print $2; exit}' anvil_output.txt)
COMPETITOR_KEY=$(awk '/Private Keys/ {flag=1; next} flag && /^\(1\)/ {print $2; exit}' anvil_output.txt)
PUBLISHER_KEY=$(awk '/Private Keys/ {flag=1; next} flag && /^\(1\)/ {print $3; exit}' anvil_output.txt)

forge create --root ../../../contracts Registry --private-key $PROVIDER_KEY >/dev/null
forge create --root ../../../contracts MockToken --private-key $PROVIDER_KEY >/dev/null

REGISTRY_ADDRESS=0x5FbDB2315678afecb367f032d93F642f64180aa3
TOKEN_ADDRESS=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512

REGION="bul-east-1"

echo "Registering 2 providers, 2 pricing tables:\n"

go run ../../../cmd/tpodserver registry register --config config.yaml --ethereum-key "$PROVIDER_KEY" --registry-contract "$REGISTRY_ADDRESS" --token-contract "$TOKEN_ADDRESS" >/dev/null
go run ../../../cmd/tpodserver registry register --config competitor.yaml --ethereum-key "$COMPETITOR_KEY" --registry-contract "$REGISTRY_ADDRESS" --token-contract "$TOKEN_ADDRESS" >/dev/null
go run ../../../cmd/trustedpods registry get --config config.yaml --ethereum-key "$PUBLISHER_KEY" --registry-contract "$REGISTRY_ADDRESS" --token-contract "$TOKEN_ADDRESS"

echo "\nUnsubscribe provider1 from the first table, subscribe to the second, and return providers in bul-east-1 region (skip tables with no subscribers)\n"

go run ../../../cmd/tpodserver registry unsubscribe 1 --ethereum-key "$PROVIDER_KEY" --registry-contract "$REGISTRY_ADDRESS" --config config.yaml
go run ../../../cmd/tpodserver registry subscribe 2 --ethereum-key "$PROVIDER_KEY" --registry-contract "$REGISTRY_ADDRESS" --config config.yaml
go run ../../../cmd/trustedpods registry get --config config.yaml --ethereum-key "$PUBLISHER_KEY" --registry-contract "$REGISTRY_ADDRESS" --token-contract "$TOKEN_ADDRESS" --region "$REGION"

pkill anvil
ipfs shutdown
rm anvil_output.txt
