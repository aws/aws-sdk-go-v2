package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/internal/rand"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/logging"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// TestCorrelatedLogger_overMiddlewareStack wires the middleware a client stack
// uses for the invocation id, the logger, retries, and request/response
// logging, over a transport that fails the first attempt, to verify that every
// entry logged for a retried request carries the identifiers of the invocation
// and attempt it belongs to, and that the response is tied to the attempt that
// produced it.
func TestCorrelatedLogger_overMiddlewareStack(t *testing.T) {
	restoreSleep := sdk.TestingUseNopSleep()
	defer restoreSleep()

	oReader := rand.Reader
	defer func() { rand.Reader = oReader }()
	// the retryer draws from the same source for jitter, so use a source that
	// never runs dry
	rand.Reader = zeroReader{}

	recorder := &recordingLogger{}

	stack := smithymiddleware.NewStack("test", smithyhttp.NewStackRequest)
	if err := smithymiddleware.AddSetLoggerMiddleware(stack, middleware.NewCorrelatedLogger(recorder)); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := stack.Initialize.Add(
		&middleware.RegisterServiceMetadata{ServiceID: "S3", OperationName: "PutObject"}, smithymiddleware.After); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := stack.Serialize.Add(serializeRequest{}, smithymiddleware.After); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := middleware.AddClientRequestIDMiddleware(stack); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	attempter := retry.NewAttemptMiddleware(retry.NewStandard(), func(i interface{}) interface{} { return i },
		func(a *retry.Attempt) { a.LogAttempts = true })
	if err := stack.Finalize.Add(attempter, smithymiddleware.Before); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := stack.Deserialize.Add(&smithyhttp.RequestResponseLogger{
		LogRequest:  true,
		LogResponse: true,
	}, smithymiddleware.After); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	handler := smithymiddleware.DecorateHandler(retryTransportHandler{}, stack)
	if _, _, err := handler.Handle(context.Background(), nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// The first attempt logs a request that gets no response, the retry
	// request and its response must carry the attempt number that produced
	// them, and every entry must carry the same invocation id.
	expected := []struct {
		prefix string
		body   string
	}{
		{putObjectPrefix(1, " "), "Request\n"},
		{putObjectPrefix(2, " "), "retrying request S3/PutObject, attempt 2"},
		{putObjectPrefix(2, " "), "Request\n"},
		{putObjectPrefix(2, " "), "Response\nHTTP/1.1 200"},
	}
	if e, a := len(expected), len(recorder.entries); e != a {
		t.Fatalf("expected %v entries, got %v: %#v", e, a, recorder.entries)
	}
	for i, c := range expected {
		entry := recorder.entries[i]
		if !strings.HasPrefix(entry, c.prefix+c.body) {
			t.Errorf("entry %v\nexpected prefix: %q\ngot:             %q", i, c.prefix+c.body, entry)
		}
	}
}

// TestCorrelatedLogger_prefixIsFirstLineOnly documents the scope of the
// correlation prefix on multi-line entries: the identifiers are on the first
// line, the dumped headers are not attributed.
func TestCorrelatedLogger_prefixIsFirstLineOnly(t *testing.T) {
	recorder := &recordingLogger{}

	logging.WithContext(putObjectContext(t, 1), middleware.NewCorrelatedLogger(recorder)).
		Logf(logging.Debug, "Request\nPUT /key HTTP/1.1\nHost: s3.amazonaws.com")

	line, _, _ := strings.Cut(recorder.entries[0], "\n")
	if e, a := putObjectPrefix(1, " ")+"Request", line; e != a {
		t.Fatalf("expected %q, got %q", e, a)
	}
	if strings.Contains(recorder.entries[0][len(line):], "invocation-id") {
		t.Fatalf("expected the prefix on the first line only, got %q", recorder.entries[0])
	}
}

// retryTransportHandler stands in for the HTTP transport: it fails the first
// attempt with a retryable error, and returns a response for the retry.
type retryTransportHandler struct{}

func (retryTransportHandler) Handle(ctx context.Context, input interface{}) (
	output interface{}, metadata smithymiddleware.Metadata, err error,
) {
	// what is logged must be what is sent
	if id, header := middleware.GetInvocationID(ctx), requestHeader(input, "Amz-Sdk-Invocation-Id"); id != header {
		panic(fmt.Sprintf("expected invocation id %q on context, header carries %q", header, id))
	}

	if middleware.GetAttemptNumber(ctx) < 2 {
		return output, metadata, retryableError{}
	}

	return &smithyhttp.Response{Response: &http.Response{
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
		Body:       http.NoBody,
	}}, metadata, err
}

// requestHeader returns the value of the given header of the stack request.
func requestHeader(input interface{}, key string) string {
	switch v := input.(type) {
	case *smithyhttp.Request:
		return v.Header.Get(key)
	case *http.Request:
		return v.Header.Get(key)
	}
	return ""
}

// zeroReader is a deterministic source of zeroes for the SDK's random reader.
type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

type retryableError struct{}

func (retryableError) Error() string        { return "mock retryable error" }
func (retryableError) RetryableError() bool { return true }

// serializeRequest turns the stack request into a dumpable HTTP request.
type serializeRequest struct{}

func (serializeRequest) ID() string { return "SerializeRequest" }

func (serializeRequest) HandleSerialize(ctx context.Context, in smithymiddleware.SerializeInput, next smithymiddleware.SerializeHandler) (
	out smithymiddleware.SerializeOutput, metadata smithymiddleware.Metadata, err error,
) {
	req := in.Request.(*smithyhttp.Request)
	req.Method = http.MethodPut
	req.URL = &url.URL{Scheme: "https", Host: "s3.amazonaws.com", Path: "/key"}
	req.Host = "s3.amazonaws.com"
	return next.HandleSerialize(ctx, in)
}
