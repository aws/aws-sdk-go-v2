package config

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

func TestSharedConfigCredentialSource(t *testing.T) {
	const configFileForWindows = "testdata/credential_source_config_for_windows"
	const configFile = "testdata/credential_source_config"

	cases := map[string]struct {
		name              string
		envProfile        string
		configProfile     string
		expectedError     error
		expectedAccessKey string
		expectedSecretKey string
		expectedChain     []string
		init              func()
		dependentOnOS     bool
		client            HTTPClient
	}{
		"credential source and source profile": {
			envProfile:    "invalid_source_and_credential_source",
			expectedError: fmt.Errorf("TODO"),
			init: func() {
				os.Setenv("AWS_ACCESS_KEY", "access_key")
				os.Setenv("AWS_SECRET_KEY", "secret_key")
			},
		},
		"env var credential source": {
			configProfile:     "env_var_credential_source",
			expectedAccessKey: "AKID",
			expectedSecretKey: "SECRET",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_env",
			},
			init: func() {
				os.Setenv("AWS_ACCESS_KEY", "access_key")
				os.Setenv("AWS_SECRET_KEY", "secret_key")
			},
		},
		"ec2metadata credential source": {
			envProfile: "ec2metadata",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ec2",
			},
			expectedAccessKey: "AKID",
			expectedSecretKey: "SECRET",
			client: mockHttpClient(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader(nil))}, nil
			}),
		},
		"ecs container credential source": {
			envProfile:        "ecscontainer",
			expectedAccessKey: "AKID",
			expectedSecretKey: "SECRET",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ecs",
			},
			init: func() {
				os.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "/ECS")
			},
		},
		"chained assume role with env creds": {
			envProfile:        "chained_assume_role",
			expectedAccessKey: "AKID",
			expectedSecretKey: "SECRET",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_chain",
				"assume_role_w_creds_role_arn_ec2",
			},
		},
		"credential process with no ARN set": {
			envProfile:        "cred_proc_no_arn_set",
			dependentOnOS:     true,
			expectedAccessKey: "cred_proc_akid",
			expectedSecretKey: "cred_proc_secret",
		},
		"credential process with ARN set": {
			envProfile:        "cred_proc_arn_set",
			dependentOnOS:     true,
			expectedAccessKey: "AKID",
			expectedSecretKey: "SECRET",
			expectedChain: []string{
				"assume_role_w_creds_proc_role_arn",
			},
		},
		"chained assume role with credential process": {
			envProfile:        "chained_cred_proc",
			dependentOnOS:     true,
			expectedAccessKey: "AKID",
			expectedSecretKey: "SECRET",
			expectedChain: []string{
				"assume_role_w_creds_proc_source_prof",
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			restoreEnv := awstesting.StashEnv()
			defer awstesting.PopEnv(restoreEnv)

			if c.dependentOnOS && runtime.GOOS == "windows" {
				os.Setenv("AWS_CONFIG_FILE", configFileForWindows)
			} else {
				os.Setenv("AWS_CONFIG_FILE", configFile)
			}

			os.Setenv("AWS_REGION", "us-east-1")
			if len(c.envProfile) != 0 {
				os.Setenv("AWS_PROFILE", c.envProfile)
			}

			if c.init != nil {
				c.init()
			}

			var credChain []string

			const (
				ECS = iota
				AssumeRole
			)
			var apiExecuting int

			configSources := []Config{
				WithAPIOptions(append([]func(*middleware.Stack) error{}, func(stack *middleware.Stack) error {
					return stack.Initialize.Add(middleware.InitializeMiddlewareFunc("GetRoleArns", func(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
					) (
						out middleware.InitializeOutput, metadata middleware.Metadata, err error,
					) {
						switch v := in.Parameters.(type) {
						case *sts.AssumeRoleInput:
							credChain = append(credChain, *v.RoleArn)
							apiExecuting = AssumeRole
						default:
							apiExecuting = ECS
						}

						return next.HandleInitialize(ctx, in)
					}), middleware.After)
				}, func(stack *middleware.Stack) error {
					return stack.Deserialize.Add(middleware.DeserializeMiddlewareFunc("Response", func(ctx context.Context, input middleware.DeserializeInput, next middleware.DeserializeHandler) (
						out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
					) {
						if apiExecuting == ECS {
							var response []byte
							if input.Request.(*http.Request).URL.Path == "/ECS" {
								response = []byte(ecsResponse)
							} else {
								response = []byte("")
							}
							out.RawResponse = &smithyhttp.Response{Response: &http.Response{
								StatusCode: 200,
								Body:       ioutil.NopCloser(bytes.NewReader(response)),
							}}
						} else if apiExecuting == AssumeRole {
							out.RawResponse = &smithyhttp.Response{Response: &http.Response{
								StatusCode: 200,
								Body: ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf(assumeRoleRespMsg,
									time.Now().Add(15*time.Minute).Format("2006-01-02T15:04:05Z"))))),
							}}
						} else {
							out, metadata, err = next.HandleDeserialize(ctx, input)
						}
						return out, metadata, err
					}), middleware.After)
				})),
			}

			if c.client != nil {
				configSources = append(configSources, WithHTTPClient{c.client})
			}

			if len(c.configProfile) != 0 {
				configSources = append(configSources, WithSharedConfigProfile(c.configProfile))
			}

			config, err := LoadDefaultConfig(configSources...)
			if e, a := c.expectedError, err; !reflect.DeepEqual(e, a) {
				t.Fatalf("expected %v, but received %v", e, a)
			}

			if c.expectedError != nil {
				return
			}

			creds, err := config.Credentials.Retrieve(context.Background())
			if err != nil {
				t.Fatalf("expected no error, but received %v", err)
			}

			if e, a := c.expectedChain, credChain; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}

			if e, a := c.expectedAccessKey, creds.AccessKeyID; e != a {
				t.Errorf("expected %v, but received %v", e, a)
			}

			if e, a := c.expectedSecretKey, creds.SecretAccessKey; e != a {
				t.Errorf("expected %v, but received %v", e, a)
			}
		})
	}
}
