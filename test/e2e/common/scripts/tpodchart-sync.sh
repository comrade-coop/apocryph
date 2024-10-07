#!/bin/bash
echo "Tpodserver Withdraw address set to $1"
echo "Tpodserver Eth Key set to $2"

kubectl delete namespace trustedpods
helmfile apply -f ../minikube -l name=trustedpods --skip-deps --set withdraw.address="$1" --set ethKey="$2"
