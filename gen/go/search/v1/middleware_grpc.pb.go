// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: search/v1/middleware.proto

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

const (
	PrivateMiddlewareService_RegisterChannel_FullMethodName = "/search.v1.PrivateMiddlewareService/RegisterChannel"
	PrivateMiddlewareService_RegisterApp_FullMethodName     = "/search.v1.PrivateMiddlewareService/RegisterApp"
	PrivateMiddlewareService_AppSend_FullMethodName         = "/search.v1.PrivateMiddlewareService/AppSend"
	PrivateMiddlewareService_AppRecv_FullMethodName         = "/search.v1.PrivateMiddlewareService/AppRecv"
)

// PrivateMiddlewareServiceClient is the client API for PrivateMiddlewareService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrivateMiddlewareServiceClient interface {
	// This is used by a requires point to start a new channel with a requirement contract.
	RegisterChannel(ctx context.Context, in *RegisterChannelRequest, opts ...grpc.CallOption) (*RegisterChannelResponse, error)
	// This is used by provider services to register their provision contract with the Registry/Broker.
	RegisterApp(ctx context.Context, in *RegisterAppRequest, opts ...grpc.CallOption) (PrivateMiddlewareService_RegisterAppClient, error)
	// This is used by the local app to communicate with other participants in an already
	// initiated or registered channel
	AppSend(ctx context.Context, in *AppSendRequest, opts ...grpc.CallOption) (*AppSendResponse, error)
	AppRecv(ctx context.Context, in *AppRecvRequest, opts ...grpc.CallOption) (*AppRecvResponse, error)
}

type privateMiddlewareServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPrivateMiddlewareServiceClient(cc grpc.ClientConnInterface) PrivateMiddlewareServiceClient {
	return &privateMiddlewareServiceClient{cc}
}

func (c *privateMiddlewareServiceClient) RegisterChannel(ctx context.Context, in *RegisterChannelRequest, opts ...grpc.CallOption) (*RegisterChannelResponse, error) {
	out := new(RegisterChannelResponse)
	err := c.cc.Invoke(ctx, PrivateMiddlewareService_RegisterChannel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privateMiddlewareServiceClient) RegisterApp(ctx context.Context, in *RegisterAppRequest, opts ...grpc.CallOption) (PrivateMiddlewareService_RegisterAppClient, error) {
	stream, err := c.cc.NewStream(ctx, &PrivateMiddlewareService_ServiceDesc.Streams[0], PrivateMiddlewareService_RegisterApp_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &privateMiddlewareServiceRegisterAppClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PrivateMiddlewareService_RegisterAppClient interface {
	Recv() (*RegisterAppResponse, error)
	grpc.ClientStream
}

type privateMiddlewareServiceRegisterAppClient struct {
	grpc.ClientStream
}

func (x *privateMiddlewareServiceRegisterAppClient) Recv() (*RegisterAppResponse, error) {
	m := new(RegisterAppResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *privateMiddlewareServiceClient) AppSend(ctx context.Context, in *AppSendRequest, opts ...grpc.CallOption) (*AppSendResponse, error) {
	out := new(AppSendResponse)
	err := c.cc.Invoke(ctx, PrivateMiddlewareService_AppSend_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *privateMiddlewareServiceClient) AppRecv(ctx context.Context, in *AppRecvRequest, opts ...grpc.CallOption) (*AppRecvResponse, error) {
	out := new(AppRecvResponse)
	err := c.cc.Invoke(ctx, PrivateMiddlewareService_AppRecv_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrivateMiddlewareServiceServer is the server API for PrivateMiddlewareService service.
// All implementations must embed UnimplementedPrivateMiddlewareServiceServer
// for forward compatibility
type PrivateMiddlewareServiceServer interface {
	// This is used by a requires point to start a new channel with a requirement contract.
	RegisterChannel(context.Context, *RegisterChannelRequest) (*RegisterChannelResponse, error)
	// This is used by provider services to register their provision contract with the Registry/Broker.
	RegisterApp(*RegisterAppRequest, PrivateMiddlewareService_RegisterAppServer) error
	// This is used by the local app to communicate with other participants in an already
	// initiated or registered channel
	AppSend(context.Context, *AppSendRequest) (*AppSendResponse, error)
	AppRecv(context.Context, *AppRecvRequest) (*AppRecvResponse, error)
	mustEmbedUnimplementedPrivateMiddlewareServiceServer()
}

// UnimplementedPrivateMiddlewareServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPrivateMiddlewareServiceServer struct {
}

func (UnimplementedPrivateMiddlewareServiceServer) RegisterChannel(context.Context, *RegisterChannelRequest) (*RegisterChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterChannel not implemented")
}
func (UnimplementedPrivateMiddlewareServiceServer) RegisterApp(*RegisterAppRequest, PrivateMiddlewareService_RegisterAppServer) error {
	return status.Errorf(codes.Unimplemented, "method RegisterApp not implemented")
}
func (UnimplementedPrivateMiddlewareServiceServer) AppSend(context.Context, *AppSendRequest) (*AppSendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppSend not implemented")
}
func (UnimplementedPrivateMiddlewareServiceServer) AppRecv(context.Context, *AppRecvRequest) (*AppRecvResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppRecv not implemented")
}
func (UnimplementedPrivateMiddlewareServiceServer) mustEmbedUnimplementedPrivateMiddlewareServiceServer() {
}

// UnsafePrivateMiddlewareServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrivateMiddlewareServiceServer will
// result in compilation errors.
type UnsafePrivateMiddlewareServiceServer interface {
	mustEmbedUnimplementedPrivateMiddlewareServiceServer()
}

func RegisterPrivateMiddlewareServiceServer(s grpc.ServiceRegistrar, srv PrivateMiddlewareServiceServer) {
	s.RegisterService(&PrivateMiddlewareService_ServiceDesc, srv)
}

func _PrivateMiddlewareService_RegisterChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateMiddlewareServiceServer).RegisterChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PrivateMiddlewareService_RegisterChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateMiddlewareServiceServer).RegisterChannel(ctx, req.(*RegisterChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivateMiddlewareService_RegisterApp_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RegisterAppRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PrivateMiddlewareServiceServer).RegisterApp(m, &privateMiddlewareServiceRegisterAppServer{stream})
}

type PrivateMiddlewareService_RegisterAppServer interface {
	Send(*RegisterAppResponse) error
	grpc.ServerStream
}

type privateMiddlewareServiceRegisterAppServer struct {
	grpc.ServerStream
}

func (x *privateMiddlewareServiceRegisterAppServer) Send(m *RegisterAppResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PrivateMiddlewareService_AppSend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppSendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateMiddlewareServiceServer).AppSend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PrivateMiddlewareService_AppSend_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateMiddlewareServiceServer).AppSend(ctx, req.(*AppSendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PrivateMiddlewareService_AppRecv_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppRecvRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrivateMiddlewareServiceServer).AppRecv(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PrivateMiddlewareService_AppRecv_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrivateMiddlewareServiceServer).AppRecv(ctx, req.(*AppRecvRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PrivateMiddlewareService_ServiceDesc is the grpc.ServiceDesc for PrivateMiddlewareService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PrivateMiddlewareService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.v1.PrivateMiddlewareService",
	HandlerType: (*PrivateMiddlewareServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterChannel",
			Handler:    _PrivateMiddlewareService_RegisterChannel_Handler,
		},
		{
			MethodName: "AppSend",
			Handler:    _PrivateMiddlewareService_AppSend_Handler,
		},
		{
			MethodName: "AppRecv",
			Handler:    _PrivateMiddlewareService_AppRecv_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RegisterApp",
			Handler:       _PrivateMiddlewareService_RegisterApp_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "search/v1/middleware.proto",
}

const (
	PublicMiddlewareService_InitChannel_FullMethodName     = "/search.v1.PublicMiddlewareService/InitChannel"
	PublicMiddlewareService_StartChannel_FullMethodName    = "/search.v1.PublicMiddlewareService/StartChannel"
	PublicMiddlewareService_MessageExchange_FullMethodName = "/search.v1.PublicMiddlewareService/MessageExchange"
)

// PublicMiddlewareServiceClient is the client API for PublicMiddlewareService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PublicMiddlewareServiceClient interface {
	// The Broker, when a new channel is registered, signals all participants with this
	InitChannel(ctx context.Context, in *InitChannelRequest, opts ...grpc.CallOption) (*InitChannelResponse, error)
	StartChannel(ctx context.Context, in *StartChannelRequest, opts ...grpc.CallOption) (*StartChannelResponse, error)
	MessageExchange(ctx context.Context, opts ...grpc.CallOption) (PublicMiddlewareService_MessageExchangeClient, error)
}

type publicMiddlewareServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPublicMiddlewareServiceClient(cc grpc.ClientConnInterface) PublicMiddlewareServiceClient {
	return &publicMiddlewareServiceClient{cc}
}

func (c *publicMiddlewareServiceClient) InitChannel(ctx context.Context, in *InitChannelRequest, opts ...grpc.CallOption) (*InitChannelResponse, error) {
	out := new(InitChannelResponse)
	err := c.cc.Invoke(ctx, PublicMiddlewareService_InitChannel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicMiddlewareServiceClient) StartChannel(ctx context.Context, in *StartChannelRequest, opts ...grpc.CallOption) (*StartChannelResponse, error) {
	out := new(StartChannelResponse)
	err := c.cc.Invoke(ctx, PublicMiddlewareService_StartChannel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicMiddlewareServiceClient) MessageExchange(ctx context.Context, opts ...grpc.CallOption) (PublicMiddlewareService_MessageExchangeClient, error) {
	stream, err := c.cc.NewStream(ctx, &PublicMiddlewareService_ServiceDesc.Streams[0], PublicMiddlewareService_MessageExchange_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &publicMiddlewareServiceMessageExchangeClient{stream}
	return x, nil
}

type PublicMiddlewareService_MessageExchangeClient interface {
	Send(*ApplicationMessageWithHeaders) error
	Recv() (*ApplicationMessageWithHeaders, error)
	grpc.ClientStream
}

type publicMiddlewareServiceMessageExchangeClient struct {
	grpc.ClientStream
}

func (x *publicMiddlewareServiceMessageExchangeClient) Send(m *ApplicationMessageWithHeaders) error {
	return x.ClientStream.SendMsg(m)
}

func (x *publicMiddlewareServiceMessageExchangeClient) Recv() (*ApplicationMessageWithHeaders, error) {
	m := new(ApplicationMessageWithHeaders)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PublicMiddlewareServiceServer is the server API for PublicMiddlewareService service.
// All implementations must embed UnimplementedPublicMiddlewareServiceServer
// for forward compatibility
type PublicMiddlewareServiceServer interface {
	// The Broker, when a new channel is registered, signals all participants with this
	InitChannel(context.Context, *InitChannelRequest) (*InitChannelResponse, error)
	StartChannel(context.Context, *StartChannelRequest) (*StartChannelResponse, error)
	MessageExchange(PublicMiddlewareService_MessageExchangeServer) error
	mustEmbedUnimplementedPublicMiddlewareServiceServer()
}

// UnimplementedPublicMiddlewareServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPublicMiddlewareServiceServer struct {
}

func (UnimplementedPublicMiddlewareServiceServer) InitChannel(context.Context, *InitChannelRequest) (*InitChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitChannel not implemented")
}
func (UnimplementedPublicMiddlewareServiceServer) StartChannel(context.Context, *StartChannelRequest) (*StartChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartChannel not implemented")
}
func (UnimplementedPublicMiddlewareServiceServer) MessageExchange(PublicMiddlewareService_MessageExchangeServer) error {
	return status.Errorf(codes.Unimplemented, "method MessageExchange not implemented")
}
func (UnimplementedPublicMiddlewareServiceServer) mustEmbedUnimplementedPublicMiddlewareServiceServer() {
}

// UnsafePublicMiddlewareServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PublicMiddlewareServiceServer will
// result in compilation errors.
type UnsafePublicMiddlewareServiceServer interface {
	mustEmbedUnimplementedPublicMiddlewareServiceServer()
}

func RegisterPublicMiddlewareServiceServer(s grpc.ServiceRegistrar, srv PublicMiddlewareServiceServer) {
	s.RegisterService(&PublicMiddlewareService_ServiceDesc, srv)
}

func _PublicMiddlewareService_InitChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicMiddlewareServiceServer).InitChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PublicMiddlewareService_InitChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicMiddlewareServiceServer).InitChannel(ctx, req.(*InitChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PublicMiddlewareService_StartChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicMiddlewareServiceServer).StartChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PublicMiddlewareService_StartChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicMiddlewareServiceServer).StartChannel(ctx, req.(*StartChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PublicMiddlewareService_MessageExchange_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PublicMiddlewareServiceServer).MessageExchange(&publicMiddlewareServiceMessageExchangeServer{stream})
}

type PublicMiddlewareService_MessageExchangeServer interface {
	Send(*ApplicationMessageWithHeaders) error
	Recv() (*ApplicationMessageWithHeaders, error)
	grpc.ServerStream
}

type publicMiddlewareServiceMessageExchangeServer struct {
	grpc.ServerStream
}

func (x *publicMiddlewareServiceMessageExchangeServer) Send(m *ApplicationMessageWithHeaders) error {
	return x.ServerStream.SendMsg(m)
}

func (x *publicMiddlewareServiceMessageExchangeServer) Recv() (*ApplicationMessageWithHeaders, error) {
	m := new(ApplicationMessageWithHeaders)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PublicMiddlewareService_ServiceDesc is the grpc.ServiceDesc for PublicMiddlewareService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PublicMiddlewareService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.v1.PublicMiddlewareService",
	HandlerType: (*PublicMiddlewareServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitChannel",
			Handler:    _PublicMiddlewareService_InitChannel_Handler,
		},
		{
			MethodName: "StartChannel",
			Handler:    _PublicMiddlewareService_StartChannel_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MessageExchange",
			Handler:       _PublicMiddlewareService_MessageExchange_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "search/v1/middleware.proto",
}
