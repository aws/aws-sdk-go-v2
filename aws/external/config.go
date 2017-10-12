package external

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

var _ s3.S3

type CustomCABundleFileProvider interface {
	GetCustomCABundleFile() (string, error)
}
type RegionProvider interface {
	GetRegion() (string, error)
}
type CredentialsValueProvider interface {
	GetCredentialsValue() (credentials.Value, error)
}
type CredentialsEndpointProvider interface {
	GetCredentialsEndpoint() (string, error)
}
type SharedConfigProfileProvider interface {
	GetSharedConfigProfile() (string, error)
}
type SharedConfigFilesProvider interface {
	GetSharedConfigFiles() ([]string, error)
}

type Config interface{}

type ConfigLoader func(Configs) (Config, error)

type AWSConfigResolver func(cfg *aws.Config, configs Configs) error

type Configs []Config

// AppendFromLoaders iterates over the slice of loaders passed in calling each
// loader function in order. The external config value returned by the loader
// will be added to the returned Configs slice.
//
// If a loader returns an error this method will stop iterating and return
// that error.
func (cs Configs) AppendFromLoaders(loaders []ConfigLoader) (Configs, error) {
	for _, fn := range loaders {
		cfg, err := fn(cs)
		if err != nil {
			return nil, err
		}

		cs = append(cs, cfg)
	}

	return cs, nil
}

// ResolveAWSConfig returns a AWS configuration populated with values by calling
// the resolvers slice passed in. Each resolver is called in order. Any resolver
// may overwrite the AWs Configuration value of a previous resolver.
//
// If an resolver returns an error this method will return that error, and stop
// iterating over the resolvers.
func (cs Configs) ResolveAWSConfig(resolvers []AWSConfigResolver) (aws.Config, error) {
	var cfg aws.Config

	for _, fn := range resolvers {
		if err := fn(&cfg, cs); err != nil {
			// TODO provide better error?
			return aws.Config{}, err
		}
	}

	return cfg, nil
}

// DefaultConfigLoaders are a slice of functions that will read external configuration
// sources for configuration values. These values are read by the AWSConfigResolvers
// using interfaces to extract specific information from the external configuration.
var DefaultConfigLoaders = []ConfigLoader{
	LoadEnvConfig,
	LoadSharedConfig,
}

// DefaultAWSConfigResolvers are a slice of functions that will resolve external
// configuration values into AWS configuration values.
//
// This will setup the AWS configuration's Region,
var DefaultAWSConfigResolvers = []AWSConfigResolver{
	ResolveDefaultAWSConfig,
	ResolveCustomCABundle,
	ResolveRegion,
	ResolveCredentialsValue,
	ResolveEndpointCredentials,
	ResolveAssumeRoleCredentials,
	ResolveFallbackEC2Credentials,
}

// LoadDefaultAWSConfig reads the SDK's default external configurations, and
// popultes an AWS Config with the values from the external configurations.
//
// The default configuration sources are:
// * Environment Variables
// * Shared Configuration and Shared Credentials files.
func LoadDefaultAWSConfig() (aws.Config, error) {
	var cfgs Configs

	cfgs, err := cfgs.AppendFromLoaders(DefaultConfigLoaders)
	if err != nil {
		return aws.Config{}, err
	}

	return cfgs.ResolveAWSConfig(DefaultAWSConfigResolvers)
}
