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
