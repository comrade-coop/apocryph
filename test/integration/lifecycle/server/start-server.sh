#!/bin/sh
ipfs daemon > /dev/null 2>&1 &
sleep 1; echo "tpodserver: waiting ipfs daemon to start"; sleep 12 ;
ipfs p2p listen /x/trusted-pods/provision-pod/0.0.1 /ip4/127.0.0.1/tcp/$SERVER_PORT  > /dev/null 2>&1
./server $SERVER_ADDRESS
