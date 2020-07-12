package main

import (
	"context"
	"flag"
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
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func printTest(client pb.BrokerClient, contract *pb.RequirementsContract) {
	log.Printf("Llamando al broker con contrato %v", contract.Contract)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	brokerresult, err := client.GetCompatibleParticipants(ctx, contract)
	if err != nil {
		log.Fatalf("%v.GetCompatibleParticipants(_) = _, %v: ", client, err)
	}
	log.Println(brokerresult.Participants)
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
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBrokerClient(conn)

	// Pruebo llamar al broker
	printTest(client, &pb.RequirementsContract{Contract: "hola mundo", Participants: []string{"p1", "p2"}})

}
