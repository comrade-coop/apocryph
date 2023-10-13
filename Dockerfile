FROM docker.io/golang:1.21-bookworm as build

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY proto ./proto
COPY pkg/proto ./pkg/proto
RUN go generate ./pkg/proto

COPY pkg ./pkg
COPY cmd/tpodserver ./cmd/tpodserver
RUN go build -v -o /usr/local/bin/tpodserver ./cmd/tpodserver

FROM docker.io/debian:bookworm-slim as final

COPY --from=build /usr/local/bin/tpodserver /usr/local/bin/tpodserver

ENTRYPOINT ["tpodserver"]
