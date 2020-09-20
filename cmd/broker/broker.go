package main

import (
	"flag"

	"github.com/clpombo/search/internal/broker"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	host       = flag.String("host", "", "The host on which the broker listens (defaults to all interfaces")
	port       = flag.Int("port", 10000, "Port number on which the broker listens")
)

func main() {
	flag.Parse()
	b := broker.NewBrokerServer()
	b.StartServer(*host, *port, *tls, *certFile, *keyFile)
}

