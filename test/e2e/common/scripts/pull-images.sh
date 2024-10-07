#!/bin/bash
docker pull registry:2
docker pull ghcr.io/foundry-rs/foundry:nightly-619f3c56302b5a665164002cb98263cd9812e4d5

docker run -d -p 5000:5000 --restart=always --name registry registry:2 || echo "Docker registry already running"
