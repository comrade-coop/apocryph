#!/usr/bin/env bash
# SPDX-License-Identifier: GPL-3.0

set -e

trap 'kill $(jobs -p) &>/dev/null' EXIT

cd "$(dirname "$0")"

which go >/dev/null
which ipfs >/dev/null
which docker >/dev/null
which jq >/dev/null

set -v

ipfs daemon >/dev/null || true &
{ while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.5; done; } 2>/dev/null

docker image pull hello-world

## 1. Test docker daemon <-> our IPDR

docker image rm hello-world-copy &>/dev/null || true

RUN_OUTPUT=$(mktemp)
go run . docker-daemon:hello-world:latest ipdr: | tee $RUN_OUTPUT
IPDR_REFERENCE=$(tail -n 1 $RUN_OUTPUT)
go run . "$IPDR_REFERENCE" docker-daemon:hello-world-copy:latest

diff <(docker image inspect hello-world | jq '.[].RootFS') <(docker image inspect hello-world-copy | jq '.[].RootFS') -q

## 2. Test docker daemon <-> our IPDR (encrypted)

docker image rm hello-world-copy &>/dev/null || true

RUN_OUTPUT=$(mktemp)
KEY_FILE=$(mktemp --suffix=".json"); echo $KEY_FILE
go run . docker-daemon:hello-world:latest ipdr: "$KEY_FILE" | tee $RUN_OUTPUT

IPDR_REFERENCE=$(tail -n 1 $RUN_OUTPUT)

go run . "$IPDR_REFERENCE" docker-daemon:hello-world-copy:latest "$KEY_FILE"

diff <(docker image inspect hello-world | jq '.[].RootFS') <(docker image inspect hello-world-copy | jq '.[].RootFS') -q

## 3. Test docker daemon -> orig IPDR -> our IPDR -> docker daemon

which ipdr >/dev/null # Rest of the test checks for compatibility with ipdr

docker image rm hello-world-copy &>/dev/null || true

IPDR_OUTPUT=$(mktemp)
ipdr push hello-world | tee $IPDR_OUTPUT
IPDR_REFERENCE_2=$(tail -n 1 $IPDR_OUTPUT)
go run . "ipdr:$IPDR_REFERENCE_2" docker-daemon:hello-world-copy:latest

diff <(docker image inspect hello-world | jq '.[].RootFS') <(docker image inspect hello-world-copy | jq '.[].RootFS') -q

## 4.  Test docker daemon -> our IPDR -> orig IPDR -> docker daemon
# ipdr pull "$IPDR_REFERENCE" - ipdr pull is rather buggy and not quite testable, sadly. Maybe when they fix it one day

## 5.

echo success!

kill $(jobs -p) &>/dev/null
