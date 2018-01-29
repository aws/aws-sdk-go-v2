// +build bench

package restxml_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bytes"
	"encoding/xml"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restxml"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	cloudfrontSvc *cloudfront.CloudFront
	s3Svc         *s3.S3
)

func TestMain(m *testing.M) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cfg := unit.Config()

	cfg.Credentials = aws.NewStaticCredentialsProvider("Key", "Secret", "Token")
	cfg.Region = endpoints.UsWest2RegionID
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	cloudfrontSvc = cloudfront.New(cfg)
	s3Svc = s3.New(cfg)
	s3Svc.ForcePathStyle = true

	c := m.Run()
	server.Close()
	os.Exit(c)
}

func BenchmarkRESTXMLBuild_Complex_CFCreateDistro(b *testing.B) {
	params := cloudfrontCreateDistributionInput()

	benchRESTXMLBuild(b, func() *aws.Request {
		req := cloudfrontSvc.CreateDistributionRequest(params)
		return req.Request
	})
}

func BenchmarkRESTXMLBuild_Simple_CFDeleteDistro(b *testing.B) {
	params := cloudfrontDeleteDistributionInput()

	benchRESTXMLBuild(b, func() *aws.Request {
		req := cloudfrontSvc.DeleteDistributionRequest(params)
		return req.Request
	})
}

func BenchmarkRESTXMLBuild_REST_S3HeadObject(b *testing.B) {
	params := s3HeadObjectInput()

	benchRESTXMLBuild(b, func() *aws.Request {
		req := s3Svc.HeadObjectRequest(params)
		return req.Request
	})
}

func BenchmarkRESTXMLBuild_XML_S3PutObjectAcl(b *testing.B) {
	params := s3PutObjectAclInput()

	benchRESTXMLBuild(b, func() *aws.Request {
		req := s3Svc.PutObjectAclRequest(params)
		return req.Request
	})
}

func BenchmarkRESTXMLRequest_Complex_CFCreateDistro(b *testing.B) {
	benchRESTXMLRequest(b, func() *aws.Request {
		req := cloudfrontSvc.CreateDistributionRequest(cloudfrontCreateDistributionInput())
		return req.Request
	})
}

func BenchmarkRESTXMLRequest_Simple_CFDeleteDistro(b *testing.B) {
	benchRESTXMLRequest(b, func() *aws.Request {
		req := cloudfrontSvc.DeleteDistributionRequest(cloudfrontDeleteDistributionInput())
		return req.Request
	})
}

func BenchmarkRESTXMLRequest_REST_S3HeadObject(b *testing.B) {
	benchRESTXMLRequest(b, func() *aws.Request {
		req := s3Svc.HeadObjectRequest(s3HeadObjectInput())
		return req.Request
	})
}

func BenchmarkRESTXMLRequest_XML_S3PutObjectAcl(b *testing.B) {
	benchRESTXMLRequest(b, func() *aws.Request {
		req := s3Svc.PutObjectAclRequest(s3PutObjectAclInput())
		return req.Request
	})
}

func BenchmarkEncodingXML_Simple(b *testing.B) {
	params := cloudfrontDeleteDistributionInput()

	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		encoder := xml.NewEncoder(buf)
		if err := encoder.Encode(params); err != nil {
			b.Fatal("Unexpected error", err)
		}
	}
}

func benchRESTXMLBuild(b *testing.B, reqFn func() *aws.Request) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := reqFn()
		restxml.Build(req)
		if req.Error != nil {
			b.Fatal("Unexpected error", req.Error)
		}
	}
}

func benchRESTXMLRequest(b *testing.B, reqFn func() *aws.Request) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := reqFn().Send()
		if err != nil {
			b.Fatal("Unexpected error", err)
		}
	}
}

func cloudfrontCreateDistributionInput() *cloudfront.CreateDistributionInput {
	return &cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{ // Required
			CallerReference: aws.String("string"), // Required
			Comment:         aws.String("string"), // Required
			DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{ // Required
				ForwardedValues: &cloudfront.ForwardedValues{ // Required
					Cookies: &cloudfront.CookiePreference{ // Required
						Forward: cloudfront.ItemSelection("ItemSelection"), // Required
						WhitelistedNames: &cloudfront.CookieNames{
							Quantity: aws.Int64(1), // Required
							Items: []string{
								"string", // Required
								// More values...
							},
						},
					},
					QueryString: aws.Bool(true), // Required
					Headers: &cloudfront.Headers{
						Quantity: aws.Int64(1), // Required
						Items: []string{
							"string", // Required
							// More values...
						},
					},
				},
				MinTTL:         aws.Int64(1),         // Required
				TargetOriginId: aws.String("string"), // Required
				TrustedSigners: &cloudfront.TrustedSigners{ // Required
					Enabled:  aws.Bool(true), // Required
					Quantity: aws.Int64(1),   // Required
					Items: []string{
						"string", // Required
						// More values...
					},
				},
				ViewerProtocolPolicy: cloudfront.ViewerProtocolPolicy("ViewerProtocolPolicy"), // Required
				AllowedMethods: &cloudfront.AllowedMethods{
					Items: []cloudfront.Method{ // Required
						cloudfront.Method("string"), // Required
						// More values...
					},
					Quantity: aws.Int64(1), // Required
					CachedMethods: &cloudfront.CachedMethods{
						Items: []cloudfront.Method{ // Required
							cloudfront.Method("string"), // Required
							// More values...
						},
						Quantity: aws.Int64(1), // Required
					},
				},
				DefaultTTL:      aws.Int64(1),
				MaxTTL:          aws.Int64(1),
				SmoothStreaming: aws.Bool(true),
			},
			Enabled: aws.Bool(true), // Required
			Origins: &cloudfront.Origins{ // Required
				Quantity: aws.Int64(1), // Required
				Items: []cloudfront.Origin{
					{ // Required
						DomainName: aws.String("string"), // Required
						Id:         aws.String("string"), // Required
						CustomOriginConfig: &cloudfront.CustomOriginConfig{
							HTTPPort:             aws.Int64(1),                                            // Required
							HTTPSPort:            aws.Int64(1),                                            // Required
							OriginProtocolPolicy: cloudfront.OriginProtocolPolicy("OriginProtocolPolicy"), // Required
						},
						OriginPath: aws.String("string"),
						S3OriginConfig: &cloudfront.S3OriginConfig{
							OriginAccessIdentity: aws.String("string"), // Required
						},
					},
					// More values...
				},
			},
			Aliases: &cloudfront.Aliases{
				Quantity: aws.Int64(1), // Required
				Items: []string{
					"string", // Required
					// More values...
				},
			},
			CacheBehaviors: &cloudfront.CacheBehaviors{
				Quantity: aws.Int64(1), // Required
				Items: []cloudfront.CacheBehavior{
					{ // Required
						ForwardedValues: &cloudfront.ForwardedValues{ // Required
							Cookies: &cloudfront.CookiePreference{ // Required
								Forward: cloudfront.ItemSelection("ItemSelection"), // Required
								WhitelistedNames: &cloudfront.CookieNames{
									Quantity: aws.Int64(1), // Required
									Items: []string{
										"string", // Required
										// More values...
									},
								},
							},
							QueryString: aws.Bool(true), // Required
							Headers: &cloudfront.Headers{
								Quantity: aws.Int64(1), // Required
								Items: []string{
									"string", // Required
									// More values...
								},
							},
						},
						MinTTL:         aws.Int64(1),         // Required
						PathPattern:    aws.String("string"), // Required
						TargetOriginId: aws.String("string"), // Required
						TrustedSigners: &cloudfront.TrustedSigners{ // Required
							Enabled:  aws.Bool(true), // Required
							Quantity: aws.Int64(1),   // Required
							Items: []string{
								"string", // Required
								// More values...
							},
						},
						ViewerProtocolPolicy: cloudfront.ViewerProtocolPolicy("ViewerProtocolPolicy"), // Required
						AllowedMethods: &cloudfront.AllowedMethods{
							Items: []cloudfront.Method{ // Required
								cloudfront.Method("string"), // Required
								// More values...
							},
							Quantity: aws.Int64(1), // Required
							CachedMethods: &cloudfront.CachedMethods{
								Items: []cloudfront.Method{ // Required
									cloudfront.Method("string"), // Required
									// More values...
								},
								Quantity: aws.Int64(1), // Required
							},
						},
						DefaultTTL:      aws.Int64(1),
						MaxTTL:          aws.Int64(1),
						SmoothStreaming: aws.Bool(true),
					},
					// More values...
				},
			},
			CustomErrorResponses: &cloudfront.CustomErrorResponses{
				Quantity: aws.Int64(1), // Required
				Items: []cloudfront.CustomErrorResponse{
					{ // Required
						ErrorCode:          aws.Int64(1), // Required
						ErrorCachingMinTTL: aws.Int64(1),
						ResponseCode:       aws.String("string"),
						ResponsePagePath:   aws.String("string"),
					},
					// More values...
				},
			},
			DefaultRootObject: aws.String("string"),
			Logging: &cloudfront.LoggingConfig{
				Bucket:         aws.String("string"), // Required
				Enabled:        aws.Bool(true),       // Required
				IncludeCookies: aws.Bool(true),       // Required
				Prefix:         aws.String("string"), // Required
			},
			PriceClass: cloudfront.PriceClass("PriceClass"),
			Restrictions: &cloudfront.Restrictions{
				GeoRestriction: &cloudfront.GeoRestriction{ // Required
					Quantity:        aws.Int64(1),                                        // Required
					RestrictionType: cloudfront.GeoRestrictionType("GeoRestrictionType"), // Required
					Items: []string{
						"string", // Required
						// More values...
					},
				},
			},
			ViewerCertificate: &cloudfront.ViewerCertificate{
				CloudFrontDefaultCertificate: aws.Bool(true),
				IAMCertificateId:             aws.String("string"),
				MinimumProtocolVersion:       cloudfront.MinimumProtocolVersion("MinimumProtocolVersion"),
				SSLSupportMethod:             cloudfront.SSLSupportMethod("SSLSupportMethod"),
			},
		},
	}
}

func cloudfrontDeleteDistributionInput() *cloudfront.DeleteDistributionInput {
	return &cloudfront.DeleteDistributionInput{
		Id:      aws.String("string"), // Required
		IfMatch: aws.String("string"),
	}
}

func s3HeadObjectInput() *s3.HeadObjectInput {
	return &s3.HeadObjectInput{
		Bucket:    aws.String("somebucketname"),
		Key:       aws.String("keyname"),
		VersionId: aws.String("someVersion"),
		IfMatch:   aws.String("IfMatch"),
	}
}

func s3PutObjectAclInput() *s3.PutObjectAclInput {
	return &s3.PutObjectAclInput{
		Bucket: aws.String("somebucketname"),
		Key:    aws.String("keyname"),
		AccessControlPolicy: &s3.AccessControlPolicy{
			Grants: []s3.Grant{
				{
					Grantee: &s3.Grantee{
						DisplayName:  aws.String("someName"),
						EmailAddress: aws.String("someAddr"),
						ID:           aws.String("someID"),
						Type:         s3.TypeCanonicalUser,
						URI:          aws.String("someURI"),
					},
					Permission: s3.PermissionWrite,
				},
			},
			Owner: &s3.Owner{
				DisplayName: aws.String("howdy"),
				ID:          aws.String("someID"),
			},
		},
	}
}
