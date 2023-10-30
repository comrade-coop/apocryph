#!/usr/bin/env bash

set -e

trap 'kill $(jobs -p) &>/dev/null' EXIT

if [ "$1" = "teardown" ]; then
   minikube delete
   exit 0
fi

cd "$(dirname "$0")"

which curl >/dev/null; which jq >/dev/null; which xargs >/dev/null; which sed >/dev/null
which go >/dev/null
which ipfs >/dev/null
which forge >/dev/null; which cast >/dev/null
which minikube >/dev/null; which helmfile >/dev/null; which helm >/dev/null; which kustomize >/dev/null; which kubectl >/dev/null

# TODO: docker pull nginxdemos/nginx-hello:latest && ipdr push nginxdemos/nginx-hello

set -v

## 1: Start a Kubernetes environment ##

[ "$(minikube status -f'{{.Host}}')" = "Running" ] || minikube start --insecure-registry='host.minikube.internal:5000'

helmfile sync || { while ! kubectl get -n keda endpoints ingress-nginx-controller -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; helmfile sync; }

## 1.1: Configure provider/in-cluster IPFS ##

{ while ! kubectl get -n ipfs endpoints ipfs-rpc -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; }

O_IPFS_PATH=$IPFS_PATH
export IPFS_PATH=$(mktemp ipfs.XXXX --tmpdir -d)

kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5
echo /ip4/127.0.0.1/tcp/5004 > $IPFS_PATH/api

SWARM_ADDRESSES=$(minikube service  -n ipfs ipfs-swarm --url | head -n 1 | sed -E 's|http://(.+):(.+)|["/ip4/\1/tcp/\2", "/ip4/\1/udp/\2/quic", "/ip4/\1/udp/\2/quic-v1", "/ip4/\1/udp/\2/quic-v1/webtransport"]|')

PROVIDER_IPFS=$(ipfs id -f '<id>'); echo $PROVIDER_IPFS

CONFIG_BEFORE=$(ipfs config Addresses.AppendAnnounce)
ipfs config Addresses.AppendAnnounce --json "$SWARM_ADDRESSES"
CONFIG_AFTER=$(ipfs config Addresses.AppendAnnounce)

[ "$CONFIG_BEFORE" = "$CONFIG_AFTER"  ] || kubectl delete -n ipfs $(kubectl get po -o name -n ipfs) # Restart ipfs daemon

export IPFS_PATH=$O_IPFS_PATH

{ while ! kubectl get -n ipfs endpoints ipfs-rpc -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; }

## 1.2: Configure publisher/local IPFS ##

ipfs id &>/dev/null || ipfs init

ipfs config --json Experimental.Libp2pStreamMounting true

ipfs daemon &
{ while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null

echo "$ADDRESSES" | jq -r '.[] + "/p2p/'"$PROVIDER_P2P"'"' | xargs -n 1 ipfs swarm connect || true

sleep 1

## 2: Deploy contracts to anvil ##

{ while ! kubectl get -n eth endpoints eth-rpc -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; }

kubectl port-forward --namespace eth svc/eth-rpc 8545:8545 &

DEPLOYER_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 #TODO= anvil.accounts[0]

~/.foundry/bin/forge create --root ../../../contracts Payment --private-key $DEPLOYER_KEY --nonce 0 --silent || true
PAYMENT_CONTRACT=0x5FbDB2315678afecb367f032d93F642f64180aa3 # TODO= result of forge create

~/.foundry/bin/forge create --root ../../../contracts MockToken --private-key $DEPLOYER_KEY --nonce 1 --silent || true
TOKEN_CONTACT=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512 # TODO= result of forge create

## 3: Deploy example manifest to cluster ##

PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
PUBLISHER_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a #TODO= anvil.accounts[2]
FUNDS=10000000000000000000000

go run ../../../cmd/trustedpods/ payment mint \
  --ethereum-key $PUBLISHER_KEY \
  --token "$TOKEN_CONTACT" \
  --funds "$FUNDS"

go run ../../../cmd/trustedpods/ pod deploy manifest-guestbook.json --format json \
  --ethereum-key "$PUBLISHER_KEY" \
  --provider "$PROVIDER_IPFS" \
  --provider-eth "$PROVIDER_ETH" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --token "$TOKEN_CONTACT" \
  --funds "$FUNDS"

WITHDRAW_ETH=0x90F79bf6EB2c4f870365E785982E1f101E93b906

echo "Provider balance before:"
cast call "$TOKEN_CONTACT" "balanceOf(address)" "$WITHDRAW_ETH" | cast to-fixed-point 18

INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true | head -n 1); echo $INGRESS_URL

MANIFEST_HOST=$(jq -r '.containers[]?.ports[]?.hostHttpHost' manifest-guestbook.json); echo $MANIFEST_HOST

set -x

curl --connect-timeout 40 -H "Host: $MANIFEST_HOST" $INGRESS_URL
curl -H "Host: $MANIFEST_HOST" $INGRESS_URL/test.html

set +x

sleep 45

echo "Provider withdrawn balance:"
cast call "$TOKEN_CONTACT" "balanceOf(address)" "$WITHDRAW_ETH" | cast to-fixed-point 18

# NOTE: you can run the following to interact with the guestbook
# kubectl port-forward --namespace keda ingress-nginx-controller 1234:80 &
# xdg-open http://guestbook.localhost:1234/
