package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocql/gocql"
	businessService "github.com/microbusinesses/AddressService/business/service"
	"github.com/microbusinesses/AddressService/config"
	dataService "github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/AddressService/endpoint"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

var consulAddress string
var consulScheme string
var listeningPort int
var cassandraHosts string
var cassandraKeyspace string
var cassandraProtoclVersion int

func main() {
	flag.StringVar(&consulAddress, "consul-address", "", "The consul address in form of host:port. The default value is empty string.")
	flag.StringVar(&consulScheme, "consul-scheme", "", "The consul scheme. The default value is empty string.")
	flag.IntVar(&listeningPort, "listening-port", 0, "The port the application is serving HTTP request on. The default is zero.")
	flag.StringVar(&cassandraHosts, "cassandra-hosts", "", "The list of cassandra hosts to connect to. The default value is empty string.")
	flag.StringVar(&cassandraKeyspace, "cassandra-keyspace", "", "The cassandra keyspace. The default value is empty string.")
	flag.IntVar(&cassandraProtoclVersion, "cassandra-protocl-version", 0, "The cassandra protocl version. The default value is zero.")
	flag.Parse()

	consulConfigurationReader := config.ConsulConfigurationReader{ConsulAddress: consulAddress, ConsulScheme: consulScheme}

	setConsulConfigurationValuesRequireToBeOverriden(&consulConfigurationReader)

	endpoint := endpoint.Endpoint{ConfigurationReader: consulConfigurationReader}

	cassandraHosts, err := consulConfigurationReader.GetCassandraHosts()

	if err != nil {
		log.Fatal(err.Error())

		return
	}

	cassandraKeyspace, err := consulConfigurationReader.GetCassandraKeyspace()

	if err != nil {
		log.Fatal(err.Error())

		return
	}

	cassandraProtocolVersion, err := consulConfigurationReader.GetCassandraProtocolVersion()

	if err != nil {
		log.Fatal(err.Error())

		return
	}

	uuidGeneratorService := system.UUIDGeneratorServiceImpl{}

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

func setConsulConfigurationValuesRequireToBeOverriden(consulConfigurationReader *config.ConsulConfigurationReader) {
	diagnostics.IsNotNil(consulConfigurationReader, "consulConfigurationReader", "consulConfigurationReader is nil.")

	if listeningPort != 0 {
		consulConfigurationReader.ListeningPortToOverride = listeningPort
	} else if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil && port != 0 {
		consulConfigurationReader.ListeningPortToOverride = port
	}

	if len(cassandraHosts) != 0 {
		consulConfigurationReader.CassandraHostsToOverride = strings.Split(cassandraHosts, ",")
	}

	if len(cassandraKeyspace) != 0 {
		consulConfigurationReader.CassandraKeyspaceToOverride = cassandraKeyspace
	}

	if cassandraProtoclVersion != 0 {
		consulConfigurationReader.CassandraProtocolVersionToOverride = cassandraProtoclVersion
	}
}
