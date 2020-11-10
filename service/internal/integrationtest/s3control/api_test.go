// +build integration

package s3control

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/s3control/types"
)

func TestInteg_PublicAccessBlock(t *testing.T) {
	ctx := context.Background()
	_, err := svc.GetPublicAccessBlock(ctx, &s3control.GetPublicAccessBlockInput{
		AccountId: aws.String(accountID),
	})
	if err != nil {
		var e *types.NoSuchPublicAccessBlockConfiguration
		if !errors.As(err, &e) {
			t.Fatalf("expect no error for GetPublicAccessBlock, got %v", err)
		}
	}

	_, err = svc.PutPublicAccessBlock(ctx, &s3control.PutPublicAccessBlockInput{
		AccountId: aws.String(accountID),
		PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
			IgnorePublicAcls: true,
		},
	})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	_, err = svc.DeletePublicAccessBlock(ctx, &s3control.DeletePublicAccessBlockInput{
		AccountId: aws.String(accountID),
	})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
}
