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

	"github.com/pmontepagano/search/contract"
	"github.com/pmontepagano/search/ent"
	"github.com/pmontepagano/search/ent/compatibilityresult"
	"github.com/pmontepagano/search/ent/registeredcontract"
	"github.com/pmontepagano/search/ent/registeredprovider"
	"github.com/pmontepagano/search/internal/searcherrors"

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

type candidateErr struct {
	participantName string
	err             error
}

func (e candidateErr) Error() string {
	return e.err.Error()
}

func (s *brokerServer) getBestCandidate(ctx context.Context, req contract.GlobalContract, p string) (*ent.RegisteredProvider, error) {
	// Get from the database the projection of the requirement contract (it should have been saved when handling BrokerChannelRequest).
	// rc *ent.RegisteredContract
	projection, err := req.GetProjection(p)
	if err != nil {
		return nil, candidateErr{p, fmt.Errorf("error getting projection for participant %s: %w", p, err)}
	}
	s.logger.Printf("running getBestCandidate for participant %s, requirement projection ID %s", p, projection.GetContractID())
	rc, err := s.dbClient.RegisteredContract.Get(ctx, projection.GetContractID())
	if err != nil {
		return nil, candidateErr{p, fmt.Errorf("error getting registered contract for participant %s", p)}
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
		return nil, candidateErr{p, fmt.Errorf("error querying compatibility results for participant %s: %w", p, err)}
	}
	if len(cachedResults) > 0 {
		// If there is at least one, return the best ranked one.
		// But before returning, send to a work queue a requests to calculate compatibility between this req and providers with
		// which compatibility has not yet been calculated.

		// Here we can have multiple contracts compatible and each contract can have more than one provider!

		// TODO: sort by ranking (first we need to add ranking to the schema).
		result, err := cachedResults[0].Edges.ProviderContract.QueryProviders().First(ctx)
		if err != nil {
			return nil, candidateErr{p, fmt.Errorf("error getting provider for participant %s: %w", p, err)}
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
		return nil, candidateErr{p, fmt.Errorf("error querying registered contracts for participant %s: %w", p, err)}
	}
	// allRegisteredProviders, err := s.dbClient.Debug().RegisteredProvider.Query().WithContract().All(ctx)
	// for _, rp := range allRegisteredProviders {
	// 	fmt.Printf("registered provider with participant name %s and contract ID: %s\n", rp.ParticipantName, rp.Edges.Contract.ID)
	// }
	if len(contractsToCalculate) == 0 {
		return nil, candidateErr{p, errors.New("no compatible provider found and no potential candidates to calculate compatibility")}
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
		return nil, candidateErr{p, fmt.Errorf("no compatible provider found for participant %s", p)}
	}

}

// given a contract and a list of participants, get best possible candidates for those candidates
// from all the candidates present in the registry
func (s *brokerServer) getBestCandidates(ctx context.Context, contract contract.GlobalContract, participants []string,
	denylistedProviders map[*pb.RemoteParticipant]struct{}) (map[string]*ent.RegisteredProvider, error) {
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

	// TODO: use denylistedProviders
	providers, err := iter.MapErr(participants, func(p *string) (*ent.RegisteredProvider, error) {
		return s.getBestCandidate(ctx, contract, *p)
	})
	if err != nil {
		failedParticipants := make([]string, 0)
		for _, e := range err.(interface{ Unwrap() []error }).Unwrap() {
			failedParticipants = append(failedParticipants, e.(candidateErr).participantName)
		}
		return nil, searcherrors.BrokerageFailedError(failedParticipants)
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

func (s *brokerServer) BrokerChannel(ctx context.Context, request *pb.BrokerChannelRequest) (*pb.BrokerChannelResponse, error) {
	if request.Contract == nil {
		return nil, status.New(codes.InvalidArgument, "contract is required").Err()
	}
	s.logger.Printf("Received broker request for contract: '%s'", request.Contract.Contract)
	reqContract, err := contract.ConvertPBGlobalContract(request.GetContract())
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "failed to parse requirement contract").Err()
	}
	presetParticipants := request.GetPresetParticipants()
	initiatorName := request.Contract.GetInitiatorName()

	for _, pName := range reqContract.GetRemoteParticipantNames() {
		if _, ok := presetParticipants[pName]; !ok {
			participantContract, err := reqContract.GetProjection(pName)
			if err != nil {
				return nil, status.New(codes.FailedPrecondition, err.Error()).Err()
			}
			_, err = s.getOrSaveContract(ctx, participantContract)
			if err != nil {
				return nil, status.New(codes.Unavailable, err.Error()).Err()
			}
			s.logger.Printf("Saved requirement projection for participant %s with ID %s", pName, participantContract.GetContractID())
		}
	}

	// TODO: split this function up in two functions: broker and initializeChannel
	// func (s *brokerServer) broker(ctx context.Context, reqContract contract.Contract, initiatorName string, presetParticipants map[string]struct{},
	//		denylistedProviders map[*ent.RegisteredProvider]struct{}) (map[string]*ent.RegisteredProviders, error)
	// func (s *brokerServer) initializeChannel(ctx context.Context, )... not sure yet what else needs

	// TODO: implement functionality to denylist providers.
	denylistedProviders := make(map[*pb.RemoteParticipant]struct{})

	// Broker participants.
	s.logger.Printf("brokerAndInitialize presetParticipants: %v", presetParticipants)
	s.logger.Printf("brokerAndInitialize initiatorName: %v", initiatorName)
	s.logger.Printf("brokerAndInitialize contract.GetParticipants(): %v", reqContract.GetParticipants())
	participantsToMatch := filterParticipants(reqContract.GetParticipants(), presetParticipants)
	s.logger.Printf("brokerAndInitialize participantsToMatch: %v", participantsToMatch)
	candidates, err := s.getBestCandidates(ctx, reqContract, participantsToMatch, denylistedProviders)
	if err != nil {
		s.logger.Printf("Could not get any provider candidates. Error: %s", err)
		return nil, err
	}

	// allParticipants will contain the RemoteParticipant object for all participants, including Service Client, preset providers and provider candidates
	// that the broker has selected. The key is the name of the participant according to the Service Client's requirement contract.
	allParticipants := make(map[string]*pb.RemoteParticipant)
	for pname, p := range candidates {
		allParticipants[pname] = getRemoteParticipant(p)
	}
	for pname, p := range presetParticipants {
		allParticipants[pname] = p
	}
	s.logger.Printf("allParticipants: %v", allParticipants)

	// TODO: why not simply accept the UUID generated by the Service Client? We should maybe change them and make them have a prefix
	// that is specific to the Service Client's ID (and authenticate this connection...?). This simplifies tracing, since we can use the
	// same channel ID everywhere and makes it very explicit who is the owner of the channel.
	channelID := uuid.New()

	initChannelCtx, cancelInitChannel := context.WithTimeout(ctx, 60*time.Second)
	defer cancelInitChannel()
	initChannelPool := pool.NewWithResults[struct {
		pname  string
		client pb.PublicMiddlewareServiceClient
		conn   *grpc.ClientConn
	}]().WithErrors().WithFirstError().WithContext(initChannelCtx).WithCancelOnError()

	// This channel is used to signal providers to which are failing so that we can exclude them and ask for new candidates.
	unresponsiveParticipants := make(chan string, len(allParticipants))

	// first round: InitChannel to all candidates and presetParticipants and wait for response
	// if any one of them does not respond, send its name through unresponsiveParticipants
	for pname, p := range allParticipants {
		if pname == initiatorName {
			continue
		}
		pname, p := pname, p // https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		initChannelPool.Go(func(ctx context.Context) (struct {
			pname  string
			client pb.PublicMiddlewareServiceClient
			conn   *grpc.ClientConn
		}, error) {
			result := struct {
				pname  string
				client pb.PublicMiddlewareServiceClient
				conn   *grpc.ClientConn
			}{pname, nil, nil}
			s.logger.Printf("Connecting to participant %s during brokerage. Provider AppId: %s", pname, p.AppId)

			// Connect to the provider's middleware.
			conn, err := grpc.DialContext(
				ctx,
				p.Url,
				grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO: use tls
				grpc.WithBlock(),
				grpc.FailOnNonTempDialError(true),
			)
			result.conn = conn
			if err != nil {
				unresponsiveParticipants <- pname
				return result, fmt.Errorf("failed to dial %s: %w", p.Url, err)
			}
			client := pb.NewPublicMiddlewareServiceClient(conn)
			result.client = client
			var participantsMapping map[string]*pb.RemoteParticipant
			participantsMapping, err = s.getParticipantMapping(reqContract, allParticipants, pname, candidates[pname])
			if err != nil {
				s.logger.Printf(err.Error())
				unresponsiveParticipants <- pname
				return result, fmt.Errorf("internal error calculating participant mapping: %w", err)
			}
			req := pb.InitChannelRequest{
				ChannelId:    channelID.String(),
				AppId:        p.AppId,
				Participants: participantsMapping,
			}
			res, err := client.InitChannel(ctx, &req)
			if err != nil {
				unresponsiveParticipants <- pname
				s.logger.Printf("Error doing InitChannel")
				return result, fmt.Errorf("unresponsive provider during InitChannel: %w", err)
			}
			if res.Result != pb.InitChannelResponse_RESULT_ACK {
				unresponsiveParticipants <- pname
				s.logger.Printf("Received non-ACK response to InitChannel")
				return result, fmt.Errorf("non-ACK received during InitChannel")
			}
			return result, nil
		})
	}

	providerMwClientsSlice, err := initChannelPool.Wait()
	if err != nil {
		s.logger.Printf("Failed to InitChannel at least one participant")
		close(unresponsiveParticipants)
		for pname := range unresponsiveParticipants {
			// TODO: what do we do if an unresponsive provider is one of the presetParticipants? Send error to Service Client.
			denylistedProviders[allParticipants[pname]] = struct{}{}
		}
		s.logger.Fatalf("TODO: failed to InitChannel at least one participant. Denylist providers and ask for new candidates.")
	}
	providerMwClients := make(map[string]struct {
		client pb.PublicMiddlewareServiceClient
		conn   *grpc.ClientConn
	})
	for _, p := range providerMwClientsSlice {
		pname, client, conn := p.pname, p.client, p.conn
		providerMwClients[pname] = struct {
			client pb.PublicMiddlewareServiceClient
			conn   *grpc.ClientConn
		}{client: client, conn: conn}
	}

	// Second round: when all responded ACK, signal them all to start choreography with StartChannelRequest.
	// Middlewares allow receiving messages on the channel after responding to the first round.
	// But the second round (StartChannel) is required to start the sender routines in all middlewares.

	// TODO: We could maybe improve this by only notifying the participants that can send a message from their start state.
	//   This would require that Middlewares also start the sender routines when they receive a message from the channel.
	s.logger.Printf("Brokering: second round...")

	startChannelCtx, startChannelCancel := context.WithTimeout(context.Background(), 60*time.Second)
	startChannelPool := pool.New().WithContext(startChannelCtx)

	for pname, p := range allParticipants {
		if pname == initiatorName {
			continue
		}
		pname, p := pname, p // https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		startChannelPool.Go(func(ctx context.Context) error {
			s.logger.Printf("Sending StartChannel to: %s", p.AppId)
			client := providerMwClients[pname].client
			defer providerMwClients[pname].conn.Close()

			req := pb.StartChannelRequest{
				ChannelId: channelID.String(),
				AppId:     p.AppId,
			}
			res, err := client.StartChannel(ctx, &req)
			if err != nil {
				s.logger.Printf("Error sending StartChannel to participant %s, with URL %s", pname, p.Url)
				return fmt.Errorf("Error sending StartChannel to participant %s, with URL %s. Error caused by %w", pname, p.Url, err)
			}
			if res.Result != pb.StartChannelResponse_RESULT_ACK {
				msg := fmt.Sprintf("Non-ACK response while sending StartChannel to participant %s, with URL %s", pname, p.Url)
				s.logger.Print(msg)
				return fmt.Errorf(msg)
			}
			return nil
		})
	}
	go func() {
		defer startChannelCancel()
		err := startChannelPool.Wait()
		if err != nil {
			// TODO: at this point, we can no longer recover. In the future we can implement a mechanism to notify the Service Client
			//   that this channel should be teared down.
			s.logger.Printf("Error sending StartChannel to at least one participant")
		}
	}()

	return &pb.BrokerChannelResponse{
		ChannelId:    channelID.String(),
		Participants: allParticipants,
	}, nil
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

// StartServer starts the broker server. If tls is true, certFile and keyFile must be provided. If notifyStartChan is not nil,
// the broker will send its public URL to the channel when it's already listening on the port.
func (s *brokerServer) StartServer(address string, tls bool, certFile string, keyFile string, notifyStartChan chan string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.Fatalf("failed to listen: %v", err)
	}
	s.PublicURL = lis.Addr().String()
	s.logger = log.New(os.Stderr, fmt.Sprintf("[BROKER] %s - ", s.PublicURL), log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	var opts []grpc.ServerOption
	if tls {
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
	if notifyStartChan != nil {
		s.logger.Print("Sending broker public URL to notifyStartChan...")
		notifyStartChan <- s.PublicURL
		s.logger.Print("Sent broker public URL to notifyStartChan...")
	}
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
