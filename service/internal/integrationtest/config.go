package integrationtest

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// LoadConfigWithDefaultRegion loads the default configuration for the SDK, and
// falls back to a default region if one is not specified.
func LoadConfigWithDefaultRegion(defaultRegion string) (cfg aws.Config, err error) {
	var lm aws.ClientLogMode

	if strings.EqualFold(os.Getenv("AWS_DEBUG_REQUEST"), "true") {
		lm |= aws.LogRequest

	} else if strings.EqualFold(os.Getenv("AWS_DEBUG_REQUEST_BODY"), "true") {
		lm |= aws.LogRequestWithBody
	}

	cfg, err = config.LoadDefaultConfig(config.WithClientLogMode(lm))
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

type logger struct{}

func (logger) Logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
