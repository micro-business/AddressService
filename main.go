package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/gocql/gocql"
	businessService "github.com/microbusinesses/AddressService/business/service"
	"github.com/microbusinesses/AddressService/config"
	dataService "github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/AddressService/endpoint"
)

var listeningPort int
var consulAddress string
var consulScheme string

func main() {
	flag.IntVar(&listeningPort, "listening-port", 80, "The port the application is serving HTTP request on - default is 80")
	flag.StringVar(&consulAddress, "consul-address", "127.0.0.1:8500", "The consul address in form of host:port. The default value is 127.0.0.1:8500")
	flag.StringVar(&consulScheme, "consul-scheme", "http", "The consul scheme. The default value is http")
	flag.Parse()

	endpoint := endpoint.Endpoint{ListeningPort: listeningPort}

	if endpoint.ListeningPort == 0 {
		if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
			endpoint.ListeningPort = port
		}
	}

	consulConfigurationReader := config.ConsulConfigurationReader{ConsulAddress: consulAddress, ConsulScheme: consulScheme}
	cassandraHosts, err := consulConfigurationReader.GetCassandraHosts()

	if err != nil {
		log.Fatal(err)

		return
	}

	cassandraKeyspace, err := consulConfigurationReader.GetCassandraKeyspace()

	if err != nil {
		log.Fatal(err)

		return
	}

	cassandraProtocolVersion, err := consulConfigurationReader.GetCassandraProtocolVersion()

	if err != nil {
		log.Fatal(err)

		return
	}

	uuidGeneratorService := dataService.UUIDGeneratorServiceImpl{}

	cluster := gocql.NewCluster()
	cluster.Hosts = cassandraHosts
	cluster.ProtoVersion = cassandraProtocolVersion
	cluster.Keyspace = cassandraKeyspace
	cluster.Consistency = gocql.Quorum

	addressDataService := dataService.AddressDataService{UUIDGeneratorService: &uuidGeneratorService, ClusterConfig: cluster}
	addressService := businessService.AddressService{AddressDataService: &addressDataService}

	endpoint.AddressService = addressService

	endpoint.StartServer()
}
