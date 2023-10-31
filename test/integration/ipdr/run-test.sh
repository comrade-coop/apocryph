#!/usr/bin/env bash

set -e

trap 'kill $(jobs -p) &>/dev/null' EXIT

cd "$(dirname "$0")"

which go >/dev/null
which ipfs >/dev/null
which docker >/dev/null

set -v

ipfs daemon >/dev/null || true &
{ while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.5; done; } 2>/dev/null

docker image pull hello-world
docker image rm hello-world-copy &>/dev/null || true

RUN_OUTPUT=$(mktemp)
go run . docker-daemon:hello-world:latest ipdr: | tee $RUN_OUTPUT

IPDR_REFERENCE=$(tail -n 1 $RUN_OUTPUT)

go run . "$IPDR_REFERENCE" docker-daemon:hello-world-copy:latest

diff <(docker image inspect hello-world) <(docker image inspect hello-world-copy) -q

which ipdr >/dev/null # Rest of the test checks for compatibility with ipdr

docker image rm hello-world-copy &>/dev/null || true

IPDR_OUTPUT=$(mktemp)
ipdr push hello-world | tee $IPDR_OUTPUT

IPDR_REFERENCE_2=$(tail -n 1 $IPDR_OUTPUT)
go run . "ipdr:$IPDR_REFERENCE_2" docker-daemon:hello-world-copy:latest

diff <(docker image inspect hello-world) <(docker image inspect hello-world-copy) -q

# ipdr pull "$IPDR_REFERENCE" - ipdr pull is rather buggy and not quite testable, sadly. Maybe when they fix it one day

kill $(jobs -p) &>/dev/null
