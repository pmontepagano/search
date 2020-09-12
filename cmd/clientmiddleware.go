package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	pb "dc.uba.ar/this/search/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	brokerAddr         = flag.String("broker_addr", "localhost", "The server address in the format of host:port")
	brokerPort         = flag.Int("broker_port", 10000, "The port in which the broker is listening")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
	// provAddr           = flag.String("provider_addr", "localhost:20000", "The provider middleware address (ONLY FOR TESTING PURPOSES)")
	// provPort           = flag.Int("provider_port", 10000, "The provider middleware port (ONLY FOR TESTING PURPOSES)")
)

func initiateChannel(client pb.BrokerClient, contract *pb.Contract) pb.BrokerChannelResponse {
	log.Printf("Llamando al broker con contrato %v", contract.Contract)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	brokerresult, err := client.BrokerChannel(ctx, &pb.BrokerChannelRequest{Contract: contract})
	if err != nil {
		log.Fatalf("%v.GetCompatibleParticipants(_) = _, %v: ", client, err)
	}
	log.Println(brokerresult.Participants)
	return *brokerresult
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithBlock())

	fmt.Println("Waiting 2 seconds for broker...")
	time.Sleep(2 * time.Second)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *brokerAddr, *brokerPort), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBrokerClient(conn)

	// Pruebo llamar al broker
	brokerRes := initiateChannel(client, &pb.Contract{Contract: "hola mundo", RemoteParticipants: []string{"p1", "p2"}})

	// Ahora intento iniciar conexión con p1
	participants := brokerRes.GetParticipants()
	p1 := participants["p1"]

	provconn, err := grpc.Dial(fmt.Sprintf("%s:10000", p1.GetUrl()), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer provconn.Close()
	provClient := pb.NewPublicMiddlewareClient(provconn)

	stream, err := provClient.MessageExchange(context.Background())
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// se terminó
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("failed to receive data: %v", err)
			}
			log.Printf("Received message %s from %s", in.Content.Body, in.SenderId)
		}
	}()
	for i := 0; i < 5; i++ {
		msg := pb.ApplicationMessageWithHeaders{
			SenderId:    "clientmw-1",
			ChannelId:   "testSESSION",
			RecipientId: "provmw-44",
			Content:     &pb.MessageContent{Body: []byte(fmt.Sprintf("hola %d", i))}}

		if err := stream.Send(&msg); err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc

}
