package restxml

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/aws"
	"io/ioutil"
	"net/http"
	"testing"
)

//var (
//	s3Svc *s3.Client
//)
//
//func TestMain(m *testing.M) {
//
//}

func BenchmarkRESTXML_Unmarshal(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(getRESTXMLResponse())
	}
}

func getRESTXMLResponse() *aws.Request {
	buf := bytes.NewReader([]byte("<OperationNameResponse><FooEnum>foo</FooEnum><ListEnums><member>0</member><member>1</member></ListEnums></OperationNameResponse>"))
	req := &aws.Request{Data: struct{}{}, HTTPResponse: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}}
	req.HTTPResponse.Header.Set("x-amz-enum", "baz")
	return req
}

