package middleware

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	pb "github.com/clpombo/search/api"
	"github.com/clpombo/search/internal/broker"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// Start Middleware that listens on localhost and then send to it a
// dummy RegisterChannel RPC with a dummy GlobalContract
func TestRegisterChannel(t *testing.T) {
	mw := NewMiddlewareServer("broker", 7777)
	go mw.StartMiddlewareServer("localhost", 4444, "localhost", 5555, false, "", "")

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

	// This checks internal state of the Middleware server. Probably not good practice.
	// check that the contract is properly saved inside the MiddleWare Server
	// in its "unbrokered" channels list.
	schan := mw.unBrokeredChannels[regResult.ChannelId]
	if schan.Contract.GetContract() != "hola" {
		t.Error("Contract from channel different from original")
	}

	// stop middleware to free-up port and resources after test run
	mw.Stop()
}

// Start a Broker and two Middleware servers. One of the middleware servers shall have a
// dummy provider registered, and the other will have an initiator app requesting a channel
// We should see brokering happen and message exchange between apps
func Test1(t *testing.T) {
	// start broker
	bs := broker.NewBrokerServer()
	go bs.StartServer("localhost", 7777, false, "", "")

	// start provider middleware
	provMw := NewMiddlewareServer("localhost", 7777)
	go provMw.StartMiddlewareServer("localhost", 4444, "localhost", 5555, false, "", "")

	// start client middleware
	clientMw := NewMiddlewareServer("localhost", 7777)
	go clientMw.StartMiddlewareServer("localhost", 8888, "localhost", 9999, false, "", "")

	// common grpc.DialOption
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	// register dummy provider app and keep waiting for a notification
	go func() {
		// connect to provider middleware
		conn, err := grpc.Dial("localhost:5555", opts...)
		if err != nil {
			t.Error("Could not contact local private middleware server.")
		}
		defer conn.Close()
		client := pb.NewPrivateMiddlewareClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req := pb.RegisterAppRequest{
			ProviderContract: &pb.Contract{
				Contract: "dummy",
				RemoteParticipants: []string{"self", "p1"},
			},
		}
		stream, err := client.RegisterApp(ctx, &req)
		if err != nil {
			t.Error("Could not Register App")
		}
		_, err = stream.Recv()
		if err != nil {
			t.Error("Could not receive ACK from RegisterApp")
		}

		for {
			new, err := stream.Recv()
			if err == io.EOF {
				t.Error("Broker unexpectedly ended connection with provider")
			}
			if err != nil {
				t.Error("Error receiving notification from RegisterApp")
			}
			log.Printf(new.GetNotification().ChannelId)
			// TODO: handle new channel
		}

	}()

	// wait a couple of seconds so that provider gets to register with broker
	time.Sleep(2 * time.Second)


	// connect to client middleware and register channel
	conn, err := grpc.Dial("localhost:9999", opts...)
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
			RemoteParticipants: []string{"self", "p2"},
		},
	}
	regResult, err := client.RegisterChannel(ctx, &req)
	if err != nil {
		t.Error("Received error from RegisterChannel")
	}

	// AppSend to p2
	_, err = client.AppSend(ctx, &pb.ApplicationMessageOut{
		ChannelId: regResult.ChannelId,
		Recipient: "p2",
		Content: &pb.MessageContent{Body: []byte("hello world")},
	})

	time.Sleep(2 * time.Second)


}