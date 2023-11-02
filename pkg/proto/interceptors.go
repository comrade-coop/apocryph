package proto

import (
	context "context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type HasCredentials interface{ GetCredentials() *Credentials }

// func (p *PodLogRequest) GetCredentials() *Credentials { return p.Credentials }

func AuthUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod != "/apocryph.proto.v0.provisionPod.ProvisionPodService/ProvisionPod" {
			fmt.Printf("Authenticating gRPC call: %v \n", info.FullMethod)
			// Extract the credentials from the request
			credentials := req.(HasCredentials).GetCredentials()
			// because the auto-genrated protobuf code only checks the whole message not the internal field
			if credentials == nil {
				return nil, status.Errorf(codes.Unauthenticated, "Empty Credentials")
			}
			// Perform authentication
			valid, err := VerifyPayload(credentials.PublisherAddress, credentials.Signature)
			if err != nil {
				log.Printf("Error verifying payload: %v\n", err)
				return nil, err
			}

			if !valid {
				return nil, status.Errorf(codes.Unauthenticated, "Invalid signature")
			}
		}
		// Call the handler function
		return handler(ctx, req)
	}
}
func SignPayload(data []byte, acc accounts.Account, psw string, ks *keystore.KeyStore) ([]byte, error) {
	hash := crypto.Keccak256(data)
	signature, err := ks.SignHashWithPassphrase(acc, psw, hash)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func VerifyPayload(message []byte, signature []byte) (bool, error) {

	pubKeyECDSA, err := ExtractPubKey(message, signature)
	if err != nil {
		return false, err
	}

	// Ensure the signed address corresponds to the public key's address in the signature.
	// The signer should exclusively sign their own address;
	// thus, only the pods associated with their address used as IDs will be affected.
	address := []byte(crypto.PubkeyToAddress(*pubKeyECDSA).Hex())

	if !reflect.DeepEqual(message, address) {
		return false, nil
	}

	pubKey := crypto.FromECDSAPub(pubKeyECDSA)
	valid := crypto.VerifySignature(pubKey, crypto.Keccak256(message), signature[:len(signature)-1])

	if valid {
		return true, nil
	}
	return false, nil
}

func ExtractPubKey(message []byte, signature []byte) (*ecdsa.PublicKey, error) {
	pubKeyECDSA, err := crypto.SigToPub(crypto.Keccak256(message), signature)
	if err != nil {
		return nil, err
	}
	return pubKeyECDSA, nil
}
