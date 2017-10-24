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
	cfg.CredentialsLoader = nil

	if err := ResolveRegion(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "mock-region", aws.StringValue(cfg.Region); e != a {
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
	cfg.CredentialsLoader = nil

	if err := ResolveCredentialsValue(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.CredentialsLoader.Provider.(aws.StaticCredentialsProvider)
	if e, a := "AKID", p.Value.AccessKeyID; e != a {
		t.Errorf("expect %v key, got %v", e, a)
	}
	if e, a := "SECRET", p.Value.SecretAccessKey; e != a {
		t.Errorf("expect %v secret, got %v", e, a)
	}
	if e, a := "valid", p.Value.Source; e != a {
		t.Errorf("expect %v provider name, got %v", e, a)
	}

	creds, err := cfg.CredentialsLoader.Get()
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

	cfg := unit.Config.Copy()
	cfg.CredentialsLoader = nil

	if err := ResolveEndpointCredentials(cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.CredentialsLoader.Provider.(*endpointcreds.Provider)
	if p.Client == nil {
		t.Errorf("expect client set")
	}
	if e, a := u, p.Client.ClientInfo.Endpoint; e != a {
		t.Errorf("expect %q endpoint, got %q", e, a)
	}
}

func TestResolveEndpointCredentials_ValidateEndpoint(t *testing.T) {
	configs := Configs{
		WithCredentialsEndpoint("http://notvalid.com"),
	}
	cfg := unit.Config.Copy()

	err := ResolveEndpointCredentials(cfg, configs)
	if err == nil {
		t.Fatalf("expect error")
	}

	if e, a := "invalid endpoint", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q to be in %q", e, a)
	}
}

func TestValidateLocalEndpointURL(t *testing.T) {
	cases := []struct {
		URL   string
		IsErr bool
	}{
		{"http://127.0.0.1", false},
		{"http://localhost", false},
		{"http://invalid.tld", true},
		{"http://164.254.170.1", true},
	}

	for i, c := range cases {
		err := validateLocalEndpointURL(c.URL)
		if e, a := c.IsErr, err != nil; e != a {
			t.Errorf("%d, expect %t err, got %t", i, e, a)
		}
	}
}

func TestResolveContainerEndpointPathCredentials(t *testing.T) {
	const u = "/some/path"

	configs := Configs{
		WithContainerCredentialsEndpointPath(u),
	}

	cfg := unit.Config.Copy()
	cfg.CredentialsLoader = nil

	if err := ResolveContainerEndpointPathCredentials(cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.CredentialsLoader.Provider.(*endpointcreds.Provider)
	if p.Client == nil {
		t.Errorf("expect client set")
	}

	expect := containerCredentialsEndpoint + u
	if e, a := expect, p.Client.ClientInfo.Endpoint; e != a {
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

	cfg := unit.Config.Copy()
	cfg.CredentialsLoader = nil

	if err := ResolveAssumeRoleCredentials(cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.CredentialsLoader.Provider.(*stscreds.AssumeRoleProvider)
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

	cfg := unit.Config.Copy()
	cfg.CredentialsLoader = nil

	if err := ResolveAssumeRoleCredentials(cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	p := cfg.CredentialsLoader.Provider.(*stscreds.AssumeRoleProvider)
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

	cfg := unit.Config.Copy()
	cfg.CredentialsLoader = nil

	err := ResolveAssumeRoleCredentials(cfg, configs)
	if err == nil {
		t.Fatalf("expect error")
	}
	if e, a := "MFA", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q error in %q", e, a)
	}
	if cfg.CredentialsLoader != nil {
		t.Errorf("expect no credentials")
	}
}

func TestResolveFallbackEC2Credentials(t *testing.T) {
	configs := Configs{}

	cfg := unit.Config.Copy()

	if err := ResolveFallbackEC2Credentials(cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if cfg.CredentialsLoader == nil {
		t.Errorf("expect credentials set")
	}

	p := cfg.CredentialsLoader.Provider.(*ec2rolecreds.Provider)
	if p.Client == nil {
		t.Errorf("expect client set")
	}
	if e, a := 5*time.Minute, p.ExpiryWindow; e != a {
		t.Errorf("expect %v expiry window, got %v", e, a)
	}
}

// TODO use these tests for endpoint and assume role credentials
//func TestHTTPCredProvider(t *testing.T) {
//	cases := []struct {
//		Host string
//		Fail bool
//	}{
//		{"localhost", false}, {"127.0.0.1", false},
//		{"www.example.com", true}, {"169.254.170.2", true},
//	}
//
//	defer os.Clearenv()
//
//	for i, c := range cases {
//		u := fmt.Sprintf("http://%s/abc/123", c.Host)
//		os.Setenv(httpProviderEnvVar, u)
//
//		provider := RemoteCredProvider(aws.Config{}, aws.Handlers{})
//		if provider == nil {
//			t.Fatalf("%d, expect provider not to be nil, but was", i)
//		}
//
//		if c.Fail {
//			creds, err := provider.Retrieve()
//			if err == nil {
//				t.Fatalf("%d, expect error but got none", i)
//			} else {
//				aerr := err.(awserr.Error)
//				if e, a := "CredentialsEndpointError", aerr.Code(); e != a {
//					t.Errorf("%d, expect %s error code, got %s", i, e, a)
//				}
//			}
//			if e, a := endpointcreds.ProviderName, creds.ProviderName; e != a {
//				t.Errorf("%d, expect %s provider name got %s", i, e, a)
//			}
//		} else {
//			httpProvider := provider.(*endpointcreds.Provider)
//			if e, a := u, httpProvider.Client.Endpoint; e != a {
//				t.Errorf("%d, expect %q endpoint, got %q", i, e, a)
//			}
//		}
//	}
//}
//
//func TestECSCredProvider(t *testing.T) {
//	defer os.Clearenv()
//	os.Setenv(ecsCredsProviderEnvVar, "/abc/123")
//
//	provider := RemoteCredProvider(aws.Config{}, aws.Handlers{})
//	if provider == nil {
//		t.Fatalf("expect provider not to be nil, but was")
//	}
//
//	httpProvider := provider.(*endpointcreds.Provider)
//	if httpProvider == nil {
//		t.Fatalf("expect provider not to be nil, but was")
//	}
//	if e, a := "http://169.254.170.2/abc/123", httpProvider.Client.Endpoint; e != a {
//		t.Errorf("expect %q endpoint, got %q", e, a)
//	}
//}
//
//func TestDefaultEC2RoleProvider(t *testing.T) {
//	provider := RemoteCredProvider(aws.Config{}, aws.Handlers{})
//	if provider == nil {
//		t.Fatalf("expect provider not to be nil, but was")
//	}
//
//	ec2Provider := provider.(*ec2rolecreds.EC2RoleProvider)
//	if ec2Provider == nil {
//		t.Fatalf("expect provider not to be nil, but was")
//	}
//	if e, a := "http://169.254.169.254/latest", ec2Provider.Client.Endpoint; e != a {
//		t.Errorf("expect %q endpoint, got %q", e, a)
//	}
//}

// TODO integrate meaningful tests cases from shared config creds provider
//func TestSharedCredentialsProvider(t *testing.T) {
//	os.Clearenv()
//
//	p := SharedCredentialsProvider{Filename: "example.ini", Profile: ""}
//	creds, err := p.Retrieve()
//	assert.Nil(t, err, "Expect no error")
//
//	assert.Equal(t, "accessKey", creds.AccessKeyID, "Expect access key ID to match")
//	assert.Equal(t, "secret", creds.SecretAccessKey, "Expect secret access key to match")
//	assert.Equal(t, "token", creds.SessionToken, "Expect session token to match")
//}
//
//func TestSharedCredentialsProviderIsExpired(t *testing.T) {
//	os.Clearenv()
//
//	p := SharedCredentialsProvider{Filename: "example.ini", Profile: ""}
//
//	assert.True(t, p.IsExpired(), "Expect creds to be expired before retrieve")
//
//	_, err := p.Retrieve()
//	assert.Nil(t, err, "Expect no error")
//
//	assert.False(t, p.IsExpired(), "Expect creds to not be expired after retrieve")
//}
//
//func TestSharedCredentialsProviderWithAWS_SHARED_CREDENTIALS_FILE(t *testing.T) {
//	os.Clearenv()
//	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "example.ini")
//	p := SharedCredentialsProvider{}
//	creds, err := p.Retrieve()
//
//	assert.Nil(t, err, "Expect no error")
//
//	assert.Equal(t, "accessKey", creds.AccessKeyID, "Expect access key ID to match")
//	assert.Equal(t, "secret", creds.SecretAccessKey, "Expect secret access key to match")
//	assert.Equal(t, "token", creds.SessionToken, "Expect session token to match")
//}
//
//func TestSharedCredentialsProviderWithAWS_SHARED_CREDENTIALS_FILEAbsPath(t *testing.T) {
//	os.Clearenv()
//	wd, err := os.Getwd()
//	assert.NoError(t, err)
//	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", filepath.Join(wd, "example.ini"))
//	p := SharedCredentialsProvider{}
//	creds, err := p.Retrieve()
//	assert.Nil(t, err, "Expect no error")
//
//	assert.Equal(t, "accessKey", creds.AccessKeyID, "Expect access key ID to match")
//	assert.Equal(t, "secret", creds.SecretAccessKey, "Expect secret access key to match")
//	assert.Equal(t, "token", creds.SessionToken, "Expect session token to match")
//}
//
//func TestSharedCredentialsProviderWithAWS_PROFILE(t *testing.T) {
//	os.Clearenv()
//	os.Setenv("AWS_PROFILE", "no_token")
//
//	p := SharedCredentialsProvider{Filename: "example.ini", Profile: ""}
//	creds, err := p.Retrieve()
//	assert.Nil(t, err, "Expect no error")
//
//	assert.Equal(t, "accessKey", creds.AccessKeyID, "Expect access key ID to match")
//	assert.Equal(t, "secret", creds.SecretAccessKey, "Expect secret access key to match")
//	assert.Empty(t, creds.SessionToken, "Expect no token")
//}
//
//func TestSharedCredentialsProviderWithoutTokenFromProfile(t *testing.T) {
//	os.Clearenv()
//
//	p := SharedCredentialsProvider{Filename: "example.ini", Profile: "no_token"}
//	creds, err := p.Retrieve()
//	assert.Nil(t, err, "Expect no error")
//
//	assert.Equal(t, "accessKey", creds.AccessKeyID, "Expect access key ID to match")
//	assert.Equal(t, "secret", creds.SecretAccessKey, "Expect secret access key to match")
//	assert.Empty(t, creds.SessionToken, "Expect no token")
//}
//
//func TestSharedCredentialsProviderColonInCredFile(t *testing.T) {
//	os.Clearenv()
//
//	p := SharedCredentialsProvider{Filename: "example.ini", Profile: "with_colon"}
//	creds, err := p.Retrieve()
//	assert.Nil(t, err, "Expect no error")
//
//	assert.Equal(t, "accessKey", creds.AccessKeyID, "Expect access key ID to match")
//	assert.Equal(t, "secret", creds.SecretAccessKey, "Expect secret access key to match")
//	assert.Empty(t, creds.SessionToken, "Expect no token")
//}
//
//func TestSharedCredentialsProvider_DefaultFilename(t *testing.T) {
//	os.Clearenv()
//	os.Setenv("USERPROFILE", "profile_dir")
//	os.Setenv("HOME", "home_dir")
//
//	// default filename and profile
//	p := SharedCredentialsProvider{}
//
//	filename, err := p.filename()
//
//	if err != nil {
//		t.Fatalf("expect no error, got %v", err)
//	}
//
//	if e, a := shareddefaults.SharedCredentialsFilename(), filename; e != a {
//		t.Errorf("expect %q filename, got %q", e, a)
//	}
//}
//
//func BenchmarkSharedCredentialsProvider(b *testing.B) {
//	os.Clearenv()
//
//	p := SharedCredentialsProvider{Filename: "example.ini", Profile: ""}
//	_, err := p.Retrieve()
//	if err != nil {
//		b.Fatal(err)
//	}
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, err := p.Retrieve()
//		if err != nil {
//			b.Fatal(err)
//		}
//	}
//}
