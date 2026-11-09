package middleware

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	internalcontext "github.com/aws/aws-sdk-go-v2/internal/context"
	"github.com/aws/smithy-go/logging"
	smithymiddleware "github.com/aws/smithy-go/middleware"
)

func TestInvocationIDContextValue(t *testing.T) {
	ctx := setInvocationID(context.Background(), "00000000-0000-4000-8000-000000000000")

	if e, a := "00000000-0000-4000-8000-000000000000", GetInvocationID(ctx); e != a {
		t.Fatalf("expected invocation id %q, got %q", e, a)
	}

	if a := GetInvocationID(context.Background()); len(a) != 0 {
		t.Fatalf("expected no invocation id on a bare context, got %q", a)
	}

	ctx = setInvocationID(ctx, "other")
	if e, a := "other", GetInvocationID(ctx); e != a {
		t.Fatalf("expected invocation id to be overwritten with %q, got %q", e, a)
	}
}

func TestAttemptNumberContextValue(t *testing.T) {
	ctx := internalcontext.SetAttemptNumberContext(context.Background(), 3)

	if e, a := 3, GetAttemptNumber(ctx); e != a {
		t.Fatalf("expected attempt number %v, got %v", e, a)
	}

	if a := GetAttemptNumber(context.Background()); a != 0 {
		t.Fatalf("expected no attempt number on a bare context, got %v", a)
	}
}

func TestCorrelationContextValuesAreStackValues(t *testing.T) {
	ctx := setInvocationID(context.Background(), "some-id")
	ctx = internalcontext.SetAttemptNumberContext(ctx, 1)

	ctx = smithymiddleware.ClearStackValues(ctx)
	if a := GetInvocationID(ctx); len(a) != 0 {
		t.Errorf("expected invocation id to be cleared, got %q", a)
	}
	if a := GetAttemptNumber(ctx); a != 0 {
		t.Errorf("expected attempt number to be cleared, got %v", a)
	}
}

func TestNewCorrelatedLogger_idempotent(t *testing.T) {
	recorder := &syncRecorder{}
	logger := NewCorrelatedLogger(recorder)

	if a := NewCorrelatedLogger(logger); a != logger {
		t.Fatalf("expected already correlated logger to be returned as is, got %T", a)
	}

	NewCorrelatedLogger(nil).Logf(logging.Debug, "no panic for %v", "nil")
	(&CorrelatedLogger{}).Logf(logging.Debug, "no panic for %v", "zero value")
}

// TestCorrelatedLogger_indexedArguments covers that prepending the prefix does
// not shift the arguments a caller refers to by index.
func TestCorrelatedLogger_indexedArguments(t *testing.T) {
	recorder := &syncRecorder{}
	logger := NewCorrelatedLogger(recorder).WithContext(invocationContextForTest("S3", "PutObject", "abc", 1))

	logger.Logf(logging.Debug, "%[2]v before %[1]v", "second", "first")

	if e, a := "[service=S3 operation=PutObject invocation-id=abc attempt=1] first before second", recorder.lines[0]; e != a {
		t.Fatalf("expected\n%q\ngot\n%q", e, a)
	}
}

// TestCorrelatedLogger_percentInIdentifiers covers that an identifier carrying
// a percent sign is rendered literally, and does not take a verb of the
// caller's format string nor shift its arguments.
func TestCorrelatedLogger_percentInIdentifiers(t *testing.T) {
	cases := map[string]struct {
		invocationID string
		format       string
		args         []interface{}
		expected     string
	}{
		"verb looking identifier": {
			invocationID: "a%vb%%c",
			format:       "Request %v/%v",
			args:         []interface{}{"1", "2"},
			expected:     "[invocation-id=a%vb%%c] Request 1/2",
		},
		"indexed verbs": {
			invocationID: "100%-sure",
			format:       "%[2]v then %[1]v",
			args:         []interface{}{"second", "first"},
			expected:     "[invocation-id=100%-sure] first then second",
		},
		"no verbs": {
			invocationID: "a%vb",
			format:       "plain message",
			expected:     "[invocation-id=a%vb] plain message",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			recorder := &syncRecorder{}
			ctx := setInvocationID(context.Background(), c.invocationID)

			logging.WithContext(ctx, NewCorrelatedLogger(recorder)).
				Logf(logging.Debug, c.format, c.args...)

			if e, a := c.expected, recorder.lines[0]; e != a {
				t.Fatalf("expected\n%q\ngot\n%q", e, a)
			}
		})
	}
}

// TestCorrelatedLogger_concurrentInvocations checks that invocations running
// at the same time are correlated with their own identifiers only.
func TestCorrelatedLogger_concurrentInvocations(t *testing.T) {
	recorder := &syncRecorder{}
	logger := NewCorrelatedLogger(recorder)

	const concurrency = 32
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer wg.Done()

			id := fmt.Sprintf("invocation-%v", i)
			ctx := invocationContextForTest("S3", "PutObject", id, 1)
			bound := logging.WithContext(ctx, logger)

			bound.Logf(logging.Debug, "Request\nPUT /%v HTTP/1.1", id)
			bound.Logf(logging.Debug, "Response\nHTTP/1.1 %v OK", 200)
		}(i)
	}
	wg.Wait()

	if e, a := 2*concurrency, len(recorder.lines); e != a {
		t.Fatalf("expected %v entries, got %v", e, a)
	}

	requests := map[string]string{}
	responses := map[string]string{}
	for _, line := range recorder.lines {
		id, kind, ok := splitCorrelatedLine(line)
		if !ok {
			t.Fatalf("expected correlated entry, got %q", line)
		}
		switch kind {
		case "Request":
			requests[id] = line
		case "Response":
			responses[id] = line
		}
	}

	// every request must have exactly one response carrying the same
	// identifiers, and the other way around
	if !reflect.DeepEqual(keysOf(requests), keysOf(responses)) {
		t.Fatalf("expected request and response identifiers to match,\n%v\n%v", keysOf(requests), keysOf(responses))
	}
	if e, a := concurrency, len(requests); e != a {
		t.Fatalf("expected %v distinct invocations, got %v", e, a)
	}
}

type syncRecorder struct {
	sync.Mutex
	lines []string
}

func (l *syncRecorder) Logf(_ logging.Classification, format string, v ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.lines = append(l.lines, fmt.Sprintf(format, v...))
}

// invocationContextForTest returns the context of an operation invocation,
// assembled the way the SDK's middleware assembles it.
func invocationContextForTest(serviceID, operationName, invocationID string, attempt int) context.Context {
	ctx := SetServiceID(context.Background(), serviceID)
	ctx = SetOperationName(ctx, operationName)
	ctx = setInvocationID(ctx, invocationID)
	return internalcontext.SetAttemptNumberContext(ctx, attempt)
}

func splitCorrelatedLine(line string) (id, kind string, ok bool) {
	const prefix = "[service=S3 operation=PutObject invocation-id="
	if !strings.HasPrefix(line, prefix) {
		return "", "", false
	}
	rest := line[len(prefix):]
	id, rest, ok = strings.Cut(rest, "] ")
	if !ok {
		return "", "", false
	}
	if kind, _, ok = strings.Cut(rest, "\n"); !ok {
		kind = rest
	}
	switch {
	case strings.HasPrefix(kind, "Request"):
		kind = "Request"
	case strings.HasPrefix(kind, "Response"):
		kind = "Response"
	default:
		return "", "", false
	}
	return id, kind, true
}

func keysOf(m map[string]string) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
