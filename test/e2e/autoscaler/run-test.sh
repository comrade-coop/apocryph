#!/bin/bash

cd "$(dirname "$0")"

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

# (NOTE: Unfortunately, we cannot use a port other than 8545, or otherwise the eth-rpc service will break)
docker run -d -p 8545:8545 --restart=always --name=anvil \
  ghcr.io/foundry-rs/foundry:nightly-619f3c56302b5a665164002cb98263cd9812e4d5 \
  -- 'anvil --host 0.0.0.0 --state /anvil-state.json' 2>/dev/null || {
    docker exec anvil ash -c 'kill 1 && rm -f /anvil-state.json' # Reset anvil state
}
sleep 5

# deploy the contracts
DEPLOYER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(0\)/ {print $2; exit}') # anvil.accounts[0]
( cd ../../../contracts; forge script script/Deploy.s.sol --private-key "$DEPLOYER_KEY" --rpc-url http://localhost:8545 --broadcast)

## 1.0 Starting the First Cluster
echo "Starting the first cluster"
minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd --driver=kvm2 -p c1
minikube profile c1
./deploy.sh

# wait until all the deployments are ready
./wait-deployments.sh

## 1.1: Register the provider in the marketplace
[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }
go run ../../../cmd/tpodserver  registry  register \
  --config ../common/config.yaml \
  --ipfs /ip4/127.0.0.1/tcp/5004 \
  --ethereum-rpc http://127.0.0.1:8545 \
  --ethereum-key 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \


## 2.0: Starting the second Cluster
echo "Starting the second cluster"
minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd --driver=kvm2 -p c2
minikube profile c2
./deploy.sh

# wait until all the deployments are ready
./wait-deployments.sh

[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }
go run ../../../cmd/tpodserver  registry  register \
  --config ../common/config.yaml \
  --ipfs /ip4/127.0.0.1/tcp/5004 \
  --ethereum-rpc http://127.0.0.1:8545 \
  --ethereum-key 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \


## 3.0: Starting the third Cluster
echo "Starting the third cluster"
minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd --driver=kvm2 -p c3
minikube profile c3
./deploy.sh

# wait until all the deployments are ready
./wait-deployments.sh

[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5005:5001 & sleep 0.5; }
go run ../../../cmd/tpodserver  registry  register \
  --config ../common/config.yaml \
  --ipfs /ip4/127.0.0.1/tcp/5005 \
  --ethereum-rpc http://127.0.0.1:8545 \
  --ethereum-key 0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6 \
  --token-contract 0x5FbDB2315678afecb367f032d93F642f64180aa3 \
  --registry-contract 0x9fe46736679d2d9a65f0992f2272de9f3c7fa6e0 \


minikube profile list
