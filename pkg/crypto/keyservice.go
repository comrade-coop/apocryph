package crypto

// TODO
type KeyService interface {
	// generateKeys()
	// getKeys()
	// signWith(data []byte)    // sign with a key (add key argument)
	EncryptWith(data []byte, key string) // encrypt with a key
}

type MockKeyService struct{}

func (ks MockKeyService) EncryptWith(data []byte, key string) {}
