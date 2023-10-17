package crypto

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

func CreateKeystore(path string) *keystore.KeyStore {
	return keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
}

func GetAccountManger(ks *keystore.KeyStore) *accounts.Manager {
	return accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, ks)
}

func CreateAccount(psw string, ks *keystore.KeyStore) (*accounts.Account, error) {
	newAcc, err := ks.NewAccount(psw)
	if err != nil {
		return nil, err
	}
	return &newAcc, nil
}

func ExportAccount(acc accounts.Account, psw, exportpsw string, ks *keystore.KeyStore) ([]byte, error) {
	keyJson, err := ks.Export(acc, psw, exportpsw)
	if err != nil {
		return nil, err
	}
	return keyJson, nil
}

func ImportAccount(psw string, importpsw string, keyJson []byte, ks *keystore.KeyStore) (*accounts.Account, error) {
	impAcc, err := ks.Import(keyJson, psw, importpsw)
	if err != nil {
		return nil, err
	}
	return &impAcc, nil

}

func DerivePrvKey(key string) (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func AccountFromECDSA(key *ecdsa.PrivateKey, psw string, ks *keystore.KeyStore) (accounts.Account, error) {
	return ks.ImportECDSA(key, psw)
}
