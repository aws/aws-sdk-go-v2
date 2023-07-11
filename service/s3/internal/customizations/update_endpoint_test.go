package customizations_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/internal/v4a"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/internal/endpoints"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/ptr"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type s3BucketTest struct {
	bucket string
	key    string
	url    string
	err    string
}

func Test_UpdateEndpointBuild(t *testing.T) {
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
					{"abc", "key", "https://s3.mock-region.amazonaws.com/abc/key?x-id=GetObject", ""},
					{"a$b$c", "key", "https://s3.mock-region.amazonaws.com/a%24b%24c/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://s3.mock-region.amazonaws.com/a.b.c/key?x-id=GetObject", ""},
					{"a..bc", "key", "https://s3.mock-region.amazonaws.com/a..bc/key?x-id=GetObject", ""},
					{"abc", "k:e,y", "https://s3.mock-region.amazonaws.com/abc/k%3Ae%2Cy?x-id=GetObject", ""},
				},
			},
			"VirtualHostStyleBucket": {
				tests: []s3BucketTest{
					{"abc", "key", "https://abc.s3.mock-region.amazonaws.com/key?x-id=GetObject", ""},
					{"a$b$c", "key", "https://s3.mock-region.amazonaws.com/a%24b%24c/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://s3.mock-region.amazonaws.com/a.b.c/key?x-id=GetObject", ""},
					{"a..bc", "key", "https://s3.mock-region.amazonaws.com/a..bc/key?x-id=GetObject", ""},
					{"abc", "k:e,y", "https://abc.s3.mock-region.amazonaws.com/k%3Ae%2Cy?x-id=GetObject", ""},
				},
			},
			"Accelerate": {
				useAccelerate: true,
				tests: []s3BucketTest{
					{"abc", "key", "https://abc.s3-accelerate.amazonaws.com/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://s3.mock-region.amazonaws.com/a.b.c/key?x-id=GetObject", "cannot be used with"},
					{"a$b$c", "key", "https://s3.mock-region.amazonaws.com/a%24b%24c/key?x-id=GetObject", "cannot be used with"},
				},
			},
			"AccelerateNoSSLTests": {
				useAccelerate: true,
				disableHTTPS:  true,
				tests: []s3BucketTest{
					{"abc", "key", "http://abc.s3-accelerate.amazonaws.com/key?x-id=GetObject", ""},
					// {"a.b.c", "key", "http://a.b.c.s3-accelerate.amazonaws.com/key?x-id=GetObject", ""},
					{"a$b$c", "key", "http://s3.mock-region.amazonaws.com/a%24b%24c/key?x-id=GetObject", "cannot be used with"},
				},
			},
			"DualStack": {
				useDualstack: true,
				tests: []s3BucketTest{
					{"abc", "key", "https://abc.s3.dualstack.mock-region.amazonaws.com/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://s3.dualstack.mock-region.amazonaws.com/a.b.c/key?x-id=GetObject", ""},
					{"a$b$c", "key", "https://s3.dualstack.mock-region.amazonaws.com/a%24b%24c/key?x-id=GetObject", ""},
				},
			},
			"DualStackWithPathStyle": {
				useDualstack: true,
				usePathStyle: true,
				tests: []s3BucketTest{
					{"abc", "key", "https://s3.dualstack.mock-region.amazonaws.com/abc/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://s3.dualstack.mock-region.amazonaws.com/a.b.c/key?x-id=GetObject", ""},
					{"a$b$c", "key", "https://s3.dualstack.mock-region.amazonaws.com/a%24b%24c/key?x-id=GetObject", ""},
				},
			},
			"AccelerateWithDualStack": {
				useAccelerate: true,
				useDualstack:  true,
				tests: []s3BucketTest{
					{"abc", "key", "https://abc.s3-accelerate.dualstack.amazonaws.com/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://s3.mock-region.dualstack.amazonaws.com/a.b.c/key?x-id=GetObject", "cannot be used with"},
					{"a$b$c", "key", "https://s3.mock-region.dualstack.amazonaws.com/a%24b%24c/key?x-id=GetObject", "cannot be used with"},
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
					{"abc", "key", "https://example.region.amazonaws.com/abc/key?x-id=GetObject", ""},
					{"a$b$c", "key", "https://example.region.amazonaws.com/a%24b%24c/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://example.region.amazonaws.com/a.b.c/key?x-id=GetObject", ""},
					{"a..bc", "key", "https://example.region.amazonaws.com/a..bc/key?x-id=GetObject", ""},
					{"abc", "k:e,y", "https://example.region.amazonaws.com/abc/k%3Ae%2Cy?x-id=GetObject", ""},
				},
			},
			"VirtualHostStyleBucket": {
				customEndpoint: &aws.Endpoint{
					URL:               "https://example.region.amazonaws.com",
					HostnameImmutable: true,
				},
				tests: []s3BucketTest{
					{"abc", "key", "https://example.region.amazonaws.com/abc/key?x-id=GetObject", ""},
					{"a$b$c", "key", "https://example.region.amazonaws.com/a%24b%24c/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://example.region.amazonaws.com/a.b.c/key?x-id=GetObject", ""},
					{"a..bc", "key", "https://example.region.amazonaws.com/a..bc/key?x-id=GetObject", ""},
					{"abc", "k:e,y", "https://example.region.amazonaws.com/abc/k%3Ae%2Cy?x-id=GetObject", ""},
				},
			},
			"Accelerate": {
				useAccelerate: true,
				customEndpoint: &aws.Endpoint{
					URL:               "https://example.region.amazonaws.com",
					HostnameImmutable: true,
				},
				tests: []s3BucketTest{
					{"abc", "key", "https://example.region.amazonaws.com/abc/key?x-id=GetObject", ""},
					{"a$b$c", "key", "https://example.region.amazonaws.com/a%24b%24c/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://example.region.amazonaws.com/a.b.c/key?x-id=GetObject", ""},
					{"a..bc", "key", "https://example.region.amazonaws.com/a..bc/key?x-id=GetObject", ""},
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
					{"abc", "key", "https://example.region.amazonaws.com/abc/key?x-id=GetObject", ""},
					{"a.b.c", "key", "https://example.region.amazonaws.com/a.b.c/key?x-id=GetObject", ""},
					{"a$b$c", "key", "https://example.region.amazonaws.com/a%24b%24c/key?x-id=GetObject", ""},
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
						Retryer:     aws.NopRetryer{},
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
							_, err := svc.GetObject(context.Background(),
								&s3.GetObjectInput{Bucket: &test.bucket, Key: &test.key},
								func(options *s3.Options) {
									options.APIOptions = append(options.APIOptions,
										func(stack *middleware.Stack) error {
											stack.Serialize.Insert(&fm,
												"OperationSerializer", middleware.After)
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

// test case struct used to test endpoint customizations
type testCaseForEndpointCustomization struct {
	options               s3.Options
	bucket                string
	operation             func(ctx context.Context, svc *s3.Client, fm *requestRetriever) (interface{}, error)
	expectedErr           string
	expectedReqURL        string
	expectedSigningName   string
	expectedSigningRegion string
	expectedHeader        map[string]string
}

func TestEndpointWithARN(t *testing.T) {
	// test cases
	cases := map[string]testCaseForEndpointCustomization{
		"Object Lambda with no UseARNRegion flag set": {
			bucket: "arn:aws:s3-object-lambda:us-west-2:123456789012:accesspoint/myap",
			options: s3.Options{
				Region: "us-west-2",
			},
			expectedReqURL:        "https://myap-123456789012.s3-object-lambda.us-west-2.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3-object-lambda",
			expectedSigningRegion: "us-west-2",
		},
		"Object Lambda with UseARNRegion flag set": {
			bucket: "arn:aws:s3-object-lambda:us-east-1:123456789012:accesspoint/myap",
			options: s3.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myap-123456789012.s3-object-lambda.us-east-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3-object-lambda",
			expectedSigningRegion: "us-east-1",
		},
		"Object Lambda with Cross-Region error": {
			bucket: "arn:aws:s3-object-lambda:us-east-1:123456789012:accesspoint/myap",
			options: s3.Options{
				Region: "us-west-2",
			},
			expectedErr: "region from ARN `us-east-1` does not match client region `us-west-2` and UseArnRegion is `false`",
		},
		"Object Lambda Pseudo-Region with UseARNRegion flag set": {
			bucket: "arn:aws:s3-object-lambda:us-east-1:123456789012:accesspoint/myap",
			options: s3.Options{
				Region:       "aws-global",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myap-123456789012.s3-object-lambda.us-east-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningRegion: "us-east-1",
			expectedSigningName:   "s3-object-lambda",
		},
		"Object Lambda Cross-Region DualStack error": {
			bucket: "arn:aws:s3-object-lambda:us-east-1:123456789012:accesspoint/myap",
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: true,
				UseARNRegion: true,
			},
			expectedErr: "S3 Object Lambda does not support Dual-stack",
		},
		"Object Lambda Cross-Partition error": {
			bucket: "arn:aws-cn:s3-object-lambda:cn-north-1:123456789012:accesspoint/myap",
			options: s3.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "Client was configured for partition `aws` but ARN (`arn:aws-cn:s3-object-lambda:cn-north-1:123456789012:accesspoint/myap`) has `aws-cn`",
		},
		"Object Lambda FIPS": {
			bucket: "arn:aws-us-gov:s3-object-lambda:us-gov-west-1:123456789012:accesspoint/myap",
			options: s3.Options{
				Region: "us-gov-west-1",
				EndpointOptions: endpoints.Options{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			expectedReqURL:        "https://myap-123456789012.s3-object-lambda-fips.us-gov-west-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningRegion: "us-gov-west-1",
			expectedSigningName:   "s3-object-lambda",
		},
		"Object Lambda FIPS (ResolvedRegion)": {
			bucket: "arn:aws-us-gov:s3-object-lambda:us-gov-west-1:123456789012:accesspoint/myap",
			options: s3.Options{
				Region: "fips-us-gov-west-1",
			},
			expectedReqURL:        "https://myap-123456789012.s3-object-lambda-fips.us-gov-west-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningRegion: "us-gov-west-1",
			expectedSigningName:   "s3-object-lambda",
		},
		"Object Lambda FIPS with UseARNRegion flag set": {
			bucket: "arn:aws-us-gov:s3-object-lambda:us-gov-west-1:123456789012:accesspoint/myap",
			options: s3.Options{
				Region:       "us-gov-west-1",
				UseARNRegion: true,
				EndpointOptions: s3.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			expectedReqURL:        "https://myap-123456789012.s3-object-lambda-fips.us-gov-west-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningRegion: "us-gov-west-1",
			expectedSigningName:   "s3-object-lambda",
		},
		"Object Lambda FIPS (ResolvedRegion) with UseARNRegion flag set": {
			bucket: "arn:aws-us-gov:s3-object-lambda:us-gov-west-1:123456789012:accesspoint/myap",
			options: s3.Options{
				Region:       "fips-us-gov-west-1",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myap-123456789012.s3-object-lambda-fips.us-gov-west-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningRegion: "us-gov-west-1",
			expectedSigningName:   "s3-object-lambda",
		},
		"Object Lambda with Accelerate": {
			bucket: "arn:aws:s3-object-lambda:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region:        "us-west-2",
				UseAccelerate: true,
			},
			expectedErr: "S3 Object Lambda does not support S3 Accelerate",
		},
		"Object Lambda with Custom Endpoint Source": {
			bucket: "arn:aws:s3-object-lambda:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           "https://my-domain.com",
						Source:        aws.EndpointSourceCustom,
						SigningName:   "custom-sign-name",
						SigningRegion: region,
					}, nil
				}),
			},
			expectedReqURL:        "https://myendpoint-123456789012.my-domain.com/testkey?x-id=GetObject",
			expectedSigningName:   "custom-sign-name",
			expectedSigningRegion: "us-west-2",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			runValidations(t, c, func(ctx context.Context, svc *s3.Client, fetcher *requestRetriever) (interface{}, error) {
				if c.operation != nil {
					return c.operation(ctx, svc, fetcher)
				}
				return svc.GetObject(ctx, &s3.GetObjectInput{
					Bucket: ptr.String(c.bucket),
					Key:    ptr.String("testkey"),
				}, addRequestRetriever(fetcher))
			})
		})
	}
}

func TestVPC_CustomEndpoint(t *testing.T) {
	cases := map[string]testCaseForEndpointCustomization{
		"standard custom endpoint url": {
			bucket: "bucketname",
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://beta.example.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
				Region: "us-west-2",
			},
			expectedReqURL:        "https://bucketname.beta.example.com/",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"AccessPoint with custom endpoint url": {
			bucket: "arn:aws:s3:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://beta.example.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
				Region: "us-west-2",
			},
			expectedReqURL:        "https://myendpoint-123456789012.beta.example.com/",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"Outpost AccessPoint with custom endpoint url": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://beta.example.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
				Region: "us-west-2",
			},
			expectedReqURL:        "https://myaccesspoint-123456789012.op-01234567890123456.beta.example.com/",
			expectedSigningName:   "s3-outposts",
			expectedSigningRegion: "us-west-2",
		},
		"ListBucket with custom endpoint url": {
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://bucket.vpce-123-abc.s3.us-west-2.vpce.amazonaws.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
				Region: "us-west-2",
			},
			operation: func(ctx context.Context, svc *s3.Client, fm *requestRetriever) (interface{}, error) {
				return svc.ListBuckets(ctx, &s3.ListBucketsInput{}, addRequestRetriever(fm))
			},
			expectedReqURL:        "https://bucket.vpce-123-abc.s3.us-west-2.vpce.amazonaws.com/",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"Path-style addressing with custom endpoint url": {
			bucket: "bucketname",
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://bucket.vpce-123-abc.s3.us-west-2.vpce.amazonaws.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
				Region:       "us-west-2",
				UsePathStyle: true,
			},
			expectedReqURL:        "https://bucket.vpce-123-abc.s3.us-west-2.vpce.amazonaws.com/bucketname",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"Virtual host addressing with custom endpoint url": {
			bucket: "bucketname",
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://bucket.vpce-123-abc.s3.us-west-2.vpce.amazonaws.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
				Region: "us-west-2",
			},
			expectedReqURL:        "https://bucketname.bucket.vpce-123-abc.s3.us-west-2.vpce.amazonaws.com/",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"Access-point with custom endpoint url and use_arn_region set": {
			bucket: "arn:aws:s3:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://accesspoint.vpce-123-abc.s3.us-west-2.vpce.amazonaws.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
				Region:       "eu-west-1",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myendpoint-123456789012.accesspoint.vpce-123-abc.s3.us-west-2.vpce.amazonaws.com/",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"Custom endpoint url with Dualstack": {
			bucket: "bucketname",
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://beta.example.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
				Region:       "us-west-2",
				UseDualstack: true,
			},
			expectedReqURL:        "https://bucketname.beta.example.com/",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"Outpost with custom endpoint url and Dualstack": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://beta.example.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
				Region:       "us-west-2",
				UseDualstack: true,
			},
			expectedErr: "S3 Outposts does not support Dual-stack",
		},
		"Standard custom endpoint url with Immutable Host": {
			bucket: "bucketname",
			options: s3.Options{
				EndpointResolver: s3.EndpointResolverFromURL("https://beta.example.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
					endpoint.HostnameImmutable = true
				}),
				Region: "us-west-2",
			},
			expectedReqURL:        "https://beta.example.com/bucketname",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			runValidations(t, c, func(ctx context.Context, svc *s3.Client, fm *requestRetriever) (interface{}, error) {
				if c.operation != nil {
					return c.operation(ctx, svc, fm)
				}

				return svc.ListObjects(ctx, &s3.ListObjectsInput{
					Bucket: ptr.String(c.bucket),
				}, addRequestRetriever(fm))
			})
		})
	}
}

func TestWriteGetObjectResponse_UpdateEndpoint(t *testing.T) {
	cases := map[string]testCaseForEndpointCustomization{
		"standard endpoint": {
			options: s3.Options{
				Region: "us-west-2",
			},
			expectedReqURL:        "https://test-route.s3-object-lambda.us-west-2.amazonaws.com/WriteGetObjectResponse?x-id=WriteGetObjectResponse",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3-object-lambda",
		},
		"fips endpoint": {
			options: s3.Options{
				Region: "us-gov-west-1",
				EndpointOptions: s3.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			expectedReqURL:        "https://test-route.s3-object-lambda-fips.us-gov-west-1.amazonaws.com/WriteGetObjectResponse?x-id=WriteGetObjectResponse",
			expectedSigningRegion: "us-gov-west-1",
			expectedSigningName:   "s3-object-lambda",
		},
		"fips endpoint (ResolvedRegion)": {
			options: s3.Options{
				Region: "fips-us-gov-west-1",
			},
			expectedReqURL:        "https://test-route.s3-object-lambda-fips.us-gov-west-1.amazonaws.com/WriteGetObjectResponse?x-id=WriteGetObjectResponse",
			expectedSigningRegion: "us-gov-west-1",
			expectedSigningName:   "s3-object-lambda",
		},
		"dualstack endpoint": {
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: true,
			},
			expectedErr: "S3 Object Lambda does not support Dual-stack",
		},
		"accelerate endpoint": {
			options: s3.Options{
				Region:        "us-west-2",
				UseAccelerate: true,
			},
			expectedErr: "S3 Object Lambda does not support S3 Accelerate",
		},
		"custom endpoint": {
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           "https://my-domain.com",
						SigningRegion: region,
						SigningName:   "s3", // incorrect signing name gets overwritten
					}, nil
				}),
			},
			expectedReqURL:        "https://test-route.my-domain.com/WriteGetObjectResponse?x-id=WriteGetObjectResponse",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3-object-lambda",
		},
		"custom endpoint immutable": {
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               "https://test-route.my-domain.com",
						SigningRegion:     region,
						SigningName:       "s3", // incorrect signing name gets overwritten
						HostnameImmutable: true,
					}, nil
				}),
			},
			expectedReqURL:        "https://test-route.my-domain.com/WriteGetObjectResponse?x-id=WriteGetObjectResponse",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3-object-lambda",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			runValidations(t, c, func(ctx context.Context, client *s3.Client, retrieverMiddleware *requestRetriever) (interface{}, error) {
				return client.WriteGetObjectResponse(context.Background(),
					&s3.WriteGetObjectResponseInput{
						RequestRoute: aws.String("test-route"),
						RequestToken: aws.String("test-token"),
					}, addRequestRetriever(retrieverMiddleware))
			})
		})
	}
}

func TestUseDualStackClientBehavior(t *testing.T) {
	cases := map[string]testCaseForEndpointCustomization{
		"client options dual-stack false, endpoint resolver dual-stack unset": {
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: false,
			},
			expectedReqURL:        "https://test-bucket.s3.us-west-2.amazonaws.com/test-key?x-id=GetObject",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack true, endpoint resolver dual-stack unset": {
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: true,
			},
			expectedReqURL:        "https://test-bucket.s3.dualstack.us-west-2.amazonaws.com/test-key?x-id=GetObject",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack off, endpoint resolver dual-stack disabled": {
			options: s3.Options{
				Region: "us-west-2",
				EndpointOptions: s3.EndpointResolverOptions{
					UseDualStackEndpoint: aws.DualStackEndpointStateDisabled,
				},
			},
			expectedReqURL:        "https://test-bucket.s3.us-west-2.amazonaws.com/test-key?x-id=GetObject",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack off, endpoint resolver dual-stack enabled": {
			options: s3.Options{
				Region: "us-west-2",
				EndpointOptions: s3.EndpointResolverOptions{
					UseDualStackEndpoint: aws.DualStackEndpointStateEnabled,
				},
			},
			expectedReqURL:        "https://test-bucket.s3.dualstack.us-west-2.amazonaws.com/test-key?x-id=GetObject",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack on, endpoint resolver dual-stack disabled": {
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: true,
				EndpointOptions: s3.EndpointResolverOptions{
					UseDualStackEndpoint: aws.DualStackEndpointStateDisabled,
				},
			},
			expectedReqURL:        "https://test-bucket.s3.us-west-2.amazonaws.com/test-key?x-id=GetObject",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack off, endpoint resolver dual-stack on": {
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: false,
				EndpointOptions: s3.EndpointResolverOptions{
					UseDualStackEndpoint: aws.DualStackEndpointStateEnabled,
				},
			},
			expectedReqURL:        "https://test-bucket.s3.dualstack.us-west-2.amazonaws.com/test-key?x-id=GetObject",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			runValidations(t, tt, func(ctx context.Context, client *s3.Client, retrieverMiddleware *requestRetriever) (interface{}, error) {
				return client.GetObject(context.Background(),
					&s3.GetObjectInput{
						Bucket: aws.String("test-bucket"),
						Key:    aws.String("test-key"),
					}, addRequestRetriever(retrieverMiddleware))
			})
		})
	}
}

func TestMultiRegionAccessPoints_UpdateEndpoint(t *testing.T) {
	cases := map[string]testCaseForEndpointCustomization{
		"region as us-east-1": {
			options: s3.Options{
				Region: "us-east-1",
			},
			bucket:         "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedReqURL: "https://mfzwi23gnjvgw.mrap.accesspoint.s3-global.amazonaws.com/",
			expectedHeader: map[string]string{
				v4a.AmzRegionSetKey: "*",
			},
			expectedSigningName:   "s3",
			expectedSigningRegion: "*",
		},
		"region as us-west-2": {
			options: s3.Options{
				Region: "us-west-2",
			},
			bucket:         "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedReqURL: "https://mfzwi23gnjvgw.mrap.accesspoint.s3-global.amazonaws.com/",
			expectedHeader: map[string]string{
				v4a.AmzRegionSetKey: "*",
			},
			expectedSigningName:   "s3",
			expectedSigningRegion: "*",
		},
		"region as aws-global": {
			options: s3.Options{
				Region: "aws-global",
			},
			bucket:         "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedReqURL: "https://mfzwi23gnjvgw.mrap.accesspoint.s3-global.amazonaws.com/",
			expectedHeader: map[string]string{
				v4a.AmzRegionSetKey: "*",
			},
			expectedSigningName:   "s3",
			expectedSigningRegion: "*",
		},
		"cn partition": {
			options: s3.Options{
				Region: "cn-north-1",
			},
			bucket:         "arn:aws-cn:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedReqURL: "https://mfzwi23gnjvgw.mrap.accesspoint.s3-global.amazonaws.com.cn/",
			expectedHeader: map[string]string{
				v4a.AmzRegionSetKey: "*",
			},
			expectedSigningName:   "s3",
			expectedSigningRegion: "*",
		},
		"cn partition arn with cross partition client region": {
			options: s3.Options{
				Region: "ap-north-1",
			},
			bucket:         "arn:aws-cn:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedReqURL: "https://mfzwi23gnjvgw.mrap.accesspoint.s3-global.amazonaws.com.cn/",
			expectedHeader: map[string]string{
				v4a.AmzRegionSetKey: "*",
			},
			expectedSigningName:   "s3",
			expectedSigningRegion: "*",
		},
		"region as us-west-2 with mrap disabled": {
			options: s3.Options{
				Region:                         "us-west-2",
				DisableMultiRegionAccessPoints: true,
			},
			bucket:      "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedErr: "Multi-Region access point ARNs are disabled",
		},
		"region as aws-global with mrap disabled": {
			options: s3.Options{
				Region:                         "aws-global",
				DisableMultiRegionAccessPoints: true,
			},
			bucket:      "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedErr: "Multi-Region access point ARNs are disabled",
		},
		"with dualstack": {
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: true,
			},
			bucket:      "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedErr: "client configured for S3 Dual-stack but is not supported with resource",
		},
		"with accelerate": {
			options: s3.Options{
				Region:        "us-west-2",
				UseAccelerate: true,
			},
			bucket:      "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedErr: "client configured for S3 Accelerate but is not supported with resource",
		},
		"access point with no region and mrap disabled": {
			options: s3.Options{
				Region:                         "us-west-2",
				DisableMultiRegionAccessPoints: true,
			},
			bucket:      "arn:aws:s3::123456789012:accesspoint:myendpoint",
			expectedErr: "Multi-Region access point ARNs are disabled",
		},
		"endpoint with no region and disabled mrap": {
			options: s3.Options{
				Region:                         "us-west-2",
				DisableMultiRegionAccessPoints: true,
			},
			bucket:      "arn:aws:s3::123456789012:accesspoint:myendpoint",
			expectedErr: "Multi-Region access point ARNs are disabled",
		},
		"endpoint with no region": {
			options: s3.Options{
				Region: "us-west-2",
			},
			bucket:         "arn:aws:s3::123456789012:accesspoint:myendpoint",
			expectedReqURL: "https://myendpoint.accesspoint.s3-global.amazonaws.com/",
			expectedHeader: map[string]string{
				v4a.AmzRegionSetKey: "*",
			},
			expectedSigningName:   "s3",
			expectedSigningRegion: "*",
		},
		"endpoint containing dot with no region": {
			options: s3.Options{
				Region: "us-west-2",
			},
			bucket:         "arn:aws:s3::123456789012:accesspoint:my.bucket",
			expectedReqURL: "https://my.bucket.accesspoint.s3-global.amazonaws.com/",
			expectedHeader: map[string]string{
				v4a.AmzRegionSetKey: "*",
			},
			expectedSigningName:   "s3",
			expectedSigningRegion: "*",
		},
		"custom endpoint": {
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: s3.EndpointResolverFromURL("https://mockendpoint.amazonaws.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
				}),
			},
			bucket:         "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedReqURL: "https://mfzwi23gnjvgw.mrap.accesspoint.s3-global.amazonaws.com/",
			expectedHeader: map[string]string{
				v4a.AmzRegionSetKey: "*",
			},
			expectedSigningName:   "s3",
			expectedSigningRegion: "*",
		},
		"custom endpoint with hostname immutable": {
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: s3.EndpointResolverFromURL("https://mockendpoint.amazonaws.com", func(endpoint *aws.Endpoint) {
					endpoint.SigningRegion = "us-west-2"
					endpoint.HostnameImmutable = true
				}),
			},
			bucket:         "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedReqURL: "https://mockendpoint.amazonaws.com/arn%3Aaws%3As3%3A%3A123456789012%3Aaccesspoint%3Amfzwi23gnjvgw.mrap",
			expectedHeader: map[string]string{
				v4a.AmzRegionSetKey: "*",
			},
			expectedSigningName:   "s3",
			expectedSigningRegion: "*",
		},
		"with fips client": {
			options: s3.Options{
				Region: "us-west-2",
				EndpointOptions: s3.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			bucket:      "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedErr: "client configured for fips but cross-region resource ARN provided",
		},
		"with fips (ResolvedRegion) client": {
			options: s3.Options{
				Region: "fips-us-west-2",
			},
			bucket:      "arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedErr: "client configured for fips but cross-region resource ARN provided",
		},
		"Accesspoint ARN with region and MRAP disabled": {
			options: s3.Options{
				Region:                         "us-west-2",
				DisableMultiRegionAccessPoints: false,
			},
			bucket:                "arn:aws:s3:us-west-2:123456789012:accesspoint:mfzwi23gnjvgw.mrap",
			expectedReqURL:        "https://mfzwi23gnjvgw.mrap-123456789012.s3-accesspoint.us-west-2.amazonaws.com/",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			runValidations(t, c, func(ctx context.Context, svc *s3.Client, reqRetriever *requestRetriever) (interface{}, error) {
				if c.operation != nil {
					return c.operation(ctx, svc, reqRetriever)
				}

				return svc.ListObjects(ctx, &s3.ListObjectsInput{
					Bucket: ptr.String(c.bucket),
				}, addRequestRetriever(reqRetriever))
			})
		})
	}
}

// addRequestRetriever provides request retriever function - that can be used to fetch request from
// various build steps. Currently we support fetching after serializing and after finalized middlewares.
var addRequestRetriever = func(fm *requestRetriever) func(options *s3.Options) {
	return func(options *s3.Options) {
		// append request retriever middleware for request inspection
		options.APIOptions = append(options.APIOptions,
			func(stack *middleware.Stack) error {
				// adds AFTER operation serializer middleware
				return stack.Serialize.Insert(fm.serializedRequest, "OperationSerializer", middleware.After)
			},
			func(stack *middleware.Stack) error {
				// adds AFTER operation finalize middleware
				return stack.Finalize.Add(fm.signedRequest, middleware.After)
			})
	}
}

// requestRetriever can be used to fetch request within various stages of request.
// currently we support fetching requests after serialization, and after signing.
type requestRetriever struct {
	// serializedRequest retriver should be used to fetch request after Operation serializers are executed.
	serializedRequest *requestRetrieverMiddleware

	// signedRequest retriever should be used to fetch request from Finalize step after
	signedRequest *requestRetrieverMiddleware
}

func runValidations(t *testing.T, c testCaseForEndpointCustomization, operation func(
	context.Context, *s3.Client, *requestRetriever) (interface{}, error)) {
	// options
	opts := c.options.Copy()
	opts.Credentials = unit.StubCredentialsProvider{}
	opts.HTTPClient = smithyhttp.NopClient{}
	opts.Retryer = aws.NopRetryer{}

	// build an s3 client
	svc := s3.New(opts)

	// initialize request fetcher to fetch after input is serialized for request
	serializedRequest := requestRetrieverMiddleware{}

	// initialize request fetcher to fetch request after it is signed
	signedRequest := requestRetrieverMiddleware{}

	ctx := context.Background()

	// call an operation
	_, err := operation(ctx, svc, &requestRetriever{
		serializedRequest: &serializedRequest,
		signedRequest:     &signedRequest,
	})

	// inspect any errors
	if len(c.expectedErr) != 0 {
		if err == nil {
			t.Fatalf("expected error, got none")
		}
		if a, e := err.Error(), c.expectedErr; !strings.Contains(a, e) {
			t.Fatalf("expect error code to contain %q, got %q", e, a)
		}
		return
	}
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	// build the captured request
	req := serializedRequest.request.Build(ctx)
	// verify the built request is as expected
	if e, a := c.expectedReqURL, req.URL.String(); e != a {
		t.Fatalf("expect url %s, got %s", e, a)
	}

	if e, a := c.expectedSigningRegion, serializedRequest.signingRegion; !strings.EqualFold(e, a) {
		t.Fatalf("expect signing region as %s, got %s", e, a)
	}

	if e, a := c.expectedSigningName, serializedRequest.signingName; !strings.EqualFold(e, a) {
		t.Fatalf("expect signing name as %s, got %s", e, a)
	}

	// fetch signed request
	signedReq := signedRequest.request
	// validate if expected headers are present in request
	for key, ev := range c.expectedHeader {
		av := signedReq.Header.Get(key)
		if len(av) == 0 {
			t.Fatalf("expected header %v to be present in %v was not", key, req.Header)
		}
		if !strings.EqualFold(ev, av) {
			t.Fatalf("expected header %v to be %v, got %v instead", key, ev, av)
		}
	}
}

// request retriever middleware is used to fetch request within a stack step.
type requestRetrieverMiddleware struct {
	request       *smithyhttp.Request
	signingRegion string
	signingName   string
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

	rm.signingName = awsmiddleware.GetSigningName(ctx)
	rm.signingRegion = awsmiddleware.GetSigningRegion(ctx)

	return next.HandleSerialize(ctx, in)
}

func (rm *requestRetrieverMiddleware) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}
	rm.request = req

	rm.signingName = awsmiddleware.GetSigningName(ctx)
	rm.signingRegion = awsmiddleware.GetSigningRegion(ctx)

	return next.HandleFinalize(ctx, in)
}
