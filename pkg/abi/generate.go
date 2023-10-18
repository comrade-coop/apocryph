// Dependencies:
// - github.com/ethereum/go-ethereum/cmd/abigen@v1.13.3
// - forge build ran in ../../contracts

//go:generate abigen --pkg abi --type MockToken --out ./MockToken.abi.go --abi ../../contracts/out/MockToken.sol/MockToken.abi.json --bin ../../contracts/out/MockToken.sol/MockToken.bin
//go:generate abigen --pkg abi --type Payment --out ./Payment.abi.go --abi ../../contracts/out/Payment.sol/Payment.abi.json --bin ../../contracts/out/Payment.sol/Payment.bin
package abi
