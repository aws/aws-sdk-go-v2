package s3_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/service/s3"
)

const (
	metaKeyPrefix = `X-Amz-Meta-`
	utf8KeySuffix = `My-Info`
	utf8Value     = "hello-世界\u0444"
)

func TestPutObjectMetadataWithUnicode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e, a := utf8Value, r.Header.Get(metaKeyPrefix+utf8KeySuffix); e != a {
			t.Errorf("expected %s, but received %s", e, a)
		}
	}))

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := s3.New(cfg)

	req := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("my_bucket"),
		Key:    aws.String("my_key"),
		Body:   strings.NewReader(""),
		Metadata: func() map[string]string {
			v := map[string]string{}
			v[utf8KeySuffix] = utf8Value
			return v
		}(),
	})

	_, err := req.Send(context.Background())
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
}

func TestGetObjectMetadataWithUnicode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(metaKeyPrefix+utf8KeySuffix, utf8Value)
	}))
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := s3.New(cfg)

	req := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("my_bucket"),
		Key:    aws.String("my_key"),
	})
	resp, err := req.Send(context.Background())

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	resp.Body.Close()

	if e, a := utf8Value, resp.Metadata[utf8KeySuffix]; e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}

}
