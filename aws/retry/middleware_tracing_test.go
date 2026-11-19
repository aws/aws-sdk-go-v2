package retry

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/tracing"
)

func TestAttemptRetryLoopSpans(t *testing.T) {
	restoreSleep := sdk.TestingUseNopSleep()
	defer restoreSleep()
	for _, parentName := range []string{"", "Operation", "RetryLoop"} {
		for _, tc := range []struct {
			name      string
			failures  int
			wantCalls int
		}{{"success", 0, 1}, {"retry", 1, 2}, {"exhaustion", 2, 2}} {
			t.Run(parentName+"/"+tc.name, func(t *testing.T) {
				tracer := &retrySpanRecorder{}
				ctx := tracing.WithOperationTracer(context.Background(), tracer)
				var parent tracing.Span
				if parentName != "" {
					ctx, parent = tracer.StartSpan(ctx, parentName)
					defer parent.End()
				}
				am := NewAttemptMiddleware(NewStandard(func(o *StandardOptions) { o.MaxAttempts = 2 }), func(v interface{}) interface{} { return v })
				calls := 0
				_, metadata, err := am.HandleFinalize(ctx, middleware.FinalizeInput{Request: testRequest{}}, middleware.FinalizeHandlerFunc(
					func(ctx context.Context, _ middleware.FinalizeInput) (middleware.FinalizeOutput, middleware.Metadata, error) {
						calls++
						current, ok := tracing.GetSpan(ctx)
						if !ok || current.Name() != "Attempt" {
							t.Error("handler must run inside an attempt span")
						}
						if calls <= tc.failures {
							return middleware.FinalizeOutput{}, middleware.Metadata{}, mockRetryableError{b: true}
						}
						return middleware.FinalizeOutput{Result: "success"}, middleware.Metadata{}, nil
					}))
				if calls != tc.wantCalls {
					t.Errorf("expected %d calls, got %d", tc.wantCalls, calls)
				}
				if tc.failures == 2 {
					var exhausted *MaxAttemptsError
					if !errors.As(err, &exhausted) {
						t.Errorf("expected exhausted retries, got %v", err)
					}
				} else if err != nil {
					t.Fatal(err)
				}
				if results, ok := GetAttemptResults(metadata); !ok || len(results.Results) != calls {
					t.Error("attempt metadata does not match calls")
				}
				var loop *retryRecordedSpan
				attempts := 0
				for _, span := range tracer.spans {
					if span.Name() == "RetryLoop" {
						if loop != nil {
							t.Error("multiple RetryLoop spans")
						}
						loop = span
					}
					if span.Name() == "Attempt" {
						attempts++
					}
				}
				if loop == nil {
					t.Fatal("missing RetryLoop span")
				}
				if attempts != calls {
					t.Errorf("expected %d attempt spans, got %d", calls, attempts)
				}
				for _, span := range tracer.spans {
					if span.Name() == "Attempt" && (span.parent != loop || span.ends != 1) {
						t.Error("attempt must be a closed direct child of RetryLoop")
					}
				}
				wantEnds := 1
				if parentName == "RetryLoop" {
					wantEnds = 0
				}
				if loop.ends != wantEnds {
					t.Errorf("expected loop ended %d times by retry middleware, got %d", wantEnds, loop.ends)
				}
			})
		}
	}
}

type retrySpanRecorder struct{ spans []*retryRecordedSpan }

func (r *retrySpanRecorder) StartSpan(ctx context.Context, name string, _ ...tracing.SpanOption) (context.Context, tracing.Span) {
	parent, _ := tracing.GetSpan(ctx)
	_, noop := tracing.NopTracerProvider{}.Tracer("").StartSpan(ctx, name)
	span := &retryRecordedSpan{Span: noop, name: name, parent: parent}
	r.spans = append(r.spans, span)
	return tracing.WithSpan(ctx, span), span
}

type retryRecordedSpan struct {
	tracing.Span
	name   string
	parent tracing.Span
	ends   int
}

func (s *retryRecordedSpan) Name() string { return s.name }
func (s *retryRecordedSpan) End()         { s.ends++ }
