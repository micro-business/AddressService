package transport

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
)

// DecodeAPIRequest decodes the request message sent by the client. The request message can be sent using GET method as part of
// URL or can be the payload of a POST HTTP message.
func DecodeAPIRequest(context context.Context, httpRequest *http.Request) (interface{}, error) {
	if httpRequest.Method == "GET" {
		return httpRequest.URL.Query()["query"][0], nil
	}

	query, err := ioutil.ReadAll(httpRequest.Body)

	if err != nil {
		return nil, err
	}

	return string(query), nil
}
