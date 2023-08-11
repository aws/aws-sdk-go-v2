package protocoltest

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// HttpClient is a mock http client used by protocol test cases to
// respond success response back
type HttpClient struct{}

func (*HttpClient) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Header:     request.Header,
		Body:       ioutil.NopCloser(strings.NewReader("")),
		Request:    request,
	}, nil
}

// NewClient returns pointer of a new HttpClient for protocol test client
func NewClient() *HttpClient {
	return &HttpClient{}
}
