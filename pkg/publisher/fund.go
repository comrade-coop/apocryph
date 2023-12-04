package publisher

import (
	"fmt"
	"math/big"
	"os"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func FundPaymentChannel(ethClient *ethclient.Client, publisherAuth *bind.TransactOpts, deployment *pb.Deployment, funds *big.Int, unlockTimeInt *big.Int, mintFunds bool) error {
	providerEth := common.BytesToAddress(deployment.Provider.EthereumAddress)
	paymentContract := common.BytesToAddress(deployment.Payment.PaymentContractAddress)
	podID := common.BytesToHash(deployment.Payment.PodID)

	if funds.Cmp(common.Big0) != 0 {
		fmt.Fprintf(os.Stderr, "Funding payment channel...\n")

		// get a payment contract instance
		payment, err := abi.NewPayment(paymentContract, ethClient)
		if err != nil {
			return fmt.Errorf("Failed instanciating payment contract: %w", err)
		}
		tokenContract, err := payment.Token(&bind.CallOpts{})
		if err != nil {
			return fmt.Errorf("Failed resolving token contract address: %w", err)
		}

		if mintFunds {
			token, err := abi.NewMockToken(tokenContract, ethClient)
			if err != nil {
				return fmt.Errorf("Failed instanciating token contract: %w", err)
			}

			tx, err := token.Mint(publisherAuth, funds)
			if err != nil {
				return fmt.Errorf("Failed minting tokens: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Token mint successful! %v\n", tx.Hash())
		}

		token, err := abi.NewIERC20(tokenContract, ethClient)
		if err != nil {
			return fmt.Errorf("Failed instanciating token contract: %w", err)
		}

		tx, err := token.Approve(publisherAuth, paymentContract, funds)
		if err != nil {
			return fmt.Errorf("Failed approving token transfer: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Token approval successful! %v\n", tx.Hash())

		channel, err := payment.Channels(&bind.CallOpts{}, publisherAuth.From, providerEth, podID)
		if err != nil {
			return err
		}

		if channel.InvestedByPublisher.Cmp(common.Big0) > 0 {
			tx, err := payment.Deposit(publisherAuth, providerEth, podID, funds)
			if err != nil {
				return fmt.Errorf("Failed depositing payment contract funds: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Payment channel funded with %d extra funds! %v\n", funds, tx.Hash())
		} else {
			tx, err := payment.CreateChannel(publisherAuth, providerEth, podID, unlockTimeInt, funds)
			if err != nil {
				return fmt.Errorf("Failed creating payment contract: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Payment channel created with %d initial funds and %d seconds unlock time! %v\n", funds, unlockTimeInt, tx.Hash())
		}
	}
	return nil
}
