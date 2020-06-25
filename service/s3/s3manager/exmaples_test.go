package s3manager_test

import (
	"bytes"
	"fmt"

	"github.com/jviney/aws-sdk-go-v2/aws/external"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/service/s3/s3manager"
)

// ExampleNewUploader_overrideReadSeekerProvider gives an example
// on a custom ReadSeekerWriteToProvider can be provided to Uploader
// to define how parts will be buffered in memory.
func ExampleNewUploader_overrideReadSeekerProvider() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load SDK config: %v", err))
	}

	uploader := s3manager.NewUploader(cfg, func(u *s3manager.Uploader) {
		// Define a strategy that will buffer 25 MiB in memory
		u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
	})

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("examplebucket"),
		Key:    aws.String("largeobject"),
		Body:   bytes.NewReader([]byte("large_multi_part_upload")),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
