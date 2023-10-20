#!/usr/bin/env bash

set -e

trap 'kill $(jobs -p) &>/dev/null' EXIT

cd "$(dirname "$0")"

which anvil >/dev/null
which go >/dev/null
which awk >/dev/null

set -v

export ANVIL_OUTPUT=$(mktemp ipfs.XXXX --tmpdir)
anvil | tee $ANVIL_OUTPUT &

sleep 2

PUBLISHER_KEY=$(awk '/Private Keys/ {flag=1; next} flag && /^\(0\)/ {print $2; exit}' $ANVIL_OUTPUT); echo $PUBLISHER_KEY
PROVIDER_KEY=$(awk '/Private Keys/ {flag=1; next} flag && /^\(1\)/ {print $2; exit}' $ANVIL_OUTPUT); echo $PROVIDER_KEY

set -x

go run . "$PUBLISHER_KEY" "$PROVIDER_KEY"

