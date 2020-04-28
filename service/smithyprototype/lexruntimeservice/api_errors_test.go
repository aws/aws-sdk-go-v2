package lexruntimeservice

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/awslabs/smithy-go"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

/*
/// Mock smithy tests
apply InvalidParametersException @httpResponseTests([
	{
		id: "RestJsonErrorDecode",
		documentation: "deserialize error response body",
        protocol: "aws.rest-json-1.1",
        code: 400,
        body: """
			{"code": "InvalidParametersException", "message": "invalid input parameters"}
             """,
        bodyMediaType: "application/json",
        headers: {
            "Content-Type": "application/json",
			"X-Amz-RequestId": "abc123"
        },
		vendorParams: {
			"RequestID": "abc123"
		},
        params: {
			Code: "InvalidParametersException",
            Message: "invalid input parameters",
			RetryAfterSeconds: "10",
        }
	}
	... Other tests
])
*/
func TestInvalidParametersException_deserialize_RestJsonErrorDecode(t *testing.T) {
	t.Skip("middleware not implemented")

	cfg := unit.Config()
	client := NewFromConfig(cfg, func(o *Options) {
		o.HTTPClient = smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Header: http.Header{
					"Content-Type":    []string{"application/json"},
					"X-Amz-RequestId": []string{"abc123"},
				},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"code": "UnknownError", "message": "invalid message"}`))),
			}, nil
		})
	})

	stack := middleware.NewStack("modeled api error deserialize", smithyhttp.NewStackRequest)
	// TODO build deserialization stack for errors.
	for _, fn := range client.options.APIOptions {
		if err := fn(stack); err != nil {
			t.Fatalf("failed to modify stack, %f", err)
		}
	}

	h := middleware.DecorateHandler(awshttp.ClientHandler{
		Client: client.options.HTTPClient,
	}, stack)

	// Mimic an API operation's deserialize behavior in the stack.
	result, _, err := h.Handle(context.Background(), struct{}{})
	if err == nil {
		t.Fatalf("expect error, got none")
	}
	if result != nil {
		t.Errorf("expect nil result, got %v", result)
	}

	// Assert params
	var apiErr *InvalidParameterException
	if !errors.As(err, &apiErr) {
		t.Fatalf("expect %T error, got %v", apiErr, err)
	}
	if e, a := "InvalidParametersException", apiErr.ErrorCode(); e != a {
		t.Errorf("expect %v code, got %v", e, a)
	}
	if e, a := "invalid message", apiErr.ErrorMessage(); e != a {
		t.Errorf("expect %v code, got %v", e, a)
	}
	if e, a := "10", apiErr.GetsRetryAfterSeconds(); e != a {
		t.Errorf("expect %v code, got %v", e, a)
	}

	// Vendor Metadata
	var reqErr *awshttp.ResponseError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expect %T error, got %v", reqErr, err)
	}
	if e, a := "abc123", reqErr.RequestID; e != a {
		t.Errorf("expect %v RequestID, got %v", e, a)
	}
}

func TestUnmodeledAPIError_deserialize(t *testing.T) {
	t.Skip("middleware not implemented")

	cfg := unit.Config()
	client := NewFromConfig(cfg, func(o *Options) {
		o.HTTPClient = smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Header: http.Header{
					"Content-Type":    []string{"application/json"},
					"X-Amz-RequestId": []string{"abc123"},
				},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"code": "UnknownError", "message": "invalid message"}`))),
			}, nil
		})
	})

	stack := middleware.NewStack("unmodeled api error deserialize", smithyhttp.NewStackRequest)
	// TODO build deserialization stack for errors.
	for _, fn := range client.options.APIOptions {
		if err := fn(stack); err != nil {
			t.Fatalf("failed to modify stack, %f", err)
		}
	}

	h := middleware.DecorateHandler(awshttp.ClientHandler{
		Client: client.options.HTTPClient,
	}, stack)

	// Mimic an API operation's deserialize behavior in the stack.
	result, _, err := h.Handle(context.Background(), struct{}{})
	if err == nil {
		t.Fatalf("expect error, got none")
	}
	if result != nil {
		t.Errorf("expect nil result, got %v", result)
	}

	// Assert params
	var apiErr *smithy.GenericAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expect %T error, got %v", apiErr, err)
	}
	if e, a := "UnknownError", apiErr.ErrorCode(); e != a {
		t.Errorf("expect %v code, got %v", e, a)
	}
	if e, a := "invalid message", apiErr.ErrorMessage(); e != a {
		t.Errorf("expect %v code, got %v", e, a)
	}

	// Vendor Metadata
	var reqErr *awshttp.ResponseError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expect %T error, got %v", reqErr, err)
	}
	if e, a := "abc123", reqErr.RequestID; e != a {
		t.Errorf("expect %v RequestID, got %v", e, a)
	}
}
