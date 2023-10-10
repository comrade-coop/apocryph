#!/bin/sh
cd "$(dirname "$0")" || exit 1

if ! which ipfs >/dev/null 2>&1 || ! which docker >/dev/null 2>&1; then
    echo "Requires: Docker, IPFS"
	exit 1
fi
if ! pgrep -x "ipfs" >/dev/null; then
    echo "IPFS is not running. Starting IPFS..."
    ipfs daemon > /dev/null 2>&1 &
fi
docker pull nginx:1.25.2
docker pull hello-world:linux
rm -r /tmp/package 2> /dev/null
go run main.go package
echo "\n"
echo "CLIENT: Pod Package:"
ls -al package
echo "PROVIDER: Decrypted pod package from IPFS"
ls -al /tmp/package
if diff -q "/tmp/package/manifest.yaml" "./package-files/manifest.yaml" > /dev/null; then
		echo "\n File decrypted succefully and matches manifest.yaml ✅"
else
		echo "File decrypted does not match original manifest.yaml ❌"
fi

rm -r package

