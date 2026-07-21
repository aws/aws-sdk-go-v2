package kitchensinktestrestjson

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// closeTrackingBody wraps a response body and records whether Close was called.
// It lets these tests assert on the observable close behavior of the deserialize
// path rather than on which middleware happens to implement it.
type closeTrackingBody struct {
	r      io.Reader
	closed int32
}

func newCloseTrackingBody(s string) *closeTrackingBody {
	return &closeTrackingBody{r: strings.NewReader(s)}
}

func (b *closeTrackingBody) Read(p []byte) (int, error) { return b.r.Read(p) }

func (b *closeTrackingBody) Close() error {
	atomic.AddInt32(&b.closed, 1)
	return nil
}

func (b *closeTrackingBody) closeCount() int { return int(atomic.LoadInt32(&b.closed)) }

type closeTestEndpointResolver struct{}

func (*closeTestEndpointResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (smithyendpoints.Endpoint, error) {
	return smithyendpoints.Endpoint{URI: url.URL{Scheme: "https", Host: "test.example.com"}}, nil
}

func closeTestClient(status int, header http.Header, body io.ReadCloser) *Client {
	return New(Options{
		Region: "us-east-1",
		HTTPClient: smithyhttp.ClientDoFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode:    status,
				Header:        header,
				Body:          body,
				ContentLength: -1,
				Request:       req,
			}, nil
		}),
		EndpointResolverV2: &closeTestEndpointResolver{},
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				s.Finalize.Clear()
				return nil
			},
		},
	})
}

// TestResponseBodyClose_StreamingSuccess asserts that a streaming (caller-owned)
// output does NOT close the response body on the success path: the body is
// handed to the caller as Output.Body and must stay open for reading.
func TestResponseBodyClose_StreamingSuccess(t *testing.T) {
	const payload = "hello streaming body"
	body := newCloseTrackingBody(payload)

	svc := closeTestClient(200, http.Header{}, body)

	out, err := svc.GetStreamingResource(context.Background(), &GetStreamingResourceInput{
		Id: strPtr("some-id"),
	})
	if err != nil {
		t.Fatalf("get streaming resource: %v", err)
	}

	if body.closeCount() != 0 {
		t.Errorf("expected streaming response body to stay open on success, but it was closed %d time(s)", body.closeCount())
	}

	// The caller-owned stream must still be readable and yield the payload.
	got, err := io.ReadAll(out.Body)
	if err != nil {
		t.Fatalf("read output body: %v", err)
	}
	if string(got) != payload {
		t.Errorf("output body = %q, want %q", string(got), payload)
	}

	if err := out.Body.Close(); err != nil {
		t.Fatalf("close output body: %v", err)
	}
	if body.closeCount() != 1 {
		t.Errorf("expected exactly one close after caller closes the stream, got %d", body.closeCount())
	}
}

// TestResponseBodyClose_NonStreamingSuccess asserts that a non-streaming output
// closes the response body on the success path.
func TestResponseBodyClose_NonStreamingSuccess(t *testing.T) {
	body := newCloseTrackingBody(`{"name":"thing"}`)

	svc := closeTestClient(200, http.Header{}, body)

	_, err := svc.GetResource(context.Background(), &GetResourceInput{
		Id: strPtr("some-id"),
	})
	if err != nil {
		t.Fatalf("get resource: %v", err)
	}

	if body.closeCount() == 0 {
		t.Errorf("expected non-streaming response body to be closed on success path, but it was not")
	}
}

// TestResponseBodyClose_StreamingError asserts that the response body IS closed
// on the error path even for a streaming operation: there is no caller-owned
// stream to hand back when the operation fails.
func TestResponseBodyClose_StreamingError(t *testing.T) {
	header := http.Header{}
	header.Set("X-Amzn-Errortype", "ResourceNotFound")
	body := newCloseTrackingBody(`{"__type":"ResourceNotFound"}`)

	svc := closeTestClient(404, header, body)

	_, err := svc.GetStreamingResource(context.Background(), &GetStreamingResourceInput{
		Id: strPtr("some-id"),
	})
	if err == nil {
		t.Fatal("expected error, got none")
	}

	if body.closeCount() == 0 {
		t.Errorf("expected response body to be closed on error path, but it was not")
	}
}

func strPtr(s string) *string { return &s }
