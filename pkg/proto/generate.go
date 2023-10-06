//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/capacity.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/provision-service.proto
package proto
