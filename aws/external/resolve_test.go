package external

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/aws/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/aws/stscreds"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestResolveCustomCABundle(t *testing.T) {
	configs := Configs{
		WithCustomCABundle(awstesting.TLSBundleCA),
	}

	cfg := aws.Config{
		HTTPClient: &http.Client{Transport: &http.Transport{}},
	}

	if err := ResolveCustomCABundle(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	transport := cfg.HTTPClient.Transport.(*http.Transport)
	if transport.TLSClientConfig.RootCAs == nil {
		t.Errorf("expect root CAs set")
	}
}

func TestResolveRegion(t *testing.T) {
	configs := Configs{
		WithRegion("mock-region"),
		WithRegion("ignored-region"),
	}

	cfg := aws.Config{}
	cfg.Credentials = nil

	if err := ResolveRegion(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "mock-region", cfg.Region; e != a {
		t.Errorf("expect %v region, got %v", e, a)
	}
}

func TestResolveCredentialsValue(t *testing.T) {
	configs := Configs{
		WithCredentialsValue(aws.Credentials{
			Source: "invalid provider",
		}),
		WithCredentialsValue(aws.Credentials{
			AccessKeyID: "AKID", SecretAccessKey: "SECRET",
			Source: "valid",
		}),
		WithCredentialsValue(aws.Credentials{
			Source: "invalid provider 2",
		}),
	}

	cfg := aws.Config{}
	cfg.Credentials = nil

	if err := ResolveCredentialsValue(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.Credentials.(aws.StaticCredentialsProvider)
	if e, a := "AKID", p.Value.AccessKeyID; e != a {
		t.Errorf("expect %v key, got %v", e, a)
	}
	if e, a := "SECRET", p.Value.SecretAccessKey; e != a {
		t.Errorf("expect %v secret, got %v", e, a)
	}
	if e, a := "valid", p.Value.Source; e != a {
		t.Errorf("expect %v provider name, got %v", e, a)
	}

	creds, err := cfg.Credentials.Retrieve()
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "valid", creds.Source; e != a {
		t.Errorf("expect %v creds, got %v", e, a)
	}
}

func TestResolveEndpointCredentials(t *testing.T) {
	const u = "https://localhost/something"

	configs := Configs{
		WithCredentialsEndpoint(u),
	}

	cfg := unit.Config()
	cfg.Credentials = nil

	if err := ResolveEndpointCredentials(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.Credentials.(*endpointcreds.Provider)
	if p.Client == nil {
		t.Errorf("expect client set")
	}

	endpoint, err := p.Client.EndpointResolver.ResolveEndpoint(endpointcreds.ProviderName, "")
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := u, endpoint.URL; e != a {
		t.Errorf("expect %q endpoint, got %q", e, a)
	}
}

func TestResolveEndpointCredentials_ValidateEndpoint(t *testing.T) {
	configs := Configs{
		WithCredentialsEndpoint("http://notvalid.com"),
	}
	cfg := unit.Config()

	err := ResolveEndpointCredentials(&cfg, configs)
	if err == nil {
		t.Fatalf("expect error")
	}

	if e, a := "invalid endpoint", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q to be in %q", e, a)
	}
}

func TestResolveContainerEndpointPathCredentials(t *testing.T) {
	const u = "/some/path"

	configs := Configs{
		WithContainerCredentialsEndpointPath(u),
	}

	cfg := unit.Config()
	cfg.Credentials = nil

	if err := ResolveContainerEndpointPathCredentials(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.Credentials.(*endpointcreds.Provider)
	if p.Client == nil {
		t.Errorf("expect client set")
	}

	endpoint, err := p.Client.EndpointResolver.ResolveEndpoint(endpointcreds.ProviderName, "")
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	expect := containerCredentialsEndpoint + u
	if e, a := expect, endpoint.URL; e != a {
		t.Errorf("expect %q endpoint, got %q", e, a)
	}
}

func TestResolveAssumeRoleCredentials(t *testing.T) {
	configs := Configs{
		WithAssumeRoleConfig(AssumeRoleConfig{
			RoleARN:    "arn",
			ExternalID: "external",
			Source: &SharedConfig{
				Profile: "source",
				Credentials: aws.Credentials{
					AccessKeyID: "AKID", SecretAccessKey: "SECRET",
				},
			},
		}),
	}

	cfg := unit.Config()
	cfg.Credentials = nil

	if err := ResolveAssumeRoleCredentials(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.Credentials.(*stscreds.AssumeRoleProvider)
	if p.Client == nil {
		t.Errorf("expect client set")
	}
	if e, a := "arn", p.RoleARN; e != a {
		t.Errorf("expect %q arn, got %q", e, a)
	}
	if e, a := "external", *p.ExternalID; e != a {
		t.Errorf("expect %q external id, got %q", e, a)
	}
	if p.SerialNumber != nil {
		t.Errorf("expect no serial number")
	}
	if p.TokenProvider != nil {
		t.Errorf("expect no token provider")
	}
}

func TestResolveAssumeRoleCredentials_WithMFAToken(t *testing.T) {
	configs := Configs{
		WithAssumeRoleConfig(AssumeRoleConfig{
			RoleARN:    "arn",
			ExternalID: "external",
			MFASerial:  "abc123",
			Source: &SharedConfig{
				Profile: "source",
				Credentials: aws.Credentials{
					AccessKeyID: "AKID", SecretAccessKey: "SECRET",
				},
			},
		}),
		WithMFATokenFunc(func() (string, error) {
			return "token", nil
		}),
	}

	cfg := unit.Config()
	cfg.Credentials = nil

	if err := ResolveAssumeRoleCredentials(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.Credentials.(*stscreds.AssumeRoleProvider)
	if p.Client == nil {
		t.Errorf("expect client set")
	}
	if e, a := "arn", p.RoleARN; e != a {
		t.Errorf("expect %q arn, got %q", e, a)
	}
	if e, a := "external", *p.ExternalID; e != a {
		t.Errorf("expect %q external id, got %q", e, a)
	}
	if e, a := "abc123", *p.SerialNumber; e != a {
		t.Errorf("expect %q serial, got %q", e, a)
	}
	if p.TokenProvider == nil {
		t.Errorf("expect token provider")
	}
}

func TestResolveAssumeRoleCredentials_WithMFATokenError(t *testing.T) {
	configs := Configs{
		WithAssumeRoleConfig(AssumeRoleConfig{
			RoleARN:    "arn",
			ExternalID: "external",
			MFASerial:  "abc123",
			Source: &SharedConfig{
				Profile: "source",
				Credentials: aws.Credentials{
					AccessKeyID: "AKID", SecretAccessKey: "SECRET",
				},
			},
		}),
	}

	cfg := unit.Config()
	cfg.Credentials = nil

	err := ResolveAssumeRoleCredentials(&cfg, configs)
	if err == nil {
		t.Fatalf("expect error")
	}
	if e, a := "MFA", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q error in %q", e, a)
	}
	if cfg.Credentials != nil {
		t.Errorf("expect no credentials")
	}
}

func TestResolveFallbackEC2Credentials(t *testing.T) {
	configs := Configs{}

	cfg := unit.Config()

	if err := ResolveFallbackEC2Credentials(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if cfg.Credentials == nil {
		t.Errorf("expect credentials set")
	}

	p := cfg.Credentials.(*ec2rolecreds.Provider)
	if p.Client == nil {
		t.Errorf("expect client set")
	}
	if e, a := 5*time.Minute, p.ExpiryWindow; e != a {
		t.Errorf("expect %v expiry window, got %v", e, a)
	}
}
