package testing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// TransferManagerLoggingClient is a mock client that can be used to record and stub responses for testing the transfer manager.
type TransferManagerLoggingClient struct {
	// params for upload test

	UploadInvocations []string
	Params            []interface{}

	ConsumeBody bool

	ignoredOperations []string

	PartNum int

	// params for download test

	Data       []byte
	PartsCount int32

	GetObjectInvocations int

	RetrievedRanges []string
	RetrievedParts  []int32

	ErrReaders []TestErrReader
	index      int

	m sync.Mutex

	PutObjectFn               func(*TransferManagerLoggingClient, *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	UploadPartFn              func(*TransferManagerLoggingClient, *s3.UploadPartInput) (*s3.UploadPartOutput, error)
	CreateMultipartUploadFn   func(*TransferManagerLoggingClient, *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error)
	CompleteMultipartUploadFn func(*TransferManagerLoggingClient, *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error)
	AbortMultipartUploadFn    func(*TransferManagerLoggingClient, *s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error)
	GetObjectFn               func(*TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

func (c *TransferManagerLoggingClient) simulateHTTPClientOption(optFns ...func(*s3.Options)) error {

	o := s3.Options{
		HTTPClient: httpDoFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Request: r,
			}, nil
		}),
	}

	for _, fn := range optFns {
		fn(&o)
	}

	_, err := o.HTTPClient.Do(&http.Request{
		URL: &url.URL{
			Scheme:   "https",
			Host:     "mock.amazonaws.com",
			Path:     "/key",
			RawQuery: "foo=bar",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

type httpDoFunc func(*http.Request) (*http.Response, error)

func (f httpDoFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

func (c *TransferManagerLoggingClient) traceOperation(name string, params interface{}) {
	if slices.Contains(c.ignoredOperations, name) {
		return
	}
	c.UploadInvocations = append(c.UploadInvocations, name)
	c.Params = append(c.Params, params)

}

// PutObject is the S3 PutObject API.
func (c *TransferManagerLoggingClient) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.ConsumeBody {
		io.Copy(ioutil.Discard, params.Body)
	}

	c.traceOperation("PutObject", params)

	if err := c.simulateHTTPClientOption(optFns...); err != nil {
		return nil, err
	}

	if c.PutObjectFn != nil {
		return c.PutObjectFn(c, params)
	}

	return &s3.PutObjectOutput{
		VersionId: aws.String("VERSION-ID"),
	}, nil
}

// UploadPart is the S3 UploadPart API.
func (c *TransferManagerLoggingClient) UploadPart(ctx context.Context, params *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.ConsumeBody {
		io.Copy(ioutil.Discard, params.Body)
	}

	c.traceOperation("UploadPart", params)

	if err := c.simulateHTTPClientOption(optFns...); err != nil {
		return nil, err
	}

	if c.UploadPartFn != nil {
		return c.UploadPartFn(c, params)
	}

	return &s3.UploadPartOutput{
		ETag: aws.String(fmt.Sprintf("ETAG%d", *params.PartNumber)),
	}, nil
}

// CreateMultipartUpload is the S3 CreateMultipartUpload API.
func (c *TransferManagerLoggingClient) CreateMultipartUpload(ctx context.Context, params *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
	c.m.Lock()
	defer c.m.Unlock()

	c.traceOperation("CreateMultipartUpload", params)

	if err := c.simulateHTTPClientOption(optFns...); err != nil {
		return nil, err
	}

	if c.CreateMultipartUploadFn != nil {
		return c.CreateMultipartUploadFn(c, params)
	}

	return &s3.CreateMultipartUploadOutput{
		UploadId: aws.String("UPLOAD-ID"),
	}, nil
}

// CompleteMultipartUpload is the S3 CompleteMultipartUpload API.
func (c *TransferManagerLoggingClient) CompleteMultipartUpload(ctx context.Context, params *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
	c.m.Lock()
	defer c.m.Unlock()

	c.traceOperation("CompleteMultipartUpload", params)

	if err := c.simulateHTTPClientOption(optFns...); err != nil {
		return nil, err
	}

	if c.CompleteMultipartUploadFn != nil {
		return c.CompleteMultipartUploadFn(c, params)
	}

	return &s3.CompleteMultipartUploadOutput{
		Location:  aws.String("http://location"),
		VersionId: aws.String("VERSION-ID"),
	}, nil
}

// AbortMultipartUpload is the S3 AbortMultipartUpload API.
func (c *TransferManagerLoggingClient) AbortMultipartUpload(ctx context.Context, params *s3.AbortMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error) {
	c.m.Lock()
	defer c.m.Unlock()

	c.traceOperation("AbortMultipartUpload", params)
	if err := c.simulateHTTPClientOption(optFns...); err != nil {
		return nil, err
	}

	if c.AbortMultipartUploadFn != nil {
		return c.AbortMultipartUploadFn(c, params)
	}

	return &s3.AbortMultipartUploadOutput{}, nil
}

func (c *TransferManagerLoggingClient) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	c.m.Lock()
	defer c.m.Unlock()

	c.GetObjectInvocations++

	if params.Range != nil {
		c.RetrievedRanges = append(c.RetrievedRanges, aws.ToString(params.Range))
	}
	if params.PartNumber != nil {
		c.RetrievedParts = append(c.RetrievedParts, aws.ToInt32(params.PartNumber))
	}

	if c.GetObjectFn != nil {
		return c.GetObjectFn(c, params)
	}

	return &s3.GetObjectOutput{}, nil
}

// NewUploadLoggingClient returns a new TransferManagerLoggingClient for upload testing.
func NewUploadLoggingClient(ignoredOps []string) (*TransferManagerLoggingClient, *[]string, *[]interface{}) {
	c := &TransferManagerLoggingClient{
		ignoredOperations: ignoredOps,
	}

	return c, &c.UploadInvocations, &c.Params
}

func NewDownloadClient() (*TransferManagerLoggingClient, *int, *[]int32, *[]string) {
	c := &TransferManagerLoggingClient{}

	return c, &c.GetObjectInvocations, &c.RetrievedParts, &c.RetrievedRanges
}

var rangeValueRegex = regexp.MustCompile(`bytes=(\d+)-(\d+)`)

func parseRange(rangeValue string) (start, fin int64) {
	rng := rangeValueRegex.FindStringSubmatch(rangeValue)
	start, _ = strconv.ParseInt(rng[1], 10, 64)
	fin, _ = strconv.ParseInt(rng[2], 10, 64)
	return start, fin
}

var RangeGetObjectFn = func(c *TransferManagerLoggingClient, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	start, fin := parseRange(aws.ToString(params.Range))
	fin++

	if fin >= int64(len(c.Data)) {
		fin = int64(len(c.Data))
	}

	bodyBytes := c.Data[start:fin]

	return &s3.GetObjectOutput{
		Body:          ioutil.NopCloser(bytes.NewReader(bodyBytes)),
		ContentRange:  aws.String(fmt.Sprintf("bytes %d-%d/%d", start, fin-1, len(c.Data))),
		ContentLength: aws.Int64(int64(len(bodyBytes))),
	}, nil
}

var ErrRangeGetObjectFn = func(c *TransferManagerLoggingClient, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	out, err := RangeGetObjectFn(c, params)
	c.index++
	if c.index > 1 {
		return &s3.GetObjectOutput{}, fmt.Errorf("s3 service error")
	}
	return out, err
}

var NonRangeGetObjectFn = func(c *TransferManagerLoggingClient, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{
		Body:          ioutil.NopCloser(bytes.NewReader(c.Data[:])),
		ContentLength: aws.Int64(int64(len(c.Data))),
	}, nil
}

var ErrReaderFn = func(c *TransferManagerLoggingClient, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	r := c.ErrReaders[c.index]
	out := &s3.GetObjectOutput{
		Body:          ioutil.NopCloser(&r),
		ContentRange:  aws.String(fmt.Sprintf("bytes %d-%d/%d", 0, r.Len-1, r.Len)),
		ContentLength: aws.Int64(r.Len),
	}
	c.index++
	return out, nil
}

var PartGetObjectFn = func(c *TransferManagerLoggingClient, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{
		Body:          ioutil.NopCloser(bytes.NewReader(c.Data)),
		ContentLength: aws.Int64(8 * 1024 * 1024),
		PartsCount:    aws.Int32(c.PartsCount),
	}, nil
}

var ErrPartGetObjectFn = func(c *TransferManagerLoggingClient, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	out, err := PartGetObjectFn(c, params)
	c.index++
	if c.index > 1 {
		return &s3.GetObjectOutput{}, fmt.Errorf("s3 service error")
	}
	return out, err
}

type TestErrReader struct {
	Buf []byte
	Err error
	Len int64

	off int
}

func (r *TestErrReader) Read(p []byte) (int, error) {
	to := len(r.Buf) - r.off

	n := copy(p, r.Buf[r.off:to])
	r.off += n

	if n < len(p) {
		return n, r.Err

	}

	return n, nil
}
