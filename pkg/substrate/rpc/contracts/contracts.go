package contracts

import (
	"errors"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/author"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	tptypes "github.com/comrade-coop/trusted-pods/pkg/substrate/types"
)

type Contracts interface {
	GetCodeHash(address types.AccountID, block *types.Hash) (types.Hash, error)
	QueryContract(address types.AccountID, fromAddress types.AccountID, inputData tptypes.ContractCallInputData, value types.U128, block *types.Hash) (*tptypes.ContractExecResult, error)
	CreateContractCallExtrinsic(address types.AccountID, fromAddress types.AccountID, inputData tptypes.ContractCallInputData, value types.U128, block *types.Hash) (*types.Extrinsic, *tptypes.ContractExecResult, error)

	CallContract(address types.AccountID, from signature.KeyringPair, inputData tptypes.ContractCallInputData, value types.U128) (*author.ExtrinsicStatusSubscription, error)
}

type contracts struct {
	api  *gsrpc.SubstrateAPI
	meta *types.Metadata
}

func NewContracts(api *gsrpc.SubstrateAPI) (Contracts, error) {
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	return &contracts{api, meta}, nil
}

func (c *contracts) GetCodeHash(address types.AccountID, block *types.Hash) (types.Hash, error) {
	result := tptypes.ContractInfo{}

	key, err := types.CreateStorageKey(c.meta, "Contracts", "ContractInfoOf", address.ToBytes())
	if err != nil {
		return types.Hash{}, err
	}

	var found bool
	if block == nil {
		found, err = c.api.RPC.State.GetStorageLatest(key, &result)
	} else {
		found, err = c.api.RPC.State.GetStorage(key, &result, *block)
	}
	if err != nil {
		return types.Hash{}, err
	}
	if !found {
		return types.Hash{}, errors.New("Contract not found")
	}

	return result.CodeHash, nil
}

func (c *contracts) QueryContract(address types.AccountID, fromAddress types.AccountID, inputData tptypes.ContractCallInputData, value types.U128, block *types.Hash) (*tptypes.ContractExecResult, error) {
	req, err := codec.EncodeToHex(tptypes.ContractCallRequest{
		Origin:    fromAddress,
		Dest:      address,
		Value:     value,
		InputData: inputData,
	})
	if err != nil {
		return nil, err
	}

	var res string
	if block == nil {
		err = c.api.Client.Call(&res, "state_call", "ContractsApi_call", req)
	} else {
		err = c.api.Client.Call(&res, "state_call", "ContractsApi_call", req, block.Hex())
	}
	if err != nil {
		return nil, err
	}

	var execResult tptypes.ContractExecResult

	err = codec.DecodeFromHex(res, &execResult)
	if err != nil {
		return nil, err
	}

	return &execResult, nil
}

func (c *contracts) CreateContractCallExtrinsic(address types.AccountID, fromAddress types.AccountID, inputData tptypes.ContractCallInputData, value types.U128, block *types.Hash) (*types.Extrinsic, *tptypes.ContractExecResult, error) {
	execResult, err := c.QueryContract(address, fromAddress, inputData, value, block)
	if err != nil {
		return nil, nil, err
	}

	dest := types.MultiAddress{IsID: true, AsID: address}
	gasPrice := execResult.GasRequired // TODO: Give a bit of margin?
	storageDepositLimit := tptypes.Option[types.UCompact]{HasValue: false}

	call, err := types.NewCall(c.meta, "Contracts.call", dest, types.NewUCompact(value.Int), gasPrice, storageDepositLimit, inputData)
	if err != nil {
		return nil, nil, err
	}

	extrinsic := types.NewExtrinsic(call)

	return &extrinsic, execResult, nil
}

func (c *contracts) CallContract(address types.AccountID, from signature.KeyringPair, inputData tptypes.ContractCallInputData, value types.U128) (*author.ExtrinsicStatusSubscription, error) {
	extrinsic, _, err := c.CreateContractCallExtrinsic(address, types.AccountID(from.PublicKey), inputData, value, nil)

	rv, err := c.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}

	genesisHash, err := c.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}

	key, err := types.CreateStorageKey(c.meta, "System", "Account", from.PublicKey)
	if err != nil {
		return nil, err
	}

	retries := 5
	var subscription *author.ExtrinsicStatusSubscription
	for {
		// Via: https://github.com/centrifuge/go-substrate-rpc-client/blob/master/teste2e/author_submit_and_watch_extrinsic_test.go
		var accountInfo types.AccountInfo
		ok, err := c.api.RPC.State.GetStorageLatest(key, &accountInfo)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("Caller account not found")
		}

		err = extrinsic.Sign(from, types.SignatureOptions{
			// BlockHash:          blockHash,
			BlockHash:          genesisHash, // BlockHash needs to == GenesisHash if era is immortal.
			Era:                types.ExtrinsicEra{IsMortalEra: false},
			GenesisHash:        genesisHash,
			Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
			SpecVersion:        rv.SpecVersion,
			Tip:                types.NewUCompactFromUInt(0),
			TransactionVersion: rv.TransactionVersion,
		})
		if err != nil {
			return nil, err
		}

		subscription, err = c.api.RPC.Author.SubmitAndWatchExtrinsic(*extrinsic)
		if err != nil {
			retries--
			if retries <= 0 {
				return nil, err
			}
			continue
		}

		return subscription, nil
	}
}
