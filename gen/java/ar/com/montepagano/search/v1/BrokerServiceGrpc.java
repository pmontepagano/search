package ar.com.montepagano.search.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.58.0)",
    comments = "Source: search/v1/broker.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class BrokerServiceGrpc {

  private BrokerServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "search.v1.BrokerService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Broker.BrokerChannelRequest,
      ar.com.montepagano.search.v1.Broker.BrokerChannelResponse> getBrokerChannelMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "BrokerChannel",
      requestType = ar.com.montepagano.search.v1.Broker.BrokerChannelRequest.class,
      responseType = ar.com.montepagano.search.v1.Broker.BrokerChannelResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Broker.BrokerChannelRequest,
      ar.com.montepagano.search.v1.Broker.BrokerChannelResponse> getBrokerChannelMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Broker.BrokerChannelRequest, ar.com.montepagano.search.v1.Broker.BrokerChannelResponse> getBrokerChannelMethod;
    if ((getBrokerChannelMethod = BrokerServiceGrpc.getBrokerChannelMethod) == null) {
      synchronized (BrokerServiceGrpc.class) {
        if ((getBrokerChannelMethod = BrokerServiceGrpc.getBrokerChannelMethod) == null) {
          BrokerServiceGrpc.getBrokerChannelMethod = getBrokerChannelMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.Broker.BrokerChannelRequest, ar.com.montepagano.search.v1.Broker.BrokerChannelResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "BrokerChannel"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Broker.BrokerChannelRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Broker.BrokerChannelResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BrokerServiceMethodDescriptorSupplier("BrokerChannel"))
              .build();
        }
      }
    }
    return getBrokerChannelMethod;
  }

  private static volatile io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Broker.RegisterProviderRequest,
      ar.com.montepagano.search.v1.Broker.RegisterProviderResponse> getRegisterProviderMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "RegisterProvider",
      requestType = ar.com.montepagano.search.v1.Broker.RegisterProviderRequest.class,
      responseType = ar.com.montepagano.search.v1.Broker.RegisterProviderResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Broker.RegisterProviderRequest,
      ar.com.montepagano.search.v1.Broker.RegisterProviderResponse> getRegisterProviderMethod() {
    io.grpc.MethodDescriptor<ar.com.montepagano.search.v1.Broker.RegisterProviderRequest, ar.com.montepagano.search.v1.Broker.RegisterProviderResponse> getRegisterProviderMethod;
    if ((getRegisterProviderMethod = BrokerServiceGrpc.getRegisterProviderMethod) == null) {
      synchronized (BrokerServiceGrpc.class) {
        if ((getRegisterProviderMethod = BrokerServiceGrpc.getRegisterProviderMethod) == null) {
          BrokerServiceGrpc.getRegisterProviderMethod = getRegisterProviderMethod =
              io.grpc.MethodDescriptor.<ar.com.montepagano.search.v1.Broker.RegisterProviderRequest, ar.com.montepagano.search.v1.Broker.RegisterProviderResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "RegisterProvider"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Broker.RegisterProviderRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  ar.com.montepagano.search.v1.Broker.RegisterProviderResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BrokerServiceMethodDescriptorSupplier("RegisterProvider"))
              .build();
        }
      }
    }
    return getRegisterProviderMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static BrokerServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BrokerServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BrokerServiceStub>() {
        @java.lang.Override
        public BrokerServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BrokerServiceStub(channel, callOptions);
        }
      };
    return BrokerServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static BrokerServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BrokerServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BrokerServiceBlockingStub>() {
        @java.lang.Override
        public BrokerServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BrokerServiceBlockingStub(channel, callOptions);
        }
      };
    return BrokerServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static BrokerServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BrokerServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BrokerServiceFutureStub>() {
        @java.lang.Override
        public BrokerServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BrokerServiceFutureStub(channel, callOptions);
        }
      };
    return BrokerServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public interface AsyncService {

    /**
     */
    default void brokerChannel(ar.com.montepagano.search.v1.Broker.BrokerChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Broker.BrokerChannelResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getBrokerChannelMethod(), responseObserver);
    }

    /**
     */
    default void registerProvider(ar.com.montepagano.search.v1.Broker.RegisterProviderRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Broker.RegisterProviderResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getRegisterProviderMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service BrokerService.
   */
  public static abstract class BrokerServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return BrokerServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service BrokerService.
   */
  public static final class BrokerServiceStub
      extends io.grpc.stub.AbstractAsyncStub<BrokerServiceStub> {
    private BrokerServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BrokerServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BrokerServiceStub(channel, callOptions);
    }

    /**
     */
    public void brokerChannel(ar.com.montepagano.search.v1.Broker.BrokerChannelRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Broker.BrokerChannelResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getBrokerChannelMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void registerProvider(ar.com.montepagano.search.v1.Broker.RegisterProviderRequest request,
        io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Broker.RegisterProviderResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getRegisterProviderMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service BrokerService.
   */
  public static final class BrokerServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<BrokerServiceBlockingStub> {
    private BrokerServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BrokerServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BrokerServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public ar.com.montepagano.search.v1.Broker.BrokerChannelResponse brokerChannel(ar.com.montepagano.search.v1.Broker.BrokerChannelRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getBrokerChannelMethod(), getCallOptions(), request);
    }

    /**
     */
    public ar.com.montepagano.search.v1.Broker.RegisterProviderResponse registerProvider(ar.com.montepagano.search.v1.Broker.RegisterProviderRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getRegisterProviderMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service BrokerService.
   */
  public static final class BrokerServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<BrokerServiceFutureStub> {
    private BrokerServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BrokerServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BrokerServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<ar.com.montepagano.search.v1.Broker.BrokerChannelResponse> brokerChannel(
        ar.com.montepagano.search.v1.Broker.BrokerChannelRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getBrokerChannelMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<ar.com.montepagano.search.v1.Broker.RegisterProviderResponse> registerProvider(
        ar.com.montepagano.search.v1.Broker.RegisterProviderRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getRegisterProviderMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_BROKER_CHANNEL = 0;
  private static final int METHODID_REGISTER_PROVIDER = 1;

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
        case METHODID_BROKER_CHANNEL:
          serviceImpl.brokerChannel((ar.com.montepagano.search.v1.Broker.BrokerChannelRequest) request,
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Broker.BrokerChannelResponse>) responseObserver);
          break;
        case METHODID_REGISTER_PROVIDER:
          serviceImpl.registerProvider((ar.com.montepagano.search.v1.Broker.RegisterProviderRequest) request,
              (io.grpc.stub.StreamObserver<ar.com.montepagano.search.v1.Broker.RegisterProviderResponse>) responseObserver);
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
          getBrokerChannelMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.Broker.BrokerChannelRequest,
              ar.com.montepagano.search.v1.Broker.BrokerChannelResponse>(
                service, METHODID_BROKER_CHANNEL)))
        .addMethod(
          getRegisterProviderMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              ar.com.montepagano.search.v1.Broker.RegisterProviderRequest,
              ar.com.montepagano.search.v1.Broker.RegisterProviderResponse>(
                service, METHODID_REGISTER_PROVIDER)))
        .build();
  }

  private static abstract class BrokerServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    BrokerServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return ar.com.montepagano.search.v1.Broker.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("BrokerService");
    }
  }

  private static final class BrokerServiceFileDescriptorSupplier
      extends BrokerServiceBaseDescriptorSupplier {
    BrokerServiceFileDescriptorSupplier() {}
  }

  private static final class BrokerServiceMethodDescriptorSupplier
      extends BrokerServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    BrokerServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (BrokerServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new BrokerServiceFileDescriptorSupplier())
              .addMethod(getBrokerChannelMethod())
              .addMethod(getRegisterProviderMethod())
              .build();
        }
      }
    }
    return result;
  }
}
