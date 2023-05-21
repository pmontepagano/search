package broker

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/clpombo/search/ent"
	pb "github.com/clpombo/search/gen/go/search/v1"
	"github.com/clpombo/search/mocks"

	"github.com/google/uuid"
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
