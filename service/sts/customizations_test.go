package sts_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func init() {
	cfg := unit.Config()
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	svc = sts.New(cfg)
}

var svc *sts.STS

func TestUnsignedRequest_AssumeRoleWithSAML(t *testing.T) {
	req := svc.AssumeRoleWithSAMLRequest(&sts.AssumeRoleWithSAMLInput{
		PrincipalArn:  aws.String("ARN01234567890123456789"),
		RoleArn:       aws.String("ARN01234567890123456789"),
		SAMLAssertion: aws.String("ASSERT"),
	})

	err := req.Sign()
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := "", req.HTTPRequest.Header.Get("Authorization"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestUnsignedRequest_AssumeRoleWithWebIdentity(t *testing.T) {
	req := svc.AssumeRoleWithWebIdentityRequest(&sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String("ARN01234567890123456789"),
		RoleSessionName:  aws.String("SESSION"),
		WebIdentityToken: aws.String("TOKEN"),
	})

	err := req.Sign()
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := "", req.HTTPRequest.Header.Get("Authorization"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}
