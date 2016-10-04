package endpoint

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	businessContract "github.com/microbusinesses/AddressService/business/contract"
	"github.com/microbusinesses/AddressService/config"
	"github.com/microbusinesses/AddressService/endpoint/transport"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"golang.org/x/net/context"
)

type Endpoint struct {
	ConfigurationReader config.ConfigurationReader
	AddressService      businessContract.AddressService
}

func (endpoint Endpoint) StartServer() {
	diagnostics.IsNotNil(endpoint.AddressService, "endpoint.AddressService", "AddressService must be provided.")
	diagnostics.IsNotNil(endpoint.ConfigurationReader, "endpoint.ConfigurationReader", "ConfigurationReader must be provided.")

	ctx := context.Background()

	handlers := getHandlers(endpoint, ctx)
	http.HandleFunc("/CheckHealth", checkHealthHandleFunc)

	for pattern, handler := range handlers {
		http.Handle(pattern, handler)
	}

	if listeningPort, err := endpoint.ConfigurationReader.GetListeningPort(); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(listeningPort), nil))
	}
}

func getHandlers(endpoint Endpoint, ctx context.Context) map[string]http.Handler {
	handlers := make(map[string]http.Handler)
	handlers["/Api"] = createApiHandler(endpoint, ctx)
	handlers["/CreateAddress"] = createCreateAddressHandler(endpoint, ctx)
	handlers["/UpdateAddress"] = createUpdateAddressHandler(endpoint, ctx)
	handlers["/ReadAllAddress"] = createReadAllAddressHandler(endpoint, ctx)
	handlers["/DeleteAddress"] = createDeleteAddressHandler(endpoint, ctx)

	return handlers
}

func checkHealthHandleFunc(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Alive")
}

func createApiHandler(endpoint Endpoint, ctx context.Context) http.Handler {
	return httptransport.NewServer(
		ctx,
		createApiEndpoint(endpoint.AddressService),
		transport.DecodeApiRequest,
		transport.EncodeApiResponse)
}

func createCreateAddressHandler(endpoint Endpoint, ctx context.Context) http.Handler {
	return httptransport.NewServer(
		ctx,
		createCreateAddressEndpoint(endpoint.AddressService),
		transport.DecodeCreateAddressRequest,
		transport.EncodeCreateAddressResponse)
}

func createUpdateAddressHandler(endpoint Endpoint, ctx context.Context) http.Handler {
	return httptransport.NewServer(
		ctx,
		createUpdateAddressEndpoint(endpoint.AddressService),
		transport.DecodeUpdateAddressRequest,
		transport.EncodeUpdateAddressResponse)
}

func createReadAllAddressHandler(endpoint Endpoint, ctx context.Context) http.Handler {
	return httptransport.NewServer(
		ctx,
		createReadAllAddressEndpoint(endpoint.AddressService),
		transport.DecodeReadAllAddressRequest,
		transport.EncodeReadAllAddressResponse)
}

func createDeleteAddressHandler(endpoint Endpoint, ctx context.Context) http.Handler {
	return httptransport.NewServer(
		ctx,
		createDeleteAddressEndpoint(endpoint.AddressService),
		transport.DecodeDeleteAddressRequest,
		transport.EncodeDeleteAddressResponse)
}
