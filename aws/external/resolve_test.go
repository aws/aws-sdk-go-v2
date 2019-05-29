package external

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
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

	cfg := defaults.Config()
	if err := ResolveCustomCABundle(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	type transportGetter interface {
		GetTransport() *http.Transport
	}

	trGetter := cfg.HTTPClient.(transportGetter)
	tr := trGetter.GetTransport()
	if tr.TLSClientConfig.RootCAs == nil {
		t.Errorf("expect root CAs set")
	}
}

func TestResolveCustomCABundle_ValidCA(t *testing.T) {
	certFile, keyFile, caFile, err := awstesting.CreateTLSBundleFiles()
	if err != nil {
		t.Fatalf("failed to create cert temp files, %v", err)
	}
	defer func() {
		awstesting.CleanupTLSBundleFiles(certFile, keyFile, caFile)
	}()

	serverAddr, err := awstesting.CreateTLSServer(certFile, keyFile, nil)
	if err != nil {
		t.Fatalf("failed to start TLS server, %v", err)
	}

	caPEM, err := ioutil.ReadFile(caFile)
	if err != nil {
		t.Fatalf("failed to read CA file, %v", err)
	}

	configs := Configs{
		WithCustomCABundle(caPEM),
	}

	cfg := defaults.Config()
	if err := ResolveCustomCABundle(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	req, _ := http.NewRequest("GET", serverAddr, nil)
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make request to TLS server, %v", err)
	}
	resp.Body.Close()

	if e, a := http.StatusOK, resp.StatusCode; e != a {
		t.Errorf("expect %v status, got %v", e, a)
	}
}

func TestResolveCustomCABundle_ErrorCustomClient(t *testing.T) {
	configs := Configs{
		WithCustomCABundle(awstesting.TLSBundleCA),
	}

	cfg := aws.Config{
		HTTPClient: &http.Client{},
	}
	if err := ResolveCustomCABundle(&cfg, configs); err == nil {
		t.Fatalf("expect error, got none")
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
	orgLookup := lookupHostFn
	defer func() {
		lookupHostFn = orgLookup
	}()

	cases := map[string]struct {
		LookupFn func(string) ([]string, error)
		Err      string
	}{
		"no addrs": {
			LookupFn: func(h string) ([]string, error) {
				return []string{}, nil
			},
			Err: "failed to resolve",
		},
		"lookup error": {
			LookupFn: func(h string) ([]string, error) {
				return []string{}, nil
			},
			Err: "failed to resolve",
		},
		"no local": {
			LookupFn: func(h string) ([]string, error) {
				return []string{"10.10.10.10"}, nil
			},
			Err: "failed to resolve",
		},
	}

	lookupHostFn = func(h string) ([]string, error) {
		return []string{}, nil
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			configs := Configs{
				WithCredentialsEndpoint("http://notvalid.com"),
			}
			cfg := unit.Config()

			err := ResolveEndpointCredentials(&cfg, configs)
			if err == nil {
				t.Fatalf("expect error")
			}

			if e, a := c.Err, err.Error(); !strings.Contains(a, e) {
				t.Errorf("expect %q to be in %q", e, a)
			}
		})
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
