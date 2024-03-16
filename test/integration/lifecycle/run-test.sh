#!/bin/sh
# SPDX-License-Identifier: GPL-3.0

[ "$(minikube status -f'{{.Host}}')" = "Running" ] || minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd

# helmfile sync

anvil > /dev/null &

MINIKUBE_IP=$(minikube ip)

SERVER_PORT=$(kubectl get svc devserver -n devspace -o jsonpath='{.spec.ports[0].nodePort}')
IPFS_PORT=$(kubectl get svc ipfs -n devspace -o jsonpath='{.spec.ports[0].nodePort}')
# known k8s issue, broken pipes when transferring large files
# kubectl port-forward --namespace devspace svc/ipfs 5001:5001 > /dev/null & sleep 0.5;
# kubectl port-forward --namespace devspace svc/ipfs 4001:4001 > /dev/null & sleep 0.5;
# kubectl port-forward --namespace devspace svc/ipfs 8080:8080 > /dev/null & sleep 0.5;

# PROVIDER_IPFS=$(curl -X POST "http://$MINIKUBE_IP:$IPFS_PORT/api/v0/id" | jq '.ID' -r); echo $PROVIDER_IPFS

go run . $MINIKUBE_IP $SERVER_PORT $IPFS_PORT

pkill anvil
