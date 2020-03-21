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
		w.Write([]byte(Status200Response))
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

func TestProtoCreateApiKeyRequestErrorUnmarshaler_Diff(t *testing.T) {
	cases := map[string]struct {
		status          int
		response        []byte
		ErrorTypeHeader string
	}{
		"FailureCase": {
			status:          500,
			response:        []byte(Status500Response),
			ErrorTypeHeader: "",
		},
		"FailureCaseWithHeader": {
			status:          500,
			response:        []byte(Status500Response),
			ErrorTypeHeader: "baz",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.status)
				w.Header().Set("X-Amzn-Errortype", c.ErrorTypeHeader)
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
			_, expectedErr := request.Send(context.TODO())
			if expectedErr == nil {
				t.Fatalf("Expected error, got none")
			}

			prototypeRequest := svc.ProtoCreateAPIKeyRequest(&input)
			_, prototypeErr := prototypeRequest.Send(context.TODO())
			if prototypeErr == nil {
				t.Fatalf("Expected error, got none")
			}

			if diff := cmp.Diff(expectedErr.Error(), prototypeErr.Error()); diff != "" {
				t.Errorf("Found diff: %v", diff)
			}
		})
	}
}

const Status200Response = `{
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
}`

const Status500Response = `{
	"code" : "500 code",
	"message" : "Internal service error"
}`
