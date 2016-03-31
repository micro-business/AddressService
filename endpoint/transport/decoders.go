package transport

import (
	"encoding/json"
	"net/http"

	"github.com/microbusinesses/AddressService/endpoint/message"
)

func DecodeCreateAddressRequest(httpRequest *http.Request) (interface{}, error) {
	var request message.CreateAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeUpdateAddressRequest(httpRequest *http.Request) (interface{}, error) {
	var request message.UpdateAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeReadAddressRequest(httpRequest *http.Request) (interface{}, error) {
	var request message.ReadAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func DecodeDeleteAddressRequest(httpRequest *http.Request) (interface{}, error) {
	var request message.DeleteAddressRequest

	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}
