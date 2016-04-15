package endpoint

import (
	"log"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	businessService "github.com/microbusinesses/AddressService/business/service"
	"github.com/microbusinesses/AddressService/endpoint/transport"
	"github.com/microbusinesses/Micro-Businesses-Core/common/diagnostics"
	"golang.org/x/net/context"
)

type Endpoint struct {
	ListeningPort  int
	AddressService businessService.AddressService
}

func (endpoint Endpoint) StartServer() {
	diagnostics.IsNotNil(endpoint.AddressService, "endpoint.AddressService", "AddressService must be provided.")

	ctx := context.Background()

	if handlers, err := getHandlers(endpoint, ctx); err != nil {
		log.Fatal(err.Error())
	} else {
		for pattern, handler := range handlers {
			http.Handle(pattern, handler)
		}

		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(endpoint.ListeningPort), nil))
	}
}

func getHandlers(endpoint Endpoint, ctx context.Context) (map[string]http.Handler, error) {
	handlers := make(map[string]http.Handler)

	if handler, err := createCreateAddressRequestHandler(endpoint, ctx); err != nil {
		return map[string]http.Handler{}, err
	} else {
		handlers["/CreateAddress"] = handler
	}

	if handler, err := createUpdateAddressRequestHandler(endpoint, ctx); err != nil {
		return map[string]http.Handler{}, err
	} else {
		handlers["/UpdateAddress"] = handler
	}

	if handler, err := createReadAddressRequestHandler(endpoint, ctx); err != nil {
		return map[string]http.Handler{}, err
	} else {
		handlers["/ReadAddress"] = handler
	}

	if handler, err := createDeleteAddressRequestHandler(endpoint, ctx); err != nil {
		return map[string]http.Handler{}, err
	} else {
		handlers["/DeleteAddress"] = handler
	}

	return handlers, nil
}

func createCreateAddressRequestHandler(endpoint Endpoint, ctx context.Context) (http.Handler, error) {
	return httptransport.NewServer(
		ctx,
		createCreateAddressEndpoint(endpoint.AddressService),
		transport.DecodeCreateAddressRequest,
		transport.EncodeCreateAddressResponse), nil
}

func createUpdateAddressRequestHandler(endpoint Endpoint, ctx context.Context) (http.Handler, error) {
	return httptransport.NewServer(
		ctx,
		createUpdateAddressEndpoint(endpoint.AddressService),
		transport.DecodeUpdateAddressRequest,
		transport.EncodeUpdateAddressResponse), nil
}

func createReadAddressRequestHandler(endpoint Endpoint, ctx context.Context) (http.Handler, error) {
	return httptransport.NewServer(
		ctx,
		createReadAddressEndpoint(endpoint.AddressService),
		transport.DecodeReadAddressRequest,
		transport.EncodeReadAddressResponse), nil
}

func createDeleteAddressRequestHandler(endpoint Endpoint, ctx context.Context) (http.Handler, error) {
	return httptransport.NewServer(
		ctx,
		createDeleteAddressEndpoint(endpoint.AddressService),
		transport.DecodeDeleteAddressRequest,
		transport.EncodeDeleteAddressResponse), nil
}
