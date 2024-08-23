package transfermanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// getReaderLength discards the bytes from reader and returns the length
func getReaderLength(r io.Reader) int64 {
	n, _ := io.Copy(ioutil.Discard, r)
	return n
}

func TestUploadOrderMulti(t *testing.T) {
	c, invocations, args := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{})

	resp, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket:               "Bucket",
		Key:                  "Key - value",
		Body:                 bytes.NewReader(buf20MB),
		ServerSideEncryption: "aws:kms",
		SSEKMSKeyID:          "KmsId",
		ContentType:          "content/type",
	})

	if err != nil {
		t.Errorf("Expected no error but received %v", err)
	}

	if diff := cmpDiff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart", "UploadPart", "CompleteMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Errorf(diff)
	}

	if "UPLOAD-ID" != resp.UploadID {
		t.Errorf("expect %q, got %q", "UPLOAD-ID", resp.UploadID)
	}

	if "VERSION-ID" != resp.VersionID {
		t.Errorf("expect %q, got %q", "VERSION-ID", resp.VersionID)
	}

	// Validate input values

	// UploadPart
	for i := 1; i < 4; i++ {
		v := aws.ToString((*args)[i].(*s3.UploadPartInput).UploadId)
		if "UPLOAD-ID" != v {
			t.Errorf("Expected %q, but received %q", "UPLOAD-ID", v)
		}
	}

	// CompleteMultipartUpload
	v := aws.ToString((*args)[4].(*s3.CompleteMultipartUploadInput).UploadId)
	if "UPLOAD-ID" != v {
		t.Errorf("Expected %q, but received %q", "UPLOAD-ID", v)
	}

	parts := (*args)[4].(*s3.CompleteMultipartUploadInput).MultipartUpload.Parts

	for i := 0; i < 3; i++ {
		num := parts[i].PartNumber
		etag := aws.ToString(parts[i].ETag)

		if int32(i+1) != aws.ToInt32(num) {
			t.Errorf("expect %d, got %d", i+1, num)
		}

		if matched, err := regexp.MatchString(`^ETAG\d+$`, etag); !matched || err != nil {
			t.Errorf("Failed regexp expression `^ETAG\\d+$`")
		}
	}

	// Custom headers
	cmu := (*args)[0].(*s3.CreateMultipartUploadInput)

	if e, a := types.ServerSideEncryption("aws:kms"), cmu.ServerSideEncryption; e != a {
		t.Errorf("expect %q, got %q", e, a)
	}

	if e, a := "KmsId", aws.ToString(cmu.SSEKMSKeyId); e != a {
		t.Errorf("expect %q, got %q", e, a)
	}

	if e, a := "content/type", aws.ToString(cmu.ContentType); e != a {
		t.Errorf("expect %q, got %q", e, a)
	}
}

func TestUploadOrderMultiDifferentPartSize(t *testing.T) {
	c, ops, args := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{
		PartSizeBytes: 1024 * 1024 * 11,
		Concurrency:   1,
	})

	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   bytes.NewReader(buf20MB),
	})

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	vals := []string{"CreateMultipartUpload", "UploadPart", "UploadPart", "CompleteMultipartUpload"}
	if !reflect.DeepEqual(vals, *ops) {
		t.Errorf("expect %v, got %v", vals, *ops)
	}

	// Part lengths
	if len := getReaderLength((*args)[1].(*s3.UploadPartInput).Body); 1024*1024*11 != len {
		t.Errorf("expect %d, got %d", 1024*1024*7, len)
	}
	if len := getReaderLength((*args)[2].(*s3.UploadPartInput).Body); 1024*1024*9 != len {
		t.Errorf("expect %d, got %d", 1024*1024*5, len)
	}
}

func TestUploadFailIfPartSizeTooSmall(t *testing.T) {
	mgr := New(s3.New(s3.Options{}), Options{},
		func(o *Options) {
			o.PartSizeBytes = 5
		})
	resp, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   bytes.NewReader(buf20MB),
	})

	if resp != nil {
		t.Errorf("Expected response to be nil, but received %v", resp)
	}

	if err == nil {
		t.Errorf("Expected error, but received nil")
	}
	if e, a := "part size must be at least", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v to be in %v", e, a)
	}
}

func TestUploadOrderSingle(t *testing.T) {
	c, invocations, params := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{})
	resp, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket:               "Bucket",
		Key:                  "Key - value",
		Body:                 bytes.NewReader(buf2MB),
		ServerSideEncryption: "aws:kms",
		SSEKMSKeyID:          "KmsId",
		ContentType:          "content/type",
	})

	if err != nil {
		t.Errorf("expect no error but received %v", err)
	}

	if diff := cmpDiff([]string{"PutObject"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	if e := "VERSION-ID"; e != resp.VersionID {
		t.Errorf("expect %q, got %q", e, resp.VersionID)
	}

	if len(resp.UploadID) > 0 {
		t.Errorf("expect empty string, got %q", resp.UploadID)
	}

	putObjectInput := (*params)[0].(*s3.PutObjectInput)

	if e, a := types.ServerSideEncryption("aws:kms"), putObjectInput.ServerSideEncryption; e != a {
		t.Errorf("expect %q, got %q", e, a)
	}

	if e, a := "KmsId", aws.ToString(putObjectInput.SSEKMSKeyId); e != a {
		t.Errorf("expect %q, got %q", e, a)
	}

	if e, a := "content/type", aws.ToString(putObjectInput.ContentType); e != a {
		t.Errorf("Expected %q, but received %q", e, a)
	}
}

func TestUploadSingleFailure(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)

	c.PutObjectFn = func(*s3testing.UploadLoggingClient, *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return nil, fmt.Errorf("put object failure")
	}

	mgr := New(c, Options{})
	resp, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   bytes.NewReader(buf2MB),
	})

	if err == nil {
		t.Error("expect error, got nil")
	}

	if diff := cmpDiff([]string{"PutObject"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	if resp != nil {
		t.Errorf("expect response to be nil, got %v", resp)
	}
}

func TestUploadOrderZero(t *testing.T) {
	c, invocations, params := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{})
	resp, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   bytes.NewReader(make([]byte, 0)),
	})

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmpDiff([]string{"PutObject"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	if len(resp.UploadID) > 0 {
		t.Errorf("expect empty string, got %q", resp.UploadID)
	}

	if e, a := int64(0), getReaderLength((*params)[0].(*s3.PutObjectInput).Body); e != a {
		t.Errorf("Expected %d, but received %d", e, a)
	}
}

func TestUploadOrderMultiFailure(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)

	c.UploadPartFn = func(u *s3testing.UploadLoggingClient, params *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
		if *params.PartNumber == 2 {
			return nil, fmt.Errorf("an unexpected error")
		}
		return &s3.UploadPartOutput{ETag: aws.String(fmt.Sprintf("ETAG%d", u.PartNum))}, nil
	}

	mgr := New(c, Options{}, func(o *Options) {
		o.Concurrency = 1
	})
	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   bytes.NewReader(buf20MB),
	})

	if err == nil {
		t.Error("expect error, got nil")
	}

	if diff := cmpDiff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart", "AbortMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderMultiFailureOnComplete(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)

	c.CompleteMultipartUploadFn = func(*s3testing.UploadLoggingClient, *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
		return nil, fmt.Errorf("complete multipart error")
	}

	mgr := New(c, Options{}, func(o *Options) {
		o.Concurrency = 1
	})
	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   bytes.NewReader(buf20MB),
	})

	if err == nil {
		t.Error("expect error, got nil")
	}

	if diff := cmpDiff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart", "UploadPart",
		"CompleteMultipartUpload", "AbortMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderMultiFailureOnCreate(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)

	c.CreateMultipartUploadFn = func(*s3testing.UploadLoggingClient, *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
		return nil, fmt.Errorf("create multipart upload failure")
	}

	mgr := New(c, Options{})
	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   bytes.NewReader(make([]byte, 1024*1024*12)),
	})

	if err == nil {
		t.Error("expect error, got nil")
	}

	if diff := cmpDiff([]string{"CreateMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}
}

type failreader struct {
	times     int
	failCount int
}

func (f *failreader) Read(b []byte) (int, error) {
	f.failCount++
	if f.failCount >= f.times {
		return 0, fmt.Errorf("random failure")
	}
	return len(b), nil
}

func TestUploadOrderReadFail1(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{})
	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   &failreader{times: 1},
	})
	if err == nil {
		t.Fatalf("expect error to not be nil")
	}

	if e, a := "random failure", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v, got %v", e, a)
	}

	if diff := cmpDiff([]string(nil), *invocations); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderReadFail2(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient([]string{"UploadPart"})
	mgr := New(c, Options{}, func(o *Options) {
		o.Concurrency = 1
	})
	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   &failreader{times: 2},
	})
	if err == nil {
		t.Fatalf("expect error to not be nil")
	}

	if e, a := "random failure", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v, got %q", e, a)
	}

	if diff := cmpDiff([]string{"CreateMultipartUpload", "AbortMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}
}

type sizedReader struct {
	size int
	cur  int
	err  error
}

func (s *sizedReader) Read(p []byte) (n int, err error) {
	if s.cur >= s.size {
		if s.err == nil {
			s.err = io.EOF
		}
		return 0, s.err
	}

	n = len(p)
	s.cur += len(p)
	if s.cur > s.size {
		n -= s.cur - s.size
	}

	return n, err
}

func TestUploadOrderMultiBufferedReader(t *testing.T) {
	c, invocations, params := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{})
	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   &sizedReader{size: 1024 * 1024 * 21},
	})
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmpDiff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart",
		"UploadPart", "CompleteMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	// Part lengths
	var parts []int64
	for i := 1; i <= 3; i++ {
		parts = append(parts, getReaderLength((*params)[i].(*s3.UploadPartInput).Body))
	}
	sort.Slice(parts, func(i, j int) bool {
		return parts[i] < parts[j]
	})

	if diff := cmpDiff([]int64{1024 * 1024 * 5, 1024 * 1024 * 8, 1024 * 1024 * 8}, parts); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderMultiBufferedReaderPartial(t *testing.T) {
	c, invocations, params := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{})
	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   &sizedReader{size: 1024 * 1024 * 21, err: io.EOF},
	})
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmpDiff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart",
		"UploadPart", "CompleteMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	// Part lengths
	var parts []int64
	for i := 1; i <= 3; i++ {
		parts = append(parts, getReaderLength((*params)[i].(*s3.UploadPartInput).Body))
	}
	sort.Slice(parts, func(i, j int) bool {
		return parts[i] < parts[j]
	})

	if diff := cmpDiff([]int64{1024 * 1024 * 5, 1024 * 1024 * 8, 1024 * 1024 * 8}, parts); len(diff) > 0 {
		t.Error(diff)
	}
}

// TestUploadOrderMultiBufferedReaderEOF tests the edge case where the
// file size is the same as part size.
func TestUploadOrderMultiBufferedReaderEOF(t *testing.T) {
	c, invocations, params := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{})
	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   &sizedReader{size: 1024 * 1024 * 16, err: io.EOF},
	})

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmpDiff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart", "CompleteMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	// Part lengths
	var parts []int64
	for i := 1; i <= 2; i++ {
		parts = append(parts, getReaderLength((*params)[i].(*s3.UploadPartInput).Body))
	}
	sort.Slice(parts, func(i, j int) bool {
		return parts[i] < parts[j]
	})

	if diff := cmpDiff([]int64{1024 * 1024 * 8, 1024 * 1024 * 8}, parts); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderSingleBufferedReader(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{})
	resp, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   &sizedReader{size: 1024 * 1024 * 2},
	})

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmpDiff([]string{"PutObject"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	if len(resp.UploadID) > 0 {
		t.Errorf("expect no value, got %q", resp.UploadID)
	}
}

func TestUploadZeroLenObject(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)

	mgr := New(c, Options{})
	resp, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   strings.NewReader(""),
	})

	if err != nil {
		t.Errorf("expect no error but received %v", err)
	}
	if diff := cmpDiff([]string{"PutObject"}, *invocations); len(diff) > 0 {
		t.Errorf("expect request to have been made, but was not, %v", diff)
	}

	if len(resp.UploadID) > 0 {
		t.Errorf("expect empty string, but received %q", resp.UploadID)
	}
}

type testIncompleteReader struct {
	Size int64
	read int64
}

func (r *testIncompleteReader) Read(p []byte) (n int, err error) {
	r.read += int64(len(p))
	if r.read >= r.Size {
		return int(r.read - r.Size), io.ErrUnexpectedEOF
	}
	return len(p), nil
}

func TestUploadUnexpectedEOF(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)
	mgr := New(c, Options{}, func(o *Options) {
		o.Concurrency = 1
		o.PartSizeBytes = minPartSizeBytes
	})
	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body: &testIncompleteReader{
			Size: minPartSizeBytes + 1,
		},
	})
	if err == nil {
		t.Error("expect error, got nil")
	}

	// Ensure upload started.
	if e, a := "CreateMultipartUpload", (*invocations)[0]; e != a {
		t.Errorf("expect %q, got %q", e, a)
	}

	// Part may or may not be sent because of timing of sending parts and
	// reading next part in upload manager. Just check for the last abort.
	if e, a := "AbortMultipartUpload", (*invocations)[len(*invocations)-1]; e != a {
		t.Errorf("expect %q, got %q", e, a)
	}
}

func TestSSE(t *testing.T) {
	c, _, _ := s3testing.NewUploadLoggingClient(nil)
	c.UploadPartFn = func(u *s3testing.UploadLoggingClient, params *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
		if params.SSECustomerAlgorithm == nil {
			t.Fatal("SSECustomerAlgoritm should not be nil")
		}
		if params.SSECustomerKey == nil {
			t.Fatal("SSECustomerKey should not be nil")
		}
		return &s3.UploadPartOutput{ETag: aws.String(fmt.Sprintf("ETAG%d", u.PartNum))}, nil
	}

	mgr := New(c, Options{}, func(o *Options) {
		o.Concurrency = 5
	})

	_, err := mgr.PutObject(context.Background(), &PutObjectInput{
		Bucket:               "Bucket",
		Key:                  "Key",
		SSECustomerAlgorithm: "AES256",
		SSECustomerKey:       "foo",
		Body:                 bytes.NewBuffer(make([]byte, 1024*1024*10)),
	})

	if err != nil {
		t.Fatal("Expected no error, but received" + err.Error())
	}
}

func TestUploadWithContextCanceled(t *testing.T) {
	c := s3.New(s3.Options{
		UsePathStyle: true,
		Region:       "mock-region",
	})
	u := New(c, Options{})

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	ctx.Error = fmt.Errorf("context canceled")
	close(ctx.DoneCh)

	_, err := u.PutObject(ctx, &PutObjectInput{
		Bucket: "Bucket",
		Key:    "Key",
		Body:   bytes.NewReader(make([]byte, 0)),
	})
	if err == nil {
		t.Fatalf("expect error, got nil")
	}

	if e, a := "canceled", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expected error message to contain %q, but did not %q", e, a)
	}
}

func TestUploadRetry(t *testing.T) {
	const part, retries = 3, 10
	testFile, testFileCleanup, err := createTempFile(t, minPartSizeBytes*part)
	if err != nil {
		t.Fatalf("failed to create test file, %v", err)
	}
	defer testFileCleanup(t)

	cases := map[string]struct {
		Body         io.Reader
		PartHandlers func(testing.TB) []http.Handler
	}{
		"bytes.Buffer": {
			Body: bytes.NewBuffer(make([]byte, minPartSizeBytes*part)),
			PartHandlers: func(tb testing.TB) []http.Handler {
				return buildFailHandlers(tb, part, retries)
			},
		},
		"bytes.Reader": {
			Body: bytes.NewReader(make([]byte, minPartSizeBytes*part)),
			PartHandlers: func(tb testing.TB) []http.Handler {
				return buildFailHandlers(tb, part, retries)
			},
		},
		"os.File": {
			Body: testFile,
			PartHandlers: func(tb testing.TB) []http.Handler {
				return buildFailHandlers(tb, part, retries)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			restoreSleep := sdk.TestingUseNopSleep()
			defer restoreSleep()

			mux := newMockS3UploadServer(t, c.PartHandlers(t))
			server := httptest.NewServer(mux)
			defer server.Close()

			client := s3.New(s3.Options{
				EndpointResolverV2: s3testing.EndpointResolverV2{URL: server.URL},
				UsePathStyle:       true,
				Retryer: retry.NewStandard(func(o *retry.StandardOptions) {
					o.MaxAttempts = retries + 1
				}),
			})

			uploader := New(client, Options{})
			_, err := uploader.PutObject(context.Background(), &PutObjectInput{
				Bucket: "bucket",
				Key:    "key",
				Body:   c.Body,
			})

			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

func newMockS3UploadServer(tb testing.TB, partHandler []http.Handler) *mockS3UploadServer {
	s := &mockS3UploadServer{
		ServeMux:     http.NewServeMux(),
		partHandlers: partHandler,
		tb:           tb,
	}

	s.HandleFunc("/", s.handleRequest)

	return s
}

func buildFailHandlers(tb testing.TB, part, retry int) []http.Handler {
	handlers := make([]http.Handler, part)

	for i := 0; i < part; i++ {
		handlers[i] = &failPartHandler{
			tb:                 tb,
			failLeft:           retry,
			successPartHandler: &successPartHandler{tb: tb},
		}
	}

	return handlers
}

type mockS3UploadServer struct {
	*http.ServeMux

	tb           testing.TB
	partHandlers []http.Handler
}

func (s mockS3UploadServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			failRequest(w, 0, "BodyCloseError", fmt.Sprintf("request body close error: %v", err))
		}
	}()

	_, hasUploads := r.URL.Query()["uploads"]

	switch {
	case r.Method == "POST" && hasUploads:
		// CreateMultipartUpload request
		w.Header().Set("Content-Length", strconv.Itoa(len(createUploadResp)))
		w.Write([]byte(createUploadResp))
	case r.Method == "PUT":
		partStr := r.URL.Query().Get("partNumber")
		part, err := strconv.ParseInt(partStr, 10, 64)
		if err != nil {
			failRequest(w, 400, "BadRequest", fmt.Sprintf("unable to parse partNumber, %q, %v", partStr, err))
			return
		}
		if part <= 0 || part > int64(len(s.partHandlers)) {
			failRequest(w, 400, "BadRequest", fmt.Sprintf("invalid partNumber %v", part))
			return
		}
		s.partHandlers[part-1].ServeHTTP(w, r)
	case r.Method == "POST":
		// CompleteMultipartUpload request
		w.Header().Set("Content-Length", strconv.Itoa(len(completeUploadResp)))
		w.Write([]byte(completeUploadResp))
	case r.Method == "DELETE":
		w.Header().Set("Content-Length", strconv.Itoa(len(abortUploadResp)))
		w.Write([]byte(abortUploadResp))
		w.WriteHeader(200)
	default:
		failRequest(w, 400, "BadRequest", fmt.Sprintf("invalid request %v %v", r.Method, r.URL))
	}
}

func createTempFile(t *testing.T, size int64) (*os.File, func(*testing.T), error) {
	file, err := ioutil.TempFile(os.TempDir(), aws.SDKName+t.Name())
	if err != nil {
		return nil, nil, err
	}
	filename := file.Name()
	if err := file.Truncate(size); err != nil {
		return nil, nil, err
	}

	return file,
		func(t *testing.T) {
			if err := file.Close(); err != nil {
				t.Errorf("failed to close temp file, %s, %v", filename, err)
			}
			if err := os.Remove(filename); err != nil {
				t.Errorf("failed to remove temp file, %s, %v", filename, err)
			}
		},
		nil
}

type failPartHandler struct {
	tb                 testing.TB
	failLeft           int
	successPartHandler http.Handler
}

func (h *failPartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			failRequest(w, 0, "BodyCloseError", fmt.Sprintf("request body close error: %v", err))
		}
	}()

	if h.failLeft == 0 && h.successPartHandler != nil {
		h.successPartHandler.ServeHTTP(w, r)
		return
	}

	io.Copy(ioutil.Discard, r.Body)
	failRequest(w, 500, "InternalException", fmt.Sprintf("mock error, partNumber %v", r.URL.Query().Get("partNumber")))
	h.failLeft--
}

func failRequest(w http.ResponseWriter, status int, code, msg string) {
	msg = fmt.Sprintf(baseRequestErrorResp, code, msg)
	w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

type successPartHandler struct {
	tb testing.TB
}

func (h *successPartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			failRequest(w, 0, "BodyCloseError", fmt.Sprintf("request body close error: %v", err))
		}
	}()

	n, err := io.Copy(ioutil.Discard, r.Body)
	if err != nil {
		failRequest(w, 400, "BadRequest", fmt.Sprintf("failed to read body, %v", err))
		return
	}
	contentLength := r.Header.Get("Content-Length")
	expectLength, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		h.tb.Logf("expect content-length, got %q, %v", contentLength, err)
		failRequest(w, 400, "BadRequest", fmt.Sprintf("unable to get content-length %v", err))
		return
	}

	if e, a := expectLength, n; e != a {
		h.tb.Logf("expect content-length to be %v, got %v", e, a)
		failRequest(w, 400, "BadRequest", fmt.Sprintf("content-length and body do not match, %v, %v", e, a))
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(uploadPartResp)))
	w.Write([]byte(uploadPartResp))
}

const createUploadResp = `<CreateMultipartUploadResponse>
  <Bucket>bucket</Bucket>
  <Key>key</Key>
  <UploadId>abc123</UploadId>
</CreateMultipartUploadResponse>`

const uploadPartResp = `<UploadPartResponse>
  <ETag>key</ETag>
</UploadPartResponse>`
const baseRequestErrorResp = `<batchItemError>
  <Code>%s</Code>
  <Message>%s</Message>
  <RequestId>request-id</RequestId>
  <HostId>host-id</HostId>
</batchItemError>`

const completeUploadResp = `<CompleteMultipartUploadResponse>
  <Bucket>bucket</Bucket>
  <Key>key</Key>
  <ETag>key</ETag>
  <Location>https://bucket.us-west-2.amazonaws.com/key</Location>
  <UploadId>abc123</UploadId>
</CompleteMultipartUploadResponse>`

const abortUploadResp = `<AbortMultipartUploadResponse></AbortMultipartUploadResponse>`

func cmpDiff(e, a interface{}) string {
	if !reflect.DeepEqual(e, a) {
		return fmt.Sprintf("%v != %v", e, a)
	}
	return ""
}
