package config

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/smithy-go/logging"
	"github.com/google/go-cmp/cmp"
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

func TestResolveAppID(t *testing.T) {
	var options LoadOptions
	optFns := []func(options *LoadOptions) error{
		WithAppID("1234"),

		WithAppID("5678"),
	}

	for _, optFn := range optFns {
		optFn(&options)
	}

	configs := configs{options}

	var cfg aws.Config

	if err := resolveAppID(context.Background(), &cfg, configs); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "5678", cfg.AppID; e != a {
		t.Errorf("expect %v app ID, got %v", e, a)
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

func TestResolveDefaultsMode(t *testing.T) {
	cases := []struct {
		Mode                       aws.DefaultsMode
		ExpectedDefaultsMode       aws.DefaultsMode
		ExpectedRuntimeEnvironment aws.RuntimeEnvironment
		WithIMDS                   func() *httptest.Server
		Env                        map[string]string
	}{
		{
			ExpectedDefaultsMode: aws.DefaultsModeLegacy,
		},
		{
			Mode:                 aws.DefaultsModeStandard,
			ExpectedDefaultsMode: aws.DefaultsModeStandard,
		},
		{
			Mode:                 aws.DefaultsModeInRegion,
			ExpectedDefaultsMode: aws.DefaultsModeInRegion,
		},
		{
			Mode:                 aws.DefaultsModeCrossRegion,
			ExpectedDefaultsMode: aws.DefaultsModeCrossRegion,
		},
		{
			Mode:                 aws.DefaultsModeMobile,
			ExpectedDefaultsMode: aws.DefaultsModeMobile,
		},
		{
			Mode: aws.DefaultsModeAuto,
			Env: map[string]string{
				"AWS_EXECUTION_ENV": "envName",
				"AWS_REGION":        "us-west-2",
			},
			WithIMDS: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Path == "/latest/dynamic/instance-identity/document" {
							out, _ := json.Marshal(&imds.InstanceIdentityDocument{
								Region: "us-west-2",
							})
							w.Write(out)
						} else if r.URL.Path == "/latest/api/token" {
							header := w.Header()
							// bounce the TTL header
							const ttlHeader = "X-Aws-Ec2-Metadata-Token-Ttl-Seconds"
							header.Set(ttlHeader, r.Header.Get(ttlHeader))
							w.Write([]byte("validToken"))
						} else {
							w.Write([]byte(""))
						}
					}))
			},
			ExpectedDefaultsMode: aws.DefaultsModeAuto,
			ExpectedRuntimeEnvironment: aws.RuntimeEnvironment{
				EnvironmentIdentifier:     "envName",
				Region:                    "us-west-2",
				EC2InstanceMetadataRegion: "us-west-2",
			},
		},
		{
			Mode: aws.DefaultsModeAuto,
			Env: map[string]string{
				"AWS_EXECUTION_ENV": "envName",
				"AWS_REGION":        "us-west-2",
			},
			WithIMDS: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(500)
					}))
			},
			ExpectedDefaultsMode: aws.DefaultsModeAuto,
			ExpectedRuntimeEnvironment: aws.RuntimeEnvironment{
				EnvironmentIdentifier:     "envName",
				Region:                    "us-west-2",
				EC2InstanceMetadataRegion: "",
			},
		},
		{
			Mode: aws.DefaultsModeAuto,
			Env: map[string]string{
				"AWS_EXECUTION_ENV":         "envName",
				"AWS_REGION":                "us-west-2",
				"AWS_EC2_METADATA_DISABLED": "true",
			},
			ExpectedDefaultsMode: aws.DefaultsModeAuto,
			ExpectedRuntimeEnvironment: aws.RuntimeEnvironment{
				EnvironmentIdentifier:     "envName",
				Region:                    "us-west-2",
				EC2InstanceMetadataRegion: "",
			},
		},
		{
			Mode: aws.DefaultsModeAuto,
			Env: map[string]string{
				"AWS_REGION":                "us-west-2",
				"AWS_DEFAULT_REGION":        "other",
				"AWS_EC2_METADATA_DISABLED": "true",
			},
			ExpectedDefaultsMode: aws.DefaultsModeAuto,
			ExpectedRuntimeEnvironment: aws.RuntimeEnvironment{
				Region: "us-west-2",
			},
		},
		{
			Mode: aws.DefaultsModeAuto,
			Env: map[string]string{
				"AWS_DEFAULT_REGION":        "us-west-2",
				"AWS_EC2_METADATA_DISABLED": "true",
			},
			ExpectedDefaultsMode: aws.DefaultsModeAuto,
			ExpectedRuntimeEnvironment: aws.RuntimeEnvironment{
				Region: "us-west-2",
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var server *httptest.Server
			if tt.WithIMDS != nil {
				server = tt.WithIMDS()
				defer server.Close()
			}
			loadOptionsFunc := func(*LoadOptions) error {
				return nil
			}
			if len(tt.Mode) != 0 {
				loadOptionsFunc = WithDefaultsMode(tt.Mode, func(options *DefaultsModeOptions) {
					if server != nil {
						options.IMDSClient = imds.New(imds.Options{
							Endpoint: server.URL,
						})
					}
				})
			}

			if len(tt.Env) > 0 {
				restoreEnv := awstesting.StashEnv()
				defer awstesting.PopEnv(restoreEnv)

				for key := range tt.Env {
					_ = os.Setenv(key, tt.Env[key])
				}
			}

			cfg, err := LoadDefaultConfig(context.Background(), loadOptionsFunc)
			if err != nil {
				t.Errorf("expect no error, got %v", err)
			}

			if diff := cmp.Diff(tt.ExpectedDefaultsMode, cfg.DefaultsMode); len(diff) > 0 {
				t.Errorf(diff)
			}

			if diff := cmp.Diff(tt.ExpectedRuntimeEnvironment, cfg.RuntimeEnvironment); len(diff) > 0 {
				t.Errorf(diff)
			}
		})
	}
}
