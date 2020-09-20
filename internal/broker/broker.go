package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
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
	savedData []*pb.RemoteParticipant // read-only after initialized
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

// given a contract and a list of participants, get best possible candidates for those candidates
// from all the candidates present in the registry
// TODO: we probably need to add a parameter to EXCLUDE candidates from regitry
func (s *brokerServer) getBestCandidates(contract pb.Contract, participants []string) map[string]*pb.RemoteParticipant {
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
	rand.Seed(time.Now().Unix())
	for _, v := range participants {
		log.Println("Received requirements contract with participant", v)
		response[v] = s.savedData[rand.Intn(len(s.savedData))]
	}

	return response
}

// contract will always have a distinguished participant called "self" that corresponds to the initiator
// self MUST be present in presetParticipants with its url and ID.
// this routine will find compatible candidates and notify them
func (s *brokerServer) brokerAndInitialize(contract *pb.Contract, presetParticipants map[string]*pb.RemoteParticipant) {
	participantsToMatch := filterParticipants(contract.GetRemoteParticipants(), presetParticipants)
	candidates := s.getBestCandidates(*contract, participantsToMatch)
	// TODO: if there is no result from getBestCandidates, we should notify error to initiator somehow

	allParticipants := make(map[string]*pb.RemoteParticipant)
	for pname, p := range candidates {
		allParticipants[pname] = p
	}
	for pname, p := range presetParticipants {
		allParticipants[pname] = p
	}

	log.Println(candidates)

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
			log.Fatalf("Couldn't contact participant")
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
			log.Fatalf("Error doing InitChannel")
		}
		if res.Result != pb.InitChannelResponse_ACK {
			// TODO: ??
			log.Fatalf("Received non-ACK response to InitChannel")
		}
	}

	// second round: when all responded ACK, signal them all to start choreography

}

func (s *brokerServer) BrokerChannel(ctx context.Context, request *pb.BrokerChannelRequest) (*pb.BrokerChannelResponse, error) {
	log.Println("brokering...")
	contract := request.GetContract()
	presetParticipants := request.GetPresetParticipants()

	go s.brokerAndInitialize(contract, presetParticipants)

	return &pb.BrokerChannelResponse{Result: pb.BrokerChannelResponse_ACK}, nil
}

// loads data from a JSON file
func (s *brokerServer) loadData(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to load default data: %v", err)
		}
	} else {
		data = exampleData
	}
	if err := json.Unmarshal(data, &s.savedData); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}
}

// NewBrokerServer brokerServer constructor
func NewBrokerServer(jsonDBFile string) *brokerServer {
	s := &brokerServer{}
	s.loadData(jsonDBFile)
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
	pb.RegisterBrokerServer(grpcServer, s)
	log.Printf("Broker server starting...")
	grpcServer.Serve(lis)
}

var exampleData = []byte(`[{
	"Url": "providermiddleware",
	"AppId": "example1FAKEuuid",
	"Contract": {
		"Contract": "unparsed",
		"RemoteParticipants": ["p2"]
	}
}, {
	"Url": "providermiddleware",
	"AppId": "example2FAKEuuid",
	"Contract": {
		"Contract": "unparsed",
		"RemoteParticipants": ["p2", "p3"]
	}
}]`)
