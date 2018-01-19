package external

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
)

var _ SharedConfigProfileProvider = (*EnvConfig)(nil)
var _ SharedConfigFilesProvider = (*EnvConfig)(nil)
var _ CustomCABundleProvider = (*EnvConfig)(nil)
var _ RegionProvider = (*EnvConfig)(nil)
var _ CredentialsValueProvider = (*EnvConfig)(nil)
var _ CredentialsEndpointProvider = (*EnvConfig)(nil)
var _ ContainerCredentialsEndpointPathProvider = (*EnvConfig)(nil)

func TestNewEnvConfig_Creds(t *testing.T) {
	env := awstesting.StashEnv()
	defer awstesting.PopEnv(env)

	cases := []struct {
		Env map[string]string
		Val aws.Credentials
	}{
		{
			Env: map[string]string{
				"AWS_ACCESS_KEY": "AKID",
			},
			Val: aws.Credentials{},
		},
		{
			Env: map[string]string{
				"AWS_ACCESS_KEY_ID": "AKID",
			},
			Val: aws.Credentials{},
		},
		{
			Env: map[string]string{
				"AWS_SECRET_KEY": "SECRET",
			},
			Val: aws.Credentials{},
		},
		{
			Env: map[string]string{
				"AWS_SECRET_ACCESS_KEY": "SECRET",
			},
			Val: aws.Credentials{},
		},
		{
			Env: map[string]string{
				"AWS_ACCESS_KEY_ID":     "AKID",
				"AWS_SECRET_ACCESS_KEY": "SECRET",
			},
			Val: aws.Credentials{
				AccessKeyID: "AKID", SecretAccessKey: "SECRET",
				Source: CredentialsSourceName,
			},
		},
		{
			Env: map[string]string{
				"AWS_ACCESS_KEY": "AKID",
				"AWS_SECRET_KEY": "SECRET",
			},
			Val: aws.Credentials{
				AccessKeyID: "AKID", SecretAccessKey: "SECRET",
				Source: CredentialsSourceName,
			},
		},
		{
			Env: map[string]string{
				"AWS_ACCESS_KEY":    "AKID",
				"AWS_SECRET_KEY":    "SECRET",
				"AWS_SESSION_TOKEN": "TOKEN",
			},
			Val: aws.Credentials{
				AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "TOKEN",
				Source: CredentialsSourceName,
			},
		},
	}

	for i, c := range cases {
		os.Clearenv()

		for k, v := range c.Env {
			os.Setenv(k, v)
		}

		cfg, err := NewEnvConfig()
		if err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}

		if !reflect.DeepEqual(c.Val, cfg.Credentials) {
			t.Errorf("%d, expect aws to match.\n%s", i,
				awstesting.SprintExpectActual(c.Val, cfg.Credentials))
		}
	}
}

func TestNewEnvConfig(t *testing.T) {
	env := awstesting.StashEnv()
	defer awstesting.PopEnv(env)

	cases := []struct {
		Env    map[string]string
		Config EnvConfig
	}{
		{
			Env: map[string]string{
				"AWS_REGION":  "region",
				"AWS_PROFILE": "profile",
			},
			Config: EnvConfig{
				Region: "region", SharedConfigProfile: "profile",
			},
		},
		{
			Env: map[string]string{
				"AWS_REGION":          "region",
				"AWS_DEFAULT_REGION":  "default_region",
				"AWS_PROFILE":         "profile",
				"AWS_DEFAULT_PROFILE": "default_profile",
			},
			Config: EnvConfig{
				Region: "region", SharedConfigProfile: "profile",
			},
		},
		{
			Env: map[string]string{
				"AWS_REGION":          "region",
				"AWS_DEFAULT_REGION":  "default_region",
				"AWS_PROFILE":         "profile",
				"AWS_DEFAULT_PROFILE": "default_profile",
			},
			Config: EnvConfig{
				Region: "region", SharedConfigProfile: "profile",
			},
		},
		{
			Env: map[string]string{
				"AWS_DEFAULT_REGION":  "default_region",
				"AWS_DEFAULT_PROFILE": "default_profile",
			},
			Config: EnvConfig{
				Region: "default_region", SharedConfigProfile: "default_profile",
			},
		},
		{
			Env: map[string]string{
				"AWS_REGION":  "region",
				"AWS_PROFILE": "profile",
			},
			Config: EnvConfig{
				Region: "region", SharedConfigProfile: "profile",
			},
		},
		{
			Env: map[string]string{
				"AWS_REGION":          "region",
				"AWS_DEFAULT_REGION":  "default_region",
				"AWS_PROFILE":         "profile",
				"AWS_DEFAULT_PROFILE": "default_profile",
			},
			Config: EnvConfig{
				Region: "region", SharedConfigProfile: "profile",
			},
		},
		{
			Env: map[string]string{
				"AWS_REGION":          "region",
				"AWS_DEFAULT_REGION":  "default_region",
				"AWS_PROFILE":         "profile",
				"AWS_DEFAULT_PROFILE": "default_profile",
			},
			Config: EnvConfig{
				Region: "region", SharedConfigProfile: "profile",
			},
		},
		{
			Env: map[string]string{
				"AWS_DEFAULT_REGION":  "default_region",
				"AWS_DEFAULT_PROFILE": "default_profile",
			},
			Config: EnvConfig{
				Region: "default_region", SharedConfigProfile: "default_profile",
			},
		},
		{
			Env: map[string]string{
				"AWS_CA_BUNDLE": "custom_ca_bundle",
			},
			Config: EnvConfig{
				CustomCABundle: "custom_ca_bundle",
			},
		},
		{
			Env: map[string]string{
				"AWS_CA_BUNDLE": "custom_ca_bundle",
			},
			Config: EnvConfig{
				CustomCABundle: "custom_ca_bundle",
			},
		},
		{
			Env: map[string]string{
				"AWS_SHARED_CREDENTIALS_FILE": "/path/to/aws/file",
				"AWS_CONFIG_FILE":             "/path/to/config/file",
			},
			Config: EnvConfig{
				SharedCredentialsFile: "/path/to/aws/file",
				SharedConfigFile:      "/path/to/config/file",
			},
		},
	}

	for _, c := range cases {
		os.Clearenv()

		for k, v := range c.Env {
			os.Setenv(k, v)
		}

		cfg, err := NewEnvConfig()
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		if !reflect.DeepEqual(c.Config, cfg) {
			t.Errorf("expect config to match.\n%s",
				awstesting.SprintExpectActual(c.Config, cfg))
		}
	}
}

func TestSetEnvValue(t *testing.T) {
	env := awstesting.StashEnv()
	defer awstesting.PopEnv(env)

	os.Setenv("empty_key", "")
	os.Setenv("second_key", "2")
	os.Setenv("third_key", "3")

	var dst string
	setFromEnvVal(&dst, []string{
		"empty_key", "first_key", "second_key", "third_key",
	})

	if e, a := "2", dst; e != a {
		t.Errorf("expect %s value from environment, got %s", e, a)
	}
}
