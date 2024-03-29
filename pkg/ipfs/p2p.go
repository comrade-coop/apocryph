// SPDX-License-Identifier: GPL-3.0

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

// A wrapper around the experimental P2P api of kubo.
type P2pApi struct {
	http     *rpc.HttpApi
	nodeAddr multiaddr.Multiaddr
}

// Create a new [P2pApi] given an HttpApi and the IP address where the kubo node resides.
// The node address is passed as a [multiaddr.Multiaddr] and is expected to be either a multiaddr that will be accepted by [manet.ToNetAddr] or the same plus a TCP component at the end which will be stripped off.
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

// A [net.Listener] which has been enhanced with a Close() method that can be used to stop listening for the given connection.
type IpfsListener struct {
	net.Listener
	*ExposedEndpoint
}

// Start listening for a certain libp2p protocol on the kubo node, and forward connections to the returned listener.
// Make sure to close the resulting listener once you are done using it, for garbage-collection.
// If another Listen request has been made to the kubo node for the same protocol, this method will produce an error.
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

// Close the local listener as well as the endpoint on the kubo node.
func (l *IpfsListener) Close() error {
	return errors.Join(l.ExposedEndpoint.Close(), l.Listener.Close())
}

// A [net.Addr] which has been enhanced with a Close() method that can be used to stop forwarding requests for the given connection.
// You may convert the address to a String() representing the TCP host and port; however, note that this address is only accessible locally and may not be resolvable on other hosts.
type IpfsAddr struct {
	net.Addr
	*ForwardedConnection
}

// Connect to a given peer over the specified libp2p protocol through the kubo node.
// Make sure to close the resulting address once you are done using it, for garbage-collection.
func (api *P2pApi) ConnectTo(protocol string, target peer.ID) (*IpfsAddr, error) {
	addrs, err := peer.AddrInfoToP2pAddrs(&peer.AddrInfo{ID: target})
	if err != nil {
		return nil, err
	}
	return api.Connect(protocol, addrs[0])
}

// Connect to a given peer (as a /p2p/... multiaddr) over the specified libp2p protocol through the kubo node.
// Make sure to close the resulting address once you are done using it, for garbage-collection
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
		return nil, fmt.Errorf("failed to convert connection address to network address: %w", err)
	}

	return &IpfsAddr{netAddr, connection}, nil
}

// Represents the low-level counterpart of a [IpfsListener]; a local endpoint that the kubo node is exposing to the world on a given libp2p protocol.
type ExposedEndpoint struct {
	P2PListenerInfoOutput
	api *P2pApi
}

// Configures how [ExposeEndpoint] should behave when there is an endpoint for the same protocol already registered by another ExposeEndpoint call.
type ExistingOption int

const (
	ConflictExistingEndpoint ExistingOption = 0 // Throw an error (default)
	ReturnExistingEndpoint   ExistingOption = 1 // Return the previously existing endpoint, at which point you can close it or compare it with the desired endpoint.
)

// Expose an endpoint accessible by the kubo node to the rest of the network over a libp2p protocol. (any peer -> us)
// The endpoint passed in should be accessible _from_ the kubo node; this methods makes no attempts to fix the network address given. (i.e. this means that passing localhost:XXX for the endpoint will result in the kubo node exposing its own localhost and not the localhost of the machine making the request.)
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

// Stop listening on the exposed endpoint, removing it from the node, and making the desired protocol available for future ExposeEndpoint calls.
// Note that this is not automatically garbage-collected in cases of process crash; if a lingering p2p connection remains, `ipfs p2p close -a` would close all forwarded and exposed connections, while `ipfs p2p close -l /p2p/<own node id> -p <protocol>` would close the specific exposed endpoint.
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

// Represents the low-level counterpart of a [IpfsAddr]; a local address that the kubo node is listening on and forwarding requests to over a given libp2p protocol.
type ForwardedConnection struct {
	P2PListenerInfoOutput
	api *P2pApi
}

// Listen on an address on the kubo node and forward connections to a particular peer over a libp2p protocol. (us -> specific peer)
// The endpoint passed in should be on the kubo node itself; this methods makes no attempts to fix the network address given. (i.e. this means that passing localhost:XXX for the endpoint will result in the kubo node listening on its own localhost -- which is typically precisely the desired effect)
// The target peer should be a /p2p/... multiaddr.
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

// Stop forwarding connections from the specified address, making the endpoint accessible for futere ForwardConnection calls.
// Note that this is not automatically garbage-collected in cases of process crash; if a lingering p2p connection remains, `ipfs p2p close -a` would close all forwarded and exposed connections, while `ipfs p2p close -t /p2p/<target node id> -p <protocol>` would close the specific exposed endpoint.
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

// Copy of [github.com/ipfs/kubo/core/commands.P2PListenerInfoOutput]
type P2PListenerInfoOutput struct {
	Protocol      string
	ListenAddress string
	TargetAddress string
}

// Copy of [github.com/ipfs/kubo/core/commands.P2PLsOutput]
type P2PLsOutput struct {
	Listeners []P2PListenerInfoOutput
}

// Look up the listener for a given protocol+targetAddress pair, used since p2p/forward ad p2p/expose don't return any info on the P2P listener registered in kubo.
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
