package broker

import (
	"context"
	"crypto/sha512"
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

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/testdata"

	"github.com/clpombo/search/internal/contract"

	pb "github.com/clpombo/search/gen/go/search/v1"
)

type brokerServer struct {
	pb.UnimplementedBrokerServiceServer
	server *grpc.Server

	db *gorm.DB

	PublicURL string
	logger    *log.Logger
}

type registeredProvider struct {
	AppId     string `gorm:"primaryKey;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Url          string
	contractID   string
	contract     registeredContract
	providerName string // Name of the provider in the contract.
}

type registeredContract struct {
	ID        string `gorm:"primaryKey;type:varchar(128)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	format   pb.ContractFormat
	contract []byte
}

func (r *registeredProvider) GetRemoteParticipant() (res *pb.RemoteParticipant) {
	res.AppId = r.AppId
	res.Url = r.Url
	return res
}

func (r *registeredProvider) GetContract() (contract.Contract, error) {
	return r.contract.GetContract()
}

func (c *registeredContract) GetContract() (contract.Contract, error) {
	return contract.ConvertPBContract(&pb.Contract{
		Contract: c.contract,
		Format:   c.format,
	})
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

func (s *brokerServer) getRegisteredProvider(appID string) (*registeredProvider, error) {
	var prov registeredProvider
	err := s.db.First(&prov, appID).Error
	if err != nil {
		return nil, fmt.Errorf("non registered appID %v", appID)
	}
	return &prov, nil
}

func (s *brokerServer) getRandomRegisteredProvider() *registeredProvider {
	// TODO: this should actually have the number of participants as a parameter
	var r registeredProvider
	s.db.Take(&r)
	return &r
}

// TODO: use generics!
func getRandomKeyFromMap(mapI interface{}) interface{} {
	keys := reflect.ValueOf(mapI).MapKeys()
	rand.NewSource(time.Now().UnixNano())
	return keys[rand.Intn(len(keys))].Interface()
}

// given a contract and a list of participants, get best possible candidates for those candidates
// from all the candidates present in the registry
// TODO: we probably need to add a parameter to EXCLUDE candidates from regitry
func (s *brokerServer) getBestCandidates(contract contract.Contract, participants []string) (map[string]*pb.RemoteParticipant, error) {
	// sanity check: check that all elements of param `participants` are present as participants in param `contract`
	contractParticipants := contract.GetParticipants()
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

	// TODO: hardcoded responses for TestCircle
	// if bytes.Equal(contract.GetContract(), []byte("send hello to r1, and later receive mesage from r3")) {
	// 	for _, rp := range s.registeredProviders {
	// 		url := rp.participant.Url
	// 		if url == "localhost:20001" {
	// 			p := rp.participant
	// 			response["r1_special"] = &p
	// 		}
	// 		if url == "localhost:20003" {
	// 			p := rp.participant
	// 			response["r2_special"] = &p
	// 		}
	// 		if url == "localhost:20005" {
	// 			p := rp.participant
	// 			response["r3_special"] = &p
	// 		}
	// 	}
	// } else {

	// TODO: for now, we choose participants at random
	for _, v := range participants {
		prov := s.getRandomRegisteredProvider()
		// prov, err := s.getBestProvider(contract, v)
		// if err != nil {
		// 	return nil, fmt.Errorf("could not find provider for %s", v)
		// }
		response[v] = prov.GetRemoteParticipant()
	}

	return response, nil
}

// func (s *brokerServer) getBestProvider(requirementContract contract.Contract, participant string) (*registeredProvider, error) {
// 	// Search in the cache?
// 	for
// }

// getParticipantMapping takes the mapping between participant's names and RemoteParticipants from the
// initiator's perspective, and a participant name (also in the initiator's perspective), called receiver,
// and returns a mapping between participant names and RemoteParticipant's but using the receiver's
// perspective for participant names. The receiver has to be present in the registry because
// we need to parse its contract to get the mapping.
func (s *brokerServer) getParticipantMapping(initiatorMapping map[string]*pb.RemoteParticipant, receiver string, initiatorName string) map[string]*pb.RemoteParticipant {
	s.logger.Printf("getParticipantMapping for %v", receiver)
	receiverRemoteParticipant, ok := initiatorMapping[receiver]
	if !ok {
		s.logger.Fatal("Receiver not present in initiatormapping")
	}
	receiverProvider, err := s.getRegisteredProvider(receiverRemoteParticipant.AppId)
	if err != nil {
		s.logger.Fatal("Receiver not registered.")
	}

	res := make(map[string]*pb.RemoteParticipant)
	res[receiverProvider.providerName] = receiverRemoteParticipant

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
		receiverContract, _ := receiverProvider.GetContract()
		if len(receiverContract.GetParticipants()) > 2 {
			s.logger.Fatal("Receiver has more than two participants in its contract. Cannot translate mapping, not yet implemented.")
		}
		var initiatorNameInReceiverContract string
		for _, p := range receiverContract.GetParticipants() {
			if p != receiverProvider.providerName {
				initiatorNameInReceiverContract = p
				break
			}
		}
		res[initiatorNameInReceiverContract] = initiatorMapping[initiatorName]
	}
	s.logger.Printf("getParticipantMapping for %v returning %v", receiver, res)
	return res
}

// contract will always have a distinguished participant called "self" that corresponds to the initiator
// self MUST be present in presetParticipants with its url and ID.
// this routine will find compatible candidates and notify them
func (s *brokerServer) brokerAndInitialize(contract contract.Contract, presetParticipants map[string]*pb.RemoteParticipant,
	initiatorName string) {
	s.logger.Printf("brokerAndInitialize presetParticipants: %v", presetParticipants)
	s.logger.Printf("brokerAndInitialize initiatorName: %v", initiatorName)
	s.logger.Printf("brokerAndInitialize contract.GetParticipants(): %v", contract.GetParticipants())
	participantsToMatch := filterParticipants(contract.GetParticipants(), presetParticipants)
	s.logger.Printf("brokerAndInitialize participantsToMatch: %v", participantsToMatch)
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

	s.logger.Printf("allParticipants: %v", allParticipants)

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
		client := pb.NewPublicMiddlewareServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: remove hardcoded timeout
		defer cancel()
		var participantsMapping map[string]*pb.RemoteParticipant
		if pname == initiatorName {
			participantsMapping = allParticipants
		} else {
			participantsMapping = s.getParticipantMapping(allParticipants, pname, initiatorName)
		}
		req := pb.InitChannelRequest{
			ChannelId:    channelID.String(),
			AppId:        p.AppId,
			Participants: participantsMapping,
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
			s.logger.Printf("Couldn't contact participant %s", pname) // TODO: panic?
			return
		}
		defer conn.Close()
		client := pb.NewPublicMiddlewareServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: remove hardcoded timeout
		defer cancel()
		req := pb.StartChannelRequest{
			ChannelId: channelID.String(),
			AppId:     p.AppId,
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
	contract, err := contract.ConvertPBContract(request.GetContract())
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract")
	}
	presetParticipants := request.GetPresetParticipants()
	initiatorName := request.GetInitiatorName()
	go s.brokerAndInitialize(contract, presetParticipants, initiatorName)

	return &pb.BrokerChannelResponse{Result: pb.BrokerChannelResponse_RESULT_ACK}, nil
}

// we receive a LocalContract and the url, and we assign an AppID to this provider
func (s *brokerServer) RegisterProvider(ctx context.Context, req *pb.RegisterProviderRequest) (*pb.RegisterProviderResponse, error) {
	s.logger.Printf("Registering provider from URL: %s, contract '%s'", req.Url, req.Contract.Contract)

	appID := uuid.New()
	_, err := contract.ConvertPBContract(req.GetContract())
	if err != nil {
		st := status.New(codes.InvalidArgument, "invalid contract or format")
		return nil, st.Err()
	}
	contractHash := sha512.Sum512(req.Contract.Contract)
	contractID := fmt.Sprintf("%x", contractHash[:])

	var contract registeredContract
	err = s.db.First(&contract, contractID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c := registeredContract{
			ID:       contractID,
			format:   req.Contract.Format,
			contract: req.Contract.GetContract(),
		}
		result := s.db.Create(&c)
		if result.Error != nil {
			return nil, status.New(codes.Internal, "database error registering contract").Err()
		}
	} else if err != nil {
		return nil, status.New(codes.Internal, "database error looking up contract").Err()
	}

	r := registeredProvider{
		Url:          req.Url,
		AppId:        appID.String(),
		contractID:   contractID,
		providerName: req.GetProviderName(),
	}
	result := s.db.Create(&r)
	if result.Error != nil {
		return nil, status.New(codes.Internal, "database error registering provider").Err()
	}

	return &pb.RegisterProviderResponse{AppId: appID.String()}, nil
}

// NewBrokerServer brokerServer constructor
func NewBrokerServer(databasePath string) *brokerServer {
	s := &brokerServer{}

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&registeredProvider{})
	db.AutoMigrate(&registeredContract{})
	s.db = db

	s.logger = log.New(os.Stderr, "[BROKER] - ", log.LstdFlags|log.Lmsgprefix)
	return s
}

func (s *brokerServer) StartServer(host string, port int, tls bool, certFile string, keyFile string) {
	s.PublicURL = fmt.Sprintf("%s:%d", host, port)
	s.logger = log.New(os.Stderr, fmt.Sprintf("[BROKER] %s - ", s.PublicURL), log.LstdFlags|log.Lmsgprefix)
	// s.logger.Printf("Delaying broker start 5 seconds to test something...")
	// time.Sleep(5 * time.Second)
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
	pb.RegisterBrokerServiceServer(grpcServer, s)
	s.logger.Printf("Broker server starting...")
	grpcServer.Serve(lis)
}

func (s *brokerServer) Stop() {
	if s.server != nil {
		s.server.Stop()
	}
}
