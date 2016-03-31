package endpoint

import (
	"log"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gocql/gocql"
	businessService "github.com/microbusinesses/AddressService/business/service"
	dataService "github.com/microbusinesses/AddressService/data/service"
	"github.com/microbusinesses/AddressService/endpoint/transport"
	"golang.org/x/net/context"
)

var ListeningPort int
var CassandraAddress string

func StartServer() {
	ctx := context.Background()

	if handlers, err := getHandlers(ctx); err != nil {
		log.Fatal(err.Error())
	} else {
		for pattern, handler := range handlers {
			http.Handle(pattern, handler)
		}

		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(ListeningPort), nil))
	}
}

func getHandlers(ctx context.Context) (map[string]http.Handler, error) {
	handlers := make(map[string]http.Handler)

	if handler, err := createCreateAddressRequestHandler(ctx); err != nil {
		return map[string]http.Handler{}, err
	} else {
		handlers["/CreateAddress"] = handler
	}

	return handlers, nil
}

func createCreateAddressRequestHandler(ctx context.Context) (http.Handler, error) {
	uuidGeneratorService := dataService.UUIDGeneratorServiceImpl{}

	cluster := gocql.NewCluster(CassandraAddress)
	cluster.Keyspace = "address"
	cluster.Consistency = gocql.Quorum

	addressDataService := dataService.AddressDataService{UUIDGeneratorService: &uuidGeneratorService, ClusterConfig: cluster}
	addressService := businessService.AddressService{AddressDataService: &addressDataService}

	return httptransport.NewServer(
		ctx,
		createCreateAddressEndpoint(addressService),
		transport.DecodeCreateAddressRequest,
		transport.EncodeCreateAddressResponse), nil
}
