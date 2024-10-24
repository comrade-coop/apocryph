#!/bin/sh

ipfs config --json Experimental.Libp2pStreamMounting true
ipfs config Addresses.Gateway /ip4/127.0.0.1/tcp/8082
