// SPDX-License-Identifier: GPL-3.0

// Dependencies:
// - github.com/ethereum/go-ethereum/cmd/abigen@v1.14.9
// - forge build ran in ../../../contracts

//go:generate abigen --pkg abi --type IERC20 --out ./IERC20.abi.go --abi ../../contracts/out/IERC20.sol/IERC20.abi.json

// disabled for now: abigen --pkg abi --type PaymentV2 --out ./PaymentV2.abi.go --abi ../../contracts/out/PaymentV2.sol/PaymentV2.abi.json

package abi
