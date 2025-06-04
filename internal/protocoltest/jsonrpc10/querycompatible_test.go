package jsonrpc10

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/smithy-go/middleware"
	smithyprivateprotocol "github.com/aws/smithy-go/private/protocol"
)

type mockHTTP struct{}

func (m mockHTTP) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 400, Body: http.NoBody}, nil
}

// SEP query-protocol-migration-compatibility.md
//
// TC5: Validate SDK does not send x-amzn-query-mode header when service
// doesn't have @awsQueryCompatible trait
func TestQueryCompatible_SEP5(t *testing.T) {
	var req http.Request
	svc := New(Options{
		HTTPClient: mockHTTP{},
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				return smithyprivateprotocol.AddCaptureRequestMiddleware(s, &req)
			},
		},
	})

	svc.SimpleScalarProperties(context.Background(), &SimpleScalarPropertiesInput{})
	if actual := req.Header.Values("X-Amzn-Query-Mode"); len(actual) != 0 {
		t.Errorf("X-Amzn-Query-Mode is set: %v", actual)
	}
}
