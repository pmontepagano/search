package broker

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sync"
	"testing"

	"github.com/pmontepagano/search/contract"
	"github.com/pmontepagano/search/ent"
	pb "github.com/pmontepagano/search/gen/go/search/v1"
	"github.com/pmontepagano/search/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestBrokerRegisterProviderRequest(t *testing.T) {
	tmpDir := t.TempDir()
	b := NewBrokerServer(fmt.Sprintf("%s/t.db", tmpDir))
	t.Cleanup(b.Stop)
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

}

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
			t.Cleanup(b.Stop)

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
	t.Cleanup(b.Stop)

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
	t.Cleanup(b.Stop)

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
	t.Cleanup(b.Stop)

	// Requirement mocks.
	participantName := "test_participant_name"
	req := mocks.NewGlobalContract(t)
	reqProjection := mocks.NewLocalContract(t)
	req.EXPECT().GetProjection(participantName).Return(reqProjection, nil)
	req.EXPECT().GetContractID().Return("5cd3d822b5cb3383b8027750e0d8d2ccfc48a5d87020500d5c2c786f934a6937b6481b17bb0baa535cc24d55e68f0e6c87ee800e8053085d0b223cc5458c63b3")
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
	b.compatChecksWaitGroup.Wait()

	// Assert
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no compatible provider found for participant test_participant_name")
	assert.Equal(t, 3, numberOfCallsToCompatFunc)
}

func TestGetBestCandidate_OnlyOneCompatibleResult(t *testing.T) {
	// Arrange
	testDir := t.TempDir()
	b := NewBrokerServer(fmt.Sprintf("%s/t.db", testDir))
	t.Cleanup(b.Stop)

	// Requirement mocks.
	participantName := "test_participant_name"
	req := mocks.NewGlobalContract(t)
	reqProjection := mocks.NewLocalContract(t)
	req.EXPECT().GetProjection(participantName).Return(reqProjection, nil)
	req.EXPECT().GetContractID().Return("5cd3d822b5cb3383b8027750e0d8d2ccfc48a5d87020500d5c2c786f934a6937b6481b17bb0baa535cc24d55e68f0e6c87ee800e8053085d0b223cc5458c63b3")
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
	provider3Contract.EXPECT().GetContractID().Return("a62c6b1ee8a142d1086ffa68acc298705a293fcb09cf58784c0e390e44feafd8457d199b2bd239ac01f42f4939658927b68a43214f38d7a5e7a703da437374b8")
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
		if prov.GetContractID() == "a62c6b1ee8a142d1086ffa68acc298705a293fcb09cf58784c0e390e44feafd8457d199b2bd239ac01f42f4939658927b68a43214f38d7a5e7a703da437374b8" {
			return true, map[string]string{"serviceClient": "dunno"}, nil
		}
		return false, nil, nil
	})

	// Act
	result, err := b.getBestCandidate(context.TODO(), req, participantName)
	t.Log("Waiting for compatibility checks to finish...")
	b.compatChecksWaitGroup.Wait()

	// Assert
	assert.Equal(t, 3, numberOfCallsToCompatFunc)
	require.Nil(t, err)
	assert.Equal(t, provider3RegisteredProvider.ID, result.ID)
	assert.Equal(t, provider3AppID, result.ID)
	assert.Equal(t, provider3Url, result.URL)
	assert.Equal(t, provider3RegisteredContract.ID, result.ContractID)
}

func TestGetBestCandidate_CachedResult(t *testing.T) {
	// Arrange
	testDir := t.TempDir()
	b := NewBrokerServer(fmt.Sprintf("%s/t.db", testDir))
	t.Cleanup(b.Stop)

	// To save a cached result we first need to save the requirement projection and the provider contract.

	// Requirement mocks.
	participantName := "test_participant_name"
	req := mocks.NewGlobalContract(t)
	reqProjection := mocks.NewLocalContract(t)
	req.EXPECT().GetProjection(participantName).Return(reqProjection, nil)
	reqProjection.EXPECT().GetContractID().Return("1692526aab84461a8aebcefddcba2b33fb5897ab180c53e8b345ae125484d0aaa35baf60487050be21ed8909a48eace93851bf139087ce1f7a87d97b6120a651")
	reqProjection.EXPECT().GetFormat().Return(pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA)
	reqProjection.EXPECT().GetBytesRepr().Return([]byte(`dummy`))

	// Save requirement projection.
	reqProjectionRegisteredContract, err := b.getOrSaveContract(context.Background(), reqProjection)
	if err != nil {
		t.Fatalf("error saving contract: %v", err)
	}

	// Provider mocks.
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
	provider1RegisteredProvider, err := b.saveProvider(context.TODO(), provider1AppID, provider1Url, provider1RegisteredContract)
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

	// Now save a cached result for only provider 1.
	// We later want to check that the compatibility function is run only for provider 2.
	_, err = b.saveCompatibilityResult(context.Background(), reqProjectionRegisteredContract, provider1RegisteredContract, true,
		map[string]string{"serviceClient": "dunno"})
	if err != nil {
		t.Fatalf("error saving compatibility result: %v", err)
	}

	// Set broker's compatibility function to a mock that returns true for both providers. We later want to assert how many times the function was called.
	numberOfCallsToCompatFunc := 0
	counterLock := sync.Mutex{}
	b.SetCompatFunc(func(ctx context.Context, req contract.LocalContract, prov contract.LocalContract) (bool, map[string]string, error) {
		counterLock.Lock()
		numberOfCallsToCompatFunc++
		counterLock.Unlock()
		t.Logf("provID: %v, prov data: %s", prov.GetContractID(), prov.GetBytesRepr())
		return true, map[string]string{"serviceClient": "dunno"}, nil
	})

	// Act
	result, err := b.getBestCandidate(context.TODO(), req, participantName)
	b.compatChecksWaitGroup.Wait()

	// Assert
	assert.Nil(t, err)
	require.Equal(t, provider1RegisteredProvider.ContractID, result.ContractID)
	// assert.Equal(t, 1, numberOfCallsToCompatFunc) // Compatibility check is only run for provider 2.
}
