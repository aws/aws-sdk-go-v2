package kitchensinktest

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// TestRequestCompression_ContentLengthMatchesCompressedBody asserts that for a
// request-compression operation the Content-Length equals the number of bytes
// actually sent (the gzip-compressed body).
//
// Request compression runs as a Serialize-step middleware that replaces the
// body with its compressed form after the body is serialized. Content length
// must therefore be computed after compression. If it is computed inline in the
// serializer (before compression runs), the request advertises the uncompressed
// length, which is larger than the bytes actually sent — net/http then hangs or
// the server rejects the request. No other test exercises a compression
// operation, so this guards against that regression.
func TestRequestCompression_ContentLengthMatchesCompressedBody(t *testing.T) {
	// Large, repetitive payload: exceeds the compression threshold and
	// compresses to meaningfully fewer bytes than the input.
	payload := strings.Repeat("kitchensink request compression content length regression guard. ", 512)

	var gotContentLength int64
	var gotBody []byte
	var gotEncoding string

	svc := New(Options{
		Region: "us-east-1",
		HTTPClient: smithyhttp.ClientDoFunc(func(req *http.Request) (*http.Response, error) {
			gotContentLength = req.ContentLength
			gotEncoding = req.Header.Get("Content-Encoding")
			if req.Body != nil {
				b, err := io.ReadAll(req.Body)
				if err != nil {
					return nil, err
				}
				gotBody = b
			}
			return &http.Response{
				StatusCode:    200,
				Header:        http.Header{},
				Body:          http.NoBody,
				ContentLength: -1,
				Request:       req,
			}, nil
		}),
		EndpointResolverV2:          &closeTestEndpointResolver{},
		DisableRequestCompression:   false,
		RequestMinCompressSizeBytes: 0,
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				s.Finalize.Clear()
				return nil
			},
		},
	})

	_, err := svc.PutCompressedData(context.Background(), &PutCompressedDataInput{
		Data: aws.String(payload),
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	// Confirm compression actually happened, otherwise the assertion below
	// would pass trivially.
	if !strings.Contains(gotEncoding, "gzip") {
		t.Fatalf("expected gzip-compressed request, Content-Encoding=%q", gotEncoding)
	}
	if len(gotBody) == 0 {
		t.Fatal("expected a non-empty request body")
	}
	if len(gotBody) >= len(payload) {
		t.Fatalf("expected compressed body (%d bytes) smaller than payload (%d bytes)",
			len(gotBody), len(payload))
	}

	if int64(len(gotBody)) != gotContentLength {
		t.Errorf("Content-Length does not match compressed body size: Content-Length=%d, actual body=%d bytes",
			gotContentLength, len(gotBody))
	}
}
