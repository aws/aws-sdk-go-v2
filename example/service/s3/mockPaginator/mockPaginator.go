// +build example

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jviney/aws-sdk-go-v2/aws/external"
	"github.com/jviney/aws-sdk-go-v2/service/s3"
	"github.com/jviney/aws-sdk-go-v2/service/s3/s3iface"
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
	keys, err := getKeys(svc, bucket)
	if err != nil {
		log.Fatalf("failed to get keys, %v", err)
	}

	fmt.Printf("keys for bucket %q,\n%v\n", bucket, keys)
}

func getKeys(svc s3iface.ClientAPI, bucket string) ([]string, error) {
	req := svc.ListObjectsRequest(&s3.ListObjectsInput{
		Bucket: &bucket,
	})
	p := s3.NewListObjectsPaginator(req)
	keys := []string{}
	for p.Next(context.Background()) {
		page := p.CurrentPage()
		for _, obj := range page.Contents {
			keys = append(keys, *obj.Key)
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
