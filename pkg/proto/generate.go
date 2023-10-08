//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/capacity.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/provision-service.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/sample.proto
//go:generate protoc -I=../../proto --go_out=../../proto --go_opt=paths=source_relative --go-grpc_out=../../proto --go-grpc_opt=paths=source_relative ../../proto/sample.proto
package proto
