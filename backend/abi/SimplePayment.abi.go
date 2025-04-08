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

// SimplePaymentMetaData contains all meta data concerning the SimplePayment contract.
var SimplePaymentMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalPaid\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"version\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"payerAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"version\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"withdrawAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561000f575f80fd5b50604051610b5e380380610b5e83398181016040528101906100319190610228565b335f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100999190610262565b60405180910390fd5b6100b1816100f860201b60201c565b508060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505061027b565b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101e6826101bd565b9050919050565b5f6101f7826101dc565b9050919050565b610207816101ed565b8114610211575f80fd5b50565b5f81519050610222816101fe565b92915050565b5f6020828403121561023d5761023c6101b9565b5b5f61024a84828501610214565b91505092915050565b61025c816101dc565b82525050565b5f6020820190506102755f830184610253565b92915050565b6108d6806102885f395ff3fe608060405234801561000f575f80fd5b5060043610610060575f3560e01c806313b6731f146100645780632cd69d6014610080578063715018a6146100b05780638da5cb5b146100ba578063f2fde38b146100d8578063fc0c546a146100f4575b5f80fd5b61007e6004803603810190610079919061067a565b610112565b005b61009a600480360381019061009591906106de565b61023d565b6040516100a7919061072b565b60405180910390f35b6100b861025d565b005b6100c2610270565b6040516100cf9190610753565b60405180910390f35b6100f260048036038101906100ed919061076c565b610297565b005b6100fc61031b565b60405161010991906107f2565b60405180910390f35b61011a610340565b8060025f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8567ffffffffffffffff1667ffffffffffffffff1681526020019081526020015f205f8282546101899190610838565b925050819055508267ffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fb283270b87db7ad5d1fbb15af2039324aa28bebf00c89e37579882f7cb261d19836040516101e1919061072b565b60405180910390a361023784838360015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166103c7909392919063ffffffff16565b50505050565b6002602052815f5260405f20602052805f5260405f205f91509150505481565b610265610340565b61026e5f610449565b565b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b61029f610340565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361030f575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016103069190610753565b60405180910390fd5b61031881610449565b50565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b61034861050a565b73ffffffffffffffffffffffffffffffffffffffff16610366610270565b73ffffffffffffffffffffffffffffffffffffffff16146103c55761038961050a565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016103bc9190610753565b60405180910390fd5b565b610443848573ffffffffffffffffffffffffffffffffffffffff166323b872dd8686866040516024016103fc9392919061086b565b604051602081830303815290604052915060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050610511565b50505050565b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f33905090565b5f8060205f8451602086015f885af180610530576040513d5f823e3d81fd5b3d92505f519150505f8214610549576001811415610564565b5f8473ffffffffffffffffffffffffffffffffffffffff163b145b156105a657836040517f5274afe700000000000000000000000000000000000000000000000000000000815260040161059d9190610753565b60405180910390fd5b50505050565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6105d9826105b0565b9050919050565b6105e9816105cf565b81146105f3575f80fd5b50565b5f81359050610604816105e0565b92915050565b5f67ffffffffffffffff82169050919050565b6106268161060a565b8114610630575f80fd5b50565b5f813590506106418161061d565b92915050565b5f819050919050565b61065981610647565b8114610663575f80fd5b50565b5f8135905061067481610650565b92915050565b5f805f8060808587031215610692576106916105ac565b5b5f61069f878288016105f6565b94505060206106b087828801610633565b93505060406106c1878288016105f6565b92505060606106d287828801610666565b91505092959194509250565b5f80604083850312156106f4576106f36105ac565b5b5f610701858286016105f6565b925050602061071285828601610633565b9150509250929050565b61072581610647565b82525050565b5f60208201905061073e5f83018461071c565b92915050565b61074d816105cf565b82525050565b5f6020820190506107665f830184610744565b92915050565b5f60208284031215610781576107806105ac565b5b5f61078e848285016105f6565b91505092915050565b5f819050919050565b5f6107ba6107b56107b0846105b0565b610797565b6105b0565b9050919050565b5f6107cb826107a0565b9050919050565b5f6107dc826107c1565b9050919050565b6107ec816107d2565b82525050565b5f6020820190506108055f8301846107e3565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61084282610647565b915061084d83610647565b92508282019050808211156108655761086461080b565b5b92915050565b5f60608201905061087e5f830186610744565b61088b6020830185610744565b610898604083018461071c565b94935050505056fea2646970667358221220c00cfbe4d8ae78d900e71284acdc0737b43478c3ac304f90686ce125728c926d64736f6c63430008170033",
}

// SimplePaymentABI is the input ABI used to generate the binding from.
// Deprecated: Use SimplePaymentMetaData.ABI instead.
var SimplePaymentABI = SimplePaymentMetaData.ABI

// SimplePaymentBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SimplePaymentMetaData.Bin instead.
var SimplePaymentBin = SimplePaymentMetaData.Bin

// DeploySimplePayment deploys a new Ethereum contract, binding an instance of SimplePayment to it.
func DeploySimplePayment(auth *bind.TransactOpts, backend bind.ContractBackend, _token common.Address) (common.Address, *types.Transaction, *SimplePayment, error) {
	parsed, err := SimplePaymentMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimplePaymentBin), backend, _token)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimplePayment{SimplePaymentCaller: SimplePaymentCaller{contract: contract}, SimplePaymentTransactor: SimplePaymentTransactor{contract: contract}, SimplePaymentFilterer: SimplePaymentFilterer{contract: contract}}, nil
}

// SimplePayment is an auto generated Go binding around an Ethereum contract.
type SimplePayment struct {
	SimplePaymentCaller     // Read-only binding to the contract
	SimplePaymentTransactor // Write-only binding to the contract
	SimplePaymentFilterer   // Log filterer for contract events
}

// SimplePaymentCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimplePaymentCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimplePaymentTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimplePaymentTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimplePaymentFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimplePaymentFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimplePaymentSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimplePaymentSession struct {
	Contract     *SimplePayment    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SimplePaymentCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimplePaymentCallerSession struct {
	Contract *SimplePaymentCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SimplePaymentTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimplePaymentTransactorSession struct {
	Contract     *SimplePaymentTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SimplePaymentRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimplePaymentRaw struct {
	Contract *SimplePayment // Generic contract binding to access the raw methods on
}

// SimplePaymentCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimplePaymentCallerRaw struct {
	Contract *SimplePaymentCaller // Generic read-only contract binding to access the raw methods on
}

// SimplePaymentTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimplePaymentTransactorRaw struct {
	Contract *SimplePaymentTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimplePayment creates a new instance of SimplePayment, bound to a specific deployed contract.
func NewSimplePayment(address common.Address, backend bind.ContractBackend) (*SimplePayment, error) {
	contract, err := bindSimplePayment(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimplePayment{SimplePaymentCaller: SimplePaymentCaller{contract: contract}, SimplePaymentTransactor: SimplePaymentTransactor{contract: contract}, SimplePaymentFilterer: SimplePaymentFilterer{contract: contract}}, nil
}

// NewSimplePaymentCaller creates a new read-only instance of SimplePayment, bound to a specific deployed contract.
func NewSimplePaymentCaller(address common.Address, caller bind.ContractCaller) (*SimplePaymentCaller, error) {
	contract, err := bindSimplePayment(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimplePaymentCaller{contract: contract}, nil
}

// NewSimplePaymentTransactor creates a new write-only instance of SimplePayment, bound to a specific deployed contract.
func NewSimplePaymentTransactor(address common.Address, transactor bind.ContractTransactor) (*SimplePaymentTransactor, error) {
	contract, err := bindSimplePayment(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimplePaymentTransactor{contract: contract}, nil
}

// NewSimplePaymentFilterer creates a new log filterer instance of SimplePayment, bound to a specific deployed contract.
func NewSimplePaymentFilterer(address common.Address, filterer bind.ContractFilterer) (*SimplePaymentFilterer, error) {
	contract, err := bindSimplePayment(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimplePaymentFilterer{contract: contract}, nil
}

// bindSimplePayment binds a generic wrapper to an already deployed contract.
func bindSimplePayment(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SimplePaymentMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimplePayment *SimplePaymentRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimplePayment.Contract.SimplePaymentCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimplePayment *SimplePaymentRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimplePayment.Contract.SimplePaymentTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimplePayment *SimplePaymentRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimplePayment.Contract.SimplePaymentTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimplePayment *SimplePaymentCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimplePayment.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimplePayment *SimplePaymentTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimplePayment.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimplePayment *SimplePaymentTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimplePayment.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimplePayment *SimplePaymentCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimplePayment.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimplePayment *SimplePaymentSession) Owner() (common.Address, error) {
	return _SimplePayment.Contract.Owner(&_SimplePayment.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimplePayment *SimplePaymentCallerSession) Owner() (common.Address, error) {
	return _SimplePayment.Contract.Owner(&_SimplePayment.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_SimplePayment *SimplePaymentCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimplePayment.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_SimplePayment *SimplePaymentSession) Token() (common.Address, error) {
	return _SimplePayment.Contract.Token(&_SimplePayment.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_SimplePayment *SimplePaymentCallerSession) Token() (common.Address, error) {
	return _SimplePayment.Contract.Token(&_SimplePayment.CallOpts)
}

// TotalPaid is a free data retrieval call binding the contract method 0x2cd69d60.
//
// Solidity: function totalPaid(address payer, uint64 version) view returns(uint256)
func (_SimplePayment *SimplePaymentCaller) TotalPaid(opts *bind.CallOpts, payer common.Address, version uint64) (*big.Int, error) {
	var out []interface{}
	err := _SimplePayment.contract.Call(opts, &out, "totalPaid", payer, version)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalPaid is a free data retrieval call binding the contract method 0x2cd69d60.
//
// Solidity: function totalPaid(address payer, uint64 version) view returns(uint256)
func (_SimplePayment *SimplePaymentSession) TotalPaid(payer common.Address, version uint64) (*big.Int, error) {
	return _SimplePayment.Contract.TotalPaid(&_SimplePayment.CallOpts, payer, version)
}

// TotalPaid is a free data retrieval call binding the contract method 0x2cd69d60.
//
// Solidity: function totalPaid(address payer, uint64 version) view returns(uint256)
func (_SimplePayment *SimplePaymentCallerSession) TotalPaid(payer common.Address, version uint64) (*big.Int, error) {
	return _SimplePayment.Contract.TotalPaid(&_SimplePayment.CallOpts, payer, version)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimplePayment *SimplePaymentTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimplePayment.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimplePayment *SimplePaymentSession) RenounceOwnership() (*types.Transaction, error) {
	return _SimplePayment.Contract.RenounceOwnership(&_SimplePayment.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimplePayment *SimplePaymentTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SimplePayment.Contract.RenounceOwnership(&_SimplePayment.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimplePayment *SimplePaymentTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SimplePayment.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimplePayment *SimplePaymentSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SimplePayment.Contract.TransferOwnership(&_SimplePayment.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimplePayment *SimplePaymentTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SimplePayment.Contract.TransferOwnership(&_SimplePayment.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x13b6731f.
//
// Solidity: function withdraw(address payerAddress, uint64 version, address withdrawAddress, uint256 amount) returns()
func (_SimplePayment *SimplePaymentTransactor) Withdraw(opts *bind.TransactOpts, payerAddress common.Address, version uint64, withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SimplePayment.contract.Transact(opts, "withdraw", payerAddress, version, withdrawAddress, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x13b6731f.
//
// Solidity: function withdraw(address payerAddress, uint64 version, address withdrawAddress, uint256 amount) returns()
func (_SimplePayment *SimplePaymentSession) Withdraw(payerAddress common.Address, version uint64, withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SimplePayment.Contract.Withdraw(&_SimplePayment.TransactOpts, payerAddress, version, withdrawAddress, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x13b6731f.
//
// Solidity: function withdraw(address payerAddress, uint64 version, address withdrawAddress, uint256 amount) returns()
func (_SimplePayment *SimplePaymentTransactorSession) Withdraw(payerAddress common.Address, version uint64, withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SimplePayment.Contract.Withdraw(&_SimplePayment.TransactOpts, payerAddress, version, withdrawAddress, amount)
}

// SimplePaymentOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SimplePayment contract.
type SimplePaymentOwnershipTransferredIterator struct {
	Event *SimplePaymentOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SimplePaymentOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimplePaymentOwnershipTransferred)
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
		it.Event = new(SimplePaymentOwnershipTransferred)
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
func (it *SimplePaymentOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimplePaymentOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimplePaymentOwnershipTransferred represents a OwnershipTransferred event raised by the SimplePayment contract.
type SimplePaymentOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimplePayment *SimplePaymentFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SimplePaymentOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimplePayment.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SimplePaymentOwnershipTransferredIterator{contract: _SimplePayment.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimplePayment *SimplePaymentFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SimplePaymentOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimplePayment.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimplePaymentOwnershipTransferred)
				if err := _SimplePayment.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimplePayment *SimplePaymentFilterer) ParseOwnershipTransferred(log types.Log) (*SimplePaymentOwnershipTransferred, error) {
	event := new(SimplePaymentOwnershipTransferred)
	if err := _SimplePayment.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimplePaymentWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the SimplePayment contract.
type SimplePaymentWithdrawIterator struct {
	Event *SimplePaymentWithdraw // Event containing the contract specifics and raw log

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
func (it *SimplePaymentWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimplePaymentWithdraw)
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
		it.Event = new(SimplePaymentWithdraw)
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
func (it *SimplePaymentWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimplePaymentWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimplePaymentWithdraw represents a Withdraw event raised by the SimplePayment contract.
type SimplePaymentWithdraw struct {
	Payer   common.Address
	Version uint64
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xb283270b87db7ad5d1fbb15af2039324aa28bebf00c89e37579882f7cb261d19.
//
// Solidity: event Withdraw(address indexed payer, uint64 indexed version, uint256 amount)
func (_SimplePayment *SimplePaymentFilterer) FilterWithdraw(opts *bind.FilterOpts, payer []common.Address, version []uint64) (*SimplePaymentWithdrawIterator, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}
	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _SimplePayment.contract.FilterLogs(opts, "Withdraw", payerRule, versionRule)
	if err != nil {
		return nil, err
	}
	return &SimplePaymentWithdrawIterator{contract: _SimplePayment.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xb283270b87db7ad5d1fbb15af2039324aa28bebf00c89e37579882f7cb261d19.
//
// Solidity: event Withdraw(address indexed payer, uint64 indexed version, uint256 amount)
func (_SimplePayment *SimplePaymentFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *SimplePaymentWithdraw, payer []common.Address, version []uint64) (event.Subscription, error) {

	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}
	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _SimplePayment.contract.WatchLogs(opts, "Withdraw", payerRule, versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimplePaymentWithdraw)
				if err := _SimplePayment.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xb283270b87db7ad5d1fbb15af2039324aa28bebf00c89e37579882f7cb261d19.
//
// Solidity: event Withdraw(address indexed payer, uint64 indexed version, uint256 amount)
func (_SimplePayment *SimplePaymentFilterer) ParseWithdraw(log types.Log) (*SimplePaymentWithdraw, error) {
	event := new(SimplePaymentWithdraw)
	if err := _SimplePayment.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
