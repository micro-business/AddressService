package main

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/microbusinesses/AddressService/endpoint"
)

func main() {
	flag.IntVar(&endpoint.ListeningPort, "listening-port", 80, "The port the application is serving HTTP request on - default is 80")
	flag.StringVar(&endpoint.CassandraAddress, "cassandra-address", "", "The address of the server the cassandra database is hosted on")
	flag.Parse()

	if endpoint.ListeningPort == 0 {
		if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
			endpoint.ListeningPort = port
		}
	}

	if len(strings.TrimSpace(endpoint.CassandraAddress)) == 0 {
		endpoint.CassandraAddress = os.Getenv("CASSANDRA_ADDRESS")
	}

	endpoint.StartServer()
}
