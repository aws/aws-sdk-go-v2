package cognitoidentity_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

var svc = cognitoidentity.New(unit.Config())

func TestUnsignedRequest_GetID(t *testing.T) {
	req := svc.GetIdRequest(&cognitoidentity.GetIdInput{
		IdentityPoolId: aws.String("IdentityPoolId"),
	})

	err := req.Sign()
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := "", req.HTTPRequest.Header.Get("Authorization"); e != a {
		t.Errorf("expected empty value '%v', but received, %v", e, a)
	}
}

func TestUnsignedRequest_GetOpenIDToken(t *testing.T) {
	req := svc.GetOpenIdTokenRequest(&cognitoidentity.GetOpenIdTokenInput{
		IdentityId: aws.String("IdentityId"),
	})

	err := req.Sign()
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := "", req.HTTPRequest.Header.Get("Authorization"); e != a {
		t.Errorf("expected empty value '%v', but received, %v", e, a)
	}
}

func TestUnsignedRequest_GetCredentialsForIdentity(t *testing.T) {
	req := svc.GetCredentialsForIdentityRequest(&cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: aws.String("IdentityId"),
	})

	err := req.Sign()
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := "", req.HTTPRequest.Header.Get("Authorization"); e != a {
		t.Errorf("expected empty value '%v', but received, %v", e, a)
	}
}
