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

	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
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
			c := &computePayloadSHA256{}

			next := middleware.BuildHandlerFunc(func(ctx context.Context, in middleware.BuildInput) (out middleware.BuildOutput, metadata middleware.Metadata, err error) {
				value, ok := ctx.Value(payloadHashKey{}).(string)
				if !ok {
					t.Fatalf("expected payload hash value to be on context")
				}
				if e, a := tt.expectedHash, value; e != a {
					t.Errorf("expected %v, got %v", e, a)
				}

				return out, metadata, err
			})

			stream, err := smithyhttp.NewStackRequest().(*smithyhttp.Request).SetStream(tt.content)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			_, _, err = c.HandleBuild(context.Background(), middleware.BuildInput{Request: stream}, next)
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

type httpSignerFunc func(ctx context.Context, credentials aws.Credentials, r *http.Request, payloadHash string, service string, region string, signingTime time.Time) error

func (f httpSignerFunc) SignHTTP(ctx context.Context, credentials aws.Credentials, r *http.Request, payloadHash string, service string, region string, signingTime time.Time) error {
	return f(ctx, credentials, r, payloadHash, service, region, signingTime)
}

func TestSignHTTPRequestMiddleware(t *testing.T) {
	cases := map[string]struct {
		creds       aws.CredentialsProvider
		hash        string
		expectedErr error
	}{
		"success": {
			creds: unit.StubCredentialsProvider{},
			hash:  "0123456789abcdef",
		},
		"error": {
			creds:       unit.StubCredentialsProvider{},
			hash:        "",
			expectedErr: &SigningError{},
		},
		"anonymous creds": {
			creds: aws.AnonymousCredentials{},
		},
		"nil creds": {
			creds: nil,
		},
	}

	const (
		signingName   = "serviceId"
		signingRegion = "regionName"
	)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			c := &SignHTTPRequest{
				credentialsProvider: tt.creds,
				signer: httpSignerFunc(
					func(ctx context.Context,
						credentials aws.Credentials, r *http.Request, payloadHash string,
						service string, region string, signingTime time.Time,
					) error {
						expectCreds, _ := unit.StubCredentialsProvider{}.Retrieve(context.Background())
						if e, a := expectCreds, credentials; e != a {
							t.Errorf("expected %v, got %v", e, a)
						}
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

			next := middleware.FinalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
				return out, metadata, err
			})

			ctx := awsmiddleware.SetSigningRegion(
				awsmiddleware.SetSigningName(context.Background(), signingName),
				signingRegion)

			if len(tt.hash) != 0 {
				ctx = context.WithValue(ctx, payloadHashKey{}, tt.hash)
			}

			_, _, err := c.HandleFinalize(ctx, middleware.FinalizeInput{
				Request: &smithyhttp.Request{Request: &http.Request{}},
			}, next)
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

var (
	_ middleware.BuildMiddleware    = &unsignedPayload{}
	_ middleware.BuildMiddleware    = &computePayloadSHA256{}
	_ middleware.BuildMiddleware    = &contentSHA256Header{}
	_ middleware.FinalizeMiddleware = &SignHTTPRequest{}
)
