// SPDX-License-Identifier: GPL-3.0

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: autoscaler.proto

package protoconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	proto "github.com/comrade-coop/apocryph/pkg/proto"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// AutoscalerServiceName is the fully-qualified name of the AutoscalerService service.
	AutoscalerServiceName = "apocryph.proto.v0.autoscaler.AutoscalerService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AutoscalerServiceConnectClusterProcedure is the fully-qualified name of the AutoscalerService's
	// ConnectCluster RPC.
	AutoscalerServiceConnectClusterProcedure = "/apocryph.proto.v0.autoscaler.AutoscalerService/ConnectCluster"
	// AutoscalerServiceTriggerNodeProcedure is the fully-qualified name of the AutoscalerService's
	// TriggerNode RPC.
	AutoscalerServiceTriggerNodeProcedure = "/apocryph.proto.v0.autoscaler.AutoscalerService/TriggerNode"
)

// AutoscalerServiceClient is a client for the apocryph.proto.v0.autoscaler.AutoscalerService
// service.
type AutoscalerServiceClient interface {
	ConnectCluster(context.Context, *connect.Request[proto.ConnectClusterRequest]) (*connect.Response[proto.ConnectClusterResponse], error)
	TriggerNode(context.Context, *connect.Request[proto.ConnectClusterRequest]) (*connect.Response[proto.TriggerNodeResponse], error)
}

// NewAutoscalerServiceClient constructs a client for the
// apocryph.proto.v0.autoscaler.AutoscalerService service. By default, it uses the Connect protocol
// with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To
// use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb()
// options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAutoscalerServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AutoscalerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &autoscalerServiceClient{
		connectCluster: connect.NewClient[proto.ConnectClusterRequest, proto.ConnectClusterResponse](
			httpClient,
			baseURL+AutoscalerServiceConnectClusterProcedure,
			opts...,
		),
		triggerNode: connect.NewClient[proto.ConnectClusterRequest, proto.TriggerNodeResponse](
			httpClient,
			baseURL+AutoscalerServiceTriggerNodeProcedure,
			opts...,
		),
	}
}

// autoscalerServiceClient implements AutoscalerServiceClient.
type autoscalerServiceClient struct {
	connectCluster *connect.Client[proto.ConnectClusterRequest, proto.ConnectClusterResponse]
	triggerNode    *connect.Client[proto.ConnectClusterRequest, proto.TriggerNodeResponse]
}

// ConnectCluster calls apocryph.proto.v0.autoscaler.AutoscalerService.ConnectCluster.
func (c *autoscalerServiceClient) ConnectCluster(ctx context.Context, req *connect.Request[proto.ConnectClusterRequest]) (*connect.Response[proto.ConnectClusterResponse], error) {
	return c.connectCluster.CallUnary(ctx, req)
}

// TriggerNode calls apocryph.proto.v0.autoscaler.AutoscalerService.TriggerNode.
func (c *autoscalerServiceClient) TriggerNode(ctx context.Context, req *connect.Request[proto.ConnectClusterRequest]) (*connect.Response[proto.TriggerNodeResponse], error) {
	return c.triggerNode.CallUnary(ctx, req)
}

// AutoscalerServiceHandler is an implementation of the
// apocryph.proto.v0.autoscaler.AutoscalerService service.
type AutoscalerServiceHandler interface {
	ConnectCluster(context.Context, *connect.Request[proto.ConnectClusterRequest]) (*connect.Response[proto.ConnectClusterResponse], error)
	TriggerNode(context.Context, *connect.Request[proto.ConnectClusterRequest]) (*connect.Response[proto.TriggerNodeResponse], error)
}

// NewAutoscalerServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAutoscalerServiceHandler(svc AutoscalerServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	autoscalerServiceConnectClusterHandler := connect.NewUnaryHandler(
		AutoscalerServiceConnectClusterProcedure,
		svc.ConnectCluster,
		opts...,
	)
	autoscalerServiceTriggerNodeHandler := connect.NewUnaryHandler(
		AutoscalerServiceTriggerNodeProcedure,
		svc.TriggerNode,
		opts...,
	)
	return "/apocryph.proto.v0.autoscaler.AutoscalerService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AutoscalerServiceConnectClusterProcedure:
			autoscalerServiceConnectClusterHandler.ServeHTTP(w, r)
		case AutoscalerServiceTriggerNodeProcedure:
			autoscalerServiceTriggerNodeHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAutoscalerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAutoscalerServiceHandler struct{}

func (UnimplementedAutoscalerServiceHandler) ConnectCluster(context.Context, *connect.Request[proto.ConnectClusterRequest]) (*connect.Response[proto.ConnectClusterResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apocryph.proto.v0.autoscaler.AutoscalerService.ConnectCluster is not implemented"))
}

func (UnimplementedAutoscalerServiceHandler) TriggerNode(context.Context, *connect.Request[proto.ConnectClusterRequest]) (*connect.Response[proto.TriggerNodeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apocryph.proto.v0.autoscaler.AutoscalerService.TriggerNode is not implemented"))
}
