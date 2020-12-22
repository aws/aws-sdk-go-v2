package customizations

import (
	"bytes"
	"context"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"io"
	"strings"
	"testing"
)

func TestTreeHashMiddleware(t *testing.T) {
	m := TreeHash{}

	cases := map[string]struct {
		Body       io.Reader
		SetRequest func(r *smithyhttp.Request) *smithyhttp.Request
		Linear     string
		Tree       string
		ExpectErr  string
	}{
		"all zeroes": {
			Body: func() io.ReadSeeker {
				buf := make([]byte, 5767168) // 5.5MB buffer
				for i := range buf {
					buf[i] = '0' // Fill with zero characters
				}
				return bytes.NewReader(buf)
			}(),
			Linear: "68aff0c5a91aa0491752bfb96e3fef33eb74953804f6a2f7b708d5bcefa8ff6b",
			Tree:   "154e26c78fd74d0c2c9b3cc4644191619dc4f2cd539ae2a74d5fd07957a3ee6a",
		},
		"non-seekable stream without precomputed hash": {
			Body: func() io.Reader {
				buf := make([]byte, 5767168)
				for i := range buf {
					buf[i] = '0'
				}
				return bytes.NewBuffer(buf)
			}(),
			ExpectErr: "tree hash",
		},
		"non-seekable stream with precomputed hash": {
			Body: func() io.Reader {
				buf := make([]byte, 5767168)
				for i := range buf {
					buf[i] = '0'
				}
				return bytes.NewBuffer(buf)
			}(),
			SetRequest: func(r *smithyhttp.Request) *smithyhttp.Request {
				r.Header.Set("X-Amz-Content-Sha256", "precomputed")
				r.Header.Set("X-Amz-Sha256-Tree-Hash", "precomputed")
				return r
			},
			// These values aren't actually correct, but they are what we're explicitly setting.
			// If the middleware were to try to compute the values anyway, the checks would fail.
			Linear: "precomputed",
			Tree:   "precomputed",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, _, err := m.HandleFinalize(context.Background(),
				middleware.FinalizeInput{
					Request: func() *smithyhttp.Request {
						r := smithyhttp.NewStackRequest()
						req, ok := r.(*smithyhttp.Request)
						if !ok || req == nil {
							t.Fatalf("expect smithy request, got %T", r)
						}
						if c.SetRequest != nil {
							req = c.SetRequest(req)
						}
						req, err := req.SetStream(c.Body)
						if err != nil {
							t.Fatalf("expect no error, got %v", err)
						}
						return req
					}(),
				},
				middleware.FinalizeHandlerFunc(
					func(ctx context.Context, input middleware.FinalizeInput) (
						output middleware.FinalizeOutput, metadata middleware.Metadata, err error,
					) {
						req, ok := input.Request.(*smithyhttp.Request)
						if !ok || req == nil {
							t.Fatalf("expect smithy request, got %T", input.Request)
						}

						actualLinear := req.Header.Get("X-Amz-Content-Sha256")
						if actualLinear != c.Linear {
							t.Fatalf("expected linear hash to be \"%s\" but was \"%s\"", c.Linear, actualLinear)
						}

						actualTree := req.Header.Get("X-Amz-Sha256-Tree-Hash")
						if actualTree != c.Tree {
							t.Fatalf("expected tree hash to be \"%s\" but was \"%s\"", c.Tree, actualTree)
						}

						if body, ok := c.Body.(io.ReadSeeker); ok {
							n, _ := body.Seek(0, io.SeekCurrent)
							if n != 0 {
								t.Fatalf("expected body to be rewound, but was at position %d", n)
							}

						}
						return output, metadata, err
					}),
			)
			if len(c.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect error to contain %v, got %v", e, a)
				}
			} else {
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
			}
		})
	}
}
