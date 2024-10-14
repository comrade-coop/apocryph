#!/bin/bash
set -e
set -v

which constellation >/dev/null; which kubectl >/dev/null; which tilt >/dev/null

WORKSPACE_PATH=$(echo "$HOME/.apocryph/constellation-"*)
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
echo -e "---\e[0m"

## 0: Build custom OS image & run the cluster ##

"$SCRIPT_DIR/build-qemu.sh"

## 1: Start the Constellation cluster ##

pushd "$WORKSPACE_PATH"
constellation terminate || true
constellation apply -y
popd

## 1.1: Wait for the setup ##

export KUBECONFIG="$WORKSPACE_PATH/constellation-admin.conf"

kubectl wait --namespace keda --for=condition=available deployment/ingress-nginx-controller
kubectl wait --namespace prometheus --for=condition=available deployment/prometheus-kube-state-metrics
kubectl wait --namespace prometheus --for=condition=available deployment/prometheus-prometheus-pushgateway
kubectl wait --namespace prometheus --for=condition=available deployment/prometheus-server
kubectl wait --namespace ipfs --for=condition=available StatefulSet/ipfs
kubectl wait --namespace trustedpods --for=condition=available deployment/tpodserver

echo "Run \`$0 example\` to also run a tilt example"
echo "Run \`$0 teardown\` to stop everything"

exit 0
## example: start tilt

pushd "$REPO_DIR"

tilt up -- --deploy-stack=False --include ./test/e2e/nginx/Tiltfile

popd

exit 0
## teardown: Stop the Constellation cluster ##

pushd "$REPO_DIR"

constellation terminate

popd

