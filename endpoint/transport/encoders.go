package transport

import (
	"encoding/json"
	"net/http"
)

func EncodeCreateAddressResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
