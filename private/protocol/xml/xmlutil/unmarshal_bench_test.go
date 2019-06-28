package xmlutil

import (
	"bytes"
	"encoding/xml"
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

func BenchmarkXMLUnmarshal_Simple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := getXMLResponse_Simple()
		UnmarshalXML(req.Data, xml.NewDecoder(req.HTTPResponse.Body), "")
	}
}

func getXMLResponse_Simple() *aws.Request {
	buf := bytes.NewReader([]byte("<OperationNameResponse><FooEnum>foo</FooEnum><ListEnums><member>0</member><member>1</member></ListEnums></OperationNameResponse>"))
	output := DataOutput{}
	req := aws.Request{Data: &output, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}
