package protocoltest

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type httpClient struct{}

func (*httpClient) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Header:     request.Header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
		Request:    request,
	}, nil
}

func NewClient() *httpClient {
	return &httpClient{}
}
