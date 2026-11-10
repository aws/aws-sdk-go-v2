package kitchensinktest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/smithy-go"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	smithytraits "github.com/aws/smithy-go/traits"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type retryableTraitEndpointResolver struct{}

func (*retryableTraitEndpointResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (smithyendpoints.Endpoint, error) {
	return smithyendpoints.Endpoint{URI: url.URL{Scheme: "https", Host: "test.example.com"}}, nil
}

// newRetryableTraitClient returns a client whose every attempt fails with the
// given modeled error, counting attempts. The retryer is the standard retryer
// with its default checks, and no backoff delay.
func newRetryableTraitClient(errorType string, attempts *int) *Client {
	return New(Options{
		Region: "us-east-1",
		Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}, nil
		}),
		EndpointResolverV2: &retryableTraitEndpointResolver{},
		Retryer: retry.NewStandard(func(o *retry.StandardOptions) {
			o.Backoff = retry.BackoffDelayerFunc(func(int, error) (time.Duration, error) {
				return 0, nil
			})
		}),
		HTTPClient: smithyhttp.ClientDoFunc(func(req *http.Request) (*http.Response, error) {
			*attempts++
			return &http.Response{
				StatusCode: 400,
				Header:     http.Header{"Content-Type": []string{"application/x-amz-json-1.0"}},
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"__type":"` + errorType + `","message":"failed"}`))),
				Request:    req,
			}, nil
		}),
	})
}

// TestRetryableTrait verifies that modeled errors report the smithy.api#retryable
// trait, and that the default retryer honors it.
func TestRetryableTrait(t *testing.T) {
	cases := map[string]struct {
		ErrorType      string
		ExpectModeled  bool
		ExpectTrait    smithytraits.Retryable
		ExpectAttempts int
		ExpectThrottle bool
	}{
		"without trait": {
			ErrorType:      "ItemNotFound",
			ExpectAttempts: 1,
		},
		"retryable": {
			ErrorType:      "RetryableTraitError",
			ExpectModeled:  true,
			ExpectAttempts: 3,
		},
		"retryable throttling": {
			ErrorType:      "ThrottlingRetryableTraitError",
			ExpectModeled:  true,
			ExpectTrait:    smithytraits.Retryable{Throttling: true},
			ExpectAttempts: 3,
			ExpectThrottle: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var attempts int
			svc := newRetryableTraitClient(c.ErrorType, &attempts)

			_, err := svc.GetItem(context.Background(), &GetItemInput{})

			var apiErr smithy.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expect API error, got %T: %v", err, err)
			}
			if e, a := c.ErrorType, apiErr.ErrorCode(); e != a {
				t.Fatalf("expect error code %q, got %q", e, a)
			}

			var modeled interface {
				RetryableTrait() (smithytraits.Retryable, bool)
			}
			if !errors.As(err, &modeled) {
				t.Fatalf("expect modeled error with RetryableTrait, got %T", err)
			}
			trait, ok := modeled.RetryableTrait()
			if e, a := c.ExpectModeled, ok; e != a {
				t.Errorf("expect modeled with trait %v, got %v", e, a)
			}
			if e, a := c.ExpectTrait, trait; e != a {
				t.Errorf("expect trait %+v, got %+v", e, a)
			}

			if e, a := c.ExpectAttempts, attempts; e != a {
				t.Errorf("expect %v attempts, got %v", e, a)
			}
			throttle := retry.IsErrorThrottles(retry.DefaultThrottles).IsErrorThrottle(err).Bool()
			if e, a := c.ExpectThrottle, throttle; e != a {
				t.Errorf("expect throttle %v, got %v", e, a)
			}
		})
	}
}
