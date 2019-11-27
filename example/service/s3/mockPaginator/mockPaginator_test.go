// +build example

package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type mockS3Client struct {
	*s3.Client
	index   int
	objects []types.ListObjectsOutput
}

func (c *mockS3Client) ListObjectsRequest(input *types.ListObjectsInput) s3.ListObjectsRequest {
	req := c.Client.ListObjectsRequest(input)
	req.Copy = func(v *types.ListObjectsInput) s3.ListObjectsRequest {
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
	objects := []types.ListObjectsOutput{
		{
			Contents: []types.Object{
				{
					Key: aws.String("1"),
				},
			},
			NextMarker:  aws.String("marker"),
			IsTruncated: aws.Bool(true),
		},
		{
			Contents: []types.Object{
				{
					Key: aws.String("2"),
				},
			},
			NextMarker:  aws.String("marker"),
			IsTruncated: aws.Bool(true),
		},
		{
			Contents: []types.Object{
				{
					Key: aws.String("3"),
				},
			},
			IsTruncated: aws.Bool(false),
		},
		{
			Contents: []types.Object{
				{
					Key: aws.String("2"),
				},
			},
			NextMarker:  aws.String("marker"),
			IsTruncated: aws.Bool(true),
		},
	}

	svc.Client = s3.New(defaults.Config())
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
