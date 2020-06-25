// +build example

package main

import (
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/defaults"
	"github.com/jviney/aws-sdk-go-v2/service/s3"
)

type mockS3Client struct {
	*s3.Client
	index   int
	objects []s3.ListObjectsOutput
}

func (c *mockS3Client) ListObjectsRequest(input *s3.ListObjectsInput) s3.ListObjectsRequest {
	req := c.Client.ListObjectsRequest(input)
	req.Copy = func(v *s3.ListObjectsInput) s3.ListObjectsRequest {
		r := c.Client.ListObjectsRequest(v)
		r.Handlers.Clear()
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

	cfg := defaults.Config()
	cfg.Region = "us-west-2"

	svc.Client = s3.New(cfg)
	svc.objects = objects

	keys, err := getKeys(svc, "foo")
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
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
