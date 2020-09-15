package customizations_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3BucketTest struct {
	bucket string
	url    string
	err    string
}

func TestUpdateEndpointBuild(t *testing.T) {
	cases := map[string]struct {
		usePathStyle bool
		tests        []s3BucketTest
	}{
		"PathStyleBucket": {
			usePathStyle: true,
			tests: []s3BucketTest{
				{"abc", "https://s3.mock-region.amazonaws.com/abc", ""},
				{"a$b$c", "https://s3.mock-region.amazonaws.com/a%24b%24c", ""},
				{"a.b.c", "https://s3.mock-region.amazonaws.com/a.b.c", ""},
				{"a..bc", "https://s3.mock-region.amazonaws.com/a..bc", ""},
			},
		},
		"VirtualHostStyleBucket": {
			tests: []s3BucketTest{
				{"abc", "https://abc.s3.mock-region.amazonaws.com/", ""},
				{"a$b$c", "https://s3.mock-region.amazonaws.com/a%24b%24c", ""},
				{"a.b.c", "https://s3.mock-region.amazonaws.com/a.b.c", ""},
				{"a..bc", "https://s3.mock-region.amazonaws.com/a..bc", ""},
			},
		},
	}

	for name, c := range cases {
		options := s3.Options{
			Credentials: aws.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION",
					Source: "unit test credentials",
				}},
			Retryer:      aws.NoOpRetryer{},
			Region:       "mock-region",
			UsePathStyle: c.usePathStyle,
		}

		svc := s3.New(options)
		for i, test := range c.tests {
			fm := requestRetrieverMiddleware{}
			_, err := svc.ListObjects(context.Background(),
				&s3.ListObjectsInput{Bucket: &test.bucket},
				func(options *s3.Options) {
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
