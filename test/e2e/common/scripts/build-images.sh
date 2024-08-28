#!/bin/bash

docker build -t comradecoop/apocryph/server:latest ../../../../ --target server

docker build -t comradecoop/apocryph/p2p-helper:latest ../../../../ --target p2p-helper

docker build -t comradecoop/apocryph/autoscaler:latest ../../../../ --target autoscaler

docker build -t comradecoop/apocryph/tpod-proxy:latest ../../../../ --target tpod-proxy

docker tag comradecoop/apocryph/server:latest localhost:5000/comradecoop/apocryph/server:latest
docker push localhost:5000/comradecoop/apocryph/server:latest

docker tag comradecoop/apocryph/p2p-helper:latest localhost:5000/comradecoop/apocryph/p2p-helper:latest
docker push localhost:5000/comradecoop/apocryph/p2p-helper:latest

docker tag comradecoop/apocryph/autoscaler:latest localhost:5000/comradecoop/apocryph/autoscaler:latest
docker push localhost:5000/comradecoop/apocryph/autoscaler:latest

# TODO tag,sign & push to a proper registry instead of ttl
# for demenstration purposes, we push tpod-proxy and sign it
docker tag comradecoop/apocryph/tpod-proxy:latest ttl.sh/comradecoop/apocryph/tpod-proxy:5h
docker push ttl.sh/comradecoop/apocryph/tpod-proxy:5h
IMAGE_DIGEST=$(docker inspect --format='{{index .RepoDigests 0}}' comradecoop/apocryph/tpod-proxy)
cosign sign $IMAGE_DIGEST
