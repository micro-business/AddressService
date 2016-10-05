package transport

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"github.com/microbusinesses/AddressService/endpoint/message"
)

func DecodeApiRequest(context context.Context, httpRequest *http.Request) (interface{}, error) {
	if httpRequest.Method == "GET" {
		return httpRequest.URL.Query()["query"][0], nil
	}

	if query, err := ioutil.ReadAll(httpRequest.Body); err != nil {
		return nil, err
	} else {
		return string(query), nil
	}
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
