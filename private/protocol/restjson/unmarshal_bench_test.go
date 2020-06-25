package restjson

import (
	"bytes"
	"github.com/jviney/aws-sdk-go-v2/aws"
	"io/ioutil"
	"net/http"
	"testing"
)

func BenchmarkRESTJSONUnmarshalError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalError(getRESTJSONError())
	}
}

func getRESTJSONError() *aws.Request {
	buf := bytes.NewReader([]byte(`{"message":"test error message"}`))
	req := aws.Request{RequestID: "b25f48e8-84fd-11e6-80d9-574e0c4664cb",
		HTTPResponse: &http.Response{StatusCode: 404, Body: ioutil.NopCloser(buf), Header: http.Header{}}}
	req.HTTPResponse.Header.Set("X-Amzn-Errortype", "baz")
	return &req
}
