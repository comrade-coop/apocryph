#!/bin/sh
ipfs daemon > /dev/null 2>&1 &
sleep 1; echo "CLIENT:waiting ipfs daemon to start"; sleep 12 ;
# or do this
# ipfs p2p forward /x/trusted-pods/provision-pod/0.0.1 /ip4/127.0.0.1/tcp/5000 /p2p/12D3KooWPcMp99mZkfdk8qfHrfjFzQbZcjxSf14egDQtq5zWFhWp
./client
