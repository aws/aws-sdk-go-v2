// +build integration

package s3_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/s3integ"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const integBucketPrefix = "aws-sdk-go-v2-integration"

var bucketName *string
var svc *s3.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	sess := integration.ConfigWithDefaultRegion("us-west-2")
	svc = s3.New(sess)
	bucketName = aws.String(s3integ.GenerateBucketName())
	if err := s3integ.SetupTest(ctx, svc, *bucketName); err != nil {
		panic(err)
	}

	var result int
	defer func() {
		if err := s3integ.CleanupTest(ctx, svc, *bucketName); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "S3 integrationt tests paniced,", r)
			result = 1
		}
		os.Exit(result)
	}()

	result = m.Run()
}

func putTestFile(ctx context.Context, t *testing.T, filename, key string) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open testfile, %v", err)
	}
	defer f.Close()

	putTestContent(ctx, t, f, key)
}

func putTestContent(ctx context.Context, t *testing.T, reader io.ReadSeeker, key string) {
	fmt.Println(bucketName, key, svc)
	req := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: bucketName,
		Key:    aws.String(key),
		Body:   reader,
	})
	if _, err := req.Send(ctx); err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
