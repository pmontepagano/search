package broker

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/clpombo/search/ent"
	"github.com/clpombo/search/ent/enttest"
	pb "github.com/clpombo/search/gen/go/search/v1"
	"github.com/clpombo/search/internal/contract"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/testing/protocmp"
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
		log.Fatalf("ERROR RegisterProvider: %v", err)
	}

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
	tests := []struct {
		name            string
		req             contract.GlobalContract
		initiatorMap    map[string]*pb.RemoteParticipant
		receiver        string
		receiverRegProv *ent.RegisteredProvider
		want            map[string]*pb.RemoteParticipant
		wantErr         bool
	}{
		{
			name: "receiver not in initiator map",
			req:  contract.GlobalContract{ContractID: 1},
			initiatorMap: map[string]*pb.RemoteParticipant{
				"a": {Name: "a"},
				"b": {Name: "b"},
			},
			receiver:        "c",
			receiverRegProv: nil,
			want:            nil,
			wantErr:         true,
		},
		{
			name: "compatibility result not found",
			req:  contract.GlobalContract{ContractID: 1},
			initiatorMap: map[string]*pb.RemoteParticipant{
				"a": {Name: "a"},
				"b": {Name: "b"},
				"c": {Name: "c"},
			},
			receiver:        "c",
			receiverRegProv: &ent.RegisteredProvider{ID: 1, ContractID: 2},
			want:            nil,
			wantErr:         true,
		},
		{
			name: "success",
			req:  contract.GlobalContract{ContractID: 1},
			initiatorMap: map[string]*pb.RemoteParticipant{
				"a": {Name: "a"},
				"b": {Name: "b"},
				"c": {Name: "c"},
			},
			receiver:        "c",
			receiverRegProv: &ent.RegisteredProvider{ID: 1, ContractID: 1},
			want: map[string]*pb.RemoteParticipant{
				"a": {Name: "a"},
				"b": {Name: "b"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := brokerServer{
				dbClient: enttest.NewClient(t),
			}
			got, err := s.getParticipantMapping(tt.req, tt.initiatorMap, tt.receiver, tt.receiverRegProv)
			if (err != nil) != tt.wantErr {
				t.Errorf("getParticipantMapping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getParticipantMapping() got = %v, want %v", got, tt.want)
			}
		})
	}
}
