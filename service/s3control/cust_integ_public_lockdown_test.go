// +build integration

package s3control_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/s3control/types"
)

func TestInteg_PublicAccessBlock(t *testing.T) {
	_, err := svc.GetPublicAccessBlockRequest(&types.GetPublicAccessBlockInput{
		AccountId: aws.String(accountID),
	}).Send(context.Background())
	if err != nil {
		aerr, ok := err.(awserr.RequestFailure)
		if !ok {
			t.Fatalf("unknown exception, %T, %v", err, err)
		}
		// Only no such configuration is valid error to receive.
		if e, a := s3control.ErrCodeNoSuchPublicAccessBlockConfiguration, aerr.Code(); e != a {
			t.Fatalf("expected no error, or no such configuration, got %v", err)
		}
	}

	_, err = svc.PutPublicAccessBlockRequest(&types.PutPublicAccessBlockInput{
		AccountId: aws.String(accountID),
		PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
			IgnorePublicAcls: aws.Bool(true),
		},
	}).Send(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	_, err = svc.DeletePublicAccessBlockRequest(&types.DeletePublicAccessBlockInput{
		AccountId: aws.String(accountID),
	}).Send(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
}
