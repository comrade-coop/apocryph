package main

import (
	"fmt"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/spf13/pflag"
)

var podFlags = &pflag.FlagSet{}
var manifestFormat string

var deploymentFlags = &pflag.FlagSet{}
var deploymentFormat string
var providerPeer string
var providerEthAddress string

var uploadFlags = &pflag.FlagSet{}
var ipfsApi string
var uploadImages bool
var uploadSecrets bool

var fundFlags = &pflag.FlagSet{}
var ethereumRpc string
var publisherKey string
var paymentContractAddress string
var podId string
var unlockTime int64
var funds string
var debugMintFunds bool

var syncFlags = &pflag.FlagSet{}

var _ = func() error {
	podFlags := podCmd.PersistentFlags()

	podFlags.StringVar(&manifestFormat, "format", "", fmt.Sprintf("Manifest format. One of %v (leave empty to auto-detect)", pb.FormatNames))

	deploymentFlags.StringVar(&manifestFormat, "deployment-format", "", fmt.Sprintf("Deployment format. One of %v (leave empty to auto-detect)", pb.FormatNames))
	deploymentFlags.StringVar(&providerPeer, "provider", "", "provider peer id")
	deploymentFlags.StringVar(&providerEthAddress, "provider-eth", "", "provider public address")

	uploadFlags.StringVar(&ipfsApi, "ipfs", "", "multiaddr where the ipfs/kubo api can be accessed (leave blank to use the daemon running in IPFS_PATH)")
	uploadFlags.BoolVar(&uploadImages, "upload-images", true, "upload images")
	uploadFlags.BoolVar(&uploadSecrets, "upload-secrets", true, "upload secrets")

	fundFlags.StringVar(&ethereumRpc, "ethereum-rpc", "http://127.0.0.1:8545", "ethereum rpc node")
	fundFlags.StringVar(&publisherKey, "ethereum-key", "", "account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")
	fundFlags.StringVar(&paymentContractAddress, "payment-contract", "", "payment contract address")
	fundFlags.StringVar(&podId, "pod-id", "", "pod id")
	fundFlags.StringVar(&funds, "funds", "0", "intial funds")
	fundFlags.BoolVar(&debugMintFunds, "mint-funds", false, "Attempt minting funds with a mint(amount) call on the token")
	fundFlags.Int64Var(&unlockTime, "unlock-time", 5*60, "time for unlocking tokens (in seconds)")

	syncFlags.AddFlag(uploadFlags.Lookup("ipfs"))

	return nil
}()
