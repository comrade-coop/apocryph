#!/usr/bin/env bash
# SPDX-License-Identifier: GPL-3.0

set -e

trap 'pkill -f "kubectl port-forward" && kill $(jobs -p) &>/dev/null' EXIT

if [ "$1" = "teardown" ]; then
  minikube delete

  if [ "$2" = "full" ]; then
    docker rm --force registry
    docker rm --force anvil
    rm -r /home/ezio/.apocryph/deployment/
    exit 0
  fi
  exit 0
fi

cd "$(dirname "$0")"

which curl >/dev/null; which jq >/dev/null; which xargs >/dev/null; which sed >/dev/null
which go >/dev/null
which ipfs >/dev/null
which forge &>/dev/null || export PATH=$PATH:~/.bin/foundry
which forge >/dev/null; which cast >/dev/null
which minikube >/dev/null; which helmfile >/dev/null; which helm >/dev/null; which kubectl >/dev/null
which docker >/dev/null

# based on https://stackoverflow.com/a/31269848 / https://bobcopeland.com/blog/2012/10/goto-in-bash/
if [ -n "$1" ]; then
  STEP=${1:-1}
  eval "set -v; $(sed -n "/## $STEP: /{:a;n;p;ba};" $0)"
  exit
fi

echo -e "\e[1;32m---"
echo "Note: To skip steps, use '$0 <number>'"
echo "  e.g. to skip ahead to configuring IPFS, run '$0 1.2'"
echo -e "---\e[0m"

set -v

sudo chmod o+rw /run/containerd/containerd.sock

## 0: Set up the external environment

## 0.1: Build/tag server and p2p-helper images

docker build -t comradecoop/apocryph/server:latest ../../.. --target server

docker build -t comradecoop/apocryph/p2p-helper:latest ../../.. --target p2p-helper

## 0.2: Create local registry and push server and p2p-helper images

docker run -d -p 5000:5000 --restart=always --name registry registry:2 || echo "Docker registry already running"

docker tag comradecoop/apocryph/server:latest localhost:5000/comradecoop/apocryph/server:latest
docker push localhost:5000/comradecoop/apocryph/server:latest

docker tag comradecoop/apocryph/p2p-helper:latest localhost:5000/comradecoop/apocryph/p2p-helper:latest
docker push localhost:5000/comradecoop/apocryph/p2p-helper:latest

## 0.3: Set up a local ethereum node and deploy contracts to it

# (NOTE: Unfortunately, we cannot use a port other than 8545, or otherwise the eth-rpc service will break)
docker run -d -p 8545:8545 --restart=always --name=anvil \
  ghcr.io/foundry-rs/foundry:nightly-619f3c56302b5a665164002cb98263cd9812e4d5 \
  -- 'anvil --host 0.0.0.0 --state /anvil-state.json' 2>/dev/null || {
    docker exec anvil ash -c 'kill 1 && rm -f /anvil-state.json' # Reset anvil state
}
sleep 5

DEPLOYER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(0\)/ {print $2; exit}') # anvil.accounts[0]

( cd ../../../contracts; forge script script/Deploy.s.sol --private-key "$DEPLOYER_KEY" --rpc-url http://localhost:8545 --broadcast)

## 1: Set up the Kubernetes environment ##


[ "$(minikube status -f'{{.Kubelet}}')" = "Running" ] || minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd

minikube addons enable metrics-server

## 1.1: Apply Helm configurations ##

kubectl delete namespace trustedpods 2>/dev/null || true

helmfile sync 

## 1.2: Register the provider in the marketplace

[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }

go run ../../../cmd/tpodserver  registry  register \
  --config ../common/config.yaml \
  --ipfs /ip4/127.0.0.1/tcp/5004 \
  --ethereum-rpc http://127.0.0.1:8545 \
  --ethereum-key 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \

## 1.3: Configure provider/in-cluster IPFS and publisher IPFS ##

{ while ! kubectl get -n ipfs endpoints ipfs-rpc -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; }

O_IPFS_PATH=$IPFS_PATH
export IPFS_PATH=$(mktemp ipfs.XXXX --tmpdir -d)

[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }
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

[ -n "$IPFS_DAEMON" ] || { IPFS_DAEMON=yes; ipfs daemon & { while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null; }

echo "$SWARM_ADDRESSES" | jq -r '.[] + "/p2p/'"$PROVIDER_IPFS"'"' | xargs -n 1 ipfs swarm connect || true

sleep 1



## 2: Deploy example manifest to cluster ##

PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
FUNDS=10000000000000000000000

[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }
[ -n "$PROVIDER_IPFS" ] || { PROVIDER_IPFS=$(curl -X POST "http://127.0.0.1:5004/api/v0/id" -s | jq '.ID' -r); echo $PROVIDER_IPFS; }
[ -n "$IPFS_DAEMON" ] || { IPFS_DAEMON=yes; ipfs daemon & { while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null; }

set +v
set -x


go run ../../../cmd/trustedpods/ pod deploy ../common/manifest-nginx.yaml \
  --ethereum-key "$PUBLISHER_KEY" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --registry-contract "$REGISTRY_CONTRACT" \
  --funds "$FUNDS" \
  --upload-images=true \
  --mint-funds

set +x
set -v

## 3: Connect and measure balances ##

WITHDRAW_ETH=0x90F79bf6EB2c4f870365E785982E1f101E93b906 #TODO copied from trustedpods/tpodserver.yml
TOKEN_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.token.value')
INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true | head -n 1); echo $INGRESS_URL
MANIFEST_HOST=example.local # From manifest-nginx.yaml

echo "Provider balance before:" $(cast call "$TOKEN_CONTRACT" "balanceOf(address)" "$WITHDRAW_ETH" | cast to-fixed-point 18)

set -x

while ! curl --connect-timeout 40 -H "Host: $MANIFEST_HOST" $INGRESS_URL --fail-with-body; do sleep 10; done
curl -H "Host: $MANIFEST_HOST" $INGRESS_URL --fail-with-body

set +x

sleep 45

echo "Provider balance after:" $(cast call "$TOKEN_CONTRACT" "balanceOf(address)" "$WITHDRAW_ETH" | cast to-fixed-point 18)

## 4: In conclusion..

set +v

echo -e "\e[1;32m---"
echo "Note: To interact with the deployed guestbook, run the following"
echo "  kubectl port-forward --namespace keda svc/ingress-nginx-controller 1234:80 &"
echo "  xdg-open http://guestbook.localhost:1234/"
echo "Note: To stop the minikube cluster/provider, use '$0 teardown'"
echo "  and to clean-up everything the script does, use '$0 teardown full'"
echo -e "---\e[0m"

exit 0;

## Env: (debug stuff)

export PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
export PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
export TOKEN_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.token.value')
export PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
export REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
export FUNDS=10000000000000000000000
export INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true | head -n 1); echo $INGRESS_URL
export MANIFEST_HOST=guestbook.localhost # From manifest-guestbook.yaml
[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }
[ -n "$PROVIDER_IPFS" ] || { PROVIDER_IPFS=$(curl -X POST "http://127.0.0.1:5004/api/v0/id" -s | jq '.ID' -r); echo $PROVIDER_IPFS; }
[ -n "$IPFS_DAEMON" ] || { IPFS_DAEMON=yes; ipfs daemon & { while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null; }

bash
