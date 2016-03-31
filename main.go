package main

import (
	"flag"

	"github.com/microbusinesses/AddressService/endpoint"
)

func main() {
	flag.IntVar(&endpoint.ListeningPort, "listening-port", 0, "The port the application is serving HTTP request on")
	flag.StringVar(&endpoint.CassandraAddress, "cassandra-address", "", "The address of the server the cassandra database is hosted on")
	flag.Parse()

	endpoint.StartServer()
}
