package external

import (
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/aws/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/aws/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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
	cfg.CredentialsLoader = aws.NewCredentialsLoader(provider)

	return nil
}

// ResolveEndpointCredentials will extract the credentials endpoint from the config
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

	if err := validateLocalEndpointURL(v); err != nil {
		return err
	}

	cfgCp := cfg.Copy()
	cfgCp.EndpointResolver = aws.ResolveWithEndpointURL(v)

	provider := endpointcreds.New(*cfgCp)
	provider.ExpiryWindow = 5 * time.Minute

	cfg.CredentialsLoader = aws.NewCredentialsLoader(provider)

	return nil
}

func validateLocalEndpointURL(v string) error {
	u, err := url.Parse(v)
	if err != nil {
		return err
	}

	if host := u.Hostname(); !(host == "localhost" || host == "127.0.0.1") {
		return fmt.Errorf("invalid endpoint credentials URL, %q, only localhost and 127.0.0.1 are valid", host)
	}

	return nil
}

const containerCredentialsEndpoint = "http://169.254.170.2"

// ResolveContainerEndpointPathCredentials will extract the container credentials
// endpoint from the config slice. Using the endpoint provided, to create a
// endpoint credential provider.
//
// Config providers used:
// * ContainerCredentialsEndpointPathProvider
func ResolveContainerEndpointPathCredentials(cfg *aws.Config, configs Configs) error {
	v, found, err := GetContainerCredentialsEndpointPath(configs)
	if err != nil {
		// TODO error handling, What is the best way to handle this?
		// capture previous errors continue. error out if all errors
		return err
	}
	if !found {
		return nil
	}

	cfgCp := cfg.Copy()

	v = containerCredentialsEndpoint + v
	cfgCp.EndpointResolver = aws.ResolveWithEndpointURL(v)

	provider := endpointcreds.New(*cfgCp)
	provider.ExpiryWindow = 5 * time.Minute

	cfg.CredentialsLoader = aws.NewCredentialsLoader(provider)

	return nil
}

// ResolveAssumeRoleCredentials extracts the assume role configuration from
// the external configurations.
//
// Config providers used:
func ResolveAssumeRoleCredentials(cfg *aws.Config, configs Configs) error {
	// TODO need implemented with assume role information from SharedConfig
	v, found, err := GetAssumeRoleConfig(configs)
	if err != nil {
		// TODO error handling, What is the best way to handle this?
		// capture previous errors continue. error out if all errors
		return err
	}
	if !found {
		return nil
	}

	cfgCp := cfg.Copy()
	// TODO support additional credential providers that are already set?
	cfgCp.CredentialsLoader = aws.NewCredentialsLoader(
		aws.StaticProvider{Value: v.Source.Credentials},
	)

	provider := &stscreds.AssumeRoleProvider{
		Client:          sts.New(cfgCp),
		RoleARN:         v.RoleARN,
		RoleSessionName: v.RoleSessionName,
	}
	if id := v.ExternalID; len(id) > 0 {
		provider.ExternalID = aws.String(id)
	}
	if len(v.MFASerial) > 0 {
		tp, foundTP, err := GetMFATokenFunc(configs)
		if err != nil {
			return err
		}
		if !foundTP {
			return fmt.Errorf("token provider required for AssumeRole with MFA")
		}
		provider.SerialNumber = aws.String(v.MFASerial)
		provider.TokenProvider = tp
	}

	cfg.CredentialsLoader = aws.NewCredentialsLoader(provider)

	return nil
}

// ResolveFallbackEC2Credentials will configure the AWS config credentials to
// use EC2 Instance Role if the config's Credentials field is not already set.
func ResolveFallbackEC2Credentials(cfg *aws.Config, configs Configs) error {
	if cfg.CredentialsLoader != nil {
		return nil
	}

	cfgCp := cfg.Copy()

	provider := &ec2rolecreds.Provider{
		Client:       ec2metadata.New(*cfgCp),
		ExpiryWindow: 5 * time.Minute,
	}
	cfg.CredentialsLoader = aws.NewCredentialsLoader(provider)

	return nil
}
