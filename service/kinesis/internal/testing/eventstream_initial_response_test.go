package testing

import (
	"context"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	"github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type closeTrackingReadCloser struct {
	io.ReadCloser
	closed atomic.Bool
}

func (r *closeTrackingReadCloser) Close() error {
	r.closed.Store(true)
	return r.ReadCloser.Close()
}

func TestSubscribeToShard_InitialResponseErrorClosesBody(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("initial-response"),
						},
					},
					Payload: []byte(`{`),
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("setup event stream: %v", err)
	}
	defer cleanupFn()

	baseClient := cfg.HTTPClient
	var body *closeTrackingReadCloser
	cfg.HTTPClient = smithyhttp.ClientDoFunc(func(req *http.Request) (*http.Response, error) {
		resp, err := baseClient.Do(req)
		if err != nil {
			return nil, err
		}

		body = &closeTrackingReadCloser{ReadCloser: resp.Body}
		resp.Body = body
		return resp, nil
	})

	svc := kinesis.NewFromConfig(cfg)
	_, err = svc.SubscribeToShard(context.Background(), &kinesis.SubscribeToShardInput{}, func(options *kinesis.Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			_, err := stack.Initialize.Remove("OperationInputValidation")
			return err
		})
	})
	if err == nil {
		t.Fatal("expected initial response deserialization error")
	}
	if body == nil {
		t.Fatal("expected response body")
	}
	if !body.closed.Load() {
		t.Error("expected response body to be closed")
	}
}
