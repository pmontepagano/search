package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/pmontepagano/search/gen/go/search/v1"
)

var middlewareURL = flag.String("middleware-url", "middleware-payments:11000", "The URL for the middleware")

const ppsContract = `
.outputs PPS
.state graph
q0 ClientApp ? CardDetailsWithTotalAmount q1
q1 ClientApp ! PaymentNonce q2
q2 Srv ? RequestChargeWithNonce q3
q3 Srv ! ChargeOK q4
q3 Srv ! ChargeFail q5
.marking q0
.end
`

func main() {
	var logger = log.New(os.Stderr, fmt.Sprintf("[PPS] - "), log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*middlewareURL, opts...)
	if err != nil {
		logger.Fatalf("Error connecting to middleware URL %s", *middlewareURL)
	}
	defer conn.Close()
	client := pb.NewPrivateMiddlewareServiceClient(conn)

	// Register provider contract with registry.
	req := pb.RegisterAppRequest{
		ProviderContract: &pb.LocalContract{
			Contract: []byte(ppsContract),
			Format:   pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA,
		},
	}
	streamCtx, streamCtxCancel := context.WithCancel(context.Background())
	defer streamCtxCancel()
	stream, err := client.RegisterApp(streamCtx, &req)
	if err != nil {
		logger.Fatal("Could not Register App")
	}
	ack, err := stream.Recv()
	if err != nil || ack.GetAppId() == "" {
		logger.Fatal("Could not receive ACK from RegisterApp")
	}
	appID := ack.GetAppId()
	logger.Printf("Finished registration. Got AppId %s", appID)

	// wait on RegisterAppResponse stream to await for new channel
	type NewSessionNotification struct {
		regappResp *pb.RegisterAppResponse
		err        error
	}
	recvChan := make(chan NewSessionNotification)
	go func(stream pb.PrivateMiddlewareService_RegisterAppClient, recvChan chan NewSessionNotification) {
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
				logger.Fatalf("Error receiving RegisterApp notification: %v", err)
			}
		}
	}(stream, recvChan)
	for mainloop := true; mainloop; {
		select {
		// case <-exitPong:
		// 	mainloop = false
		// 	log.Printf("Exiting pongProgram...")
		case newSess := <-recvChan:
			err := newSess.err
			if err == io.EOF {
				logger.Printf("Broker unexpectedly ended connection with provider?")
				mainloop = false
				break
			}
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					logger.Fatalf("Unknown error attempting to receive RegisterApp notification: %v", err)
				}
				logger.Fatalf("Error receiving notification from RegisterApp: %v", st)
			}

			channelID := newSess.regappResp.GetNotification().GetChannelId()
			logger.Printf("Received Notification. ChannelID: %s", channelID)
			go func(channelID string, client pb.PrivateMiddlewareServiceClient) {
				// This is the actual program for PPS (the part that implements the contract).
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				for loop := true; loop; {
					// Receive ping, finished or bye request.
					logger.Print("Waiting for ping/bye/finished...")
					recvResponse, err := client.AppRecv(ctx, &pb.AppRecvRequest{
						ChannelId:   channelID,
						Participant: "Other",
					})
					if err != nil {
						logger.Fatal("Failed to receive ping/bye/finished from Other.")
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
							logger.Fatal("Failed to AppSend")
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
							logger.Fatal("Failed to AppSend")
						}
						loop = false
					case "finished":
						loop = false
					default:
						logger.Fatalf("Received invalid message of type %v", recvResponse.Message.Type)
					}
				}

			}(channelID, client)
		}
	}

}
