## common: ##

FROM docker.io/library/golang@sha256:52362e252f452df17c24131b021bf2ebf1c9869f65c28f88ddb326191defea9c as build-common
# 1.21-bookworm

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y protobuf-compiler libgpgme-dev && rm -rf /var/lib/apt/lists/*
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && go install github.com/ethereum/go-ethereum/cmd/abigen@v1.13.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY pkg ./pkg

FROM docker.io/debian@sha256:2bc5c236e9b262645a323e9088dfa3bb1ecb16cc75811daf40a23a824d665be9 as run-common
# bookworm-slim, bookworm-20231120-slim matching golang:1.21-bookworm

RUN apt-get update && apt-get install -y libgpgme11 curl jq && rm -rf /var/lib/apt/lists/*

## p2p-helper: ##

FROM build-common as build-p2p-helper

COPY cmd/ipfs-p2p-helper ./cmd/ipfs-p2p-helper
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/ipfs-p2p-helper ./cmd/ipfs-p2p-helper

FROM run-common as p2p-helper

COPY --from=build-p2p-helper /usr/local/bin/ipfs-p2p-helper /usr/local/bin/ipfs-p2p-helper

ENTRYPOINT ["ipfs-p2p-helper"]

## server: ##

FROM build-common as build-server

COPY cmd/tpodserver ./cmd/tpodserver
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/tpodserver ./cmd/tpodserver

FROM run-common as server

COPY --from=build-server /usr/local/bin/tpodserver /usr/local/bin/tpodserver

ENTRYPOINT ["tpodserver"]
