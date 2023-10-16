#!/usr/bin/env bash

set -e

cd "$(dirname "$0")"

which docker >/dev/null

set -v

cd ../../..

docker build -t comradecoop/trusted-pods/server .

docker build -t comradecoop/trusted-pods/p2p-helper . -f ./Dockerfile.p2p-helper

docker run -d -p 5000:5000 --restart=always --name registry registry:2 || echo "Docker registry already running"

docker tag comradecoop/trusted-pods/server localhost:5000/comradecoop/trusted-pods/server
docker push localhost:5000/comradecoop/trusted-pods/server

docker tag comradecoop/trusted-pods/p2p-helper localhost:5000/comradecoop/trusted-pods/p2p-helper
docker push localhost:5000/comradecoop/trusted-pods/p2p-helper
