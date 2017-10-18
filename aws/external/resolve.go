package external

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/aws/endpointcreds"
)

// ResolveDefaultAWSConfig will write default configuration values into the cfg
// value. It will write the default values, overwriting any previous value.
//
// This should be used as the first resolver in the slice of resolvers when
// resolving external configuration.
func ResolveDefaultAWSConfig(cfg *aws.Config, configs Configs) error {
	*cfg = defaults.Config()
	return nil
}

// ResolveCustomCABundle extracts the first instance of a custom CA bundle filename
// from the external configurations. It will update the HTTP Client's builder
// to be configured with the custom CA bundle.
//
// Config provider used:
// * CustomCABundleFileProvider
func ResolveCustomCABundle(cfg *aws.Config, configs Configs) error {
	v, found, err := GetCustomCABundleFile(configs)
	if err != nil {
		// TODO error handling, What is the best way to handle this?
		// capture previous errors continue. error out if all errors
		return err
	}
	if !found {
		return nil
	}

	// TODO need to suport custom CA bundle. Adding it to the TLs cert pool.
	return fmt.Errorf("ResolveCustomeCABundle pending HTTP Client builder, %s", v)
}

// ResolveRegion extracts the first instance of a Region from the Configs slice.
//
// Config providers used:
// * RegionProvider
func ResolveRegion(cfg *aws.Config, configs Configs) error {
	v, found, err := GetRegion(configs)
	if err != nil {
		// TODO error handling, What is the best way to handle this?
		// capture previous errors continue. error out if all errors
		return err
	}
	if !found {
		return nil
	}

	cfg.Region = aws.String(v)
	return nil
}

// ResolveCredentialsValue extracts the first instance of Credentials from the
// config slices.
//
// Config providers used:
// * CredentialsValueProvider
func ResolveCredentialsValue(cfg *aws.Config, configs Configs) error {
	v, found, err := GetCredentialsValue(configs)
	if err != nil {
		// TODO error handling, What is the best way to handle this?
		// capture previous errors continue. error out if all errors
		return err
	}
	if !found {
		return nil
	}

	provider := aws.StaticProvider{Value: v}
	cfg.Credentials = aws.NewCredentials(provider)

	return nil
}

// ResolveEndpointCredentials will extra the credentials endpoint from the config
// slice. Using the endpoint, provided, to create a endpoint credential provider.
//
// Config providers used:
// * CredentialsEndpointProvider
func ResolveEndpointCredentials(cfg *aws.Config, configs Configs) error {
	v, found, err := GetCredentialsEndpoint(configs)
	if err != nil {
		// TODO error handling, What is the best way to handle this?
		// capture previous errors continue. error out if all errors
		return err
	}
	if !found {
		return nil
	}

	// TODO validate endpoint URL (localhost, 127/8, ect)
	cfgCp := cfg.Copy()
	cfgCp.EndpointResolver = aws.ResolveStaticEndpointURL(v)

	provider := endpointcreds.New(*cfgCp)
	provider.ExpiryWindow = 5 * time.Minute

	cfg.Credentials = aws.NewCredentials(provider)

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
	cfg.Credentials = aws.NewCredentials(provider)

	return nil
}
