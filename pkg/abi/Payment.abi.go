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
	Bin: "0x608060405234801561000f575f80fd5b5061222d8061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610091575f3560e01c8063995cee2611610064578063995cee2614610105578063a94d241514610121578063c7f149261461013d578063f428e28a14610159578063f940e3851461018e57610091565b8063165c326f1461009557806331f60646146100b157806370e6c993146100cd5780638fbb2a90146100e9575b5f80fd5b6100af60048036038101906100aa9190611913565b6101aa565b005b6100cb60048036038101906100c69190611913565b6103dd565b005b6100e760048036038101906100e29190611963565b610664565b005b61010360048036038101906100fe9190611913565b610a99565b005b61011f600480360381019061011a91906119ec565b610bd9565b005b61013b60048036038101906101369190611913565b610fe4565b005b610157600480360381019061015291906119ec565b6113b6565b005b610173600480360381019061016e9190611a2a565b61166f565b60405161018596959493929190611a89565b60405180910390f35b6101a860048036038101906101a391906119ec565b6116bc565b005b5f805f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015481526020016001820154815260200160028201548152602001600382015481526020016004820154815260200160058201548152505090505f8160600151836102b69190611b15565b9050815f01518260200151826102cc9190611b56565b111561030d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161030490611c09565b60405180910390fd5b81602001518161031d9190611b56565b5f808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20600101819055505050505050565b5f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015481526020016001820154815260200160028201548152602001600382015481526020016004820154815260200160058201548152505090508060a00151821161051e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161051590611c71565b60405180910390fd5b80604001518211610564576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161055b90611cff565b60405180910390fd5b81816040018181525050805f803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f820151815f01556020820151816001015560408201518160020155606082015181600301556080820151816004015560a0820151816005015590505050505050565b835f81116106a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161069e90611d67565b60405180910390fd5b8533865f8373ffffffffffffffffffffffffffffffffffffffff1663dd62ed3e84306040518363ffffffff1660e01b81526004016106e6929190611d94565b602060405180830381865afa158015610701573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107259190611dcf565b9050818114610769576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161076090611e6a565b60405180910390fd5b4288116107ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107a290611ed2565b60405180910390fd5b5f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8d73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8c73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015481526020016001820154815260200160028201548152602001600382015481526020016004820154815260200160058201548152505090505f815f0151146108eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108e290611f3a565b60405180910390fd5b6040518060c001604052808b81526020015f81526020018a81526020018881526020015f8152602001898152509050805f803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8e73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8d73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f820151815f01556020820151816001015560408201518160020155606082015181600301556080820151816004015560a082015181600501559050508a73ffffffffffffffffffffffffffffffffffffffff166323b872dd3330845f01516040518463ffffffff1660e01b8152600401610a4a93929190611f58565b6020604051808303815f875af1158015610a66573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610a8a9190611fc2565b50505050505050505050505050565b805f8111610adc576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ad390611d67565b60405180910390fd5b815f808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20600401819055507f321c73989a5585f52e30bb8aabd608154125bf4413611a3b4d8b97105d9c6dc084338585604051610bcb9493929190611fed565b60405180910390a150505050565b5f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015481526020016001820154815260200160028201548152602001600382015481526020016004820154815260200160058201548152505090508060400151421015610d1b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d129061207a565b60405180910390fd5b5f81602001511115610da8578273ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8383602001516040518363ffffffff1660e01b8152600401610d66929190612098565b6020604051808303815f875af1158015610d82573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610da69190611fc2565b505b5f815f015111610ded576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610de490612109565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff1663a9059cbb33835f01516040518363ffffffff1660e01b8152600401610e2b929190612098565b6020604051808303815f875af1158015610e47573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e6b9190611fc2565b505f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f01819055505f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2060010181905550505050565b805f8111611027576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161101e90611d67565b60405180910390fd5b3384845f805f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015481526020016001820154815260200160028201548152602001600382015481526020016004820154815260200160058201548152505090508060400151421061116b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161116290612171565b60405180910390fd5b5f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f820154815260200160018201548152602001600282015481526020016003820154815260200160048201548152602001600582015481525050905086815f01516112759190611b56565b5f803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f01819055508773ffffffffffffffffffffffffffffffffffffffff166323b872dd33308a6040518463ffffffff1660e01b815260040161136a93929190611f58565b6020604051808303815f875af1158015611386573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906113aa9190611fc2565b50505050505050505050565b5f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f820154815260200160018201548152602001600282015481526020016003820154815260200160048201548152602001600582015481525050905080608001515f803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20600301819055505f805f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20600401819055507f8c86180e4992276e0f056bb49eed13ea7b192a0d9e795a7f219b5b281c22d78e33338484606001516040516116629493929190611fed565b60405180910390a1505050565b5f602052825f5260405f20602052815f5260405f20602052805f5260405f205f925092505050805f0154908060010154908060020154908060030154908060040154908060050154905086565b5f805f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f206040518060c00160405290815f82015481526020016001820154815260200160028201548152602001600382015481526020016004820154815260200160058201548152505090505f8160200151116117fd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016117f4906121d9565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff1663a9059cbb3383602001516040518363ffffffff1660e01b815260040161183c929190612098565b6020604051808303815f875af1158015611858573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061187c9190611fc2565b50505050565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6118af82611886565b9050919050565b6118bf816118a5565b81146118c9575f80fd5b50565b5f813590506118da816118b6565b92915050565b5f819050919050565b6118f2816118e0565b81146118fc575f80fd5b50565b5f8135905061190d816118e9565b92915050565b5f805f6060848603121561192a57611929611882565b5b5f611937868287016118cc565b9350506020611948868287016118cc565b9250506040611959868287016118ff565b9150509250925092565b5f805f805f8060c0878903121561197d5761197c611882565b5b5f61198a89828a016118cc565b965050602061199b89828a016118cc565b95505060406119ac89828a016118ff565b94505060606119bd89828a016118ff565b93505060806119ce89828a016118ff565b92505060a06119df89828a016118ff565b9150509295509295509295565b5f8060408385031215611a0257611a01611882565b5b5f611a0f858286016118cc565b9250506020611a20858286016118cc565b9150509250929050565b5f805f60608486031215611a4157611a40611882565b5b5f611a4e868287016118cc565b9350506020611a5f868287016118cc565b9250506040611a70868287016118cc565b9150509250925092565b611a83816118e0565b82525050565b5f60c082019050611a9c5f830189611a7a565b611aa96020830188611a7a565b611ab66040830187611a7a565b611ac36060830186611a7a565b611ad06080830185611a7a565b611add60a0830184611a7a565b979650505050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f611b1f826118e0565b9150611b2a836118e0565b9250828202611b38816118e0565b91508282048414831517611b4f57611b4e611ae8565b5b5092915050565b5f611b60826118e0565b9150611b6b836118e0565b9250828201905080821115611b8357611b82611ae8565b5b92915050565b5f82825260208201905092915050565b7f4f776564416d6f756e7420697320626967676572207468616e206368616e6e655f8201527f6c277320617661696c61626c652066756e647300000000000000000000000000602082015250565b5f611bf3603383611b89565b9150611bfe82611b99565b604082019050919050565b5f6020820190508181035f830152611c2081611be7565b9050919050565b7f4e657720446561646c696e6520746f6f2073686f7274000000000000000000005f82015250565b5f611c5b601683611b89565b9150611c6682611c27565b602082019050919050565b5f6020820190508181035f830152611c8881611c4f565b9050919050565b7f4e657720446561646c696e65206973206c657373207468616e2063757272656e5f8201527f7420646561646c696e6500000000000000000000000000000000000000000000602082015250565b5f611ce9602a83611b89565b9150611cf482611c8f565b604082019050919050565b5f6020820190508181035f830152611d1681611cdd565b9050919050565b7f63616e277420616363706574207a65726f20617320612076616c7565000000005f82015250565b5f611d51601c83611b89565b9150611d5c82611d1d565b602082019050919050565b5f6020820190508181035f830152611d7e81611d45565b9050919050565b611d8e816118a5565b82525050565b5f604082019050611da75f830185611d85565b611db46020830184611d85565b9392505050565b5f81519050611dc9816118e9565b92915050565b5f60208284031215611de457611de3611882565b5b5f611df184828501611dbb565b91505092915050565b7f616c6c6f77616e636520646f6573206e6f74206d6174636820737065636966695f8201527f656420616d6f756e740000000000000000000000000000000000000000000000602082015250565b5f611e54602983611b89565b9150611e5f82611dfa565b604082019050919050565b5f6020820190508181035f830152611e8181611e48565b9050919050565b7f446561646c696e652045787069726564000000000000000000000000000000005f82015250565b5f611ebc601083611b89565b9150611ec782611e88565b602082019050919050565b5f6020820190508181035f830152611ee981611eb0565b9050919050565b7f4368616e6e656c20616c726561647920637265617465640000000000000000005f82015250565b5f611f24601783611b89565b9150611f2f82611ef0565b602082019050919050565b5f6020820190508181035f830152611f5181611f18565b9050919050565b5f606082019050611f6b5f830186611d85565b611f786020830185611d85565b611f856040830184611a7a565b949350505050565b5f8115159050919050565b611fa181611f8d565b8114611fab575f80fd5b50565b5f81519050611fbc81611f98565b92915050565b5f60208284031215611fd757611fd6611882565b5b5f611fe484828501611fae565b91505092915050565b5f6080820190506120005f830187611d85565b61200d6020830186611d85565b61201a6040830185611d85565b6120276060830184611a7a565b95945050505050565b7f446561646c696e65206e6f7420726561636865642079657400000000000000005f82015250565b5f612064601883611b89565b915061206f82612030565b602082019050919050565b5f6020820190508181035f83015261209181612058565b9050919050565b5f6040820190506120ab5f830185611d85565b6120b86020830184611a7a565b9392505050565b7f456d707479204368616e6e656c000000000000000000000000000000000000005f82015250565b5f6120f3600d83611b89565b91506120fe826120bf565b602082019050919050565b5f6020820190508181035f830152612120816120e7565b9050919050565b7f4368616e6e656c204578706972656400000000000000000000000000000000005f82015250565b5f61215b600f83611b89565b915061216682612127565b602082019050919050565b5f6020820190508181035f8301526121888161214f565b9050919050565b7f5a65726f204f776e6572736869700000000000000000000000000000000000005f82015250565b5f6121c3600e83611b89565b91506121ce8261218f565b602082019050919050565b5f6020820190508181035f8301526121f0816121b7565b905091905056fea2646970667358221220ebb42c2f18f34ce97597de7340ead0aeebf524a5aac63a8becff6d82b0cf912664736f6c63430008150033",
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
