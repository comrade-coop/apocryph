#!/usr/bin/env bash
cd "$(dirname "$0")"
anvil > anvil_output.txt &

sleep 2

private_key1=$(awk '/Private Keys/ {flag=1; next} flag && /^\(0\)/ {print $2; exit}' anvil_output.txt)
private_key2=$(awk '/Private Keys/ {flag=1; next} flag && /^\(1\)/ {print $2; exit}' anvil_output.txt)
echo "client private key: $private_key1"
echo "provider private key: $private_key2"
go run . ./keystore $private_key1 $private_key2

rm -r ./keystore > /dev/null
pkill anvil
rm anvil_output.txt
