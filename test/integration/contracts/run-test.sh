#!/usr/bin/env bash
# SPDX-License-Identifier: GPL-3.0

set -e

trap 'kill $(jobs -p) &>/dev/null' EXIT

cd "$(dirname "$0")"

which anvil >/dev/null
which go >/dev/null
which awk >/dev/null

set -v

export ANVIL_OUTPUT=$(mktemp --tmpdir)
anvil | tee $ANVIL_OUTPUT &

sleep 2

DEPLOYER_KEY=$(awk '/Private Keys/ {flag=1; next} flag && /^\(0\)/ {print $2; exit}' $ANVIL_OUTPUT); echo $DEPLOYER_KEY
PUBLISHER_KEY=$(awk '/Private Keys/ {flag=1; next} flag && /^\(0\)/ {print $2; exit}' $ANVIL_OUTPUT); echo $PUBLISHER_KEY
PROVIDER_KEY=$(awk '/Private Keys/ {flag=1; next} flag && /^\(1\)/ {print $2; exit}' $ANVIL_OUTPUT); echo $PROVIDER_KEY

set -x

( cd ../../../contracts; forge script script/Deploy.s.sol --private-key "$DEPLOYER_KEY" --rpc-url http://127.0.0.1:8545 --broadcast)

PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')

go run . "$PUBLISHER_KEY" "$PROVIDER_KEY" "$PAYMENT_CONTRACT"

