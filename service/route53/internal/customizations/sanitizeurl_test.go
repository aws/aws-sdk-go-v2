package customizations_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/transport/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func TestSanitizeURLMiddleware(t *testing.T) {
	cases := map[string]struct {
		Given       string
		ExpectedURL string
	}{
		"includes hostedzone": {
			Given:       "hostedzone/ABCDEFG",
			ExpectedURL: "https://route53.amazonaws.com/2013-04-01/delegationset/ABCDEFG",
		},
		"excludes hostedzone": {
			Given:       "ABCDEFG",
			ExpectedURL: "https://route53.amazonaws.com/2013-04-01/delegationset/ABCDEFG",
		},
		"includes leading / in hostedzone": {
			Given:       "/hostedzone/ABCDEFG",
			ExpectedURL: "https://route53.amazonaws.com/2013-04-01/delegationset/ABCDEFG",
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

			fm := requestRetrieverMiddleware{}
			svc := route53.NewFromConfig(cfg)
			svc.DeleteReusableDelegationSet(context.Background(), &route53.DeleteReusableDelegationSetInput{
				Id: &c.Given,
			}, func(options *route53.Options) {
				options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
					stack.Serialize.Insert(&fm, "OperationSerializer", middleware.After)
					return nil
				})
			})

			if fm.request == nil {
				t.Fatalf("expected request to be serialized, got none")
			}

			if e, a := c.ExpectedURL, fm.request.URL.String(); !strings.EqualFold(e, a) {
				t.Fatalf("Expected url to be serialized as %v, got %v", e, a)
			}

		})
	}
}

type requestRetrieverMiddleware struct {
	request *http.Request
}

func (*requestRetrieverMiddleware) ID() string { return "Route53:requestRetrieverMiddleware" }

func (rm *requestRetrieverMiddleware) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*http.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}
	rm.request = req
	return next.HandleSerialize(ctx, in)
}
