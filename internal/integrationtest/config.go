package integrationtest

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// LoadConfigWithDefaultRegion loads the default configuration for the SDK, and
// falls back to a default region if one is not specified.
func LoadConfigWithDefaultRegion(defaultRegion string) (cfg aws.Config, err error) {
	cfg, err = config.LoadDefaultConfig()
	if err != nil {
		return cfg, err
	}

	if len(cfg.Region) == 0 {
		cfg.Region = defaultRegion
	}

	cfg.APIOptions = append(cfg.APIOptions,
		RemoveOperationInputValidationMiddleware)

	return cfg, nil
}
