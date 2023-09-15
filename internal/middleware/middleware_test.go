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
	"go.uber.org/goleak"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Start Middleware that listens on localhost and then send to it a
// dummy RegisterChannel RPC with a dummy GlobalContract
func TestRegisterChannel(t *testing.T) {
	defer goleak.VerifyNone(t)
	mw := NewMiddlewareServer("broker:7777")
	var wg sync.WaitGroup
	mw.StartMiddlewareServer(&wg, "localhost:4444", "localhost:5555", false, "", "", nil)

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

// Test scenarios where the Broker cannot find any provider that satisfies the requirements.
func TestNoCompatibleProviders(t *testing.T) {
	type testcase struct {
		name          string
		req           []byte // requirement contract
		initiatorName string
		action        func(*testing.T, pb.PrivateMiddlewareServiceClient, context.Context, string) // action to trigger search for providers
	}
	tests := []testcase{
		{
			name:          "onlysendonce",
			initiatorName: "0",
			req: []byte(`--
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
			action: func(t *testing.T, client pb.PrivateMiddlewareServiceClient, ctx context.Context, channelID string) {
				// Send a message to FooBar to trigger brokerage of channel. This send will succeed.
				appSendResponse, err := client.AppSend(ctx, &pb.AppSendRequest{
					ChannelId: channelID,
					Recipient: "FooBar",
					Message: &pb.AppMessage{
						Type: "hello",
						Body: []byte("hello")},
				})
				if err != nil {
					t.Errorf("Received error when sending message to FooBar")
				}
				if appSendResponse.Result != pb.AppSendResponse_RESULT_OK {
					t.Errorf("Received non-OK result when sending message to FooBar")
				}

				// Close the channel. This should fail because the brokerage failed.
				closeResponse, err := client.CloseChannel(ctx, &pb.CloseChannelRequest{ChannelId: channelID})
				if err == nil {
					t.Errorf("Received no error when closing channel")
				}
				if closeResponse != nil && closeResponse.Result == pb.CloseChannelResponse_RESULT_CLOSED {
					t.Errorf("Received proper close when no provider is compatible...")
				}
			},
		},
		{
			name:          "failtorecv",
			initiatorName: "Bar",
			req: []byte(`--
				.outputs Foo
				.state graph
				q0 1 ! hello q0
				.marking q0
				.end

				.outputs Bar
				.state graph
				q0 0 ? hello q0
				.marking q0
				.end
			`),
			action: func(t *testing.T, client pb.PrivateMiddlewareServiceClient, ctx context.Context, channelID string) {
				// Try to receive "hello" from Foo. This should fail because brokerage should fail.
				_, err := client.AppRecv(ctx, &pb.AppRecvRequest{
					ChannelId:   channelID,
					Participant: "Foo",
				})
				if err == nil {
					t.Error("AppRecv was expected to fail because there are no compatible providers.")
				} else {
					st := status.Convert(err)
					if st.Code() != codes.NotFound {
						t.Error("AppRecv was expected to fail with gRPC Unknown error code.")
					}
					details := st.Details()
					if len(details) != 1 {
						t.Error("AppRecv was expected to fail with one error detail.")
					}
					detail := details[0].(*errdetails.ErrorInfo)
					if detail.GetReason() != "CHANNEL_BROKERAGE_FAILED" {
						t.Errorf("AppRecv was expected to fail with CHANNEL_BROKERAGE_FAILED, but got %s", detail.GetReason())
					}
					if detail.Metadata["failed_participants"] != "Foo" {
						t.Errorf("Unexpected list of failed participants.")
					}

				}

				// Close the channel. This should fail because the brokerage failed.
				closeResponse, err := client.CloseChannel(ctx, &pb.CloseChannelRequest{ChannelId: channelID})
				if err == nil {
					t.Errorf("Received no error when closing channel")
				}
				if closeResponse != nil && closeResponse.Result == pb.CloseChannelResponse_RESULT_CLOSED {
					t.Errorf("Received proper close when no provider is compatible...")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			tmpDir := t.TempDir()
			bs := broker.NewBrokerServer(fmt.Sprintf("%s/testnocompatibleproviders-%s-%s.db", tmpDir, tt.name, time.Now().Format("2006-01-02T15:04:05")))
			t.Cleanup(bs.Stop)
			brokerStartedURL := make(chan string, 1)
			go bs.StartServer("localhost:", false, "", "", brokerStartedURL)

			var wgServiceClient sync.WaitGroup
			// start middleware for Service Client
			t.Log("waiting for broker...")
			brokerURL := <-brokerStartedURL // wait until the Broker has selected a free TCP port and started listening on it.
			serviceClientMw := NewMiddlewareServer(brokerURL)
			serviceClientMw.StartMiddlewareServer(&wgServiceClient, "localhost:", "localhost:", false, "", "", nil)
			t.Cleanup(serviceClientMw.Stop)

			// Connect to Service Client's middleware.
			var opts []grpc.DialOption
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
			opts = append(opts, grpc.WithBlock())
			conn, err := grpc.Dial(serviceClientMw.PrivateURL, opts...)
			if err != nil {
				t.Error("Could not contact local private middleware server.")
			}
			defer conn.Close()
			client := pb.NewPrivateMiddlewareServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			req := pb.RegisterChannelRequest{
				RequirementsContract: &pb.GlobalContract{
					Contract:      tt.req,
					InitiatorName: tt.initiatorName,
					Format:        pb.GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA,
				},
			}
			regResult, err := client.RegisterChannel(ctx, &req)
			if err != nil {
				t.Errorf("Received error from RegisterChannel: %v", err)
			}
			t.Logf("Received ChannelID: %s", regResult.ChannelId)

			tt.action(t, client, ctx, regResult.ChannelId)
		})
	}

}

func TestCircle(t *testing.T) {
	p1Port, p2Port, p3Port, initiatorPort := 20001, 20003, 20005, 20007

	tmpDir := t.TempDir()
	bs := broker.NewBrokerServer(fmt.Sprintf("%s/testcircle-%s.db", tmpDir, time.Now().Format("2006-01-02T15:04:05")))
	t.Cleanup(bs.Stop)
	bs.SetCompatFunc(circleContractCompatChecker)
	brokerStartedURL := make(chan string, 1)
	go bs.StartServer("localhost:", false, "", "", brokerStartedURL)
	brokerURL := <-brokerStartedURL // wait until the Broker has selected a free TCP port and started listening on it.

	var wgProviders sync.WaitGroup
	var wgInitiator sync.WaitGroup
	// start middlewares
	p1Mw := NewMiddlewareServer(brokerURL)
	p1Mw.StartMiddlewareServer(&wgProviders, fmt.Sprintf("localhost:%d", p1Port), fmt.Sprintf("localhost:%d", p1Port+1), false, "", "", nil)
	p2Mw := NewMiddlewareServer(brokerURL)
	p2Mw.StartMiddlewareServer(&wgProviders, fmt.Sprintf("localhost:%d", p2Port), fmt.Sprintf("localhost:%d", p2Port+1), false, "", "", nil)
	p3Mw := NewMiddlewareServer(brokerURL)
	p3Mw.StartMiddlewareServer(&wgProviders, fmt.Sprintf("localhost:%d", p3Port), fmt.Sprintf("localhost:%d", p3Port+1), false, "", "", nil)
	initiatorMw := NewMiddlewareServer(brokerURL)
	initiatorMw.StartMiddlewareServer(&wgInitiator, fmt.Sprintf("localhost:%d", initiatorPort), fmt.Sprintf("localhost:%d", initiatorPort+1), false, "", "", nil)
	defer initiatorMw.Stop()

	// common grpc.DialOption
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	// launch 3 provider apps that simply pass the message to next member adding their name...?
	for idx, mw := range []*MiddlewareServer{p1Mw, p2Mw, p3Mw} {
		go func(t *testing.T, mw *MiddlewareServer, idx int) {
			// this function is for provider app
			defer mw.Stop()
			// connect to provider middleware
			conn, err := grpc.Dial(mw.PrivateURL, opts...)
			if err != nil {
				t.Error("Could not contact local private middleware server.")
			}
			defer conn.Close()
			client := pb.NewPrivateMiddlewareServiceClient(conn)

			// Have a slightly different provider contract for each Service Provider to be able
			// to differentiate them in the compatibility check function.
			providerContract := fmt.Sprintf(`
			.outputs msg_passer_%v
			.state graph
			q0 sender ? word q1
			q1 receiver ! word q0
			.marking q0
			.end
			`, idx)
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
			log.Printf("[PROVIDER msg_passer_%v - %s] - Received Notification. ChannelID: %s", idx, appID, channelID)

			// await message from sender, then add a word to the message and relay it to receiver
			log.Printf("[PROVIDER msg_passer_%v - %s] - Awaiting message from sender. ChannelID: %s", idx, appID, channelID)
			res, err := client.AppRecv(context.Background(), &pb.AppRecvRequest{
				ChannelId:   channelID,
				Participant: "sender",
			})
			if err != nil {
				t.Errorf("[PROVIDER msg_passer_%v - %s] - Error reading AppRecv. Error: %v", idx, appID, err)
			}
			log.Printf("[PROVIDER msg_passer_%v - %s] - Received message from sender: %s", idx, appID, res.Message.GetBody())
			msg := string(res.Message.GetBody())
			msg = msg + " dummy"
			appSendResp, err := client.AppSend(context.Background(), &pb.AppSendRequest{
				ChannelId: channelID,
				Recipient: "receiver",
				Message: &pb.AppMessage{
					Body: []byte(msg),
				},
			})
			if err != nil || appSendResp.Result != pb.AppSendResponse_RESULT_OK {
				t.Errorf("[PROVIDER msg_passer_%v - %s] - Error sending message to receiver. Error: %v", idx, appID, err)
			}
			log.Printf("[PROVIDER msg_passer_%v - %s] - Sent message to receiver: %s", idx, appID, msg)
			log.Printf("[PROVIDER msg_passer_%v - %s] - Closing channel...", idx, appID)

			closeChannelResponse, err := client.CloseChannel(context.Background(), &pb.CloseChannelRequest{ChannelId: channelID})
			if err != nil {
				t.Errorf("error closing channel: %v", err)
			}
			if closeChannelResponse.Result != pb.CloseChannelResponse_RESULT_CLOSED {
				t.Error("channel was not closed")
			}
			log.Printf("[PROVIDER msg_passer_%v - %s] - Closed channel, exiting... ChannelID: %s", idx, appID, channelID)

		}(t, mw, idx)
	}

	// wait so that providers get to register with broker
	// TODO: fix this with a syncroization channel.
	t.Log("Waiting for providers to register with broker...")
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
	t.Log("Sending message to r1...")
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	sendRespR1, err := client.AppSend(ctx, &pb.AppSendRequest{
		ChannelId: regResult.ChannelId,
		Recipient: "r1_special",
		Message:   &pb.AppMessage{Body: []byte("hola")},
	})
	if err != nil || sendRespR1.Result != pb.AppSendResponse_RESULT_OK {
		t.Error("Could not send message to r1")
	}
	t.Log("Sent message to r1.")

	// receive from r3
	t.Log("Awaiting message from r3...")
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

	// time.Sleep(5 * time.Second)
	log.Printf("Waiting for providers to exit...")
	wgProviders.Wait()
	initiatorMw.Stop()
}

// Mock function for checking contracts in TestCircle.
func circleContractCompatChecker(ctx context.Context, req contract.LocalContract, prov contract.LocalContract) (bool, map[string]string, error) {
	log.Printf("[TestCircle] - Checking req ID: %s, prov ID: %s, req participants: %v, prov participants: %v", req.GetContractID(), prov.GetContractID(), req.GetRemoteParticipantNames(), prov.GetRemoteParticipantNames())
	mapping := make(map[string]string)
	if req.GetContractID() == "c6c64de47a8ca6293d5bfc108b47c8982c6ebc178b7d6cb5bcd22581685b1ce7b512c62e78ea8badaff14df89075d06186da982ae0a0719b651174f7031bcb92" && prov.GetContractID() == "bc15e8c84c7ef522602bafb1da221735198067c825dc5ee2d2af2ba5d3f5e87d5b0c3385b65fd49752d0c18025177db1cdc13d831909ae0b6a59aa3b84df9ee2" {
		mapping["sender"] = "self"
		mapping["receiver"] = "r2_special"
	}
	if req.GetContractID() == "0b719e3ea36cac5dbbd0d3efdb36ade1de42ceaf267dc978c49018be7dc4aa51aa60133c756ca552cb378970220dd59d6c70a650eabb9f97f80f176027ef7c44" && prov.GetContractID() == "af79e694939e2ea7fb836ebbe96c4aadb07177dd0646d9ab4b482741065dba7c3dd7556a03617d55201595f3a19d7f5de174eedfa7e8dd5981ffb5bdd76d007d" {
		mapping["sender"] = "r1_special"
		mapping["receiver"] = "r3_special"
	}
	if req.GetContractID() == "cbe93c0bf4e2e1f8340e8febdc1b7bea9290aec3cfa66ca616e1c815f9efa34a1ae8b762c3118703ad47f0cee341ced5a2949a07c42e2afecb088ed4c8852642" && prov.GetContractID() == "f83914e510d991a1da253701150c9cc51f9d85743d9d30bfecfe9fcfd24f1e39e87643689b23cf2ccfa08987bb4cdea90e30caa1eeb6ba5f1ab58eb094e97e72" {
		mapping["sender"] = "r2_special"
		mapping["receiver"] = "self"
	}
	_, ok := mapping["receiver"]
	if !ok {
		return false, nil, nil
	}
	return true, mapping, nil
}

func TestPingPongFullExample(t *testing.T) {
	// In this test we'll create several entities that will interact in this example:
	// 1. Middleware for Ping (initiator). Runs private and public middleware servers in gorotines.
	// 2. Middleware for Pong (provider). Runs private and public middleware servers in gorotines.
	// 3. goroutine for Pong.
	// 4. Broker. Runs in gorotine.

	brokerPort, pingPrivPort, pingPubPort, pongPrivPort, pongPubPort := 30000, 30001, 30002, 30003, 30004

	// start broker
	tmpDir := t.TempDir()
	bs := broker.NewBrokerServer(fmt.Sprintf("%s/testpingpongfullexample.db", tmpDir))
	t.Cleanup(bs.Stop)
	bs.SetCompatFunc(pingPongContractCompatChecker)
	go bs.StartServer(fmt.Sprintf("localhost:%d", brokerPort), false, "", "", nil)
	defer bs.Stop()

	// start middlewares
	var wg sync.WaitGroup
	pingMiddleware := NewMiddlewareServer(fmt.Sprintf("localhost:%d", brokerPort))
	pingMiddleware.StartMiddlewareServer(&wg, fmt.Sprintf("localhost:%d", pingPubPort), fmt.Sprintf("localhost:%d", pingPrivPort), false, "", "", nil)
	pongMiddleware := NewMiddlewareServer(fmt.Sprintf("localhost:%d", brokerPort))
	pongMiddleware.StartMiddlewareServer(&wg, fmt.Sprintf("localhost:%d", pongPubPort), fmt.Sprintf("localhost:%d", pongPrivPort), false, "", "", nil)
	defer pingMiddleware.Stop()
	defer pongMiddleware.Stop()

	pongRegistered := make(chan bool, 1) // used to signal that Pong has already registered with broker.
	exitPong := make(chan bool, 1)       // used to signal Pong to gracefully exit.
	go pongProgram(t, pongMiddleware.PrivateURL, pongRegistered, exitPong)
	pingProgram(t, pingMiddleware.PrivateURL, pongRegistered)

	// Signal pongProgram to exit.
	exitPong <- true
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

// Auxiliary function for TestPingPongFullExample.
func pongProgram(t *testing.T, middlewareURL string, registeredNotify chan bool, exitPong chan bool) {
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
	go func(t *testing.T, stream pb.PrivateMiddlewareService_RegisterAppClient, recvChan chan NewSessionNotification) {
		// We make a goroutine to have a channel interface instead of a blocking Recv()
		// https://github.com/grpc/grpc-go/issues/465#issuecomment-179414474
		// This goroutine simply waits for any new RegisterAppResponse in the stream and sends it
		// to the recvChan.
		for {
			newResponse, err := stream.Recv()
			recvChan <- NewSessionNotification{
				regappResp: newResponse,
				err:        err,
			}
			if err != nil {
				t.Errorf("Error receiving RegisterApp notification in pongProgram: %v", err)
				return
			}
		}
	}(t, stream, recvChan)
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
					t.Log("Waiting for ping/bye/finished...")
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
						if err != nil || sendResponse.Result != pb.AppSendResponse_RESULT_OK {
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
						if err != nil || sendResponse.Result != pb.AppSendResponse_RESULT_OK {
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
	if err != nil || sendResponse.Result != pb.AppSendResponse_RESULT_OK {
		t.Errorf("Failed to AppSend")
	}

	// Receive a pong message.
	t.Log("Sent ping, waiting for pong reply...")
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
	t.Log("Successfully received pong reply.")

	// Send bye message.
	sendResponse, err = client.AppSend(ctx, &pb.AppSendRequest{
		ChannelId: channelID,
		Recipient: "Pong",
		Message: &pb.AppMessage{
			Body: []byte("Hasta la vista, baby."), // body is unimportant.
			Type: "bye",
		},
	})
	if err != nil || sendResponse.Result != pb.AppSendResponse_RESULT_OK {
		t.Errorf("Failed to AppSend")
	}

	// Recieve bye response.
	t.Log("Successfully sent bye message. Waiting for response...")
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
	t.Log("Successfully received bye response.")

	// Tell pong program to stop
}
