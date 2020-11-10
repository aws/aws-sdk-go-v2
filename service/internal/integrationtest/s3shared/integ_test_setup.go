// +build integration

package s3shared

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
)

// BucketPrefix is the root prefix of integration test buckets.
const BucketPrefix = "aws-sdk-go-v2-integration"

// GenerateBucketName returns a unique bucket name.
func GenerateBucketName() string {
	return fmt.Sprintf("%s-%s",
		BucketPrefix, integrationtest.UniqueID())
}

// SetupBucket returns a test bucket created for the integration tests.
func SetupBucket(ctx context.Context, svc *s3.Client, bucketName string) (err error) {
	fmt.Println("Setup: Creating test bucket,", bucketName)
	_, err = svc.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintUsWest2,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket %s, %v", bucketName, err)
	}

	// TODO: change this to use waiter to wait until BucketExists instead of loop
	// 	svc.WaitUntilBucketExists(HeadBucketInput)

	// HeadBucket to determine if bucket exists
	var attempt = 0
	params := &s3.HeadBucketInput{
		Bucket: &bucketName,
	}
pt:
	_, err = svc.HeadBucket(ctx, params)
	// increment an attempt
	attempt++

	// retry till 10 attempt
	if err != nil {
		if attempt < 10 {
			goto pt
		}
		// fail if not succeed after 10 attempts
		return fmt.Errorf("failed to determine if a bucket %s exists and you have permission to access it", bucketName)
	}

	return nil
}

// CleanupBucket deletes the contents of a S3 bucket, before deleting the bucket
// it self.
// TODO: list and delete methods should use paginators
func CleanupBucket(ctx context.Context, svc *s3.Client, bucketName string) (err error) {
	var errs = make([]error, 0)

	fmt.Println("TearDown: Deleting objects from test bucket,", bucketName)
	listObjectsResp, err := svc.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return fmt.Errorf("failed to list objects, %w", err)
	}

	for _, o := range listObjectsResp.Contents {
		_, err := svc.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: &bucketName,
			Key:    o.Key,
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("failed to delete objects, %s", errs)
	}

	fmt.Println("TearDown: Deleting partial uploads from test bucket,", bucketName)
	multipartUploadResp, err := svc.ListMultipartUploads(ctx, &s3.ListMultipartUploadsInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return fmt.Errorf("failed to list multipart objects, %w", err)
	}

	for _, u := range multipartUploadResp.Uploads {
		_, err = svc.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket:   &bucketName,
			Key:      u.Key,
			UploadId: u.UploadId,
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("failed to delete multipart upload objects, %s", errs)
	}

	fmt.Println("TearDown: Deleting test bucket,", bucketName)
	_, err = svc.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return fmt.Errorf("failed to delete bucket, %s", bucketName)
	}

	return nil
}

// SetupAccessPoint returns an access point for the given bucket for testing
func SetupAccessPoint(ctx context.Context, svc *s3control.Client, accountID, bucketName, accessPoint string) error {
	fmt.Printf("Setup: creating access point %q for bucket %q\n", accessPoint, bucketName)
	_, err := svc.CreateAccessPoint(ctx, &s3control.CreateAccessPointInput{
		AccountId: &accountID,
		Bucket:    &bucketName,
		Name:      &accessPoint,
	})
	if err != nil {
		return fmt.Errorf("failed to create access point : %w", err)
	}
	return nil
}

// CleanupAccessPoint deletes the given access point
func CleanupAccessPoint(ctx context.Context, svc *s3control.Client, accountID, accessPoint string) error {
	fmt.Printf("TearDown: Deleting access point %q\n", accessPoint)
	_, err := svc.DeleteAccessPoint(ctx, &s3control.DeleteAccessPointInput{
		AccountId: &accountID,
		Name:      &accessPoint,
	})
	if err != nil {
		return fmt.Errorf("failed to delete access point: %w", err)
	}
	return nil
}
