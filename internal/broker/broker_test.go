package broker

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	pb "github.com/clpombo/search/gen/go/search/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestBrokerChannel_Request(t *testing.T) {
	// TODO: change this test and mock provider?
	b := NewBrokerServer("file::memory:?cache=shared")
	go b.StartServer("localhost", 3333, false, "", "")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
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
		Contract: &pb.Contract{
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
			Format:           pb.ContractFormat_CONTRACT_FORMAT_FSA,
			LocalParticipant: "FooBar",
		},
		ProviderName: "FooBar",
		Url:          "fakeurl",
	})
	if err != nil {
		log.Fatalf("ERROR RegisterProvider: %v", err)
	}

	// ask for channel brokerage
	c := pb.Contract{
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
		Format:           pb.ContractFormat_CONTRACT_FORMAT_FSA,
		LocalParticipant: "0",
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
	b := NewBrokerServer("file::memory:?cache=shared")
	go b.StartServer("localhost", 3333, false, "", "")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", 3333), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBrokerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// register dummy provider
	registrationResponse, err := client.RegisterProvider(ctx, &pb.RegisterProviderRequest{
		Contract: &pb.Contract{
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
			Format:           pb.ContractFormat_CONTRACT_FORMAT_FSA,
			LocalParticipant: "0",
		},
		ProviderName: "0",
		Url:          "fakeurl",
	})
	if err != nil {
		log.Fatalf("ERROR RegisterProvider: %v", err)
	}

	// test initiator mapping
	initiatorMapping := map[string]*pb.RemoteParticipant{
		"self": {
			Url:   "initiator_fake_url",
			AppId: "initiator_fake_appid",
		},
		"other": {
			Url:   "fakeurl",
			AppId: registrationResponse.GetAppId(),
		},
	}
	mapping := b.getParticipantMapping(initiatorMapping, "other", "self")
	expected := map[string]*pb.RemoteParticipant{
		"0": {
			Url:   "fakeurl",
			AppId: registrationResponse.GetAppId(),
		},
		"FooBar": {
			Url:   "initiator_fake_url",
			AppId: "initiator_fake_appid",
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
