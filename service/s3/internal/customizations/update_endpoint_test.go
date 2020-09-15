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

var (
	virtualHostStyleTests = []s3BucketTest{
		{"abc", "https://abc.s3.mock-region.amazonaws.com/", ""},
		{"a$b$c", "https://s3.mock-region.amazonaws.com/a%24b%24c", ""},
		{"a.b.c", "https://s3.mock-region.amazonaws.com/a.b.c", ""},
		{"a..bc", "https://s3.mock-region.amazonaws.com/a..bc", ""},
	}

	pathStyleTests = []s3BucketTest{
		{"abc", "https://s3.mock-region.amazonaws.com/abc", ""},
		{"a$b$c", "https://s3.mock-region.amazonaws.com/a%24b%24c", ""},
		{"a.b.c", "https://s3.mock-region.amazonaws.com/a.b.c", ""},
		{"a..bc", "https://s3.mock-region.amazonaws.com/a..bc", ""},
	}
)

var unitcreds = aws.StaticCredentialsProvider{
	Value: aws.Credentials{
		AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION",
		Source: "unit test credentials",
	},
}

func TestPathStyleBucketBuild(t *testing.T) {
	options := s3.Options{
		Credentials:  unitcreds,
		Retryer:      aws.NoOpRetryer{},
		Region:       "mock-region",
		UsePathStyle: true,
	}

	s := s3.New(options)
	runTests(t, s, pathStyleTests)
}

func TestHostStyleBucketBuild(t *testing.T) {
	options := s3.Options{
		Credentials: unitcreds,
		Retryer:     aws.NoOpRetryer{},
		Region:      "mock-region",
	}

	s := s3.New(options)
	runTests(t, s, virtualHostStyleTests)
}

func runTests(t *testing.T, svc *s3.Client, tests []s3BucketTest) {
	for i, test := range tests {
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
				t.Fatalf("%d, expect error code to contain %q, got %q", i, e, a)
			}
		}

		req := fm.request

		if e, a := test.url, req.URL.String(); e != a {
			t.Fatalf("%d, expect url %s, got %s", i, e, a)
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
