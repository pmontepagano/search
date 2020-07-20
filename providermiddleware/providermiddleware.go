package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "dc.uba.ar/this/search/protobuf"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 20000, "The server port")
)

type publicMiddlewareServer struct {
	pb.UnimplementedPublicMiddlewareServer
}

// ApplicationMessaging is the main function that allows the two apps connected to speak to each other
func (s *publicMiddlewareServer) MessageExchange(stream pb.PublicMiddleware_MessageExchangeServer) error {
	// I need a channel that is exposed to the local app
	// for now we'll send mock messages from the MW itself
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Println("Received message from", in.SenderId, ":", string(in.Content.Body))

		// TODO: send to local app replacing sender name with local name

		ack := pb.ApplicationMessageWithHeaders{
			ChannelId:   in.ChannelId,
			RecipientId: in.SenderId,
			SenderId:    "provmwID-44",
			Content:     &pb.MessageContent{Body: []byte("ack")}}

		if err := stream.Send(&ack); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)

	var pms publicMiddlewareServer
	pb.RegisterPublicMiddlewareServer(grpcServer, &pms)
	grpcServer.Serve(lis)
}
