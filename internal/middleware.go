package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/google/uuid"

	pb "dc.uba.ar/this/search/protobuf"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	"github.com/vishalkuo/bimap"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	caFile     = flag.String("ca_file", "", "The file containing the CA root cert file")
	port       = flag.Int("port", 10000, "The server port")
	brokerAddr = flag.String("broker_addr", "localhost", "The server address in the format of host:port")
	brokerPort = flag.Int("broker_port", 10000, "The port in which the broker is listening")
	bufferSize = flag.Int("buffer_size", 100, "The size of each connection's  buffer")
)

type middlewareServer struct {
	pb.UnimplementedPublicMiddlewareServer
	pb.UnimplementedPrivateMiddlewareServer

	// local provider apps. key: appID, value: RegisterAppServer (connection is kept open
	// with each local app so as to notify new channels)
	registeredApps map[string]registeredApp

	// channels already brokered. key: ChannelID
	brokeredChannels map[string]*SEARCHChannel

	// channels registered by local apps that have not yet been used. key: LocalID
	unBrokeredChannels map[string]*SEARCHChannel

	// mapping of channels' LocalID <--> ID (global)
	localChannels *bimap.BiMap
}

func newMiddlewareServer() *middlewareServer {
	var s middlewareServer
	s.localChannels = bimap.NewBiMap() // mapping between local channelID and global channelID. When initiator not local, they are equal
	s.registeredApps = make(map[string]registeredApp)
	s.brokeredChannels = make(map[string]*SEARCHChannel)   // channels already brokered (locally initiated or not)
	s.unBrokeredChannels = make(map[string]*SEARCHChannel) // channels locally registered but not yet brokered

	return &s
}

type registeredApp struct {
	Contract   pb.Contract
	NotifyChan *chan pb.InitChannelNotification
}

type SEARCHChannel struct {
	LocalID  uuid.UUID // the identifier the local app uses to identify the channel
	ID       uuid.UUID // channel identifier assigned by the broker
	Contract pb.Contract

	addresses map[string]*pb.RemoteParticipant                      // map participant names to remote URLs and AppIDs
	streams   map[string]*pb.PublicMiddleware_MessageExchangeClient // map participant names to streams

	// buffers for incoming/outgoing messages from/to each participant
	Outgoing map[string]chan pb.MessageContent
	Incoming map[string]chan pb.MessageContent
}

// TODO: maybe refactor and make this function a middlewareServer method?
func newSEARCHChannel(contract pb.Contract) *SEARCHChannel {
	var r SEARCHChannel
	r.LocalID = uuid.New()
	r.Contract = contract
	r.addresses = make(map[string]*pb.RemoteParticipant)
	r.streams = make(map[string]*pb.PublicMiddleware_MessageExchangeClient)

	r.Outgoing = make(map[string]chan pb.MessageContent)
	r.Incoming = make(map[string]chan pb.MessageContent)
	for _, p := range contract.GetRemoteParticipants() {
		r.Outgoing[p] = make(chan pb.MessageContent, *bufferSize)
		r.Incoming[p] = make(chan pb.MessageContent, *bufferSize)
	}

	return &r
}

// connect to the broker, send contract, wait for result and save data in the channel
func (r *SEARCHChannel) broker(mw *middlewareServer) {
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		// creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		// if err != nil {
		// 	log.Fatalf("Failed to create TLS credentials %v", err)
		//}
		//opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *brokerAddr, *brokerPort), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBrokerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	brokerresult, err := client.BrokerChannel(ctx, &pb.BrokerChannelRequest{Contract: &r.Contract})
	if err != nil {
		log.Fatalf("%v.BrokerChannel(_) = _, %v: ", client, err)
	}
	if brokerresult.Result != pb.BrokerChannelResponse_ACK {
		log.Fatalf("Non ACK return code when trying to broker channel.")
	}

	// TODO: refactor this part into new function (also used on InitChannel)
	// r.addresses = brokerresult.GetParticipants()
	// r.ID = uuid.MustParse(brokerresult.GetChannelId())

	// TODO: use mutex to handle maps
	// mw.brokeredChannels[r.ID.String()] = r
	// delete(mw.unBrokeredChannels, r.LocalID.String())
	// mw.localChannels.Insert(r.LocalID.String(), r.ID.String())

	// for _, p := range r.Contract.GetRemoteParticipants() {
	// 	go r.sender(p)
	// }
}

// invoked by local provider app with a provision contract
func (s *middlewareServer) RegisterApp(req *pb.RegisterAppRequest, stream pb.PrivateMiddleware_RegisterAppServer) error {
	// TODO: talk to the Registry to get my app_id
	mockUUID := uuid.New() // TODO: this should be generated by the Registry
	ack := pb.RegisterAppResponse{
		AckOrNew: &pb.RegisterAppResponse_AppId{
			AppId: mockUUID.String()}}
	if err := stream.Send(&ack); err != nil {
		return err
	}

	notifyInitChannel := make(chan pb.InitChannelNotification, *bufferSize)
	s.registeredApps[mockUUID.String()] = registeredApp{Contract: *req.ProviderContract, NotifyChan: &notifyInitChannel}
	for {
		newChan := <-notifyInitChannel
		msg := pb.RegisterAppResponse{
			AckOrNew: &pb.RegisterAppResponse_Notification{
				Notification: &newChan}}
		if err := stream.Send(&msg); err != nil {
			return err
		}
		// TODO: in which case do I break the loop?
	}
	return nil
}

// invoked by local initiator app with a requirements contract
func (s *middlewareServer) RegisterChannel(ctx context.Context, in *pb.RegisterChannelRequest) (*pb.RegisterChannelResponse, error) {
	c := newSEARCHChannel(*in.GetRequirementsContract())
	s.unBrokeredChannels[c.LocalID.String()] = c // this has to be changed when brokering
	return &pb.RegisterChannelResponse{ChannelId: c.LocalID.String()}, nil
}

// routine that actually sends messages to remote particpant on a channel
func (r *SEARCHChannel) sender(participant string) {
	log.Println("Started sender routine for channel %s, participant %s", r.ID, participant)
	// TODO: there has to be a way to stop this goroutine
	for {
		// wait for first message
		msg := <-r.Outgoing[participant]
		// check if connected to remote participant
		stream, ok := r.streams[participant]
		if !ok {
			// if not connected, connect and save stream in r.streams
			var opts []grpc.DialOption
			opts = append(opts, grpc.WithBlock())
			provconn, err := grpc.Dial(fmt.Sprintf("%s:10000", r.addresses[participant].GetUrl()), opts...)
			if err != nil {
				log.Fatalf("fail to dial: %v", err)
			}
			defer provconn.Close()
			provClient := pb.NewPublicMiddlewareClient(provconn)

			stream, err := provClient.MessageExchange(context.Background())
			if err != nil {
				log.Fatalf("failed to establish MessageExchange")
			}
			r.streams[participant] = &stream
		}
		messageWithHeaders := pb.ApplicationMessageWithHeaders{
			SenderId:    "initiator", // TODO: what should we use here? participant name in contract is enough
			ChannelId:   r.ID.String(),
			RecipientId: r.addresses[participant].AppId,
			Content:     &msg}
		s := *stream
		s.Send(&messageWithHeaders)
	}
}

// simple (g)rpc the local app uses when sending a message to a remote participant on an already registered channel
func (s *middlewareServer) AppSend(ctx context.Context, req *pb.ApplicationMessageOut) (*pb.AppSendResponse, error) {
	localID := req.ChannelId // when local app is initiator, localID != channelID (global, set by broker)
	c, ok := s.unBrokeredChannels[localID]
	if ok {
		// channel has not been brokered
		go c.broker(s)
	} else {
		channelID, ok := s.localChannels.Get(localID)
		if !ok {
			log.Fatalf("AppSend invoked on channel ID %s: there's no localChannel with that ID.", localID)
		}
		c = s.brokeredChannels[channelID.(string)]
	}
	c.Outgoing[req.Recipient] <- *req.Content // enqueue message in outgoing buffer

	// TODO: reply with error code in case there is an error. eg buffer full.
	return &pb.AppSendResponse{Result: pb.Result_OK}, nil
}

// When the middleware receives a message in its public interface, it must enqueue it so that
// the corresponding local app can receive it
func (s *middlewareServer) MessageExchange(stream pb.PublicMiddleware_MessageExchangeServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Println("Received message from", in.SenderId, ":", string(in.Content.Body))

		// TODO: send to local app replacing sender name with local name

		ack := pb.ApplicationMessageWithHeaders{
			ChannelId:   in.ChannelId,
			RecipientId: in.SenderId,
			SenderId:    "provmwID-44",
			Content:     &pb.MessageContent{Body: []byte("ack")}}

		if err := stream.Send(&ack); err != nil {
			return err
		}
	}
}

// rpc invoked by broker to provider when initializing channel
func (s *middlewareServer) InitChannel(ctx context.Context, icr *pb.InitChannelRequest) (*pb.InitChannelResponse, error) {
	log.Printf("Running InitChannel")
	// InitChannelRequest: app_id, channel_id, participants (map[string]RemoteParticipant)
	if regapp, ok := s.registeredApps[icr.GetAppId()]; ok {
		// create registered channel with channel_id
		r := newSEARCHChannel(regapp.Contract)

		// TODO: refactor this section (repeated in broker func)
		r.addresses = icr.GetParticipants()
		r.ID = uuid.MustParse(icr.GetChannelId())
		r.LocalID = r.ID // in non locally initiated channels, there is a single channel_id

		// TODO: use mutex to handle maps
		s.brokeredChannels[r.ID.String()] = r
		s.localChannels.Insert(r.LocalID.String(), r.ID.String())

		for _, p := range r.Contract.GetRemoteParticipants() {
			go r.sender(p)
		}

		// notify local provider app
		*regapp.NotifyChan <- pb.InitChannelNotification{ChannelId: icr.GetChannelId()}
	} else {
		log.Printf("InitChannel with inexistent app_id %s", icr.GetAppId())
		// TODO: return error
	}
	return &pb.InitChannelResponse{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)

	pms := newMiddlewareServer()
	pb.RegisterPublicMiddlewareServer(grpcServer, pms)
	grpcServer.Serve(lis)
}
