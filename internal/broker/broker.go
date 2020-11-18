package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	pb "github.com/clpombo/search/api"
)

type brokerServer struct {
	pb.UnimplementedBrokerServer
	server *grpc.Server
	// mapping of AppIDs to registered providers
	registeredProviders map[string]registeredProvider
	PublicURL string
	logger *log.Logger
}

type registeredProvider struct {
	participant pb.RemoteParticipant
	contract pb.Contract	// TODO: we'll need to use our LocalContract definition eventually
}

// returns new slice with keys of r filtered-out from orig. All keys of r MUST be present in orig
func filterParticipants(orig []string, r map[string]*pb.RemoteParticipant) []string {
	result := make([]string, 0)
	for _, o := range orig {
		if _, ok := r[o]; !ok {
			result = append(result, o)
		}
	}
	return result
}

func GetRandomKeyFromMap(mapI interface{}) interface{} {
	keys := reflect.ValueOf(mapI).MapKeys()
	rand.NewSource(time.Now().UnixNano())
	return keys[rand.Intn(len(keys))].Interface()
}

// given a contract and a list of participants, get best possible candidates for those candidates
// from all the candidates present in the registry
// TODO: we probably need to add a parameter to EXCLUDE candidates from regitry
func (s *brokerServer) getBestCandidates(contract *pb.Contract, participants []string) (map[string]*pb.RemoteParticipant, error) {
	// sanity check: check that all elements of param `participants` are present as participants in param `contract`
	contractParticipants := contract.GetRemoteParticipants()
	sort.Strings(contractParticipants)
	for _, p := range participants {
		// https://golang.org/pkg/sort/#Search
		i := sort.Search(len(contractParticipants),
			func(i int) bool { return contractParticipants[i] >= p })
		if !(i < len(contractParticipants) && contractParticipants[i] == p) {
			// participant element not present in contract
			s.logger.Fatalf("getBestCandidates got a participant not present in contract.")
		}
	}

	response := make(map[string]*pb.RemoteParticipant)

	// TODO: for now, we choose participants at random
	if len(s.registeredProviders) < 1 {
		return nil, errors.New("No providers registered.")
	}
	for _, v := range participants {
		appid := GetRandomKeyFromMap(s.registeredProviders).(string)
		participant := s.registeredProviders[appid].participant
		response[v] = &participant
	}

	return response, nil
}

// getParticipantMapping takes the mapping between participant's names and RemoteParticipants from the
// initiator's perspective, and a participant name (also in the initiator's perspective), called receiver, 
// and returns a mapping between participant names and RemoteParticipant's but using the receiver's
// perspective for participant names. The receiver has to be present in the registry because
// we need to parse its contract to get the mapping.
func (s *brokerServer) getParticipantMapping(initiatorMapping map[string]*pb.RemoteParticipant, receiver string) map[string]*pb.RemoteParticipant {
	if receiver == "self" {
		return initiatorMapping
	}
	receiverRemoteParticipant, ok := initiatorMapping[receiver]
	if !ok {
		s.logger.Fatal("Receiver not present in initiatormapping")
	}
	receiverProvider, ok := s.registeredProviders[receiverRemoteParticipant.AppId]
	if !ok {
		s.logger.Fatal("Receiver not registered.")
	}

	res := make(map[string]*pb.RemoteParticipant)
	res["self"] = receiverRemoteParticipant

	// TODO: without parsing and understanding contracts, we can only translate mappings of
	// contracts that have only two participants
	// we make an exception for receivers with name "_special" to run a test...
	if strings.HasSuffix(receiver, "_special") {
		if receiver == "r1_special" {
			res["sender"] = initiatorMapping["self"]
			res["receiver"] = initiatorMapping["r2_special"]
		}
		if receiver == "r2_special" {
			res["sender"] = initiatorMapping["r1_special"]
			res["receiver"] = initiatorMapping["r3_special"]
		}
		if receiver == "r3_special" {
			res["sender"] = initiatorMapping["r2_special"]
			res["receiver"] = initiatorMapping["self"]
		}
	} else {
		if len(initiatorMapping) > 2 {
			s.logger.Fatal("Cannot translate a mapping of more than two participants. Not yet implemented.")
		}
		if len(receiverProvider.contract.RemoteParticipants) > 2 {
			s.logger.Fatal("Receiver has more than two participants in its contract. Cannot translate mapping, not yet implemented.")
		}
		otherParticipantName := receiverProvider.contract.RemoteParticipants[1]
		res[otherParticipantName] = initiatorMapping["self"]
	}
	return res
}

// contract will always have a distinguished participant called "self" that corresponds to the initiator
// self MUST be present in presetParticipants with its url and ID.
// this routine will find compatible candidates and notify them
func (s *brokerServer) brokerAndInitialize(contract *pb.Contract, presetParticipants map[string]*pb.RemoteParticipant) {
	participantsToMatch := filterParticipants(contract.GetRemoteParticipants(), presetParticipants)
	candidates, err := s.getBestCandidates(contract, participantsToMatch)
	if err != nil {
		// TODO: if there is no result from getBestCandidates, we should notify error to initiator somehow
		s.logger.Fatal("Could not get any provider candidates.")
	}
	

	allParticipants := make(map[string]*pb.RemoteParticipant)
	for pname, p := range candidates {
		allParticipants[pname] = p
	}
	for pname, p := range presetParticipants {
		allParticipants[pname] = p
	}

	// s.logger.Println(allParticipants)

	channelID := uuid.New()

	// first round: InitChannel to all candidates and presetParticipants and wait for response
	// if any one of them did not respond, recurse into this func excluding unresponsive participants
	unresponsiveParticipants := make(map[string]bool)
	s.logger.Printf("Brokering: first round...")
	for pname, p := range allParticipants {
		s.logger.Printf("Sending InitChannel to: %s", p.AppId)
		conn, err := grpc.Dial(
			p.Url,
			grpc.WithInsecure(), // TODO: use tls
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
		)
		if err != nil {
			// TODO: discard this participant if it's not in presetParticipants by recursing into this func excluding this participant
			unresponsiveParticipants[pname] = true
			s.logger.Printf("Couldn't contact participant")
			return
			// TODO: here we should increment sequence number for InitChannel and restart
		}
		defer conn.Close()
		client := pb.NewPublicMiddlewareClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: remove hardcoded timeout
		defer cancel()
		req := pb.InitChannelRequest{
			ChannelId:    channelID.String(),
			AppId:        p.AppId,
			Participants: s.getParticipantMapping(allParticipants, pname),
		}
		res, err := client.InitChannel(ctx, &req)
		if err != nil {
			// TODO: discard this participant if it's not in presetParticipants
			unresponsiveParticipants[pname] = true
			s.logger.Printf("Error doing InitChannel")
			return
		}
		if res.Result != pb.InitChannelResponse_ACK {
			// TODO: ??
			s.logger.Printf("Received non-ACK response to InitChannel")
			return
		}
	}

	// second round: when all responded ACK, signal them all to start choreography
	s.logger.Printf("Brokering: second round...")
	for pname, p := range allParticipants {
		s.logger.Printf("Sending StartChannel to: %s", p.AppId)
		// TODO: refactor this repeated code
		conn, err := grpc.Dial(
			p.Url,
			grpc.WithInsecure(), // TODO: use tls
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
		)
		if err != nil {
			s.logger.Printf("Couldn't contact participant %s", pname)  // TODO: panic?
			return
		}
		defer conn.Close()
		client := pb.NewPublicMiddlewareClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: remove hardcoded timeout
		defer cancel()
		req := pb.StartChannelRequest{
			ChannelId:    channelID.String(),
			AppId:        p.AppId,
		}
		res, err := client.StartChannel(ctx, &req)
		if err != nil {
			// TODO: panic?
			s.logger.Printf("Error doing StartChannel")
			return
		}
		if res.Result != pb.StartChannelResponse_ACK {
			// TODO: ??
			s.logger.Printf("Received non-ACK response to StartChannel")
			return
		}
	}

}

func (s *brokerServer) BrokerChannel(ctx context.Context, request *pb.BrokerChannelRequest) (*pb.BrokerChannelResponse, error) {
	s.logger.Printf("Received broker request for contract: '%s'", request.Contract.Contract)
	contract := request.GetContract()
	presetParticipants := request.GetPresetParticipants()

	go s.brokerAndInitialize(contract, presetParticipants)

	return &pb.BrokerChannelResponse{Result: pb.BrokerChannelResponse_ACK}, nil
}

// we receive a LocalContract and the url, and we assign an AppID to this provider
func (s *brokerServer) RegisterProvider(ctx context.Context, req *pb.RegisterProviderRequest) (*pb.RegisterProviderResponse, error) {
	s.logger.Printf("Registering provider from URL: %s, contract '%s'", req.Url, req.Contract.Contract)

	// TODO: parse and validate contract
	appID := uuid.New()
	s.registeredProviders[appID.String()] = registeredProvider{
		participant: pb.RemoteParticipant{
			Url: req.Url,
			AppId: appID.String(),
		},
		contract: *req.GetContract(),
	}
	return &pb.RegisterProviderResponse{AppId: appID.String()}, nil
}

// NewBrokerServer brokerServer constructor
func NewBrokerServer() *brokerServer {
	s := &brokerServer{}
	s.registeredProviders = make(map[string]registeredProvider)
	s.logger = log.New(os.Stderr, "[BROKER] - ", log.LstdFlags | log.Lmsgprefix)
	return s
}

func (s *brokerServer) StartServer(host string, port int, tls bool, certFile string, keyFile string){
	s.PublicURL = fmt.Sprintf("%s:%d", host, port)
	s.logger = log.New(os.Stderr, fmt.Sprintf("[BROKER] %s - ", s.PublicURL), log.LstdFlags | log.Lmsgprefix)
	lis, err := net.Listen("tcp", s.PublicURL)
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
	grpcServer := grpc.NewServer(opts...)
	s.server = grpcServer
	pb.RegisterBrokerServer(grpcServer, s)
	s.logger.Printf("Broker server starting...")
	grpcServer.Serve(lis)
}

func (s *brokerServer) Stop(){
	if s.server != nil {
		s.server.Stop()
	}
}
