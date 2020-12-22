package config

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
)

func TestAssumeRole(t *testing.T) {
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_CONFIG_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_w_creds")

	client := mockHTTPClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf(assumeRoleRespMsg,
				time.Now().Add(15*time.Minute).Format("2006-01-02T15:04:05Z"))))),
		}, nil
	})

	config, err := LoadDefaultConfig(context.Background(), WithHTTPClient(client))
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	creds, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SESSION_TOKEN", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "AssumeRoleProvider", creds.Source; !strings.Contains(a, e) {
		t.Errorf("expect %v, to be in %v", e, a)
	}
}

func TestAssumeRole_WithMFA(t *testing.T) {
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_CONFIG_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_w_creds")

	client := mockHTTPClient(func(r *http.Request) (*http.Response, error) {
		t.Helper()

		if e, a := r.FormValue("SerialNumber"), "0123456789"; e != a {
			t.Errorf("expect %v, got %v", e, a)
		}
		if e, a := r.FormValue("TokenCode"), "tokencode"; e != a {
			t.Errorf("expect %v, got %v", e, a)
		}
		if e, a := "900", r.FormValue("DurationSeconds"); e != a {
			t.Errorf("expect %v, got %v", e, a)
		}

		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf(assumeRoleRespMsg,
				time.Now().Add(15*time.Minute).Format("2006-01-02T15:04:05Z"))))),
		}, nil
	})

	customProviderCalled := false
	config, err := LoadDefaultConfig(context.Background(),
		WithHTTPClient(client),
		WithRegion("us-east-1"),
		WithSharedConfigProfile("assume_role_w_mfa"),
		WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
			options.TokenProvider = func() (string, error) {
				customProviderCalled = true
				return "tokencode", nil
			}
		}),
	)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	creds, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if !customProviderCalled {
		t.Errorf("expect true")
	}

	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SESSION_TOKEN", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "AssumeRoleProvider", creds.Source; !strings.Contains(a, e) {
		t.Errorf("expect %v, to be in %v", e, a)
	}
}

func TestAssumeRole_WithMFA_NoTokenProvider(t *testing.T) {
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_CONFIG_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_w_creds")

	_, err := LoadDefaultConfig(context.Background(), WithSharedConfigProfile("assume_role_w_mfa"))
	if e, a := (AssumeRoleTokenProviderNotSetError{}), err; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestAssumeRole_InvalidSourceProfile(t *testing.T) {
	// Backwards compatibility with Shared config disabled
	// assume role should not be built into the config.
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_CONFIG_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_invalid_source_profile")

	_, err := LoadDefaultConfig(context.Background())
	if err == nil {
		t.Fatalf("expect error, got none")
	}

	expectMsg := "failed to load assume role"
	if e, a := expectMsg, err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v, to be in %v", e, a)
	}
}

func TestAssumeRole_ExtendedDuration(t *testing.T) {
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_CONFIG_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_w_creds_ext_dur")

	client := mockHTTPClient(func(r *http.Request) (*http.Response, error) {
		t.Helper()

		if e, a := "1800", r.FormValue("DurationSeconds"); e != a {
			t.Errorf("expect %v, got %v", e, a)
		}

		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf(
				assumeRoleRespMsg,
				time.Now().Add(15*time.Minute).Format("2006-01-02T15:04:05Z"))))),
		}, nil
	})

	config, err := LoadDefaultConfig(context.Background(), WithHTTPClient(client))
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	creds, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SESSION_TOKEN", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "AssumeRoleProvider", creds.Source; !strings.Contains(a, e) {
		t.Errorf("expect %v, to be in %v", e, a)
	}
}
