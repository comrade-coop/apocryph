#!/bin/sh

gensecret() {
  head -c 10 /dev/random | od -tx -An | tr -d ' '
}

mkdir -p /secrets/; # TODO: Secure better?

cat <<EOF >/secrets/backend.env
export ACCESS_KEY="$(gensecret)"
export SECRET_KEY="$(gensecret)"
EOF

cat <<EOF >/secrets/minio.env
export MINIO_ROOT_USER="$(gensecret)"
export MINIO_ROOT_PASSWORD="$(gensecret)"
EOF
