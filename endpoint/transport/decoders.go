package transport

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/microbusinesses/AddressService/endpoint/message"
)

func DecodeApiRequest(context context.Context, httpRequest *http.Request) (interface{}, error) {
	var request message.ApiRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeCreateAddressRequest(context context.Context, httpRequest *http.Request) (interface{}, error) {
	var request message.CreateAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeUpdateAddressRequest(context context.Context, httpRequest *http.Request) (interface{}, error) {
	var request message.UpdateAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeReadAllAddressRequest(context context.Context, httpRequest *http.Request) (interface{}, error) {
	var request message.ReadAllAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeDeleteAddressRequest(context context.Context, httpRequest *http.Request) (interface{}, error) {
	var request message.DeleteAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}
