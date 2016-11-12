package endpoint

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/micro-business/AddressService/business/contract"
	"github.com/micro-business/AddressService/config"
	"github.com/micro-business/AddressService/endpoint/transport"
	"github.com/micro-business/Micro-Business-Core/common/diagnostics"
	"golang.org/x/net/context"
)

// Endpoint implements method to start the service. The structure contains all the dependencies required by the Endpoint service.
type Endpoint struct {
	ConfigurationReader config.ConfigurationReader
	AddressService      contract.AddressService
}

// StartServer creates all the endpoints and starts the server.
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
	handlers["/Api"] = createAPIHandler(endpoint, ctx)

	return handlers
}

func checkHealthHandleFunc(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Alive")
}

func createAPIHandler(endpoint Endpoint, ctx context.Context) http.Handler {
	return httptransport.NewServer(
		ctx,
		createAPIEndpoint(endpoint.AddressService),
		transport.DecodeAPIRequest,
		transport.EncodeAPIResponse)
}
