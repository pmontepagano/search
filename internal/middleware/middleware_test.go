package middleware

import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/clpombo/search/api"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)


var mw = NewMiddlewareServer("broker", 7777)


func init() {
	go mw.StartMiddlewareServer(4444, 5555, false, "", "")
}
func TestRegisterChannel(t *testing.T) {
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

}