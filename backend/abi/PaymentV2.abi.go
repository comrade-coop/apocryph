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

// PaymentV2MetaData contains all meta data concerning the PaymentV2 contract.
var PaymentV2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"_unlockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorize\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissions\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"reservedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"extraLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"channelAuthorizations\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"reservation\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"limit\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"unlocksAt\",\"type\":\"uint240\",\"internalType\":\"uint240\"},{\"name\":\"permissions\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"channels\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"available\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"reserved\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"create\",\"inputs\":[{\"name\":\"channelDiscriminator\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"initialAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createAndAuthorize\",\"inputs\":[{\"name\":\"channelDiscriminator\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"initialAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"permissions\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"reservedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getChannelId\",\"inputs\":[{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"channelDiscriminator\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unlock\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"transferAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"transferAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Authorize\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"permissions\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"reservedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"limit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unlock\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"channelId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AuthorizationLocked\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAuthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// PaymentV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use PaymentV2MetaData.ABI instead.
var PaymentV2ABI = PaymentV2MetaData.ABI

// PaymentV2 is an auto generated Go binding around an Ethereum contract.
type PaymentV2 struct {
	PaymentV2Caller     // Read-only binding to the contract
	PaymentV2Transactor // Write-only binding to the contract
	PaymentV2Filterer   // Log filterer for contract events
}

// PaymentV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type PaymentV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type PaymentV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PaymentV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PaymentV2Session struct {
	Contract     *PaymentV2        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PaymentV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PaymentV2CallerSession struct {
	Contract *PaymentV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// PaymentV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PaymentV2TransactorSession struct {
	Contract     *PaymentV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// PaymentV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type PaymentV2Raw struct {
	Contract *PaymentV2 // Generic contract binding to access the raw methods on
}

// PaymentV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PaymentV2CallerRaw struct {
	Contract *PaymentV2Caller // Generic read-only contract binding to access the raw methods on
}

// PaymentV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PaymentV2TransactorRaw struct {
	Contract *PaymentV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewPaymentV2 creates a new instance of PaymentV2, bound to a specific deployed contract.
func NewPaymentV2(address common.Address, backend bind.ContractBackend) (*PaymentV2, error) {
	contract, err := bindPaymentV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PaymentV2{PaymentV2Caller: PaymentV2Caller{contract: contract}, PaymentV2Transactor: PaymentV2Transactor{contract: contract}, PaymentV2Filterer: PaymentV2Filterer{contract: contract}}, nil
}

// NewPaymentV2Caller creates a new read-only instance of PaymentV2, bound to a specific deployed contract.
func NewPaymentV2Caller(address common.Address, caller bind.ContractCaller) (*PaymentV2Caller, error) {
	contract, err := bindPaymentV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentV2Caller{contract: contract}, nil
}

// NewPaymentV2Transactor creates a new write-only instance of PaymentV2, bound to a specific deployed contract.
func NewPaymentV2Transactor(address common.Address, transactor bind.ContractTransactor) (*PaymentV2Transactor, error) {
	contract, err := bindPaymentV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentV2Transactor{contract: contract}, nil
}

// NewPaymentV2Filterer creates a new log filterer instance of PaymentV2, bound to a specific deployed contract.
func NewPaymentV2Filterer(address common.Address, filterer bind.ContractFilterer) (*PaymentV2Filterer, error) {
	contract, err := bindPaymentV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PaymentV2Filterer{contract: contract}, nil
}

// bindPaymentV2 binds a generic wrapper to an already deployed contract.
func bindPaymentV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PaymentV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PaymentV2 *PaymentV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PaymentV2.Contract.PaymentV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PaymentV2 *PaymentV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PaymentV2.Contract.PaymentV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PaymentV2 *PaymentV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PaymentV2.Contract.PaymentV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PaymentV2 *PaymentV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PaymentV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PaymentV2 *PaymentV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PaymentV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PaymentV2 *PaymentV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PaymentV2.Contract.contract.Transact(opts, method, params...)
}

// ChannelAuthorizations is a free data retrieval call binding the contract method 0x61f02b61.
//
// Solidity: function channelAuthorizations(bytes32 , address ) view returns(uint128 reservation, uint128 limit, uint240 unlocksAt, uint16 permissions)
func (_PaymentV2 *PaymentV2Caller) ChannelAuthorizations(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (struct {
	Reservation *big.Int
	Limit       *big.Int
	UnlocksAt   *big.Int
	Permissions uint16
}, error) {
	var out []interface{}
	err := _PaymentV2.contract.Call(opts, &out, "channelAuthorizations", arg0, arg1)

	outstruct := new(struct {
		Reservation *big.Int
		Limit       *big.Int
		UnlocksAt   *big.Int
		Permissions uint16
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reservation = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Limit = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.UnlocksAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Permissions = *abi.ConvertType(out[3], new(uint16)).(*uint16)

	return *outstruct, err

}

// ChannelAuthorizations is a free data retrieval call binding the contract method 0x61f02b61.
//
// Solidity: function channelAuthorizations(bytes32 , address ) view returns(uint128 reservation, uint128 limit, uint240 unlocksAt, uint16 permissions)
func (_PaymentV2 *PaymentV2Session) ChannelAuthorizations(arg0 [32]byte, arg1 common.Address) (struct {
	Reservation *big.Int
	Limit       *big.Int
	UnlocksAt   *big.Int
	Permissions uint16
}, error) {
	return _PaymentV2.Contract.ChannelAuthorizations(&_PaymentV2.CallOpts, arg0, arg1)
}

// ChannelAuthorizations is a free data retrieval call binding the contract method 0x61f02b61.
//
// Solidity: function channelAuthorizations(bytes32 , address ) view returns(uint128 reservation, uint128 limit, uint240 unlocksAt, uint16 permissions)
func (_PaymentV2 *PaymentV2CallerSession) ChannelAuthorizations(arg0 [32]byte, arg1 common.Address) (struct {
	Reservation *big.Int
	Limit       *big.Int
	UnlocksAt   *big.Int
	Permissions uint16
}, error) {
	return _PaymentV2.Contract.ChannelAuthorizations(&_PaymentV2.CallOpts, arg0, arg1)
}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels(bytes32 ) view returns(uint128 available, uint128 reserved)
func (_PaymentV2 *PaymentV2Caller) Channels(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Available *big.Int
	Reserved  *big.Int
}, error) {
	var out []interface{}
	err := _PaymentV2.contract.Call(opts, &out, "channels", arg0)

	outstruct := new(struct {
		Available *big.Int
		Reserved  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Available = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserved = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels(bytes32 ) view returns(uint128 available, uint128 reserved)
func (_PaymentV2 *PaymentV2Session) Channels(arg0 [32]byte) (struct {
	Available *big.Int
	Reserved  *big.Int
}, error) {
	return _PaymentV2.Contract.Channels(&_PaymentV2.CallOpts, arg0)
}

// Channels is a free data retrieval call binding the contract method 0x7a7ebd7b.
//
// Solidity: function channels(bytes32 ) view returns(uint128 available, uint128 reserved)
func (_PaymentV2 *PaymentV2CallerSession) Channels(arg0 [32]byte) (struct {
	Available *big.Int
	Reserved  *big.Int
}, error) {
	return _PaymentV2.Contract.Channels(&_PaymentV2.CallOpts, arg0)
}

// GetChannelId is a free data retrieval call binding the contract method 0xbd654669.
//
// Solidity: function getChannelId(address creator, bytes32 channelDiscriminator) pure returns(bytes32 channelId)
func (_PaymentV2 *PaymentV2Caller) GetChannelId(opts *bind.CallOpts, creator common.Address, channelDiscriminator [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _PaymentV2.contract.Call(opts, &out, "getChannelId", creator, channelDiscriminator)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetChannelId is a free data retrieval call binding the contract method 0xbd654669.
//
// Solidity: function getChannelId(address creator, bytes32 channelDiscriminator) pure returns(bytes32 channelId)
func (_PaymentV2 *PaymentV2Session) GetChannelId(creator common.Address, channelDiscriminator [32]byte) ([32]byte, error) {
	return _PaymentV2.Contract.GetChannelId(&_PaymentV2.CallOpts, creator, channelDiscriminator)
}

// GetChannelId is a free data retrieval call binding the contract method 0xbd654669.
//
// Solidity: function getChannelId(address creator, bytes32 channelDiscriminator) pure returns(bytes32 channelId)
func (_PaymentV2 *PaymentV2CallerSession) GetChannelId(creator common.Address, channelDiscriminator [32]byte) ([32]byte, error) {
	return _PaymentV2.Contract.GetChannelId(&_PaymentV2.CallOpts, creator, channelDiscriminator)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_PaymentV2 *PaymentV2Caller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PaymentV2.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_PaymentV2 *PaymentV2Session) Token() (common.Address, error) {
	return _PaymentV2.Contract.Token(&_PaymentV2.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_PaymentV2 *PaymentV2CallerSession) Token() (common.Address, error) {
	return _PaymentV2.Contract.Token(&_PaymentV2.CallOpts)
}

// UnlockTime is a free data retrieval call binding the contract method 0x251c1aa3.
//
// Solidity: function unlockTime() view returns(uint256)
func (_PaymentV2 *PaymentV2Caller) UnlockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PaymentV2.contract.Call(opts, &out, "unlockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnlockTime is a free data retrieval call binding the contract method 0x251c1aa3.
//
// Solidity: function unlockTime() view returns(uint256)
func (_PaymentV2 *PaymentV2Session) UnlockTime() (*big.Int, error) {
	return _PaymentV2.Contract.UnlockTime(&_PaymentV2.CallOpts)
}

// UnlockTime is a free data retrieval call binding the contract method 0x251c1aa3.
//
// Solidity: function unlockTime() view returns(uint256)
func (_PaymentV2 *PaymentV2CallerSession) UnlockTime() (*big.Int, error) {
	return _PaymentV2.Contract.UnlockTime(&_PaymentV2.CallOpts)
}

// Authorize is a paid mutator transaction binding the contract method 0x0221d6a3.
//
// Solidity: function authorize(bytes32 channelId, address recipient, uint16 permissions, uint256 reservedAmount, uint256 extraLimit) returns()
func (_PaymentV2 *PaymentV2Transactor) Authorize(opts *bind.TransactOpts, channelId [32]byte, recipient common.Address, permissions uint16, reservedAmount *big.Int, extraLimit *big.Int) (*types.Transaction, error) {
	return _PaymentV2.contract.Transact(opts, "authorize", channelId, recipient, permissions, reservedAmount, extraLimit)
}

// Authorize is a paid mutator transaction binding the contract method 0x0221d6a3.
//
// Solidity: function authorize(bytes32 channelId, address recipient, uint16 permissions, uint256 reservedAmount, uint256 extraLimit) returns()
func (_PaymentV2 *PaymentV2Session) Authorize(channelId [32]byte, recipient common.Address, permissions uint16, reservedAmount *big.Int, extraLimit *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.Authorize(&_PaymentV2.TransactOpts, channelId, recipient, permissions, reservedAmount, extraLimit)
}

// Authorize is a paid mutator transaction binding the contract method 0x0221d6a3.
//
// Solidity: function authorize(bytes32 channelId, address recipient, uint16 permissions, uint256 reservedAmount, uint256 extraLimit) returns()
func (_PaymentV2 *PaymentV2TransactorSession) Authorize(channelId [32]byte, recipient common.Address, permissions uint16, reservedAmount *big.Int, extraLimit *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.Authorize(&_PaymentV2.TransactOpts, channelId, recipient, permissions, reservedAmount, extraLimit)
}

// Create is a paid mutator transaction binding the contract method 0xa042c132.
//
// Solidity: function create(bytes32 channelDiscriminator, uint256 initialAmount) returns(bytes32 channelId)
func (_PaymentV2 *PaymentV2Transactor) Create(opts *bind.TransactOpts, channelDiscriminator [32]byte, initialAmount *big.Int) (*types.Transaction, error) {
	return _PaymentV2.contract.Transact(opts, "create", channelDiscriminator, initialAmount)
}

// Create is a paid mutator transaction binding the contract method 0xa042c132.
//
// Solidity: function create(bytes32 channelDiscriminator, uint256 initialAmount) returns(bytes32 channelId)
func (_PaymentV2 *PaymentV2Session) Create(channelDiscriminator [32]byte, initialAmount *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.Create(&_PaymentV2.TransactOpts, channelDiscriminator, initialAmount)
}

// Create is a paid mutator transaction binding the contract method 0xa042c132.
//
// Solidity: function create(bytes32 channelDiscriminator, uint256 initialAmount) returns(bytes32 channelId)
func (_PaymentV2 *PaymentV2TransactorSession) Create(channelDiscriminator [32]byte, initialAmount *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.Create(&_PaymentV2.TransactOpts, channelDiscriminator, initialAmount)
}

// CreateAndAuthorize is a paid mutator transaction binding the contract method 0xe333df11.
//
// Solidity: function createAndAuthorize(bytes32 channelDiscriminator, uint256 initialAmount, address recipient, uint16 permissions, uint256 reservedAmount, uint256 limit) returns()
func (_PaymentV2 *PaymentV2Transactor) CreateAndAuthorize(opts *bind.TransactOpts, channelDiscriminator [32]byte, initialAmount *big.Int, recipient common.Address, permissions uint16, reservedAmount *big.Int, limit *big.Int) (*types.Transaction, error) {
	return _PaymentV2.contract.Transact(opts, "createAndAuthorize", channelDiscriminator, initialAmount, recipient, permissions, reservedAmount, limit)
}

// CreateAndAuthorize is a paid mutator transaction binding the contract method 0xe333df11.
//
// Solidity: function createAndAuthorize(bytes32 channelDiscriminator, uint256 initialAmount, address recipient, uint16 permissions, uint256 reservedAmount, uint256 limit) returns()
func (_PaymentV2 *PaymentV2Session) CreateAndAuthorize(channelDiscriminator [32]byte, initialAmount *big.Int, recipient common.Address, permissions uint16, reservedAmount *big.Int, limit *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.CreateAndAuthorize(&_PaymentV2.TransactOpts, channelDiscriminator, initialAmount, recipient, permissions, reservedAmount, limit)
}

// CreateAndAuthorize is a paid mutator transaction binding the contract method 0xe333df11.
//
// Solidity: function createAndAuthorize(bytes32 channelDiscriminator, uint256 initialAmount, address recipient, uint16 permissions, uint256 reservedAmount, uint256 limit) returns()
func (_PaymentV2 *PaymentV2TransactorSession) CreateAndAuthorize(channelDiscriminator [32]byte, initialAmount *big.Int, recipient common.Address, permissions uint16, reservedAmount *big.Int, limit *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.CreateAndAuthorize(&_PaymentV2.TransactOpts, channelDiscriminator, initialAmount, recipient, permissions, reservedAmount, limit)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(bytes32 channelId, uint256 amount) returns()
func (_PaymentV2 *PaymentV2Transactor) Deposit(opts *bind.TransactOpts, channelId [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _PaymentV2.contract.Transact(opts, "deposit", channelId, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(bytes32 channelId, uint256 amount) returns()
func (_PaymentV2 *PaymentV2Session) Deposit(channelId [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.Deposit(&_PaymentV2.TransactOpts, channelId, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(bytes32 channelId, uint256 amount) returns()
func (_PaymentV2 *PaymentV2TransactorSession) Deposit(channelId [32]byte, amount *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.Deposit(&_PaymentV2.TransactOpts, channelId, amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x27978c85.
//
// Solidity: function unlock(bytes32 channelId, address recipient) returns()
func (_PaymentV2 *PaymentV2Transactor) Unlock(opts *bind.TransactOpts, channelId [32]byte, recipient common.Address) (*types.Transaction, error) {
	return _PaymentV2.contract.Transact(opts, "unlock", channelId, recipient)
}

// Unlock is a paid mutator transaction binding the contract method 0x27978c85.
//
// Solidity: function unlock(bytes32 channelId, address recipient) returns()
func (_PaymentV2 *PaymentV2Session) Unlock(channelId [32]byte, recipient common.Address) (*types.Transaction, error) {
	return _PaymentV2.Contract.Unlock(&_PaymentV2.TransactOpts, channelId, recipient)
}

// Unlock is a paid mutator transaction binding the contract method 0x27978c85.
//
// Solidity: function unlock(bytes32 channelId, address recipient) returns()
func (_PaymentV2 *PaymentV2TransactorSession) Unlock(channelId [32]byte, recipient common.Address) (*types.Transaction, error) {
	return _PaymentV2.Contract.Unlock(&_PaymentV2.TransactOpts, channelId, recipient)
}

// Withdraw is a paid mutator transaction binding the contract method 0xa6fb97d1.
//
// Solidity: function withdraw(bytes32 channelId, address transferAddress, uint256 transferAmount) returns()
func (_PaymentV2 *PaymentV2Transactor) Withdraw(opts *bind.TransactOpts, channelId [32]byte, transferAddress common.Address, transferAmount *big.Int) (*types.Transaction, error) {
	return _PaymentV2.contract.Transact(opts, "withdraw", channelId, transferAddress, transferAmount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xa6fb97d1.
//
// Solidity: function withdraw(bytes32 channelId, address transferAddress, uint256 transferAmount) returns()
func (_PaymentV2 *PaymentV2Session) Withdraw(channelId [32]byte, transferAddress common.Address, transferAmount *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.Withdraw(&_PaymentV2.TransactOpts, channelId, transferAddress, transferAmount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xa6fb97d1.
//
// Solidity: function withdraw(bytes32 channelId, address transferAddress, uint256 transferAmount) returns()
func (_PaymentV2 *PaymentV2TransactorSession) Withdraw(channelId [32]byte, transferAddress common.Address, transferAmount *big.Int) (*types.Transaction, error) {
	return _PaymentV2.Contract.Withdraw(&_PaymentV2.TransactOpts, channelId, transferAddress, transferAmount)
}

// PaymentV2AuthorizeIterator is returned from FilterAuthorize and is used to iterate over the raw logs and unpacked data for Authorize events raised by the PaymentV2 contract.
type PaymentV2AuthorizeIterator struct {
	Event *PaymentV2Authorize // Event containing the contract specifics and raw log

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
func (it *PaymentV2AuthorizeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentV2Authorize)
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
		it.Event = new(PaymentV2Authorize)
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
func (it *PaymentV2AuthorizeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentV2AuthorizeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentV2Authorize represents a Authorize event raised by the PaymentV2 contract.
type PaymentV2Authorize struct {
	ChannelId      [32]byte
	Recipient      common.Address
	Permissions    uint16
	ReservedAmount *big.Int
	Limit          *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterAuthorize is a free log retrieval operation binding the contract event 0xf96c5460e590096c681d8c8377cc1d01c036860746afa93ec6cb0260b424674e.
//
// Solidity: event Authorize(bytes32 indexed channelId, address indexed recipient, uint16 permissions, uint256 reservedAmount, uint256 limit)
func (_PaymentV2 *PaymentV2Filterer) FilterAuthorize(opts *bind.FilterOpts, channelId [][32]byte, recipient []common.Address) (*PaymentV2AuthorizeIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentV2.contract.FilterLogs(opts, "Authorize", channelIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PaymentV2AuthorizeIterator{contract: _PaymentV2.contract, event: "Authorize", logs: logs, sub: sub}, nil
}

// WatchAuthorize is a free log subscription operation binding the contract event 0xf96c5460e590096c681d8c8377cc1d01c036860746afa93ec6cb0260b424674e.
//
// Solidity: event Authorize(bytes32 indexed channelId, address indexed recipient, uint16 permissions, uint256 reservedAmount, uint256 limit)
func (_PaymentV2 *PaymentV2Filterer) WatchAuthorize(opts *bind.WatchOpts, sink chan<- *PaymentV2Authorize, channelId [][32]byte, recipient []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentV2.contract.WatchLogs(opts, "Authorize", channelIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentV2Authorize)
				if err := _PaymentV2.contract.UnpackLog(event, "Authorize", log); err != nil {
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

// ParseAuthorize is a log parse operation binding the contract event 0xf96c5460e590096c681d8c8377cc1d01c036860746afa93ec6cb0260b424674e.
//
// Solidity: event Authorize(bytes32 indexed channelId, address indexed recipient, uint16 permissions, uint256 reservedAmount, uint256 limit)
func (_PaymentV2 *PaymentV2Filterer) ParseAuthorize(log types.Log) (*PaymentV2Authorize, error) {
	event := new(PaymentV2Authorize)
	if err := _PaymentV2.contract.UnpackLog(event, "Authorize", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentV2DepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the PaymentV2 contract.
type PaymentV2DepositIterator struct {
	Event *PaymentV2Deposit // Event containing the contract specifics and raw log

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
func (it *PaymentV2DepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentV2Deposit)
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
		it.Event = new(PaymentV2Deposit)
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
func (it *PaymentV2DepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentV2DepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentV2Deposit represents a Deposit event raised by the PaymentV2 contract.
type PaymentV2Deposit struct {
	ChannelId [32]byte
	Payer     common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x182fa52899142d44ff5c45a6354d3b3e868d5b07db6a65580b39bd321bdaf8ac.
//
// Solidity: event Deposit(bytes32 indexed channelId, address indexed payer, uint256 amount)
func (_PaymentV2 *PaymentV2Filterer) FilterDeposit(opts *bind.FilterOpts, channelId [][32]byte, payer []common.Address) (*PaymentV2DepositIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PaymentV2.contract.FilterLogs(opts, "Deposit", channelIdRule, payerRule)
	if err != nil {
		return nil, err
	}
	return &PaymentV2DepositIterator{contract: _PaymentV2.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x182fa52899142d44ff5c45a6354d3b3e868d5b07db6a65580b39bd321bdaf8ac.
//
// Solidity: event Deposit(bytes32 indexed channelId, address indexed payer, uint256 amount)
func (_PaymentV2 *PaymentV2Filterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *PaymentV2Deposit, channelId [][32]byte, payer []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var payerRule []interface{}
	for _, payerItem := range payer {
		payerRule = append(payerRule, payerItem)
	}

	logs, sub, err := _PaymentV2.contract.WatchLogs(opts, "Deposit", channelIdRule, payerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentV2Deposit)
				if err := _PaymentV2.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x182fa52899142d44ff5c45a6354d3b3e868d5b07db6a65580b39bd321bdaf8ac.
//
// Solidity: event Deposit(bytes32 indexed channelId, address indexed payer, uint256 amount)
func (_PaymentV2 *PaymentV2Filterer) ParseDeposit(log types.Log) (*PaymentV2Deposit, error) {
	event := new(PaymentV2Deposit)
	if err := _PaymentV2.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentV2UnlockIterator is returned from FilterUnlock and is used to iterate over the raw logs and unpacked data for Unlock events raised by the PaymentV2 contract.
type PaymentV2UnlockIterator struct {
	Event *PaymentV2Unlock // Event containing the contract specifics and raw log

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
func (it *PaymentV2UnlockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentV2Unlock)
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
		it.Event = new(PaymentV2Unlock)
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
func (it *PaymentV2UnlockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentV2UnlockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentV2Unlock represents a Unlock event raised by the PaymentV2 contract.
type PaymentV2Unlock struct {
	ChannelId [32]byte
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnlock is a free log retrieval operation binding the contract event 0x3228d53aec01c740899a507719f566916ac058e41fb369ffae1a9664089c6639.
//
// Solidity: event Unlock(bytes32 indexed channelId, address indexed recipient)
func (_PaymentV2 *PaymentV2Filterer) FilterUnlock(opts *bind.FilterOpts, channelId [][32]byte, recipient []common.Address) (*PaymentV2UnlockIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentV2.contract.FilterLogs(opts, "Unlock", channelIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PaymentV2UnlockIterator{contract: _PaymentV2.contract, event: "Unlock", logs: logs, sub: sub}, nil
}

// WatchUnlock is a free log subscription operation binding the contract event 0x3228d53aec01c740899a507719f566916ac058e41fb369ffae1a9664089c6639.
//
// Solidity: event Unlock(bytes32 indexed channelId, address indexed recipient)
func (_PaymentV2 *PaymentV2Filterer) WatchUnlock(opts *bind.WatchOpts, sink chan<- *PaymentV2Unlock, channelId [][32]byte, recipient []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentV2.contract.WatchLogs(opts, "Unlock", channelIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentV2Unlock)
				if err := _PaymentV2.contract.UnpackLog(event, "Unlock", log); err != nil {
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

// ParseUnlock is a log parse operation binding the contract event 0x3228d53aec01c740899a507719f566916ac058e41fb369ffae1a9664089c6639.
//
// Solidity: event Unlock(bytes32 indexed channelId, address indexed recipient)
func (_PaymentV2 *PaymentV2Filterer) ParseUnlock(log types.Log) (*PaymentV2Unlock, error) {
	event := new(PaymentV2Unlock)
	if err := _PaymentV2.contract.UnpackLog(event, "Unlock", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentV2WithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the PaymentV2 contract.
type PaymentV2WithdrawIterator struct {
	Event *PaymentV2Withdraw // Event containing the contract specifics and raw log

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
func (it *PaymentV2WithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentV2Withdraw)
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
		it.Event = new(PaymentV2Withdraw)
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
func (it *PaymentV2WithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentV2WithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentV2Withdraw represents a Withdraw event raised by the PaymentV2 contract.
type PaymentV2Withdraw struct {
	ChannelId [32]byte
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xe7284ffe0c70ad2f3b0aa15cde1cfe95f736935651a138725b21fd168edc5d6a.
//
// Solidity: event Withdraw(bytes32 indexed channelId, address indexed recipient, uint256 amount)
func (_PaymentV2 *PaymentV2Filterer) FilterWithdraw(opts *bind.FilterOpts, channelId [][32]byte, recipient []common.Address) (*PaymentV2WithdrawIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentV2.contract.FilterLogs(opts, "Withdraw", channelIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PaymentV2WithdrawIterator{contract: _PaymentV2.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xe7284ffe0c70ad2f3b0aa15cde1cfe95f736935651a138725b21fd168edc5d6a.
//
// Solidity: event Withdraw(bytes32 indexed channelId, address indexed recipient, uint256 amount)
func (_PaymentV2 *PaymentV2Filterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *PaymentV2Withdraw, channelId [][32]byte, recipient []common.Address) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PaymentV2.contract.WatchLogs(opts, "Withdraw", channelIdRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentV2Withdraw)
				if err := _PaymentV2.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xe7284ffe0c70ad2f3b0aa15cde1cfe95f736935651a138725b21fd168edc5d6a.
//
// Solidity: event Withdraw(bytes32 indexed channelId, address indexed recipient, uint256 amount)
func (_PaymentV2 *PaymentV2Filterer) ParseWithdraw(log types.Log) (*PaymentV2Withdraw, error) {
	event := new(PaymentV2Withdraw)
	if err := _PaymentV2.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
