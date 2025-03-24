#!/bin/sh

genprivkey() {
  head -c 32 /dev/urandom | od -tx -An | tr -d ' \n'
  # Technically, if this generates a number more than fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141 or equal to zero, it won't be a valid private key
  # But chance of that happening is <0.2%, and we are already going to be monitoring deployments anyway
}

mkdir -p /shared_secrets

[ -f /shared_secrets/backend_key.env ] || cat <<EOF >/shared_secrets/backend_key.env
export PRIVATE_KEY="$(genprivkey)"
EOF
