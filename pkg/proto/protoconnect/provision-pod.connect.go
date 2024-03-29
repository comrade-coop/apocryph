// SPDX-License-Identifier: GPL-3.0

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: provision-pod.proto

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
	// ProvisionPodServiceName is the fully-qualified name of the ProvisionPodService service.
	ProvisionPodServiceName = "apocryph.proto.v0.provisionPod.ProvisionPodService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ProvisionPodServiceProvisionPodProcedure is the fully-qualified name of the ProvisionPodService's
	// ProvisionPod RPC.
	ProvisionPodServiceProvisionPodProcedure = "/apocryph.proto.v0.provisionPod.ProvisionPodService/ProvisionPod"
	// ProvisionPodServiceUpdatePodProcedure is the fully-qualified name of the ProvisionPodService's
	// UpdatePod RPC.
	ProvisionPodServiceUpdatePodProcedure = "/apocryph.proto.v0.provisionPod.ProvisionPodService/UpdatePod"
	// ProvisionPodServiceDeletePodProcedure is the fully-qualified name of the ProvisionPodService's
	// DeletePod RPC.
	ProvisionPodServiceDeletePodProcedure = "/apocryph.proto.v0.provisionPod.ProvisionPodService/DeletePod"
	// ProvisionPodServiceGetPodLogsProcedure is the fully-qualified name of the ProvisionPodService's
	// GetPodLogs RPC.
	ProvisionPodServiceGetPodLogsProcedure = "/apocryph.proto.v0.provisionPod.ProvisionPodService/GetPodLogs"
)

// ProvisionPodServiceClient is a client for the apocryph.proto.v0.provisionPod.ProvisionPodService
// service.
type ProvisionPodServiceClient interface {
	ProvisionPod(context.Context, *connect.Request[proto.ProvisionPodRequest]) (*connect.Response[proto.ProvisionPodResponse], error)
	UpdatePod(context.Context, *connect.Request[proto.UpdatePodRequest]) (*connect.Response[proto.ProvisionPodResponse], error)
	DeletePod(context.Context, *connect.Request[proto.DeletePodRequest]) (*connect.Response[proto.DeletePodResponse], error)
	GetPodLogs(context.Context, *connect.Request[proto.PodLogRequest]) (*connect.ServerStreamForClient[proto.PodLogResponse], error)
}

// NewProvisionPodServiceClient constructs a client for the
// apocryph.proto.v0.provisionPod.ProvisionPodService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewProvisionPodServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ProvisionPodServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &provisionPodServiceClient{
		provisionPod: connect.NewClient[proto.ProvisionPodRequest, proto.ProvisionPodResponse](
			httpClient,
			baseURL+ProvisionPodServiceProvisionPodProcedure,
			opts...,
		),
		updatePod: connect.NewClient[proto.UpdatePodRequest, proto.ProvisionPodResponse](
			httpClient,
			baseURL+ProvisionPodServiceUpdatePodProcedure,
			opts...,
		),
		deletePod: connect.NewClient[proto.DeletePodRequest, proto.DeletePodResponse](
			httpClient,
			baseURL+ProvisionPodServiceDeletePodProcedure,
			opts...,
		),
		getPodLogs: connect.NewClient[proto.PodLogRequest, proto.PodLogResponse](
			httpClient,
			baseURL+ProvisionPodServiceGetPodLogsProcedure,
			opts...,
		),
	}
}

// provisionPodServiceClient implements ProvisionPodServiceClient.
type provisionPodServiceClient struct {
	provisionPod *connect.Client[proto.ProvisionPodRequest, proto.ProvisionPodResponse]
	updatePod    *connect.Client[proto.UpdatePodRequest, proto.ProvisionPodResponse]
	deletePod    *connect.Client[proto.DeletePodRequest, proto.DeletePodResponse]
	getPodLogs   *connect.Client[proto.PodLogRequest, proto.PodLogResponse]
}

// ProvisionPod calls apocryph.proto.v0.provisionPod.ProvisionPodService.ProvisionPod.
func (c *provisionPodServiceClient) ProvisionPod(ctx context.Context, req *connect.Request[proto.ProvisionPodRequest]) (*connect.Response[proto.ProvisionPodResponse], error) {
	return c.provisionPod.CallUnary(ctx, req)
}

// UpdatePod calls apocryph.proto.v0.provisionPod.ProvisionPodService.UpdatePod.
func (c *provisionPodServiceClient) UpdatePod(ctx context.Context, req *connect.Request[proto.UpdatePodRequest]) (*connect.Response[proto.ProvisionPodResponse], error) {
	return c.updatePod.CallUnary(ctx, req)
}

// DeletePod calls apocryph.proto.v0.provisionPod.ProvisionPodService.DeletePod.
func (c *provisionPodServiceClient) DeletePod(ctx context.Context, req *connect.Request[proto.DeletePodRequest]) (*connect.Response[proto.DeletePodResponse], error) {
	return c.deletePod.CallUnary(ctx, req)
}

// GetPodLogs calls apocryph.proto.v0.provisionPod.ProvisionPodService.GetPodLogs.
func (c *provisionPodServiceClient) GetPodLogs(ctx context.Context, req *connect.Request[proto.PodLogRequest]) (*connect.ServerStreamForClient[proto.PodLogResponse], error) {
	return c.getPodLogs.CallServerStream(ctx, req)
}

// ProvisionPodServiceHandler is an implementation of the
// apocryph.proto.v0.provisionPod.ProvisionPodService service.
type ProvisionPodServiceHandler interface {
	ProvisionPod(context.Context, *connect.Request[proto.ProvisionPodRequest]) (*connect.Response[proto.ProvisionPodResponse], error)
	UpdatePod(context.Context, *connect.Request[proto.UpdatePodRequest]) (*connect.Response[proto.ProvisionPodResponse], error)
	DeletePod(context.Context, *connect.Request[proto.DeletePodRequest]) (*connect.Response[proto.DeletePodResponse], error)
	GetPodLogs(context.Context, *connect.Request[proto.PodLogRequest], *connect.ServerStream[proto.PodLogResponse]) error
}

// NewProvisionPodServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewProvisionPodServiceHandler(svc ProvisionPodServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	provisionPodServiceProvisionPodHandler := connect.NewUnaryHandler(
		ProvisionPodServiceProvisionPodProcedure,
		svc.ProvisionPod,
		opts...,
	)
	provisionPodServiceUpdatePodHandler := connect.NewUnaryHandler(
		ProvisionPodServiceUpdatePodProcedure,
		svc.UpdatePod,
		opts...,
	)
	provisionPodServiceDeletePodHandler := connect.NewUnaryHandler(
		ProvisionPodServiceDeletePodProcedure,
		svc.DeletePod,
		opts...,
	)
	provisionPodServiceGetPodLogsHandler := connect.NewServerStreamHandler(
		ProvisionPodServiceGetPodLogsProcedure,
		svc.GetPodLogs,
		opts...,
	)
	return "/apocryph.proto.v0.provisionPod.ProvisionPodService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ProvisionPodServiceProvisionPodProcedure:
			provisionPodServiceProvisionPodHandler.ServeHTTP(w, r)
		case ProvisionPodServiceUpdatePodProcedure:
			provisionPodServiceUpdatePodHandler.ServeHTTP(w, r)
		case ProvisionPodServiceDeletePodProcedure:
			provisionPodServiceDeletePodHandler.ServeHTTP(w, r)
		case ProvisionPodServiceGetPodLogsProcedure:
			provisionPodServiceGetPodLogsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedProvisionPodServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedProvisionPodServiceHandler struct{}

func (UnimplementedProvisionPodServiceHandler) ProvisionPod(context.Context, *connect.Request[proto.ProvisionPodRequest]) (*connect.Response[proto.ProvisionPodResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apocryph.proto.v0.provisionPod.ProvisionPodService.ProvisionPod is not implemented"))
}

func (UnimplementedProvisionPodServiceHandler) UpdatePod(context.Context, *connect.Request[proto.UpdatePodRequest]) (*connect.Response[proto.ProvisionPodResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apocryph.proto.v0.provisionPod.ProvisionPodService.UpdatePod is not implemented"))
}

func (UnimplementedProvisionPodServiceHandler) DeletePod(context.Context, *connect.Request[proto.DeletePodRequest]) (*connect.Response[proto.DeletePodResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apocryph.proto.v0.provisionPod.ProvisionPodService.DeletePod is not implemented"))
}

func (UnimplementedProvisionPodServiceHandler) GetPodLogs(context.Context, *connect.Request[proto.PodLogRequest], *connect.ServerStream[proto.PodLogResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("apocryph.proto.v0.provisionPod.ProvisionPodService.GetPodLogs is not implemented"))
}
