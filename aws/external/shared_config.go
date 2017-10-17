package external

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-ini/ini"
)

const (
	// Static Credentials group
	accessKeyIDKey  = `aws_access_key_id`     // group required
	secretAccessKey = `aws_secret_access_key` // group required
	sessionTokenKey = `aws_session_token`     // optional

	// Assume Role Credentials group
	roleArnKey         = `role_arn`          // group required
	sourceProfileKey   = `source_profile`    // group required
	externalIDKey      = `external_id`       // optional
	mfaSerialKey       = `mfa_serial`        // optional
	roleSessionNameKey = `role_session_name` // optional

	// Additional Config fields
	regionKey = `region`
)

// DefaultSharedConfigProfile is the default profile to be used when
// loading configuration from the config files if another profile name
// is not provided.
const DefaultSharedConfigProfile = `default`

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
	SourceProfile   string
	ExternalID      string
	MFASerial       string
	RoleSessionName string
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
	Creds aws.Value

	// TODO need good way to expose these in Provider interface
	AssumeRole       AssumeRoleConfig
	AssumeRoleSource *SharedConfig

	// Region is the region the SDK should use for looking up AWS service endpoints
	// and signing requests.
	//
	//	region
	Region string
}

// GetRegion returns the region for the profile if a region is set.
func (c SharedConfig) GetRegion() (string, error) {
	return c.Region, nil
}

// GetCredentialsValue returns the credentials for a profile if they were set.
func (c SharedConfig) GetCredentialsValue() (aws.Value, error) {
	return c.Creds, nil
}

// StaticSharedConfigProfile wraps a strings to satisfy the SharedConfigProfileProvider
// interface so a slice of custom shared config files ared used when loading the
// SharedConfig.
type StaticSharedConfigProfile string

// GetSharedConfigProfile returns the shared config profile.
func (c StaticSharedConfigProfile) GetSharedConfigProfile() (string, error) {
	return string(c), nil
}

// StaticSharedConfigFiles wraps a slice of strings to satisfy the
// SharedConfigFilesProvider interface so a slice of custom shared config files
// ared used when loading the SharedConfig.
type StaticSharedConfigFiles []string

// GetSharedConfigFiles returns the slice of shared config files.
func (c StaticSharedConfigFiles) GetSharedConfigFiles() ([]string, error) {
	return []string(c), nil
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
func LoadSharedConfig(cfgs Configs) (Config, error) {
	var profile string
	var files []string

	for _, cfg := range cfgs {
		if len(profile) == 0 {
			if g, ok := cfg.(SharedConfigProfileProvider); ok {
				profile, _ = g.GetSharedConfigProfile()
				// TODO error handling...
			}

		}
		if len(files) == 0 {
			if g, ok := cfg.(SharedConfigFilesProvider); ok {
				files, _ = g.GetSharedConfigFiles()
				// TODO error handling...
			}
		}
	}
	if len(profile) == 0 {
		profile = DefaultSharedConfigProfile
	}
	if len(files) == 0 {
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
//
// See SharedConfig.setFromFile for information how the config files
// will be loaded.
func NewSharedConfig(profile string, filenames []string) (SharedConfig, error) {
	if len(profile) == 0 {
		profile = DefaultSharedConfigProfile
	}

	files, err := loadSharedConfigIniFiles(filenames)
	if err != nil {
		return SharedConfig{}, err
	}

	cfg := SharedConfig{}
	if err = cfg.setFromIniFiles(profile, files); err != nil {
		return SharedConfig{}, err
	}

	if len(cfg.AssumeRole.SourceProfile) > 0 {
		if err := cfg.setAssumeRoleSource(profile, files); err != nil {
			return SharedConfig{}, err
		}
	}

	return cfg, nil
}

type sharedConfigFile struct {
	Filename string
	IniData  *ini.File
}

func loadSharedConfigIniFiles(filenames []string) ([]sharedConfigFile, error) {
	files := make([]sharedConfigFile, 0, len(filenames))

	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			// Skip files which can't be opened and read for whatever reason
			continue
		}

		f, err := ini.Load(b)
		if err != nil {
			return nil, SharedConfigLoadError{Filename: filename, Err: err}
		}

		files = append(files, sharedConfigFile{
			Filename: filename, IniData: f,
		})
	}

	return files, nil
}

func (cfg *SharedConfig) setAssumeRoleSource(origProfile string, files []sharedConfigFile) error {
	var assumeRoleSrc SharedConfig

	// Multiple level assume role chains are not support
	if cfg.AssumeRole.SourceProfile == origProfile {
		assumeRoleSrc = *cfg
		assumeRoleSrc.AssumeRole = AssumeRoleConfig{}
	} else {
		err := assumeRoleSrc.setFromIniFiles(cfg.AssumeRole.SourceProfile, files)
		if err != nil {
			return err
		}
	}

	if len(assumeRoleSrc.Creds.AccessKeyID) == 0 {
		return SharedConfigAssumeRoleError{RoleARN: cfg.AssumeRole.RoleARN}
	}

	cfg.AssumeRoleSource = &assumeRoleSrc

	return nil
}

func (cfg *SharedConfig) setFromIniFiles(profile string, files []sharedConfigFile) error {
	cfg.Profile = profile

	// Trim files from the list that don't exist.
	for _, f := range files {
		if err := cfg.setFromIniFile(profile, f); err != nil {
			if _, ok := err.(SharedConfigProfileNotExistsError); ok {
				// Ignore proviles missings
				continue
			}
			return err
		}
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
func (cfg *SharedConfig) setFromIniFile(profile string, file sharedConfigFile) error {
	section, err := file.IniData.GetSection(profile)
	if err != nil {
		// Fallback to to alternate profile name: profile <name>
		section, err = file.IniData.GetSection(fmt.Sprintf("profile %s", profile))
		if err != nil {
			return SharedConfigProfileNotExistsError{Profile: profile, Err: err}
		}
	}

	// Shared Credentials
	akid := section.Key(accessKeyIDKey).String()
	secret := section.Key(secretAccessKey).String()
	if len(akid) > 0 && len(secret) > 0 {
		cfg.Creds = aws.Value{
			AccessKeyID:     akid,
			SecretAccessKey: secret,
			SessionToken:    section.Key(sessionTokenKey).String(),
			ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", file.Filename),
		}
	}

	// Assume Role
	roleArn := section.Key(roleArnKey).String()
	srcProfile := section.Key(sourceProfileKey).String()
	if len(roleArn) > 0 && len(srcProfile) > 0 {
		cfg.AssumeRole = AssumeRoleConfig{
			RoleARN:         roleArn,
			SourceProfile:   srcProfile,
			ExternalID:      section.Key(externalIDKey).String(),
			MFASerial:       section.Key(mfaSerialKey).String(),
			RoleSessionName: section.Key(roleSessionNameKey).String(),
		}
	}

	// Region
	if v := section.Key(regionKey).String(); len(v) > 0 {
		cfg.Region = v
	}

	return nil
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
	return fmt.Sprintf("failed to load config file, %s", e.Filename)
}

// SharedConfigProfileNotExistsError is an error for the shared config when
// the profile was not find in the config file.
type SharedConfigProfileNotExistsError struct {
	Profile string
	Err     error
}

func (e SharedConfigProfileNotExistsError) Error() string {
	return fmt.Sprintf("failed to get profile, %s", e.Profile)
}

// SharedConfigAssumeRoleError is an error for the shared config when the
// profile contains assume role information, but that information is invalid
// or not complete.
type SharedConfigAssumeRoleError struct {
	RoleARN string
}

func (e SharedConfigAssumeRoleError) Error() string {
	return fmt.Sprintf("failed to load assume role for %s, source profile has no shared credentials",
		e.RoleARN)
}

func userHomeDir() string {
	if runtime.GOOS == "windows" { // Windows
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}
