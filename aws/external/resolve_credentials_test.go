package external

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/awserr"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting"
	"github.com/jviney/aws-sdk-go-v2/service/sts"
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
			if r.URL.Path == "/meta-data/iam/security-credentials/RoleName" {
				w.Write([]byte(ec2MetadataResponse))
			} else if r.URL.Path == "/meta-data/iam/security-credentials/" {
				w.Write([]byte("RoleName"))
			} else {
				w.Write([]byte(""))
			}
		}))

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
			case "ec2metadata":
				return aws.Endpoint{
					URL: ec2MetadataServer.URL,
				}, nil
			case "sts":
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
	}{
		"credential source and source profile": {
			envProfile:    "invalid_source_and_credential_source",
			expectedError: awserr.New(ErrCodeSharedConfig, "only source profile or credential source can be specified, not both", nil),
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

			configSources := []Config{
				WithHandlersFunc(func(handlers aws.Handlers) aws.Handlers {
					handlers.Sign.PushBack(func(r *aws.Request) {
						if r.Config.Credentials == aws.AnonymousCredentials {
							return
						}
						params := r.Params.(*sts.AssumeRoleInput)
						credChain = append(credChain, *params.RoleArn)
					})
					return handlers
				}),
				WithEndpointResolverFunc(func(resolver aws.EndpointResolver) aws.EndpointResolver {
					return endpointResolver
				}),
			}

			if len(c.configProfile) != 0 {
				configSources = append(configSources, WithSharedConfigProfile(c.configProfile))
			}

			config, err := LoadDefaultAWSConfig(configSources...)
			if e, a := c.expectedError, err; !reflect.DeepEqual(e, a) {
				t.Fatalf("expected %v, but received %v", e, a)
			}

			// TODO: Should add a Logger Resolver?
			//config.Logger = t

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

const ecsResponse = `{
	  "Code": "Success",
	  "Type": "AWS-HMAC",
	  "AccessKeyId" : "ecs-access-key",
	  "SecretAccessKey" : "ecs-secret-key",
	  "Token" : "token",
	  "Expiration" : "2100-01-01T00:00:00Z",
	  "LastUpdated" : "2009-11-23T0:00:00Z"
	}`

const ec2MetadataResponse = `{
	  "Code": "Success",
	  "Type": "AWS-HMAC",
	  "AccessKeyId" : "ec2-access-key",
	  "SecretAccessKey" : "ec2-secret-key",
	  "Token" : "token",
	  "Expiration" : "2100-01-01T00:00:00Z",
	  "LastUpdated" : "2009-11-23T0:00:00Z"
	}`

const assumeRoleRespMsg = `
<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <AssumeRoleResult>
    <AssumedRoleUser>
      <Arn>arn:aws:sts::account_id:assumed-role/role/session_name</Arn>
      <AssumedRoleId>AKID:session_name</AssumedRoleId>
    </AssumedRoleUser>
    <Credentials>
      <AccessKeyId>AKID</AccessKeyId>
      <SecretAccessKey>SECRET</SecretAccessKey>
      <SessionToken>SESSION_TOKEN</SessionToken>
      <Expiration>%s</Expiration>
    </Credentials>
  </AssumeRoleResult>
  <ResponseMetadata>
    <RequestId>request-id</RequestId>
  </ResponseMetadata>
</AssumeRoleResponse>
`

func initConfigTestEnv() (oldEnv []string) {
	oldEnv = awstesting.StashEnv()
	os.Setenv("AWS_CONFIG_FILE", "file_not_exists")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "file_not_exists")

	return oldEnv
}

func TestAssumeRole(t *testing.T) {
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_w_creds")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(
			assumeRoleRespMsg,
			time.Now().Add(15*time.Minute).Format("2006-01-02T15:04:05Z"))))
	}))
	defer server.Close()

	config, err := LoadDefaultAWSConfig(
		WithEndpointResolverFunc(func(resolver aws.EndpointResolver) aws.EndpointResolver {
			return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: server.URL}, nil
			})
		}),
	)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	creds, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SESSION_TOKEN", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "AssumeRoleProvider", creds.Source; !strings.Contains(a, e) {
		t.Errorf("expect %v, to be in %v", e, a)
	}
}

func TestAssumeRole_WithMFA(t *testing.T) {
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_w_creds")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e, a := r.FormValue("SerialNumber"), "0123456789"; e != a {
			t.Errorf("expect %v, got %v", e, a)
		}
		if e, a := r.FormValue("TokenCode"), "tokencode"; e != a {
			t.Errorf("expect %v, got %v", e, a)
		}
		if e, a := "900", r.FormValue("DurationSeconds"); e != a {
			t.Errorf("expect %v, got %v", e, a)
		}

		w.Write([]byte(fmt.Sprintf(
			assumeRoleRespMsg,
			time.Now().Add(15*time.Minute).Format("2006-01-02T15:04:05Z"))))
	}))
	defer server.Close()

	customProviderCalled := false
	config, err := LoadDefaultAWSConfig(
		WithRegion("us-east-1"),
		WithSharedConfigProfile("assume_role_w_mfa"),
		WithMFATokenFunc(func() (string, error) {
			customProviderCalled = true

			return "tokencode", nil
		}),
		WithEndpointResolverFunc(func(resolver aws.EndpointResolver) aws.EndpointResolver {
			return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: server.URL}, nil
			})
		}),
	)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	creds, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if !customProviderCalled {
		t.Errorf("expect true")
	}

	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SESSION_TOKEN", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "AssumeRoleProvider", creds.Source; !strings.Contains(a, e) {
		t.Errorf("expect %v, to be in %v", e, a)
	}
}

func TestAssumeRole_WithMFA_NoTokenProvider(t *testing.T) {
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_w_creds")

	_, err := LoadDefaultAWSConfig(WithSharedConfigProfile("assume_role_w_mfa"))
	if e, a := (AssumeRoleTokenProviderNotSetError{}), err; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestAssumeRole_InvalidSourceProfile(t *testing.T) {
	// Backwards compatibility with Shared config disabled
	// assume role should not be built into the config.
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_invalid_source_profile")

	_, err := LoadDefaultAWSConfig()
	if err == nil {
		t.Fatalf("expect error, got none")
	}

	expectMsg := "failed to load assume role"
	if e, a := expectMsg, err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v, to be in %v", e, a)
	}
}

func TestAssumeRole_ExtendedDuration(t *testing.T) {
	restoreEnv := initConfigTestEnv()
	defer awstesting.PopEnv(restoreEnv)

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", testConfigFilename)
	os.Setenv("AWS_PROFILE", "assume_role_w_creds_ext_dur")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e, a := "1800", r.FormValue("DurationSeconds"); e != a {
			t.Errorf("expect %v, got %v", e, a)
		}

		w.Write([]byte(fmt.Sprintf(
			assumeRoleRespMsg,
			time.Now().Add(15*time.Minute).Format("2006-01-02T15:04:05Z"))))
	}))
	defer server.Close()

	config, err := LoadDefaultAWSConfig(
		WithEndpointResolverFunc(func(resolver aws.EndpointResolver) aws.EndpointResolver {
			return aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: server.URL}, nil
			})
		}),
	)
	// TODO: Set Assume Role Duration
	//	AssumeRoleDuration: 30 * time.Minute,
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	creds, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SESSION_TOKEN", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "AssumeRoleProvider", creds.Source; !strings.Contains(a, e) {
		t.Errorf("expect %v, to be in %v", e, a)
	}
}
