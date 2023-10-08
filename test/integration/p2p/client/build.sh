#!/bin/sh
# might replace this with nginx to use with k8s
docker pull hello-world:latest > /dev/null
mkdir client-pod
cp manifest.yaml client-pod
docker save -o client-pod/hello-world.tar hello-world > /dev/null
go build client.go
docker build -t publisher . > /dev/null
rm -r client-pod
rm client
