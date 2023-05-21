package broker

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sync"
	"testing"

	"github.com/clpombo/search/contract"
	"github.com/clpombo/search/ent"
	pb "github.com/clpombo/search/gen/go/search/v1"
	"github.com/clpombo/search/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestBrokerChannel_Request(t *testing.T) {
	// TODO: change this test and mock provider?
	tmpDir := t.TempDir()
	b := NewBrokerServer(fmt.Sprintf("%s/t.db", tmpDir))
	go b.StartServer("localhost", 3333, false, "", "")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", 3333), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBrokerServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// register dummy provider
	_, err = client.RegisterProvider(ctx, &pb.RegisterProviderRequest{
		Contract: &pb.LocalContract{
			Contract: []byte(`--
.outputs FooBar
.state graph
q0 0 ? hello q0
.marking q0
.end
`),
			Format: pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA,
		},
		Url: "fakeurl",
	})
	if err != nil {
		t.Fatalf("ERROR RegisterProvider: %v", err)
	}

	/* TODO: fix this test? Or move to an integration test.
	// ask for channel brokerage
	c := pb.GlobalContract{
		Contract: []byte(`--
		.outputs
		.state graph
		q0 1 ! hello q0
		.marking q0
		.end

		.outputs FooBar
		.state graph
		q0 0 ? hello q0
		.marking q0
		.end
		`),
		Format:        pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA,
		InitiatorName: "0",
	}
	req := pb.BrokerChannelRequest{
		Contract: &c,
		PresetParticipants: map[string]*pb.RemoteParticipant{
			"0": {
				Url:   "fake",
				AppId: "fake",
			},
		},
	}
	brokerresult, err := client.BrokerChannel(ctx, &req)
	if err != nil {
		t.Fatal("Received error from broker.")
	}
	if brokerresult.Result != pb.BrokerChannelResponse_RESULT_ACK {
		t.Fatal("Non ACK return code from broker")
	}

	*/
	b.Stop()
}

/* broken test
func TestGetParticipantMappingIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	b := NewBrokerServer(fmt.Sprintf("%s/t.db", tmpDir))
	go b.StartServer("localhost", 3333, false, "", "")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", 3333), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBrokerServiceClient(conn)
	ctx := context.Background()

	// register dummy provider
	var dummyProvContract = []byte(`--
.outputs FooBar
.state graph
q0 0 ? hello q0
.marking q0
.end
`)
	var dummyReqContract = []byte(`--
.outputs Dummy
.state graph
q0 FooBar ! hello q0
.marking q0
.end

.outputs FooBar
.state graph
q0 Dummy ? hello q0
.marking q0
.end
`)
	registrationResponse, err := client.RegisterProvider(ctx, &pb.RegisterProviderRequest{
		Contract: &pb.LocalContract{
			Contract: dummyProvContract,
			Format:   pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA,
		},
		Url: "fakeurl_for_provider",
	})
	if err != nil {
		log.Fatalf("ERROR RegisterProvider: %v", err)
	}
	provider, err := b.getRegisteredProvider(registrationResponse.AppId)
	if err != nil {
		t.Fatal("error getting registered provider")
	}

	b.SetCompatFunc(mockTestGetParticipantMappingContractCompatChecker)
	// test initiator mapping
	_, err = client.BrokerChannel(ctx, &pb.BrokerChannelRequest{
		Contract: &pb.GlobalContract{
			Contract:      dummyReqContract,
			Format:        pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA,
			InitiatorName: "Dummy",
		},
		PresetParticipants: map[string]*pb.RemoteParticipant{
			"Dummy": {
				Url:   "dummy_fake_url",
				AppId: "dummy_fake_appid",
			},
		},
	})
	if err != nil {
		t.Error("error brokering channel")
	}
	initiatorMapping := map[string]*pb.RemoteParticipant{
		"FooBar": {
			Url:   "foobar_fake_url",
			AppId: "foobar_fake_appid",
		},
		"0": {
			Url:   "fakeurl_for_provider",
			AppId: registrationResponse.GetAppId(),
		},
	}
	reqContract, err := contract.ConvertPBGlobalContract(&pb.GlobalContract{
		Contract:      dummyReqContract,
		InitiatorName: "FooBar",
		Format:        pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA,
	})
	if err != nil {
		t.Error("error converting contract")
	}

	_, err = b.getBestCandidate(context.Background(), reqContract, "FooBar")
	if err != nil {
		t.Fatal("error running getBestCandidate")
	}
	mapping, err := b.getParticipantMapping(reqContract, initiatorMapping, "0", provider)
	if err != nil {
		t.Fatal("error getting participant mapping")
	}
	expected := map[string]*pb.RemoteParticipant{
		"0": {
			Url:   "fakeurl_for_provider",
			AppId: registrationResponse.GetAppId(),
		},
		"FooBar": {
			Url:   "foobar_fake_url",
			AppId: "foobar_fake_appid",
		},
	}

	cmpOpts := []cmp.Option{
		protocmp.Transform(),
	}
	if !cmp.Equal(mapping, expected, cmpOpts...) {
		t.Errorf("Received erroneous response from getParticipantMapping: %v", mapping)
	}
	b.Stop()
}

func mockTestGetParticipantMappingContractCompatChecker(ctx context.Context, req contract.LocalContract, prov contract.LocalContract) (bool, map[string]string, error) {
	mapping := map[string]string{
		"0": "Dummy",
	}
	return true, mapping, nil
}
*/

func TestFilterParticipants(t *testing.T) {
	tests := []struct {
		name     string
		orig     []string
		r        map[string]*pb.RemoteParticipant
		expected []string
	}{
		{
			name:     "Filter some",
			orig:     []string{"a", "b", "c", "d"},
			r:        map[string]*pb.RemoteParticipant{"a": nil, "c": nil},
			expected: []string{"b", "d"},
		},
		{
			name:     "Empty",
			orig:     []string{},
			r:        map[string]*pb.RemoteParticipant{},
			expected: []string{},
		},
		{
			name:     "All filtered",
			orig:     []string{"a", "b", "c"},
			r:        map[string]*pb.RemoteParticipant{"a": nil, "b": nil, "c": nil},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := filterParticipants(tt.orig, tt.r)
			require.ElementsMatch(t, tt.expected, actual)
		})
	}
}

func TestGetParticipantMapping(t *testing.T) {
	type testcase struct {
		name                string
		req                 *mocks.GlobalContract
		reqProjection       *mocks.LocalContract
		initiatorMap        map[string]*pb.RemoteParticipant
		receiver            string
		providerReg         *ent.RegisteredProvider
		providerContract    *mocks.LocalContract
		want                map[string]*pb.RemoteParticipant
		wantErr             bool
		setup               func(*testcase) // We use this function to create the mocked contracts.
		compatResult        bool
		compatResultMapping map[string]string
	}
	tests := []testcase{
		{
			name: "receiver not in initiator map",
			initiatorMap: map[string]*pb.RemoteParticipant{
				"a": {Url: "a"},
				"b": {Url: "b"},
			},
			receiver: "c",
			want:     nil,
			wantErr:  true,
		},
		{
			name: "compatibility result not found",
			initiatorMap: map[string]*pb.RemoteParticipant{
				"a": {Url: "a"},
				"b": {Url: "b"},
				"c": {Url: "c"},
			},
			receiver:    "c",
			providerReg: &ent.RegisteredProvider{ID: uuid.UUID{}, ContractID: "1"},
			want:        nil,
			wantErr:     true,
			setup: func(tt *testcase) {
				tt.reqProjection.EXPECT().GetContractID().Return("83d95d53489d9fdd3d10c23da9a88020e0047e942bca77486a73b374b5d1d35bdf12dc998f4bc8dfe888931dd2fa5a48e6c71de0b08a9f184614b3f8679bb826")
				tt.reqProjection.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
				tt.reqProjection.EXPECT().GetBytesRepr().Return([]byte(`dummy`))
				tt.req.EXPECT().GetProjection(tt.receiver).Return(tt.reqProjection, nil)
			},
		},
		{
			name: "success",
			initiatorMap: map[string]*pb.RemoteParticipant{
				"TravelClient":            {Url: "travelclient.example.org/tc"},
				"HotelService":            {Url: "hotelservice.example.org"},
				"PaymentProcessorService": {Url: "pps.example.org"},
			},
			receiver:    "PaymentProcessorService",
			providerReg: &ent.RegisteredProvider{ID: uuid.UUID{}, ContractID: "8eb42c28c03d9d0160243f774fe4ed514b734ac4e24c80387b3d380df68cf6316758f47d7af317e59c43209a2cd674e284956136ab6568977e8c80841dfe548f"},
			want: map[string]*pb.RemoteParticipant{
				"Merchant": {Url: "travelclient.example.org/tc"},
			},
			wantErr: false,
			setup: func(tt *testcase) {
				tt.reqProjection.EXPECT().GetContractID().Return("83d95d53489d9fdd3d10c23da9a88020e0047e942bca77486a73b374b5d1d35bdf12dc998f4bc8dfe888931dd2fa5a48e6c71de0b08a9f184614b3f8679bb826")
				tt.reqProjection.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
				tt.reqProjection.EXPECT().GetBytesRepr().Return([]byte(`dummy`))
				tt.req.EXPECT().GetProjection(tt.receiver).Return(tt.reqProjection, nil)

				tt.providerContract.EXPECT().GetContractID().Return("8eb42c28c03d9d0160243f774fe4ed514b734ac4e24c80387b3d380df68cf6316758f47d7af317e59c43209a2cd674e284956136ab6568977e8c80841dfe548f")
				tt.providerContract.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
				tt.providerContract.EXPECT().GetBytesRepr().Return([]byte(`PaymentProcessorServiceFSA`))
			},
			compatResult: true,
			compatResultMapping: map[string]string{
				"Merchant": "TravelClient",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.req = mocks.NewGlobalContract(t)
				tt.reqProjection = mocks.NewLocalContract(t)
				tt.providerContract = mocks.NewLocalContract(t)
				tt.setup(&tt)
			}

			testDir := t.TempDir()
			b := NewBrokerServer(fmt.Sprintf("%s/t.db", testDir))

			if tt.req != nil {
				// Project requirement from Global Contract.
				projReq, err := tt.req.GetProjection(tt.receiver)
				if err != nil {
					t.Fatalf("error getting projection: %v", err)
				}
				// Save projected contract in the database.
				projReqReg, err := b.getOrSaveContract(context.Background(), projReq)
				if err != nil {
					t.Fatalf("error saving contract: %v", err)
				}
				// Save provider contract in the database.
				if tt.providerContract.IsMethodCallable(t, "GetBytesRepr") {
					provReg, err := b.getOrSaveContract(context.Background(), tt.providerContract)
					if err != nil {
						t.Fatalf("error saving contract: %v", err)
					}
					// Save compatibility result in the database.
					_, err = b.saveCompatibilityResult(context.Background(), projReqReg, provReg, tt.compatResult, tt.compatResultMapping)
					if err != nil {
						t.Fatalf("error saving compatibility result: %v", err)
					}
				}
			}

			got, err := b.getParticipantMapping(tt.req, tt.initiatorMap, tt.receiver, tt.providerReg)
			if (err != nil) != tt.wantErr {
				t.Errorf("getParticipantMapping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range got {
				wantValue, ok := tt.want[k]
				if !ok {
					t.Errorf("getParticipantMapping() got = %v, want %v", got, tt.want)
				}
				if v.GetUrl() != wantValue.GetUrl() {
					t.Errorf("getParticipantMapping() got = %v, want %v", got, tt.want)
				}
			}
			for k, v := range tt.want {
				gotValue, ok := got[k]
				if !ok {
					t.Errorf("getParticipantMapping() got = %v, want %v", got, tt.want)
				}
				if v.GetUrl() != gotValue.GetUrl() {
					t.Errorf("getParticipantMapping() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// Test getBestCandidate with an empty database.
func TestGetBestCandidate_UnregisteredRequirementContract(t *testing.T) {
	// Arrange
	testDir := t.TempDir()
	b := NewBrokerServer(fmt.Sprintf("%s/t.db", testDir))

	participantName := "test_participant_name"
	req := mocks.NewGlobalContract(t)
	reqProjection := mocks.NewLocalContract(t)
	req.EXPECT().GetProjection(participantName).Return(reqProjection, nil)
	reqProjection.EXPECT().GetContractID().Return("req_projection_contract_id")

	// Act
	result, err := b.getBestCandidate(context.TODO(), req, participantName)

	// Assert
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ent: registered_contract not found")
}

func TestGetBestCandidate_NoCompatibilityResults(t *testing.T) {
	// Arrange
	testDir := t.TempDir()
	b := NewBrokerServer(fmt.Sprintf("%s/t.db", testDir))

	participantName := "test_participant_name"
	req := mocks.NewGlobalContract(t)
	reqProjection := mocks.NewLocalContract(t)
	req.EXPECT().GetProjection(participantName).Return(reqProjection, nil)
	reqProjection.EXPECT().GetContractID().Return("1692526aab84461a8aebcefddcba2b33fb5897ab180c53e8b345ae125484d0aaa35baf60487050be21ed8909a48eace93851bf139087ce1f7a87d97b6120a651")
	reqProjection.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	reqProjection.EXPECT().GetBytesRepr().Return([]byte(`dummy`))

	_, err := b.getOrSaveContract(context.Background(), reqProjection)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}

	// Act
	result, err := b.getBestCandidate(context.TODO(), req, participantName)

	// Assert
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no compatible provider found and no potential candidates to calculate compatibility")
}

const (
	provider1FSA = `
		.outputs provider1
		.state graph
		q0 serviceClient ? hello q1
		q1 serviceClient ! ack q0
		.marking q0
		.end`
	provider2FSA = `
		.outputs provider2
		.state graph
		q0 serviceClient ? hello q1
		q1 serviceClient ! ack q0
		.marking q0
		.end`
	provider3FSA = `
		.outputs provider3
		.state graph
		q0 serviceClient ? hello q1
		q1 serviceClient ! ack q0
		.marking q0
		.end`
)

func TestGetBestCandidate_NoCompatibleResultsWithRegisteredProviders(t *testing.T) {
	// Arrange
	testDir := t.TempDir()
	b := NewBrokerServer(fmt.Sprintf("%s/t.db", testDir))

	// Requirement mocks.
	participantName := "test_participant_name"
	req := mocks.NewGlobalContract(t)
	reqProjection := mocks.NewLocalContract(t)
	req.EXPECT().GetProjection(participantName).Return(reqProjection, nil)
	reqProjection.EXPECT().GetContractID().Return("1692526aab84461a8aebcefddcba2b33fb5897ab180c53e8b345ae125484d0aaa35baf60487050be21ed8909a48eace93851bf139087ce1f7a87d97b6120a651")
	reqProjection.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	reqProjection.EXPECT().GetBytesRepr().Return([]byte(`dummy`))

	// Save requirement projection.
	_, err := b.getOrSaveContract(context.Background(), reqProjection)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}

	// Save three different providers in the database but make them incompatible.
	provider1Contract := mocks.NewLocalContract(t)
	provider1Contract.EXPECT().GetContractID().Return("62c4660d4cb1cd6074379e99b7859bfd7f0ecd9796a0923a0ba7355be5e65f81939dd62c2dcebf0e866028026f9a4817570bdf5d6b5512c05967ad828e515836")
	provider1Contract.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	provider1Contract.EXPECT().GetBytesRepr().Return([]byte(provider1FSA))
	provider1RegisteredContract, err := b.getOrSaveContract(context.Background(), provider1Contract)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}
	provider1AppID, _ := uuid.Parse("285b25a8-6621-4abe-a18e-90cc896f3476")
	provider1Url, _ := url.Parse("provider1.example.org:8080")
	_, err = b.saveProvider(context.TODO(), provider1AppID, provider1Url, provider1RegisteredContract)
	if err != nil {
		t.Fatalf("error saving provider: %v", err)
	}

	provider2Contract := mocks.NewLocalContract(t)
	provider2Contract.EXPECT().GetContractID().Return("97760560dbca16c5626ea91b81058b63fb94d0f4e0b4d93867af52df276b80210e249d3a0274b23251682ac27ed3d016fb8911f336c567b0c33639ae706fb880")
	provider2Contract.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	provider2Contract.EXPECT().GetBytesRepr().Return([]byte(provider2FSA))
	provider2RegisteredContract, err := b.getOrSaveContract(context.Background(), provider2Contract)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}
	provider2AppID, _ := uuid.Parse("4cda7849-2288-4d61-86df-1d98d273eb1e")
	provider2Url, _ := url.Parse("provider2.example.org:8080")
	_, err = b.saveProvider(context.TODO(), provider2AppID, provider2Url, provider2RegisteredContract)
	if err != nil {
		t.Fatalf("error saving provider: %v", err)
	}

	provider3Contract := mocks.NewLocalContract(t)
	provider3Contract.EXPECT().GetContractID().Return("bd128995fe64da655ccc27c5d07714cb76478ae749df31f867cdc7cb83a0cb868b076679dc3665142540359aa96faf8667fd015d846537c6403c695574a34fa0")
	provider3Contract.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	provider3Contract.EXPECT().GetBytesRepr().Return([]byte(provider3FSA))
	provider3RegisteredContract, err := b.getOrSaveContract(context.Background(), provider3Contract)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}
	provider3AppID, _ := uuid.Parse("5aaca746-3008-4c0a-984c-a6a5c85ca6d3")
	provider3Url, _ := url.Parse("provider3.example.org:8080")
	_, err = b.saveProvider(context.TODO(), provider3AppID, provider3Url, provider3RegisteredContract)
	if err != nil {
		t.Fatalf("error saving provider: %v", err)
	}

	// Set broker's compatibility function to a mock that returns false. We later want to assert how many times the function was called.
	numberOfCallsToCompatFunc := 0
	counterLock := sync.Mutex{}
	b.SetCompatFunc(func(ctx context.Context, req contract.LocalContract, prov contract.LocalContract) (bool, map[string]string, error) {
		counterLock.Lock()
		numberOfCallsToCompatFunc++
		counterLock.Unlock()
		return false, nil, nil
	})

	// Act
	result, err := b.getBestCandidate(context.TODO(), req, participantName)

	// Assert
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no compatible provider found for participant test_participant_name")
	assert.Equal(t, 3, numberOfCallsToCompatFunc)
	b.Stop()
}

func TestGetBestCandidate_OnlyOneCompatibleResult(t *testing.T) {
	// Arrange
	testDir := t.TempDir()
	b := NewBrokerServer(fmt.Sprintf("%s/t.db", testDir))

	// Requirement mocks.
	participantName := "test_participant_name"
	req := mocks.NewGlobalContract(t)
	reqProjection := mocks.NewLocalContract(t)
	req.EXPECT().GetProjection(participantName).Return(reqProjection, nil)
	reqProjection.EXPECT().GetContractID().Return("1692526aab84461a8aebcefddcba2b33fb5897ab180c53e8b345ae125484d0aaa35baf60487050be21ed8909a48eace93851bf139087ce1f7a87d97b6120a651")
	reqProjection.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	reqProjection.EXPECT().GetBytesRepr().Return([]byte(`dummy`))

	// Save requirement projection.
	_, err := b.getOrSaveContract(context.Background(), reqProjection)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}

	// Save three different providers in the database.
	provider1Contract := mocks.NewLocalContract(t)
	provider1Contract.EXPECT().GetContractID().Return("826516eb343f93d8f1ca808ef34b30deb024bffd05c1630e0eee81243faeca243e9ce0cf431d691ce63b2ebaf60e48e1220f1d4ab85731b5c2f591fe1e958043")
	provider1Contract.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	provider1Contract.EXPECT().GetBytesRepr().Return([]byte(provider1FSA))
	provider1RegisteredContract, err := b.getOrSaveContract(context.Background(), provider1Contract)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}
	provider1AppID, _ := uuid.Parse("285b25a8-6621-4abe-a18e-90cc896f3476")
	provider1Url, _ := url.Parse("provider1.example.org:8080")
	_, err = b.saveProvider(context.TODO(), provider1AppID, provider1Url, provider1RegisteredContract)
	if err != nil {
		t.Fatalf("error saving provider: %v", err)
	}

	provider2Contract := mocks.NewLocalContract(t)
	provider2Contract.EXPECT().GetContractID().Return("64867f97a221cd291084fc7c27516fca3c3dfebaae7ec99f5f3f836c540fcec6978a3f8d7a98c3ca6ee8039e550ea125b88d369c033cc203820caa562f3c7986")
	provider2Contract.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	provider2Contract.EXPECT().GetBytesRepr().Return([]byte(provider2FSA))
	provider2RegisteredContract, err := b.getOrSaveContract(context.Background(), provider2Contract)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}
	provider2AppID, _ := uuid.Parse("4cda7849-2288-4d61-86df-1d98d273eb1e")
	provider2Url, _ := url.Parse("provider2.example.org:8080")
	_, err = b.saveProvider(context.TODO(), provider2AppID, provider2Url, provider2RegisteredContract)
	if err != nil {
		t.Fatalf("error saving provider: %v", err)
	}

	provider3Contract := mocks.NewLocalContract(t)
	provider3Contract.EXPECT().GetContractID().Return("bd128995fe64da655ccc27c5d07714cb76478ae749df31f867cdc7cb83a0cb868b076679dc3665142540359aa96faf8667fd015d846537c6403c695574a34fa0")
	provider3Contract.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	provider3Contract.EXPECT().GetBytesRepr().Return([]byte(provider3FSA))
	provider3RegisteredContract, err := b.getOrSaveContract(context.Background(), provider3Contract)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}
	provider3AppID, _ := uuid.Parse("5aaca746-3008-4c0a-984c-a6a5c85ca6d3")
	provider3Url, _ := url.Parse("provider3.example.org:8080")
	provider3RegisteredProvider, err := b.saveProvider(context.TODO(), provider3AppID, provider3Url, provider3RegisteredContract)
	if err != nil {
		t.Fatalf("error saving provider: %v", err)
	}

	// Set broker's compatibility function to a mock that returns true for provider 3. We later want to assert how many times the function was called.
	numberOfCallsToCompatFunc := 0
	counterLock := sync.Mutex{}
	b.SetCompatFunc(func(ctx context.Context, req contract.LocalContract, prov contract.LocalContract) (bool, map[string]string, error) {
		counterLock.Lock()
		numberOfCallsToCompatFunc++
		counterLock.Unlock()
		t.Logf("provID: %v, prov data: %s", prov.GetContractID(), prov.GetBytesRepr())
		if prov.GetContractID() == "bd128995fe64da655ccc27c5d07714cb76478ae749df31f867cdc7cb83a0cb868b076679dc3665142540359aa96faf8667fd015d846537c6403c695574a34fa0" {
			return true, map[string]string{"serviceClient": "dunno"}, nil
		}
		return false, nil, nil
	})

	// Act
	result, err := b.getBestCandidate(context.TODO(), req, participantName)

	// Assert
	assert.Equal(t, provider3RegisteredProvider.ID, result.ID)
	assert.Equal(t, provider3AppID, result.ID)
	assert.Equal(t, provider3Url, result.URL)
	assert.Equal(t, provider3RegisteredContract.ID, result.ContractID)
	assert.Nil(t, err)
	assert.Equal(t, 3, numberOfCallsToCompatFunc)
	b.Stop()
}
