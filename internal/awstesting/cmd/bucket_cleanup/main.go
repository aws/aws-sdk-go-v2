// +build integration

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Searches the buckets of an account that match the prefix, and deletes
// those buckets, and all objects within. Before deleting will prompt user
// to confirm bucket should be deleted. Positive confirmation is required.
//
// Usage:
//    go run deleteBuckets.go <bucketPrefix>
func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(err)
	}

	svc := s3.New(cfg)
	req := svc.ListBucketsRequest(&s3.ListBucketsInput{})
	buckets, err := req.Send()
	if err != nil {
		panic(fmt.Sprintf("failed to list buckets, %v", err))
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "bucket prefix required")
		os.Exit(1)
	}
	bucketPrefix := os.Args[1]

	var failed bool
	for _, b := range buckets.Buckets {
		bucket := aws.StringValue(b.Name)

		if !strings.HasPrefix(bucket, bucketPrefix) {
			continue
		}

		fmt.Printf("Delete bucket %q? [y/N]: ", bucket)
		var v string
		if _, err := fmt.Scanln(&v); err != nil || !(v == "Y" || v == "y") {
			fmt.Println("\tSkipping")
			continue
		}

		fmt.Println("\tDeleting")
		if err := deleteBucket(svc, bucket); err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete bucket %q, %v", bucket, err)
			failed = true
		}
	}

	if failed {
		os.Exit(1)
	}
}

func deleteBucket(svc *s3.S3, bucket string) error {
	bucketName := &bucket

	listReq := svc.ListObjectsRequest(&s3.ListObjectsInput{Bucket: bucketName})
	objs, err := listReq.Send()
	if err != nil {
		return fmt.Errorf("failed to list bucket %q objects, %v", bucketName, err)
	}

	for _, o := range objs.Contents {
		delReq := svc.DeleteObjectRequest(&s3.DeleteObjectInput{Bucket: bucketName, Key: o.Key})
		delReq.Send()
	}

	listMulReq := svc.ListMultipartUploadsRequest(&s3.ListMultipartUploadsInput{Bucket: bucketName})
	uploads, err := listMulReq.Send()
	if err != nil {
		return fmt.Errorf("failed to list bucket %q multipart objects, %v", bucketName, err)
	}

	for _, u := range uploads.Uploads {
		abortReq := svc.AbortMultipartUploadRequest(&s3.AbortMultipartUploadInput{
			Bucket:   bucketName,
			Key:      u.Key,
			UploadId: u.UploadId,
		})
		abortReq.Send()
	}

	_, err = svc.DeleteBucketRequest(&s3.DeleteBucketInput{Bucket: bucketName}).Send()
	if err != nil {
		return fmt.Errorf("failed to delete bucket %q, %v", bucketName, err)
	}

	return nil
}
