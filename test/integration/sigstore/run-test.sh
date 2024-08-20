#!/bin/sh
set -e 
set -v

sudo chmod o+rw /run/containerd/containerd.sock

cd "$(dirname "$0")"

trap 'kill $(jobs -p) &>/dev/null' EXIT

ipfs daemon >/dev/null &
sleep 2

docker tag hello-world ttl.sh/hello-world:1h
docker push ttl.sh/hello-world:1h

go run ../../../cmd/trustedpods pod upload ./manifest.yaml
