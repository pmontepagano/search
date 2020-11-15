package main

import (
	"flag"
	"sync"

	"github.com/clpombo/search/internal/middleware"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	caFile     = flag.String("ca_file", "", "The file containing the CA root cert file")
	publicHost = flag.String("public_host", "", "Host IP on which public service listens (defaults to all)")
	publicPort       = flag.Int("public_port", 10000, "The port for public facing middleware")
	privateHost = flag.String("private_host", "localhost", "Host IP on which private service listens")
	privatePort      = flag.Int("private_port", 11000, "The port for private services")
	brokerAddr = flag.String("broker_addr", "localhost", "The server address in the format of host:port")
	brokerPort = flag.Int("broker_port", 10000, "The port in which the broker is listening")
)

func main() {
	flag.Parse()
	mw := middleware.NewMiddlewareServer(*brokerAddr, *brokerPort)
	var wg sync.WaitGroup
	mw.StartMiddlewareServer(&wg, *publicHost, *publicPort, *privateHost, *privatePort, *tls, *certFile, *keyFile)
	wg.Wait()
}
