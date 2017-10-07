package external

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/credentials"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/request"
)

func DefaultEndpointResolver() endpoints.Resolver {
	// TODO move Resolver interface to aws package
	return endpoints.DefaultResolver()
}
func DefaultLogger() aws.Logger {
	// TODO move default logger setup into this package
	return aws.NewDefaultLogger()
}
func DefaultRequestRetrier() request.Retryer {
	// TODO default request retrier
	return nil
}
func DefaultHandlers() request.Handlers {
	// TODO move default handlers into this package?
	return defaults.Handlers()
}

func ResolveDefaultAWSConfig(cfg *aws.Config, configs Configs) error {
	cfg.EndpointResolver = DefaultEndpointResolver()
	cfg.Logger = DefaultLogger()
	cfg.Retryer = DefaultRequestRetrier()
	//	cfg.Handlers = DefaultHandlers()
	return nil
}

func ResolveRegion(cfg *aws.Config, configs Configs) error {
	for _, extCfg := range configs {
		if p, ok := extCfg.(RegionProvider); ok {
			if v, err := p.GetRegion(); err != nil {
				// TODO error handling, What is the best way to handle this?
				// capture previous errors continue. error out if all errors
			} else if len(v) > 0 {
				cfg.Region = aws.String(v)
				break
			}
		}
	}

	return nil
}

func ResolveStaticCredentials(cfg *aws.Config, configs Configs) error {
	for _, extCfg := range configs {
		if p, ok := extCfg.(CredentialsValueProvider); ok {
			if v, err := p.GetCredentialsValue(); err != nil {
				// TODO error handling, What is the best way to handle this?
				// capture previous errors continue. error out if all errors
			} else if len(v.AccessKeyID) > 0 && len(v.SecretAccessKey) > 0 {
				cfg.Credentials = credentials.NewStaticCredentialsFromCreds(v)
				break
			}
		}
	}

	return nil
}
