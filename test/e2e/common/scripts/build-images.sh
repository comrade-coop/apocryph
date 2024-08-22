#!/bin/bash

docker build -t comradecoop/apocryph/server:latest ../../../../ --target server

docker build -t comradecoop/apocryph/p2p-helper:latest ../../../../ --target p2p-helper

docker build -t comradecoop/apocryph/autoscaler:latest ../../../../ --target autoscaler

docker tag comradecoop/apocryph/server:latest localhost:5000/comradecoop/apocryph/server:latest
docker push localhost:5000/comradecoop/apocryph/server:latest

docker tag comradecoop/apocryph/p2p-helper:latest localhost:5000/comradecoop/apocryph/p2p-helper:latest
docker push localhost:5000/comradecoop/apocryph/p2p-helper:latest

docker tag comradecoop/apocryph/autoscaler:latest localhost:5000/comradecoop/apocryph/autoscaler:latest
docker push localhost:5000/comradecoop/apocryph/autoscaler:latest

