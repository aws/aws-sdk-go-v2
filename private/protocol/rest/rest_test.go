package rest_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	metadata "github.com/aws/aws-sdk-go-v2/aws"
	request "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/private/protocol/rest"
)

func TestUnsetHeaders(t *testing.T) {
	cfg := unit.Config()
	cfg.Region = "us-west-2"

	svc := aws.NewClient(
		cfg,
		metadata.Metadata{
			ServiceName: "testService",
		},
	)

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(rest.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(rest.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(rest.UnmarshalMetaHandler)
	op := &request.Operation{
		Name:     "test-operation",
		HTTPPath: "/",
	}

	input := &struct {
		Foo aws.JSONValue `location:"header" locationName:"x-amz-foo" type:"jsonvalue"`
		Bar aws.JSONValue `location:"header" locationName:"x-amz-bar" type:"jsonvalue"`
	}{}

	output := &struct {
		Foo aws.JSONValue `location:"header" locationName:"x-amz-foo" type:"jsonvalue"`
		Bar aws.JSONValue `location:"header" locationName:"x-amz-bar" type:"jsonvalue"`
	}{}

	req := svc.NewRequest(op, input, output)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(nil)), Header: http.Header{}}
	req.HTTPResponse.Header.Set("X-Amz-Foo", "e30=")

	// unmarshal response
	rest.UnmarshalMeta(req)
	rest.Unmarshal(req)
	if req.Error != nil {
		t.Fatal(req.Error)
	}
}
