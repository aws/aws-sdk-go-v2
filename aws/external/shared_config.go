package external

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/internal/ini"
)

const (
	// Static Credentials group
	accessKeyIDKey  = `aws_access_key_id`     // group required
	secretAccessKey = `aws_secret_access_key` // group required
	sessionTokenKey = `aws_session_token`     // optional

	// Assume Role Credentials group
	roleArnKey          = `role_arn`          // group required
	sourceProfileKey    = `source_profile`    // group required
	credentialSourceKey = `credential_source` // group required (or source_profile)
	externalIDKey       = `external_id`       // optional
	mfaSerialKey        = `mfa_serial`        // optional
	roleSessionNameKey  = `role_session_name` // optional

	// Additional Config fields
	regionKey = `region`

	// endpoint discovery group
	enableEndpointDiscoveryKey = `endpoint_discovery_enabled` // optional

	// External Credential P/Crocess
	credentialProcessKey = `credential_process` // optional

	// Web Identity Token File
	webIdentityTokenFileKey = `web_identity_token_file` // optional

	// S3 ARN Region Usage
	s3UseARNRegionKey = "s3_use_arn_region"

	// ErrCodeSharedConfig AWS SDK Error Code for Shared Configuration Errors
	ErrCodeSharedConfig = "SharedConfigErr"
)

// DefaultSharedConfigProfile is the default profile to be used when
// loading configuration from the config files if another profile name
// is not provided.
var DefaultSharedConfigProfile = `default`

// ErrSharedConfigSourceCollision will be returned if a section contains both
// source_profile and credential_source
var ErrSharedConfigSourceCollision = awserr.New(ErrCodeSharedConfig, "only source profile or credential source can be specified, not both", nil)

// ErrSharedConfigECSContainerEnvVarEmpty will be returned if the environment
// variables are empty and Environment was set as the credential source
var ErrSharedConfigECSContainerEnvVarEmpty = awserr.New(ErrCodeSharedConfig, "EcsContainer was specified as the credential_source, but 'AWS_CONTAINER_CREDENTIALS_RELATIVE_URI' was not set", nil)

// ErrSharedConfigInvalidCredSource will be returned if an invalid credential source was provided
var ErrSharedConfigInvalidCredSource = awserr.New(ErrCodeSharedConfig, "credential source values must be EcsContainer, Ec2InstanceMetadata, or Environment", nil)

// DefaultSharedCredentialsFilename returns the SDK's default file path
// for the shared credentials file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.aws/credentials
//   - Windows: %USERPROFILE%\.aws\credentials
func DefaultSharedCredentialsFilename() string {
	return filepath.Join(userHomeDir(), ".aws", "credentials")
}

// DefaultSharedConfigFilename returns the SDK's default file path for
// the shared config file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.aws/config
//   - Windows: %USERPROFILE%\.aws\config
func DefaultSharedConfigFilename() string {
	return filepath.Join(userHomeDir(), ".aws", "config")
}

// DefaultSharedConfigFiles is a slice of the default shared config files that
// the will be used in order to load the SharedConfig.
var DefaultSharedConfigFiles = []string{
	DefaultSharedCredentialsFilename(),
	DefaultSharedConfigFilename(),
}

// AssumeRoleConfig provides the values defining the configuration for an IAM
// assume role.
type AssumeRoleConfig struct {
	RoleARN         string
	ExternalID      string
	MFASerial       string
	RoleSessionName string

	SourceProfileName string
	Source            *SharedConfig
}

// SharedConfig represents the configuration fields of the SDK config files.
type SharedConfig struct {
	Profile string

	// Credentials values from the config file. Both aws_access_key_id
	// and aws_secret_access_key must be provided together in the same file
	// to be considered valid. The values will be ignored if not a complete group.
	// aws_session_token is an optional field that can be provided if both of the
	// other two fields are also provided.
	//
	//	aws_access_key_id
	//	aws_secret_access_key
	//	aws_session_token
	Credentials aws.Credentials

	CredentialSource     string
	CredentialProcess    string
	WebIdentityTokenFile string

	AssumeRole AssumeRoleConfig

	// Region is the region the SDK should use for looking up AWS service endpoints
	// and signing requests.
	//
	//	region
	Region string

	// EnableEndpointDiscovery can be enabled in the shared config by setting
	// endpoint_discovery_enabled to true
	//
	//	endpoint_discovery_enabled = true
	EnableEndpointDiscovery *bool

	// Specifies if the S3 service should allow ARNs to direct the region
	// the client's requests are sent to.
	//
	// s3_use_arn_region=true
	S3UseARNRegion *bool
}

// GetEnableEndpointDiscovery returns whether to enable service endpoint discovery
func (c *SharedConfig) GetEnableEndpointDiscovery() (value, ok bool, err error) {
	if c.EnableEndpointDiscovery == nil {
		return false, false, nil
	}

	return *c.EnableEndpointDiscovery, true, nil
}

// GetS3UseARNRegion retions if the S3 service should allow ARNs to direct the region
// the client's requests are sent to.
func (c *SharedConfig) GetS3UseARNRegion() (value, ok bool, err error) {
	if c.S3UseARNRegion == nil {
		return false, false, nil
	}

	return *c.S3UseARNRegion, true, nil
}

// GetRegion returns the region for the profile if a region is set.
func (c SharedConfig) GetRegion() (string, error) {
	return c.Region, nil
}

// GetCredentialsValue returns the credentials for a profile if they were set.
func (c SharedConfig) GetCredentialsValue() (aws.Credentials, error) {
	return c.Credentials, nil
}

// GetAssumeRoleConfig returns the assume role config for a profile. Will be
// a zero value if not set.
func (c SharedConfig) GetAssumeRoleConfig() (AssumeRoleConfig, error) {
	return c.AssumeRole, nil
}

// LoadSharedConfigIgnoreNotExist is an alias for LoadSharedConfig with the
// addition of ignoring when none of the files exist or when the profile
// is not found in any of the files.
func LoadSharedConfigIgnoreNotExist(configs Configs) (Config, error) {
	cfg, err := LoadSharedConfig(configs)
	if err != nil {
		if _, ok := err.(SharedConfigNotExistErrors); ok {
			return SharedConfig{}, nil
		}
		return nil, err
	}

	return cfg, nil
}

// LoadSharedConfig uses the Configs passed in to load the SharedConfig from file
// The file names and profile name are sourced from the Configs.
//
// If profile name is not provided DefaultSharedConfigProfile (default) will
// be used.
//
// If shared config filenames are not provided DefaultSharedConfigFiles will
// be used.
//
// Config providers used:
// * SharedConfigProfileProvider
// * SharedConfigFilesProvider
func LoadSharedConfig(configs Configs) (Config, error) {
	var profile string
	var files []string
	var ok bool
	var err error

	profile, ok, err = GetSharedConfigProfile(configs)
	if err != nil {
		return nil, err
	}
	if !ok {
		profile = DefaultSharedConfigProfile
	}

	files, ok, err = GetSharedConfigFiles(configs)
	if err != nil {
		return nil, err
	}
	if !ok {
		files = DefaultSharedConfigFiles
	}

	return NewSharedConfig(profile, files)
}

// NewSharedConfig retrieves the configuration from the list of files
// using the profile provided. The order the files are listed will determine
// precedence. Values in subsequent files will overwrite values defined in
// earlier files.
//
// For example, given two files A and B. Both define credentials. If the order
// of the files are A then B, B's credential values will be used instead of A's.
func NewSharedConfig(profile string, filenames []string) (SharedConfig, error) {
	if len(filenames) == 0 {
		return SharedConfig{}, fmt.Errorf("no shared config files provided")
	}

	files, err := loadSharedConfigIniFiles(filenames)
	if err != nil {
		return SharedConfig{}, err
	}

	cfg := SharedConfig{}
	profiles := map[string]struct{}{}
	if err = cfg.setFromIniFiles(profiles, profile, files); err != nil {
		return SharedConfig{}, err
	}

	return cfg, nil
}

type sharedConfigFile struct {
	Filename string
	IniData  ini.Sections
}

func loadSharedConfigIniFiles(filenames []string) ([]sharedConfigFile, error) {
	files := make([]sharedConfigFile, 0, len(filenames))

	var errs SharedConfigNotExistErrors
	for _, filename := range filenames {
		sections, err := ini.OpenFile(filename)
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == ini.ErrCodeUnableToReadFile {
			errs = append(errs,
				SharedConfigFileNotExistError{Filename: filename, Err: err},
			)
			// Skip files which can't be opened and read for whatever reason
			continue
		} else if err != nil {
			return nil, SharedConfigLoadError{Filename: filename, Err: err}
		}

		files = append(files, sharedConfigFile{
			Filename: filename, IniData: sections,
		})
	}

	if len(files) == 0 {
		return nil, errs
	}

	return files, nil
}

// Returns an error if all of the files fail to load. If at least one file is
// successfully loaded and contains the profile, no error will be returned.
func (c *SharedConfig) setFromIniFiles(profiles map[string]struct{}, profile string, files []sharedConfigFile) error {
	c.Profile = profile

	// Trim files from the list that don't exist.
	existErrs := SharedConfigNotExistErrors{}
	for _, f := range files {
		if err := c.setFromIniFile(profile, f); err != nil {
			if _, ok := err.(SharedConfigProfileNotExistError); ok {
				existErrs = append(existErrs, err)
				continue
			}
			return err
		}
	}

	if len(existErrs) == len(files) {
		return existErrs
	}

	if _, ok := profiles[profile]; ok {
		// if this is the second instance of the profile the Assume Role
		// options must be cleared because they are only valid for the
		// first reference of a profile. The self linked instance of the
		// profile only have credential provider options.
		c.AssumeRole = AssumeRoleConfig{}
	} else {
		// First time a profile has been seen, It must either be a assume role
		// or credentials. Assert if the credential type requires a role ARN,
		// the ARN is also set.
		if err := c.validateCredentialsRequireARN(profile); err != nil {
			return err
		}
	}
	profiles[profile] = struct{}{}

	if err := c.validateCredentialType(); err != nil {
		return err
	}

	// Link source profiles for assume roles
	if len(c.AssumeRole.SourceProfileName) != 0 {
		// Linked profile via source_profile ignore credential provider
		// options, the source profile must provide the credentials.
		c.clearCredentialOptions()

		srcCfg := &SharedConfig{}
		err := srcCfg.setFromIniFiles(profiles, c.AssumeRole.SourceProfileName, files)
		if err != nil {
			// SourceProfileName that doesn't exist is an error in configuration.
			if _, ok := err.(SharedConfigNotExistErrors); ok {
				err = SharedConfigAssumeRoleError{
					RoleARN: c.AssumeRole.RoleARN,
					Profile: c.AssumeRole.SourceProfileName,
					Err:     err,
				}
			}
			return err
		}

		if !srcCfg.hasCredentials() {
			return SharedConfigAssumeRoleError{
				RoleARN: c.AssumeRole.RoleARN,
				Profile: c.AssumeRole.SourceProfileName,
			}
		}

		c.AssumeRole.Source = srcCfg
	}

	return nil
}

// setFromFile loads the configuration from the file using
// the profile provided. A SharedConfig pointer type value is used so that
// multiple config file loadings can be chained.
//
// Only loads complete logically grouped values, and will not set fields in cfg
// for incomplete grouped values in the config. Such as credentials. For example
// if a config file only includes aws_access_key_id but no aws_secret_access_key
// the aws_access_key_id will be ignored.
func (c *SharedConfig) setFromIniFile(profile string, file sharedConfigFile) error {
	section, ok := file.IniData.GetSection(profile)
	if !ok {
		// Fallback to to alternate profile name: profile <name>
		section, ok = file.IniData.GetSection(fmt.Sprintf("profile %s", profile))
		if !ok {
			return SharedConfigProfileNotExistError{
				Filename: file.Filename,
				Profile:  profile,
				Err:      nil,
			}
		}
	}

	// Shared Credentials
	akid := section.String(accessKeyIDKey)
	secret := section.String(secretAccessKey)
	if len(akid) > 0 && len(secret) > 0 {
		c.Credentials = aws.Credentials{
			AccessKeyID:     akid,
			SecretAccessKey: secret,
			SessionToken:    section.String(sessionTokenKey),
			Source:          fmt.Sprintf("SharedConfigCredentials: %s", file.Filename),
		}
	}

	// Assume Role
	roleArn := section.String(roleArnKey)
	srcProfile := section.String(sourceProfileKey)
	if len(roleArn) > 0 && len(srcProfile) > 0 {
		c.AssumeRole = AssumeRoleConfig{
			RoleARN:           roleArn,
			ExternalID:        section.String(externalIDKey),
			MFASerial:         section.String(mfaSerialKey),
			RoleSessionName:   section.String(roleSessionNameKey),
			SourceProfileName: srcProfile,
		}
	}

	if section.Has(credentialProcessKey) {
		c.CredentialProcess = section.String(credentialProcessKey)
	}

	if section.Has(webIdentityTokenFileKey) {
		c.WebIdentityTokenFile = section.String(webIdentityTokenFileKey)
	}

	// Region
	if v := section.String(regionKey); len(v) > 0 {
		c.Region = v
	}

	// S3 Use ARN Region
	if section.Has(s3UseARNRegionKey) {
		v := section.Bool(s3UseARNRegionKey)
		c.S3UseARNRegion = &v
	}

	// Endpoint discovery
	if section.Has(enableEndpointDiscoveryKey) {
		v := section.Bool(enableEndpointDiscoveryKey)
		c.EnableEndpointDiscovery = &v
	}

	return nil
}

func (c *SharedConfig) validateCredentialsRequireARN(profile string) error {
	var credSource string

	switch {
	case len(c.AssumeRole.SourceProfileName) != 0:
		credSource = sourceProfileKey
	case len(c.CredentialSource) != 0:
		credSource = credentialSourceKey
	case len(c.WebIdentityTokenFile) != 0:
		credSource = webIdentityTokenFileKey
	}

	if len(credSource) != 0 && len(c.AssumeRole.RoleARN) == 0 {
		return CredentialRequiresARNError{
			Type:    credSource,
			Profile: profile,
		}
	}

	return nil
}

func (c *SharedConfig) validateCredentialType() error {
	// Only one or no credential type can be defined.
	if !oneOrNone(
		len(c.AssumeRole.SourceProfileName) != 0,
		len(c.CredentialSource) != 0,
		len(c.CredentialProcess) != 0,
		len(c.WebIdentityTokenFile) != 0,
	) {
		return ErrSharedConfigSourceCollision
	}

	return nil
}

func (c *SharedConfig) hasCredentials() bool {
	switch {
	case len(c.AssumeRole.SourceProfileName) != 0:
	case len(c.CredentialSource) != 0:
	case len(c.CredentialProcess) != 0:
	case len(c.WebIdentityTokenFile) != 0:
	case c.Credentials.HasKeys():
	default:
		return false
	}

	return true
}

func (c *SharedConfig) clearCredentialOptions() {
	c.CredentialSource = ""
	c.CredentialProcess = ""
	c.WebIdentityTokenFile = ""
	c.Credentials = aws.Credentials{}
}

// SharedConfigNotExistErrors provides an error type for failure to load shared
// config because resources do not exist.
type SharedConfigNotExistErrors []error

func (es SharedConfigNotExistErrors) Error() string {
	msg := "failed to load shared config\n"
	for _, e := range es {
		msg += "\t" + e.Error()
	}
	return msg
}

// SharedConfigLoadError is an error for the shared config file failed to load.
type SharedConfigLoadError struct {
	Filename string
	Err      error
}

// Cause is the underlying error that caused the failure.
func (e SharedConfigLoadError) Cause() error {
	return e.Err
}

func (e SharedConfigLoadError) Error() string {
	return fmt.Sprintf("failed to load shared config file, %s, %v", e.Filename, e.Err)
}

// SharedConfigFileNotExistError is an error for the shared config when
// the filename does not exist.
type SharedConfigFileNotExistError struct {
	Filename string
	Profile  string
	Err      error
}

// Cause is the underlying error that caused the failure.
func (e SharedConfigFileNotExistError) Cause() error {
	return e.Err
}

func (e SharedConfigFileNotExistError) Error() string {
	return fmt.Sprintf("failed to open shared config file, %s, %v", e.Filename, e.Err)
}

// SharedConfigProfileNotExistError is an error for the shared config when
// the profile was not find in the config file.
type SharedConfigProfileNotExistError struct {
	Filename string
	Profile  string
	Err      error
}

// Cause is the underlying error that caused the failure.
func (e SharedConfigProfileNotExistError) Cause() error {
	return e.Err
}

func (e SharedConfigProfileNotExistError) Error() string {
	return fmt.Sprintf("failed to get shared config profile, %s, in %s, %v", e.Profile, e.Filename, e.Err)
}

// SharedConfigAssumeRoleError is an error for the shared config when the
// profile contains assume role information, but that information is invalid
// or not complete.
type SharedConfigAssumeRoleError struct {
	Profile string
	RoleARN string
	Err     error
}

func (e SharedConfigAssumeRoleError) Error() string {
	return fmt.Sprintf("failed to load assume role %s, of profile %s, %v",
		e.RoleARN, e.Profile, e.Err)
}

// CredentialRequiresARNError provides the error for shared config credentials
// that are incorrectly configured in the shared config or credentials file.
type CredentialRequiresARNError struct {
	// type of credentials that were configured.
	Type string

	// Profile name the credentials were in.
	Profile string
}

// Code is the short id of the error.
func (e CredentialRequiresARNError) Code() string {
	return "CredentialRequiresARNError"
}

// Message is the description of the error
func (e CredentialRequiresARNError) Message() string {
	return fmt.Sprintf(
		"credential type %s requires role_arn, profile %s",
		e.Type, e.Profile,
	)
}

// OrigErr is the underlying error that caused the failure.
func (e CredentialRequiresARNError) OrigErr() error {
	return nil
}

// Error satisfies the error interface.
func (e CredentialRequiresARNError) Error() string {
	return awserr.SprintError(e.Code(), e.Message(), "", nil)
}

func userHomeDir() string {
	if runtime.GOOS == "windows" { // Windows
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}

func oneOrNone(bs ...bool) bool {
	var count int

	for _, b := range bs {
		if b {
			count++
			if count > 1 {
				return false
			}
		}
	}

	return true
}
