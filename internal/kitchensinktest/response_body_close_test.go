package kitchensinktest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"

	"github.com/aws/aws-sdk-go-v2/internal/kitchensinktest/types"
)

// trackedBody wraps a response body and records whether Close was called and
// how many bytes were read out of it before closing. It lets the tests assert
// on the observable close/consume behavior of the deserialize path rather than
// on which middleware happens to implement it.
type trackedBody struct {
	r        io.Reader
	closed   atomic.Int32
	readN    atomic.Int64
	closeErr error
}

func newTrackedBody(b []byte) *trackedBody {
	return &trackedBody{r: bytes.NewReader(b)}
}

func (t *trackedBody) Read(p []byte) (int, error) {
	n, err := t.r.Read(p)
	t.readN.Add(int64(n))
	return n, err
}

func (t *trackedBody) Close() error {
	t.closed.Add(1)
	return t.closeErr
}

func (t *trackedBody) wasClosed() bool  { return t.closed.Load() > 0 }
func (t *trackedBody) closeCount() int  { return int(t.closed.Load()) }
func (t *trackedBody) bytesRead() int64 { return t.readN.Load() }

type closeTestEndpointResolver struct{}

func (*closeTestEndpointResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (smithyendpoints.Endpoint, error) {
	return smithyendpoints.Endpoint{URI: url.URL{Scheme: "https", Host: "test.example.com"}}, nil
}

// closeTestClient builds a client whose transport returns a response backed by
// the given tracked body, with auth/finalize stripped so the test exercises the
// deserialize path in isolation.
func closeTestClient(status int, body *trackedBody) *Client {
	return New(Options{
		Region: "us-east-1",
		HTTPClient: smithyhttp.ClientDoFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode:    status,
				Header:        http.Header{},
				Body:          body,
				ContentLength: -1,
				Request:       req,
			}, nil
		}),
		EndpointResolverV2: &closeTestEndpointResolver{},
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				s.Finalize.Clear()
				s.Initialize.Remove("OperationInputValidation")
				return nil
			},
		},
	})
}

// TestResponseBodyClose_Success asserts that a non-streaming operation closes
// the response body on the success path.
func TestResponseBodyClose_Success(t *testing.T) {
	body := newTrackedBody([]byte(`{}`))
	svc := closeTestClient(200, body)

	_, err := svc.GetItem(context.Background(), &GetItemInput{})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !body.wasClosed() {
		t.Errorf("expected response body to be closed on success path")
	}
	if body.closeCount() != 1 {
		t.Errorf("expected body closed exactly once, got %d", body.closeCount())
	}
}

// TestResponseBodyClose_Error asserts that a non-streaming operation closes the
// response body on the error path.
func TestResponseBodyClose_Error(t *testing.T) {
	body := newTrackedBody([]byte(`{"__type":"ItemNotFound"}`))
	svc := closeTestClient(400, body)

	_, err := svc.GetItem(context.Background(), &GetItemInput{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var notFound *types.ItemNotFound
	if !errors.As(err, &notFound) {
		t.Fatalf("expected types.ItemNotFound, got %v", err)
	}
	if !body.wasClosed() {
		t.Errorf("expected response body to be closed on error path")
	}
	if body.closeCount() != 1 {
		t.Errorf("expected body closed exactly once, got %d", body.closeCount())
	}
}
