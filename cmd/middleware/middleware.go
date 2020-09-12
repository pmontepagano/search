package main

import (
	"flag"

	"github.com/clpombo/search/internal/middleware"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	caFile     = flag.String("ca_file", "", "The file containing the CA root cert file")
	port       = flag.Int("port", 10000, "The server port")
	brokerAddr = flag.String("broker_addr", "localhost", "The server address in the format of host:port")
	brokerPort = flag.Int("broker_port", 10000, "The port in which the broker is listening")
)

func main() {
	flag.Parse()
	middleware.StartServer(port, tls, certFile, keyFile, *brokerAddr, *brokerPort)
}
