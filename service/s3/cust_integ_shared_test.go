// +build integration

package s3_test

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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration/s3integ"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const integBucketPrefix = "aws-sdk-go-integration"

var integMetadata = struct {
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

var s3Svc *s3.Client
var s3ControlSvc *s3control.Client
var stsSvc *sts.Client
var httpClient *http.Client

// TODO: (Westeros) Remove Custom Resolver Usage Before Launch
type customResolver struct {
	endpoint string
	withTLS  bool
	region   string
}

func (r customResolver) ResolveEndpoint(service, _ string) (aws.Endpoint, error) {
	switch strings.ToLower(service) {
	case "s3-control":
	case "s3":
	default:
		return aws.Endpoint{}, fmt.Errorf("unsupported in custom resolver")
	}

	scheme := "https"
	if !r.withTLS {
		scheme = "http"
	}

	return aws.Endpoint{
		PartitionID:   "aws",
		SigningRegion: r.region,
		SigningName:   "s3",
		SigningMethod: "s3v4",
		URL:           fmt.Sprintf("%s://%s", scheme, r.endpoint),
	}, nil
}

type LoggerFunc func(...interface{})

func (l LoggerFunc) Log(v ...interface{}) {
	l(v...)
}

func TestMain(m *testing.M) {
	var result int
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "S3 integration tests paniced,", r)
			result = 1
		}
		os.Exit(result)
	}()

	var verifyTLS, logDebug bool
	var s3Endpoint, s3ControlEndpoint string
	var s3EnableTLS, s3ControlEnableTLS bool

	flag.StringVar(&s3Endpoint, "s3-endpoint", "", "integration endpoint for S3")
	flag.BoolVar(&s3EnableTLS, "s3-tls", true, "enable TLS for S3 endpoint")

	flag.StringVar(&s3ControlEndpoint, "s3-control-endpoint", "", "integration endpoint for S3")
	flag.BoolVar(&s3ControlEnableTLS, "s3-control-tls", true, "enable TLS for S3 control endpoint")

	flag.StringVar(&integMetadata.AccountID, "account", "", "integration account id")
	flag.BoolVar(&verifyTLS, "verify-tls", true, "verify server TLS certificate")
	flag.BoolVar(&logDebug, "log-debug", false, "verify server TLS certificate")
	flag.Parse()

	cfg, err := external.LoadDefaultAWSConfig(external.WithRegion("us-west-2"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load external config: %v\n", err)
		result = 1
		return
	}

	if logDebug {
		cfg.Logger = LoggerFunc(func(v ...interface{}) {
			fmt.Printf(v[0].(string), v[1:]...)
		})
		cfg.LogLevel = aws.LogDebugWithRequestErrors
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: verifyTLS},
		},
	}

	cfg.HTTPClient = httpClient

	s3Cfg := cfg.Copy()
	if len(s3Endpoint) != 0 {
		s3Cfg.EndpointResolver = customResolver{
			endpoint: s3Endpoint,
			withTLS:  s3EnableTLS,
			region:   cfg.Region,
		}
	}
	s3Svc = s3.New(s3Cfg)
	s3Svc.UseARNRegion = true

	s3ControlCfg := cfg.Copy()
	if len(s3Endpoint) != 0 {
		s3ControlCfg.EndpointResolver = customResolver{
			endpoint: s3ControlEndpoint,
			withTLS:  s3ControlEnableTLS,
			region:   cfg.Region,
		}
	}

	s3ControlSvc = s3control.New(s3ControlCfg)
	stsSvc = sts.New(cfg)

	integMetadata.AccountID, err = getAccountID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get integration aws account id: %v\n", err)
		result = 1
		return
	}

	bucketCleanup, err := setupBuckets()
	defer bucketCleanup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup integration test buckets: %v\n", err)
		result = 1
		return
	}

	accessPointsCleanup, err := setupAccessPoints()
	defer accessPointsCleanup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup integration test access points: %v\n", err)
		result = 1
		return
	}

	result = m.Run()
}

func getAccountID() (string, error) {
	if len(integMetadata.AccountID) != 0 {
		return integMetadata.AccountID, nil
	}

	req := stsSvc.GetCallerIdentityRequest(nil)
	output, err := req.Send(context.Background())
	if err != nil {
		return "", fmt.Errorf("faield to get sts caller identity: %v", err)
	}

	return *output.Account, nil
}

func setupBuckets() (func(), error) {
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
		{name: &integMetadata.Buckets.Source.Name, arn: &integMetadata.Buckets.Source.ARN},
		{name: &integMetadata.Buckets.Target.Name, arn: &integMetadata.Buckets.Target.ARN},
	}

	for _, bucket := range bucketCreates {
		*bucket.name = s3integ.GenerateBucketName()

		if err := s3integ.SetupBucket(context.Background(), s3Svc, *bucket.name); err != nil {
			return cleanup, err
		}

		// Compute ARN
		bARN := arn.ARN{
			Partition: "aws",
			Service:   "s3",
			Region:    s3Svc.Metadata.SigningRegion,
			AccountID: integMetadata.AccountID,
			Resource:  fmt.Sprintf("bucket_name:%s", *bucket.name),
		}.String()

		*bucket.arn = bARN

		bucketName := *bucket.name
		cleanups = append(cleanups, func() {
			if err := s3integ.CleanupBucket(context.Background(), s3Svc, bucketName); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		})
	}

	return cleanup, nil
}

func setupAccessPoints() (func(), error) {
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
		{bucket: integMetadata.Buckets.Source.Name, name: &integMetadata.AccessPoints.Source.Name, arn: &integMetadata.AccessPoints.Source.ARN},
		{bucket: integMetadata.Buckets.Target.Name, name: &integMetadata.AccessPoints.Target.Name, arn: &integMetadata.AccessPoints.Target.ARN},
	}

	for _, ap := range creates {
		*ap.name = integration.UniqueID()

		err := s3integ.SetupAccessPoint(s3ControlSvc, integMetadata.AccountID, ap.bucket, *ap.name)
		if err != nil {
			return cleanup, err
		}

		// Compute ARN
		apARN := arn.ARN{
			Partition: "aws",
			Service:   "s3",
			Region:    s3ControlSvc.Metadata.SigningRegion,
			AccountID: integMetadata.AccountID,
			Resource:  fmt.Sprintf("accesspoint/%s", *ap.name),
		}.String()

		*ap.arn = apARN

		apName := *ap.name
		cleanups = append(cleanups, func() {
			err := s3integ.CleanupAccessPoint(s3ControlSvc, integMetadata.AccountID, apName)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		})
	}

	return cleanup, nil
}

func putTestFile(t *testing.T, filename, key string) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open testfile, %v", err)
	}
	defer f.Close()

	putTestContent(t, f, key)
}

func putTestContent(t *testing.T, reader io.ReadSeeker, key string) {
	t.Logf("uploading test file %s/%s", integMetadata.Buckets.Source.Name, key)
	req := s3Svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: &integMetadata.Buckets.Source.Name,
		Key:    aws.String(key),
		Body:   reader,
	})
	_, err := req.Send(context.Background())
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func testWriteToObject(t *testing.T, bucket string) {
	key := integration.UniqueID()

	req := s3Svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader([]byte("hello world")),
	})
	_, err := req.Send(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	getReq := s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	resp, err := getReq.Send(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	if e, a := []byte("hello world"), b; !bytes.Equal(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func testPresignedGetPut(t *testing.T, bucket string) {
	key := integration.UniqueID()

	putReq := s3Svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	var err error

	// Presign a PUT aws
	var puturl string
	puturl, err = putReq.Presign(5 * time.Minute)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	// PUT to the presigned URL with a body
	var putHTTPReq *http.Request
	buf := bytes.NewReader([]byte("hello world"))
	putHTTPReq, err = http.NewRequest("PUT", puturl, buf)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	var putresp *http.Response
	putresp, err = httpClient.Do(putHTTPReq)
	if err != nil {
		t.Errorf("expect put with presign url no error, got %v", err)
	}
	if e, a := 200, putresp.StatusCode; e != a {
		t.Fatalf("expect %v, got %v", e, a)
	}

	// Presign a GET on the same URL
	getReq := s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	var geturl string
	geturl, err = getReq.Presign(300 * time.Second)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	// Get the body
	var getresp *http.Response
	getresp, err = httpClient.Get(geturl)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	var b []byte
	defer getresp.Body.Close()
	b, err = ioutil.ReadAll(getresp.Body)
	if e, a := "hello world", string(b); e != a {
		t.Fatalf("expect %v, got %v", e, a)
	}
}

func testCopyObject(t *testing.T, sourceBucket string, targetBucket string, opts ...aws.Option) {
	key := integration.UniqueID()

	putReq := s3Svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: &sourceBucket,
		Key:    &key,
		Body:   bytes.NewReader([]byte("hello world")),
	})
	_, err := putReq.Send(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	copyReq := s3Svc.CopyObjectRequest(&s3.CopyObjectInput{
		Bucket:     &targetBucket,
		CopySource: aws.String("/" + sourceBucket + "/" + key),
		Key:        &key,
	})
	_, err = copyReq.Send(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	getReq := s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &targetBucket,
		Key:    &key,
	})
	resp, err := getReq.Send(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	if e, a := []byte("hello world"), b; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}
