package config

import (
	"fmt"
	"strings"
)

// ClientEnableState provides an enumeration if the client is enabled,
// disabled, or default behavior.
type ClientEnableState uint

// Enumeration values for ClientEnableState
const (
	ClientDefaultEnableState ClientEnableState = iota
	ClientDisabled
	ClientEnabled
)

// EndpointMode is the EC2 IMDS Endpoint Configuration Mode
type EndpointMode uint

// SetFromString sets the EndpointMode based on the provided string value. Unknown values will default to EndpointModeUnset
func (e *EndpointMode) SetFromString(v string) error {
	switch {
	case len(v) == 0:
		*e = EndpointModeUnset
	case strings.EqualFold(v, "IPv6"):
		*e = EndpointModeIPv6
	case strings.EqualFold(v, "IPv4"):
		*e = EndpointModeIPv4
	default:
		return fmt.Errorf("unknown EC2 IMDS endpoint mode")
	}
	return nil
}

// Enumeration values for ClientEnableState
const (
	EndpointModeUnset EndpointMode = iota
	EndpointModeIPv4
	EndpointModeIPv6
)

// ClientEnableStateResolver is a config resolver interface for retrieving whether the IMDS client is disabled.
type ClientEnableStateResolver interface {
	GetEC2IMDSClientEnableState() (ClientEnableState, bool, error)
}

// EndpointModeResolver is a config resolver interface for retrieving the EndpointMode configuration.
type EndpointModeResolver interface {
	GetEC2IMDSEndpointMode() (EndpointMode, bool, error)
}

// EndpointResolver is a config resolver interface for retrieving the endpoint.
type EndpointResolver interface {
	GetEC2IMDSEndpoint() (string, bool, error)
}

// ResolveClientEnableState resolves the ClientEnableState from a list of configuration sources.
func ResolveClientEnableState(sources []interface{}) (value ClientEnableState, found bool, err error) {
	for _, source := range sources {
		if resolver, ok := source.(ClientEnableStateResolver); ok {
			value, found, err = resolver.GetEC2IMDSClientEnableState()
			if err != nil || found {
				return value, found, err
			}
		}
	}
	return value, found, err
}

// ResolveEndpointModeConfig resolves the EndpointMode from a list of configuration sources.
func ResolveEndpointModeConfig(sources []interface{}) (value EndpointMode, found bool, err error) {
	for _, source := range sources {
		if resolver, ok := source.(EndpointModeResolver); ok {
			value, found, err = resolver.GetEC2IMDSEndpointMode()
			if err != nil || found {
				return value, found, err
			}
		}
	}
	return value, found, err
}

// ResolveEndpointConfig resolves the endpoint from a list of configuration sources.
func ResolveEndpointConfig(sources []interface{}) (value string, found bool, err error) {
	for _, source := range sources {
		if resolver, ok := source.(EndpointResolver); ok {
			value, found, err = resolver.GetEC2IMDSEndpoint()
			if err != nil || found {
				return value, found, err
			}
		}
	}
	return value, found, err
}
