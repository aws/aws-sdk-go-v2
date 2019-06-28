package restjson

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/aws"
	"io/ioutil"
	"net/http"
	"testing"
)

type DataOutput struct {
	_ struct{} `type:"structure"`

	FooEnum string `type:"string" enum:"true"`

	ListEnums []string `type:"list"`
}

func BenchmarkRESTJSONUnmarshal_Simple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(getRESTJSONResponse_Simple())
	}
}

func getRESTJSONResponse_Simple() *aws.Request {
	buf := bytes.NewReader([]byte(`{"FooEnum": "foo", "ListEnums": ["0", "1"]}`))
	output := DataOutput{}
	req := aws.Request{Data: &output, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}
