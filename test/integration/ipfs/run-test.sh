#!/usr/bin/env bash

set -e

trap 'kill $(jobs -p) &>/dev/null' EXIT

cd "$(dirname "$0")"

which go >/dev/null
which ipfs >/dev/null

export IPFS_BASE=$(mktemp ipfs.XXXX --tmpdir -d)
function as {
  export IPFS_PATH="$IPFS_BASE/$1"
}

as publisher
ipfs init --profile randomports > /dev/null
ipfs config --json Experimental.Libp2pStreamMounting true > /dev/null
ipfs config Addresses.API /ip4/127.0.0.1/tcp/0 > /dev/null
ipfs config Addresses.Gateway /ip4/127.0.0.1/tcp/0 > /dev/null
ipfs daemon > /dev/null &

as provider
ipfs init --profile randomports > /dev/null
ipfs config --json Experimental.Libp2pStreamMounting true > /dev/null
ipfs config Addresses.API /ip4/127.0.0.1/tcp/0 > /dev/null
ipfs config Addresses.Gateway /ip4/127.0.0.1/tcp/0 > /dev/null
PROVIDER_ID=$(ipfs id -f '<id>')
ipfs daemon > /dev/null &
{ while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.5; done; } 2>/dev/null
go run . provider &

as publisher
{ while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.5; done; } 2>/dev/null
sleep 3
go run . publisher "$PROVIDER_ID"

kill $(jobs -p) &>/dev/null
