package jsonutil_test

import (
	"bytes"
	"encoding/json"
	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/private/protocol/json/jsonutil"
	"github.com/jviney/aws-sdk-go-v2/service/dynamodb"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	simpleJSON  = []byte(`{"FooEnum": "foo", "ListEnums": ["0", "1"]}`)
	complexJSON = []byte(`{"Table":{"AttributeDefinitions":[{"AttributeName":"1","AttributeType":"N"}],"CreationDateTime":1.562054355238E9,"ItemCount":0,"KeySchema":[{"AttributeName":"1","KeyType":"HASH"}],"ProvisionedThroughput":{"NumberOfDecreasesToday":0,"ReadCapacityUnits":5,"WriteCapacityUnits":5},"TableArn":"arn:aws:dynamodb:us-west-2:183557167593:table/TestTable","TableId":"575d0be6-34e3-4843-838c-8e8e8d4ea2f7","TableName":"TestTable","TableSizeBytes":0,"TableStatus":"ACTIVE"}}`)
)

type DataOutput struct {
	_ struct{} `type:"structure"`

	FooEnum string `type:"string" enum:"true"`

	ListEnums []string `type:"list"`
}

func BenchmarkJSONUnmarshal_Simple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := getJSONResponseSimple()
		jsonutil.UnmarshalJSON(req.Data, req.HTTPResponse.Body)
	}
}

func BenchmarkJSONUnmarshal_Complex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := getJSONResponseComplex()
		jsonutil.UnmarshalJSON(req.Data, req.HTTPResponse.Body)
	}
}

func BenchmarkStdlibJSON_Unmarshal_Simple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Unmarshal(simpleJSON, &DataOutput{})
	}
}

func BenchmarkStdlibJSON_Unmarshal_Complex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Unmarshal(complexJSON, &dynamodb.DescribeTableOutput{})
	}
}

func getJSONResponseSimple() *aws.Request {
	buf := bytes.NewReader(simpleJSON)
	req := aws.Request{Data: &DataOutput{}, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}

func getJSONResponseComplex() *aws.Request {
	buf := bytes.NewReader(complexJSON)
	req := aws.Request{Data: &dynamodb.DescribeTableOutput{}, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}
