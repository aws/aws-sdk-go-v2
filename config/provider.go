package config

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/credentials/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/credentials/processcreds"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/ec2imds"
	"github.com/awslabs/smithy-go/middleware"
)

// SharedConfigProfileProvider provides access to the shared config profile
// name external configuration value.
type SharedConfigProfileProvider interface {
	GetSharedConfigProfile() (string, error)
}

// WithSharedConfigProfile wraps a strings to satisfy the SharedConfigProfileProvider
// interface so a slice of custom shared config files ared used when loading the
// SharedConfig.
type WithSharedConfigProfile string

// GetSharedConfigProfile returns the shared config profile.
func (c WithSharedConfigProfile) GetSharedConfigProfile() (string, error) {
	return string(c), nil
}

// GetSharedConfigProfile searches the Configs for a SharedConfigProfileProvider
// and returns the value if found. Returns an error if a provider fails before a
// value is found.
func GetSharedConfigProfile(configs Configs) (string, bool, error) {
	for _, cfg := range configs {
		if p, ok := cfg.(SharedConfigProfileProvider); ok {
			v, err := p.GetSharedConfigProfile()
			if err != nil {
				return "", false, err
			}
			if len(v) > 0 {
				return v, true, nil
			}
		}
	}

	return "", false, nil
}

// SharedConfigFilesProvider provides access to the shared config filesnames
// external configuration value.
type SharedConfigFilesProvider interface {
	GetSharedConfigFiles() ([]string, error)
}

// WithSharedConfigFiles wraps a slice of strings to satisfy the
// SharedConfigFilesProvider interface so a slice of custom shared config files
// ared used when loading the SharedConfig.
type WithSharedConfigFiles []string

// GetSharedConfigFiles returns the slice of shared config files.
func (c WithSharedConfigFiles) GetSharedConfigFiles() ([]string, error) {
	return []string(c), nil
}

// GetSharedConfigFiles searchds the Configs for a SharedConfigFilesProvider
// and returns the value if found. Returns an error if a provider fails before a
// value is found.
func GetSharedConfigFiles(configs Configs) ([]string, bool, error) {
	for _, cfg := range configs {
		if p, ok := cfg.(SharedConfigFilesProvider); ok {
			v, err := p.GetSharedConfigFiles()
			if err != nil {
				return nil, false, err
			}
			if len(v) > 0 {
				return v, true, nil
			}
		}
	}

	return nil, false, nil
}

// CustomCABundleProvider provides access to the custom CA bundle PEM bytes.
type CustomCABundleProvider interface {
	GetCustomCABundle() ([]byte, error)
}

// WithCustomCABundle provides wrapping of a region string to satisfy the
// CustomCABundleProvider interface.
type WithCustomCABundle []byte

// GetCustomCABundle returns the CA bundle PEM bytes.
func (v WithCustomCABundle) GetCustomCABundle() ([]byte, error) {
	return []byte(v), nil
}

// GetCustomCABundle searchds the Configs for a CustomCABundleProvider
// and returns the value if found. Returns an error if a provider fails before a
// value is found.
func GetCustomCABundle(configs Configs) ([]byte, bool, error) {
	for _, cfg := range configs {
		if p, ok := cfg.(CustomCABundleProvider); ok {
			v, err := p.GetCustomCABundle()
			if err != nil {
				return nil, false, err
			}
			if len(v) > 0 {
				return v, true, nil
			}
		}
	}

	return nil, false, nil
}

// RegionProvider provides access to the region external configuration value.
type RegionProvider interface {
	GetRegion() (string, error)
}

// WithRegion provides wrapping of a region string to satisfy the RegionProvider
// interface.
type WithRegion string

// GetRegion returns the region string.
func (v WithRegion) GetRegion() (string, error) {
	return string(v), nil
}

// GetRegion searchds the Configs for a RegionProvider and returns the value
// if found. Returns an error if a provider fails before a value is found.
func GetRegion(configs Configs) (string, bool, error) {
	for _, cfg := range configs {
		if p, ok := cfg.(RegionProvider); ok {
			v, err := p.GetRegion()
			if err != nil {
				return "", false, err
			}
			if len(v) > 0 {
				return v, true, nil
			}
		}
	}

	return "", false, nil
}

// WithEC2IMDSRegion provides a RegionProvider that retrieves the region
// from the EC2 Metadata service.
//
// TODO should this provider be added to the default config loading?
type WithEC2IMDSRegion struct {
	// If unset will be defaulted to Background context
	Context context.Context

	// If unset will default to generic EC2 IMDS client.
	Client *ec2imds.Client
}

// GetRegion attempts to retrieve the region from EC2 Metadata service.
func (p WithEC2IMDSRegion) GetRegion() (string, error) {
	ctx := p.Context
	if ctx == nil {
		ctx = context.Background()
	}

	client := p.Client
	if client == nil {
		client = ec2imds.New(ec2imds.Options{})
	}

	result, err := p.Client.GetRegion(ctx, nil)
	if err != nil {
		return "", err
	}

	return result.Region, nil
}

// CredentialsProviderProvider provides access to the credentials external
// configuration value.
type CredentialsProviderProvider interface {
	GetCredentialsProvider() (aws.CredentialsProvider, bool, error)
}

// WithCredentialsProvider provides wrapping of a credentials Value to satisfy the
// CredentialsProviderProvider interface.
type WithCredentialsProvider struct {
	aws.CredentialsProvider
}

// GetCredentialsProvider returns the credentials value.
func (v WithCredentialsProvider) GetCredentialsProvider() (aws.CredentialsProvider, bool, error) {
	if v.CredentialsProvider == nil {
		return nil, false, nil
	}

	return v.CredentialsProvider, true, nil
}

// GetCredentialsProvider searches the Configs for a CredentialsProviderProvider
// and returns the value if found. Returns an error if a provider fails before a
// value is found.
func GetCredentialsProvider(configs Configs) (p aws.CredentialsProvider, found bool, err error) {
	for _, cfg := range configs {
		if provider, ok := cfg.(CredentialsProviderProvider); ok {
			p, found, err = provider.GetCredentialsProvider()
			if err != nil {
				return nil, false, err
			}
			if found {
				break
			}
		}
	}

	return p, found, err
}

// MFATokenFuncProvider provides access to the MFA token function needed for
// Assume Role with MFA.
type MFATokenFuncProvider interface {
	GetMFATokenFunc() (func() (string, error), error)
}

// WithMFATokenFunc provides wrapping of a string to satisfy the
// MFATokenFuncProvider interface.
type WithMFATokenFunc func() (string, error)

// GetMFATokenFunc returns the MFA Token function.
func (p WithMFATokenFunc) GetMFATokenFunc() (func() (string, error), error) {
	return p, nil
}

// GetMFATokenFunc searches the Configs for a MFATokenFuncProvider
// and returns the value if found. Returns an error if a provider fails before a
// value is found.
func GetMFATokenFunc(configs Configs) (func() (string, error), bool, error) {
	for _, cfg := range configs {
		if p, ok := cfg.(MFATokenFuncProvider); ok {
			v, err := p.GetMFATokenFunc()
			if err != nil {
				return nil, false, err
			}
			if v != nil {
				return v, true, nil
			}
		}
	}

	return nil, false, nil
}

// WithAssumeRoleDuration provides a wrapping type of a time.Duration to satisfy
type WithAssumeRoleDuration time.Duration

// GetAssumeRoleDuration returns the wrapped time.Duration value to use when setting
// the assume role credentials duration.
func (w WithAssumeRoleDuration) GetAssumeRoleDuration() (time.Duration, bool, error) {
	return time.Duration(w), true, nil
}

// ProcessCredentialOptions is an interface for retrieving a function for setting
// the processcreds.Options.
type ProcessCredentialOptions interface {
	GetProcessCredentialOptions() (func(*processcreds.Options), bool, error)
}

// WithProcessCredentialOptions wraps a function and satisfies the
// ProcessCredentialOptions interface
type WithProcessCredentialOptions func(*processcreds.Options)

// GetProcessCredentialOptions returns the wrapped function
func (w WithProcessCredentialOptions) GetProcessCredentialOptions() (func(*processcreds.Options), bool, error) {
	return w, true, nil
}

// GetProcessCredentialOptions searches the slice of configs and returns the first function found
func GetProcessCredentialOptions(configs Configs) (f func(*processcreds.Options), found bool, err error) {
	for _, config := range configs {
		if p, ok := config.(ProcessCredentialOptions); ok {
			f, found, err = p.GetProcessCredentialOptions()
			if err != nil {
				return nil, false, err
			}
			if found {
				break
			}
		}
	}
	return f, found, err
}

// EC2RoleCredentialProviderOptions is an interface for retrieving a function
// for setting the ec2rolecreds.Provider options.
type EC2RoleCredentialProviderOptions interface {
	GetEC2RoleCredentialProviderOptions() (func(*ec2rolecreds.Options), bool, error)
}

// WithEC2RoleCredentialProviderOptions wraps a function and satisfies the
// EC2RoleCredentialProviderOptions interface
type WithEC2RoleCredentialProviderOptions func(*ec2rolecreds.Options)

// GetEC2RoleCredentialProviderOptions returns the wrapped function
func (w WithEC2RoleCredentialProviderOptions) GetEC2RoleCredentialProviderOptions() (func(*ec2rolecreds.Options), bool, error) {
	return w, true, nil
}

// GetEC2RoleCredentialProviderOptions searches the slice of configs and returns the first function found
func GetEC2RoleCredentialProviderOptions(configs Configs) (f func(*ec2rolecreds.Options), found bool, err error) {
	for _, config := range configs {
		if p, ok := config.(EC2RoleCredentialProviderOptions); ok {
			f, found, err = p.GetEC2RoleCredentialProviderOptions()
			if err != nil {
				return nil, false, err
			}
			if found {
				break
			}
		}
	}
	return f, found, err
}

// DefaultRegionProvider is an interface for retrieving a default region if a region was not resolved from other sources
type DefaultRegionProvider interface {
	GetDefaultRegion() (string, bool, error)
}

// WithDefaultRegion wraps a string and satisfies the DefaultRegionProvider interface
type WithDefaultRegion string

// GetDefaultRegion returns wrapped fallback region
func (w WithDefaultRegion) GetDefaultRegion() (string, bool, error) {
	return string(w), true, nil
}

// GetDefaultRegion searches the slice of configs and returns the first fallback region found
func GetDefaultRegion(configs Configs) (value string, found bool, err error) {
	for _, config := range configs {
		if p, ok := config.(DefaultRegionProvider); ok {
			value, found, err = p.GetDefaultRegion()
			if err != nil {
				return "", false, err
			}
			if found {
				break
			}
		}
	}

	return value, found, err
}

// EndpointCredentialProviderOptions is an interface for retrieving a function for setting
// the endpointcreds.ProviderOptions.
type EndpointCredentialProviderOptions interface {
	GetEndpointCredentialProviderOptions() (func(*endpointcreds.Options), bool, error)
}

// WithEndpointCredentialProviderOptions wraps a function and satisfies the EC2RoleCredentialProviderOptions interface
type WithEndpointCredentialProviderOptions func(*endpointcreds.Options)

// GetEndpointCredentialProviderOptions returns the wrapped function
func (w WithEndpointCredentialProviderOptions) GetEndpointCredentialProviderOptions() (func(*endpointcreds.Options), bool, error) {
	return w, true, nil
}

// GetEndpointCredentialProviderOptions searches the slice of configs and returns the first function found
func GetEndpointCredentialProviderOptions(configs Configs) (f func(*endpointcreds.Options), found bool, err error) {
	for _, config := range configs {
		if p, ok := config.(EndpointCredentialProviderOptions); ok {
			f, found, err = p.GetEndpointCredentialProviderOptions()
			if err != nil {
				return nil, false, err
			}
			if found {
				break
			}
		}
	}
	return f, found, err
}

// WebIdentityCredentialProviderOptions is an interface for retrieving a function for setting
// the stscreds.WebIdentityCredentialProviderOptions.
type WebIdentityCredentialProviderOptions interface {
	GetWebIdentityCredentialProviderOptions() (func(*stscreds.WebIdentityRoleOptions), bool, error)
}

// WithWebIdentityCredentialProviderOptions wraps a function and satisfies the EC2RoleCredentialProviderOptions interface
type WithWebIdentityCredentialProviderOptions func(*stscreds.WebIdentityRoleOptions)

// GetWebIdentityCredentialProviderOptions returns the wrapped function
func (w WithWebIdentityCredentialProviderOptions) GetWebIdentityCredentialProviderOptions() (func(*stscreds.WebIdentityRoleOptions), bool, error) {
	return w, true, nil
}

// GetWebIdentityCredentialProviderOptions searches the slice of configs and returns the first function found
func GetWebIdentityCredentialProviderOptions(configs Configs) (f func(*stscreds.WebIdentityRoleOptions), found bool, err error) {
	for _, config := range configs {
		if p, ok := config.(WebIdentityCredentialProviderOptions); ok {
			f, found, err = p.GetWebIdentityCredentialProviderOptions()
			if err != nil {
				return nil, false, err
			}
			if found {
				break
			}
		}
	}
	return f, found, err
}

// AssumeRoleCredentialProviderOptions is an interface for retrieving a function for setting
// the stscreds.AssumeRoleProviderOptions.
type AssumeRoleCredentialProviderOptions interface {
	GetAssumeRoleCredentialProviderOptions() (func(*stscreds.AssumeRoleOptions), bool, error)
}

// WithAssumeRoleCredentialProviderOptions wraps a function and satisfies the EC2RoleCredentialProviderOptions interface
type WithAssumeRoleCredentialProviderOptions func(*stscreds.AssumeRoleOptions)

// GetAssumeRoleCredentialProviderOptions returns the wrapped function
func (w WithAssumeRoleCredentialProviderOptions) GetAssumeRoleCredentialProviderOptions() (func(*stscreds.AssumeRoleOptions), bool, error) {
	return w, true, nil
}

// GetAssumeRoleCredentialProviderOptions searches the slice of configs and returns the first function found
func GetAssumeRoleCredentialProviderOptions(configs Configs) (f func(*stscreds.AssumeRoleOptions), found bool, err error) {
	for _, config := range configs {
		if p, ok := config.(AssumeRoleCredentialProviderOptions); ok {
			f, found, err = p.GetAssumeRoleCredentialProviderOptions()
			if err != nil {
				return nil, false, err
			}
			if found {
				break
			}
		}
	}
	return f, found, err
}

// HTTPClient is an HTTP client implementation
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// HTTPClientProvider is an interface for retrieving an HTTPClient.
type HTTPClientProvider interface {
	GetHTTPClient() (HTTPClient, bool, error)
}

// WithHTTPClient wraps a HTTPClient and satisfies the HTTPClientProvider interface
type WithHTTPClient struct {
	HTTPClient
}

// GetHTTPClient returns the wrapped HTTPClient. Returns an error if the wrapped client is nil.
func (w WithHTTPClient) GetHTTPClient() (HTTPClient, bool, error) {
	if w.HTTPClient == nil {
		return nil, false, fmt.Errorf("http client must not be nil")
	}
	return w.HTTPClient, true, nil
}

// GetHTTPClient searches the slice of configs and returns the first HTTPClient found.
func GetHTTPClient(configs Configs) (c HTTPClient, found bool, err error) {
	for _, config := range configs {
		if p, ok := config.(HTTPClientProvider); ok {
			c, found, err = p.GetHTTPClient()
			if err != nil {
				return nil, false, err
			}
			if found {
				break
			}
		}
	}
	return c, found, err
}

// APIOptionsProvider is an interface for retrieving APIOptions.
type APIOptionsProvider interface {
	GetAPIOptions() ([]func(*middleware.Stack) error, bool, error)
}

// WithAPIOptions wraps a slice of middlewares stack mutators and satisfies the APIOptionsProvider interface.
type WithAPIOptions []func(*middleware.Stack) error

// GetAPIOptions returns the wrapped middleware stack mutators.
func (w WithAPIOptions) GetAPIOptions() ([]func(*middleware.Stack) error, bool, error) {
	return w, true, nil
}

// GetAPIOptions searches the slice of configs and returns the first APIOptions found.
func GetAPIOptions(configs Configs) (o []func(*middleware.Stack) error, found bool, err error) {
	for _, config := range configs {
		if p, ok := config.(APIOptionsProvider); ok {
			o, found, err = p.GetAPIOptions()
			if err != nil {
				return nil, false, err
			}
			if found {
				break
			}
		}
	}
	return o, found, err
}

// EndpointResolverFuncProvider is an interface for retrieving an aws.EndpointResolver from a configuration source
type EndpointResolverFuncProvider interface {
	GetEndpointResolver() (aws.EndpointResolver, bool, error)
}

// WithEndpointResolver wraps a aws.EndpointResolver value to satisfy the EndpointResolverFuncProvider interface
type WithEndpointResolver struct {
	aws.EndpointResolver
}

// GetEndpointResolver returns the wrapped EndpointResolver
func (w WithEndpointResolver) GetEndpointResolver() (aws.EndpointResolver, bool, error) {
	return w.EndpointResolver, true, nil
}

// GetEndpointResolver searches the provided config sources for a EndpointResolverFunc that can be used
// to configure the aws.Config.EndpointResolver value.
func GetEndpointResolver(configs Configs) (f aws.EndpointResolver, found bool, err error) {
	for _, c := range configs {
		if p, ok := c.(EndpointResolverFuncProvider); ok {
			f, found, err = p.GetEndpointResolver()
			if err != nil {
				return nil, false, err
			}
		}
	}
	return f, found, err
}
