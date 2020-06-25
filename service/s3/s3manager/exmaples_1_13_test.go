// +build go1.13

package s3manager_test

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/jviney/aws-sdk-go-v2/aws/external"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/service/s3/s3manager"
)

// ExampleNewUploader_overrideTransport gives an example
// on how to override the default HTTP transport. This can
// be used to tune timeouts such as response headers, or
// write / read buffer usage (go1.13) when writing or reading respectively
// from the net/http transport.
func ExampleNewUploader_overrideTransport() {
	// Create Transport
	tr := &http.Transport{
		ResponseHeaderTimeout: 1 * time.Second,
		WriteBufferSize:       1024 * 1024,
		ReadBufferSize:        1024 * 1024,
	}

	cfg, err := external.LoadDefaultAWSConfig(aws.Config{HTTPClient: &http.Client{Transport: tr}})
	if err != nil {
		panic(fmt.Sprintf("failed to load SDK config: %v", err))
	}

	uploader := s3manager.NewUploader(cfg)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("examplebucket"),
		Key:    aws.String("largeobject"),
		Body:   bytes.NewReader([]byte("large_multi_part_upload")),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
