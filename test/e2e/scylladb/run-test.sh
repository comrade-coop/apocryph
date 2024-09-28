#!/bin/sh

set -e

which curl >/dev/null; which jq >/dev/null; which xargs >/dev/null; which sed >/dev/null
which go >/dev/null
which ipfs >/dev/null
which forge &>/dev/null || export PATH=$PATH:~/.bin/foundry
which forge >/dev/null; which cast >/dev/null
which minikube >/dev/null; which helmfile >/dev/null; which helm >/dev/null; which kubectl >/dev/null
which docker >/dev/null

if [ "$1" = "undeploy" ]; then
  set -v
  
  PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
  PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
  PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
  REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
  FUNDS=10000000000000000000000
  
  go run ../../../cmd/trustedpods/ pod delete ./manifest-lighthouse.yaml ./deployments/lighthouse.json \
    --ethereum-key "$PUBLISHER_KEY"
  go run ../../../cmd/trustedpods/ pod delete ./manifest-node.yaml ./deployments/db1.json \
    --ethereum-key "$PUBLISHER_KEY"
  go run ../../../cmd/trustedpods/ pod delete ./manifest-node.yaml ./deployments/db2.json \
    --ethereum-key "$PUBLISHER_KEY"
  
  
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


echo -e "\e[1;32m---"
echo "Note: This script expects the the minikube e2e test is already running"
echo "  (this might change in the future)"
echo "  (You might need to run: minikube config set cpus 4; minikube config set memory 5000M)"
echo -e "---\e[0m"

# sudo tee /etc/containerd/certs.d/host.minikube.internal:5000/hosts.toml
# server = "https://registry-1.docker.io"
# [host."http://host.minikube.internal:5000"]
# skip_verify = true


set -v

## 2: Prepare nebula configuration ##

mkdir -p ./nebula-config

if ! which nebula-cert > /dev/null; then
  echo "Using nebula-cert via docker!"
  alias nebula-cert="docker run --rm --entrypoint /nebula-cert -v .:/nc -w /nc -u $(id -u):$(id -g) nebulaoss/nebula"
fi

pushd ./nebula-config/

[ -f ./ca.key ] || nebula-cert ca -name "Test ScyllaDB cluster"   
[ -f ./lighthouse.key ] || nebula-cert sign -name "lighthouse" -ip "172.16.123.1/8"
[ -f ./db1.key ] || nebula-cert sign -name "db1" -ip "172.16.123.2/8"
[ -f ./db2.key ] ||  nebula-cert sign -name "db2" -ip "172.16.123.3/8"
[ -f ./local.key ] || nebula-cert sign -name "local" -ip "172.16.123.101/8"

popd

## 3: Deploy everything ##

mkdir -p ./deployments

## 3.1: Deploy Nebula lighthouse ##

PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
FUNDS=10000000000000000000000

#[ -n "$IPFS_DAEMON" ] || { IPFS_DAEMON=yes; ipfs daemon & { while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null; }

go run ../../../cmd/trustedpods/ pod deploy ./manifest-lighthouse.yaml ./deployments/lighthouse.json \
  --ethereum-key "$PUBLISHER_KEY" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --registry-contract "$REGISTRY_CONTRACT" \
  --funds "$FUNDS" \
  --upload-images=false \
  --mint-funds

## 3.2: Deploy first Scylla+Nebula node ##

PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
FUNDS=10000000000000000000000

[ -n "$LIGHTHOUSE_ADDR" ] || { LIGHTHOUSE_ADDR=$(minikube ip):$(jq -r '.deployed.addresses[0].multiaddr' ./deployments/lighthouse.json | sed 's|udp/||'); echo $LIGHTHOUSE_ADDR; }

sed -E 's|\$SERVICEPORT\$|30001|g' <./manifest-node-tmpl.yaml >./manifest-node.yaml # TODO: Would be nice if 
sed -E 's|\$LIGHTHOUSE\$|'$LIGHTHOUSE_ADDR'|g;s|\$SERVICEADDR\$|'$SERVICE_ADDR':30001|g' <./nebula-conf-tmpl.yaml >./nebula-config/config.yml
cp nebula-config/db1.crt nebula-config/host.crt; cp nebula-config/db1.key nebula-config/host.key

go run ../../../cmd/trustedpods/ pod deploy ./manifest-node.yaml ./deployments/db1.json \
  --ethereum-key "$PUBLISHER_KEY" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --registry-contract "$REGISTRY_CONTRACT" \
  --funds "$FUNDS" \
  --upload-images=false \
  --pod-id=0x0000000000000000000000000000000000000000000000000000000000000001 \
  --mint-funds

## 3.3: Deploy second Scylla+Nebula node ##

PROVIDER_ETH=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 #TODO= anvil.accounts[1]
PUBLISHER_KEY=$(docker logs anvil | awk '/Private Keys/ {flag=1; next} flag && /^\(2\)/ {print $2; exit}')
PAYMENT_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.payment.value')
REGISTRY_CONTRACT=$(cat ../../../contracts/broadcast/Deploy.s.sol/31337/run-latest.json | jq -r '.returns.registry.value')
FUNDS=10000000000000000000000

[ -n "$LIGHTHOUSE_ADDR" ] || { LIGHTHOUSE_ADDR=$(minikube ip):$(jq -r '.deployed.addresses[0].multiaddr' ./deployments/lighthouse.json | sed 's|udp/||'); echo $LIGHTHOUSE_ADDR; }
[ -n "$SERVICE_ADDR" ] || { SERVICE_ADDR=$(minikube ip); }

sed -E 's|\$SERVICEPORT\$|30002|g' <./manifest-node-tmpl.yaml >./manifest-node.yaml
sed -E 's|\$LIGHTHOUSE\$|'$LIGHTHOUSE_ADDR'|g;s|\$SERVICEADDR\$|'$SERVICE_ADDR':30002|g' <./nebula-conf-tmpl.yaml >./nebula-config/config.yml
cp nebula-config/db2.crt nebula-config/host.crt; cp nebula-config/db2.key nebula-config/host.key

go run ../../../cmd/trustedpods/ pod deploy ./manifest-node.yaml ./deployments/db2.json \
  --ethereum-key "$PUBLISHER_KEY" \
  --payment-contract "$PAYMENT_CONTRACT" \
  --registry-contract "$REGISTRY_CONTRACT" \
  --funds "$FUNDS" \
  --upload-images=false \
  --pod-id=0x0000000000000000000000000000000000000000000000000000000000000002 \
  --mint-funds
  
rm nebula-config/host.crt nebula-config/host.key

exit 1

## 4: Connect local machine to Nebula ##

if ! which nebula > /dev/null; then
  echo "Using nebula via docker!"
  alias nebula="docker run --name nebula --rm --network host --cap-add NET_ADMIN -v .:/config -w /config nebulaoss/nebula"
fi

[ -n "$LIGHTHOUSE_ADDR" ] || { LIGHTHOUSE_ADDR=$(minikube ip):$(jq -r '.deployed.addresses[0].multiaddr' ./deployments/lighthouse.json | sed 's|udp/||'); echo $LIGHTHOUSE_ADDR; }

sed -E 's|\$LIGHTHOUSE\$|'$LIGHTHOUSE_ADDR'|g;s|\$SERVICEADDR\$|'$(minikube ip | sed -E 's|\.[0-9]*||')'.1:0|g' <./nebula-conf-tmpl.yaml >./nebula-config/config.yml
cp nebula-config/local.crt nebula-config/host.crt; cp nebula-config/local.key nebula-config/host.key
chmod a+r nebula-config/host.crt nebula-config/host.key # .. Docker permissions problem

pushd nebula-config

nebula --config /config/config.yml & #TODO

popd

exit 1
## 4.1: Connect using CQLSH, create some sample values ##

if ! which cqlsh > /dev/null; then
  echo "Using cqlsh via docker! (consider installing it with 'pip install cqlsh')"
  docker start cqlsh-container || docker run --rm -d --name cqlsh-container --network host --entrypoint sleep scylladb/scylla 10000
  # https://stackoverflow.com/a/44739847
  alias cqlsh="docker exec cqlsh-container cqlsh"
fi

cqlsh 172.16.123.2 -e "CREATE KEYSPACE IF NOT EXISTS test WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 2} AND tablets = { 'enabled': false }; CREATE TABLE IF NOT EXISTS test.test1 (id uuid PRIMARY KEY, count counter);"

cqlsh 172.16.123.3 -e "UPDATE test.test1 SET count = count + 4 where id = 6c89c24d-0690-4dac-8f6e-71e98ad2b165;"

cqlsh 172.16.123.2 -e "SELECT * FROM test.test1"
