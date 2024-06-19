package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type reader struct {
	p    []byte
	read bool
}

func (r *reader) Read(p []byte) (int, error) {
	if r.read {
		return 0, io.EOF
	}

	r.read = true
	copy(p, r.p)
	return len(r.p), nil
}

func TestValidatePutObjectContentLength(t *testing.T) {
	for name, cs := range map[string]struct {
		Input     *PutObjectInput
		ExpectErr bool
	}{
		"noseek,nolen": {
			Input: &PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   &reader{p: []byte("foo")},
			},
			ExpectErr: true,
		},
		"noseek,len": {
			Input: &PutObjectInput{
				Bucket:        aws.String("bucket"),
				Key:           aws.String("key"),
				Body:          &reader{p: []byte("foo")},
				ContentLength: aws.Int64(3),
			},
			ExpectErr: false,
		},
		"seek,nolen": {
			Input: &PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   bytes.NewReader([]byte("foo")),
			},
			ExpectErr: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			svc := New(Options{
				Region: "us-east-1",
			})

			_, err := svc.PutObject(context.Background(), cs.Input)
			if cs.ExpectErr && !errors.Is(err, errNoContentLength) {
				t.Errorf("expected errNoContentLength, got %v", err)
			}
			if !cs.ExpectErr && errors.Is(err, errNoContentLength) {
				t.Errorf("expected no errNoContentLength but got it")
			}
		})
	}
}
