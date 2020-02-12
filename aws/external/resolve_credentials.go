package external

import (
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/aws/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/processcreds"
	"github.com/aws/aws-sdk-go-v2/aws/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const (
	// valid credential source values
	credSourceEc2Metadata  = "Ec2InstanceMetadata"
	credSourceEnvironment  = "Environment"
	credSourceECSContainer = "EcsContainer"
)

var (
	ecsContainerEndpoint = "http://169.254.170.2" // not constant to allow for swapping during unit-testing

	// WebIdentityEmptyRoleARNErr will occur if 'AWS_WEB_IDENTITY_TOKEN_FILE' was set but
	// 'AWS_ROLE_ARN' was not set.
	WebIdentityEmptyRoleARNErr = awserr.New(stscreds.ErrCodeWebIdentity, "role ARN is not set", nil)

	// WebIdentityEmptyTokenFilePathErr will occur if 'AWS_ROLE_ARN' was set but
	// 'AWS_WEB_IDENTITY_TOKEN_FILE' was not set.
	WebIdentityEmptyTokenFilePathErr = awserr.New(stscreds.ErrCodeWebIdentity, "token file path is not set", nil)
)

// ResolveCredentials TODO: Public Documentation
func ResolveCredentials(cfg *aws.Config, configs Configs) error {
	found, err := ResolveCredentialProvider(cfg, configs)
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	err = ResolveCredentialChain(cfg, configs)
	if err != nil {
		return err
	}

	return nil
}

// ResolveCredentialProvider extracts the first instance of Credentials from the
// config slices.
//
// Config providers used:
// * CredentialsProviderProvider
func ResolveCredentialProvider(cfg *aws.Config, configs Configs) (bool, error) {
	credentials, found, err := GetCredentialsProvider(configs)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	cfg.Credentials = credentials

	return true, nil
}

// ResolveCredentialChain TODO: Public Documentation
func ResolveCredentialChain(cfg *aws.Config, configs Configs) (err error) {
	envConfig, sharedConfig, other := getAWSExternalConfigs(configs)

	_, sharedProfileSet, err := GetSharedConfigProfile(other)
	if err != nil {
		return err
	}

	switch {
	case sharedProfileSet:
		err = resolveCredsFromProfile(cfg, envConfig, sharedConfig, other)
	case envConfig.Credentials.HasKeys():
		cfg.Credentials = aws.StaticCredentialsProvider{Value: envConfig.Credentials}
	case len(envConfig.WebIdentityTokenFilePath) > 0:
		err = assumeWebIdentity(cfg, envConfig.WebIdentityTokenFilePath, envConfig.RoleARN, envConfig.RoleSessionName)
	default:
		err = resolveCredsFromProfile(cfg, envConfig, sharedConfig, other)
	}

	return err
}

func resolveCredsFromProfile(cfg *aws.Config, envConfig *EnvConfig, sharedConfig *SharedConfig, configs Configs) (err error) {
	switch {
	case sharedConfig.Source != nil:
		// Assume IAM role with credentials source from a different profile.
		err = resolveCredsFromProfile(cfg, envConfig, sharedConfig.Source, configs)
	case sharedConfig.Credentials.HasKeys():
		// Static Credentials from Shared Config/Credentials file.
		cfg.Credentials = aws.StaticCredentialsProvider{
			Value: sharedConfig.Credentials,
		}
	case len(sharedConfig.CredentialProcess) != 0:
		// Get credentials from CredentialProcess
		cfg.Credentials = processcreds.NewProvider(sharedConfig.CredentialProcess)
	case len(sharedConfig.CredentialSource) != 0:
		err = resolveCredsFromSource(cfg, envConfig, sharedConfig, configs)
	case len(sharedConfig.WebIdentityTokenFile) != 0:
		// Credentials from Assume Web Identity token require an IAM Role, and
		// that roll will be assumed. May be wrapped with another assume role
		// via SourceProfile.
		err = assumeWebIdentity(cfg, sharedConfig.WebIdentityTokenFile, sharedConfig.RoleARN, sharedConfig.RoleSessionName)
	case len(envConfig.ContainerCredentialsEndpoint) != 0:
		err = resolveLocalHTTPCredProvider(cfg, envConfig.ContainerCredentialsEndpoint, envConfig.ContainerAuthorizationToken)
	case len(envConfig.ContainerCredentialsRelativePath) != 0:
		err = resolveHTTPCredProvider(cfg, ecsContainerURI(envConfig.ContainerCredentialsRelativePath), envConfig.ContainerAuthorizationToken)
	default:
		err = resolveEC2RoleCredentials(cfg)
	}
	if err != nil {
		return err
	}

	if len(sharedConfig.RoleARN) > 0 {
		return credsFromAssumeRole(cfg, sharedConfig, configs)
	}

	return nil
}

func ecsContainerURI(path string) string {
	return fmt.Sprintf("%s%s", ecsContainerEndpoint, path)
}

func resolveLocalHTTPCredProvider(cfg *aws.Config, endpointURL, authToken string) error {
	var errMsg string

	parsed, err := url.Parse(endpointURL)
	if err != nil {
		errMsg = fmt.Sprintf("invalid URL, %v", err)
	} else {
		host := parsed.Hostname()
		if len(host) == 0 {
			errMsg = "unable to parse host from local HTTP cred provider URL"
		} else if isLoopback, loopbackErr := isLoopbackHost(host); loopbackErr != nil {
			errMsg = fmt.Sprintf("failed to resolve host %q, %v", host, loopbackErr)
		} else if !isLoopback {
			errMsg = fmt.Sprintf("invalid endpoint host, %q, only loopback hosts are allowed.", host)
		}
	}

	if len(errMsg) > 0 {
		if cfg.Logger != nil {
			cfg.Logger.Log("Ignoring, HTTP credential provider", errMsg, err)
		}
		return awserr.New("CredentialsEndpointError", errMsg, err)
	}

	return resolveHTTPCredProvider(cfg, endpointURL, authToken)
}

func resolveHTTPCredProvider(cfg *aws.Config, url, authToken string) error {
	cfgCopy := cfg.Copy()

	cfgCopy.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{URL: url}, nil
	})

	provider := endpointcreds.New(cfgCopy, func(options *endpointcreds.ProviderOptions) {
		options.ExpiryWindow = 5 * time.Minute
		if len(authToken) != 0 {
			options.AuthorizationToken = authToken
		}
	})

	cfg.Credentials = provider

	return nil
}

func resolveEC2RoleCredentials(cfg *aws.Config) error {
	resolver := cfg.EndpointResolver
	if resolver == nil {
		resolver = endpoints.NewDefaultResolver()
	}

	cfgCpy := *cfg

	cfgCpy.EndpointResolver = resolver

	provider := ec2rolecreds.New(ec2metadata.New(cfgCpy), func(options *ec2rolecreds.ProviderOptions) {
		options.ExpiryWindow = 5 * time.Minute
	})

	cfg.Credentials = provider

	return nil
}

func resolveCredsFromSource(cfg *aws.Config, envConfig *EnvConfig, sharedCfg *SharedConfig, configs Configs) (err error) {
	switch sharedCfg.CredentialSource {
	case credSourceEc2Metadata:
		err = resolveEC2RoleCredentials(cfg)
	case credSourceEnvironment:
		cfg.Credentials = aws.StaticCredentialsProvider{Value: envConfig.Credentials}
	case credSourceECSContainer:
		if len(envConfig.ContainerCredentialsRelativePath) == 0 {
			return ErrSharedConfigECSContainerEnvVarEmpty
		}
		return resolveHTTPCredProvider(cfg, ecsContainerURI(envConfig.ContainerCredentialsRelativePath), envConfig.ContainerAuthorizationToken)
	default:
		return ErrSharedConfigInvalidCredSource
	}

	return nil
}

func assumeWebIdentity(cfg *aws.Config, filepath string, roleARN, sessionName string) error {
	if len(filepath) == 0 {
		return WebIdentityEmptyTokenFilePathErr
	}

	if len(roleARN) == 0 {
		return WebIdentityEmptyRoleARNErr
	}

	provider := stscreds.NewWebIdentityRoleProvider(
		sts.New(*cfg),
		roleARN,
		sessionName,
		filepath,
	)

	cfg.Credentials = provider

	return nil
}

func credsFromAssumeRole(cfg *aws.Config, sharedCfg *SharedConfig, configs Configs) (err error) {
	var tokenFunc func() (string, error)
	if len(sharedCfg.MFASerial) != 0 {
		var found bool
		tokenFunc, found, err = GetMFATokenFunc(configs)
		if err != nil {
			return err
		}

		if !found {
			// AssumeRole Token provider is required if doing Assume Role
			// with MFA.
			return AssumeRoleTokenProviderNotSetError{}
		}
	}

	sts := sts.New(*cfg)

	cfg.Credentials = stscreds.NewAssumeRoleProvider(sts, sharedCfg.RoleARN,
		func(options *stscreds.AssumeRoleProviderOptions) {
			options.RoleSessionName = sharedCfg.RoleSessionName
			if sharedCfg.RoleDurationSeconds != nil {
				if *sharedCfg.RoleDurationSeconds/time.Second < 900 {
					options.Duration = time.Second * 900
				} else {
					options.Duration = *sharedCfg.RoleDurationSeconds
				}
			}

			// Assume role with external ID
			if len(sharedCfg.ExternalID) > 0 {
				options.ExternalID = aws.String(sharedCfg.ExternalID)
			}

			// Assume role with MFA
			if len(sharedCfg.MFASerial) != 0 {
				options.SerialNumber = aws.String(sharedCfg.MFASerial)
				options.TokenProvider = tokenFunc
			}
		},
	)

	return nil
}

func getAWSExternalConfigs(configs Configs) (*EnvConfig, *SharedConfig, Configs) {
	var (
		envConfig    *EnvConfig
		sharedConfig *SharedConfig
		other        Configs
	)

	for i := range configs {
		switch c := configs[i].(type) {
		case EnvConfig:
			if envConfig == nil {
				envConfig = &c
			}
		case *EnvConfig:
			if envConfig == nil {
				envConfig = c
			}
		case SharedConfig:
			if sharedConfig == nil {
				sharedConfig = &c
			}
		case *SharedConfig:
			if envConfig == nil {
				sharedConfig = c
			}
		default:
			other = append(other, c)
		}
	}

	if envConfig == nil {
		envConfig = &EnvConfig{}
	}

	if sharedConfig == nil {
		sharedConfig = &SharedConfig{}
	}

	return envConfig, sharedConfig, other
}

// AssumeRoleTokenProviderNotSetError is an error returned when creating a
// session when the MFAToken option is not set when shared config is configured
// load assume a role with an MFA token.
type AssumeRoleTokenProviderNotSetError struct{}

// Code is the short id of the error.
func (e AssumeRoleTokenProviderNotSetError) Code() string {
	return "AssumeRoleTokenProviderNotSetError"
}

// Message is the description of the error
func (e AssumeRoleTokenProviderNotSetError) Message() string {
	return fmt.Sprintf("assume role with MFA enabled, but AssumeRoleTokenProvider session option not set.")
}

// OrigErr is the underlying error that caused the failure.
func (e AssumeRoleTokenProviderNotSetError) OrigErr() error {
	return nil
}

// Error satisfies the error interface.
func (e AssumeRoleTokenProviderNotSetError) Error() string {
	return awserr.SprintError(e.Code(), e.Message(), "", nil)
}
