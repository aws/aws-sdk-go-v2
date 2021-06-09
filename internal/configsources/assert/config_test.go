package assert

import (
	"github.com/aws/aws-sdk-go-v2/config"
	internalConfig "github.com/aws/aws-sdk-go-v2/internal/configsources"
)

// EnableEndpointDiscoveryProvider Assertions
var (
	_ internalConfig.EnableEndpointDiscoveryProvider = &config.EnvConfig{}
	_ internalConfig.EnableEndpointDiscoveryProvider = &config.SharedConfig{}
)
