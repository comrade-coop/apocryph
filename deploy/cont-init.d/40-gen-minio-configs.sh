#!/bin/sh

gensecret() {
  head -c 10 /dev/urandom | od -tx -An | tr -d ' '
  # Technically, if this generates a number more than fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141 or equal to zero, it won't be a valid private key
  # But chance of that happening is <0.2%, and we are already going to be monitoring deployments anyway
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
