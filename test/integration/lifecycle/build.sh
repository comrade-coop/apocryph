#!/bin/sh
# SPDX-License-Identifier: GPL-3.0

cd "$(dirname "$0")"

cd server
go build .

docker build . -t comradecoop/apocryph/devserver:latest

docker run -d -p 5000:5000 --restart=always --name registry registry:3 || echo "Docker registry already running"

docker tag comradecoop/apocryph/devserver:latest localhost:5000/comradecoop/apocryph/devserver:latest
docker push localhost:5000/comradecoop/apocryph/devserver:latest

rm server
