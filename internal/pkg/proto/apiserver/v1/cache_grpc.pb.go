// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.20.3
// source: cache.proto

package v1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Cache_ListSecrets_FullMethodName  = "/proto.Cache/ListSecrets"
	Cache_ListPolicies_FullMethodName = "/proto.Cache/ListPolicies"
)

// CacheClient is the client API for Cache service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Cache implements a cache rpc service.
type CacheClient interface {
	ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error)
	ListPolicies(ctx context.Context, in *ListPoliciesRequest, opts ...grpc.CallOption) (*ListPoliciesResponse, error)
}

type cacheClient struct {
	cc grpc.ClientConnInterface
}

func NewCacheClient(cc grpc.ClientConnInterface) CacheClient {
	return &cacheClient{cc}
}

func (c *cacheClient) ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListSecretsResponse)
	err := c.cc.Invoke(ctx, Cache_ListSecrets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cacheClient) ListPolicies(ctx context.Context, in *ListPoliciesRequest, opts ...grpc.CallOption) (*ListPoliciesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPoliciesResponse)
	err := c.cc.Invoke(ctx, Cache_ListPolicies_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CacheServer is the server API for Cache service.
// All implementations must embed UnimplementedCacheServer
// for forward compatibility
//
// Cache implements a cache rpc service.
type CacheServer interface {
	ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error)
	ListPolicies(context.Context, *ListPoliciesRequest) (*ListPoliciesResponse, error)
	mustEmbedUnimplementedCacheServer()
}

// UnimplementedCacheServer must be embedded to have forward compatible implementations.
type UnimplementedCacheServer struct {
}

func (UnimplementedCacheServer) ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSecrets not implemented")
}
func (UnimplementedCacheServer) ListPolicies(context.Context, *ListPoliciesRequest) (*ListPoliciesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPolicies not implemented")
}
func (UnimplementedCacheServer) mustEmbedUnimplementedCacheServer() {}

// UnsafeCacheServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CacheServer will
// result in compilation errors.
type UnsafeCacheServer interface {
	mustEmbedUnimplementedCacheServer()
}

func RegisterCacheServer(s grpc.ServiceRegistrar, srv CacheServer) {
	s.RegisterService(&Cache_ServiceDesc, srv)
}

func _Cache_ListSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSecretsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacheServer).ListSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cache_ListSecrets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacheServer).ListSecrets(ctx, req.(*ListSecretsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cache_ListPolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPoliciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacheServer).ListPolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cache_ListPolicies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacheServer).ListPolicies(ctx, req.(*ListPoliciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Cache_ServiceDesc is the grpc.ServiceDesc for Cache service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cache_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Cache",
	HandlerType: (*CacheServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSecrets",
			Handler:    _Cache_ListSecrets_Handler,
		},
		{
			MethodName: "ListPolicies",
			Handler:    _Cache_ListPolicies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cache.proto",
}
