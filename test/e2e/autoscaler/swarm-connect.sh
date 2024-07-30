#!/bin/bash

## Configure provider/in-cluster IPFS and publisher IPFS ##
pkill ipfs
pkill -f "kubectl port-forward"

PORT_5004=""
PROVIDER_IPFS=""
IPFS_DAEMON=""

{ while ! kubectl get -n ipfs endpoints ipfs-rpc -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; }

O_IPFS_PATH=$IPFS_PATH
export IPFS_PATH=$(mktemp ipfs.XXXX --tmpdir -d)

kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 &

echo /ip4/127.0.0.1/tcp/5004 > $IPFS_PATH/api

SWARM_ADDRESSES=$(minikube service  -n ipfs ipfs-swarm --url | head -n 1 | sed -E 's|http://(.+):(.+)|["/ip4/\1/tcp/\2", "/ip4/\1/udp/\2/quic", "/ip4/\1/udp/\2/quic-v1", "/ip4/\1/udp/\2/quic-v1/webtransport"]|')

PROVIDER_IPFS=$(curl -X POST "http://127.0.0.1:5004/api/v0/id" | jq '.ID' -r); echo $PROVIDER_IPFS

CONFIG_BEFORE=$(ipfs config Addresses.AppendAnnounce)
ipfs config Addresses.AppendAnnounce --json "$SWARM_ADDRESSES"
CONFIG_AFTER=$(ipfs config Addresses.AppendAnnounce)

[ "$CONFIG_BEFORE" = "$CONFIG_AFTER"  ] || kubectl delete -n ipfs $(kubectl get po -o name -n ipfs) # Restart ipfs daemon

export IPFS_PATH=$O_IPFS_PATH

{ while ! kubectl get -n ipfs endpoints ipfs-rpc -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; }

ipfs id &>/dev/null || ipfs init

ipfs config --json Experimental.Libp2pStreamMounting true

ipfs daemon & { while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; }

echo "$SWARM_ADDRESSES" | jq -r '.[] + "/p2p/'"$PROVIDER_IPFS"'"' | xargs -n 1 ipfs swarm connect || true

sleep 5
