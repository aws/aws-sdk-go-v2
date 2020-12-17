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
		expectedErr                string
		expectedReqURL             string
		expectedSigningName        string
		expectedSigningRegion      string
		expectedHeaderForOutpostID string
		expectedHeaderForAccountID bool
	}{
		"Outpost AccessPoint with no S3UseARNRegion flag set": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedReqURL:             "https://s3-outposts.us-west-2.amazonaws.com/v20180820/bucket/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-west-2",
			expectedHeaderForAccountID: true,
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"Outpost AccessPoint Cross-Region Enabled": {
			bucket: "arn:aws:s3-outposts:us-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts.us-east-1.amazonaws.com/v20180820/bucket/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-east-1",
			expectedHeaderForAccountID: true,
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"Outpost AccessPoint Cross-Region Disabled": {
			bucket: "arn:aws:s3-outposts:us-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "client region does not match provided ARN region",
		},
		"Outpost AccessPoint other partition": {
			bucket: "arn:aws-cn:s3-outposts:cn-north-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "ConfigurationError : client partition does not match provided ARN partition",
		},
		"Outpost AccessPoint us-gov region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-gov-east-1",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts.us-gov-east-1.amazonaws.com/v20180820/bucket/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForAccountID: true,
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"Outpost AccessPoint with client region as Fips": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region: "us-gov-east-1-fips",
			},
			expectedErr: "InvalidARNError : resource ARN not supported for FIPS region",
		},
		"Outpost AccessPoint with client Fips region and use arn region enabled": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-gov-east-1-fips",
				UseARNRegion: true,
			},
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedReqURL:             "https://s3-outposts.us-gov-east-1.amazonaws.com/v20180820/bucket/myaccesspoint",
			expectedHeaderForAccountID: true,
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"Outpost AccessPoint Fips region in Arn": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1-fips:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-gov-east-1-fips",
				UseARNRegion: true,
			},
			expectedErr: "InvalidARNError : resource ARN not supported for FIPS region",
		},
		"Outpost AccessPoint Fips region with valid ARN region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-gov-east-1-fips",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts.us-gov-east-1.amazonaws.com/v20180820/bucket/myaccesspoint",
			expectedSigningName:        "s3-outposts",
			expectedSigningRegion:      "us-gov-east-1",
			expectedHeaderForAccountID: true,
			expectedHeaderForOutpostID: "op-01234567890123456",
		},
		"Outpost AccessPoint with DualStack": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:myaccesspoint",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
				UseDualstack: true,
			},
			expectedErr: "ConfigurationError : client configured for S3 Dual-stack but is not supported with resource ARN",
		},
		"Invalid outpost resource format": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "outpost resource-id not set",
		},
		"Missing access point for outpost resource": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "incomplete outpost resource type",
		},
		"access point": {
			bucket: "myaccesspoint",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedReqURL:             "https://123456789012.s3-control.us-west-2.amazonaws.com/v20180820/bucket/myaccesspoint",
			expectedHeaderForAccountID: true,
			expectedSigningRegion:      "us-west-2",
			expectedSigningName:        "s3",
		},
		"outpost access point with unsupported sub-resource": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:accesspoint:mybucket:object:foo",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "sub resource not supported",
		},
		"Missing outpost identifiers in outpost access point arn": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:accesspoint:myendpoint",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "invalid Amazon s3-outposts ARN",
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
			expectedErr: "client region does not match provided ARN region",
		},
		"Outpost Bucket other partition": {
			bucket: "arn:aws-cn:s3-outposts:cn-north-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-west-2",
				UseARNRegion: true,
			},
			expectedErr: "ConfigurationError : client partition does not match provided ARN partition",
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
		"Outpost Bucket Fips region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region: "us-gov-west-1-fips",
			},
			expectedErr: "ConfigurationError : client region does not match provided ARN region",
		},
		"Outpost Bucket Fips region in Arn": {
			bucket: "arn:aws-us-gov:s3-outposts:fips-us-gov-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-gov-east-1-fips",
				UseARNRegion: true,
			},
			expectedErr: "InvalidARNError : resource ARN not supported for FIPS region",
		},
		"Outpost Bucket Fips region with valid ARN region": {
			bucket: "arn:aws-us-gov:s3-outposts:us-gov-east-1:123456789012:outpost:op-01234567890123456:bucket:mybucket",
			options: s3control.Options{
				Region:       "us-gov-east-1-fips",
				UseARNRegion: true,
			},
			expectedReqURL:             "https://s3-outposts.us-gov-east-1.amazonaws.com/v20180820/bucket/mybucket",
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
			expectedErr: "ConfigurationError : client configured for S3 Dual-stack but is not supported with resource ARN",
		},
		"Missing bucket id": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:outpost:op-01234567890123456:bucket",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "invalid Amazon s3-outposts ARN",
		},
		"Invalid ARN": {
			bucket: "arn:aws:s3-outposts:us-west-2:123456789012:bucket:mybucket",
			options: s3control.Options{
				Region: "us-west-2",
			},
			expectedErr: "invalid Amazon s3-outposts ARN, unknown resource type",
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

			// call an operation
			_, err := svc.GetBucket(ctx, &s3control.GetBucketInput{
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
	cases := map[string]struct {
		options                    s3control.Options
		operation                  func(context.Context, *s3control.Client, *requestRetrieverMiddleware) (interface{}, error)
		expectedReqURL             string
		expectedSigningName        string
		expectedSigningRegion      string
		expectedHeaderForOutpostID string
		expectedErr                string
		expectedHeaderForAccountID bool
	}{
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
			expectedErr: "invalid Amazon s3 ARN, unknown resource type",
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
			req := fm.request
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
