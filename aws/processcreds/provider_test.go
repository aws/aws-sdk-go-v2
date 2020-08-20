package processcreds_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/aws/processcreds"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
)

func TestProcessProviderFromSessionCfg(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_CONFIG_FILE", "testdata\\shconfig_win.ini")
	} else {
		os.Setenv("AWS_CONFIG_FILE", "testdata/shconfig.ini")
	}

	config, err := external.LoadDefaultAWSConfig(external.WithRegion("region"))
	if err != nil {
		t.Errorf("error loading default config: %v", err)
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

func TestProcessProviderFromSessionWithProfileCfg(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_PROFILE", "non_expire")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_CONFIG_FILE", "testdata\\shconfig_win.ini")
	} else {
		os.Setenv("AWS_CONFIG_FILE", "testdata/shconfig.ini")
	}

	config, err := external.LoadDefaultAWSConfig(external.WithRegion("region"))
	if err != nil {
		t.Errorf("error loading default config: %v", err)
	}

	v, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "nonDefaultToken", v.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessProviderNotFromCredProcCfg(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_PROFILE", "not_alone")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_CONFIG_FILE", "testdata\\shconfig_win.ini")
	} else {
		os.Setenv("AWS_CONFIG_FILE", "testdata/shconfig.ini")
	}

	config, err := external.LoadDefaultAWSConfig(external.WithRegion("region"))
	if err != nil {
		t.Errorf("error loading default config: %v", err)
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

func TestProcessProviderFromSessionCred(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	if runtime.GOOS == "windows" {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata\\shcred_win.ini")
	} else {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata/shcred.ini")
	}

	config, err := external.LoadDefaultAWSConfig(external.WithRegion("region"))
	if err != nil {
		t.Errorf("error loading default config: %v", err)
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

func TestProcessProviderFromSessionWithProfileCred(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_PROFILE", "non_expire")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata\\shcred_win.ini")
	} else {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata/shcred.ini")
	}

	config, err := external.LoadDefaultAWSConfig(external.WithRegion("region"))
	if err != nil {
		t.Errorf("error loading default config: %v", err)
	}

	v, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "nonDefaultToken", v.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessProviderNotFromCredProcCrd(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_PROFILE", "not_alone")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata\\shcred_win.ini")
	} else {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata/shcred.ini")
	}

	config, err := external.LoadDefaultAWSConfig(external.WithRegion("region"))
	if err != nil {
		t.Errorf("error loading default config: %v", err)
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

func TestProcessProviderBadCommand(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	creds := processcreds.NewProvider("/bad/process")
	_, err := creds.Retrieve(context.Background())
	var pe *processcreds.ProviderError
	if ok := errors.As(err, &pe); !ok {
		t.Fatalf("expect error to be of type %T", pe)
	}
	if e, a := "error in credential_process", pe.Error(); !strings.Contains(a, e) {
		t.Errorf("expected %v, got %v", e, a)
	}
}

func TestProcessProviderMoreEmptyCommands(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	provider := processcreds.NewProvider("")
	_, err := provider.Retrieve(context.Background())
	var pe *processcreds.ProviderError
	if ok := errors.As(err, &pe); !ok {
		t.Fatalf("expect error to be of type %T", pe)
	}
	if e, a := "failed to prepare command", pe.Error(); !strings.Contains(a, e) {
		t.Errorf("expected %v, got %v", e, a)
	}
}

func TestProcessProviderExpectErrors(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	provider := processcreds.NewProvider(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "malformed.json"},
				string(os.PathSeparator))))
	_, err := provider.Retrieve(context.Background())
	var pe *processcreds.ProviderError
	if ok := errors.As(err, &pe); !ok {
		t.Fatalf("expect error to be of type %T", pe)
	}
	if e, a := "parse failed of credential_process output", pe.Error(); !strings.Contains(a, e) {
		t.Errorf("expected %v, got %v", e, a)
	}

	provider = processcreds.NewProvider(
		fmt.Sprintf("%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "wrongversion.json"},
				string(os.PathSeparator))))
	_, err = provider.Retrieve(context.Background())
	if ok := errors.As(err, &pe); !ok {
		t.Fatalf("expect error to be of type %T", pe)
	}
	if e, a := "wrong version in process output", pe.Error(); !strings.Contains(a, e) {
		t.Errorf("expected %v, got %v", e, a)
	}

	provider = processcreds.NewProvider(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "missingkey.json"},
				string(os.PathSeparator))))
	_, err = provider.Retrieve(context.Background())
	if ok := errors.As(err, &pe); !ok {
		t.Fatalf("expect error to be of type %T", pe)
	}
	if e, a := "missing AccessKeyId", pe.Error(); !strings.Contains(a, e) {
		t.Errorf("expected %v, got %v", e, a)
	}

	provider = processcreds.NewProvider(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "missingsecret.json"},
				string(os.PathSeparator))))
	_, err = provider.Retrieve(context.Background())
	if ok := errors.As(err, &pe); !ok {
		t.Fatalf("expect error to be of type %T", pe)
	}
	if e, a := "missing SecretAccessKey", pe.Error(); !strings.Contains(a, e) {
		t.Errorf("expected %v, got %v", e, a)
	}
}

func TestProcessProviderTimeout(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	command := "/bin/sleep 2"
	if runtime.GOOS == "windows" {
		// "timeout" command does not work due to pipe redirection
		command = "ping -n 2 127.0.0.1>nul"
	}

	provider := processcreds.NewProvider(command, func(options *processcreds.ProviderOptions) {
		options.Timeout = time.Duration(1) * time.Second
	})
	_, err := provider.Retrieve(context.Background())
	var pe *processcreds.ProviderError
	if ok := errors.As(err, &pe); !ok {
		t.Fatalf("expect error to be of type %T", pe)
	}
	if e, a := "credential process timed out", pe.Error(); !strings.Contains(a, e) {
		t.Errorf("expected %v, got %v", e, a)
	}
}

func TestProcessProviderWithLongSessionToken(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	provider := processcreds.NewProvider(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "longsessiontoken.json"},
				string(os.PathSeparator))))
	v, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}

	// Text string same length as session token returned by AWS for AssumeRoleWithWebIdentity
	e := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	if a := v.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}
}

type credentialTest struct {
	Version         int
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string
	Expiration      string
}

func TestProcessProviderStatic(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	// static
	provider := processcreds.NewProvider(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "static.json"},
				string(os.PathSeparator))))
	v, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if v.CanExpire != false {
		t.Errorf("expected %v, got %v", "static credentials/not expired", "can expire")
	}

}

func TestProcessProviderNotExpired(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	// non-static, not expired
	exp := &credentialTest{}
	exp.Version = 1
	exp.AccessKeyID = "accesskey"
	exp.SecretAccessKey = "secretkey"
	exp.Expiration = time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)
	b, err := json.Marshal(exp)
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmp_expiring")
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if _, err = io.Copy(tmpFile, bytes.NewReader(b)); err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	defer func() {
		if err = tmpFile.Close(); err != nil {
			t.Errorf("expected %v, got %v", "no error", err)
		}
		if err = os.Remove(tmpFile.Name()); err != nil {
			t.Errorf("expected %v, got %v", "no error", err)
		}
	}()
	provider := processcreds.NewProvider(
		fmt.Sprintf("%s %s", getOSCat(), tmpFile.Name()))
	v, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if v.Expired() {
		t.Errorf("expected %v, got %v", "not expired", "expired")
	}
}

func TestProcessProviderExpired(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	// non-static, expired
	exp := &credentialTest{}
	exp.Version = 1
	exp.AccessKeyID = "accesskey"
	exp.SecretAccessKey = "secretkey"
	exp.Expiration = time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	b, err := json.Marshal(exp)
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmp_expired")
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if _, err = io.Copy(tmpFile, bytes.NewReader(b)); err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	defer func() {
		if err = tmpFile.Close(); err != nil {
			t.Errorf("expected %v, got %v", "no error", err)
		}
		if err = os.Remove(tmpFile.Name()); err != nil {
			t.Errorf("expected %v, got %v", "no error", err)
		}
	}()
	provider := processcreds.NewProvider(
		fmt.Sprintf("%s %s", getOSCat(), tmpFile.Name()))
	v, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if !v.Expired() {
		t.Errorf("expected %v, got %v", "expired", "not expired")
	}
}

func TestProcessProviderForceExpire(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	// non-static, not expired

	// setup test credentials file
	exp := &credentialTest{}
	exp.Version = 1
	exp.AccessKeyID = "accesskey"
	exp.SecretAccessKey = "secretkey"
	exp.Expiration = time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)
	b, err := json.Marshal(exp)
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmp_force_expire")
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if _, err = io.Copy(tmpFile, bytes.NewReader(b)); err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	defer func() {
		if err = tmpFile.Close(); err != nil {
			t.Errorf("expected %v, got %v", "no error", err)
		}
		if err = os.Remove(tmpFile.Name()); err != nil {
			t.Errorf("expected %v, got %v", "no error", err)
		}
	}()

	// get credentials from file
	provider := processcreds.NewProvider(
		fmt.Sprintf("%s %s", getOSCat(), tmpFile.Name()))
	v, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if v.Expired() {
		t.Errorf("expected %v, got %v", "not expired", "expired")
	}

	// force expire creds
	provider.Invalidate()

	// renew creds
	v, err = provider.Retrieve(context.Background())
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if v.Expired() {
		t.Errorf("expected %v, got %v", "not expired", "expired")
	}
}

func TestProcessProviderAltConstruct(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	// constructing with exec.Cmd instead of string
	myCommand := exec.Command(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "static.json"},
				string(os.PathSeparator))))
	provider := processcreds.NewProviderCommand(myCommand, func(options *processcreds.ProviderOptions) {
		options.Timeout = time.Duration(1) * time.Second
	})
	v, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if v.CanExpire != false {
		t.Errorf("expected %v, got %v", "static credentials/not expired", "expired")
	}
}

func BenchmarkProcessProvider(b *testing.B) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	provider := processcreds.NewProvider(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "static.json"},
				string(os.PathSeparator))))
	_, err := provider.Retrieve(context.Background())
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err := provider.Retrieve(context.Background())
		if err != nil {
			b.Fatal(err)
		}
		b.StopTimer()
		provider.Invalidate()
	}
}

func getOSCat() string {
	if runtime.GOOS == "windows" {
		return "type"
	}
	return "cat"
}
