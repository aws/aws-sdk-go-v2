package s3manager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type listObjectsV2Paginator struct {
	client ListObjectsV2APIClient
	params *s3.ListObjectsV2Input

	done      bool
	nextToken *string
}

func newListObjectsV2Paginator(client ListObjectsV2APIClient, params *s3.ListObjectsV2Input) *listObjectsV2Paginator {
	return &listObjectsV2Paginator{
		client: client,
		params: params,
	}
}

func (p *listObjectsV2Paginator) HasMorePages() bool {
	return !p.done
}

func (p *listObjectsV2Paginator) NextPage(ctx context.Context, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	inCpy := &s3.ListObjectsV2Input{}
	awsutil.Copy(inCpy, p.params)

	if p.nextToken != nil {
		inCpy.ContinuationToken = p.nextToken
	}

	o, err := p.client.ListObjectsV2(ctx, inCpy, optFns...)
	if err != nil {
		return nil, err
	}

	if o.NextContinuationToken != nil {
		var nextToken string
		nextToken = *o.NextContinuationToken
		p.nextToken = &nextToken
	} else {
		p.done = true
	}

	return o, nil
}
