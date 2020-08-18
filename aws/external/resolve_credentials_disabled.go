// +build disabled

package external

import "fmt"

func resolveHTTPCredProvider(cfg *aws.Config, url, authToken string, configs Configs) error {
	cfgCopy := cfg.Copy()

	cfgCopy.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{URL: url}, nil
	})

	opts := []func(*endpointcreds.ProviderOptions){
		func(options *endpointcreds.ProviderOptions) {
			options.ExpiryWindow = 5 * time.Minute
			if len(authToken) != 0 {
				options.AuthorizationToken = authToken
			}
		},
	}

	options, found, err := GetEndpointCredentialProviderOptions(configs)
	if err != nil {
		return err
	}
	if found {
		opts = append(opts, options)
	}

	provider := endpointcreds.New(cfgCopy, opts...)

	cfg.Credentials = provider

	return nil
}

func resolveEC2RoleCredentials(cfg *aws.Config, configs Configs) error {
	cfgCpy := *cfg

	opts := []func(*ec2rolecreds.ProviderOptions){
		func(options *ec2rolecreds.ProviderOptions) {
			options.ExpiryWindow = 5 * time.Minute
		},
	}

	options, found, err := GetEC2RoleCredentialProviderOptions(configs)
	if err != nil {
		return err
	}
	if found {
		opts = append(opts, options)
	}

	provider := ec2rolecreds.New(ec2metadata.New(cfgCpy), opts...)

	cfg.Credentials = provider

	return nil
}

func assumeWebIdentity(cfg *aws.Config, filepath string, roleARN, sessionName string, configs Configs) error {
	if len(filepath) == 0 {
		return awserr.New(stscreds.ErrCodeWebIdentity, "token file path is not set", nil)
	}

	if len(roleARN) == 0 {
		return awserr.New(stscreds.ErrCodeWebIdentity, "role ARN is not set", nil)
	}

	var opts []func(*stscreds.WebIdentityRoleProviderOptions)

	options, found, err := GetWebIdentityCredentialProviderOptions(configs)
	if err != nil {
		return err
	}
	if found {
		opts = append(opts, options)
	}

	provider := stscreds.NewWebIdentityRoleProvider(sts.New(*cfg), roleARN, sessionName, stscreds.IdentityTokenFile(filepath), opts...)

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
			return fmt.Errorf("")
		}
	}

	sts := sts.New(*cfg)

	opts := []func(*stscreds.AssumeRoleProviderOptions){
		func(options *stscreds.AssumeRoleProviderOptions) {
			options.RoleSessionName = sharedCfg.RoleSessionName
			if sharedCfg.RoleDurationSeconds != nil {
				if *sharedCfg.RoleDurationSeconds/time.Minute > 15 {
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
	}

	options, found, err := GetAssumeRoleCredentialProviderOptions(configs)
	if err != nil {
		return err
	}
	if found {
		opts = append(opts, options)
	}

	cfg.Credentials = stscreds.NewAssumeRoleProvider(sts, sharedCfg.RoleARN, opts...)

	return nil
}
