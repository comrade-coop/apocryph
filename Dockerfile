# SPDX-License-Identifier: GPL-3.0

## Base images

FROM docker.io/library/golang:1.23.4-bookworm@sha256:5c3223fcb23efeccf495739c9fd9bbfe76cee51caea90591860395057eab3113 as go-build-base
#RUN go install github.com/ethereum/go-ethereum/cmd/abigen@v1.14.9 # ...

FROM docker.io/debian:bookworm-20241202-slim@sha256:7a81508cbf1a03e25076ea3ba9f0800321bad64c2a757defa320475dc09d3ec2 as go-run-base
# version should match golang:1.23.4-bookworm above


FROM docker.io/library/node:22-bookworm-slim@sha256:cc993f948cbd77c7cdfa0a9cc5b05e9ec9554c6ecf8cf98b90a2012156a4b998 as js-base
#RUN go install github.com/ethereum/go-ethereum/cmd/abigen@v1.14.9 # ...

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

# Same image used at runtime and at buildtime
FROM js-base as js-build-base
FROM js-base as js-run-base

## Backend

FROM go-build-base AS build-backend

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download && go mod verify

COPY backend/ ./backend
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/apocryph-s3-backend ./backend/minio-manager


FROM go-run-base AS run-backend

COPY --from=build /usr/local/bin/apocryph-s3-backend /usr/local/bin/apocryph-s3-backend

ENTRYPOINT ["/usr/local/bin/apocryph-s3-backend"]

FROM go-run-base AS run-backend-copy-local

COPY ./bin/apocryph-s3-backend /usr/local/bin/apocryph-s3-backend

ENTRYPOINT ["/usr/local/bin/apocryph-s3-backend"]

## Backend

FROM go-build-base AS build-backend

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download && go mod verify

COPY backend/ ./backend
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o /usr/local/bin/apocryph-s3-dns ./backend/dns-build


FROM go-run-base AS run-dns

COPY --from=build /usr/local/bin/apocryph-s3-dns /usr/local/bin/apocryph-s3-dns

ENTRYPOINT ["/usr/local/bin/apocryph-s3-dns"]

FROM go-run-base AS run-dns-copy-local

COPY ./bin/apocryph-s3-dns /usr/local/bin/apocryph-s3-dns

ENTRYPOINT ["/usr/local/bin/apocryph-s3-dns"]

## Frontend

FROM js-build-base AS build-frontend

COPY . /app
WORKDIR /app

FROM build-frontend AS build-frontend-prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM build-frontend AS build-frontend-build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM js-run-base AS frontend
COPY --from=prod-deps /app/node_modules /app/node_modules
COPY --from=build /app/dist /app/dist
EXPOSE 5173
ENTRYPOINT [ "pnpm", "start" ]

WORKDIR /app
