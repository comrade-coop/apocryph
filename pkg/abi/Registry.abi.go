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

// RegistryMetaData contains all meta data concerning the Registry contract.
var RegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"isSubscribed\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tableId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pricingTableId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providers\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerPricingTable\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"Prices\",\"type\":\"uint256[5]\",\"internalType\":\"uint256[5]\"},{\"name\":\"cpumodel\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"teeType\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerProvider\",\"inputs\":[{\"name\":\"cid\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"subscribe\",\"inputs\":[{\"name\":\"tableId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"subscription\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unsubscribe\",\"inputs\":[{\"name\":\"tableId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"NewPricingTable\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"Id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"CpuPrice\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"RamPrice\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"StoragePrice\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"BandwidthEgressPrice\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"BandwidthIngressPrice\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"Cpumodel\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"TeeType\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Subscribed\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b50610b22806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80634c5be7791161005b5780634c5be779146101425780634ebdde7d14610155578063ad0b27fb1461016c578063d9ac40221461017f57600080fd5b80630787bc271461008d5780630f574ba7146100b657806336867ae6146100cb57806344d14ff414610114575b600080fd5b6100a061009b366004610633565b610192565b6040516100ad919061069b565b60405180910390f35b6100c96100c43660046106ae565b61022c565b005b6101046100d93660046106c7565b6001600160a01b03919091166000908152600160209081526040808320938352929052205460ff1690565b60405190151581526020016100ad565b6101046101223660046106c7565b600160209081526000928352604080842090915290825290205460ff1681565b6100c96101503660046107bd565b610341565b61015e60025481565b6040519081526020016100ad565b6100c961017a3660046106ae565b6104d4565b6100c961018d366004610880565b6105b3565b600060208190529081526040902080546101ab906108f2565b80601f01602080910402602001604051908101604052809291908181526020018280546101d7906108f2565b80156102245780601f106101f957610100808354040283529160200191610224565b820191906000526020600020905b81548152906001019060200180831161020757829003601f168201915b505050505081565b3360009081526020819052604081208054610246906108f2565b80601f0160208091040260200160405190810160405280929190818152602001828054610272906108f2565b80156102bf5780601f10610294576101008083540402835291602001916102bf565b820191906000526020600020905b8154815290600101906020018083116102a257829003601f168201915b5050505050905060008151116102f05760405162461bcd60e51b81526004016102e79061092c565b60405180910390fd5b336000818152600160208181526040808420878552909152808320805460ff19169092179091555184917f5db0e562b58e88ae25b795493b5a9c538bb02bd38430aa3194dbf8c68f619f5491a35050565b336000908152602081905260408120805461035b906108f2565b80601f0160208091040260200160405190810160405280929190818152602001828054610387906108f2565b80156103d45780601f106103a9576101008083540402835291602001916103d4565b820191906000526020600020905b8154815290600101906020018083116103b757829003601f168201915b5050505050905060008151116103fc5760405162461bcd60e51b81526004016102e79061092c565b835160208501516040860151606087015160808801516002805490600061042283610963565b91905055506002548a6001600160a01b03167fbc10dcb42f82b1cc23891ab8033eb256ab9f9b75d3cb52bc3064807a79d7faf187878787878f8f60405161046f979695949392919061098a565b60405180910390a33360008181526001602081815260408084206002805486529252808420805460ff191690931790925554905190917f5db0e562b58e88ae25b795493b5a9c538bb02bd38430aa3194dbf8c68f619f5491a350505050505050505050565b33600090815260208190526040812080546104ee906108f2565b80601f016020809104026020016040519081016040528092919081815260200182805461051a906108f2565b80156105675780601f1061053c57610100808354040283529160200191610567565b820191906000526020600020905b81548152906001019060200180831161054a57829003601f168201915b50505050509050600081511161058f5760405162461bcd60e51b81526004016102e79061092c565b5033600090815260016020908152604080832093835292905220805460ff19169055565b806105f85760405162461bcd60e51b8152602060048201526015602482015274636964206d757374206e6f7420626520656d70747960581b60448201526064016102e7565b336000908152602081905260409020610612828483610a2b565b505050565b80356001600160a01b038116811461062e57600080fd5b919050565b60006020828403121561064557600080fd5b61064e82610617565b9392505050565b6000815180845260005b8181101561067b5760208185018101518683018201520161065f565b506000602082860101526020601f19601f83011685010191505092915050565b60208152600061064e6020830184610655565b6000602082840312156106c057600080fd5b5035919050565b600080604083850312156106da57600080fd5b6106e383610617565b946020939093013593505050565b634e487b7160e01b600052604160045260246000fd5b60405160a0810167ffffffffffffffff8111828210171561072a5761072a6106f1565b60405290565b600082601f83011261074157600080fd5b813567ffffffffffffffff8082111561075c5761075c6106f1565b604051601f8301601f19908116603f01168101908282118183101715610784576107846106f1565b8160405283815286602085880101111561079d57600080fd5b836020870160208301376000602085830101528094505050505092915050565b60008060008061010085870312156107d457600080fd5b6107dd85610617565b9350602086603f8701126107f057600080fd5b6107f8610707565b8060c088018981111561080a57600080fd5b602089015b81811015610826578035845292840192840161080f565b5090955035915067ffffffffffffffff90508082111561084557600080fd5b61085188838901610730565b935060e087013591508082111561086757600080fd5b5061087487828801610730565b91505092959194509250565b6000806020838503121561089357600080fd5b823567ffffffffffffffff808211156108ab57600080fd5b818501915085601f8301126108bf57600080fd5b8135818111156108ce57600080fd5b8660208285010111156108e057600080fd5b60209290920196919550909350505050565b600181811c9082168061090657607f821691505b60208210810361092657634e487b7160e01b600052602260045260246000fd5b50919050565b6020808252601b908201527f50726f7669646572204d75737420626520726567697374657265640000000000604082015260600190565b60006001820161098357634e487b7160e01b600052601160045260246000fd5b5060010190565b87815286602082015285604082015284606082015283608082015260e060a082015260006109bb60e0830185610655565b82810360c08401526109cd8185610655565b9a9950505050505050505050565b601f821115610612576000816000526020600020601f850160051c81016020861015610a045750805b601f850160051c820191505b81811015610a2357828155600101610a10565b505050505050565b67ffffffffffffffff831115610a4357610a436106f1565b610a5783610a5183546108f2565b836109db565b6000601f841160018114610a8b5760008515610a735750838201355b600019600387901b1c1916600186901b178355610ae5565b600083815260209020601f19861690835b82811015610abc5786850135825560209485019460019092019101610a9c565b5086821015610ad95760001960f88860031b161c19848701351681555b505060018560011b0183555b505050505056fea2646970667358221220b6ad7e7c376c6eaa8dd37ac8be3f528cd11dc0d1a6f6b19ab90cd215ca4c4f1d64736f6c63430008170033",
}

// RegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMetaData.ABI instead.
var RegistryABI = RegistryMetaData.ABI

// RegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RegistryMetaData.Bin instead.
var RegistryBin = RegistryMetaData.Bin

// DeployRegistry deploys a new Ethereum contract, binding an instance of Registry to it.
func DeployRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Registry, error) {
	parsed, err := RegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// IsSubscribed is a free data retrieval call binding the contract method 0x36867ae6.
//
// Solidity: function isSubscribed(address provider, uint256 tableId) view returns(bool)
func (_Registry *RegistryCaller) IsSubscribed(opts *bind.CallOpts, provider common.Address, tableId *big.Int) (bool, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "isSubscribed", provider, tableId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSubscribed is a free data retrieval call binding the contract method 0x36867ae6.
//
// Solidity: function isSubscribed(address provider, uint256 tableId) view returns(bool)
func (_Registry *RegistrySession) IsSubscribed(provider common.Address, tableId *big.Int) (bool, error) {
	return _Registry.Contract.IsSubscribed(&_Registry.CallOpts, provider, tableId)
}

// IsSubscribed is a free data retrieval call binding the contract method 0x36867ae6.
//
// Solidity: function isSubscribed(address provider, uint256 tableId) view returns(bool)
func (_Registry *RegistryCallerSession) IsSubscribed(provider common.Address, tableId *big.Int) (bool, error) {
	return _Registry.Contract.IsSubscribed(&_Registry.CallOpts, provider, tableId)
}

// PricingTableId is a free data retrieval call binding the contract method 0x4ebdde7d.
//
// Solidity: function pricingTableId() view returns(uint256)
func (_Registry *RegistryCaller) PricingTableId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "pricingTableId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PricingTableId is a free data retrieval call binding the contract method 0x4ebdde7d.
//
// Solidity: function pricingTableId() view returns(uint256)
func (_Registry *RegistrySession) PricingTableId() (*big.Int, error) {
	return _Registry.Contract.PricingTableId(&_Registry.CallOpts)
}

// PricingTableId is a free data retrieval call binding the contract method 0x4ebdde7d.
//
// Solidity: function pricingTableId() view returns(uint256)
func (_Registry *RegistryCallerSession) PricingTableId() (*big.Int, error) {
	return _Registry.Contract.PricingTableId(&_Registry.CallOpts)
}

// Providers is a free data retrieval call binding the contract method 0x0787bc27.
//
// Solidity: function providers(address ) view returns(string)
func (_Registry *RegistryCaller) Providers(opts *bind.CallOpts, arg0 common.Address) (string, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "providers", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Providers is a free data retrieval call binding the contract method 0x0787bc27.
//
// Solidity: function providers(address ) view returns(string)
func (_Registry *RegistrySession) Providers(arg0 common.Address) (string, error) {
	return _Registry.Contract.Providers(&_Registry.CallOpts, arg0)
}

// Providers is a free data retrieval call binding the contract method 0x0787bc27.
//
// Solidity: function providers(address ) view returns(string)
func (_Registry *RegistryCallerSession) Providers(arg0 common.Address) (string, error) {
	return _Registry.Contract.Providers(&_Registry.CallOpts, arg0)
}

// Subscription is a free data retrieval call binding the contract method 0x44d14ff4.
//
// Solidity: function subscription(address , uint256 ) view returns(bool)
func (_Registry *RegistryCaller) Subscription(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "subscription", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Subscription is a free data retrieval call binding the contract method 0x44d14ff4.
//
// Solidity: function subscription(address , uint256 ) view returns(bool)
func (_Registry *RegistrySession) Subscription(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _Registry.Contract.Subscription(&_Registry.CallOpts, arg0, arg1)
}

// Subscription is a free data retrieval call binding the contract method 0x44d14ff4.
//
// Solidity: function subscription(address , uint256 ) view returns(bool)
func (_Registry *RegistryCallerSession) Subscription(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _Registry.Contract.Subscription(&_Registry.CallOpts, arg0, arg1)
}

// RegisterPricingTable is a paid mutator transaction binding the contract method 0x4c5be779.
//
// Solidity: function registerPricingTable(address token, uint256[5] Prices, string cpumodel, string teeType) returns()
func (_Registry *RegistryTransactor) RegisterPricingTable(opts *bind.TransactOpts, token common.Address, Prices [5]*big.Int, cpumodel string, teeType string) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "registerPricingTable", token, Prices, cpumodel, teeType)
}

// RegisterPricingTable is a paid mutator transaction binding the contract method 0x4c5be779.
//
// Solidity: function registerPricingTable(address token, uint256[5] Prices, string cpumodel, string teeType) returns()
func (_Registry *RegistrySession) RegisterPricingTable(token common.Address, Prices [5]*big.Int, cpumodel string, teeType string) (*types.Transaction, error) {
	return _Registry.Contract.RegisterPricingTable(&_Registry.TransactOpts, token, Prices, cpumodel, teeType)
}

// RegisterPricingTable is a paid mutator transaction binding the contract method 0x4c5be779.
//
// Solidity: function registerPricingTable(address token, uint256[5] Prices, string cpumodel, string teeType) returns()
func (_Registry *RegistryTransactorSession) RegisterPricingTable(token common.Address, Prices [5]*big.Int, cpumodel string, teeType string) (*types.Transaction, error) {
	return _Registry.Contract.RegisterPricingTable(&_Registry.TransactOpts, token, Prices, cpumodel, teeType)
}

// RegisterProvider is a paid mutator transaction binding the contract method 0xd9ac4022.
//
// Solidity: function registerProvider(string cid) returns()
func (_Registry *RegistryTransactor) RegisterProvider(opts *bind.TransactOpts, cid string) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "registerProvider", cid)
}

// RegisterProvider is a paid mutator transaction binding the contract method 0xd9ac4022.
//
// Solidity: function registerProvider(string cid) returns()
func (_Registry *RegistrySession) RegisterProvider(cid string) (*types.Transaction, error) {
	return _Registry.Contract.RegisterProvider(&_Registry.TransactOpts, cid)
}

// RegisterProvider is a paid mutator transaction binding the contract method 0xd9ac4022.
//
// Solidity: function registerProvider(string cid) returns()
func (_Registry *RegistryTransactorSession) RegisterProvider(cid string) (*types.Transaction, error) {
	return _Registry.Contract.RegisterProvider(&_Registry.TransactOpts, cid)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(uint256 tableId) returns()
func (_Registry *RegistryTransactor) Subscribe(opts *bind.TransactOpts, tableId *big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "subscribe", tableId)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(uint256 tableId) returns()
func (_Registry *RegistrySession) Subscribe(tableId *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Subscribe(&_Registry.TransactOpts, tableId)
}

// Subscribe is a paid mutator transaction binding the contract method 0x0f574ba7.
//
// Solidity: function subscribe(uint256 tableId) returns()
func (_Registry *RegistryTransactorSession) Subscribe(tableId *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Subscribe(&_Registry.TransactOpts, tableId)
}

// Unsubscribe is a paid mutator transaction binding the contract method 0xad0b27fb.
//
// Solidity: function unsubscribe(uint256 tableId) returns()
func (_Registry *RegistryTransactor) Unsubscribe(opts *bind.TransactOpts, tableId *big.Int) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "unsubscribe", tableId)
}

// Unsubscribe is a paid mutator transaction binding the contract method 0xad0b27fb.
//
// Solidity: function unsubscribe(uint256 tableId) returns()
func (_Registry *RegistrySession) Unsubscribe(tableId *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Unsubscribe(&_Registry.TransactOpts, tableId)
}

// Unsubscribe is a paid mutator transaction binding the contract method 0xad0b27fb.
//
// Solidity: function unsubscribe(uint256 tableId) returns()
func (_Registry *RegistryTransactorSession) Unsubscribe(tableId *big.Int) (*types.Transaction, error) {
	return _Registry.Contract.Unsubscribe(&_Registry.TransactOpts, tableId)
}

// RegistryNewPricingTableIterator is returned from FilterNewPricingTable and is used to iterate over the raw logs and unpacked data for NewPricingTable events raised by the Registry contract.
type RegistryNewPricingTableIterator struct {
	Event *RegistryNewPricingTable // Event containing the contract specifics and raw log

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
func (it *RegistryNewPricingTableIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistryNewPricingTable)
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
		it.Event = new(RegistryNewPricingTable)
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
func (it *RegistryNewPricingTableIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistryNewPricingTableIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistryNewPricingTable represents a NewPricingTable event raised by the Registry contract.
type RegistryNewPricingTable struct {
	Token                 common.Address
	Id                    *big.Int
	CpuPrice              *big.Int
	RamPrice              *big.Int
	StoragePrice          *big.Int
	BandwidthEgressPrice  *big.Int
	BandwidthIngressPrice *big.Int
	Cpumodel              string
	TeeType               string
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterNewPricingTable is a free log retrieval operation binding the contract event 0xbc10dcb42f82b1cc23891ab8033eb256ab9f9b75d3cb52bc3064807a79d7faf1.
//
// Solidity: event NewPricingTable(address indexed token, uint256 indexed Id, uint256 CpuPrice, uint256 RamPrice, uint256 StoragePrice, uint256 BandwidthEgressPrice, uint256 BandwidthIngressPrice, string Cpumodel, string TeeType)
func (_Registry *RegistryFilterer) FilterNewPricingTable(opts *bind.FilterOpts, token []common.Address, Id []*big.Int) (*RegistryNewPricingTableIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var IdRule []interface{}
	for _, IdItem := range Id {
		IdRule = append(IdRule, IdItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "NewPricingTable", tokenRule, IdRule)
	if err != nil {
		return nil, err
	}
	return &RegistryNewPricingTableIterator{contract: _Registry.contract, event: "NewPricingTable", logs: logs, sub: sub}, nil
}

// WatchNewPricingTable is a free log subscription operation binding the contract event 0xbc10dcb42f82b1cc23891ab8033eb256ab9f9b75d3cb52bc3064807a79d7faf1.
//
// Solidity: event NewPricingTable(address indexed token, uint256 indexed Id, uint256 CpuPrice, uint256 RamPrice, uint256 StoragePrice, uint256 BandwidthEgressPrice, uint256 BandwidthIngressPrice, string Cpumodel, string TeeType)
func (_Registry *RegistryFilterer) WatchNewPricingTable(opts *bind.WatchOpts, sink chan<- *RegistryNewPricingTable, token []common.Address, Id []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var IdRule []interface{}
	for _, IdItem := range Id {
		IdRule = append(IdRule, IdItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "NewPricingTable", tokenRule, IdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistryNewPricingTable)
				if err := _Registry.contract.UnpackLog(event, "NewPricingTable", log); err != nil {
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

// ParseNewPricingTable is a log parse operation binding the contract event 0xbc10dcb42f82b1cc23891ab8033eb256ab9f9b75d3cb52bc3064807a79d7faf1.
//
// Solidity: event NewPricingTable(address indexed token, uint256 indexed Id, uint256 CpuPrice, uint256 RamPrice, uint256 StoragePrice, uint256 BandwidthEgressPrice, uint256 BandwidthIngressPrice, string Cpumodel, string TeeType)
func (_Registry *RegistryFilterer) ParseNewPricingTable(log types.Log) (*RegistryNewPricingTable, error) {
	event := new(RegistryNewPricingTable)
	if err := _Registry.contract.UnpackLog(event, "NewPricingTable", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RegistrySubscribedIterator is returned from FilterSubscribed and is used to iterate over the raw logs and unpacked data for Subscribed events raised by the Registry contract.
type RegistrySubscribedIterator struct {
	Event *RegistrySubscribed // Event containing the contract specifics and raw log

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
func (it *RegistrySubscribedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RegistrySubscribed)
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
		it.Event = new(RegistrySubscribed)
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
func (it *RegistrySubscribedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RegistrySubscribedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RegistrySubscribed represents a Subscribed event raised by the Registry contract.
type RegistrySubscribed struct {
	Id       *big.Int
	Provider common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubscribed is a free log retrieval operation binding the contract event 0x5db0e562b58e88ae25b795493b5a9c538bb02bd38430aa3194dbf8c68f619f54.
//
// Solidity: event Subscribed(uint256 indexed id, address indexed provider)
func (_Registry *RegistryFilterer) FilterSubscribed(opts *bind.FilterOpts, id []*big.Int, provider []common.Address) (*RegistrySubscribedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Registry.contract.FilterLogs(opts, "Subscribed", idRule, providerRule)
	if err != nil {
		return nil, err
	}
	return &RegistrySubscribedIterator{contract: _Registry.contract, event: "Subscribed", logs: logs, sub: sub}, nil
}

// WatchSubscribed is a free log subscription operation binding the contract event 0x5db0e562b58e88ae25b795493b5a9c538bb02bd38430aa3194dbf8c68f619f54.
//
// Solidity: event Subscribed(uint256 indexed id, address indexed provider)
func (_Registry *RegistryFilterer) WatchSubscribed(opts *bind.WatchOpts, sink chan<- *RegistrySubscribed, id []*big.Int, provider []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Registry.contract.WatchLogs(opts, "Subscribed", idRule, providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RegistrySubscribed)
				if err := _Registry.contract.UnpackLog(event, "Subscribed", log); err != nil {
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

// ParseSubscribed is a log parse operation binding the contract event 0x5db0e562b58e88ae25b795493b5a9c538bb02bd38430aa3194dbf8c68f619f54.
//
// Solidity: event Subscribed(uint256 indexed id, address indexed provider)
func (_Registry *RegistryFilterer) ParseSubscribed(log types.Log) (*RegistrySubscribed, error) {
	event := new(RegistrySubscribed)
	if err := _Registry.contract.UnpackLog(event, "Subscribed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
