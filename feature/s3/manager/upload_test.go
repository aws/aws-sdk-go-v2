package manager_test

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
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/manager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/go-cmp/cmp"
)

// getReaderLength discards the bytes from reader and returns the length
func getReaderLength(r io.Reader) int64 {
	n, _ := io.Copy(ioutil.Discard, r)
	return n
}

func TestUploadOrderMulti(t *testing.T) {
	c, invocations, args := s3testing.NewUploadLoggingClient(nil)
	u := manager.NewUploader(c)

	resp, err := u.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:               aws.String("Bucket"),
		Key:                  aws.String("Key - value"),
		Body:                 bytes.NewReader(buf12MB),
		ServerSideEncryption: "aws:kms",
		SSEKMSKeyId:          aws.String("KmsId"),
		ContentType:          aws.String("content/type"),
	})

	if err != nil {
		t.Errorf("Expected no error but received %v", err)
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart",
		"UploadPart", "CompleteMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(err)
	}

	if e, a := `https://mock.amazonaws.com/key`, resp.Location; e != a {
		t.Errorf("expect %q, got %q", e, a)
	}

	if "UPLOAD-ID" != resp.UploadID {
		t.Errorf("expect %q, got %q", "UPLOAD-ID", resp.UploadID)
	}

	if "VERSION-ID" != *resp.VersionID {
		t.Errorf("expect %q, got %q", "VERSION-ID", *resp.VersionID)
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

		if int32(i+1) != num {
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
	s, ops, args := s3testing.NewUploadLoggingClient(nil)
	mgr := manager.NewUploader(s, func(u *manager.Uploader) {
		u.PartSize = 1024 * 1024 * 7
		u.Concurrency = 1
	})
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf12MB),
	})

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	vals := []string{"CreateMultipartUpload", "UploadPart", "UploadPart", "CompleteMultipartUpload"}
	if !reflect.DeepEqual(vals, *ops) {
		t.Errorf("expect %v, got %v", vals, *ops)
	}

	// Part lengths
	if len := getReaderLength((*args)[1].(*s3.UploadPartInput).Body); 1024*1024*7 != len {
		t.Errorf("expect %d, got %d", 1024*1024*7, len)
	}
	if len := getReaderLength((*args)[2].(*s3.UploadPartInput).Body); 1024*1024*5 != len {
		t.Errorf("expect %d, got %d", 1024*1024*5, len)
	}
}

func TestUploadIncreasePartSize(t *testing.T) {
	s, invocations, args := s3testing.NewUploadLoggingClient(nil)
	mgr := manager.NewUploader(s, func(u *manager.Uploader) {
		u.Concurrency = 1
		u.MaxUploadParts = 2
	})
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf12MB),
	})

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if int64(manager.DefaultDownloadPartSize) != mgr.PartSize {
		t.Errorf("expect %d, got %d", manager.DefaultDownloadPartSize, mgr.PartSize)
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart", "CompleteMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	// Part lengths
	if len := getReaderLength((*args)[1].(*s3.UploadPartInput).Body); (1024*1024*6)+1 != len {
		t.Errorf("expect %d, got %d", (1024*1024*6)+1, len)
	}

	if len := getReaderLength((*args)[2].(*s3.UploadPartInput).Body); (1024*1024*6)-1 != len {
		t.Errorf("expect %d, got %d", (1024*1024*6)-1, len)
	}
}

func TestUploadFailIfPartSizeTooSmall(t *testing.T) {
	mgr := manager.NewUploader(s3.New(s3.Options{}), func(u *manager.Uploader) {
		u.PartSize = 5
	})
	resp, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf12MB),
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
	client, invocations, params := s3testing.NewUploadLoggingClient(nil)
	mgr := manager.NewUploader(client)
	resp, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:               aws.String("Bucket"),
		Key:                  aws.String("Key - value"),
		Body:                 bytes.NewReader(buf2MB),
		ServerSideEncryption: "aws:kms",
		SSEKMSKeyId:          aws.String("KmsId"),
		ContentType:          aws.String("content/type"),
	})

	if err != nil {
		t.Errorf("expect no error but received %v", err)
	}

	if diff := cmp.Diff([]string{"PutObject"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	if e, a := `https://mock.amazonaws.com/key`, resp.Location; e != a {
		t.Errorf("expect %q, got %q", e, a)
	}

	if e := "VERSION-ID"; e != *resp.VersionID {
		t.Errorf("expect %q, got %q", e, *resp.VersionID)
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

func TestUploadOrderSingleFailure(t *testing.T) {
	client, ops, _ := s3testing.NewUploadLoggingClient(nil)

	client.PutObjectFn = func(*s3testing.UploadLoggingClient, *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		return nil, fmt.Errorf("put object failure")
	}

	mgr := manager.NewUploader(client)
	resp, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf2MB),
	})

	if err == nil {
		t.Error("expect error, got nil")
	}

	if diff := cmp.Diff([]string{"PutObject"}, *ops); len(diff) > 0 {
		t.Error(diff)
	}

	if resp != nil {
		t.Errorf("expect response to be nil, got %v", resp)
	}
}

func TestUploadOrderZero(t *testing.T) {
	c, invocations, params := s3testing.NewUploadLoggingClient(nil)
	mgr := manager.NewUploader(c)
	resp, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(make([]byte, 0)),
	})

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmp.Diff([]string{"PutObject"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	if len(resp.Location) == 0 {
		t.Error("expect Location to not be empty")
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

	c.UploadPartFn = func(u *s3testing.UploadLoggingClient, _ *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
		if u.PartNum == 2 {
			return nil, fmt.Errorf("an unexpected error")
		}
		return &s3.UploadPartOutput{ETag: aws.String(fmt.Sprintf("ETAG%d", u.PartNum))}, nil
	}

	mgr := manager.NewUploader(c, func(u *manager.Uploader) {
		u.Concurrency = 1
	})
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf12MB),
	})

	if err == nil {
		t.Error("expect error, got nil")
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart", "AbortMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderMultiFailureOnComplete(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)

	c.CompleteMultipartUploadFn = func(*s3testing.UploadLoggingClient, *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
		return nil, fmt.Errorf("complete multipart error")
	}

	mgr := manager.NewUploader(c, func(u *manager.Uploader) {
		u.Concurrency = 1
	})
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(buf12MB),
	})

	if err == nil {
		t.Error("expect error, got nil")
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart", "UploadPart",
		"CompleteMultipartUpload", "AbortMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderMultiFailureOnCreate(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)

	c.CreateMultipartUploadFn = func(*s3testing.UploadLoggingClient, *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
		return nil, fmt.Errorf("create multipart upload failure")
	}

	mgr := manager.NewUploader(c)
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(make([]byte, 1024*1024*12)),
	})

	if err == nil {
		t.Error("expect error, got nil")
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderMultiFailureLeaveParts(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)

	c.UploadPartFn = func(u *s3testing.UploadLoggingClient, _ *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
		if u.PartNum == 2 {
			return nil, fmt.Errorf("upload part failure")
		}
		return &s3.UploadPartOutput{ETag: aws.String(fmt.Sprintf("ETAG%d", u.PartNum))}, nil
	}

	mgr := manager.NewUploader(c, func(u *manager.Uploader) {
		u.Concurrency = 1
		u.LeavePartsOnError = true
	})
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(make([]byte, 1024*1024*12)),
	})

	if err == nil {
		t.Error("expect error, got nil")
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart"}, *invocations); len(diff) > 0 {
		t.Error(err)
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
	mgr := manager.NewUploader(c)
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   &failreader{times: 1},
	})
	if err == nil {
		t.Fatalf("expect error to not be nil")
	}

	if e, a := "random failure", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v, got %v", e, a)
	}

	if diff := cmp.Diff([]string(nil), *invocations); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderReadFail2(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient([]string{"UploadPart"})
	mgr := manager.NewUploader(c, func(u *manager.Uploader) {
		u.Concurrency = 1
	})
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   &failreader{times: 2},
	})
	if err == nil {
		t.Fatalf("expect error to not be nil")
	}

	if e, a := "random failure", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %v, got %q", e, a)
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "AbortMultipartUpload"}, *invocations); len(diff) > 0 {
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
	mgr := manager.NewUploader(c)
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   &sizedReader{size: 1024 * 1024 * 12},
	})
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart",
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

	if diff := cmp.Diff([]int64{1024 * 1024 * 2, 1024 * 1024 * 5, 1024 * 1024 * 5}, parts); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderMultiBufferedReaderPartial(t *testing.T) {
	c, invocations, params := s3testing.NewUploadLoggingClient(nil)
	mgr := manager.NewUploader(c)
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   &sizedReader{size: 1024 * 1024 * 12, err: io.EOF},
	})
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart",
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

	if diff := cmp.Diff([]int64{1024 * 1024 * 2, 1024 * 1024 * 5, 1024 * 1024 * 5}, parts); len(diff) > 0 {
		t.Error(diff)
	}
}

// TestUploadOrderMultiBufferedReaderEOF tests the edge case where the
// file size is the same as part size.
func TestUploadOrderMultiBufferedReaderEOF(t *testing.T) {
	c, invocations, params := s3testing.NewUploadLoggingClient(nil)
	mgr := manager.NewUploader(c)
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   &sizedReader{size: 1024 * 1024 * 10, err: io.EOF},
	})

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "UploadPart", "UploadPart", "CompleteMultipartUpload"}, *invocations); len(diff) > 0 {
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

	if diff := cmp.Diff([]int64{1024 * 1024 * 5, 1024 * 1024 * 5}, parts); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestUploadOrderMultiBufferedReaderExceedTotalParts(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient([]string{"UploadPart"})
	mgr := manager.NewUploader(c, func(u *manager.Uploader) {
		u.Concurrency = 1
		u.MaxUploadParts = 2
	})
	resp, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   &sizedReader{size: 1024 * 1024 * 12},
	})
	if err == nil {
		t.Fatal("expect error, got nil")
	}

	if resp != nil {
		t.Errorf("expect nil, got %v", resp)
	}

	if diff := cmp.Diff([]string{"CreateMultipartUpload", "AbortMultipartUpload"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	if !strings.Contains(err.Error(), "configured MaxUploadParts (2)") {
		t.Errorf("expect 'configured MaxUploadParts (2)', got %q", err.Error())
	}
}

func TestUploadOrderSingleBufferedReader(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)
	mgr := manager.NewUploader(c)
	resp, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   &sizedReader{size: 1024 * 1024 * 2},
	})

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if diff := cmp.Diff([]string{"PutObject"}, *invocations); len(diff) > 0 {
		t.Error(diff)
	}

	if len(resp.Location) == 0 {
		t.Error("expect a value in Location")
	}

	if len(resp.UploadID) > 0 {
		t.Errorf("expect no value, got %q", resp.UploadID)
	}
}

func TestUploadZeroLenObject(t *testing.T) {
	client, invocations, _ := s3testing.NewUploadLoggingClient(nil)

	mgr := manager.NewUploader(client)
	resp, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   strings.NewReader(""),
	})

	if err != nil {
		t.Errorf("expect no error but received %v", err)
	}
	if diff := cmp.Diff([]string{"PutObject"}, *invocations); len(diff) > 0 {
		t.Errorf("expect request to have been made, but was not, %v", diff)
	}

	// TODO: not needed?
	if len(resp.Location) == 0 {
		t.Error("expect a non-empty string value for Location")
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
	mgr := manager.NewUploader(c, func(u *manager.Uploader) {
		u.Concurrency = 1
		u.PartSize = manager.MinUploadPartSize
	})
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body: &testIncompleteReader{
			Size: manager.MinUploadPartSize + 1,
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
	client, _, _ := s3testing.NewUploadLoggingClient(nil)
	client.UploadPartFn = func(u *s3testing.UploadLoggingClient, params *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
		if params.SSECustomerAlgorithm == nil {
			t.Fatal("SSECustomerAlgoritm should not be nil")
		}
		if params.SSECustomerKey == nil {
			t.Fatal("SSECustomerKey should not be nil")
		}
		return &s3.UploadPartOutput{ETag: aws.String(fmt.Sprintf("ETAG%d", u.PartNum))}, nil
	}

	mgr := manager.NewUploader(client, func(u *manager.Uploader) {
		u.Concurrency = 5
	})

	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:               aws.String("Bucket"),
		Key:                  aws.String("Key"),
		SSECustomerAlgorithm: aws.String("AES256"),
		SSECustomerKey:       aws.String("foo"),
		Body:                 bytes.NewBuffer(make([]byte, 1024*1024*10)),
	})

	if err != nil {
		t.Fatal("Expected no error, but received" + err.Error())
	}
}

func TestUploadWithContextCanceled(t *testing.T) {
	u := manager.NewUploader(s3.New(s3.Options{
		UsePathStyle: true,
		Region:       "mock-region",
	}))

	params := s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   bytes.NewReader(make([]byte, 0)),
	}

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	ctx.Error = fmt.Errorf("context canceled")
	close(ctx.DoneCh)

	_, err := u.Upload(ctx, &params)
	if err == nil {
		t.Fatalf("expect error, got nil")
	}

	if e, a := "canceled", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expected error message to contain %q, but did not %q", e, a)
	}
}

// S3 Uploader incorrectly fails an upload if the content being uploaded
// has a size of MinPartSize * MaxUploadParts.
// Github:  aws/aws-sdk-go#2557
func TestUploadMaxPartsEOF(t *testing.T) {
	c, invocations, _ := s3testing.NewUploadLoggingClient(nil)
	mgr := manager.NewUploader(c, func(u *manager.Uploader) {
		u.Concurrency = 1
		u.PartSize = manager.DefaultUploadPartSize
		u.MaxUploadParts = 2
	})
	f := bytes.NewReader(make([]byte, int(mgr.PartSize)*int(mgr.MaxUploadParts)))

	r1 := io.NewSectionReader(f, 0, manager.DefaultUploadPartSize)
	r2 := io.NewSectionReader(f, manager.DefaultUploadPartSize, 2*manager.DefaultUploadPartSize)
	body := io.MultiReader(r1, r2)

	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("Bucket"),
		Key:    aws.String("Key"),
		Body:   body,
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	expectOps := []string{
		"CreateMultipartUpload",
		"UploadPart",
		"UploadPart",
		"CompleteMultipartUpload",
	}
	if diff := cmp.Diff(expectOps, *invocations); len(diff) > 0 {
		t.Error(diff)
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

func buildFailHandlers(tb testing.TB, parts, retry int) []http.Handler {
	handlers := make([]http.Handler, parts)
	for i := 0; i < len(handlers); i++ {
		handlers[i] = &failPartHandler{
			tb:             tb,
			failsRemaining: retry,
			successHandler: successPartHandler{tb: tb},
		}
	}

	return handlers
}

func TestUploadRetry(t *testing.T) {
	const numParts, retries = 3, 10

	testFile, testFileCleanup, err := createTempFile(t, manager.DefaultUploadPartSize*numParts)
	if err != nil {
		t.Fatalf("failed to create test file, %v", err)
	}
	defer testFileCleanup(t)

	cases := map[string]struct {
		Body         io.Reader
		PartHandlers func(testing.TB) []http.Handler
	}{
		"bytes.Buffer": {
			Body: bytes.NewBuffer(make([]byte, manager.DefaultUploadPartSize*numParts)),
			PartHandlers: func(tb testing.TB) []http.Handler {
				return buildFailHandlers(tb, numParts, retries)
			},
		},
		"bytes.Reader": {
			Body: bytes.NewReader(make([]byte, manager.DefaultUploadPartSize*numParts)),
			PartHandlers: func(tb testing.TB) []http.Handler {
				return buildFailHandlers(tb, numParts, retries)
			},
		},
		"os.File": {
			Body: testFile,
			PartHandlers: func(tb testing.TB) []http.Handler {
				return buildFailHandlers(tb, numParts, retries)
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
				EndpointResolver: s3testing.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL: server.URL,
					}, nil
				}),
				UsePathStyle: true,
				Retryer: retry.NewStandard(func(o *retry.StandardOptions) {
					o.MaxAttempts = retries + 1
				}),
			})

			uploader := manager.NewUploader(client)
			_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   c.Body,
			})

			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

func TestUploadBufferStrategy(t *testing.T) {
	cases := map[string]struct {
		PartSize  int64
		Size      int64
		Strategy  manager.ReadSeekerWriteToProvider
		callbacks int
	}{
		"NoBuffer": {
			PartSize: manager.DefaultUploadPartSize,
			Strategy: nil,
		},
		"SinglePart": {
			PartSize:  manager.DefaultUploadPartSize,
			Size:      manager.DefaultUploadPartSize,
			Strategy:  &recordedBufferProvider{size: int(manager.DefaultUploadPartSize)},
			callbacks: 1,
		},
		"MultiPart": {
			PartSize:  manager.DefaultUploadPartSize,
			Size:      manager.DefaultUploadPartSize * 2,
			Strategy:  &recordedBufferProvider{size: int(manager.DefaultUploadPartSize)},
			callbacks: 2,
		},
	}

	for name, tCase := range cases {
		t.Run(name, func(t *testing.T) {
			client, _, _ := s3testing.NewUploadLoggingClient(nil)
			client.ConsumeBody = true

			uploader := manager.NewUploader(client, func(u *manager.Uploader) {
				u.PartSize = tCase.PartSize
				u.BufferProvider = tCase.Strategy
				u.Concurrency = 1
			})

			expected := s3testing.GetTestBytes(int(tCase.Size))
			_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   bytes.NewReader(expected),
			})
			if err != nil {
				t.Fatalf("failed to upload file: %v", err)
			}

			switch strat := tCase.Strategy.(type) {
			case *recordedBufferProvider:
				if !bytes.Equal(expected, strat.content) {
					t.Errorf("content buffered did not match expected")
				}
				if tCase.callbacks != strat.callbackCount {
					t.Errorf("expected %v, got %v callbacks", tCase.callbacks, strat.callbackCount)
				}
			}
		})
	}
}

func TestUploaderValidARN(t *testing.T) {
	cases := map[string]struct {
		input   s3.PutObjectInput
		wantErr bool
	}{
		"standard bucket": {
			input: s3.PutObjectInput{
				Bucket: aws.String("test-bucket"),
				Key:    aws.String("test-key"),
				Body:   bytes.NewReader([]byte("test body content")),
			},
		},
		"accesspoint": {
			input: s3.PutObjectInput{
				Bucket: aws.String("arn:aws:s3:us-west-2:123456789012:accesspoint/myap"),
				Key:    aws.String("test-key"),
				Body:   bytes.NewReader([]byte("test body content")),
			},
		},
		"outpost accesspoint": {
			input: s3.PutObjectInput{
				Bucket: aws.String("arn:aws:s3-outposts:us-west-2:012345678901:outpost/op-1234567890123456/accesspoint/myaccesspoint"),
				Key:    aws.String("test-key"),
				Body:   bytes.NewReader([]byte("test body content")),
			},
		},
		"s3-object-lambda accesspoint": {
			input: s3.PutObjectInput{
				Bucket: aws.String("arn:aws:s3-object-lambda:us-west-2:123456789012:accesspoint/myap"),
				Key:    aws.String("test-key"),
				Body:   bytes.NewReader([]byte("test body content")),
			},
			wantErr: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			client, _, _ := s3testing.NewUploadLoggingClient(nil)
			client.ConsumeBody = true

			uploader := manager.NewUploader(client)

			_, err := uploader.Upload(context.Background(), &tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("err: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

type mockS3UploadServer struct {
	*http.ServeMux

	tb          testing.TB
	partHandler []http.Handler
}

func newMockS3UploadServer(tb testing.TB, partHandler []http.Handler) *mockS3UploadServer {
	s := &mockS3UploadServer{
		ServeMux:    http.NewServeMux(),
		partHandler: partHandler,
		tb:          tb,
	}

	s.HandleFunc("/", s.handleRequest)

	return s
}

func (s mockS3UploadServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	defer func() {
		closeErr := r.Body.Close()
		if closeErr != nil {
			failRequest(w, 0, "BodyCloseError",
				fmt.Sprintf("request body close error: %v", closeErr))
		}
	}()

	_, hasUploads := r.URL.Query()["uploads"]

	switch {
	case r.Method == "POST" && hasUploads:
		// CreateMultipartUpload
		w.Header().Set("Content-Length", strconv.Itoa(len(createUploadResp)))
		w.Write([]byte(createUploadResp))

	case r.Method == "PUT":
		// UploadPart
		partNumStr := r.URL.Query().Get("partNumber")
		id, err := strconv.Atoi(partNumStr)
		if err != nil {
			failRequest(w, 400, "BadRequest",
				fmt.Sprintf("unable to parse partNumber, %q, %v",
					partNumStr, err))
			return
		}
		id--
		if id < 0 || id >= len(s.partHandler) {
			failRequest(w, 400, "BadRequest",
				fmt.Sprintf("invalid partNumber %v", id))
			return
		}
		s.partHandler[id].ServeHTTP(w, r)

	case r.Method == "POST":
		// CompleteMultipartUpload
		w.Header().Set("Content-Length", strconv.Itoa(len(completeUploadResp)))
		w.Write([]byte(completeUploadResp))

	case r.Method == "DELETE":
		// AbortMultipartUpload
		w.Header().Set("Content-Length", strconv.Itoa(len(abortUploadResp)))
		w.WriteHeader(200)
		w.Write([]byte(abortUploadResp))

	default:
		failRequest(w, 400, "BadRequest",
			fmt.Sprintf("invalid request %v %v", r.Method, r.URL))
	}
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

func (h successPartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		closeErr := r.Body.Close()
		if closeErr != nil {
			failRequest(w, 0, "BodyCloseError",
				fmt.Sprintf("request body close error: %v", closeErr))
		}
	}()

	n, err := io.Copy(ioutil.Discard, r.Body)
	if err != nil {
		failRequest(w, 400, "BadRequest",
			fmt.Sprintf("failed to read body, %v", err))
		return
	}

	contLenStr := r.Header.Get("Content-Length")
	expectLen, err := strconv.ParseInt(contLenStr, 10, 64)
	if err != nil {
		h.tb.Logf("expect content-length, got %q, %v", contLenStr, err)
		failRequest(w, 400, "BadRequest",
			fmt.Sprintf("unable to get content-length %v", err))
		return
	}
	if e, a := expectLen, n; e != a {
		h.tb.Logf("expect %v read, got %v", e, a)
		failRequest(w, 400, "BadRequest",
			fmt.Sprintf(
				"content-length and body do not match, %v, %v", e, a))
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(uploadPartResp)))
	w.Write([]byte(uploadPartResp))
}

type failPartHandler struct {
	tb testing.TB

	failsRemaining int
	successHandler http.Handler
}

func (h *failPartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		closeErr := r.Body.Close()
		if closeErr != nil {
			failRequest(w, 0, "BodyCloseError",
				fmt.Sprintf("request body close error: %v", closeErr))
		}
	}()

	if h.failsRemaining == 0 && h.successHandler != nil {
		h.successHandler.ServeHTTP(w, r)
		return
	}

	io.Copy(ioutil.Discard, r.Body)

	failRequest(w, 500, "InternalException",
		fmt.Sprintf("mock error, partNumber %v", r.URL.Query().Get("partNumber")))

	h.failsRemaining--
}

type recordedBufferProvider struct {
	content       []byte
	size          int
	callbackCount int
}

func (r *recordedBufferProvider) GetWriteTo(seeker io.ReadSeeker) (manager.ReadSeekerWriteTo, func()) {
	b := make([]byte, r.size)
	w := &manager.BufferedReadSeekerWriteTo{BufferedReadSeeker: manager.NewBufferedReadSeeker(seeker, b)}

	return w, func() {
		r.content = append(r.content, b...)
		r.callbackCount++
	}
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
