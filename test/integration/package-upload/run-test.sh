#!/bin/sh
docker pull nginx:1.25.2
docker pull hello-world:linux
rm -r /tmp/package 2> /dev/null
go run main.go
echo "Pod Package:"
ls -al mypod
echo "Retreived package from IPFS"
ls -al /tmp/package
rm -r mypod
