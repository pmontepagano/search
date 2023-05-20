package middleware

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"

	pb "github.com/clpombo/search/gen/go/search/v1"
	"github.com/clpombo/search/internal/contract"
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

	// local provider apps. key: appID, value: RegisterAppServer (connection is kept open
	// with each local app so as to notify new channels)
	registeredApps map[string]registeredApp

	// channels already brokered. key: ChannelID
	brokeredChannels map[string]*SEARCHChannel

	// channels registered by local apps that have not yet been brokered. key: LocalID
	unBrokeredChannels map[string]*SEARCHChannel

	// channels being brokered
	brokeringChannels map[string]*SEARCHChannel
	channelLock       *sync.Mutex

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
	s.channelLock = &sync.Mutex{}

	s.brokerAddr = brokerAddr
	s.brokerPort = brokerPort
	s.logger = log.New(os.Stderr, "[MIDDLEWARE] - ", log.LstdFlags|log.Lmsgprefix)
	return &s
}

// Struct to represent each Service Provider that is registered in the middleware.
// We notify to the local Service Provider that a new channel is being initiatied for it using the NotifyChan.
type registeredApp struct {
	Contract   contract.LocalContract
	NotifyChan *chan pb.InitChannelNotification
}

type messageExchangeStream interface {
	Send(*pb.ApplicationMessageWithHeaders) error
	Recv() (*pb.ApplicationMessageWithHeaders, error)
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

	outStreams map[string]messageExchangeStream // map participant names to outgoing streams
	inStreams  map[string]messageExchangeStream // map participant names to incoming streams

	// buffers for incoming/outgoing messages from/to each participant
	Outgoing map[string]chan *pb.AppMessage
	Incoming map[string]chan *pb.AppMessage

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
	r.outStreams = make(map[string]messageExchangeStream)
	r.inStreams = make(map[string]messageExchangeStream)

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

// invoked by local provider app with a provision contract
func (s *MiddlewareServer) RegisterApp(req *pb.RegisterAppRequest, stream pb.PrivateMiddlewareService_RegisterAppServer) error {
	client, conn := s.connectBroker()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.RegisterProvider(ctx, &pb.RegisterProviderRequest{
		Contract: req.ProviderContract,
		Url:      s.PublicURL,
	})
	if err != nil {
		s.logger.Fatalf("ERROR RegisterApp: %v", err)
	}
	ack := pb.RegisterAppResponse{
		AckOrNew: &pb.RegisterAppResponse_AppId{
			AppId: res.GetAppId()}}
	if err := stream.Send(&ack); err != nil {
		return err
	}

	notifyInitChannel := make(chan pb.InitChannelNotification, bufferSize)
	provContract, err := contract.ConvertPBLocalContract(req.ProviderContract)
	if err != nil {
		s.logger.Printf("Error parsing protobuf contract. %v", err)
		return err
	}
	s.registeredApps[res.GetAppId()] = registeredApp{
		Contract:   provContract,
		NotifyChan: &notifyInitChannel,
	}
	for {
		newChan := <-notifyInitChannel
		msg := pb.RegisterAppResponse{
			AckOrNew: &pb.RegisterAppResponse_Notification{
				Notification: &newChan}}
		if err := stream.Send(&msg); err != nil {
			return err
		}
		// TODO: in which case do I break the loop? We need to listen for disconnects from the provider app and/or add a msg to unregister.
	}
	return nil
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
func (r *SEARCHChannel) sender(participant string) {
	r.mw.logger.Printf("Started sender routine for channel %s, participant %s\n", r.ID, participant)
	// wait for first message
	msg := <-r.Outgoing[participant]
	// connect and save stream in r.outStreams
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())) // TODO: use tls
	provconn, err := grpc.Dial(r.addresses[participant].GetUrl(), opts...)
	if err != nil {
		r.mw.logger.Fatalf("fail to dial: %v", err)
	}
	defer provconn.Close()
	provClient := pb.NewPublicMiddlewareServiceClient(provconn)

	stream, err := provClient.MessageExchange(context.Background())
	if err != nil {
		r.mw.logger.Fatalf("failed to establish MessageExchange")
	}
	r.outStreams[participant] = stream
	for {
		// TODO: there has to be a way to stop this goroutine
		messageWithHeaders := pb.ApplicationMessageWithHeaders{
			SenderId:    r.AppID.String(),
			ChannelId:   r.ID.String(),
			RecipientId: r.addresses[participant].AppId,
			Content:     msg}
		if err := stream.Send(&messageWithHeaders); err != nil {
			r.mw.logger.Fatalf("failed to send message: %v", err)
		}
		msg = <-r.Outgoing[participant]
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

	// TODO: reply with error code in case there is an error. eg buffer full.
	return &pb.AppSendResponse{Result: pb.Result_OK}, nil
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

// When the middleware receives a message in its public interface, it must enqueue it so that
// the corresponding local app can receive it
func (s *MiddlewareServer) MessageExchange(stream pb.PublicMiddlewareService_MessageExchangeServer) error {
	s.logger.Print("Received MessageExchange...")
	in, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err // TODO: what should we do here?
	}
	s.channelLock.Lock()
	c, ok := s.brokeredChannels[in.GetChannelId()]
	if !ok {
		s.logger.Printf("Received MessageExchange with ChannelID %s but it is not a brokered Channel in this middleware.", in.GetChannelId())
	}
	// TODO: must check in.RecipientId... we could be hosting two different apps from same channel
	participantName := c.participants[in.SenderId]
	c.inStreams[participantName] = stream
	s.channelLock.Unlock()

	s.logger.Printf("Received message from %s: %s", in.SenderId, string(in.Content.Body))
	c.Incoming[participantName] <- in.Content

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		s.logger.Printf("Received message from %s: %s", in.SenderId, string(in.Content.Body))
		c.Incoming[participantName] <- in.Content

		// ack := pb.ApplicationMessageWithHeaders{
		// 	ChannelId:   in.ChannelId,
		// 	RecipientId: in.SenderId,
		// 	SenderId:    "provmwID-44",
		// 	Content:     &pb.MessageContent{Body: []byte("ack")}}

		// if err := stream.Send(&ack); err != nil {
		// 	return err
		// }
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
		*regapp.NotifyChan <- pb.InitChannelNotification{ChannelId: icr.GetChannelId()}
	} else {
		r, ok = s.brokeringChannels[icr.AppId]
		if !ok {
			s.logger.Panicf("InitChannel with inexistent app_id: %s", icr.GetAppId())
			// TODO: return error
		}
		s.logger.Printf("Received InitChannel for channel being brokered. %v", r.ID)
		delete(s.brokeringChannels, icr.AppId)
	}
	r.addresses = icr.GetParticipants()
	for k, v := range r.addresses {
		r.participants[v.AppId] = k
	}
	r.ID = uuid.MustParse(icr.GetChannelId())
	s.brokeredChannels[icr.GetChannelId()] = r
	s.localChannels.Insert(r.LocalID.String(), r.ID.String())

	return &pb.InitChannelResponse{Result: pb.InitChannelResponse_ACK}, nil
}

// StartChannel is fired by the broker and sent to all participants (initiator included) to signal that
// all participants are ready and communication can start.
func (s *MiddlewareServer) StartChannel(ctx context.Context, req *pb.StartChannelRequest) (*pb.StartChannelResponse, error) {
	s.logger.Printf("StartChannel. ChannelID: %s. AppID: %s", req.ChannelId, req.AppId)
	s.channelLock.Lock()
	c, ok := s.brokeredChannels[req.ChannelId]
	s.channelLock.Unlock()
	if !ok {
		s.logger.Panicf("Received StartChannel on non brokered ChannelID: %s", req.ChannelId)
	}
	for _, p := range c.Contract.GetRemoteParticipantNames() {
		go c.sender(p)
	}
	return &pb.StartChannelResponse{Result: pb.StartChannelResponse_ACK}, nil
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
