package middleware_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	internalcontext "github.com/aws/aws-sdk-go-v2/internal/context"
	"github.com/aws/aws-sdk-go-v2/internal/rand"
	"github.com/aws/smithy-go/logging"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// invocationID is the id the SDK's ClientRequestID middleware generates when
// the SDK's random source is seeded with zeroes, as the tests below do.
const invocationID = "00000000-0000-4000-8000-000000000000"

// recordingLogger captures every entry it is given.
type recordingLogger struct {
	entries []string
}

func (l *recordingLogger) Logf(_ logging.Classification, format string, v ...interface{}) {
	l.entries = append(l.entries, fmt.Sprintf(format, v...))
}

// stdoutLogger writes entries to stdout, and discards the classification.
type stdoutLogger struct{}

func (stdoutLogger) Logf(_ logging.Classification, format string, v ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", v...)
}

// serviceOperationContext returns the context of an operation invocation as
// the client stack's initialize step leaves it.
func serviceOperationContext(serviceID, operationName string) (context.Context, error) {
	var invocationCtx context.Context

	_, _, err := (&middleware.RegisterServiceMetadata{
		ServiceID:     serviceID,
		OperationName: operationName,
	}).HandleInitialize(context.Background(), smithymiddleware.InitializeInput{},
		smithymiddleware.InitializeHandlerFunc(func(ctx context.Context, _ smithymiddleware.InitializeInput) (
			out smithymiddleware.InitializeOutput, m smithymiddleware.Metadata, err error,
		) {
			invocationCtx = ctx
			return out, m, err
		}))

	return invocationCtx, err
}

// withInvocationID returns the context as the build step of the stack, and the
// given retry attempt of the retryer, leave it. The invocation id is generated
// by the SDK's own middleware, from a random source seeded with zeroes. An
// attempt of 0 stands for a stack that records no attempt number.
func withInvocationID(ctx context.Context, attempt int) (context.Context, error) {
	oReader := rand.Reader
	defer func() { rand.Reader = oReader }()
	rand.Reader = bytes.NewReader(make([]byte, 16))

	in := smithymiddleware.BuildInput{
		Request: &smithyhttp.Request{Request: &http.Request{Header: make(http.Header)}},
	}

	var built context.Context
	_, _, err := middleware.ClientRequestID{}.HandleBuild(ctx, in,
		smithymiddleware.BuildHandlerFunc(func(ctx context.Context, _ smithymiddleware.BuildInput) (
			out smithymiddleware.BuildOutput, m smithymiddleware.Metadata, err error,
		) {
			// HandleBuild hands the context it derived to its next handler,
			// which is where the retryer's attempt number is added.
			built = internalcontext.SetAttemptNumberContext(ctx, attempt)
			return out, m, err
		}))

	return built, err
}

// putObjectContext returns the context of a S3 PutObject invocation at the
// given retry attempt, the identifiers being what request/response logging
// observes at runtime.
func putObjectContext(t *testing.T, attempt int) context.Context {
	t.Helper()

	ctx, err := serviceOperationContext("S3", "PutObject")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	ctx, err = withInvocationID(ctx, attempt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id := middleware.GetInvocationID(ctx); id != invocationID {
		t.Fatalf("expected invocation id %q, got %q", invocationID, id)
	}

	return ctx
}

// putObjectPrefix returns the prefix the identifiers of a S3 PutObject
// invocation are logged as. An attempt of 0 leaves the attempt identifier out,
// and the trailing separator, e.g. the space that separates the prefix from the
// message, or nothing for CorrelationPrefix, is appended as requested.
func putObjectPrefix(attempt int, separator string) string {
	prefix := "[service=S3 operation=PutObject invocation-id=" + invocationID
	if attempt != 0 {
		prefix += fmt.Sprintf(" attempt=%v", attempt)
	}
	return prefix + "]" + separator
}

func TestCorrelatedLogger_correlatesRequestAndResponse(t *testing.T) {
	recorder := &recordingLogger{}
	logger := middleware.NewCorrelatedLogger(recorder)

	if _, ok := interface{}(logger).(logging.ContextLogger); !ok {
		t.Fatalf("expected %T to implement logging.ContextLogger", logger)
	}

	bound := logging.WithContext(putObjectContext(t, 2), logger)
	bound.Logf(logging.Debug, "Request\nPUT /%v HTTP/1.1", "my-key")
	bound.Logf(logging.Debug, "Response\nHTTP/1.1 %v Bad Request", 400)

	expected := []string{
		putObjectPrefix(2, " ") + "Request\nPUT /my-key HTTP/1.1",
		putObjectPrefix(2, " ") + "Response\nHTTP/1.1 400 Bad Request",
	}
	if e, a := len(expected), len(recorder.entries); e != a {
		t.Fatalf("expected %v entries, got %v: %#v", e, a, recorder.entries)
	}
	for i := range expected {
		if e, a := expected[i], recorder.entries[i]; e != a {
			t.Errorf("entry %v\nexpected: %q\ngot:      %q", i, e, a)
		}
	}
}

func TestCorrelatedLogger_partialIdentifiers(t *testing.T) {
	serviceOperationCtx, err := serviceOperationContext("S3", "PutObject")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	noAttemptCtx, err := withInvocationID(serviceOperationCtx, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cases := map[string]struct {
		ctx      context.Context
		expected string
	}{
		"nothing available": {
			ctx:      context.Background(),
			expected: "cache warmed",
		},
		"service and operation only": {
			ctx:      serviceOperationCtx,
			expected: "[service=S3 operation=PutObject] cache warmed",
		},
		"no attempt available": {
			ctx:      noAttemptCtx,
			expected: putObjectPrefix(0, " ") + "cache warmed",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			recorder := &recordingLogger{}

			logging.WithContext(c.ctx, middleware.NewCorrelatedLogger(recorder)).
				Logf(logging.Debug, "cache warmed")

			if e, a := c.expected, recorder.entries[0]; e != a {
				t.Fatalf("expected\n%q\ngot\n%q", e, a)
			}
		})
	}
}

func TestCorrelatedLogger_nilLogger(t *testing.T) {
	logger := middleware.NewCorrelatedLogger(nil)

	// Must not panic.
	logging.WithContext(putObjectContext(t, 1), logger).Logf(logging.Debug, "hello %v", "world")
}

func TestCorrelatedLogger_withoutContext(t *testing.T) {
	recorder := &recordingLogger{}

	middleware.NewCorrelatedLogger(recorder).Logf(logging.Debug, "logged outside an invocation")

	if e, a := "logged outside an invocation", recorder.entries[0]; e != a {
		t.Fatalf("expected %q, got %q", e, a)
	}
}

func TestCorrelatedLogger_composesWithContextLogger(t *testing.T) {
	recorder := &recordingLogger{}

	// A wrapped logger that implements logging.ContextLogger must still be
	// given the operation context.
	logger := middleware.NewCorrelatedLogger(contextAwareLogger{inner: recorder})

	logging.WithContext(putObjectContext(t, 2), logger).Logf(logging.Debug, "hello")

	if e, a := "bound "+putObjectPrefix(2, " ")+"hello", recorder.entries[0]; e != a {
		t.Fatalf("expected\n%q\ngot\n%q", e, a)
	}
}

func TestCorrelatedLogger_rebindsContext(t *testing.T) {
	recorder := &recordingLogger{}
	bound := logging.WithContext(putObjectContext(t, 2), middleware.NewCorrelatedLogger(recorder))

	// A bound logger handed to a nested invocation must stamp the nested
	// invocation's identifiers, not the ones it was originally bound to.
	nested := bound.(logging.ContextLogger).WithContext(putObjectContext(t, 1))
	nested.Logf(logging.Debug, "Request")

	if e, a := putObjectPrefix(1, " ")+"Request", recorder.entries[0]; e != a {
		t.Fatalf("expected\n%q\ngot\n%q", e, a)
	}

	// Once the stack values are gone there is nothing to correlate with, and
	// the entry is passed through as it is.
	cleared := bound.(logging.ContextLogger).WithContext(smithymiddleware.ClearStackValues(putObjectContext(t, 2)))
	cleared.Logf(logging.Debug, "Request")

	if e, a := "Request", recorder.entries[1]; e != a {
		t.Fatalf("expected\n%q\ngot\n%q", e, a)
	}
}

func TestCorrelationPrefix(t *testing.T) {
	cases := map[string]struct {
		ctx      context.Context
		expected string
	}{
		"empty": {
			ctx:      context.Background(),
			expected: "",
		},
		"all identifiers": {
			ctx:      putObjectContext(t, 2),
			expected: putObjectPrefix(2, ""),
		},
		"first attempt": {
			ctx: func() context.Context {
				ctx := middleware.SetServiceID(context.Background(), "S3")
				return internalcontext.SetAttemptNumberContext(ctx, 1)
			}(),
			expected: "[service=S3 attempt=1]",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.expected, middleware.CorrelationPrefix(c.ctx); e != a {
				t.Fatalf("expected %q, got %q", e, a)
			}
		})
	}
}

func ExampleNewCorrelatedLogger() {
	ctx, err := serviceOperationContext("S3", "PutObject")
	if err != nil {
		panic(err)
	}
	ctx, err = withInvocationID(ctx, 2)
	if err != nil {
		panic(err)
	}

	logger := middleware.NewCorrelatedLogger(stdoutLogger{})

	logging.WithContext(ctx, logger).Logf(logging.Debug, "Request\nPUT /%v HTTP/1.1", "my-key")

	// Entries logged without an operation context bound, e.g. while resolving
	// configuration, are passed through as they are.
	logger.Logf(logging.Debug, "resolved region %v", "us-east-2")

	// Output:
	// [service=S3 operation=PutObject invocation-id=00000000-0000-4000-8000-000000000000 attempt=2] Request
	// PUT /my-key HTTP/1.1
	// resolved region us-east-2
}

// contextAwareLogger is a logging.ContextLogger that tags entries once it has
// been bound to a context.
type contextAwareLogger struct {
	inner *recordingLogger
}

func (l contextAwareLogger) Logf(_ logging.Classification, format string, v ...interface{}) {
	l.inner.entries = append(l.inner.entries, fmt.Sprintf(format, v...))
}

func (l contextAwareLogger) WithContext(context.Context) logging.Logger {
	return boundLogger{inner: l.inner}
}

type boundLogger struct {
	inner *recordingLogger
}

func (l boundLogger) Logf(_ logging.Classification, format string, v ...interface{}) {
	l.inner.entries = append(l.inner.entries, fmt.Sprintf("bound "+format, v...))
}
