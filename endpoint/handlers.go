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

	if handlers, err := getHandlers(endpoint, ctx); err != nil {
		log.Fatal(err.Error())
	} else {
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
}

func getHandlers(endpoint Endpoint, ctx context.Context) (map[string]http.Handler, error) {
	handlers := make(map[string]http.Handler)

	if handler, err := createCreateAddressHandler(endpoint, ctx); err != nil {
		return map[string]http.Handler{}, err
	} else {
		handlers["/CreateAddress"] = handler
	}

	if handler, err := createUpdateAddressHandler(endpoint, ctx); err != nil {
		return map[string]http.Handler{}, err
	} else {
		handlers["/UpdateAddress"] = handler
	}

	if handler, err := createReadAddressHandler(endpoint, ctx); err != nil {
		return map[string]http.Handler{}, err
	} else {
		handlers["/ReadAddress"] = handler
	}

	if handler, err := createDeleteAddressHandler(endpoint, ctx); err != nil {
		return map[string]http.Handler{}, err
	} else {
		handlers["/DeleteAddress"] = handler
	}

	return handlers, nil
}

func checkHealthHandleFunc(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Alive")
}

func createCreateAddressHandler(endpoint Endpoint, ctx context.Context) (http.Handler, error) {
	return httptransport.NewServer(
		ctx,
		createCreateAddressEndpoint(endpoint.AddressService),
		transport.DecodeCreateAddressRequest,
		transport.EncodeCreateAddressResponse), nil
}

func createUpdateAddressHandler(endpoint Endpoint, ctx context.Context) (http.Handler, error) {
	return httptransport.NewServer(
		ctx,
		createUpdateAddressEndpoint(endpoint.AddressService),
		transport.DecodeUpdateAddressRequest,
		transport.EncodeUpdateAddressResponse), nil
}

func createReadAddressHandler(endpoint Endpoint, ctx context.Context) (http.Handler, error) {
	return httptransport.NewServer(
		ctx,
		createReadAddressEndpoint(endpoint.AddressService),
		transport.DecodeReadAddressRequest,
		transport.EncodeReadAddressResponse), nil
}

func createDeleteAddressHandler(endpoint Endpoint, ctx context.Context) (http.Handler, error) {
	return httptransport.NewServer(
		ctx,
		createDeleteAddressEndpoint(endpoint.AddressService),
		transport.DecodeDeleteAddressRequest,
		transport.EncodeDeleteAddressResponse), nil
}
