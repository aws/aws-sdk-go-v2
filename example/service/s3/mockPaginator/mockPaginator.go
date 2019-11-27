// +build example

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {
	if len(os.Args) < 2 {
		panic("you must specify a bucket")
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config, %v\n", err))
	}

	bucket := os.Args[1]
	svc := s3.New(cfg)
	keys := getKeys(svc, bucket)

	fmt.Printf("keys for bucket %q,\n%v\n", bucket, keys)
}

func getKeys(svc s3iface.ClientAPI, bucket string) []string {
	req := svc.ListObjectsRequest(&types.ListObjectsInput{
		Bucket: &bucket,
	})
	p := s3.NewListObjectsPaginator(req)
	var keys []string
	for p.Next(context.Background()) {
		page := p.CurrentPage()
		for _, obj := range page.Contents {
			keys = append(keys, *obj.Key)
		}
	}
	return keys
}
