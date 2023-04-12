package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/middleware"
	"testing"
)

var limit int32

type testListMPUMiddleware struct {
	id int
}

func (m *testListMPUMiddleware) ID() string {
	return fmt.Sprintf("mock middleware %d", m.id)
}

type testListOVMiddleware struct {
	id int
}

func (m *testListMPUMiddleware) HandleInitialize(ctx context.Context, input middleware.InitializeInput, next middleware.InitializeHandler) (
	output middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	limit = input.Parameters.(*ListMultipartUploadsInput).MaxUploads
	return middleware.InitializeOutput{Result: &ListMultipartUploadsOutput{}}, metadata, nil
}

func (m *testListOVMiddleware) ID() string {
	return fmt.Sprintf("mock middleware %d", m.id)
}

func (m *testListOVMiddleware) HandleInitialize(ctx context.Context, input middleware.InitializeInput, next middleware.InitializeHandler) (
	output middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	limit = input.Parameters.(*ListObjectVersionsInput).MaxKeys
	return middleware.InitializeOutput{Result: &ListObjectVersionsOutput{}}, metadata, nil
}

func TestListMultipartUploadsPaginator(t *testing.T) {
	cases := map[string]struct {
		limit int32
	}{
		"page limit 5": {
			limit: 5,
		},
		"page limit 10": {
			limit: 10,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := NewFromConfig(aws.Config{})

			paginator := NewListMultipartUploadsPaginator(client, &ListMultipartUploadsInput{
				Bucket: aws.String("test-bucket"),
			}, func(options *ListMultipartUploadsPaginatorOptions) {
				options.Limit = c.limit
			})
			if !paginator.HasMorePages() {
				t.Errorf("Expect paginator has more page, got not")
			}

			paginator.NextPage(context.TODO(), initializeMiddlewareFn(&testListMPUMiddleware{1}))

			testNextPageResult(c.limit, paginator.keyMarker, paginator.uploadIDMarker, t)
		})
	}
}

func TestListObjectVersionsPaginator(t *testing.T) {
	cases := map[string]struct {
		limit int32
	}{
		"page limit 5": {
			limit: 5,
		},
		"page limit 10": {
			limit: 10,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := NewFromConfig(aws.Config{})

			paginator := NewListObjectVersionsPaginator(client, &ListObjectVersionsInput{
				Bucket: aws.String("test-bucket"),
			}, func(options *ListObjectVersionsPaginatorOptions) {
				options.Limit = c.limit
			})
			if !paginator.HasMorePages() {
				t.Errorf("Expect paginator has more page, got not")
			}

			paginator.NextPage(context.TODO(), initializeMiddlewareFn(&testListOVMiddleware{1}))

			testNextPageResult(c.limit, paginator.keyMarker, paginator.versionIDMarker, t)
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
func testNextPageResult(expectLimit int32, keyMarker *string, idMarker *string, t *testing.T) {
	if expectLimit != limit {
		t.Errorf("Expect page limit to be %d, got %d", expectLimit, limit)
	}
	if keyMarker != nil || idMarker != nil {
		t.Errorf("Expect paginator keyMarker and idMarker to be nil, got %s and %s", *keyMarker, *idMarker)
	}
}
