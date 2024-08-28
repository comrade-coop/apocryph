#!/bin/sh

set -e
echo "USAGE: $0 CERTIFICATE_IDENTITY CERTIFICATE_OIDC_ISSUER STEP"
echo "EXAMPLE: $0 example@email.com https://github.com/login/oauth"
set -v
# NOTE: The oidc-issuer for Google is https://accounts.google.com, Microsoft is https://login.microsoftonline.com, GitHub is https://github.com/login/oauth, and GitLab is https://gitlab.com."

# based on https://stackoverflow.com/a/31269848 / https://bobcopeland.com/blog/2012/10/goto-in-bash/
if [ -n "$3" ]; then
  STEP=${3:-1}
  eval "set -v; $(sed -n "/## $STEP: /{:a;n;p;ba};" $0)"
  exit
fi

CERTIFICATE_IDENTITY=$1
CERTIFICATE_OIDC_ISSUER=$2

IMAGE_DIGEST=$(docker inspect --format='{{index .RepoDigests 0}}' comradecoop/apocryph/tpod-proxy)
echo $IMAGE_DIGEST

## 1: configure the cluster to support image validation
../common/scripts/redeploy-images.sh --set policy.issuer=$CERTIFICATE_OIDC_ISSUER --set policy.subject=$CERTIFICATE_IDENTITY --set policy.image=$IMAGE_DIGEST

## 2: deploy app
docker pull nginxdemos/nginx-hello@sha256:2ab1f0bef4461020a1aabee4260a1fe93b03ed69d7f72908acca3a7ec33cb1c0
docker tag docker.io/nginxdemos/nginx-hello:latest ttl.sh/nginx-hello:1h
docker push ttl.sh/nginx-hello:1h

./deploy-pod.sh ../common/manifests/manifest-attestation-nginx.yaml --certificate-identity $CERTIFICATE_IDENTITY --certificate-oidc-issuer $CERTIFICATE_OIDC_ISSUER

## 3: Get Application info
INGRESS_URL=$(minikube service  -n keda ingress-nginx-controller --url=true -p c1 | head -n 1); echo $INGRESS_URL
MANIFEST_HOST=example.local.info # From manifest-nginx.yaml

while ! curl --connect-timeout 40 -H "Host: $MANIFEST_HOST" $INGRESS_URL --fail-with-body; do sleep 10; done
curl -H "Host: $MANIFEST_HOST" $INGRESS_URL --fail-with-body
