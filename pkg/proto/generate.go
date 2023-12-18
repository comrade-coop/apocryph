// SPDX-License-Identifier: GPL-3.0

// Dependencies:
// - protoc
// - google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
// - connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.12.0

//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/pod.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/pricing.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/deployment.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. --connect-go_out=paths=source_relative:. ../../proto/provision-pod.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/provisioning-capacity.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/registry.proto
package proto
