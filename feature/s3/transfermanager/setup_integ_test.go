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
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
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

	largeObjectBuf = make([]byte, 100*1024*1024)
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

type downloadObjectTestData struct {
	Body        io.Reader
	ExpectBody  []byte
	ExpectError string
	OptFns      []func(*Options)
}

type getObjectTestData struct {
	Body            io.Reader
	ExpectBody      []byte
	ExpectGetError  string
	ExpectReadError string
	OptFns          []func(*Options)
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

	_, err := s3TransferManagerClient.UploadObject(context.Background(),
		&UploadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
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

func testGetObject(t *testing.T, bucket string, testData getObjectTestData) {
	key := UniqueID()

	_, err := s3Client.PutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   testData.Body,
		})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	out, err := s3TransferManagerClient.GetObject(context.Background(),
		&GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}, testData.OptFns...)

	if err != nil {
		if len(testData.ExpectGetError) == 0 {
			t.Fatalf("expect no error when getting object, got %v", err)
		}
		if e, a := testData.ExpectGetError, err.Error(); !strings.Contains(a, e) {
			t.Fatalf("expect error to contain %v, got %v", e, a)
		}
	} else {
		if e := testData.ExpectGetError; len(e) != 0 {
			t.Fatalf("expect error when getting object: %v, got none", e)
		}
	}
	if len(testData.ExpectGetError) != 0 {
		return
	}

	b, err := io.ReadAll(out.Body)
	if err != nil {
		if len(testData.ExpectReadError) == 0 {
			t.Fatalf("expect no error when reading responses, got %v", err)
		}
		if e, a := testData.ExpectReadError, err.Error(); !strings.Contains(a, e) {
			t.Fatalf("expect error to contain %v, got %v", e, a)
		}
	} else {
		if e := testData.ExpectReadError; len(e) != 0 {
			t.Fatalf("expect error when reading responses: %v, got none", e)
		}
	}
	if len(testData.ExpectReadError) != 0 {
		return
	}
	if e, a := testData.ExpectBody, b; !bytes.EqualFold(e, a) {
		t.Errorf("expect %s, got %s", e, a)
	}
}

func testDownloadObject(t *testing.T, bucket string, testData downloadObjectTestData) {
	key := UniqueID()

	_, err := s3Client.PutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   testData.Body,
		})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	w := types.NewWriteAtBuffer(make([]byte, 0))
	_, err = s3TransferManagerClient.DownloadObject(context.Background(),
		&DownloadObjectInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			WriterAt: w,
		}, testData.OptFns...)
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

	if e, a := testData.ExpectBody, w.Bytes(); !bytes.EqualFold(e, a) {
		t.Errorf("expect %s, got %s", e, a)
	}
}

type uploadDirectoryTestData struct {
	FilesSize           map[string]int64
	Source              string
	Recursive           bool
	KeyPrefix           string
	ExpectFilesUploaded int64
	ExpectKeys          []string
	ExpectError         string
}

func testUploadDirectory(t *testing.T, bucket string, testData uploadDirectoryTestData) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "testdata")
	delimiter := "/"
	expectObjects := map[string][]byte{}
	source := filepath.Join(root, testData.Source)
	if err := os.MkdirAll(source, os.ModePerm); err != nil {
		t.Fatalf("error when creating test folder %v", err)
	}
	defer os.RemoveAll(source)
	for f, size := range testData.FilesSize {
		path := filepath.Join(source, strings.Replace(f, "/", string(os.PathSeparator), -1))
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			t.Fatalf("error when creating directory for file %s", path)
		}
		objectBuf := make([]byte, size)
		_, err := rand.Read(objectBuf)
		if err != nil {
			t.Fatalf("error when mocking test data for file %s", path)
		}
		file, err := os.Create(path)
		if err != nil {
			t.Fatalf("error when opening test file %s: %v", path, err)
		}
		_, err = file.Write(objectBuf)
		if err != nil {
			t.Fatalf("error when writing test file %s: %v", path, err)
		}
		key := strings.Replace(f, "/", delimiter, -1)
		if testData.KeyPrefix != "" {
			key = testData.KeyPrefix + delimiter + key
		}
		expectObjects[key] = objectBuf
	}

	out, err := s3TransferManagerClient.UploadDirectory(context.Background(), &UploadDirectoryInput{
		Bucket:    aws.String(bucket),
		Source:    aws.String(source),
		Recursive: aws.Bool(testData.Recursive),
		KeyPrefix: aws.String(testData.KeyPrefix),
	})
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

	if e, a := testData.ExpectFilesUploaded, out.ObjectsUploaded; e != a {
		t.Errorf("expect %d files uploaded, got %d", e, a)
	}
	for _, key := range testData.ExpectKeys {
		resp, err := s3Client.GetObject(context.Background(),
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			})
		if err != nil {
			t.Fatalf("error when getting object %s", key)
		}

		b, _ := ioutil.ReadAll(resp.Body)
		expectData, ok := expectObjects[key]
		if !ok {
			t.Errorf("no data recorded for object %s", key)
		}
		if e, a := expectData, b; !bytes.EqualFold(e, a) {
			t.Errorf("for object %s, expect %s, got %s", key, e, a)
		}
	}
}

type downloadDirectoryTestData struct {
	ObjectsSize             map[string]int64
	KeyPrefix               string
	ExpectObjectsDownloaded int64
	ExpectFiles             []string
	ExpectError             string
}

func testDownloadDirectory(t *testing.T, bucket string, testData downloadDirectoryTestData) {
	_, filename, _, _ := runtime.Caller(0)
	dst := filepath.Join(filepath.Dir(filename), "testdata", "integ")
	defer os.RemoveAll(dst)

	delimiter := "/"
	keyprefix := testData.KeyPrefix
	if keyprefix != "" && !strings.HasSuffix(keyprefix, delimiter) {
		keyprefix = keyprefix + delimiter
	}
	expectFiles := map[string][]byte{}
	for key, size := range testData.ObjectsSize {
		fileBuf := make([]byte, size)
		_, err := rand.Read(fileBuf)
		if err != nil {
			t.Fatalf("error when mocking test data for object %s", key)
		}
		_, err = s3Client.PutObject(context.Background(),
			&s3.PutObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Body:   bytes.NewReader(fileBuf),
			})
		if err != nil {
			t.Fatalf("error when putting object %s", key)
		}
		file := strings.ReplaceAll(strings.TrimPrefix(key, keyprefix), delimiter, string(os.PathSeparator))
		expectFiles[file] = fileBuf
	}

	out, err := s3TransferManagerClient.DownloadDirectory(context.Background(), &DownloadDirectoryInput{
		Bucket:      aws.String(bucket),
		Destination: aws.String(dst),
		KeyPrefix:   aws.String(testData.KeyPrefix),
	})
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

	if e, a := testData.ExpectObjectsDownloaded, out.ObjectsDownloaded; e != a {
		t.Errorf("expect %d objects downloaded, got %d", e, a)
	}
	for _, file := range testData.ExpectFiles {
		f := strings.ReplaceAll(file, delimiter, string(os.PathSeparator))
		path := filepath.Join(dst, f)
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("error when reading downloaded file %s: %v", path, err)
		}
		expectData, ok := expectFiles[f]
		if !ok {
			t.Errorf("no data recorded for file %s", path)
			continue
		}
		if e, a := expectData, b; !bytes.EqualFold(e, a) {
			t.Errorf("for file %s, expect %s, got %s", f, e, a)
		}
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
		CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
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
		CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
			Location: &s3types.LocationInfo{
				Name: aws.String(expressAZID),
				Type: s3types.LocationTypeAvailabilityZone,
			},
			Bucket: &s3types.BucketInfo{
				DataRedundancy: s3types.DataRedundancySingleAvailabilityZone,
				Type:           s3types.BucketTypeDirectory,
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
