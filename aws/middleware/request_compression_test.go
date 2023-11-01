package middleware

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	awsutil2 "github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"io"
	"strings"
	"testing"
)

func TestRequestCompression(t *testing.T) {
	cases := map[string]struct {
		DisableRequestCompression   bool
		RequestMinCompressSizeBytes int64
		Body                        io.ReadCloser
		ContentLength               int64
		Header                      map[string][]string
		Stream                      io.Reader
		ExpectedBody                []byte
		ExpectedStream              []byte
		ExpectedHeader              map[string][]string
	}{
		"GZip request stream": {
			Stream:         strings.NewReader("Hi, world!"),
			ExpectedStream: []byte("Hi, world!"),
			ExpectedHeader: map[string][]string{
				"Content-Encoding": {"gzip"},
			},
		},
		"GZip request body": {
			RequestMinCompressSizeBytes: 0,
			Body:                        io.NopCloser(strings.NewReader("Hello, world!")),
			ContentLength:               13,
			ExpectedBody:                []byte("Hello, world!"),
			ExpectedHeader: map[string][]string{
				"Content-Encoding": {"gzip"},
			},
		},
		"GZip request body with existing encoding header": {
			RequestMinCompressSizeBytes: 0,
			Body:                        io.NopCloser(strings.NewReader("Hello, world!")),
			ContentLength:               13,
			Header: map[string][]string{
				"Content-Encoding": {"custom"},
			},
			ExpectedBody: []byte("Hello, world!"),
			ExpectedHeader: map[string][]string{
				"Content-Encoding": {"custom,gzip"},
			},
		},
		"GZip request stream ignoring min compress request size": {
			RequestMinCompressSizeBytes: 100,
			Stream:                      strings.NewReader("Hi, world!"),
			ExpectedStream:              []byte("Hi, world!"),
			ExpectedHeader: map[string][]string{
				"Content-Encoding": {"gzip"},
			},
		},
		"Disable GZip request stream": {
			DisableRequestCompression: true,
			Stream:                    strings.NewReader("Hi, world!"),
			ExpectedStream:            []byte("Hi, world!"),
			ExpectedHeader:            map[string][]string{},
		},
		"Disable GZip request body": {
			DisableRequestCompression:   true,
			RequestMinCompressSizeBytes: 0,
			Body:                        io.NopCloser(strings.NewReader("Hello, world!")),
			ContentLength:               13,
			ExpectedBody:                []byte("Hello, world!"),
			ExpectedHeader:              map[string][]string{},
		},
		"Disable Gzip request body due to size threshold": {
			RequestMinCompressSizeBytes: 14,
			Body:                        io.NopCloser(strings.NewReader("Hello, world!")),
			ContentLength:               13,
			ExpectedBody:                []byte("Hello, world!"),
			ExpectedHeader:              map[string][]string{},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var err error
			req := smithyhttp.NewStackRequest().(*smithyhttp.Request)
			req.ContentLength = c.ContentLength
			req.Body = c.Body
			req, _ = req.SetStream(c.Stream)
			if c.Header != nil {
				req.Header = c.Header
			}
			var updatedRequest *smithyhttp.Request

			m := requestCompression{
				disableRequestCompression:   c.DisableRequestCompression,
				requestMinCompressSizeBytes: c.RequestMinCompressSizeBytes,
			}
			_, _, err = m.HandleBuild(context.Background(),
				middleware.BuildInput{Request: req},
				middleware.BuildHandlerFunc(func(ctx context.Context, input middleware.BuildInput) (
					out middleware.BuildOutput, metadata middleware.Metadata, err error) {
					updatedRequest = input.Request.(*smithyhttp.Request)
					return out, metadata, nil
				}),
			)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if stream := updatedRequest.GetStream(); stream != nil {
				if err := testUnzipContent(stream, c.ExpectedStream, c.DisableRequestCompression); err != nil {
					t.Errorf("error while checking request stream: %q", err)
				}
			}

			if body := updatedRequest.Body; body != nil {
				if c.RequestMinCompressSizeBytes > c.ContentLength {
					bodyBytes, err := io.ReadAll(body)
					if err != nil {
						t.Errorf("error while reading request body")
					}
					if e, a := c.ExpectedBody, bodyBytes; !bytes.Equal(e, a) {
						t.Errorf("expect body to be %s, got %s", e, a)
					}
				} else if err := testUnzipContent(body, c.ExpectedBody, c.DisableRequestCompression); err != nil {
					t.Errorf("error while checking request body: %q", err)
				}
			}

			if e, a := c.ExpectedHeader, map[string][]string(updatedRequest.Header); !awsutil2.DeepEqual(e, a) {
				t.Errorf("expect request header to be %q, got %q", e, a)
			}
		})
	}
}

func testUnzipContent(content io.Reader, expect []byte, disableRequestCompression bool) error {
	if disableRequestCompression {
		b, err := io.ReadAll(content)
		if err != nil {
			return fmt.Errorf("error while reading request")
		}
		if e, a := expect, b; !bytes.Equal(e, a) {
			return fmt.Errorf("expect content to be %s, got %s", e, a)
		}
	} else {
		r, err := gzip.NewReader(content)
		if err != nil {
			return fmt.Errorf("error while reading request")
		}

		var actualBytes bytes.Buffer
		_, err = actualBytes.ReadFrom(r)
		if err != nil {
			return fmt.Errorf("error while unzipping request payload")
		}

		if e, a := expect, actualBytes.Bytes(); !bytes.Equal(e, a) {
			return fmt.Errorf("expect unzipped content to be %s, got %s", e, a)
		}
	}

	return nil
}
