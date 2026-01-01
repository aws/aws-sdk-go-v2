package config

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	smithytime "github.com/aws/smithy-go/time"
)

// HTTPClient implementation that captures the User-Agent header and
// the features associated with it.
// why STS? because it's an already existing dependency to this package, but in reality
// any client that uses Credentials would work
type stsCaptureUserAgentFeatures struct {
	ua       string
	features []string
	// lost of credential providers use an intermediate STS client to fetch credentials.
	// capture which features were emitted by that client
	intermediateFeatures []string
}

func (u *stsCaptureUserAgentFeatures) Do(r *http.Request) (*http.Response, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	form := r.Form
	switch form.Get("Action") {
	// We use "DecodeAuthorization" as the sentinel value to say "this is the target operation, not an intermediate one"
	case "DecodeAuthorizationMessage":
		u.features = u.captureFeatures(r)
		return &http.Response{StatusCode: 503, Body: http.NoBody}, nil
	case "AssumeRoleWithWebIdentity":
		u.features = u.captureFeatures(r)
		return &http.Response{StatusCode: 200, Body: replaceDateInTemplate(assumeRoleWithWebIdentityResponse)}, nil
	}
	// if none of the above matched, assume his was an intermediate call and return the standard assume role response
	u.intermediateFeatures = u.captureFeatures(r)
	return &http.Response{StatusCode: 200, Body: replaceDateInTemplate(assumeRoleRespMsg)}, nil
}

func (u *stsCaptureUserAgentFeatures) captureFeatures(r *http.Request) []string {
	u.ua = r.Header.Get("User-Agent")
	parts := strings.Split(u.ua, " ")

	var features string
	for _, part := range parts {
		if after, hasPrefix := strings.CutPrefix(part, "m/"); hasPrefix {
			features = after
			break
		}
	}
	allFeatures := strings.Split(features, ",")
	result := make([]string, 0)
	// calling Split("x", ",") results in `["x"]`, so treat this as a special case
	if len(allFeatures) == 1 {
		result = append(result, allFeatures[0])
		return result
	}

	for _, f := range allFeatures {
		asRune := f[0]
		// only capture features related to credentials providers
		if (asRune >= 'e' && asRune <= 'z') || (asRune == '0') {
			result = append(result, string(asRune))
		}
	}

	return result
}

// bunch of our templates have "%s" as a placeholder for expiration. Using this to quickly create an io readcloser from those
func replaceDateInTemplate(template string) io.ReadCloser {
	t := fmt.Sprintf(template, smithytime.FormatDateTime(time.Now().Add(15*time.Minute)))
	return io.NopCloser(strings.NewReader(t))
}

type mockAssumeRole struct {
	TestInput func(*sts.AssumeRoleInput)
}

func (s mockAssumeRole) AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	if s.TestInput != nil {
		s.TestInput(params)
	}
	expiry := time.Now().Add(60 * time.Minute)

	return &sts.AssumeRoleOutput{
		AssumedRoleUser: &types.AssumedRoleUser{
			Arn: aws.String("arn:aws:sts::111111111111:assumed-role/assume-role-integration-test-role/Name"),
		},
		Credentials: &types.Credentials{
			AccessKeyId:     params.RoleArn,
			SecretAccessKey: aws.String("assumedSecretAccessKey"),
			SessionToken:    aws.String("assumedSessionToken"),
			Expiration:      &expiry,
		},
	}, nil
}

// intercepts calls to IMDS and handles them appropriately
// forwards any that don't match known paths to the inner HTTPClient
type imdsForwarder struct {
	innerClient HTTPClient
}

func (f imdsForwarder) Do(r *http.Request) (*http.Response, error) {
	if r.URL.Path == "/latest/api/token" {
		header := http.Header{}
		// bounce the TTL header
		const ttlHeader = "X-Aws-Ec2-Metadata-Token-Ttl-Seconds"
		header.Set(ttlHeader, r.Header.Get(ttlHeader))
		return &http.Response{StatusCode: 200, Header: header, Body: io.NopCloser(strings.NewReader("validToken"))}, nil
	}
	if r.URL.Path == "/latest/meta-data/iam/security-credentials/" {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("RoleName"))}, nil
	}
	if r.URL.Path == "/latest/meta-data/iam/security-credentials/RoleName" {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(ecsResponse))}, nil
	}
	return f.innerClient.Do(r)
}

type ssoForwarder struct {
	innerClient HTTPClient
}

func (f ssoForwarder) Do(r *http.Request) (*http.Response, error) {
	if strings.Contains(r.URL.Host, "sso") {
		body := fmt.Sprintf(
			getRoleCredentialsResponse,
			time.Now().
				Add(15*time.Minute).
				UnixNano()/int64(time.Millisecond))
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
	}
	return f.innerClient.Do(r)
}

func TestUserAgentCredentials(t *testing.T) {
	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	cases := map[string]struct {
		// given
		ExtraLoadFunctions []func(*LoadOptions) error
		Init               func(*testing.T)
		HTTPClientOverride func(existing HTTPClient) HTTPClient

		// assert
		Expect                []middleware.UserAgentFeature
		ExpectIntermediateSts []middleware.UserAgentFeature
	}{
		"hardcoded credentials": {
			Expect:             []middleware.UserAgentFeature{middleware.UserAgentFeatureCredentialsCode},
			ExtraLoadFunctions: []func(*LoadOptions) error{WithCredentialsProvider(credentials.NewStaticCredentialsProvider("key", "id", "session"))},
		},
		"environment credentials": {
			Init: func(t *testing.T) {
				t.Setenv("AWS_ACCESS_KEY_ID", "key")
				t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
			},
			Expect: []middleware.UserAgentFeature{middleware.UserAgentFeatureCredentialsEnvVars},
		},
		"web identity token from profile": {
			Init: func(t *testing.T) {
				t.Setenv("AWS_WEB_IDENTITY_TOKEN_FILE", "testdata/wit.txt")
				t.Setenv("AWS_REGION", "us-east-1")
			},
			ExtraLoadFunctions: []func(*LoadOptions) error{
				WithWebIdentityRoleCredentialOptions(func(options *stscreds.WebIdentityRoleOptions) {
					options.RoleARN = "test-arn"
				})},
			Expect: []middleware.UserAgentFeature{middleware.UserAgentFeatureCredentialsEnvVarsStsWebIDToken},
		},
		"sts assume role": {
			ExtraLoadFunctions: []func(*LoadOptions) error{
				WithCredentialsProvider(stscreds.NewAssumeRoleProvider(mockAssumeRole{}, "role-arn", func(o *stscreds.AssumeRoleOptions) {}))},
			Expect: []middleware.UserAgentFeature{middleware.UserAgentFeatureCredentialsStsAssumeRole},
		},
		"credentials profile": {
			ExtraLoadFunctions: []func(*LoadOptions) error{
				WithSharedCredentialsFiles([]string{"testdata/config_ua_credential_provider"}),
			},
			Expect: []middleware.UserAgentFeature{middleware.UserAgentFeatureCredentialsProfile},
		},
		"assume role with source profile": {
			ExtraLoadFunctions: []func(*LoadOptions) error{
				WithSharedConfigFiles([]string{"testdata/config_ua_credential_provider"}),
				WithSharedConfigProfile("A"),
			},
			Expect: []middleware.UserAgentFeature{
				middleware.UserAgentFeatureCredentialsProfile,
				middleware.UserAgentFeatureCredentialsProfileSourceProfile,
				middleware.UserAgentFeatureCredentialsStsAssumeRole,
			},
			ExpectIntermediateSts: []middleware.UserAgentFeature{
				middleware.UserAgentFeatureCredentialsProfile,
				middleware.UserAgentFeatureCredentialsProfileSourceProfile,
			},
		},
		"named credentials profile http": {
			ExtraLoadFunctions: []func(*LoadOptions) error{
				WithSharedConfigFiles([]string{"testdata/config_ua_credential_provider"}),
				WithSharedConfigProfile("ecscontainer"),
			},
			Init: func(t *testing.T) {
				ecsMetadataServer := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte(ecsResponse))
					}))
				t.Setenv("AWS_CONTAINER_CREDENTIALS_FULL_URI", ecsMetadataServer.URL)
				t.Cleanup(ecsMetadataServer.Close)
			},
			Expect: []middleware.UserAgentFeature{
				middleware.UserAgentFeatureCredentialsStsAssumeRole,
				middleware.UserAgentFeatureCredentialsHTTP,
				middleware.UserAgentFeatureCredentialsProfileNamedProvider,
			},
			ExpectIntermediateSts: []middleware.UserAgentFeature{
				middleware.UserAgentFeatureCredentialsHTTP,
				middleware.UserAgentFeatureCredentialsProfileNamedProvider,
			},
		},
		"sso credentials": {
			Init: func(t *testing.T) {
				// filename is the sha1 of the sso_sesion value
				ssoSetup(t, "d033e22ae348aeb5660fc2140aec35850c4da997.json")
			},
			HTTPClientOverride: func(existing HTTPClient) HTTPClient {
				return ssoForwarder{existing}
			},
			Expect: []middleware.UserAgentFeature{
				middleware.UserAgentFeatureCredentialsProfileSso,
				middleware.UserAgentFeatureCredentialsSso,
			},
			ExtraLoadFunctions: []func(*LoadOptions) error{
				WithSharedConfigFiles([]string{"testdata/config_ua_credential_provider"}),
				WithSharedConfigProfile("sso-token"),
			},
		},
		"sso credentials legacy": {
			Init: func(t *testing.T) {
				// filename is sha1 of the host name
				ssoSetup(t, "239e5a291d940f9e746e1c722f1b936d0dc3814d.json")
			},
			HTTPClientOverride: func(existing HTTPClient) HTTPClient {
				return ssoForwarder{existing}
			},
			Expect: []middleware.UserAgentFeature{
				middleware.UserAgentFeatureCredentialsProfileSsoLegacy,
				middleware.UserAgentFeatureCredentialsSsoLegacy,
			},
			ExtraLoadFunctions: []func(*LoadOptions) error{
				WithSharedConfigFiles([]string{"testdata/config_ua_credential_provider"}),
				WithSharedConfigProfile("sso_creds"),
			},
		},
		"profile process": {
			Expect: []middleware.UserAgentFeature{
				middleware.UserAgentFeatureCredentialsProfileProcess,
				middleware.UserAgentFeatureCredentialsProcess,
			},
			ExtraLoadFunctions: []func(*LoadOptions) error{
				func(options *LoadOptions) error {
					if runtime.GOOS == "windows" {
						options.SharedConfigFiles = []string{"testdata/config_ua_credential_provider_windows"}
					} else {
						options.SharedConfigFiles = []string{"testdata/config_ua_credential_provider"}
					}
					return nil
				},
				WithSharedConfigProfile("process"),
			},
		},
		"imds": {
			Init: func(t *testing.T) {
				// Actual value doesn't matter, as long as the env variable is set
				t.Setenv("AWS_EC2_METADATA_SERVICE_ENDPOINT", "http://non-routable-just-setting-to-signal-ec2.invalid")
			},
			HTTPClientOverride: func(existing HTTPClient) HTTPClient {
				return imdsForwarder{existing}
			},
			Expect: []middleware.UserAgentFeature{
				middleware.UserAgentFeatureCredentialsIMDS,
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var ua stsCaptureUserAgentFeatures
			optFns := []func(*LoadOptions) error{
				WithHTTPClient(&ua),
				// If any credentials or config files are set up locally, ignore them
				// to prevent accidentally picking them up and messing up the tests
				WithSharedConfigFiles([]string{""}),
				WithSharedCredentialsFiles([]string{""}),
				WithRetryer(func() aws.Retryer {
					return aws.NopRetryer{}
				}),
				// some endpoint rules seem to require a region set up
				WithRegion("us-east-1"),
			}
			if c.ExtraLoadFunctions != nil {
				optFns = append(optFns, c.ExtraLoadFunctions...)
			}
			if c.Init != nil {
				c.Init(t)
			}
			if c.HTTPClientOverride != nil {
				optFns = append(optFns, WithHTTPClient(c.HTTPClientOverride(&ua)))
			}
			cfg, err := LoadDefaultConfig(context.TODO(), optFns...)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			client := sts.NewFromConfig(cfg)
			// doesn't matter, we're just making a call to record user agent
			in := "test"
			client.DecodeAuthorizationMessage(context.TODO(), &sts.DecodeAuthorizationMessageInput{EncodedMessage: &in})
			compareFeatures(t, c.Expect, ua.features)
			if c.ExpectIntermediateSts != nil {
				compareFeatures(t, c.ExpectIntermediateSts, ua.intermediateFeatures)
			}
		})
	}
}

func compareFeatures(t *testing.T, expected []middleware.UserAgentFeature, actual []string) {
	expectedStr := make([]string, 0, len(expected))
	for _, e := range expected {
		expectedStr = append(expectedStr, string(e))
	}

	slices.Sort(expectedStr)
	slices.Sort(actual)
	ok := slices.Equal(actual, expectedStr)
	if !ok {
		t.Errorf("expect %v, got %v", expectedStr, actual)
	}
}

func ssoSetup(t *testing.T, filename string) {
	dir := t.TempDir()
	cacheDir := filepath.Join(dir, ".aws", "sso", "cache")
	err := os.MkdirAll(cacheDir, 0750)
	if err != nil {
		t.Fatal(err)
	}

	tokenFile, err := os.Create(filepath.Join(cacheDir, filename))
	if err != nil {
		t.Fatal(err)
	}

	defer tokenFile.Close()

	tokenContent := fmt.Sprintf(ssoTokenCacheFile, smithytime.FormatDateTime(time.Now().Add(15*time.Minute)))
	_, err = tokenFile.WriteString(tokenContent)
	if err != nil {
		t.Fatal(err)
	}

	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", dir)
	} else {
		t.Setenv("HOME", dir)
	}
}
