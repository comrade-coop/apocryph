# IPFS P2P connection test

This integration test assesses the functionality of IPFS peer-to-peer communication between a client and a provider, showcasing that the client sends gRPC requests to a different port within a container using its own IPFS node, while the provider listens on a completely different port to prove that the P2P connection is routed through the IPFS nodes.

## Client

- building:
  - pulls a hello-world oci image
  - save it in a directory (alongside a manifest.yaml file) called client-pod
  - build the publisher-client image (also builds client.go)
  - cleanup (remove client-pod & client)
- executing:
  - forward the client grpc requests to localhost:5000 to the provider IPFS id for provision-pod protocol
  - start the IPFS daemon
  - Add the client-pod to the container IPFS node
  - prints the pod package content

## Provider

- Run the grpc server
- Create an IPFS p2p service that routes all the requests for provision-pod protocol to the grpc server in localhost:6000
- when receiving a SendPod rpc call, it pulls the pod package from container IPFS node
- prints the pod package content

Run the following commands:

```
go generate ./prg/proto
cd test/integration/p2p/
./build.sh
./run-test.sh
```

expected output:

```
2023/10/08 16:40:49 PROVIDER: server listening at [::]:6000
CLIENT:waiting ipfs daemon to start
CLIENT: pod package cid: /ipfs/QmdzwMXxbSTqG9xizYCKfnugJUWx7WkcCfuDZZfM1jH6Nd
CLIENT: Pod Package: total 36
drwxrwxr-x    2 root     root          4096 Oct  8 15:40 .
drwxr-xr-x    1 root     root          4096 Oct  8 15:40 ..
-rw-------    1 root     root         24064 Oct  8 15:40 hello-world.tar
-rw-rw-r--    1 root     root          1396 Oct  8 15:40 manifest.yaml

PROVIDER: pod package cid: /ipfs/QmdzwMXxbSTqG9xizYCKfnugJUWx7WkcCfuDZZfM1jH6Nd
PROVIDER: Pod Package: total 44
drwxrwxr-x  2 ezio ezio  4096 Oct  8 16:41 .
drwxrwxrwt 49 root root 12288 Oct  8 16:41 ..
-rw-rw-r--  1 ezio ezio 24064 Oct  8 16:41 hello-world.tar
-rw-rw-r--  1 ezio ezio  1396 Oct  8 16:41 manifest.yaml

2023/10/08 15:41:02 CLIENT: Pod Endpoint: http://provider.com/mypod
```
