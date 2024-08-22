#!/bin/sh
set -e 
echo "USAGE: $0 CERTIFICATE_IDENTITY CERTIFICATE_OIDC_ISSUER"
echo "EXAMPLE: $0 example@email.com https://github.com/login/oauth"
echo "NOTE: The oidc-issuer for Google is https://accounts.google.com, Microsoft is https://login.microsoftonline.com, GitHub is https://github.com/login/oauth, and GitLab is https://gitlab.com."

trap 'kill $(jobs -p) &>/dev/null' EXIT
set -v

cd "$(dirname "$0")"

sudo chmod o+rw /run/containerd/containerd.sock

ipfs shutdown || true
ipfs daemon >/dev/null &
sleep 2



CERTIFICATE_IDENTITY=$1
CERTIFICATE_OIDC_ISSUER=$2

# docker tag hello-world ttl.sh/hello-world:1h
# docker push ttl.sh/hello-world:1h
#
# go run ../../../cmd/trustedpods pod upload ../../e2e/common/manifests/manifest-attestation.yaml \
#   --sign-images --certificate-identity $CERTIFICATE_IDENTITY --certificate-oidc-issuer $CERTIFICATE_OIDC_ISSUER
#
# go run ../../../cmd/trustedpods pod verify ../../e2e/common/manifests/manifest-attestation.yaml

go run ../../../cmd/trustedpods verify  ttl.sh/hello-world@sha256:d37ada95d47ad12224c205a938129df7a3e52345828b4fa27b03a98825d1e2e7 \
  --certificate-identity $CERTIFICATE_IDENTITY --certificate-oidc-issuer $CERTIFICATE_OIDC_ISSUER \
  --signature ~/.apocryph/signatures/ttl.sh_hello-world@sha256:d37ada95d47ad12224c205a938129df7a3e52345828b4fa27b03a98825d1e2e7.sig
