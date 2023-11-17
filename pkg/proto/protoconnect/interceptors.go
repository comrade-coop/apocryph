package protoconnect

import (
	context "context"
	"crypto/ecdsa"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	tpPrefix      = "tpod"
	authorization = "Authorization"
)

type SignFunc func(data []byte) ([]byte, error)

type HasPubAddress interface{ GetPublisherAddress() []byte }

// func (p *PodLogRequest) GetCredentials() *Credentials { return p.Credentials }

/*
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
}*/

type authInterceptor struct {
	k8Client client.Client
}

func NewAuthInterceptor(c client.Client) connect.Interceptor {
	return authInterceptor{k8Client: c}
}

func (i authInterceptor) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		fmt.Printf("Authenticating gRPC call: %v \n", req.Spec())

		err := i.authenticate(req.Header())
		if err != nil {
			return nil, err
		}

		// verify that pod exists for the rest of the methods
		if req.Spec().Procedure != ProvisionPodServiceProvisionPodProcedure {
			p := &v1.Namespace{}
			namespace := req.Header().Get("namespace")
			err = i.k8Client.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: namespace}, p)
			if err != nil {
				return nil, err
			}
		}

		// Call the handler function
		return handler(ctx, req)
	}
}

func (i authInterceptor) WrapStreamingClient(handler connect.StreamingClientFunc) connect.StreamingClientFunc {
	return handler
}

func (i authInterceptor) WrapStreamingHandler(handler connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, stream connect.StreamingHandlerConn) error {
		fmt.Printf("Authenticating gRPC call: %v \n", stream.Spec())

		err := i.authenticate(stream.RequestHeader())
		if err != nil {
			return err
		}

		return handler(ctx, stream)
	}
}

func (a authInterceptor) authenticate(header http.Header) error {
	auth := header.Get(authorization)
	if len(auth) == 0 {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("Empty Authentication field"))
	}
	signature, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return err
	}
	tokenMd := header.Get("token")
	jsonData, err := base64.StdEncoding.DecodeString(tokenMd)
	if err != nil {
		return err
	}

	// verify if token has exired or not
	token := &Token{}
	err = json.Unmarshal(jsonData, token)
	if err != nil {
		return connect.NewError(connect.CodeDataLoss, fmt.Errorf("Failed Unmarshalling token"))
	}
	if time.Now().After(token.ExpirationTime) {
		return connect.NewError(connect.CodeDeadlineExceeded, fmt.Errorf("Token Expired"))
	}

	// Verify Signature
	valid, err := VerifyPayload(token.Publisher, jsonData, signature)
	if err != nil {
		return fmt.Errorf("Error verifying payload: %w", err)
	}

	// verify publisherAddress in namespace is same one signed in token
	ns := header.Get("namespace")
	nsExpected := namespaceFromTokenParts(token.Publisher, token.PodId)
	if ns != nsExpected {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("Inavlid namespace"))
	}
	// header.Set("namespace", nsExpected)?

	if !valid {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("Inavlid signature"))
	}

	return nil
}

type AuthInterceptorClient struct {
	Token            Token
	Signature        string
	Sign             SignFunc
	ExpirationOffset time.Duration
}

func NewAuthInterceptorClient(deployment *pb.Deployment, operation string, expirationOffset int64, sign SignFunc) *AuthInterceptorClient {
	return &AuthInterceptorClient{
		Token:            newToken(deployment, operation, expirationOffset),
		Sign:             sign,
		ExpirationOffset: time.Duration(expirationOffset) * time.Second,
	}
}

type Token struct {
	PodId          common.Hash
	Operation      string
	ExpirationTime time.Time
	Publisher      common.Address
}

func namespaceFromTokenParts(publisher common.Address, podId common.Hash) string {
	namespaceParts := []byte{}
	namespaceParts = append(namespaceParts, publisher[:]...)
	namespaceParts = append(namespaceParts, podId[:]...)
	namespacePartsHash := crypto.Keccak256(namespaceParts)
	return "tpod-" + strings.TrimRight(strings.ToLower(base32.StdEncoding.EncodeToString(namespacePartsHash)), "=")
}

func newToken(deployment *pb.Deployment, operation string, expirationTime int64) Token {

	return Token{
		PodId:          common.BytesToHash(deployment.Payment.PodID),
		Operation:      operation,
		ExpirationTime: time.Now().Add(time.Duration(expirationTime) * time.Second),
		Publisher:      common.BytesToAddress(deployment.Payment.PublisherAddress),
	}
}

func (a *AuthInterceptorClient) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		a.updateContext(req.Header())
		return handler(ctx, req)
	}
}

func (a *AuthInterceptorClient) WrapStreamingClient(handler connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		conn := handler(ctx, spec)
		a.updateContext(conn.RequestHeader())

		return conn
	}
}

func (a *AuthInterceptorClient) WrapStreamingHandler(handler connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return handler
}

func (a *AuthInterceptorClient) token() []byte {
	jsonData, err := json.Marshal(a.Token)
	if err != nil {
		return nil
	}
	return jsonData
}

func (a *AuthInterceptorClient) updateContext(header http.Header) error {

	// check if reached expiring date or first call to create a new signature
	if time.Now().After(a.Token.ExpirationTime) || a.Signature == "" {
		fmt.Println("Token Expired, Signing a New one ...")
		a.Token.ExpirationTime = time.Now().Add(a.ExpirationOffset)
		signature, err := a.Sign(a.token())
		if err != nil {
			return err
		}
		a.Signature = base64.StdEncoding.EncodeToString(signature)
	}

	token := base64.StdEncoding.EncodeToString(a.token())
	ns := namespaceFromTokenParts(a.Token.Publisher, a.Token.PodId)

	// Append custom metadata to the outgoing context
	header.Add(authorization, a.Signature)
	header.Add("token", token)
	header.Add("namespace", ns)
	return nil
}

func VerifyPayload(publisher common.Address, message []byte, signature []byte) (bool, error) {

	pubKeyECDSA, err := ExtractPubKey(message, signature)
	if err != nil {
		return false, err
	}

	// Ensure the signed address corresponds to the public key's address in the signature.
	// The signer should exclusively sign their own address;
	// thus, only the pods associated with their address used as IDs will be affected.
	address := crypto.PubkeyToAddress(*pubKeyECDSA)

	if publisher != address {
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
