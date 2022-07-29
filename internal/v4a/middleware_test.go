package v4a

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/google/go-cmp/cmp"
)

type stubCredentialsProviderFunc func(context.Context) (Credentials, error)

func (f stubCredentialsProviderFunc) RetrievePrivateKey(ctx context.Context) (Credentials, error) {
	return f(ctx)
}

type httpSignerFunc func(ctx context.Context, credentials Credentials, r *http.Request, payloadHash string, service string, regionSet []string, signingTime time.Time, optFns ...func(*SignerOptions)) error

func (f httpSignerFunc) SignHTTP(ctx context.Context, credentials Credentials, r *http.Request, payloadHash string, service string, regionSet []string, signingTime time.Time, optFns ...func(*SignerOptions)) error {
	return f(ctx, credentials, r, payloadHash, service, regionSet, signingTime, optFns...)
}

func TestSignHTTPRequestMiddleware(t *testing.T) {
	cases := map[string]struct {
		creds       CredentialsProvider
		hash        string
		logSigning  bool
		expectedErr interface{}
	}{
		"success": {
			creds: stubCredentials,
			hash:  "0123456789abcdef",
		},
		"error": {
			creds: stubCredentialsProviderFunc(func(ctx context.Context) (Credentials, error) {
				return Credentials{}, fmt.Errorf("credential error")
			}),
			hash:        "",
			expectedErr: &SigningError{},
		},
		"nil creds": {
			creds: nil,
		},
		"with log signing": {
			creds:      stubCredentials,
			hash:       "0123456789abcdef",
			logSigning: true,
		},
	}

	const (
		signingName   = "serviceId"
		signingRegion = "regionName"
	)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			c := &SignHTTPRequestMiddleware{
				credentials: tt.creds,
				signer: httpSignerFunc(
					func(ctx context.Context,
						credentials Credentials, r *http.Request, payloadHash string,
						service string, regionSet []string, signingTime time.Time,
						optFns ...func(*SignerOptions),
					) error {
						var options SignerOptions
						for _, fn := range optFns {
							fn(&options)
						}
						if options.Logger == nil {
							t.Errorf("expect logger, got none")
						}
						if options.LogSigning {
							options.Logger.Logf(logging.Debug, t.Name())
						}

						expectCreds, _ := tt.creds.RetrievePrivateKey(ctx)
						if diff := cmp.Diff(expectCreds, credentials); len(diff) > 0 {
							t.Error(diff)
						}
						if e, a := tt.hash, payloadHash; e != a {
							t.Errorf("expected %v, got %v", e, a)
						}
						if e, a := signingName, service; e != a {
							t.Errorf("expected %v, got %v", e, a)
						}
						if diff := cmp.Diff([]string{signingRegion}, regionSet); len(diff) > 0 {
							t.Error(diff)
						}
						return nil
					}),
				logSigning: tt.logSigning,
			}

			next := middleware.FinalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
				return out, metadata, err
			})

			ctx := awsmiddleware.SetSigningRegion(
				awsmiddleware.SetSigningName(context.Background(), signingName),
				signingRegion)

			var loggerBuf bytes.Buffer
			logger := logging.NewStandardLogger(&loggerBuf)
			ctx = middleware.SetLogger(ctx, logger)

			if len(tt.hash) != 0 {
				ctx = v4.SetPayloadHash(ctx, tt.hash)
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

			if tt.logSigning {
				if e, a := t.Name(), loggerBuf.String(); !strings.Contains(a, e) {
					t.Errorf("expect %v logged in %v", e, a)
				}
			} else {
				if loggerBuf.Len() != 0 {
					t.Errorf("expect no log, got %v", loggerBuf.String())
				}
			}
		})
	}
}

var (
	_ middleware.FinalizeMiddleware = &SignHTTPRequestMiddleware{}
)
