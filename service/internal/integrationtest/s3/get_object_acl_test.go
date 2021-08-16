//go:build integration
// +build integration

package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestInteg_XSIType(t *testing.T) {
	key := integrationtest.UniqueID()

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &setupMetadata.Buckets.Source.Name,
		Key:    &key,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
		Bucket: &setupMetadata.Buckets.Source.Name,
		Key:    &key,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Grants) == 0 {
		t.Fatalf("expect Grants to not be empty")
	}

	if len(resp.Grants[0].Grantee.Type) == 0 {
		t.Errorf("expect grantee type to not be empty")
	}
}
