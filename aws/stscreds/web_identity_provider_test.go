package stscreds_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/awserr"
	"github.com/jviney/aws-sdk-go-v2/aws/defaults"
	"github.com/jviney/aws-sdk-go-v2/aws/stscreds"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/internal/sdk"
	"github.com/jviney/aws-sdk-go-v2/service/sts"
)

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

	var reqCount int
	cases := map[string]struct {
		onSendReq         func(*testing.T, *aws.Request)
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
			onSendReq: func(t *testing.T, r *aws.Request) {
				input := r.Params.(*sts.AssumeRoleWithWebIdentityInput)
				if e, a := "foo", *input.RoleSessionName; !reflect.DeepEqual(e, a) {
					t.Errorf("expected %v, but received %v", e, a)
				}

				data := r.Data.(*sts.AssumeRoleWithWebIdentityOutput)
				*data = sts.AssumeRoleWithWebIdentityOutput{
					Credentials: &sts.Credentials{
						Expiration:      aws.Time(sdk.NowTime()),
						AccessKeyId:     aws.String("access-key-id"),
						SecretAccessKey: aws.String("secret-access-key"),
						SessionToken:    aws.String("session-token"),
					},
				}
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
		"invalid token retry": {
			roleARN:       "arn01234567890123456789",
			tokenFilepath: "testdata/token.jwt",
			sessionName:   "foo",
			onSendReq: func(t *testing.T, r *aws.Request) {
				input := r.Params.(*sts.AssumeRoleWithWebIdentityInput)
				if e, a := "foo", *input.RoleSessionName; !reflect.DeepEqual(e, a) {
					t.Errorf("expected %v, but received %v", e, a)
				}

				if reqCount == 0 {
					r.HTTPResponse.StatusCode = 400
					r.Error = awserr.New(sts.ErrCodeInvalidIdentityTokenException,
						"some error message", nil)
					return
				}

				data := r.Data.(*sts.AssumeRoleWithWebIdentityOutput)
				*data = sts.AssumeRoleWithWebIdentityOutput{
					Credentials: &sts.Credentials{
						Expiration:      aws.Time(sdk.NowTime()),
						AccessKeyId:     aws.String("access-key-id"),
						SecretAccessKey: aws.String("secret-access-key"),
						SessionToken:    aws.String("session-token"),
					},
				}
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
			reqCount = 0

			cfg := unit.Config()
			cfg.Logger = t

			client := sts.New(cfg)
			client.Handlers.Send.Swap(defaults.SendHandler.Name, aws.NamedHandler{
				Name: "custom send stub handler",
				Fn: func(r *aws.Request) {
					r.HTTPResponse = &http.Response{
						StatusCode: 200, Header: http.Header{},
					}
					c.onSendReq(t, r)
					reqCount++
				},
			})
			client.Handlers.UnmarshalMeta.Clear()
			client.Handlers.Unmarshal.Clear()
			client.Handlers.UnmarshalError.Clear()

			p := stscreds.NewWebIdentityRoleProvider(client, c.roleARN, c.sessionName, stscreds.IdentityTokenFile(c.tokenFilepath))
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
