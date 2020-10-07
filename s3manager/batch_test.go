package s3manager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestHasParity(t *testing.T) {
	cases := []struct {
		o1       *s3.DeleteObjectsInput
		o2       BatchDeleteObject
		expected bool
	}{
		{
			&s3.DeleteObjectsInput{},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{},
			},
			true,
		},
		{
			&s3.DeleteObjectsInput{
				Bucket: aws.String("foo"),
			},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{
					Bucket: aws.String("bar"),
				},
			},
			false,
		},
		{
			&s3.DeleteObjectsInput{},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{
					Bucket: aws.String("foo"),
				},
			},
			false,
		},
		{
			&s3.DeleteObjectsInput{
				Bucket: aws.String("foo"),
			},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{},
			},
			false,
		},
		{
			&s3.DeleteObjectsInput{
				MFA: aws.String("foo"),
			},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{
					MFA: aws.String("bar"),
				},
			},
			false,
		},
		{
			&s3.DeleteObjectsInput{},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{
					MFA: aws.String("foo"),
				},
			},
			false,
		},
		{
			&s3.DeleteObjectsInput{
				MFA: aws.String("foo"),
			},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{},
			},
			false,
		},
		{
			&s3.DeleteObjectsInput{
				RequestPayer: "foo",
			},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{
					RequestPayer: "bar",
				},
			},
			false,
		},
		{
			&s3.DeleteObjectsInput{},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{
					RequestPayer: "foo",
				},
			},
			false,
		},
		{
			&s3.DeleteObjectsInput{
				RequestPayer: "foo",
			},
			BatchDeleteObject{
				Object: &s3.DeleteObjectInput{},
			},
			false,
		},
	}

	for i, c := range cases {
		if result := hasParity(c.o1, c.o2); result != c.expected {
			t.Errorf("Case %d: expected %t, but received %t\n", i, c.expected, result)
		}
	}
}

func TestBatchDelete(t *testing.T) {
	cases := []struct {
		objects  []BatchDeleteObject
		size     int
		expected int
	}{
		{
			[]BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("2"),
						Bucket: aws.String("bucket2"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("3"),
						Bucket: aws.String("bucket3"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("4"),
						Bucket: aws.String("bucket4"),
					},
				},
			},
			1,
			4,
		},
		{
			[]BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("2"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("3"),
						Bucket: aws.String("bucket3"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("4"),
						Bucket: aws.String("bucket3"),
					},
				},
			},
			1,
			4,
		},
		{
			[]BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("2"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("3"),
						Bucket: aws.String("bucket3"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("4"),
						Bucket: aws.String("bucket3"),
					},
				},
			},
			4,
			2,
		},
		{
			[]BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("2"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("3"),
						Bucket: aws.String("bucket3"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("4"),
						Bucket: aws.String("bucket3"),
					},
				},
			},
			10,
			2,
		},
		{
			[]BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("2"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("3"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("4"),
						Bucket: aws.String("bucket3"),
					},
				},
			},
			2,
			3,
		},
	}

	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		count++
	}))
	defer server.Close()

	svc := &mockS3Client{Client: buildS3SvcClient(server.URL)}
	for i, c := range cases {
		batcher := BatchDelete{
			Client:    svc,
			BatchSize: c.size,
		}

		if err := batcher.Delete(context.Background(), &DeleteObjectsIterator{Objects: c.objects}); err != nil {
			panic(err)
		}

		if count != c.expected {
			t.Errorf("Case %d: expected %d, but received %d", i, c.expected, count)
		}

		count = 0
	}
}

func TestBatchDeleteError(t *testing.T) {
	cases := []struct {
		objects            []BatchDeleteObject
		output             s3.DeleteObjectsOutput
		size               int
		expectedErrMessage string
	}{
		{
			objects: []BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
			},
			output: s3.DeleteObjectsOutput{
				Errors: []*types.Error{
					{
						Code:    aws.String("foo code"),
						Message: aws.String("foo error"),
					},
				},
			},
			size:               1,
			expectedErrMessage: "foo error",
		},
		{
			objects: []BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
			},
			output: s3.DeleteObjectsOutput{
				Errors: []*types.Error{{}},
			},
			size:               1,
			expectedErrMessage: errDefaultDeleteBatchMessage,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	index := 0
	svc := &mockS3Client{
		Client: buildS3SvcClient(server.URL),
		deleteObjects: func() (*s3.DeleteObjectsOutput, error) {
			output := &cases[index].output
			index++
			return output, nil
		},
	}
	for _, c := range cases {
		batcher := BatchDelete{
			Client:    svc,
			BatchSize: c.size,
		}

		err := batcher.Delete(context.Background(), &DeleteObjectsIterator{Objects: c.objects})
		if err == nil {
			t.Errorf("expect error, but got nil")
		}

		var bErr BatchError
		if !errors.As(err, &bErr) {
			t.Fatalf("expect %T, got %T", bErr, err)
		}

		errs := bErr.Errors()
		if len(errs) != 1 {
			t.Errorf("expect 1 error, but received %d", len(errs))
		}

		msg := errs[0].Error()
		if e, a := c.expectedErrMessage, msg; !strings.Contains(a, e) {
			t.Errorf("expected %q, but received %q", e, a)
		}
	}
}

type mockS3Client struct {
	*s3.Client

	ListObjectsV2Invocations int

	index         int
	objects       []*s3.ListObjectsV2Output
	deleteObjects func() (*s3.DeleteObjectsOutput, error)
}

func (m *mockS3Client) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	m.ListObjectsV2Invocations++
	object := m.objects[m.index]
	m.index++
	return object, nil
}

func (m *mockS3Client) DeleteObjects(ctx context.Context, input *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
	if m.deleteObjects == nil {
		return m.Client.DeleteObjects(ctx, input, optFns...)
	}

	return m.deleteObjects()
}

func TestNilOrigError(t *testing.T) {
	err := batchItemError{
		bucket: "bucket",
		key:    "key",
	}
	errStr := err.Error()
	const expected1 = `failed to perform batch operation on "key" to "bucket"`
	if errStr != expected1 {
		t.Errorf("Expected %s, but received %s", expected1, errStr)
	}

	err = batchItemError{
		err:    errors.New("foo"),
		bucket: "bucket",
		key:    "key",
	}
	errStr = err.Error()
	const expected2 = "failed to perform batch operation on \"key\" to \"bucket\":\nfoo"
	if errStr != expected2 {
		t.Errorf("Expected %s, but received %s", expected2, errStr)
	}

}

func TestBatchDeleteList(t *testing.T) {
	count := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		count++
	}))
	defer server.Close()

	objects := []*s3.ListObjectsV2Output{
		{
			Contents: []*types.Object{
				{
					Key: aws.String("1"),
				},
			},
			NextContinuationToken: aws.String("marker"),
			IsTruncated:           aws.Bool(true),
		},
		{
			Contents: []*types.Object{
				{
					Key: aws.String("2"),
				},
			},
			NextContinuationToken: aws.String("marker"),
			IsTruncated:           aws.Bool(true),
		},
		{
			Contents: []*types.Object{
				{
					Key: aws.String("3"),
				},
			},
			IsTruncated: aws.Bool(false),
		},
	}

	client := &mockS3Client{Client: buildS3SvcClient(server.URL), objects: objects}
	batcher := BatchDelete{
		Client:    client,
		BatchSize: 1,
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String("bucket"),
	}
	iter := &DeleteListIterator{
		bucket:    input.Bucket,
		paginator: newListObjectsV2Paginator(client, input),
	}

	if err := batcher.Delete(context.Background(), iter); err != nil {
		t.Error(err)
	}

	if count != len(objects) {
		t.Errorf("Expected %d, but received %d", len(objects), count)
	}
}

type mockEndpointResolver func(region string, options s3.ResolverOptions) (aws.Endpoint, error)

func (m mockEndpointResolver) ResolveEndpoint(region string, options s3.ResolverOptions) (aws.Endpoint, error) {
	return m(region, options)
}

func buildS3SvcClient(u string) *s3.Client {
	return s3.New(s3.Options{
		EndpointResolver: mockEndpointResolver(func(region string, options s3.ResolverOptions) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           u,
				SigningRegion: region,
			}, nil
		}),
		UsePathStyle: true,
	})
}

func TestBatchDeleteList_EmptyListObjects(t *testing.T) {
	count := 0

	mockClient := &mockS3Client{}
	mockClient.objects = append(mockClient.objects, &s3.ListObjectsV2Output{Contents: []*types.Object{}})

	batcher := BatchDelete{
		Client: mockClient,
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String("bucket"),
	}

	// Test DeleteListIterator in the case when the ListObjectsRequest responds
	// with an empty listing.

	// We need a new iterator with a fresh Pagination since
	// Pagination.HasNextPage() is always true the first time Pagination.Next()
	// called on it
	iter := NewDeleteListIterator(mockClient, input)

	if err := batcher.Delete(context.Background(), iter); err != nil {
		t.Error(err)
	}
	if mockClient.ListObjectsV2Invocations != 1 {
		t.Errorf("expect count to be 1, got %d", count)
	}
}

func TestBatchDownload(t *testing.T) {
	count := 0
	expected := []struct {
		bucket, key string
	}{
		{
			key:    "1",
			bucket: "bucket1",
		},
		{
			key:    "2",
			bucket: "bucket2",
		},
		{
			key:    "3",
			bucket: "bucket3",
		},
		{
			key:    "4",
			bucket: "bucket4",
		},
	}

	received := []struct {
		bucket, key string
	}{}

	payload := []string{
		"1",
		"2",
		"3",
		"4",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")
		received = append(received, struct{ bucket, key string }{urlParts[1], urlParts[2]})
		w.Write([]byte(payload[count]))
		count++
	}))
	defer server.Close()

	svc := NewDownloader(buildS3SvcClient(server.URL))

	objects := []BatchDownloadObject{
		{
			Object: &s3.GetObjectInput{
				Key:    aws.String("1"),
				Bucket: aws.String("bucket1"),
			},
			Writer: aws.NewWriteAtBuffer(make([]byte, 128)),
		},
		{
			Object: &s3.GetObjectInput{
				Key:    aws.String("2"),
				Bucket: aws.String("bucket2"),
			},
			Writer: aws.NewWriteAtBuffer(make([]byte, 128)),
		},
		{
			Object: &s3.GetObjectInput{
				Key:    aws.String("3"),
				Bucket: aws.String("bucket3"),
			},
			Writer: aws.NewWriteAtBuffer(make([]byte, 128)),
		},
		{
			Object: &s3.GetObjectInput{
				Key:    aws.String("4"),
				Bucket: aws.String("bucket4"),
			},
			Writer: aws.NewWriteAtBuffer(make([]byte, 128)),
		},
	}

	iter := &DownloadObjectsIterator{Objects: objects}
	if err := svc.DownloadWithIterator(context.Background(), iter); err != nil {
		panic(err)
	}

	if count != len(objects) {
		t.Errorf("Expected %d, but received %d", len(objects), count)
	}

	if len(expected) != len(received) {
		t.Errorf("Expected %d, but received %d", len(expected), len(received))
	}

	for i := 0; i < len(expected); i++ {
		if expected[i].key != received[i].key {
			t.Errorf("Expected %q, but received %q", expected[i].key, received[i].key)
		}

		if expected[i].bucket != received[i].bucket {
			t.Errorf("Expected %q, but received %q", expected[i].bucket, received[i].bucket)
		}
	}

	for i, p := range payload {
		b := iter.Objects[i].Writer.(*aws.WriteAtBuffer).Bytes()
		b = bytes.Trim(b, "\x00")

		if string(b) != p {
			t.Errorf("Expected %q, but received %q", p, b)
		}
	}
}

func TestBatchUpload(t *testing.T) {
	count := 0
	expected := []struct {
		bucket, key string
		reqBody     string
	}{
		{
			key:     "1",
			bucket:  "bucket1",
			reqBody: "1",
		},
		{
			key:     "2",
			bucket:  "bucket2",
			reqBody: "2",
		},
		{
			key:     "3",
			bucket:  "bucket3",
			reqBody: "3",
		},
		{
			key:     "4",
			bucket:  "bucket4",
			reqBody: "4",
		},
	}

	received := []struct {
		bucket, key, reqBody string
	}{}

	payload := []string{
		"a",
		"b",
		"c",
		"d",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.Path, "/")

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		received = append(received, struct{ bucket, key, reqBody string }{urlParts[1], urlParts[2], string(b)})
		w.Write([]byte(payload[count]))

		count++
	}))
	defer server.Close()

	svc := NewUploader(buildS3SvcClient(server.URL))

	objects := []BatchUploadObject{
		{
			Object: &s3.PutObjectInput{
				Key:    aws.String("1"),
				Bucket: aws.String("bucket1"),
				Body:   bytes.NewBuffer([]byte("1")),
			},
		},
		{
			Object: &s3.PutObjectInput{
				Key:    aws.String("2"),
				Bucket: aws.String("bucket2"),
				Body:   bytes.NewBuffer([]byte("2")),
			},
		},
		{
			Object: &s3.PutObjectInput{
				Key:    aws.String("3"),
				Bucket: aws.String("bucket3"),
				Body:   bytes.NewBuffer([]byte("3")),
			},
		},
		{
			Object: &s3.PutObjectInput{
				Key:    aws.String("4"),
				Bucket: aws.String("bucket4"),
				Body:   bytes.NewBuffer([]byte("4")),
			},
		},
	}

	iter := &UploadObjectsIterator{Objects: objects}
	if err := svc.UploadWithIterator(context.Background(), iter); err != nil {
		panic(err)
	}

	if count != len(objects) {
		t.Errorf("Expected %d, but received %d", len(objects), count)
	}

	if len(expected) != len(received) {
		t.Errorf("Expected %d, but received %d", len(expected), len(received))
	}

	for i := 0; i < len(expected); i++ {
		if expected[i].key != received[i].key {
			t.Errorf("Expected %q, but received %q", expected[i].key, received[i].key)
		}

		if expected[i].bucket != received[i].bucket {
			t.Errorf("Expected %q, but received %q", expected[i].bucket, received[i].bucket)
		}

		if expected[i].reqBody != received[i].reqBody {
			t.Errorf("Expected %q, but received %q", expected[i].reqBody, received[i].reqBody)
		}
	}
}

type mockClient struct {
	*s3.Client
	Put    func() (*s3.PutObjectOutput, error)
	Get    func() (*s3.GetObjectOutput, error)
	List   func() (*s3.ListObjectsV2Output, error)
	Delete func() (*s3.DeleteObjectsOutput, error)
}

type response struct {
	out interface{}
	err error
}

func (client *mockClient) DeleteObjects(context.Context, *s3.DeleteObjectsInput, ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
	return client.Delete()
}

func (client *mockClient) GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return client.Get()
}

func (client *mockClient) PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return client.Put()
}

func (client *mockClient) ListObjectsV2(context.Context, *s3.ListObjectsV2Input, ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return client.List()
}

func TestBatchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer server.Close()

	index := 0
	responses := []response{
		{
			&s3.PutObjectOutput{},
			errors.New("Foo"),
		},
		{
			&s3.PutObjectOutput{},
			nil,
		},
		{
			&s3.PutObjectOutput{},
			nil,
		},
		{
			&s3.PutObjectOutput{},
			errors.New("Bar"),
		},
	}

	client := &mockClient{
		Client: buildS3SvcClient(server.URL),
		Put: func() (*s3.PutObjectOutput, error) {
			resp := responses[index]
			index++
			return resp.out.(*s3.PutObjectOutput), resp.err
		},
		List: func() (*s3.ListObjectsV2Output, error) {
			resp := responses[index]
			index++
			return resp.out.(*s3.ListObjectsV2Output), resp.err
		},
	}
	uploader := NewUploader(client)

	objects := []BatchUploadObject{
		{
			Object: &s3.PutObjectInput{
				Key:    aws.String("1"),
				Bucket: aws.String("bucket1"),
				Body:   bytes.NewBuffer([]byte("1")),
			},
		},
		{
			Object: &s3.PutObjectInput{
				Key:    aws.String("2"),
				Bucket: aws.String("bucket2"),
				Body:   bytes.NewBuffer([]byte("2")),
			},
		},
		{
			Object: &s3.PutObjectInput{
				Key:    aws.String("3"),
				Bucket: aws.String("bucket3"),
				Body:   bytes.NewBuffer([]byte("3")),
			},
		},
		{
			Object: &s3.PutObjectInput{
				Key:    aws.String("4"),
				Bucket: aws.String("bucket4"),
				Body:   bytes.NewBuffer([]byte("4")),
			},
		},
	}

	iter := &UploadObjectsIterator{Objects: objects}
	if err := uploader.UploadWithIterator(context.Background(), iter); err != nil {
		var bErr BatchError

		if !errors.As(err, &bErr) {
			t.Errorf("expect BatchError, got %T", err)
		} else {
			be := bErr.Errors()

			if len(be) != 2 {
				t.Errorf("expect 2 errors, got %d", len(be))
			}

			expected := []struct {
				bucket, key string
			}{
				{
					"bucket1",
					"1",
				},
				{
					"bucket4",
					"4",
				},
			}
			for i, expect := range expected {
				var bi BatchItemError
				if !errors.As(be[i], &bi) {
					t.Errorf("expect BatchItemError, got %T", be[i])
				}

				if bi.Bucket() != expect.bucket {
					t.Errorf("case %d: invalid bucket expect %s, but received %s", i, expect.bucket, bi.Bucket())
				}

				if bi.Key() != expect.key {
					t.Errorf("case %d: invalid key expect %s, but received %s", i, expect.key, bi.Key())
				}
			}
		}
	} else {
		t.Error("Expected error, but received nil")
	}

	if index != len(objects) {
		t.Errorf("Expected %d, but received %d", len(objects), index)
	}

}

type testAfterDeleteIter struct {
	afterDelete bool
	next        bool
}

func (iter *testAfterDeleteIter) Next() bool {
	next := !iter.next
	iter.next = !iter.next
	return next
}

func (iter *testAfterDeleteIter) Err() error {
	return nil
}

func (iter *testAfterDeleteIter) DeleteObject() BatchDeleteObject {
	return BatchDeleteObject{
		Object: &s3.DeleteObjectInput{
			Bucket: aws.String("foo"),
			Key:    aws.String("foo"),
		},
		After: func() error {
			iter.afterDelete = true
			return nil
		},
	}
}

type testAfterDownloadIter struct {
	afterDownload bool
	next          bool
}

func (iter *testAfterDownloadIter) Next() bool {
	next := !iter.next
	iter.next = !iter.next
	return next
}

func (iter *testAfterDownloadIter) Err() error {
	return nil
}

func (iter *testAfterDownloadIter) DownloadObject() BatchDownloadObject {
	return BatchDownloadObject{
		Object: &s3.GetObjectInput{
			Bucket: aws.String("foo"),
			Key:    aws.String("foo"),
		},
		Writer: aws.NewWriteAtBuffer([]byte{}),
		After: func() error {
			iter.afterDownload = true
			return nil
		},
	}
}

type testAfterUploadIter struct {
	afterUpload bool
	next        bool
}

func (iter *testAfterUploadIter) Next() bool {
	next := !iter.next
	iter.next = !iter.next
	return next
}

func (iter *testAfterUploadIter) Err() error {
	return nil
}

func (iter *testAfterUploadIter) UploadObject() BatchUploadObject {
	return BatchUploadObject{
		Object: &s3.PutObjectInput{
			Bucket: aws.String("foo"),
			Key:    aws.String("foo"),
			Body:   strings.NewReader("bar"),
		},
		After: func() error {
			iter.afterUpload = true
			return nil
		},
	}
}

func TestAfter(t *testing.T) {
	index := 0
	responses := []response{
		{
			&s3.PutObjectOutput{},
			nil,
		},
		{
			&s3.GetObjectOutput{
				ContentLength: aws.Int64(4),
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("test"))),
			},
			nil,
		},
		{
			&s3.DeleteObjectsOutput{},
			nil,
		},
	}

	client := &mockClient{
		Put: func() (*s3.PutObjectOutput, error) {
			resp := responses[index]
			index++
			return resp.out.(*s3.PutObjectOutput), resp.err
		},
		Get: func() (*s3.GetObjectOutput, error) {
			resp := responses[index]
			index++
			return resp.out.(*s3.GetObjectOutput), resp.err
		},
		List: func() (*s3.ListObjectsV2Output, error) {
			resp := responses[index]
			index++
			return resp.out.(*s3.ListObjectsV2Output), resp.err
		},
		Delete: func() (*s3.DeleteObjectsOutput, error) {
			resp := responses[index]
			index++
			return resp.out.(*s3.DeleteObjectsOutput), resp.err
		},
	}
	uploader := NewUploader(client)
	downloader := NewDownloader(client)
	deleter := NewBatchDelete(client)

	deleteIter := &testAfterDeleteIter{}
	downloadIter := &testAfterDownloadIter{}
	uploadIter := &testAfterUploadIter{}

	if err := uploader.UploadWithIterator(context.Background(), uploadIter); err != nil {
		t.Error(err)
	}

	if err := downloader.DownloadWithIterator(context.Background(), downloadIter); err != nil {
		t.Error(err)
	}

	if err := deleter.Delete(context.Background(), deleteIter); err != nil {
		t.Error(err)
	}

	if !deleteIter.afterDelete {
		t.Error("expect 'afterDelete' to be true, but received false")
	}

	if !downloadIter.afterDownload {
		t.Error("expect 'afterDownload' to be true, but received false")
	}

	if !uploadIter.afterUpload {
		t.Error("expect 'afterUpload' to be true, but received false")
	}
}

// #1790 bug
func TestBatchDeleteContext(t *testing.T) {
	cases := []struct {
		objects     []BatchDeleteObject
		batchSize   int
		expected    int
		earlyCancel bool
		checkError  func(error) error
	}{
		0: {
			objects: []BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("2"),
						Bucket: aws.String("bucket2"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("3"),
						Bucket: aws.String("bucket3"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("4"),
						Bucket: aws.String("bucket4"),
					},
				},
			},
			batchSize:   1,
			expected:    0,
			earlyCancel: true,
			checkError: func(err error) error {
				var bErr BatchError
				if !errors.As(err, &bErr) {
					return fmt.Errorf("expect %T, got %T, %v", bErr, err, err)
				}

				errs := bErr.Errors()
				if len(errs) != 4 {
					return fmt.Errorf("expect 4 batch errors, got %d", len(errs))
				}

				for _, err := range errs {
					var iErr BatchItemError
					if !errors.As(err, &iErr) {
						return fmt.Errorf("expect %T, got %T, %v", iErr, err, err)
					}

					if e, a := "context canceled", iErr.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %v, got %v", e, a)
					}
				}

				return nil
			},
		},
		1: {
			objects: []BatchDeleteObject{
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("1"),
						Bucket: aws.String("bucket1"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("2"),
						Bucket: aws.String("bucket2"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("3"),
						Bucket: aws.String("bucket3"),
					},
				},
				{
					Object: &s3.DeleteObjectInput{
						Key:    aws.String("4"),
						Bucket: aws.String("bucket4"),
					},
				},
			},
			batchSize: 1,
			expected:  4,
			checkError: func(err error) error {
				if err != nil {
					return fmt.Errorf("expect no error, got %v", err)
				}
				return nil
			},
		},
	}

	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		count++
	}))
	defer server.Close()

	client := &mockS3Client{Client: buildS3SvcClient(server.URL)}
	for i, c := range cases {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if c.earlyCancel {
			cancel()
		}

		batcher := BatchDelete{
			Client:    client,
			BatchSize: c.batchSize,
		}

		err := batcher.Delete(ctx, &DeleteObjectsIterator{Objects: c.objects})
		if terr := c.checkError(err); terr != nil {
			t.Fatalf("%d, %s", i, terr)
		}

		if count != c.expected {
			t.Errorf("Case %d: expected %d, but received %d", i, c.expected, count)
		}

		count = 0
	}
}
