// +build example

package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
)

const (
	bucketName  = "myBucketName"
	keyName     = "myKeyName"
	accountID   = "123456789012"
	accessPoint = "accesspointname"
)

func main() {
	config, err := external.LoadDefaultAWSConfig(nil)
	if err != nil {
		panic(fmt.Errorf("failed to load aws config: %v", err))
	}

	s3Svc := s3.New(config)
	s3ControlSvc := s3control.New(config)

	// Create an S3 Bucket
	fmt.Println("create s3 bucket")
	req := s3Svc.CreateBucketRequest(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	_, err = req.Send(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to create bucket: %v", err))
	}

	// Wait for S3 Bucket to Exist
	fmt.Println("wait for s3 bucket to exist")
	err = s3Svc.WaitUntilBucketExists(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}, nil)
	if err != nil {
		panic(fmt.Sprintf("bucket failed to materialize: %v", err))
	}

	// Create an Access Point referring to the bucket
	fmt.Println("create an access point")
	createAccessPoint := s3ControlSvc.CreateAccessPointRequest(&s3control.CreateAccessPointInput{
		AccountId: aws.String(accountID),
		Bucket:    aws.String(bucketName),
		Name:      aws.String(accessPoint),
	})
	_, err = createAccessPoint.Send(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to create access point: %v", err))
	}

	// Use the SDK's ARN builder to create an ARN for the Access Point.
	apARN := arn.ARN{
		Partition: "aws",
		Service:   "s3",
		Region:    config.Region,
		AccountID: accountID,
		Resource:  "accesspoint/" + accessPoint,
	}

	// And Use Access Point ARN where bucket parameters are accepted
	fmt.Println("get object using access point")
	getObject := s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(apARN.String()),
		Key:    aws.String("somekey"),
	})
	getObjectResponse, err := getObject.Send(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed get object request: %v", err))
	}

	_, err = ioutil.ReadAll(getObjectResponse.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read object body: %v", err))
	}
}
