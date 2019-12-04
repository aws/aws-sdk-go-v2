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
	cfEnums "github.com/aws/aws-sdk-go-v2/service/cloudfront/enums"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Enums "github.com/aws/aws-sdk-go-v2/service/s3/enums"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	cloudfrontSvc *cloudfront.Client
	s3Svc         *s3.Client
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

func cloudfrontCreateDistributionInput() *cfTypes.CreateDistributionInput {
	return &cfTypes.CreateDistributionInput{
		DistributionConfig: &cfTypes.DistributionConfig{ // Required
			CallerReference: aws.String("string"), // Required
			Comment:         aws.String("string"), // Required
			DefaultCacheBehavior: &cfTypes.DefaultCacheBehavior{ // Required
				ForwardedValues: &cfTypes.ForwardedValues{ // Required
					Cookies: &cfTypes.CookiePreference{ // Required
						Forward: cfEnums.ItemSelection("ItemSelection"), // Required
						WhitelistedNames: &cfTypes.CookieNames{
							Quantity: aws.Int64(1), // Required
							Items: []string{
								"string", // Required
								// More values...
							},
						},
					},
					QueryString: aws.Bool(true), // Required
					Headers: &cfTypes.Headers{
						Quantity: aws.Int64(1), // Required
						Items: []string{
							"string", // Required
							// More values...
						},
					},
				},
				MinTTL:         aws.Int64(1),         // Required
				TargetOriginId: aws.String("string"), // Required
				TrustedSigners: &cfTypes.TrustedSigners{ // Required
					Enabled:  aws.Bool(true), // Required
					Quantity: aws.Int64(1),   // Required
					Items: []string{
						"string", // Required
						// More values...
					},
				},
				ViewerProtocolPolicy: cfEnums.ViewerProtocolPolicy("ViewerProtocolPolicy"), // Required
				AllowedMethods: &cfTypes.AllowedMethods{
					Items: []cfEnums.Method{ // Required
						cfEnums.Method("string"), // Required
						// More values...
					},
					Quantity: aws.Int64(1), // Required
					CachedMethods: &cfTypes.CachedMethods{
						Items: []cfEnums.Method{ // Required
							cfEnums.Method("string"), // Required
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
			Origins: &cfTypes.Origins{ // Required
				Quantity: aws.Int64(1), // Required
				Items: []cfTypes.Origin{
					{ // Required
						DomainName: aws.String("string"), // Required
						Id:         aws.String("string"), // Required
						CustomOriginConfig: &cfTypes.CustomOriginConfig{
							HTTPPort:             aws.Int64(1),                                         // Required
							HTTPSPort:            aws.Int64(1),                                         // Required
							OriginProtocolPolicy: cfEnums.OriginProtocolPolicy("OriginProtocolPolicy"), // Required
						},
						OriginPath: aws.String("string"),
						S3OriginConfig: &cfTypes.S3OriginConfig{
							OriginAccessIdentity: aws.String("string"), // Required
						},
					},
					// More values...
				},
			},
			Aliases: &cfTypes.Aliases{
				Quantity: aws.Int64(1), // Required
				Items: []string{
					"string", // Required
					// More values...
				},
			},
			CacheBehaviors: &cfTypes.CacheBehaviors{
				Quantity: aws.Int64(1), // Required
				Items: []cfTypes.CacheBehavior{
					{ // Required
						ForwardedValues: &cfTypes.ForwardedValues{ // Required
							Cookies: &cfTypes.CookiePreference{ // Required
								Forward: cfEnums.ItemSelection("ItemSelection"), // Required
								WhitelistedNames: &cfTypes.CookieNames{
									Quantity: aws.Int64(1), // Required
									Items: []string{
										"string", // Required
										// More values...
									},
								},
							},
							QueryString: aws.Bool(true), // Required
							Headers: &cfTypes.Headers{
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
						TrustedSigners: &cfTypes.TrustedSigners{ // Required
							Enabled:  aws.Bool(true), // Required
							Quantity: aws.Int64(1),   // Required
							Items: []string{
								"string", // Required
								// More values...
							},
						},
						ViewerProtocolPolicy: cfEnums.ViewerProtocolPolicy("ViewerProtocolPolicy"), // Required
						AllowedMethods: &cfTypes.AllowedMethods{
							Items: []cfEnums.Method{ // Required
								cfEnums.Method("string"), // Required
								// More values...
							},
							Quantity: aws.Int64(1), // Required
							CachedMethods: &cfTypes.CachedMethods{
								Items: []cfEnums.Method{ // Required
									cfEnums.Method("string"), // Required
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
			CustomErrorResponses: &cfTypes.CustomErrorResponses{
				Quantity: aws.Int64(1), // Required
				Items: []cfTypes.CustomErrorResponse{
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
			Logging: &cfTypes.LoggingConfig{
				Bucket:         aws.String("string"), // Required
				Enabled:        aws.Bool(true),       // Required
				IncludeCookies: aws.Bool(true),       // Required
				Prefix:         aws.String("string"), // Required
			},
			PriceClass: cfEnums.PriceClass("PriceClass"),
			Restrictions: &cfTypes.Restrictions{
				GeoRestriction: &cfTypes.GeoRestriction{ // Required
					Quantity:        aws.Int64(1),                                     // Required
					RestrictionType: cfEnums.GeoRestrictionType("GeoRestrictionType"), // Required
					Items: []string{
						"string", // Required
						// More values...
					},
				},
			},
			ViewerCertificate: &cfTypes.ViewerCertificate{
				CloudFrontDefaultCertificate: aws.Bool(true),
				IAMCertificateId:             aws.String("string"),
				MinimumProtocolVersion:       cfEnums.MinimumProtocolVersion("MinimumProtocolVersion"),
				SSLSupportMethod:             cfEnums.SSLSupportMethod("SSLSupportMethod"),
			},
		},
	}
}

func cloudfrontDeleteDistributionInput() *cfTypes.DeleteDistributionInput {
	return &cfTypes.DeleteDistributionInput{
		Id:      aws.String("string"), // Required
		IfMatch: aws.String("string"),
	}
}

func s3HeadObjectInput() *s3Types.HeadObjectInput {
	return &s3Types.HeadObjectInput{
		Bucket:    aws.String("somebucketname"),
		Key:       aws.String("keyname"),
		VersionId: aws.String("someVersion"),
		IfMatch:   aws.String("IfMatch"),
	}
}

func s3PutObjectAclInput() *s3Types.PutObjectAclInput {
	return &s3Types.PutObjectAclInput{
		Bucket: aws.String("somebucketname"),
		Key:    aws.String("keyname"),
		AccessControlPolicy: &s3Types.AccessControlPolicy{
			Grants: []s3Types.Grant{
				{
					Grantee: &s3Types.Grantee{
						DisplayName:  aws.String("someName"),
						EmailAddress: aws.String("someAddr"),
						ID:           aws.String("someID"),
						Type:         s3Enums.TypeCanonicalUser,
						URI:          aws.String("someURI"),
					},
					Permission: s3Enums.PermissionWrite,
				},
			},
			Owner: &s3Types.Owner{
				DisplayName: aws.String("howdy"),
				ID:          aws.String("someID"),
			},
		},
	}
}
