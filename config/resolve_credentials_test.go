package config

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/awslabs/smithy-go/middleware"
)

func swapECSContainerURI(path string) func() {
	o := ecsContainerEndpoint
	ecsContainerEndpoint = path
	return func() {
		ecsContainerEndpoint = o
	}
}

func setupCredentialsEndpoints(t *testing.T) (aws.EndpointResolver, func()) {
	ecsMetadataServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ECS" {
				w.Write([]byte(ecsResponse))
			} else {
				w.Write([]byte(""))
			}
		}))
	resetECSEndpoint := swapECSContainerURI(ecsMetadataServer.URL)

	ec2MetadataServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/latest/meta-data/iam/security-credentials/RoleName" {
				w.Write([]byte(ec2MetadataResponse))
			} else if r.URL.Path == "/latest/meta-data/iam/security-credentials/" {
				w.Write([]byte("RoleName"))
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

	os.Setenv("AWS_EC2_METADATA_SERVICE_ENDPOINT", ec2MetadataServer.URL)

	stsServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf(
				assumeRoleRespMsg,
				time.Now().
					Add(15*time.Minute).
					Format("2006-01-02T15:04:05Z"))))
		}))

	resolver := aws.EndpointResolverFunc(
		func(service, region string) (aws.Endpoint, error) {
			switch service {
			case sts.ServiceID:
				return aws.Endpoint{
					URL: stsServer.URL,
				}, nil
			default:
				return aws.Endpoint{},
					fmt.Errorf("unknown service endpoint, %s", service)
			}
		})

	return resolver, func() {
		resetECSEndpoint()
		ecsMetadataServer.Close()
		ec2MetadataServer.Close()
		stsServer.Close()
	}
}

func TestSharedConfigCredentialSource(t *testing.T) {
	var configFileForWindows = filepath.Join("testdata", "credential_source_config_for_windows")
	var configFile = filepath.Join("testdata", "credential_source_config")

	cases := map[string]struct {
		name              string
		envProfile        string
		configProfile     string
		expectedError     string
		expectedAccessKey string
		expectedSecretKey string
		expectedChain     []string
		init              func()
		dependentOnOS     bool
	}{
		"credential source and source profile": {
			envProfile:    "invalid_source_and_credential_source",
			expectedError: "nly source profile or credential source can be specified",
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

			endpointResolver, cleanupFn := setupCredentialsEndpoints(t)
			defer cleanupFn()

			if c.init != nil {
				c.init()
			}

			var credChain []string

			loadOptions := []func(*LoadOptions) error{
				WithEndpointResolver(endpointResolver),
				WithAPIOptions([]func(*middleware.Stack) error{
					func(stack *middleware.Stack) error {
						return stack.Initialize.Add(middleware.InitializeMiddlewareFunc("GetRoleArns", func(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
						) (
							out middleware.InitializeOutput, metadata middleware.Metadata, err error,
						) {
							switch v := in.Parameters.(type) {
							case *sts.AssumeRoleInput:
								credChain = append(credChain, *v.RoleArn)
							}

							return next.HandleInitialize(ctx, in)
						}), middleware.After)
					},
				}),
			}

			if len(c.configProfile) != 0 {
				loadOptions = append(loadOptions, WithSharedConfigProfile(c.configProfile))
			}

			config, err := LoadDefaultConfig(context.Background(), loadOptions...)
			if err != nil {
				if len(c.expectedError) > 0 {
					if e, a := c.expectedError, err.Error(); !strings.Contains(a, e) {
						t.Fatalf("expect %v, but got %v", e, a)
					}
					return
				}
				t.Fatalf("expect no error, got %v", err)
			} else if len(c.expectedError) > 0 {
				t.Fatalf("expect error, got none")
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
