package jsonutil_test

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/json/jsonutil"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"io/ioutil"
	"net/http"
	"testing"
)

type DataOutput struct {
	_ struct{} `type:"structure"`

	FooEnum string `type:"string" enum:"true"`

	ListEnums []string `type:"list"`
}

func BenchmarkJSONUnmarshal_Simple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := getJSONResponse_Simple()
		jsonutil.UnmarshalJSON(req.Data, req.HTTPResponse.Body)
	}
}

func BenchmarkJSONUnmarshal_Complex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := getJSONResponse_Complex()
		jsonutil.UnmarshalJSON(req.Data, req.HTTPResponse.Body)
	}
}

func getJSONResponse_Simple() *aws.Request {
	buf := bytes.NewReader([]byte(`{"FooEnum": "foo", "ListEnums": ["0", "1"]}`))
	req := aws.Request{Data: &DataOutput{}, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}

func getJSONResponse_Complex() *aws.Request {
	buf := bytes.NewReader([]byte(`{"Table":{"AttributeDefinitions":[{"AttributeName":"1","AttributeType":"N"}],"CreationDateTime":1.562054355238E9,"ItemCount":0,"KeySchema":[{"AttributeName":"1","KeyType":"HASH"}],"ProvisionedThroughput":{"NumberOfDecreasesToday":0,"ReadCapacityUnits":5,"WriteCapacityUnits":5},"TableArn":"arn:aws:dynamodb:us-west-2:183557167593:table/TestTable","TableId":"575d0be6-34e3-4843-838c-8e8e8d4ea2f7","TableName":"TestTable","TableSizeBytes":0,"TableStatus":"ACTIVE"}}`))
	req := aws.Request{Data: &dynamodb.DescribeTableOutput{}, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}
