// +build integration

package support

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/support"
	"github.com/awslabs/smithy-go"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
)

func TestInteg_00_DescribeServices(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-east-1")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := support.NewFromConfig(cfg)
	params := &support.DescribeServicesInput{}
	_, err = client.DescribeServices(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestInteg_01_CreateCase(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-east-1")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := support.NewFromConfig(cfg)
	params := &support.CreateCaseInput{
		CategoryCode:      aws.String("category"),
		CommunicationBody: aws.String("communication"),
		ServiceCode:       aws.String("amazon-dynamodb"),
		SeverityCode:      aws.String("low"),
		Subject:           aws.String("subject"),
	}
	_, err = client.CreateCase(ctx, params)
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
