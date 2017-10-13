package external

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/credentials"
)

func TestResolveCABundle(t *testing.T) {
	t.Errorf("not implemented")
}

func TestResolveRegion(t *testing.T) {
	configs := Configs{
		WithRegion("mock-region"),
		WithRegion("ignored-region"),
	}
	cfg := aws.Config{}

	if err := ResolveRegion(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "mock-region", aws.StringValue(cfg.Region); e != a {
		t.Errorf("expect %v region, got %v", e, a)
	}
}

func TestResolveCredentialsValue(t *testing.T) {
	configs := Configs{
		WithCredentialsValue(credentials.Value{
			ProviderName: "invalid provider",
		}),
		WithCredentialsValue(credentials.Value{
			AccessKeyID: "AKID", SecretAccessKey: "SECRET",
			ProviderName: "valid",
		}),
		WithCredentialsValue(credentials.Value{
			ProviderName: "invalid provider 2",
		}),
	}
	cfg := aws.Config{}

	if err := ResolveCredentialsValue(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	creds, err := cfg.Credentials.Get()
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "valid", creds.ProviderName; e != a {
		t.Errorf("expect %v creds, got %v", e, a)
	}
}

func TestResolveEndpointCredentilas(t *testing.T) {
	t.Errorf("not implemented")
}

func TestResolveAssumeRoleCredentilas(t *testing.T) {
	t.Errorf("not implemented")
}

func TestResolveFallbackEC2Credentials(t *testing.T) {
	configs := Configs{}
	cfg := aws.Config{}

	if err := ResolveFallbackEC2Credentials(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if cfg.Credentials == nil {
		t.Errorf("expect credentials set")
	}
}
