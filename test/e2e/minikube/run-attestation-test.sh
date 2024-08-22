#!/bin/sh

set -e
echo "USAGE: $0 CERTIFICATE_IDENTITY CERTIFICATE_OIDC_ISSUER"
echo "EXAMPLE: $0 example@email.com https://github.com/login/oauth"
echo "NOTE: The oidc-issuer for Google is https://accounts.google.com, Microsoft is https://login.microsoftonline.com, GitHub is https://github.com/login/oauth, and GitLab is https://gitlab.com."

docker pull nginxdemos/nginx-hello@sha256:2ab1f0bef4461020a1aabee4260a1fe93b03ed69d7f72908acca3a7ec33cb1c0
docker tag docker.io/nginxdemos/nginx-hello:latest ttl.sh/nginx-hello:1h
docker push ttl.sh/nginx-hello:1h

CERTIFICATE_IDENTITY=$1
CERTIFICATE_OIDC_ISSUER=$2

./deploy-pod.sh ../common/manifests/manifest-attestation-nginx.yaml --sign-images --certificate-identity $CERTIFICATE_IDENTITY --certificate-oidc-issuer $CERTIFICATE_OIDC_ISSUER
