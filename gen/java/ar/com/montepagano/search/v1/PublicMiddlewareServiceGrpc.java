package ar.com.montepagano.search.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 *This service is what a Middleware exposes to external components (other participants and the broker)
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.62.2)",
    comments = "Source: search/v1/middleware.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class PublicMiddlewareServiceGrpc {

  private PublicMiddlewareServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "search.v1.PublicMiddlewareService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.InitChannelRequest,
      ar.com.montepagano.search.v1.Middleware.InitChannelResponse> getInitChannelMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "InitChannel",
      requestType = ar.com.montepagano.search.v1.Middleware.InitChannelRequest.class,
      responseType = ar.com.montepagano.search.v1.Middleware.InitChannelResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.InitChannelRequest,
      ar.com.montepagano.search.v1.Middleware.InitChannelResponse> getInitChannelMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.InitChannelRequest, ar.com.montepagano.search.v1.Middleware.InitChannelResponse> getInitChannelMethod;
    if ((getInitChannelMethod = PublicMiddlewareServiceGrpc.getInitChannelMethod) == null) {
      synchronized (PublicMiddlewareServiceGrpc.class) {
        if ((getInitChannelMethod = PublicMiddlewareServiceGrpc.getInitChannelMethod) == null) {
          PublicMiddlewareServiceGrpc.getInitChannelMethod = getInitChannelMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.Middleware.InitChannelRequest, ar.com.montepagano.search.v1.Middleware.InitChannelResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "InitChannel"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.InitChannelRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.InitChannelResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PublicMiddlewareServiceMethodDescriptorSupplier("InitChannel"))
              .build();
        }
      }
    }
    return getInitChannelMethod;
  }

  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.StartChannelRequest,
      ar.com.montepagano.search.v1.Middleware.StartChannelResponse> getStartChannelMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "StartChannel",
      requestType = ar.com.montepagano.search.v1.Middleware.StartChannelRequest.class,
      responseType = ar.com.montepagano.search.v1.Middleware.StartChannelResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.StartChannelRequest,
      ar.com.montepagano.search.v1.Middleware.StartChannelResponse> getStartChannelMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Middleware.StartChannelRequest, ar.com.montepagano.search.v1.Middleware.StartChannelResponse> getStartChannelMethod;
    if ((getStartChannelMethod = PublicMiddlewareServiceGrpc.getStartChannelMethod) == null) {
      synchronized (PublicMiddlewareServiceGrpc.class) {
        if ((getStartChannelMethod = PublicMiddlewareServiceGrpc.getStartChannelMethod) == null) {
          PublicMiddlewareServiceGrpc.getStartChannelMethod = getStartChannelMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.Middleware.StartChannelRequest, ar.com.montepagano.search.v1.Middleware.StartChannelResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "StartChannel"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.StartChannelRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.StartChannelResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PublicMiddlewareServiceMethodDescriptorSupplier("StartChannel"))
              .build();
        }
      }
    }
    return getStartChannelMethod;
  }

  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.AppMessageOuterClass.MessageExchangeRequest,
      ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse> getMessageExchangeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "MessageExchange",
      requestType = ar.com.montepagano.search.v1.AppMessageOuterClass.MessageExchangeRequest.class,
      responseType = ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.AppMessageOuterClass.MessageExchangeRequest,
      ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse> getMessageExchangeMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.AppMessageOuterClass.MessageExchangeRequest, ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse> getMessageExchangeMethod;
    if ((getMessageExchangeMethod = PublicMiddlewareServiceGrpc.getMessageExchangeMethod) == null) {
      synchronized (PublicMiddlewareServiceGrpc.class) {
        if ((getMessageExchangeMethod = PublicMiddlewareServiceGrpc.getMessageExchangeMethod) == null) {
          PublicMiddlewareServiceGrpc.getMessageExchangeMethod = getMessageExchangeMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.AppMessageOuterClass.MessageExchangeRequest, ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "MessageExchange"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.AppMessageOuterClass.MessageExchangeRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PublicMiddlewareServiceMethodDescriptorSupplier("MessageExchange"))
              .build();
        }
      }
    }
    return getMessageExchangeMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static PublicMiddlewareServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PublicMiddlewareServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PublicMiddlewareServiceStub>() {
        @java.lang.Override
        public PublicMiddlewareServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PublicMiddlewareServiceStub(channel, callOptions);
        }
      };
    return PublicMiddlewareServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static PublicMiddlewareServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PublicMiddlewareServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PublicMiddlewareServiceBlockingStub>() {
        @java.lang.Override
        public PublicMiddlewareServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PublicMiddlewareServiceBlockingStub(channel, callOptions);
        }
      };
    return PublicMiddlewareServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static PublicMiddlewareServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PublicMiddlewareServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PublicMiddlewareServiceFutureStub>() {
        @java.lang.Override
        public PublicMiddlewareServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PublicMiddlewareServiceFutureStub(channel, callOptions);
        }
      };
    return PublicMiddlewareServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   *This service is what a Middleware exposes to external components (other participants and the broker)
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * The Broker, when a new channel is registered, signals all providers with this
     * </pre>
     */
    default void initChannel(ar.com.montepagano.search.v1.Middleware.InitChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.InitChannelResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getInitChannelMethod(), responseObserver);
    }

    /**
     */
    default void startChannel(ar.com.montepagano.search.v1.Middleware.StartChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.StartChannelResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getStartChannelMethod(), responseObserver);
    }

    /**
     */
    default io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.AppMessageOuterClass.MessageExchangeRequest> messageExchange(
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse> responseObserver) {
      return io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall(getMessageExchangeMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service PublicMiddlewareService.
   * <pre>
   *This service is what a Middleware exposes to external components (other participants and the broker)
   * </pre>
   */
  public static abstract class PublicMiddlewareServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return PublicMiddlewareServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service PublicMiddlewareService.
   * <pre>
   *This service is what a Middleware exposes to external components (other participants and the broker)
   * </pre>
   */
  public static final class PublicMiddlewareServiceStub
      extends io.grpc.stub.AbstractAsyncStub<PublicMiddlewareServiceStub> {
    private PublicMiddlewareServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PublicMiddlewareServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PublicMiddlewareServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * The Broker, when a new channel is registered, signals all providers with this
     * </pre>
     */
    public void initChannel(ar.com.montepagano.search.v1.Middleware.InitChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.InitChannelResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getInitChannelMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void startChannel(ar.com.montepagano.search.v1.Middleware.StartChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.StartChannelResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getStartChannelMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.AppMessageOuterClass.MessageExchangeRequest> messageExchange(
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse> responseObserver) {
      return io.grpc.stub.ClientCalls.asyncClientStreamingCall(
          getChannel().newCall(getMessageExchangeMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service PublicMiddlewareService.
   * <pre>
   *This service is what a Middleware exposes to external components (other participants and the broker)
   * </pre>
   */
  public static final class PublicMiddlewareServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<PublicMiddlewareServiceBlockingStub> {
    private PublicMiddlewareServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PublicMiddlewareServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PublicMiddlewareServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * The Broker, when a new channel is registered, signals all providers with this
     * </pre>
     */
    public ar.com.montepagano.search.v1.Middleware.InitChannelResponse initChannel(ar.com.montepagano.search.v1.Middleware.InitChannelRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getInitChannelMethod(), getCallOptions(), request);
    }

    /**
     */
    public ar.com.montepagano.search.v1.Middleware.StartChannelResponse startChannel(ar.com.montepagano.search.v1.Middleware.StartChannelRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getStartChannelMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service PublicMiddlewareService.
   * <pre>
   *This service is what a Middleware exposes to external components (other participants and the broker)
   * </pre>
   */
  public static final class PublicMiddlewareServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<PublicMiddlewareServiceFutureStub> {
    private PublicMiddlewareServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PublicMiddlewareServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PublicMiddlewareServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * The Broker, when a new channel is registered, signals all providers with this
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<ar.com.montepagano.search.v1.Middleware.InitChannelResponse> initChannel(
        ar.com.montepagano.search.v1.Middleware.InitChannelRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getInitChannelMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<ar.com.montepagano.search.v1.Middleware.StartChannelResponse> startChannel(
        ar.com.montepagano.search.v1.Middleware.StartChannelRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getStartChannelMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_INIT_CHANNEL = 0;
  private static final int METHODID_START_CHANNEL = 1;
  private static final int METHODID_MESSAGE_EXCHANGE = 2;

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
        case METHODID_INIT_CHANNEL:
          serviceImpl.initChannel((ar.com.montepagano.search.v1.Middleware.InitChannelRequest) request,
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.InitChannelResponse>) responseObserver);
          break;
        case METHODID_START_CHANNEL:
          serviceImpl.startChannel((ar.com.montepagano.search.v1.Middleware.StartChannelRequest) request,
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.StartChannelResponse>) responseObserver);
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
        case METHODID_MESSAGE_EXCHANGE:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.messageExchange(
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  public static final io.grpc.ServerServiceDefinition bindService(AsyncService service) {
    return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
        .addMethod(
          getInitChannelMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.Middleware.InitChannelRequest,
              ar.com.montepagano.search.v1.Middleware.InitChannelResponse>(
                service, METHODID_INIT_CHANNEL)))
        .addMethod(
          getStartChannelMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.Middleware.StartChannelRequest,
              ar.com.montepagano.search.v1.Middleware.StartChannelResponse>(
                service, METHODID_START_CHANNEL)))
        .addMethod(
          getMessageExchangeMethod(),
          io.grpc.stub.ServerCalls.asyncClientStreamingCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.AppMessageOuterClass.MessageExchangeRequest,
              ar.com.montepagano.search.v1.Middleware.MessageExchangeResponse>(
                service, METHODID_MESSAGE_EXCHANGE)))
        .build();
  }

  private static abstract class PublicMiddlewareServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    PublicMiddlewareServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return ar.com.montepagano.search.v1.Middleware.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("PublicMiddlewareService");
    }
  }

  private static final class PublicMiddlewareServiceFileDescriptorSupplier
      extends PublicMiddlewareServiceBaseDescriptorSupplier {
    PublicMiddlewareServiceFileDescriptorSupplier() {}
  }

  private static final class PublicMiddlewareServiceMethodDescriptorSupplier
      extends PublicMiddlewareServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    PublicMiddlewareServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (PublicMiddlewareServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new PublicMiddlewareServiceFileDescriptorSupplier())
              .addMethod(getInitChannelMethod())
              .addMethod(getStartChannelMethod())
              .addMethod(getMessageExchangeMethod())
              .build();
        }
      }
    }
    return result;
  }
}
