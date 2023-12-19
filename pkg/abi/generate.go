// SPDX-License-Identifier: GPL-3.0

// Dependencies:
// - github.com/ethereum/go-ethereum/cmd/abigen@v1.13.3
// - forge build ran in ../../contracts

//go:generate abigen --pkg abi --type MockToken --out ./MockToken.abi.go --abi ../../contracts/out/MockToken.sol/MockToken.abi.json
//go:generate abigen --pkg abi --type IERC20 --out ./IERC20.abi.go --abi ../../contracts/out/IERC20.sol/IERC20.abi.json
//go:generate abigen --pkg abi --type Payment --out ./Payment.abi.go --abi ../../contracts/out/Payment.sol/Payment.abi.json
//go:generate abigen --pkg abi --type Registry --out ./Registry.abi.go --abi ../../contracts/out/Registry.sol/Registry.abi.json
package abi
