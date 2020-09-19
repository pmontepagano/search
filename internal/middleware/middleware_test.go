package middleware

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/clpombo/search/api"
	"github.com/clpombo/search/internal/broker"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)


func TestRegisterChannel(t *testing.T) {
	mw := NewMiddlewareServer("broker", 7777)
	go mw.StartMiddlewareServer("localhost", 4444, "localhost", 5555, false, "", "")
	log.Println("hello world")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial("localhost:5555", opts...)
	if err != nil {
		t.Error("Could not contact local private middleware server.")
	}
	defer conn.Close()
	client := pb.NewPrivateMiddlewareClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := pb.RegisterChannelRequest{
		RequirementsContract: &pb.Contract{
			Contract: "hola",
			RemoteParticipants: []string{"self", "p1", "p2"},
		},
	}
	regResult, err := client.RegisterChannel(ctx, &req)
	if err != nil {
		t.Error("Received error from RegisterChannel")
	}
	_, err = uuid.Parse(regResult.ChannelId)
	if err != nil {
		t.Error("Received a non UUID ChannelID from RegisterChannel")
	}

	schan := mw.unBrokeredChannels[regResult.ChannelId]
	if schan.Contract.GetContract() != "hola" {
		t.Error("Contract from channel different from original")
	}
	mw.Stop()
}

func TestAppSend(t *testing.T) {
	bs := broker.NewBrokerServer("")
	go bs.StartServer("localhost", 7777, false, "", "")

}