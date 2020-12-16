package imds

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

func TestClientEndpoint(t *testing.T) {
	cases := map[string]struct {
		Endpoint       string
		EndpointEnvVar string
		Expect         string
	}{
		"default": {
			Expect: defaultEndpoint,
		},
		"from option": {
			Endpoint: "http://endpoint.localhost",
			Expect:   "http://endpoint.localhost",
		},
		"from option with environment": {
			Endpoint:       "http://endpoint.localhost",
			EndpointEnvVar: "http://[::1]",
			Expect:         "http://endpoint.localhost",
		},
		"from environment": {
			EndpointEnvVar: "http://[::1]",
			Expect:         "http://[::1]",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			envs := awstesting.StashEnv()
			defer awstesting.PopEnv(envs)

			if v := c.EndpointEnvVar; len(v) != 0 {
				os.Setenv(endpointEnvVar, v)
			}
			endpoint := c.Endpoint

			client := New(Options{
				disableAPIToken: true,
				Endpoint:        endpoint,
				HTTPClient: smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
					if e, a := c.Expect+getMetadataPath, r.URL.String(); e != a {
						return nil, fmt.Errorf("expect %v endpoint, got %v", e, a)
					}

					return newMockResponse(), nil
				}),
			})

			_, err := client.GetMetadata(context.Background(), nil)
			if err != nil {
				t.Fatal(err)
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

func newMockResponse() *http.Response {
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
	}
}
