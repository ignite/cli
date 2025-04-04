// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: ignite/services/plugin/grpc/v1/service.proto

package v1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	InterfaceService_Manifest_FullMethodName           = "/ignite.services.plugin.grpc.v1.InterfaceService/Manifest"
	InterfaceService_Execute_FullMethodName            = "/ignite.services.plugin.grpc.v1.InterfaceService/Execute"
	InterfaceService_ExecuteHookPre_FullMethodName     = "/ignite.services.plugin.grpc.v1.InterfaceService/ExecuteHookPre"
	InterfaceService_ExecuteHookPost_FullMethodName    = "/ignite.services.plugin.grpc.v1.InterfaceService/ExecuteHookPost"
	InterfaceService_ExecuteHookCleanUp_FullMethodName = "/ignite.services.plugin.grpc.v1.InterfaceService/ExecuteHookCleanUp"
)

// InterfaceServiceClient is the client API for InterfaceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// InterfaceService defines the interface that must be implemented by all plugins.
type InterfaceServiceClient interface {
	// Manifest declares the plugin's Command(s) and Hook(s).
	Manifest(ctx context.Context, in *ManifestRequest, opts ...grpc.CallOption) (*ManifestResponse, error)
	// Execute will be invoked by ignite when a plugin Command is executed.
	// It is global for all commands declared in Manifest, if you have declared
	// multiple commands, use cmd.Path to distinguish them.
	Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (*ExecuteResponse, error)
	// ExecuteHookPre is invoked by ignite when a command specified by the Hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookPre(ctx context.Context, in *ExecuteHookPreRequest, opts ...grpc.CallOption) (*ExecuteHookPreResponse, error)
	// ExecuteHookPost is invoked by ignite when a command specified by the hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookPost(ctx context.Context, in *ExecuteHookPostRequest, opts ...grpc.CallOption) (*ExecuteHookPostResponse, error)
	// ExecuteHookCleanUp is invoked by ignite when a command specified by the
	// hook path is invoked. Unlike ExecuteHookPost, it is invoked regardless of
	// execution status of the command and hooks.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookCleanUp(ctx context.Context, in *ExecuteHookCleanUpRequest, opts ...grpc.CallOption) (*ExecuteHookCleanUpResponse, error)
}

type interfaceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInterfaceServiceClient(cc grpc.ClientConnInterface) InterfaceServiceClient {
	return &interfaceServiceClient{cc}
}

func (c *interfaceServiceClient) Manifest(ctx context.Context, in *ManifestRequest, opts ...grpc.CallOption) (*ManifestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ManifestResponse)
	err := c.cc.Invoke(ctx, InterfaceService_Manifest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *interfaceServiceClient) Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (*ExecuteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExecuteResponse)
	err := c.cc.Invoke(ctx, InterfaceService_Execute_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *interfaceServiceClient) ExecuteHookPre(ctx context.Context, in *ExecuteHookPreRequest, opts ...grpc.CallOption) (*ExecuteHookPreResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExecuteHookPreResponse)
	err := c.cc.Invoke(ctx, InterfaceService_ExecuteHookPre_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *interfaceServiceClient) ExecuteHookPost(ctx context.Context, in *ExecuteHookPostRequest, opts ...grpc.CallOption) (*ExecuteHookPostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExecuteHookPostResponse)
	err := c.cc.Invoke(ctx, InterfaceService_ExecuteHookPost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *interfaceServiceClient) ExecuteHookCleanUp(ctx context.Context, in *ExecuteHookCleanUpRequest, opts ...grpc.CallOption) (*ExecuteHookCleanUpResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExecuteHookCleanUpResponse)
	err := c.cc.Invoke(ctx, InterfaceService_ExecuteHookCleanUp_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InterfaceServiceServer is the server API for InterfaceService service.
// All implementations must embed UnimplementedInterfaceServiceServer
// for forward compatibility.
//
// InterfaceService defines the interface that must be implemented by all plugins.
type InterfaceServiceServer interface {
	// Manifest declares the plugin's Command(s) and Hook(s).
	Manifest(context.Context, *ManifestRequest) (*ManifestResponse, error)
	// Execute will be invoked by ignite when a plugin Command is executed.
	// It is global for all commands declared in Manifest, if you have declared
	// multiple commands, use cmd.Path to distinguish them.
	Execute(context.Context, *ExecuteRequest) (*ExecuteResponse, error)
	// ExecuteHookPre is invoked by ignite when a command specified by the Hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookPre(context.Context, *ExecuteHookPreRequest) (*ExecuteHookPreResponse, error)
	// ExecuteHookPost is invoked by ignite when a command specified by the hook
	// path is invoked.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookPost(context.Context, *ExecuteHookPostRequest) (*ExecuteHookPostResponse, error)
	// ExecuteHookCleanUp is invoked by ignite when a command specified by the
	// hook path is invoked. Unlike ExecuteHookPost, it is invoked regardless of
	// execution status of the command and hooks.
	// It is global for all hooks declared in Manifest, if you have declared
	// multiple hooks, use hook.Name to distinguish them.
	ExecuteHookCleanUp(context.Context, *ExecuteHookCleanUpRequest) (*ExecuteHookCleanUpResponse, error)
	mustEmbedUnimplementedInterfaceServiceServer()
}

// UnimplementedInterfaceServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedInterfaceServiceServer struct{}

func (UnimplementedInterfaceServiceServer) Manifest(context.Context, *ManifestRequest) (*ManifestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Manifest not implemented")
}
func (UnimplementedInterfaceServiceServer) Execute(context.Context, *ExecuteRequest) (*ExecuteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedInterfaceServiceServer) ExecuteHookPre(context.Context, *ExecuteHookPreRequest) (*ExecuteHookPreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteHookPre not implemented")
}
func (UnimplementedInterfaceServiceServer) ExecuteHookPost(context.Context, *ExecuteHookPostRequest) (*ExecuteHookPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteHookPost not implemented")
}
func (UnimplementedInterfaceServiceServer) ExecuteHookCleanUp(context.Context, *ExecuteHookCleanUpRequest) (*ExecuteHookCleanUpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteHookCleanUp not implemented")
}
func (UnimplementedInterfaceServiceServer) mustEmbedUnimplementedInterfaceServiceServer() {}
func (UnimplementedInterfaceServiceServer) testEmbeddedByValue()                          {}

// UnsafeInterfaceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InterfaceServiceServer will
// result in compilation errors.
type UnsafeInterfaceServiceServer interface {
	mustEmbedUnimplementedInterfaceServiceServer()
}

func RegisterInterfaceServiceServer(s grpc.ServiceRegistrar, srv InterfaceServiceServer) {
	// If the following call pancis, it indicates UnimplementedInterfaceServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&InterfaceService_ServiceDesc, srv)
}

func _InterfaceService_Manifest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ManifestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterfaceServiceServer).Manifest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InterfaceService_Manifest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterfaceServiceServer).Manifest(ctx, req.(*ManifestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InterfaceService_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterfaceServiceServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InterfaceService_Execute_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterfaceServiceServer).Execute(ctx, req.(*ExecuteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InterfaceService_ExecuteHookPre_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteHookPreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterfaceServiceServer).ExecuteHookPre(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InterfaceService_ExecuteHookPre_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterfaceServiceServer).ExecuteHookPre(ctx, req.(*ExecuteHookPreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InterfaceService_ExecuteHookPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteHookPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterfaceServiceServer).ExecuteHookPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InterfaceService_ExecuteHookPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterfaceServiceServer).ExecuteHookPost(ctx, req.(*ExecuteHookPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InterfaceService_ExecuteHookCleanUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteHookCleanUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterfaceServiceServer).ExecuteHookCleanUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InterfaceService_ExecuteHookCleanUp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterfaceServiceServer).ExecuteHookCleanUp(ctx, req.(*ExecuteHookCleanUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InterfaceService_ServiceDesc is the grpc.ServiceDesc for InterfaceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InterfaceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ignite.services.plugin.grpc.v1.InterfaceService",
	HandlerType: (*InterfaceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Manifest",
			Handler:    _InterfaceService_Manifest_Handler,
		},
		{
			MethodName: "Execute",
			Handler:    _InterfaceService_Execute_Handler,
		},
		{
			MethodName: "ExecuteHookPre",
			Handler:    _InterfaceService_ExecuteHookPre_Handler,
		},
		{
			MethodName: "ExecuteHookPost",
			Handler:    _InterfaceService_ExecuteHookPost_Handler,
		},
		{
			MethodName: "ExecuteHookCleanUp",
			Handler:    _InterfaceService_ExecuteHookCleanUp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ignite/services/plugin/grpc/v1/service.proto",
}

const (
	ClientAPIService_GetChainInfo_FullMethodName = "/ignite.services.plugin.grpc.v1.ClientAPIService/GetChainInfo"
)

// ClientAPIServiceClient is the client API for ClientAPIService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// ClientAPIService defines the interface that allows plugins to get chain app analysis info.
type ClientAPIServiceClient interface {
	// GetChainInfo returns basic chain info for the configured app
	GetChainInfo(ctx context.Context, in *GetChainInfoRequest, opts ...grpc.CallOption) (*GetChainInfoResponse, error)
}

type clientAPIServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewClientAPIServiceClient(cc grpc.ClientConnInterface) ClientAPIServiceClient {
	return &clientAPIServiceClient{cc}
}

func (c *clientAPIServiceClient) GetChainInfo(ctx context.Context, in *GetChainInfoRequest, opts ...grpc.CallOption) (*GetChainInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetChainInfoResponse)
	err := c.cc.Invoke(ctx, ClientAPIService_GetChainInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientAPIServiceServer is the server API for ClientAPIService service.
// All implementations must embed UnimplementedClientAPIServiceServer
// for forward compatibility.
//
// ClientAPIService defines the interface that allows plugins to get chain app analysis info.
type ClientAPIServiceServer interface {
	// GetChainInfo returns basic chain info for the configured app
	GetChainInfo(context.Context, *GetChainInfoRequest) (*GetChainInfoResponse, error)
	mustEmbedUnimplementedClientAPIServiceServer()
}

// UnimplementedClientAPIServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedClientAPIServiceServer struct{}

func (UnimplementedClientAPIServiceServer) GetChainInfo(context.Context, *GetChainInfoRequest) (*GetChainInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChainInfo not implemented")
}
func (UnimplementedClientAPIServiceServer) mustEmbedUnimplementedClientAPIServiceServer() {}
func (UnimplementedClientAPIServiceServer) testEmbeddedByValue()                          {}

// UnsafeClientAPIServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientAPIServiceServer will
// result in compilation errors.
type UnsafeClientAPIServiceServer interface {
	mustEmbedUnimplementedClientAPIServiceServer()
}

func RegisterClientAPIServiceServer(s grpc.ServiceRegistrar, srv ClientAPIServiceServer) {
	// If the following call pancis, it indicates UnimplementedClientAPIServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ClientAPIService_ServiceDesc, srv)
}

func _ClientAPIService_GetChainInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChainInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientAPIServiceServer).GetChainInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClientAPIService_GetChainInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientAPIServiceServer).GetChainInfo(ctx, req.(*GetChainInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ClientAPIService_ServiceDesc is the grpc.ServiceDesc for ClientAPIService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClientAPIService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ignite.services.plugin.grpc.v1.ClientAPIService",
	HandlerType: (*ClientAPIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetChainInfo",
			Handler:    _ClientAPIService_GetChainInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ignite/services/plugin/grpc/v1/service.proto",
}
