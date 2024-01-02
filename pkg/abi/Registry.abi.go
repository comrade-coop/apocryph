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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"Id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"CpuPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"RamPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"StoragePrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"BandwidthEgressPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"BandwidthIngressPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"Cpumodel\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"TeeType\",\"type\":\"string\"}],\"name\":\"NewPricingTable\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"}],\"name\":\"Subscribed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tableId\",\"type\":\"uint256\"}],\"name\":\"isSubscribed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pricingTableId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"providers\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256[5]\",\"name\":\"Prices\",\"type\":\"uint256[5]\"},{\"internalType\":\"string\",\"name\":\"cpumodel\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"teeType\",\"type\":\"string\"}],\"name\":\"registerPricingTable\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"}],\"name\":\"registerProvider\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tableId\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"subscription\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tableId\",\"type\":\"uint256\"}],\"name\":\"unsubscribe\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMetaData.ABI instead.
var RegistryABI = RegistryMetaData.ABI

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
