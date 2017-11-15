// +build example

package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type client struct {
	s3Client *s3.S3
	bucket   *string
}

// concatenate will contenate key1's object to key2's object under the key testKey
func (c *client) concatenate(key1, key2, key3 string, uploadID *string) (*string, *string, error) {
	// The first part to be uploaded which is represented as part number 1
	req := c.s3Client.UploadPartCopyRequest(&s3.UploadPartCopyInput{
		Bucket:     c.bucket,
		CopySource: aws.String(url.QueryEscape(*c.bucket + "/" + key1)),
		PartNumber: aws.Int64(1),
		Key:        &key3,
		UploadId:   uploadID,
	})
	foo, err := req.Send()
	if err != nil {
		return nil, nil, err
	}

	// The second part that is going to be appended to the newly created testKey
	// object.
	req = c.s3Client.UploadPartCopyRequest(&s3.UploadPartCopyInput{
		Bucket:     c.bucket,
		CopySource: aws.String(url.QueryEscape(*c.bucket + "/" + key2)),
		PartNumber: aws.Int64(2),
		Key:        &key3,
		UploadId:   uploadID,
	})
	bar, err := req.Send()
	if err != nil {
		return nil, nil, err
	}
	// The ETags are needed to complete the process
	return foo.CopyPartResult.ETag, bar.CopyPartResult.ETag, nil
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 4 {
		log.Println("USAGE ERROR: AWS_REGION=us-east-1 go run concatenateObjects.go <bucket> <key for object 1> <key for object 2> <key for output>")
		return
	}

	bucket := os.Args[1]
	key1 := os.Args[2]
	key2 := os.Args[3]
	key3 := os.Args[4]

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}
	svc := s3.New(cfg)

	c := client{svc, &bucket}

	// We let the service know that we want to do a multipart upload
	createReq := c.s3Client.CreateMultipartUploadRequest(&s3.CreateMultipartUploadInput{
		Bucket: &bucket,
		Key:    &key3,
	})
	output, err := createReq.Send()

	if err != nil {
		exitErrorf("failed to create upload, %v", err)
	}

	foo, bar, err := c.concatenate(key1, key2, key3, output.UploadId)
	if err != nil {
		exitErrorf("failed to concatenate, %v", err)
	}

	// We finally complete the multipart upload.
	compReq := c.s3Client.CompleteMultipartUploadRequest(&s3.CompleteMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &key3,
		UploadId: output.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: []s3.CompletedPart{
				{
					ETag:       foo,
					PartNumber: aws.Int64(1),
				},
				{
					ETag:       bar,
					PartNumber: aws.Int64(2),
				},
			},
		},
	})
	if _, err := compReq.Send(); err != nil {
		exitErrorf("failed to complete CompleteMultipartUpload, %v", err)
	}
}
