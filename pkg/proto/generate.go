//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/pod.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. ../../proto/pricing.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ../../proto/provision-pod.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ../../proto/provisioning-capacity.proto
//go:generate protoc -I=../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ../../proto/sample.proto
//go:generate abigen --abi ../../contracts/payment/build/MockToken.abi --pkg abi --type MockToken --out  ../abi/MockToken.abi.go --bin ../../contracts/payment/build/MockToken.bin
//go:generate abigen --abi ../../contracts/payment/build/Payment.abi --pkg abi --type Payment --out  ../abi/Payment.abi.go --bin ../../contracts/payment/build/Payment.bin
package proto
