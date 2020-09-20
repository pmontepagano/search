package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"reflect"
	"sort"
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
	registeredProviders map[string]registeredProvider
}

type registeredProvider struct {
	participant pb.RemoteParticipant
	contract pb.Contract	// TODO: we'll need to use our LocalContract definition eventually
}

// returns new slice with keys of r filtered-out from orig. All keys of r MUST be present in orig
func filterParticipants(orig []string, r map[string]*pb.RemoteParticipant) []string {
	result := make([]string, len(orig)-len(r))
	sort.Strings(orig)
	nOrig := 0
	nDest := 0
	for p := range r {
		if orig[nOrig] != p {
			result[nDest] = orig[nOrig]
			nDest++
		}
		nOrig++
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
			log.Fatalf("getBestCandidates got a participant not present in contract.")
		}
	}

	response := make(map[string]*pb.RemoteParticipant)

	// TODO: for now, we choose participants at random
	if len(s.registeredProviders) < 1 {
		return nil, errors.New("No providers registered.")
	}
	for _, v := range participants {
		log.Println("Received requirements contract with participant", v)
		
		appid := GetRandomKeyFromMap(s.registeredProviders).(string)
		participant := s.registeredProviders[appid].participant
		response[v] = &participant
	}

	return response, nil
}

// contract will always have a distinguished participant called "self" that corresponds to the initiator
// self MUST be present in presetParticipants with its url and ID.
// this routine will find compatible candidates and notify them
func (s *brokerServer) brokerAndInitialize(contract *pb.Contract, presetParticipants map[string]*pb.RemoteParticipant) {
	participantsToMatch := filterParticipants(contract.GetRemoteParticipants(), presetParticipants)
	candidates, err := s.getBestCandidates(contract, participantsToMatch)
	if err != nil {
		// TODO: if there is no result from getBestCandidates, we should notify error to initiator somehow
		log.Fatal("Could not get any provider candidates.")
	}
	

	allParticipants := make(map[string]*pb.RemoteParticipant)
	for pname, p := range candidates {
		allParticipants[pname] = p
	}
	for pname, p := range presetParticipants {
		allParticipants[pname] = p
	}

	log.Println(allParticipants)

	channelID := uuid.New()

	// first round: InitChannel to all candidates and presetParticipants and wait for response
	// if any one of them did not respond, recurse into this func excluding unresponsive participants
	unresponsiveParticipants := make(map[string]bool)
	for pname, p := range allParticipants {
		log.Printf("Brokering, first round. Contacting %s", p.AppId)
		conn, err := grpc.Dial(
			p.Url,
			grpc.WithInsecure(), // TODO: use tls
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
		)
		if err != nil {
			// TODO: discard this participant if it's not in presetParticipants by recursing into this func excluding this participant
			unresponsiveParticipants[pname] = true
			log.Printf("Couldn't contact participant")
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
			Participants: allParticipants, // TODO: this should be a different map for each participant, using their local participants' names as keys
		}
		res, err := client.InitChannel(ctx, &req)
		if err != nil {
			// TODO: discard this participant if it's not in presetParticipants
			unresponsiveParticipants[pname] = true
			log.Printf("Error doing InitChannel")
			return
		}
		if res.Result != pb.InitChannelResponse_ACK {
			// TODO: ??
			log.Printf("Received non-ACK response to InitChannel")
			return
		}
	}

	// second round: when all responded ACK, signal them all to start choreography

}

func (s *brokerServer) BrokerChannel(ctx context.Context, request *pb.BrokerChannelRequest) (*pb.BrokerChannelResponse, error) {
	log.Println("BROKER: brokering...")
	contract := request.GetContract()
	presetParticipants := request.GetPresetParticipants()

	go s.brokerAndInitialize(contract, presetParticipants)

	return &pb.BrokerChannelResponse{Result: pb.BrokerChannelResponse_ACK}, nil
}

// we receive a LocalContract and the url, and we assign an AppID to this provider
func (s *brokerServer) RegisterProvider(ctx context.Context, req *pb.RegisterProviderRequest) (*pb.RegisterProviderResponse, error) {
	log.Printf("BROKER: Registering provider from URL: %s", req.Url)

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
	return s
}

func (s *brokerServer) StartServer(host string, port int, tls bool, certFile string, keyFile string){
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host ,port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
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
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	s.server = grpcServer
	pb.RegisterBrokerServer(grpcServer, s)
	log.Printf("Broker server starting...")
	grpcServer.Serve(lis)
}

func (s *brokerServer) Stop(){
	if s.server != nil {
		s.server.Stop()
	}
}
