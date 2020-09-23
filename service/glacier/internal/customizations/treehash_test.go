package customizations

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
	"io"
	"testing"
)

func TestTreeHashMiddleware(t *testing.T) {
	m := TreeHashMiddleware{}

	cases := map[string]struct {
		Body   io.ReadSeeker
		Linear string
		Tree   string
	}{
		"all zeroes": {
			Body:   createTestReader(),
			Linear: "68aff0c5a91aa0491752bfb96e3fef33eb74953804f6a2f7b708d5bcefa8ff6b",
			Tree:   "154e26c78fd74d0c2c9b3cc4644191619dc4f2cd539ae2a74d5fd07957a3ee6a",
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

						n, _ := c.Body.Seek(0, io.SeekCurrent)
						if n != 0 {
							t.Fatalf("expected body to be rewound, but was at position %d", n)
						}

						return output, metadata, err
					}),
			)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

func TestComputeHashes(t *testing.T) {
	cases := map[string]struct {
		Body   io.ReadSeeker
		Linear string
		Tree   string
	}{
		"all zeroes": {
			Body:   createTestReader(),
			Linear: "68aff0c5a91aa0491752bfb96e3fef33eb74953804f6a2f7b708d5bcefa8ff6b",
			Tree:   "154e26c78fd74d0c2c9b3cc4644191619dc4f2cd539ae2a74d5fd07957a3ee6a",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := ComputeHashes(c.Body)
			actualLinear := fmt.Sprintf("%x", actual.LinearHash)
			if actualLinear != c.Linear {
				t.Fatalf("expected linear hash to be \"%s\" but was \"%s\"", c.Linear, actualLinear)
			}

			actualTree := fmt.Sprintf("%x", actual.TreeHash)
			if actualTree != c.Tree {
				t.Fatalf("expected tree hash to be \"%s\" but was \"%s\"", c.Tree, actualTree)
			}

			n, _ := c.Body.Seek(0, io.SeekCurrent)
			if n != 0 {
				t.Fatalf("expected body to be rewound, but was at position %d", n)
			}
		})
	}
}

func TestComputeTreeHash(t *testing.T) {
	r := createTestReader()

	const chunkSize = 1024 * 1024 // 1MB
	buf := make([]byte, chunkSize)
	hashes := [][]byte{}

	for {
		// Reach 1MB chunks from reader to generate hashes from
		n, err := io.ReadAtLeast(r, buf, chunkSize)
		if n == 0 {
			break
		}

		tmpHash := sha256.Sum256(buf[:n])
		hashes = append(hashes, tmpHash[:])
		if err != nil {
			break // last chunk
		}
	}

	actual := fmt.Sprintf("%x", ComputeTreeHash(hashes))
	expected := "154e26c78fd74d0c2c9b3cc4644191619dc4f2cd539ae2a74d5fd07957a3ee6a"

	if actual != expected {
		t.Fatalf("expected tree hash to be \"%s\" but was \"%s\"", expected, actual)
	}
}

func createTestReader() io.ReadSeeker {
	buf := make([]byte, 5767168) // 5.5MB buffer
	for i := range buf {
		buf[i] = '0' // Fill with zero characters
	}

	return bytes.NewReader(buf)
}
