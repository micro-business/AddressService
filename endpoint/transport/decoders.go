package transport

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/microbusinesses/AddressService/endpoint/message"
)

func DecodeCreateAddressRequest(_ context.Context, httpRequest *http.Request) (interface{}, error) {
	var request message.CreateAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeUpdateAddressRequest(_ context.Context, httpRequest *http.Request) (interface{}, error) {
	var request message.UpdateAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeReadAddressRequest(_ context.Context, httpRequest *http.Request) (interface{}, error) {
	var request message.ReadAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeDeleteAddressRequest(_ context.Context, httpRequest *http.Request) (interface{}, error) {
	var request message.DeleteAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}
