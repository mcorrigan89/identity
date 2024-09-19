// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: serviceapis/identity/v1/identity.proto

package identityv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/mcorrigan89/identity/internal/api/serviceapis/identity/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// IdentityServiceName is the fully-qualified name of the IdentityService service.
	IdentityServiceName = "serviceapis.identity.v1.IdentityService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// IdentityServiceGetUserByIdProcedure is the fully-qualified name of the IdentityService's
	// GetUserById RPC.
	IdentityServiceGetUserByIdProcedure = "/serviceapis.identity.v1.IdentityService/GetUserById"
	// IdentityServiceCreateUserProcedure is the fully-qualified name of the IdentityService's
	// CreateUser RPC.
	IdentityServiceCreateUserProcedure = "/serviceapis.identity.v1.IdentityService/CreateUser"
	// IdentityServiceAuthenticateWithGoogleCodeProcedure is the fully-qualified name of the
	// IdentityService's AuthenticateWithGoogleCode RPC.
	IdentityServiceAuthenticateWithGoogleCodeProcedure = "/serviceapis.identity.v1.IdentityService/AuthenticateWithGoogleCode"
	// IdentityServiceAuthenticateWithPasswordProcedure is the fully-qualified name of the
	// IdentityService's AuthenticateWithPassword RPC.
	IdentityServiceAuthenticateWithPasswordProcedure = "/serviceapis.identity.v1.IdentityService/AuthenticateWithPassword"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	identityServiceServiceDescriptor                          = v1.File_serviceapis_identity_v1_identity_proto.Services().ByName("IdentityService")
	identityServiceGetUserByIdMethodDescriptor                = identityServiceServiceDescriptor.Methods().ByName("GetUserById")
	identityServiceCreateUserMethodDescriptor                 = identityServiceServiceDescriptor.Methods().ByName("CreateUser")
	identityServiceAuthenticateWithGoogleCodeMethodDescriptor = identityServiceServiceDescriptor.Methods().ByName("AuthenticateWithGoogleCode")
	identityServiceAuthenticateWithPasswordMethodDescriptor   = identityServiceServiceDescriptor.Methods().ByName("AuthenticateWithPassword")
)

// IdentityServiceClient is a client for the serviceapis.identity.v1.IdentityService service.
type IdentityServiceClient interface {
	GetUserById(context.Context, *connect.Request[v1.GetUserByIdRequest]) (*connect.Response[v1.GetUserByIdResponse], error)
	CreateUser(context.Context, *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreateUserResponse], error)
	AuthenticateWithGoogleCode(context.Context, *connect.Request[v1.AuthenticateWithGoogleCodeRequest]) (*connect.Response[v1.AuthenticateWithGoogleCodeResponse], error)
	AuthenticateWithPassword(context.Context, *connect.Request[v1.AuthenticateWithPasswordRequest]) (*connect.Response[v1.AuthenticateWithPasswordResponse], error)
}

// NewIdentityServiceClient constructs a client for the serviceapis.identity.v1.IdentityService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewIdentityServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) IdentityServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &identityServiceClient{
		getUserById: connect.NewClient[v1.GetUserByIdRequest, v1.GetUserByIdResponse](
			httpClient,
			baseURL+IdentityServiceGetUserByIdProcedure,
			connect.WithSchema(identityServiceGetUserByIdMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createUser: connect.NewClient[v1.CreateUserRequest, v1.CreateUserResponse](
			httpClient,
			baseURL+IdentityServiceCreateUserProcedure,
			connect.WithSchema(identityServiceCreateUserMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		authenticateWithGoogleCode: connect.NewClient[v1.AuthenticateWithGoogleCodeRequest, v1.AuthenticateWithGoogleCodeResponse](
			httpClient,
			baseURL+IdentityServiceAuthenticateWithGoogleCodeProcedure,
			connect.WithSchema(identityServiceAuthenticateWithGoogleCodeMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		authenticateWithPassword: connect.NewClient[v1.AuthenticateWithPasswordRequest, v1.AuthenticateWithPasswordResponse](
			httpClient,
			baseURL+IdentityServiceAuthenticateWithPasswordProcedure,
			connect.WithSchema(identityServiceAuthenticateWithPasswordMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// identityServiceClient implements IdentityServiceClient.
type identityServiceClient struct {
	getUserById                *connect.Client[v1.GetUserByIdRequest, v1.GetUserByIdResponse]
	createUser                 *connect.Client[v1.CreateUserRequest, v1.CreateUserResponse]
	authenticateWithGoogleCode *connect.Client[v1.AuthenticateWithGoogleCodeRequest, v1.AuthenticateWithGoogleCodeResponse]
	authenticateWithPassword   *connect.Client[v1.AuthenticateWithPasswordRequest, v1.AuthenticateWithPasswordResponse]
}

// GetUserById calls serviceapis.identity.v1.IdentityService.GetUserById.
func (c *identityServiceClient) GetUserById(ctx context.Context, req *connect.Request[v1.GetUserByIdRequest]) (*connect.Response[v1.GetUserByIdResponse], error) {
	return c.getUserById.CallUnary(ctx, req)
}

// CreateUser calls serviceapis.identity.v1.IdentityService.CreateUser.
func (c *identityServiceClient) CreateUser(ctx context.Context, req *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreateUserResponse], error) {
	return c.createUser.CallUnary(ctx, req)
}

// AuthenticateWithGoogleCode calls
// serviceapis.identity.v1.IdentityService.AuthenticateWithGoogleCode.
func (c *identityServiceClient) AuthenticateWithGoogleCode(ctx context.Context, req *connect.Request[v1.AuthenticateWithGoogleCodeRequest]) (*connect.Response[v1.AuthenticateWithGoogleCodeResponse], error) {
	return c.authenticateWithGoogleCode.CallUnary(ctx, req)
}

// AuthenticateWithPassword calls serviceapis.identity.v1.IdentityService.AuthenticateWithPassword.
func (c *identityServiceClient) AuthenticateWithPassword(ctx context.Context, req *connect.Request[v1.AuthenticateWithPasswordRequest]) (*connect.Response[v1.AuthenticateWithPasswordResponse], error) {
	return c.authenticateWithPassword.CallUnary(ctx, req)
}

// IdentityServiceHandler is an implementation of the serviceapis.identity.v1.IdentityService
// service.
type IdentityServiceHandler interface {
	GetUserById(context.Context, *connect.Request[v1.GetUserByIdRequest]) (*connect.Response[v1.GetUserByIdResponse], error)
	CreateUser(context.Context, *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreateUserResponse], error)
	AuthenticateWithGoogleCode(context.Context, *connect.Request[v1.AuthenticateWithGoogleCodeRequest]) (*connect.Response[v1.AuthenticateWithGoogleCodeResponse], error)
	AuthenticateWithPassword(context.Context, *connect.Request[v1.AuthenticateWithPasswordRequest]) (*connect.Response[v1.AuthenticateWithPasswordResponse], error)
}

// NewIdentityServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewIdentityServiceHandler(svc IdentityServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	identityServiceGetUserByIdHandler := connect.NewUnaryHandler(
		IdentityServiceGetUserByIdProcedure,
		svc.GetUserById,
		connect.WithSchema(identityServiceGetUserByIdMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	identityServiceCreateUserHandler := connect.NewUnaryHandler(
		IdentityServiceCreateUserProcedure,
		svc.CreateUser,
		connect.WithSchema(identityServiceCreateUserMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	identityServiceAuthenticateWithGoogleCodeHandler := connect.NewUnaryHandler(
		IdentityServiceAuthenticateWithGoogleCodeProcedure,
		svc.AuthenticateWithGoogleCode,
		connect.WithSchema(identityServiceAuthenticateWithGoogleCodeMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	identityServiceAuthenticateWithPasswordHandler := connect.NewUnaryHandler(
		IdentityServiceAuthenticateWithPasswordProcedure,
		svc.AuthenticateWithPassword,
		connect.WithSchema(identityServiceAuthenticateWithPasswordMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/serviceapis.identity.v1.IdentityService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case IdentityServiceGetUserByIdProcedure:
			identityServiceGetUserByIdHandler.ServeHTTP(w, r)
		case IdentityServiceCreateUserProcedure:
			identityServiceCreateUserHandler.ServeHTTP(w, r)
		case IdentityServiceAuthenticateWithGoogleCodeProcedure:
			identityServiceAuthenticateWithGoogleCodeHandler.ServeHTTP(w, r)
		case IdentityServiceAuthenticateWithPasswordProcedure:
			identityServiceAuthenticateWithPasswordHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedIdentityServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedIdentityServiceHandler struct{}

func (UnimplementedIdentityServiceHandler) GetUserById(context.Context, *connect.Request[v1.GetUserByIdRequest]) (*connect.Response[v1.GetUserByIdResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("serviceapis.identity.v1.IdentityService.GetUserById is not implemented"))
}

func (UnimplementedIdentityServiceHandler) CreateUser(context.Context, *connect.Request[v1.CreateUserRequest]) (*connect.Response[v1.CreateUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("serviceapis.identity.v1.IdentityService.CreateUser is not implemented"))
}

func (UnimplementedIdentityServiceHandler) AuthenticateWithGoogleCode(context.Context, *connect.Request[v1.AuthenticateWithGoogleCodeRequest]) (*connect.Response[v1.AuthenticateWithGoogleCodeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("serviceapis.identity.v1.IdentityService.AuthenticateWithGoogleCode is not implemented"))
}

func (UnimplementedIdentityServiceHandler) AuthenticateWithPassword(context.Context, *connect.Request[v1.AuthenticateWithPasswordRequest]) (*connect.Response[v1.AuthenticateWithPasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("serviceapis.identity.v1.IdentityService.AuthenticateWithPassword is not implemented"))
}
