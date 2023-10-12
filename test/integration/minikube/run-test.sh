#!/usr/bin/env bash

set -e

trap 'kill $(jobs -p) &>/dev/null' EXIT

if [ "$1" = "teardown" ]; then
   minikube delete
   exit 0
fi

cd "$(dirname "$0")"

which minikube >/dev/null
which helmfile >/dev/null; which helm >/dev/null; which kustomize >/dev/null
which go >/dev/null
which curl >/dev/null; which jq >/dev/null; which xargs >/dev/null
which ipfs >/dev/null

# TODO: docker pull nginxdemos/nginx-hello:latest && ipdr push nginxdemos/nginx-hello

set -v

[ "$(minikube status -f'{{.Host}}')" = "Running" ] || minikube start

#helmfile sync || { sleep 10; helmfile sync; }

function as_ { X=$1; shift; IPFS_PATH="$IPFS_ROOT/$X" "$@"; }
IPFS_ROOT=$(mktemp ipfs.XXXX --tmpdir -d); echo $IPFS_ROOT

as_ provider ipfs init --profile test,randomports
as_ provider ipfs config --json Experimental.Libp2pStreamMounting true

as_ publisher ipfs init --profile test,randomports
as_ publisher ipfs config --json Experimental.Libp2pStreamMounting true

as_ provider ipfs daemon & as_ publisher ipfs daemon & sleep 1

as_ provider ipfs id | jq -r '.Addresses[]' | as_ publisher xargs -n 1 ipfs swarm connect

PROVIDERID=$(as_ provider ipfs id | jq -r '.ID')

# go run ../../../cmd/tpodserver/ manifest apply manifest.json --format json
as_ provider go run ../../../cmd/tpodserver/ manifest serve & sleep 5

as_ publisher go run ../../../cmd/trustedpods/ pod deploy manifest.json --format json --provider $PROVIDERID

INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true | head -n 1); echo $INGRESS_URL

MANIFEST_HOST=$(jq -r '.containers[].ports[].hostHttpHost' manifest.json); echo $MANIFEST_HOST

set -x

curl --connect-timeout 40 -H "Host: $MANIFEST_HOST" $INGRESS_URL

sleep 10

kubectl port-forward --namespace prometheus service/prometheus-server 19090:80 &

go run ../../../cmd/tpodserver/ metrics get --prometheus http://127.0.0.1:19090/ --pricing pricing.json
