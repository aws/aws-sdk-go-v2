// +build integration

package marketplacecommerceanalytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/marketplacecommerceanalytics"
	"github.com/aws/aws-sdk-go-v2/service/marketplacecommerceanalytics/types"
	"github.com/awslabs/smithy-go"
	smithytime "github.com/awslabs/smithy-go/time"
)

func TestInteg_00_GenerateDataSet(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-east-1")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := marketplacecommerceanalytics.NewFromConfig(cfg)
	params := &marketplacecommerceanalytics.GenerateDataSetInput{
		DataSetPublicationDate:  aws.Time(smithytime.ParseEpochSeconds(0.000000)),
		DataSetType:             types.DataSetTypeDailyBusinessFees,
		DestinationS3BucketName: aws.String("fake-bucket"),
		RoleNameArn:             aws.String("fake-arn"),
		SnsTopicArn:             aws.String("fake-arn"),
	}
	_, err = client.GenerateDataSet(ctx, params)
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
