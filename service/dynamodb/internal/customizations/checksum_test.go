package customizations

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestCRC32ChecksumValidate(t *testing.T) {
	cases := map[string]struct {
		Reader      io.ReadCloser
		ExpectCRC32 uint32
		ExpectErr   string
	}{
		"empty reader": {
			Reader:      io.NopCloser(&bytes.Buffer{}),
			ExpectCRC32: 0,
		},
		"wrong checksum": {
			Reader:      io.NopCloser(bytes.NewBuffer([]byte("abc123"))),
			ExpectCRC32: 123456,
			ExpectErr:   "did not match",
		},
		"with closer": {
			Reader: &wasClosedReadCloser{
				Reader: bytes.NewBuffer([]byte("abc123")),
			},
			ExpectCRC32: 0xcf02bb5c,
		},
		"without closer": {
			Reader:      io.NopCloser(bytes.NewBuffer([]byte("abc123"))),
			ExpectCRC32: 0xcf02bb5c,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			reader := wrapCRC32ChecksumValidate(c.ExpectCRC32, c.Reader)
			// Asserts
			io.Copy(io.Discard, reader)

			err := reader.Close()
			if len(c.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Errorf("expect error to contain %v, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if c, ok := c.Reader.(interface{ WasClosed() bool }); ok {
				if !c.WasClosed() {
					t.Errorf("expect original reader closed, but was not")
				}
			}
		})
	}
}

func TestChecksumHandleDeserialize(t *testing.T) {
	cases := map[string]struct {
		Header     http.Header
		ExpectWrap bool
	}{
		"header present": {
			Header: http.Header{
				crc32ChecksumHeader: []string{"3162747320"},
			},
			ExpectWrap: true,
		},
		"header absent": {
			Header:     http.Header{},
			ExpectWrap: false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			m := &Checksum{}

			output, _, err := m.HandleDeserialize(context.Background(),
				middleware.DeserializeInput{},
				middleware.DeserializeHandlerFunc(
					func(ctx context.Context, input middleware.DeserializeInput) (
						output middleware.DeserializeOutput, metadata middleware.Metadata, err error,
					) {
						output.RawResponse = &smithyhttp.Response{
							Response: &http.Response{
								StatusCode: 200,
								Header:     c.Header,
								Body:       io.NopCloser(bytes.NewBufferString("abc123")),
							},
						}
						return output, metadata, err
					}),
			)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			resp, ok := output.RawResponse.(*smithyhttp.Response)
			if !ok || resp == nil {
				t.Fatalf("expect smithy response, got %T", output.RawResponse)
			}

			_, wrapped := resp.Body.(*crc32ChecksumValidate)
			if e, a := c.ExpectWrap, wrapped; e != a {
				t.Errorf("expect wrap %v, got %v", e, a)
			}
		})
	}
}

type wasClosedReadCloser struct {
	io.Reader
	closed bool
}

func (c *wasClosedReadCloser) WasClosed() bool {
	return c.closed
}

func (c *wasClosedReadCloser) Close() error {
	c.closed = true
	if v, ok := c.Reader.(io.Closer); ok {
		return v.Close()
	}
	return nil
}
