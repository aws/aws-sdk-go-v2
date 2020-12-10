package config

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/ini"
	"github.com/awslabs/smithy-go/ptr"
)

var _ regionProvider = (*SharedConfig)(nil)

var (
	testConfigFilename      = filepath.Join("testdata", "shared_config")
	testConfigOtherFilename = filepath.Join("testdata", "shared_config_other")
)

func TestNewSharedConfig(t *testing.T) {
	cases := []struct {
		Filenames []string
		Profile   string
		Expected  SharedConfig
		Err       error
	}{
		0: {
			Filenames: []string{"file_not_exist"},
			Profile:   "default",
			Err:       fmt.Errorf("failed to open shared config file, file_not_exist"),
		},
		1: {
			Filenames: []string{testConfigFilename},
			Profile:   "default",
			Expected: SharedConfig{
				Profile: "default",
				Region:  "default_region",
			},
		},
		2: {
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
		3: {
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
		4: {
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
		5: {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_invalid_source_profile",
			Err: SharedConfigAssumeRoleError{
				Profile: "profile_not_exists",
				RoleARN: "assume_role_invalid_source_profile_role_arn",
				Err: SharedConfigNotExistErrors{
					SharedConfigProfileNotExistError{
						Profile:  "profile_not_exists",
						Filename: testConfigOtherFilename,
						Err:      nil,
					},
					SharedConfigProfileNotExistError{
						Profile:  "profile_not_exists",
						Filename: testConfigFilename,
						Err:      nil,
					},
				},
			},
		},
		6: {
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
		7: {
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
		8: {
			Filenames: []string{filepath.Join("testdata", "shared_config_invalid_ini")},
			Profile:   "profile_name",
			Err: SharedConfigLoadError{
				Filename: filepath.Join("testdata", "shared_config_invalid_ini"),
				Err:      fmt.Errorf("invalid state"),
			},
		},
		9: {
			Profile:   "valid_arn_region",
			Filenames: []string{testConfigFilename},
			Expected: SharedConfig{
				Profile:        "valid_arn_region",
				S3UseARNRegion: ptr.Bool(true),
			},
		},
		10: {
			Profile:   "endpoint_discovery",
			Filenames: []string{testConfigFilename},
			Expected: SharedConfig{
				Profile:                 "endpoint_discovery",
				EnableEndpointDiscovery: ptr.Bool(true),
			},
		},
		11: {
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_with_credential_source",
			Expected: SharedConfig{
				Profile:          "assume_role_with_credential_source",
				RoleARN:          "assume_role_with_credential_source_role_arn",
				CredentialSource: credSourceEc2Metadata,
			},
		},
		12: {
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
		13: {
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
		14: {
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

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cfg, err := NewSharedConfig(c.Profile, c.Filenames)
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
				t.Errorf(" expect %v, got %v", e, a)
			}
		})
	}
}

func TestLoadSharedConfigFromFile(t *testing.T) {
	filename := testConfigFilename
	sections, err := ini.OpenFile(filename)
	if err != nil {
		t.Fatalf("failed to load test config file, %s, %v", filename, err)
	}
	iniFile := sharedConfigFile{IniData: sections, Filename: filename}

	cases := []struct {
		Profile  string
		Expected SharedConfig
		Err      error
	}{
		0: {
			Profile:  "default",
			Expected: SharedConfig{Region: "default_region"},
		},
		1: {
			Profile:  "alt_profile_name",
			Expected: SharedConfig{Region: "alt_profile_name_region"},
		},
		2: {
			Profile:  "short_profile_name_first",
			Expected: SharedConfig{Region: "short_profile_name_first_short"},
		},
		3: {
			Profile:  "partial_creds",
			Expected: SharedConfig{},
		},
		4: {
			Profile: "complete_creds",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "complete_creds_akid",
					SecretAccessKey: "complete_creds_secret",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
			},
		},
		5: {
			Profile: "complete_creds_with_token",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "complete_creds_with_token_akid",
					SecretAccessKey: "complete_creds_with_token_secret",
					SessionToken:    "complete_creds_with_token_token",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
			},
		},
		6: {
			Profile: "full_profile",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "full_profile_akid",
					SecretAccessKey: "full_profile_secret",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
				Region: "full_profile_region",
			},
		},
		7: {
			Profile: "partial_assume_role",
			Expected: SharedConfig{
				RoleARN: "partial_assume_role_role_arn",
			},
		},
		8: {
			Profile: "assume_role",
			Expected: SharedConfig{
				RoleARN:           "assume_role_role_arn",
				SourceProfileName: "complete_creds",
			},
		},
		9: {
			Profile: "assume_role_w_mfa",
			Expected: SharedConfig{
				RoleARN:           "assume_role_role_arn",
				SourceProfileName: "complete_creds",
				MFASerial:         "0123456789",
			},
		},
		10: {
			Profile: "does_not_exist",
			Err: SharedConfigProfileNotExistError{
				Filename: filepath.Join("testdata", "shared_config"),
				Profile:  "does_not_exist",
				Err:      nil,
			},
		},
		{
			Profile: "with_mixed_case_keys",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "accessKey",
					SecretAccessKey: "secret",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var cfg SharedConfig
			err := cfg.setFromIniFile(c.Profile, iniFile)
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

func TestLoadSharedConfigIniFiles(t *testing.T) {
	cases := []struct {
		Filenames []string
		Expected  []sharedConfigFile
	}{
		{
			Filenames: []string{"not_exist", testConfigFilename},
			Expected: []sharedConfigFile{
				{Filename: testConfigFilename},
			},
		},
		{
			Filenames: []string{testConfigFilename, testConfigOtherFilename},
			Expected: []sharedConfigFile{
				{Filename: testConfigFilename},
				{Filename: testConfigOtherFilename},
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			files, err := loadSharedConfigIniFiles(c.Filenames)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := len(c.Expected), len(files); e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if e, a := c.Expected, files; !cmpFiles(e, a) {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func cmpFiles(expects, actuals []sharedConfigFile) bool {
	for i, expect := range expects {
		if expect.Filename != actuals[i].Filename {
			return false
		}
	}
	return true
}

func TestLoadSharedConfig(t *testing.T) {
	origProf := defaultSharedConfigProfile
	origFiles := DefaultSharedConfigFiles
	defer func() {
		defaultSharedConfigProfile = origProf
		DefaultSharedConfigFiles = origFiles
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
			Err:    "failed to open shared config file, file_not_exist",
		},
		{
			LoadOptionFn: WithSharedConfigProfile("profile_not_exist"),
			Files: []string{
				filepath.Join("testdata", "shared_config"),
			},
			LoadFn: loadSharedConfig,
			Err:    "failed to get shared config profile, profile_not_exist",
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
			DefaultSharedConfigFiles = origFiles

			if len(c.Profile) > 0 {
				defaultSharedConfigProfile = c.Profile
			}
			if len(c.Files) > 0 {
				DefaultSharedConfigFiles = c.Files
			}

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
