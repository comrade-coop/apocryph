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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"NewPrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"NewPriceUpdate\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"acceptNewPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"channels\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"owedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"suggestedPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minAdvanceDuration\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minAdvanceDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"createChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"lockFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"}],\"name\":\"unclockFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newDeadline\",\"type\":\"uint256\"}],\"name\":\"updateDeadline\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"updatePrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"units\",\"type\":\"uint256\"}],\"name\":\"uploadMetrics\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506110d5806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063995cee2611610066578063995cee26146100e6578063a94d2415146100f9578063c7f149261461010c578063f428e28a1461011f578063f940e385146101a357600080fd5b8063165c326f1461009857806331f60646146100ad57806370e6c993146100c05780638fbb2a90146100d3575b600080fd5b6100ab6100a6366004610ed9565b6101b6565b005b6100ab6100bb366004610ed9565b610302565b6100ab6100ce366004610f15565b610484565b6100ab6100e1366004610ed9565b6107fe565b6100ab6100f4366004610f6a565b6108a2565b6100ab610107366004610ed9565b610ad1565b6100ab61011a366004610f6a565b610ccb565b61017261012d366004610f9d565b60006020818152938152604080822085529281528281209093528252902080546001820154600283015460038401546004850154600590950154939492939192909186565b604080519687526020870195909552938501929092526060840152608083015260a082015260c00160405180910390f35b6100ab6101b1366004610f6a565b610d95565b6001600160a01b0380841660009081526020818152604080832033845282528083209386168352928152828220835160c081018552815481526001820154928101929092526002810154938201939093526003830154606082018190526004840154608083015260059093015460a0820152916102339084610ff6565b82516020840151919250906102489083611013565b11156102b75760405162461bcd60e51b815260206004820152603360248201527f4f776564416d6f756e7420697320626967676572207468616e206368616e6e656044820152726c277320617661696c61626c652066756e647360681b60648201526084015b60405180910390fd5b60208201516102c69082611013565b6001600160a01b039586166000908152602081815260408083203384528252808320979098168252959095529490932060010193909355505050565b336000908152602081815260408083206001600160a01b0387811685529083528184209086168452825291829020825160c08101845281548152600182015492810192909252600281015492820192909252600382015460608201526004820154608082015260059091015460a0820181905282116103bc5760405162461bcd60e51b815260206004820152601660248201527513995dc8111958591b1a5b99481d1bdbc81cda1bdc9d60521b60448201526064016102ae565b806040015182116104225760405162461bcd60e51b815260206004820152602a60248201527f4e657720446561646c696e65206973206c657373207468616e2063757272656e6044820152697420646561646c696e6560b01b60648201526084016102ae565b6040808201928352336000908152602081815282822081528282206001600160a01b03969096168252948552208151815592810151600184015590516002830155606081015160038301556080810151600483015560a0015160059091015550565b83600081116104a55760405162461bcd60e51b81526004016102ae90611026565b604051636eb1769f60e11b81523360048201819052306024830152879187906000906001600160a01b0385169063dd62ed3e90604401602060405180830381865afa1580156104f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061051c919061105d565b905081811461057f5760405162461bcd60e51b815260206004820152602960248201527f616c6c6f77616e636520646f6573206e6f74206d617463682073706563696669604482015268195908185b5bdd5b9d60ba1b60648201526084016102ae565b4288116105c15760405162461bcd60e51b815260206004820152601060248201526f111958591b1a5b9948115e1c1a5c995960821b60448201526064016102ae565b336000908152602081815260408083206001600160a01b038f81168552908352818420908e168452825291829020825160c081018452815480825260018301549382019390935260028201549381019390935260038101546060840152600481015460808401526005015460a08301521561067e5760405162461bcd60e51b815260206004820152601760248201527f4368616e6e656c20616c7265616479206372656174656400000000000000000060448201526064016102ae565b6040518060c001604052808b8152602001600081526020018a81526020018881526020016000815260200189815250905080600080336001600160a01b03166001600160a01b0316815260200190815260200160002060008e6001600160a01b03166001600160a01b0316815260200190815260200160002060008d6001600160a01b03166001600160a01b03168152602001908152602001600020600082015181600001556020820151816001015560408201518160020155606082015181600301556080820151816004015560a082015181600501559050508a6001600160a01b03166323b872dd333084600001516040518463ffffffff1660e01b81526004016107ac939291906001600160a01b039384168152919092166020820152604081019190915260600190565b6020604051808303816000875af11580156107cb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107ef9190611076565b50505050505050505050505050565b806000811161081f5760405162461bcd60e51b81526004016102ae90611026565b6001600160a01b038481166000818152602081815260408083203380855290835281842095891680855295835292819020600401879055805193845290830191909152810191909152606081018390527f321c73989a5585f52e30bb8aabd608154125bf4413611a3b4d8b97105d9c6dc09060800160405180910390a150505050565b336000908152602081815260408083206001600160a01b0385811685529083528184209086168452825291829020825160c08101845281548152600182015492810192909252600281015492820183905260038101546060830152600481015460808301526005015460a0820152904210156109605760405162461bcd60e51b815260206004820152601860248201527f446561646c696e65206e6f74207265616368656420796574000000000000000060448201526064016102ae565b6020810151156109e657602081015160405163a9059cbb60e01b81526001600160a01b03848116600483015260248201929092529084169063a9059cbb906044016020604051808303816000875af11580156109c0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109e49190611076565b505b8051610a245760405162461bcd60e51b815260206004820152600d60248201526c115b5c1d1e4810da185b9b995b609a1b60448201526064016102ae565b805160405163a9059cbb60e01b815233600482015260248101919091526001600160a01b0384169063a9059cbb906044016020604051808303816000875af1158015610a74573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a989190611076565b5050336000908152602081815260408083206001600160a01b039485168452825280832094909316825292909252812081815560010155565b8060008111610af25760405162461bcd60e51b81526004016102ae90611026565b336000818152602081815260408083206001600160a01b03808a1685529083528184209088168452825291829020825160c08101845281548152600182015492810192909252600281015492820183905260038101546060830152600481015460808301526005015460a082015286918691904210610ba55760405162461bcd60e51b815260206004820152600f60248201526e10da185b9b995b08115e1c1a5c9959608a1b60448201526064016102ae565b336000908152602081815260408083206001600160a01b038c81168552908352818420908b168452825291829020825160c081018452815480825260018301549382019390935260028201549381019390935260038101546060840152600481015460808401526005015460a0830152610c20908890611013565b336000818152602081815260408083206001600160a01b038f81168552908352818420908e1680855292529182902093909355516323b872dd60e01b81526004810191909152306024820152604481018990526323b872dd906064016020604051808303816000875af1158015610c9b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cbf9190611076565b50505050505050505050565b336000818152602081815260408083206001600160a01b038781168552908352818420908616808552818452828520835160c081018552815481526001820154818701526002820154818601526003820180546060808401918252600485018054608080870182905260059097015460a0870152909355858a5295885297905595518451888152958601979097529284015282019390935290917f8c86180e4992276e0f056bb49eed13ea7b192a0d9e795a7f219b5b281c22d78e910160405180910390a1505050565b6001600160a01b038082166000908152602081815260408083203384528252808320938616835292815290829020825160c08101845281548152600182015492810183905260028201549381019390935260038101546060840152600481015460808401526005015460a0830152610e405760405162461bcd60e51b815260206004820152600e60248201526d05a65726f204f776e6572736869760941b60448201526064016102ae565b602081015160405163a9059cbb60e01b815233600482015260248101919091526001600160a01b0384169063a9059cbb906044016020604051808303816000875af1158015610e93573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610eb79190611076565b50505050565b80356001600160a01b0381168114610ed457600080fd5b919050565b600080600060608486031215610eee57600080fd5b610ef784610ebd565b9250610f0560208501610ebd565b9150604084013590509250925092565b60008060008060008060c08789031215610f2e57600080fd5b610f3787610ebd565b9550610f4560208801610ebd565b95989597505050506040840135936060810135936080820135935060a0909101359150565b60008060408385031215610f7d57600080fd5b610f8683610ebd565b9150610f9460208401610ebd565b90509250929050565b600080600060608486031215610fb257600080fd5b610fbb84610ebd565b9250610fc960208501610ebd565b9150610fd760408501610ebd565b90509250925092565b634e487b7160e01b600052601160045260246000fd5b808202811582820484141761100d5761100d610fe0565b92915050565b8082018082111561100d5761100d610fe0565b6020808252601c908201527f63616e277420616363706574207a65726f20617320612076616c756500000000604082015260600190565b60006020828403121561106f57600080fd5b5051919050565b60006020828403121561108857600080fd5b8151801515811461109857600080fd5b939250505056fea26469706673582212207c87ef023f7379d05bbe413082f17c2232d455ff2761bf81e8773c591157cf3964736f6c63430008150033",
}

// PaymentABI is the input ABI used to generate the binding from.
// Deprecated: Use PaymentMetaData.ABI instead.
var PaymentABI = PaymentMetaData.ABI

// PaymentBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PaymentMetaData.Bin instead.
var PaymentBin = PaymentMetaData.Bin

// DeployPayment deploys a new Ethereum contract, binding an instance of Payment to it.
func DeployPayment(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Payment, error) {
	parsed, err := PaymentMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PaymentBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Payment{PaymentCaller: PaymentCaller{contract: contract}, PaymentTransactor: PaymentTransactor{contract: contract}, PaymentFilterer: PaymentFilterer{contract: contract}}, nil
}

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

// Channels is a free data retrieval call binding the contract method 0xf428e28a.
//
// Solidity: function channels(address , address , address ) view returns(uint256 total, uint256 owedAmount, uint256 deadline, uint256 price, uint256 suggestedPrice, uint256 minAdvanceDuration)
func (_Payment *PaymentCaller) Channels(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 common.Address) (struct {
	Total              *big.Int
	OwedAmount         *big.Int
	Deadline           *big.Int
	Price              *big.Int
	SuggestedPrice     *big.Int
	MinAdvanceDuration *big.Int
}, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "channels", arg0, arg1, arg2)

	outstruct := new(struct {
		Total              *big.Int
		OwedAmount         *big.Int
		Deadline           *big.Int
		Price              *big.Int
		SuggestedPrice     *big.Int
		MinAdvanceDuration *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Total = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.OwedAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Deadline = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Price = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.SuggestedPrice = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.MinAdvanceDuration = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Channels is a free data retrieval call binding the contract method 0xf428e28a.
//
// Solidity: function channels(address , address , address ) view returns(uint256 total, uint256 owedAmount, uint256 deadline, uint256 price, uint256 suggestedPrice, uint256 minAdvanceDuration)
func (_Payment *PaymentSession) Channels(arg0 common.Address, arg1 common.Address, arg2 common.Address) (struct {
	Total              *big.Int
	OwedAmount         *big.Int
	Deadline           *big.Int
	Price              *big.Int
	SuggestedPrice     *big.Int
	MinAdvanceDuration *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2)
}

// Channels is a free data retrieval call binding the contract method 0xf428e28a.
//
// Solidity: function channels(address , address , address ) view returns(uint256 total, uint256 owedAmount, uint256 deadline, uint256 price, uint256 suggestedPrice, uint256 minAdvanceDuration)
func (_Payment *PaymentCallerSession) Channels(arg0 common.Address, arg1 common.Address, arg2 common.Address) (struct {
	Total              *big.Int
	OwedAmount         *big.Int
	Deadline           *big.Int
	Price              *big.Int
	SuggestedPrice     *big.Int
	MinAdvanceDuration *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2)
}

// AcceptNewPrice is a paid mutator transaction binding the contract method 0xc7f14926.
//
// Solidity: function acceptNewPrice(address provider, address token) returns()
func (_Payment *PaymentTransactor) AcceptNewPrice(opts *bind.TransactOpts, provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "acceptNewPrice", provider, token)
}

// AcceptNewPrice is a paid mutator transaction binding the contract method 0xc7f14926.
//
// Solidity: function acceptNewPrice(address provider, address token) returns()
func (_Payment *PaymentSession) AcceptNewPrice(provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.AcceptNewPrice(&_Payment.TransactOpts, provider, token)
}

// AcceptNewPrice is a paid mutator transaction binding the contract method 0xc7f14926.
//
// Solidity: function acceptNewPrice(address provider, address token) returns()
func (_Payment *PaymentTransactorSession) AcceptNewPrice(provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.AcceptNewPrice(&_Payment.TransactOpts, provider, token)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x70e6c993.
//
// Solidity: function createChannel(address provider, address token, uint256 amount, uint256 deadline, uint256 minAdvanceDuration, uint256 price) returns()
func (_Payment *PaymentTransactor) CreateChannel(opts *bind.TransactOpts, provider common.Address, token common.Address, amount *big.Int, deadline *big.Int, minAdvanceDuration *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "createChannel", provider, token, amount, deadline, minAdvanceDuration, price)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x70e6c993.
//
// Solidity: function createChannel(address provider, address token, uint256 amount, uint256 deadline, uint256 minAdvanceDuration, uint256 price) returns()
func (_Payment *PaymentSession) CreateChannel(provider common.Address, token common.Address, amount *big.Int, deadline *big.Int, minAdvanceDuration *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateChannel(&_Payment.TransactOpts, provider, token, amount, deadline, minAdvanceDuration, price)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x70e6c993.
//
// Solidity: function createChannel(address provider, address token, uint256 amount, uint256 deadline, uint256 minAdvanceDuration, uint256 price) returns()
func (_Payment *PaymentTransactorSession) CreateChannel(provider common.Address, token common.Address, amount *big.Int, deadline *big.Int, minAdvanceDuration *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateChannel(&_Payment.TransactOpts, provider, token, amount, deadline, minAdvanceDuration, price)
}

// LockFunds is a paid mutator transaction binding the contract method 0xa94d2415.
//
// Solidity: function lockFunds(address provider, address token, uint256 amount) returns()
func (_Payment *PaymentTransactor) LockFunds(opts *bind.TransactOpts, provider common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "lockFunds", provider, token, amount)
}

// LockFunds is a paid mutator transaction binding the contract method 0xa94d2415.
//
// Solidity: function lockFunds(address provider, address token, uint256 amount) returns()
func (_Payment *PaymentSession) LockFunds(provider common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.LockFunds(&_Payment.TransactOpts, provider, token, amount)
}

// LockFunds is a paid mutator transaction binding the contract method 0xa94d2415.
//
// Solidity: function lockFunds(address provider, address token, uint256 amount) returns()
func (_Payment *PaymentTransactorSession) LockFunds(provider common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.LockFunds(&_Payment.TransactOpts, provider, token, amount)
}

// UnclockFunds is a paid mutator transaction binding the contract method 0x995cee26.
//
// Solidity: function unclockFunds(address token, address provider) returns()
func (_Payment *PaymentTransactor) UnclockFunds(opts *bind.TransactOpts, token common.Address, provider common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "unclockFunds", token, provider)
}

// UnclockFunds is a paid mutator transaction binding the contract method 0x995cee26.
//
// Solidity: function unclockFunds(address token, address provider) returns()
func (_Payment *PaymentSession) UnclockFunds(token common.Address, provider common.Address) (*types.Transaction, error) {
	return _Payment.Contract.UnclockFunds(&_Payment.TransactOpts, token, provider)
}

// UnclockFunds is a paid mutator transaction binding the contract method 0x995cee26.
//
// Solidity: function unclockFunds(address token, address provider) returns()
func (_Payment *PaymentTransactorSession) UnclockFunds(token common.Address, provider common.Address) (*types.Transaction, error) {
	return _Payment.Contract.UnclockFunds(&_Payment.TransactOpts, token, provider)
}

// UpdateDeadline is a paid mutator transaction binding the contract method 0x31f60646.
//
// Solidity: function updateDeadline(address provider, address token, uint256 newDeadline) returns()
func (_Payment *PaymentTransactor) UpdateDeadline(opts *bind.TransactOpts, provider common.Address, token common.Address, newDeadline *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "updateDeadline", provider, token, newDeadline)
}

// UpdateDeadline is a paid mutator transaction binding the contract method 0x31f60646.
//
// Solidity: function updateDeadline(address provider, address token, uint256 newDeadline) returns()
func (_Payment *PaymentSession) UpdateDeadline(provider common.Address, token common.Address, newDeadline *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UpdateDeadline(&_Payment.TransactOpts, provider, token, newDeadline)
}

// UpdateDeadline is a paid mutator transaction binding the contract method 0x31f60646.
//
// Solidity: function updateDeadline(address provider, address token, uint256 newDeadline) returns()
func (_Payment *PaymentTransactorSession) UpdateDeadline(provider common.Address, token common.Address, newDeadline *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UpdateDeadline(&_Payment.TransactOpts, provider, token, newDeadline)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0x8fbb2a90.
//
// Solidity: function updatePrice(address client, address token, uint256 price) returns()
func (_Payment *PaymentTransactor) UpdatePrice(opts *bind.TransactOpts, client common.Address, token common.Address, price *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "updatePrice", client, token, price)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0x8fbb2a90.
//
// Solidity: function updatePrice(address client, address token, uint256 price) returns()
func (_Payment *PaymentSession) UpdatePrice(client common.Address, token common.Address, price *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UpdatePrice(&_Payment.TransactOpts, client, token, price)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0x8fbb2a90.
//
// Solidity: function updatePrice(address client, address token, uint256 price) returns()
func (_Payment *PaymentTransactorSession) UpdatePrice(client common.Address, token common.Address, price *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UpdatePrice(&_Payment.TransactOpts, client, token, price)
}

// UploadMetrics is a paid mutator transaction binding the contract method 0x165c326f.
//
// Solidity: function uploadMetrics(address client, address token, uint256 units) returns()
func (_Payment *PaymentTransactor) UploadMetrics(opts *bind.TransactOpts, client common.Address, token common.Address, units *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "uploadMetrics", client, token, units)
}

// UploadMetrics is a paid mutator transaction binding the contract method 0x165c326f.
//
// Solidity: function uploadMetrics(address client, address token, uint256 units) returns()
func (_Payment *PaymentSession) UploadMetrics(client common.Address, token common.Address, units *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UploadMetrics(&_Payment.TransactOpts, client, token, units)
}

// UploadMetrics is a paid mutator transaction binding the contract method 0x165c326f.
//
// Solidity: function uploadMetrics(address client, address token, uint256 units) returns()
func (_Payment *PaymentTransactorSession) UploadMetrics(client common.Address, token common.Address, units *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UploadMetrics(&_Payment.TransactOpts, client, token, units)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(address token, address client) returns()
func (_Payment *PaymentTransactor) Withdraw(opts *bind.TransactOpts, token common.Address, client common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdraw", token, client)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(address token, address client) returns()
func (_Payment *PaymentSession) Withdraw(token common.Address, client common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, token, client)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf940e385.
//
// Solidity: function withdraw(address token, address client) returns()
func (_Payment *PaymentTransactorSession) Withdraw(token common.Address, client common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, token, client)
}

// PaymentNewPriceIterator is returned from FilterNewPrice and is used to iterate over the raw logs and unpacked data for NewPrice events raised by the Payment contract.
type PaymentNewPriceIterator struct {
	Event *PaymentNewPrice // Event containing the contract specifics and raw log

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
func (it *PaymentNewPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentNewPrice)
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
		it.Event = new(PaymentNewPrice)
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
func (it *PaymentNewPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentNewPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentNewPrice represents a NewPrice event raised by the Payment contract.
type PaymentNewPrice struct {
	Client   common.Address
	Provider common.Address
	Token    common.Address
	Price    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewPrice is a free log retrieval operation binding the contract event 0x8c86180e4992276e0f056bb49eed13ea7b192a0d9e795a7f219b5b281c22d78e.
//
// Solidity: event NewPrice(address client, address provider, address token, uint256 price)
func (_Payment *PaymentFilterer) FilterNewPrice(opts *bind.FilterOpts) (*PaymentNewPriceIterator, error) {

	logs, sub, err := _Payment.contract.FilterLogs(opts, "NewPrice")
	if err != nil {
		return nil, err
	}
	return &PaymentNewPriceIterator{contract: _Payment.contract, event: "NewPrice", logs: logs, sub: sub}, nil
}

// WatchNewPrice is a free log subscription operation binding the contract event 0x8c86180e4992276e0f056bb49eed13ea7b192a0d9e795a7f219b5b281c22d78e.
//
// Solidity: event NewPrice(address client, address provider, address token, uint256 price)
func (_Payment *PaymentFilterer) WatchNewPrice(opts *bind.WatchOpts, sink chan<- *PaymentNewPrice) (event.Subscription, error) {

	logs, sub, err := _Payment.contract.WatchLogs(opts, "NewPrice")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentNewPrice)
				if err := _Payment.contract.UnpackLog(event, "NewPrice", log); err != nil {
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

// ParseNewPrice is a log parse operation binding the contract event 0x8c86180e4992276e0f056bb49eed13ea7b192a0d9e795a7f219b5b281c22d78e.
//
// Solidity: event NewPrice(address client, address provider, address token, uint256 price)
func (_Payment *PaymentFilterer) ParseNewPrice(log types.Log) (*PaymentNewPrice, error) {
	event := new(PaymentNewPrice)
	if err := _Payment.contract.UnpackLog(event, "NewPrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentNewPriceUpdateIterator is returned from FilterNewPriceUpdate and is used to iterate over the raw logs and unpacked data for NewPriceUpdate events raised by the Payment contract.
type PaymentNewPriceUpdateIterator struct {
	Event *PaymentNewPriceUpdate // Event containing the contract specifics and raw log

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
func (it *PaymentNewPriceUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentNewPriceUpdate)
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
		it.Event = new(PaymentNewPriceUpdate)
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
func (it *PaymentNewPriceUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentNewPriceUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentNewPriceUpdate represents a NewPriceUpdate event raised by the Payment contract.
type PaymentNewPriceUpdate struct {
	Client   common.Address
	Provider common.Address
	Token    common.Address
	Price    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewPriceUpdate is a free log retrieval operation binding the contract event 0x321c73989a5585f52e30bb8aabd608154125bf4413611a3b4d8b97105d9c6dc0.
//
// Solidity: event NewPriceUpdate(address client, address provider, address token, uint256 price)
func (_Payment *PaymentFilterer) FilterNewPriceUpdate(opts *bind.FilterOpts) (*PaymentNewPriceUpdateIterator, error) {

	logs, sub, err := _Payment.contract.FilterLogs(opts, "NewPriceUpdate")
	if err != nil {
		return nil, err
	}
	return &PaymentNewPriceUpdateIterator{contract: _Payment.contract, event: "NewPriceUpdate", logs: logs, sub: sub}, nil
}

// WatchNewPriceUpdate is a free log subscription operation binding the contract event 0x321c73989a5585f52e30bb8aabd608154125bf4413611a3b4d8b97105d9c6dc0.
//
// Solidity: event NewPriceUpdate(address client, address provider, address token, uint256 price)
func (_Payment *PaymentFilterer) WatchNewPriceUpdate(opts *bind.WatchOpts, sink chan<- *PaymentNewPriceUpdate) (event.Subscription, error) {

	logs, sub, err := _Payment.contract.WatchLogs(opts, "NewPriceUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentNewPriceUpdate)
				if err := _Payment.contract.UnpackLog(event, "NewPriceUpdate", log); err != nil {
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

// ParseNewPriceUpdate is a log parse operation binding the contract event 0x321c73989a5585f52e30bb8aabd608154125bf4413611a3b4d8b97105d9c6dc0.
//
// Solidity: event NewPriceUpdate(address client, address provider, address token, uint256 price)
func (_Payment *PaymentFilterer) ParseNewPriceUpdate(log types.Log) (*PaymentNewPriceUpdate, error) {
	event := new(PaymentNewPriceUpdate)
	if err := _Payment.contract.UnpackLog(event, "NewPriceUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
