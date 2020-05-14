package apigateway_test

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/google/go-cmp/cmp"
)

func assertJSON(t testing.TB, a, b io.ReadCloser) {
	t.Helper()
	buf, err := ioutil.ReadAll(a)
	if err != nil {
		t.Fatal("Error reading a", err)
	}
	protoBuff, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal("Error reading b", err)
	}
	var av, bv interface{}

	if err := json.Unmarshal(buf, &av); err != nil {
		t.Fatalf("assertJSON: unable to unmarshal a, %v", err)
	}

	if err := json.Unmarshal(protoBuff, &bv); err != nil {
		t.Fatalf("assertJSON: unable to unmarshal b, %v", err)
	}

	if !reflect.DeepEqual(av, bv) {
		t.Fatalf("JSON are not equal\nexpect:\n%v\nactual:\n%v\n", a, b)
	}
}

func TestProtoCreateApiKeyRequestMarshaler_Diff(t *testing.T) {
	svc := apigateway.New(mock.Config())
	input := apigateway.CreateApiKeyInput{
		CustomerId:         aws.String("mock id"),
		Description:        aws.String("mock operation description"),
		Enabled:            aws.Bool(true),
		GenerateDistinctId: aws.Bool(true),
		Name:               aws.String("mock name"),
		StageKeys: []apigateway.StageKey{apigateway.StageKey{
			RestApiId: aws.String("mock rest api id"),
			StageName: aws.String("mock stage name"),
		}},
		Tags:  map[string]string{"a": "1", "b": "2"},
		Value: aws.String("mock value"),
	}

	request := svc.CreateApiKeyRequest(&input)
	_, err := request.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	prototypeRequest := svc.ProtoCreateAPIKeyRequest(&input)
	_, err = prototypeRequest.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(request.HTTPRequest.Header, prototypeRequest.HTTPRequest.Header); diff != "" {
		t.Errorf("Found diff: %v", diff)
	}
	if diff := cmp.Diff(request.HTTPRequest.URL, prototypeRequest.HTTPRequest.URL); diff != "" {
		t.Errorf("Found diff: %v", diff)
	}

	reqBody, err := request.HTTPRequest.GetBody()
	if err != nil {
		t.Fatal("failed to read body from request", err)
	}

	protoBody, err := prototypeRequest.HTTPRequest.GetBody()
	if err != nil {
		t.Fatal("failed to read body from prototyped request", err)
	}
	assertJSON(t, reqBody, protoBody)
}
  
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
