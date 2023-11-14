package proto

import (
	context "context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	reflect "reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	tpPrefix      = "tpod"
	CreatePod     = "/apocryph.proto.v0.provisionPod.ProvisionPodService/ProvisionPod"
	UpdatePod     = "/apocryph.proto.v0.provisionPod.ProvisionPodService/UpdatePod"
	DeletePod     = "/apocryph.proto.v0.provisionPod.ProvisionPodService/DeletePod"
	GetPodLogs    = "/apocryph.proto.v0.provisionPod.ProvisionPodService/GetPodLogs"
	authorization = "Authorization"
)

type HasPubAddress interface{ GetPublisherAddress() []byte }

// func (p *PodLogRequest) GetCredentials() *Credentials { return p.Credentials }

func NoCrashUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
	defer func() {
		if errRecover := recover(); errRecover != nil {
			fmt.Printf("Caught panic while processing GRPC call! %v %v\n", info, errRecover)
			err = fmt.Errorf("panic: %v", errRecover)
		}
	}()
	res, err = handler(ctx, req)
	return
}

func NoCrashStreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if errRecover := recover(); errRecover != nil {
			fmt.Printf("Caught panic while processing GRPC call! %v %v\n", info, errRecover)
			err = fmt.Errorf("panic: %v", errRecover)
		}
	}()
	err = handler(srv, ss)
	return
}

func AuthStreamServerInterceptor(c client.Client) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		fmt.Printf("Authenticating gRPC call: %v \n", info.FullMethod)

		// Extract metadata from the incoming context
		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			return status.Errorf(codes.Unauthenticated, "metadata not found")
		}

		err := authenticate(md, c)
		if err != nil {
			return err
		}

		// Call the handler function
		err = handler(srv, stream)
		if err != nil {
			log.Printf("Error handling stream: %v\n", err)
		}
		return err
	}
}
func AuthUnaryServerInterceptor(c client.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("Authenticating gRPC call: %v \n", info.FullMethod)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
		}

		err := authenticate(md, c)
		if err != nil {
			return nil, err
		}

		// verify that pod exists for the rest of the methods
		if info.FullMethod != CreatePod {
			p := &v1.Namespace{}
			namespace := md.Get("namespace")[0]
			err = c.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: namespace}, p)
			if err != nil {
				return nil, err
			}
		}

		// Call the handler function
		return handler(ctx, req)
	}
}

func authenticate(md metadata.MD, c client.Client) error {
	auth := md.Get(authorization)
	if len(auth) == 0 {
		return status.Errorf(codes.Unauthenticated, "Empty Authentication field")
	}
	signature, err := base64.StdEncoding.DecodeString(auth[0])
	if err != nil {
		return err
	}
	tokenMd := md.Get("token")[0]
	jsonData, err := base64.StdEncoding.DecodeString(tokenMd)
	if err != nil {
		return err
	}

	// verify if token has exired or not
	token := &Token{}
	err = json.Unmarshal(jsonData, token)
	if err != nil {
		return status.Errorf(codes.DataLoss, "Failed Unmarshalling token")
	}
	if time.Now().After(token.ExpirationTime) {
		return status.Errorf(codes.DeadlineExceeded, "Token Expired")
	}

	// Verify Signature
	valid, err := VerifyPayload(token.Publisher, jsonData, signature)
	if err != nil {
		log.Printf("Error verifying payload: %v\n", err)
		return err
	}

	// verify publisherAddress in namespace is same one signed in token
	namespace := md.Get("namespace")[0]
	publisherAddress := strings.Split(namespace, "-")[1]
	tokenAddress := strings.ToLower(common.BytesToAddress(token.Publisher).String())
	if publisherAddress != tokenAddress {
		return status.Errorf(codes.Unauthenticated, "Inavlid namespace")
	}
	if !valid {
		return status.Errorf(codes.Unauthenticated, "Invalid signature")
	}

	return nil
}

type AuthInterceptorClient struct {
	Token            Token
	Account          accounts.Account
	Signature        string
	ExpirationOffset time.Duration
	Keystore         *keystore.KeyStore
}

type Token struct {
	PodId          string
	Operation      string
	ExpirationTime time.Time
	Publisher      []byte
}

func (a *AuthInterceptorClient) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		a.updateContext(&ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (a *AuthInterceptorClient) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		a.updateContext(&ctx)
		// Call the streamer function to obtain a client stream
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}

		return clientStream, nil
	}
}

func (a *AuthInterceptorClient) token() []byte {
	jsonData, err := json.Marshal(a.Token)
	if err != nil {
		return nil
	}
	return jsonData
}

func (a *AuthInterceptorClient) updateContext(ctx *context.Context) error {

	// check if reached expiring date or first call to create a new signature
	if time.Now().After(a.Token.ExpirationTime) || a.Signature == "" {
		fmt.Println("Token Expired, Signing a New one ...")
		a.Token.ExpirationTime = time.Now().Add(a.ExpirationOffset)
		signature, err := SignPayload(a.token(), a.Account, "123", a.Keystore)
		if err != nil {
			return err
		}
		a.Signature = base64.StdEncoding.EncodeToString(signature)
	}
	token := base64.StdEncoding.EncodeToString(a.token())
	namespace := "tpod-" + strings.ToLower(a.Account.Address.String()+"-"+a.Token.PodId)

	// Append custom metadata to the outgoing context
	*ctx = metadata.AppendToOutgoingContext(*ctx, authorization, a.Signature)
	*ctx = metadata.AppendToOutgoingContext(*ctx, "token", token)
	*ctx = metadata.AppendToOutgoingContext(*ctx, "namespace", namespace)
	return nil
}

func SignPayload(data []byte, acc accounts.Account, psw string, ks *keystore.KeyStore) ([]byte, error) {
	hash := crypto.Keccak256(data)
	signature, err := ks.SignHashWithPassphrase(acc, psw, hash)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func VerifyPayload(publisher, message []byte, signature []byte) (bool, error) {

	pubKeyECDSA, err := ExtractPubKey(message, signature)
	if err != nil {
		return false, err
	}

	// Ensure the signed address corresponds to the public key's address in the signature.
	// The signer should exclusively sign their own address;
	// thus, only the pods associated with their address used as IDs will be affected.
	address := crypto.PubkeyToAddress(*pubKeyECDSA).Bytes()

	if !reflect.DeepEqual(publisher, address) {
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
