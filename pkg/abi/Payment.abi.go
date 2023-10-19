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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AmountRequired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChannelLocked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"ChannelClosed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockedAt\",\"type\":\"uint256\"}],\"name\":\"UnlockTimerStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockedAmount\",\"type\":\"uint256\"}],\"name\":\"Unlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawnAmount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"available\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"channels\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"investedByPublisher\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawnByProvider\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockedAt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"closeChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialAmount\",\"type\":\"uint256\"}],\"name\":\"createChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transferAddress\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"withdrawUnlocked\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"totalWithdrawlAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transferAddress\",\"type\":\"address\"}],\"name\":\"withdrawUpTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610d00806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c80638340f549116100665780638340f549146100f9578063bed630411461010c578063e0f33dd41461011f578063e674a0bd14610132578063f428e28a1461014557600080fd5b80630e917f76146100985780631086b75d146100ad57806343ea3f36146100c057806367d07056146100d3575b600080fd5b6100ab6100a6366004610ac9565b6101aa565b005b6100ab6100bb366004610b1c565b6101f8565b6100ab6100ce366004610b62565b6102e5565b6100e66100e1366004610b9b565b6103e2565b6040519081526020015b60405180910390f35b6100ab610107366004610be6565b61042d565b6100ab61011a366004610b62565b6104ec565b6100ab61012d366004610ac9565b610666565b6100ab610140366004610b62565b61078e565b61018a610153366004610b9b565b6000602081815293815260408082208552928152828120909352825290208054600182015460028301546003909301549192909184565b6040805194855260208501939093529183015260608201526080016100f0565b6001600160a01b038085166000908152602081815260408083203384528252808320938716835292905220600101546101f290859085906101ec908690610c3d565b84610666565b50505050565b806000036102195760405163d43e607560e01b815260040160405180910390fd5b336000818152602081815260408083206001600160a01b038981168552908352818420908816845290915290208054156102665760405163119b4fd360e11b815260040160405180910390fd5b60018101541561027857610278610c50565b828155600281018490556040518381526001600160a01b0386811691888216918516907f4174a9435a04d04d274c76779cad136a41fde6937c56241c09ab9d3c7064a1a99060200160405180910390a46102dd6001600160a01b038616333086610873565b505050505050565b336000818152602081815260408083206001600160a01b03878116855290835281842090861684529091529020600381015415806103265750806003015442105b1561034457604051635652d88760e01b815260040160405180910390fd5b6001810154815460009161035791610c66565b90508060000361037a5760405163d43e607560e01b815260040160405180910390fd5b600182015482556040518181526001600160a01b0385811691878216918616907f6c5b0a34e5e78b423db3d15d7f4f72f3beb727025c0e535bf2bac21d0227c4819060200160405180910390a46103db6001600160a01b03851684836108da565b5050505050565b6001600160a01b03808416600090815260208181526040808320868516845282528083209385168352929052908120600181015481546104229190610c66565b9150505b9392505050565b8060000361044e5760405163d43e607560e01b815260040160405180910390fd5b336000818152602081815260408083206001600160a01b038881168552908352818420908716845290915290208054610488908490610c3d565b8155600060038201556040518381526001600160a01b0385811691878216918516907f4174a9435a04d04d274c76779cad136a41fde6937c56241c09ab9d3c7064a1a99060200160405180910390a46103db6001600160a01b038516333086610873565b336000818152602081815260408083206001600160a01b038781168552908352818420908616845290915290206003810154158061052d5750806003015442105b1561054b57604051635652d88760e01b815260040160405180910390fd5b6001810154815460009161055e91610c66565b6001600160a01b038085166000908152602081815260408083208a85168452825280832093891683529290529081208181556001810182905560028101829055600301559050801561060257836001600160a01b0316856001600160a01b0316846001600160a01b03167f6c5b0a34e5e78b423db3d15d7f4f72f3beb727025c0e535bf2bac21d0227c481846040516105f991815260200190565b60405180910390a45b836001600160a01b0316856001600160a01b0316846001600160a01b03167f4509c78f8c652633a480f64f4f56b8187ec70d118e774af3139ec088e6f3ea9160405160405180910390a480156103db576103db6001600160a01b03851684836108da565b6001600160a01b0381166106775750335b6001600160a01b0380851660009081526020818152604080832033808552908352818420948816845293909152902080548411156106c85760405163356680b760e01b815260040160405180910390fd5b806001015484116106ec5760405163d43e607560e01b815260040160405180910390fd5b60008160010154856106fe9190610c66565b9050848260010181905550856001600160a01b0316836001600160a01b0316886001600160a01b03167fa4195c37c2947bbe89165f03e320b6903116f0b10d8cfdb522330f7ce6f9fa248460405161075891815260200190565b60405180910390a4600382015415610771574260038301555b6107856001600160a01b03871685836108da565b50505050505050565b336000818152602081815260408083206001600160a01b0387811685529083528184209086168452909152812080549091036107dd5760405163b0ce759160e01b815260040160405180910390fd5b60008160020154426107ef9190610c3d565b90508160030154600014806108075750808260030154105b156103db57808260030181905550836001600160a01b0316856001600160a01b0316846001600160a01b03167f9c21c91a443e1aeab0e24df34b134bf134e2c0c9fecd918faa11306a0adfa62e8460405161086491815260200190565b60405180910390a45050505050565b6040516001600160a01b0384811660248301528381166044830152606482018390526101f29186918216906323b872dd906084015b604051602081830303815290604052915060e01b6020820180516001600160e01b038381831617835250505050610910565b6040516001600160a01b0383811660248301526044820183905261090b91859182169063a9059cbb906064016108a8565b505050565b60006109256001600160a01b03841683610978565b9050805160001415801561094a5750808060200190518101906109489190610c79565b155b1561090b57604051635274afe760e01b81526001600160a01b03841660048201526024015b60405180910390fd5b60606109868383600061098f565b90505b92915050565b6060814710156109b45760405163cd78605960e01b815230600482015260240161096f565b600080856001600160a01b031684866040516109d09190610c9b565b60006040518083038185875af1925050503d8060008114610a0d576040519150601f19603f3d011682016040523d82523d6000602084013e610a12565b606091505b5091509150610a22868383610a2c565b9695505050505050565b606082610a4157610a3c82610a88565b610426565b8151158015610a5857506001600160a01b0384163b155b15610a8157604051639996b31560e01b81526001600160a01b038516600482015260240161096f565b5080610426565b805115610a985780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b50565b6001600160a01b0381168114610ab157600080fd5b60008060008060808587031215610adf57600080fd5b8435610aea81610ab4565b93506020850135610afa81610ab4565b9250604085013591506060850135610b1181610ab4565b939692955090935050565b60008060008060808587031215610b3257600080fd5b8435610b3d81610ab4565b93506020850135610b4d81610ab4565b93969395505050506040820135916060013590565b60008060408385031215610b7557600080fd5b8235610b8081610ab4565b91506020830135610b9081610ab4565b809150509250929050565b600080600060608486031215610bb057600080fd5b8335610bbb81610ab4565b92506020840135610bcb81610ab4565b91506040840135610bdb81610ab4565b809150509250925092565b600080600060608486031215610bfb57600080fd5b8335610c0681610ab4565b92506020840135610c1681610ab4565b929592945050506040919091013590565b634e487b7160e01b600052601160045260246000fd5b8082018082111561098957610989610c27565b634e487b7160e01b600052600160045260246000fd5b8181038181111561098957610989610c27565b600060208284031215610c8b57600080fd5b8151801515811461042657600080fd5b6000825160005b81811015610cbc5760208186018101518583015201610ca2565b50600092019182525091905056fea26469706673582212202619efd1572139080fe32564a4167eb3fa89a58f251618394e1b109b2a2754b264736f6c63430008150033",
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

// Available is a free data retrieval call binding the contract method 0x67d07056.
//
// Solidity: function available(address publisher, address provider, address token) view returns(uint256)
func (_Payment *PaymentCaller) Available(opts *bind.CallOpts, publisher common.Address, provider common.Address, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "available", publisher, provider, token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Available is a free data retrieval call binding the contract method 0x67d07056.
//
// Solidity: function available(address publisher, address provider, address token) view returns(uint256)
func (_Payment *PaymentSession) Available(publisher common.Address, provider common.Address, token common.Address) (*big.Int, error) {
	return _Payment.Contract.Available(&_Payment.CallOpts, publisher, provider, token)
}

// Available is a free data retrieval call binding the contract method 0x67d07056.
//
// Solidity: function available(address publisher, address provider, address token) view returns(uint256)
func (_Payment *PaymentCallerSession) Available(publisher common.Address, provider common.Address, token common.Address) (*big.Int, error) {
	return _Payment.Contract.Available(&_Payment.CallOpts, publisher, provider, token)
}

// Channels is a free data retrieval call binding the contract method 0xf428e28a.
//
// Solidity: function channels(address , address , address ) view returns(uint256 investedByPublisher, uint256 withdrawnByProvider, uint256 unlockTime, uint256 unlockedAt)
func (_Payment *PaymentCaller) Channels(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 common.Address) (struct {
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

// Channels is a free data retrieval call binding the contract method 0xf428e28a.
//
// Solidity: function channels(address , address , address ) view returns(uint256 investedByPublisher, uint256 withdrawnByProvider, uint256 unlockTime, uint256 unlockedAt)
func (_Payment *PaymentSession) Channels(arg0 common.Address, arg1 common.Address, arg2 common.Address) (struct {
	InvestedByPublisher *big.Int
	WithdrawnByProvider *big.Int
	UnlockTime          *big.Int
	UnlockedAt          *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2)
}

// Channels is a free data retrieval call binding the contract method 0xf428e28a.
//
// Solidity: function channels(address , address , address ) view returns(uint256 investedByPublisher, uint256 withdrawnByProvider, uint256 unlockTime, uint256 unlockedAt)
func (_Payment *PaymentCallerSession) Channels(arg0 common.Address, arg1 common.Address, arg2 common.Address) (struct {
	InvestedByPublisher *big.Int
	WithdrawnByProvider *big.Int
	UnlockTime          *big.Int
	UnlockedAt          *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2)
}

// CloseChannel is a paid mutator transaction binding the contract method 0xbed63041.
//
// Solidity: function closeChannel(address provider, address token) returns()
func (_Payment *PaymentTransactor) CloseChannel(opts *bind.TransactOpts, provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "closeChannel", provider, token)
}

// CloseChannel is a paid mutator transaction binding the contract method 0xbed63041.
//
// Solidity: function closeChannel(address provider, address token) returns()
func (_Payment *PaymentSession) CloseChannel(provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.CloseChannel(&_Payment.TransactOpts, provider, token)
}

// CloseChannel is a paid mutator transaction binding the contract method 0xbed63041.
//
// Solidity: function closeChannel(address provider, address token) returns()
func (_Payment *PaymentTransactorSession) CloseChannel(provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.CloseChannel(&_Payment.TransactOpts, provider, token)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x1086b75d.
//
// Solidity: function createChannel(address provider, address token, uint256 unlockTime, uint256 initialAmount) returns()
func (_Payment *PaymentTransactor) CreateChannel(opts *bind.TransactOpts, provider common.Address, token common.Address, unlockTime *big.Int, initialAmount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "createChannel", provider, token, unlockTime, initialAmount)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x1086b75d.
//
// Solidity: function createChannel(address provider, address token, uint256 unlockTime, uint256 initialAmount) returns()
func (_Payment *PaymentSession) CreateChannel(provider common.Address, token common.Address, unlockTime *big.Int, initialAmount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateChannel(&_Payment.TransactOpts, provider, token, unlockTime, initialAmount)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x1086b75d.
//
// Solidity: function createChannel(address provider, address token, uint256 unlockTime, uint256 initialAmount) returns()
func (_Payment *PaymentTransactorSession) CreateChannel(provider common.Address, token common.Address, unlockTime *big.Int, initialAmount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateChannel(&_Payment.TransactOpts, provider, token, unlockTime, initialAmount)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(address provider, address token, uint256 amount) returns()
func (_Payment *PaymentTransactor) Deposit(opts *bind.TransactOpts, provider common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "deposit", provider, token, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(address provider, address token, uint256 amount) returns()
func (_Payment *PaymentSession) Deposit(provider common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.Deposit(&_Payment.TransactOpts, provider, token, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(address provider, address token, uint256 amount) returns()
func (_Payment *PaymentTransactorSession) Deposit(provider common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.Deposit(&_Payment.TransactOpts, provider, token, amount)
}

// Unlock is a paid mutator transaction binding the contract method 0xe674a0bd.
//
// Solidity: function unlock(address provider, address token) returns()
func (_Payment *PaymentTransactor) Unlock(opts *bind.TransactOpts, provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "unlock", provider, token)
}

// Unlock is a paid mutator transaction binding the contract method 0xe674a0bd.
//
// Solidity: function unlock(address provider, address token) returns()
func (_Payment *PaymentSession) Unlock(provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Unlock(&_Payment.TransactOpts, provider, token)
}

// Unlock is a paid mutator transaction binding the contract method 0xe674a0bd.
//
// Solidity: function unlock(address provider, address token) returns()
func (_Payment *PaymentTransactorSession) Unlock(provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Unlock(&_Payment.TransactOpts, provider, token)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0e917f76.
//
// Solidity: function withdraw(address publisher, address token, uint256 amount, address transferAddress) returns()
func (_Payment *PaymentTransactor) Withdraw(opts *bind.TransactOpts, publisher common.Address, token common.Address, amount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdraw", publisher, token, amount, transferAddress)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0e917f76.
//
// Solidity: function withdraw(address publisher, address token, uint256 amount, address transferAddress) returns()
func (_Payment *PaymentSession) Withdraw(publisher common.Address, token common.Address, amount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, publisher, token, amount, transferAddress)
}

// Withdraw is a paid mutator transaction binding the contract method 0x0e917f76.
//
// Solidity: function withdraw(address publisher, address token, uint256 amount, address transferAddress) returns()
func (_Payment *PaymentTransactorSession) Withdraw(publisher common.Address, token common.Address, amount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, publisher, token, amount, transferAddress)
}

// WithdrawUnlocked is a paid mutator transaction binding the contract method 0x43ea3f36.
//
// Solidity: function withdrawUnlocked(address provider, address token) returns()
func (_Payment *PaymentTransactor) WithdrawUnlocked(opts *bind.TransactOpts, provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdrawUnlocked", provider, token)
}

// WithdrawUnlocked is a paid mutator transaction binding the contract method 0x43ea3f36.
//
// Solidity: function withdrawUnlocked(address provider, address token) returns()
func (_Payment *PaymentSession) WithdrawUnlocked(provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUnlocked(&_Payment.TransactOpts, provider, token)
}

// WithdrawUnlocked is a paid mutator transaction binding the contract method 0x43ea3f36.
//
// Solidity: function withdrawUnlocked(address provider, address token) returns()
func (_Payment *PaymentTransactorSession) WithdrawUnlocked(provider common.Address, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUnlocked(&_Payment.TransactOpts, provider, token)
}

// WithdrawUpTo is a paid mutator transaction binding the contract method 0xe0f33dd4.
//
// Solidity: function withdrawUpTo(address publisher, address token, uint256 totalWithdrawlAmount, address transferAddress) returns()
func (_Payment *PaymentTransactor) WithdrawUpTo(opts *bind.TransactOpts, publisher common.Address, token common.Address, totalWithdrawlAmount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdrawUpTo", publisher, token, totalWithdrawlAmount, transferAddress)
}

// WithdrawUpTo is a paid mutator transaction binding the contract method 0xe0f33dd4.
//
// Solidity: function withdrawUpTo(address publisher, address token, uint256 totalWithdrawlAmount, address transferAddress) returns()
func (_Payment *PaymentSession) WithdrawUpTo(publisher common.Address, token common.Address, totalWithdrawlAmount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUpTo(&_Payment.TransactOpts, publisher, token, totalWithdrawlAmount, transferAddress)
}

// WithdrawUpTo is a paid mutator transaction binding the contract method 0xe0f33dd4.
//
// Solidity: function withdrawUpTo(address publisher, address token, uint256 totalWithdrawlAmount, address transferAddress) returns()
func (_Payment *PaymentTransactorSession) WithdrawUpTo(publisher common.Address, token common.Address, totalWithdrawlAmount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUpTo(&_Payment.TransactOpts, publisher, token, totalWithdrawlAmount, transferAddress)
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
	Token     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelClosed is a free log retrieval operation binding the contract event 0x4509c78f8c652633a480f64f4f56b8187ec70d118e774af3139ec088e6f3ea91.
//
// Solidity: event ChannelClosed(address indexed publisher, address indexed provider, address indexed token)
func (_Payment *PaymentFilterer) FilterChannelClosed(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, token []common.Address) (*PaymentChannelClosedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "ChannelClosed", publisherRule, providerRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PaymentChannelClosedIterator{contract: _Payment.contract, event: "ChannelClosed", logs: logs, sub: sub}, nil
}

// WatchChannelClosed is a free log subscription operation binding the contract event 0x4509c78f8c652633a480f64f4f56b8187ec70d118e774af3139ec088e6f3ea91.
//
// Solidity: event ChannelClosed(address indexed publisher, address indexed provider, address indexed token)
func (_Payment *PaymentFilterer) WatchChannelClosed(opts *bind.WatchOpts, sink chan<- *PaymentChannelClosed, publisher []common.Address, provider []common.Address, token []common.Address) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "ChannelClosed", publisherRule, providerRule, tokenRule)
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

// ParseChannelClosed is a log parse operation binding the contract event 0x4509c78f8c652633a480f64f4f56b8187ec70d118e774af3139ec088e6f3ea91.
//
// Solidity: event ChannelClosed(address indexed publisher, address indexed provider, address indexed token)
func (_Payment *PaymentFilterer) ParseChannelClosed(log types.Log) (*PaymentChannelClosed, error) {
	event := new(PaymentChannelClosed)
	if err := _Payment.contract.UnpackLog(event, "ChannelClosed", log); err != nil {
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
	Token         common.Address
	DepositAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x4174a9435a04d04d274c76779cad136a41fde6937c56241c09ab9d3c7064a1a9.
//
// Solidity: event Deposited(address indexed publisher, address indexed provider, address indexed token, uint256 depositAmount)
func (_Payment *PaymentFilterer) FilterDeposited(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, token []common.Address) (*PaymentDepositedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "Deposited", publisherRule, providerRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PaymentDepositedIterator{contract: _Payment.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x4174a9435a04d04d274c76779cad136a41fde6937c56241c09ab9d3c7064a1a9.
//
// Solidity: event Deposited(address indexed publisher, address indexed provider, address indexed token, uint256 depositAmount)
func (_Payment *PaymentFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *PaymentDeposited, publisher []common.Address, provider []common.Address, token []common.Address) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "Deposited", publisherRule, providerRule, tokenRule)
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

// ParseDeposited is a log parse operation binding the contract event 0x4174a9435a04d04d274c76779cad136a41fde6937c56241c09ab9d3c7064a1a9.
//
// Solidity: event Deposited(address indexed publisher, address indexed provider, address indexed token, uint256 depositAmount)
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
	Token      common.Address
	UnlockedAt *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUnlockTimerStarted is a free log retrieval operation binding the contract event 0x9c21c91a443e1aeab0e24df34b134bf134e2c0c9fecd918faa11306a0adfa62e.
//
// Solidity: event UnlockTimerStarted(address indexed publisher, address indexed provider, address indexed token, uint256 unlockedAt)
func (_Payment *PaymentFilterer) FilterUnlockTimerStarted(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, token []common.Address) (*PaymentUnlockTimerStartedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "UnlockTimerStarted", publisherRule, providerRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PaymentUnlockTimerStartedIterator{contract: _Payment.contract, event: "UnlockTimerStarted", logs: logs, sub: sub}, nil
}

// WatchUnlockTimerStarted is a free log subscription operation binding the contract event 0x9c21c91a443e1aeab0e24df34b134bf134e2c0c9fecd918faa11306a0adfa62e.
//
// Solidity: event UnlockTimerStarted(address indexed publisher, address indexed provider, address indexed token, uint256 unlockedAt)
func (_Payment *PaymentFilterer) WatchUnlockTimerStarted(opts *bind.WatchOpts, sink chan<- *PaymentUnlockTimerStarted, publisher []common.Address, provider []common.Address, token []common.Address) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "UnlockTimerStarted", publisherRule, providerRule, tokenRule)
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

// ParseUnlockTimerStarted is a log parse operation binding the contract event 0x9c21c91a443e1aeab0e24df34b134bf134e2c0c9fecd918faa11306a0adfa62e.
//
// Solidity: event UnlockTimerStarted(address indexed publisher, address indexed provider, address indexed token, uint256 unlockedAt)
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
	Token          common.Address
	UnlockedAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUnlocked is a free log retrieval operation binding the contract event 0x6c5b0a34e5e78b423db3d15d7f4f72f3beb727025c0e535bf2bac21d0227c481.
//
// Solidity: event Unlocked(address indexed publisher, address indexed provider, address indexed token, uint256 unlockedAmount)
func (_Payment *PaymentFilterer) FilterUnlocked(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, token []common.Address) (*PaymentUnlockedIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "Unlocked", publisherRule, providerRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PaymentUnlockedIterator{contract: _Payment.contract, event: "Unlocked", logs: logs, sub: sub}, nil
}

// WatchUnlocked is a free log subscription operation binding the contract event 0x6c5b0a34e5e78b423db3d15d7f4f72f3beb727025c0e535bf2bac21d0227c481.
//
// Solidity: event Unlocked(address indexed publisher, address indexed provider, address indexed token, uint256 unlockedAmount)
func (_Payment *PaymentFilterer) WatchUnlocked(opts *bind.WatchOpts, sink chan<- *PaymentUnlocked, publisher []common.Address, provider []common.Address, token []common.Address) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "Unlocked", publisherRule, providerRule, tokenRule)
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

// ParseUnlocked is a log parse operation binding the contract event 0x6c5b0a34e5e78b423db3d15d7f4f72f3beb727025c0e535bf2bac21d0227c481.
//
// Solidity: event Unlocked(address indexed publisher, address indexed provider, address indexed token, uint256 unlockedAmount)
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
	Token           common.Address
	WithdrawnAmount *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0xa4195c37c2947bbe89165f03e320b6903116f0b10d8cfdb522330f7ce6f9fa24.
//
// Solidity: event Withdrawn(address indexed publisher, address indexed provider, address indexed token, uint256 withdrawnAmount)
func (_Payment *PaymentFilterer) FilterWithdrawn(opts *bind.FilterOpts, publisher []common.Address, provider []common.Address, token []common.Address) (*PaymentWithdrawnIterator, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "Withdrawn", publisherRule, providerRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PaymentWithdrawnIterator{contract: _Payment.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0xa4195c37c2947bbe89165f03e320b6903116f0b10d8cfdb522330f7ce6f9fa24.
//
// Solidity: event Withdrawn(address indexed publisher, address indexed provider, address indexed token, uint256 withdrawnAmount)
func (_Payment *PaymentFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *PaymentWithdrawn, publisher []common.Address, provider []common.Address, token []common.Address) (event.Subscription, error) {

	var publisherRule []interface{}
	for _, publisherItem := range publisher {
		publisherRule = append(publisherRule, publisherItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "Withdrawn", publisherRule, providerRule, tokenRule)
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

// ParseWithdrawn is a log parse operation binding the contract event 0xa4195c37c2947bbe89165f03e320b6903116f0b10d8cfdb522330f7ce6f9fa24.
//
// Solidity: event Withdrawn(address indexed publisher, address indexed provider, address indexed token, uint256 withdrawnAmount)
func (_Payment *PaymentFilterer) ParseWithdrawn(log types.Log) (*PaymentWithdrawn, error) {
	event := new(PaymentWithdrawn)
	if err := _Payment.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
