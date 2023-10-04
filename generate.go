//go:generate protoc -I=./proto --go_out=paths=source_relative:./proto ./proto/capacity.proto
//go:generate protoc -I=./proto --go_out=paths=source_relative:./proto ./proto/provision-service.proto
package trustedpods
