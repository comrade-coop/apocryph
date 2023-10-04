#!/usr/bin/env bash

set -e

if [ "$1" = "teardown" ]; then
   minikube delete
   exit 0
fi

cd "$(dirname "$0")"

which minikube >/dev/null
which helmfile >/dev/null; which helm >/dev/null; which kustomize >/dev/null
which go >/dev/null
which curl >/dev/null; which jq >/dev/null

# TODO: docker pull nginxdemos/nginx-hello:latest && ipdr push nginxdemos/nginx-hello

set -v

[ "$(minikube status -f'{{.Host}}')" = "Running" ] || minikube start

helmfile sync || wait 10; helmfile sync

go run ../../../cmd/tpodserver/ manifest apply manifest.json --config config.yaml --format json

INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true | head -n 1); echo $INGRESS_URL

MANIFEST_HOST=$(jq -r '.podManifest.containers[].ports[].hostHttpHost' manifest.json); echo $MANIFEST_HOST

set -x

curl --connect-timeout 40 -H "Host: $MANIFEST_HOST" $INGRESS_URL
