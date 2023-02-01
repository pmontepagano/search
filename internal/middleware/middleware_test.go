package middleware

import (
	"context"
	"io"
	"log"
	"sync"
	"testing"
	"time"

	pb "github.com/clpombo/search/gen/go/search/v1"
	"github.com/clpombo/search/internal/broker"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// Start Middleware that listens on localhost and then send to it a
// dummy RegisterChannel RPC with a dummy GlobalContract
func TestRegisterChannel(t *testing.T) {
	mw := NewMiddlewareServer("broker", 7777)
	var wg sync.WaitGroup
	mw.StartMiddlewareServer(&wg, "localhost", 4444, "localhost", 5555, false, "", "")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial("localhost:5555", opts...)
	if err != nil {
		t.Error("Could not contact local private middleware server.")
	}
	defer conn.Close()
	client := pb.NewPrivateMiddlewareServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := pb.RegisterChannelRequest{
		RequirementsContract: &pb.Contract{
			Contract:           "hola",
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
	wg.Wait()
}

// Start a Broker and two Middleware servers. One of the middleware servers shall have a
// dummy provider registered, and the other will have an initiator app requesting a channel
// We should see brokering happen and message exchange between apps
func TestPingPong(t *testing.T) {
	// start broker
	bs := broker.NewBrokerServer()
	go bs.StartServer("localhost", 7777, false, "", "")

	var wg sync.WaitGroup
	// start provider middleware
	provMw := NewMiddlewareServer("localhost", 7777)
	provMw.StartMiddlewareServer(&wg, "localhost", 4444, "localhost", 5555, false, "", "")

	// start client middleware
	clientMw := NewMiddlewareServer("localhost", 7777)
	clientMw.StartMiddlewareServer(&wg, "localhost", 8888, "localhost", 9999, false, "", "")

	// common grpc.DialOption
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	// register dummy provider app and keep waiting for a notification
	go func() {
		// this function is for provider app

		// connect to provider middleware
		conn, err := grpc.Dial("localhost:5555", opts...)
		if err != nil {
			t.Error("Could not contact local private middleware server.")
		}

		client := pb.NewPrivateMiddlewareServiceClient(conn)

		// register dummy app with provider middleware
		req := pb.RegisterAppRequest{
			ProviderContract: &pb.Contract{
				Contract:           "dummy provider contract",
				RemoteParticipants: []string{"self", "p1"},
			},
		}

		stream, err := client.RegisterApp(context.Background(), &req)
		if err != nil {
			t.Error("Could not Register App")
		}
		ack, err := stream.Recv()
		if err != nil || ack.GetAppId() == "" {
			t.Error("Could not receive ACK from RegisterApp")
		}

		// wait on RegisterAppResponse stream to await for new channel (once only for this test)
		new, err := stream.Recv()
		if err == io.EOF {
			t.Error("Broker unexpectedly ended connection with provider")
		}
		if err != nil {
			t.Errorf("Error receiving notification from RegisterApp: %v", err)
		}
		channelID := new.GetNotification().GetChannelId()
		log.Printf("[PROVIDER] - Received Notification. ChannelID: %s", channelID)

		// reply "ping!" messages with "pong!" until we receive a different message, then exit
		go func(channelID string, conn *grpc.ClientConn) {
			defer conn.Close()
			defer provMw.Stop()
			for {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				res, err := client.AppRecv(ctx, &pb.AppRecvRequest{
					ChannelId:   channelID,
					Participant: "p1",
				})
				if err != nil {
					t.Errorf("[PROVIDER] - Error reading AppRecv. Error: %v", err)
				}
				log.Printf("[PROVIDER] - Received message from p1: %s", res.Content.GetBody())
				if string(res.Content.GetBody()) == "ping!" {
					ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					client.AppSend(ctx, &pb.ApplicationMessageOut{
						ChannelId: channelID,
						Recipient: "p1",
						Content: &pb.MessageContent{
							Body: []byte("pong!"),
						},
					})
				} else {
					log.Printf("[PROVIDER] - Exiting...")
					break
				}
			}
		}(channelID, conn)

	}()

	// connect to client middleware and register channel
	conn, err := grpc.Dial("localhost:9999", opts...)
	if err != nil {
		t.Error("Could not contact local private middleware server.")
	}
	defer conn.Close()
	client := pb.NewPrivateMiddlewareServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := pb.RegisterChannelRequest{
		RequirementsContract: &pb.Contract{
			Contract:           "client example requirement contract",
			RemoteParticipants: []string{"self", "p2"},
		},
	}
	regResult, err := client.RegisterChannel(ctx, &req)
	if err != nil {
		t.Error("Received error from RegisterChannel")
	}

	// AppSend to p2
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = client.AppSend(ctx, &pb.ApplicationMessageOut{
		ChannelId: regResult.ChannelId,
		Recipient: "p2",
		Content:   &pb.MessageContent{Body: []byte("ping!")},
	})
	if err != nil {
		t.Error(err)
	}

	// receive echo from p2
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.AppRecv(ctx, &pb.AppRecvRequest{
		ChannelId:   regResult.ChannelId,
		Participant: "p2",
	})
	if err != nil {
		t.Error("Could not receive message from p2")
	}
	log.Printf("Received message from p2: %s", resp.Content)

	// AppSend goodbye to p2 so that it exits
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = client.AppSend(ctx, &pb.ApplicationMessageOut{
		ChannelId: regResult.ChannelId,
		Recipient: "p2",
		Content:   &pb.MessageContent{Body: []byte("goodbye!")},
	})

	clientMw.Stop()
	wg.Wait()
}

func TestCircle(t *testing.T) {
	brokerPort, p1Port, p2Port, p3Port, initiatorPort := 20000, 20001, 20003, 20005, 20007

	// start broker
	bs := broker.NewBrokerServer()
	go bs.StartServer("localhost", brokerPort, false, "", "")

	var wg sync.WaitGroup
	// start middlewares
	p1Mw := NewMiddlewareServer("localhost", brokerPort)
	p1Mw.StartMiddlewareServer(&wg, "localhost", p1Port, "localhost", p1Port+1, false, "", "")
	p2Mw := NewMiddlewareServer("localhost", brokerPort)
	p2Mw.StartMiddlewareServer(&wg, "localhost", p2Port, "localhost", p2Port+1, false, "", "")
	p3Mw := NewMiddlewareServer("localhost", brokerPort)
	p3Mw.StartMiddlewareServer(&wg, "localhost", p3Port, "localhost", p3Port+1, false, "", "")
	initiatorMw := NewMiddlewareServer("localhost", brokerPort)
	initiatorMw.StartMiddlewareServer(&wg, "localhost", initiatorPort, "localhost", initiatorPort+1, false, "", "")

	// common grpc.DialOption
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	// launch 3 provider apps that simply pass the message to next member adding their name...?
	for _, mw := range []*MiddlewareServer{p1Mw, p2Mw, p3Mw} {
		go func(mw *MiddlewareServer) {
			// this function is for provider app

			// connect to provider middleware
			conn, err := grpc.Dial(mw.PrivateURL, opts...)
			if err != nil {
				t.Error("Could not contact local private middleware server.")
			}
			client := pb.NewPrivateMiddlewareServiceClient(conn)

			// register dummy app with provider middleware
			req := pb.RegisterAppRequest{
				ProviderContract: &pb.Contract{
					Contract:           "dummy provider contract",
					RemoteParticipants: []string{"self", "sender", "receiver"},
				},
			}

			streamCtx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
			defer cancel()
			stream, err := client.RegisterApp(streamCtx, &req)
			if err != nil {
				t.Error("Could not Register App")
			}
			ack, err := stream.Recv()
			if err != nil || ack.GetAppId() == "" {
				t.Error("Could not receive ACK from RegisterApp")
			}
			appID := ack.GetAppId()

			// wait on RegisterAppResponse stream to await for new channel (once only for this test)
			new, err := stream.Recv()
			if err == io.EOF {
				t.Error("Broker unexpectedly ended connection with provider")
			}
			if err != nil {
				t.Errorf("Error receiving notification from RegisterApp: %v", err)
			}
			channelID := new.GetNotification().GetChannelId()
			log.Printf("[PROVIDER %s] - Received Notification. ChannelID: %s", appID, channelID)

			// await message from sender, then add a word to the message and relay it to receiver
			go func(channelID string, conn *grpc.ClientConn) {
				defer conn.Close()
				defer mw.Stop()
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				res, err := client.AppRecv(ctx, &pb.AppRecvRequest{
					ChannelId:   channelID,
					Participant: "sender",
				})
				if err != nil {
					t.Errorf("[PROVIDER] - Error reading AppRecv. Error: %v", err)
				}
				log.Printf("[PROVIDER] - Received message from sender: %s", res.Content.GetBody())
				msg := string(res.Content.GetBody())
				msg = msg + " dummy"
				ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				client.AppSend(ctx, &pb.ApplicationMessageOut{
					ChannelId: channelID,
					Recipient: "receiver",
					Content: &pb.MessageContent{
						Body: []byte(msg),
					},
				})
			}(channelID, conn)

		}(mw)
	}

	// wait so that providers get to register with broker
	time.Sleep(2 * time.Second)

	// connect to initiator's middleware and register channel
	conn, err := grpc.Dial(initiatorMw.PrivateURL, opts...)
	if err != nil {
		t.Error("Could not contact local private middleware server.")
	}
	defer conn.Close()
	client := pb.NewPrivateMiddlewareServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := pb.RegisterChannelRequest{
		RequirementsContract: &pb.Contract{
			Contract:           "send hello to r1, and later receive mesage from r3",
			RemoteParticipants: []string{"self", "r1_special", "r2_special", "r3_special"},
		},
	}
	regResult, err := client.RegisterChannel(ctx, &req)
	if err != nil {
		t.Error("Received error from RegisterChannel")
	}

	// AppSend to r1
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = client.AppSend(ctx, &pb.ApplicationMessageOut{
		ChannelId: regResult.ChannelId,
		Recipient: "r1_special",
		Content:   &pb.MessageContent{Body: []byte("hola")},
	})

	// receive from r3
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.AppRecv(ctx, &pb.AppRecvRequest{
		ChannelId:   regResult.ChannelId,
		Participant: "r3_special",
	})
	if err != nil {
		t.Error("Could not receive message from r3")
	}
	log.Printf("Received message from r3: %s", resp.Content)

	initiatorMw.Stop()
	wg.Wait()
}