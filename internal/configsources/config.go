package configsources

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// EnableEndpointDiscoveryProvider is an interface for retrieving external configuration value
// for Enable Endpoint Discovery
type EnableEndpointDiscoveryProvider interface {
	GetEnableEndpointDiscovery(ctx context.Context) (value aws.EndpointDiscoveryEnableState, found bool, err error)
}

// ResolveEnableEndpointDiscovery extracts the first instance of a EnableEndpointDiscoveryProvider from the config slice.
// Additionally returns a aws.EndpointDiscoveryEnableState to indicate if the value was found in provided configs,
// and error if one is encountered.
func ResolveEnableEndpointDiscovery(ctx context.Context, configs []interface{}) (value aws.EndpointDiscoveryEnableState, found bool, err error) {
	for _, cfg := range configs {
		if p, ok := cfg.(EnableEndpointDiscoveryProvider); ok {
			value, found, err = p.GetEnableEndpointDiscovery(ctx)
			if err != nil || found {
				break
			}
		}
	}
	return
}
