//go:build go1.16
// +build go1.16

package ssocreds

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	smithybearer "github.com/aws/smithy-go/auth/bearer"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSSOTokenProvider(t *testing.T) {
	restoreTime := sdk.TestingUseReferenceTime(time.Date(2021, 12, 21, 12, 21, 1, 0, time.UTC))
	defer restoreTime()

	tempDir, err := ioutil.TempDir(os.TempDir(), "aws-sdk-go-v2-"+t.Name())
	if err != nil {
		t.Fatalf("failed to create temporary test directory, %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("failed to cleanup temporary test directory, %v", err)
		}
	}()

	cases := map[string]struct {
		setup         func() error
		postRetrieve  func() error
		client        CreateTokenAPIClient
		cacheFilePath string
		optFns        []func(*SSOTokenProviderOptions)

		expectToken smithybearer.Token
		expectErr   string
	}{
		"no cache file": {
			cacheFilePath: filepath.Join("testdata", "file_not_exists"),
			expectErr:     "failed to read cached SSO token file",
		},
		"invalid json cache file": {
			cacheFilePath: filepath.Join("testdata", "invalid_json.json"),
			expectErr:     "failed to parse cached SSO token file",
		},
		"missing accessToken": {
			cacheFilePath: filepath.Join("testdata", "missing_accessToken.json"),
			expectErr:     "must contain accessToken and expiresAt fields",
		},
		"missing expiresAt": {
			cacheFilePath: filepath.Join("testdata", "missing_expiresAt.json"),
			expectErr:     "must contain accessToken and expiresAt fields",
		},
		"expired no clientSecret": {
			cacheFilePath: filepath.Join("testdata", "missing_clientSecret.json"),
			expectErr:     "cached SSO token is expired, or not present",
		},
		"expired no clientId": {
			cacheFilePath: filepath.Join("testdata", "missing_clientId.json"),
			expectErr:     "cached SSO token is expired, or not present",
		},
		"expired no refreshToken": {
			cacheFilePath: filepath.Join("testdata", "missing_refreshToken.json"),
			expectErr:     "cached SSO token is expired, or not present",
		},
		"valid sso token": {
			cacheFilePath: filepath.Join("testdata", "valid_token.json"),
			expectToken: smithybearer.Token{
				Value:     "dGhpcyBpcyBub3QgYSByZWFsIHZhbHVl",
				CanExpire: true,
				Expires:   time.Date(2044, 4, 4, 7, 0, 1, 0, time.UTC),
			},
		},
		"refresh expired token": {
			setup: func() error {
				testFile, err := os.ReadFile(filepath.Join("testdata", "expired_token.json"))
				if err != nil {
					return err
				}

				return os.WriteFile(filepath.Join(tempDir, "expired_token.json"), testFile, 0600)
			},
			postRetrieve: func() error {
				actual, err := loadCachedToken(filepath.Join(tempDir, "expired_token.json"))
				if err != nil {
					return err

				}
				expect := token{
					tokenKnownFields: tokenKnownFields{
						AccessToken: "updated access token",
						ExpiresAt:   (*rfc3339)(aws.Time(time.Date(2021, 12, 21, 12, 31, 1, 0, time.UTC))),

						RefreshToken: "updated refresh token",
						ClientID:     "client id",
						ClientSecret: "client secret",
					},
					UnknownFields: map[string]interface{}{
						"unknownField": "some value",
					},
				}

				if diff := cmp.Diff(expect, actual, tokenCmpOptions...); diff != "" {
					return fmt.Errorf("expect token file match\n%s", diff)
				}
				return nil
			},
			cacheFilePath: filepath.Join(tempDir, "expired_token.json"),
			client: &mockCreateTokenAPIClient{
				expectInput: &ssooidc.CreateTokenInput{
					ClientId:     aws.String("client id"),
					ClientSecret: aws.String("client secret"),
					RefreshToken: aws.String("refresh token"),
					GrantType:    aws.String("refresh_token"),
				},
				output: &ssooidc.CreateTokenOutput{
					AccessToken:  aws.String("updated access token"),
					ExpiresIn:    600,
					RefreshToken: aws.String("updated refresh token"),
				},
			},
			expectToken: smithybearer.Token{
				Value:     "updated access token",
				CanExpire: true,
				Expires:   time.Date(2021, 12, 21, 12, 31, 1, 0, time.UTC),
			},
		},
		"fail refresh expired token": {
			setup: func() error {
				testFile, err := os.ReadFile(filepath.Join("testdata", "expired_token.json"))
				if err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(tempDir, "expired_token.json"), testFile, 0600)
			},
			postRetrieve: func() error {
				actual, err := loadCachedToken(filepath.Join(tempDir, "expired_token.json"))
				if err != nil {
					return err

				}
				expect := token{
					tokenKnownFields: tokenKnownFields{
						AccessToken: "access token",
						ExpiresAt:   (*rfc3339)(aws.Time(time.Date(2021, 12, 21, 12, 21, 1, 0, time.UTC))),

						RefreshToken: "refresh token",
						ClientID:     "client id",
						ClientSecret: "client secret",
					},
				}

				if diff := cmp.Diff(expect, actual, tokenCmpOptions...); diff != "" {
					return fmt.Errorf("expect token file match\n%s", diff)
				}
				return nil
			},
			cacheFilePath: filepath.Join(tempDir, "expired_token.json"),
			client: &mockCreateTokenAPIClient{
				err: fmt.Errorf("sky is falling"),
			},
			expectErr: "unable to refresh SSO token, sky is falling",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if c.setup != nil {
				if err := c.setup(); err != nil {
					t.Fatalf("failed to setup test, %v", err)
				}
			}
			provider := NewSSOTokenProvider(c.client, c.cacheFilePath, c.optFns...)

			token, err := provider.RetrieveBearerToken(context.Background())
			if c.expectErr != "" {
				if err == nil {
					t.Fatalf("expect %v error, got none", c.expectErr)
				}
				if e, a := c.expectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %v error, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if diff := cmp.Diff(c.expectToken, token, tokenCmpOptions...); diff != "" {
				t.Errorf("expect token match\n%s", diff)
			}

			if c.postRetrieve != nil {
				if err := c.postRetrieve(); err != nil {
					t.Fatalf("post retrieve failed, %v", err)
				}
			}
		})
	}
}

type mockCreateTokenAPIClient struct {
	expectInput *ssooidc.CreateTokenInput
	output      *ssooidc.CreateTokenOutput
	err         error
}

func (c *mockCreateTokenAPIClient) CreateToken(
	ctx context.Context, input *ssooidc.CreateTokenInput, optFns ...func(*ssooidc.Options)) (
	*ssooidc.CreateTokenOutput, error,
) {
	if c.expectInput != nil {
		opts := cmp.Options{
			cmpopts.IgnoreUnexported(ssooidc.CreateTokenInput{}),
		}
		if diff := cmp.Diff(c.expectInput, input, opts...); diff != "" {
			return nil, fmt.Errorf("expect input match\n%s", diff)
		}
	}

	return c.output, c.err
}
