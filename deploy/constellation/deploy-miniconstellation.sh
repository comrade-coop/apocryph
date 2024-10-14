#!/usr/bin/env bash
# SPDX-License-Identifier: GPL-3.0

set -e

which sed >/dev/null
which helmfile >/dev/null; which helm >/dev/null; which kubectl >/dev/null
which constellation >/dev/null || { echo "Install Constellation, https://docs.edgeless.systems/constellation/getting-started/first-steps-local#software-installation-on-ubuntu"; exit 1; }

WORKSPACE_PATH="$HOME/.apocryph/constellation-mini"
SCRIPT_DIR=$(realpath $(dirname "$0"))
REPO_DIR="$SCRIPT_DIR/../../"

# based on https://stackoverflow.com/a/31269848 / https://bobcopeland.com/blog/2012/10/goto-in-bash/
if [ -n "$1" ]; then
  STEP=${1:-1}
  eval "set -v; $(sed -n "/## $STEP: /{:a;n;p;ba};" $0)"
  exit
fi

echo -e "\e[1;32m---"
echo "Note: To skip steps, use '$0 <number>'"
echo "  e.g. to skip ahead to configuring IPFS, run '$0 1.2'"
echo -e "---\e[0m"

set -v

## 1: Start Miniconstellation ##

echo 'CONSTELLATION_PATH='$CONSTELLATION_PATH

mkdir -p $WORKSPACE_PATH
pushd "$WORKSPACE_PATH"

constellation mini up || true

kubectl patch -n kube-system configmap ip-masq-agent --type merge -p '{"data":{"config": "{\"masqLinkLocal\":true,\"nonMasqueradeCIDRs\":[]}"}}'
kubectl rollout restart -n kube-system daemonset cilium
kubectl delete pod -l k8s-app=join-service -n kube-system

popd

## 2: Apply the Helm configuration ##

pushd "$SCRIPT_DIR"

# IMAGE_PREFIX=$(uuidgen)
# docker build -t ttl.sh/$IMAGE_PREFIX-apocryph-server:1h . --target server-copy-local
# docker push ttl.sh/$IMAGE_PREFIX-apocryph-server:1h
## replace ghcr.io/comrade-coop/apocryph/server:master with ttl.sh/$IMAGE_PREFIX-apocryph-server:1h in helmfile.yaml

helmfile sync || kubectl wait --namespace keda --for=condition=available deployment/ingress-nginx-controller && helmfile sync

popd

echo "Run \`$0 example\` to also run a tilt example"
echo "Run \`$0 teardown\` to stop everything"
exit 0

## example: Start a tilt example ##

pushd "$REPO_DIR"

tilt up -- --deploy-stack=False --include ./test/e2e/nginx/Tiltfile --allow-context 'mini-qemu-admin@mini-qemu'

popd

exit 0
## teardown: Stop Miniconstellation ##

pushd "$REPO_DIR"

constellation mini down

popd

