package external

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-ini/ini"
)

var _ RegionProvider = (*SharedConfig)(nil)
var _ CredentialsValueProvider = (*SharedConfig)(nil)
var _ AssumeRoleConfigProvider = (*SharedConfig)(nil)

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
		{
			Filenames: []string{"file_not_exists"},
			Profile:   "default",
		},
		{
			Filenames: []string{testConfigFilename},
			Profile:   "default",
			Expected: SharedConfig{
				Profile: "default",
				Region:  "default_region",
			},
		},
		{
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "config_file_load_order",
			Expected: SharedConfig{
				Profile: "config_file_load_order",
				Region:  "shared_config_region",
				Credentials: aws.Credentials{
					AccessKeyID:     "shared_config_akid",
					SecretAccessKey: "shared_config_secret",
					ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
			},
		},
		{
			Filenames: []string{testConfigFilename, testConfigOtherFilename},
			Profile:   "config_file_load_order",
			Expected: SharedConfig{
				Profile: "config_file_load_order",
				Region:  "shared_config_other_region",
				Credentials: aws.Credentials{
					AccessKeyID:     "shared_config_other_akid",
					SecretAccessKey: "shared_config_other_secret",
					ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", testConfigOtherFilename),
				},
			},
		},
		{
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role",
			Expected: SharedConfig{
				Profile: "assume_role",
				AssumeRole: AssumeRoleConfig{
					RoleARN:       "assume_role_role_arn",
					sourceProfile: "complete_creds",
					Source: &SharedConfig{
						Profile: "complete_creds",
						Credentials: aws.Credentials{
							AccessKeyID:     "complete_creds_akid",
							SecretAccessKey: "complete_creds_secret",
							ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
						},
					},
				},
			},
		},
		{
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_invalid_source_profile",
			Expected: SharedConfig{
				Profile: "assume_role_invalid_source_profile",
				AssumeRole: AssumeRoleConfig{
					RoleARN:       "assume_role_invalid_source_profile_role_arn",
					sourceProfile: "profile_not_exists",
				},
			},
			Err: SharedConfigAssumeRoleError{RoleARN: "assume_role_invalid_source_profile_role_arn"},
		},
		{
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_w_creds",
			Expected: SharedConfig{
				Profile: "assume_role_w_creds",
				Credentials: aws.Credentials{
					AccessKeyID:     "assume_role_w_creds_akid",
					SecretAccessKey: "assume_role_w_creds_secret",
					ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
				AssumeRole: AssumeRoleConfig{
					RoleARN:         "assume_role_w_creds_role_arn",
					ExternalID:      "1234",
					RoleSessionName: "assume_role_w_creds_session_name",
					sourceProfile:   "assume_role_w_creds",
					Source: &SharedConfig{
						Profile: "assume_role_w_creds",
						Credentials: aws.Credentials{
							AccessKeyID:     "assume_role_w_creds_akid",
							SecretAccessKey: "assume_role_w_creds_secret",
							ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
						},
					},
				},
			},
		},
		{
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_wo_creds",
			Expected: SharedConfig{
				Profile: "assume_role_wo_creds",
				AssumeRole: AssumeRoleConfig{
					RoleARN:       "assume_role_wo_creds_role_arn",
					sourceProfile: "assume_role_wo_creds",
				},
			},
			Err: SharedConfigAssumeRoleError{RoleARN: "assume_role_wo_creds_role_arn"},
		},
		{
			Filenames: []string{filepath.Join("testdata", "shared_config_invalid_ini")},
			Profile:   "profile_name",
			Err:       SharedConfigLoadError{Filename: filepath.Join("testdata", "shared_config_invalid_ini")},
		},
	}

	for i, c := range cases {
		cfg, err := NewSharedConfig(c.Profile, c.Filenames)
		if c.Err != nil {
			if e, a := c.Err.Error(), err.Error(); !strings.Contains(a, e) {
				t.Errorf("expect %q to be in %q", e, a)
			}
			continue
		}

		if err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}
		if e, a := c.Expected, cfg; !reflect.DeepEqual(e, a) {
			t.Errorf("%d, expect %v, got %v", i, e, a)
		}
	}
}

func TestLoadSharedConfigFromFile(t *testing.T) {
	filename := testConfigFilename
	f, err := ini.Load(filename)
	if err != nil {
		t.Fatalf("failed to load test config file, %s, %v", filename, err)
	}
	iniFile := sharedConfigFile{IniData: f, Filename: filename}

	cases := []struct {
		Profile  string
		Expected SharedConfig
		Err      error
	}{
		{
			Profile:  "default",
			Expected: SharedConfig{Region: "default_region"},
		},
		{
			Profile:  "alt_profile_name",
			Expected: SharedConfig{Region: "alt_profile_name_region"},
		},
		{
			Profile:  "short_profile_name_first",
			Expected: SharedConfig{Region: "short_profile_name_first_short"},
		},
		{
			Profile:  "partial_creds",
			Expected: SharedConfig{},
		},
		{
			Profile: "complete_creds",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "complete_creds_akid",
					SecretAccessKey: "complete_creds_secret",
					ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
			},
		},
		{
			Profile: "complete_creds_with_token",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "complete_creds_with_token_akid",
					SecretAccessKey: "complete_creds_with_token_secret",
					SessionToken:    "complete_creds_with_token_token",
					ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
			},
		},
		{
			Profile: "full_profile",
			Expected: SharedConfig{
				Credentials: aws.Credentials{
					AccessKeyID:     "full_profile_akid",
					SecretAccessKey: "full_profile_secret",
					ProviderName:    fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
				Region: "full_profile_region",
			},
		},
		{
			Profile:  "partial_assume_role",
			Expected: SharedConfig{},
		},
		{
			Profile: "assume_role",
			Expected: SharedConfig{
				AssumeRole: AssumeRoleConfig{
					RoleARN:       "assume_role_role_arn",
					sourceProfile: "complete_creds",
				},
			},
		},
		{
			Profile: "assume_role_w_mfa",
			Expected: SharedConfig{
				AssumeRole: AssumeRoleConfig{
					RoleARN:       "assume_role_role_arn",
					sourceProfile: "complete_creds",
					MFASerial:     "0123456789",
				},
			},
		},
		{
			Profile: "does_not_exists",
			Err:     SharedConfigProfileNotExistsError{Profile: "does_not_exists"},
		},
	}

	for i, c := range cases {
		cfg := SharedConfig{}

		err := cfg.setFromIniFile(c.Profile, iniFile)
		if c.Err != nil {
			if e, a := c.Err.Error(), err.Error(); !strings.Contains(a, e) {
				t.Errorf("expect %q to be in %q", e, a)
			}
			continue
		}

		if err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}
		if e, a := c.Expected, cfg; !reflect.DeepEqual(e, a) {
			t.Errorf("%d, expect %v, got %v", i, e, a)
		}
	}
}

func TestLoadSharedConfigIniFiles(t *testing.T) {
	cases := []struct {
		Filenames []string
		Expected  []sharedConfigFile
	}{
		{
			Filenames: []string{"not_exists", testConfigFilename},
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
		files, err := loadSharedConfigIniFiles(c.Filenames)

		if err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}
		if e, a := len(c.Expected), len(files); e != a {
			t.Errorf("%d, expect %v, got %v", i, e, a)
		}
		if e, a := c.Expected, files; !cmpFiles(e, a) {
			t.Errorf("%d, expect %v, got %v", i, e, a)
		}
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
	t.Errorf("not tested")
}
