package external

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestResolveCustomCABundle(t *testing.T) {
	configs := Configs{
		WithCustomCABundle(awstesting.TLSBundleCA),
	}

	cfg := aws.Config{}
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

	cfg := aws.Config{}
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

func TestResolveCredentialsProvider(t *testing.T) {
	configs := Configs{
		WithCredentialsProvider{aws.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "AKID",
				SecretAccessKey: "SECRET",
				Source:          "valid",
			}},
		},
	}

	cfg := aws.Config{}
	cfg.Credentials = nil

	if found, err := ResolveCredentialProvider(&cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	} else if e, a := true, found; e != a {
		t.Fatalf("expected %v, got %v", e, a)
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

	creds, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "valid", creds.Source; e != a {
		t.Errorf("expect %v creds, got %v", e, a)
	}
}

func TestDefaultRegion(t *testing.T) {
	configs := Configs{
		WithDefaultRegion("foo-region"),
	}

	cfg := unit.Config()

	err := ResolveDefaultRegion(&cfg, configs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if e, a := "mock-region", cfg.Region; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	cfg.Region = ""

	err = ResolveDefaultRegion(&cfg, configs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if e, a := "foo-region", cfg.Region; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}
}
