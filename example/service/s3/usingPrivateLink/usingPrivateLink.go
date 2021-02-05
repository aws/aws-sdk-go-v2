package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
)

const (
	bucketName  = "myBucketName"
	accountID   = "123456789012"
	accessPoint = "accesspointname"

	// vpcBucketEndpoint will be used by the SDK to resolve an endpoint, when making a call to
	// access `bucket` data using s3 interface endpoint. This endpoint may be mutated by the SDK,
	// as per the input provided to work with ARNs.
	vpcBucketEndpoint = "https://bucket.vpce-0xxxxxxx-xxx8xxg.s3.us-west-2.vpce.amazonaws.com"

	// vpcAccesspointEndpoint will be used by the SDK to resolve an endpoint, when making a call to
	// access `access-point` data using s3 interface endpoint. This endpoint may be mutated by the SDK,
	// as per the input provided to work with ARNs.
	vpcAccesspointEndpoint = "https://accesspoint.vpce-0xxxxxxx-xxx8xxg.s3.us-west-2.vpce.amazonaws.com"

	// vpcControlEndpoint will be used by the SDK to resolve an endpoint, when making a call to
	// access `control` data using s3 interface endpoint. This endpoint may be mutated by the SDK,
	// as per the input provided to work with ARNs.
	vpcControlEndpoint = "https://control.vpce-0xxxxxxx-xxx8xxg.s3.us-west-2.vpce.amazonaws.com"
)

func main() {
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	// Load the SDK's configuration from environment and shared config, and
	// create the client with this.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	s3controlClient := s3control.NewFromConfig(cfg)

	// Create an S3 Bucket
	fmt.Println("create s3 bucket")

	setVPCBucketEndpoint := s3.WithEndpointResolver(s3.EndpointResolverFromURL(vpcBucketEndpoint))
	createBucketParams := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	_, err = s3Client.CreateBucket(context.TODO(), createBucketParams, setVPCBucketEndpoint)
	if err != nil {
		panic(fmt.Errorf("failed to create bucket: %v", err))
	}

	// Wait for S3 Bucket to Exist
	fmt.Println("wait for s3 bucket to exist")
	waiter := s3.NewBucketExistsWaiter(s3Client)
	err = waiter.Wait(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}, 120*time.Second)
	if err != nil {
		panic(fmt.Sprintf("bucket failed to materialize: %v", err))
	}

	// Create an Access Point referring to the bucket
	fmt.Println("create an access point")

	setVpcControlEndpoint := s3control.WithEndpointResolver(s3control.EndpointResolverFromURL(vpcControlEndpoint))
	createAccesspointInput := &s3control.CreateAccessPointInput{
		AccountId: aws.String(accountID),
		Bucket:    aws.String(bucketName),
		Name:      aws.String(accessPoint),
	}
	_, err = s3controlClient.CreateAccessPoint(context.TODO(), createAccesspointInput, setVpcControlEndpoint)
	if err != nil {
		panic(fmt.Sprintf("failed to create access point: %v", err))
	}

	// build an arn
	apARN := arn.ARN{
		Partition: "aws",
		Service:   "s3",
		Region:    cfg.Region,
		AccountID: accountID,
		Resource:  "accesspoint/" + accessPoint,
	}

	// get object using access point ARN
	fmt.Println("get object using access point")

	setVPCAccesspointEndpoint := s3.WithEndpointResolver(s3.EndpointResolverFromURL(vpcAccesspointEndpoint))
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(apARN.String()),
		Key:    aws.String("somekey"),
	}

	getObjectOutput, err := s3Client.GetObject(context.TODO(), getObjectInput, setVPCAccesspointEndpoint)
	if err != nil {
		panic(fmt.Sprintf("failed get object request: %v", err))
	}

	_, err = ioutil.ReadAll(getObjectOutput.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read object body: %v", err))
	}
}
