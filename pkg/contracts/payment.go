package contracts

import (
	"log"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func DeployPaymentContract(auth *bind.TransactOpts, c *ethclient.Client) (*common.Address, *abi.Payment, error) {

	address, tx, instance, err := abi.DeployPayment(auth, c)
	if err != nil {
		log.Printf("Could not deploy contract: %v", err)
		return nil, nil, err
	}
	log.Println("Payment contract address:", address)

	log.Printf("Transaction hash: 0x%x\n\n", tx.Hash())
	return &address, instance, nil
}

func GetContractInstance(client *ethclient.Client, address common.Address) (*abi.Payment, error) {
	return abi.NewPayment(address, client)
}
