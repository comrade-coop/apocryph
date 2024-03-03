#!/bin/sh

# Check the number of arguments
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <helm-chart-path> <constellation-path>"
    exit 1
fi

current_dir=$(pwd)

# Generate helm template and inject it into constellation base image
helmfile template -f $1 > $2/image/base/mkosi.skeleton/usr/lib/helmfile-template

# cleanup leftover
cd $current_dir
cd constell-cluster
constellation terminate
cd $current_dir
sudo rm -r constell-cluster 2>/dev/null
mkdir constell-cluster

cd $2
# Get go dependecies
bazel run //:tidy

# Build the image
bazel build //image/system:qemu_stable

# Get the new image measurements
link=$(readlink -f bazel-out/k8-opt/bin/image/system/qemu_qemu-vtpm_stable) 
output=$(bazel run --run_under="sudo -E" //image/measured-boot/cmd $link/constellation.raw $current_dir/custom-measurements.json 2>&1)

# Extract relevant PCR values
PCR4=$(echo "$output" | sed -n '/PCR\[ *4\]/ s/.*: \(.*\)/\1/p')
PCR9=$(echo "$output" | sed -n '/PCR\[ *9\]/ s/.*: \(.*\)/\1/p')
PCR11=$(echo "$output" | sed -n '/PCR\[ *11\]/ s/.*: \(.*\)/\1/p')

echo "PCR4: $PCR4"
echo "PCR9: $PCR9"
echo "PCR11: $PCR11"

cd "$current_dir"
cd constell-cluster

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

# Start the constellation cluster
constellation apply -y
