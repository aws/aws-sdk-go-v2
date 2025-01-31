//go:build integration
// +build integration

package transfermanager

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var setupMetadata = struct {
	AccountID string
	Region    string
	Buckets   struct {
		Source struct {
			Name string
			ARN  string
		}
	}
}{}

// s3 client to use for integ testing
var s3Client *s3.Client

// s3TransferManagerClient to use for integ testing
var s3TransferManagerClient *Client

// sts client to use for integ testing
var stsClient *sts.Client

// http client setting to use for integ testing
var httpClient *http.Client

var region = "us-west-2"

// large object buffer to test multipart upload
var largeObjectBuf []byte

// TestMain executes at start of package tests
func TestMain(m *testing.M) {
	flag.Parse()
	flag.CommandLine.Visit(func(f *flag.Flag) {
		if !(f.Name == "run" || f.Name == "test.run") {
			return
		}
		value := f.Value.String()
		if value == `NONE` {
			os.Exit(0)
		}
	})

	var result int
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "S3 TransferManager integration tests panic,", r)
			result = 1
		}
		os.Exit(result)
	}()

	var verifyTLS bool
	var s3Endpoint string

	flag.StringVar(&s3Endpoint, "s3-endpoint", "", "integration endpoint for S3")

	flag.StringVar(&setupMetadata.AccountID, "account", "", "integration account id")
	flag.BoolVar(&verifyTLS, "verify-tls", true, "verify server TLS certificate")
	flag.Parse()

	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: verifyTLS},
		},
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred while loading config with region %v, %v", region, err)
		result = 1
		return
	}

	// assign the http client
	cfg.HTTPClient = httpClient

	// create a s3 client
	s3cfg := cfg.Copy()
	if len(s3Endpoint) != 0 {
		s3cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           s3Endpoint,
				PartitionID:   "aws",
				SigningName:   "s3",
				SigningRegion: region,
				SigningMethod: "s3v4",
			}, nil
		})
	}

	// build s3 client from config
	s3Client = s3.NewFromConfig(s3cfg)

	// build s3 transfermanager client from config
	s3TransferManagerClient = NewFromConfig(s3Client, s3cfg)

	// build sts client from config
	stsClient = sts.NewFromConfig(cfg)

	// context
	ctx := context.Background()

	setupMetadata.AccountID, err = getAccountID(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get integration aws account id: %v\n", err)
		result = 1
		return
	}

	bucketCleanup, err := setupBuckets(ctx)
	defer bucketCleanup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup integration test buckets: %v\n", err)
		result = 1
		return
	}

	largeObjectBuf = make([]byte, 20*1024*1024)
	_, err = rand.Read(largeObjectBuf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate large object for multipart upload: %v\n", err)
		result = 1
		return
	}

	result = m.Run()
}

// getAccountID retrieves account id
func getAccountID(ctx context.Context) (string, error) {
	if len(setupMetadata.AccountID) != 0 {
		return setupMetadata.AccountID, nil
	}
	identity, err := stsClient.GetCallerIdentity(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("error fetching caller identity, %w", err)
	}
	return *identity.Account, nil
}

// setupBuckets creates buckets needed for integration test
func setupBuckets(ctx context.Context) (func(), error) {
	var cleanups []func()

	cleanup := func() {
		for i := range cleanups {
			cleanups[i]()
		}
	}

	bucketCreates := []struct {
		name *string
		arn  *string
	}{
		{name: &setupMetadata.Buckets.Source.Name, arn: &setupMetadata.Buckets.Source.ARN},
	}

	for _, bucket := range bucketCreates {
		*bucket.name = GenerateBucketName()

		if err := SetupBucket(ctx, s3Client, *bucket.name); err != nil {
			return cleanup, err
		}

		// Compute ARN
		bARN := arn.ARN{
			Partition: "aws",
			Service:   "s3",
			Region:    region,
			AccountID: setupMetadata.AccountID,
			Resource:  fmt.Sprintf("bucket_name:%s", *bucket.name),
		}.String()

		*bucket.arn = bARN

		bucketName := *bucket.name
		cleanups = append(cleanups, func() {
			if err := CleanupBucket(ctx, s3Client, bucketName); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		})
	}

	return cleanup, nil
}

type putObjectTestData struct {
	Body        io.Reader
	ExpectBody  []byte
	ExpectError string
}

// UniqueID returns a unique UUID-like identifier for use in generating
// resources for integration tests.
//
// TODO: duped from service/internal/integrationtest, remove after beta.
func UniqueID() string {
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	return fmt.Sprintf("%x", uuid)
}

func testPutObject(t *testing.T, bucket string, testData putObjectTestData, opts ...func(options *Options)) {
	key := UniqueID()

	_, err := s3TransferManagerClient.PutObject(context.Background(),
		&PutObjectInput{
			Bucket: bucket,
			Key:    key,
			Body:   testData.Body,
		}, opts...)
	if err != nil {
		if len(testData.ExpectError) == 0 {
			t.Fatalf("expect no error, got %v", err)
		}
		if e, a := testData.ExpectError, err.Error(); !strings.Contains(a, e) {
			t.Fatalf("expect error to contain %v, got %v", e, a)
		}
	} else {
		if e := testData.ExpectError; len(e) != 0 {
			t.Fatalf("expect error: %v, got none", e)
		}
	}

	if len(testData.ExpectError) != 0 {
		return
	}

	resp, err := s3Client.GetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	if e, a := testData.ExpectBody, b; !bytes.EqualFold(e, a) {
		t.Errorf("expect %s, got %s", e, a)
	}
}

// TODO: duped from service/internal/integrationtest, remove after beta.
const expressAZID = "usw2-az3"

// TODO: duped from service/internal/integrationtest, remove after beta.
const expressSuffix = "--usw2-az3--x-s3"

// BucketPrefix is the root prefix of integration test buckets.
//
// TODO: duped from service/internal/integrationtest, remove after beta.
const BucketPrefix = "aws-sdk-go-v2-integration"

// GenerateBucketName returns a unique bucket name.
//
// TODO: duped from service/internal/integrationtest, remove after beta.
func GenerateBucketName() string {
	return fmt.Sprintf("%s-%s",
		BucketPrefix, UniqueID())
}

// GenerateBucketName returns a unique express-formatted bucket name.
//
// TODO: duped from service/internal/integrationtest, remove after beta.
func GenerateExpressBucketName() string {
	return fmt.Sprintf(
		"%s-%s%s",
		BucketPrefix,
		UniqueID()[0:8], // express suffix adds length, regain that here
		expressSuffix,
	)
}

// SetupBucket returns a test bucket created for the integration tests.
//
// TODO: duped from service/internal/integrationtest, remove after beta.
func SetupBucket(ctx context.Context, svc *s3.Client, bucketName string) (err error) {
	fmt.Println("Setup: Creating test bucket,", bucketName)
	_, err = svc.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: "us-west-2",
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
		return fmt.Errorf("failed to determine if a bucket %s exists and you have permission to access it %v", bucketName, err)
	}

	return nil
}

// CleanupBucket deletes the contents of a S3 bucket, before deleting the bucket
// it self.
// TODO: list and delete methods should use paginators
//
// TODO: duped from service/internal/integrationtest, remove after beta.
func CleanupBucket(ctx context.Context, svc *s3.Client, bucketName string) (err error) {
	var errs = make([]error, 0)

	fmt.Println("TearDown: Deleting objects from test bucket,", bucketName)
	listObjectsResp, err := svc.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
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

// SetupExpressBucket returns an express bucket for testing.
//
// TODO: duped from service/internal/integrationtest, remove after beta.
func SetupExpressBucket(ctx context.Context, svc *s3.Client, bucketName string) error {
	if !strings.HasSuffix(bucketName, expressSuffix) {
		return fmt.Errorf("bucket name %s is missing required suffix %s", bucketName, expressSuffix)
	}

	fmt.Println("Setup: Creating test express bucket,", bucketName)
	_, err := svc.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			Location: &types.LocationInfo{
				Name: aws.String(expressAZID),
				Type: types.LocationTypeAvailabilityZone,
			},
			Bucket: &types.BucketInfo{
				DataRedundancy: types.DataRedundancySingleAvailabilityZone,
				Type:           types.BucketTypeDirectory,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create express bucket %s: %v", bucketName, err)
	}

	w := s3.NewBucketExistsWaiter(svc)
	err = w.Wait(ctx, &s3.HeadBucketInput{
		Bucket: &bucketName,
	}, 10*time.Second)
	if err != nil {
		return fmt.Errorf("wait for express bucket %s: %v", bucketName, err)
	}

	return nil
}
