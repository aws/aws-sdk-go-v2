package imds

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestUnauthorizedErrorUnwrapsToResponseError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.Write([]byte("token"))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := New(Options{
		Endpoint: server.URL,
		Retryer:  aws.NopRetryer{},
	})

	_, err := client.GetMetadata(context.Background(), &GetMetadataInput{Path: "iam/security-credentials/"})
	if err == nil {
		t.Fatalf("expect error, got none")
	}

	var respErr *smithyhttp.ResponseError
	if !errors.As(err, &respErr) {
		t.Fatalf("expect error to unwrap to %T, got %v", respErr, err)
	}
	if e, a := http.StatusUnauthorized, respErr.HTTPStatusCode(); e != a {
		t.Errorf("expect %v status code, got %v", e, a)
	}
}
