// +build integration

package s3integ

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/s3control/s3controliface"
)

// BucketPrefix is the root prefix of integration test buckets.
const BucketPrefix = "aws-sdk-go-v2-integration"

// GenerateBucketName returns a unique bucket name.
func GenerateBucketName() string {
	return fmt.Sprintf("%s-%s",
		BucketPrefix, integration.UniqueID())
}

// SetupBucket returns a test bucket created for the integration tests.
func SetupBucket(ctx context.Context, svc s3iface.ClientAPI, bucketName string) (err error) {
	fmt.Println("Setup: Creating test bucket,", bucketName)
	_, err = svc.CreateBucketRequest(&s3.CreateBucketInput{Bucket: &bucketName}).Send(ctx)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s, %v", bucketName, err)
	}

	fmt.Println("Setup: Waiting for bucket to exist,", bucketName)
	err = svc.WaitUntilBucketExists(ctx, &s3.HeadBucketInput{Bucket: &bucketName})
	if err != nil {
		return fmt.Errorf("failed waiting for bucket %s to be created, %v",
			bucketName, err)
	}

	return nil
}

// CleanupBucket deletes the contents of a S3 bucket, before deleting the bucket
// it self.
func CleanupBucket(ctx context.Context, svc s3iface.ClientAPI, bucketName string) error {
	var errs []error

	fmt.Println("TearDown: Deleting objects from test bucket,", bucketName)
	listReq := svc.ListObjectsRequest(
		&s3.ListObjectsInput{Bucket: &bucketName},
	)

	listObjPager := s3.NewListObjectsPaginator(listReq)
	for listObjPager.Next(ctx) {
		for _, o := range listObjPager.CurrentPage().Contents {
			_, err := svc.DeleteObjectRequest(&s3.DeleteObjectInput{
				Bucket: &bucketName,
				Key:    o.Key,
			}).Send(ctx)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	if err := listObjPager.Err(); err != nil {
		return fmt.Errorf("failed to list objects, %s, %v", bucketName, err)
	}

	fmt.Println("TearDown: Deleting partial uploads from test bucket,",
		bucketName)
	listMPReq := svc.ListMultipartUploadsRequest(
		&s3.ListMultipartUploadsInput{Bucket: &bucketName},
	)

	listMPPager := s3.NewListMultipartUploadsPaginator(listMPReq)
	for listMPPager.Next(ctx) {
		for _, u := range listMPPager.CurrentPage().Uploads {
			svc.AbortMultipartUploadRequest(&s3.AbortMultipartUploadInput{
				Bucket:   &bucketName,
				Key:      u.Key,
				UploadId: u.UploadId,
			}).Send(ctx)
		}
	}
	if err := listMPPager.Err(); err != nil {
		return fmt.Errorf("failed to list multipart objects, %s, %v",
			bucketName, err)
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to delete objects, %s", errs)
	}

	fmt.Println("TearDown: Deleting test bucket,", bucketName)
	if _, err := svc.DeleteBucketRequest(&s3.DeleteBucketInput{
		Bucket: &bucketName,
	}).Send(ctx); err != nil {
		return fmt.Errorf("failed to delete test bucket, %s", bucketName)
	}

	return nil
}

// SetupAccessPoint returns an access point for the given bucket for testing
func SetupAccessPoint(svc s3controliface.ClientAPI, account, bucket, accessPoint string) error {
	fmt.Printf("Setup: creating access point %q for bucket %q\n", accessPoint, bucket)
	req := svc.CreateAccessPointRequest(&s3control.CreateAccessPointInput{
		AccountId: &account,
		Bucket:    &bucket,
		Name:      &accessPoint,
	})
	_, err := req.Send(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create access point: %v", err)
	}
	return nil
}

// CleanupAccessPoint deletes the given access point
func CleanupAccessPoint(svc s3controliface.ClientAPI, account, accessPoint string) error {
	fmt.Printf("TearDown: Deleting access point %q\n", accessPoint)
	req := svc.DeleteAccessPointRequest(&s3control.DeleteAccessPointInput{
		AccountId: &account,
		Name:      &accessPoint,
	})
	_, err := req.Send(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete access point: %v", err)
	}
	return nil
}
