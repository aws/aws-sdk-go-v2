package s3

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type mockHeadObject struct {
	expires string
}

func (m *mockHeadObject) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Expires": {m.expires},
		},
		Body: http.NoBody,
	}, nil
}

func TestInvalidExpires(t *testing.T) {
	expires := "2023-11-01"
	svc := New(Options{
		HTTPClient: &mockHeadObject{expires},
		Region:     "us-east-1",
	})

	out, err := svc.HeadObject(context.Background(), &HeadObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if out.Expires != nil {
		t.Errorf("out.Expires should be nil, is %s", *out.Expires)
	}
	if aws.ToString(out.ExpiresString) != expires {
		t.Errorf("out.ExpiresString should be %s, is %s", expires, *out.ExpiresString)
	}
}

func TestValidExpires(t *testing.T) {
	exs := "Mon, 02 Jan 2006 15:04:05 GMT"
	ext, err := time.Parse(exs, exs)
	if err != nil {
		t.Fatal(err)
	}

	svc := New(Options{
		HTTPClient: &mockHeadObject{exs},
		Region:     "us-east-1",
	})

	out, err := svc.HeadObject(context.Background(), &HeadObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if aws.ToTime(out.Expires) != ext {
		t.Errorf("out.Expires should be %s, is %s", ext, *out.Expires)
	}
	if aws.ToString(out.ExpiresString) != exs {
		t.Errorf("out.ExpiresString should be %s, is %s", exs, *out.ExpiresString)
	}
}
