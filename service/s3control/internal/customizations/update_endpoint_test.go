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
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/ptr"
	smithyhttp "github.com/aws/smithy-go/transport/http"
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
		},
	}

	for suitName, cs := range cases {
		t.Run(suitName, func(t *testing.T) {
			for unitName, c := range cs {
				t.Run(unitName, func(t *testing.T) {

					options := s3control.Options{
						Credentials: unit.StubCredentialsProvider{},
						Retryer:     aws.NopRetryer{},
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

func TestEndpointWithARN(t *testing.T) {
	// test cases
	cases := map[string]struct {
		options                    s3control.Options
		bucket                     string
		accessPoint                string
		expectedErr                string
		expectedReqURL             string
		expectedSigningName        string
		expectedSigningRegion      string
		expectedHeaderForOutpostID string
		expectedHeaderForAccountID bool
	}{
		"Outpost AccessPoint with no S3UseARNRegion flag set": {
			accessPoint: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedReqURL:             "https://s3-outposts.us-west-2.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForAccountID: true,
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"Outpost AccessPoint Cross-Region Enabled": {
			accessPoint: "arn:aws:s3-outposts:us-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts.us-east-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-east-1",
			expectedHeaderForAccountID: true,
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"Outpost AccessPoint Cross-Region Disabled": {
			accessPoint: "arn:aws:s3-outposts:us-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "Invalid configuration: region from ARN `us-east-1` does not match client region `us-west-2` and UseArnRegion is `false`",
		},
		"Outpost AccessPoint other partition": {
			accessPoint: "arn:aws-cn:s3-outposts:cn-north-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "Client was configured for partition `aws` but ARN has `aws-cn`",
		},
		"Outpost AccessPoint us-gov region": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-gov-east-1",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts.us-gov-east-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForAccountID: true,
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"Outpost AccessPoint with client region as FIPS": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region: "us-gov-east-1",
				EndpointOptions: s3control.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-east-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost AccessPoint with client region as FIPS (ResolvedRegion)": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region: "us-gov-east-1-fips",
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-east-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost AccessPoint with client FIPS and use arn region enabled": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region: "us-gov-east-1",
				EndpointOptions: s3control.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-east-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost AccessPoint with client FIPS (ResolvedRegion) and use arn region enabled": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-gov-east-1-fips",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-east-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost AccessPoint client FIPS and cross region ARN": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-gov-west-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-gov-east-1",
				UseARNRegion: true,
				EndpointOptions: s3control.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-west-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-west-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost AccessPoint client FIPS (ResolvedRegion) and cross region ARN": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-gov-west-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-gov-east-1-fips",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-west-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-west-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost AccessPoint client FIPS with valid ARN region": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region: "us-gov-east-1",
				EndpointOptions: s3control.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-east-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost AccessPoint client FIPS (ResolvedRegion) with valid ARN region": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-gov-east-1-fips",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-east-1.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost AccessPoint with DualStack": {
			accessPoint: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
				UseDualstack: true,
			},
			expectedErr: "Invalid configuration: Outpost Access Points do not support dual-stack",
		},
		"Invalid outpost resource format": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "Invalid ARN: The Outpost Id was not set",
		},
		"Missing access point for outpost resource": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "Invalid ARN: Expected a 4-component resource",
		},
		"access point": {
			accessPoint: "myaccesspoint",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedReqURL:             "https://123456789012.s3-control.us-west-2.amazonaws.com/v20180820/accesspoint/myaccesspoint",
			expectedHeaderForAccountID: true,
			expectedSigningRegion:      "us-west-2",
			expectedSigningName:        "s3",
		},
		"outpost access point with unsupported sub-resource": {
			accessPoint: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:mybucket:object:foo",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "sub resource not supported",
		},
		"Missing outpost identifiers in outpost access point arn": {
			accessPoint: "arn:aws:s3-outposts:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "Invalid ARN: Expected a 4-component resource",
		},
		"Outpost Bucket with no S3UseARNRegion flag set": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedReqURL:             "https://s3-outposts.us-west-2.amazonaws.com/v20180820/bucket/mybucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost Bucket Cross-Region Enabled": {
			bucket: "arn:aws:s3-outposts:us-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts.us-east-1.amazonaws.com/v20180820/bucket/mybucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost Bucket Cross-Region Disabled": {
			bucket: "arn:aws:s3-outposts:us-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "Invalid configuration: region from ARN `us-east-1` does not match client region `us-west-2` and UseArnRegion is `false`",
		},
		"Outpost Bucket other partition": {
			bucket: "arn:aws-cn:s3-outposts:cn-north-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "Client was configured for partition `aws` but ARN has `aws-cn`",
		},
		"Outpost Bucket us-gov region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-gov-east-1",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts.us-gov-east-1.amazonaws.com/v20180820/bucket/mybucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost Bucket client FIPS, cross-region ARN": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region: "us-gov-west-1",
				EndpointOptions: s3control.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			expectedErr: "Invalid configuration: region from ARN `us-gov-east-1` does not match client region `us-gov-west-1` and UseArnRegion is `false`",
		},
		"Outpost Bucket client FIPS (ResolvedRegion), cross-region ARN": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region: "us-gov-west-1-fips",
			},
			expectedErr: "Invalid configuration: region from ARN `us-gov-east-1` does not match client region `us-gov-west-1` and UseArnRegion is `false`",
		},
		"Outpost Bucket client FIPS with non cross-region ARN region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-gov-east-1",
				UseARNRegion: true,
				EndpointOptions: s3control.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-east-1.amazonaws.com/v20180820/bucket/mybucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost Bucket client FIPS (ResolvedRegion) with non cross-region ARN region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-gov-east-1-fips",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-east-1.amazonaws.com/v20180820/bucket/mybucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"Outpost Bucket with DualStack": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-west-2",
				UseDualstack: true,
			},
			expectedErr: "Invalid configuration: Outpost buckets do not support dual-stack",
		},
		"Missing bucket id": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:bucket",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "Invalid ARN: expected a bucket name",
		},
		"Invalid ARN": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:bucket:mybucket",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "Invalid ARN: Expected a 4-component resource",
		},
		"Invalid Outpost Bucket ARN with FIPS pseudo-region (prefix)": {
			bucket: "arn:aws:s3-outposts:fips-us-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "FIPS region not allowed in ARN",
		},
		"Invalid Outpost Bucket ARN with FIPS pseudo-region (suffix)": {
			bucket: "arn:aws:s3-outposts:us-east-1-fips:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "FIPS region not allowed in ARN",
		},
		"Invalid Outpost AccessPoint ARN with FIPS pseudo-region (prefix)": {
			accessPoint: "arn:aws-us-gov:s3-outposts:fips-us-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "FIPS region not allowed in ARN",
		},
		"Invalid Outpost AccessPoint ARN with FIPS pseudo-region (suffix)": {
			accessPoint: "arn:aws-us-gov:s3-outposts:us-east-1-fips:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "FIPS region not allowed in ARN",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			// options
			opts := c.options.Copy()
			opts.Credentials = unit.StubCredentialsProvider{}
			opts.HTTPClient = smithyhttp.NopClient{}
			opts.Retryer = aws.NopRetryer{}

			// build an s3control client
			svc := s3control.New(opts)
			// setup a request retriever middleware
			fm := requestRetrieverMiddleware{}

			ctx := context.Background()

			var err error
			if len(c.accessPoint) > 0 {
				_, err = svc.GetAccessPoint(ctx, &s3control.GetAccessPointInput{
					Name:      ptr.String(c.accessPoint),
					AccountId: ptr.String("123456789012"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(&fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			} else {
				_, err = svc.GetBucket(ctx, &s3control.GetBucketInput{
					Bucket:    ptr.String(c.bucket),
					AccountId: ptr.String("123456789012"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(&fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			}

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

			if c.expectedHeaderForAccountID {
				if e, a := "123456789012", req.Header.Get("x-amz-account-id"); e != a {
					t.Fatalf("expect account id header value to be %v, got %v", e, a)
				}
			}

			if e, a := c.expectedHeaderForOutpostID, req.Header.Get("x-amz-outpost-id"); e != a {
				t.Fatalf("expect outpost id header value to be %v, got %v", e, a)
			}
		})

	}
}

type requestRetrieverMiddleware struct {
	request       *smithyhttp.Request
	signingRegion string
	signingName   string
}

func TestCustomEndpoint_SpecialOperations(t *testing.T) {
	cases := map[string]testCaseForEndpointCustomization{
		"CreateBucketOperation": {
			options: s3control.Options{
				Region: "us-west-2",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.CreateBucket(ctx, &s3control.CreateBucketInput{
					Bucket:    aws.String("mockBucket"),
					OutpostId: aws.String("op-01234567890123456"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedReqURL:             "https://s3-outposts.us-west-2.amazonaws.com/v20180820/bucket/mockBucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: false,
		},
		"ListRegionalBucketsOperation": {
			options: s3control.Options{
				Region: "us-west-2",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.ListRegionalBuckets(ctx, &s3control.ListRegionalBucketsInput{
					AccountId: aws.String("123456789012"),
					OutpostId: aws.String("op-01234567890123456"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedReqURL:             "https://s3-outposts.us-west-2.amazonaws.com/v20180820/bucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"ListRegionalBucketsOperation with client FIPS": {
			options: s3control.Options{
				Region: "us-west-2",
				EndpointOptions: s3control.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.ListRegionalBuckets(ctx, &s3control.ListRegionalBucketsInput{
					AccountId: aws.String("123456789012"),
					OutpostId: aws.String("op-01234567890123456"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedReqURL:             "https://s3-outposts-fips.us-west-2.amazonaws.com/v20180820/bucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"ListRegionalBucketsOperation with client FIPS (ResolvedRegion)": {
			options: s3control.Options{
				Region: "us-west-2-fips",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.ListRegionalBuckets(ctx, &s3control.ListRegionalBucketsInput{
					AccountId: aws.String("123456789012"),
					OutpostId: aws.String("op-01234567890123456"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedReqURL:             "https://s3-outposts-fips.us-west-2.amazonaws.com/v20180820/bucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: true,
		},
		"CreateBucketOperation with client FIPS": {
			options: s3control.Options{
				Region: "us-gov-west-1",
				EndpointOptions: s3control.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.CreateBucket(ctx, &s3control.CreateBucketInput{
					Bucket:    aws.String("mockBucket"),
					OutpostId: aws.String("op-01234567890123456"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-west-1.amazonaws.com/v20180820/bucket/mockBucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-west-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: false,
		},
		"CreateBucketOperation with client FIPS (ResolvedRegion)": {
			options: s3control.Options{
				Region: "us-gov-west-1-fips",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.CreateBucket(ctx, &s3control.CreateBucketInput{
					Bucket:    aws.String("mockBucket"),
					OutpostId: aws.String("op-01234567890123456"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedReqURL:             "https://s3-outposts-fips.us-gov-west-1.amazonaws.com/v20180820/bucket/mockBucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-west-1",
			expectedHeaderForOutpostID: "op-01234567890123456",
			expectedHeaderForAccountID: false,
		},
		"CreateAccessPoint bucket arn": {
			options: s3control.Options{
				Region: "us-west-2",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.CreateAccessPoint(ctx, &s3control.CreateAccessPointInput{
					AccountId: aws.String("123456789012"),
					Bucket:    aws.String("arn:aws:s3:us-west-2:123456789012:bucket:mockBucket"),
					Name:      aws.String("mockName"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedErr: "Endpoint resolution failed. Invalid operation or environment input",
		},
		"CreateAccessPoint outpost bucket arn": {
			options: s3control.Options{
				Region: "us-west-2",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.CreateAccessPoint(ctx, &s3control.CreateAccessPointInput{
					AccountId: aws.String("123456789012"),
					Bucket:    aws.String("arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:bucket:mockBucket"),
					Name:      aws.String("mockName"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedReqURL:             "https://s3-outposts.us-west-2.amazonaws.com/v20180820/accesspoint/mockName",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"CreateAccessPoint outpost bucket arn, client FIPS": {
			options: s3control.Options{
				Region: "us-west-2",
				EndpointOptions: s3control.EndpointResolverOptions{
					UseFIPSEndpoint: aws.FIPSEndpointStateEnabled,
				},
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.CreateAccessPoint(ctx, &s3control.CreateAccessPointInput{
					AccountId: aws.String("123456789012"),
					Bucket:    aws.String("arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:bucket:mockBucket"),
					Name:      aws.String("mockName"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedReqURL:             "https://s3-outposts-fips.us-west-2.amazonaws.com/v20180820/accesspoint/mockName",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"CreateAccessPoint outpost bucket arn, client FIPS (ResolvedRegion)": {
			options: s3control.Options{
				Region: "us-west-2-fips",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.CreateAccessPoint(ctx, &s3control.CreateAccessPointInput{
					AccountId: aws.String("123456789012"),
					Bucket:    aws.String("arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:bucket:mockBucket"),
					Name:      aws.String("mockName"),
				}, func(options *s3control.Options) {
					// append request retriever middleware for request inspection
					options.APIOptions = append(options.APIOptions,
						func(stack *middleware.Stack) error {
							// adds AFTER operation serializer middleware
							stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
							return nil
						})
				})
			},
			expectedReqURL:             "https://s3-outposts-fips.us-west-2.amazonaws.com/v20180820/accesspoint/mockName",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			runValidations(t, c)
		})
	}
}

func runValidations(t *testing.T, c testCaseForEndpointCustomization) {
	// options
	opts := c.options.Copy()
	opts.Credentials = unit.StubCredentialsProvider{}
	opts.HTTPClient = smithyhttp.NopClient{}
	opts.Retryer = aws.NopRetryer{}

	// build an s3control client
	svc := s3control.New(opts)
	// setup a request retriever middleware
	fm := requestRetrieverMiddleware{}

	ctx := context.Background()

	// call an operation
	_, err := c.operation(ctx, svc, &fm)

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

	if c.expectedHeaderForAccountID {
		if e, a := "123456789012", req.Header.Get("x-amz-account-id"); e != a {
			t.Fatalf("expect account id header value to be %v, got %v", e, a)
		}
	}

	if e, a := c.expectedHeaderForOutpostID, req.Header.Get("x-amz-outpost-id"); e != a {
		t.Fatalf("expect outpost id header value to be %v, got %v", e, a)
	}
}

type testCaseForEndpointCustomization struct {
	options                    s3control.Options
	operation                  func(context.Context, *s3control.Client, *requestRetrieverMiddleware) (interface{}, error)
	expectedReqURL             string
	expectedSigningName        string
	expectedSigningRegion      string
	expectedHeaderForOutpostID string
	expectedErr                string
	expectedHeaderForAccountID bool
}

func TestVPC_CustomEndpoint(t *testing.T) {
	account := "123456789012"
	cases := map[string]testCaseForEndpointCustomization{
		"standard GetAccesspoint with custom endpoint url": {
			options: s3control.Options{
				EndpointResolver: s3control.EndpointResolverFromURL("https://beta.example.com"),
				Region:           "us-west-2",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.GetAccessPoint(ctx, &s3control.GetAccessPointInput{
					AccountId: aws.String(account),
					Name:      aws.String("apname"),
				}, addRequestRetriever(fm))
			},
			expectedReqURL:        "https://123456789012.beta.example.com/v20180820/accesspoint/apname",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"Outpost Accesspoint ARN with GetAccesspoint and custom endpoint url": {
			options: s3control.Options{
				EndpointResolver: s3control.EndpointResolverFromURL(
					"https://beta.example.com",
				),
				Region: "us-west-2",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.GetAccessPoint(ctx, &s3control.GetAccessPointInput{
					AccountId: aws.String(account),
					Name:      aws.String("arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint"),
				}, addRequestRetriever(fm))
			},
			expectedReqURL:             "https://beta.example.com/v20180820/accesspoint/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"standard CreateBucket with custom endpoint url": {
			options: s3control.Options{
				EndpointResolver: s3control.EndpointResolverFromURL("https://beta.example.com"),
				Region:           "us-west-2",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.CreateBucket(ctx, &s3control.CreateBucketInput{
					Bucket:    aws.String("bucketname"),
					OutpostId: aws.String("op-01234567890123456"),
				}, addRequestRetriever(fm))
			},
			expectedReqURL:             "https://beta.example.com/v20180820/bucket/bucketname",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"Outpost Accesspoint for GetBucket with custom endpoint url": {
			options: s3control.Options{
				EndpointResolver: s3control.EndpointResolverFromURL("https://beta.example.com"),
				Region:           "us-west-2",
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.GetBucket(ctx, &s3control.GetBucketInput{
					Bucket: aws.String("arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:bucket:mybucket"),
				}, addRequestRetriever(fm))
			},
			expectedReqURL:             "https://beta.example.com/v20180820/bucket/mybucket",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"GetAccesspoint with dualstack and custom endpoint url": {
			options: s3control.Options{
				EndpointResolver: s3control.EndpointResolverFromURL("https://beta.example.com"),
				Region:           "us-west-2",
				UseDualstack:     true,
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.GetAccessPoint(ctx, &s3control.GetAccessPointInput{
					AccountId: aws.String(account),
					Name:      aws.String("apname"),
				}, addRequestRetriever(fm))
			},
			expectedReqURL:        "https://123456789012.beta.example.com/v20180820/accesspoint/apname",
			expectedSigningName:   "s3",
			expectedSigningRegion: "us-west-2",
		},
		"GetAccesspoint with Outposts accesspoint ARN and dualstack": {
			options: s3control.Options{
				EndpointResolver: s3control.EndpointResolverFromURL("https://beta.example.com"),
				Region:           "us-west-2",
				UseDualstack:     true,
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.GetAccessPoint(ctx, &s3control.GetAccessPointInput{
					AccountId: aws.String(account),
					Name:      aws.String("arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint"),
				}, addRequestRetriever(fm))
			},
			expectedErr: "client configured for S3 Dual-stack but is not supported with resource ARN",
		},
		"standard CreateBucket with dualstack": {
			options: s3control.Options{
				EndpointResolver: s3control.EndpointResolverFromURL("https://beta.example.com"),
				Region:           "us-west-2",
				UseDualstack:     true,
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.CreateBucket(ctx, &s3control.CreateBucketInput{
					Bucket:    aws.String("bucketname"),
					OutpostId: aws.String("op-1234567890123456"),
				}, addRequestRetriever(fm))
			},
			expectedErr: " dualstack is not supported for outposts request",
		},
		"GetBucket with Outpost bucket ARN": {
			options: s3control.Options{
				EndpointResolver: s3control.EndpointResolverFromURL("https://beta.example.com"),
				Region:           "us-west-2",
				UseDualstack:     true,
			},
			operation: func(ctx context.Context, svc *s3control.Client, fm *requestRetrieverMiddleware) (interface{}, error) {
				return svc.GetBucket(ctx, &s3control.GetBucketInput{
					Bucket: aws.String("arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:bucket:mybucket"),
				}, addRequestRetriever(fm))
			},
			expectedErr: "client configured for S3 Dual-stack but is not supported with resource ARN",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			runValidations(t, c)
		})
	}
}

func TestInputIsNotModified(t *testing.T) {
	inputBucket := "arn:aws:s3-outposts:us-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket"

	// build options
	opts := s3control.Options{}
	opts.Credentials = unit.StubCredentialsProvider{}
	opts.HTTPClient = smithyhttp.NopClient{}
	opts.Retryer = aws.NopRetryer{}
	opts.Region = "us-west-2"
	opts.UseARNRegion = true

	ctx := context.Background()
	fm := requestRetrieverMiddleware{}
	svc := s3control.New(opts)
	params := s3control.DeleteBucketInput{Bucket: ptr.String(inputBucket)}
	_, err := svc.DeleteBucket(ctx, &params, func(options *s3control.Options) {
		// append request retriever middleware for request inspection
		options.APIOptions = append(options.APIOptions,
			func(stack *middleware.Stack) error {
				// adds AFTER operation serializer middleware
				stack.Serialize.Insert(&fm, "OperationSerializer", middleware.After)
				return nil
			})
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err.Error())
	}

	// check if req params were modified
	if e, a := *params.Bucket, inputBucket; !strings.EqualFold(e, a) {
		t.Fatalf("expected no modification for operation input, "+
			"expected %v, got %v as bucket input", e, a)
	}

	if params.AccountId != nil {
		t.Fatalf("expected original input to be unmodified, but account id was backfilled")
	}

	req := fm.request.Build(ctx)
	modifiedAccountID := req.Header.Get("x-amz-account-id")
	if len(modifiedAccountID) == 0 {
		t.Fatalf("expected account id to be backfilled/modified, was not")
	}
	if e, a := "123456789012", modifiedAccountID; !strings.EqualFold(e, a) {
		t.Fatalf("unexpected diff in account id backfilled from arn, expected %v, got %v", e, a)
	}
}

func TestUseDualStackClientBehavior(t *testing.T) {
	cases := map[string]testCaseForEndpointCustomization{
		"client options dual-stack false, endpoint resolver dual-stack unset": {
			options: s3control.Options{
				Region:       "us-west-2",
				UseDualstack: false,
			},
			expectedReqURL:        "https://012345678901.s3-control.us-west-2.amazonaws.com/v20180820/bucket/test-bucket",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack true, endpoint resolver dual-stack unset": {
			options: s3control.Options{
				Region:       "us-west-2",
				UseDualstack: true,
			},
			expectedReqURL:        "https://012345678901.s3-control.dualstack.us-west-2.amazonaws.com/v20180820/bucket/test-bucket",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack off, endpoint resolver dual-stack disabled": {
			options: s3control.Options{
				Region: "us-west-2",
				EndpointOptions: s3control.EndpointResolverOptions{
					UseDualStackEndpoint: aws.DualStackEndpointStateDisabled,
				},
			},
			expectedReqURL:        "https://012345678901.s3-control.us-west-2.amazonaws.com/v20180820/bucket/test-bucket",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack off, endpoint resolver dual-stack enabled": {
			options: s3control.Options{
				Region: "us-west-2",
				EndpointOptions: s3control.EndpointResolverOptions{
					UseDualStackEndpoint: aws.DualStackEndpointStateEnabled,
				},
			},
			expectedReqURL:        "https://012345678901.s3-control.dualstack.us-west-2.amazonaws.com/v20180820/bucket/test-bucket",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack on, endpoint resolver dual-stack disabled": {
			options: s3control.Options{
				Region:       "us-west-2",
				UseDualstack: true,
				EndpointOptions: s3control.EndpointResolverOptions{
					UseDualStackEndpoint: aws.DualStackEndpointStateDisabled,
				},
			},
			expectedReqURL:        "https://012345678901.s3-control.us-west-2.amazonaws.com/v20180820/bucket/test-bucket",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
		"client options dual-stack off, endpoint resolver dual-stack on": {
			options: s3control.Options{
				Region:       "us-west-2",
				UseDualstack: false,
				EndpointOptions: s3control.EndpointResolverOptions{
					UseDualStackEndpoint: aws.DualStackEndpointStateEnabled,
				},
			},
			expectedReqURL:        "https://012345678901.s3-control.dualstack.us-west-2.amazonaws.com/v20180820/bucket/test-bucket",
			expectedSigningRegion: "us-west-2",
			expectedSigningName:   "s3",
		},
	}
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			tt.operation = func(ctx context.Context, client *s3control.Client, retrieverMiddleware *requestRetrieverMiddleware) (interface{}, error) {
				return client.GetBucket(ctx, &s3control.GetBucketInput{
					AccountId: aws.String("012345678901"),
					Bucket:    aws.String("test-bucket"),
				}, addRequestRetriever(retrieverMiddleware))
			}
			runValidations(t, tt)
		})
	}
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

var addRequestRetriever = func(fm *requestRetrieverMiddleware) func(options *s3control.Options) {
	return func(options *s3control.Options) {
		// append request retriever middleware for request inspection
		options.APIOptions = append(options.APIOptions,
			func(stack *middleware.Stack) error {
				// adds AFTER operation serializer middleware
				stack.Serialize.Insert(fm, "OperationSerializer", middleware.After)
				return nil
			})
	}
}
