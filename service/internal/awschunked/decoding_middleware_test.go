package awschunked

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// decodeStep is a minimal Deserialize middleware used only by this test to
// exercise the decoding reader end-to-end through a smithy middleware stack. It
// mirrors how a response-path decode step engages: it inspects the HTTP
// response and wraps the body with the decoding reader so the caller reads only
// decoded payload bytes. It carries no knowledge of what the trailer contains.
type decodeStep struct {
	decodedLength int64
	// capture, if set, receives the reader so the test can inspect it mid-stream.
	capture func(*Reader)
}

func (decodeStep) ID() string { return "test:DecodeAWSChunked" }

func (m decodeStep) HandleDeserialize(
	ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler,
) (middleware.DeserializeOutput, middleware.Metadata, error) {
	out, md, err := next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, md, err
	}

	resp, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok || resp.Body == nil {
		return out, md, err
	}

	reader := NewReader(resp.Body, func(o *Options) {
		o.DecodedLength = m.decodedLength
	})
	if m.capture != nil {
		m.capture(reader)
	}
	resp.Body = reader
	return out, md, err
}

func runDecodeStep(
	t *testing.T, step decodeStep, rawBody string, oneByte bool,
) *smithyhttp.Response {
	t.Helper()

	var body io.Reader = strings.NewReader(rawBody)
	if oneByte {
		body = iotest.OneByteReader(body)
	}

	response := &smithyhttp.Response{
		Response: &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
			Body:       io.NopCloser(body),
		},
	}

	out, _, err := step.HandleDeserialize(context.Background(),
		middleware.DeserializeInput{},
		middleware.DeserializeHandlerFunc(
			func(context.Context, middleware.DeserializeInput) (
				o middleware.DeserializeOutput, md middleware.Metadata, err error,
			) {
				o.RawResponse = response
				return o, md, nil
			},
		),
	)
	if err != nil {
		t.Fatalf("expect no deserialize error, got %v", err)
	}
	return out.RawResponse.(*smithyhttp.Response)
}

// TestAWSChunkedDecodingThroughMiddleware drives the decoding reader end-to-end
// through a smithy Deserialize step against a synthetic HTTP response, rather
// than calling the reader directly. The caller reads the decoded payload off
// the response body exactly as it would after a real client invocation, and
// the captured trailers are observed once the body reaches EOF.
func TestAWSChunkedDecodingThroughMiddleware(t *testing.T) {
	cases := map[string]struct {
		rawBody       string
		decodedLength int64
		oneByte       bool
		expectPayload string
		expectTrailer map[string]string
	}{
		"single chunk with trailer": {
			rawBody:       "b\r\nHello world\r\n0\r\ntrailer-key:abc123==\r\n\r\n",
			decodedLength: 11,
			expectPayload: "Hello world",
			expectTrailer: map[string]string{"Trailer-Key": "abc123=="},
		},
		"multiple chunks no trailer": {
			rawBody:       "5\r\nhello\r\n6\r\n world\r\n0\r\n\r\n",
			decodedLength: 11,
			expectPayload: "hello world",
		},
		"one byte at a time": {
			rawBody:       "b\r\nHello world\r\n0\r\ntrailer-key:abc123==\r\n\r\n",
			decodedLength: 11,
			oneByte:       true,
			expectPayload: "Hello world",
			expectTrailer: map[string]string{"Trailer-Key": "abc123=="},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var reader *Reader
			step := decodeStep{
				decodedLength: c.decodedLength,
				capture:       func(r *Reader) { reader = r },
			}

			resp := runDecodeStep(t, step, c.rawBody, c.oneByte)

			payload, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("expect no read error, got %v", err)
			}
			if got, want := string(payload), c.expectPayload; got != want {
				t.Errorf("payload = %q, want %q", got, want)
			}
			if err := resp.Body.Close(); err != nil {
				t.Errorf("expect no close error, got %v", err)
			}

			// Trailers are only available after the body is fully consumed.
			for k, want := range c.expectTrailer {
				if got := reader.Trailers().Get(k); got != want {
					t.Errorf("trailer %q = %q, want %q", k, got, want)
				}
			}
		})
	}
}

// TestAWSChunkedDecodingTrailersUnavailableUntilEOF asserts that trailer fields
// are not observable until the decoded body has been fully consumed: a partial
// read exposes no trailers. Once the remaining bytes are read through EOF, the
// trailers become available.
func TestAWSChunkedDecodingTrailersUnavailableUntilEOF(t *testing.T) {
	const raw = "b\r\nHello world\r\n0\r\ntrailer-key:abc123==\r\n\r\n"

	var reader *Reader
	step := decodeStep{
		decodedLength: 11,
		capture:       func(r *Reader) { reader = r },
	}

	resp := runDecodeStep(t, step, raw, false)

	// Partial read: consume only the first 5 bytes of the 11-byte payload.
	buf := make([]byte, 5)
	if _, err := io.ReadFull(resp.Body, buf); err != nil {
		t.Fatalf("unexpected error on partial read: %v", err)
	}
	if got := string(buf); got != "Hello" {
		t.Fatalf("partial payload = %q, want %q", got, "Hello")
	}

	// Before EOF the trailers must not be available.
	if trailers := reader.Trailers(); len(trailers) != 0 {
		t.Errorf("trailers available before EOF: %v", trailers)
	}

	// Consume the remainder through EOF.
	rest, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error reading remainder: %v", err)
	}
	if got := string(rest); got != " world" {
		t.Fatalf("remaining payload = %q, want %q", got, " world")
	}

	// After EOF the trailers are available.
	if got, want := reader.Trailers().Get("trailer-key"), "abc123=="; got != want {
		t.Errorf("trailer after EOF = %q, want %q", got, want)
	}

	if err := resp.Body.Close(); err != nil {
		t.Errorf("expect no close error, got %v", err)
	}
}
