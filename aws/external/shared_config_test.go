package external

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
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
			Filenames: []string{"file_not_exist"},
			Profile:   "default",
			Err:       fmt.Errorf("failed to open shared config file, file_not_exist"),
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
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
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
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigOtherFilename),
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
							Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
						},
					},
				},
			},
		},
		{
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_invalid_source_profile",
			Err: SharedConfigAssumeRoleError{
				Profile: "assume_role_invalid_source_profile",
				RoleARN: "assume_role_invalid_source_profile_role_arn",
				Err: SharedConfigNotExistErrors{
					SharedConfigProfileNotExistError{
						Profile:  "profile_not_exists",
						Filename: testConfigOtherFilename,
						Err:      fmt.Errorf("section 'profile profile_not_exists' does not exist"),
					},
					SharedConfigProfileNotExistError{
						Profile:  "profile_not_exists",
						Filename: testConfigFilename,
						Err:      fmt.Errorf("section 'profile profile_not_exists' does not exist"),
					},
				},
			},
		},
		{
			Filenames: []string{testConfigOtherFilename, testConfigFilename},
			Profile:   "assume_role_w_creds",
			Expected: SharedConfig{
				Profile: "assume_role_w_creds",
				Credentials: aws.Credentials{
					AccessKeyID:     "assume_role_w_creds_akid",
					SecretAccessKey: "assume_role_w_creds_secret",
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
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
							Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
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
			Err: SharedConfigAssumeRoleError{
				Profile: "assume_role_wo_creds",
				RoleARN: "assume_role_wo_creds_role_arn",
				Err:     fmt.Errorf("source profile has no shared credentials"),
			},
		},
		{
			Filenames: []string{filepath.Join("testdata", "shared_config_invalid_ini")},
			Profile:   "profile_name",
			Err: SharedConfigLoadError{
				Filename: filepath.Join("testdata", "shared_config_invalid_ini"),
				Err:      fmt.Errorf("unclosed section: [profile_nam"),
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
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
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
					Source:          fmt.Sprintf("SharedConfigCredentials: %s", testConfigFilename),
				},
			},
		},
		{
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
			Profile: "does_not_exist",
			Err: SharedConfigProfileNotExistError{
				Filename: "testdata/shared_config",
				Profile:  "does_not_exist",
				Err:      fmt.Errorf("section 'profile does_not_exist' does not exist"),
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cfg := SharedConfig{}

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
	origProf := DefaultSharedConfigProfile
	origFiles := DefaultSharedConfigFiles
	defer func() {
		DefaultSharedConfigProfile = origProf
		DefaultSharedConfigFiles = origFiles
	}()

	cases := []struct {
		Configs Configs
		Files   []string
		Profile string
		LoadFn  func(Configs) (Config, error)
		Expect  SharedConfig
		Err     string
	}{
		{
			Configs: Configs{
				WithSharedConfigProfile("alt_profile_name"),
			},
			Files: []string{
				filepath.Join("testdata", "shared_config"),
			},
			LoadFn: LoadSharedConfig,
			Expect: SharedConfig{
				Profile: "alt_profile_name",
				Region:  "alt_profile_name_region",
			},
		},
		{
			Configs: Configs{
				WithSharedConfigFiles([]string{
					filepath.Join("testdata", "shared_config"),
				}),
			},
			Profile: "alt_profile_name",
			LoadFn:  LoadSharedConfig,
			Expect: SharedConfig{
				Profile: "alt_profile_name",
				Region:  "alt_profile_name_region",
			},
		},
		{
			Configs: Configs{
				WithSharedConfigProfile("default"),
			},
			Files: []string{
				filepath.Join("file_not_exist"),
			},
			LoadFn: LoadSharedConfig,
			Err:    "failed to open shared config file, file_not_exist",
		},
		{
			Configs: Configs{
				WithSharedConfigProfile("profile_not_exist"),
			},
			Files: []string{
				filepath.Join("testdata", "shared_config"),
			},
			LoadFn: LoadSharedConfig,
			Err:    "failed to get shared config profile, profile_not_exist",
		},
		{
			Configs: Configs{
				WithSharedConfigProfile("default"),
			},
			Files: []string{
				filepath.Join("file_not_exist"),
			},
			LoadFn: LoadSharedConfigIgnoreNotExist,
		},
		{
			Configs: Configs{
				WithSharedConfigProfile("assume_role_invalid_source_profile"),
			},
			Files: []string{
				testConfigOtherFilename, testConfigFilename,
			},
			LoadFn: LoadSharedConfig,
			Err:    "section 'profile profile_not_exists' does not exist",
		},
		{
			Configs: Configs{
				WithSharedConfigProfile("assume_role_invalid_source_profile"),
			},
			Files: []string{
				testConfigOtherFilename, testConfigFilename,
			},
			LoadFn: LoadSharedConfigIgnoreNotExist,
			Err:    "section 'profile profile_not_exists' does not exist",
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			DefaultSharedConfigProfile = origProf
			DefaultSharedConfigFiles = origFiles

			if len(c.Profile) > 0 {
				DefaultSharedConfigProfile = c.Profile
			}
			if len(c.Files) > 0 {
				DefaultSharedConfigFiles = c.Files
			}

			cfg, err := c.LoadFn(c.Configs)
			if len(c.Err) > 0 {
				if e, a := c.Err, err.Error(); !strings.Contains(a, e) {
					t.Errorf("expect %q to be in %q", e, a)
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
