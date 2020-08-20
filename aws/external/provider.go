package external

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/processcreds"
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

// ProcessCredentialProviderOptions is an interface for retrieving a function for setting
// the processcreds.ProviderOptions.
type ProcessCredentialProviderOptions interface {
	GetProcessCredentialProviderOptions() (func(*processcreds.ProviderOptions), bool, error)
}

// WithProcessCredentialProviderOptions wraps a function and satisfies the EC2RoleCredentialProviderOptions interface
type WithProcessCredentialProviderOptions func(*processcreds.ProviderOptions)

// GetProcessCredentialProviderOptions returns the wrapped function
func (w WithProcessCredentialProviderOptions) GetProcessCredentialProviderOptions() (func(*processcreds.ProviderOptions), bool, error) {
	return w, true, nil
}

// GetProcessCredentialProviderOptions searches the slice of configs and returns the first function found
func GetProcessCredentialProviderOptions(configs Configs) (f func(*processcreds.ProviderOptions), found bool, err error) {
	for _, config := range configs {
		if p, ok := config.(ProcessCredentialProviderOptions); ok {
			f, found, err = p.GetProcessCredentialProviderOptions()
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
