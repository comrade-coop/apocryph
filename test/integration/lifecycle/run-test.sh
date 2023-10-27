#!/bin/sh
[ "$(minikube status -f'{{.Host}}')" = "Running" ] || minikube start --insecure-registry='host.minikube.internal:5000'

# helmfile sync

MINIKUBE_IP=$(minikube ip)

SERVER_PORT=$(kubectl get svc tpodserver -n trustedpods -o jsonpath='{.spec.ports[0].nodePort}')

go run publisher.go $MINIKUBE_IP:$SERVER_PORT
