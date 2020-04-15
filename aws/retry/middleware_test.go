package retry

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/awslabs/smithy-go/middleware"
	smithyHTTP "github.com/awslabs/smithy-go/transport/http"
)

type mockFinalizeHandler func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error)

func (f mockFinalizeHandler) HandleFinalize(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
	return f(ctx, in)
}

func TestMetricsHeaderMiddleware(t *testing.T) {
	cases := []struct {
		input          middleware.FinalizeInput
		ctx            context.Context
		expectedHeader string
		expectedErr    string
	}{
		{
			input: middleware.FinalizeInput{Request: &smithyHTTP.Request{Request: &http.Request{Header: make(http.Header)}}},
			ctx: func() context.Context {
				return setRetryMetadata(context.Background(), retryMetadata{
					AttemptNum:       0,
					AttemptTime:      time.Date(2020, 01, 02, 03, 04, 05, 0, time.UTC),
					MaxAttempts:      5,
					AttemptClockSkew: 0,
				})
			}(),
			expectedHeader: "attempt=0; max=5",
		},
		{
			input: middleware.FinalizeInput{Request: &smithyHTTP.Request{Request: &http.Request{Header: make(http.Header)}}},
			ctx: func() context.Context {
				attemptTime := time.Date(2020, 01, 02, 03, 04, 05, 0, time.UTC)
				ctx, cancel := context.WithDeadline(context.Background(), attemptTime.Add(time.Minute))
				defer cancel()
				return setRetryMetadata(ctx, retryMetadata{
					AttemptNum:       1,
					AttemptTime:      attemptTime,
					MaxAttempts:      5,
					AttemptClockSkew: time.Second * 1,
				})
			}(),
			expectedHeader: "attempt=1; max=5; ttl=20200102T030506Z",
		},
		{
			ctx: func() context.Context {
				return setRetryMetadata(context.Background(), retryMetadata{})
			}(),
			expectedErr: "unknown transport type",
		},
	}

	retryMiddleware := MetricsHeaderMiddleware{}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := tt.ctx
			_, _, err := retryMiddleware.HandleFinalize(ctx, tt.input, mockFinalizeHandler(func(ctx context.Context, in middleware.FinalizeInput) (
				out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
			) {
				req := in.Request.(*smithyHTTP.Request)

				if e, a := tt.expectedHeader, req.Header.Get("amz-sdk-request"); e != a {
					t.Errorf("expected %v, got %v", e, a)
				}

				return out, metadata, err
			}))
			if err != nil && len(tt.expectedErr) == 0 {
				t.Fatalf("expected no error, got %q", err)
			} else if err != nil && len(tt.expectedErr) != 0 {
				if e, a := tt.expectedErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expected %q, got %q", e, a)
				}
			} else if err == nil && len(tt.expectedErr) != 0 {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}
