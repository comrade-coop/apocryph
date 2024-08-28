#!/bin/bash 

cd "$(dirname "$0")"
./build-images.sh
minikube profile c1
kubectl delete namespace trustedpods

# delete old image policy
kubectl delete ClusterImagePolicy tpod-proxy-policy
# will use default withdraw address & eth keys specefied in values.yaml
helmfile apply -f ../../minikube -l name=trustedpods --skip-deps $1 $2 $3 $4 $5 $6
