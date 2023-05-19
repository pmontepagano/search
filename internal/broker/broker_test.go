package broker

import (
	"context"
	"fmt"
	"log"
	"testing"

	pb "github.com/clpombo/search/gen/go/search/v1"
	"github.com/clpombo/search/internal/contract"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestBrokerChannel_Request(t *testing.T) {
	// TODO: change this test and mock provider?
	b := NewBrokerServer("file:ent?mode=memory&_fk=1")
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
		t.Error("Received error from broker.")
	}
	if brokerresult.Result != pb.BrokerChannelResponse_RESULT_ACK {
		t.Error("Non ACK return code from broker")
	}

	b.Stop()
}

func TestGetParticipantMapping(t *testing.T) {
	b := NewBrokerServer("file:ent?mode=memory&_fk=1")
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
		t.Error("error getting registered provider")
	}

	b.SetCompatFunc(mockTestGetParticipantMappingContractCompatChecker)
	// test initiator mapping
	_, err = client.BrokerChannel(ctx, &pb.BrokerChannelRequest{
		Contract: &pb.GlobalContract{
			Contract:      dummyReqContract,
			Format:        pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA,
			InitiatorName: "FooBar",
		},
		PresetParticipants: map[string]*pb.RemoteParticipant{
			"FooBar": {
				Url:   "foobar_fake_url",
				AppId: "foobar_fake_appid",
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
		t.Error("error running getBestCandidate")
	}
	mapping, err := b.getParticipantMapping(reqContract, initiatorMapping, "0", "FooBar", provider)
	if err != nil {
		t.Error("error getting participant mapping")
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
