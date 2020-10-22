package customizations_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3control"

	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"
)

type s3controlEndpointTest struct {
	bucket    string
	accountID string
	url       string
	err       string
}

func TestUpdateEndpointBuild(t *testing.T) {
	cases := map[string]struct {
		tests        []s3controlEndpointTest
		useDualstack bool
	}{
		"DualStack": {
			useDualstack: true,
			tests: []s3controlEndpointTest{
				{"abc", "123456789012", "https://s3-control.dualstack.mock-region.amazonaws.com/v20180820/bucket/abc", ""},
				{"a.b.c", "123456789012", "https://s3-control.dualstack.mock-region.amazonaws.com/v20180820/bucket/a.b.c", ""},
				{"a$b$c", "123456789012", "https://s3-control.dualstack.mock-region.amazonaws.com/v20180820/bucket/a%24b%24c", ""},
			},
		},
	}

	for name, c := range cases {
		options := s3control.Options{
			Credentials:  unit.StubCredentialsProvider{},
			Retryer:      aws.NoOpRetryer{},
			Region:       "mock-region",
			UseDualstack: c.useDualstack,
		}

		svc := s3control.New(options)
		for i, test := range c.tests {
			fm := requestRetrieverMiddleware{}
			_, err := svc.DeleteBucket(context.Background(),
				&s3control.DeleteBucketInput{
					Bucket:    &test.bucket,
					AccountId: &test.accountID,
				},
				func(options *s3control.Options) {
					options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
						stack.Serialize.Insert(&fm, "OperationSerializer", middleware.Before)
						return nil
					})
				},
			)

			if test.err != "" {
				if err == nil {
					t.Fatalf("test %d: expected error, got none", i)
				}
				if a, e := err.Error(), test.err; !strings.Contains(a, e) {
					t.Fatalf("%s: %d, expect error code to contain %q, got %q", name, i, e, a)
				}
			}

			req := fm.request

			if e, a := test.url, req.URL.String(); e != a {
				t.Fatalf("%s: %d, expect url %s, got %s", name, i, e, a)
			}
		}
	}
}

type requestRetrieverMiddleware struct {
	request *http.Request
}

func (*requestRetrieverMiddleware) ID() string { return "S3:requestRetrieverMiddleware" }

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
