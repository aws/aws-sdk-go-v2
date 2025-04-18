// Code generated by smithy-go-codegen DO NOT EDIT.

package endpoints

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	endpoints "github.com/aws/aws-sdk-go-v2/internal/endpoints/v2"
	"github.com/aws/smithy-go/logging"
	"regexp"
)

// Options is the endpoint resolver configuration options
type Options struct {
	// Logger is a logging implementation that log events should be sent to.
	Logger logging.Logger

	// LogDeprecated indicates that deprecated endpoints should be logged to the
	// provided logger.
	LogDeprecated bool

	// ResolvedRegion is used to override the region to be resolved, rather then the
	// using the value passed to the ResolveEndpoint method. This value is used by the
	// SDK to translate regions like fips-us-east-1 or us-east-1-fips to an alternative
	// name. You must not set this value directly in your application.
	ResolvedRegion string

	// DisableHTTPS informs the resolver to return an endpoint that does not use the
	// HTTPS scheme.
	DisableHTTPS bool

	// UseDualStackEndpoint specifies the resolver must resolve a dual-stack endpoint.
	UseDualStackEndpoint aws.DualStackEndpointState

	// UseFIPSEndpoint specifies the resolver must resolve a FIPS endpoint.
	UseFIPSEndpoint aws.FIPSEndpointState
}

func (o Options) GetResolvedRegion() string {
	return o.ResolvedRegion
}

func (o Options) GetDisableHTTPS() bool {
	return o.DisableHTTPS
}

func (o Options) GetUseDualStackEndpoint() aws.DualStackEndpointState {
	return o.UseDualStackEndpoint
}

func (o Options) GetUseFIPSEndpoint() aws.FIPSEndpointState {
	return o.UseFIPSEndpoint
}

func transformToSharedOptions(options Options) endpoints.Options {
	return endpoints.Options{
		Logger:               options.Logger,
		LogDeprecated:        options.LogDeprecated,
		ResolvedRegion:       options.ResolvedRegion,
		DisableHTTPS:         options.DisableHTTPS,
		UseDualStackEndpoint: options.UseDualStackEndpoint,
		UseFIPSEndpoint:      options.UseFIPSEndpoint,
	}
}

// Resolver Snowball endpoint resolver
type Resolver struct {
	partitions endpoints.Partitions
}

// ResolveEndpoint resolves the service endpoint for the given region and options
func (r *Resolver) ResolveEndpoint(region string, options Options) (endpoint aws.Endpoint, err error) {
	if len(region) == 0 {
		return endpoint, &aws.MissingRegionError{}
	}

	opt := transformToSharedOptions(options)
	return r.partitions.ResolveEndpoint(region, opt)
}

// New returns a new Resolver
func New() *Resolver {
	return &Resolver{
		partitions: defaultPartitions,
	}
}

var partitionRegexp = struct {
	Aws      *regexp.Regexp
	AwsCn    *regexp.Regexp
	AwsEusc  *regexp.Regexp
	AwsIso   *regexp.Regexp
	AwsIsoB  *regexp.Regexp
	AwsIsoE  *regexp.Regexp
	AwsIsoF  *regexp.Regexp
	AwsUsGov *regexp.Regexp
}{

	Aws:      regexp.MustCompile("^(us|eu|ap|sa|ca|me|af|il|mx)\\-\\w+\\-\\d+$"),
	AwsCn:    regexp.MustCompile("^cn\\-\\w+\\-\\d+$"),
	AwsEusc:  regexp.MustCompile("^eusc\\-(de)\\-\\w+\\-\\d+$"),
	AwsIso:   regexp.MustCompile("^us\\-iso\\-\\w+\\-\\d+$"),
	AwsIsoB:  regexp.MustCompile("^us\\-isob\\-\\w+\\-\\d+$"),
	AwsIsoE:  regexp.MustCompile("^eu\\-isoe\\-\\w+\\-\\d+$"),
	AwsIsoF:  regexp.MustCompile("^us\\-isof\\-\\w+\\-\\d+$"),
	AwsUsGov: regexp.MustCompile("^us\\-gov\\-\\w+\\-\\d+$"),
}

var defaultPartitions = endpoints.Partitions{
	{
		ID: "aws",
		Defaults: map[endpoints.DefaultKey]endpoints.Endpoint{
			{
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname:          "snowball.{region}.api.aws",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname:          "snowball-fips.{region}.amazonaws.com",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname:          "snowball-fips.{region}.api.aws",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: 0,
			}: {
				Hostname:          "snowball.{region}.amazonaws.com",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
		},
		RegionRegex:    partitionRegexp.Aws,
		IsRegionalized: true,
		Endpoints: endpoints.Endpoints{
			endpoints.EndpointKey{
				Region: "af-south-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "af-south-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.af-south-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "af-south-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.af-south-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "af-south-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.af-south-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "ap-east-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "ap-east-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.ap-east-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "ap-east-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.ap-east-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "ap-east-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.ap-east-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "ap-northeast-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "ap-northeast-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.ap-northeast-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "ap-northeast-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.ap-northeast-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "ap-northeast-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.ap-northeast-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "ap-northeast-2",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "ap-northeast-2",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.ap-northeast-2.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "ap-northeast-2",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.ap-northeast-2.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "ap-northeast-2",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.ap-northeast-2.api.aws",
			},
			endpoints.EndpointKey{
				Region: "ap-northeast-3",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "ap-northeast-3",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.ap-northeast-3.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "ap-northeast-3",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.ap-northeast-3.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "ap-northeast-3",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.ap-northeast-3.api.aws",
			},
			endpoints.EndpointKey{
				Region: "ap-south-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "ap-south-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.ap-south-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "ap-south-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.ap-south-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "ap-south-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.ap-south-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "ap-southeast-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "ap-southeast-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.ap-southeast-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "ap-southeast-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.ap-southeast-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "ap-southeast-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.ap-southeast-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "ap-southeast-2",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "ap-southeast-2",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.ap-southeast-2.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "ap-southeast-2",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.ap-southeast-2.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "ap-southeast-2",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.ap-southeast-2.api.aws",
			},
			endpoints.EndpointKey{
				Region: "ap-southeast-3",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "ap-southeast-3",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.ap-southeast-3.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "ap-southeast-3",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.ap-southeast-3.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "ap-southeast-3",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.ap-southeast-3.api.aws",
			},
			endpoints.EndpointKey{
				Region: "ca-central-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "ca-central-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.ca-central-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "ca-central-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.ca-central-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "ca-central-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.ca-central-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "eu-central-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "eu-central-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.eu-central-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "eu-central-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.eu-central-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "eu-central-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.eu-central-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "eu-north-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "eu-north-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.eu-north-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "eu-north-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.eu-north-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "eu-north-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.eu-north-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "eu-south-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "eu-south-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.eu-south-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "eu-south-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.eu-south-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "eu-south-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.eu-south-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "eu-west-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "eu-west-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.eu-west-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "eu-west-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.eu-west-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "eu-west-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.eu-west-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "eu-west-2",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "eu-west-2",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.eu-west-2.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "eu-west-2",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.eu-west-2.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "eu-west-2",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.eu-west-2.api.aws",
			},
			endpoints.EndpointKey{
				Region: "eu-west-3",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "eu-west-3",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.eu-west-3.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "eu-west-3",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.eu-west-3.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "eu-west-3",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.eu-west-3.api.aws",
			},
			endpoints.EndpointKey{
				Region: "fips-af-south-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.af-south-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "af-south-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-ap-east-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.ap-east-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "ap-east-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-ap-northeast-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.ap-northeast-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "ap-northeast-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-ap-northeast-2",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.ap-northeast-2.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "ap-northeast-2",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-ap-northeast-3",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.ap-northeast-3.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "ap-northeast-3",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-ap-south-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.ap-south-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "ap-south-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-ap-southeast-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.ap-southeast-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "ap-southeast-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-ap-southeast-2",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.ap-southeast-2.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "ap-southeast-2",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-ap-southeast-3",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.ap-southeast-3.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "ap-southeast-3",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-ca-central-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.ca-central-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "ca-central-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-eu-central-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.eu-central-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "eu-central-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-eu-north-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.eu-north-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "eu-north-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-eu-south-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.eu-south-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "eu-south-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-eu-west-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.eu-west-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "eu-west-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-eu-west-2",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.eu-west-2.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "eu-west-2",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-eu-west-3",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.eu-west-3.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "eu-west-3",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-il-central-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.il-central-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "il-central-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-me-central-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.me-central-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "me-central-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-sa-east-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.sa-east-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "sa-east-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-us-east-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.us-east-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-east-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-us-east-2",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.us-east-2.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-east-2",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-us-west-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.us-west-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-west-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-us-west-2",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.us-west-2.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-west-2",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "il-central-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "il-central-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.il-central-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "il-central-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.il-central-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "il-central-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.il-central-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "me-central-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "me-central-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.me-central-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "me-central-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.me-central-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "me-central-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.me-central-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "sa-east-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "sa-east-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.sa-east-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "sa-east-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.sa-east-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "sa-east-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.sa-east-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "us-east-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "us-east-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.us-east-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "us-east-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.us-east-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "us-east-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.us-east-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "us-east-2",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "us-east-2",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.us-east-2.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "us-east-2",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.us-east-2.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "us-east-2",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.us-east-2.api.aws",
			},
			endpoints.EndpointKey{
				Region: "us-west-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "us-west-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.us-west-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "us-west-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.us-west-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "us-west-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.us-west-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "us-west-2",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "us-west-2",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.us-west-2.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "us-west-2",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.us-west-2.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "us-west-2",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.us-west-2.api.aws",
			},
		},
	},
	{
		ID: "aws-cn",
		Defaults: map[endpoints.DefaultKey]endpoints.Endpoint{
			{
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname:          "snowball.{region}.api.amazonwebservices.com.cn",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname:          "snowball-fips.{region}.amazonaws.com.cn",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname:          "snowball-fips.{region}.api.amazonwebservices.com.cn",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: 0,
			}: {
				Hostname:          "snowball.{region}.amazonaws.com.cn",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
		},
		RegionRegex:    partitionRegexp.AwsCn,
		IsRegionalized: true,
		Endpoints: endpoints.Endpoints{
			endpoints.EndpointKey{
				Region: "cn-north-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "cn-north-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.cn-north-1.amazonaws.com.cn",
			},
			endpoints.EndpointKey{
				Region: "cn-northwest-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "cn-northwest-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.cn-northwest-1.amazonaws.com.cn",
			},
			endpoints.EndpointKey{
				Region: "fips-cn-north-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.cn-north-1.amazonaws.com.cn",
				CredentialScope: endpoints.CredentialScope{
					Region: "cn-north-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-cn-northwest-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.cn-northwest-1.amazonaws.com.cn",
				CredentialScope: endpoints.CredentialScope{
					Region: "cn-northwest-1",
				},
				Deprecated: aws.TrueTernary,
			},
		},
	},
	{
		ID: "aws-eusc",
		Defaults: map[endpoints.DefaultKey]endpoints.Endpoint{
			{
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname:          "snowball-fips.{region}.amazonaws.eu",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: 0,
			}: {
				Hostname:          "snowball.{region}.amazonaws.eu",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
		},
		RegionRegex:    partitionRegexp.AwsEusc,
		IsRegionalized: true,
	},
	{
		ID: "aws-iso",
		Defaults: map[endpoints.DefaultKey]endpoints.Endpoint{
			{
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname:          "snowball-fips.{region}.c2s.ic.gov",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: 0,
			}: {
				Hostname:          "snowball.{region}.c2s.ic.gov",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
		},
		RegionRegex:    partitionRegexp.AwsIso,
		IsRegionalized: true,
		Endpoints: endpoints.Endpoints{
			endpoints.EndpointKey{
				Region: "us-iso-east-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region: "us-iso-west-1",
			}: endpoints.Endpoint{},
		},
	},
	{
		ID: "aws-iso-b",
		Defaults: map[endpoints.DefaultKey]endpoints.Endpoint{
			{
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname:          "snowball-fips.{region}.sc2s.sgov.gov",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: 0,
			}: {
				Hostname:          "snowball.{region}.sc2s.sgov.gov",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
		},
		RegionRegex:    partitionRegexp.AwsIsoB,
		IsRegionalized: true,
		Endpoints: endpoints.Endpoints{
			endpoints.EndpointKey{
				Region: "us-isob-east-1",
			}: endpoints.Endpoint{},
		},
	},
	{
		ID: "aws-iso-e",
		Defaults: map[endpoints.DefaultKey]endpoints.Endpoint{
			{
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname:          "snowball-fips.{region}.cloud.adc-e.uk",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: 0,
			}: {
				Hostname:          "snowball.{region}.cloud.adc-e.uk",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
		},
		RegionRegex:    partitionRegexp.AwsIsoE,
		IsRegionalized: true,
	},
	{
		ID: "aws-iso-f",
		Defaults: map[endpoints.DefaultKey]endpoints.Endpoint{
			{
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname:          "snowball-fips.{region}.csp.hci.ic.gov",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: 0,
			}: {
				Hostname:          "snowball.{region}.csp.hci.ic.gov",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
		},
		RegionRegex:    partitionRegexp.AwsIsoF,
		IsRegionalized: true,
	},
	{
		ID: "aws-us-gov",
		Defaults: map[endpoints.DefaultKey]endpoints.Endpoint{
			{
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname:          "snowball.{region}.api.aws",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname:          "snowball-fips.{region}.amazonaws.com",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname:          "snowball-fips.{region}.api.aws",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
			{
				Variant: 0,
			}: {
				Hostname:          "snowball.{region}.amazonaws.com",
				Protocols:         []string{"https"},
				SignatureVersions: []string{"v4"},
			},
		},
		RegionRegex:    partitionRegexp.AwsUsGov,
		IsRegionalized: true,
		Endpoints: endpoints.Endpoints{
			endpoints.EndpointKey{
				Region: "fips-us-gov-east-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.us-gov-east-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-gov-east-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "fips-us-gov-west-1",
			}: endpoints.Endpoint{
				Hostname: "snowball-fips.us-gov-west-1.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-gov-west-1",
				},
				Deprecated: aws.TrueTernary,
			},
			endpoints.EndpointKey{
				Region: "us-gov-east-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "us-gov-east-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.us-gov-east-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "us-gov-east-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.us-gov-east-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "us-gov-east-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.us-gov-east-1.api.aws",
			},
			endpoints.EndpointKey{
				Region: "us-gov-west-1",
			}: endpoints.Endpoint{},
			endpoints.EndpointKey{
				Region:  "us-gov-west-1",
				Variant: endpoints.FIPSVariant,
			}: {
				Hostname: "snowball-fips.us-gov-west-1.amazonaws.com",
			},
			endpoints.EndpointKey{
				Region:  "us-gov-west-1",
				Variant: endpoints.FIPSVariant | endpoints.DualStackVariant,
			}: {
				Hostname: "snowball-fips.us-gov-west-1.api.aws",
			},
			endpoints.EndpointKey{
				Region:  "us-gov-west-1",
				Variant: endpoints.DualStackVariant,
			}: {
				Hostname: "snowball.us-gov-west-1.api.aws",
			},
		},
	},
}
