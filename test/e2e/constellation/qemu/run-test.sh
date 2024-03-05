#!/bin/sh
set -e

chart_path="$1"
constellation_path="$2"
workspace_path="$3"
current_dir=$(pwd)

if [ -n "$4" ]; then
  STEP=${4:-1}
  eval "set -v; $(sed -n "/## $STEP: /{:a;n;p;ba};" $0)"
  exit
fi

# Check the number of arguments
if [ "$#" -lt 3 ]; then
    echo "Usage: $0 <helm-chart-path> <constellation-path> <workspace-path>"
    exit 1
fi

## 0: Generate helm template and inject it into constellation base image
helmfile template -f "$chart_path" > "$constellation_path/image/base/mkosi.skeleton/usr/lib/helmfile-template"

## 1: Build modified image
cd $current_dir
cd "$constellation_path"
bazel run //:tidy

## 1.1: Build the image
cd $current_dir
cd "$constellation_path"
bazel build //image/system:qemu_stable

## 1.2: Get the new image measurements
link=$(readlink -f bazel-out/k8-opt/bin/image/system/qemu_qemu-vtpm_stable) 
output=$(bazel run --run_under="sudo -E" //image/measured-boot/cmd $link/constellation.raw "$workspace_path/custom-measurements.json" 2>&1)

## 1.3: Extract relevant PCR values
PCR4=$(echo "$output" | sed -n '/PCR\[ *4\]/ s/.*: \(.*\)/\1/p')
PCR9=$(echo "$output" | sed -n '/PCR\[ *9\]/ s/.*: \(.*\)/\1/p')
PCR11=$(echo "$output" | sed -n '/PCR\[ *11\]/ s/.*: \(.*\)/\1/p')

echo "PCR4:  $PCR4"
echo "PCR9:  $PCR9"
echo "PCR11: $PCR11"

## 2: terminate clutser
cd "$workspace_path"
constellation terminate
sudo rm -r "$workspace_path" 2>/dev/null
mkdir "$workspace_path"


## 3: create,configure, run workspace
cd "$current_dir"
cd "$workspace_path"

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

## 3.1 Start the constellation cluster
constellation apply -y
