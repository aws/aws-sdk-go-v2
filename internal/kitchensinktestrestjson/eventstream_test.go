package kitchensinktestrestjson

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type eventStreamEndpointResolver struct{}

func (*eventStreamEndpointResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (smithyendpoints.Endpoint, error) {
	return smithyendpoints.Endpoint{URI: url.URL{Scheme: "https", Host: "test.example.com"}}, nil
}

// encodeEventStream serializes the given messages into the vnd.amazon.eventstream
// wire format so they can be served as a fake response body.
func encodeEventStream(t *testing.T, msgs []eventstream.Message) []byte {
	t.Helper()
	var buf bytes.Buffer
	enc := eventstream.NewEncoder()
	for _, m := range msgs {
		if err := enc.Encode(&buf, m); err != nil {
			t.Fatalf("encode event: %v", err)
		}
	}
	return buf.Bytes()
}

// TestSubscribeEvents_Close verifies the caller-owned event stream behavior on
// the success path: the response body is handed to the event stream reader
// rather than closed by the deserialize path, events are readable, and closing
// the stream via GetStream().Close() closes the underlying body exactly once.
// This locks in the pre-existing behavior so a middleware-stack change can be
// shown not to alter it.
func TestSubscribeEvents_Close(t *testing.T) {
	payload := encodeEventStream(t, []eventstream.Message{
		{
			Headers: eventstream.Headers{
				{Name: eventstreamapi.MessageTypeHeader, Value: eventstream.StringValue(eventstreamapi.EventMessageType)},
				{Name: eventstreamapi.EventTypeHeader, Value: eventstream.StringValue("initial-response")},
			},
			Payload: []byte(`{}`),
		},
		{
			Headers: eventstream.Headers{
				{Name: eventstreamapi.MessageTypeHeader, Value: eventstream.StringValue(eventstreamapi.EventMessageType)},
				{Name: eventstreamapi.EventTypeHeader, Value: eventstream.StringValue("message")},
			},
			Payload: []byte(`{"body":"hello"}`),
		},
	})

	body := &closeTrackingBody{r: bytes.NewReader(payload)}

	svc := New(Options{
		Region: "us-east-1",
		HTTPClient: smithyhttp.ClientDoFunc(func(req *http.Request) (*http.Response, error) {
			h := http.Header{}
			h.Set("Content-Type", "application/vnd.amazon.eventstream")
			return &http.Response{
				StatusCode:    200,
				Header:        h,
				Body:          body,
				ContentLength: -1,
				Request:       req,
			}, nil
		}),
		EndpointResolverV2: &eventStreamEndpointResolver{},
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				s.Finalize.Clear()
				return nil
			},
		},
	})

	resp, err := svc.SubscribeEvents(context.Background(), &SubscribeEventsInput{
		Id: strPtr("some-id"),
	})
	if err != nil {
		t.Fatalf("subscribe events: %v", err)
	}

	// The deserialize path must NOT close the response body on the success path:
	// for an event stream the body is caller-owned and drained by the stream's
	// background copy loop, so it stays open for the reader to consume.
	if got := body.closeCount(); got != 0 {
		t.Errorf("expected response body to stay open on event stream success path, but it was closed %d time(s)", got)
	}

	// The stream must be readable.
	event := <-resp.GetStream().Events()
	if event == nil {
		t.Fatal("expect an event, got nil")
	}

	// Closing the caller-owned stream must succeed without error.
	if err := resp.GetStream().Close(); err != nil {
		t.Errorf("expect no error closing stream, got %v", err)
	}
	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no stream error, got %v", err)
	}
}
