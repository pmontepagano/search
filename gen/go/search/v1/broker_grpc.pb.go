// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: search/v1/broker.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BrokerServiceClient is the client API for BrokerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BrokerServiceClient interface {
	BrokerChannel(ctx context.Context, in *BrokerChannelRequest, opts ...grpc.CallOption) (*BrokerChannelResponse, error)
	RegisterProvider(ctx context.Context, in *RegisterProviderRequest, opts ...grpc.CallOption) (*RegisterProviderResponse, error)
}

type brokerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBrokerServiceClient(cc grpc.ClientConnInterface) BrokerServiceClient {
	return &brokerServiceClient{cc}
}

func (c *brokerServiceClient) BrokerChannel(ctx context.Context, in *BrokerChannelRequest, opts ...grpc.CallOption) (*BrokerChannelResponse, error) {
	out := new(BrokerChannelResponse)
	err := c.cc.Invoke(ctx, "/search.v1.BrokerService/BrokerChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brokerServiceClient) RegisterProvider(ctx context.Context, in *RegisterProviderRequest, opts ...grpc.CallOption) (*RegisterProviderResponse, error) {
	out := new(RegisterProviderResponse)
	err := c.cc.Invoke(ctx, "/search.v1.BrokerService/RegisterProvider", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BrokerServiceServer is the server API for BrokerService service.
// All implementations must embed UnimplementedBrokerServiceServer
// for forward compatibility
type BrokerServiceServer interface {
	BrokerChannel(context.Context, *BrokerChannelRequest) (*BrokerChannelResponse, error)
	RegisterProvider(context.Context, *RegisterProviderRequest) (*RegisterProviderResponse, error)
	mustEmbedUnimplementedBrokerServiceServer()
}

// UnimplementedBrokerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBrokerServiceServer struct {
}

func (UnimplementedBrokerServiceServer) BrokerChannel(context.Context, *BrokerChannelRequest) (*BrokerChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BrokerChannel not implemented")
}
func (UnimplementedBrokerServiceServer) RegisterProvider(context.Context, *RegisterProviderRequest) (*RegisterProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterProvider not implemented")
}
func (UnimplementedBrokerServiceServer) mustEmbedUnimplementedBrokerServiceServer() {}

// UnsafeBrokerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BrokerServiceServer will
// result in compilation errors.
type UnsafeBrokerServiceServer interface {
	mustEmbedUnimplementedBrokerServiceServer()
}

func RegisterBrokerServiceServer(s grpc.ServiceRegistrar, srv BrokerServiceServer) {
	s.RegisterService(&BrokerService_ServiceDesc, srv)
}

func _BrokerService_BrokerChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BrokerChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServiceServer).BrokerChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/search.v1.BrokerService/BrokerChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServiceServer).BrokerChannel(ctx, req.(*BrokerChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrokerService_RegisterProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterProviderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrokerServiceServer).RegisterProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/search.v1.BrokerService/RegisterProvider",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrokerServiceServer).RegisterProvider(ctx, req.(*RegisterProviderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BrokerService_ServiceDesc is the grpc.ServiceDesc for BrokerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BrokerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.v1.BrokerService",
	HandlerType: (*BrokerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BrokerChannel",
			Handler:    _BrokerService_BrokerChannel_Handler,
		},
		{
			MethodName: "RegisterProvider",
			Handler:    _BrokerService_RegisterProvider_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search/v1/broker.proto",
}