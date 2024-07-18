#!/bin/bash

kubectl wait --namespace keda --for=condition=available deployment/ingress-nginx-controller
kubectl wait --namespace prometheus --for=condition=available deployment/prometheus-kube-state-metrics
kubectl wait --namespace prometheus --for=condition=available deployment/prometheus-prometheus-pushgateway
kubectl wait --namespace prometheus --for=condition=available deployment/prometheus-server
kubectl wait --namespace ipfs --for=condition=available StatefulSet/ipfs
kubectl wait --namespace trustedpods --for=condition=available deployment/tpodserver
