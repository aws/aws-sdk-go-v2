// +build example

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jviney/aws-sdk-go-v2/aws/external"
	"github.com/jviney/aws-sdk-go-v2/service/s3"
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
	p := s3.NewListObjectsPaginator(req)
	for p.Next(context.TODO()) {
		page := p.CurrentPage()
		for _, obj := range page.Contents {
			fmt.Println("Object: ", *obj.Key)
		}
	}

	if err := p.Err(); err != nil {
		exitErrorf("failed to list objects, %v", err)
	}
}
