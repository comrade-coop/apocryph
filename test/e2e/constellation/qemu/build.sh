#!/bin/sh
set -e
set -v
CHART_PATH="$1"
CONSTELLATION_PATH="../../../../constellation"
# sometimes if the cluster did no terminate correctly, try modifying this
WORKSPACE_PATH="$HOME/.apocryph/constellation"
CURRENT_DIR=$(pwd)

if [ -n "$2" ]; then
  STEP=${2:-1}
  eval "set -v; $(sed -n "/## $STEP: /{:a;n;p;ba};" $0)"
  exit
fi

# Check the number of arguments
if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <CHART_PATH> <STEP(optional)>"
  exit 1
fi

if [ "$1" = "teardown" ]; then
   ( cd $WORKSPACE_PATH; constellation terminate )
   exit 0
fi

## 0: Generate helm template and inject it into constellation base image
helmfile template -f "$CHART_PATH" > "$CONSTELLATION_PATH/image/base/mkosi.skeleton/usr/lib/helmfile-template"

## 1: Build modified image
cd $CURRENT_DIR
cd "$CONSTELLATION_PATH"
bazel run //:tidy

## 1.1: Build the image
cd $CURRENT_DIR
cd "$CONSTELLATION_PATH"
bazel build //image/system:qemu_stable


## 2: create,configure, run workspace
cd "$CURRENT_DIR"
cd "$CONSTELLATION_PATH"
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

if [ -d "$WORKSPACE_PATH" ]; then
    cd "$WORKSPACE_PATH" && constellation terminate
    sudo rm -r "$WORKSPACE_PATH"
fi

mkdir -p "$WORKSPACE_PATH"
cd "$WORKSPACE_PATH"

constellation config generate qemu

# Replace the values in the configuration file
sed -i.bak -e "/^\s*4:/ {n;s/\( *expected: \).*$/\1$PCR4/}" \
           -e "/^\s*9:/ {n;s/\( *expected: \).*$/\1$PCR9/}" \
           -e "/^\s*11:/ {n;s/\( *expected: \).*$/\1$PCR11/}" constellation-conf.yaml
# Replace the control plance count & nodes to 1
sed -i 's/initialCount: [0-9]*/initialCount: 1/' "constellation-conf.yaml"

output=$(constellation version)
version=$(echo "$output" | grep -oP 'Version:\s+\K\S+' | head -n 1)

# Copy the image & rename it to the current constellation version to bypass downloading upstream image
cp $link/constellation.raw "$version.raw" 




