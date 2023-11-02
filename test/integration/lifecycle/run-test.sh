#!/bin/sh
[ "$(minikube status -f'{{.Host}}')" = "Running" ] || minikube start --insecure-registry='host.minikube.internal:5000'

# helmfile sync

MINIKUBE_IP=$(minikube ip)

SERVER_PORT=$(kubectl get svc devserver -n devspace -o jsonpath='{.spec.ports[0].nodePort}')
go run publisher.go $MINIKUBE_IP:$SERVER_PORT

rm -r ./keystore
