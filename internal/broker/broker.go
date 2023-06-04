package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"sort"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/sourcegraph/conc"
	"github.com/sourcegraph/conc/iter"
	"github.com/sourcegraph/conc/pool"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/testdata"

	"github.com/pmontepagano/search/contract"
	"github.com/pmontepagano/search/ent"
	"github.com/pmontepagano/search/ent/compatibilityresult"
	"github.com/pmontepagano/search/ent/registeredcontract"
	"github.com/pmontepagano/search/ent/registeredprovider"

	_ "github.com/mattn/go-sqlite3"
	pb "github.com/pmontepagano/search/gen/go/search/v1"
)

type brokerServer struct {
	pb.UnimplementedBrokerServiceServer
	server *grpc.Server

	dbClient   *ent.Client
	compatFunc ContractCompatibilityFunc

	PublicURL             string
	logger                *log.Logger
	compatChecksWaitGroup *conc.WaitGroup
}

// TODO: move this into an ent template? https://entgo.io/docs/templates
func getRemoteParticipant(r *ent.RegisteredProvider) *pb.RemoteParticipant {
	var res pb.RemoteParticipant
	res.AppId = r.ID.String()
	res.Url = r.URL.String()
	return &res
}

func getContract(c *ent.RegisteredContract) (contract.LocalContract, error) {
	return contract.ConvertPBLocalContract(&pb.LocalContract{
		Contract: c.Contract,
		Format:   pb.LocalContractFormat(c.Format),
	})
}

type ContractCompatibilityFunc func(ctx context.Context, req contract.LocalContract, prov contract.LocalContract) (bool, map[string]string, error)

func computeCompatibility(ctx context.Context, req contract.LocalContract, prov contract.LocalContract) (bool, map[string]string, error) {
	return false, nil, nil
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

func (s *brokerServer) getBestCandidate(ctx context.Context, req contract.GlobalContract, p string) (*ent.RegisteredProvider, error) {
	// Get from the database the projection of the requirement contract (it should have been saved when handling BrokerChannelRequest).
	// rc *ent.RegisteredContract
	projection, err := req.GetProjection(p)
	if err != nil {
		return nil, err
	}
	s.logger.Printf("running getBestCandidate for participant %s, requirement projection ID %s", p, projection.GetContractID())
	rc, err := s.dbClient.RegisteredContract.Get(ctx, projection.GetContractID())
	if err != nil {
		return nil, err
	}
	// Search in database if there are any compatible providers.
	cachedResults, err := s.dbClient.CompatibilityResult.
		Query().
		WithProviderContract().
		Where(compatibilityresult.And(
			compatibilityresult.Result(true),
			compatibilityresult.HasRequirementContractWith(registeredcontract.ID(rc.ID)),
		)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	if len(cachedResults) > 0 {
		// If there is at least one, return the best ranked one.
		// But before returning, send to a work queue a requests to calculate compatibility between this req and providers with
		// which compatibility has not yet been calculated.

		// Here we can have multiple contracts compatible and each contract can have more than one provider!

		// TODO: sort by ranking (first we need to add ranking to the schema).
		result, err := cachedResults[0].Edges.ProviderContract.QueryProviders().First(ctx)
		if err != nil {
			return nil, err
		}
		s.logger.Printf("found cached result for participant %s, returning provider with contract ID %s", p, result.ContractID)
		// TODO: enqueue jobs to calculate compatibility with all providers for which we have no data
		return result, nil
	}

	// If there is no result in the cache, also enqueue jobs to calculate compatibility with all providers for which we have not
	// calculated compatibility, but in this case wait until either they all fail to find
	// a compatible provider or we find one (and in that case we return the first one).
	// FALTA FILTRAR POR LOS QUE TENGAN CONTRATO DE REQUERIMIENTO IGUAL AL QUE ESTAMOS PROCESANDO.
	contractsToCalculate, err := s.dbClient.RegisteredContract.Query().
		Where(func(s *sql.Selector) {
			rpt := sql.Table(registeredprovider.Table)
			crt := sql.Table(compatibilityresult.Table)
			s.Join(rpt).On(s.C(registeredcontract.FieldID), rpt.C(registeredprovider.ContractColumn))
			s.LeftJoin(crt).On(s.C(registeredcontract.FieldID), crt.C(compatibilityresult.FieldProviderContractID))
		}).WithProviders().
		All(ctx)

	// for _, cDebug := range contractsToCalculate {
	// 	s.logger.Printf("provider contract to calculate with ID %s and content %s", cDebug.ID, cDebug.Contract)
	// }

	// providersToCompute, err := s.dbClient.Debug().RegisteredProvider.Query().
	// 	Where(
	// 		registeredprovider.Not(registeredprovider.HasContractWith(registeredcontract)
	if err != nil {
		return nil, err
	}
	// allRegisteredProviders, err := s.dbClient.Debug().RegisteredProvider.Query().WithContract().All(ctx)
	// for _, rp := range allRegisteredProviders {
	// 	fmt.Printf("registered provider with participant name %s and contract ID: %s\n", rp.ParticipantName, rp.Edges.Contract.ID)
	// }
	if len(contractsToCalculate) == 0 {
		return nil, errors.New("no compatible provider found and no potential candidates to calculate compatibility")
	}
	firstResult := make(chan *ent.RegisteredProvider, 1)
	workersFinished := make(chan struct{}, 1)
	s.compatChecksWaitGroup.Go(func() {
		workers := pool.New()
		for _, c := range contractsToCalculate {
			for _, prov := range c.Edges.Providers {
				thisContract := c
				thisProv := prov
				workers.Go(func() {
					provContract, err := getContract(thisContract)
					if err != nil {
						s.logger.Printf("error parsing registeredContract with ID %s", thisContract.ID)
						// TODO: handle error properly
						return
					}
					reqProjection, err := req.GetProjection(p)
					if err != nil {
						s.logger.Printf("error getting projection for participant %s", p)
						// TODO: handle error properly
						return
					}
					isCompatible, mapping, err := s.compatFunc(ctx, reqProjection, provContract)
					if err != nil {
						s.logger.Printf("error calculating compatibility: %v", err)
						// TODO: proper error handler
						return
					}
					_, err = s.saveCompatibilityResult(ctx, rc, thisContract, isCompatible, mapping)
					if err != nil {
						s.logger.Printf("error saving compatibility result: %v", err)
						// TODO: proper error handler
						return
					}
					if isCompatible {
						s.logger.Printf("found compatible provider %s, returning early...", thisProv.ID)
						firstResult <- thisProv
					}
				})
			}
		}
		s.logger.Printf("waiting for workers to finish...")
		workers.Wait()
		s.logger.Print("workers finished")
		workersFinished <- struct{}{}
		s.logger.Printf("Exiting routine for getBestCandidate with req %s and participant %s", req.GetContractID(), p)
	})

	select {
	case r := <-firstResult:
		return r, nil
	case <-workersFinished:
		return nil, fmt.Errorf("no compatible provider found for participant %s", p)
	}

}

// given a contract and a list of participants, get best possible candidates for those candidates
// from all the candidates present in the registry
func (s *brokerServer) getBestCandidates(ctx context.Context, contract contract.GlobalContract, participants []string,
	blacklistedProviders map[*pb.RemoteParticipant]struct{}) (map[string]*ent.RegisteredProvider, error) {
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

// getParticipantMapping takes the mapping between participant's names and RemoteParticipants from the
// initiator's perspective, and a participant name (also in the initiator's perspective), called receiver,
// and returns a mapping between participant names and RemoteParticipant's but using the receiver's
// perspective for participant names. The receiver has to be present in the registry because
// we need to parse its contract to get the mapping.
func (s *brokerServer) getParticipantMapping(req contract.GlobalContract, initiatorMapping map[string]*pb.RemoteParticipant, receiver string,
	receiverRegisteredProvider *ent.RegisteredProvider) (map[string]*pb.RemoteParticipant, error) {

	s.logger.Printf("getParticipantMapping for %v", receiver)
	_, ok := initiatorMapping[receiver]
	if !ok {
		return nil, fmt.Errorf("receiver %s not present in initiator's mapping", receiver)
	}

	projectedReq, err := req.GetProjection(receiver)
	if err != nil {
		return nil, fmt.Errorf("failed to get projection for %s", receiver)
	}

	compatResult, err := s.dbClient.CompatibilityResult.Query().Where(
		compatibilityresult.And(
			compatibilityresult.Result(true),
			compatibilityresult.HasRequirementContractWith(registeredcontract.ID(projectedReq.GetContractID())),
			compatibilityresult.HasProviderContractWith(registeredcontract.ID(receiverRegisteredProvider.ContractID)),
		),
	).First(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve mapping for %s", receiverRegisteredProvider.ID)
	}
	res := make(map[string]*pb.RemoteParticipant)
	for a, b := range compatResult.Mapping {
		res[a] = initiatorMapping[b]
	}

	s.logger.Printf("getParticipantMapping for %v returning %v", receiver, res)
	return res, nil
}

// this routine will find compatible candidates and notify them
func (s *brokerServer) brokerAndInitialize(reqContract contract.GlobalContract, presetParticipants map[string]*pb.RemoteParticipant,
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
		s.logger.Fatalf("Could not get any provider candidates. Error: %s", err)
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
			grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO: use tls
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
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var participantsMapping map[string]*pb.RemoteParticipant
		if pname == initiatorName {
			participantsMapping = allParticipants
		} else {
			participantsMapping, err = s.getParticipantMapping(reqContract, allParticipants, pname, candidates[pname])
			if err != nil {
				// TODO: proper error handling
				s.logger.Panicf(err.Error())
			}
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
		if res.Result != pb.InitChannelResponse_RESULT_ACK {
			// TODO: ??
			s.logger.Printf("Received non-ACK response to InitChannel")
			return
		}
	}

	// Second round: when all responded ACK, signal them all to start choreography with StartChannelRequest.
	// Middlewares allow receiving messages on the channel after responding to the first round.
	// But the second round (StartChannel) is required to start the sender routines in all middlewares.

	// TODO: We could maybe improve this by only notifying the participants that can send a message from their start state.
	//   This would require that Middlewares also start the sender routines when they receive a message from the channel.
	s.logger.Printf("Brokering: second round...")
	// TODO: Do this concurrently.
	for pname, p := range allParticipants {
		s.logger.Printf("Sending StartChannel to: %s", p.AppId)
		// TODO: refactor this repeated code
		conn, err := grpc.Dial(
			p.Url,
			grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO: use tls
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
		if res.Result != pb.StartChannelResponse_RESULT_ACK {
			// TODO: ??
			s.logger.Printf("Received non-ACK response to StartChannel")
			return
		}
	}

}

func (s *brokerServer) BrokerChannel(ctx context.Context, request *pb.BrokerChannelRequest) (*pb.BrokerChannelResponse, error) {
	s.logger.Printf("Received broker request for contract: '%s'", request.Contract.Contract)
	reqContract, err := contract.ConvertPBGlobalContract(request.GetContract())
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract")
	}
	presetParticipants := request.GetPresetParticipants()
	initiatorName := request.Contract.GetInitiatorName()

	for _, pName := range reqContract.GetRemoteParticipantNames() {
		if _, ok := presetParticipants[pName]; !ok {
			participantContract, err := reqContract.GetProjection(pName)
			if err != nil {
				return nil, status.New(codes.Internal, err.Error()).Err()
			}
			_, err = s.getOrSaveContract(ctx, participantContract)
			if err != nil {
				return nil, status.New(codes.Internal, err.Error()).Err()
			}
			s.logger.Printf("Saved requirement projection for participant %s with ID %s", pName, participantContract.GetContractID())
		}
	}

	go s.brokerAndInitialize(reqContract, presetParticipants, initiatorName, make(map[*pb.RemoteParticipant]struct{}))

	return &pb.BrokerChannelResponse{Result: pb.BrokerChannelResponse_RESULT_ACK}, nil
}

func (s *brokerServer) getOrSaveContract(ctx context.Context, c contract.LocalContract) (*ent.RegisteredContract, error) {
	contractID := c.GetContractID()

	tx, err := s.dbClient.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting a transaction: %w", err)
	}
	rc, err := tx.RegisteredContract.Get(ctx, contractID)
	if err != nil {
		if ent.IsNotFound(err) {
			rc, err = tx.RegisteredContract.Create().
				SetID(contractID).
				SetFormat(int(c.GetFormat().Number())).
				SetContract(c.GetBytesRepr()).
				Save(ctx)
			if err != nil {
				if rerr := tx.Rollback(); rerr != nil {
					err = fmt.Errorf("%w: %v", err, rerr)
				}
				return nil, fmt.Errorf("database error registering contract: %v", err)
			}
		} else {
			if rerr := tx.Rollback(); rerr != nil {
				err = fmt.Errorf("%w: %v", err, rerr)
			}
			return nil, fmt.Errorf("database error fetching existing contract: %v", err)
		}
	}
	return rc, tx.Commit()
}

// Save compatibility result in the database.
func (s *brokerServer) saveCompatibilityResult(ctx context.Context, req, prov *ent.RegisteredContract, areCompatible bool, mapping map[string]string) (*ent.CompatibilityResult, error) {
	tx, err := s.dbClient.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting a transaction: %w", err)
	}
	cr, err := tx.CompatibilityResult.Query().Where(compatibilityresult.And(
		compatibilityresult.HasRequirementContractWith(registeredcontract.ID(req.ID)),
		compatibilityresult.HasProviderContractWith(registeredcontract.ID(prov.ID))),
	).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			cr, err = tx.CompatibilityResult.Create().
				SetRequirementContract(req).
				SetProviderContract(prov).
				SetResult(areCompatible).
				SetMapping(mapping).
				Save(ctx)
			if err != nil {
				if rerr := tx.Rollback(); rerr != nil {
					err = fmt.Errorf("%w: %v", err, rerr)
				}
				return nil, fmt.Errorf("database error saving compatibility result: %v", err)
			}
		} else {
			if rerr := tx.Rollback(); rerr != nil {
				err = fmt.Errorf("%w: %v", err, rerr)
			}
			return nil, fmt.Errorf("database error fetching existing compatibility result: %v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("database error commiting compatibility result: %v", err)
	}
	s.logger.Printf("saved compatibility result between requirement projection %s and provider contract %s", req.ID, prov.ID)
	return cr, nil
}

// Save provider in the database.
func (s *brokerServer) saveProvider(ctx context.Context, appID uuid.UUID, url *url.URL, rc *ent.RegisteredContract) (*ent.RegisteredProvider, error) {
	prov, err := s.dbClient.RegisteredProvider.Create().
		SetID(appID).
		SetURL(url).
		SetContract(rc).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error registering provider: %v", err)
	}
	return prov, nil
}

// we receive a LocalContract and the url, and we assign an AppID to this provider
func (s *brokerServer) RegisterProvider(ctx context.Context, req *pb.RegisterProviderRequest) (*pb.RegisterProviderResponse, error) {
	s.logger.Printf("Registering provider from URL: %s, contract '%s'", req.Url, req.Contract.Contract)

	provContract, err := contract.ConvertPBLocalContract(req.GetContract())
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

	// Generate AppID for this new Service Provider.
	appID := uuid.New()

	// Save provider in the database.
	_, err = s.saveProvider(ctx, appID, providerUrl, rc)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	s.logger.Printf("Registered provider with URL %s and contract ID %v", providerUrl.String(), rc.ID)

	return &pb.RegisterProviderResponse{AppId: appID.String()}, nil
}

func setUpDB(databasePath string) *ent.Client {
	dbClient, err := openDB(databasePath)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	if err := dbClient.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return dbClient
}

// https://entgo.io/docs/sql-integration/
func openDB(databaseFilePath string) (*ent.Client, error) {
	dsn := fmt.Sprintf("file:%s?_fk=1&mode=rwc&busy_timeout=1000&cache=shared", databaseFilePath)
	drv, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	// Get the underlying sql.DB object of the driver.
	db := drv.DB()
	db.SetMaxOpenConns(1)
	return ent.NewClient(ent.Driver(drv)), nil
}

// NewBrokerServer brokerServer constructor
func NewBrokerServer(databasePath string) *brokerServer {
	s := &brokerServer{}

	s.dbClient = setUpDB(databasePath)
	s.compatFunc = computeCompatibility
	s.compatChecksWaitGroup = conc.NewWaitGroup()

	s.logger = log.New(os.Stderr, "[BROKER] - ", log.LstdFlags|log.Lmsgprefix)
	return s
}

func (s *brokerServer) StartServer(host string, port int, tls bool, certFile string, keyFile string) {
	s.PublicURL = fmt.Sprintf("%s:%d", host, port)
	s.logger = log.New(os.Stderr, fmt.Sprintf("[BROKER] %s - ", s.PublicURL), log.LstdFlags|log.Lmsgprefix|log.Lshortfile)

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
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *brokerServer) Stop() {
	s.logger.Printf("Broker server stopping...")
	s.compatChecksWaitGroup.Wait()
	if s.server != nil {
		s.server.Stop()
	}
	s.dbClient.Close()
}

func (s *brokerServer) SetCompatFunc(f ContractCompatibilityFunc) {
	s.compatFunc = f
}
