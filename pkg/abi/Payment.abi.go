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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AmountRequired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChannelLocked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFunds\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"ChannelClosed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"ChannelCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"depositAmount\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockedAt\",\"type\":\"uint256\"}],\"name\":\"UnlockTimerStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockedAmount\",\"type\":\"uint256\"}],\"name\":\"Unlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawnAmount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"available\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"channels\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"investedByPublisher\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawnByProvider\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockedAt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"closeChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialAmount\",\"type\":\"uint256\"}],\"name\":\"createChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transferAddress\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"withdrawUnlocked\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"publisher\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"podId\",\"type\":\"bytes32\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"totalWithdrawlAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"transferAddress\",\"type\":\"address\"}],\"name\":\"withdrawUpTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610d7f806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c80636c24538c116100665780636c24538c146100e6578063708c2cd11461010c578063915d1b511461011f5780639ae9a30614610132578063dbf03dea146101a057600080fd5b8063013c210f14610098578063221e6b5a146100ad5780634b4d7db0146100c05780635f23c0f7146100d3575b600080fd5b6100ab6100a6366004610b18565b6101b3565b005b6100ab6100bb366004610b6a565b6102b0565b6100ab6100ce366004610bc9565b6103d6565b6100ab6100e1366004610b6a565b6104e3565b6100f96100f4366004610c0b565b61053c565b6040519081526020015b60405180910390f35b6100ab61011a366004610c5e565b61058e565b6100ab61012d366004610bc9565b61065c565b610180610140366004610c0b565b6000602081815294815260408082208652938152838120855291825282822090935291825290208054600182015460028301546003909301549192909184565b604080519485526020850193909352918301526060820152608001610103565b6100ab6101ae366004610bc9565b61073e565b806000036101d45760405163d43e607560e01b815260040160405180910390fd5b336000818152602081815260408083206001600160a01b038a811685529083528184208985528352818420908816845290915290208054156102295760405163119b4fd360e11b815260040160405180910390fd5b60018101541561023b5761023b610ca6565b82815560028101849055604080516001600160a01b0387811682526020820186905288928a821692918616917fd5c35171b66326051ea3dc2d87454588db25b3bf09b40a4a98d8186df125fbba910160405180910390a46102a76001600160a01b0386163330866108ba565b50505050505050565b6001600160a01b0381166102c15750335b6001600160a01b038086166000908152602081815260408083203380855290835281842089855283528184209488168452939091529020805484111561031a5760405163356680b760e01b815260040160405180910390fd5b8060010154841161033e5760405163d43e607560e01b815260040160405180910390fd5b60008160010154856103509190610cd2565b60018301869055604080516001600160a01b03898116825260208201849052929350899280871692908c16917f474d5f22a0a93e8d83fd31d4373b4395c9aba71c7917f347f37eea5f66b0adde910160405180910390a46003820154156103b8574260038301555b6103cc6001600160a01b0387168583610927565b5050505050505050565b336000818152602081815260408083206001600160a01b0388811685529083528184208785528352818420908616845290915290206003810154158061041f5750806003015442105b1561043d57604051635652d88760e01b815260040160405180910390fd5b6001810154815460009161045091610cd2565b9050806000036104735760405163d43e607560e01b815260040160405180910390fd5b60018201548255604080516001600160a01b038681168252602082018490528792818a1692918716917f7bee1a6180b28dc049128fe38ca1e5605a4c8db0a26ce3ae2122013579560c09910160405180910390a46104db6001600160a01b0385168483610927565b505050505050565b6001600160a01b0380861660009081526020818152604080832033845282528083208884528252808320938716835292905220600101546105359086908690869061052f908790610ce5565b856102b0565b5050505050565b6001600160a01b038085166000908152602081815260408083208785168452825280832086845282528083209385168352929052908120600181015481546105849190610cd2565b9695505050505050565b806000036105af5760405163d43e607560e01b815260040160405180910390fd5b336000818152602081815260408083206001600160a01b03898116855290835281842088855283528184209087168452909152902080546105f1908490610ce5565b815560006003820155604080516001600160a01b03868116825260208201869052879289821692918616917fd5c35171b66326051ea3dc2d87454588db25b3bf09b40a4a98d8186df125fbba910160405180910390a46104db6001600160a01b0385163330866108ba565b336000818152602081815260408083206001600160a01b03888116855290835281842087855283528184209086168452909152812080549091036106b35760405163b0ce759160e01b815260040160405180910390fd5b60008160020154426106c59190610ce5565b90508160030154600014806106dd5750808260030154105b156104db5760038201819055604080516001600160a01b038681168252602082018490528792818a1692918716917facecda61b3afbd99f7c0b0e5e952c3130ee683634ea648c6863f358fded5f053910160405180910390a4505050505050565b336000818152602081815260408083206001600160a01b038881168552908352818420878552835281842090861684529091529020600381015415806107875750806003015442105b156107a557604051635652d88760e01b815260040160405180910390fd5b600181015481546000916107b891610cd2565b6001600160a01b038085166000908152602081815260408083208b8516845282528083208a8452825280832093891683529290529081208181556001810182905560028101829055600301559050801561085a57604080516001600160a01b038681168252602082018490528792818a1692918716917f7bee1a6180b28dc049128fe38ca1e5605a4c8db0a26ce3ae2122013579560c09910160405180910390a45b6040516001600160a01b0385811682528691818916918616907f5cd428d43ab363aba520572bd20250cdcb20f91ca7c4feab5dc81c8c33ad5ba79060200160405180910390a480156104db576104db6001600160a01b0385168483610927565b6040516001600160a01b0384811660248301528381166044830152606482018390526109219186918216906323b872dd906084015b604051602081830303815290604052915060e01b6020820180516001600160e01b03838183161783525050505061095d565b50505050565b6040516001600160a01b0383811660248301526044820183905261095891859182169063a9059cbb906064016108ef565b505050565b60006109726001600160a01b038416836109c5565b905080516000141580156109975750808060200190518101906109959190610cf8565b155b1561095857604051635274afe760e01b81526001600160a01b03841660048201526024015b60405180910390fd5b60606109d3838360006109dc565b90505b92915050565b606081471015610a015760405163cd78605960e01b81523060048201526024016109bc565b600080856001600160a01b03168486604051610a1d9190610d1a565b60006040518083038185875af1925050503d8060008114610a5a576040519150601f19603f3d011682016040523d82523d6000602084013e610a5f565b606091505b5091509150610a6f868383610a7b565b925050505b9392505050565b606082610a9057610a8b82610ad7565b610a74565b8151158015610aa757506001600160a01b0384163b155b15610ad057604051639996b31560e01b81526001600160a01b03851660048201526024016109bc565b5080610a74565b805115610ae75780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b50565b6001600160a01b0381168114610b0057600080fd5b600080600080600060a08688031215610b3057600080fd5b8535610b3b81610b03565b9450602086013593506040860135610b5281610b03565b94979396509394606081013594506080013592915050565b600080600080600060a08688031215610b8257600080fd5b8535610b8d81610b03565b9450602086013593506040860135610ba481610b03565b9250606086013591506080860135610bbb81610b03565b809150509295509295909350565b600080600060608486031215610bde57600080fd5b8335610be981610b03565b9250602084013591506040840135610c0081610b03565b809150509250925092565b60008060008060808587031215610c2157600080fd5b8435610c2c81610b03565b93506020850135610c3c81610b03565b9250604085013591506060850135610c5381610b03565b939692955090935050565b60008060008060808587031215610c7457600080fd5b8435610c7f81610b03565b9350602085013592506040850135610c9681610b03565b9396929550929360600135925050565b634e487b7160e01b600052600160045260246000fd5b634e487b7160e01b600052601160045260246000fd5b818103818111156109d6576109d6610cbc565b808201808211156109d6576109d6610cbc565b600060208284031215610d0a57600080fd5b81518015158114610a7457600080fd5b6000825160005b81811015610d3b5760208186018101518583015201610d21565b50600092019182525091905056fea264697066735822122032e43ae0e1d5280ccf68a543427134fc5eab507e0327e72abc270915b850f3ea64736f6c63430008150033",
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

// Available is a free data retrieval call binding the contract method 0x6c24538c.
//
// Solidity: function available(address publisher, address provider, bytes32 podId, address token) view returns(uint256)
func (_Payment *PaymentCaller) Available(opts *bind.CallOpts, publisher common.Address, provider common.Address, podId [32]byte, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "available", publisher, provider, podId, token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Available is a free data retrieval call binding the contract method 0x6c24538c.
//
// Solidity: function available(address publisher, address provider, bytes32 podId, address token) view returns(uint256)
func (_Payment *PaymentSession) Available(publisher common.Address, provider common.Address, podId [32]byte, token common.Address) (*big.Int, error) {
	return _Payment.Contract.Available(&_Payment.CallOpts, publisher, provider, podId, token)
}

// Available is a free data retrieval call binding the contract method 0x6c24538c.
//
// Solidity: function available(address publisher, address provider, bytes32 podId, address token) view returns(uint256)
func (_Payment *PaymentCallerSession) Available(publisher common.Address, provider common.Address, podId [32]byte, token common.Address) (*big.Int, error) {
	return _Payment.Contract.Available(&_Payment.CallOpts, publisher, provider, podId, token)
}

// Channels is a free data retrieval call binding the contract method 0x9ae9a306.
//
// Solidity: function channels(address , address , bytes32 , address ) view returns(uint256 investedByPublisher, uint256 withdrawnByProvider, uint256 unlockTime, uint256 unlockedAt)
func (_Payment *PaymentCaller) Channels(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 [32]byte, arg3 common.Address) (struct {
	InvestedByPublisher *big.Int
	WithdrawnByProvider *big.Int
	UnlockTime          *big.Int
	UnlockedAt          *big.Int
}, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "channels", arg0, arg1, arg2, arg3)

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

// Channels is a free data retrieval call binding the contract method 0x9ae9a306.
//
// Solidity: function channels(address , address , bytes32 , address ) view returns(uint256 investedByPublisher, uint256 withdrawnByProvider, uint256 unlockTime, uint256 unlockedAt)
func (_Payment *PaymentSession) Channels(arg0 common.Address, arg1 common.Address, arg2 [32]byte, arg3 common.Address) (struct {
	InvestedByPublisher *big.Int
	WithdrawnByProvider *big.Int
	UnlockTime          *big.Int
	UnlockedAt          *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2, arg3)
}

// Channels is a free data retrieval call binding the contract method 0x9ae9a306.
//
// Solidity: function channels(address , address , bytes32 , address ) view returns(uint256 investedByPublisher, uint256 withdrawnByProvider, uint256 unlockTime, uint256 unlockedAt)
func (_Payment *PaymentCallerSession) Channels(arg0 common.Address, arg1 common.Address, arg2 [32]byte, arg3 common.Address) (struct {
	InvestedByPublisher *big.Int
	WithdrawnByProvider *big.Int
	UnlockTime          *big.Int
	UnlockedAt          *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2, arg3)
}

// CloseChannel is a paid mutator transaction binding the contract method 0xdbf03dea.
//
// Solidity: function closeChannel(address provider, bytes32 podId, address token) returns()
func (_Payment *PaymentTransactor) CloseChannel(opts *bind.TransactOpts, provider common.Address, podId [32]byte, token common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "closeChannel", provider, podId, token)
}

// CloseChannel is a paid mutator transaction binding the contract method 0xdbf03dea.
//
// Solidity: function closeChannel(address provider, bytes32 podId, address token) returns()
func (_Payment *PaymentSession) CloseChannel(provider common.Address, podId [32]byte, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.CloseChannel(&_Payment.TransactOpts, provider, podId, token)
}

// CloseChannel is a paid mutator transaction binding the contract method 0xdbf03dea.
//
// Solidity: function closeChannel(address provider, bytes32 podId, address token) returns()
func (_Payment *PaymentTransactorSession) CloseChannel(provider common.Address, podId [32]byte, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.CloseChannel(&_Payment.TransactOpts, provider, podId, token)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x013c210f.
//
// Solidity: function createChannel(address provider, bytes32 podId, address token, uint256 unlockTime, uint256 initialAmount) returns()
func (_Payment *PaymentTransactor) CreateChannel(opts *bind.TransactOpts, provider common.Address, podId [32]byte, token common.Address, unlockTime *big.Int, initialAmount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "createChannel", provider, podId, token, unlockTime, initialAmount)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x013c210f.
//
// Solidity: function createChannel(address provider, bytes32 podId, address token, uint256 unlockTime, uint256 initialAmount) returns()
func (_Payment *PaymentSession) CreateChannel(provider common.Address, podId [32]byte, token common.Address, unlockTime *big.Int, initialAmount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateChannel(&_Payment.TransactOpts, provider, podId, token, unlockTime, initialAmount)
}

// CreateChannel is a paid mutator transaction binding the contract method 0x013c210f.
//
// Solidity: function createChannel(address provider, bytes32 podId, address token, uint256 unlockTime, uint256 initialAmount) returns()
func (_Payment *PaymentTransactorSession) CreateChannel(provider common.Address, podId [32]byte, token common.Address, unlockTime *big.Int, initialAmount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateChannel(&_Payment.TransactOpts, provider, podId, token, unlockTime, initialAmount)
}

// Deposit is a paid mutator transaction binding the contract method 0x708c2cd1.
//
// Solidity: function deposit(address provider, bytes32 podId, address token, uint256 amount) returns()
func (_Payment *PaymentTransactor) Deposit(opts *bind.TransactOpts, provider common.Address, podId [32]byte, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "deposit", provider, podId, token, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x708c2cd1.
//
// Solidity: function deposit(address provider, bytes32 podId, address token, uint256 amount) returns()
func (_Payment *PaymentSession) Deposit(provider common.Address, podId [32]byte, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.Deposit(&_Payment.TransactOpts, provider, podId, token, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x708c2cd1.
//
// Solidity: function deposit(address provider, bytes32 podId, address token, uint256 amount) returns()
func (_Payment *PaymentTransactorSession) Deposit(provider common.Address, podId [32]byte, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.Deposit(&_Payment.TransactOpts, provider, podId, token, amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x915d1b51.
//
// Solidity: function unlock(address provider, bytes32 podId, address token) returns()
func (_Payment *PaymentTransactor) Unlock(opts *bind.TransactOpts, provider common.Address, podId [32]byte, token common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "unlock", provider, podId, token)
}

// Unlock is a paid mutator transaction binding the contract method 0x915d1b51.
//
// Solidity: function unlock(address provider, bytes32 podId, address token) returns()
func (_Payment *PaymentSession) Unlock(provider common.Address, podId [32]byte, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Unlock(&_Payment.TransactOpts, provider, podId, token)
}

// Unlock is a paid mutator transaction binding the contract method 0x915d1b51.
//
// Solidity: function unlock(address provider, bytes32 podId, address token) returns()
func (_Payment *PaymentTransactorSession) Unlock(provider common.Address, podId [32]byte, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Unlock(&_Payment.TransactOpts, provider, podId, token)
}

// Withdraw is a paid mutator transaction binding the contract method 0x5f23c0f7.
//
// Solidity: function withdraw(address publisher, bytes32 podId, address token, uint256 amount, address transferAddress) returns()
func (_Payment *PaymentTransactor) Withdraw(opts *bind.TransactOpts, publisher common.Address, podId [32]byte, token common.Address, amount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdraw", publisher, podId, token, amount, transferAddress)
}

// Withdraw is a paid mutator transaction binding the contract method 0x5f23c0f7.
//
// Solidity: function withdraw(address publisher, bytes32 podId, address token, uint256 amount, address transferAddress) returns()
func (_Payment *PaymentSession) Withdraw(publisher common.Address, podId [32]byte, token common.Address, amount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, publisher, podId, token, amount, transferAddress)
}

// Withdraw is a paid mutator transaction binding the contract method 0x5f23c0f7.
//
// Solidity: function withdraw(address publisher, bytes32 podId, address token, uint256 amount, address transferAddress) returns()
func (_Payment *PaymentTransactorSession) Withdraw(publisher common.Address, podId [32]byte, token common.Address, amount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, publisher, podId, token, amount, transferAddress)
}

// WithdrawUnlocked is a paid mutator transaction binding the contract method 0x4b4d7db0.
//
// Solidity: function withdrawUnlocked(address provider, bytes32 podId, address token) returns()
func (_Payment *PaymentTransactor) WithdrawUnlocked(opts *bind.TransactOpts, provider common.Address, podId [32]byte, token common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdrawUnlocked", provider, podId, token)
}

// WithdrawUnlocked is a paid mutator transaction binding the contract method 0x4b4d7db0.
//
// Solidity: function withdrawUnlocked(address provider, bytes32 podId, address token) returns()
func (_Payment *PaymentSession) WithdrawUnlocked(provider common.Address, podId [32]byte, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUnlocked(&_Payment.TransactOpts, provider, podId, token)
}

// WithdrawUnlocked is a paid mutator transaction binding the contract method 0x4b4d7db0.
//
// Solidity: function withdrawUnlocked(address provider, bytes32 podId, address token) returns()
func (_Payment *PaymentTransactorSession) WithdrawUnlocked(provider common.Address, podId [32]byte, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUnlocked(&_Payment.TransactOpts, provider, podId, token)
}

// WithdrawUpTo is a paid mutator transaction binding the contract method 0x221e6b5a.
//
// Solidity: function withdrawUpTo(address publisher, bytes32 podId, address token, uint256 totalWithdrawlAmount, address transferAddress) returns()
func (_Payment *PaymentTransactor) WithdrawUpTo(opts *bind.TransactOpts, publisher common.Address, podId [32]byte, token common.Address, totalWithdrawlAmount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdrawUpTo", publisher, podId, token, totalWithdrawlAmount, transferAddress)
}

// WithdrawUpTo is a paid mutator transaction binding the contract method 0x221e6b5a.
//
// Solidity: function withdrawUpTo(address publisher, bytes32 podId, address token, uint256 totalWithdrawlAmount, address transferAddress) returns()
func (_Payment *PaymentSession) WithdrawUpTo(publisher common.Address, podId [32]byte, token common.Address, totalWithdrawlAmount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUpTo(&_Payment.TransactOpts, publisher, podId, token, totalWithdrawlAmount, transferAddress)
}

// WithdrawUpTo is a paid mutator transaction binding the contract method 0x221e6b5a.
//
// Solidity: function withdrawUpTo(address publisher, bytes32 podId, address token, uint256 totalWithdrawlAmount, address transferAddress) returns()
func (_Payment *PaymentTransactorSession) WithdrawUpTo(publisher common.Address, podId [32]byte, token common.Address, totalWithdrawlAmount *big.Int, transferAddress common.Address) (*types.Transaction, error) {
	return _Payment.Contract.WithdrawUpTo(&_Payment.TransactOpts, publisher, podId, token, totalWithdrawlAmount, transferAddress)
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
	Token     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelClosed is a free log retrieval operation binding the contract event 0x5cd428d43ab363aba520572bd20250cdcb20f91ca7c4feab5dc81c8c33ad5ba7.
//
// Solidity: event ChannelClosed(address indexed publisher, address indexed provider, bytes32 indexed podId, address token)
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

// WatchChannelClosed is a free log subscription operation binding the contract event 0x5cd428d43ab363aba520572bd20250cdcb20f91ca7c4feab5dc81c8c33ad5ba7.
//
// Solidity: event ChannelClosed(address indexed publisher, address indexed provider, bytes32 indexed podId, address token)
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

// ParseChannelClosed is a log parse operation binding the contract event 0x5cd428d43ab363aba520572bd20250cdcb20f91ca7c4feab5dc81c8c33ad5ba7.
//
// Solidity: event ChannelClosed(address indexed publisher, address indexed provider, bytes32 indexed podId, address token)
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
	Token     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChannelCreated is a free log retrieval operation binding the contract event 0xc62979fe49c7a61465cd3d1cb0627a16036ac39287d686ae063b483aa1f68df7.
//
// Solidity: event ChannelCreated(address indexed publisher, address indexed provider, bytes32 indexed podId, address token)
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

// WatchChannelCreated is a free log subscription operation binding the contract event 0xc62979fe49c7a61465cd3d1cb0627a16036ac39287d686ae063b483aa1f68df7.
//
// Solidity: event ChannelCreated(address indexed publisher, address indexed provider, bytes32 indexed podId, address token)
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

// ParseChannelCreated is a log parse operation binding the contract event 0xc62979fe49c7a61465cd3d1cb0627a16036ac39287d686ae063b483aa1f68df7.
//
// Solidity: event ChannelCreated(address indexed publisher, address indexed provider, bytes32 indexed podId, address token)
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
	Token         common.Address
	DepositAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0xd5c35171b66326051ea3dc2d87454588db25b3bf09b40a4a98d8186df125fbba.
//
// Solidity: event Deposited(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 depositAmount)
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

// WatchDeposited is a free log subscription operation binding the contract event 0xd5c35171b66326051ea3dc2d87454588db25b3bf09b40a4a98d8186df125fbba.
//
// Solidity: event Deposited(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 depositAmount)
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

// ParseDeposited is a log parse operation binding the contract event 0xd5c35171b66326051ea3dc2d87454588db25b3bf09b40a4a98d8186df125fbba.
//
// Solidity: event Deposited(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 depositAmount)
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
	Token      common.Address
	UnlockedAt *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUnlockTimerStarted is a free log retrieval operation binding the contract event 0xacecda61b3afbd99f7c0b0e5e952c3130ee683634ea648c6863f358fded5f053.
//
// Solidity: event UnlockTimerStarted(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 unlockedAt)
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

// WatchUnlockTimerStarted is a free log subscription operation binding the contract event 0xacecda61b3afbd99f7c0b0e5e952c3130ee683634ea648c6863f358fded5f053.
//
// Solidity: event UnlockTimerStarted(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 unlockedAt)
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

// ParseUnlockTimerStarted is a log parse operation binding the contract event 0xacecda61b3afbd99f7c0b0e5e952c3130ee683634ea648c6863f358fded5f053.
//
// Solidity: event UnlockTimerStarted(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 unlockedAt)
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
	Token          common.Address
	UnlockedAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUnlocked is a free log retrieval operation binding the contract event 0x7bee1a6180b28dc049128fe38ca1e5605a4c8db0a26ce3ae2122013579560c09.
//
// Solidity: event Unlocked(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 unlockedAmount)
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

// WatchUnlocked is a free log subscription operation binding the contract event 0x7bee1a6180b28dc049128fe38ca1e5605a4c8db0a26ce3ae2122013579560c09.
//
// Solidity: event Unlocked(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 unlockedAmount)
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

// ParseUnlocked is a log parse operation binding the contract event 0x7bee1a6180b28dc049128fe38ca1e5605a4c8db0a26ce3ae2122013579560c09.
//
// Solidity: event Unlocked(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 unlockedAmount)
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
	Token           common.Address
	WithdrawnAmount *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x474d5f22a0a93e8d83fd31d4373b4395c9aba71c7917f347f37eea5f66b0adde.
//
// Solidity: event Withdrawn(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 withdrawnAmount)
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

// WatchWithdrawn is a free log subscription operation binding the contract event 0x474d5f22a0a93e8d83fd31d4373b4395c9aba71c7917f347f37eea5f66b0adde.
//
// Solidity: event Withdrawn(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 withdrawnAmount)
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

// ParseWithdrawn is a log parse operation binding the contract event 0x474d5f22a0a93e8d83fd31d4373b4395c9aba71c7917f347f37eea5f66b0adde.
//
// Solidity: event Withdrawn(address indexed publisher, address indexed provider, bytes32 indexed podId, address token, uint256 withdrawnAmount)
func (_Payment *PaymentFilterer) ParseWithdrawn(log types.Log) (*PaymentWithdrawn, error) {
	event := new(PaymentWithdrawn)
	if err := _Payment.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
