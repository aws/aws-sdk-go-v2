package auth

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

type dbTokenTestCase struct {
	endpoint            string
	region              string
	expires             time.Duration
	credsExpireIn       time.Duration
	expectedHost        string
	expectedQueryParams []string
	expectedError       string
}

type tokenGenFunc func(ctx context.Context, endpoint, region string, creds aws.CredentialsProvider, optFns ...func(options *TokenOptions)) (string, error)

func TestGenerateDbConnectAuthToken(t *testing.T) {
	cases := map[string]dbTokenTestCase{
		"no region": {
			endpoint:      "https://oo0bar1baz2quux3quuux4.dsql.us-east-1.on.aws",
			expectedError: "no region",
		},
		"no endpoint": {
			region:        "us-west-2",
			expectedError: "endpoint is required",
		},
		"endpoint with scheme": {
			endpoint:            "https://oo0bar1baz2quux3quuux4.dsql.us-east-1.on.aws",
			region:              "us-east-1",
			expectedHost:        "oo0bar1baz2quux3quuux4.dsql.us-east-1.on.aws",
			expectedQueryParams: []string{"Action=DbConnect"},
		},
		"endpoint without scheme": {
			endpoint:            "oo0bar1baz2quux3quuux4.dsql.us-east-1.on.aws",
			region:              "us-east-1",
			expectedHost:        "oo0bar1baz2quux3quuux4.dsql.us-east-1.on.aws",
			expectedQueryParams: []string{"Action=DbConnect"},
		},
		"endpoint with region and expires": {
			endpoint:     "peccy.dsql.us-east-1.on.aws",
			region:       "us-east-1",
			expires:      time.Second * 450,
			expectedHost: "peccy.dsql.us-east-1.on.aws",
			expectedQueryParams: []string{
				"Action=DbConnect",
				"X-Amz-Algorithm=AWS4-HMAC-SHA256",
				"X-Amz-Credential=akid/20240827/us-east-1/dsql/aws4_request",
				"X-Amz-Date=20240827T000000Z",
				"X-Amz-Expires=450"},
		},
		"pick credential expires when less than expires": {
			endpoint:      "peccy.dsql.us-east-1.on.aws",
			region:        "us-east-1",
			credsExpireIn: time.Second * 100,
			expires:       time.Second * 450,
			expectedHost:  "peccy.dsql.us-east-1.on.aws",
			expectedQueryParams: []string{
				"Action=DbConnect",
				"X-Amz-Algorithm=AWS4-HMAC-SHA256",
				"X-Amz-Credential=akid/20240827/us-east-1/dsql/aws4_request",
				"X-Amz-Date=20240827T000000Z",
				"X-Amz-Expires=100"},
		},
	}

	for _, c := range cases {
		creds := &staticCredentials{AccessKey: "akid", SecretKey: "secret", expiresIn: c.credsExpireIn}
		defer withTempGlobalTime(time.Date(2024, time.August, 27, 0, 0, 0, 0, time.UTC))()
		optFns := func(options *TokenOptions) {}
		if c.expires != 0 {
			optFns = func(options *TokenOptions) {
				options.ExpiresIn = c.expires
			}
		}
		verifyTestCase(GenerateDbConnectAuthToken, c, creds, optFns, t)

		// Update the test case to use Admin variant
		updated := []string{}
		for _, part := range c.expectedQueryParams {
			if part == "Action=DbConnect" {
				part = "Action=DbConnectAdmin"
			}
			updated = append(updated, part)
		}
		c.expectedQueryParams = updated

		verifyTestCase(GenerateDBConnectAdminAuthToken, c, creds, optFns, t)
	}
}

func verifyTestCase(f tokenGenFunc, c dbTokenTestCase, creds aws.CredentialsProvider, optFns func(options *TokenOptions), t *testing.T) {
	token, err := f(context.Background(), c.endpoint, c.region, creds, optFns)
	isErrorExpected := len(c.expectedError) > 0
	if err != nil && !isErrorExpected {
		t.Fatalf("expect no err, got: %v", err)
	} else if err == nil && isErrorExpected {
		t.Fatalf("Expected error %v got none", c.expectedError)
	}
	// adding a scheme so we can parse it back as a URL. This is because comparing
	// just direct string comparison was failing since "Action=DbConnect" is a substring or
	// "Action=DBConnectAdmin"
	parsed, err := url.Parse("http://" + token)
	if err != nil {
		t.Fatalf("Couldn't parse the token %v to URL after adding a scheme, got: %v", token, err)
	}
	if parsed.Host != c.expectedHost {
		t.Errorf("expect host %v, got %v", c.expectedHost, parsed.Host)
	}

	q := parsed.Query()
	queryValuePair := map[string]any{}
	for k, v := range q {
		pair := k + "=" + v[0]
		queryValuePair[pair] = struct{}{}
	}

	for _, part := range c.expectedQueryParams {
		if _, ok := queryValuePair[part]; !ok {
			t.Errorf("expect part %s to be present at token %s", part, token)
		}
	}
	if token != "" && c.expires == 0 {
		if !strings.Contains(token, "X-Amz-Expires=900") {
			t.Errorf("expect token to contain default X-Amz-Expires value of 900, got %v", token)
		}
	}
}

type staticCredentials struct {
	AccessKey, SecretKey, Session string
	expiresIn                     time.Duration
}

func (s *staticCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	c := aws.Credentials{
		AccessKeyID:     s.AccessKey,
		SecretAccessKey: s.SecretKey,
		SessionToken:    s.Session,
	}
	if s.expiresIn != 0 {
		c.CanExpire = true
		c.Expires = sdk.NowTime().Add(s.expiresIn)
	}
	return c, nil
}

func withTempGlobalTime(t time.Time) func() {
	sdk.NowTime = func() time.Time { return t }
	return func() { sdk.NowTime = time.Now }
}
