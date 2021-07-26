package checksum

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/google/go-cmp/cmp"
)

func TestValidateOutputPayloadChecksum(t *testing.T) {
	cases := map[string]struct {
		response                 *smithyhttp.Response
		validateOptions          func(*validateOutputPayloadChecksum)
		expectHaveAlgorithmsUsed bool
		expectAlgorithmsUsed     []string
		expectErr                string
		expectReadErr            string
		expectLogged             string
		expectPayload            []byte
	}{
		"success": {
			response: &smithyhttp.Response{
				Response: &http.Response{
					StatusCode: 200,
					Header: func() http.Header {
						h := http.Header{}
						h.Set(AlgorithmHTTPHeader(AlgorithmCRC32), "DUoRhQ==")
						return h
					}(),
					Body: ioutil.NopCloser(strings.NewReader("hello world")),
				},
			},
			expectHaveAlgorithmsUsed: true,
			expectAlgorithmsUsed:     []string{"CRC32"},
			expectPayload:            []byte("hello world"),
		},
		"failure": {
			response: &smithyhttp.Response{
				Response: &http.Response{
					StatusCode: 200,
					Header: func() http.Header {
						h := http.Header{}
						h.Set(AlgorithmHTTPHeader(AlgorithmCRC32), "AAAAAA==")
						return h
					}(),
					Body: ioutil.NopCloser(strings.NewReader("hello world")),
				},
			},
			expectReadErr: "checksum did not match",
		},
		"read error": {
			response: &smithyhttp.Response{
				Response: &http.Response{
					StatusCode: 200,
					Header: func() http.Header {
						h := http.Header{}
						h.Set(AlgorithmHTTPHeader(AlgorithmCRC32), "AAAAAA==")
						return h
					}(),
					Body: ioutil.NopCloser(iotest.ErrReader(fmt.Errorf("some read error"))),
				},
			},
			expectReadErr: "some read error",
		},
		"unsupported algorithm": {
			response: &smithyhttp.Response{
				Response: &http.Response{
					StatusCode: 200,
					Header: func() http.Header {
						h := http.Header{}
						h.Set(AlgorithmHTTPHeader("unsupported"), "AAAAAA==")
						return h
					}(),
					Body: ioutil.NopCloser(strings.NewReader("hello world")),
				},
			},
			expectLogged:  "no supported checksum",
			expectPayload: []byte("hello world"),
		},
		"success ignore multipart checksum": {
			response: &smithyhttp.Response{
				Response: &http.Response{
					StatusCode: 200,
					Header: func() http.Header {
						h := http.Header{}
						h.Set(AlgorithmHTTPHeader(AlgorithmCRC32), "DUoRhQ==")
						return h
					}(),
					Body: ioutil.NopCloser(strings.NewReader("hello world")),
				},
			},
			validateOptions: func(o *validateOutputPayloadChecksum) {
				o.IgnoreMultipartValidation = true
			},
			expectHaveAlgorithmsUsed: true,
			expectAlgorithmsUsed:     []string{"CRC32"},
			expectPayload:            []byte("hello world"),
		},
		"success skip ignore multipart checksum": {
			response: &smithyhttp.Response{
				Response: &http.Response{
					StatusCode: 200,
					Header: func() http.Header {
						h := http.Header{}
						h.Set(AlgorithmHTTPHeader(AlgorithmCRC32), "DUoRhQ==-12")
						return h
					}(),
					Body: ioutil.NopCloser(strings.NewReader("hello world")),
				},
			},
			validateOptions: func(o *validateOutputPayloadChecksum) {
				o.IgnoreMultipartValidation = true
			},
			expectLogged:  "Skipped validation of multipart checksum",
			expectPayload: []byte("hello world"),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var logged bytes.Buffer
			ctx := middleware.SetLogger(context.Background(), logging.LoggerFunc(
				func(classification logging.Classification, format string, v ...interface{}) {
					fmt.Fprintf(&logged, format, v...)
				}))

			validateOutput := validateOutputPayloadChecksum{
				Algorithms: []Algorithm{
					AlgorithmSHA1, AlgorithmCRC32, AlgorithmCRC32C,
				},
				LogValidationSkipped:          true,
				LogMultipartValidationSkipped: true,
			}
			if c.validateOptions != nil {
				c.validateOptions(&validateOutput)
			}

			out, meta, err := validateOutput.HandleDeserialize(ctx,
				middleware.DeserializeInput{},
				middleware.DeserializeHandlerFunc(
					func(ctx context.Context, input middleware.DeserializeInput) (
						out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
					) {
						out.RawResponse = c.response
						return out, metadata, nil
					},
				),
			)
			if err == nil && len(c.expectErr) != 0 {
				t.Fatalf("expect error %v, got none", c.expectErr)
			}
			if err != nil && len(c.expectErr) == 0 {
				t.Fatalf("expect no error, got %v", err)
			}
			if err != nil && !strings.Contains(err.Error(), c.expectErr) {
				t.Fatalf("expect error to contain %v, got %v", c.expectErr, err)
			}
			if c.expectErr != "" {
				return
			}

			response := out.RawResponse.(*smithyhttp.Response)

			actualPayload, err := ioutil.ReadAll(response.Body)
			if err == nil && len(c.expectReadErr) != 0 {
				t.Fatalf("expected read error: %v, got none", c.expectReadErr)
			}
			if err != nil && len(c.expectReadErr) == 0 {
				t.Fatalf("expect no read error, got %v", err)
			}
			if err != nil && !strings.Contains(err.Error(), c.expectReadErr) {
				t.Fatalf("expected read error %v to contain %v", err, c.expectReadErr)
			}
			if c.expectReadErr != "" {
				return
			}

			if c.expectLogged != "" {
				if e, a := c.expectLogged, logged.String(); !strings.Contains(a, e) {
					t.Errorf("expected %q logged in:\n%s", e, a)
				}
			}

			if diff := cmp.Diff(string(c.expectPayload), string(actualPayload)); diff != "" {
				t.Errorf("expect payload match:\n%s", diff)
			}

			if err = response.Body.Close(); err != nil {
				t.Errorf("expect no close error, got %v", err)
			}

			values, ok := GetOutputValidationAlgorithmsUsed(meta)
			if ok != c.expectHaveAlgorithmsUsed {
				t.Errorf("expect metadata to contain algorithms used, %t", c.expectHaveAlgorithmsUsed)
			}
			if diff := cmp.Diff(c.expectAlgorithmsUsed, values); diff != "" {
				t.Errorf("expect algorithms used to match\n%s", diff)
			}
		})
	}
}
