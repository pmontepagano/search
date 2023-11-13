package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
	stub := pb.NewPrivateMiddlewareServiceClient(conn)

	// Register provider contract with registry.
	req := pb.RegisterAppRequest{
		ProviderContract: &pb.LocalContract{
			Contract: []byte(ppsContract),
			Format:   pb.LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA,
		},
	}
	streamCtx, streamCtxCancel := context.WithCancel(context.Background())
	defer streamCtxCancel()
	stream, err := stub.RegisterApp(streamCtx, &req)
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
	for {
		newSess := <-recvChan
		err := newSess.err
		if err == io.EOF {
			logger.Printf("middleware unexpectedly ended connection!")
			break
		}
		if err != nil {
			logger.Fatalf("Error receiving notification from RegisterApp: %v", err)
		}

		channelID := newSess.regappResp.GetNotification().GetChannelId()
		logger.Printf("Received Notification. ChannelID: %s", channelID)
		go func(channelID string, client pb.PrivateMiddlewareServiceClient) {
			// This is the actual program for PPS (the part that implements the contract).
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Receive CardDetailsWithTotalAmount from ClientApp.
			logger.Print("Waiting for CardDetailsWithTotalAmount...")
			recvResponse, err := client.AppRecv(ctx, &pb.AppRecvRequest{
				ChannelId:   channelID,
				Participant: "ClientApp",
			})
			if err != nil {
				logger.Fatal("Failed to receive CardDetailsWithTotalAmount from ClientApp.")
			}
			if recvResponse.Message.Type != "CardDetailsWithTotalAmount" {
				logger.Fatalf("Received invalid message of type %v", recvResponse.Message.Type)
			}
			// We don't validate the card details or the amount, this is just a demo.
			type CardDetailsWithTotalAmount struct {
				CardNumber         string  `json:"card_number"`
				CardExpirationDate string  `json:"card_expirationdate"`
				CardCVV            string  `json:"card_cvv"`
				TotalAmount        float64 `json:"total_amount"`
			}
			var cardDetails CardDetailsWithTotalAmount
			err = json.Unmarshal(recvResponse.Message.Body, &cardDetails)
			if err != nil {
				logger.Fatalf("Error unmarshalling CardDetailsWithTotalAmount: %v", err)
			}
			logger.Printf("Received CardDetailsWithTotalAmount: %v", cardDetails)

			// Send PaymentNonce to ClientApp.
			logger.Print("Sending PaymentNonce to ClientApp...")
			sendResponse, err := client.AppSend(ctx, &pb.AppSendRequest{
				ChannelId: channelID,
				Recipient: "ClientApp",
				Message: &pb.AppMessage{
					Type: "PaymentNonce",
					Body: []byte("1234567890"),
				},
			})
			if err != nil {
				logger.Fatal("Failed to send PaymentNonce to ClientApp.")
			}
			if sendResponse.Result != pb.AppSendResponse_RESULT_OK {
				logger.Fatalf("Failed to send PaymentNonce to ClientApp. Result: %v", sendResponse.Result)
			}

			// Receive RequestChargeWithNonce from Srv.
			logger.Print("Waiting for RequestChargeWithNonce...")
			recvResponse, err = client.AppRecv(ctx, &pb.AppRecvRequest{
				ChannelId:   channelID,
				Participant: "Srv",
			})
			if err != nil {
				logger.Fatal("Failed to receive RequestChargeWithNonce from Srv.")
			}
			if recvResponse.Message.Type != "RequestChargeWithNonce" {
				logger.Fatalf("Received invalid message of type %v", recvResponse.Message.Type)
			}
			type RequestChargeWithNonce struct {
				Nonce  string  `json:"nonce"`
				Amount float64 `json:"amount"`
			}
			var requestCharge RequestChargeWithNonce
			err = json.Unmarshal(recvResponse.Message.Body, &requestCharge)
			if err != nil {
				logger.Fatalf("Error unmarshalling RequestChargeWithNonce: %v", err)
			}
			logger.Printf("Received RequestChargeWithNonce: %v", requestCharge)
			approvePurchase := true
			reasonForRejection := ""
			if requestCharge.Nonce != "1234567890" {
				logger.Fatalf("Received invalid nonce: %v", requestCharge.Nonce)
				approvePurchase = false
				reasonForRejection = "Invalid nonce"
			}
			if requestCharge.Amount != cardDetails.TotalAmount {
				logger.Fatalf("Received invalid amount: %v", requestCharge.Amount)
				approvePurchase = false
				reasonForRejection = "Invalid amount"
			}

			// Send ChargeOK or ChargeFail to Srv.
			if approvePurchase {
				logger.Print("Sending ChargeOK to Srv...")
				sendResponse, err = client.AppSend(ctx, &pb.AppSendRequest{
					ChannelId: channelID,
					Recipient: "Srv",
					Message: &pb.AppMessage{
						Type: "ChargeOK",
						Body: []byte(""),
					},
				})
				if err != nil {
					logger.Fatal("Failed to send ChargeOK to Srv.")
				}
				if sendResponse.Result != pb.AppSendResponse_RESULT_OK {
					logger.Fatalf("Failed to send ChargeOK to Srv. Result: %v", sendResponse.Result)
				}
			} else {
				logger.Print("Sending ChargeFail to Srv...")
				sendResponse, err = client.AppSend(ctx, &pb.AppSendRequest{
					ChannelId: channelID,
					Recipient: "Srv",
					Message: &pb.AppMessage{
						Type: "ChargeFail",
						Body: []byte(reasonForRejection),
					},
				})
				if err != nil {
					logger.Fatal("Failed to send ChargeFail to Srv.")
				}
				if sendResponse.Result != pb.AppSendResponse_RESULT_OK {
					logger.Fatalf("Failed to send ChargeFail to Srv. Result: %v", sendResponse.Result)
				}
			}
			logger.Printf("Finished PPS program for channel %s", channelID)

		}(channelID, stub)
	}
}
