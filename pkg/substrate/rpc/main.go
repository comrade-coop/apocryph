package rpc

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc"
	"github.com/comrade-coop/trusted-pods/pkg/substrate/rpc/contracts"
)

type RPC struct {
	rpc.RPC
	Contracts contracts.Contracts
}

func NewRPC(api *gsrpc.SubstrateAPI) (*RPC, error) {
	c, err := contracts.NewContracts(api)
	if err != nil {
		return nil, err
	}
	return &RPC{
		*api.RPC,
		c,
	}, nil
}
