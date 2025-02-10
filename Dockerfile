# SPDX-License-Identifier: GPL-3.0

## Base images

FROM docker.io/library/golang:1.23-bookworm@sha256:441f59f8a2104b99320e1f5aaf59a81baabbc36c81f4e792d5715ef09dd29355 AS go-build-base
#RUN go install github.com/ethereum/go-ethereum/cmd/abigen@v1.14.9 # ...

FROM docker.io/debian:bookworm-20250203-slim@sha256:40b107342c492725bc7aacbe93a49945445191ae364184a6d24fedb28172f6f7 AS go-run-base

RUN apt update && apt install -y ca-certificates && rm -rf /var/lib/apt/lists/*

FROM docker.io/library/node:lts-bookworm-slim@sha256:83fdfa2a4de32d7f8d79829ea259bd6a4821f8b2d123204ac467fbe3966450fc AS js-build-base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable && corepack use pnpm@9.15.3

FROM docker.io/library/nginx:latest@sha256:0a399eb16751829e1af26fea27b20c3ec28d7ab1fb72182879dcae1cca21206a AS js-serve-base

## Backend

FROM go-build-base AS build-backend

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download && go mod verify

COPY backend/ ./backend
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/apocryph-s3-backend ./backend/minio-manager


FROM go-run-base AS run-backend

COPY --from=build-backend /usr/local/bin/apocryph-s3-backend /usr/local/bin/apocryph-s3-backend

ENTRYPOINT ["/usr/local/bin/apocryph-s3-backend"]

FROM go-run-base AS run-backend-copy-local

COPY ./bin/apocryph-s3-backend /usr/local/bin/apocryph-s3-backend

ENTRYPOINT ["/usr/local/bin/apocryph-s3-backend"]

## Backend

FROM go-build-base AS build-dns

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download && go mod verify

COPY backend/ ./backend
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/apocryph-s3-dns ./backend/dns-build


FROM go-run-base AS run-dns

COPY --from=build-dns /usr/local/bin/apocryph-s3-dns /usr/local/bin/apocryph-s3-dns

ENTRYPOINT ["/usr/local/bin/apocryph-s3-dns"]

FROM go-run-base AS run-dns-copy-local

COPY ./bin/apocryph-s3-dns /usr/local/bin/apocryph-s3-dns

ENTRYPOINT ["/usr/local/bin/apocryph-s3-dns"]

## Frontend

FROM js-build-base AS build-frontend

COPY ./pnpm-lock.yaml ./pnpm-workspace.yaml ./package.json /app/
COPY ./frontend/package.json /app/frontend/package.json
COPY ./frontend/abi/package.json /app/frontend/abi/package.json
WORKDIR /app

# HACK! https://github.com/nodejs/corepack/issues/612
ENV COREPACK_INTEGRITY_KEYS=0

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

ARG VITE_TOKEN=""
ARG VITE_STORAGE_SYSTEM=""
COPY ./frontend/ /app/frontend/
RUN cd /app/frontend && pnpm run build

FROM js-serve-base AS serve-frontend
COPY --from=build-frontend /app/frontend/dist /usr/share/nginx/html/

FROM js-serve-base AS serve-frontend-copy-local
COPY ./frontend/dist /usr/share/nginx/html/
