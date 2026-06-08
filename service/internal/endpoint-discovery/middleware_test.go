package endpointdiscovery

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func Test_Discover_Endpoint(t *testing.T) {
	cases := map[string]struct {
		region      string
		requestHost string
		address     WeightedAddress
		expectHost  string
	}{
		"endpoint overwritten by discovered value with valid dns suffix": {
			region:      "cn-north-1",
			requestHost: "initialhost.amazonaws.com.cn",
			address: WeightedAddress{
				URL: &url.URL{
					Host: "cachedep.amazonaws.com.cn",
				},
			},
			expectHost: "cachedep.amazonaws.com.cn",
		},
		"endpoint overwritten by discovered value with valid dual stack dns suffix": {
			region:      "us-east-1",
			requestHost: "initialhost.amazonaws.com",
			address: WeightedAddress{
				URL: &url.URL{
					Host: "cachedep.api.aws",
				},
			},
			expectHost: "cachedep.api.aws",
		},
		"discover endpoint omitted due to out of bound partition suffix": {
			region:      "us-east-1",
			requestHost: "initialhost.amazonaws.com",
			address: WeightedAddress{
				URL: &url.URL{
					Host: "cachedep.amazonaws.com.cn",
				},
			},
			expectHost: "initialhost.amazonaws.com",
		},
		"discover endpoint omitted due to invalid partition suffix": {
			region:      "us-east-1",
			requestHost: "initialhost.amazonaws.com",
			address: WeightedAddress{
				URL: &url.URL{
					Host: "cachedep.invalidamazonaws.com",
				},
			},
			expectHost: "initialhost.amazonaws.com",
		},
		"nil endpoint discovered and omitted": {
			region:      "us-east-1",
			requestHost: "initialhost.amazonaws.com",
			address:     WeightedAddress{},
			expectHost:  "initialhost.amazonaws.com",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			d := &DiscoverEndpoint{
				DiscoverOperation: func(ctx context.Context, region string, options ...func(*DiscoverEndpointOptions)) (WeightedAddress, error) {
					return c.address, nil
				},
				EndpointDiscoveryEnableState: aws.EndpointDiscoveryEnabled,
				Region:                       c.region,
			}
			_, _, err := d.HandleFinalize(context.TODO(), middleware.FinalizeInput{
				Request: &smithyhttp.Request{
					Request: &http.Request{
						URL: &url.URL{
							Host: c.requestHost,
						},
					},
				},
			}, middleware.FinalizeHandlerFunc(func(ctx context.Context, in middleware.FinalizeInput) (out middleware.FinalizeOutput, metadata middleware.Metadata, error error) {
				req, ok := in.Request.(*smithyhttp.Request)
				if !ok {
					t.Fatal("invalid request type")
				}
				if e, a := c.expectHost, req.URL.Host; e != a {
					t.Errorf("expected host %s, got %s", e, a)
				}
				return
			}))

			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
