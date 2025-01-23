# SPDX-License-Identifier: GPL-3.0

## Base images

FROM docker.io/library/golang:1.23.4-bookworm@sha256:5c3223fcb23efeccf495739c9fd9bbfe76cee51caea90591860395057eab3113 AS go-build-base
#RUN go install github.com/ethereum/go-ethereum/cmd/abigen@v1.14.9 # ...

FROM docker.io/debian:bookworm-20241202-slim@sha256:7a81508cbf1a03e25076ea3ba9f0800321bad64c2a757defa320475dc09d3ec2 AS go-run-base
# version should match golang:1.23.4-bookworm above


FROM docker.io/library/node:22-bookworm-slim@sha256:cc993f948cbd77c7cdfa0a9cc5b05e9ec9554c6ecf8cf98b90a2012156a4b998 AS js-build-base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

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
COPY ./frontend /app/frontend
WORKDIR /app

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN cd /app/frontend && pnpm run build

FROM js-serve-base AS serve-frontend
COPY --from=build-frontend /app/frontend/dist /usr/share/nginx/html/

FROM js-serve-base AS serve-frontend-copy-local
COPY ./frontend/dist /usr/share/nginx/html/
