// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: rate/v1/rate_service.proto

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
	RateLimiterService_CheckRateLimit_FullMethodName      = "/rateLimiter.RateLimiterService/CheckRateLimit"
	RateLimiterService_GetUserRateLimit_FullMethodName    = "/rateLimiter.RateLimiterService/GetUserRateLimit"
	RateLimiterService_UpdateUserRateLimit_FullMethodName = "/rateLimiter.RateLimiterService/UpdateUserRateLimit"
)

// RateLimiterServiceClient is the client API for RateLimiterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// RateLimiter service definition
type RateLimiterServiceClient interface {
	// Check if a request is allowed based on the user's rate limit
	CheckRateLimit(ctx context.Context, in *CheckRateLimitRequest, opts ...grpc.CallOption) (*CheckRateLimitResponse, error)
	// Get the current rate limit for a specific user
	GetUserRateLimit(ctx context.Context, in *GetUserRateLimitRequest, opts ...grpc.CallOption) (*GetUserRateLimitResponse, error)
	// Update the rate limit for a specific user (e.g., admins can increase or decrease the limit)
	UpdateUserRateLimit(ctx context.Context, in *UpdateUserRateLimitRequest, opts ...grpc.CallOption) (*UpdateUserRateLimitResponse, error)
}

type rateLimiterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRateLimiterServiceClient(cc grpc.ClientConnInterface) RateLimiterServiceClient {
	return &rateLimiterServiceClient{cc}
}

func (c *rateLimiterServiceClient) CheckRateLimit(ctx context.Context, in *CheckRateLimitRequest, opts ...grpc.CallOption) (*CheckRateLimitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckRateLimitResponse)
	err := c.cc.Invoke(ctx, RateLimiterService_CheckRateLimit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rateLimiterServiceClient) GetUserRateLimit(ctx context.Context, in *GetUserRateLimitRequest, opts ...grpc.CallOption) (*GetUserRateLimitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserRateLimitResponse)
	err := c.cc.Invoke(ctx, RateLimiterService_GetUserRateLimit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rateLimiterServiceClient) UpdateUserRateLimit(ctx context.Context, in *UpdateUserRateLimitRequest, opts ...grpc.CallOption) (*UpdateUserRateLimitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateUserRateLimitResponse)
	err := c.cc.Invoke(ctx, RateLimiterService_UpdateUserRateLimit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RateLimiterServiceServer is the server API for RateLimiterService service.
// All implementations must embed UnimplementedRateLimiterServiceServer
// for forward compatibility
//
// RateLimiter service definition
type RateLimiterServiceServer interface {
	// Check if a request is allowed based on the user's rate limit
	CheckRateLimit(context.Context, *CheckRateLimitRequest) (*CheckRateLimitResponse, error)
	// Get the current rate limit for a specific user
	GetUserRateLimit(context.Context, *GetUserRateLimitRequest) (*GetUserRateLimitResponse, error)
	// Update the rate limit for a specific user (e.g., admins can increase or decrease the limit)
	UpdateUserRateLimit(context.Context, *UpdateUserRateLimitRequest) (*UpdateUserRateLimitResponse, error)
	mustEmbedUnimplementedRateLimiterServiceServer()
}

// UnimplementedRateLimiterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRateLimiterServiceServer struct {
}

func (UnimplementedRateLimiterServiceServer) CheckRateLimit(context.Context, *CheckRateLimitRequest) (*CheckRateLimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckRateLimit not implemented")
}
func (UnimplementedRateLimiterServiceServer) GetUserRateLimit(context.Context, *GetUserRateLimitRequest) (*GetUserRateLimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserRateLimit not implemented")
}
func (UnimplementedRateLimiterServiceServer) UpdateUserRateLimit(context.Context, *UpdateUserRateLimitRequest) (*UpdateUserRateLimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserRateLimit not implemented")
}
func (UnimplementedRateLimiterServiceServer) mustEmbedUnimplementedRateLimiterServiceServer() {}

// UnsafeRateLimiterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RateLimiterServiceServer will
// result in compilation errors.
type UnsafeRateLimiterServiceServer interface {
	mustEmbedUnimplementedRateLimiterServiceServer()
}

func RegisterRateLimiterServiceServer(s grpc.ServiceRegistrar, srv RateLimiterServiceServer) {
	s.RegisterService(&RateLimiterService_ServiceDesc, srv)
}

func _RateLimiterService_CheckRateLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRateLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimiterServiceServer).CheckRateLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimiterService_CheckRateLimit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimiterServiceServer).CheckRateLimit(ctx, req.(*CheckRateLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RateLimiterService_GetUserRateLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRateLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimiterServiceServer).GetUserRateLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimiterService_GetUserRateLimit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimiterServiceServer).GetUserRateLimit(ctx, req.(*GetUserRateLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RateLimiterService_UpdateUserRateLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRateLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimiterServiceServer).UpdateUserRateLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimiterService_UpdateUserRateLimit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimiterServiceServer).UpdateUserRateLimit(ctx, req.(*UpdateUserRateLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RateLimiterService_ServiceDesc is the grpc.ServiceDesc for RateLimiterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RateLimiterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rateLimiter.RateLimiterService",
	HandlerType: (*RateLimiterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckRateLimit",
			Handler:    _RateLimiterService_CheckRateLimit_Handler,
		},
		{
			MethodName: "GetUserRateLimit",
			Handler:    _RateLimiterService_GetUserRateLimit_Handler,
		},
		{
			MethodName: "UpdateUserRateLimit",
			Handler:    _RateLimiterService_UpdateUserRateLimit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rate/v1/rate_service.proto",
}
