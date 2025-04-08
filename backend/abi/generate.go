// SPDX-License-Identifier: GPL-3.0

// Dependencies:
// - github.com/ethereum/go-ethereum/cmd/abigen@v1.14.9
// - forge build ran in ../../../contracts

//go:generate abigen --pkg abi --type IERC20 --out ./IERC20.abi.go --abi ../../contracts/out/IERC20.sol/IERC20.abi.json
//go:generate abigen --pkg abi --type SimplePayment --out ./SimplePayment.abi.go --abi ../../contracts/out/SimplePayment.sol/SimplePayment.abi.json --bin ../../contracts/out/SimplePayment.sol/SimplePayment.bin

package abi
