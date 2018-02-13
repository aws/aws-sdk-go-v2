// +build example

package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockS3Client struct {
	*s3.S3
	index   int
	objects []s3.ListObjectsOutput
}

func (c *mockS3Client) ListObjectsRequest(input *s3.ListObjectsInput) s3.ListObjectsRequest {
	req := c.S3.ListObjectsRequest(input)
	req.Copy = func(v *s3.ListObjectsInput) s3.ListObjectsRequest {
		r := c.S3.ListObjectsRequest(v)
		r.Handlers.Send.Clear()
		r.Handlers.Unmarshal.Clear()
		r.Handlers.UnmarshalMeta.Clear()
		r.Handlers.UnmarshalError.Clear()
		r.Handlers.ValidateResponse.Clear()
		r.Handlers.Send.PushBack(func(r *aws.Request) {
			object := c.objects[c.index]
			r.Data = &object
			c.index++
		})
		return r
	}

	return req
}

func TestListObjectsPagination(t *testing.T) {
	svc := &mockS3Client{}
	objects := []s3.ListObjectsOutput{
		{
			Contents: []s3.Object{
				{
					Key: aws.String("1"),
				},
			},
			NextMarker:  aws.String("marker"),
			IsTruncated: aws.Bool(true),
		},
		{
			Contents: []s3.Object{
				{
					Key: aws.String("2"),
				},
			},
			NextMarker:  aws.String("marker"),
			IsTruncated: aws.Bool(true),
		},
		{
			Contents: []s3.Object{
				{
					Key: aws.String("3"),
				},
			},
			IsTruncated: aws.Bool(false),
		},
		{
			Contents: []s3.Object{
				{
					Key: aws.String("2"),
				},
			},
			NextMarker:  aws.String("marker"),
			IsTruncated: aws.Bool(true),
		},
	}

	svc.S3 = s3.New(unit.Config())
	svc.objects = objects

	keys := getKeys(svc, "foo")
	expected := []string{"1", "2", "3"}

	if e, a := 3, len(keys); e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}

	for i := 0; i < 3; i++ {
		if keys[i] != expected[i] {
			t.Errorf("expected %q, but received %q", expected[i], keys[i])
		}
	}
}
