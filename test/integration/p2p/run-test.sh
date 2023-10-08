rm -r /tmp/pod-package 2> /dev/null

if ! which ipfs >/dev/null 2>&1 || ! which docker >/dev/null 2>&1; then
		echo "make sure the following packages are installed and available in your path: ipfs, docker, jq"
		exit 1
fi
pkill provider
ipfs config --json Experimental.Libp2pStreamMounting true > /dev/null 2>&1
# route all ipfs p2p connections of the provios-pod protocol to the grpc server
# ipfs p2p listen /x/trusted-pods/provision-pod/0.0.1 /ip4/127.0.0.1/tcp/6000  > /dev/null 2>&1
go run provider.go &

ipfs_id_output=$(ipfs id)
provider_id=$(echo $ipfs_id_output | jq -r .ID)
docker run -d --rm -e PROVIDER_ID=$provider_id --name publisher-client publisher > /dev/null; docker logs publisher-client -f

