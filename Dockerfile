# common:

# FROM ghcr.io/foundry-rs/foundry:v1.0.0 as build-common-solc
#
# WORKDIR /app
# COPY contracts/payment ./contracts/payment
# RUN forge build --root contracts/payment

FROM docker.io/golang:1.21-bookworm as build-common

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && go install github.com/ethereum/go-ethereum/cmd/abigen@v1.13.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# COPY proto ./proto
# COPY pkg/proto ./pkg/proto
# RUN go generate ./pkg/proto/generate.go
# COPY pkg/abi ./pkg/abi
# COPY --from=build-common-solc /app/contracts/payment contracts/payment
# RUN go generate ./pkg/abi/generate.go

COPY pkg ./pkg

# p2p-helper:

FROM build-common as build-p2p-helper

COPY cmd/ipfs-p2p-helper ./cmd/ipfs-p2p-helper
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/ipfs-p2p-helper ./cmd/ipfs-p2p-helper

FROM docker.io/debian:bookworm-slim as p2p-helper

COPY --from=build-p2p-helper /usr/local/bin/ipfs-p2p-helper /usr/local/bin/ipfs-p2p-helper

ENTRYPOINT ["ipfs-p2p-helper"]

# server:

FROM build-common as build-server

COPY cmd/tpodserver ./cmd/tpodserver
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/tpodserver ./cmd/tpodserver

FROM docker.io/debian:bookworm-slim as server

COPY --from=build-server /usr/local/bin/tpodserver /usr/local/bin/tpodserver

ENTRYPOINT ["tpodserver"]

