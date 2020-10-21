package customizations_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3/internal/endpoints"

	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3BucketTest struct {
	bucket string
	url    string
	err    string
}

func TestUpdateEndpointBuild(t *testing.T) {
	cases := map[string]map[string]struct {
		tests          []s3BucketTest
		useAccelerate  bool
		useDualstack   bool
		usePathStyle   bool
		disableHTTPS   bool
		customEndpoint *aws.Endpoint
	}{
		"default endpoint": {
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
			"Accelerate": {
				useAccelerate: true,
				tests: []s3BucketTest{
					{"abc", "https://abc.s3-accelerate.amazonaws.com/", ""},
					{"a.b.c", "https://s3.mock-region.amazonaws.com/a.b.c", "not compatible"},
					{"a$b$c", "https://s3.mock-region.amazonaws.com/a%24b%24c", "not compatible"},
				},
			},
			"AccelerateNoSSLTests": {
				useAccelerate: true,
				disableHTTPS:  true,
				tests: []s3BucketTest{
					{"abc", "http://abc.s3-accelerate.amazonaws.com/", ""},
					{"a.b.c", "http://a.b.c.s3-accelerate.amazonaws.com/", ""},
					{"a$b$c", "http://s3.mock-region.amazonaws.com/a%24b%24c", "not compatible"},
				},
			},
			"DualStack": {
				useDualstack: true,
				tests: []s3BucketTest{
					{"abc", "https://abc.s3.dualstack.mock-region.amazonaws.com/", ""},
					{"a.b.c", "https://s3.dualstack.mock-region.amazonaws.com/a.b.c", ""},
					{"a$b$c", "https://s3.dualstack.mock-region.amazonaws.com/a%24b%24c", ""},
				},
			},
			"DualStackWithPathStyle": {
				useDualstack: true,
				usePathStyle: true,
				tests: []s3BucketTest{
					{"abc", "https://s3.dualstack.mock-region.amazonaws.com/abc", ""},
					{"a.b.c", "https://s3.dualstack.mock-region.amazonaws.com/a.b.c", ""},
					{"a$b$c", "https://s3.dualstack.mock-region.amazonaws.com/a%24b%24c", ""},
				},
			},
			"AccelerateWithDualStack": {
				useAccelerate: true,
				useDualstack:  true,
				tests: []s3BucketTest{
					{"abc", "https://abc.s3-accelerate.dualstack.amazonaws.com/", ""},
					{"a.b.c", "https://s3.dualstack.mock-region.amazonaws.com/a.b.c", "not compatible"},
					{"a$b$c", "https://s3.dualstack.mock-region.amazonaws.com/a%24b%24c", "not compatible"},
				},
			},
		},

		"immutable endpoint": {
			"PathStyleBucket": {
				usePathStyle: true,
				customEndpoint: &aws.Endpoint{
					URL:               "https://example.region.amazonaws.com",
					HostnameImmutable: true,
				},
				tests: []s3BucketTest{
					{"abc", "https://example.region.amazonaws.com/abc", ""},
					{"a$b$c", "https://example.region.amazonaws.com/a%24b%24c", ""},
					{"a.b.c", "https://example.region.amazonaws.com/a.b.c", ""},
					{"a..bc", "https://example.region.amazonaws.com/a..bc", ""},
				},
			},
			"VirtualHostStyleBucket": {
				customEndpoint: &aws.Endpoint{
					URL:               "https://example.region.amazonaws.com",
					HostnameImmutable: true,
				},
				tests: []s3BucketTest{
					{"abc", "https://example.region.amazonaws.com/abc", ""},
					{"a$b$c", "https://example.region.amazonaws.com/a%24b%24c", ""},
					{"a.b.c", "https://example.region.amazonaws.com/a.b.c", ""},
					{"a..bc", "https://example.region.amazonaws.com/a..bc", ""},
				},
			},
			"Accelerate": {
				useAccelerate: true,
				customEndpoint: &aws.Endpoint{
					URL:               "https://example.region.amazonaws.com",
					HostnameImmutable: true,
				},
				tests: []s3BucketTest{
					{"abc", "https://example.region.amazonaws.com/abc", ""},
					{"a$b$c", "https://example.region.amazonaws.com/a%24b%24c", ""},
					{"a.b.c", "https://example.region.amazonaws.com/a.b.c", ""},
					{"a..bc", "https://example.region.amazonaws.com/a..bc", ""},
				},
			},
			"AccelerateNoSSLTests": {
				useAccelerate: true,
				disableHTTPS:  true,
				customEndpoint: &aws.Endpoint{
					URL:               "https://example.region.amazonaws.com",
					HostnameImmutable: true,
				},
				tests: []s3BucketTest{
					{"abc", "https://example.region.amazonaws.com/abc", ""},
					{"a.b.c", "https://example.region.amazonaws.com/a.b.c", ""},
					{"a$b$c", "https://example.region.amazonaws.com/a%24b%24c", ""},
				},
			},
			"DualStack": {
				useDualstack: true,
				customEndpoint: &aws.Endpoint{
					URL:               "https://example.region.amazonaws.com",
					HostnameImmutable: true,
				},
				tests: []s3BucketTest{
					{"abc", "https://example.region.amazonaws.com/abc", ""},
					{"a.b.c", "https://example.region.amazonaws.com/a.b.c", ""},
					{"a$b$c", "https://example.region.amazonaws.com/a%24b%24c", ""},
				},
			},
		},
	}

	for suitName, cs := range cases {
		t.Run(suitName, func(t *testing.T) {
			for unitName, c := range cs {
				t.Run(unitName, func(t *testing.T) {
					options := s3.Options{
						Credentials: unit.StubCredentialsProvider{},
						Retryer:     aws.NoOpRetryer{},
						Region:      "mock-region",

						HTTPClient: smithyhttp.NopClient{},

						EndpointOptions: endpoints.Options{
							DisableHTTPS: c.disableHTTPS,
						},

						UsePathStyle:  c.usePathStyle,
						UseAccelerate: c.useAccelerate,
						UseDualstack:  c.useDualstack,
					}

					if c.customEndpoint != nil {
						options.EndpointResolver = s3.EndpointResolverFunc(
							func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
								return *c.customEndpoint, nil
							})
					}

					svc := s3.New(options)
					for i, test := range c.tests {
						t.Run(strconv.Itoa(i), func(t *testing.T) {
							fm := requestRetrieverMiddleware{}
							_, err := svc.ListObjects(context.Background(),
								&s3.ListObjectsInput{Bucket: &test.bucket},
								func(options *s3.Options) {
									options.APIOptions = append(options.APIOptions,
										func(stack *middleware.Stack) error {
											stack.Serialize.Insert(&fm,
												"OperationSerializer", middleware.Before)
											return nil
										})
								},
							)

							if test.err != "" {
								if err == nil {
									t.Fatalf("test %d: expected error, got none", i)
								}
								if a, e := err.Error(), test.err; !strings.Contains(a, e) {
									t.Fatalf("expect error code to contain %q, got %q", e, a)
								}
								return
							}
							if err != nil {
								t.Fatalf("expect no error, got %v", err)
							}

							req := fm.request.Build(context.Background())
							if e, a := test.url, req.URL.String(); e != a {
								t.Fatalf("expect url %s, got %s", e, a)
							}
						})
					}
				})
			}
		})
	}
}

type requestRetrieverMiddleware struct {
	request *smithyhttp.Request
}

func (*requestRetrieverMiddleware) ID() string { return "S3:requestRetrieverMiddleware" }

func (rm *requestRetrieverMiddleware) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}
	rm.request = req
	return next.HandleSerialize(ctx, in)
}
