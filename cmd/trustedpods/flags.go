// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"math/big"

	"github.com/comrade-coop/apocryph/pkg/abi"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/spf13/pflag"
)

var podFlags = &pflag.FlagSet{}
var manifestFormat string

var deploymentFlags = &pflag.FlagSet{}
var deploymentFormat string
var providerPeer string
var providerEthAddress string
var registryContractAddress string
var tokenContractAddress string
var expirationOffset int64
var authorize bool

var uploadFlags = &pflag.FlagSet{}
var ipfsApi string
var uploadImages bool
var uploadSecrets bool
var uploadSignatures bool

var signImages bool
var verify bool

var verifyFlags = &pflag.FlagSet{}
var signaturePath string
var hostHeader string

var imageCertificateFlags = &pflag.FlagSet{}
var certificateIdentity string
var certificateOidcIssuer string

var fundFlags = &pflag.FlagSet{}
var ethereumRpc string
var publisherKey string
var paymentContractAddress string
var podId string
var unlockTime int64
var funds string
var debugMintFunds bool

var syncFlags = &pflag.FlagSet{}

var registryFlags = &pflag.FlagSet{}
var cpuPrice string
var ramPrice string
var storagePrice string
var bandwidthEPrice string
var bandwidthInPrice string
var cpuModel string
var teeType string
var tableId string
var region string

var _ = func() error {
	podFlags := podCmd.PersistentFlags()

	podFlags.StringVar(&manifestFormat, "format", "", fmt.Sprintf("Manifest format. One of %v (leave empty to auto-detect)", pb.FormatNames))

	deploymentFlags.StringVar(&manifestFormat, "deployment-format", "", fmt.Sprintf("Deployment format. One of %v (leave empty to auto-detect)", pb.FormatNames))
	deploymentFlags.StringVar(&providerPeer, "provider", "", "provider peer id")
	deploymentFlags.StringVar(&providerEthAddress, "provider-eth", "", "provider public address")
	deploymentFlags.Int64Var(&expirationOffset, "token-expiration", 10, "authentication token expires after token-expiration seconds (expired after 10 seconds by default)")
	deploymentFlags.StringVar(&ipfsApi, "ipfs", "/ip4/127.0.0.1/tcp/5001", "multiaddr where the ipfs/kubo api can be accessed")
	deploymentFlags.BoolVar(&authorize, "authorize", false, "Create a key pair for the application and authorize the returned addresses to control the payment channel")
	deploymentFlags.BoolVar(&verify, "verify", false, "verify the pod images (requires certificate-identity & certificate-oidc-issuer flags)")
	deploymentFlags.BoolVar(&uploadSignatures, "upload-signatures", false, "skip uploading signatures to the registry")
	deploymentFlags.BoolVar(&signImages, "sign-images", false, "sign pod images")

	imageCertificateFlags.StringVar(&certificateIdentity, "certificate-identity", "", "identity used for signing the image")
	imageCertificateFlags.StringVar(&certificateOidcIssuer, "certificate-oidc-issuer", "", "issuer of the oidc")

	uploadFlags.StringVar(&ipfsApi, "ipfs", "/ip4/127.0.0.1/tcp/5001", "multiaddr where the ipfs/kubo api can be accessed")
	uploadFlags.BoolVar(&uploadImages, "upload-images", true, "upload images")
	uploadFlags.BoolVar(&uploadSecrets, "upload-secrets", true, "upload secrets")
	uploadFlags.BoolVar(&signImages, "sign-images", false, "sign pod images (requires certificate identity & issuer flags)")
	uploadFlags.AddFlagSet(imageCertificateFlags)
	uploadFlags.BoolVar(&uploadSignatures, "upload-signatures", false, "skip uploading signatures to the registry")

	verifyFlags.AddFlagSet(imageCertificateFlags)
	verifyFlags.StringVar(&signaturePath, "signature", "", "path to the signature you want to verify")
	verifyFlags.StringVar(&hostHeader, "host-header", "", "the verification host header when passing a tpod ip endpoint to verify")

	fundFlags.StringVar(&ethereumRpc, "ethereum-rpc", "http://127.0.0.1:8545", "ethereum rpc node")
	fundFlags.StringVar(&publisherKey, "ethereum-key", "", "account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")
	fundFlags.StringVar(&paymentContractAddress, "payment-contract", "", "payment contract address")
	fundFlags.StringVar(&podId, "pod-id", "", "pod id")
	fundFlags.StringVar(&funds, "funds", "0", "initial funds")
	fundFlags.BoolVar(&debugMintFunds, "mint-funds", false, "Attempt minting funds with a mint(amount) call on the token")
	fundFlags.Int64Var(&unlockTime, "unlock-time", 5*60, "time for unlocking tokens (in seconds)")

	syncFlags.AddFlag(uploadFlags.Lookup("ipfs"))
	syncFlags.StringVar(&publisherKey, "ethereum-key", "", "account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")

	registryFlags.StringVar(&ipfsApi, "ipfs", "/ip4/127.0.0.1/tcp/5001", "multiaddr where the ipfs/kubo api can be accessed")
	registryFlags.StringVar(&registryContractAddress, "registry-contract", "", "registry contract address")
	registryFlags.StringVar(&tokenContractAddress, "token-contract", "", "token contract address")
	registryFlags.AddFlag(fundFlags.Lookup("payment-contract"))
	registryFlags.StringVar(&cpuPrice, "cpu-price", "", "CPU price")
	registryFlags.StringVar(&ramPrice, "ram-price", "", "RAM price")
	registryFlags.StringVar(&storagePrice, "storage-price", "", "Storage price")
	registryFlags.StringVar(&bandwidthEPrice, "bandwidthE-price", "", "Egress Bandiwdth price")
	registryFlags.StringVar(&bandwidthInPrice, "bandwidthIn-price", "", "Ingress Bandiwdth price")
	registryFlags.StringVar(&cpuModel, "cpu-Model", "", "cpu Model")
	registryFlags.StringVar(&teeType, "tee-Type", "", "tee Type")
	registryFlags.StringVar(&tableId, "id", "", "table id")
	registryFlags.StringVar(&region, "region", "", "filter providers by region, Ex: us-east-8")
	registryFlags.AddFlag(fundFlags.Lookup("ethereum-key"))

	return nil
}()

func getRegistryTableFilter() (*abi.RegistryNewPricingTable, error) {
	result := &abi.RegistryNewPricingTable{}
	ok := true
	if cpuPrice != "" && ok {
		result.CpuPrice, ok = (&big.Int{}).SetString(cpuPrice, 10)
	}
	if ramPrice != "" && ok {
		result.RamPrice, ok = (&big.Int{}).SetString(ramPrice, 10)
	}
	if storagePrice != "" && ok {
		result.StoragePrice, ok = (&big.Int{}).SetString(storagePrice, 10)
	}
	if bandwidthEPrice != "" && ok {
		result.BandwidthEgressPrice, ok = (&big.Int{}).SetString(bandwidthEPrice, 10)
	}
	if bandwidthInPrice != "" && ok {
		result.BandwidthIngressPrice, ok = (&big.Int{}).SetString(bandwidthInPrice, 10)
	}
	if tableId != "" && ok {
		result.Id, ok = (&big.Int{}).SetString(tableId, 10)
	}
	if !ok {
		return nil, fmt.Errorf("Badly-formatted integer flag") // -- it'd be nice to give more feedback as to which one
	}
	result.Cpumodel = cpuModel
	result.TeeType = teeType

	return result, nil
}
