package contracts

import (
	"log"
	"math/big"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//	func CreatePaymentChannel(contract string) error {
//		client, err := ConnectToLocalNode()
//		if err != nil {
//			return err
//		}
//		return nil
//	}
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

func CreatePaymentChannel(auth *bind.TransactOpts, instance *abi.Payment, p common.Address, t common.Address, a *big.Int, d *big.Int, dur *big.Int, price *big.Int) (*types.Transaction, error) {
	tx, err := instance.CreateChannel(auth, p, t, a, d, dur, price)
	if err != nil {
		return nil, err
	}
	log.Printf("Payment Channel Created Successfully, tx hash: %v", tx.Hash())
	return tx, nil
}

func Withdraw(auth *bind.TransactOpts, instance *abi.Payment, client common.Address, id *big.Int, token common.Address) (*types.Transaction, error) {
	tx, err := instance.Withdraw(auth, client, id, token)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func LockFunds(auth *bind.TransactOpts, instance *abi.Payment, provider common.Address, id *big.Int, token common.Address, amount *big.Int) (*types.Transaction, error) {
	tx, err := instance.LockFunds(auth, provider, id, token, amount)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
func UploadMetrics(auth *bind.TransactOpts, instance *abi.Payment, client common.Address, id *big.Int, token common.Address, amount *big.Int) (*types.Transaction, error) {
	tx, err := instance.UploadMetrics(auth, client, id, token, amount)
	if err != nil {
		return nil, err
	}
	return tx, nil

}
func UpdateDeadline(auth *bind.TransactOpts, instance *abi.Payment, provider common.Address, id *big.Int, token common.Address, newDeadline *big.Int) (*types.Transaction, error) {
	tx, err := instance.UpdateDeadline(auth, provider, id, token, newDeadline)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func UpdatePrice(auth *bind.TransactOpts, instance *abi.Payment, client common.Address, id *big.Int, token common.Address, price *big.Int) (*types.Transaction, error) {
	tx, err := instance.UpdatePrice(auth, client, id, token, price)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func AcceptPrice(auth *bind.TransactOpts, instance *abi.Payment, provider common.Address, id *big.Int, token common.Address) (*types.Transaction, error) {
	tx, err := instance.AcceptNewPrice(auth, provider, id, token)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func GetContractInstance(client *ethclient.Client, address common.Address) (*abi.Payment, error) {
	return abi.NewPayment(address, client)
}
