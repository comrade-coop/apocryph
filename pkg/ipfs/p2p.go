// this package contains ipfs helper functions
package ipfs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/ipfs/kubo/client/rpc"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

type P2pApi struct {
	http     *rpc.HttpApi
	nodeAddr multiaddr.Multiaddr
}

func NewP2pApi(http *rpc.HttpApi, nodeAddr multiaddr.Multiaddr) *P2pApi {
	nodeBaseAddr := nodeAddr
	for {
		var lastComponent *multiaddr.Component
		nodeBaseAddr, lastComponent = multiaddr.SplitLast(nodeBaseAddr)
		if lastComponent.Protocol().Code == multiaddr.P_TCP {
			break
		}
		if nodeBaseAddr == nil {
			nodeBaseAddr = nodeAddr
			break
		}
	}

	return &P2pApi{
		http:     http,
		nodeAddr: nodeBaseAddr,
	}
}

type IpfsListener struct {
	net.Listener
	*ExposedEndpoint
}

func (api *P2pApi) Listen(protocol string) (*IpfsListener, error) {
	fakeConn, err := manet.Dial(api.nodeAddr.Encapsulate(multiaddr.StringCast("/udp/12345")))
	if err != nil {
		return nil, fmt.Errorf("failed to fake-dial ipfs node: %w", err)
	}
	localAddr := fakeConn.LocalAddr()
	fakeConn.Close()

	localHost, _, err := net.SplitHostPort(localAddr.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse local address: %w", err)
	}

	listener, err := net.Listen("tcp", localHost+":0")
	if err != nil {
		return nil, fmt.Errorf("failed to start tcp listener: %w", err)
	}

	multiaddr, err := manet.FromNetAddr(listener.Addr())
	if err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to convert local tcp address to multiaddr: %w", err)
	}

	service, err := api.ExposeEndpoint(protocol, multiaddr, ConflictExistingEndpoint)
	if err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to expose local listener as a service: %w", err)
	}

	fmt.Printf("Listening for requests on: %s\n", multiaddr.String())
	fmt.Printf("Listening for requests on: %s\n", service.ListenAddress)

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
	connection, err := api.ForwardConnection(protocol, api.nodeAddr.Encapsulate(multiaddr.StringCast("/tcp/0")), target)
	if err != nil {
		return nil, fmt.Errorf("failed start forwarding connection: %w", err)
	}

	multiaddr, err := multiaddr.NewMultiaddr(connection.ListenAddress)
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to parse connection address: %w", err)
	}

	netAddr, err := manet.ToNetAddr(multiaddr)
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to convert connection address to network addresss: %w", err)
	}

	return &IpfsAddr{netAddr, connection}, nil
}

type ExposedEndpoint struct {
	P2PListenerInfoOutput
	api *P2pApi
}

type ExistingOption int

const (
	ConflictExistingEndpoint ExistingOption = 0
	ReturnExistingEndpoint   ExistingOption = 1
)

func (api *P2pApi) ExposeEndpoint(protocol string, endpoint multiaddr.Multiaddr, existing ExistingOption) (*ExposedEndpoint, error) {
	ctx := context.Background()
	if existing == ReturnExistingEndpoint {
		listener, err := api.findListenerForAddress(protocol, endpoint)
		if err == nil {
			return &ExposedEndpoint{listener, api}, nil
		}
	}
	// Ipfs automatically handles ConflictExistingEndpoint for us

	request := api.http.Request("p2p/listen", protocol, endpoint.String())
	response, err := request.Send(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to expose endpoint: %w", err)
	}
	if response.Error != nil {
		return nil, fmt.Errorf("failed to expose endpoint: %w", response.Error)
	}

	listener, err := api.findListenerForAddress(protocol, endpoint)
	if err != nil {
		return nil, err
	}
	return &ExposedEndpoint{listener, api}, nil
}

func (i *ExposedEndpoint) Close() error {
	ctx := context.Background()
	request := i.api.http.Request("p2p/close").Option("target-address", i.TargetAddress)
	response, err := request.Send(ctx)
	if err != nil {
		return fmt.Errorf("failed to close endpoint: %w", err)
	}
	if response.Error != nil {
		return fmt.Errorf("failed to close endpoint: %w", response.Error)
	}
	return nil
}

type ForwardedConnection struct {
	P2PListenerInfoOutput
	api *P2pApi
}

func (api *P2pApi) ForwardConnection(protocol string, endpoint multiaddr.Multiaddr, target multiaddr.Multiaddr) (*ForwardedConnection, error) {

	ctx := context.Background()

	request := api.http.Request("p2p/forward", protocol, endpoint.String(), target.String())
	response, err := request.Send(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to forward connection: %w", err)
	}
	if response.Error != nil {
		return nil, fmt.Errorf("failed to forward connection: %w", response.Error)
	}

	listener, err := api.findListenerForAddress(protocol, target)
	if err != nil {
		return nil, err
	}
	return &ForwardedConnection{listener, api}, nil
}

func (f *ForwardedConnection) Close() error {
	ctx := context.Background()
	request := f.api.http.Request("p2p/close").Option("listen-address", f.ListenAddress)
	response, err := request.Send(ctx)
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	if response.Error != nil {
		return fmt.Errorf("failed to close connection: %w", response.Error)
	}
	return nil
}

type P2PListenerInfoOutput struct {
	Protocol      string
	ListenAddress string
	TargetAddress string
}

type P2PLsOutput struct {
	Listeners []P2PListenerInfoOutput
}

func (api *P2pApi) findListenerForAddress(protocol string, targetAddress multiaddr.Multiaddr) (P2PListenerInfoOutput, error) {
	ctx := context.Background()

	requestLs := api.http.Request("p2p/ls")
	responseLs, err := requestLs.Send(ctx)
	if err != nil {
		return P2PListenerInfoOutput{}, fmt.Errorf("failed to list listeners: %w", err)
	}
	if responseLs.Error != nil {
		return P2PListenerInfoOutput{}, fmt.Errorf("failed to list listeners: %w", responseLs.Error)
	}

	listeners := &P2PLsOutput{}
	err = json.NewDecoder(responseLs.Output).Decode(listeners)
	if err != nil {
		return P2PListenerInfoOutput{}, fmt.Errorf("failed to decode listeners: %w", err)
	}

	for _, listener := range listeners.Listeners {
		if listener.Protocol == protocol && listener.TargetAddress == targetAddress.String() {
			return listener, nil
		}
	}

	return P2PListenerInfoOutput{}, errors.New("could not find target address in the list of listeners")
}
