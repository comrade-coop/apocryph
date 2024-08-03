// SPDX-License-Identifier: GPL-3.0

package publisher

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/comrade-coop/apocryph/pkg/abi"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitchellh/go-homedir"
)

var DefaultDeploymentPath = "~/.apocryph/deployment"
var DefaultPodFile = "manifest.yaml"

const PrivateKeySize = 256 / 8

func GenerateDeploymentFilename(podFile string, deploymentFormat string) (deploymentFile string, relPodFile string, err error) {
	if deploymentFormat == "" {
		deploymentFormat = "yaml"
	}

	absPodFile, err := filepath.Abs(podFile)
	if err != nil {
		return
	}

	podFileHash := sha256.Sum256([]byte(absPodFile))
	podFileHashHex := hex.EncodeToString(podFileHash[:])
	deploymentFilename := fmt.Sprintf("%s.%s", podFileHashHex, deploymentFormat)
	deploymentRoot, err := homedir.Expand(DefaultDeploymentPath)
	if err != nil {
		return
	}

	err = os.MkdirAll(deploymentRoot, 0755)
	if err != nil {
		return
	}

	deploymentFile = filepath.Join(deploymentRoot, deploymentFilename)

	absDeploymentFile, err := filepath.Abs(deploymentFile)
	if err != nil {
		return
	}
	relPodFile, err = filepath.Rel(filepath.Dir(absDeploymentFile), absPodFile)

	return
}

func ReadPodAndDeployment(args []string, manifestFormat string, deploymentFormat string) (podFile string, deploymentFile string, pod *pb.Pod, deployment *pb.Deployment, err error) {
	deployment = &pb.Deployment{}
	readDeployment := false

	switch len(args) {
	case 0:
		podFile = DefaultPodFile
	case 1:
		podFile = args[0]
	case 2:
		podFile = args[0]
		deploymentFile = args[1]
	default:
		err = fmt.Errorf("Wrong number of arguments passed to ReadPodAndDeployment")
		return
	}

	// Get the name of the deployment file if it was not passed in the args
	if deploymentFile == "" {
		deploymentFile, deployment.PodManifestFile, err = GenerateDeploymentFilename(podFile, deploymentFormat)
		if err != nil {
			err = fmt.Errorf("Failed resolving deployment file path: %w", err)
			return
		}
	}

	if !readDeployment {
		err = pb.UnmarshalFile(deploymentFile, deploymentFormat, deployment)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			err = fmt.Errorf("Failed reading deployment file %s: %w", deploymentFile, err)
			return
		}
	}

	pod = &pb.Pod{}
	err = pb.UnmarshalFile(podFile, manifestFormat, pod)
	if err != nil {
		err = fmt.Errorf("Failed reading manifest file %s: %w", podFile, err)
	}

	return
}

func SaveDeployment(deploymentFile string, deploymentFormat string, deployment *pb.Deployment) error {
	err := pb.MarshalFile(deploymentFile, deploymentFormat, deployment)
	if err != nil {
		return fmt.Errorf("Failed saving deployment file: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Stored deployment data in %s\n", deploymentFile)
	// if deployment == nil {
	// 	err := os.Remove(deploymentFile)
	// 	if err != nil && !errors.Is(err, os.ErrNotExist) {
	// 		return fmt.Errorf("Failed to remove deployment file: %w", err)
	// 	}
	// 	fmt.Fprintf(os.Stderr, "Removed deployment data from %s\n", deploymentFile)
	// }
	return nil
}

func AuthorizeAndFundApplication(ctx context.Context, response *pb.ProvisionPodResponse, deployment *pb.Deployment, ethClient *ethclient.Client, publisherAuth *bind.TransactOpts, publisherKey string, amount int64) error {
	if deployment.KeyPair != nil {
		// authorize the public address to control the payment channel
		// get a payment contract instance
		pubAddress := common.HexToAddress(deployment.KeyPair.PubAddress)
		payment, err := abi.NewPayment(common.Address(deployment.Payment.PaymentContractAddress), ethClient)
		if err != nil {
			return fmt.Errorf("Failed instantiating payment contract: %w", err)
		}
		_, err = payment.Authorize(publisherAuth, pubAddress, common.Address(deployment.Provider.EthereumAddress), [32]byte(deployment.Payment.PodID))
		if err != nil {
			return fmt.Errorf("Failed Authorizing Address: %w", err)
		}

		// NOTE: the deployed application must be funded with the base
		// currency (eth) in order for it to be able to make transactions
		// (like creating subchannels)
		amount := big.NewInt(amount)
		err = FundAddress(ctx, publisherKey, publisherAuth.From, pubAddress, ethClient, amount)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Authorized And funded address %v\n", pubAddress)
	}
	return nil
}

func FundAddress(ctx context.Context, key string, from, to common.Address, ethClient *ethclient.Client, amount *big.Int) error {

	// Estimate gas limit
	gasLimit, err := ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  from,
		To:    &to,
		Value: amount,
	})
	if err != nil {
		return fmt.Errorf("Failed to estimate gas: %w", err)
	}

	// Get gas price
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get gas price: %w", err)
	}

	nonce, err := ethClient.PendingNonceAt(ctx, from)
	if err != nil {
		return fmt.Errorf("Failed to get gas price: %w", err)
	}

	// Create a transaction
	tx := types.NewTransaction(
		nonce,
		to,
		amount,
		gasLimit,
		gasPrice,
		nil,
	)
	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get chain ID: %w", err)
	}

	var privateKey *ecdsa.PrivateKey
	if privKey, ok := strings.CutPrefix(key, "0x"); ok && len(privKey) == PrivateKeySize*2 {
		privateKey, err = crypto.HexToECDSA(privKey)
		if err != nil {
			return err
		}
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		return fmt.Errorf("Failed to sign transaction: %w", err)
	}

	// Send the transaction
	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("Failed to send transaction: %w", err)
	}
	return nil
}
