#!/bin/bash

cd "$(dirname "$0")"

trap 'kill $(jobs -p) &>/dev/null' EXIT

which curl >/dev/null; which jq >/dev/null; which xargs >/dev/null; which sed >/dev/null
which go >/dev/null
which ipfs >/dev/null
which forge &>/dev/null || export PATH=$PATH:~/.bin/foundry
which forge >/dev/null; which cast >/dev/null
which minikube >/dev/null; which helmfile >/dev/null; which helm >/dev/null; which kubectl >/dev/null
which docker >/dev/null

set -v

# based on https://stackoverflow.com/a/31269848 / https://bobcopeland.com/blog/2012/10/goto-in-bash/
if [ -n "$1" ]; then
  STEP=${1:-1}
  eval "set -v; $(sed -n "/## $STEP: /{:a;n;p;ba};" $0)"
  exit
fi

sudo chmod o+rw /run/containerd/containerd.sock

## 0: Set up the external environment

## 0.1: Build/tag server and p2p-helper and autoscaler images

docker build -t comradecoop/apocryph/server:latest ../../.. --target server

docker build -t comradecoop/apocryph/p2p-helper:latest ../../.. --target p2p-helper

docker build -t comradecoop/apocryph/autoscaler:latest ../../.. --target autoscaler

## 0.2: Create local registry and push server and p2p-helper images

docker run -d -p 5000:5000 --restart=always --name registry registry:2 || echo "Docker registry already running"

docker tag comradecoop/apocryph/server:latest localhost:5000/comradecoop/apocryph/server:latest
docker push localhost:5000/comradecoop/apocryph/server:latest

docker tag comradecoop/apocryph/p2p-helper:latest localhost:5000/comradecoop/apocryph/p2p-helper:latest
docker push localhost:5000/comradecoop/apocryph/p2p-helper:latest

## 0.3: Set up a local ethereum node and deploy contracts to it

./redeploy-contracts.sh

## 1.0 Starting the First Cluster
minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd --driver=virtualbox -p c1
minikube profile c1
helmfile sync -f ../minikube || { while ! kubectl get -n keda endpoints ingress-nginx-controller -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; helmfile apply; }

# wait until all the deployments are ready
./wait-deployments.sh


## 2.0: Starting the second Cluster
minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd --driver=virtualbox -p c2
minikube profile c2
helmfile sync -f ../minikube || { while ! kubectl get -n keda endpoints ingress-nginx-controller -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; helmfile apply; }

# wait until all the deployments are ready
./wait-deployments.sh


## 3.0: Starting the third Cluster
minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd --driver=virtualbox -p c3
minikube profile c3
helmfile sync -f ../minikube || { while ! kubectl get -n keda endpoints ingress-nginx-controller -o json | jq '.subsets[].addresses[].ip' &>/dev/null; do sleep 1; done; helmfile apply; }

# wait until all the deployments are ready
./wait-deployments.sh


minikube profile list

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

# Connect the three ipfs nodes

## 4.1: Get the tables and the providers  


pkill -f "kubectl port-forward"

ipfs daemon >/dev/null &
go run ../../../cmd/trustedpods registry get --config ../../integration/registry/config.yaml config.yaml \
  --ethereum-key 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
