package apigateway_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/google/go-cmp/cmp"
)

func TestProtoCreateApiKeyRequestUnmarshaler_Diff(t *testing.T) {
	// server is the mock server that simply writes a 200 status back to the client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
	"customerId": "mock customer id",
	"Name": "mock name",
	"createdDate": 1494359783.453,
	"description": "mock operation description",
	"enabled": true,
	"id": "mockid",
	"lastUpdatedDate": 1494359783.453,
	"stageKeys": ["mock stage key"],
	"tags": {
		"a": "1",
		"b": "2"
	},
	"value": "mock value"
}`))
	}))

	defer server.Close()

	config := mock.Config()
	config.EndpointResolver = aws.ResolveWithEndpoint(aws.Endpoint{
		URL:           server.URL,
		SigningRegion: config.Region,
	})

	svc := apigateway.New(config)

	input := apigateway.CreateApiKeyInput{}

	request := svc.CreateApiKeyRequest(&input)
	expectedResponse, err := request.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	prototypeRequest := svc.ProtoCreateAPIKeyRequest(&input)
	prototypeResponse, err := prototypeRequest.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(expectedResponse.CreateApiKeyOutput, prototypeResponse.CreateApiKeyOutput); diff != "" {
		t.Errorf("Found diff: %v", diff)
	}
}
