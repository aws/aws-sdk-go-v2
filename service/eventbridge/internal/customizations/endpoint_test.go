package customizations_test

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"strings"
	"testing"
)

func TestPutEventsUpdateEndpoint(t *testing.T) {
	tests := map[string]struct {
		DisableHTTPS     bool
		CustomEndpoint   *aws.Endpoint
		UseDualStack     aws.DualStackEndpointState
		UseFIPS          aws.FIPSEndpointState
		EndpointId       *string
		Region           string
		WantErr          bool
		WantEndpoint     string
		WantSignedRegion string
		WantSignedName   string
		WantV4a          bool
	}{
		"standard aws endpoint": {
			Region:           "us-mock-1",
			WantEndpoint:     "https://events.us-mock-1.amazonaws.com/",
			WantSignedName:   "events",
			WantSignedRegion: "us-mock-1",
		},
		"dualstack aws endpoint": {
			Region:           "us-mock-1",
			UseDualStack:     aws.DualStackEndpointStateEnabled,
			WantEndpoint:     "https://events.us-mock-1.api.aws/",
			WantSignedName:   "events",
			WantSignedRegion: "us-mock-1",
		},
		"fips aws endpoint": {
			Region:           "us-mock-1",
			UseFIPS:          aws.FIPSEndpointStateEnabled,
			WantEndpoint:     "https://events-fips.us-mock-1.amazonaws.com/",
			WantSignedName:   "events",
			WantSignedRegion: "us-mock-1",
		},
		"dualstack & fips aws endpoint": {
			Region:           "us-mock-1",
			UseDualStack:     aws.DualStackEndpointStateEnabled,
			UseFIPS:          aws.FIPSEndpointStateEnabled,
			WantEndpoint:     "https://events-fips.us-mock-1.api.aws/",
			WantSignedName:   "events",
			WantSignedRegion: "us-mock-1",
		},
		"custom endpoint": {
			Region: "us-mock-1",
			CustomEndpoint: &aws.Endpoint{
				URL:           "https://custom.amazonaws.com",
				SigningRegion: "us-mock-1",
				Source:        aws.EndpointSourceCustom,
			},
			WantEndpoint:     "https://custom.amazonaws.com/",
			WantSignedName:   "events",
			WantSignedRegion: "us-mock-1",
		},
		"multi-region aws endpoint": {
			Region:           "us-mock-1",
			EndpointId:       aws.String("abc123.456def"),
			WantEndpoint:     "https://abc123.456def.endpoint.events.amazonaws.com/",
			WantSignedName:   "events",
			WantSignedRegion: "*",
			WantV4a:          true,
		},
		"multi-region dualstack aws endpoint": {
			Region:           "us-mock-1",
			EndpointId:       aws.String("abc123.456def"),
			UseDualStack:     aws.DualStackEndpointStateEnabled,
			WantEndpoint:     "https://abc123.456def.endpoint.events.api.aws/",
			WantSignedName:   "events",
			WantSignedRegion: "*",
			WantV4a:          true,
		},
		"multi-region fips aws endpoint": {
			Region:     "us-mock-1",
			EndpointId: aws.String("abc123.456def"),
			UseFIPS:    aws.FIPSEndpointStateEnabled,
			WantErr:    true,
		},
		"multi-region dualstack & fips aws endpoint": {
			Region:       "us-mock-1",
			EndpointId:   aws.String("abc123.456def"),
			UseDualStack: aws.DualStackEndpointStateEnabled,
			UseFIPS:      aws.FIPSEndpointStateEnabled,
			WantErr:      true,
		},
		"multi-region custom endpoint not service source": {
			Region:     "us-mock-1",
			EndpointId: aws.String("abc123.456def"),
			CustomEndpoint: &aws.Endpoint{
				URL:           "https://custom.amazonaws.com",
				SigningRegion: "us-mock-1",
				Source:        aws.EndpointSourceCustom,
			},
			WantEndpoint:     "https://custom.amazonaws.com/",
			WantSignedName:   "events",
			WantSignedRegion: "*",
			WantV4a:          true,
		},
		"multi-region custom endpoint service source": {
			Region:     "us-mock-1",
			EndpointId: aws.String("abc123.456def"),
			CustomEndpoint: &aws.Endpoint{
				URL:           "https://custom.amazonaws.com",
				SigningRegion: "us-mock-1",
				Source:        aws.EndpointSourceServiceMetadata,
			},
			WantEndpoint:     "https://abc123.456def.endpoint.events.amazonaws.com/",
			WantSignedName:   "events",
			WantSignedRegion: "*",
			WantV4a:          true,
		},
		"multi-region custom endpoint service source alt signing region for alt partition": {
			Region:     "us-mock-1",
			EndpointId: aws.String("abc123.456def"),
			CustomEndpoint: &aws.Endpoint{
				URL:           "https://custom.amazonaws.com",
				SigningRegion: "us-iso-mock-1",
				Source:        aws.EndpointSourceServiceMetadata,
			},
			WantEndpoint:     "https://abc123.456def.endpoint.events.c2s.ic.gov/",
			WantSignedName:   "events",
			WantSignedRegion: "*",
			WantV4a:          true,
		},
		"multi-region aws endpoint no ssl": {
			Region:           "us-mock-1",
			DisableHTTPS:     true,
			EndpointId:       aws.String("abc123.456def"),
			WantEndpoint:     "http://abc123.456def.endpoint.events.amazonaws.com/",
			WantSignedName:   "events",
			WantSignedRegion: "*",
			WantV4a:          true,
		},
		"multi-region aws endpoint empty endpoint id": {
			Region:     "us-mock-1",
			EndpointId: aws.String(""),
			WantErr:    true,
		},
		"multi-region aws endpoint bad host label": {
			Region:     "us-mock-1",
			EndpointId: aws.String("badactor.com?foo=bar"),
			WantErr:    true,
		},
		"multi-region us-iso-mock-1": {
			Region:           "us-iso-mock-1",
			EndpointId:       aws.String("abc123.456def"),
			WantEndpoint:     "https://abc123.456def.endpoint.events.c2s.ic.gov/",
			WantSignedName:   "events",
			WantSignedRegion: "*",
			WantV4a:          true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			options := eventbridge.Options{
				Credentials: unit.StubCredentialsProvider{},
				Retryer:     aws.NopRetryer{},
				Region:      tt.Region,
				HTTPClient:  smithyhttp.NopClient{},
				EndpointOptions: eventbridge.EndpointResolverOptions{
					DisableHTTPS:         tt.DisableHTTPS,
					UseDualStackEndpoint: tt.UseDualStack,
					UseFIPSEndpoint:      tt.UseFIPS,
				},
			}

			if tt.CustomEndpoint != nil {
				options.EndpointResolver = eventbridge.EndpointResolverFunc(
					func(region string, options eventbridge.EndpointResolverOptions) (aws.Endpoint, error) {
						return *tt.CustomEndpoint, nil
					})
			}

			client := eventbridge.New(options)

			var (
				request       *smithyhttp.Request
				signingRegion string
				signingName   string
			)

			_, err := client.PutEvents(context.TODO(), &eventbridge.PutEventsInput{
				Entries: []types.PutEventsRequestEntry{{
					Detail: aws.String("{}"),
				}},
				EndpointId: tt.EndpointId,
			}, func(o *eventbridge.Options) {
				o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
					return stack.Deserialize.Add(middleware.DeserializeMiddlewareFunc("loopback", func(
						ctx context.Context, input middleware.DeserializeInput, handler middleware.DeserializeHandler,
					) (
						out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
					) {
						request = input.Request.(*smithyhttp.Request)
						signingName = awsmiddleware.GetSigningName(ctx)
						signingRegion = awsmiddleware.GetSigningRegion(ctx)
						out.Result = &eventbridge.PutEventsOutput{}
						return out, metadata, nil
					}), middleware.Before)
				})
			})

			if e, a := tt.WantErr, err != nil; e != a {
				t.Fatalf("WantErr(%v) got %v", e, err)
			}

			if tt.WantErr {
				return
			}

			req := request.Build(context.Background())
			if e, a := tt.WantEndpoint, req.URL.String(); e != a {
				t.Errorf("expect url %s, got %s", e, a)
			}

			if e, a := tt.WantSignedRegion, signingRegion; e != a {
				t.Errorf("expect %s, got %s", e, a)
			}

			if e, a := tt.WantSignedName, signingName; e != a {
				t.Errorf("expect %s, got %s", e, a)
			}

			authValue := strings.SplitN(req.Header.Get("Authorization"), " ", 2)

			wantAuth := "AWS4-HMAC-SHA256"
			if tt.WantV4a {
				wantAuth = "AWS4-ECDSA-P256-SHA256"
			}

			if e, a := wantAuth, authValue[0]; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
