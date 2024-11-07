//go:build integration
// +build integration

package s3transfermanager

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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	tm "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest/s3shared"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
var s3TransferManagerClient *tm.Client

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

	cfg, err := integrationtest.LoadConfigWithDefaultRegion(region)
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
	s3TransferManagerClient = tm.NewFromConfig(s3Client, s3cfg)

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
		*bucket.name = s3shared.GenerateBucketName()

		if err := s3shared.SetupBucket(ctx, s3Client, *bucket.name); err != nil {
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
			if err := s3shared.CleanupBucket(ctx, s3Client, bucketName); err != nil {
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

func testPutObject(t *testing.T, bucket string, testData putObjectTestData, opts ...func(options *tm.Options)) {
	key := integrationtest.UniqueID()

	_, err := s3TransferManagerClient.PutObject(context.Background(),
		&tm.PutObjectInput{
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
