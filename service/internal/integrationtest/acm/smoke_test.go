// +build integration

package acm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/awslabs/smithy-go"
)

func TestInteg_00_ListCertificates(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := acm.NewFromConfig(cfg)

	params := &acm.ListCertificatesInput{}
	_, err = client.ListCertificates(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
func TestInteg_01_GetCertificate(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := acm.NewFromConfig(cfg)

	params := &acm.GetCertificateInput{
		CertificateArn: aws.String("arn:aws:acm:region:123456789012:certificate/12345678-1234-1234-1234-123456789012"),
	}

	_, err = client.GetCertificate(ctx, params)
	if err == nil {
		t.Fatalf("expect request to fail")
	}

	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expect error to be API error, was not, %v", err)
	}
	if len(apiErr.ErrorCode()) == 0 {
		t.Errorf("expect non-empty error code")
	}
	if len(apiErr.ErrorMessage()) == 0 {
		t.Errorf("expect non-empty error message")
	}
}
