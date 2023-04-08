package middleware

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"testing"

	pb "github.com/clpombo/search/gen/go/search/v1"
	"github.com/clpombo/search/internal/broker"
	"google.golang.org/grpc"

	"github.com/nickng/cfsm"
)

func TestTravelClient(t *testing.T) {
	// This example system of 3 CFSMs is a reproduction
	// of the system described in the paper "Communicating
	// machines as a dynamic binding mechanism of services", by
	// Vissani, López Pombo, Tuosto (p. 94)
	sys := cfsm.NewSystem()
	TravelClient := sys.NewMachine()
	HotelService := sys.NewMachine()
	PaymentProcessorService := sys.NewMachine()

	// TravelClient states
	tcStart := TravelClient.NewState()
	tc1 := TravelClient.NewState()
	tc2 := TravelClient.NewState()
	tc3 := TravelClient.NewState()
	tc4 := TravelClient.NewState()
	tc5 := TravelClient.NewState()
	TravelClient.Start = tcStart

	// HotelService states
	hsStart := HotelService.NewState()
	hs1 := HotelService.NewState()
	hs2 := HotelService.NewState()
	hs3 := HotelService.NewState()
	hs4 := HotelService.NewState()
	hs5 := HotelService.NewState()
	hs6 := HotelService.NewState()
	HotelService.Start = hsStart

	// PaymentProcessorService
	ppsStart := PaymentProcessorService.NewState()
	pps1 := PaymentProcessorService.NewState()
	pps2 := PaymentProcessorService.NewState()
	pps3 := PaymentProcessorService.NewState()
	PaymentProcessorService.Start = ppsStart

	// TravelClient transitions
	tBookHotelsSend := cfsm.NewSend(HotelService, "bookHotels")
	tBookHotelsSend.SetNext(tc1)
	tcStart.AddTransition(tBookHotelsSend)

	tHotelsRecv := cfsm.NewRecv(HotelService, "hotels")
	tHotelsRecv.SetNext(tc2)
	tc1.AddTransition(tHotelsRecv)

	tAcceptSend := cfsm.NewSend(HotelService, "accept")
	tAcceptSend.SetNext(tc3)
	tc2.AddTransition(tAcceptSend)

	tDeclineSend := cfsm.NewSend(HotelService, "decline")
	tDeclineSend.SetNext(tcStart)
	tc2.AddTransition(tDeclineSend)

	tPleasePayRecv := cfsm.NewRecv(PaymentProcessorService, "pleasePay")
	tPleasePayRecv.SetNext(tc4)
	tc3.AddTransition(tPleasePayRecv)

	tPaymentDataSend := cfsm.NewSend(PaymentProcessorService, "paymentData")
	tPaymentDataSend.SetNext(tc5)
	tc4.AddTransition(tPaymentDataSend)

	tPaymentRejectedRecv := cfsm.NewRecv(HotelService, "paymentRejected")
	tPaymentRejectedRecv.SetNext(tcStart)
	tc5.AddTransition(tPaymentRejectedRecv)

	tReservationsRecv := cfsm.NewRecv(HotelService, "reservations")
	tReservationsRecv.SetNext(tcStart)
	tc5.AddTransition(tReservationsRecv)

	// HotelService transitions
	tBookHotelsRecv := cfsm.NewRecv(TravelClient, "bookHotels")
	tBookHotelsRecv.SetNext(hs1)
	hsStart.AddTransition(tBookHotelsRecv)

	tHotelsSend := cfsm.NewSend(TravelClient, "hotels")
	tHotelsSend.SetNext(hs2)
	hs1.AddTransition(tHotelsSend)

	tAcceptRecv := cfsm.NewRecv(TravelClient, "accept")
	tAcceptRecv.SetNext(hs3)
	hs2.AddTransition(tAcceptRecv)

	tDeclineRecv := cfsm.NewRecv(TravelClient, "decline")
	tDeclineRecv.SetNext(hsStart)
	hs2.AddTransition(tDeclineRecv)

	tPaymentRejectedSend := cfsm.NewSend(TravelClient, "paymentRejected")
	tPaymentRejectedSend.SetNext(hsStart)
	hs5.AddTransition(tPaymentRejectedSend)

	tReservationsSend := cfsm.NewSend(TravelClient, "reservations")
	tReservationsSend.SetNext(hsStart)
	hs6.AddTransition(tReservationsSend)

	tAskForPaymentSend := cfsm.NewSend(PaymentProcessorService, "askForPayment")
	tAskForPaymentSend.SetNext(hs4)
	hs3.AddTransition(tAskForPaymentSend)

	tAcceptedRecv := cfsm.NewRecv(PaymentProcessorService, "accepted")
	tAcceptedRecv.SetNext(hs6)
	hs4.AddTransition(tAcceptedRecv)

	tRejectedRecv := cfsm.NewRecv(PaymentProcessorService, "rejected")
	tRejectedRecv.SetNext(hs5)
	hs4.AddTransition(tRejectedRecv)

	// PaymentProcessorService transitions
	tPleasePaySend := cfsm.NewSend(TravelClient, "pleasePay")
	tPleasePaySend.SetNext(pps2)
	pps1.AddTransition(tPleasePaySend)

	tPaymentDataRecv := cfsm.NewRecv(TravelClient, "paymentData")
	tPaymentDataRecv.SetNext(pps3)
	pps2.AddTransition(tPaymentDataRecv)

	tAskForPaymentRecv := cfsm.NewRecv(HotelService, "askForPayment")
	tAskForPaymentRecv.SetNext(pps1)
	ppsStart.AddTransition(tAskForPaymentRecv)

	tAcceptedSend := cfsm.NewSend(HotelService, "accepted")
	tAcceptedSend.SetNext(ppsStart)
	pps3.AddTransition(tAcceptedSend)

	tRejectedSend := cfsm.NewSend(HotelService, "rejected")
	tRejectedSend.SetNext(ppsStart)
	pps3.AddTransition(tRejectedSend)

	// Now we'll start a broker, a middleware for each CFSM
	// and try to get them to speak to each other.

	var wg sync.WaitGroup
	brokerPort := 7777
	hotelServicePublicPort := 4444
	hotelServicePrivatePort := 4445
	paymentProcesorServicePublicPort := 4446
	paymentProcesorServicePrivatePort := 4447
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	// Start broker
	bs := broker.NewBrokerServer()
	go bs.StartServer("localhost", brokerPort, false, "", "")
	defer bs.Stop()

	// Start HotelService middleware
	hotelServiceMiddleware := NewMiddlewareServer("localhost", brokerPort)
	hotelServiceMiddleware.StartMiddlewareServer(&wg, "localhost", hotelServicePublicPort, "localhost", hotelServicePrivatePort, false, "", "")

	// Start PaymentProcessorService middleware
	paymentProcessorMiddleware := NewMiddlewareServer("localhost", brokerPort)
	paymentProcessorMiddleware.StartMiddlewareServer(&wg, "localhost", paymentProcesorServicePublicPort, "localhost", paymentProcesorServicePrivatePort, false, "", "")

	go func() {
		// This runs the HotelService
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", hotelServicePrivatePort))
		if err != nil {
			t.Error("Could not connect to HotelService middleware.")
		}
		client := pb.NewPrivateMiddlewareServiceClient(conn)

		// Register the HotelService
		req := pb.RegisterAppRequest{
			ProviderContract: &pb.Contract{
				Contract: []byte("Serialized HotelService CFSM."), // TODO: replace with fsa
			},
		}
		stream, err := client.RegisterApp(context.Background(), &req)
		if err != nil {
			t.Error("Could not register HotelService with broker.")
		}
		ack, err := stream.Recv()
		if err != nil || ack.GetAppId() == "" {
			t.Error("Could not register HotelService with broker.")
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
		log.Printf("[HotelService - AppID %s] - Received Notification. ChannelID: %s", appID, channelID)

		// Start HotelService protocol
		recvReq := pb.AppRecvRequest{
			ChannelId:   channelID,
			Participant: "tc",
		}
		client.AppRecv(context.Background(), &recvReq)
	}()

}

func TestChorgramConversion(t *testing.T) {
	// First we create a Global Choreography for our requirement.
	const pingPongGC = `
Ping -> Pong : finished
   +
   *{
      Ping -> Pong : ping ; Pong -> Ping : pong
   } @ Ping ; Ping -> Pong : finished
`
	// Then we convert this to FSA format using Chorgram's gc2fsa
	const expectedPingPongFSA = `
.outputs Ping
.state graph
0 1 ! ping 5
2 1 ? bye 1
3 1 ! bye 2
3 1 ! finished 2
4 1 ! *<1 0
4 1 ! >*1 3
5 1 ? pong 4
.marking 0
.end



.outputs Pong
.state graph
0 0 ? ping 5
2 0 ! bye 1
3 0 ? bye 2
3 0 ? finished 2
4 0 ? *<1 0
4 0 ? >*1 3
5 0 ! pong 4
.marking 0
.end
`
	// TODO: call to Chorgram. Something like this:
	// convertedPingPongFSA := gc2fsa(pingPongGC)
	// if convertedPingPongFSA != expectedPingPongFSA { t.Error("Failed") }
}
