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
	Bin: "0x60806040523480156200001157600080fd5b506040518060400160405280600981526020016826b7b1b5aa37b5b2b760b91b815250604051806040016040528060038152602001621352d560ea1b8152508160039081620000619190620002ad565b506004620000708282620002ad565b5050506200008d30670de0b6b3a76400006200009360201b60201c565b620003a1565b6001600160a01b038216620000c35760405163ec442f0560e01b8152600060048201526024015b60405180910390fd5b620000d160008383620000d5565b5050565b6001600160a01b03831662000104578060026000828254620000f8919062000379565b90915550620001789050565b6001600160a01b03831660009081526020819052604090205481811015620001595760405163391434e360e21b81526001600160a01b03851660048201526024810182905260448101839052606401620000ba565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b0382166200019657600280548290039055620001b5565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620001fb91815260200190565b60405180910390a3505050565b634e487b7160e01b600052604160045260246000fd5b600181811c908216806200023357607f821691505b6020821081036200025457634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002a857600081815260208120601f850160051c81016020861015620002835750805b601f850160051c820191505b81811015620002a4578281556001016200028f565b5050505b505050565b81516001600160401b03811115620002c957620002c962000208565b620002e181620002da84546200021e565b846200025a565b602080601f831160018114620003195760008415620003005750858301515b600019600386901b1c1916600185901b178555620002a4565b600085815260208120601f198616915b828110156200034a5788860151825594840194600190910190840162000329565b5085821015620003695787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b808201808211156200039b57634e487b7160e01b600052601160045260246000fd5b92915050565b6107d180620003b16000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c806370a082311161006657806370a082311461011857806395d89b4114610141578063a9059cbb14610149578063c06ed0c21461015c578063dd62ed3e1461017157600080fd5b806306fdde03146100a3578063095ea7b3146100c157806318160ddd146100e457806323b872dd146100f6578063313ce56714610109575b600080fd5b6100ab6101aa565b6040516100b89190610602565b60405180910390f35b6100d46100cf36600461066c565b61023c565b60405190151581526020016100b8565b6002545b6040519081526020016100b8565b6100d4610104366004610696565b610256565b604051601281526020016100b8565b6100e86101263660046106d2565b6001600160a01b031660009081526020819052604090205490565b6100ab61027a565b6100d461015736600461066c565b610289565b61016f61016a3660046106f4565b610297565b005b6100e861017f36600461070d565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6060600380546101b990610740565b80601f01602080910402602001604051908101604052809291908181526020018280546101e590610740565b80156102325780601f1061020757610100808354040283529160200191610232565b820191906000526020600020905b81548152906001019060200180831161021557829003601f168201915b5050505050905090565b60003361024a818585610314565b60019150505b92915050565b600033610264858285610326565b61026f8585856103a4565b506001949350505050565b6060600480546101b990610740565b60003361024a8185856103a4565b3060009081526020819052604090205481106103065760405162461bcd60e51b8152602060048201526024808201527f436f6e747261637420646f6573206e6f74206861766520456e6f75676820746f6044820152636b656e7360e01b60648201526084015b60405180910390fd5b6103113033836103a4565b50565b6103218383836001610403565b505050565b6001600160a01b03838116600090815260016020908152604080832093861683529290522054600019811461039e578181101561038f57604051637dc7a0d960e11b81526001600160a01b038416600482015260248101829052604481018390526064016102fd565b61039e84848484036000610403565b50505050565b6001600160a01b0383166103ce57604051634b637e8f60e11b8152600060048201526024016102fd565b6001600160a01b0382166103f85760405163ec442f0560e01b8152600060048201526024016102fd565b6103218383836104d8565b6001600160a01b03841661042d5760405163e602df0560e01b8152600060048201526024016102fd565b6001600160a01b03831661045757604051634a1406b160e11b8152600060048201526024016102fd565b6001600160a01b038085166000908152600160209081526040808320938716835292905220829055801561039e57826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516104ca91815260200190565b60405180910390a350505050565b6001600160a01b0383166105035780600260008282546104f8919061077a565b909155506105759050565b6001600160a01b038316600090815260208190526040902054818110156105565760405163391434e360e21b81526001600160a01b038516600482015260248101829052604481018390526064016102fd565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b038216610591576002805482900390556105b0565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516105f591815260200190565b60405180910390a3505050565b600060208083528351808285015260005b8181101561062f57858101830151858201604001528201610613565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b038116811461066757600080fd5b919050565b6000806040838503121561067f57600080fd5b61068883610650565b946020939093013593505050565b6000806000606084860312156106ab57600080fd5b6106b484610650565b92506106c260208501610650565b9150604084013590509250925092565b6000602082840312156106e457600080fd5b6106ed82610650565b9392505050565b60006020828403121561070657600080fd5b5035919050565b6000806040838503121561072057600080fd5b61072983610650565b915061073760208401610650565b90509250929050565b600181811c9082168061075457607f821691505b60208210810361077457634e487b7160e01b600052602260045260246000fd5b50919050565b8082018082111561025057634e487b7160e01b600052601160045260246000fdfea264697066735822122000bad3def6545adca488fba426dc7ca32326c64389cff674d833418d626402a164736f6c63430008150033",
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
