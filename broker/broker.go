package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	pb "dc.uba.ar/this/search/protobuf"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type brokerServer struct {
	pb.UnimplementedBrokerServer
	savedData []*pb.RemoteParticipant // read-only after initialized
}

func (s *brokerServer) BrokerChannel(ctx context.Context, request *pb.BrokerChannelRequest) (*pb.BrokerChannelResponse, error) {
	rand.Seed(time.Now().Unix())
	contract := request.GetContract()
	res := make(map[string]*pb.RemoteParticipant)
	for _, v := range contract.GetRemoteParticipants() {
		log.Println("Received requirements contract with participant", v)
		res[v] = s.savedData[rand.Intn(len(s.savedData))]
	}
	sessionID, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Failed to generate UUID")
	}
	return &pb.BrokerChannelResponse{Participants: res, SessionId: sessionID.String()}, nil
}

// loads data from a JSON file
func (s *brokerServer) loadData(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to load default data: %v", err)
		}
	} else {
		data = exampleData
	}
	if err := json.Unmarshal(data, &s.savedData); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}
}

func newServer() *brokerServer {
	s := &brokerServer{}
	s.loadData(*jsonDBFile)
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
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
	pb.RegisterBrokerServer(grpcServer, newServer())
	fmt.Println("Broker server starting...")
	grpcServer.Serve(lis)
}

// exampleData is a copy of testdata/route_guide_db.json. It's to avoid
// specifying file path with `go run`.
var exampleData = []byte(`[{
	"Url": "providermiddleware",
	"AppId": "example1FAKEuuid"
}, {
	"Url": "providermiddleware",
	"AppId": "example2FAKEuuid"
}]`)
