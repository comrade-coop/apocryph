#!/bin/bash

cd "$(dirname "$0")"
set -v

sudo chmod o+rw /run/containerd/containerd.sock

trap 'kill $(jobs -p) &>/dev/null' EXIT

which curl >/dev/null; which jq >/dev/null; which xargs >/dev/null; which sed >/dev/null
which go >/dev/null
which ipfs >/dev/null
which forge &>/dev/null || export PATH=$PATH:~/.bin/foundry
which forge >/dev/null; which cast >/dev/null
which minikube >/dev/null; which helmfile >/dev/null; which helm >/dev/null; which kubectl >/dev/null
which docker >/dev/null
which virtualbox >/dev/null


# based on https://stackoverflow.com/a/31269848 / https://bobcopeland.com/blog/2012/10/goto-in-bash/
if [ -n "$1" ]; then
  STEP=${1:-1}
  eval "set -v; $(sed -n "/## $STEP: /{:a;n;p;ba};" $0)"
  exit
fi


## 0: Set up the external environment

## 0.1: Build/tag server and p2p-helper and autoscaler images
./redeploy-images.sh

## 0.2: Set up a local ethereum node and deploy contracts to it

./redeploy-contracts.sh

## 0.3: start clusters
./start-clusters.sh

## 1.0 Starting the First Cluster
minikube profile c1
helmfile sync -f ../minikube || { while ! kubectl get -n keda endpoints ingress-nginx-controller -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; helmfile apply -f ../minikube; }

# wait until all the deployments are ready
./wait-deployments.sh


## 2.0: Starting the second Cluster
minikube profile c2
helmfile sync -f ../minikube || { while ! kubectl get -n keda endpoints ingress-nginx-controller -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; helmfile apply -f ../minikube; }

# wait until all the deployments are ready
./wait-deployments.sh


## 3.0: Starting the third Cluster
minikube profile c3
helmfile sync -f ../minikube || { while ! kubectl get -n keda endpoints ingress-nginx-controller -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; helmfile apply -f ../minikube; }

# wait until all the deployments are ready
./wait-deployments.sh


minikube profile list

sleep 5

## 4.0: Register the providers in the marketplace

minikube profile c1
pkill -f "kubectl port-forward"
kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5;
go run ../../../cmd/tpodserver  registry  register \
  --config ../common/config.yaml \
  --ipfs /ip4/127.0.0.1/tcp/5004 \
  --ethereum-rpc http://127.0.0.1:8545 \
  --ethereum-key 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \


minikube profile c2
pkill -f "kubectl port-forward"
kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5;
go run ../../../cmd/tpodserver  registry  register \
  --config ../common/config2.yaml \
  --ipfs /ip4/127.0.0.1/tcp/5004 \
  --ethereum-rpc http://127.0.0.1:8545 \
  --ethereum-key 0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97 \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \


minikube profile c3
pkill -f "kubectl port-forward"
kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5;

go run ../../../cmd/tpodserver  registry  register \
  --config ../common/config3.yaml \
  --ipfs /ip4/127.0.0.1/tcp/5004 \
  --ethereum-rpc http://127.0.0.1:8545 \
  --ethereum-key 0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6 \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \

## 4.1: Get the tables and the providers  

pkill -f "kubectl port-forward"
kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5;

go run ../../../cmd/trustedpods registry get --config ../../integration/registry/config.yaml config.yaml \
  --ethereum-key 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --ipfs /ip4/127.0.0.1/tcp/5004 \

sleep 5
## 5.0: Deploy the autoscaler to the providers using their p2p multiaddr

set -v
rm -f ~/.apocryph/deployment/*

PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
FUNDS=10000000000000000000000

minikube profile c1
source swarm-connect.sh

PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
echo $PROVIDER_IPFS

go run ../../../cmd/trustedpods/ pod deploy ../common/manifest-autoscaler.yaml \
  --ethereum-key "$PUBLISHER_KEY" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --registry-contract "$REGISTRY_CONTRACT" \
  --funds "$FUNDS" \
  --upload-images=true \
  --mint-funds \
  --provider /p2p/"$PROVIDER_IPFS" \
  --provider-eth "$PROVIDER_ETH" \
  --authorize
sleep 5

## 5.2: deploy to the second cluster
minikube profile c2
source swarm-connect.sh
# for now just remove the deployment file to avoid uploading instead of deploying
rm -f ~/.apocryph/deployment/*

PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
FUNDS=10000000000000000000000

PROVIDER_ETH=0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f #TODO= anvil.accounts[7]
echo $PROVIDER_IPFS

go run ../../../cmd/trustedpods/ pod deploy ../common/manifest-autoscaler.yaml \
  --ethereum-key "$PUBLISHER_KEY" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --registry-contract "$REGISTRY_CONTRACT" \
  --funds "$FUNDS" \
  --upload-images=true \
  --mint-funds \
  --provider /p2p/"$PROVIDER_IPFS" \
  --provider-eth "$PROVIDER_ETH" \
  --authorize

sleep 5

## 5.3: deploy to the third cluster
minikube profile c3
source swarm-connect.sh
rm -f ~/.apocryph/deployment/*

PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
FUNDS=10000000000000000000000

PROVIDER_ETH=0xa0Ee7A142d267C1f36714E4a8F75612F20a79720 #TODO= anvil.accounts[8]
FUNDS=10000000000000000000000

echo $PROVIDER_IPFS

go run ../../../cmd/trustedpods/ pod deploy ../common/manifest-autoscaler.yaml \
  --ethereum-key "$PUBLISHER_KEY" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --registry-contract "$REGISTRY_CONTRACT" \
  --funds "$FUNDS" \
  --upload-images=true \
  --mint-funds \
  --provider /p2p/"$PROVIDER_IPFS" \
  --provider-eth "$PROVIDER_ETH" \
  --authorize

## 6.0: Connect the cluster

minikube profile c1
C1_INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true | head -n 1); echo $C1_INGRESS_URL
minikube profile c2
C2_INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true | head -n 1); echo $C2_INGRESS_URL
minikube profile c3
C3_INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true | head -n 1); echo $C3_INGRESS_URL

go run ../../../cmd/trustedpods/ autoscale --url "$C1_INGRESS_URL" --providers "$C1_INGRESS_URL","$C2_INGRESS_URL","$C3_INGRESS_URL"

## 6.1: Get Last deployed Autoscaler Logs
PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
go run ../../../cmd/trustedpods/ pod log ../common/manifest-autoscaler.yaml --ethereum-key "$PUBLISHER_KEY"
