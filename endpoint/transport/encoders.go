package transport

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

// EncodeAPIResponse encodes the response message before sending back to the client
func EncodeAPIResponse(context context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	return json.NewEncoder(writer).Encode(response)
}
