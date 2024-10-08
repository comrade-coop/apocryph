#!/usr/bin/env bash
# SPDX-License-Identifier: GPL-3.0

set -e

trap 'pkill -f "kubectl port-forward" && kill $(jobs -p) &>/dev/null' EXIT

cd "$(dirname "$0")"

which curl >/dev/null; which jq >/dev/null; which xargs >/dev/null; which sed >/dev/null
which go >/dev/null
which ipfs >/dev/null
which forge &>/dev/null || export PATH=$PATH:~/.bin/foundry
which forge >/dev/null; which cast >/dev/null
which helmfile >/dev/null; which helm >/dev/null; which kubectl >/dev/null
which constellation >/dev/null || { echo "Install Constellation, https://docs.edgeless.systems/constellation/getting-started/first-steps-local#software-installation-on-ubuntu"; exit 1; }

CONSTELLATION_PATH=${CONSTELLATION_PATH:-~/.constellation-root}
mkdir -p $CONSTELLATION_PATH

if [ "$1" = "teardown" ]; then
   ( cd $CONSTELLATION_PATH; constellation mini down )
   exit 0
fi

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

## 1: Set up the Kubernetes environment ##

echo 'CONSTELLATION_PATH='$CONSTELLATION_PATH

( cd $CONSTELLATION_PATH; constellation mini up || true )

kubectl patch -n kube-system configmap ip-masq-agent --type merge -p '{"data":{"config": "{\"masqLinkLocal\":true,\"nonMasqueradeCIDRs\":[]}"}}'
kubectl rollout  restart -n kube-system daemonset cilium
kubectl delete pod -l k8s-app=join-service -n kube-system

## 1.0: Deploy contracts to anvil ##

helmfile apply -l name=eth

{ while ! kubectl get -n eth endpoints eth-rpc -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; }

[ "$PORT_8545" == "" ] && { PORT_8545="yes" ; kubectl port-forward --namespace eth svc/eth-rpc 8545:8545 & }

DEPLOYER_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 #TODO= anvil.accounts[0]

( cd ../../../contracts; forge script script/Deploy.s.sol --private-key "$DEPLOYER_KEY" --rpc-url http://localhost:8545 --broadcast)

## 1.1: Apply the rest of the Helm configuration ##

helmfile apply 

## 1.2: Configure provider/in-cluster IPFS and publisher IPFS ##

{ while ! kubectl get -n ipfs endpoints ipfs-rpc -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; }

[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }

NODE_ADDRESS=$(kubectl get no -o json | jq -r '.items[].status.addresses[] | select(.type == "InternalIP") | .address' | head -n 1)
SWARM_PORT=$(kubectl get svc -n ipfs ipfs-swarm -o json | jq -r '.spec.ports[].nodePort' | head -n 1)

SWARM_ADDRESSES="[\"/ip4/$NODE_ADDRESS/tcp/$SWARM_PORT\", \"/ip4/$NODE_ADDRESS/udp/$SWARM_PORT/quic\", \"/ip4/$NODE_ADDRESS/udp/$SWARM_PORT/quic-v1\", \"/ip4/$NODE_ADDRESS/udp/$SWARM_PORT/quic-v1/webtransport\"]"

PROVIDER_IPFS=$(curl -X POST "http://127.0.0.1:5004/api/v0/id" | jq '.ID' -r); echo $PROVIDER_IPFS

# Unfortunately, we can't restart the ipfs daemon since we don't have persistent storage in miniconstellation. Swarm addresses have been hardcoded.
#O_IPFS_PATH=$IPFS_PATH
#export IPFS_PATH=$(mktemp ipfs.XXXX --tmpdir -d)
#echo /ip4/127.0.0.1/tcp/5004 > $IPFS_PATH/api

#CONFIG_BEFORE=$(ipfs config Addresses.AppendAnnounce)
#ipfs config Addresses.AppendAnnounce --json "$SWARM_ADDRESSES"
#CONFIG_AFTER=$(ipfs config Addresses.AppendAnnounce)

#[ "$CONFIG_BEFORE" = "$CONFIG_AFTER"  ] || false # Restart ipfs daemon
#export IPFS_PATH=$O_IPFS_PATH

ipfs id &>/dev/null || ipfs init

ipfs config --json Experimental.Libp2pStreamMounting true

[ -n "$IPFS_DAEMON" ] || { IPFS_DAEMON=yes; ipfs daemon & { while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null; }

echo "$SWARM_ADDRESSES" | jq -r '.[] + "/p2p/'"$PROVIDER_IPFS"'"' | xargs -n 1 ipfs swarm connect || true

sleep 1

## 1.3: Register the provider

go run ../../../cmd/tpodserver/  registry  register \
  --config ../../common/configs/config.yaml \
  --ipfs /ip4/127.0.0.1/tcp/5001 \
  --ethereum-rpc http://127.0.0.1:8545 \
  --ethereum-key 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \

## 2: Deploy example manifest to cluster ##

DEPLOYER_ETH=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 #TODO= anvil.accounts[0]
PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
PUBLISHER_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a #TODO= anvil.accounts[2]
PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
FUNDS=10000000000000000000000

[ "$PORT_8545" == "" ] && { PORT_8545="yes" ; kubectl port-forward --namespace eth svc/eth-rpc 8545:8545 & }
[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }
[ -n "$PROVIDER_IPFS" ] || { PROVIDER_IPFS=$(curl -X POST "http://127.0.0.1:5004/api/v0/id" -s | jq '.ID' -r); echo $PROVIDER_IPFS; }
[ -n "$IPFS_DAEMON" ] || { IPFS_DAEMON=yes; ipfs daemon & { while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null; }

set +v
set -x

go run ../../../cmd/trustedpods/ pod deploy ../common/manifest-guestbook-nostorage.yaml \
  --ethereum-key "$PUBLISHER_KEY" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --registry-contract "$REGISTRY_CONTRACT" \
  --funds "$FUNDS" \
  --upload-images=false \
  --mint-funds

set +x
set -v

## 3: Connect and measure balances ##

DEPLOYER_ETH=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 #TODO= anvil.accounts[0]
WITHDRAW_ETH=0x90F79bf6EB2c4f870365E785982E1f101E93b906 # From trustedpods/tpodserver.yml
TOKEN_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.token.value')
NODE_ADDRESS=$(kubectl get no -o json | jq -r '.items[].status.addresses[] | select(.type == "InternalIP") | .address' | head -n 1)
INGRESS_PORT=$(kubectl get svc -n keda ingress-nginx-controller -o json | jq -r '.spec.ports[] | select(.name == "http") | .nodePort' | head -n 1)
INGRESS_URL="http://$NODE_ADDRESS:$INGRESS_PORT"; echo $INGRESS_URL
MANIFEST_HOST=guestbook.localhost # From manifest-guestbook.yaml

[ "$PORT_8545" == "" ] && { PORT_8545="yes" ; kubectl port-forward --namespace eth svc/eth-rpc 8545:8545 & }

echo "Provider balance before:" $(cast call "$TOKEN_CONTRACT" "balanceOf(address)" "$WITHDRAW_ETH" | cast to-fixed-point 18)

set -x

while ! curl --connect-timeout 40 -H "Host: $MANIFEST_HOST" $INGRESS_URL --fail-with-body; do sleep 10; done
curl -H "Host: $MANIFEST_HOST" $INGRESS_URL/test.html --fail-with-body

set +x

sleep 45

echo "Provider balance after:" $(cast call "$TOKEN_CONTRACT" "balanceOf(address)" "$WITHDRAW_ETH" | cast to-fixed-point 18)

## 4: In conclusion.. ##
 
set +v

echo -e "\e[1;32m---"
echo "Note: To interact with the deployed guestbook, run the following"
echo "  kubectl port-forward --namespace keda svc/ingress-nginx-controller 1234:80 &"
echo "  xdg-open http://guestbook.localhost:1234/"
echo "Note: To stop the minikube cluster/provider, use '$0 teardown'"
echo "  and to clean-up everything the script does, use '$0 teardown full'"
echo -e "---\e[0m"
