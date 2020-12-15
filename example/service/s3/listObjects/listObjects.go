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
	maxKeys         int
)

func init() {
	flag.StringVar(&bucketName, "bucket", "", "The `name` of the S3 bucket to list objects from.")
	flag.StringVar(&objectPrefix, "prefix", "", "The optional `object prefix` of the S3 Object keys to list.")
	flag.StringVar(&objectDelimiter, "delimiter", "",
		"The optional `object key delimiter` used by S3 List objects to group object keys.")
	flag.IntVar(&maxKeys, "max-keys", 0,
		"The maximum number of `keys per page` to retrieve at once.")
}

// Lists all objects in a bucket using pagination
func main() {
	flag.Parse()
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	// Load the SDK's configuration from environment and shared config, and
	// create the client with this.
	cfg, err := config.LoadDefaultConfig(context.TODO())
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
	p := NewS3ListObjectsV2Paginator(client, params, func(o *S3ListObjectsV2PaginatorOptions) {
		if v := int32(maxKeys); v != 0 {
			o.Limit = v
		}
	})

	// Iterate through the S3 object pages, printing each object returned.
	var i int
	log.Println("Objects:")
	for p.HasMorePages() {
		i++

		// Next Page takes a new context for each page retrieval. This is where
		// you could add timeouts or deadlines.
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

// S3ListObjectsV2PaginatorOptions provides the options for configuring the
// ListObjectsV2Paginator.
type S3ListObjectsV2PaginatorOptions struct {
	// The maximum number of keys to return per page.
	Limit int32
}

// S3ListObjectsV2Paginator provides the paginator to paginate S3 ListObjectsV2
// response pages.
type S3ListObjectsV2Paginator struct {
	options S3ListObjectsV2PaginatorOptions

	client S3ListObjectsV2APIClient
	params s3.ListObjectsV2Input

	nextToken *string
	firstPage bool
}

// NewS3ListObjectsV2Paginator initializes a new S3 ListObjectsV2 Paginator for
// paginating the ListObjectsV2 respones.
func NewS3ListObjectsV2Paginator(client S3ListObjectsV2APIClient, params *s3.ListObjectsV2Input, optFns ...func(*S3ListObjectsV2PaginatorOptions)) *S3ListObjectsV2Paginator {
	var options S3ListObjectsV2PaginatorOptions
	for _, fn := range optFns {
		fn(&options)
	}
	p := &S3ListObjectsV2Paginator{
		options:   options,
		client:    client,
		firstPage: true,
	}
	if params != nil {
		p.params = *params
	}
	return p
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

	params := p.params
	if v := p.options.Limit; v != 0 {
		params.MaxKeys = v
	}
	result, err := p.client.ListObjectsV2(ctx, &params)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	if !result.IsTruncated {
		p.nextToken = nil
	} else {
		p.nextToken = result.NextContinuationToken
	}
	p.params.ContinuationToken = p.nextToken

	return result, nil
}
