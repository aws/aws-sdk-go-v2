package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketName      string
	objectPrefix    string
	objectDelimiter string
)

func init() {
	flag.StringVar(&bucketName, "bucket", "", "The `name` of the S3 bucket to list objects from.")
	flag.StringVar(&objectPrefix, "prefix", "", "The optional `object prefix` of the S3 Object keys to list.")
	flag.StringVar(&objectDelimiter, "delimiter", "/",
		"The optional `object key delimiter` used by S3 List objects to group object keys.")
}

// Lists all objects in a bucket using pagination
//
// Usage:
// listObjects <bucket>
func main() {
	flag.Parse()
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	// Load the SDK's configuration from environment and shared config, and
	// create the client with this.
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	// Set the parameters based on the CLI flag inputs.
	params := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}
	if len(objectPrefix) != 0 {
		params.Prefix = &objectPrefix
	}
	if len(objectDelimiter) != 0 {
		params.Delimiter = &objectDelimiter
	}

	// TODO replace this with the code generate paginator when available
	// s3.NewListObjectsV2Paginator()
	p := NewS3ListObjectsV2Paginator(client, params)

	// Iterate through the S3 object pages, printing each object returned.
	var i int
	log.Println("Objects:")
	for p.HasMorePages() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		// Log the objects found
		for _, obj := range page.Contents {
			fmt.Println("Object:", *obj.Key)
		}
	}
}

// S3ListObjectsV2APIClient provides interface for the S3 API client
// ListObjectsV2 operation call.
type S3ListObjectsV2APIClient interface {
	ListObjectsV2(context.Context, *s3.ListObjectsV2Input, ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// S3ListObjectsV2Paginator provides the paginator to paginate S3 ListObjectsV2
// response pages.
type S3ListObjectsV2Paginator struct {
	client S3ListObjectsV2APIClient
	params *s3.ListObjectsV2Input

	nextToken *string
	firstPage bool
}

// NewS3ListObjectsV2Paginator initializes a new S3 ListObjectsV2 Paginator for
// paginating the ListObjectsV2 respones.
func NewS3ListObjectsV2Paginator(client S3ListObjectsV2APIClient, params *s3.ListObjectsV2Input) *S3ListObjectsV2Paginator {
	return &S3ListObjectsV2Paginator{
		client:    client,
		params:    params,
		firstPage: true,
	}
}

// HasMorePages returns true if there are more pages or if the first page has
// not been retrieved yet.
func (p *S3ListObjectsV2Paginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage attempts to retrieve the next page, or returns error if unable to.
func (p *S3ListObjectsV2Paginator) NextPage(ctx context.Context) (
	*s3.ListObjectsV2Output, error,
) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	result, err := p.client.ListObjectsV2(ctx, &params)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.nextToken = result.NextContinuationToken

	return result, nil
}
