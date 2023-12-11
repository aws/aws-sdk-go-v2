package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// see https://github.com/aws/aws-sdk-go-v2/issues/2015
func TestResolveWebIdentityWithOptions(t *testing.T) {

	t.Run("token from env", func(t *testing.T) {
		restoreEnv := initConfigTestEnv()
		defer awstesting.PopEnv(restoreEnv)

		var tokenFile = filepath.Join("testdata", "wit.txt")
		os.Setenv("AWS_WEB_IDENTITY_TOKEN_FILE", tokenFile)
		os.Setenv("AWS_REGION", "us-east-1")

		_, err := LoadDefaultConfig(context.Background(),
			WithWebIdentityRoleCredentialOptions(func(options *stscreds.WebIdentityRoleOptions) {
				options.RoleARN = "test-arn"
			}),
		)

		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
	})

	t.Run("token from profile", func(t *testing.T) {
		// profile is still required to fully specify web identity properties for consistency with other SDKs/SEP
		restoreEnv := initConfigTestEnv()
		defer awstesting.PopEnv(restoreEnv)

		var configFileForWindows = filepath.Join("testdata", "config_source_shared_for_windows")
		var configFile = filepath.Join("testdata", "config_source_shared")

		os.Setenv("AWS_REGION", "us-east-1")
		os.Setenv("AWS_PROFILE", "webident-partial")

		if runtime.GOOS == "windows" {
			os.Setenv("AWS_CONFIG_FILE", configFileForWindows)
		} else {
			os.Setenv("AWS_CONFIG_FILE", configFile)
		}

		_, err := LoadDefaultConfig(context.Background())

		if err == nil || !strings.Contains(err.Error(), "web_identity_token_file requires role_arn") {
			t.Fatalf("expected profile parsing error, got %v", err)
		}
	})

	t.Run("token supplied directly from loadOptions", func(t *testing.T) {
		restoreEnv := initConfigTestEnv()
		defer awstesting.PopEnv(restoreEnv)

		var tokenFile = filepath.Join("testdata", "wit.txt")
		os.Setenv("AWS_REGION", "us-east-1")

		config, err := LoadDefaultConfig(context.Background(),
			WithEC2IMDSClientEnableState(imds.ClientDisabled),
			WithWebIdentityRoleCredentialOptions(func(options *stscreds.WebIdentityRoleOptions) {
				options.TokenRetriever = stscreds.IdentityTokenFile(tokenFile)
				options.RoleARN = "test-arn"
				options.RoleSessionName = "test-session"
			}),
		)

		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		target := stscreds.WebIdentityRoleProvider{}
		if !aws.IsCredentialsProvider(config.Credentials, &target) {
			t.Fatalf("expected type %T", target)
		}
	})
}
