package broker

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sourcegraph/conc/iter"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/testdata"

	"github.com/clpombo/search/ent"
	"github.com/clpombo/search/internal/contract"

	pb "github.com/clpombo/search/gen/go/search/v1"
	_ "github.com/mattn/go-sqlite3"
)

type brokerServer struct {
	pb.UnimplementedBrokerServiceServer
	server *grpc.Server

	// db       *gorm.DB
	dbClient *ent.Client

	PublicURL string
	logger    *log.Logger
}

// TODO: move this into an ent template? https://entgo.io/docs/templates
func getRemoteParticipant(r *ent.RegisteredProvider) (res *pb.RemoteParticipant) {
	res.AppId = r.ID.String()
	res.Url = r.URL.String()
	return res
}

func getProviderContract(r *ent.RegisteredProvider) (contract.Contract, error) {
	return getContract(r.Edges.Contract)
}

func getContract(c *ent.RegisteredContract) (contract.Contract, error) {
	return contract.ConvertPBContract(&pb.Contract{
		Contract: c.Contract,
		Format:   pb.ContractFormat(c.Format),
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

func (s *brokerServer) getRegisteredProvider(appID string) (*ent.RegisteredProvider, error) {
	// TODO: add context as param
	// TODO: replace string param with uuid
	prov, err := s.dbClient.RegisteredProvider.Get(context.TODO(), uuid.MustParse(appID))
	if err != nil {
		return nil, fmt.Errorf("non registered appID %v", appID)
	}
	return prov, nil
}

func (s *brokerServer) getRandomRegisteredProvider() *ent.RegisteredProvider {
	// TODO: this should actually have the number of participants as a parameter
	// TODO: error handling
	// TODO: context as param
	r, _ := s.dbClient.RegisteredProvider.Query().First(context.TODO())
	return r
}

func (s *brokerServer) getBestCandidate(ctx context.Context, req contract.Contract, p string) (*ent.RegisteredProvider, error) {
	// Search in database if there are any compatible providers.
	// s.dbClient.CompatibilityResult.Query().Where(compatibilityresult.HasRequirementContractWith(registeredcontract.))

	// If there is at least one, return the best ranked one.
	// But before returning, send to a work queue a requests to calculate compatibility between this req and providers with
	// which compatibility has not yet been calculated.

	// It there is no result in the cache, also enqueue jobs to calculate compatibility with all providers for which we have not
	// calculated compatibility, but in this case wait until either they all fail to find
	// a compatible provider or we find one (and in that case we return the first one).

	return nil, nil
}

func (s *brokerServer) computeCompatibility(ctx context.Context, req contract.Contract, prov contract.Contract, reqParticipant string, provParticipant string) (bool, map[string]string) {
	return false, nil
}

// given a contract and a list of participants, get best possible candidates for those candidates
// from all the candidates present in the registry
func (s *brokerServer) getBestCandidates(ctx context.Context, contract contract.Contract, participants []string, blacklistedProviders map[*pb.RemoteParticipant]struct{}) (map[string]*ent.RegisteredProvider, error) {
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

	response := make(map[string]*ent.RegisteredProvider)

	// TODO: use blacklistedProviders
	providers, err := iter.MapErr(participants, func(p *string) (*ent.RegisteredProvider, error) {
		return s.getBestCandidate(ctx, contract, *p)
	})
	if err != nil {
		// TODO: better error handling.
		return nil, err
	}
	for idx, p := range participants {
		response[p] = providers[idx]
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
	res[receiverProvider.ParticipantName] = receiverRemoteParticipant

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
		receiverContract, err := getProviderContract(receiverProvider)
		if err != nil {
			s.logger.Fatal("error retrieving contract from provider")
		}
		if len(receiverContract.GetParticipants()) > 2 {
			s.logger.Fatal("Receiver has more than two participants in its contract. Cannot translate mapping, not yet implemented.")
		}
		var initiatorNameInReceiverContract string
		for _, p := range receiverContract.GetParticipants() {
			if p != receiverProvider.ParticipantName {
				initiatorNameInReceiverContract = p
				break
			}
		}
		res[initiatorNameInReceiverContract] = initiatorMapping[initiatorName]
	}
	s.logger.Printf("getParticipantMapping for %v returning %v", receiver, res)
	return res
}

// this routine will find compatible candidates and notify them
func (s *brokerServer) brokerAndInitialize(reqContract contract.Contract, presetParticipants map[string]*pb.RemoteParticipant,
	initiatorName string, blacklistedProviders map[*pb.RemoteParticipant]struct{}) {
	// TODO: split this function up in two functions: broker and initializeChannel
	// func (s *brokerServer) broker(ctx context.Context, reqContract contract.Contract, initiatorName string, presetParticipants map[string]struct{},
	//		blacklistedProviders map[*ent.RegisteredProvider]struct{}) (map[string]*ent.RegisteredProviders, error)
	// func (s *brokerServer) initializeChannel(ctx context.Context, )... not sure yet what else needs

	// Broker participants.
	s.logger.Printf("brokerAndInitialize presetParticipants: %v", presetParticipants)
	s.logger.Printf("brokerAndInitialize initiatorName: %v", initiatorName)
	s.logger.Printf("brokerAndInitialize contract.GetParticipants(): %v", reqContract.GetParticipants())
	participantsToMatch := filterParticipants(reqContract.GetParticipants(), presetParticipants)
	s.logger.Printf("brokerAndInitialize participantsToMatch: %v", participantsToMatch)
	candidates, err := s.getBestCandidates(context.TODO(), reqContract, participantsToMatch, blacklistedProviders)
	if err != nil {
		// TODO: if there is no result from getBestCandidates, we should notify error to initiator somehow
		s.logger.Fatal("Could not get any provider candidates.")
	}

	// Initialize channel.
	allParticipants := make(map[string]*pb.RemoteParticipant)
	for pname, p := range candidates {
		allParticipants[pname] = getRemoteParticipant(p)
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
	reqContract, err := contract.ConvertPBContract(request.GetContract())
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract")
	}
	presetParticipants := request.GetPresetParticipants()
	initiatorName := request.GetInitiatorName()

	_, err = s.getOrSaveContract(ctx, reqContract)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	go s.brokerAndInitialize(reqContract, presetParticipants, initiatorName, make(map[*pb.RemoteParticipant]struct{}))

	return &pb.BrokerChannelResponse{Result: pb.BrokerChannelResponse_RESULT_ACK}, nil
}

func (s *brokerServer) getOrSaveContract(ctx context.Context, c contract.Contract) (*ent.RegisteredContract, error) {
	contractHash := sha512.Sum512(c.GetBytesRepr())
	contractID := fmt.Sprintf("%x", contractHash[:])

	rc, err := s.dbClient.RegisteredContract.Get(ctx, contractID)
	if err != nil {
		if ent.IsNotFound(err) {
			rc, err = s.dbClient.RegisteredContract.Create().
				SetID(contractID).
				SetFormat(int(c.GetFormat().Number())).
				SetContract(c.GetBytesRepr()).
				Save(ctx)
			if err != nil {
				return nil, errors.New("database error registering contract")
			}
		} else {
			return nil, errors.New("database error fetching existing contract")
		}
	}
	return rc, nil
}

// we receive a LocalContract and the url, and we assign an AppID to this provider
func (s *brokerServer) RegisterProvider(ctx context.Context, req *pb.RegisterProviderRequest) (*pb.RegisterProviderResponse, error) {
	s.logger.Printf("Registering provider from URL: %s, contract '%s'", req.Url, req.Contract.Contract)

	appID := uuid.New()
	provContract, err := contract.ConvertPBContract(req.GetContract())
	if err != nil {
		st := status.New(codes.InvalidArgument, "invalid contract or format")
		return nil, st.Err()
	}
	providerUrl, err := url.Parse(req.Url)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid url for provider").Err()
	}
	rc, err := s.getOrSaveContract(ctx, provContract)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	// Save provider in the database.
	_, err = s.dbClient.RegisteredProvider.Create().
		SetID(appID).
		SetURL(providerUrl).
		SetParticipantName(req.ProviderName).
		SetContract(rc).
		Save(ctx)
	if err != nil {
		return nil, status.New(codes.Internal, "database error registering provider").Err()
	}

	return &pb.RegisterProviderResponse{AppId: appID.String()}, nil
}

// NewBrokerServer brokerServer constructor
func NewBrokerServer(databasePath string) *brokerServer {
	s := &brokerServer{}

	dbClient, err := ent.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	s.dbClient = dbClient
	if err := dbClient.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

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
	s.dbClient.Close()
}
