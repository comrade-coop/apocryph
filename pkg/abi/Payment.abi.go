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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"NewPrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"NewPriceUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"podID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"PaymentChannel\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"podID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"acceptNewPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"channels\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"owedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"suggestedPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minAdvanceDuration\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minAdvanceDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"createChannel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"incrementId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"podID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"lockFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"podID\",\"type\":\"uint256\"}],\"name\":\"unclockFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"podID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newDeadline\",\"type\":\"uint256\"}],\"name\":\"updateDeadline\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"podID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"updatePrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"podID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"units\",\"type\":\"uint256\"}],\"name\":\"uploadMetrics\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"client\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"podID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000805534801561001457600080fd5b506113f8806100246000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c806369328dec1161006657806369328dec1461017d57806370e6c99314610190578063bda208d1146101a3578063d871122e146101ab578063d90d8bc8146101be57600080fd5b8063159d80e6146100a35780632805a455146100b857806348190c5e146100cb5780634dc63aa8146101575780635e012db01461016a575b600080fd5b6100b66100b136600461117c565b6101d1565b005b6100b66100c636600461117c565b61043f565b6101266100d93660046111c0565b600160208181526000958652604080872082529486528486208152928552838520909252835291208054918101546002820154600383015460048401546005909401549293919290919086565b604080519687526020870195909552938501929092526060840152608083015260a082015260c00160405180910390f35b6100b661016536600461117c565b610591565b6100b661017836600461117c565b6106ec565b6100b661018b36600461120d565b610883565b6100b661019e366004611249565b6109f0565b6100b6610d96565b6100b66101b936600461120d565b610da9565b6100b66101cc36600461129e565b610f1d565b80600081116101fb5760405162461bcd60e51b81526004016101f2906112da565b60405180910390fd5b3360008181526001602081815260408084206001600160a01b03808c1686529083528185208a865283528185209089168552825292839020835160c081018552815481529281015491830191909152600281015492820183905260038101546060830152600481015460808301526005015460a082015287918791879142106102965760405162461bcd60e51b81526004016101f290611311565b3360009081526001602081815260408084206001600160a01b038f811686529083528185208e86528352818520908d168552825292839020835160c08101855281548082529382015492810192909252600281015493820193909352600383015460608201526004830154608082015260059092015460a083015261031c908990611350565b60016000336001600160a01b03166001600160a01b0316815260200190815260200160002060008d6001600160a01b03166001600160a01b0316815260200190815260200160002060008c815260200190815260200160002060008b6001600160a01b03166001600160a01b0316815260200190815260200160002060000181905550886001600160a01b03166323b872dd33308b6040518463ffffffff1660e01b81526004016103ee939291906001600160a01b039384168152919092166020820152604081019190915260600190565b6020604051808303816000875af115801561040d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104319190611369565b505050505050505050505050565b80600081116104605760405162461bcd60e51b81526004016101f2906112da565b6001600160a01b038086166000908152600160208181526040808420338086529083528185208a86528352818520958916855294825292839020835160c081018552815481529281015491830191909152600281015492820183905260038101546060830152600481015460808301526005015460a08201528792918791879142106104fe5760405162461bcd60e51b81526004016101f290611311565b6001600160a01b038a81166000818152600160209081526040808320338085529083528184208f85528352818420958e16808552958352928190206004018c9055805193845290830191909152810191909152606081018890527f321c73989a5585f52e30bb8aabd608154125bf4413611a3b4d8b97105d9c6dc09060800160405180910390a150505050505050505050565b6001600160a01b038085166000908152600160208181526040808420338552825280842088855282528084209487168452938152838320845160c0810186528154815292810154918301919091526002810154938201939093526003830154606082018190526004840154608083015260059093015460a0820152916106179084611392565b825160208401519192509061062c9083611350565b11156106965760405162461bcd60e51b815260206004820152603360248201527f4f776564416d6f756e7420697320626967676572207468616e206368616e6e656044820152726c277320617661696c61626c652066756e647360681b60648201526084016101f2565b60208201516106a59082611350565b6001600160a01b0396871660009081526001602081815260408084203385528252808420998452988152888320979099168252959097529490952090920192909255505050565b3360009081526001602081815260408084206001600160a01b03898116865290835281852088865283528185209087168552825292839020835160c081018552815481529281015491830191909152600281015492820192909252600382015460608201526004820154608082015260059091015460a0820181905282116107af5760405162461bcd60e51b815260206004820152601660248201527513995dc8111958591b1a5b99481d1bdbc81cda1bdc9d60521b60448201526064016101f2565b806040015182116108155760405162461bcd60e51b815260206004820152602a60248201527f4e657720446561646c696e65206973206c657373207468616e2063757272656e6044820152697420646561646c696e6560b01b60648201526084016101f2565b60408082019283523360009081526001602081815283832081528383209783529687528282206001600160a01b0396909616825294865220815181559381015192840192909255516002830155606081015160038301556080810151600483015560a0015160059091015550565b6001600160a01b03808416600090815260016020818152604080842033855282528084208785528252808420948616845293815291839020835160c0810185528154815291810154928201839052600281015493820193909352600383015460608201526004830154608082015260059092015460a08301526109395760405162461bcd60e51b815260206004820152600e60248201526d05a65726f204f776e6572736869760941b60448201526064016101f2565b602081015160405163a9059cbb60e01b815233600482015260248101919091526001600160a01b0383169063a9059cbb906044016020604051808303816000875af115801561098c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109b09190611369565b50506001600160a01b0392831660009081526001602081815260408084203385528252808420958452948152848320939095168252919093529082200155565b8360008111610a115760405162461bcd60e51b81526004016101f2906112da565b604051636eb1769f60e11b81523360048201819052306024830152879187906000906001600160a01b0385169063dd62ed3e90604401602060405180830381865afa158015610a64573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a8891906113a9565b9050818114610aeb5760405162461bcd60e51b815260206004820152602960248201527f616c6c6f77616e636520646f6573206e6f74206d617463682073706563696669604482015268195908185b5bdd5b9d60ba1b60648201526084016101f2565b428811610b2d5760405162461bcd60e51b815260206004820152601060248201526f111958591b1a5b9948115e1c1a5c995960821b60448201526064016101f2565b600060016000336001600160a01b03166001600160a01b0316815260200190815260200160002060008d6001600160a01b03166001600160a01b03168152602001908152602001600020600080546001610b879190611350565b815260200190815260200160002060008c6001600160a01b03166001600160a01b031681526020019081526020016000206040518060c0016040529081600082015481526020016001820154815260200160028201548152602001600382015481526020016004820154815260200160058201548152505090506040518060c001604052808b8152602001600081526020018a8152602001888152602001600081526020018981525090508060016000336001600160a01b03166001600160a01b0316815260200190815260200160002060008e6001600160a01b03166001600160a01b03168152602001908152602001600020600080546001610c8b9190611350565b815260200190815260200160002060008d6001600160a01b03166001600160a01b03168152602001908152602001600020600082015181600001556020820151816001015560408201518160020155606082015181600301556080820151816004015560a082015181600501559050508a6001600160a01b03166323b872dd333084600001516040518463ffffffff1660e01b8152600401610d4e939291906001600160a01b039384168152919092166020820152604081019190915260600190565b6020604051808303816000875af1158015610d6d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d919190611369565b506104315b600054610da4906001611350565b600055565b3360008181526001602081815260408084206001600160a01b03808a16865290835281852088865283528185209087168552825292839020835160c081018552815481529281015491830191909152600281015492820183905260038101546060830152600481015460808301526005015460a08201528591859185914210610e445760405162461bcd60e51b81526004016101f290611311565b3360008181526001602081815260408084206001600160a01b038e811686529083528185208d86528352818520908c16808652818452828620835160c08101855281548152958101548686015260028101548685015260038101805460608089019182526004840180546080808c0182905260059096015460a08c0152909355848a52948752979055955183518881529485019790975283830152820194909452925190927f8c86180e4992276e0f056bb49eed13ea7b192a0d9e795a7f219b5b281c22d78e92908290030190a1505050505050505050565b3360009081526001602081815260408084206001600160a01b03878116865290835281852086865283528185209088168552825292839020835160c081018552815481529281015491830191909152600281015492820183905260038101546060830152600481015460808301526005015460a082015290421015610fe45760405162461bcd60e51b815260206004820152601860248201527f446561646c696e65206e6f74207265616368656420796574000000000000000060448201526064016101f2565b60208101511561106a57602081015160405163a9059cbb60e01b81526001600160a01b03858116600483015260248201929092529085169063a9059cbb906044016020604051808303816000875af1158015611044573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110689190611369565b505b80516110a85760405162461bcd60e51b815260206004820152600d60248201526c115b5c1d1e4810da185b9b995b609a1b60448201526064016101f2565b805160405163a9059cbb60e01b815233600482015260248101919091526001600160a01b0385169063a9059cbb906044016020604051808303816000875af11580156110f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061111c9190611369565b50503360009081526001602081815260408084206001600160a01b039687168552825280842094845293815283832095909416825293909252812081815590910155565b80356001600160a01b038116811461117757600080fd5b919050565b6000806000806080858703121561119257600080fd5b61119b85611160565b9350602085013592506111b060408601611160565b9396929550929360600135925050565b600080600080608085870312156111d657600080fd5b6111df85611160565b93506111ed60208601611160565b92506040850135915061120260608601611160565b905092959194509250565b60008060006060848603121561122257600080fd5b61122b84611160565b92506020840135915061124060408501611160565b90509250925092565b60008060008060008060c0878903121561126257600080fd5b61126b87611160565b955061127960208801611160565b95989597505050506040840135936060810135936080820135935060a0909101359150565b6000806000606084860312156112b357600080fd5b6112bc84611160565b92506112ca60208501611160565b9150604084013590509250925092565b6020808252601c908201527f63616e277420616363706574207a65726f20617320612076616c756500000000604082015260600190565b6020808252600f908201526e10da185b9b995b08115e1c1a5c9959608a1b604082015260600190565b634e487b7160e01b600052601160045260246000fd5b808201808211156113635761136361133a565b92915050565b60006020828403121561137b57600080fd5b8151801515811461138b57600080fd5b9392505050565b80820281158282048414176113635761136361133a565b6000602082840312156113bb57600080fd5b505191905056fea26469706673582212208e774ad99cae02b06cd479eea330633baa27604c792adce3c38d53c94380f16964736f6c63430008150033",
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

// Channels is a free data retrieval call binding the contract method 0x48190c5e.
//
// Solidity: function channels(address , address , uint256 , address ) view returns(uint256 total, uint256 owedAmount, uint256 deadline, uint256 price, uint256 suggestedPrice, uint256 minAdvanceDuration)
func (_Payment *PaymentCaller) Channels(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 common.Address) (struct {
	Total              *big.Int
	OwedAmount         *big.Int
	Deadline           *big.Int
	Price              *big.Int
	SuggestedPrice     *big.Int
	MinAdvanceDuration *big.Int
}, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "channels", arg0, arg1, arg2, arg3)

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

// Channels is a free data retrieval call binding the contract method 0x48190c5e.
//
// Solidity: function channels(address , address , uint256 , address ) view returns(uint256 total, uint256 owedAmount, uint256 deadline, uint256 price, uint256 suggestedPrice, uint256 minAdvanceDuration)
func (_Payment *PaymentSession) Channels(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 common.Address) (struct {
	Total              *big.Int
	OwedAmount         *big.Int
	Deadline           *big.Int
	Price              *big.Int
	SuggestedPrice     *big.Int
	MinAdvanceDuration *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2, arg3)
}

// Channels is a free data retrieval call binding the contract method 0x48190c5e.
//
// Solidity: function channels(address , address , uint256 , address ) view returns(uint256 total, uint256 owedAmount, uint256 deadline, uint256 price, uint256 suggestedPrice, uint256 minAdvanceDuration)
func (_Payment *PaymentCallerSession) Channels(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 common.Address) (struct {
	Total              *big.Int
	OwedAmount         *big.Int
	Deadline           *big.Int
	Price              *big.Int
	SuggestedPrice     *big.Int
	MinAdvanceDuration *big.Int
}, error) {
	return _Payment.Contract.Channels(&_Payment.CallOpts, arg0, arg1, arg2, arg3)
}

// AcceptNewPrice is a paid mutator transaction binding the contract method 0xd871122e.
//
// Solidity: function acceptNewPrice(address provider, uint256 podID, address token) returns()
func (_Payment *PaymentTransactor) AcceptNewPrice(opts *bind.TransactOpts, provider common.Address, podID *big.Int, token common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "acceptNewPrice", provider, podID, token)
}

// AcceptNewPrice is a paid mutator transaction binding the contract method 0xd871122e.
//
// Solidity: function acceptNewPrice(address provider, uint256 podID, address token) returns()
func (_Payment *PaymentSession) AcceptNewPrice(provider common.Address, podID *big.Int, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.AcceptNewPrice(&_Payment.TransactOpts, provider, podID, token)
}

// AcceptNewPrice is a paid mutator transaction binding the contract method 0xd871122e.
//
// Solidity: function acceptNewPrice(address provider, uint256 podID, address token) returns()
func (_Payment *PaymentTransactorSession) AcceptNewPrice(provider common.Address, podID *big.Int, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.AcceptNewPrice(&_Payment.TransactOpts, provider, podID, token)
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

// IncrementId is a paid mutator transaction binding the contract method 0xbda208d1.
//
// Solidity: function incrementId() returns()
func (_Payment *PaymentTransactor) IncrementId(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "incrementId")
}

// IncrementId is a paid mutator transaction binding the contract method 0xbda208d1.
//
// Solidity: function incrementId() returns()
func (_Payment *PaymentSession) IncrementId() (*types.Transaction, error) {
	return _Payment.Contract.IncrementId(&_Payment.TransactOpts)
}

// IncrementId is a paid mutator transaction binding the contract method 0xbda208d1.
//
// Solidity: function incrementId() returns()
func (_Payment *PaymentTransactorSession) IncrementId() (*types.Transaction, error) {
	return _Payment.Contract.IncrementId(&_Payment.TransactOpts)
}

// LockFunds is a paid mutator transaction binding the contract method 0x159d80e6.
//
// Solidity: function lockFunds(address provider, uint256 podID, address token, uint256 amount) returns()
func (_Payment *PaymentTransactor) LockFunds(opts *bind.TransactOpts, provider common.Address, podID *big.Int, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "lockFunds", provider, podID, token, amount)
}

// LockFunds is a paid mutator transaction binding the contract method 0x159d80e6.
//
// Solidity: function lockFunds(address provider, uint256 podID, address token, uint256 amount) returns()
func (_Payment *PaymentSession) LockFunds(provider common.Address, podID *big.Int, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.LockFunds(&_Payment.TransactOpts, provider, podID, token, amount)
}

// LockFunds is a paid mutator transaction binding the contract method 0x159d80e6.
//
// Solidity: function lockFunds(address provider, uint256 podID, address token, uint256 amount) returns()
func (_Payment *PaymentTransactorSession) LockFunds(provider common.Address, podID *big.Int, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.LockFunds(&_Payment.TransactOpts, provider, podID, token, amount)
}

// UnclockFunds is a paid mutator transaction binding the contract method 0xd90d8bc8.
//
// Solidity: function unclockFunds(address token, address provider, uint256 podID) returns()
func (_Payment *PaymentTransactor) UnclockFunds(opts *bind.TransactOpts, token common.Address, provider common.Address, podID *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "unclockFunds", token, provider, podID)
}

// UnclockFunds is a paid mutator transaction binding the contract method 0xd90d8bc8.
//
// Solidity: function unclockFunds(address token, address provider, uint256 podID) returns()
func (_Payment *PaymentSession) UnclockFunds(token common.Address, provider common.Address, podID *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UnclockFunds(&_Payment.TransactOpts, token, provider, podID)
}

// UnclockFunds is a paid mutator transaction binding the contract method 0xd90d8bc8.
//
// Solidity: function unclockFunds(address token, address provider, uint256 podID) returns()
func (_Payment *PaymentTransactorSession) UnclockFunds(token common.Address, provider common.Address, podID *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UnclockFunds(&_Payment.TransactOpts, token, provider, podID)
}

// UpdateDeadline is a paid mutator transaction binding the contract method 0x5e012db0.
//
// Solidity: function updateDeadline(address provider, uint256 podID, address token, uint256 newDeadline) returns()
func (_Payment *PaymentTransactor) UpdateDeadline(opts *bind.TransactOpts, provider common.Address, podID *big.Int, token common.Address, newDeadline *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "updateDeadline", provider, podID, token, newDeadline)
}

// UpdateDeadline is a paid mutator transaction binding the contract method 0x5e012db0.
//
// Solidity: function updateDeadline(address provider, uint256 podID, address token, uint256 newDeadline) returns()
func (_Payment *PaymentSession) UpdateDeadline(provider common.Address, podID *big.Int, token common.Address, newDeadline *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UpdateDeadline(&_Payment.TransactOpts, provider, podID, token, newDeadline)
}

// UpdateDeadline is a paid mutator transaction binding the contract method 0x5e012db0.
//
// Solidity: function updateDeadline(address provider, uint256 podID, address token, uint256 newDeadline) returns()
func (_Payment *PaymentTransactorSession) UpdateDeadline(provider common.Address, podID *big.Int, token common.Address, newDeadline *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UpdateDeadline(&_Payment.TransactOpts, provider, podID, token, newDeadline)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0x2805a455.
//
// Solidity: function updatePrice(address client, uint256 podID, address token, uint256 price) returns()
func (_Payment *PaymentTransactor) UpdatePrice(opts *bind.TransactOpts, client common.Address, podID *big.Int, token common.Address, price *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "updatePrice", client, podID, token, price)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0x2805a455.
//
// Solidity: function updatePrice(address client, uint256 podID, address token, uint256 price) returns()
func (_Payment *PaymentSession) UpdatePrice(client common.Address, podID *big.Int, token common.Address, price *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UpdatePrice(&_Payment.TransactOpts, client, podID, token, price)
}

// UpdatePrice is a paid mutator transaction binding the contract method 0x2805a455.
//
// Solidity: function updatePrice(address client, uint256 podID, address token, uint256 price) returns()
func (_Payment *PaymentTransactorSession) UpdatePrice(client common.Address, podID *big.Int, token common.Address, price *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UpdatePrice(&_Payment.TransactOpts, client, podID, token, price)
}

// UploadMetrics is a paid mutator transaction binding the contract method 0x4dc63aa8.
//
// Solidity: function uploadMetrics(address client, uint256 podID, address token, uint256 units) returns()
func (_Payment *PaymentTransactor) UploadMetrics(opts *bind.TransactOpts, client common.Address, podID *big.Int, token common.Address, units *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "uploadMetrics", client, podID, token, units)
}

// UploadMetrics is a paid mutator transaction binding the contract method 0x4dc63aa8.
//
// Solidity: function uploadMetrics(address client, uint256 podID, address token, uint256 units) returns()
func (_Payment *PaymentSession) UploadMetrics(client common.Address, podID *big.Int, token common.Address, units *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UploadMetrics(&_Payment.TransactOpts, client, podID, token, units)
}

// UploadMetrics is a paid mutator transaction binding the contract method 0x4dc63aa8.
//
// Solidity: function uploadMetrics(address client, uint256 podID, address token, uint256 units) returns()
func (_Payment *PaymentTransactorSession) UploadMetrics(client common.Address, podID *big.Int, token common.Address, units *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.UploadMetrics(&_Payment.TransactOpts, client, podID, token, units)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address client, uint256 podID, address token) returns()
func (_Payment *PaymentTransactor) Withdraw(opts *bind.TransactOpts, client common.Address, podID *big.Int, token common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "withdraw", client, podID, token)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address client, uint256 podID, address token) returns()
func (_Payment *PaymentSession) Withdraw(client common.Address, podID *big.Int, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, client, podID, token)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address client, uint256 podID, address token) returns()
func (_Payment *PaymentTransactorSession) Withdraw(client common.Address, podID *big.Int, token common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Withdraw(&_Payment.TransactOpts, client, podID, token)
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

// PaymentPaymentChannelIterator is returned from FilterPaymentChannel and is used to iterate over the raw logs and unpacked data for PaymentChannel events raised by the Payment contract.
type PaymentPaymentChannelIterator struct {
	Event *PaymentPaymentChannel // Event containing the contract specifics and raw log

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
func (it *PaymentPaymentChannelIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentPaymentChannel)
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
		it.Event = new(PaymentPaymentChannel)
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
func (it *PaymentPaymentChannelIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentPaymentChannelIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentPaymentChannel represents a PaymentChannel event raised by the Payment contract.
type PaymentPaymentChannel struct {
	Client   common.Address
	Provider common.Address
	PodID    *big.Int
	Token    common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPaymentChannel is a free log retrieval operation binding the contract event 0x38c88eaeba9f4f82e722c3d38b2caa14fbd5155db0e5f80b3ffb64c9678226f0.
//
// Solidity: event PaymentChannel(address client, address provider, uint256 podID, address token)
func (_Payment *PaymentFilterer) FilterPaymentChannel(opts *bind.FilterOpts) (*PaymentPaymentChannelIterator, error) {

	logs, sub, err := _Payment.contract.FilterLogs(opts, "PaymentChannel")
	if err != nil {
		return nil, err
	}
	return &PaymentPaymentChannelIterator{contract: _Payment.contract, event: "PaymentChannel", logs: logs, sub: sub}, nil
}

// WatchPaymentChannel is a free log subscription operation binding the contract event 0x38c88eaeba9f4f82e722c3d38b2caa14fbd5155db0e5f80b3ffb64c9678226f0.
//
// Solidity: event PaymentChannel(address client, address provider, uint256 podID, address token)
func (_Payment *PaymentFilterer) WatchPaymentChannel(opts *bind.WatchOpts, sink chan<- *PaymentPaymentChannel) (event.Subscription, error) {

	logs, sub, err := _Payment.contract.WatchLogs(opts, "PaymentChannel")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentPaymentChannel)
				if err := _Payment.contract.UnpackLog(event, "PaymentChannel", log); err != nil {
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

// ParsePaymentChannel is a log parse operation binding the contract event 0x38c88eaeba9f4f82e722c3d38b2caa14fbd5155db0e5f80b3ffb64c9678226f0.
//
// Solidity: event PaymentChannel(address client, address provider, uint256 podID, address token)
func (_Payment *PaymentFilterer) ParsePaymentChannel(log types.Log) (*PaymentPaymentChannel, error) {
	event := new(PaymentPaymentChannel)
	if err := _Payment.contract.UnpackLog(event, "PaymentChannel", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
