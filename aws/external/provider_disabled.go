// +build disabled

package external

// WithEC2MetadataRegion provides a RegionProvider that retrieves the region
// from the EC2 Metadata service.
//
// TODO add this provider to the default config loading?
type WithEC2MetadataRegion struct {
	ctx    context.Context
	client *ec2metadata.Client
}

// NewWithEC2MetadataRegion function takes in a context and an ec2metadataClient,
// returns a WithEC2MetadataRegion region provider
//
// Usage:
// ec2metaClient := ec2metadata.New(defaults.Config())
//
// cfg, err := external.LoadDefaultAWSConfig(
//    external.NewWithEC2MetadataRegion(ctx, ec2metaClient),
// )
//
func NewWithEC2MetadataRegion(ctx context.Context, client *ec2metadata.Client) WithEC2MetadataRegion {
	return WithEC2MetadataRegion{
		ctx:    ctx,
		client: client,
	}
}

// GetRegion attempts to retrieve the region from EC2 Metadata service.
func (p WithEC2MetadataRegion) GetRegion() (string, error) {
	return p.client.Region(p.ctx)
}

// EC2RoleCredentialProviderOptions is an interface for retrieving a function for setting
// the ec2rolecreds.Provider options.
type EC2RoleCredentialProviderOptions interface {
	GetEC2RoleCredentialProviderOptions() (func(*ec2rolecreds.ProviderOptions), bool, error)
}

// WithEC2RoleCredentialProviderOptions wraps a function and satisfies the EC2RoleCredentialProviderOptions interface
type WithEC2RoleCredentialProviderOptions func(*ec2rolecreds.ProviderOptions)

// GetEC2RoleCredentialProviderOptions returns the wrapped function
func (w WithEC2RoleCredentialProviderOptions) GetEC2RoleCredentialProviderOptions() (func(*ec2rolecreds.ProviderOptions), bool, error) {
	return w, true, nil
}

// GetEC2RoleCredentialProviderOptions searches the slice of configs and returns the first function found
func GetEC2RoleCredentialProviderOptions(configs Configs) (f func(*ec2rolecreds.ProviderOptions), found bool, err error) {
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

// EndpointCredentialProviderOptions is an interface for retrieving a function for setting
// the endpointcreds.ProviderOptions.
type EndpointCredentialProviderOptions interface {
	GetEndpointCredentialProviderOptions() (func(*endpointcreds.ProviderOptions), bool, error)
}

// WithEndpointCredentialProviderOptions wraps a function and satisfies the EC2RoleCredentialProviderOptions interface
type WithEndpointCredentialProviderOptions func(*endpointcreds.ProviderOptions)

// GetEndpointCredentialProviderOptions returns the wrapped function
func (w WithEndpointCredentialProviderOptions) GetEndpointCredentialProviderOptions() (func(*endpointcreds.ProviderOptions), bool, error) {
	return w, true, nil
}

// GetEndpointCredentialProviderOptions searches the slice of configs and returns the first function found
func GetEndpointCredentialProviderOptions(configs Configs) (f func(*endpointcreds.ProviderOptions), found bool, err error) {
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

// AssumeRoleCredentialProviderOptions is an interface for retrieving a function for setting
// the stscreds.AssumeRoleProviderOptions.
type AssumeRoleCredentialProviderOptions interface {
	GetAssumeRoleCredentialProviderOptions() (func(*stscreds.AssumeRoleProviderOptions), bool, error)
}

// WithAssumeRoleCredentialProviderOptions wraps a function and satisfies the EC2RoleCredentialProviderOptions interface
type WithAssumeRoleCredentialProviderOptions func(*stscreds.AssumeRoleProviderOptions)

// GetAssumeRoleCredentialProviderOptions returns the wrapped function
func (w WithAssumeRoleCredentialProviderOptions) GetAssumeRoleCredentialProviderOptions() (func(*stscreds.AssumeRoleProviderOptions), bool, error) {
	return w, true, nil
}

// GetAssumeRoleCredentialProviderOptions searches the slice of configs and returns the first function found
func GetAssumeRoleCredentialProviderOptions(configs Configs) (f func(*stscreds.AssumeRoleProviderOptions), found bool, err error) {
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

// WebIdentityCredentialProviderOptions is an interface for retrieving a function for setting
// the stscreds.WebIdentityCredentialProviderOptions.
type WebIdentityCredentialProviderOptions interface {
	GetWebIdentityCredentialProviderOptions() (func(*stscreds.WebIdentityRoleProviderOptions), bool, error)
}

// WithWebIdentityCredentialProviderOptions wraps a function and satisfies the EC2RoleCredentialProviderOptions interface
type WithWebIdentityCredentialProviderOptions func(*stscreds.WebIdentityRoleProviderOptions)

// GetWebIdentityCredentialProviderOptions returns the wrapped function
func (w WithWebIdentityCredentialProviderOptions) GetWebIdentityCredentialProviderOptions() (func(*stscreds.WebIdentityRoleProviderOptions), bool, error) {
	return w, true, nil
}

// GetWebIdentityCredentialProviderOptions searches the slice of configs and returns the first function found
func GetWebIdentityCredentialProviderOptions(configs Configs) (f func(*stscreds.WebIdentityRoleProviderOptions), found bool, err error) {
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
