package v4

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/awslabs/smithy-go/middleware"
	smithyHTTP "github.com/awslabs/smithy-go/transport/http"
)

func TestComputePayloadHashMiddleware(t *testing.T) {
	cases := []struct {
		content      io.Reader
		expectedHash string
		expectedErr  error
	}{
		0: {
			content: func() io.Reader {
				br := bytes.NewReader([]byte("some content"))
				return br
			}(),
			expectedHash: "290f493c44f5d63d06b374d0a5abd292fae38b92cab2fae5efefe1b0e9347f56",
		},
		1: {
			content: func() io.Reader {
				return &nonSeeker{}
			}(),
			expectedErr: &HashComputationError{},
		},
		2: {
			content: func() io.Reader {
				return &semiSeekable{}
			}(),
			expectedErr: &HashComputationError{},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c := &ComputePayloadSHA256Middleware{}

			next := finalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
				value, ok := ctx.Value(payloadHashKey{}).(string)
				if !ok {
					t.Fatalf("expected payload hash value to be on context")
				}
				if e, a := tt.expectedHash, value; e != a {
					t.Errorf("expected %v, got %v", e, a)
				}

				return out, metadata, err
			})

			stream, err := smithyHTTP.NewStackRequest().(*smithyHTTP.Request).SetStream(tt.content)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			_, _, err = c.HandleFinalize(context.Background(), middleware.FinalizeInput{Request: stream}, next)
			if err != nil && tt.expectedErr == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err != nil && tt.expectedErr != nil {
				e, a := tt.expectedErr, err
				if !errors.As(a, &e) {
					t.Errorf("expected error type %T, got %T", e, a)
				}
			} else if err == nil && tt.expectedErr != nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

type mockSigner func(ctx context.Context, r *http.Request, payloadHash string, service string, region string, signingTime time.Time) error

func (f mockSigner) SignHTTP(ctx context.Context, r *http.Request, payloadHash string, service string, region string, signingTime time.Time) error {
	return f(ctx, r, payloadHash, service, region, signingTime)
}

func TestSignHTTPRequestMiddleware(t *testing.T) {
	cases := []struct {
		hash        string
		expectedErr error
	}{
		0: {
			hash: "0123456789abcdef",
		},
		1: {
			hash:        "",
			expectedErr: &SigningError{},
		},
	}

	const (
		signingName   = "serviceId"
		signingRegion = "regionName"
	)

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c := &SignHTTPRequestMiddleware{
				signer: mockSigner(func(ctx context.Context, r *http.Request, payloadHash string, service string, region string, signingTime time.Time) error {
					if e, a := tt.hash, payloadHash; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					if e, a := signingName, service; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					if e, a := signingRegion, region; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return nil
				}),
			}

			next := finalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
				return out, metadata, err
			})

			ctx := SetSigningMetadata(context.Background(), SigningMetadata{
				SigningName:   signingName,
				SigningRegion: signingRegion,
			})

			if len(tt.hash) != 0 {
				ctx = context.WithValue(ctx, payloadHashKey{}, tt.hash)
			}

			_, _, err := c.HandleFinalize(ctx, middleware.FinalizeInput{Request: &smithyHTTP.Request{Request: &http.Request{}}}, next)
			if err != nil && tt.expectedErr == nil {
				t.Errorf("expected no error, got %v", err)
			} else if err != nil && tt.expectedErr != nil {
				e, a := tt.expectedErr, err
				if !errors.As(a, &e) {
					t.Errorf("expected error type %T, got %T", e, a)
				}
			} else if err == nil && tt.expectedErr != nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

type nonSeeker struct{}

func (nonSeeker) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

type semiSeekable struct {
	hasSeeked bool
}

func (s *semiSeekable) Seek(offset int64, whence int) (int64, error) {
	if !s.hasSeeked {
		s.hasSeeked = true
		return 0, nil
	}
	return 0, fmt.Errorf("io seek error")
}

func (*semiSeekable) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

type finalizeHandlerFunc func(ctx context.Context, in middleware.FinalizeInput) (middleware.FinalizeOutput, middleware.Metadata, error)

func (f finalizeHandlerFunc) HandleFinalize(ctx context.Context, in middleware.FinalizeInput) (middleware.FinalizeOutput, middleware.Metadata, error) {
	return f(ctx, in)
}

var (
	_ middleware.FinalizeMiddleware = &ComputePayloadSHA256Middleware{}
	_ middleware.FinalizeMiddleware = &SignHTTPRequestMiddleware{}
)
