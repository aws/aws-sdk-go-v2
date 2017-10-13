package external

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/credentials"
	"github.com/aws/aws-sdk-go-v2/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/aws/credentials/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
)

// ResolveDefaultAWSConfig will write default configuration values into the cfg
// value. It will write the default values, overwriting any previous value.
//
// This should be used as the first resolver in the slice of resolvers when
// resolving external configuration.
func ResolveDefaultAWSConfig(cfg *aws.Config, configs Configs) error {
	cfg.EndpointResolver = endpoints.DefaultResolver()
	cfg.HTTPClient = &http.Client{} // TODO replace with a Sender not HTTP specific
	cfg.Logger = aws.NewDefaultLogger()
	cfg.Retryer = nil // TODO need expose defaulte retrier
	//	TODO cfg.Handlers = defaults.Handlers()
	return nil
}

// ResolveCustomCABundle extracts the first instance of a custom CA bundle filename
// from the external configurations. It will update the HTTP Client's builder
// to be configured with the custom CA bundle.
//
// Config provider used:
// * CustomCABundleFileProvider
func ResolveCustomCABundle(cfg *aws.Config, configs Configs) error {
	for _, extCfg := range configs {
		if p, ok := extCfg.(CustomCABundleFileProvider); ok {
			if v, err := p.GetCustomCABundleFile(); err == nil && len(v) > 0 {
				// TODO need to suport custom CA bundle. Adding it to the
				// Sender's TLs cert pool
				return fmt.Errorf("ResolveCustomeCABundle not implemented")
			}
			// TODO error handling, What is the best way to handle this?
			// capture previous errors continue. error out if all errors
		}
	}

	return nil
}

// ResolveRegion extracts the first instance of a Region from the Configs slice.
//
// Config providers used:
// * RegionProvider
func ResolveRegion(cfg *aws.Config, configs Configs) error {
	for _, extCfg := range configs {
		if p, ok := extCfg.(RegionProvider); ok {
			if v, err := p.GetRegion(); err == nil && len(v) > 0 {
				cfg.Region = aws.String(v)
				break
			}
			// TODO error handling, What is the best way to handle this?
			// capture previous errors continue. error out if all errors
		}
	}

	return nil
}

// ResolveCredentialsValue extracts the first instance of Credentials from the
// config slices.
//
// Config providers used:
// * CredentialsValueProvider
func ResolveCredentialsValue(cfg *aws.Config, configs Configs) error {
	for _, extCfg := range configs {
		if p, ok := extCfg.(CredentialsValueProvider); ok {
			if v, err := p.GetCredentialsValue(); err == nil && v.Valid() {
				provider := credentials.StaticProvider{
					Value: v,
				}
				cfg.Credentials = credentials.NewCredentials(provider)
				break
			}
			// TODO error handling, What is the best way to handle this?
			// capture previous errors continue. error out if all errors
		}
	}

	return nil
}

// ResolveEndpointCredentials will extra the credentials endpoint from the config
// slice. Using the endpoint, provided, to create a endpoint credential provider.
//
// Config providers used:
// * CredentialsEndpointProvider
func ResolveEndpointCredentials(cfg *aws.Config, configs Configs) error {
	for _, extCfg := range configs {
		if p, ok := extCfg.(CredentialsEndpointProvider); ok {
			if v, err := p.GetCredentialsEndpoint(); err == nil && len(v) > 0 {
				provider := &endpointcreds.Provider{
					//					AWSConfig:    *cfg,
					//					Endpoint:     v,
					ExpiryWindow: 5 * time.Minute,
				}
				cfg.Credentials = credentials.NewCredentials(provider)

				break
			}
			// TODO error handling, What is the best way to handle this?
			// capture previous errors continue. error out if all errors
		}
	}

	return nil
}

// ResolveAssumeRoleCredentials extracts the assume role configuration from
// the external configurations.
//
// Config providers used:
func ResolveAssumeRoleCredentials(cfg *aws.Config, configs Configs) error {
	// TODO need implemented with assume role information from SharedConfig
	return nil
}

// ResolveFallbackEC2Credentials will configure the AWS config credentials to
// use EC2 Instance Role if the config's Credentials field is not already set.
func ResolveFallbackEC2Credentials(cfg *aws.Config, configs Configs) error {
	if cfg.Credentials != nil {
		return nil
	}

	provider := &ec2rolecreds.EC2RoleProvider{
		//		AWSConfig:    *cfg,
		ExpiryWindow: 5 * time.Minute,
	}
	cfg.Credentials = credentials.NewCredentials(provider)

	return nil
}
