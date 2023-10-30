#!/bin/sh

cd "$(dirname "$0")"
cd server
go build .

docker build . -t comradecoop/trusted-pods/devserver:latest

docker run -d -p 5000:5000 --restart=always --name registry registry:3 || echo "Docker registry already running"

docker tag comradecoop/trusted-pods/devserver:latest localhost:5000/comradecoop/trusted-pods/devserver:latest
docker push localhost:5000/comradecoop/trusted-pods/devserver:latest

rm server
