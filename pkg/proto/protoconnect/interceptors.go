// SPDX-License-Identifier: GPL-3.0

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
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	tpPrefix            = "tpod"
	authorizationHeader = "Authorization"
	NamespaceHeader     = "X-Namespace"
)

type Token struct {
	PodId          common.Hash
	Operation      string
	ExpirationTime time.Time
	Publisher      common.Address
}

type SignFunc func(data []byte) ([]byte, error)

type HasPaymentChannel interface{ GetPayment() *pb.PaymentChannel }

// func (p *PodLogRequest) GetCredentials() *Credentials { return p.Credentials }

type authInterceptor struct {
	k8Client client.Client
}

func NewAuthInterceptor(c client.Client) connect.Interceptor {
	return authInterceptor{k8Client: c}
}

func (i authInterceptor) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		fmt.Printf("Authenticating gRPC call: %v \n", req.Spec())
		if req.Spec().Procedure == ProvisionPodServiceGetPodInfosProcedure {
			return handler(ctx, req)
		}
		expectedPublisher, err := i.authenticate(req.Header())
		if err != nil {
			return nil, err
		}

		if hasPayment, ok := req.Any().(HasPaymentChannel); ok {
			publisherAddress := common.BytesToAddress(hasPayment.GetPayment().GetPublisherAddress())
			if expectedPublisher != publisherAddress {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("Payment channel created for unauthorized address: (expected: %v, found: %v)", expectedPublisher, publisherAddress))
			}
		}

		// verify that pod exists for the rest of the methods
		if req.Spec().Procedure != ProvisionPodServiceProvisionPodProcedure {
			p := &v1.Namespace{}
			namespace := GetNamespace(req)
			err = i.k8Client.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: namespace}, p)
			if err != nil {
				return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("Namespace not found: %w", err))
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

		_, err := i.authenticate(stream.RequestHeader())
		if err != nil {
			return err
		}

		return handler(ctx, stream)
	}
}

func (a authInterceptor) authenticate(header http.Header) (common.Address, error) {
	auth := header.Get(authorizationHeader)
	if len(auth) == 0 {
		return common.Address{}, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("Empty Authentication field"))
	}
	tokenString, ok := strings.CutPrefix(auth, "Bearer ")
	if !ok {
		return common.Address{}, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("Expected Bearer Authentication"))
	}
	tokenParts := strings.Split(tokenString, ".")
	if len(tokenParts) != 2 {
		return common.Address{}, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("Invalid token (wrong number of parts)"))
	}

	tokenData, err := base64.StdEncoding.DecodeString(tokenParts[0])
	if err != nil {
		return common.Address{}, err
	}
	tokenSignature, err := base64.StdEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return common.Address{}, err
	}

	// verify if token has exired or not
	token := &Token{}
	err = json.Unmarshal(tokenData, token)
	if err != nil {
		return common.Address{}, connect.NewError(connect.CodeDataLoss, fmt.Errorf("Failed Unmarshalling token"))
	}
	if time.Now().UTC().After(token.ExpirationTime) {
		return common.Address{}, connect.NewError(connect.CodeDeadlineExceeded, fmt.Errorf("Token Expired"))
	}

	// Verify Signature
	valid, err := VerifyPayload(token.Publisher, tokenData, tokenSignature)
	if err != nil {
		return common.Address{}, fmt.Errorf("Error verifying payload: %w", err)
	}
	if !valid {
		return common.Address{}, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("Invalid signature"))
	}

	// verify publisherAddress in namespace is same one signed in token
	ns := header.Get(NamespaceHeader)
	nsExpected := NamespaceFromTokenParts(token.Publisher, token.PodId)
	if ns == "" {
		header.Set(NamespaceHeader, nsExpected)
	} else if ns != nsExpected {
		return common.Address{}, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("Invalid namespace"))
	}

	return token.Publisher, nil
}

func NamespaceFromTokenParts(publisher common.Address, podId common.Hash) string {
	namespaceParts := []byte{}
	namespaceParts = append(namespaceParts, publisher[:]...)
	namespaceParts = append(namespaceParts, podId[:]...)
	namespacePartsHash := crypto.Keccak256(namespaceParts)
	return "tpod-" + strings.TrimRight(strings.ToLower(base32.StdEncoding.EncodeToString(namespacePartsHash)), "=")
}

func GetNamespace(req connect.AnyRequest) string {
	return req.Header().Get(NamespaceHeader)
}

type AuthInterceptorClient struct {
	tokens           map[string]serializedToken
	podId            common.Hash
	publisher        common.Address
	sign             SignFunc
	expirationOffset time.Duration
}

type serializedToken struct {
	expirationTime time.Time
	bearer         string
}

func NewAuthInterceptorClient(deployment *pb.Deployment, expirationOffset int64, sign SignFunc) *AuthInterceptorClient {
	return &AuthInterceptorClient{
		tokens:           make(map[string]serializedToken),
		sign:             sign,
		expirationOffset: time.Duration(expirationOffset) * time.Second,
		podId:            common.BytesToHash(deployment.Payment.PodID),
		publisher:        common.BytesToAddress(deployment.Payment.PublisherAddress),
	}
}

func (a *AuthInterceptorClient) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if req.Spec().Procedure == ProvisionPodServiceGetPodInfosProcedure {
			return handler(ctx, req)
		}
		a.authorize(req.Spec().Procedure, req.Header())
		return handler(ctx, req)
	}
}

func (a *AuthInterceptorClient) WrapStreamingClient(handler connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		conn := handler(ctx, spec)
		a.authorize(spec.Procedure, conn.RequestHeader())

		return conn
	}
}

func (a *AuthInterceptorClient) WrapStreamingHandler(handler connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return handler
}

func (a *AuthInterceptorClient) getOrCreateToken(operation string) (serializedToken, error) {
	token := a.tokens[operation]
	// check if reached expiring date or first call to create a new signature
	if time.Now().After(token.expirationTime) || token.bearer == "" {
		fmt.Println("Token Expired, Signing a New one ...")
		tokenData := Token{
			PodId:          a.podId,
			Operation:      operation,
			ExpirationTime: time.Now().UTC().Add(a.expirationOffset),
			Publisher:      a.publisher,
		}
		tokenDataBytes, err := json.Marshal(tokenData)
		if err != nil {
			return serializedToken{}, err
		}
		signature, err := a.sign(tokenDataBytes)
		tokenDataEncoded := base64.StdEncoding.EncodeToString(tokenDataBytes)
		signatureEncoded := base64.StdEncoding.EncodeToString(signature)
		token = serializedToken{
			expirationTime: tokenData.ExpirationTime,
			bearer:         fmt.Sprintf("%s.%s", tokenDataEncoded, signatureEncoded),
		}
		a.tokens[operation] = token
	}
	return token, nil
}

func (a *AuthInterceptorClient) authorize(operation string, header http.Header) error {
	token, err := a.getOrCreateToken(operation)
	if err != nil {
		return err
	}

	// Append custom metadata to the outgoing context
	header.Add(authorizationHeader, fmt.Sprintf("Bearer %s", token.bearer))
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
