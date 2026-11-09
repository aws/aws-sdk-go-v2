package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// BenchmarkClientRequestID measures the cost of the build step that generates
// the invocation id, records it as the Amz-Sdk-Invocation-Id header, and mirrors
// it onto the context. The context write happens once per operation, and the
// retryer records one comparable context value per retry attempt.
func BenchmarkClientRequestID(b *testing.B) {
	ctx := SetServiceID(context.Background(), "S3")
	in := middleware.BuildInput{
		Request: &smithyhttp.Request{Request: &http.Request{Header: make(http.Header)}},
	}
	next := middleware.BuildHandlerFunc(func(ctx context.Context, in middleware.BuildInput) (
		middleware.BuildOutput, middleware.Metadata, error,
	) {
		return middleware.BuildOutput{Result: in.Request}, middleware.Metadata{}, nil
	})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, _, err := (ClientRequestID{}).HandleBuild(ctx, in, next); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCorrelatedLoggerLogf measures the per entry cost of decorating an
// entry, on top of whatever the wrapped logger does.
func BenchmarkCorrelatedLoggerLogf(b *testing.B) {
	logger := NewCorrelatedLogger(discardLogger{}).
		WithContext(invocationContextForTest("S3", "PutObject", invocationIDForBenchmark, 1)).(*CorrelatedLogger)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Logf(logging.Debug, "Request\nPUT /%v HTTP/1.1", "my-key")
	}
}

// BenchmarkCorrelatedLoggerLogf_unbound measures the cost of an entry that
// carries no correlation, i.e. a logger that was never bound to a context.
func BenchmarkCorrelatedLoggerLogf_unbound(b *testing.B) {
	logger := NewCorrelatedLogger(discardLogger{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Logf(logging.Debug, "Request\nPUT /%v HTTP/1.1", "my-key")
	}
}

const invocationIDForBenchmark = "00000000-0000-4000-8000-000000000000"

type discardLogger struct{}

func (discardLogger) Logf(logging.Classification, string, ...interface{}) {}
