#!/bin/sh
set -e

which helmfile >/dev/null
which basel >/dev/null || { echo "Install Bazel, ideally through Bazelisk, https://bazel.build/install/bazelisk"; exit 1; }
which nix >/dev/null || { echo "Install Nix, https://nixos.org/download/"; exit 1; }

set -v
SUFFIX=$RANDOM
WORKSPACE_PATH="$HOME/.apocryph/constellation-$SUFFIX"
echo "WORKSPACE_PATH:: $WORKSPACE_PATH"

SCRIPT_DIR=$(realpath $(dirname "$0"))
HELMFILE_PATH="$SCRIPT_DIR/helmfile.yaml"
CONSTELLATION_SRC="$SCRIPT_DIR/constellation-src"

if [ -n "$1" ]; then
  STEP=${1:-1}
  eval "set -v; $(sed -n "/## $STEP: /{:a;n;p;ba};" $0)"
  exit
fi

if [ "$1" = "teardown" ]; then
   sudo rm -r "$WORKSPACE_PATH"
   exit 0
fi

## 0: Generate helm template and inject it into constellation base image
helmfile template -f "$HELMFILE_PATH" --kube-version 1.30.0 > "$CONSTELLATION_SRC/image/base/mkosi.skeleton/usr/lib/helmfile-template"

## 1: Build modified image
pushd "$CONSTELLATION_SRC"
bazel run //:tidy
popd

## 1.1: Build the image
pushd "$CONSTELLATION_SRC"
bazel build //image/system:qemu_stable
popd

## 2: create & configure constellation workspace
set -e
set -v

pushd "$CONSTELLATION_SRC"
# Get the new image measurements
link=$(readlink -f bazel-out/k8-opt/bin/image/system/qemu_qemu-vtpm_stable) 
output=$(bazel run --run_under="sudo -E" //image/measured-boot/cmd $link/constellation.raw "/tmp/custom-measurements.json" 2>&1) # second arg needed

# Extract relevant PCR values
PCR4=$(echo "$output" | sed -n '/PCR\[ *4\]/ s/.*: \(.*\)/\1/p')
PCR9=$(echo "$output" | sed -n '/PCR\[ *9\]/ s/.*: \(.*\)/\1/p')
PCR11=$(echo "$output" | sed -n '/PCR\[ *11\]/ s/.*: \(.*\)/\1/p')

echo "PCR4:  $PCR4"
echo "PCR9:  $PCR9"
echo "PCR11: $PCR11"

popd

mkdir -p "$WORKSPACE_PATH"
pushd "$WORKSPACE_PATH"

constellation config generate qemu

# Replace the values in the configuration file
sed -i.bak -e "/^\s*4:/ {n;s/\( *expected: \).*$/\1$PCR4/}" \
           -e "/^\s*9:/ {n;s/\( *expected: \).*$/\1$PCR9/}" \
           -e "/^\s*11:/ {n;s/\( *expected: \).*$/\1$PCR11/}" constellation-conf.yaml
# Replace the control plance count & nodes to 1
sed -i 's/initialCount: [0-9]*/initialCount: 1/' "constellation-conf.yaml"
sed -i 's/vcpus: [0-9]*/vcpus: 4/' "constellation-conf.yaml"
sed -i 's/memory: [0-9]*/memory: 4096/' "constellation-conf.yaml"

output=$(constellation version)
version=$(echo "$output" | grep -oP 'Version:\s+\K\S+' | head -n 1)

# Copy the image & rename it to the current constellation version to bypass downloading upstream image
cp $link/constellation.raw "$version.raw" 
popd

export KUBECONFIG="$WORKSPACE_PATH/constellation-admin.conf"




