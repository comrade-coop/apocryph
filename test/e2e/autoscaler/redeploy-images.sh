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
helmfile apply -f ../minikube -l name=trustedpods --skip-deps --set withdraw.address="0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f" --set ethKey="0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97"


minikube profile c3
kubectl delete namespace trustedpods
helmfile apply -f ../minikube -l name=trustedpods --skip-deps --set withdraw.address=0xa0Ee7A142d267C1f36714E4a8F75612F20a79720 --set ethKey="0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6"

