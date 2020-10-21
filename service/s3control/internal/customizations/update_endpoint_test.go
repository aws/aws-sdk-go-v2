package customizations_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3control"

	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

type s3controlEndpointTest struct {
	bucket    string
	accountID string
	url       string
	err       string
}

func TestUpdateEndpointBuild(t *testing.T) {
	cases := map[string]map[string]struct {
		tests          []s3controlEndpointTest
		useDualstack   bool
		customEndpoint *aws.Endpoint
	}{
		"default endpoint": {
			"default": {
				tests: []s3controlEndpointTest{
					{"abc", "123456789012", "https://123456789012.s3-control.mock-region.amazonaws.com/v20180820/bucket/abc", ""},
					{"a.b.c", "123456789012", "https://123456789012.s3-control.mock-region.amazonaws.com/v20180820/bucket/a.b.c", ""},
					{"a$b$c", "123456789012", "https://123456789012.s3-control.mock-region.amazonaws.com/v20180820/bucket/a%24b%24c", ""},
				},
			},
			"DualStack": {
				useDualstack: true,
				tests: []s3controlEndpointTest{
					{"abc", "123456789012", "https://123456789012.s3-control.dualstack.mock-region.amazonaws.com/v20180820/bucket/abc", ""},
					{"a.b.c", "123456789012", "https://123456789012.s3-control.dualstack.mock-region.amazonaws.com/v20180820/bucket/a.b.c", ""},
					{"a$b$c", "123456789012", "https://123456789012.s3-control.dualstack.mock-region.amazonaws.com/v20180820/bucket/a%24b%24c", ""},
				},
			},
		},

		"immutable endpoint": {
			"default": {
				customEndpoint: &aws.Endpoint{
					URL:               "https://example.region.amazonaws.com",
					HostnameImmutable: true,
				},
				tests: []s3controlEndpointTest{
					{"abc", "123456789012", "https://example.region.amazonaws.com/v20180820/bucket/abc", ""},
					{"a.b.c", "123456789012", "https://example.region.amazonaws.com/v20180820/bucket/a.b.c", ""},
					{"a$b$c", "123456789012", "https://example.region.amazonaws.com/v20180820/bucket/a%24b%24c", ""},
				},
			},
			"DualStack": {
				useDualstack: true,
				customEndpoint: &aws.Endpoint{
					URL:               "https://example.region.amazonaws.com",
					HostnameImmutable: true,
				},
				tests: []s3controlEndpointTest{
					{"abc", "123456789012", "https://example.region.amazonaws.com/v20180820/bucket/abc", ""},
					{"a.b.c", "123456789012", "https://example.region.amazonaws.com/v20180820/bucket/a.b.c", ""},
					{"a$b$c", "123456789012", "https://example.region.amazonaws.com/v20180820/bucket/a%24b%24c", ""},
				},
			},
		},
	}

	for suitName, cs := range cases {
		t.Run(suitName, func(t *testing.T) {
			for unitName, c := range cs {
				t.Run(unitName, func(t *testing.T) {

					options := s3control.Options{
						Credentials: unit.StubCredentialsProvider{},
						Retryer:     aws.NoOpRetryer{},
						Region:      "mock-region",

						HTTPClient: smithyhttp.NopClient{},

						UseDualstack: c.useDualstack,
					}

					if c.customEndpoint != nil {
						options.EndpointResolver = s3control.EndpointResolverFunc(
							func(region string, options s3control.EndpointResolverOptions) (aws.Endpoint, error) {
								return *c.customEndpoint, nil
							})
					}

					svc := s3control.New(options)
					for i, test := range c.tests {
						t.Run(strconv.Itoa(i), func(t *testing.T) {
							fm := requestRetrieverMiddleware{}
							_, err := svc.DeleteBucket(context.Background(),
								&s3control.DeleteBucketInput{
									Bucket:    &test.bucket,
									AccountId: &test.accountID,
								},
								func(options *s3control.Options) {
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
								t.Fatalf("expect URL %s, got %s", e, a)
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
