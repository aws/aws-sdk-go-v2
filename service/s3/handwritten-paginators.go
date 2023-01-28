package s3

import (
	"context"
	"fmt"
)

// ListObjectVersionsAPIClient is a client that implements the ListObjectVersions
// operation
type ListObjectVersionsAPIClient interface {
	ListObjectVersions(context.Context, *ListObjectVersionsInput, ...func(*Options)) (*ListObjectVersionsOutput, error)
}

var _ ListObjectVersionsAPIClient = (*Client)(nil)

// ListObjectVersionsPaginatorOptions is the paginator options for ListObjectVersions
type ListObjectVersionsPaginatorOptions struct {
	// (Optional) The maximum number of ResourceRecordSets that you want Amazon Route 53 to
	// return.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListObjectVersionsPaginator is a paginator for ListObjectVersions
type ListObjectVersionsPaginator struct {
	options             ListObjectVersionsPaginatorOptions
	client              ListObjectVersionsAPIClient
	params              *ListObjectVersionsInput
	firstPage           bool
	nextKeyMarker       *string
	nextVersionIdMarker *string
}

// NewListObjectVersionsPaginator returns a new ListObjectVersionsPaginator
func NewListObjectVersionsPaginator(client ListObjectVersionsAPIClient, params *ListObjectVersionsInput, optFns ...func(*ListObjectVersionsPaginatorOptions)) *ListObjectVersionsPaginator {
	if params == nil {
		params = &ListObjectVersionsInput{}
	}

	options := ListObjectVersionsPaginatorOptions{}

	options.Limit = params.MaxKeys

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListObjectVersionsPaginator{
		options:             options,
		client:              client,
		params:              params,
		firstPage:           true,
		nextKeyMarker:       params.KeyMarker,
		nextVersionIdMarker: params.VersionIdMarker,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListObjectVersionsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextKeyMarker != nil && len(*p.nextKeyMarker) != 0)
}

// NextPage retrieves the next ListObjectVersions page.
func (p *ListObjectVersionsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListObjectVersionsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.KeyMarker = p.nextKeyMarker

	params.VersionIdMarker = p.nextVersionIdMarker

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxKeys = *limit

	result, err := p.client.ListObjectVersions(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextKeyMarker
	p.nextKeyMarker = result.NextKeyMarker

	p.nextVersionIdMarker = result.NextVersionIdMarker

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextKeyMarker != nil &&
		*prevToken == *p.nextKeyMarker {
		p.nextKeyMarker = nil
	}

	return result, nil
}
