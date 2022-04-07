package v4a

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
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

type httpPresignerFunc func(
	ctx context.Context, credentials Credentials, r *http.Request,
	payloadHash string, service string, regionSet []string, signingTime time.Time,
	optFns ...func(*SignerOptions),
) (url string, signedHeader http.Header, err error)

func (f httpPresignerFunc) PresignHTTP(
	ctx context.Context, credentials Credentials, r *http.Request,
	payloadHash string, service string, regionSet []string, signingTime time.Time,
	optFns ...func(*SignerOptions),
) (
	url string, signedHeader http.Header, err error,
) {
	return f(ctx, credentials, r, payloadHash, service, regionSet, signingTime, optFns...)
}

func TestPresignHTTPRequestMiddleware(t *testing.T) {
	cases := map[string]struct {
		Request      *http.Request
		Creds        CredentialsProvider
		PayloadHash  string
		LogSigning   bool
		ExpectResult *v4.PresignedHTTPRequest
		ExpectErr    string
	}{
		"success": {
			Request: &http.Request{
				URL: func() *url.URL {
					u, _ := url.Parse("https://example.aws/path?query=foo")
					return u
				}(),
				Header: http.Header{},
			},
			Creds:       stubCredentials,
			PayloadHash: "0123456789abcdef",
			ExpectResult: &v4.PresignedHTTPRequest{
				URL:          "https://example.aws/path?query=foo",
				SignedHeader: http.Header{},
			},
		},
		"error": {
			Request: func() *http.Request {
				return &http.Request{}
			}(),
			Creds:       stubCredentials,
			PayloadHash: "",
			ExpectErr:   "failed to sign request",
		},
		"anonymous creds": {
			Request: &http.Request{
				URL: func() *url.URL {
					u, _ := url.Parse("https://example.aws/path?query=foo")
					return u
				}(),
				Header: http.Header{},
			},
			Creds:       stubCredentials,
			PayloadHash: "",
			ExpectErr:   "failed to sign request",
			ExpectResult: &v4.PresignedHTTPRequest{
				URL:          "https://example.aws/path?query=foo",
				SignedHeader: http.Header{},
			},
		},
		"nil creds": {
			Request: &http.Request{
				URL: func() *url.URL {
					u, _ := url.Parse("https://example.aws/path?query=foo")
					return u
				}(),
				Header: http.Header{},
			},
			Creds: nil,
			ExpectResult: &v4.PresignedHTTPRequest{
				URL:          "https://example.aws/path?query=foo",
				SignedHeader: http.Header{},
			},
		},
		"with log signing": {
			Request: &http.Request{
				URL: func() *url.URL {
					u, _ := url.Parse("https://example.aws/path?query=foo")
					return u
				}(),
				Header: http.Header{},
			},
			Creds:       stubCredentials,
			PayloadHash: "0123456789abcdef",
			ExpectResult: &v4.PresignedHTTPRequest{
				URL:          "https://example.aws/path?query=foo",
				SignedHeader: http.Header{},
			},

			LogSigning: true,
		},
	}

	const (
		signingName   = "serviceId"
		signingRegion = "regionName"
	)

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			m := &PresignHTTPRequestMiddleware{
				credentialsProvider: tt.Creds,

				presigner: httpPresignerFunc(func(
					ctx context.Context, credentials Credentials, r *http.Request,
					payloadHash string, service string, regionSet []string, signingTime time.Time,
					optFns ...func(*SignerOptions),
				) (url string, signedHeader http.Header, err error) {
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

					if !hasCredentialProvider(tt.Creds) {
						t.Errorf("expect presigner not to be called for not credentials provider")
					}

					expectCreds, _ := tt.Creds.RetrievePrivateKey(context.Background())
					if diff := cmp.Diff(expectCreds, credentials); len(diff) > 0 {
						t.Error(diff)
					}
					if e, a := tt.PayloadHash, payloadHash; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					if e, a := signingName, service; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					if diff := cmp.Diff([]string{signingRegion}, regionSet); len(diff) > 0 {
						t.Error(diff)
					}

					return tt.ExpectResult.URL, tt.ExpectResult.SignedHeader, nil
				}),
				logSigning: tt.LogSigning,
			}

			next := middleware.FinalizeHandlerFunc(
				func(ctx context.Context, in middleware.FinalizeInput) (
					out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
				) {
					t.Errorf("expect next handler not to be called")
					return out, metadata, err
				})

			ctx := awsmiddleware.SetSigningRegion(
				awsmiddleware.SetSigningName(context.Background(), signingName),
				signingRegion)

			var loggerBuf bytes.Buffer
			logger := logging.NewStandardLogger(&loggerBuf)
			ctx = middleware.SetLogger(ctx, logger)

			if len(tt.PayloadHash) != 0 {
				ctx = v4.SetPayloadHash(ctx, tt.PayloadHash)
			}

			result, _, err := m.HandleFinalize(ctx, middleware.FinalizeInput{
				Request: &smithyhttp.Request{
					Request: tt.Request,
				},
			}, next)
			if len(tt.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := tt.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect error to contain %v, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if diff := cmp.Diff(tt.ExpectResult, result.Result); len(diff) != 0 {
				t.Errorf("expect result match\n%v", diff)
			}

			if tt.LogSigning {
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
	_ middleware.FinalizeMiddleware = &PresignHTTPRequestMiddleware{}
)
