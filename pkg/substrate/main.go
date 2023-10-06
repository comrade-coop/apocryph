package extrpc

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/comrade-coop/trusted-pods/pkg/substrate/rpc"
)

type SubstrateAPI struct {
	gsrpc.SubstrateAPI
	RPC *rpc.RPC
}

func WrapSubstrateAPI(origApi *gsrpc.SubstrateAPI) (*SubstrateAPI, error) {
	newRPC, err := rpc.NewRPC(origApi)
	if err != nil {
		return nil, err
	}

	return &SubstrateAPI{
		*origApi,
		newRPC,
	}, nil
}

func NewSubstrateAPI(url string) (*SubstrateAPI, error) {
	api, err := gsrpc.NewSubstrateAPI(url)
	if err != nil {
		return nil, err
	}

	return WrapSubstrateAPI(api)
}
