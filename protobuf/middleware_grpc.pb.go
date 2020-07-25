// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PrivateMiddlewareClient is the client API for PrivateMiddleware service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrivateMiddlewareClient interface {
	// This is used by the initiator to register a new channel
	RegisterChannel(ctx context.Context, in *RegisterChannelRequest, opts ...grpc.CallOption) (*RegisterChannelResponse, error)
	// This is used by the local app to register an app with the registry
	RegisterApp(ctx context.Context, in *RegisterAppRequest, opts ...grpc.CallOption) (*RegisterAppResponse, error)
	// This is used by the local app to communicate with other participants in an already
	// initiated or registered channel
	UseChannel(ctx context.Context, opts ...grpc.CallOption) (PrivateMiddleware_UseChannelClient, error)
}

type privateMiddlewareClient struct {
	cc grpc.ClientConnInterface
}

func NewPrivateMiddlewareClient(cc grpc.ClientConnInterface) PrivateMiddlewareClient {
	return &privateMiddlewareClient{cc}
}

func (c *privateMiddlewareClient) RegisterChannel(ctx context.Context, in *RegisterChannelRequest, opts ...grpc.CallOption) (*RegisterChannelResponse, error) {
	out := new(RegisterChannelResponse)
	err := c.cc.Invoke(ctx, "/protobuf.PrivateMiddleware/RegisterChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privateMiddlewareClient) RegisterApp(ctx context.Context, in *RegisterAppRequest, opts ...grpc.CallOption) (*RegisterAppResponse, error) {
	out := new(RegisterAppResponse)
	err := c.cc.Invoke(ctx, "/protobuf.PrivateMiddleware/RegisterApp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privateMiddlewareClient) UseChannel(ctx context.Context, opts ...grpc.CallOption) (PrivateMiddleware_UseChannelClient, error) {
	stream, err := c.cc.NewStream(ctx, &_PrivateMiddleware_serviceDesc.Streams[0], "/protobuf.PrivateMiddleware/UseChannel", opts...)
	if err != nil {
		return nil, err
	}
	x := &privateMiddlewareUseChannelClient{stream}
	return x, nil
}

type PrivateMiddleware_UseChannelClient interface {
	Send(*ApplicationMessageOut) error
	Recv() (*ApplicationMessageIn, error)
	grpc.ClientStream
}

type privateMiddlewareUseChannelClient struct {
	grpc.ClientStream
}

func (x *privateMiddlewareUseChannelClient) Send(m *ApplicationMessageOut) error {
	return x.ClientStream.SendMsg(m)
}

func (x *privateMiddlewareUseChannelClient) Recv() (*ApplicationMessageIn, error) {
	m := new(ApplicationMessageIn)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PrivateMiddlewareServer is the server API for PrivateMiddleware service.
// All implementations must embed UnimplementedPrivateMiddlewareServer
// for forward compatibility
type PrivateMiddlewareServer interface {
	// This is used by the initiator to register a new channel
	RegisterChannel(context.Context, *RegisterChannelRequest) (*RegisterChannelResponse, error)
	// This is used by the local app to register an app with the registry
	RegisterApp(context.Context, *RegisterAppRequest) (*RegisterAppResponse, error)
	// This is used by the local app to communicate with other participants in an already
	// initiated or registered channel
	UseChannel(PrivateMiddleware_UseChannelServer) error
	mustEmbedUnimplementedPrivateMiddlewareServer()
}

// UnimplementedPrivateMiddlewareServer must be embedded to have forward compatible implementations.
type UnimplementedPrivateMiddlewareServer struct {
}

func (*UnimplementedPrivateMiddlewareServer) RegisterChannel(context.Context, *RegisterChannelRequest) (*RegisterChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterChannel not implemented")
}
func (*UnimplementedPrivateMiddlewareServer) RegisterApp(context.Context, *RegisterAppRequest) (*RegisterAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterApp not implemented")
}
func (*UnimplementedPrivateMiddlewareServer) UseChannel(PrivateMiddleware_UseChannelServer) error {
	return status.Errorf(codes.Unimplemented, "method UseChannel not implemented")
}
func (*UnimplementedPrivateMiddlewareServer) mustEmbedUnimplementedPrivateMiddlewareServer() {}

func RegisterPrivateMiddlewareServer(s *grpc.Server, srv PrivateMiddlewareServer) {
	s.RegisterService(&_PrivateMiddleware_serviceDesc, srv)
}

func _PrivateMiddleware_RegisterChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateMiddlewareServer).RegisterChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.PrivateMiddleware/RegisterChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateMiddlewareServer).RegisterChannel(ctx, req.(*RegisterChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivateMiddleware_RegisterApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateMiddlewareServer).RegisterApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.PrivateMiddleware/RegisterApp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateMiddlewareServer).RegisterApp(ctx, req.(*RegisterAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivateMiddleware_UseChannel_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PrivateMiddlewareServer).UseChannel(&privateMiddlewareUseChannelServer{stream})
}

type PrivateMiddleware_UseChannelServer interface {
	Send(*ApplicationMessageIn) error
	Recv() (*ApplicationMessageOut, error)
	grpc.ServerStream
}

type privateMiddlewareUseChannelServer struct {
	grpc.ServerStream
}

func (x *privateMiddlewareUseChannelServer) Send(m *ApplicationMessageIn) error {
	return x.ServerStream.SendMsg(m)
}

func (x *privateMiddlewareUseChannelServer) Recv() (*ApplicationMessageOut, error) {
	m := new(ApplicationMessageOut)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _PrivateMiddleware_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.PrivateMiddleware",
	HandlerType: (*PrivateMiddlewareServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterChannel",
			Handler:    _PrivateMiddleware_RegisterChannel_Handler,
		},
		{
			MethodName: "RegisterApp",
			Handler:    _PrivateMiddleware_RegisterApp_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UseChannel",
			Handler:       _PrivateMiddleware_UseChannel_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "protobuf/middleware.proto",
}

// PublicMiddlewareClient is the client API for PublicMiddleware service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PublicMiddlewareClient interface {
	// The Broker, when a new channel is registered, signals all participants (except initiator) with this
	InitChannel(ctx context.Context, in *InitChannelRequest, opts ...grpc.CallOption) (*InitChannelResponse, error)
	MessageExchange(ctx context.Context, opts ...grpc.CallOption) (PublicMiddleware_MessageExchangeClient, error)
}

type publicMiddlewareClient struct {
	cc grpc.ClientConnInterface
}

func NewPublicMiddlewareClient(cc grpc.ClientConnInterface) PublicMiddlewareClient {
	return &publicMiddlewareClient{cc}
}

func (c *publicMiddlewareClient) InitChannel(ctx context.Context, in *InitChannelRequest, opts ...grpc.CallOption) (*InitChannelResponse, error) {
	out := new(InitChannelResponse)
	err := c.cc.Invoke(ctx, "/protobuf.PublicMiddleware/InitChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicMiddlewareClient) MessageExchange(ctx context.Context, opts ...grpc.CallOption) (PublicMiddleware_MessageExchangeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_PublicMiddleware_serviceDesc.Streams[0], "/protobuf.PublicMiddleware/MessageExchange", opts...)
	if err != nil {
		return nil, err
	}
	x := &publicMiddlewareMessageExchangeClient{stream}
	return x, nil
}

type PublicMiddleware_MessageExchangeClient interface {
	Send(*ApplicationMessageWithHeaders) error
	Recv() (*ApplicationMessageWithHeaders, error)
	grpc.ClientStream
}

type publicMiddlewareMessageExchangeClient struct {
	grpc.ClientStream
}

func (x *publicMiddlewareMessageExchangeClient) Send(m *ApplicationMessageWithHeaders) error {
	return x.ClientStream.SendMsg(m)
}

func (x *publicMiddlewareMessageExchangeClient) Recv() (*ApplicationMessageWithHeaders, error) {
	m := new(ApplicationMessageWithHeaders)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PublicMiddlewareServer is the server API for PublicMiddleware service.
// All implementations must embed UnimplementedPublicMiddlewareServer
// for forward compatibility
type PublicMiddlewareServer interface {
	// The Broker, when a new channel is registered, signals all participants (except initiator) with this
	InitChannel(context.Context, *InitChannelRequest) (*InitChannelResponse, error)
	MessageExchange(PublicMiddleware_MessageExchangeServer) error
	mustEmbedUnimplementedPublicMiddlewareServer()
}

// UnimplementedPublicMiddlewareServer must be embedded to have forward compatible implementations.
type UnimplementedPublicMiddlewareServer struct {
}

func (*UnimplementedPublicMiddlewareServer) InitChannel(context.Context, *InitChannelRequest) (*InitChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitChannel not implemented")
}
func (*UnimplementedPublicMiddlewareServer) MessageExchange(PublicMiddleware_MessageExchangeServer) error {
	return status.Errorf(codes.Unimplemented, "method MessageExchange not implemented")
}
func (*UnimplementedPublicMiddlewareServer) mustEmbedUnimplementedPublicMiddlewareServer() {}

func RegisterPublicMiddlewareServer(s *grpc.Server, srv PublicMiddlewareServer) {
	s.RegisterService(&_PublicMiddleware_serviceDesc, srv)
}

func _PublicMiddleware_InitChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicMiddlewareServer).InitChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.PublicMiddleware/InitChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicMiddlewareServer).InitChannel(ctx, req.(*InitChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PublicMiddleware_MessageExchange_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PublicMiddlewareServer).MessageExchange(&publicMiddlewareMessageExchangeServer{stream})
}

type PublicMiddleware_MessageExchangeServer interface {
	Send(*ApplicationMessageWithHeaders) error
	Recv() (*ApplicationMessageWithHeaders, error)
	grpc.ServerStream
}

type publicMiddlewareMessageExchangeServer struct {
	grpc.ServerStream
}

func (x *publicMiddlewareMessageExchangeServer) Send(m *ApplicationMessageWithHeaders) error {
	return x.ServerStream.SendMsg(m)
}

func (x *publicMiddlewareMessageExchangeServer) Recv() (*ApplicationMessageWithHeaders, error) {
	m := new(ApplicationMessageWithHeaders)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _PublicMiddleware_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.PublicMiddleware",
	HandlerType: (*PublicMiddlewareServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitChannel",
			Handler:    _PublicMiddleware_InitChannel_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MessageExchange",
			Handler:       _PublicMiddleware_MessageExchange_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "protobuf/middleware.proto",
}
