// +build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// BucketPrefix is the root prefix of integration test buckets.
const BucketPrefix = "aws-sdk-go-v2-s3manager-integration"

// GenerateBucketName returns a unique bucket name.
func GenerateBucketName() string {
	return fmt.Sprintf("%s-%s",
		BucketPrefix, MustUUID())
}

// SetupBucket returns a test bucket created for the integration tests.
func SetupBucket(client *s3.Client, bucketName string) (err error) {
	fmt.Println("Setup: Creating test bucket,", bucketName)
	_, err = client.CreateBucket(context.Background(), &s3.CreateBucketInput{Bucket: &bucketName})
	if err != nil {
		return fmt.Errorf("failed to create bucket %s, %v", bucketName, err)
	}

	fmt.Println("Setup: Waiting for bucket to exist,", bucketName)
	err = waitUntilBucketExists(context.Background(), client, &s3.HeadBucketInput{Bucket: &bucketName})
	if err != nil {
		return fmt.Errorf("failed waiting for bucket %s to be created, %v",
			bucketName, err)
	}

	return nil
}

func waitUntilBucketExists(ctx context.Context, client *s3.Client, params *s3.HeadBucketInput) error {
	for i := 0; i < 20; i++ {
		_, err := client.HeadBucket(ctx, params)
		if err == nil {
			return nil
		}

		var httpErr interface{ HTTPStatusCode() int }

		if !errors.As(err, &httpErr) {
			return err
		}

		if httpErr.HTTPStatusCode() == http.StatusMovedPermanently || httpErr.HTTPStatusCode() == http.StatusForbidden {
			return nil
		}

		if httpErr.HTTPStatusCode() != http.StatusNotFound {
			return err
		}

		time.Sleep(5 * time.Second)
	}
	return nil
}

// CleanupBucket deletes the contents of a S3 bucket, before deleting the bucket
// it self.
func CleanupBucket(client *s3.Client, bucketName string) error {
	var errs []error

	{
		fmt.Println("TearDown: Deleting objects from test bucket,", bucketName)
		input := &s3.ListObjectsV2Input{Bucket: &bucketName}
		for {
			listObjectsV2, err := client.ListObjectsV2(context.Background(), input)
			if err != nil {
				return fmt.Errorf("failed to list objects, %w", err)
			}

			var delete types.Delete
			for _, content := range listObjectsV2.Contents {
				obj := content
				delete.Objects = append(delete.Objects, &types.ObjectIdentifier{Key: obj.Key})
			}

			deleteObjects, err := client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
				Bucket: &bucketName,
				Delete: &delete,
			})
			if err != nil {
				errs = append(errs, err)
			}
			for _, deleteError := range deleteObjects.Errors {
				errs = append(errs, fmt.Errorf("failed to delete %s, %s", aws.ToString(deleteError.Key), aws.ToString(deleteError.Message)))
			}

			if aws.ToBool(listObjectsV2.IsTruncated) {
				input.ContinuationToken = listObjectsV2.NextContinuationToken
			} else {
				break
			}
		}
	}

	{
		fmt.Println("TearDown: Deleting partial uploads from test bucket,", bucketName)

		input := &s3.ListMultipartUploadsInput{Bucket: &bucketName}
		for {
			uploads, err := client.ListMultipartUploads(context.Background(), input)
			if err != nil {
				return fmt.Errorf("failed to list multipart objects, %w", err)
			}

			for _, upload := range uploads.Uploads {
				client.AbortMultipartUpload(context.Background(), &s3.AbortMultipartUploadInput{
					Bucket:   &bucketName,
					Key:      upload.Key,
					UploadId: upload.UploadId,
				})
			}

			if aws.ToBool(uploads.IsTruncated) {
				input.KeyMarker = uploads.NextKeyMarker
				input.UploadIdMarker = uploads.NextUploadIdMarker
			} else {
				break
			}
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to delete objects, %s", errs)
	}

	fmt.Println("TearDown: Deleting test bucket,", bucketName)
	if _, err := client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{Bucket: &bucketName}); err != nil {
		return fmt.Errorf("failed to delete test bucket %s, %w", bucketName, err)
	}

	return nil
}
