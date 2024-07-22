#!/bin/bash 

docker build -t comradecoop/apocryph/server:latest ../../.. --target server

docker build -t comradecoop/apocryph/p2p-helper:latest ../../.. --target p2p-helper

docker build -t comradecoop/apocryph/autoscaler:latest ../../.. --target autoscaler

docker run -d -p 5000:5000 --restart=always --name registry registry:2 || echo "Docker registry already running"

docker tag comradecoop/apocryph/server:latest localhost:5000/comradecoop/apocryph/server:latest
docker push localhost:5000/comradecoop/apocryph/server:latest

docker tag comradecoop/apocryph/p2p-helper:latest localhost:5000/comradecoop/apocryph/p2p-helper:latest
docker push localhost:5000/comradecoop/apocryph/p2p-helper:latest

docker tag comradecoop/apocryph/autoscaler:latest localhost:5000/comradecoop/apocryph/autoscaler:latest
docker push localhost:5000/comradecoop/apocryph/autoscaler:latest

minikube profile c1
kubectl delete namespace trustedpods
helmfile apply -f ../minikube -l name=trustedpods --skip-deps

minikube profile c2
kubectl delete namespace trustedpods
helmfile apply -f ../minikube -l name=trustedpods --skip-deps

minikube profile c3
kubectl delete namespace trustedpods
helmfile apply -f ../minikube -l name=trustedpods --skip-deps
