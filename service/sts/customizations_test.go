package sts_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

func init() {
	cfg := unit.Config()
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	svc = sts.New(cfg)
}

var svc *sts.Client

func TestUnsignedRequest_AssumeRoleWithSAML(t *testing.T) {
	req := svc.AssumeRoleWithSAMLRequest(&types.AssumeRoleWithSAMLInput{
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
	req := svc.AssumeRoleWithWebIdentityRequest(&types.AssumeRoleWithWebIdentityInput{
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

func TestSTSCustomRetryErrorCodes(t *testing.T) {
	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 1
	})

	svc := sts.New(cfg)
	svc.Handlers.Validate.Clear()

	const xmlErr = `<ErrorResponse><Error><Code>%s</Code><Message>some error message</Message></Error></ErrorResponse>`
	var reqCount int
	resps := []*http.Response{
		{
			StatusCode: 400,
			Header:     http.Header{},
			Body: ioutil.NopCloser(bytes.NewReader(
				[]byte(fmt.Sprintf(xmlErr, sts.ErrCodeIDPCommunicationErrorException)),
			)),
		},
		{
			StatusCode: 200,
			Header:     http.Header{},
			Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
		},
	}

	req := svc.AssumeRoleWithWebIdentityRequest(&types.AssumeRoleWithWebIdentityInput{})
	req.Handlers.Send.Swap(defaults.SendHandler.Name, aws.NamedHandler{
		Name: "custom send handler",
		Fn: func(r *aws.Request) {
			r.HTTPResponse = resps[reqCount]
			reqCount++
		},
	})

	if _, err := req.Send(context.Background()); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := 2, reqCount; e != a {
		t.Errorf("expect %v requests, got %v", e, a)
	}
}
