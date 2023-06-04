package middleware

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/sourcegraph/conc/pool"

	"github.com/pmontepagano/search/contract"
	pb "github.com/pmontepagano/search/gen/go/search/v1"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/testdata"

	"github.com/vishalkuo/bimap"
)

var (
	bufferSize = 100
)

type MiddlewareServer struct {
	pb.UnimplementedPublicMiddlewareServiceServer
	pb.UnimplementedPrivateMiddlewareServiceServer

	// Local Service Providers. key: appID, value: RegisterAppServer (connection is kept open
	// with each local app so as to notify new channels)
	registeredApps map[string]registeredApp
	providersLock  *sync.Mutex

	// channels already brokered. key: ChannelID
	brokeredChannels map[string]*SEARCHChannel

	// channels registered by local apps that have not yet been brokered. key: LocalID
	unBrokeredChannels map[string]*SEARCHChannel

	// channels being brokered
	brokeringChannels map[string]*SEARCHChannel
	channelLock       *sync.RWMutex

	// mapping of channels' LocalID <--> ID (global)
	localChannels *bimap.BiMap

	brokerAddr string
	brokerPort int

	// pointers to gRPC servers
	publicServer  *grpc.Server
	privateServer *grpc.Server
	PublicURL     string
	PrivateURL    string

	// logger with prefix
	logger *log.Logger
}

// NewMiddlewareServer is MiddlewareServer's constructor
func NewMiddlewareServer(brokerAddr string, brokerPort int) *MiddlewareServer {
	var s MiddlewareServer
	s.localChannels = bimap.NewBiMap() // mapping between local channelID and global channelID. When initiator not local, they are equal
	s.registeredApps = make(map[string]registeredApp)

	s.brokeredChannels = make(map[string]*SEARCHChannel)   // channels already brokered (locally initiated or not)
	s.unBrokeredChannels = make(map[string]*SEARCHChannel) // channels locally registered but not yet brokered
	s.brokeringChannels = make(map[string]*SEARCHChannel)  // channels being brokered
	s.channelLock = new(sync.RWMutex)
	s.providersLock = new(sync.Mutex)

	s.brokerAddr = brokerAddr
	s.brokerPort = brokerPort
	s.logger = log.New(os.Stderr, "[MIDDLEWARE] - ", log.LstdFlags|log.Lmsgprefix)
	return &s
}

// Struct to represent each Service Provider that is registered in the middleware.
// We notify to the local Service Provider that a new channel is being initiatied for it using the NotifyChan.
type registeredApp struct {
	Contract   contract.LocalContract
	NotifyChan *chan initChannelNotification // Each message in the channel is a Channel ID (extracted from pb.InitChannelNotification).
}

// We use this struct internally to avoid copying the protobuf structure with its lock.
type initChannelNotification struct {
	ChannelID string
}

// SEARCHChannel represents what we call a "channel" in SEARCH
type SEARCHChannel struct {
	LocalID    uuid.UUID // the identifier the local app uses to identify the channel
	ID         uuid.UUID // channel identifier assigned by the broker
	Contract   contract.Contract
	ContractPB *pb.GlobalContract // We only save this when this middleware is representing the Service Client.
	AppID      uuid.UUID

	// localParticipant string                           // Name of the local participant in the contract.
	addresses    map[string]*pb.RemoteParticipant // map participant names to remote URLs and AppIDs, indexed by participant name
	participants map[string]string                // participant names indexed by AppID

	// buffers for incoming/outgoing messages from/to each participant
	Outgoing map[string]chan *pb.AppMessage
	Incoming map[string]chan *pb.AppMessage

	sendersPool           *pool.ContextPool  // pool of workers for sending outbound messages.
	sendersPoolCancelFunc context.CancelFunc // cancel function for the pool of workers that send outbound messages.

	// pointer to middleware
	mw *MiddlewareServer
}

// newSEARCHChannel returns a pointer to a new channel created for the MiddlewareServer.
// localParticipant is the name of the local app in the contract.
func (mw *MiddlewareServer) newSEARCHChannel(c contract.Contract, pbContract *pb.GlobalContract) (*SEARCHChannel, error) {
	// TODO: review ContractPB. Instead of keeping a reference to the protobuf, maybe keep only a copy
	// of the contract bytes?
	var r SEARCHChannel
	r.mw = mw
	r.LocalID = uuid.New()
	r.Contract = c
	r.ContractPB = pbContract
	r.addresses = make(map[string]*pb.RemoteParticipant)
	r.participants = make(map[string]string)

	r.Outgoing = make(map[string]chan *pb.AppMessage)
	r.Incoming = make(map[string]chan *pb.AppMessage)
	for _, p := range c.GetRemoteParticipantNames() {
		r.Outgoing[p] = make(chan *pb.AppMessage, bufferSize)
		r.Incoming[p] = make(chan *pb.AppMessage, bufferSize)
	}

	return &r, nil
}

func (s *MiddlewareServer) connectBroker() (pb.BrokerServiceClient, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", s.brokerAddr, s.brokerPort), opts...)
	if err != nil {
		s.logger.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewBrokerServiceClient(conn)

	return client, conn
}

// connect to the broker, send contract, wait for result and save data in the channel
func (r *SEARCHChannel) broker() {
	r.mw.logger.Printf("Requesting brokerage of contract: '%v'", r.Contract)
	client, conn := r.mw.connectBroker()
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.TODO()) // TODO: add deadline/timeout or inherit from request context.
	defer cancel()
	brokerresult, err := client.BrokerChannel(ctx, &pb.BrokerChannelRequest{
		Contract: r.ContractPB,
		PresetParticipants: map[string]*pb.RemoteParticipant{
			r.ContractPB.InitiatorName: {
				Url:   r.mw.PublicURL,
				AppId: r.LocalID.String(), // we use channels LocalID as AppID for initiator apps
			},
		},
	})
	if err != nil {
		r.mw.logger.Fatalf("%v.BrokerChannel(_) = _, %v: ", client, err)
	}
	if brokerresult.Result != pb.BrokerChannelResponse_RESULT_ACK {
		r.mw.logger.Fatalf("Non ACK return code when trying to broker channel.")
	}
}

// RegisterApp is invoked by Service Providers on the private interface in order to register their service with
// the broker. This function keeps the connection open between the middleware and the local Service Provider
// to notify it of new channels being initiated for it.
func (s *MiddlewareServer) RegisterApp(req *pb.RegisterAppRequest, stream pb.PrivateMiddlewareService_RegisterAppServer) error {
	// Connect to broker.
	client, conn := s.connectBroker()
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We lock here because we need this blocked in case the broker sends an InitChannel
	// before we have saved the local app in our internal structures.
	s.providersLock.Lock()

	// Send contract to broker.
	res, err := client.RegisterProvider(ctx, &pb.RegisterProviderRequest{
		Contract: req.ProviderContract,
		Url:      s.PublicURL,
	})
	if err != nil {
		s.logger.Fatalf("ERROR RegisterApp: %v", err)
	}
	// Send ACK to local app via the stream.
	ack := pb.RegisterAppResponse{
		AckOrNew: &pb.RegisterAppResponse_AppId{
			AppId: res.GetAppId()}}
	if err := stream.Send(&ack); err != nil {
		return err
	}

	// Save local app in middleware's internal structures.
	notifyInitChannel := make(chan initChannelNotification, bufferSize)
	provContract, err := contract.ConvertPBLocalContract(req.ProviderContract)
	if err != nil {
		s.logger.Printf("Error parsing protobuf contract. %v", err)
		return err
	}
	s.registeredApps[res.GetAppId()] = registeredApp{
		Contract:   provContract,
		NotifyChan: &notifyInitChannel,
	}
	s.providersLock.Unlock()

	for {
		newChan := <-notifyInitChannel
		msg := pb.RegisterAppResponse{
			AckOrNew: &pb.RegisterAppResponse_Notification{
				Notification: &pb.InitChannelNotification{
					ChannelId: newChan.ChannelID,
				},
			},
		}
		if err := stream.Send(&msg); err != nil {
			return err
		}
		// TODO: in which case do I break the loop? We need to listen for disconnects from the provider app and/or add a msg to unregister.
	}
}

// Invoked by local Service Client with a requirements contract.
func (s *MiddlewareServer) RegisterChannel(ctx context.Context, in *pb.RegisterChannelRequest) (*pb.RegisterChannelResponse, error) {
	// TODO: create a monitor for this channel.
	requirementsContract, err := contract.ConvertPBGlobalContract(in.GetRequirementsContract())
	if err != nil {
		s.logger.Printf("Error parsing protobuf contract. %v", err)
		return nil, err
	}
	c, err := s.newSEARCHChannel(requirementsContract, in.GetRequirementsContract())
	if err != nil {
		// TODO: newSEARCHChannel may fail because of other kinds of errors in the future.
		err := status.Error(codes.InvalidArgument, "invalid contract")
		return nil, err
	}
	c.AppID = c.LocalID
	s.channelLock.Lock()
	s.unBrokeredChannels[c.LocalID.String()] = c // this has to be changed when brokering
	s.channelLock.Unlock()
	return &pb.RegisterChannelResponse{ChannelId: c.LocalID.String()}, nil
}

// routine that actually sends messages to remote particpant on a channel
func (r *SEARCHChannel) sender(ctx context.Context, participant string) error {
	r.mw.logger.Printf("Started sender routine for channel %s, participant %s\n", r.ID, participant)
	// wait for first message
	var msg *pb.AppMessage
	select {
	case <-ctx.Done():
		r.mw.logger.Printf("Context cancelled (cause: %s). Closing sender routine for channel %s, participant %s. No connection established to remote Service Provider.", context.Cause(ctx), r.ID, participant)
		return nil
	case msg = <-r.Outgoing[participant]:
		r.mw.logger.Printf("Received first outbound message to send on channel %s for participant %s. Opening connection to remote Service Provider...", r.ID, participant)
		break
	}
	// connect and save stream in r.outStreams
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())) // TODO: use tls
	provconn, err := grpc.Dial(r.addresses[participant].GetUrl(), opts...)
	if err != nil {
		return fmt.Errorf("fail to dial: %v", err)
	}
	defer provconn.Close()
	provClient := pb.NewPublicMiddlewareServiceClient(provconn)

	provClientctx, provClientCancel := context.WithCancel(context.Background())
	defer provClientCancel()
	stream, err := provClient.MessageExchange(provClientctx)
	if err != nil {
		return fmt.Errorf("failed to establish MessageExchange")
	}
	r.mw.logger.Printf("Established connection to remote Service Provider for channel %s, participant %s\n", r.ID, participant)
	for {
		messageWithHeaders := pb.MessageExchangeRequest{
			SenderId:    r.AppID.String(),
			ChannelId:   r.ID.String(),
			RecipientId: r.addresses[participant].AppId,
			Content:     msg}
		if err := stream.Send(&messageWithHeaders); err != nil {
			return fmt.Errorf("failed to send message: %v", err)
		}
		r.mw.logger.Printf("Sent message to remote Service Provider for channel %s, participant %s\n", r.ID, participant)
		select {
		case msg = <-r.Outgoing[participant]:
			r.mw.logger.Printf("Received outbound message to send on channel %s for participant %s\n", r.ID, participant)
			break
		case <-ctx.Done():
			r.mw.logger.Printf("Context cancelled (cause: %s). Closing sender routine for channel %s, participant %s\n", context.Cause(ctx), r.ID, participant)
			close(r.Outgoing[participant])
			reply, err := stream.CloseAndRecv()
			if err != nil {
				r.mw.logger.Printf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
				return err
			}
			if reply.Result != pb.MessageExchangeResponse_RESULT_OK {
				r.mw.logger.Printf("Error in MessageExchange when attempting to close stream: %s", reply.ErrorMessage)
				return fmt.Errorf("Error in MessageExchange when attempting to close stream: %s", reply.ErrorMessage)
			}
			r.mw.logger.Printf("Closed sender stream for channel %s, participant %s\n", r.ID, participant)
			return nil
		}
	}
}

// if the SEARCHChannel is not yet brokered, we launch goroutine to do that
// and also update middleware's internal structures to reflect that change
func (s *MiddlewareServer) getChannelForUsage(localID string) *SEARCHChannel {
	s.channelLock.Lock()
	defer s.channelLock.Unlock()
	c, ok := s.unBrokeredChannels[localID]
	if ok {
		// channel has not been brokered
		s.brokeringChannels[localID] = c
		delete(s.unBrokeredChannels, localID)
		go c.broker()
	} else {
		// check if channel is being brokered
		c, ok = s.brokeringChannels[localID]
		if !ok {
			// channel must already be brokered
			channelID, ok := s.localChannels.Get(localID)
			if !ok {
				s.logger.Fatalf("getChannelForUsage invoked on channel ID %s: there's no localChannel with that ID.", localID)
			}
			c = s.brokeredChannels[channelID.(string)]
		}
	}
	return c
}

// simple (g)rpc the local app uses when sending a message to a remote participant on an already registered channel
func (s *MiddlewareServer) AppSend(ctx context.Context, req *pb.AppSendRequest) (*pb.AppSendResponse, error) {
	c := s.getChannelForUsage(req.ChannelId)
	c.Outgoing[req.Recipient] <- req.GetMessage() // enqueue message in outgoing buffer
	s.logger.Printf("Enqueued message of type %s in outgoing buffer for channel %s, participant %s\n", req.Message.GetType(), req.ChannelId, req.Recipient)
	// TODO: reply with error code in case there is an error. eg buffer full.
	return &pb.AppSendResponse{Result: pb.AppSendResponse_RESULT_OK}, nil
}

func (s *MiddlewareServer) AppRecv(ctx context.Context, req *pb.AppRecvRequest) (*pb.AppRecvResponse, error) {
	c := s.getChannelForUsage(req.ChannelId)
	msg := <-c.Incoming[req.Participant]

	return &pb.AppRecvResponse{
		ChannelId: req.ChannelId,
		Sender:    req.Participant,
		Message:   msg,
	}, nil
}

func (s *MiddlewareServer) CloseChannel(ctx context.Context, req *pb.CloseChannelRequest) (*pb.CloseChannelResponse, error) {
	// We have to check if the channel has messages in the buffers.
	s.logger.Printf("CloseChannel. ChannelID: %s", req.ChannelId)
	s.channelLock.RLock()

	c, ok := s.brokeredChannels[req.ChannelId]
	if !ok {
		s.channelLock.RUnlock()
		return nil, fmt.Errorf("CloseChannel invoked on an inexistent or unbrokered channel ID %s", req.ChannelId)
	}
	s.channelLock.RUnlock()

	// Check that there are no inbound messages in the buffers of the channel.
	participantsWithIncoming := make([]string, 0)
	for participant, chMsg := range c.Incoming {
		if len(chMsg) > 0 {
			participantsWithIncoming = append(participantsWithIncoming, participant)
		} else {
			// close channel
			close(chMsg)
			delete(c.Incoming, participant)
		}
	}
	if len(participantsWithIncoming) > 0 {
		return &pb.CloseChannelResponse{
			Result:                         pb.CloseChannelResponse_RESULT_PENDING_INBOUND,
			ErrorMessage:                   "There are still messages inbound messages to receive!",
			ParticipantsWithPendingInbound: participantsWithIncoming,
		}, nil
	}

	// We cancel the sender goroutines.
	// TODO: it would be nice if we could wait (with a deadline) for them to finish?
	if c.sendersPoolCancelFunc != nil {
		c.sendersPoolCancelFunc()
	}
	if c.sendersPool != nil {
		err := c.sendersPool.Wait()
		if err != nil {
			return &pb.CloseChannelResponse{
				// TODO: is this the right error code? Does it make sense to indicate that we have pending outbound?
				Result:       pb.CloseChannelResponse_RESULT_PENDING_OUTBOUND,
				ErrorMessage: err.Error(),
			}, err
		}
	}

	return &pb.CloseChannelResponse{
		Result: pb.CloseChannelResponse_RESULT_CLOSED,
	}, nil
}

// When the middleware receives a message in its public interface, it must enqueue it so that
// the corresponding local app can receive it
func (s *MiddlewareServer) MessageExchange(stream pb.PublicMiddlewareService_MessageExchangeServer) error {
	s.logger.Print("Received MessageExchange...")
	// Acá se está colgando cuando se cierra el stream desde el otro lado.
	in, err := stream.Recv()
	if err == io.EOF {
		s.logger.Print("Received EOF in MessageExchange without receiving any message...")
		return err
	}
	if err != nil {
		s.logger.Printf("Error in MessageExchange when attempting to recv from stream: %s", err)
		return err // TODO: what should we do here?
	}
	s.logger.Print("Attempting to obtain channelLock...")
	s.channelLock.RLock()
	s.logger.Print("Obtained the channelLock...")
	c, ok := s.brokeredChannels[in.GetChannelId()]
	if !ok {
		s.logger.Printf("Received MessageExchange with ChannelID %s but it is not a brokered Channel in this middleware.", in.GetChannelId())
		s.channelLock.RUnlock()
		s.logger.Print("Released channelLock...")
		return fmt.Errorf("Received MessageExchange with ChannelID %s but it is not a brokered Channel in this middleware.", in.GetChannelId())
	}
	s.channelLock.RUnlock()
	s.logger.Print("Released channelLock...")
	// TODO: must check in.RecipientId... we could be hosting two different apps from same channel
	participantName := c.participants[in.SenderId]

	s.logger.Printf("Received message from %s of type %s", in.SenderId, in.Content.Type)
	c.Incoming[participantName] <- in.Content

	for {
		s.logger.Printf("Waiting for next message from %s...", participantName)
		in, err := stream.Recv()
		if err == io.EOF {
			s.logger.Print("Received EOF in MessageExchange...")
			err = stream.SendAndClose(&pb.MessageExchangeResponse{Result: pb.MessageExchangeResponse_RESULT_OK})
			if err != nil {
				s.logger.Printf("Error when closing MessageExchange after receiving EOF: %s", err)
				return err
			}
			s.logger.Print("Successful close of MessageExchange after receiving EOF. Exiting MessageExchange...")
			return nil
		}
		if err != nil {
			s.logger.Printf("Error attempting to receive msg from %s in MessageExchange: %s", participantName, err)
			return err
		}

		s.logger.Printf("Received message from %s: %s", in.SenderId, string(in.Content.Body))
		c.Incoming[participantName] <- in.Content
	}
}

// rpc invoked by broker when initializing channel
// There are two options: a) the broker has matched one of our registered (provider) apps with a requirements contract
// or b) the broker is responding to a brokerage request we sent to it. If this is the case, then AppID should
// match the LocalID we generated for that channel of which we requested brokerage.
func (s *MiddlewareServer) InitChannel(ctx context.Context, icr *pb.InitChannelRequest) (*pb.InitChannelResponse, error) {
	s.logger.Printf("Received InitChannel. ChannelID: %s. AppID: %s", icr.ChannelId, icr.AppId)
	s.logger.Printf("InitChannel mapping received:%v", icr.GetParticipants())
	var r *SEARCHChannel
	s.channelLock.Lock()
	defer s.channelLock.Unlock()
	// InitChannelRequest: app_id, channel_id, participants (map[string]RemoteParticipant)
	s.providersLock.Lock()
	defer s.providersLock.Unlock()
	if regapp, ok := s.registeredApps[icr.GetAppId()]; ok {
		// create registered channel with channel_id
		var err error
		r, err = s.newSEARCHChannel(regapp.Contract, nil)
		if err != nil {
			// TODO: return proper gRPC error
			return nil, err
		}
		r.AppID = uuid.MustParse(icr.GetAppId())

		r.LocalID = uuid.MustParse(icr.GetChannelId()) // in non locally initiated channels, LocalID == ID

		// notify local provider app
		*regapp.NotifyChan <- initChannelNotification{ChannelID: icr.GetChannelId()}
	} else {
		r, ok = s.brokeringChannels[icr.AppId]
		if !ok {
			s.logger.Panicf("InitChannel with inexistent app_id: %s", icr.GetAppId())
			// TODO: return error
		}
		s.logger.Printf("Received InitChannel for channel being brokered. %v", r.LocalID)
		delete(s.brokeringChannels, icr.AppId)
	}
	r.addresses = icr.GetParticipants()
	for k, v := range r.addresses {
		r.participants[v.AppId] = k
	}
	r.ID = uuid.MustParse(icr.GetChannelId())
	s.brokeredChannels[icr.GetChannelId()] = r
	s.localChannels.Insert(r.LocalID.String(), r.ID.String())

	return &pb.InitChannelResponse{Result: pb.InitChannelResponse_RESULT_ACK}, nil
}

// StartChannel is fired by the broker and sent to all participants (initiator included) to signal that
// all participants are ready and communication can start.
func (s *MiddlewareServer) StartChannel(ctx context.Context, req *pb.StartChannelRequest) (*pb.StartChannelResponse, error) {
	s.logger.Printf("StartChannel. ChannelID: %s. AppID: %s", req.ChannelId, req.AppId)
	s.channelLock.RLock()
	c, ok := s.brokeredChannels[req.ChannelId]
	s.channelLock.RUnlock()
	if !ok {
		// TODO: change this and return error.
		s.logger.Panicf("Received StartChannel on non brokered ChannelID: %s", req.ChannelId)
	}
	sendersCtx, cancel := context.WithCancel(context.Background())
	c.sendersPoolCancelFunc = cancel
	c.sendersPool = pool.New().WithContext(sendersCtx).WithCancelOnError()
	for _, p := range c.Contract.GetRemoteParticipantNames() {
		thisP := p
		c.sendersPool.Go(func(cont context.Context) error {
			return c.sender(cont, thisP)
		})
	}
	return &pb.StartChannelResponse{Result: pb.StartChannelResponse_RESULT_ACK}, nil
}

// StartServer starts gRPC middleware server
func (s *MiddlewareServer) StartMiddlewareServer(wg *sync.WaitGroup, publicHost string, publicPort int, privateHost string, privatePort int, tls bool, certFile string, keyFile string) {
	s.PublicURL = fmt.Sprintf("%s:%d", publicHost, publicPort)
	s.logger = log.New(os.Stderr, fmt.Sprintf("[MIDDLEWARE] %s - ", s.PublicURL), log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	lisPub, err := net.Listen("tcp", s.PublicURL)
	if err != nil {
		s.logger.Fatalf("failed to listen: %v", err)
	}
	s.PrivateURL = fmt.Sprintf("%s:%d", privateHost, privatePort)
	lisPriv, err := net.Listen("tcp", s.PrivateURL)
	if err != nil {
		s.logger.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if tls {
		if certFile == "" {
			certFile = testdata.Path("server1.pem")
		}
		if keyFile == "" {
			keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			s.logger.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	publicGrpcServer := grpc.NewServer(opts...)
	privateGrpcServer := grpc.NewServer(opts...)

	s.publicServer = publicGrpcServer
	s.privateServer = privateGrpcServer

	pb.RegisterPublicMiddlewareServiceServer(publicGrpcServer, s)
	pb.RegisterPrivateMiddlewareServiceServer(privateGrpcServer, s)

	wg.Add(2)

	go func() {
		defer wg.Done()
		// s.logger.Println("Waiting for SIGINT or SIGTERM on public server.")
		if err := publicGrpcServer.Serve(lisPub); err != nil {
			s.logger.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		// s.logger.Println("Waiting for SIGINT or SIGTERM on private server.")
		if err := privateGrpcServer.Serve(lisPriv); err != nil {
			s.logger.Fatalf("failed to serve: %v", err)
		}
	}()

}

func (s *MiddlewareServer) Stop() {
	// TODO: flush all buffers before stopping!
	if s.publicServer != nil {
		s.publicServer.Stop()
	}
	if s.privateServer != nil {
		s.privateServer.Stop()
	}
}
