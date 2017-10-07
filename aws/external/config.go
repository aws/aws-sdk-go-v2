package external

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/credentials"
)

type RegionProvider interface {
	GetRegion() (string, error)
}
type CredentialsProvider interface {
	GetCredentialsProvider() (credentials.Provider, error)
}
type CredentialsValueProvider interface {
	GetCredentialsValue() (credentials.Value, error)
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

func NewConfigFromLoaders(loaders ...ConfigLoader) (Configs, error) {
	return Configs{}.AppendFromLoaders(loaders...)
}

func (cs Configs) AppendFromLoaders(loaders ...ConfigLoader) (Configs, error) {
	for _, fn := range loaders {
		cfg, err := fn(cs)
		if err != nil {
			return nil, err
		}

		cs = append(cs, cfg)
	}

	return cs, nil
}
func (cs Configs) ResolveAWSConfig(resolvers ...AWSConfigResolver) (aws.Config, error) {
	var cfg aws.Config

	for _, fn := range resolvers {
		if err := fn(&cfg, cs); err != nil {
			// TODO provide better error?
			return aws.Config{}, err
		}
	}

	return cfg, nil
}

var DefaultConfigLoaders = []ConfigLoader{
	LoadEnvConfig,
	LoadSharedConfig,
}

var DefaultAWSConfigResolvers = []AWSConfigResolver{
	ResolveDefaultAWSConfig,
	ResolveRegion,
	ResolveStaticCredentials,
}

func LoadDefaultAWSConfig() (aws.Config, error) {
	cfgs, err := NewConfigFromLoaders(DefaultConfigLoaders...)
	if err != nil {
		return aws.Config{}, err
	}
	return cfgs.ResolveAWSConfig(DefaultAWSConfigResolvers...)
}
