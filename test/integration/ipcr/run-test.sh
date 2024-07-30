#!/bin/sh
sudo chmod o+rw /run/containerd/containerd.sock
ipfs daemon &
sleep 2
go run .
pkill ipfs
