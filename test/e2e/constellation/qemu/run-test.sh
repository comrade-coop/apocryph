#!/bin/bash
set -e
set -v

WORKSPACE_PATH="$HOME/.apocryph/constellation-"
CURRENT_DIR=$(pwd)

trap 'pkill -f "kubectl port-forward" && kill $(jobs -p) &>/dev/null' EXIT

# based on https://stackoverflow.com/a/31269848 / https://bobcopeland.com/blog/2012/10/goto-in-bash/
if [ -n "$1" ]; then
  STEP=${1:-1}
  eval "set -v; $(sed -n "/## $STEP: /{:a;n;p;ba};" $0)"
  exit
fi

echo -e "\e[1;32m---"
echo "Note: To skip steps, use '$0 <number>'"
echo -e "---\e[0m"

## 0: Build build custom OS image & run the cluster

# use the miniconstellation chart
. ./build.sh ../miniconstellation

## 1: Start the constellation cluster
cd "$WORKSPACE_PATH"*
constellation apply -y

## 1.1: wait for the setup
cd $WORKSPACE_PATH*
CONF_DIR=$(pwd)
export KUBECONFIG="$CONF_DIR/constellation-admin.conf"
cd $CURRENT_DIR
kubectl wait --namespace keda --for=condition=available deployment/ingress-nginx-controller
kubectl wait --namespace prometheus --for=condition=available deployment/prometheus-kube-state-metrics
kubectl wait --namespace prometheus --for=condition=available deployment/prometheus-prometheus-pushgateway
kubectl wait --namespace prometheus --for=condition=available deployment/prometheus-server
kubectl wait --namespace eth --for=condition=available deployment/anvil
kubectl wait --namespace ipfs --for=condition=available deployment/ipfs
kubectl wait --namespace trustedpods --for=condition=available deployment/tpodserver

## 2: Deploy sample App

[ "$PORT_8545" == "" ] && { PORT_8545="yes" ; kubectl port-forward --namespace eth svc/eth-rpc 8545:8545 & }

DEPLOYER_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 #TODO= anvil.accounts[0]

TOKEN_CONTRACT=0x5fbdb2315678afecb367f032d93f642f64180aa3 # TODO= result of forge create

forge create --root ../../../../contracts MockToken --private-key $DEPLOYER_KEY
forge create --root ../../../../contracts Payment --private-key $DEPLOYER_KEY --constructor-args $TOKEN_CONTRACT
forge create --root ../../../../contracts Registry --private-key $DEPLOYER_KEY


## 2.1: Register the provider
[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }

go run ../../../../cmd/tpodserver  registry  register \
  --config ../../common/config.yaml \
  --ipfs /ip4/127.0.0.1/tcp/5004 \
  --ethereum-rpc http://127.0.0.1:8545 \
  --ethereum-key 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \

## 3: Configure provider/in-cluster IPFS and publisher IPFS ##

{ while ! kubectl get -n ipfs endpoints ipfs-rpc -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; }

[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }

NODE_ADDRESS=$(kubectl get no -o json | jq -r '.items[].status.addresses[] | select(.type == "InternalIP") | .address' | head -n 1)
SWARM_PORT=$(kubectl get svc -n ipfs ipfs-swarm -o json | jq -r '.spec.ports[].nodePort' | head -n 1)

SWARM_ADDRESSES="[\"/ip4/$NODE_ADDRESS/tcp/$SWARM_PORT\", \"/ip4/$NODE_ADDRESS/udp/$SWARM_PORT/quic\", \"/ip4/$NODE_ADDRESS/udp/$SWARM_PORT/quic-v1\", \"/ip4/$NODE_ADDRESS/udp/$SWARM_PORT/quic-v1/webtransport\"]"

PROVIDER_IPFS=$(curl -X POST "http://127.0.0.1:5004/api/v0/id" | jq '.ID' -r); echo $PROVIDER_IPFS

ipfs id &>/dev/null || ipfs init

ipfs config --json Experimental.Libp2pStreamMounting true

IPFS_CONFIG="$HOME/.ipfs/config"
# default gateway is already used by libvirt container
NEW_GATEWAY="/ip4/127.0.0.1/tcp/8090"
jq --arg new_gateway "$NEW_GATEWAY" '.Addresses.Gateway = $new_gateway' "$IPFS_CONFIG" > tmp.json && mv tmp.json "$IPFS_CONFIG"

[ -n "$IPFS_DAEMON" ] || { IPFS_DAEMON=yes; ipfs daemon & { while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null; }

echo "$SWARM_ADDRESSES" | jq -r '.[] + "/p2p/'"$PROVIDER_IPFS"'"' | xargs -n 1 ipfs swarm connect || true

sleep 1

# Deploy example manifest to cluster ##

DEPLOYER_ETH=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 #TODO= anvil.accounts[0]
PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
PUBLISHER_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a #TODO= anvil.accounts[2]

TOKEN_CONTRACT=0x5fbdb2315678afecb367f032d93f642f64180aa3 # TODO= result of forge create
PAYMENT_CONTRACT=0xe7f1725e7734ce288f8367e1bb143e90bb3f0512 # TODO= result of forge create
REGISTRY_CONTRACT=0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 # TODO= result of forge create


# PAYMENT_CONTRACT=$(cat ../../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
# REGISTRY_CONTRACT=$(cat ../../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
FUNDS=10000000000000000000000


[ "$PORT_8545" == "" ] && { PORT_8545="yes" ; kubectl port-forward --namespace eth svc/eth-rpc 8545:8545 & }
[ -n "$IPFS_DAEMON" ] || { IPFS_DAEMON=yes; ipfs daemon & { while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null; }

set +v
set -x


sudo chmod o+rw /run/containerd/containerd.sock


go run ../../../../cmd/trustedpods/ pod deploy ../../common/manifest-redmine-nostorage.yaml \
  --ethereum-key "$PUBLISHER_KEY" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --registry-contract "$REGISTRY_CONTRACT" \
  --funds "$FUNDS" \
  --upload-images=true \
  --mint-funds

## 4: Connect and measure balances ##

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
curl -H "Host: $MANIFEST_HOST" $INGRESS_URL --fail-with-body

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
