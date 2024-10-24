# SPDX-License-Identifier: GPL-3.0

## common: ##

FROM docker.io/library/golang:1.23.1-bookworm@sha256:1a5326b07cbab12f4fd7800425f2cf25ff2bd62c404ef41b56cb99669a710a83 as build-dependencies

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y protobuf-compiler libgpgme-dev && rm -rf /var/lib/apt/lists/*
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && go install github.com/ethereum/go-ethereum/cmd/abigen@v1.14.9

FROM build-dependencies AS build-common

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download && go mod verify

COPY pkg ./pkg

FROM docker.io/debian:bookworm-20240904-slim@sha256:a629e796d77a7b2ff82186ed15d01a493801c020eed5ce6adaa2704356f15a1c as run-common
#  matching golang:1.23.1-bookworm above

RUN apt-get update && apt-get install -y libgpgme11 curl jq && rm -rf /var/lib/apt/lists/*

## p2p-helper: ##

FROM build-common AS build-p2p-helper

COPY cmd/ipfs-p2p-helper ./cmd/ipfs-p2p-helper
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/ipfs-p2p-helper ./cmd/ipfs-p2p-helper

FROM run-common AS p2p-helper

COPY --from=build-p2p-helper /usr/local/bin/ipfs-p2p-helper /usr/local/bin/ipfs-p2p-helper

ENTRYPOINT ["ipfs-p2p-helper"]

FROM run-common AS p2p-helper-copy-local

COPY ./bin/ipfs-p2p-helper /usr/local/bin/ipfs-p2p-helper

ENTRYPOINT ["ipfs-p2p-helper"]

## server: ##

FROM build-common AS build-server

COPY cmd/tpodserver ./cmd/tpodserver
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/tpodserver ./cmd/tpodserver
# RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=bind,source=.,target=/app go build -v -o /usr/local/bin/tpodserver ./cmd/tpodserver

FROM run-common AS server

COPY --from=build-server /usr/local/bin/tpodserver /usr/local/bin/tpodserver

ENTRYPOINT ["tpodserver"]

FROM run-common AS server-copy-local

COPY ./bin/tpodserver /usr/local/bin/tpodserver

ENTRYPOINT ["tpodserver"]

## autoscaler: ##

FROM build-common AS build-autoscaler

COPY autoscaler ./autoscaler
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/autoscaler ./autoscaler

FROM run-common AS autoscaler 

COPY --from=build-autoscaler /usr/local/bin/autoscaler /usr/local/bin/autoscaler

ENTRYPOINT ["autoscaler"]

## tpod-proxy: ##

FROM build-common as build-tpod-proxy

COPY pkg/proxy/ ./proxy
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/tpod-proxy ./proxy

FROM run-common as tpod-proxy

COPY --from=build-tpod-proxy /usr/local/bin/tpod-proxy /usr/local/bin/tpod-proxy

ENTRYPOINT ["tpod-proxy"]

FROM run-common AS tpod-proxy-copy-local

COPY ./bin/proxy /usr/local/bin/tpod-proxy

ENTRYPOINT ["tpod-proxy"]
