package external

import "github.com/aws/aws-sdk-go-v2/aws"

// CustomCABundleFileProvider provides access to the custom CA bundle external
// configuration value.
type CustomCABundleFileProvider interface {
	GetCustomCABundleFile() (string, error)
}

// RegionProvider provides access to the region external configuration value.
type RegionProvider interface {
	GetRegion() (string, error)
}

// CredentialsValueProvider provides access to the credentials external
// configuration value.
type CredentialsValueProvider interface {
	GetCredentialsValue() (aws.Value, error)
}

// CredentialsEndpointProvider provides access to the credentials endpoint
// external configuration value.
type CredentialsEndpointProvider interface {
	GetCredentialsEndpoint() (string, error)
}

// SharedConfigProfileProvider provides access to the shared config profile
// name external configuration value.
type SharedConfigProfileProvider interface {
	GetSharedConfigProfile() (string, error)
}

// SharedConfigFilesProvider provides access to the shared config filesnames
// external configuration value.
type SharedConfigFilesProvider interface {
	GetSharedConfigFiles() ([]string, error)
}
