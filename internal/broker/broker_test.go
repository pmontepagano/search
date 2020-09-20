package broker

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	pb "github.com/clpombo/search/api"
	"google.golang.org/grpc"
)

func TestBrokerChannel_Request(t *testing.T) {
	b := NewBrokerServer()
	go b.StartServer("localhost", 3333, false, "", "")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", 3333), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBrokerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// register dummy provider
	_, err = client.RegisterProvider(ctx, &pb.RegisterProviderRequest{
		Contract: &pb.Contract{
			Contract: "dummy",
			RemoteParticipants: []string{"self", "p0"},
		},
		Url: "fakeurl",
	})
	if err != nil {
		log.Fatalf("ERROR RegisterProvider: %v", err)
	}

	// ask for channel brokerage
	c := pb.Contract{
		Contract: "hola",
		RemoteParticipants: []string{"self", "p1"},
	}
	req := pb.BrokerChannelRequest{
		Contract: &c,
		PresetParticipants: map[string]*pb.RemoteParticipant{
			"self": {
				Url: "fake",
				AppId: "fake",
			},
		},
	}
	brokerresult, err := client.BrokerChannel(ctx, &req)
	if err != nil {
		t.Error("Received error from broker.")
	}
	if brokerresult.Result != pb.BrokerChannelResponse_ACK {
		t.Error("Non ACK return code from broker")
	}

	b.Stop()
}
