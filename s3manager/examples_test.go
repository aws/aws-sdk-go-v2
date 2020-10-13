package s3manager_test

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/s3manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ExampleNewUploader_overrideReadSeekerProvider gives an example
// on a custom ReadSeekerWriteToProvider can be provided to Uploader
// to define how parts will be buffered in memory.
func ExampleNewUploader_overrideReadSeekerProvider() {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		panic(err)
	}

	uploader := s3manager.NewUploader(s3.NewFromConfig(cfg), func(u *s3manager.Uploader) {
		// Define a strategy that will buffer 25 MiB in memory
		u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
	})

	_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("examplebucket"),
		Key:    aws.String("largeobject"),
		Body:   bytes.NewReader([]byte("large_multi_part_upload")),
	})
	if err != nil {
		panic(err)
	}
}

// ExampleNewUploader_overrideTransport gives an example
// on how to override the default HTTP transport. This can
// be used to tune timeouts such as response headers, or
// write / read buffer usage when writing or reading respectively
// from the net/http transport.
func ExampleNewUploader_overrideTransport() {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// Override Default Transport Values
		o.HTTPClient = aws.NewBuildableHTTPClient().WithTransportOptions(func(tr *http.Transport) {
			tr.ResponseHeaderTimeout = 1 * time.Second
			tr.WriteBufferSize = 1024 * 1024
			tr.ReadBufferSize = 1024 * 1024
		})
	})

	uploader := s3manager.NewUploader(client)

	_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("examplebucket"),
		Key:    aws.String("largeobject"),
		Body:   bytes.NewReader([]byte("large_multi_part_upload")),
	})
	if err != nil {
		panic(err)
	}
}
