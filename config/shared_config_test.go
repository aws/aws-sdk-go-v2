package config

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/ini"
	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/ptr"
)

var _ regionProvider = (*SharedConfig)(nil)

var (
	testConfigFilename      = filepath.Join("testdata", "shared_config")
	testConfigOtherFilename = filepath.Join("testdata", "shared_config_other")
)

func TestNewSharedConfig(t *testing.T) {
	cases := map[string]struct {
		Filenames []string
		Profile   string
		Expected  SharedConfig
		Err       error
	}{
		"file not exist": {
			Filenames: []string{"file_not_exist"},
			Profile:   "default",
			Err:       fmt.Errorf("failed to get shared config profile"),
		},
		"default profile": {
			Filenames: []string{testConfigFilename},
			Profile:   "default",
			Expected: SharedConfig{
				Profile: "default",
				Region:  "default_region",
			},
		},
		"multiple config files": {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "config_file_load_order",
			Expected: SharedConfig{
				Profile: "config_file_load_order",
				Region:  "shared_config_region",
				Credentials: aws.Credentials{
					AccessKeyID:     "shared_config_akid",
					SecretAccessKey: "shared_config_secret",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
			},
		},
		"mutliple config files reverse order": {
			Filenames: []string{testConfigFilename, testConfigOtherFilename},
			Profile:   "config_file_load_order",
			Expected: SharedConfig{
				Profile: "config_file_load_order",
				Region:  "shared_config_other_region",
				Credentials: aws.Credentials{
					AccessKeyID:     "shared_config_other_akid",
					SecretAccessKey: "shared_config_other_secret",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigOtherFilename),
				},
			},
		},
		"Assume role": {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role",
			Expected: SharedConfig{
				Profile:           "assume_role",
				RoleARN:           "assume_role_role_arn",
				SourceProfileName: "complete_creds",
				Source: &SharedConfig{
					Profile: "complete_creds",
					Credentials: aws.Credentials{
						AccessKeyID:     "complete_creds_akid",
						SecretAccessKey: "complete_creds_secret",
						Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
					},
				},
			},
		},
		"Assume role with invalid source profile": {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_invalid_source_profile",
			Err: SharedConfigAssumeRoleError{
				Profile: "profile_not_exists",
				RoleARN: "assume_role_invalid_source_profile_role_arn",
				Err: SharedConfigProfileNotExistError{
					Profile: "profile_not_exists",
					Err:     nil,
				},
			},
		},
		"Assume role with creds": {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_w_creds",
			Expected: SharedConfig{
				Profile:           "assume_role_w_creds",
				RoleARN:           "assume_role_w_creds_role_arn",
				ExternalID:        "1234",
				RoleSessionName:   "assume_role_w_creds_session_name",
				SourceProfileName: "assume_role_w_creds",
				Source: &SharedConfig{
					Profile: "assume_role_w_creds",
					Credentials: aws.Credentials{
						AccessKeyID:     "assume_role_w_creds_akid",
						SecretAccessKey: "assume_role_w_creds_secret",
						Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
					},
				},
			},
		},
		"Assume role without creds": {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_wo_creds",
			Expected: SharedConfig{
				Profile:           "assume_role_wo_creds",
				RoleARN:           "assume_role_wo_creds_role_arn",
				SourceProfileName: "assume_role_wo_creds",
			},
			Err: SharedConfigAssumeRoleError{
				Profile: "assume_role_wo_creds",
				RoleARN: "assume_role_wo_creds_role_arn",
			},
		},
		"Invalid INI file": {
			Filenames: []string{filepath.Join("testdata", "shared_config_invalid_ini")},
			Profile:   "profile_name",
			Err: SharedConfigLoadError{
				Filename: filepath.Join("testdata", "shared_config_invalid_ini"),
				Err:      fmt.Errorf("invalid state"),
			},
		},
		"S3UseARNRegion property on profile": {
			Profile:   "valid_arn_region",
			Filenames: []string{testConfigFilename},
			Expected: SharedConfig{
				Profile:        "valid_arn_region",
				S3UseARNRegion: ptr.Bool(true),
			},
		},
		"EndpointDiscovery property on profile": {
			Profile:   "endpoint_discovery",
			Filenames: []string{testConfigFilename},
			Expected: SharedConfig{
				Profile:                 "endpoint_discovery",
				EnableEndpointDiscovery: ptr.Bool(true),
			},
		},
		"Assume role with credential source Ec2Metadata": {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_with_credential_source",
			Expected: SharedConfig{
				Profile:          "assume_role_with_credential_source",
				RoleARN:          "assume_role_with_credential_source_role_arn",
				CredentialSource: credSourceEc2Metadata,
			},
		},
		"Assume role chained with creds": {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "multiple_assume_role",
			Expected: SharedConfig{
				Profile:           "multiple_assume_role",
				RoleARN:           "multiple_assume_role_role_arn",
				SourceProfileName: "assume_role",
				Source: &SharedConfig{
					Profile:           "assume_role",
					RoleARN:           "assume_role_role_arn",
					SourceProfileName: "complete_creds",
					Source: &SharedConfig{
						Profile: "complete_creds",
						Credentials: aws.Credentials{
							AccessKeyID:     "complete_creds_akid",
							SecretAccessKey: "complete_creds_secret",
							Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
						},
					},
				},
			},
		},
		"Assume role chained with credential source": {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "multiple_assume_role_with_credential_source",
			Expected: SharedConfig{
				Profile:           "multiple_assume_role_with_credential_source",
				RoleARN:           "multiple_assume_role_with_credential_source_role_arn",
				SourceProfileName: "assume_role_with_credential_source",
				Source: &SharedConfig{
					Profile:          "assume_role_with_credential_source",
					RoleARN:          "assume_role_with_credential_source_role_arn",
					CredentialSource: credSourceEc2Metadata,
				},
			},
		},
		"Assume role chained with credential source reversed order": {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "multiple_assume_role_with_credential_source2",
			Expected: SharedConfig{
				Profile:           "multiple_assume_role_with_credential_source2",
				RoleARN:           "multiple_assume_role_with_credential_source2_role_arn",
				SourceProfileName: "multiple_assume_role_with_credential_source",
				Source: &SharedConfig{
					Profile:           "multiple_assume_role_with_credential_source",
					RoleARN:           "multiple_assume_role_with_credential_source_role_arn",
					SourceProfileName: "assume_role_with_credential_source",
					Source: &SharedConfig{
						Profile:          "assume_role_with_credential_source",
						RoleARN:          "assume_role_with_credential_source_role_arn",
						CredentialSource: credSourceEc2Metadata,
					},
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			cfg, err := LoadSharedConfigProfile(context.TODO(), c.Profile, func(o *LoadSharedConfigOptions) {
				o.ConfigFiles = c.Filenames
				o.CredentialsFiles = []string{filepath.Join("testdata", "empty_creds_config")}
			})
			if c.Err != nil && err != nil {
				if e, a := c.Err.Error(), err.Error(); !strings.Contains(a, e) {
					t.Errorf("expect %q to be in %q", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if c.Err != nil {
				t.Errorf("expect error: %v, got none", c.Err)
			}
			if e, a := c.Expected, cfg; !reflect.DeepEqual(e, a) {
				t.Errorf(" expect %v, got %v", e, a)
			}
		})
	}
}

func TestLoadSharedConfigFromSection(t *testing.T) {
	filename := testConfigFilename
	sections, err := ini.OpenFile(filename)

	if err != nil {
		t.Fatalf("failed to load test config file, %s, %v", filename, err)
	}
	cases := map[string]struct {
		Profile  string
		Expected SharedConfig
		Err      error
	}{
		"Default as profile": {
			Profile:  "default",
			Expected: SharedConfig{Region: "default_region"},
		},
		"prefixed profile": {
			Profile:  "profile alt_profile_name",
			Expected: SharedConfig{Region: "alt_profile_name_region"},
		},
		"prefixed profile 2": {
			Profile:  "profile short_profile_name_first",
			Expected: SharedConfig{Region: "short_profile_name_first_alt"},
		},
		"profile with partial creds": {
			Profile:  "profile partial_creds",
			Expected: SharedConfig{},
		},
		"profile with complete creds": {
			Profile: "profile complete_creds",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "complete_creds_akid",
					SecretAccessKey: "complete_creds_secret",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", filename),
				},
			},
		},
		"profile with complete creds and token": {
			Profile: "profile complete_creds_with_token",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "complete_creds_with_token_akid",
					SecretAccessKey: "complete_creds_with_token_secret",
					SessionToken:    "complete_creds_with_token_token",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", filename),
				},
			},
		},
		"complete profile": {
			Profile: "profile full_profile",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "full_profile_akid",
					SecretAccessKey: "full_profile_secret",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", filename),
				},
				Region: "full_profile_region",
			},
		},
		"profile with partial assume role": {
			Profile: "profile partial_assume_role",
			Expected: SharedConfig{
				RoleARN: "partial_assume_role_role_arn",
			},
		},
		"profile using assume role": {
			Profile: "profile assume_role",
			Expected: SharedConfig{
				RoleARN:           "assume_role_role_arn",
				SourceProfileName: "complete_creds",
			},
		},
		"profile with assume role and MFA": {
			Profile: "profile assume_role_w_mfa",
			Expected: SharedConfig{
				RoleARN:           "assume_role_role_arn",
				SourceProfileName: "complete_creds",
				MFASerial:         "0123456789",
			},
		},
		"does not exist": {
			Profile: "does_not_exist",
			Err: SharedConfigProfileNotExistError{
				Filename: []string{filename},
				Profile:  "does_not_exist",
				Err:      nil,
			},
		},
		"profile with mixed casing": {
			Profile: "profile with_mixed_case_keys",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "accessKey",
					SecretAccessKey: "secret",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", filename),
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var cfg SharedConfig

			section, ok := sections.GetSection(c.Profile)
			if !ok {
				if c.Err == nil {
					t.Fatalf("expected section to be present, was not")
				} else {
					if e, a := c.Err.Error(), "failed to get shared config profile"; !strings.Contains(e, a) {
						t.Fatalf("expect %q to be in %q", a, e)
					}
					return
				}
			}

			err := cfg.setFromIniSection(c.Profile, section)
			if c.Err != nil {
				if e, a := c.Err.Error(), err.Error(); !strings.Contains(a, e) {
					t.Errorf("expect %q to be in %q", e, a)
				}
				return
			}

			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := c.Expected, cfg; !reflect.DeepEqual(e, a) {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestLoadSharedConfig(t *testing.T) {
	origProf := defaultSharedConfigProfile
	origConfigFiles := DefaultSharedConfigFiles
	origCredentialFiles := DefaultSharedCredentialsFiles
	defer func() {
		defaultSharedConfigProfile = origProf
		DefaultSharedConfigFiles = origConfigFiles
		DefaultSharedCredentialsFiles = origCredentialFiles
	}()

	cases := []struct {
		LoadOptionFn func(*LoadOptions) error
		Files        []string
		Profile      string
		LoadFn       func(context.Context, configs) (Config, error)
		Expect       SharedConfig
		Err          string
	}{
		{
			LoadOptionFn: WithSharedConfigProfile("alt_profile_name"),
			Files: []string{
				filepath.Join("testdata", "shared_config"),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "alt_profile_name",
				Region:  "alt_profile_name_region",
			},
		},
		{
			LoadOptionFn: WithSharedConfigFiles([]string{
				filepath.Join("testdata", "shared_config"),
			}),
			Profile: "alt_profile_name",
			LoadFn:  loadSharedConfig,
			Expect: SharedConfig{
				Profile: "alt_profile_name",
				Region:  "alt_profile_name_region",
			},
		},
		{
			LoadOptionFn: WithSharedConfigProfile("default"),
			Files: []string{
				filepath.Join("file_not_exist"),
			},
			LoadFn: loadSharedConfig,
			Err:    "failed to get shared config profile",
		},
		{
			LoadOptionFn: WithSharedConfigProfile("profile_not_exist"),
			Files: []string{
				filepath.Join("testdata", "shared_config"),
			},
			LoadFn: loadSharedConfig,
			Err:    "failed to get shared config profile",
		},
		{
			LoadOptionFn: WithSharedConfigProfile("default"),
			Files: []string{
				filepath.Join("file_not_exist"),
			},
			LoadFn: loadSharedConfigIgnoreNotExist,
		},
		{
			LoadOptionFn: WithSharedConfigProfile("assume_role_invalid_source_profile"),
			Files: []string{
				testConfigOtherFilename, testConfigFilename,
			},
			LoadFn: loadSharedConfig,
			Err:    "failed to get shared config profile",
		},
		{
			LoadOptionFn: WithSharedConfigProfile("assume_role_invalid_source_profile"),
			Files: []string{
				testConfigOtherFilename, testConfigFilename,
			},
			LoadFn: loadSharedConfigIgnoreNotExist,
			Err:    "failed to get shared config profile",
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defaultSharedConfigProfile = origProf
			DefaultSharedConfigFiles = origConfigFiles
			DefaultSharedCredentialsFiles = origCredentialFiles

			if len(c.Profile) > 0 {
				defaultSharedConfigProfile = c.Profile
			}
			if len(c.Files) > 0 {
				DefaultSharedConfigFiles = c.Files
			}

			DefaultSharedCredentialsFiles = []string{}

			var options LoadOptions
			c.LoadOptionFn(&options)

			cfg, err := c.LoadFn(context.Background(), configs{options})
			if len(c.Err) > 0 {
				if err == nil {
					t.Fatalf("expected error %v, got none", c.Err)
				}
				if e, a := c.Err, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %q to be in %q", e, a)
				}
				return
			} else if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := c.Expect, cfg; !reflect.DeepEqual(e, a) {
				t.Errorf("expect %v got %v", e, a)
			}
		})
	}
}

func TestSharedConfigLoading(t *testing.T) {
	// initialize a logger
	var loggerBuf bytes.Buffer
	logger := logging.NewStandardLogger(&loggerBuf)

	cases := map[string]struct {
		LoadOptionFns []func(*LoadOptions) error
		LoadFn        func(context.Context, configs) (Config, error)
		Expect        SharedConfig
		ExpectLog     string
		Err           string
	}{
		"duplicate profiles in the configuration files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("duplicate-profile"),
				WithSharedConfigFiles([]string{filepath.Join("testdata", "load_config")}),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "duplicate-profile",
				Region:  "us-west-2",
			},
			ExpectLog: "For profile: profile duplicate-profile, overriding region value, with a region value found in a " +
				"duplicate profile defined later in the same file",
		},

		"profile prefix not used in the configuration files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("no-such-profile"),
				WithSharedConfigFiles([]string{filepath.Join("testdata", "load_config")}),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{},
			Err:    "failed to get shared config profile",
		},

		"profile prefix overrides default": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigFiles([]string{filepath.Join("testdata", "load_config")}),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "default",
				Region:  "ap-north-1",
			},
			ExpectLog: "overrided non-prefixed default profile",
		},

		"duplicate profiles in credentials file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("duplicate-profile"),
				WithSharedConfigFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "load_credentials")}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "duplicate-profile",
				Region:  "us-west-2",
			},
			ExpectLog: "overriding region value, with a region value found in a duplicate profile defined later in the same file",
			Err:       "",
		},

		"profile prefix used in credentials files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("unused-profile"),
				WithSharedConfigFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "load_credentials")}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn:    loadSharedConfig,
			ExpectLog: "profile defined with name `profile unused-profile` is ignored.",
			Err:       "failed to get shared config profile, unused-profile",
		},
		"partial credentials in configuration files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("partial-creds-1"),
				WithSharedConfigFiles([]string{filepath.Join("testdata", "load_config")}),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "partial-creds-1",
			},
			Err: "partial credentials",
		},
		"parital credentials in the credentials files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("partial-creds-1"),
				WithSharedConfigFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "load_credentials")}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "partial-creds-1",
			},
			Err: "partial credentials found for profile partial-creds-1",
		},
		"credentials override configuration profile": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("complete"),
				WithSharedConfigFiles([]string{filepath.Join("testdata", "load_config")}),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "load_credentials")}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "complete",
				Credentials: aws.Credentials{
					AccessKeyID:     "credsAccessKey",
					SecretAccessKey: "credsSecretKey",
					Source: fmt.Sprintf("SharedConfigCredentials: %v",
						filepath.Join("testdata", "load_credentials")),
				},
				Region: "us-west-2",
			},
		},
		"credentials profile has complete credentials": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("complete"),
				WithSharedConfigFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "load_credentials")}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "complete",
				Credentials: aws.Credentials{
					AccessKeyID:     "credsAccessKey",
					SecretAccessKey: "credsSecretKey",
					Source:          fmt.Sprintf("SharedConfigCredentials: %v", filepath.Join("testdata", "load_credentials")),
				},
			},
		},
		"credentials split between multiple credentials files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("partial-creds-1"),
				WithSharedConfigFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "load_credentials"),
					filepath.Join("testdata", "load_credentials_secondary"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "partial-creds-1",
			},
			Err: "partial credentials",
		},
		"credentials split between multiple configuration files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("partial-creds-1"),
				WithSharedCredentialsFiles([]string{filepath.Join("testdata", "empty_creds_config")}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
					filepath.Join("testdata", "load_config_secondary"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "partial-creds-1",
				Region:  "us-west-2",
			},
			ExpectLog: "",
			Err:       "partial credentials",
		},
		"credentials split between credentials and config files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("partial-creds-1"),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
				}),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "load_credentials"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "partial-creds-1",
			},
			ExpectLog: "",
			Err:       "partial credentials",
		},
		"replaced profile with prefixed profile in config files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("replaced-profile"),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
				}),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "replaced-profile",
				Region:  "eu-west-1",
			},
			ExpectLog: "profile defined with name `replaced-profile` is ignored.",
		},
		"replaced profile with prefixed profile in credentials files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("replaced-profile"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "load_credentials"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "replaced-profile",
				Region:  "us-west-2",
			},
			ExpectLog: "profile defined with name `profile replaced-profile` is ignored.",
		},
		"ignored profiles w/o prefixed profile across credentials and config files": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("replaced-profile"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "load_credentials"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "replaced-profile",
				Region:  "us-west-2",
			},
			ExpectLog: "profile defined with name `profile replaced-profile` is ignored.",
		},
		"1. profile with name as `profile` in config file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("profile"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn:    loadSharedConfig,
			Err:       "failed to get shared config profile, profile",
			ExpectLog: "profile defined with name `profile` is ignored",
		},
		"2. profile with name as `profile ` in config file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("profile "),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn:    loadSharedConfig,
			Err:       "failed to get shared config profile, profile",
			ExpectLog: "profile defined with name `profile` is ignored",
		},
		"3. profile with name as `profile\t` in config file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("profile"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn:    loadSharedConfig,
			Err:       "failed to get shared config profile, profile",
			ExpectLog: "profile defined with name `profile` is ignored",
		},
		"profile with tabs as delimiter for profile prefix in config file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("with-tab"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "with-tab",
				Region:  "cn-north-1",
			},
		},
		"profile with tabs as delimiter for profile prefix in credentials file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("with-tab"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "load_credentials"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn:    loadSharedConfig,
			Err:       "failed to get shared config profile, with-tab",
			ExpectLog: "profile defined with name `profile with-tab` is ignored",
		},
		"profile with name as `profile profile` in credentials file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("profile"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "load_credentials"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn:    loadSharedConfig,
			Err:       "failed to get shared config profile, profile",
			ExpectLog: "profile defined with name `profile profile` is ignored",
		},
		"profile with name profile-bar in credentials file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("profile-bar"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "load_credentials"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Expect: SharedConfig{
				Profile: "profile-bar",
				Region:  "us-west-2",
			},
		},
		"profile with name profile-bar in config file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("profile-bar"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "empty_creds_config"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn:    loadSharedConfig,
			Err:       "failed to get shared config profile, profile-bar",
			ExpectLog: "profile defined with name `profile-bar` is ignored",
		},
		"profile ignored in credentials and config file": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedConfigProfile("ignored-profile"),
				WithSharedCredentialsFiles([]string{
					filepath.Join("testdata", "load_credentials"),
				}),
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "load_config"),
				}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn:    loadSharedConfig,
			Err:       "failed to get shared config profile, ignored-profile",
			ExpectLog: "profile defined with name `ignored-profile` is ignored.",
		},
		"config and creds files explicitly set to empty slice": {
			LoadOptionFns: []func(*LoadOptions) error{
				WithSharedCredentialsFiles([]string{}),
				WithSharedConfigFiles([]string{}),
				WithLogConfigurationWarnings(true),
				WithLogger(logger),
			},
			LoadFn: loadSharedConfig,
			Err:    "failed to get shared config profile, default",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			defer loggerBuf.Reset()

			var options LoadOptions

			for _, fn := range c.LoadOptionFns {
				fn(&options)
			}

			cfg, err := c.LoadFn(context.Background(), configs{options})

			if e, a := c.ExpectLog, loggerBuf.String(); !strings.Contains(a, e) {
				t.Errorf("expect %v logged in %v", e, a)
			}
			if loggerBuf.Len() == 0 && len(c.ExpectLog) != 0 {
				t.Errorf("expected log, got none")
			}

			if len(c.Err) > 0 {
				if err == nil {
					t.Fatalf("expected error %v, got none", c.Err)
				}
				if e, a := c.Err, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %q to be in %q", e, a)
				}
				return
			} else if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := c.Expect, cfg; !reflect.DeepEqual(e, a) {
				t.Errorf("expect %v got %v", e, a)
			}

		})
	}
}
