package rest

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"io/ioutil"
	"net/http"
	"testing"
)

type DataOutput struct {
	_ struct{} `type:"structure"`

	FooEnum string `type:"string" enum:"true"`

	HeaderEnum string `location:"header" locationName:"x-amz-enum" type:"string" enum:"true"`

	ListEnums []string `type:"list"`
}

func BenchmarkRESTUnmarshalMeta(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalMeta(getRESTResponse())
	}
}

func getRESTResponse() *aws.Request {
	output := DataOutput{}
	req := aws.Request{Data: &output, HTTPResponse: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(nil), Header: http.Header{}}}
	req.HTTPResponse.Header.Set("x-amz-enum", "baz")
	return &req
}

//type DataOutput2 struct {
//	_ struct{} `type:"structure"`
//
//	_ *string `type:"string" payload:"*string"`
//}

//func BenchmarkRESTUnmarshalBody(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		out := getRESTResponse2()
//		Unmarshal(out)
//		fmt.Println(out.Data)
//	}
//}

//func getRESTResponse2() *aws.Request {
//	buf := bytes.NewReader([]byte("String"))
//	output := DataOutput2{}
//	req := aws.Request{Data: &output, HTTPResponse: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf)}}
//	return &req
//}
