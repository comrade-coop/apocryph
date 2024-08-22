#!/bin/bash 
./build-images.sh
minikube profile c1
kubectl delete namespace trustedpods
# will use default withdraw address & eth keys specefied in values.yaml
helmfile apply -f ../../minikube -l name=trustedpods --skip-deps 
