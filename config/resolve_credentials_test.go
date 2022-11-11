package config

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithytime "github.com/aws/smithy-go/time"
)

func swapECSContainerURI(path string) func() {
	o := ecsContainerEndpoint
	ecsContainerEndpoint = path
	return func() {
		ecsContainerEndpoint = o
	}
}

func setupCredentialsEndpoints(t *testing.T) (aws.EndpointResolverWithOptions, func()) {
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
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(500)
				return
			}

			form := r.Form

			switch form.Get("Action") {
			case "AssumeRole":
				w.Write([]byte(fmt.Sprintf(
					assumeRoleRespMsg,
					smithytime.FormatDateTime(time.Now().
						Add(15*time.Minute)))))
				return
			case "AssumeRoleWithWebIdentity":
				w.Write([]byte(fmt.Sprintf(assumeRoleWithWebIdentityResponse,
					smithytime.FormatDateTime(time.Now().
						Add(15*time.Minute)))))
				return
			default:
				w.WriteHeader(404)
				return
			}
		}))

	ssoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(
			getRoleCredentialsResponse,
			time.Now().
				Add(15*time.Minute).
				UnixNano()/int64(time.Millisecond))))
	}))

	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			switch service {
			case sts.ServiceID:
				return aws.Endpoint{
					URL: stsServer.URL,
				}, nil
			case sso.ServiceID:
				return aws.Endpoint{
					URL: ssoServer.URL,
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
		ssoServer.Close()
		stsServer.Close()
	}
}

func ssoTestSetup() (fn func(), err error) {
	dir, err := ioutil.TempDir(os.TempDir(), "sso-test")
	if err != nil {
		return nil, err
	}

	cleanupTestDir := func() {
		os.RemoveAll(dir)
	}
	defer func() {
		if err != nil {
			cleanupTestDir()
		}
	}()

	cacheDir := filepath.Join(dir, ".aws", "sso", "cache")
	err = os.MkdirAll(cacheDir, 0750)
	if err != nil {
		return nil, err
	}

	tokenFile, err := os.Create(filepath.Join(cacheDir, "eb5e43e71ce87dd92ec58903d76debd8ee42aefd.json"))
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := tokenFile.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("close error: %v, original error: %w", closeErr, err)
		}
	}()

	_, err = tokenFile.WriteString(fmt.Sprintf(ssoTokenCacheFile, time.Now().
		Add(15*time.Minute).
		Format(time.RFC3339)))
	if err != nil {
		return nil, err
	}

	if runtime.GOOS == "windows" {
		os.Setenv("USERPROFILE", dir)
	} else {
		os.Setenv("HOME", dir)
	}

	return cleanupTestDir, nil
}

func TestSharedConfigCredentialSource(t *testing.T) {
	var configFileForWindows = filepath.Join("testdata", "config_source_shared_for_windows")
	var configFile = filepath.Join("testdata", "config_source_shared")

	var credFileForWindows = filepath.Join("testdata", "credentials_source_shared_for_windows")
	var credFile = filepath.Join("testdata", "credentials_source_shared")

	cases := map[string]struct {
		name                 string
		envProfile           string
		configProfile        string
		expectedError        string
		expectedAccessKey    string
		expectedSecretKey    string
		expectedSessionToken string
		expectedChain        []string
		init                 func() (func(), error)
		dependentOnOS        bool
	}{
		"credential source and source profile": {
			envProfile:    "invalid_source_and_credential_source",
			expectedError: "only one credential type may be specified per profile",
			init: func() (func(), error) {
				os.Setenv("AWS_ACCESS_KEY", "access_key")
				os.Setenv("AWS_SECRET_KEY", "secret_key")
				return func() {}, nil
			},
		},
		"env var credential source": {
			configProfile:        "env_var_credential_source",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_env",
			},
			init: func() (func(), error) {
				os.Setenv("AWS_ACCESS_KEY", "access_key")
				os.Setenv("AWS_SECRET_KEY", "secret_key")
				return func() {}, nil
			},
		},
		"ec2metadata credential source": {
			envProfile: "ec2metadata",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ec2",
			},
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
		},
		"ecs container credential source": {
			envProfile:           "ecscontainer",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ecs",
			},
			init: func() (func(), error) {
				os.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "/ECS")
				return func() {}, nil
			},
		},
		"chained assume role with env creds": {
			envProfile:           "chained_assume_role",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
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
			envProfile:           "cred_proc_arn_set",
			dependentOnOS:        true,
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_proc_role_arn",
			},
		},
		"chained assume role with credential process": {
			envProfile:           "chained_cred_proc",
			dependentOnOS:        true,
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_proc_source_prof",
			},
		},
		"credential source overrides config source": {
			envProfile:           "credentials_overide",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ec2",
			},
			init: func() (func(), error) {
				os.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "/ECS")
				return func() {}, nil
			},
		},
		"only credential source": {
			envProfile:           "only_credentials_source",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_role_w_creds_role_arn_ecs",
			},
			init: func() (func(), error) {
				os.Setenv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "/ECS")
				return func() {}, nil
			},
		},
		"sso credentials": {
			envProfile:           "sso_creds",
			expectedAccessKey:    "SSO_AKID",
			expectedSecretKey:    "SSO_SECRET_KEY",
			expectedSessionToken: "SSO_SESSION_TOKEN",
			init: func() (func(), error) {
				return ssoTestSetup()
			},
		},
		"chained assume role with sso credentials": {
			envProfile:           "source_sso_creds",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"source_sso_creds_arn",
			},
			init: func() (func(), error) {
				return ssoTestSetup()
			},
		},
		"chained assume role with sso and static credentials": {
			envProfile:           "assume_sso_and_static",
			expectedAccessKey:    "AKID",
			expectedSecretKey:    "SECRET",
			expectedSessionToken: "SESSION_TOKEN",
			expectedChain: []string{
				"assume_sso_and_static_arn",
			},
		},
		"invalid sso configuration": {
			envProfile:    "sso_invalid",
			expectedError: "profile \"sso_invalid\" is configured to use SSO but is missing required configuration: sso_region, sso_start_url",
		},
		"environment credentials with invalid sso": {
			envProfile:        "sso_invalid",
			expectedAccessKey: "access_key",
			expectedSecretKey: "secret_key",
			init: func() (func(), error) {
				os.Setenv("AWS_ACCESS_KEY", "access_key")
				os.Setenv("AWS_SECRET_KEY", "secret_key")
				return func() {}, nil
			},
		},
		"sso mixed with credential process provider": {
			envProfile:           "sso_mixed_credproc",
			expectedAccessKey:    "SSO_AKID",
			expectedSecretKey:    "SSO_SECRET_KEY",
			expectedSessionToken: "SSO_SESSION_TOKEN",
			init: func() (func(), error) {
				return ssoTestSetup()
			},
		},
		"sso mixed with web identity token provider": {
			envProfile:           "sso_mixed_webident",
			expectedAccessKey:    "WEB_IDENTITY_AKID",
			expectedSecretKey:    "WEB_IDENTITY_SECRET",
			expectedSessionToken: "WEB_IDENTITY_SESSION_TOKEN",
		},
		"SSO Session missing region": {
			envProfile:    "sso-session-missing-region",
			expectedError: "profile \"sso-session-missing-region\" is configured to use SSO but is missing required configuration: sso_region",
		},
		"SSO Session mismatched region": {
			envProfile:    "sso-session-mismatched-region",
			expectedError: "sso_region in profile \"sso-session-mismatched-region\" must match sso_region in sso-session",
		},
		"web identity": {
			envProfile:           "webident",
			expectedAccessKey:    "WEB_IDENTITY_AKID",
			expectedSecretKey:    "WEB_IDENTITY_SECRET",
			expectedSessionToken: "WEB_IDENTITY_SESSION_TOKEN",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			restoreEnv := awstesting.StashEnv()
			defer awstesting.PopEnv(restoreEnv)

			if c.dependentOnOS && runtime.GOOS == "windows" {
				os.Setenv("AWS_CONFIG_FILE", configFileForWindows)
				os.Setenv("AWS_SHARED_CREDENTIALS_FILE", credFileForWindows)
			} else {
				os.Setenv("AWS_CONFIG_FILE", configFile)
				os.Setenv("AWS_SHARED_CREDENTIALS_FILE", credFile)
			}

			os.Setenv("AWS_REGION", "us-east-1")
			if len(c.envProfile) != 0 {
				os.Setenv("AWS_PROFILE", c.envProfile)
			}

			endpointResolver, cleanupFn := setupCredentialsEndpoints(t)
			defer cleanupFn()

			var cleanup func()
			if c.init != nil {
				var err error
				cleanup, err = c.init()
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
				defer cleanup()
			}

			var credChain []string

			loadOptions := []func(*LoadOptions) error{
				WithEndpointResolverWithOptions(endpointResolver),
				WithAPIOptions([]func(*middleware.Stack) error{
					func(stack *middleware.Stack) error {
						return stack.Initialize.Add(middleware.InitializeMiddlewareFunc("GetRoleArns",
							func(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
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
				t.Errorf("expect %v, but received %v", e, a)
			}

			if e, a := c.expectedSessionToken, creds.SessionToken; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestResolveCredentialsCacheOptions(t *testing.T) {
	var cfg aws.Config
	var optionsFnCalled bool

	err := resolveCredentials(context.Background(), &cfg, configs{LoadOptions{
		CredentialsCacheOptions: func(o *aws.CredentialsCacheOptions) {
			optionsFnCalled = true
			o.ExpiryWindow = time.Minute * 5
		},
	}})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if !optionsFnCalled {
		t.Errorf("expect options to be called")
	}
}

func TestResolveCredentialsIMDSClient(t *testing.T) {
	expectEnabled := func(t *testing.T, err error) {
		if err == nil {
			t.Fatalf("expect error got none")
		}
		if e, a := "expected HTTP client error", err.Error(); !strings.Contains(a, e) {
			t.Fatalf("expected %v error in %v", e, a)
		}
	}

	expectDisabled := func(t *testing.T, err error) {
		var oe *smithy.OperationError
		if !errors.As(err, &oe) {
			t.Fatalf("unexpected error: %v", err)
		} else {
			e := errors.Unwrap(oe)
			if e == nil {
				t.Fatalf("unexpected empty operation error: %v", oe)
			} else {
				if !strings.HasPrefix(e.Error(), "access disabled to EC2 IMDS") {
					t.Fatalf("unexpected operation error: %v", oe)
				}
			}
		}
	}

	testcases := map[string]struct {
		enabledState  imds.ClientEnableState
		envvar        string
		expectedState imds.ClientEnableState
		expectedError func(*testing.T, error)
	}{
		"default no options": {
			expectedState: imds.ClientDefaultEnableState,
			expectedError: expectEnabled,
		},

		"state enabled": {
			enabledState:  imds.ClientEnabled,
			expectedState: imds.ClientEnabled,
			expectedError: expectEnabled,
		},
		"state disabled": {
			enabledState:  imds.ClientDisabled,
			expectedState: imds.ClientDisabled,
			expectedError: expectDisabled,
		},

		"env var DISABLED true": {
			envvar:        "true",
			expectedState: imds.ClientDisabled,
			expectedError: expectDisabled,
		},
		"env var DISABLED false": {
			envvar:        "false",
			expectedState: imds.ClientEnabled,
			expectedError: expectEnabled,
		},

		"option state enabled overrides env var DISABLED true": {
			enabledState:  imds.ClientEnabled,
			envvar:        "true",
			expectedState: imds.ClientEnabled,
			expectedError: expectEnabled,
		},
		"option state disabled overrides env var DISABLED false": {
			enabledState:  imds.ClientDisabled,
			envvar:        "false",
			expectedState: imds.ClientDisabled,
			expectedError: expectDisabled,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			restoreEnv := awstesting.StashEnv()
			defer awstesting.PopEnv(restoreEnv)

			var httpClient HTTPClient
			if tc.expectedState == imds.ClientDisabled {
				httpClient = stubErrorClient{err: fmt.Errorf("expect HTTP client not to be called")}
			} else {
				httpClient = stubErrorClient{err: fmt.Errorf("expected HTTP client error")}
			}

			opts := []func(*LoadOptions) error{
				WithRetryer(func() aws.Retryer { return aws.NopRetryer{} }),
				WithHTTPClient(httpClient),
			}

			if tc.enabledState != imds.ClientDefaultEnableState {
				opts = append(opts,
					WithEC2IMDSClientEnableState(tc.enabledState),
				)
			}

			if tc.envvar != "" {
				os.Setenv("AWS_EC2_METADATA_DISABLED", tc.envvar)
			}

			c, err := LoadDefaultConfig(context.TODO(), opts...)
			if err != nil {
				t.Fatalf("could not load config: %s", err)
			}

			creds := c.Credentials

			_, err = creds.Retrieve(context.TODO())
			tc.expectedError(t, err)
		})
	}
}

type stubErrorClient struct {
	err error
}

func (c stubErrorClient) Do(*http.Request) (*http.Response, error) { return nil, c.err }
