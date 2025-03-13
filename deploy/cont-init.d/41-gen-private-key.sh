#!/bin/sh

genprivkey() {
  head -c 32 /dev/urandom | od -tx -An | tr -d ' \n'
}

mkdir -p /shared_secrets

[ -f /shared_secrets/backend_key.env ] || cat <<EOF >/shared_secrets/backend_key.env
export PRIVATE_KEY="$(genprivkey)"
EOF
