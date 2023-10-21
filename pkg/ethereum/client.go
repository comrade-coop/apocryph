package ethereum

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitchellh/go-homedir"
)

func GetClient(url string) (*ethclient.Client, error) {
	if url == "" {
		url = "http://127.0.0.1:8545"
	}
	conn, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("could not connect to ethereum node: %w", err)
	}
	return conn, nil
}

const PrivateKeySize = 256 / 8
var DefaultKeystore = "~/.ethereum/keystore"

func GetAccount(accountString string, client *ethclient.Client) (*bind.TransactOpts, error) {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	if privKey, ok := strings.CutPrefix(accountString, "0x"); ok && len(privKey) == PrivateKeySize * 2 {
		key, err := crypto.HexToECDSA(privKey)
		if err != nil {
			return nil, err
		}

		auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
		if err != nil {
			return nil, err
		}
		return auth, nil
	}

	uri, accountAddress, ok := strings.Cut(accountString, "#")

	if !ok && strings.HasPrefix(accountString, "0x") {
		accountAddress = accountString
		var err error
		uri, err = homedir.Expand(DefaultKeystore)
		if err != nil {
			return nil, err
		}
	}

	if strings.Contains(uri, "://") {
		signer, err := external.NewExternalSigner(uri)
		if err != nil {
			return nil, err
		}

		for _, account := range signer.Accounts() {
			if account.Address.Hex() == accountAddress || account.URL.Path == accountAddress || accountAddress == "" {
				return bind.NewClefTransactor(signer, account), nil
			}
		}
		return nil, fmt.Errorf("Account %s not found in external signer %s", accountAddress, uri)
	}

	{
		ks := keystore.NewKeyStore(uri, keystore.StandardScryptN, keystore.StandardScryptP)
		accountAddress, accountPassphrase, _ := strings.Cut(accountAddress, ":")

		for _, account := range ks.Accounts() {
			if account.Address.Hex() == accountAddress || account.URL.Path == accountAddress || accountAddress == "" {
				if accountPassphrase == "" {
					accountPassphrase = utils.GetPassPhrase("", false)
				}
				ks.Unlock(account, accountPassphrase)
				return bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
			}
		}

		return nil, fmt.Errorf("Account %s not found in keystore %s", accountAddress, uri)
	}
}
