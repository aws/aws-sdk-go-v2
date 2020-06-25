package restxml

import (
	"bytes"
	"github.com/jviney/aws-sdk-go-v2/aws"
	"io/ioutil"
	"net/http"
	"testing"
)

func BenchmarkRESTXMLUnmarshalError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalError(getRESTXMLError())
	}
}

func getRESTXMLError() *aws.Request {
	buf := bytes.NewReader([]byte(`
		<ErrorResponse>
		  <Error>
			<Code>baz</Code>
			<Message>test error message</Message>
		  </Error>
		  <RequestId>b25f48e8-84fd-11e6-80d9-574e0c4664cb</RequestId>
		</ErrorResponse>`))
	req := aws.Request{HTTPResponse: &http.Response{StatusCode: 404, Body: ioutil.NopCloser(buf), Header: http.Header{}}}
	req.HTTPResponse.Header.Set("X-Amzn-Errortype", "baz")
	return &req
}
