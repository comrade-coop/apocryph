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

// PaymentMetaData contains all meta data concerning the Payment contract.
var PaymentMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorize\",\"inputs\":[{\"name\":\"_authorized\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"available\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"channels\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"investedByPublisher\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawnByProvider\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unlockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unlockedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"closeChannel\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createChannel\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"unlockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"initialAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createSubChannel\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"newProvider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newPodId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isAuthorized\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unlock\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transferAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawUnlocked\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawUpTo\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"totalWithdrawAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transferAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawn\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ChannelClosed\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ChannelCreated\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"depositAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnlockTimerStarted\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"unlockedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unlocked\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"unlockedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"publisher\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"podId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"withdrawnAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AmountRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChannelLocked\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAuthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// PaymentABI is the input ABI used to generate the binding from.
// Deprecated: Use PaymentMetaData.ABI instead.
var PaymentABI = PaymentMetaData.ABI

// Payment is an auto generated Go binding around an Ethereum contract.
type Payment struct {
	PaymentCaller     // Read-only binding to the contract
	PaymentTransactor // Write-only binding to the contract
	PaymentFilterer   // Log filterer for contract events
}

// PaymentCaller is an auto generated read-only Go binding around an Ethereum contract.
type PaymentCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PaymentTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PaymentFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PaymentSession struct {
	Contract     *Payment          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PaymentCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PaymentCallerSession struct {
	Contract *PaymentCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// PaymentTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PaymentTransactorSession struct {
	Contract     *PaymentTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// PaymentRaw is an auto generated low-level Go binding around an Ethereum contract.
type PaymentRaw struct {
	Contract *Payment // Generic contract binding to access the raw methods on
}

// PaymentCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PaymentCallerRaw struct {
	Contract *PaymentCaller // Generic read-only contract binding to access the raw methods on
}

// PaymentTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PaymentTransactorRaw struct {
	Contract *PaymentTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPayment creates a new instance of Payment, bound to a specific deployed contract.
func NewPayment(address common.Address, backend bind.ContractBackend) (*Payment, error) {
	contract, err := bindPayment(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Payment{PaymentCaller: PaymentCaller{contract: contract}, PaymentTransactor: PaymentTransactor{contract: contract}, PaymentFilterer: PaymentFilterer{contract: contract}}, nil
}

// NewPaymentCaller creates a new read-only instance of Payment, bound to a specific deployed contract.
func NewPaymentCaller(address common.Address, caller bind.ContractCaller) (*PaymentCaller, error) {
	contract, err := bindPayment(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentCaller{contract: contract}, nil
}

// NewPaymentTransactor creates a new write-only instance of Payment, bound to a specific deployed contract.
func NewPaymentTransactor(address common.Address, transactor bind.ContractTransactor) (*PaymentTransactor, error) {
	contract, err := bindPayment(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentTransactor{contract: contract}, nil
}

// NewPaymentFilterer creates a new log filterer instance of Payment, bound to a specific deployed contract.
func NewPaymentFilterer(address common.Address, filterer bind.ContractFilterer) (*PaymentFilterer, error) {
	contract, err := bindPayment(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PaymentFilterer{contract: contract}, nil
}

// bindPayment binds a generic wrapper to an already deployed contract.
func bindPayment(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PaymentMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Payment *PaymentRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Payment.Contract.PaymentCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Payment *PaymentRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Payment.Contract.PaymentTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Payment *PaymentRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Payment.Contract.PaymentTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Payment *PaymentCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Payment.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Payment *PaymentTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Payment.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Payment *PaymentTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Payment.Contract.contract.Transact(opts, method, params...)
}

// Available is a free data retrieval call binding the contract method 0x8247820d.
//
// Solidity: function available(address publisher, address provider, bytes32 podId) view returns(uint256)
func (_Payment *PaymentCaller) Available(opts *bind.CallOpts, publisher common.Address, provider common.Address, podId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "available", publisher, provider, podId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Available is a free data retrieval call binding the contract method 0x8247820d.
//
// Solidity: function available(address publisher, address provider, bytes32 podId) view returns(uint256)
func (_Payment *PaymentSession) Available(publisher common.Address, provider common.Address, podId [32]byte) (*big.Int, error) {
	return _Payment.Contract.Available(&_Payment.CallOpts, publisher, provider, podId)
}

// Available is a free data retrieval call binding the contract method 0x8247820d.
//
// Solidity: function available(address publisher, address provider, bytes32 podId) view returns(uint256)
func (_Payment *PaymentCallerSession) Available(publisher common.Address, provider common.Address, podId [32]byte) (*big.Int, error) {
	return _Payment.Contract.Available(&_Payment.CallOpts, publisher, provider, podId)
}

// Channels is a free data retrieval call binding the contract method 0x098f26b5.
//
// Solidity: function channels(address , address , bytes32 ) view returns(uint256 investedByPublisher, uint256 withdrawnByProvider, uint256 unlockTime, uint256 unlockedAt)
func (_Payment *PaymentCaller) Channels(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 [32]byte) (struct {
	InvestedByPublisher *big.Int
	WithdrawnByProvider *big.Int
	UnlockTime          *big.Int
	UnlockedAt          *big.Int
}, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "channels", arg0, arg1, arg2)

	outstruct := new(struct {
		InvestedByPublisher *big.Int
		WithdrawnByProvider *big.Int
		UnlockTime          *big.Int
		UnlockedAt          *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.InvestedByPublisher = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.WithdrawnByProvider = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.UnlockTime = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UnlockedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Channels is a free data retrieval call binding the contract method 0x098f26b5.
//
// Solidity: function channels(address , address , bytes32 ) view returns(uint256 investedByPublisher, uint256 withdrawnByProvider, uint256 unlockTime, uint256 unlockedAt)
func (_Payment *PaymentSession) Channels(arg0 common.Address, arg1 common.Address, arg2 [32]byte) (struct {
	InvestedByPublisher *big.Int
	WithdrawnByProvider *big.Int
	UnlockTime          *big.Int
	UnlockedAt          *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2)
}

// Channels is a free data retrieval call binding the contract method 0x098f26b5.
//
// Solidity: function channels(address , address , bytes32 ) view returns(uint256 investedByPublisher, uint256 withdrawnByProvider, uint256 unlockTime, uint256 unlockedAt)
func (_Payment *PaymentCallerSession) Channels(arg0 common.Address, arg1 common.Address, arg2 [32]byte) (struct {
	InvestedByPublisher *big.Int
	WithdrawnByProvider *big.Int
	UnlockTime          *big.Int
	UnlockedAt          *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2)
}

// IsAuthorized is a free data retrieval call binding the contract method 0x9a99b662.
//
// Solidity: function isAuthorized(address publisher, address provider, bytes32 podId, address _address) view returns(bool)
func (_Payment *PaymentCaller) IsAuthorized(opts *bind.CallOpts, publisher common.Address, provider common.Address, podId [32]byte, _address common.Address) (bool, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "isAuthorized", publisher, provider, podId, _address)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAuthorized is a free data retrieval call binding the contract method 0x9a99b662.
//
// Solidity: function isAuthorized(address publisher, address provider, bytes32 podId, address _address) view returns(bool)
func (_Payment *PaymentSession) IsAuthorized(publisher common.Address, provider common.Address, podId [32]byte, _address common.Address) (bool, error) {
	return _Payment.Contract.IsAuthorized(&_Payment.CallOpts, publisher, provider, podId, _address)
}

// IsAuthorized is a free data retrieval call binding the contract method 0x9a99b662.
//
// Solidity: function isAuthorized(address publisher, address provider, bytes32 podId, address _address) view returns(bool)
func (_Payment *PaymentCallerSession) IsAuthorized(publisher common.Address, provider common.Address, podId [32]byte, _address common.Address) (bool, error) {
	return _Payment.Contract.IsAuthorized(&_Payment.CallOpts, publisher, provider, podId, _address)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Payment *PaymentCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Payment *PaymentSession) Token() (common.Address, error) {
	return _Payment.Contract.Token(&_Payment.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Payment *PaymentCallerSession) Token() (common.Address, error) {
	return _Payment.Contract.Token(&_Payment.CallOpts)
}

// Withdrawn is a free data retrieval call binding the contract method 0x6ee6a5da.
//
// Solidity: function withdrawn(address publisher, address provider, bytes32 podId) view returns(uint256)
func (_Payment *PaymentCaller) Withdrawn(opts *bind.CallOpts, publisher common.Address, provider common.Address, podId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "withdrawn", publisher, provider, podId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Withdrawn is a free data retrieval call binding the contract method 0x6ee6a5da.
//
// Solidity: function withdrawn(address publisher, address provider, bytes32 podId) view returns(uint256)
func (_Payment *PaymentSession) Withdrawn(publisher common.Address, provider common.Address, podId [32]byte) (*big.Int, error) {
	return _Payment.Contract.Withdrawn(&_Payment.CallOpts, publisher, provider, podId)
}

// Withdrawn is a free data retrieval call binding the contract method 0x6ee6a5da.
//
// Solidity: function withdrawn(address publisher, address provider, bytes32 podId) view returns(uint256)
func (_Payment *PaymentCallerSession) Withdrawn(publisher common.Address, provider common.Address, podId [32]byte) (*big.Int, error) {
	return _Payment.Contract.Withdrawn(&_Payment.CallOpts, publisher, provider, podId)
}

// Authorize is a paid mutator transaction binding the contract method 0xc970584a.
//
// Solidity: function authorize(address _authorized, address provider, bytes32 podId) returns()
func (_Payment *PaymentTransactor) Authorize(opts *bind.TransactOpts, _authorized common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "authorize", _authorized, provider, podId)
}

// Authorize is a paid mutator transaction binding the contract method 0xc970584a.
//
// Solidity: function authorize(address _authorized, address provider, bytes32 podId) returns()
func (_Payment *PaymentSession) Authorize(_authorized common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.Contract.Authorize(&_Payment.TransactOpts, _authorized, provider, podId)
}

// Authorize is a paid mutator transaction binding the contract method 0xc970584a.
//
// Solidity: function authorize(address _authorized, address provider, bytes32 podId) returns()
func (_Payment *PaymentTransactorSession) Authorize(_authorized common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.Contract.Authorize(&_Payment.TransactOpts, _authorized, provider, podId)
}

// CloseChannel is a paid mutator transaction binding the contract method 0x0bbbd884.
//
// Solidity: function closeChannel(address publisher, address provider, bytes32 podId) returns()
func (_Payment *PaymentTransactor) CloseChannel(opts *bind.TransactOpts, publisher common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "closeChannel", publisher, provider, podId)
}

// CloseChannel is a paid mutator transaction binding the contract method 0x0bbbd884.
//
// Solidity: function closeChannel(address publisher, address provider, bytes32 podId) returns()
func (_Payment *PaymentSession) CloseChannel(publisher common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.Contract.CloseChannel(&_Payment.TransactOpts, publisher, provider, podId)
}

// CloseChannel is a paid mutator transaction binding the contract method 0x0bbbd884.
//
// Solidity: function closeChannel(address publisher, address provider, bytes32 podId) returns()
func (_Payment *PaymentTransactorSession) CloseChannel(publisher common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.Contract.CloseChannel(&_Payment.TransactOpts, publisher, provider, podId)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x744122b0.
//
// Solidity: function createChannel(address provider, bytes32 podId, uint256 unlockTime, uint256 initialAmount) returns()
func (_Payment *PaymentTransactor) CreateChannel(opts *bind.TransactOpts, provider common.Address, podId [32]byte, unlockTime *big.Int, initialAmount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "createChannel", provider, podId, unlockTime, initialAmount)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x744122b0.
//
// Solidity: function createChannel(address provider, bytes32 podId, uint256 unlockTime, uint256 initialAmount) returns()
func (_Payment *PaymentSession) CreateChannel(provider common.Address, podId [32]byte, unlockTime *big.Int, initialAmount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateChannel(&_Payment.TransactOpts, provider, podId, unlockTime, initialAmount)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x744122b0.
//
// Solidity: function createChannel(address provider, bytes32 podId, uint256 unlockTime, uint256 initialAmount) returns()
func (_Payment *PaymentTransactorSession) CreateChannel(provider common.Address, podId [32]byte, unlockTime *big.Int, initialAmount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateChannel(&_Payment.TransactOpts, provider, podId, unlockTime, initialAmount)
}

// CreateSubChannel is a paid mutator transaction binding the contract method 0x67ca76a6.
//
// Solidity: function createSubChannel(address publisher, address provider, bytes32 podId, address newProvider, bytes32 newPodId, uint256 amount) returns()
func (_Payment *PaymentTransactor) CreateSubChannel(opts *bind.TransactOpts, publisher common.Address, provider common.Address, podId [32]byte, newProvider common.Address, newPodId [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "createSubChannel", publisher, provider, podId, newProvider, newPodId, amount)
}

// CreateSubChannel is a paid mutator transaction binding the contract method 0x67ca76a6.
//
// Solidity: function createSubChannel(address publisher, address provider, bytes32 podId, address newProvider, bytes32 newPodId, uint256 amount) returns()
func (_Payment *PaymentSession) CreateSubChannel(publisher common.Address, provider common.Address, podId [32]byte, newProvider common.Address, newPodId [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateSubChannel(&_Payment.TransactOpts, publisher, provider, podId, newProvider, newPodId, amount)
}

// CreateSubChannel is a paid mutator transaction binding the contract method 0x67ca76a6.
//
// Solidity: function createSubChannel(address publisher, address provider, bytes32 podId, address newProvider, bytes32 newPodId, uint256 amount) returns()
func (_Payment *PaymentTransactorSession) CreateSubChannel(publisher common.Address, provider common.Address, podId [32]byte, newProvider common.Address, newPodId [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateSubChannel(&_Payment.TransactOpts, publisher, provider, podId, newProvider, newPodId, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xeb2243f8.
//
// Solidity: function deposit(address provider, bytes32 podId, uint256 amount) returns()
func (_Payment *PaymentTransactor) Deposit(opts *bind.TransactOpts, provider common.Address, podId [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "deposit", provider, podId, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xeb2243f8.
//
// Solidity: function deposit(address provider, bytes32 podId, uint256 amount) returns()
func (_Payment *PaymentSession) Deposit(provider common.Address, podId [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.Deposit(&_Payment.TransactOpts, provider, podId, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xeb2243f8.
//
// Solidity: function deposit(address provider, bytes32 podId, uint256 amount) returns()
func (_Payment *PaymentTransactorSession) Deposit(provider common.Address, podId [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.Deposit(&_Payment.TransactOpts, provider, podId, amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x77c5f23c.
//
// Solidity: function unlock(address publisher, address provider, bytes32 podId) returns()
func (_Payment *PaymentTransactor) Unlock(opts *bind.TransactOpts, publisher common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "unlock", publisher, provider, podId)
}

// Unlock is a paid mutator transaction binding the contract method 0x77c5f23c.
//
// Solidity: function unlock(address publisher, address provider, bytes32 podId) returns()
func (_Payment *PaymentSession) Unlock(publisher common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.Contract.Unlock(&_Payment.TransactOpts, publisher, provider, podId)
}

// Unlock is a paid mutator transaction binding the contract method 0x77c5f23c.
//
// Solidity: function unlock(address publisher, address provider, bytes32 podId) returns()
func (_Payment *PaymentTransactorSession) Unlock(publisher common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.Contract.Unlock(&_Payment.TransactOpts, publisher, provider, podId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x92a453a5.
//
// Solidity: function withdraw(address publisher, bytes32 podId, uint256 amount, address transferAddress) returns()
func (_Payment *PaymentTransactor) Withdraw(opts *bind.TransactOpts, publisher common.Address, podId [32]byte, amount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdraw", publisher, podId, amount, transferAddress)
}

// Withdraw is a paid mutator transaction binding the contract method 0x92a453a5.
//
// Solidity: function withdraw(address publisher, bytes32 podId, uint256 amount, address transferAddress) returns()
func (_Payment *PaymentSession) Withdraw(publisher common.Address, podId [32]byte, amount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, publisher, podId, amount, transferAddress)
}

// Withdraw is a paid mutator transaction binding the contract method 0x92a453a5.
//
// Solidity: function withdraw(address publisher, bytes32 podId, uint256 amount, address transferAddress) returns()
func (_Payment *PaymentTransactorSession) Withdraw(publisher common.Address, podId [32]byte, amount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, publisher, podId, amount, transferAddress)
}

// WithdrawUnlocked is a paid mutator transaction binding the contract method 0x8e1577f4.
//
// Solidity: function withdrawUnlocked(address publisher, address provider, bytes32 podId) returns()
func (_Payment *PaymentTransactor) WithdrawUnlocked(opts *bind.TransactOpts, publisher common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdrawUnlocked", publisher, provider, podId)
}

// WithdrawUnlocked is a paid mutator transaction binding the contract method 0x8e1577f4.
//
// Solidity: function withdrawUnlocked(address publisher, address provider, bytes32 podId) returns()
func (_Payment *PaymentSession) WithdrawUnlocked(publisher common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUnlocked(&_Payment.TransactOpts, publisher, provider, podId)
}

// WithdrawUnlocked is a paid mutator transaction binding the contract method 0x8e1577f4.
//
// Solidity: function withdrawUnlocked(address publisher, address provider, bytes32 podId) returns()
func (_Payment *PaymentTransactorSession) WithdrawUnlocked(publisher common.Address, provider common.Address, podId [32]byte) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUnlocked(&_Payment.TransactOpts, publisher, provider, podId)
}

// WithdrawUpTo is a paid mutator transaction binding the contract method 0x494f8587.
//
// Solidity: function withdrawUpTo(address publisher, bytes32 podId, uint256 totalWithdrawAmount, address transferAddress) returns()
func (_Payment *PaymentTransactor) WithdrawUpTo(opts *bind.TransactOpts, publisher common.Address, podId [32]byte, totalWithdrawAmount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdrawUpTo", publisher, podId, totalWithdrawAmount, transferAddress)
}

// WithdrawUpTo is a paid mutator transaction binding the contract method 0x494f8587.
//
// Solidity: function withdrawUpTo(address publisher, bytes32 podId, uint256 totalWithdrawAmount, address transferAddress) returns()
func (_Payment *PaymentSession) WithdrawUpTo(publisher common.Address, podId [32]byte, totalWithdrawAmount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUpTo(&_Payment.TransactOpts, publisher, podId, totalWithdrawAmount, transferAddress)
}

// WithdrawUpTo is a paid mutator transaction binding the contract method 0x494f8587.
//
// Solidity: function withdrawUpTo(address publisher, bytes32 podId, uint256 totalWithdrawAmount, address transferAddress) returns()
func (_Payment *PaymentTransactorSession) WithdrawUpTo(publisher common.Address, podId [32]byte, totalWithdrawAmount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUpTo(&_Payment.TransactOpts, publisher, podId, totalWithdrawAmount, transferAddress)
}

// PaymentChannelClosedIterator is returned from FilterChannelClosed and is used to iterate over the raw logs and unpacked data for ChannelClosed events raised by the Payment contract.
type PaymentChannelClosedIterator struct {
	Event *PaymentChannelClosed // Event containing the contract specifics and raw log

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
func (it *PaymentChannelClosedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentChannelClosed)
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
		it.Event = new(PaymentChannelClosed)
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
func (it *PaymentChannelClosedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentChannelClosedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentChannelClosed represents a ChannelClosed event raised by the Payment contract.
type PaymentChannelClosed struct {
	Publisher common.Address
	Provider  common.Address
	PodId     [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelClosed is a free log retrieval operation binding the contract event 0xa1de5c5c82ac5d38cba67d62238ee6b160c22d9f7f697de78ca03072f271224d.
//
// Solidity: event ChannelClosed(address indexed publisher, address indexed provider, bytes32 indexed podId)
func (_Payment *PaymentFilterer) FilterChannelClosed(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, podId [][32]byte) (*PaymentChannelClosedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "ChannelClosed", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelClosedIterator{contract: _Payment.contract, event: "ChannelClosed", logs: logs, sub: sub}, nil
}

// WatchChannelClosed is a free log subscription operation binding the contract event 0xa1de5c5c82ac5d38cba67d62238ee6b160c22d9f7f697de78ca03072f271224d.
//
// Solidity: event ChannelClosed(address indexed publisher, address indexed provider, bytes32 indexed podId)
func (_Payment *PaymentFilterer) WatchChannelClosed(opts *bind.WatchOpts, sink chan<- *PaymentChannelClosed, publisher []common.Address, provider []common.Address, podId [][32]byte) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "ChannelClosed", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentChannelClosed)
				if err := _Payment.contract.UnpackLog(event, "ChannelClosed", log); err != nil {
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

// ParseChannelClosed is a log parse operation binding the contract event 0xa1de5c5c82ac5d38cba67d62238ee6b160c22d9f7f697de78ca03072f271224d.
//
// Solidity: event ChannelClosed(address indexed publisher, address indexed provider, bytes32 indexed podId)
func (_Payment *PaymentFilterer) ParseChannelClosed(log types.Log) (*PaymentChannelClosed, error) {
	event := new(PaymentChannelClosed)
	if err := _Payment.contract.UnpackLog(event, "ChannelClosed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentChannelCreatedIterator is returned from FilterChannelCreated and is used to iterate over the raw logs and unpacked data for ChannelCreated events raised by the Payment contract.
type PaymentChannelCreatedIterator struct {
	Event *PaymentChannelCreated // Event containing the contract specifics and raw log

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
func (it *PaymentChannelCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentChannelCreated)
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
		it.Event = new(PaymentChannelCreated)
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
func (it *PaymentChannelCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentChannelCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentChannelCreated represents a ChannelCreated event raised by the Payment contract.
type PaymentChannelCreated struct {
	Publisher common.Address
	Provider  common.Address
	PodId     [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelCreated is a free log retrieval operation binding the contract event 0xfa5bcf5a0e05d18398db4fa9c5f1f3ab260ff27cec7504aef225d23ad65db460.
//
// Solidity: event ChannelCreated(address indexed publisher, address indexed provider, bytes32 indexed podId)
func (_Payment *PaymentFilterer) FilterChannelCreated(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, podId [][32]byte) (*PaymentChannelCreatedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "ChannelCreated", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelCreatedIterator{contract: _Payment.contract, event: "ChannelCreated", logs: logs, sub: sub}, nil
}

// WatchChannelCreated is a free log subscription operation binding the contract event 0xfa5bcf5a0e05d18398db4fa9c5f1f3ab260ff27cec7504aef225d23ad65db460.
//
// Solidity: event ChannelCreated(address indexed publisher, address indexed provider, bytes32 indexed podId)
func (_Payment *PaymentFilterer) WatchChannelCreated(opts *bind.WatchOpts, sink chan<- *PaymentChannelCreated, publisher []common.Address, provider []common.Address, podId [][32]byte) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "ChannelCreated", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentChannelCreated)
				if err := _Payment.contract.UnpackLog(event, "ChannelCreated", log); err != nil {
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

// ParseChannelCreated is a log parse operation binding the contract event 0xfa5bcf5a0e05d18398db4fa9c5f1f3ab260ff27cec7504aef225d23ad65db460.
//
// Solidity: event ChannelCreated(address indexed publisher, address indexed provider, bytes32 indexed podId)
func (_Payment *PaymentFilterer) ParseChannelCreated(log types.Log) (*PaymentChannelCreated, error) {
	event := new(PaymentChannelCreated)
	if err := _Payment.contract.UnpackLog(event, "ChannelCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the Payment contract.
type PaymentDepositedIterator struct {
	Event *PaymentDeposited // Event containing the contract specifics and raw log

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
func (it *PaymentDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentDeposited)
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
		it.Event = new(PaymentDeposited)
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
func (it *PaymentDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentDeposited represents a Deposited event raised by the Payment contract.
type PaymentDeposited struct {
	Publisher     common.Address
	Provider      common.Address
	PodId         [32]byte
	DepositAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x1765ac22a11675a51476924171a95254d041fd475340f385cb8f27335aa80ee7.
//
// Solidity: event Deposited(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 depositAmount)
func (_Payment *PaymentFilterer) FilterDeposited(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, podId [][32]byte) (*PaymentDepositedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "Deposited", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return &PaymentDepositedIterator{contract: _Payment.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x1765ac22a11675a51476924171a95254d041fd475340f385cb8f27335aa80ee7.
//
// Solidity: event Deposited(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 depositAmount)
func (_Payment *PaymentFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *PaymentDeposited, publisher []common.Address, provider []common.Address, podId [][32]byte) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "Deposited", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentDeposited)
				if err := _Payment.contract.UnpackLog(event, "Deposited", log); err != nil {
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

// ParseDeposited is a log parse operation binding the contract event 0x1765ac22a11675a51476924171a95254d041fd475340f385cb8f27335aa80ee7.
//
// Solidity: event Deposited(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 depositAmount)
func (_Payment *PaymentFilterer) ParseDeposited(log types.Log) (*PaymentDeposited, error) {
	event := new(PaymentDeposited)
	if err := _Payment.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentUnlockTimerStartedIterator is returned from FilterUnlockTimerStarted and is used to iterate over the raw logs and unpacked data for UnlockTimerStarted events raised by the Payment contract.
type PaymentUnlockTimerStartedIterator struct {
	Event *PaymentUnlockTimerStarted // Event containing the contract specifics and raw log

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
func (it *PaymentUnlockTimerStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentUnlockTimerStarted)
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
		it.Event = new(PaymentUnlockTimerStarted)
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
func (it *PaymentUnlockTimerStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentUnlockTimerStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentUnlockTimerStarted represents a UnlockTimerStarted event raised by the Payment contract.
type PaymentUnlockTimerStarted struct {
	Publisher  common.Address
	Provider   common.Address
	PodId      [32]byte
	UnlockedAt *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUnlockTimerStarted is a free log retrieval operation binding the contract event 0xd5acf6c94da297d8c63f389d5c3a926a68be8206e29d85d0372009d088b2b2b5.
//
// Solidity: event UnlockTimerStarted(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 unlockedAt)
func (_Payment *PaymentFilterer) FilterUnlockTimerStarted(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, podId [][32]byte) (*PaymentUnlockTimerStartedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "UnlockTimerStarted", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return &PaymentUnlockTimerStartedIterator{contract: _Payment.contract, event: "UnlockTimerStarted", logs: logs, sub: sub}, nil
}

// WatchUnlockTimerStarted is a free log subscription operation binding the contract event 0xd5acf6c94da297d8c63f389d5c3a926a68be8206e29d85d0372009d088b2b2b5.
//
// Solidity: event UnlockTimerStarted(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 unlockedAt)
func (_Payment *PaymentFilterer) WatchUnlockTimerStarted(opts *bind.WatchOpts, sink chan<- *PaymentUnlockTimerStarted, publisher []common.Address, provider []common.Address, podId [][32]byte) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "UnlockTimerStarted", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentUnlockTimerStarted)
				if err := _Payment.contract.UnpackLog(event, "UnlockTimerStarted", log); err != nil {
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

// ParseUnlockTimerStarted is a log parse operation binding the contract event 0xd5acf6c94da297d8c63f389d5c3a926a68be8206e29d85d0372009d088b2b2b5.
//
// Solidity: event UnlockTimerStarted(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 unlockedAt)
func (_Payment *PaymentFilterer) ParseUnlockTimerStarted(log types.Log) (*PaymentUnlockTimerStarted, error) {
	event := new(PaymentUnlockTimerStarted)
	if err := _Payment.contract.UnpackLog(event, "UnlockTimerStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentUnlockedIterator is returned from FilterUnlocked and is used to iterate over the raw logs and unpacked data for Unlocked events raised by the Payment contract.
type PaymentUnlockedIterator struct {
	Event *PaymentUnlocked // Event containing the contract specifics and raw log

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
func (it *PaymentUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentUnlocked)
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
		it.Event = new(PaymentUnlocked)
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
func (it *PaymentUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentUnlocked represents a Unlocked event raised by the Payment contract.
type PaymentUnlocked struct {
	Publisher      common.Address
	Provider       common.Address
	PodId          [32]byte
	UnlockedAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUnlocked is a free log retrieval operation binding the contract event 0xfa21f248f6235facbe9f2ea2b2b65d8d35d9db4dffab0370f04371e697a391b3.
//
// Solidity: event Unlocked(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 unlockedAmount)
func (_Payment *PaymentFilterer) FilterUnlocked(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, podId [][32]byte) (*PaymentUnlockedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "Unlocked", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return &PaymentUnlockedIterator{contract: _Payment.contract, event: "Unlocked", logs: logs, sub: sub}, nil
}

// WatchUnlocked is a free log subscription operation binding the contract event 0xfa21f248f6235facbe9f2ea2b2b65d8d35d9db4dffab0370f04371e697a391b3.
//
// Solidity: event Unlocked(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 unlockedAmount)
func (_Payment *PaymentFilterer) WatchUnlocked(opts *bind.WatchOpts, sink chan<- *PaymentUnlocked, publisher []common.Address, provider []common.Address, podId [][32]byte) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "Unlocked", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentUnlocked)
				if err := _Payment.contract.UnpackLog(event, "Unlocked", log); err != nil {
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

// ParseUnlocked is a log parse operation binding the contract event 0xfa21f248f6235facbe9f2ea2b2b65d8d35d9db4dffab0370f04371e697a391b3.
//
// Solidity: event Unlocked(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 unlockedAmount)
func (_Payment *PaymentFilterer) ParseUnlocked(log types.Log) (*PaymentUnlocked, error) {
	event := new(PaymentUnlocked)
	if err := _Payment.contract.UnpackLog(event, "Unlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the Payment contract.
type PaymentWithdrawnIterator struct {
	Event *PaymentWithdrawn // Event containing the contract specifics and raw log

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
func (it *PaymentWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentWithdrawn)
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
		it.Event = new(PaymentWithdrawn)
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
func (it *PaymentWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentWithdrawn represents a Withdrawn event raised by the Payment contract.
type PaymentWithdrawn struct {
	Publisher       common.Address
	Provider        common.Address
	PodId           [32]byte
	WithdrawnAmount *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x1a96470c27c63e239afca174a1d71aa5726bbc196ca104fc8f5a9be2ed5080c4.
//
// Solidity: event Withdrawn(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 withdrawnAmount)
func (_Payment *PaymentFilterer) FilterWithdrawn(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, podId [][32]byte) (*PaymentWithdrawnIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "Withdrawn", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return &PaymentWithdrawnIterator{contract: _Payment.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x1a96470c27c63e239afca174a1d71aa5726bbc196ca104fc8f5a9be2ed5080c4.
//
// Solidity: event Withdrawn(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 withdrawnAmount)
func (_Payment *PaymentFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *PaymentWithdrawn, publisher []common.Address, provider []common.Address, podId [][32]byte) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var podIdRule []interface{}
	for _, podIdItem := range podId {
		podIdRule = append(podIdRule, podIdItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "Withdrawn", publisherRule, providerRule, podIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentWithdrawn)
				if err := _Payment.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0x1a96470c27c63e239afca174a1d71aa5726bbc196ca104fc8f5a9be2ed5080c4.
//
// Solidity: event Withdrawn(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 withdrawnAmount)
func (_Payment *PaymentFilterer) ParseWithdrawn(log types.Log) (*PaymentWithdrawn, error) {
	event := new(PaymentWithdrawn)
	if err := _Payment.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
