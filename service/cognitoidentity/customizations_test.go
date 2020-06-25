// +build go1.8

package cognitoidentity_test

import (
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/service/cognitoidentity"
)

var svc = cognitoidentity.New(unit.Config())

func TestUnsignedRequests(t *testing.T) {
	type reqSigner interface {
		Sign() error
	}

	cases := map[string]struct {
		ReqFn func() reqSigner
	}{
		"GetId": {
			ReqFn: func() reqSigner {
				req := svc.GetIdRequest(&cognitoidentity.GetIdInput{
					IdentityPoolId: aws.String("IdentityPoolId"),
				})
				return req
			},
		},
		"GetOpenIdToken": {
			ReqFn: func() reqSigner {
				req := svc.GetOpenIdTokenRequest(&cognitoidentity.GetOpenIdTokenInput{
					IdentityId: aws.String("IdentityId"),
				})
				return req
			},
		},
		"UnlinkIdentity": {
			ReqFn: func() reqSigner {
				req := svc.UnlinkIdentityRequest(&cognitoidentity.UnlinkIdentityInput{
					IdentityId:     aws.String("IdentityId"),
					Logins:         map[string]string{},
					LoginsToRemove: []string{},
				})
				return req
			},
		},
		"GetCredentialsForIdentity": {
			ReqFn: func() reqSigner {
				req := svc.GetCredentialsForIdentityRequest(&cognitoidentity.GetCredentialsForIdentityInput{
					IdentityId: aws.String("IdentityId"),
				})
				return req
			},
		},
	}

	for cn, c := range cases {
		t.Run(cn, func(t *testing.T) {
			req := c.ReqFn()
			err := req.Sign()
			if err != nil {
				t.Errorf("expected no error, but received %v", err)
			}

			switch tv := req.(type) {
			case cognitoidentity.GetIdRequest:
				if e, a := aws.AnonymousCredentials, tv.Config.Credentials; e != a {
					t.Errorf("expect request to use anonymous credentias, %v", a)
				}
				if e, a := "", tv.HTTPRequest.Header.Get("Authorization"); e != a {
					t.Errorf("expected empty value '%v', but received, %v", e, a)
				}
			}

		})
	}
}
