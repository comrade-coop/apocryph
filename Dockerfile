# SPDX-License-Identifier: GPL-3.0

# HACK: TARGETARCH_x86 should be the same as TARGETARCH, except with amd64 being x86_64
ARG TARGETARCH_x86="x86_64"

# RUN_IMAGE is the image used for running the whole Minio+Nginx+Go Backend stack.
ARG RUN_IMAGE=docker.io/debian:bookworm-20250203-slim@sha256:40b107342c492725bc7aacbe93a49945445191ae364184a6d24fedb28172f6f7

# S6_MODE is always "download" - specifying whether S6 be downloaded from https://github.com/just-containers/s6-overlay/releases
ARG S6_MODE=download
# S6_OVERLAY_VERSION is the version of S6 to use
ARG S6_OVERLAY_VERSION=3.2.0.2

# BACKEND_MODE is one of "build" or "copy" - specifying whether the Go backend should be built by the Docker build or copied from the ./bin/ folder
ARG BACKEND_MODE=build

# FRONTEND_MODE is one of "build", "copy", or "none" - specifying whether the JS frontend should be built by the Docker build, copied from the ./frontend/dist/ folder, or omitted
ARG FRONTEND_MODE=build

# MINIO_MODE is one of "build" or "download" - specifying whether Minio (and Mc) be built by the Docker build, or downloaded from https://dl.min.io
ARG MINIO_MODE=download
# MINIO_RELEASE is the version of Minio to download (from e.g. https://dl.min.io/server/minio/release/linux-amd64/ *.deb)
ARG MINIO_RELEASE=20250218162555.0.0
# MINIO_RELEASE is the version of the Minio client to download (from e.g. https://dl.min.io/client/mc/release/linux-amd64/ *.deb)
ARG MINIO_MC_RELEASE=20250215103616.0.0
# MINIO_GOVER is the version of Minio to use when building (12-letter commit hash or latest)
ARG MINIO_GOVER=latest # 90f5e1e5f62c
# MINIO_MC_GOVER is the version of the Minio client to use when building (12-letter commit hash or latest)
ARG MINIO_MC_GOVER=latest # 383560b1c3d6

# PROMETHEUS_MODE is one of "build" or "download" - specifying whether Prometheus should be built by the Docker build, or downloaded from https://github.com/prometheus/prometheus/releases/
ARG PROMETHEUS_MODE=download
# PROMETHEUS_VERSION is the version of Prometheus to download from https://github.com/prometheus/prometheus/releases/
ARG PROMETHEUS_VERSION=3.2.0

# GO_BUILD_IMAGE is the image used for building the Go backend. _Must_ have the same glibc as RUN_IMAGE.
ARG GO_BUILD_IMAGE=docker.io/library/golang:1.23-bookworm@sha256:441f59f8a2104b99320e1f5aaf59a81baabbc36c81f4e792d5715ef09dd29355

# JS_BUILD_IMAGE is the image used for building the JS/React/Vite frontend.
ARG JS_BUILD_IMAGE=docker.io/library/node:lts-bookworm-slim@sha256:83fdfa2a4de32d7f8d79829ea259bd6a4821f8b2d123204ac467fbe3966450fc
# NGINX_IMAGE is the image used for serving static files for the serve-frontend target, and is otherwise unused.
ARG NGINX_IMAGE=docker.io/library/nginx:latest@sha256:0a399eb16751829e1af26fea27b20c3ec28d7ab1fb72182879dcae1cca21206a
# PNPM_VERSION is the current version of pnpm used. Should be kept in sync with package.json.
ARG PNPM_VERSION=9.15.3
# VITE_TOKEN is the on-chain address of the token used.
ARG VITE_TOKEN=""
# VITE_STORAGE_SYSTEM is the on-chain address corresponding to the backend's private key.
ARG VITE_STORAGE_SYSTEM=""
# VITE_CHAIN_CONFIG is an optional viem chain object, as described in https://viem.sh/docs/chains/introduction
ARG VITE_CHAIN_CONFIG=""
# VITE_GLOBAL_HOST is the address of the hosted aapp.
ARG VITE_GLOBAL_HOST=s3-aapp.localhost
# VITE_GLOBAL_HOST_CONSOLE is the address of the minio console.
ARG VITE_GLOBAL_HOST_CONSOLE=console-s3-aapp.localhost
# VITE_GLOBAL_HOST_APP is the address of the hosted aapp.
ARG VITE_GLOBAL_HOST_APP=console-aapp.localhost
# VITE_PUBLIC_ATTESTATION_URL is the address for the attestation link.
ARG VITE_PUBLIC_ATTESTATION_URL=

# SERF_VERSION is currently unused
ARG SERF_VERSION=0.10.1
# SERF_MODE is currently unused
ARG SERF_MODE=build

## Builder

FROM $GO_BUILD_IMAGE as builder-go

FROM $JS_BUILD_IMAGE as builder-js

ARG PNPM_VERSION

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable && corepack use pnpm@$PNPM_VERSION

## Backend

FROM builder-go AS backend-build

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download && go mod verify

COPY backend/ ./backend
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /bin/ ./backend/dns-build ./backend/minio-manager


FROM scratch AS backend-copy

COPY ./bin/dns-build ./bin/minio-manager /bin/


FROM backend-$BACKEND_MODE AS backend

FROM $RUN_IMAGE AS run-backend

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=backend /bin/minio-manager /usr/local/bin/apocryph-s3-backend
COPY --from=backend /bin/dns-build /usr/local/bin/apocryph-s3-dns

CMD ["/usr/local/bin/apocryph-s3-backend"]


## Frontend

FROM builder-js AS frontend-build


COPY ./pnpm-lock.yaml ./pnpm-workspace.yaml ./package.json /app/
COPY ./frontend/package.json /app/frontend/package.json
COPY ./frontend/abi/package.json /app/frontend/abi/package.json

# HACK! https://github.com/nodejs/corepack/issues/612
ARG COREPACK_INTEGRITY_KEYS=0

RUN --mount=type=cache,id=pnpm,target=/pnpm/store cd /app && pnpm install --frozen-lockfile

ARG VITE_TOKEN
ARG VITE_STORAGE_SYSTEM
ARG VITE_CHAIN_CONFIG
ARG VITE_GLOBAL_HOST
ARG VITE_GLOBAL_HOST_CONSOLE
ARG VITE_GLOBAL_HOST_APP
ARG VITE_PUBLIC_ATTESTATION_URL

COPY ./frontend/ /app/frontend/
RUN cd /app/frontend && pnpm run build


FROM scratch AS frontend-copy
COPY ./frontend/dist /app/frontend/dist


FROM scratch AS frontend-none
# Do nothing


FROM frontend-$FRONTEND_MODE AS frontend

FROM $NGINX_IMAGE AS serve-frontend
COPY --from=frontend /app/frontend/dist /usr/share/nginx/html/


## Minio

FROM builder-go AS minio-download

ARG TARGETARCH
ARG MINIO_RELEASE
ARG MINIO_MC_RELEASE

RUN curl \
    https://dl.min.io/server/minio/release/linux-${TARGETARCH}/archive/minio_${MINIO_RELEASE}_${TARGETARCH}.deb -o minio.deb \
    https://dl.min.io/client/mc/release/linux-${TARGETARCH}/archive/mcli_${MINIO_MC_RELEASE}_${TARGETARCH}.deb -o mcli.deb && \
    dpkg -i minio.deb && dpkg -i mcli.deb && \
    rm minio.deb mcli.deb
RUN mv /usr/local/bin/mcli /usr/local/bin/mc


FROM builder-go AS minio-build

ARG MINIO_GOVER
ARG MINIO_MC_GOVER

ENV GOBIN=/usr/bin/local
RUN --mount=type=cache,target=/root/.cache/go-build go install -v github.com/minio/minio@${MINIO_GOVER}
RUN --mount=type=cache,target=/root/.cache/go-build go install -v github.com/minio/mc@${MINIO_MC_GOVER}


FROM minio-$MINIO_MODE AS minio

## Prometheus

FROM builder-go AS prometheus-download

ARG TARGETARCH
ARG PROMETHEUS_VERSION

ADD https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-${TARGETARCH}.tar.gz /prometheus.tar.gz
RUN mkdir /prometheus && tar -xzf /prometheus.tar.gz --strip-components 1 -C /prometheus/


FROM builder-go AS prometheus-build

ARG PROMETHEUS_VERSION

RUN git clone https://github.com/prometheus/prometheus/ -b v${PROMETHEUS_VERSION} --depth 1 /prometheus
RUN --mount=type=cache,target=/root/.cache/go-build cd /prometheus && go build -O . ./cmd/prometheus/


FROM prometheus-$PROMETHEUS_MODE AS prometheus

## Serf

FROM builder-go AS serf-build

ARG SERF_VERSION

ENV GOBIN=/usr/bin/local
RUN go install github.com/hashicorp/serf/cmd/serf@v${SERF_VERSION}


FROM serf-$SERF_MODE AS serf

## s6

FROM $RUN_IMAGE AS s6-download

ARG S6_OVERLAY_VERSION
ARG TARGETARCH_x86

ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-${TARGETARCH_x86}.tar.xz /

FROM s6-${S6_MODE} AS s6

## Uberimage

FROM $RUN_IMAGE AS run-all-singlenode

ARG TARGETARCH_x86

RUN apt-get update && \
    apt-get install -y ca-certificates nginx xz-utils && \
    rm -rf /var/lib/apt/lists/*

RUN --mount=type=bind,from=s6,target=/s6 tar -C / -Jxpf /s6/s6-overlay-noarch.tar.xz && tar -C / -Jxpf /s6/s6-overlay-${TARGETARCH_x86}.tar.xz

COPY ./deploy/ /etc/
    
COPY --from=frontend /app/frontend/dist /usr/share/nginx/html/
COPY --from=backend /bin/minio-manager /usr/local/bin/apocryph-s3-backend
COPY --from=prometheus /prometheus/prometheus /usr/local/bin/
COPY --from=minio /usr/local/bin/minio /usr/local/bin/mc /usr/local/bin/

ENTRYPOINT ["/init"]
CMD ["sleep", "infinity"]

ARG VITE_GLOBAL_HOST
ARG VITE_GLOBAL_HOST_CONSOLE
ARG VITE_GLOBAL_HOST_APP
ARG VITE_BACKEND_ETH_WITHDRAW
ENV GLOBAL_HOST=$VITE_GLOBAL_HOST
ENV GLOBAL_HOST_CONSOLE=$VITE_GLOBAL_HOST_CONSOLE
ENV GLOBAL_HOST_APP=$VITE_GLOBAL_HOST_APP
ENV BACKEND_ETH_WITHDRAW=$VITE_BACKEND_ETH_WITHDRAW
ENV BACKEND_ETH_TOKEN=$VITE_TOKEN
ENV BACKEND_ETH_RPC=https://sepolia.base.org
ENV BACKEND_ETH_CHAIN_ID=84532
ENV BACKEND_EXTERNAL_URL=http://example.invalid
ENV BACKEND_REPLICATE_SITES=""
VOLUME /data
VOLUME /shared_secrets

## Default output

FROM run-all-singlenode
