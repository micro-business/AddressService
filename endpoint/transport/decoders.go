package transport

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
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
