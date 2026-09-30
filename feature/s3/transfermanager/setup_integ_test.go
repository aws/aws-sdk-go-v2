//go:build integration

package transfermanager

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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
var testHTTPClient *http.Client

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

	testHTTPClient = &http.Client{
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
	cfg.HTTPClient = testHTTPClient

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
	s3TransferManagerClient = New(s3Client)

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
	Body              io.Reader
	ChecksumAlgorithm types.ChecksumAlgorithm
	ChecksumType      types.ChecksumType
	ExpectBody        []byte
	ExpectError       string
}

type downloadObjectTestData struct {
	Body        io.Reader
	Range       string
	ExpectBody  []byte
	ExpectError string
	OptFns      []func(*Options)
	// PartSizes, when set, uploads Body as a multipart object whose parts
	// have exactly these byte sizes (all but the last must be >= 5MB per the
	// S3 minimum). Used to exercise downloads of objects with unequal part
	// sizes (#3526).
	PartSizes []int64
}

type getObjectTestData struct {
	Body            io.Reader
	Range           string
	ExpectBody      []byte
	ExpectGetError  string
	ExpectReadError string
	OptFns          []func(*Options)
	// PartSizes, when set, uploads Body as a multipart object whose parts have
	// exactly these byte sizes (all but the last must be >= 5MB per the S3
	// minimum). Used to exercise GetObject of objects with unequal part sizes
	// (#3526).
	PartSizes []int64
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
			Bucket:            aws.String(bucket),
			Key:               aws.String(key),
			Body:              testData.Body,
			ChecksumAlgorithm: testData.ChecksumAlgorithm,
			ChecksumType:      testData.ChecksumType,
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

	b, _ := io.ReadAll(resp.Body)
	if e, a := testData.ExpectBody, b; !bytes.EqualFold(e, a) {
		t.Errorf("expect %s, got %s", e, a)
	}

	if e, a := string(testData.ChecksumType), string(resp.ChecksumType); e != "" && e != a {
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
			Range:  aws.String(testData.Range),
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

// testGetObjectWithChangingPartSize uploads testData.Body as a multipart object
// whose parts have unequal sizes (testData.PartSizes), then reads it back through
// the transfer manager's GetObject and asserts the streamed bytes exactly equal
// what was uploaded. This exercises the parts-mode concurrent reader, which must
// not assume all parts share the first part's size (#3526). S3 requires every
// part except the last to be at least 5MB, so PartSizes must respect that.
func testGetObjectWithChangingPartSize(t *testing.T, bucket string, testData getObjectTestData) {
	key := UniqueID()

	body, err := io.ReadAll(testData.Body)
	if err != nil {
		t.Fatalf("expect no error reading test body, got %v", err)
	}

	var total int64
	for _, s := range testData.PartSizes {
		total += s
	}
	if total != int64(len(body)) {
		t.Fatalf("PartSizes sum to %d but body is %d bytes; they must match", total, len(body))
	}

	createOut, err := s3Client.CreateMultipartUpload(context.Background(),
		&s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		t.Fatalf("expect no error creating multipart upload, got %v", err)
	}
	uploadID := createOut.UploadId

	abort := func() {
		_, _ = s3Client.AbortMultipartUpload(context.Background(),
			&s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucket),
				Key:      aws.String(key),
				UploadId: uploadID,
			})
	}

	var completedParts []s3types.CompletedPart
	var offset int64
	for i, size := range testData.PartSizes {
		partNum := int32(i + 1)
		partOut, err := s3Client.UploadPart(context.Background(),
			&s3.UploadPartInput{
				Bucket:     aws.String(bucket),
				Key:        aws.String(key),
				UploadId:   uploadID,
				PartNumber: aws.Int32(partNum),
				Body:       bytes.NewReader(body[offset : offset+size]),
			})
		if err != nil {
			abort()
			t.Fatalf("expect no error uploading part %d, got %v", partNum, err)
		}
		completedParts = append(completedParts, s3types.CompletedPart{
			ETag:       partOut.ETag,
			PartNumber: aws.Int32(partNum),
		})
		offset += size
	}

	if _, err = s3Client.CompleteMultipartUpload(context.Background(),
		&s3.CompleteMultipartUploadInput{
			Bucket:          aws.String(bucket),
			Key:             aws.String(key),
			UploadId:        uploadID,
			MultipartUpload: &s3types.CompletedMultipartUpload{Parts: completedParts},
		}); err != nil {
		abort()
		t.Fatalf("expect no error completing multipart upload, got %v", err)
	}

	out, err := s3TransferManagerClient.GetObject(context.Background(),
		&GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Range:  aws.String(testData.Range),
		}, testData.OptFns...)
	if err != nil {
		if len(testData.ExpectGetError) == 0 {
			t.Fatalf("expect no error when getting object, got %v", err)
		}
		if e, a := testData.ExpectGetError, err.Error(); !strings.Contains(a, e) {
			t.Fatalf("expect error to contain %v, got %v", e, a)
		}
		return
	} else if e := testData.ExpectGetError; len(e) != 0 {
		t.Fatalf("expect error when getting object: %v, got none", e)
	}

	got, err := io.ReadAll(out.Body)
	if err != nil {
		if len(testData.ExpectReadError) == 0 {
			t.Fatalf("expect no error when reading responses, got %v", err)
		}
		if e, a := testData.ExpectReadError, err.Error(); !strings.Contains(a, e) {
			t.Fatalf("expect error to contain %v, got %v", e, a)
		}
		return
	} else if e := testData.ExpectReadError; len(e) != 0 {
		t.Fatalf("expect error when reading responses: %v, got none", e)
	}

	if e, a := testData.ExpectBody, got; !bytes.Equal(e, a) {
		t.Errorf("expect streamed object to equal uploaded object: uploaded %d bytes, got %d bytes", len(e), len(a))
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
			Range:    aws.String(testData.Range),
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

// testDownloadObjectWithChangingPartSize uploads testData.Body as a multipart
// object whose parts have unequal sizes (testData.PartSizes), then downloads it
// through the transfer manager and asserts the downloaded bytes exactly equal
// what was uploaded. This exercises the download path that must not assume all
// parts share the first part's size (#3526). S3 requires every part except the
// last to be at least 5MB, so PartSizes must respect that.
func testDownloadObjectWithChangingPartSize(t *testing.T, bucket string, testData downloadObjectTestData) {
	key := UniqueID()

	body, err := io.ReadAll(testData.Body)
	if err != nil {
		t.Fatalf("expect no error reading test body, got %v", err)
	}

	var total int64
	for _, s := range testData.PartSizes {
		total += s
	}
	if total != int64(len(body)) {
		t.Fatalf("PartSizes sum to %d but body is %d bytes; they must match", total, len(body))
	}

	createOut, err := s3Client.CreateMultipartUpload(context.Background(),
		&s3.CreateMultipartUploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		t.Fatalf("expect no error creating multipart upload, got %v", err)
	}
	uploadID := createOut.UploadId

	abort := func() {
		_, _ = s3Client.AbortMultipartUpload(context.Background(),
			&s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucket),
				Key:      aws.String(key),
				UploadId: uploadID,
			})
	}

	var completedParts []s3types.CompletedPart
	var offset int64
	for i, size := range testData.PartSizes {
		partNum := int32(i + 1)
		partOut, err := s3Client.UploadPart(context.Background(),
			&s3.UploadPartInput{
				Bucket:     aws.String(bucket),
				Key:        aws.String(key),
				UploadId:   uploadID,
				PartNumber: aws.Int32(partNum),
				Body:       bytes.NewReader(body[offset : offset+size]),
			})
		if err != nil {
			abort()
			t.Fatalf("expect no error uploading part %d, got %v", partNum, err)
		}
		completedParts = append(completedParts, s3types.CompletedPart{
			ETag:       partOut.ETag,
			PartNumber: aws.Int32(partNum),
		})
		offset += size
	}

	if _, err = s3Client.CompleteMultipartUpload(context.Background(),
		&s3.CompleteMultipartUploadInput{
			Bucket:          aws.String(bucket),
			Key:             aws.String(key),
			UploadId:        uploadID,
			MultipartUpload: &s3types.CompletedMultipartUpload{Parts: completedParts},
		}); err != nil {
		abort()
		t.Fatalf("expect no error completing multipart upload, got %v", err)
	}

	w := types.NewWriteAtBuffer(make([]byte, 0))
	_, err = s3TransferManagerClient.DownloadObject(context.Background(),
		&DownloadObjectInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			WriterAt: w,
			Range:    aws.String(testData.Range),
		}, testData.OptFns...)
	if err != nil {
		if len(testData.ExpectError) == 0 {
			t.Fatalf("expect no error, got %v", err)
		}
		if e, a := testData.ExpectError, err.Error(); !strings.Contains(a, e) {
			t.Fatalf("expect error to contain %v, got %v", e, a)
		}
	} else if e := testData.ExpectError; len(e) != 0 {
		t.Fatalf("expect error: %v, got none", e)
	}
	if len(testData.ExpectError) != 0 {
		return
	}

	if e, a := testData.ExpectBody, w.Bytes(); !bytes.Equal(e, a) {
		t.Errorf("expect downloaded object to equal uploaded object: uploaded %d bytes, got %d bytes", len(e), len(a))
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

		b, _ := io.ReadAll(resp.Body)
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

// ---------------------------------------------------------------------------
// Fuzzy/Property integration test suites for transfer manager v2 download. In
// each test of a fuzzy test iteration, an object with random size is uploaded
// in part(s) with random parts sizes via core S3. Then it is downloaded by
// DownloadObject/GetObject and check if the data downloaded matches the initial
// random source. Other main config options like concurrency, download type, max
// retries etc. are also randomly chosen for each case.
// ---------------------------------------------------------------------------

const randMiB = 1024 * 1024

// s3MinPartSize is S3's minimum size for every multipart part except the last.
const s3MinPartSize = 5 * randMiB

// maxRandomParts bounds the number of parts per generated object. Kept small so
// the suite's real uploads/downloads stay within a reasonable time/data budget
// (each non-last part is >= 5MiB).
const maxRandomParts = 5

// defaultReassemblyIterations is the number of randomized cases per run when
// TM_REASSEMBLY_ITERATIONS is unset.
const defaultReassemblyIterations = 10

// reassemblyCase is one generated reassembly scenario, fully determined by the
// seed. It fixes the object layout and the download options but not the API
// (GetObject vs DownloadObject), which is chosen by the calling test.
type reassemblyCase struct {
	partSizes       []int64 // len 1 => single-part (non-multipart) object
	useRanges       bool    // GetObjectRanges vs GetObjectParts
	concurrency     int
	bufferSize      int64
	partSize        int64 // 0 => leave default
	maxRetries      int
	disableChecksum bool
	rangeHeader     string // "" => whole object
	rangeStart      int64
	rangeEnd        int64
}

func (c reassemblyCase) total() int64 {
	var t int64
	for _, s := range c.partSizes {
		t += s
	}
	return t
}

func (c reassemblyCase) String() string {
	mode := "parts"
	if c.useRanges {
		mode = "ranges"
	}
	return fmt.Sprintf("mode=%s parts=%v total=%d concurrency=%d bufferSize=%d partSize=%d retries=%d disableChecksum=%t range=%q",
		mode, c.partSizes, c.total(), c.concurrency, c.bufferSize, c.partSize, c.maxRetries, c.disableChecksum, c.rangeHeader)
}

// optFns renders the case's download options.
func (c reassemblyCase) optFns() []func(*Options) {
	return []func(*Options){func(o *Options) {
		if c.useRanges {
			o.GetObjectType = types.GetObjectRanges
		} else {
			o.GetObjectType = types.GetObjectParts
		}
		o.Concurrency = c.concurrency
		o.GetObjectBufferSize = c.bufferSize
		if c.partSize != 0 {
			o.PartSizeBytes = c.partSize
		}
		o.PartBodyMaxRetries = c.maxRetries
		o.DisableChecksumValidation = c.disableChecksum
	}}
}

// rangePtr returns the case's Range header as a pointer, or nil for whole
// object.
func (c reassemblyCase) rangePtr() *string {
	if c.rangeHeader == "" {
		return nil
	}
	return aws.String(c.rangeHeader)
}

// generateReassemblyCase draws one fully-randomized scenario from rng.
func generateReassemblyCase(rng *mrand.Rand) reassemblyCase {
	c := reassemblyCase{}

	n := 1 + rng.Intn(maxRandomParts)
	if n == 1 {
		c.partSizes = []int64{int64(1 + rng.Intn(9*randMiB))}
	} else {
		sizes := make([]int64, n)
		for i := 0; i < n-1; i++ {
			sizes[i] = int64(s3MinPartSize + rng.Intn(3*randMiB+1))
		}
		// Last part: anything from 1 byte up to ~6MiB, including sub-part sizes.
		sizes[n-1] = int64(1 + rng.Intn(6*randMiB))
		c.partSizes = sizes
	}

	c.useRanges = rng.Intn(2) == 0
	c.concurrency = 1 + rng.Intn(8)
	// Buffer size: bias toward small buffers, which stress the streaming
	// in-order frontier in GetObject; occasionally use the large default.
	switch rng.Intn(3) {
	case 0:
		c.bufferSize = int64(randMiB + rng.Intn(randMiB)) // ~1MiB, forces many refills
	case 1:
		c.bufferSize = int64(6 * randMiB)
	default:
		c.bufferSize = defaultGetBufferSize
	}
	// Download part/range chunk size: sometimes override, always >= 5MiB.
	if rng.Intn(2) == 0 {
		c.partSize = int64(s3MinPartSize + rng.Intn(4*randMiB))
	}
	c.maxRetries = 1 + rng.Intn(4)
	c.disableChecksum = rng.Intn(2) == 0

	// Occasionally issue a sub-range request. Range input is only exercised in
	// ranges mode, matching the supported combinations in the hand-written
	// integration cases.
	total := c.total()
	if c.useRanges && total > 1 && rng.Intn(3) == 0 {
		start := int64(rng.Intn(int(total)))
		end := start + int64(rng.Intn(int(total-start)))
		c.rangeStart = start
		c.rangeEnd = end
		c.rangeHeader = fmt.Sprintf("bytes=%d-%d", start, end)
	}
	return c
}

// reassemblyDownloadFn performs the transfer-manager download for one case and
// returns the reassembled bytes. GetObject and DownloadObject each supply one.
type reassemblyDownloadFn func(t *testing.T, ctx context.Context, bucket, key string, c reassemblyCase) []byte

// runRandomizedReassembly is the shared property loop for downloader fuzzy test. It generates
// randomized layouts, uploads them (real multipart for >1 part), downloads them
// via downloadFn under randomized options, and asserts byte-exact reassembly.
// label distinguishes the calling test in logs.
func runRandomizedReassembly(t *testing.T, label string, downloadFn reassemblyDownloadFn) {
	t.Helper()

	seed := time.Now().UnixNano()
	if s := os.Getenv("TM_REASSEMBLY_SEED"); s != "" {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			t.Fatalf("invalid TM_REASSEMBLY_SEED %q: %v", s, err)
		}
		seed = v
	}
	iterations := defaultReassemblyIterations
	if s := os.Getenv("TM_REASSEMBLY_ITERATIONS"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v <= 0 {
			t.Fatalf("invalid TM_REASSEMBLY_ITERATIONS %q", s)
		}
		iterations = v
	}
	t.Logf("%s randomized reassembly: seed=%d iterations=%d (set TM_REASSEMBLY_SEED=%d to reproduce)", label, seed, iterations, seed)

	rng := mrand.New(mrand.NewSource(seed))
	bucket := setupMetadata.Buckets.Source.Name
	ctx := context.Background()

	for i := 0; i < iterations; i++ {
		c := generateReassemblyCase(rng)
		// Derive the body deterministically from the same rng so a seed fully
		// reproduces the case, content included.
		body := make([]byte, c.total())
		rng.Read(body)

		t.Run(fmt.Sprintf("iter-%02d", i), func(t *testing.T) {
			t.Logf("case: %s", c)
			key := UniqueID()

			if len(c.partSizes) == 1 {
				if _, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(key),
					Body:   bytes.NewReader(body),
				}); err != nil {
					t.Fatalf("put single-part object: %v", err)
				}
			} else {
				uploadMultipartExact(t, ctx, bucket, key, body, c.partSizes)
			}

			want := body
			if c.rangeHeader != "" {
				want = body[c.rangeStart : c.rangeEnd+1]
			}

			got := downloadFn(t, ctx, bucket, key, c)
			if !bytes.Equal(want, got) {
				t.Errorf("reassembly mismatch: want %d bytes, got %d bytes; first differing byte at %d\ncase: %s",
					len(want), len(got), firstDiff(want, got), c)
			}
		})
	}
}

// uploadMultipartExact uploads body as a multipart object whose parts have
// exactly the given byte sizes (all but the last must be >= 5MiB).
func uploadMultipartExact(t *testing.T, ctx context.Context, bucket, key string, body []byte, partSizes []int64) {
	t.Helper()

	createOut, err := s3Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("create multipart upload: %v", err)
	}
	uploadID := createOut.UploadId
	abort := func() {
		_, _ = s3Client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
			Bucket: aws.String(bucket), Key: aws.String(key), UploadId: uploadID,
		})
	}

	var completed []s3types.CompletedPart
	var offset int64
	for i, size := range partSizes {
		partNum := int32(i + 1)
		partOut, err := s3Client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(key),
			UploadId:   uploadID,
			PartNumber: aws.Int32(partNum),
			Body:       bytes.NewReader(body[offset : offset+size]),
		})
		if err != nil {
			abort()
			t.Fatalf("upload part %d: %v", partNum, err)
		}
		completed = append(completed, s3types.CompletedPart{
			ETag:       partOut.ETag,
			PartNumber: aws.Int32(partNum),
		})
		offset += size
	}

	if _, err := s3Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(key),
		UploadId:        uploadID,
		MultipartUpload: &s3types.CompletedMultipartUpload{Parts: completed},
	}); err != nil {
		abort()
		t.Fatalf("complete multipart upload: %v", err)
	}
}

// firstDiff returns the index of the first differing byte between a and b, or
// -1 if one is a prefix of the other (a length-only difference).
func firstDiff(a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return -1
}
