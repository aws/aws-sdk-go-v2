//go:build integration
// +build integration

package s3

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// From the SEP:
// Service-specific endpoint configuration MUST be resolved with an endpoint URL provider chain with the following precedence:
//   - The value provided through code to an AWS SDK or tool via a command line
//     parameter or a client or configuration constructor; for example the
//     --endpoint-url command line parameter or the endpoint_url parameter
//     provided to the Python SDK client.
//   - The value provided by a service-specific environment variable.
//   - The value provided by the global endpoint environment variable
//     (AWS_ENDPOINT_URL).
//   - The value provided by a service-specific parameter from a services
//     definition section referenced in a profile in the shared configuration
//     file.
//   - The value provided by the global parameter from a profile in the shared
//     configuration file.
//   - The value resolved through the methods provided by the SDK or tool when
//     no explicit endpoint URL is provided.

func TestInteg_EndpointURL(t *testing.T) {
	for name, tt := range map[string]struct {
		Env          map[string]string
		SharedConfig string
		LoadOpts     []func(*config.LoadOptions) error
		ClientOpts   []func(*s3.Options)
		Expect       string
	}{
		"no values": {
			SharedConfig: `
[default]
`,
			Expect: "",
		},

		"precedence 0: in-code, set via s3.Options": {
			Env: map[string]string{
				"AWS_ENDPOINT_URL":    "https://global-env.com",
				"AWS_ENDPOINT_URL_S3": "https://service-env.com",
			},
			SharedConfig: `
[default]
endpoint_url = https://global-cfg.com
services = service_cfg

[services service_cfg]
s3 =
  endpoint_url = https://service-cfg.com
`,
			LoadOpts: []func(*config.LoadOptions) error{
				config.WithBaseEndpoint("https://loadopts.com"),
			},
			ClientOpts: []func(*s3.Options){
				func(o *s3.Options) {
					o.BaseEndpoint = aws.String("https://clientopts.com")
				},
			},
			Expect: "https://clientopts.com",
		},

		"precedence 0: in-code, set via config.LoadOptions": {
			Env: map[string]string{
				"AWS_ENDPOINT_URL":    "https://global-env.com",
				"AWS_ENDPOINT_URL_S3": "https://service-env.com",
			},
			SharedConfig: `
			[default]
			endpoint_url = https://global-cfg.com
			services = service_cfg

			[services service_cfg]
			s3 =
			  endpoint_url = https://service-cfg.com
			`,
			LoadOpts: []func(*config.LoadOptions) error{
				config.WithBaseEndpoint("https://loadopts.com"),
			},
			Expect: "https://loadopts.com",
		},

		"precedence 1: service env": {
			Env: map[string]string{
				"AWS_ENDPOINT_URL":    "https://global-env.com",
				"AWS_ENDPOINT_URL_S3": "https://service-env.com",
			},
			SharedConfig: `
[default]
endpoint_url = https://global-cfg.com
services = service_cfg

[services service_cfg]
s3 =
  endpoint_url = https://service-cfg.com
`,
			Expect: "https://service-env.com",
		},

		"precedence 2: global env": {
			Env: map[string]string{
				"AWS_ENDPOINT_URL": "https://global-env.com",
			},
			SharedConfig: `
[default]
endpoint_url = https://global-cfg.com
services = service_cfg

[services service_cfg]
s3 =
  endpoint_url = https://service-cfg.com
`,
			Expect: "https://global-env.com",
		},

		"precedence 3: service cfg": {
			SharedConfig: `
[default]
endpoint_url = https://global-cfg.com
services = service_cfg

[services service_cfg]
s3 =
  endpoint_url = https://service-cfg.com
`,
			Expect: "https://service-cfg.com",
		},

		"precedence 4: global cfg": {
			SharedConfig: `
[default]
endpoint_url = https://global-cfg.com
`,
			Expect: "https://global-cfg.com",
		},
	} {
		t.Run(name, func(t *testing.T) {
			reset, err := mockEnvironment(tt.Env, tt.SharedConfig)
			if err != nil {
				t.Fatalf("mock environment: %v", err)
			}
			defer reset()

			loadopts := append(tt.LoadOpts,
				config.WithSharedConfigFiles([]string{"test_shared_config"}),
				config.WithSharedConfigProfile("default"))

			cfg, err := config.LoadDefaultConfig(context.Background(), loadopts...)
			if err != nil {
				t.Fatalf("load config: %v", err)
			}

			svc := s3.NewFromConfig(cfg, tt.ClientOpts...)
			actual := aws.ToString(svc.Options().BaseEndpoint)
			if tt.Expect != actual {
				t.Errorf("expect endpoint: %q != %q", tt.Expect, actual)
			}
		})
	}
}

func mockEnvironment(env map[string]string, sharedCfg string) (func(), error) {
	for k, v := range env {
		os.Setenv(k, v)
	}
	f, err := os.Create("test_shared_config")
	if err != nil {
		return nil, err
	}
	if _, err := f.Write([]byte(sharedCfg)); err != nil {
		return nil, err
	}

	return func() {
		for k := range env {
			os.Unsetenv(k)
		}
		if err := os.Remove("test_shared_config"); err != nil {
			panic(err)
		}
	}, nil
}
