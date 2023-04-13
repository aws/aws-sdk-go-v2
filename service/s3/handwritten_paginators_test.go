package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/middleware"
	"testing"
)

var bucket *string
var limit int32
var keyMarker *string
var idMarker *string

type testListMPUMiddleware struct {
	id int
}

type testListOVMiddleware struct {
	id int
}

func (m *testListMPUMiddleware) ID() string {
	return fmt.Sprintf("mock middleware %d", m.id)
}

func (m *testListMPUMiddleware) HandleInitialize(ctx context.Context, input middleware.InitializeInput, next middleware.InitializeHandler) (
	output middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	params := input.Parameters.(*ListMultipartUploadsInput)
	bucket = params.Bucket
	limit = params.MaxUploads
	keyMarker = params.KeyMarker
	idMarker = params.UploadIdMarker
	return middleware.InitializeOutput{Result: &ListMultipartUploadsOutput{}}, metadata, nil
}

func (m *testListOVMiddleware) ID() string {
	return fmt.Sprintf("mock middleware %d", m.id)
}

func (m *testListOVMiddleware) HandleInitialize(ctx context.Context, input middleware.InitializeInput, next middleware.InitializeHandler) (
	output middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	params := input.Parameters.(*ListObjectVersionsInput)
	bucket = params.Bucket
	limit = params.MaxKeys
	keyMarker = params.KeyMarker
	idMarker = params.VersionIdMarker
	return middleware.InitializeOutput{Result: &ListObjectVersionsOutput{}}, metadata, nil
}

type testCase struct {
	bucket    *string
	limit     int32
	keyMarker *string
	idMarker  *string
}

func TestListMultipartUploadsPaginator(t *testing.T) {
	cases := map[string]testCase{
		"page limit 5 without marker": {
			bucket: aws.String("testBucket1"),
			limit:  5,
		},
		"page limit 10 with marker": {
			bucket:    aws.String("testBucket2"),
			limit:     10,
			keyMarker: aws.String("testKey1"),
			idMarker:  aws.String("abc"),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := NewFromConfig(aws.Config{})

			paginator := NewListMultipartUploadsPaginator(client, &ListMultipartUploadsInput{
				Bucket:         c.bucket,
				KeyMarker:      c.keyMarker,
				UploadIdMarker: c.idMarker,
			}, func(options *ListMultipartUploadsPaginatorOptions) {
				options.Limit = c.limit
			})
			if !paginator.HasMorePages() {
				t.Errorf("Expect paginator has more page, got not")
			}

			paginator.NextPage(context.TODO(), initializeMiddlewareFn(&testListMPUMiddleware{1}))

			testNextPageResult(c, paginator.keyMarker, paginator.uploadIDMarker, t)
		})
	}
}

func TestListObjectVersionsPaginator(t *testing.T) {
	cases := map[string]testCase{
		"page limit 5": {
			bucket: aws.String("testBucket3"),
			limit:  5,
		},
		"page limit 10 with marker": {
			bucket:    aws.String("testBucket4"),
			limit:     10,
			keyMarker: aws.String("testKey2"),
			idMarker:  aws.String("def"),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := NewFromConfig(aws.Config{})

			paginator := NewListObjectVersionsPaginator(client, &ListObjectVersionsInput{
				Bucket:          c.bucket,
				KeyMarker:       c.keyMarker,
				VersionIdMarker: c.idMarker,
			}, func(options *ListObjectVersionsPaginatorOptions) {
				options.Limit = c.limit
			})
			if !paginator.HasMorePages() {
				t.Errorf("Expect paginator has more page, got not")
			}

			paginator.NextPage(context.TODO(), initializeMiddlewareFn(&testListOVMiddleware{1}))

			testNextPageResult(c, paginator.keyMarker, paginator.versionIDMarker, t)
		})
	}
}

// insert middleware at the beginning of initialize step to see if page limit
// can be passed to API call's stack input
func initializeMiddlewareFn(initializeMiddleware middleware.InitializeMiddleware) func(*Options) {
	return func(options *Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			return stack.Initialize.Add(initializeMiddleware, middleware.Before)
		})
	}
}

// unit test can not control client API call's output, so just check marker's default nil value
func testNextPageResult(c testCase, pKeyMarker *string, pIdMarker *string, t *testing.T) {
	if c.limit != limit {
		t.Errorf("Expect page limit to be %d, got %d", c.limit, limit)
	}
	if *c.bucket != *bucket {
		t.Errorf("Expect bucket to be %s, got %s", *c.bucket, *bucket)
	}
	if c.keyMarker != nil && *c.keyMarker != *keyMarker {
		t.Errorf("Expect keyMarker to be %s, got %s", *c.keyMarker, *keyMarker)
	}
	if c.idMarker != nil && *c.idMarker != *idMarker {
		t.Errorf("Expect idMarker to be %s, got %s", *c.idMarker, *idMarker)
	}
	if pKeyMarker != nil || pIdMarker != nil {
		t.Errorf("Expect paginator keyMarker and idMarker to be nil, got %s and %s", *pKeyMarker, *pIdMarker)
	}
}
