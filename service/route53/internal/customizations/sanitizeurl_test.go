package customizations_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func TestSanitizeURLMiddleware(t *testing.T) {
	cases := map[string]struct {
		Given              string
		ExpectedPathSuffix string
	}{
		"includes hostedzone": {
			Given:              "hostedzone/ABCDEFG",
			ExpectedPathSuffix: "/delegationset/ABCDEFG",
		},
		"excludes hostedzone": {
			Given:              "ABCDEFG",
			ExpectedPathSuffix: "/delegationset/ABCDEFG",
		},
		"includes leading / in hostedzone": {
			Given:              "/hostedzone/ABCDEFG",
			ExpectedPathSuffix: "/delegationset/ABCDEFG",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			cfg := aws.Config{
				Credentials: unit.StubCredentialsProvider{},
				Retryer: func() aws.Retryer {
					return aws.NopRetryer{}
				},
				Region: "mock-region",
			}

			captured := &captureRequest{}
			svc := route53.NewFromConfig(cfg)
			svc.DeleteReusableDelegationSet(context.Background(), &route53.DeleteReusableDelegationSetInput{
				Id: &c.Given,
			}, func(options *route53.Options) {
				options.HTTPClient = captured
			})

			if captured.request == nil {
				t.Fatalf("expected request to be serialized, got none")
			}

			if !strings.HasSuffix(captured.request.URL.Path, c.ExpectedPathSuffix) {
				t.Fatalf("Expected URL path to end with %v, got %v", c.ExpectedPathSuffix, captured.request.URL.Path)
			}

		})
	}
}

type captureRequest struct {
	request *http.Request
}

func (c *captureRequest) Do(r *http.Request) (*http.Response, error) {
	c.request = r
	return &http.Response{
		Body: http.NoBody,
	}, nil
}
