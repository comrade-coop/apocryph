package contracts

import (
	"log"
	"math/big"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func DeployTokenContract(auth *bind.TransactOpts, c *ethclient.Client) (*common.Address, *abi.MockToken, error) {
	address, tx, instance, err := abi.DeployMockToken(auth, c)
	if err != nil {
		log.Printf("Could not deploy Token contract: %v", err)
		return nil, nil, err
	}
	log.Println("Token contract address:", address)

	log.Printf("Transaction hash: 0x%x\n\n", tx.Hash())
	return &address, instance, nil
}

func Balance(auth *bind.TransactOpts, c *ethclient.Client, instance *abi.MockToken, a common.Address) (*big.Int, error) {
	balance, err := instance.BalanceOf(&bind.CallOpts{Pending: false}, a)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
func GetTokenContractInstance(client *ethclient.Client, address common.Address) (*abi.MockToken, error) {
	return abi.NewMockToken(address, client)
}
