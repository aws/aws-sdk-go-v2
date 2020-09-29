package stscreds_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type mockAssumeRoleWithWebIdentity func(ctx context.Context, params *sts.AssumeRoleWithWebIdentityInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleWithWebIdentityOutput, error)

func (m mockAssumeRoleWithWebIdentity) AssumeRoleWithWebIdentity(ctx context.Context, params *sts.AssumeRoleWithWebIdentityInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleWithWebIdentityOutput, error) {
	return m(ctx, params, optFns...)
}

type mockErrorCode string

func (m mockErrorCode) ErrorCode() string {
	return string(m)
}

func (m mockErrorCode) Error() string {
	return "error code: " + string(m)
}

func TestWebIdentityProviderRetrieve(t *testing.T) {
	defer func() func() {
		o := sdk.NowTime
		sdk.NowTime = func() time.Time {
			return time.Time{}
		}
		return func() {
			sdk.NowTime = o
		}
	}()()

	cases := map[string]struct {
		mockClient        mockAssumeRoleWithWebIdentity
		roleARN           string
		tokenFilepath     string
		sessionName       string
		expectedError     error
		expectedCredValue aws.Credentials
	}{
		"session name case": {
			roleARN:       "arn01234567890123456789",
			tokenFilepath: "testdata/token.jwt",
			sessionName:   "foo",
			mockClient: func(ctx context.Context, params *sts.AssumeRoleWithWebIdentityInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleWithWebIdentityOutput, error) {
				if e, a := "foo", *params.RoleSessionName; e != a {
					return nil, fmt.Errorf("expected %v, but received %v", e, a)
				}
				return &sts.AssumeRoleWithWebIdentityOutput{
					Credentials: &types.Credentials{
						Expiration:      aws.Time(sdk.NowTime()),
						AccessKeyId:     aws.String("access-key-id"),
						SecretAccessKey: aws.String("secret-access-key"),
						SessionToken:    aws.String("session-token"),
					},
				}, nil
			},
			expectedCredValue: aws.Credentials{
				AccessKeyID:     "access-key-id",
				SecretAccessKey: "secret-access-key",
				SessionToken:    "session-token",
				Source:          stscreds.WebIdentityProviderName,
				CanExpire:       true,
				Expires:         sdk.NowTime(),
			},
		},
		"configures token retry": {
			roleARN:       "arn01234567890123456789",
			tokenFilepath: "testdata/token.jwt",
			sessionName:   "foo",
			mockClient: func(ctx context.Context, params *sts.AssumeRoleWithWebIdentityInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleWithWebIdentityOutput, error) {
				o := sts.Options{}
				for _, fn := range optFns {
					fn(&o)
				}

				if o.Retryer.IsErrorRetryable(mockErrorCode("InvalidIdentityTokenException")) != true {
					return nil, fmt.Errorf("expected InvalidIdentityTokenException to be retryable")
				}

				return &sts.AssumeRoleWithWebIdentityOutput{
					Credentials: &types.Credentials{
						Expiration:      aws.Time(sdk.NowTime()),
						AccessKeyId:     aws.String("access-key-id"),
						SecretAccessKey: aws.String("secret-access-key"),
						SessionToken:    aws.String("session-token"),
					},
				}, nil
			},
			expectedCredValue: aws.Credentials{
				AccessKeyID:     "access-key-id",
				SecretAccessKey: "secret-access-key",
				SessionToken:    "session-token",
				Source:          stscreds.WebIdentityProviderName,
				CanExpire:       true,
				Expires:         sdk.NowTime(),
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			p := stscreds.NewWebIdentityRoleProvider(c.mockClient, c.roleARN, stscreds.IdentityTokenFile(c.tokenFilepath),
				func(o *stscreds.WebIdentityRoleOptions) {
					o.RoleSessionName = c.sessionName
				})
			credValue, err := p.Retrieve(context.Background())
			if e, a := c.expectedError, err; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}

			if e, a := c.expectedCredValue, credValue; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, but received %v", e, a)
			}
		})
	}
}
