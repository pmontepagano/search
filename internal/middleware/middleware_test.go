package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pmontepagano/search/contract"
	pb "github.com/pmontepagano/search/gen/go/search/v1"
	"github.com/pmontepagano/search/internal/broker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Start Middleware that listens on localhost and then send to it a
// dummy RegisterChannel RPC with a dummy GlobalContract
func TestRegisterChannel(t *testing.T) {
	mw := NewMiddlewareServer("broker", 7777)
	var wg sync.WaitGroup
	mw.StartMiddlewareServer(&wg, "localhost", 4444, "localhost", 5555, false, "", "")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial("localhost:5555", opts...)
	if err != nil {
		t.Error("Could not contact local private middleware server.")
	}
	defer conn.Close()
	client := pb.NewPrivateMiddlewareServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dummyContract := []byte(`--
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
	`)

	req := pb.RegisterChannelRequest{
		RequirementsContract: &pb.GlobalContract{
			Contract:      dummyContract,
			InitiatorName: "0",
			Format:        pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA,
		},
	}
	regResult, err := client.RegisterChannel(ctx, &req)
	if err != nil {
		t.Errorf("Received error from RegisterChannel: %v", err)
	}
	_, err = uuid.Parse(regResult.ChannelId)
	if err != nil {
		t.Error("Received a non UUID ChannelID from RegisterChannel")
	}

	// This checks internal state of the Middleware server. Probably not good practice.
	// check that the contract is properly saved inside the MiddleWare Server
	// in its "unbrokered" channels list.
	schan := mw.unBrokeredChannels[regResult.ChannelId]
	if !bytes.Equal(schan.ContractPB.GetContract(), dummyContract) {
		t.Error("Contract from channel different from original")
	}

	// stop middleware to free-up port and resources after test run
	mw.Stop()
	wg.Wait()
}

/*
func TestCircle(t *testing.T) {
	brokerPort, p1Port, p2Port, p3Port, initiatorPort := 20000, 20001, 20003, 20005, 20007

	tmpDir := t.TempDir()
	bs := broker.NewBrokerServer(fmt.Sprintf("%s/testcircle-%s.db", tmpDir, time.Now().Format("2006-01-02T15:04:05")))
	t.Cleanup(bs.Stop)
	bs.SetCompatFunc(circleContractCompatChecker)
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
	defer p1Mw.Stop()
	defer p2Mw.Stop()
	defer p3Mw.Stop()
	defer initiatorMw.Stop()

	// common grpc.DialOption
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	// launch 3 provider apps that simply pass the message to next member adding their name...?
	for idx, mw := range []*MiddlewareServer{p1Mw, p2Mw, p3Mw} {
		go func(mw *MiddlewareServer, idx int) {
			// this function is for provider app

			// connect to provider middleware
			conn, err := grpc.Dial(mw.PrivateURL, opts...)
			if err != nil {
				t.Error("Could not contact local private middleware server.")
			}
			client := pb.NewPrivateMiddlewareServiceClient(conn)

			// Have a slightly different provider contract for each Service Provider to be able
			// to differentiate them in the compatibility check function.
			providerContract := fmt.Sprintf(`
			.outputs msg_passer_%v
			.state graph
			q0 sender_%v ? word q1
			q1 receiver ! word q0
			.marking q0
			.end
			`, idx, idx)
			// register dummy app with provider middleware
			req := pb.RegisterAppRequest{
				ProviderContract: &pb.LocalContract{
					Contract: []byte(providerContract),
					Format:   pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA,
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
			defer conn.Close()
			defer mw.Stop()
			res, err := client.AppRecv(context.Background(), &pb.AppRecvRequest{
				ChannelId:   channelID,
				Participant: "sender",
			})
			if err != nil {
				t.Errorf("[PROVIDER] - Error reading AppRecv. Error: %v", err)
			}
			log.Printf("[PROVIDER] - Received message from sender: %s", res.Message.GetBody())
			msg := string(res.Message.GetBody())
			msg = msg + " dummy"
			appSendResp, err := client.AppSend(context.Background(), &pb.AppSendRequest{
				ChannelId: channelID,
				Recipient: "receiver",
				Message: &pb.AppMessage{
					Body: []byte(msg),
				},
			})
			if err != nil || appSendResp.Result != pb.Result_OK {
				t.Errorf("[PROVIDER] - Error sending message to receiver. Error: %v", err)
			}

		}(mw, idx)
	}

	// wait so that providers get to register with broker
	time.Sleep(100 * time.Millisecond)

	// connect to initiator's middleware and register channel
	conn, err := grpc.Dial(initiatorMw.PrivateURL, opts...)
	if err != nil {
		t.Error("Could not contact local private middleware server.")
	}
	defer conn.Close()
	client := pb.NewPrivateMiddlewareServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := pb.RegisterChannelRequest{
		RequirementsContract: &pb.GlobalContract{
			Contract: []byte(`
			.outputs self
			.state graph
			q0 r1_special ! word q1
			q1 r3_special ? word q2
			.marking q0
			.end

			.outputs r1_special
			.state graph
			q0 self ? word q1
			q1 r2_special ! word qf
			.marking q0
			.end

			.outputs r2_special
			.state graph
			q0 r1_special ? word q1
			q1 r3_special ! word qf
			.marking q0
			.end

			.outputs r3_special
			.state graph
			q0 r2_special ? word q1
			q1 self ! word qf
			.marking q0
			.end
			`),
			InitiatorName: "self",
			Format:        pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA,
		},
	}
	regResult, err := client.RegisterChannel(ctx, &req)
	if err != nil {
		t.Error("Received error from RegisterChannel")
	}

	// AppSend to r1
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	sendRespR1, err := client.AppSend(ctx, &pb.AppSendRequest{
		ChannelId: regResult.ChannelId,
		Recipient: "r1_special",
		Message:   &pb.AppMessage{Body: []byte("hola")},
	})
	if err != nil || sendRespR1.Result != pb.Result_OK {
		t.Error("Could not send message to r1")
	}

	// receive from r3
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	resp, err := client.AppRecv(ctx, &pb.AppRecvRequest{
		ChannelId:   regResult.ChannelId,
		Participant: "r3_special",
	})
	if err != nil {
		t.Error("Could not receive message from r3")
	}
	log.Printf("Received message from r3: %s", resp.Message)

	initiatorMw.Stop()
	wg.Wait()
}

// Mock function for checking contracts in TestCircle.
func circleContractCompatChecker(ctx context.Context, req contract.LocalContract, prov contract.LocalContract) (bool, map[string]string, error) {
	log.Printf("[TestCircle] - Checking req ID: %s, prov ID: %s, req participants: %v, prov participants: %v", req.GetContractID(), prov.GetContractID(), req.GetRemoteParticipantNames(), prov.GetRemoteParticipantNames())
	mapping := make(map[string]string)
	if req.GetContractID() == "c6c64de47a8ca6293d5bfc108b47c8982c6ebc178b7d6cb5bcd22581685b1ce7b512c62e78ea8badaff14df89075d06186da982ae0a0719b651174f7031bcb92" && prov.GetContractID() == "6f970689b5f60afc5c6d5b5244f6e4fe83cf3af1e8cc153fe1877a1d8b76c004206a6d264a3e2d4d83dd3277eb28da4e0a57ff5b38061c5b7f533771904a8425" {
		mapping["sender_0"] = "self"
		mapping["receiver"] = "r2_special"
	}
	if req.GetContractID() == "0b719e3ea36cac5dbbd0d3efdb36ade1de42ceaf267dc978c49018be7dc4aa51aa60133c756ca552cb378970220dd59d6c70a650eabb9f97f80f176027ef7c44" && prov.GetContractID() == "a5f8cb87e049fa7bb4e6c4342f0c023a47525837b0fe7eb461a3e58f40a624b66c1e8ea58391e6caef48aa1fedee9b857ce69aaf3bec2a30dec1a9164d385ce3" {
		mapping["sender_1"] = "r1_special"
		mapping["receiver"] = "r3_special"
	}
	if req.GetContractID() == "cbe93c0bf4e2e1f8340e8febdc1b7bea9290aec3cfa66ca616e1c815f9efa34a1ae8b762c3118703ad47f0cee341ced5a2949a07c42e2afecb088ed4c8852642" && prov.GetContractID() == "db02a17429e4f057337d1a4da9a668dfa50b7574948dc2bff9a8b6548f7500abdc5ad9bc1072ddbb34d144cd1a3469ba332d0a36786ce2f8ad88bd5e281237e6" {
		mapping["sender_2"] = "r2_special"
		mapping["receiver"] = "self"
	}
	_, ok := mapping["receiver"]
	if !ok {
		return false, nil, nil
	}
	return true, mapping, nil
}
*/

func TestPingPongFullExample(t *testing.T) {

	// TODO: I don't think we can accept GC format in the middleware. Because if the conversion
	// to FSA introduces new messages that are not present in the GC, then the programmer needs
	// to explicitly send those messages to the middleware!

	// In this section we'll create several entities that will interact in this example:
	// 1. Middleware for Ping (initiator). Runs private and public middleware servers in gorotines.
	// 2. goroutine for Ping.
	// 3. Middleware for Pong (provider). Runs private and public middleware servers in gorotines.
	// 4. goroutine for Pong.
	// 5. Broker. Runs in gorotine.

	brokerPort, pingPrivPort, pingPubPort, pongPrivPort, pongPubPort := 20000, 20001, 20002, 20003, 20004

	// start broker
	tmpDir := t.TempDir()
	bs := broker.NewBrokerServer(fmt.Sprintf("%s/testpingpongfullexample.db", tmpDir))
	t.Cleanup(bs.Stop)
	bs.SetCompatFunc(pingPongContractCompatChecker)
	go bs.StartServer("localhost", brokerPort, false, "", "")
	defer bs.Stop()

	var wg sync.WaitGroup
	// start middlewares
	pingMiddleware := NewMiddlewareServer("localhost", brokerPort)
	pingMiddleware.StartMiddlewareServer(&wg, "localhost", pingPubPort, "localhost", pingPrivPort, false, "", "")

	pongMiddleware := NewMiddlewareServer("localhost", brokerPort)
	pongMiddleware.StartMiddlewareServer(&wg, "localhost", pongPubPort, "localhost", pongPrivPort, false, "", "")

	defer pingMiddleware.Stop()
	defer pongMiddleware.Stop()

	pongRegistered := make(chan bool, 1) // used to signal that Pong has already registered with broker.
	exitPong := make(chan bool, 1)       // used to signal Pong to gracefully exit.
	go pongProgram(t, pongMiddleware.PrivateURL, pongRegistered, exitPong)
	pingProgram(t, pingMiddleware.PrivateURL, pongRegistered)

	// Signal pongProgram to exit.
	exitPong <- true

	// Signal both middlewares to exit.
	pingMiddleware.Stop()
	pongMiddleware.Stop()

	// Signal broker to exit.
	bs.Stop()
}

const pongContractFSA = `
.outputs Pong
.state graph
0 Other ? ping 1
1 Other ! pong 0
0 Other ? bye 2
2 Other ! bye 3
0 Other ? finished 3
.marking 0
.end

.outputs Other
.state graph
0 Other ! ping 1
0 Other ! bye 2
0 Other ! finished 3
2 Other ? bye 3
1 Other ? pong 0
.marking 0
.end
`

const pingContractFSA = `
.outputs Ping
.state graph
0 Pong ! ping 1
1 Pong ? pong 0
0 Pong ! bye 2
0 Pong ! finished 3
2 Pong ? bye 3
.marking 0
.end

.outputs Pong
.state graph
0 Ping ? ping 1
1 Ping ! pong 0
0 Ping ? bye 2
2 Ping ! bye 3
0 Ping ? finished 3
.marking 0
.end
`

// Mock function for checking contracts in TestPingPongFullExample.
func pingPongContractCompatChecker(ctx context.Context, req contract.LocalContract, prov contract.LocalContract) (bool, map[string]string, error) {
	log.Printf("Checking with pingPongContractCompatChecker...")
	mapping := map[string]string{
		"Other": "Ping",
		"Pong":  "Pong",
	}
	return true, mapping, nil
}

func pongProgram(t *testing.T, middlewareURL string, registeredNotify chan bool, exitPong chan bool) {
	// Auxiliary function for TestPingPongFullExample.

	// TODO: is it valid to send a single CFSM? The 'Other' machine is undefined in this example.
	// For now, I'll send both CFSMs, because otherwise the FSA parser fails.

	// Connect to the middleware and instantiate client.
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(middlewareURL, opts...)
	if err != nil {
		t.Errorf("Error in pongProgram connecting to middleware URL %s", middlewareURL)
	}
	defer conn.Close()
	client := pb.NewPrivateMiddlewareServiceClient(conn)

	// Register provider contract with registry.
	req := pb.RegisterAppRequest{
		ProviderContract: &pb.LocalContract{
			Contract: []byte(pongContractFSA),
			Format:   pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA,
		},
	}
	streamCtx, streamCtxCancel := context.WithCancel(context.Background())
	defer streamCtxCancel()
	stream, err := client.RegisterApp(streamCtx, &req)
	if err != nil {
		t.Error("Could not Register App")
	}
	ack, err := stream.Recv()
	if err != nil || ack.GetAppId() == "" {
		t.Error("Could not receive ACK from RegisterApp")
	}
	appID := ack.GetAppId()
	t.Logf("pongProgram finished registration. Got AppId %s", appID)
	registeredNotify <- true

	// wait on RegisterAppResponse stream to await for new channel
	type NewSessionNotification struct {
		regappResp *pb.RegisterAppResponse
		err        error
	}
	recvChan := make(chan NewSessionNotification)
	go func(stream pb.PrivateMiddlewareService_RegisterAppClient, recvChan chan NewSessionNotification) {
		// We make a goroutine to have a channel interface instead of a blocking Recv()
		// https://github.com/grpc/grpc-go/issues/465#issuecomment-179414474
		for {
			newResponse, err := stream.Recv()
			recvChan <- NewSessionNotification{
				regappResp: newResponse,
				err:        err,
			}
		}
	}(stream, recvChan)
	for mainloop := true; mainloop; {
		select {
		case <-exitPong:
			mainloop = false
			log.Printf("Exiting pongProgram...")
		case newSess := <-recvChan:
			err := newSess.err
			if err == io.EOF {
				t.Error("Broker unexpectedly ended connection with provider?")
				mainloop = false
				break
			}
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("Unknown error attempting to receive RegisterApp notification: %v", err)
					mainloop = false
					break
				}
				t.Errorf("Error receiving notification from RegisterApp: %v", st)
				mainloop = false
				break
			}

			channelID := newSess.regappResp.GetNotification().GetChannelId()
			log.Printf("[PROVIDER %s] - Received Notification. ChannelID: %s", appID, channelID)
			go func(channelID string, client pb.PrivateMiddlewareServiceClient) {
				// This is the actual program for pong (the part that implements the contract).

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				for loop := true; loop; {
					// Receive ping, finished or bye request.
					recvResponse, err := client.AppRecv(ctx, &pb.AppRecvRequest{
						ChannelId:   channelID,
						Participant: "Other",
					})
					if err != nil {
						t.Error("Failed to receive ping/bye/finished from Other.")
					}
					switch recvResponse.Message.Type {
					case "ping":
						// Send back pong response.
						sendResponse, err := client.AppSend(ctx, &pb.AppSendRequest{
							ChannelId: channelID,
							Recipient: "Other",
							Message: &pb.AppMessage{
								Body: recvResponse.Message.Body,
								Type: "pong",
							},
						})
						if err != nil || sendResponse.Result != pb.Result_OK {
							t.Errorf("Failed to AppSend")
						}
					case "bye":
						// Send back bye response and break out of the loop.
						sendResponse, err := client.AppSend(ctx, &pb.AppSendRequest{
							ChannelId: channelID,
							Recipient: "Other",
							Message: &pb.AppMessage{
								Type: "bye",
								Body: []byte("exiting..."),
							},
						})
						if err != nil || sendResponse.Result != pb.Result_OK {
							t.Errorf("Failed to AppSend")
						}
						loop = false
					case "finished":
						loop = false
					default:
						t.Errorf("Received invalid message of type %v", recvResponse.Message.Type)
					}
				}

			}(channelID, client)
		}

	}

}

func pingProgram(t *testing.T, middlewareURL string, registeredPong chan bool) {
	// Auxiliary function for TestPingPongFullExample.

	// Connect to the middleware and instantiate client.
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(middlewareURL, opts...)
	if err != nil {
		t.Errorf("Error in pingProgram connecting to middleware URL %s", middlewareURL)
	}
	defer conn.Close()
	client := pb.NewPrivateMiddlewareServiceClient(conn)

	// Register channel and obtain channelID for the Channel.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := pb.RegisterChannelRequest{
		RequirementsContract: &pb.GlobalContract{
			Contract:      []byte(pingContractFSA),
			InitiatorName: "Ping",
			Format:        pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA,
		},
	}
	regResult, err := client.RegisterChannel(ctx, &req)
	if err != nil {
		t.Error("Received error from RegisterChannel")
	}
	channelID := regResult.ChannelId
	log.Printf("[pingProgram]: Obtained channel with ID: %s", channelID)

	// We need to wait until Pong has finished registering with the broker.
	<-registeredPong

	// TODO: keep using the same ctx? Has time elapsed on this one?
	// Send a ping message.
	sendResponse, err := client.AppSend(ctx, &pb.AppSendRequest{
		ChannelId: channelID,
		Recipient: "Pong",
		Message: &pb.AppMessage{
			Body: []byte("hello"), // we send whatever content and expect it reflected back to us.
			Type: "ping",
		},
	})
	if err != nil || sendResponse.Result != pb.Result_OK {
		t.Errorf("Failed to AppSend")
	}

	// Receive a pong message.
	recvResponse, err := client.AppRecv(ctx, &pb.AppRecvRequest{
		ChannelId:   channelID,
		Participant: "Pong",
	})
	if err != nil {
		t.Error("Failed to AppRecv")
	}
	if recvResponse.Message.Type != "pong" {
		t.Errorf("Received unexpected message of type %v", recvResponse.Message.Type)
	}
	if !bytes.Equal(recvResponse.Message.GetBody(), []byte("hello")) {
		t.Error("Received different pong body that what we originally sent.")
	}

	// Send bye message.
	sendResponse, err = client.AppSend(ctx, &pb.AppSendRequest{
		ChannelId: channelID,
		Recipient: "Pong",
		Message: &pb.AppMessage{
			Body: []byte("Hasta la vista, baby."), // body is unimportant.
			Type: "bye",
		},
	})
	if err != nil || sendResponse.Result != pb.Result_OK {
		t.Errorf("Failed to AppSend")
	}

	// Recieve bye response.
	recvResponse, err = client.AppRecv(ctx, &pb.AppRecvRequest{
		ChannelId:   channelID,
		Participant: "Pong",
	})
	if err != nil {
		t.Error("Failed to AppRecv")
	}
	if recvResponse.Message.Type != "bye" {
		t.Errorf("Received unexpected message of type %v", recvResponse.Message.Type)
	}

	// Tell pong program to stop
}
