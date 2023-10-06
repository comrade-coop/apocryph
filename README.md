# Trusted Pods

Trusted Pods is a decentralized compute marketplace where developers can run container pods securely and confidentially through small and medium cloud providers. It is tailored towards providing a one-stop service for finding infrastructure to run pods in a cost efficient, serverless manner.

## Building

To build the Go projects in `./cmd/`, use the following command:

```bash
go generate ./pkg/proto && go build -o bin ./cmd/*/
```

Alternatively, you can run the commands directly with `go run`:

```bash
go generate ./pkg/proto

go run ./cmd/trustedpods/ # ..
```

## Testing

To run the minikube integration test, run the following command:

```bash
test/integration/minikube/run-test.sh
```

The command will proceed to start a minikube cluster, deploy all necessary prerequisites into the cluster, apply a [manifest file](spec/MANIFEST.md) into the cluster, and finally query the started pod over HTTP. It should display the curl command used to query the pod, and you should be able to run it yourself after the script is finished.

When you are done playing around with the test, simply run the following command to delete the minikube cluster:

```bash
test/integration/minikube/run-test.sh teardown
```

## Architecture

See [`ARCHITECTURE.md`](spec/ARCHITECTURE.md).
