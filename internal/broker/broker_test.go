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
	b := NewBrokerServer("")
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

	c := pb.Contract{
		Contract: "hola",
		RemoteParticipants: []string{"self", "p1"},
	}
	req := pb.BrokerChannelRequest{
		Contract: &c,
		PresetParticipants: map[string]*pb.RemoteParticipant{
			"self": &pb.RemoteParticipant{
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
	time.Sleep(5 * time.Second)	// FIXME
}
