package middleware

import (
	"testing"

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
}
