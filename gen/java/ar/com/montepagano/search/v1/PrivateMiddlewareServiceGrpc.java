package ar.com.montepagano.search.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * 
 *This service is what a Middleware exposes to its local users (not on the internet.)
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.58.0)",
    comments = "Source: search/v1/middleware.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class PrivateMiddlewareServiceGrpc {

  private PrivateMiddlewareServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "search.v1.PrivateMiddlewareService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest,
      ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse> getRegisterChannelMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "RegisterChannel",
      requestType = ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest.class,
      responseType = ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest,
      ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse> getRegisterChannelMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest, ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse> getRegisterChannelMethod;
    if ((getRegisterChannelMethod = PrivateMiddlewareServiceGrpc.getRegisterChannelMethod) == null) {
      synchronized (PrivateMiddlewareServiceGrpc.class) {
        if ((getRegisterChannelMethod = PrivateMiddlewareServiceGrpc.getRegisterChannelMethod) == null) {
          PrivateMiddlewareServiceGrpc.getRegisterChannelMethod = getRegisterChannelMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest, ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "RegisterChannel"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PrivateMiddlewareServiceMethodDescriptorSupplier("RegisterChannel"))
              .build();
        }
      }
    }
    return getRegisterChannelMethod;
  }

  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.RegisterAppRequest,
      ar.com.montepagano.search.v1.Middleware.RegisterAppResponse> getRegisterAppMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "RegisterApp",
      requestType = ar.com.montepagano.search.v1.Middleware.RegisterAppRequest.class,
      responseType = ar.com.montepagano.search.v1.Middleware.RegisterAppResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.RegisterAppRequest,
      ar.com.montepagano.search.v1.Middleware.RegisterAppResponse> getRegisterAppMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.RegisterAppRequest, ar.com.montepagano.search.v1.Middleware.RegisterAppResponse> getRegisterAppMethod;
    if ((getRegisterAppMethod = PrivateMiddlewareServiceGrpc.getRegisterAppMethod) == null) {
      synchronized (PrivateMiddlewareServiceGrpc.class) {
        if ((getRegisterAppMethod = PrivateMiddlewareServiceGrpc.getRegisterAppMethod) == null) {
          PrivateMiddlewareServiceGrpc.getRegisterAppMethod = getRegisterAppMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.Middleware.RegisterAppRequest, ar.com.montepagano.search.v1.Middleware.RegisterAppResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "RegisterApp"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.RegisterAppRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.RegisterAppResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PrivateMiddlewareServiceMethodDescriptorSupplier("RegisterApp"))
              .build();
        }
      }
    }
    return getRegisterAppMethod;
  }

  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.CloseChannelRequest,
      ar.com.montepagano.search.v1.Middleware.CloseChannelResponse> getCloseChannelMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CloseChannel",
      requestType = ar.com.montepagano.search.v1.Middleware.CloseChannelRequest.class,
      responseType = ar.com.montepagano.search.v1.Middleware.CloseChannelResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.CloseChannelRequest,
      ar.com.montepagano.search.v1.Middleware.CloseChannelResponse> getCloseChannelMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.CloseChannelRequest, ar.com.montepagano.search.v1.Middleware.CloseChannelResponse> getCloseChannelMethod;
    if ((getCloseChannelMethod = PrivateMiddlewareServiceGrpc.getCloseChannelMethod) == null) {
      synchronized (PrivateMiddlewareServiceGrpc.class) {
        if ((getCloseChannelMethod = PrivateMiddlewareServiceGrpc.getCloseChannelMethod) == null) {
          PrivateMiddlewareServiceGrpc.getCloseChannelMethod = getCloseChannelMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.Middleware.CloseChannelRequest, ar.com.montepagano.search.v1.Middleware.CloseChannelResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CloseChannel"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.CloseChannelRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.CloseChannelResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PrivateMiddlewareServiceMethodDescriptorSupplier("CloseChannel"))
              .build();
        }
      }
    }
    return getCloseChannelMethod;
  }

  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest,
      ar.com.montepagano.search.v1.Middleware.AppSendResponse> getAppSendMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "AppSend",
      requestType = ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest.class,
      responseType = ar.com.montepagano.search.v1.Middleware.AppSendResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest,
      ar.com.montepagano.search.v1.Middleware.AppSendResponse> getAppSendMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest, ar.com.montepagano.search.v1.Middleware.AppSendResponse> getAppSendMethod;
    if ((getAppSendMethod = PrivateMiddlewareServiceGrpc.getAppSendMethod) == null) {
      synchronized (PrivateMiddlewareServiceGrpc.class) {
        if ((getAppSendMethod = PrivateMiddlewareServiceGrpc.getAppSendMethod) == null) {
          PrivateMiddlewareServiceGrpc.getAppSendMethod = getAppSendMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest, ar.com.montepagano.search.v1.Middleware.AppSendResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "AppSend"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.AppSendResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PrivateMiddlewareServiceMethodDescriptorSupplier("AppSend"))
              .build();
        }
      }
    }
    return getAppSendMethod;
  }

  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.AppRecvRequest,
      ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse> getAppRecvMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "AppRecv",
      requestType = ar.com.montepagano.search.v1.Middleware.AppRecvRequest.class,
      responseType = ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.AppRecvRequest,
      ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse> getAppRecvMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.AppRecvRequest, ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse> getAppRecvMethod;
    if ((getAppRecvMethod = PrivateMiddlewareServiceGrpc.getAppRecvMethod) == null) {
      synchronized (PrivateMiddlewareServiceGrpc.class) {
        if ((getAppRecvMethod = PrivateMiddlewareServiceGrpc.getAppRecvMethod) == null) {
          PrivateMiddlewareServiceGrpc.getAppRecvMethod = getAppRecvMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.Middleware.AppRecvRequest, ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "AppRecv"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.AppRecvRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PrivateMiddlewareServiceMethodDescriptorSupplier("AppRecv"))
              .build();
        }
      }
    }
    return getAppRecvMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static PrivateMiddlewareServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PrivateMiddlewareServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PrivateMiddlewareServiceStub>() {
        @java.lang.Override
        public PrivateMiddlewareServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PrivateMiddlewareServiceStub(channel, callOptions);
        }
      };
    return PrivateMiddlewareServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static PrivateMiddlewareServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PrivateMiddlewareServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PrivateMiddlewareServiceBlockingStub>() {
        @java.lang.Override
        public PrivateMiddlewareServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PrivateMiddlewareServiceBlockingStub(channel, callOptions);
        }
      };
    return PrivateMiddlewareServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static PrivateMiddlewareServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PrivateMiddlewareServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PrivateMiddlewareServiceFutureStub>() {
        @java.lang.Override
        public PrivateMiddlewareServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PrivateMiddlewareServiceFutureStub(channel, callOptions);
        }
      };
    return PrivateMiddlewareServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * 
   *This service is what a Middleware exposes to its local users (not on the internet.)
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * This is used by a requires point to start a new channel with a requirement contract.
     * </pre>
     */
    default void registerChannel(ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getRegisterChannelMethod(), responseObserver);
    }

    /**
     * <pre>
     * This is used by provider services to register their provision contract with the Registry/Broker.
     * </pre>
     */
    default void registerApp(ar.com.montepagano.search.v1.Middleware.RegisterAppRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.RegisterAppResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getRegisterAppMethod(), responseObserver);
    }

    /**
     * <pre>
     * This is used by local app (be it a Service Client or a Service Provider) to close a channel.
     * </pre>
     */
    default void closeChannel(ar.com.montepagano.search.v1.Middleware.CloseChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.CloseChannelResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCloseChannelMethod(), responseObserver);
    }

    /**
     * <pre>
     * This is used by the local app to communicate with other participants in an already
     * initiated or registered channel
     * </pre>
     */
    default void appSend(ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.AppSendResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getAppSendMethod(), responseObserver);
    }

    /**
     */
    default void appRecv(ar.com.montepagano.search.v1.Middleware.AppRecvRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getAppRecvMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service PrivateMiddlewareService.
   * <pre>
   * 
   *This service is what a Middleware exposes to its local users (not on the internet.)
   * </pre>
   */
  public static abstract class PrivateMiddlewareServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return PrivateMiddlewareServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service PrivateMiddlewareService.
   * <pre>
   * 
   *This service is what a Middleware exposes to its local users (not on the internet.)
   * </pre>
   */
  public static final class PrivateMiddlewareServiceStub
      extends io.grpc.stub.AbstractAsyncStub<PrivateMiddlewareServiceStub> {
    private PrivateMiddlewareServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PrivateMiddlewareServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PrivateMiddlewareServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * This is used by a requires point to start a new channel with a requirement contract.
     * </pre>
     */
    public void registerChannel(ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getRegisterChannelMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * This is used by provider services to register their provision contract with the Registry/Broker.
     * </pre>
     */
    public void registerApp(ar.com.montepagano.search.v1.Middleware.RegisterAppRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.RegisterAppResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncServerStreamingCall(
          getChannel().newCall(getRegisterAppMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * This is used by local app (be it a Service Client or a Service Provider) to close a channel.
     * </pre>
     */
    public void closeChannel(ar.com.montepagano.search.v1.Middleware.CloseChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.CloseChannelResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCloseChannelMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * This is used by the local app to communicate with other participants in an already
     * initiated or registered channel
     * </pre>
     */
    public void appSend(ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.AppSendResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getAppSendMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void appRecv(ar.com.montepagano.search.v1.Middleware.AppRecvRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getAppRecvMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service PrivateMiddlewareService.
   * <pre>
   * 
   *This service is what a Middleware exposes to its local users (not on the internet.)
   * </pre>
   */
  public static final class PrivateMiddlewareServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<PrivateMiddlewareServiceBlockingStub> {
    private PrivateMiddlewareServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PrivateMiddlewareServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PrivateMiddlewareServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * This is used by a requires point to start a new channel with a requirement contract.
     * </pre>
     */
    public ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse registerChannel(ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getRegisterChannelMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * This is used by provider services to register their provision contract with the Registry/Broker.
     * </pre>
     */
    public java.util.Iterator<ar.com.montepagano.search.v1.Middleware.RegisterAppResponse> registerApp(
        ar.com.montepagano.search.v1.Middleware.RegisterAppRequest request) {
      return io.grpc.stub.ClientCalls.blockingServerStreamingCall(
          getChannel(), getRegisterAppMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * This is used by local app (be it a Service Client or a Service Provider) to close a channel.
     * </pre>
     */
    public ar.com.montepagano.search.v1.Middleware.CloseChannelResponse closeChannel(ar.com.montepagano.search.v1.Middleware.CloseChannelRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCloseChannelMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * This is used by the local app to communicate with other participants in an already
     * initiated or registered channel
     * </pre>
     */
    public ar.com.montepagano.search.v1.Middleware.AppSendResponse appSend(ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getAppSendMethod(), getCallOptions(), request);
    }

    /**
     */
    public ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse appRecv(ar.com.montepagano.search.v1.Middleware.AppRecvRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getAppRecvMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service PrivateMiddlewareService.
   * <pre>
   * 
   *This service is what a Middleware exposes to its local users (not on the internet.)
   * </pre>
   */
  public static final class PrivateMiddlewareServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<PrivateMiddlewareServiceFutureStub> {
    private PrivateMiddlewareServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PrivateMiddlewareServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PrivateMiddlewareServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * This is used by a requires point to start a new channel with a requirement contract.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse> registerChannel(
        ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getRegisterChannelMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * This is used by local app (be it a Service Client or a Service Provider) to close a channel.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<ar.com.montepagano.search.v1.Middleware.CloseChannelResponse> closeChannel(
        ar.com.montepagano.search.v1.Middleware.CloseChannelRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCloseChannelMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * This is used by the local app to communicate with other participants in an already
     * initiated or registered channel
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<ar.com.montepagano.search.v1.Middleware.AppSendResponse> appSend(
        ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getAppSendMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse> appRecv(
        ar.com.montepagano.search.v1.Middleware.AppRecvRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getAppRecvMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_REGISTER_CHANNEL = 0;
  private static final int METHODID_REGISTER_APP = 1;
  private static final int METHODID_CLOSE_CHANNEL = 2;
  private static final int METHODID_APP_SEND = 3;
  private static final int METHODID_APP_RECV = 4;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AsyncService serviceImpl;
    private final int methodId;

    MethodHandlers(AsyncService serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_REGISTER_CHANNEL:
          serviceImpl.registerChannel((ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest) request,
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse>) responseObserver);
          break;
        case METHODID_REGISTER_APP:
          serviceImpl.registerApp((ar.com.montepagano.search.v1.Middleware.RegisterAppRequest) request,
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.RegisterAppResponse>) responseObserver);
          break;
        case METHODID_CLOSE_CHANNEL:
          serviceImpl.closeChannel((ar.com.montepagano.search.v1.Middleware.CloseChannelRequest) request,
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.CloseChannelResponse>) responseObserver);
          break;
        case METHODID_APP_SEND:
          serviceImpl.appSend((ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest) request,
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.AppSendResponse>) responseObserver);
          break;
        case METHODID_APP_RECV:
          serviceImpl.appRecv((ar.com.montepagano.search.v1.Middleware.AppRecvRequest) request,
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  public static final io.grpc.ServerServiceDefinition bindService(AsyncService service) {
    return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
        .addMethod(
          getRegisterChannelMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest,
              ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse>(
                service, METHODID_REGISTER_CHANNEL)))
        .addMethod(
          getRegisterAppMethod(),
          io.grpc.stub.ServerCalls.asyncServerStreamingCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.Middleware.RegisterAppRequest,
              ar.com.montepagano.search.v1.Middleware.RegisterAppResponse>(
                service, METHODID_REGISTER_APP)))
        .addMethod(
          getCloseChannelMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.Middleware.CloseChannelRequest,
              ar.com.montepagano.search.v1.Middleware.CloseChannelResponse>(
                service, METHODID_CLOSE_CHANNEL)))
        .addMethod(
          getAppSendMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest,
              ar.com.montepagano.search.v1.Middleware.AppSendResponse>(
                service, METHODID_APP_SEND)))
        .addMethod(
          getAppRecvMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.Middleware.AppRecvRequest,
              ar.com.montepagano.search.v1.AppMessageOuterClass.AppRecvResponse>(
                service, METHODID_APP_RECV)))
        .build();
  }

  private static abstract class PrivateMiddlewareServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    PrivateMiddlewareServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return ar.com.montepagano.search.v1.Middleware.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("PrivateMiddlewareService");
    }
  }

  private static final class PrivateMiddlewareServiceFileDescriptorSupplier
      extends PrivateMiddlewareServiceBaseDescriptorSupplier {
    PrivateMiddlewareServiceFileDescriptorSupplier() {}
  }

  private static final class PrivateMiddlewareServiceMethodDescriptorSupplier
      extends PrivateMiddlewareServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    PrivateMiddlewareServiceMethodDescriptorSupplier(java.lang.String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (PrivateMiddlewareServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new PrivateMiddlewareServiceFileDescriptorSupplier())
              .addMethod(getRegisterChannelMethod())
              .addMethod(getRegisterAppMethod())
              .addMethod(getCloseChannelMethod())
              .addMethod(getAppSendMethod())
              .addMethod(getAppRecvMethod())
              .build();
        }
      }
    }
    return result;
  }
}
