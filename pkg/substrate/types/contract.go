package types

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
)

type Balance = types.U128
type ContractSelector = types.U32

type ContractInfo struct {
	TrieId             types.Bytes
	DepositAccount     types.AccountID // NOTE: Gone in https://github.com/paritytech/polkadot-sdk/commit/9c5315961e2f75b092c431bebe8cc1b4481853db
	CodeHash           types.Hash
	StorageBytes       types.U32
	StorageItems       types.U32
	StorageByteDeposit Balance
	StorageItemDeposit Balance
	StorageBaseDeposit Balance
}

type ContractCallRequest struct {
	Origin              types.AccountID
	Dest                types.AccountID
	Value               Balance
	GasLimit            Option[types.Weight]
	StorageDepositLimit Option[Balance]
	InputData           ContractCallInputData
}

type ContractCallInputData []byte

func EncodeInputData(selector ContractSelector, args ...interface{}) (ContractCallInputData, error) {
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, uint32(selector))
	enc := scale.NewEncoder(&buffer)
	for _, arg := range args {
		err := enc.Encode(arg)
		if err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

type ContractExecResult struct {
	GasConsumed            types.Weight
	GasRequired            types.Weight
	StorageDepositIsCharge types.Bool
	StorageDepositBalance  Balance
	DebugMessage           types.Bytes
	Result                 Result[ExecReturnValue, types.DispatchError]
}

type ExecReturnValue struct {
	ReturnFlags types.U32
	Data        types.Bytes
}

func (r ContractExecResult) DecodeResult(result interface{}) error {
	if r.Result.IsError {
		return fmt.Errorf("Contract execution error: %v", r.Result.Error)
	}
	return codec.Decode(r.Result.Value.Data, result)
}
