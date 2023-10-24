package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
)

// ServiceBaseEndpointProvider is needed to search for all providers
// that provide a configured service endpoint
type ServiceBaseEndpointProvider interface {
	GetServiceBaseEndpoint(ctx context.Context, sdkID string) (string, bool, error)
}

// ResolveServiceBaseEndpoint is used to retrieve service endpoints from configured sources
// while allowing for configured endpoints to be disabled
func ResolveServiceBaseEndpoint(ctx context.Context, sdkID string, configs []config.Config) (value string, found bool, err error) {
	if val, found, _ := config.GetIgnoreConfiguredEndpoints(ctx, configs); found && val {
		return "", false, nil
	}

	for _, cs := range configs {
		if p, ok := cs.(ServiceBaseEndpointProvider); ok {
			value, found, err = p.GetServiceBaseEndpoint(context.Background(), sdkID)
			if err != nil || found {
				break
			}
		}
	}
	return
}
