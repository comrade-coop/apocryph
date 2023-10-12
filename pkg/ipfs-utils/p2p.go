// this package contains ipfs helper functions
package ipfs_utils

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"github.com/ipfs/kubo/client/rpc"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

type P2pApi rpc.HttpApi

func NewP2pApi(h *rpc.HttpApi) *P2pApi {
	return (*P2pApi)(h)
}

type IpfsListener struct {
	net.Listener
	*ExposedEndpoint
}

func (api *P2pApi) Listen(protocol string) (*IpfsListener, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	multiaddr, err := manet.FromNetAddr(listener.Addr())
	if err != nil {
		listener.Close()
		return nil, err
	}

	service, err := api.ExposeEndpoint(protocol, multiaddr)
	if err != nil {
		listener.Close()
		return nil, err
	}

	return &IpfsListener{listener, service}, nil
}

func (l *IpfsListener) Close() error {
	return errors.Join(l.ExposedEndpoint.Close(), l.Listener.Close())
}

type IpfsAddr struct {
	net.Addr
	*ForwardedConnection
}

func (api *P2pApi) ConnectTo(protocol string, target peer.ID) (*IpfsAddr, error) {
	addrs, err := peer.AddrInfoToP2pAddrs(&peer.AddrInfo{ID: target})
	if err != nil {
		return nil, err
	}
	return api.Connect(protocol, addrs[0])
}

func (api *P2pApi) Connect(protocol string, target multiaddr.Multiaddr) (*IpfsAddr, error) {
	connection, err := api.ForwardConnection(protocol, multiaddr.StringCast("/ip4/127.0.0.1/tcp/0"), target)
	if err != nil {
		return nil, err
	}

	multiaddr, err := multiaddr.NewMultiaddr(connection.ListenAddress)
	if err != nil {
		connection.Close()
		return nil, err
	}

	netAddr, err := manet.ToNetAddr(multiaddr)
	if err != nil {
		connection.Close()
		return nil, err
	}

	return &IpfsAddr{netAddr, connection}, nil
}

type ExposedEndpoint struct {
	TargetAddress multiaddr.Multiaddr
	api           *P2pApi
}

func (api *P2pApi) ExposeEndpoint(protocol string, endpoint multiaddr.Multiaddr) (*ExposedEndpoint, error) {
	ctx := context.Background()
	request := (*rpc.HttpApi)(api).Request("p2p/listen", protocol, endpoint.String())
	response, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Unwrap()
	}
	return &ExposedEndpoint{TargetAddress: endpoint, api: api}, nil
}

func (i *ExposedEndpoint) Close() error {
	ctx := context.Background()
	request := (*rpc.HttpApi)(i.api).Request("p2p/close").Option("target-address", i.TargetAddress.String())
	response, err := request.Send(ctx)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Unwrap()
	}
	return nil
}

type ForwardedConnection struct {
	P2PListenerInfoOutput
	api *P2pApi
}

type P2PListenerInfoOutput struct {
	Protocol      string
	ListenAddress string
	TargetAddress string
}

type P2PLsOutput struct {
	Listeners []P2PListenerInfoOutput
}

func (api *P2pApi) ForwardConnection(protocol string, endpoint multiaddr.Multiaddr, target multiaddr.Multiaddr) (*ForwardedConnection, error) {
	ctx := context.Background()

	request := (*rpc.HttpApi)(api).Request("p2p/forward", protocol, endpoint.String(), target.String())
	response, err := request.Send(ctx)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Unwrap()
	}

	requestLs := (*rpc.HttpApi)(api).Request("p2p/ls")
	responseLs, err := requestLs.Send(ctx)
	if err != nil {
		return nil, err
	}
	if responseLs.Error != nil {
		return nil, responseLs.Error.Unwrap()
	}

	listeners := &P2PLsOutput{}
	err = json.NewDecoder(responseLs.Output).Decode(listeners)
	if err != nil {
		return nil, err
	}

	for _, listener := range listeners.Listeners {
		if listener.Protocol == protocol && listener.TargetAddress == target.String() {
			return &ForwardedConnection{listener, api}, nil
		}
	}

	return nil, errors.New("Could not find target address after creating it")
}

func (f *ForwardedConnection) Close() error {
	ctx := context.Background()
	request := (*rpc.HttpApi)(f.api).Request("p2p/close").Option("listen-address", f.ListenAddress)
	response, err := request.Send(ctx)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Unwrap()
	}
	return nil
}
