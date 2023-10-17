package contracts

import (
	"context"
	"log"
	"strings"

	"github.com/comrade-coop/trusted-pods/pkg/crypto"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ConnectToLocalNode() (*ethclient.Client, error) {
	return Connect("http://127.0.0.1:8545")
}

func Connect(url string) (*ethclient.Client, error) {
	conn, err := ethclient.Dial(url)
	if err != nil {
		log.Printf("could not connect to ethereum node: %v", err)
		return nil, err
	}
	return conn, nil
}

func GetAccounts(m *accounts.Manager) []common.Address {
	return m.Accounts()
}

func CreateTransactor(key string, psw string, c *ethclient.Client) (*bind.TransactOpts, error) {

	context := context.Background()
	chainID, err := c.ChainID(context)
	if err != nil {
		log.Printf("could not retreive chain id from client: %v", err)
		return nil, err
	}
	auth, err := bind.NewTransactorWithChainID(strings.NewReader(key), psw, chainID)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func DeriveAccountConfigs(privatekey string, psw string, exportpsw string, client *ethclient.Client, ks *keystore.KeyStore) (*accounts.Account, *bind.TransactOpts, error) {

	key, err := crypto.DerivePrvKey(privatekey[2:])
	if err != nil {
		log.Printf("could not derive private key: %v", err)
		return nil, nil, err
	}

	// derive the account from the key
	acc, err := crypto.AccountFromECDSA(key, psw, ks)
	if err != nil {
		log.Printf("error creating account: %v \n", err)
	}

	//export the account to get the json string
	keyjson, err := crypto.ExportAccount(acc, psw, exportpsw, ks)
	if err != nil {
		log.Printf("error exporting account: %v", err)
	}
	// get a transactor
	providerAuth, err := CreateTransactor(string(keyjson), psw, client)
	if err != nil {
		log.Println("Failed to create authorized transactor: %v", err)
		return nil, nil, err
	}
	return &acc, providerAuth, nil

}

//	type wallet struct {
//		privateKey *ecdsa.PrivateKey
//		pubaddress common.Address
//	}
//
//	func GenerateWallet(prvkey ...string) (*wallet, error) {
//		var privateKey *ecdsa.PrivateKey
//		var err error
//		if len(prvkey) > 0 {
//			privateKey, err = DerivePrvKey(prvkey[0])
//			if err != nil {
//				log.Println("could not derive private key:", err)
//				return nil, err
//			}
//		} else {
//			privateKey, err = crypto.GenerateKey()
//			if err != nil {
//				log.Println("could not generate key:", err)
//				return nil, err
//			}
//		}
//
//		// privateKeyBytes := crypto.FromECDSA(privateKey)
//		// prvkey := hexutil.Encode(privateKeyBytes)[2:]
//
//		publicKey := privateKey.Public()
//		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
//		if !ok {
//			log.Fatal("error casting public key to ECDSA")
//		}
//
//		publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
//		fmt.Println(hexutil.Encode(publicKeyBytes)[4:])
//
//		address := crypto.PubkeyToAddress(*publicKeyECDSA)
//
//		// hash := sha3.NewLegacyKeccak256()
//		// hash.Write(publicKeyBytes[1:])
//		// pubaddr := hexutil.Encode(hash.Sum(nil)[12:])
//		return &wallet{
//			privateKey: privateKey,
//			pubaddress: address,
//		}, nil
//	}
//
// func (w *wallet) GetNonce(c *ethclient.Client) (uint64, error) {
//
// 	context := context.Background()
//
// 	nonce, err := c.PendingNonceAt(context, w.pubaddress)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return nonce, nil
//
// }
//
// func (w *wallet) CreateSigner(client *ethclient.Client) (*bind.TransactOpts, error) {
// 	context := context.Background()
// 	chainID, err := client.ChainID(context)
// 	if err != nil {
// 		log.Printf("could not retreive chain id from client: %v", err)
// 		return nil, err
// 	}
// 	auth, err := bind.NewKeyedTransactorWithChainID(w.privateKey, chainID)
// 	if err != nil {
// 		log.Printf("could not create signer: %v", err)
// 		return nil, err
// 	}
// 	return auth, nil
// }
