package ssocreds

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/google/go-cmp/cmp"
)

type mockClient struct {
	t *testing.T

	Output *sso.GetRoleCredentialsOutput
	Err    error

	ExpectedAccountID    string
	ExpectedAccessToken  string
	ExpectedRoleName     string

	Response func(mockClient) (*sso.GetRoleCredentialsOutput, error)
}

func (m mockClient) GetRoleCredentials(ctx context.Context, params *sso.GetRoleCredentialsInput, optFns ...func(options *sso.Options)) (out *sso.GetRoleCredentialsOutput, err error) {
	m.t.Helper()

	if len(m.ExpectedAccountID) > 0 {
		if diff := cmp.Diff(m.ExpectedAccountID, aws.ToString(params.AccountId)); len(diff) > 0 {
			m.t.Error(diff)
		}
	}

	if len(m.ExpectedAccessToken) > 0 {
		if diff := cmp.Diff(m.ExpectedAccessToken, aws.ToString(params.AccessToken)); len(diff) > 0 {
			m.t.Error(diff)
		}
	}

	if len(m.ExpectedRoleName) > 0 {
		if diff := cmp.Diff(m.ExpectedRoleName, aws.ToString(params.RoleName)); len(diff) > 0 {
			m.t.Error(diff)
		}
	}

	if m.Response == nil {
		return out, err
	}
	return m.Response(m)
}

func swapCacheLocation(dir string) func() {
	original := defaultCacheLocation
	defaultCacheLocation = func() string {
		return dir
	}
	return func() {
		defaultCacheLocation = original
	}
}

func TestProvider(t *testing.T) {
	restoreCache := swapCacheLocation("testdata")
	defer restoreCache()

	restoreTime := sdk.TestingUseReferenceTime(time.Date(2021, 01, 19, 19, 50, 0, 0, time.UTC))
	defer restoreTime()

	cases := map[string]struct {
		Client    mockClient
		AccountID string
		Region    string
		RoleName  string
		StartURL  string
		Options   []func(*Options)

		ExpectedErr         bool
		ExpectedCredentials aws.Credentials
	}{
		"missing required parameter values": {
			StartURL:    "https://invalid-required",
			ExpectedErr: true,
		},
		"valid required parameter values": {
			Client: mockClient{
				ExpectedAccountID:    "012345678901",
				ExpectedRoleName:     "TestRole",
				ExpectedAccessToken:  "dGhpcyBpcyBub3QgYSByZWFsIHZhbHVl",
				Response: func(mock mockClient) (*sso.GetRoleCredentialsOutput, error) {
					return &sso.GetRoleCredentialsOutput{
						RoleCredentials: &types.RoleCredentials{
							AccessKeyId:     aws.String("AccessKey"),
							SecretAccessKey: aws.String("SecretKey"),
							SessionToken:    aws.String("SessionToken"),
							Expiration:      1611177743123,
						},
					}, nil
				},
			},
			AccountID: "012345678901",
			Region:    "us-west-2",
			RoleName:  "TestRole",
			StartURL:  "https://valid-required-only",
			ExpectedCredentials: aws.Credentials{
				AccessKeyID:     "AccessKey",
				SecretAccessKey: "SecretKey",
				SessionToken:    "SessionToken",
				CanExpire:       true,
				Expires:         time.Date(2021, 01, 20, 21, 22, 23, 0.123e9, time.UTC),
				Source:          ProviderName,
			},
		},
		"expired access token": {
			StartURL:    "https://expired",
			ExpectedErr: true,
		},
		"api error": {
			Client: mockClient{
				ExpectedAccountID:    "012345678901",
				ExpectedRoleName:     "TestRole",
				ExpectedAccessToken:  "dGhpcyBpcyBub3QgYSByZWFsIHZhbHVl",
				Response: func(mock mockClient) (*sso.GetRoleCredentialsOutput, error) {
					return nil, fmt.Errorf("api error")
				},
			},
			AccountID:   "012345678901",
			Region:      "us-west-2",
			RoleName:    "TestRole",
			StartURL:    "https://valid-required-only",
			ExpectedErr: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			tt.Client.t = t

			provider := New(tt.Client, tt.AccountID, tt.RoleName, tt.StartURL, tt.Options...)

			credentials, err := provider.Retrieve(context.Background())
			if (err != nil) != tt.ExpectedErr {
				t.Errorf("expect error: %v", tt.ExpectedErr)
			}

			if diff := cmp.Diff(tt.ExpectedCredentials, credentials); len(diff) > 0 {
				t.Errorf(diff)
			}
		})
	}
}
