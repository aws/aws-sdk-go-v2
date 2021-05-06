---
title: "Unit Testing with the AWS SDK for Go V2"
linkTitle: "Testing"
description: "How to mock the AWS SDK for Go V2 when unit testing your application."
weight: 9
---

When using the SDK in your application, you'll want to mock out the SDK for
your application's unit test. Mocking out the SDK allows your test to be
focused on what you want to test, not the internals of the SDK.

To support mocking use Go interfaces instead of concrete service client,
paginators, and waiter types, (e.g. `s3.Client`). This allows your application
to use patterns like dependency injection to test your application logic.

## Mocking Client Operations

In this example, `S3GetObjectAPI` is an interface that defines the set of
{{% alias service=S3 %}} API operations required by the `GetObjectFromS3`
function. `s3GetObjectAPI` is satisfied by the {{% alias service=S3 %}}
client's [GetObject]({{< apiref "service/s3#Client.GetObject" >}}) method.

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

type S3GetObjectAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func GetObjectFromS3(ctx context.Context, api S3GetObjectAPI, bucket, key string) ([]byte, error) {
	object, err := api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer object.Body.Close()

	return ioutil.ReadAll(object.Body)
}
```

To test the `GetObjectFromS3` function, use the `mockGetObjectAPI` to satisfy
the `S3GetObjectAPI` interface definition. Then use the `mockGetObjectAPI` type to mock output
and error responses returned from the service client.

```go
import "testing"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

type mockGetObjectAPI func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)

func (m mockGetObjectAPI) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetObjectFromS3(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) S3GetObjectAPI
		bucket string
		key	string
		expect []byte
	}{
		{
			client: func(t *testing.T) S3GetObjectAPI {
				return mockGetObjectAPI(func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
					t.Helper()
					if params.Bucket == nil {
						t.Fatal("expect bucket to not be nil")
					}
					if e, a := "fooBucket", *params.Bucket; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}
					if params.Key == nil {
						t.Fatal("expect key to not be nil")
					}
					if e, a := "barKey", *params.Key; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}

					return &s3.GetObjectOutput{
						Body: ioutil.NopCloser(bytes.NewReader([]byte("this is the body foo bar baz"))),
					}, nil
				})
			},
			bucket: "fooBucket",
			key:	"barKey",
			expect: []byte("this is the body foo bar baz"),
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content, err := GetObjectFromS3(ctx, tt.client(t), tt.bucket, tt.key)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := tt.expect, content; bytes.Compare(e, a) != 0 {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
```

## Mocking Paginators

Similar to service clients, paginators can be mocked by defining a Go interface
for the paginator. That interface would be used by your application's code. This
allows the SDK's implementation to be used when your application is running,
and a mocked implementation for testing.

In the following example, `ListObjectsV2Pager` is an interface that defines the
behaviors for the {{% alias service=S3 %}}
[ListObjectsV2Paginator]({{< apiref "service/s3#ListObjectsV2Paginator" >}}).
required by `CountObjects` function.

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

type ListObjectsV2Pager interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func CountObjects(ctx context.Context, pager ListObjectsV2Pager) (count int, err error) {
	for pager.HasMorePages() {
		var output *s3.ListObjectsV2Output
		output, err = pager.NextPage(ctx)
		if err != nil {
			return count, err
		}
		count += len(output.Contents)
	}
	return count, nil
}
```

To test `CountObjects`, create the `mockListObjectsV2Pager` type to
satisfy the `ListObjectsV2Pager` interface definition. Then use `mockListObjectsV2Pager`
to replicate the paging behavior of output and error responses from the service
operation paginator.

```go
import "context"
import	"fmt"
import "testing"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

type mockListObjectsV2Pager struct {
	PageNum int
	Pages   []*s3.ListObjectsV2Output
}

func (m *mockListObjectsV2Pager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}

func (m *mockListObjectsV2Pager) NextPage(ctx context.Context, f ...func(*s3.Options)) (output *s3.ListObjectsV2Output, err error) {
	if m.PageNum >= len(m.Pages) {
		return nil, fmt.Errorf("no more pages")
	}
	output = m.Pages[m.PageNum]
	m.PageNum++
	return output, nil
}

func TestCountObjects(t *testing.T) {
	pager := &mockListObjectsV2Pager{
		Pages: []*s3.ListObjectsV2Output{
			{
				KeyCount: 5,
			},
			{
				KeyCount: 10,
			},
			{
				KeyCount: 15,
			},
		},
	}
	objects, err := CountObjects(context.TODO(), pager)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if expect, actual := int64(30), objects; expect != actual {
		t.Errorf("expect %v, got %v", expect, actual)
	}
}
```

