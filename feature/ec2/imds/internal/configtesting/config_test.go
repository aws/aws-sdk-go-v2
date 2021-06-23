package configtesting

import (
	"github.com/aws/aws-sdk-go-v2/config"
	imdsconfig "github.com/aws/aws-sdk-go-v2/feature/ec2/imds/internal/config"
)

var _ imdsconfig.EndpointModeResolver = (*config.LoadOptions)(nil)
var _ imdsconfig.EndpointModeResolver = (*config.EnvConfig)(nil)
var _ imdsconfig.EndpointModeResolver = (*config.SharedConfig)(nil)

var _ imdsconfig.EndpointResolver = (*config.LoadOptions)(nil)
var _ imdsconfig.EndpointResolver = (*config.EnvConfig)(nil)
var _ imdsconfig.EndpointResolver = (*config.SharedConfig)(nil)

var _ imdsconfig.ClientEnableStateResolver = (*config.LoadOptions)(nil)
var _ imdsconfig.ClientEnableStateResolver = (*config.EnvConfig)(nil)
