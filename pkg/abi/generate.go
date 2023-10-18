// Dependencies:
// - github.com/ethereum/go-ethereum/cmd/abigen@v1.13.3
// - forge build ran in ../../contracts/payment

//go:generate abigen --pkg abi --type MockToken --out ./MockToken.abi.go --abi ../../contracts/payment/out/MockToken.sol/MockToken.abi.json --bin ../../contracts/payment/out/MockToken.sol/MockToken.bin
//go:generate abigen --pkg abi --type Payment --out ./Payment.abi.go --abi ../../contracts/payment/out/Payment.sol/Payment.abi.json --bin ../../contracts/payment/out/Payment.sol/Payment.bin
package abi
