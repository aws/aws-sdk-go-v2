// +build integration

package s3

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestInteg_PresignURL_PutObject(t *testing.T) {
	key := integrationtest.UniqueID()

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	params := &s3.PutObjectInput{
		Bucket: &setupMetadata.Buckets.Source.Name,
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(`Hello-world`)),
	}

	presignerClient := s3.NewPresignClient(client, func(options *s3.PresignOptions) {
		options.Expires = 600 * time.Second
	})

	presignRequest, err := presignerClient.PresignPutObject(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	// create a http request
	req, err := http.NewRequest(presignRequest.Method, presignRequest.URL, nil)
	if err != nil {
		t.Fatalf("failed to build presigned request, %v", err)
	}

	// assign the signed headers onto the http request
	for k, vs := range presignRequest.SignedHeader {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	// Need to ensure that the content length member is set of the HTTP Request
	// or the request will NOT be transmitted correctly with a content length
	// value across the wire.
	if contLen := req.Header.Get("Content-Length"); len(contLen) > 0 {
		req.ContentLength, _ = strconv.ParseInt(contLen, 10, 64)
	}

	// assign the request body if not nil
	if params.Body != nil {
		req.Body = ioutil.NopCloser(params.Body)
	}

	// Upload the object to S3.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to do PUT request, %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to put S3 object, %d:%s", resp.StatusCode, resp.Status)
	}
}
