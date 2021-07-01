package imds

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds/internal/config"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestClientEndpoint(t *testing.T) {
	cases := map[string]struct {
		Endpoint     string
		EndpointMode EndpointModeState
		EnvVar       map[string]string
		Expect       string
		WantErr      bool
	}{
		"default": {
			Expect: defaultIPv4Endpoint,
		},
		"options endpoint mode IPv4": {
			EndpointMode: EndpointModeStateIPv4,
			Expect:       defaultIPv4Endpoint,
		},
		"options endpoint mode IPv6": {
			EndpointMode: EndpointModeStateIPv6,
			Expect:       defaultIPv6Endpoint,
		},
		"options endpoint mode IPv6 AND options endpoint": {
			Endpoint:     "http://endpoint.localhost",
			EndpointMode: EndpointModeStateIPv6,
			Expect:       "http://endpoint.localhost",
		},
		"options endpoint mode IPv4 AND options endpoint": {
			Endpoint:     "http://endpoint.localhost",
			EndpointMode: EndpointModeStateIPv4,
			Expect:       "http://endpoint.localhost",
		},
		"options endpoint": {
			Endpoint: "http://endpoint.localhost",
			Expect:   "http://endpoint.localhost",
		},
		"options endpoint AND env endpoint": {
			Endpoint: "http://endpoint.localhost",
			EnvVar: map[string]string{
				endpointEnvVar: "http://[::1]",
			},
			Expect: "http://endpoint.localhost",
		},
		"env endpoint": {
			EnvVar: map[string]string{
				endpointEnvVar: "http://[::1]",
			},
			Expect: "http://[::1]",
		},
		"env endpoint missing scheme": {
			EnvVar: map[string]string{
				endpointEnvVar: "[::1]",
			},
			WantErr: true,
		},
		"options endpoint missing scheme": {
			Endpoint: "[::1]",
			WantErr:  true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			envs := awstesting.StashEnv()
			defer awstesting.PopEnv(envs)

			if v := c.EnvVar; len(v) != 0 {
				for k, v := range c.EnvVar {
					if err := os.Setenv(k, v); err != nil {
						t.Errorf("expect no error, got %v", err)
					}
				}
			}
			client := New(Options{
				disableAPIToken: true,
				Endpoint:        c.Endpoint,
				EndpointMode:    c.EndpointMode,
				HTTPClient: smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
					if e, a := c.Expect+getMetadataPath, r.URL.String(); e != a {
						return nil, fmt.Errorf("expect %v endpoint, got %v", e, a)
					}

					return newMockResponse(), nil
				}),
			})

			_, err := client.GetMetadata(context.Background(), nil)
			if (err != nil) != c.WantErr {
				t.Fatalf("WantErr=%v, got err=%v", c.WantErr, err)
			}
		})
	}
}

func TestClientEnableState(t *testing.T) {
	cases := map[string]struct {
		EnvironmentVar    string
		ClientEnableState ClientEnableState
		ExpectErr         string
	}{
		"default": {},
		"disabled environment var": {
			EnvironmentVar: "true",
			ExpectErr:      "access disabled",
		},
		"unknown environment value": {
			EnvironmentVar: "blah",
		},
		"disabled option": {
			ClientEnableState: ClientDisabled,
			ExpectErr:         "access disabled",
		},
		"enabled option": {
			EnvironmentVar:    "true",
			ClientEnableState: ClientEnabled,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			envs := awstesting.StashEnv()
			defer awstesting.PopEnv(envs)

			if v := c.EnvironmentVar; len(v) != 0 {
				os.Setenv(disableClientEnvVar, v)
			}

			client := New(Options{
				disableAPIToken:   true,
				ClientEnableState: c.ClientEnableState,
				HTTPClient: smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
					return newMockResponse(), nil
				}),
			})

			_, err := client.GetMetadata(context.Background(), nil)
			if len(c.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect error to contain %v, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

type WithEndpointModeSource EndpointModeState

func (w WithEndpointModeSource) GetEC2IMDSEndpointMode() (config.EndpointModeState, bool, error) {
	return config.EndpointModeState(w), true, nil
}

type WithEndpoint string

func (w WithEndpoint) GetEC2IMDSEndpoint() (string, bool, error) {
	return string(w), true, nil
}

func TestNewFromConfig(t *testing.T) {
	cases := map[string]struct {
		Sources []interface{}
		Expect  string
	}{
		"default": {
			Expect: defaultIPv4Endpoint,
		},
		"non-implementing sources": {
			Sources: []interface{}{
				struct{}{},
			},
			Expect: defaultIPv4Endpoint,
		},
		"endpoint mode IPv6": {
			Sources: []interface{}{
				WithEndpointModeSource(EndpointModeStateIPv6),
			},
			Expect: defaultIPv6Endpoint,
		},
		"endpoint mode IPv4": {
			Sources: []interface{}{
				WithEndpointModeSource(EndpointModeStateIPv4),
			},
			Expect: defaultIPv4Endpoint,
		},
		"endpoint mode unknown": {
			Sources: []interface{}{
				WithEndpointModeSource(func() (v EndpointModeState) {
					v.SetFromString("foobar")
					return v
				}()),
			},
			Expect: defaultIPv4Endpoint,
		},
		"endpoint": {
			Sources: []interface{}{
				WithEndpoint("http://endpoint.localhost"),
			},
			Expect: "http://endpoint.localhost",
		},
		"endpoint mode && endpoint": {
			Sources: []interface{}{
				WithEndpointModeSource(EndpointModeStateIPv6),
				WithEndpoint("http://endpoint.localhost"),
			},
			Expect: "http://endpoint.localhost",
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			envs := awstesting.StashEnv()
			defer awstesting.PopEnv(envs)

			client := NewFromConfig(aws.Config{
				ConfigSources: tt.Sources,
			}, func(options *Options) {
				options.disableAPIToken = true
				options.HTTPClient = smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
					if e, a := tt.Expect+getMetadataPath, r.URL.String(); e != a {
						return nil, fmt.Errorf("expect %v endpoint, got %v", e, a)
					}
					return newMockResponse(), nil
				})
			})

			_, err := client.GetMetadata(context.Background(), nil)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func newMockResponse() *http.Response {
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
	}
}
