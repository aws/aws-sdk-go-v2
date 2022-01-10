package config

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
)

func setupEnvForProcesscredsConfigFile() {
	filename := "proc_creds_config.ini"
	if runtime.GOOS == "windows" {
		filename = "proc_creds_config_win.ini"
	}

	os.Setenv("AWS_CONFIG_FILE", filepath.Join("testdata", filename))
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", filepath.Join("testdata", "empty_creds_config"))
}

func setupEnvForProcesscredsCredentialsFile() {
	filename := "proc_creds_credentials.ini"
	if runtime.GOOS == "windows" {
		filename = "proc_creds_credentials_win.ini"
	}

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", filepath.Join("testdata", filename))
}

func TestProcessCredentialsProvider_FromConfig(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	setupEnvForProcesscredsConfigFile()

	config, err := LoadDefaultConfig(context.Background(), WithRegion("region"))
	if err != nil {
		t.Fatalf("error loading default config: %v", err)
	}

	v, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "accessKey", v.AccessKeyID; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "secret", v.SecretAccessKey; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "tokenDefault", v.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessCredentialsProvider_FromConfigWithProfile(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_PROFILE", "not_expire")
	setupEnvForProcesscredsConfigFile()

	config, err := LoadDefaultConfig(context.Background(), WithRegion("region"))
	if err != nil {
		t.Fatalf("error loading default config: %v", err)
	}

	v, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "nonDefaultToken", v.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessCredentialsProvider_FromConfigWithStaticCreds(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_PROFILE", "not_alone")
	setupEnvForProcesscredsConfigFile()

	config, err := LoadDefaultConfig(context.Background(), WithRegion("region"))
	if err != nil {
		t.Fatalf("error loading default config: %v", err)
	}

	v, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "notFromCredProcAccess", v.AccessKeyID; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "notFromCredProcSecret", v.SecretAccessKey; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessCredentialsProvider_FromCredentials(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	setupEnvForProcesscredsCredentialsFile()

	config, err := LoadDefaultConfig(context.Background(), WithRegion("region"))
	if err != nil {
		t.Fatalf("error loading default config: %v", err)
	}

	v, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "accessKey", v.AccessKeyID; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "secret", v.SecretAccessKey; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "tokenDefault", v.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessCredentialsProvider_FromCredentialsWithProfile(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_PROFILE", "not_expire")
	setupEnvForProcesscredsCredentialsFile()

	config, err := LoadDefaultConfig(context.Background(), WithRegion("region"))
	if err != nil {
		t.Fatalf("error loading default config: %v", err)
	}

	v, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "nonDefaultToken", v.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessCredentialsProvider_FromCredentialsWithStaticCreds(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_PROFILE", "not_alone")
	setupEnvForProcesscredsCredentialsFile()

	config, err := LoadDefaultConfig(context.Background(), WithRegion("region"))
	if err != nil {
		t.Fatalf("error loading default config: %v", err)
	}

	v, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "notFromCredProcAccess", v.AccessKeyID; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "notFromCredProcSecret", v.SecretAccessKey; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}
}
