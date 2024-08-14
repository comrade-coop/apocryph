package autoscaler

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/comrade-coop/apocryph/pkg/abi"
	"github.com/comrade-coop/apocryph/pkg/constants"
	"github.com/comrade-coop/apocryph/pkg/ethereum"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	pbcon "github.com/comrade-coop/apocryph/pkg/proto/protoconnect"
	tpraft "github.com/comrade-coop/apocryph/pkg/raft"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/raft"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const RAFT_P2P_PORT = 32500
const RAFT_PATH = ""

type AutoScalerServer struct {
	pbcon.UnimplementedAutoscalerServiceHandler
	node           *tpraft.RaftNode
	store          *tpraft.KVStore
	peers          []string
	started        bool
	p2pHost        host.Host
	nodeGateway    string
	ChannelManager *PaymentChannelManager
	MainLoop       func(*AutoScalerServer)
}

type PaymentChannelManager struct {
	Publisher  common.Address
	Provider   common.Address
	PodId      common.Hash
	EthClient  *ethclient.Client
	Transactor *bind.TransactOpts
	Payment    *abi.Payment
}

func NewAutoSalerServer(ethereumRpc string, p2pHost host.Host) (*AutoScalerServer, error) {

	ethClient, err := ethereum.GetClient(ethereumRpc)
	if err != nil {
		return nil, err
	}

	key := os.Getenv(constants.PRIVATE_KEY)
	paymentAddress := common.HexToAddress(os.Getenv(constants.PAYMENT_ADDR_KEY))
	publisherAddress := common.HexToAddress(os.Getenv(constants.PUBLISHER_ADDR_KEY))
	providerAddress := common.HexToAddress(os.Getenv(constants.PROVIDER_ADDR_KEY))
	podId := common.HexToHash((os.Getenv(constants.POD_ID_KEY)))

	log.Printf("ENV Variables: Payment_Address: %v, Publisher Address: %v, ProviderAddress: %v, podId: %v\n", paymentAddress, publisherAddress, providerAddress, podId)

	privateKey, err := ethereum.DecodePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("Failed decoding private Key: %v", err)
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}

	payment, err := abi.NewPayment(paymentAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &AutoScalerServer{
		p2pHost: p2pHost,
		ChannelManager: &PaymentChannelManager{
			Publisher:  publisherAddress,
			Provider:   providerAddress,
			PodId:      podId,
			EthClient:  ethClient,
			Transactor: transactor,
			Payment:    payment,
		},
	}, nil
}

func (s *AutoScalerServer) TriggerNode(c context.Context, req *connect.Request[pb.ConnectClusterRequest]) (*connect.Response[pb.TriggerNodeResponse], error) {
	if s.started == false {
		log.Println("Node Triggered")
		go s.BoostrapCluster(req)
		s.started = true
	}
	return connect.NewResponse(&pb.TriggerNodeResponse{PeerID: s.p2pHost.ID().String()}), nil
}

func (s *AutoScalerServer) BoostrapCluster(req *connect.Request[pb.ConnectClusterRequest]) error {

	peerIDs, err := s.FetchPeerIDsFromServers(req)
	if err != nil {
		return fmt.Errorf("Failed to fetch PeerIDs: %v", err)
	}

	var peers []*peer.AddrInfo
	for serverAddr, addr := range peerIDs {
		addr := GetMultiAddr(serverAddr, addr)
		if addr == "" {
			return fmt.Errorf("Failed to parse server address: %v", serverAddr)
		}

		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return fmt.Errorf("Failed to parse multiaddr: %v: %v", addr, err)
		}

		peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			return fmt.Errorf("Failed to extract peer info from %v, Error: %v", addr, err)
		}

		peers = append(peers, peerInfo)
		log.Printf("Added Peer %v ID:%s \n", peerInfo.Addrs, peerInfo.ID.String())
	}

	node, err := tpraft.NewRaftNode(s.p2pHost, peers, RAFT_PATH)
	if err != nil {
		log.Println("Error:Could not Creat Raft Node")
		return fmt.Errorf("Failed Creating Raft Node %v\n", err)
	}

	// create the KVStore
	store, err := tpraft.NewKVStore(node)
	s.node = node
	s.store = store
	s.peers = req.Msg.Servers
	s.nodeGateway = req.Msg.NodeGateway

	err = s.waitLeaderElection(req.Msg.Timeout)
	if err != nil {
		return err
	}

	return nil
}

func (s *AutoScalerServer) ConnectCluster(c context.Context, req *connect.Request[pb.ConnectClusterRequest]) (*connect.Response[pb.ConnectClusterResponse], error) {
	log.Println("Forming a Raft Cluster with the following providers:", req.Msg.Servers)

	if s.started == true {
		return connect.NewResponse(&pb.ConnectClusterResponse{
			Success: false,
			Error:   "Server Already Started\n",
		}), nil
	}

	err := s.BoostrapCluster(req)
	if err != nil {
		return connect.NewResponse(&pb.ConnectClusterResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed Bootstraping Cluster: %v\n", err),
		}), nil
	}

	s.started = true
	response := &pb.ConnectClusterResponse{Success: true}
	return connect.NewResponse(response), nil
}

func (s *AutoScalerServer) FetchPeerIDsFromServers(req *connect.Request[pb.ConnectClusterRequest]) (map[string]string, error) {
	peerIDs := make(map[string]string)

	for _, addr := range req.Msg.Servers {
		client := pbcon.NewAutoscalerServiceClient(
			http.DefaultClient,
			addr)

		req.Header().Set("Host", "autoscaler.local")

		resp, err := client.TriggerNode(context.Background(), req)
		if err != nil {
			log.Printf("failed to connect to PeerID from server %v: %v", addr, err)
			continue
		}
		peerIDs[addr] = resp.Msg.PeerID
	}
	log.Printf("PeerIDs collected: %v", peerIDs)
	return peerIDs, nil
}

func (s *AutoScalerServer) waitLeaderElection(timeout uint32) error {

	log.Printf("Waiting for leader election with %v seoncds timout ...", timeout)
	time.Sleep(5 * time.Second)

	go s.watchNewStates()

	// using Raft.leaderCh() wont help because it does not count
	// first leader election as a leadership change, therefore from the docs,
	// this is the way to detect the new leader
	obsCh := make(chan raft.Observation, 1)
	observer := raft.NewObserver(obsCh, false, nil)
	s.node.Raft.RegisterObserver(observer)
	defer s.node.Raft.DeregisterObserver(observer)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeoutCh := time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case obs := <-obsCh:
			switch obs.Data.(type) {
			case raft.RaftState:
				if leaderAddr, _ := s.node.Raft.LeaderWithID(); leaderAddr != "" {
					log.Printf("Leader Elected: %v\n", leaderAddr)
					go s.MainLoop(s)
					return nil
				}
			}
		case <-ticker.C:
			if leaderAddr, _ := s.node.Raft.LeaderWithID(); leaderAddr != "" {
				log.Printf("Leader Elected: %v\n", leaderAddr)
				go s.MainLoop(s)
				return nil
			}
		case <-timeoutCh:
			log.Println("timed out waiting for leader")
			return fmt.Errorf("Timed out waiting for leadership election")
		}
	}
}

func (s *AutoScalerServer) watchNewStates() {
	for range s.node.Consensus.Subscribe() {
		newState, err := s.node.Consensus.GetCurrentState()
		if err != nil {
			log.Printf("Failed getting current state %v\n", err)
		}
		log.Printf("State Changed, New State: %v\n", newState)
	}
}

// example of main loop creating a subchannel, then setting the value of a test domain
// with the current node gateway every 10 seconds
func SetAppGatewayExample(s *AutoScalerServer) {
	_, err := s.ChannelManager.Payment.CreateSubChannel(s.ChannelManager.Transactor, s.ChannelManager.Publisher, s.ChannelManager.Provider, s.ChannelManager.PodId, s.ChannelManager.Provider, s.ChannelManager.PodId, big.NewInt(200))
	if err != nil {
		log.Printf("Failed to create subchannel: %v\n", err)
	} else {
		log.Printf("SubChannel Created Succefully \n")
	}
	log.Println("Starting Main Loop:")
	for {
		if s.node.Raft.State() == raft.Leader {
			log.Println("Setting the domain value to the current apocryph node gateway")
			err := s.store.Set("www.test.com", "http://localhost:8080")
			// leadership could be changed right before setting the value
			if _, ok := err.(*tpraft.NotLeaderError); ok {
				newLeaderAddr, newLeaderID := s.node.Raft.LeaderWithID()
				log.Printf("Leadership changed to %v:%v", newLeaderAddr, newLeaderID)
			} else if err != nil {
				log.Printf("Failed setting key:%v: %v\n", "www.test.com", err)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func GetMultiAddr(addr, peerID string) string {
	// Parse the URL
	parsedURL, err := url.Parse(addr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}
	hostPort := parsedURL.Host
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		log.Println("Invalid host:port format")
		return ""
	}
	ip := parts[0]
	// Change the port to the RAFT_P2P_PORT
	return fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", ip, RAFT_P2P_PORT, peerID)
}
