// +build integration

package s3

import (
	"bytes"
	"context"
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
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest/s3shared"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
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
		Target struct {
			Name string
			ARN  string
		}
	}

	AccessPoints struct {
		Source struct {
			Name string
			ARN  string
		}
		Target struct {
			Name string
			ARN  string
		}
	}
}{}

// s3 client to use for integ testing
var s3client *s3.Client

// s3-control client to use for integ testing
var s3ControlClient *s3control.Client

// sts client to use for integ testing
var stsClient *sts.Client

// http client setting to use for integ testing
var httpClient *http.Client

var region = "us-west-2"

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
			fmt.Fprintln(os.Stderr, "S3 integration tests panic,", r)
			result = 1
		}
		os.Exit(result)
	}()

	var verifyTLS bool
	var s3Endpoint, s3ControlEndpoint string
	var s3EnableTLS, s3ControlEnableTLS bool

	flag.StringVar(&s3Endpoint, "s3-endpoint", "", "integration endpoint for S3")
	flag.BoolVar(&s3EnableTLS, "s3-tls", true, "enable TLS for S3 endpoint")

	flag.StringVar(&s3ControlEndpoint, "s3-control-endpoint", "", "integration endpoint for S3")
	flag.BoolVar(&s3ControlEnableTLS, "s3-control-tls", true, "enable TLS for S3 control endpoint")

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
	s3client = s3.NewFromConfig(s3cfg)

	// create a s3-control client
	s3ControlCfg := cfg.Copy()
	if len(s3ControlEndpoint) != 0 {
		s3ControlCfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           s3ControlEndpoint,
				PartitionID:   "aws",
				SigningName:   "s3-control",
				SigningRegion: region,
			}, nil
		})
	}

	// build s3-control client from config
	s3ControlClient = s3control.NewFromConfig(s3ControlCfg)

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

	accessPointsCleanup, err := setupAccessPoints(ctx)
	defer accessPointsCleanup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup integration test access points: %v\n", err)
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
		{name: &setupMetadata.Buckets.Target.Name, arn: &setupMetadata.Buckets.Target.ARN},
	}

	for _, bucket := range bucketCreates {
		*bucket.name = s3shared.GenerateBucketName()

		if err := s3shared.SetupBucket(ctx, s3client, *bucket.name); err != nil {
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
			if err := s3shared.CleanupBucket(ctx, s3client, bucketName); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		})
	}

	return cleanup, nil

}

func setupAccessPoints(ctx context.Context) (func(), error) {
	var cleanups []func()

	cleanup := func() {
		for i := range cleanups {
			cleanups[i]()
		}
	}

	creates := []struct {
		bucket string
		name   *string
		arn    *string
	}{
		{bucket: setupMetadata.Buckets.Source.Name, name: &setupMetadata.AccessPoints.Source.Name, arn: &setupMetadata.AccessPoints.Source.ARN},
		{bucket: setupMetadata.Buckets.Target.Name, name: &setupMetadata.AccessPoints.Target.Name, arn: &setupMetadata.AccessPoints.Target.ARN},
	}

	for _, ap := range creates {
		*ap.name = integrationtest.UniqueID()

		err := s3shared.SetupAccessPoint(ctx, s3ControlClient, setupMetadata.AccountID, ap.bucket, *ap.name)
		if err != nil {
			return cleanup, err
		}

		// Compute ARN
		apARN := arn.ARN{
			Partition: "aws",
			Service:   "s3",
			Region:    region,
			AccountID: setupMetadata.AccountID,
			Resource:  fmt.Sprintf("accesspoint/%s", *ap.name),
		}.String()

		*ap.arn = apARN

		apName := *ap.name
		cleanups = append(cleanups, func() {
			err := s3shared.CleanupAccessPoint(ctx, s3ControlClient, setupMetadata.AccountID, apName)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		})
	}

	return cleanup, nil
}

func putTestContent(t *testing.T, reader io.ReadSeeker, key string, opts func(options *s3.Options)) {
	t.Logf("uploading test file %s/%s", setupMetadata.Buckets.Source.Name, key)
	_, err := s3client.PutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket: &setupMetadata.Buckets.Source.Name,
			Key:    aws.String(key),
			Body:   reader,
		}, opts)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

type writeToObjectTestData struct {
	Body        io.Reader
	ExpectBody  []byte
	ExpectError string
}

func testWriteToObject(t *testing.T, bucket string, testData writeToObjectTestData, opts ...func(options *s3.Options)) {
	key := integrationtest.UniqueID()

	// put object
	_, err := s3client.PutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket: &bucket,
			Key:    &key,
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
		if len(testData.ExpectError) != 0 {
			t.Fatalf("expected error: %v, got none", err)
		}
	}

	// stop if expected error writing object
	if len(testData.ExpectError) != 0 {
		return
	}

	// get object
	resp, err := s3client.GetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		}, opts...)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	if e, a := testData.ExpectBody, b; !bytes.EqualFold(e, a) {
		t.Errorf("expect %s, got %s", e, a)
	}
}

func testCopyObject(t *testing.T, sourceBucket, targetBucket string, opts func(options *s3.Options)) {
	key := integrationtest.UniqueID()

	if opts == nil {
		opts = func(options *s3.Options) {
		}
	}

	// put object
	_, err := s3client.PutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket: &sourceBucket,
			Key:    &key,
			Body:   bytes.NewReader([]byte(`hello world`)),
		}, opts)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	// copy object
	_, err = s3client.CopyObject(context.Background(),
		&s3.CopyObjectInput{
			Bucket:     &targetBucket,
			Key:        &key,
			CopySource: aws.String("/" + sourceBucket + "/" + key),
		}, opts)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	// get object
	resp, err := s3client.GetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: &targetBucket,
			Key:    &key,
		}, opts)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	if e, a := []byte("hello world"), b; !bytes.EqualFold(e, a) {
		t.Errorf("expect %s, got %s", e, a)
	}
}
