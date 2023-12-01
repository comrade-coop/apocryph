#!/usr/bin/env bash

set -e

trap 'pkill -f "kubectl port-forward" && kill $(jobs -p) &>/dev/null' EXIT

cd "$(dirname "$0")"

which curl >/dev/null; which jq >/dev/null; which xargs >/dev/null; which sed >/dev/null
which minikube >/dev/null; which helmfile >/dev/null; which helm >/dev/null; which kustomize >/dev/null; which kubectl >/dev/null
which npm >/dev/null

echo "=============================================================================================="
echo "NOTE: This test assumes that the ../minikube/run-test.sh test has been run and wasn't teardown" # FIXME
echo "=============================================================================================="
echo ""

set -v

[ "$PORT_5004" == "" ] && { PORT_5004="yes" ; kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 & sleep 0.5; }
[ "$PORT_8545" == "" ]  && { PORT_8545="yes" ; kubectl port-forward --namespace eth svc/eth-rpc 8545:8545 & }
[ "$PORT_1234" == "" ]  && { PORT_1234="yes" ; kubectl port-forward --namespace keda svc/ingress-nginx-controller 1234:80 & }

export VITE_PROVIDER_MULTIADDR=$(curl -X POST "http://127.0.0.1:5004/api/v0/id" -s | jq '.Addresses[] | select(test("192.168.+webtransport"))' -r); echo $VITE_PROVIDER_MULTIADDR

npm run dev
