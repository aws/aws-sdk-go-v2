// +build example

package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// Lists all objects in a bucket using pagination
//
// Usage:
// listObjects <bucket>
func main() {
	if len(os.Args) < 2 {
		exitErrorf("you must specify a bucket")
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}

	svc := s3.New(cfg)

	req := svc.ListObjectsRequest(&s3.ListObjectsInput{Bucket: &os.Args[1]})
	p := req.Paginate()
	for p.Next() {
		page := p.CurrentPage()
		for _, obj := range page.Contents {
			fmt.Println("Object: ", *obj.Key)
		}
	}

	if err := p.Err(); err != nil {
		exitErrorf("failed to list objects, %v", err)
	}
}
