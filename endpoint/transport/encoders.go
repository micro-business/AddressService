package transport

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func EncodeCreateAddressResponse(context context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	return json.NewEncoder(writer).Encode(response)
}

func EncodeUpdateAddressResponse(context context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	return json.NewEncoder(writer).Encode(response)
}

func EncodeReadAddressResponse(context context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	return json.NewEncoder(writer).Encode(response)
}

func EncodeDeleteAddressResponse(context context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	return json.NewEncoder(writer).Encode(response)
}
