package customizations_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/awslabs/smithy-go/ptr"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3/internal/endpoints"

	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3BucketTest struct {
	bucket string
	key    string
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
					{"a.b.c", "key", "https://s3.mock-region.amazonaws.com/a.b.c/key?x-id=GetObject", "not compatible"},
					{"a$b$c", "key", "https://s3.mock-region.amazonaws.com/a%24b%24c/key?x-id=GetObject", "not compatible"},
				},
			},
			"AccelerateNoSSLTests": {
				useAccelerate: true,
				disableHTTPS:  true,
				tests: []s3BucketTest{
					{"abc", "key", "http://abc.s3-accelerate.amazonaws.com/key?x-id=GetObject", ""},
					{"a.b.c", "key", "http://a.b.c.s3-accelerate.amazonaws.com/key?x-id=GetObject", ""},
					{"a$b$c", "key", "http://s3.mock-region.amazonaws.com/a%24b%24c/key?x-id=GetObject", "not compatible"},
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
					{"a.b.c", "key", "https://s3.dualstack.mock-region.amazonaws.com/a.b.c/key?x-id=GetObject", "not compatible"},
					{"a$b$c", "key", "https://s3.dualstack.mock-region.amazonaws.com/a%24b%24c/key?x-id=GetObject", "not compatible"},
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
			"DualStack": {
				useDualstack: true,
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

func TestEndpointWithARN(t *testing.T) {
	// test cases
	cases := map[string]struct {
		options               s3.Options
		bucket                string
		expectedErr           string
		expectedReqURL        string
		expectedSigningName   string
		expectedSigningRegion string
	}{
		"Outpost AccessPoint with no S3UseARNRegion flag set": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region: "us-west-2",
			},
			expectedReqURL:        "https://myaccesspoint-123456789012.op-01234567890123456.s3-outposts.us-west-2.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3-outposts",
			expectedSigningRegion: "us-west-2",
		},
		"Outpost AccessPoint Cross-Region Enabled": {
			bucket: "arn:aws:s3-outposts:us-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myaccesspoint-123456789012.op-01234567890123456.s3-outposts.us-east-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3-outposts",
			expectedSigningRegion: "us-east-1",
		},
		"Outpost AccessPoint Cross-Region Disabled": {
			bucket: "arn:aws:s3-outposts:us-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region: "us-west-2",
			},
			expectedErr: "client region does not match provided ARN region",
		},
		"Outpost AccessPoint other partition": {
			bucket: "arn:aws-cn:s3-outposts:cn-north-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "ConfigurationError : client partition does not match provided ARN partition",
		},
		"Outpost AccessPoint cn partition": {
			bucket: "arn:aws-cn:s3-outposts:cn-north-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region: "cn-north-1",
			},
			expectedReqURL:        "https://myaccesspoint-123456789012.op-01234567890123456.s3-outposts.cn-north-1.amazonaws.com.cn/testkey?x-id=GetObject",
			expectedSigningName:   "s3-outposts",
			expectedSigningRegion: "cn-north-1",
		},
		"Outpost AccessPoint us-gov region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region:       "us-gov-east-1",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myaccesspoint-123456789012.op-01234567890123456.s3-outposts.us-gov-east-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3-outposts",
			expectedSigningRegion: "us-gov-east-1",
		},
		"Outpost AccessPoint Fips region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region: "fips-us-gov-west-1",
			},
			expectedErr: "ConfigurationError : client region does not match provided ARN region",
		},
		"Outpost AccessPoint Fips region in Arn": {
			bucket: "arn:aws-us-gov:s3-outposts:fips-us-gov-west-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region:          "fips-us-gov-west-1",
				EndpointOptions: endpoints.Options{DisableHTTPS: true},
				UseARNRegion:    true,
			},
			expectedErr: "InvalidARNError : resource ARN not supported for FIPS region",
		},
		"Outpost AccessPoint Fips region with valid ARN region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region:       "fips-us-gov-west-1",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myaccesspoint-123456789012.op-01234567890123456.s3-outposts.us-gov-east-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3-outposts",
			expectedSigningRegion: "us-gov-east-1",
		},
		"Outpost AccessPoint with DualStack": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: true,
			},
			expectedErr: "ConfigurationError : client configured for S3 Dual-stack but is not supported with resource ARN",
		},
		"Outpost AccessPoint with Accelerate": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3.Options{
				Region:        "us-west-2",
				UseAccelerate: true,
			},
			expectedErr: "ConfigurationError : client configured for S3 Accelerate but is not supported with resource ARN",
		},
		"AccessPoint": {
			bucket: "arn:aws:s3:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region: "us-west-2",
			},
			expectedReqURL:        "https://myendpoint-123456789012.s3-accesspoint.us-west-2.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"AccessPoint slash delimiter": {
			bucket: "arn:aws:s3:us-west-2:123456789012:accesspoint/myendpoint",
			options: s3.Options{
				Region: "us-west-2",
			},
			expectedReqURL:        "https://myendpoint-123456789012.s3-accesspoint.us-west-2.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"AccessPoint other partition": {
			bucket: "arn:aws-cn:s3:cn-north-1:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region: "cn-north-1",
			},
			expectedReqURL:        "https://myendpoint-123456789012.s3-accesspoint.cn-north-1.amazonaws.com.cn/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "cn-north-1",
		},
		"AccessPoint Cross-Region Disabled": {
			bucket: "arn:aws:s3:ap-south-1:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region: "us-west-2",
			},
			expectedErr: "client region does not match provided ARN region",
		},
		"AccessPoint Cross-Region Enabled": {
			bucket: "arn:aws:s3:ap-south-1:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myendpoint-123456789012.s3-accesspoint.ap-south-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "ap-south-1",
		},
		"AccessPoint us-east-1": {
			bucket: "arn:aws:s3:us-east-1:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region:       "us-east-1",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myendpoint-123456789012.s3-accesspoint.us-east-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-east-1",
		},
		"AccessPoint us-east-1 cross region": {
			bucket: "arn:aws:s3:us-east-1:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myendpoint-123456789012.s3-accesspoint.us-east-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-east-1",
		},
		"AccessPoint Cross-Partition not supported": {
			bucket: "arn:aws-cn:s3:cn-north-1:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: true,
				UseARNRegion: true,
			},
			expectedErr: "client partition does not match provided ARN partition",
		},
		"AccessPoint DualStack": {
			bucket: "arn:aws:s3:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region:       "us-west-2",
				UseDualstack: true,
			},
			expectedReqURL:        "https://myendpoint-123456789012.s3-accesspoint.dualstack.us-west-2.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"AccessPoint FIPS same region with cross region disabled": {
			bucket: "arn:aws-us-gov:s3:us-gov-west-1:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region: "fips-us-gov-west-1",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					switch region {
					case "fips-us-gov-west-1":
						return aws.Endpoint{
							URL:           "https://s3-fips.us-gov-west-1.amazonaws.com",
							PartitionID:   "aws-us-gov",
							SigningRegion: "us-gov-west-1",
							SigningName:   "s3",
							SigningMethod: "s3v4",
						}, nil
					}
					return aws.Endpoint{}, nil
				}),
			},
			expectedReqURL:        "https://myendpoint-123456789012.s3-accesspoint-fips.us-gov-west-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-gov-west-1",
		},
		"AccessPoint FIPS same region with cross region enabled": {
			bucket: "arn:aws-us-gov:s3:us-gov-west-1:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region: "fips-us-gov-west-1",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					switch region {
					case "fips-us-gov-west-1":
						return aws.Endpoint{
							URL:           "https://s3-fips.us-gov-west-1.amazonaws.com",
							PartitionID:   "aws-us-gov",
							SigningRegion: "us-gov-west-1",
							SigningMethod: "s3v4",
						}, nil
					}
					return aws.Endpoint{}, nil
				}),
				UseARNRegion: true,
			},
			expectedReqURL:        "https://myendpoint-123456789012.s3-accesspoint-fips.us-gov-west-1.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-gov-west-1",
		},
		"AccessPoint FIPS cross region not supported": {
			bucket: "arn:aws-us-gov:s3:us-gov-east-1:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region:       "fips-us-gov-west-1",
				UseARNRegion: true,
			},
			expectedErr: "client configured for FIPS",
		},
		"AccessPoint Accelerate not supported": {
			bucket: "arn:aws:s3:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region:        "us-west-2",
				UseAccelerate: true,
			},
			expectedErr: "client configured for S3 Accelerate",
		},
		"Custom Resolver Without PartitionID in ClientInfo": {
			bucket: "arn:aws:s3:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					switch region {
					case "us-west-2":
						return aws.Endpoint{
							URL:           "https://s3.us-west-2.amazonaws.com",
							SigningRegion: "us-west-2",
							SigningName:   "s3",
							SigningMethod: "s3v4",
						}, nil
					}
					return aws.Endpoint{}, nil
				}),
			},
			expectedErr: "partition id was not found for provided request region",
		},
		"Custom Resolver Without PartitionID in Cross-Region Target": {
			bucket: "arn:aws:s3:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3.Options{
				Region:       "us-east-1",
				UseARNRegion: true,
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					switch region {
					case "us-west-2":
						return aws.Endpoint{
							URL:           "https://s3.us-west-2.amazonaws.com",
							PartitionID:   "aws",
							SigningRegion: "us-west-2",
							SigningName:   "s3",
							SigningMethod: "s3v4",
						}, nil
					case "us-east-1":
						return aws.Endpoint{
							URL:           "https://s3.us-east-1.amazonaws.com",
							SigningRegion: "us-east-1",
							SigningName:   "s3",
							SigningMethod: "s3v4",
						}, nil
					}
					return aws.Endpoint{}, nil
				}),
			},
			expectedErr: "partition id was not found for provided request region",
		},
		"bucket host-style": {
			bucket: "mock-bucket",
			options: s3.Options{
				Region: "us-west-2",
			},
			expectedReqURL:        "https://mock-bucket.s3.us-west-2.amazonaws.com/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"bucket path-style": {
			bucket: "mock-bucket",
			options: s3.Options{
				Region:       "us-west-2",
				UsePathStyle: true,
			},
			expectedReqURL:        "https://s3.us-west-2.amazonaws.com/mock-bucket/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"bucket host-style endpoint with default port": {
			bucket: "mock-bucket",
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           "https://s3.us-west-2.amazonaws.com:443",
						SigningRegion: "us-west-2",
					}, nil
				}),
			},
			expectedReqURL:        "https://mock-bucket.s3.us-west-2.amazonaws.com:443/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"bucket host-style endpoint with non-default port": {
			bucket: "mock-bucket",
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           "https://s3.us-west-2.amazonaws.com:8443",
						SigningRegion: "us-west-2",
					}, nil
				}),
			},
			expectedReqURL:        "https://mock-bucket.s3.us-west-2.amazonaws.com:8443/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"bucket path-style endpoint with default port": {
			bucket: "mock-bucket",
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           "https://s3.us-west-2.amazonaws.com:443",
						SigningRegion: "us-west-2",
					}, nil
				}),
				UsePathStyle: true,
			},
			expectedReqURL:        "https://s3.us-west-2.amazonaws.com:443/mock-bucket/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"bucket path-style endpoint with non-default port": {
			bucket: "mock-bucket",
			options: s3.Options{
				Region: "us-west-2",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           "https://s3.us-west-2.amazonaws.com:8443",
						SigningRegion: "us-west-2",
					}, nil
				}),
				UsePathStyle: true,
			},
			expectedReqURL:        "https://s3.us-west-2.amazonaws.com:8443/mock-bucket/testkey?x-id=GetObject",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			// options
			opts := c.options.Copy()
			opts.Credentials = unit.StubCredentialsProvider{}
			opts.HTTPClient = smithyhttp.NopClient{}
			opts.Retryer = aws.NopRetryer{}

			// build an s3 client
			svc := s3.New(opts)
			// setup a request retriever middleware
			fm := requestRetrieverMiddleware{}

			ctx := context.Background()

			// call an operation
			_, err := svc.GetObject(ctx, &s3.GetObjectInput{
				Bucket: ptr.String(c.bucket),
				Key:    ptr.String("testkey"),
			}, func(options *s3.Options) {
				// append request retriever middleware for request inspection
				options.APIOptions = append(options.APIOptions,
					func(stack *middleware.Stack) error {
						// adds AFTER operation serializer middleware
						stack.Serialize.Insert(&fm, "OperationSerializer", middleware.After)
						return nil
					})
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
			req := fm.request.Build(ctx)
			// verify the built request is as expected
			if e, a := c.expectedReqURL, req.URL.String(); e != a {
				t.Fatalf("expect url %s, got %s", e, a)
			}

			if e, a := c.expectedSigningRegion, fm.signingRegion; !strings.EqualFold(e, a) {
				t.Fatalf("expect signing region as %s, got %s", e, a)
			}

			if e, a := c.expectedSigningName, fm.signingName; !strings.EqualFold(e, a) {
				t.Fatalf("expect signing name as %s, got %s", e, a)
			}
		})

	}
}

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
