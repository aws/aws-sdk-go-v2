package config

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/awslabs/smithy-go/ptr"
)

var _ sharedConfigProfileProvider = (*EnvConfig)(nil)
var _ sharedConfigFilesProvider = (*EnvConfig)(nil)
var _ customCABundleProvider = (*EnvConfig)(nil)
var _ regionProvider = (*EnvConfig)(nil)

func TestNewEnvConfig_Creds(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

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
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	cases := []struct {
		Env    map[string]string
		Config EnvConfig
	}{
		0: {
			Env: map[string]string{
				"AWS_REGION":  "region",
				"AWS_PROFILE": "profile",
			},
			Config: EnvConfig{
				Region: "region", SharedConfigProfile: "profile",
			},
		},
		1: {
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
		2: {
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
		3: {
			Env: map[string]string{
				"AWS_DEFAULT_REGION":  "default_region",
				"AWS_DEFAULT_PROFILE": "default_profile",
			},
			Config: EnvConfig{
				Region: "default_region", SharedConfigProfile: "default_profile",
			},
		},
		4: {
			Env: map[string]string{
				"AWS_REGION":  "region",
				"AWS_PROFILE": "profile",
			},
			Config: EnvConfig{
				Region: "region", SharedConfigProfile: "profile",
			},
		},
		5: {
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
		6: {
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
		7: {
			Env: map[string]string{
				"AWS_DEFAULT_REGION":  "default_region",
				"AWS_DEFAULT_PROFILE": "default_profile",
			},
			Config: EnvConfig{
				Region: "default_region", SharedConfigProfile: "default_profile",
			},
		},
		8: {
			Env: map[string]string{
				"AWS_CA_BUNDLE": "custom_ca_bundle",
			},
			Config: EnvConfig{
				CustomCABundle: "custom_ca_bundle",
			},
		},
		9: {
			Env: map[string]string{
				"AWS_CA_BUNDLE": "custom_ca_bundle",
			},
			Config: EnvConfig{
				CustomCABundle: "custom_ca_bundle",
			},
		},
		10: {
			Env: map[string]string{
				"AWS_SHARED_CREDENTIALS_FILE": "/path/to/aws/file",
				"AWS_CONFIG_FILE":             "/path/to/config/file",
			},
			Config: EnvConfig{
				SharedCredentialsFile: "/path/to/aws/file",
				SharedConfigFile:      "/path/to/config/file",
			},
		},
		11: {
			Env: map[string]string{
				"AWS_S3_USE_ARN_REGION": "true",
			},
			Config: EnvConfig{
				S3UseARNRegion: ptr.Bool(true),
			},
		},
		12: {
			Env: map[string]string{
				"AWS_ENABLE_ENDPOINT_DISCOVERY": "true",
			},
			Config: EnvConfig{
				EnableEndpointDiscovery: ptr.Bool(true),
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
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
		})
	}
}

func TestSetEnvValue(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("empty_key", "")
	os.Setenv("second_key", "2")
	os.Setenv("third_key", "3")

	var dst string
	setStringFromEnvVal(&dst, []string{
		"empty_key", "first_key", "second_key", "third_key",
	})

	if e, a := "2", dst; e != a {
		t.Errorf("expect %s value from environment, got %s", e, a)
	}
}
