#!/usr/bin/env bash

set -e

trap 'kill $(jobs -p) &>/dev/null' EXIT

if [ "$1" = "teardown" ]; then
   minikube delete
   exit 0
fi

cd "$(dirname "$0")"

which minikube >/dev/null
which helmfile >/dev/null; which helm >/dev/null; which kustomize >/dev/null; which kubectl >/dev/null
which go >/dev/null
which curl >/dev/null; which jq >/dev/null; which xargs >/dev/null; which sed >/dev/null
which ipfs >/dev/null

# TODO: docker pull nginxdemos/nginx-hello:latest && ipdr push nginxdemos/nginx-hello

set -v

[ "$(minikube status -f'{{.Host}}')" = "Running" ] || minikube start --insecure-registry='host.minikube.internal:5000'

helmfile sync || { sleep 16; helmfile sync; }

O_IPFS_PATH=$IPFS_PATH
export IPFS_PATH=$(mktemp ipfs.XXXX --tmpdir -d)
{ while true; do kubectl port-forward --namespace ipfs svc/ipfs-rpc 5004:5001 || true; sleep 1; done; } &
echo /ip4/127.0.0.1/tcp/5004 > $IPFS_PATH/api
ADDRESSES=$(minikube service  -n ipfs ipfs-swarm --url | head -n 1 | sed -E 's|http://(.+):(.+)|["/ip4/\1/tcp/\2", "/ip4/\1/udp/\2/quic", "/ip4/\1/udp/\2/quic-v1", "/ip4/\1/udp/\2/quic-v1/webtransport"]|')
PROVIDERID=$(ipfs id -f '<id>')

ipfs config Addresses.AppendAnnounce --json "$ADDRESSES"

kubectl delete -n ipfs $(kubectl get po -o name -n ipfs) # HACK: Restart ipfs daemon after config change

{ while ! ipfs id -f '<id>' &>/dev/null; do sleep 0.5; done; } 2>/dev/null

kubectl delete -n trustedpods $(kubectl get po -o name -n trustedpods) # HACK: Restart trustedpods after ipfs restart

export IPFS_PATH=$O_IPFS_PATH

ipfs id -f '<id>'

ipfs id &>/dev/null || ipfs init

ipfs config --json Experimental.Libp2pStreamMounting true

ipfs daemon &
{ while ! [ -f ${IPFS_PATH:-~/.ipfs}/api ]; do sleep 0.1; done; } 2>/dev/null

echo "$ADDRESSES" | jq -r '.[] + "/p2p/'"$PROVIDERID"'"' | xargs -n 1 ipfs swarm connect || true

{ while ! kubectl get -n trustedpods endpoints -o json | jq '.items[].subsets[].ports' &>/dev/null; do sleep 0.5; done; }

sleep 1

go run ../../../cmd/trustedpods/ pod deploy manifest.json --format json --provider $PROVIDERID

INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true | head -n 1); echo $INGRESS_URL

MANIFEST_HOST=$(jq -r '.containers[].ports[].hostHttpHost' manifest.json); echo $MANIFEST_HOST

set -x

curl --connect-timeout 40 -H "Host: $MANIFEST_HOST" $INGRESS_URL

sleep 16

kubectl port-forward --namespace prometheus service/prometheus-server 19090:80 &

go run ../../../cmd/tpodserver/ metrics get --prometheus http://127.0.0.1:19090/ --pricing pricing.json
