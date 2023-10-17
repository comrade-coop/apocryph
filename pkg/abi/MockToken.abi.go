// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// MockTokenMetaData contains all meta data concerning the MockToken contract.
var MockTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ClaimTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801562000010575f80fd5b506040518060400160405280600981526020017f4d6f636b546f6b656e00000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f4d4b54000000000000000000000000000000000000000000000000000000000081525081600390816200008e9190620005d5565b508060049081620000a09190620005d5565b505050620000bd30670de0b6b3a7640000620000c360201b60201c565b620007e5565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160362000136575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016200012d9190620006fc565b60405180910390fd5b620001495f83836200014d60201b60201c565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603620001a1578060025f82825462000194919062000744565b9250508190555062000272565b5f805f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050818110156200022d578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040162000224939291906200078f565b60405180910390fd5b8181035f808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603620002bb578060025f828254039250508190555062000305565b805f808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620003649190620007ca565b60405180910390a3505050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680620003ed57607f821691505b602082108103620004035762000402620003a8565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302620004677fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826200042a565b6200047386836200042a565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f620004bd620004b7620004b1846200048b565b62000494565b6200048b565b9050919050565b5f819050919050565b620004d8836200049d565b620004f0620004e782620004c4565b84845462000436565b825550505050565b5f90565b62000506620004f8565b62000513818484620004cd565b505050565b5b818110156200053a576200052e5f82620004fc565b60018101905062000519565b5050565b601f8211156200058957620005538162000409565b6200055e846200041b565b810160208510156200056e578190505b620005866200057d856200041b565b83018262000518565b50505b505050565b5f82821c905092915050565b5f620005ab5f19846008026200058e565b1980831691505092915050565b5f620005c583836200059a565b9150826002028217905092915050565b620005e08262000371565b67ffffffffffffffff811115620005fc57620005fb6200037b565b5b620006088254620003d5565b620006158282856200053e565b5f60209050601f8311600181146200064b575f841562000636578287015190505b620006428582620005b8565b865550620006b1565b601f1984166200065b8662000409565b5f5b8281101562000684578489015182556001820191506020850194506020810190506200065d565b86831015620006a45784890151620006a0601f8916826200059a565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f620006e482620006b9565b9050919050565b620006f681620006d8565b82525050565b5f602082019050620007115f830184620006eb565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f62000750826200048b565b91506200075d836200048b565b925082820190508082111562000778576200077762000717565b5b92915050565b62000789816200048b565b82525050565b5f606082019050620007a45f830186620006eb565b620007b360208301856200077e565b620007c260408301846200077e565b949350505050565b5f602082019050620007df5f8301846200077e565b92915050565b610f3380620007f35f395ff3fe608060405234801561000f575f80fd5b506004361061009c575f3560e01c806370a082311161006457806370a082311461015a57806395d89b411461018a578063a9059cbb146101a8578063c06ed0c2146101d8578063dd62ed3e146101f45761009c565b806306fdde03146100a0578063095ea7b3146100be57806318160ddd146100ee57806323b872dd1461010c578063313ce5671461013c575b5f80fd5b6100a8610224565b6040516100b59190610af3565b60405180910390f35b6100d860048036038101906100d39190610ba4565b6102b4565b6040516100e59190610bfc565b60405180910390f35b6100f66102d6565b6040516101039190610c24565b60405180910390f35b61012660048036038101906101219190610c3d565b6102df565b6040516101339190610bfc565b60405180910390f35b61014461030d565b6040516101519190610ca8565b60405180910390f35b610174600480360381019061016f9190610cc1565b610315565b6040516101819190610c24565b60405180910390f35b61019261035a565b60405161019f9190610af3565b60405180910390f35b6101c260048036038101906101bd9190610ba4565b6103ea565b6040516101cf9190610bfc565b60405180910390f35b6101f260048036038101906101ed9190610cec565b61040c565b005b61020e60048036038101906102099190610d17565b610464565b60405161021b9190610c24565b60405180910390f35b60606003805461023390610d82565b80601f016020809104026020016040519081016040528092919081815260200182805461025f90610d82565b80156102aa5780601f10610281576101008083540402835291602001916102aa565b820191905f5260205f20905b81548152906001019060200180831161028d57829003601f168201915b5050505050905090565b5f806102be6104e6565b90506102cb8185856104ed565b600191505092915050565b5f600254905090565b5f806102e96104e6565b90506102f68582856104ff565b610301858585610591565b60019150509392505050565b5f6012905090565b5f805f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b60606004805461036990610d82565b80601f016020809104026020016040519081016040528092919081815260200182805461039590610d82565b80156103e05780601f106103b7576101008083540402835291602001916103e0565b820191905f5260205f20905b8154815290600101906020018083116103c357829003601f168201915b5050505050905090565b5f806103f46104e6565b9050610401818585610591565b600191505092915050565b8061041630610315565b11610456576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161044d90610e22565b60405180910390fd5b610461303383610591565b50565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b5f33905090565b6104fa8383836001610681565b505050565b5f61050a8484610464565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811461058b578181101561057c578281836040517ffb8f41b200000000000000000000000000000000000000000000000000000000815260040161057393929190610e4f565b60405180910390fd5b61058a84848484035f610681565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610601575f6040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016105f89190610e84565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610671575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016106689190610e84565b60405180910390fd5b61067c838383610850565b505050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16036106f1575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016106e89190610e84565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610761575f6040517f94280d620000000000000000000000000000000000000000000000000000000081526004016107589190610e84565b60405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550801561084a578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516108419190610c24565b60405180910390a35b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036108a0578060025f8282546108949190610eca565b9250508190555061096e565b5f805f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610929578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161092093929190610e4f565b60405180910390fd5b8181035f808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036109b5578060025f82825403925050819055506109ff565b805f808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610a5c9190610c24565b60405180910390a3505050565b5f81519050919050565b5f82825260208201905092915050565b5f5b83811015610aa0578082015181840152602081019050610a85565b5f8484015250505050565b5f601f19601f8301169050919050565b5f610ac582610a69565b610acf8185610a73565b9350610adf818560208601610a83565b610ae881610aab565b840191505092915050565b5f6020820190508181035f830152610b0b8184610abb565b905092915050565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610b4082610b17565b9050919050565b610b5081610b36565b8114610b5a575f80fd5b50565b5f81359050610b6b81610b47565b92915050565b5f819050919050565b610b8381610b71565b8114610b8d575f80fd5b50565b5f81359050610b9e81610b7a565b92915050565b5f8060408385031215610bba57610bb9610b13565b5b5f610bc785828601610b5d565b9250506020610bd885828601610b90565b9150509250929050565b5f8115159050919050565b610bf681610be2565b82525050565b5f602082019050610c0f5f830184610bed565b92915050565b610c1e81610b71565b82525050565b5f602082019050610c375f830184610c15565b92915050565b5f805f60608486031215610c5457610c53610b13565b5b5f610c6186828701610b5d565b9350506020610c7286828701610b5d565b9250506040610c8386828701610b90565b9150509250925092565b5f60ff82169050919050565b610ca281610c8d565b82525050565b5f602082019050610cbb5f830184610c99565b92915050565b5f60208284031215610cd657610cd5610b13565b5b5f610ce384828501610b5d565b91505092915050565b5f60208284031215610d0157610d00610b13565b5b5f610d0e84828501610b90565b91505092915050565b5f8060408385031215610d2d57610d2c610b13565b5b5f610d3a85828601610b5d565b9250506020610d4b85828601610b5d565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680610d9957607f821691505b602082108103610dac57610dab610d55565b5b50919050565b7f436f6e747261637420646f6573206e6f74206861766520456e6f75676820746f5f8201527f6b656e7300000000000000000000000000000000000000000000000000000000602082015250565b5f610e0c602483610a73565b9150610e1782610db2565b604082019050919050565b5f6020820190508181035f830152610e3981610e00565b9050919050565b610e4981610b36565b82525050565b5f606082019050610e625f830186610e40565b610e6f6020830185610c15565b610e7c6040830184610c15565b949350505050565b5f602082019050610e975f830184610e40565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610ed482610b71565b9150610edf83610b71565b9250828201905080821115610ef757610ef6610e9d565b5b9291505056fea26469706673582212208576de40f2c8a0b2cdaef4c98a982c6059f800b86428d5016d6fcd230037892164736f6c63430008150033",
}

// MockTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use MockTokenMetaData.ABI instead.
var MockTokenABI = MockTokenMetaData.ABI

// MockTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockTokenMetaData.Bin instead.
var MockTokenBin = MockTokenMetaData.Bin

// DeployMockToken deploys a new Ethereum contract, binding an instance of MockToken to it.
func DeployMockToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MockToken, error) {
	parsed, err := MockTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockToken{MockTokenCaller: MockTokenCaller{contract: contract}, MockTokenTransactor: MockTokenTransactor{contract: contract}, MockTokenFilterer: MockTokenFilterer{contract: contract}}, nil
}

// MockToken is an auto generated Go binding around an Ethereum contract.
type MockToken struct {
	MockTokenCaller     // Read-only binding to the contract
	MockTokenTransactor // Write-only binding to the contract
	MockTokenFilterer   // Log filterer for contract events
}

// MockTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockTokenSession struct {
	Contract     *MockToken        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockTokenCallerSession struct {
	Contract *MockTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// MockTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockTokenTransactorSession struct {
	Contract     *MockTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MockTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockTokenRaw struct {
	Contract *MockToken // Generic contract binding to access the raw methods on
}

// MockTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockTokenCallerRaw struct {
	Contract *MockTokenCaller // Generic read-only contract binding to access the raw methods on
}

// MockTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockTokenTransactorRaw struct {
	Contract *MockTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockToken creates a new instance of MockToken, bound to a specific deployed contract.
func NewMockToken(address common.Address, backend bind.ContractBackend) (*MockToken, error) {
	contract, err := bindMockToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockToken{MockTokenCaller: MockTokenCaller{contract: contract}, MockTokenTransactor: MockTokenTransactor{contract: contract}, MockTokenFilterer: MockTokenFilterer{contract: contract}}, nil
}

// NewMockTokenCaller creates a new read-only instance of MockToken, bound to a specific deployed contract.
func NewMockTokenCaller(address common.Address, caller bind.ContractCaller) (*MockTokenCaller, error) {
	contract, err := bindMockToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockTokenCaller{contract: contract}, nil
}

// NewMockTokenTransactor creates a new write-only instance of MockToken, bound to a specific deployed contract.
func NewMockTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*MockTokenTransactor, error) {
	contract, err := bindMockToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockTokenTransactor{contract: contract}, nil
}

// NewMockTokenFilterer creates a new log filterer instance of MockToken, bound to a specific deployed contract.
func NewMockTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*MockTokenFilterer, error) {
	contract, err := bindMockToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockTokenFilterer{contract: contract}, nil
}

// bindMockToken binds a generic wrapper to an already deployed contract.
func bindMockToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockToken *MockTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockToken.Contract.MockTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockToken *MockTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockToken.Contract.MockTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockToken *MockTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockToken.Contract.MockTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockToken *MockTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockToken *MockTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockToken *MockTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockToken *MockTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockToken *MockTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MockToken.Contract.Allowance(&_MockToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MockToken *MockTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MockToken.Contract.Allowance(&_MockToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockToken *MockTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockToken *MockTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MockToken.Contract.BalanceOf(&_MockToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MockToken *MockTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MockToken.Contract.BalanceOf(&_MockToken.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockToken *MockTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MockToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockToken *MockTokenSession) Decimals() (uint8, error) {
	return _MockToken.Contract.Decimals(&_MockToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MockToken *MockTokenCallerSession) Decimals() (uint8, error) {
	return _MockToken.Contract.Decimals(&_MockToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockToken *MockTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockToken *MockTokenSession) Name() (string, error) {
	return _MockToken.Contract.Name(&_MockToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MockToken *MockTokenCallerSession) Name() (string, error) {
	return _MockToken.Contract.Name(&_MockToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockToken *MockTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MockToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockToken *MockTokenSession) Symbol() (string, error) {
	return _MockToken.Contract.Symbol(&_MockToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MockToken *MockTokenCallerSession) Symbol() (string, error) {
	return _MockToken.Contract.Symbol(&_MockToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockToken *MockTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockToken *MockTokenSession) TotalSupply() (*big.Int, error) {
	return _MockToken.Contract.TotalSupply(&_MockToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MockToken *MockTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _MockToken.Contract.TotalSupply(&_MockToken.CallOpts)
}

// ClaimTokens is a paid mutator transaction binding the contract method 0xc06ed0c2.
//
// Solidity: function ClaimTokens(uint256 amount) returns()
func (_MockToken *MockTokenTransactor) ClaimTokens(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _MockToken.contract.Transact(opts, "ClaimTokens", amount)
}

// ClaimTokens is a paid mutator transaction binding the contract method 0xc06ed0c2.
//
// Solidity: function ClaimTokens(uint256 amount) returns()
func (_MockToken *MockTokenSession) ClaimTokens(amount *big.Int) (*types.Transaction, error) {
	return _MockToken.Contract.ClaimTokens(&_MockToken.TransactOpts, amount)
}

// ClaimTokens is a paid mutator transaction binding the contract method 0xc06ed0c2.
//
// Solidity: function ClaimTokens(uint256 amount) returns()
func (_MockToken *MockTokenTransactorSession) ClaimTokens(amount *big.Int) (*types.Transaction, error) {
	return _MockToken.Contract.ClaimTokens(&_MockToken.TransactOpts, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockToken *MockTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockToken *MockTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockToken.Contract.Approve(&_MockToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MockToken *MockTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockToken.Contract.Approve(&_MockToken.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockToken *MockTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockToken *MockTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockToken.Contract.Transfer(&_MockToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MockToken *MockTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockToken.Contract.Transfer(&_MockToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockToken *MockTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockToken *MockTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockToken.Contract.TransferFrom(&_MockToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MockToken *MockTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MockToken.Contract.TransferFrom(&_MockToken.TransactOpts, from, to, value)
}

// MockTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the MockToken contract.
type MockTokenApprovalIterator struct {
	Event *MockTokenApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockTokenApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockTokenApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockTokenApproval represents a Approval event raised by the MockToken contract.
type MockTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockToken *MockTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*MockTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MockToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &MockTokenApprovalIterator{contract: _MockToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockToken *MockTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *MockTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MockToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockTokenApproval)
				if err := _MockToken.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MockToken *MockTokenFilterer) ParseApproval(log types.Log) (*MockTokenApproval, error) {
	event := new(MockTokenApproval)
	if err := _MockToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the MockToken contract.
type MockTokenTransferIterator struct {
	Event *MockTokenTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MockTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockTokenTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MockTokenTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MockTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockTokenTransfer represents a Transfer event raised by the MockToken contract.
type MockTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockToken *MockTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MockTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MockTokenTransferIterator{contract: _MockToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockToken *MockTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *MockTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MockToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockTokenTransfer)
				if err := _MockToken.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MockToken *MockTokenFilterer) ParseTransfer(log types.Log) (*MockTokenTransfer, error) {
	event := new(MockTokenTransfer)
	if err := _MockToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
