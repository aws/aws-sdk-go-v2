package config

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/awslabs/smithy-go/logging"
)

func TestResolveCustomCABundle(t *testing.T) {
	var options LoadOptions
	var cfg aws.Config
	cfg.HTTPClient = awshttp.NewBuildableClient()

	WithCustomCABundle(bytes.NewReader(awstesting.TLSBundleCA))(&options)
	configs := configs{options}

	if err := resolveCustomCABundle(context.Background(), &cfg, configs); err != nil {
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

	var options LoadOptions
	var cfg aws.Config
	cfg.HTTPClient = awshttp.NewBuildableClient()

	WithCustomCABundle(bytes.NewReader(caPEM))(&options)
	configs := configs{options}

	if err := resolveCustomCABundle(context.Background(), &cfg, configs); err != nil {
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
	var options LoadOptions
	var cfg aws.Config

	cfg.HTTPClient = &http.Client{}

	WithCustomCABundle(bytes.NewReader(awstesting.TLSBundleCA))(&options)
	configs := configs{options}

	if err := resolveCustomCABundle(context.Background(), &cfg, configs); err == nil {
		t.Fatalf("expect error, got none")
	}
}

func TestResolveRegion(t *testing.T) {
	var options LoadOptions
	optFns := []func(options *LoadOptions) error{
		WithRegion("ignored-region"),

		WithRegion("mock-region"),
	}

	for _, optFn := range optFns {
		optFn(&options)
	}

	configs := configs{options}

	var cfg aws.Config

	if err := resolveRegion(context.Background(), &cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "mock-region", cfg.Region; e != a {
		t.Errorf("expect %v region, got %v", e, a)
	}
}

func TestResolveCredentialsProvider(t *testing.T) {
	var options LoadOptions
	optFns := []func(options *LoadOptions) error{
		WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "AKID",
				SecretAccessKey: "SECRET",
				Source:          "valid",
			}},
		),
	}

	for _, optFn := range optFns {
		optFn(&options)
	}

	configs := configs{options}

	var cfg aws.Config
	cfg.Credentials = nil

	if found, err := resolveCredentialProvider(context.Background(), &cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	} else if e, a := true, found; e != a {
		t.Fatalf("expected %v, got %v", e, a)
	}

	_, ok := cfg.Credentials.(*aws.CredentialsCache)
	if !ok {
		t.Fatalf("expect resolved credentials to be wrapped in cache, was not, %T", cfg.Credentials)
	}

	creds, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v key, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v secret, got %v", e, a)
	}
	if e, a := "valid", creds.Source; e != a {
		t.Errorf("expect %v provider name, got %v", e, a)
	}
}

func TestDefaultRegion(t *testing.T) {
	ctx := context.Background()

	var options LoadOptions
	WithDefaultRegion("foo-region")(&options)

	configs := configs{options}
	cfg := unit.Config()

	err := resolveDefaultRegion(ctx, &cfg, configs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if e, a := "mock-region", cfg.Region; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	cfg.Region = ""

	err = resolveDefaultRegion(ctx, &cfg, configs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if e, a := "foo-region", cfg.Region; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}
}

func TestResolveLogger(t *testing.T) {
	cfg, err := LoadDefaultConfig(context.Background(), func(o *LoadOptions) error {
		o.Logger = logging.Nop{}
		return nil
	})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	_, ok := cfg.Logger.(logging.Nop)
	if !ok {
		t.Error("unexpected logger type")
	}
}
