package main

import (
	"flag"
	"fmt"

	"github.com/pmontepagano/search/internal/broker"
)

var (
	tls          = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile     = flag.String("cert_file", "", "The TLS cert file")
	keyFile      = flag.String("key_file", "", "The TLS key file")
	host         = flag.String("host", "", "The host on which the broker listens (defaults to all interfaces")
	port         = flag.Int("port", 10000, "Port number on which the broker listens")
	databaseFile = flag.String("database_file", "SEARCHbroker.db", "The path for the broker database")
)

func main() {
	flag.Parse()
	b := broker.NewBrokerServer(*databaseFile)
	address := fmt.Sprintf("%s:%d", *host, *port)
	b.StartServer(address, *tls, *certFile, *keyFile, nil)
}
