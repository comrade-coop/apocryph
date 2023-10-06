package types

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
)

func NewAccountIDFromSS58(accountString string) (uint16, *types.AccountID, error) {
	network, accountIdBytes, err := subkey.SS58Decode(accountString)
	if err != nil {
		return 0, nil, err
	}

	accountId, err := types.NewAccountID(accountIdBytes)

	return network, accountId, err
}

func AccountIDToSS58(network uint16, accountId *types.AccountID) (string, error) {
	return subkey.SS58Encode(accountId.ToBytes(), network), nil
}
