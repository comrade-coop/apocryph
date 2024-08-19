#!/bin/bash

minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd --driver=virtualbox -p c1
minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd --driver=virtualbox -p c2
minikube start --insecure-registry='host.minikube.internal:5000' --container-runtime=containerd --driver=virtualbox -p c3
