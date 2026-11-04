package ec2

import (
	"context"
	"fmt"
)

// DescribeVpcEncryptionControlsAPIClient is a client that implements the
// DescribeVpcEncryptionControls operation.
type DescribeVpcEncryptionControlsAPIClient interface {
	DescribeVpcEncryptionControls(context.Context, *DescribeVpcEncryptionControlsInput, ...func(*Options)) (*DescribeVpcEncryptionControlsOutput, error)
}

var _ DescribeVpcEncryptionControlsAPIClient = (*Client)(nil)

// DescribeVpcEncryptionControlsPaginatorOptions is the paginator options for
// DescribeVpcEncryptionControls
type DescribeVpcEncryptionControlsPaginatorOptions struct {
	// The maximum number of items to return for this request. To get the next page of
	// items, make another request with the token returned in the output. For more
	// information, see [Pagination].
	//
	// [Pagination]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Query-Requests.html#api-pagination
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribeVpcEncryptionControlsPaginator is a paginator for
// DescribeVpcEncryptionControls
type DescribeVpcEncryptionControlsPaginator struct {
	options   DescribeVpcEncryptionControlsPaginatorOptions
	client    DescribeVpcEncryptionControlsAPIClient
	params    *DescribeVpcEncryptionControlsInput
	nextToken *string
	firstPage bool
}

// NewDescribeVpcEncryptionControlsPaginator returns a new
// DescribeVpcEncryptionControlsPaginator
func NewDescribeVpcEncryptionControlsPaginator(client DescribeVpcEncryptionControlsAPIClient, params *DescribeVpcEncryptionControlsInput, optFns ...func(*DescribeVpcEncryptionControlsPaginatorOptions)) *DescribeVpcEncryptionControlsPaginator {
	if params == nil {
		params = &DescribeVpcEncryptionControlsInput{}
	}

	options := DescribeVpcEncryptionControlsPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribeVpcEncryptionControlsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribeVpcEncryptionControlsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next DescribeVpcEncryptionControls page.
func (p *DescribeVpcEncryptionControlsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*DescribeVpcEncryptionControlsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	optFns = append([]func(*Options){
		addIsPaginatorUserAgent,
	}, optFns...)
	result, err := p.client.DescribeVpcEncryptionControls(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

// GetVpcResourcesBlockingEncryptionEnforcementAPIClient is a client that
// implements the GetVpcResourcesBlockingEncryptionEnforcement operation.
type GetVpcResourcesBlockingEncryptionEnforcementAPIClient interface {
	GetVpcResourcesBlockingEncryptionEnforcement(context.Context, *GetVpcResourcesBlockingEncryptionEnforcementInput, ...func(*Options)) (*GetVpcResourcesBlockingEncryptionEnforcementOutput, error)
}

var _ GetVpcResourcesBlockingEncryptionEnforcementAPIClient = (*Client)(nil)

// GetVpcResourcesBlockingEncryptionEnforcementPaginatorOptions is the paginator
// options for GetVpcResourcesBlockingEncryptionEnforcement
type GetVpcResourcesBlockingEncryptionEnforcementPaginatorOptions struct {
	// The maximum number of items to return for this request. To get the next page of
	// items, make another request with the token returned in the output. For more
	// information, see [Pagination].
	//
	// [Pagination]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Query-Requests.html#api-pagination
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// GetVpcResourcesBlockingEncryptionEnforcementPaginator is a paginator for
// GetVpcResourcesBlockingEncryptionEnforcement
type GetVpcResourcesBlockingEncryptionEnforcementPaginator struct {
	options   GetVpcResourcesBlockingEncryptionEnforcementPaginatorOptions
	client    GetVpcResourcesBlockingEncryptionEnforcementAPIClient
	params    *GetVpcResourcesBlockingEncryptionEnforcementInput
	nextToken *string
	firstPage bool
}

// NewGetVpcResourcesBlockingEncryptionEnforcementPaginator returns a new
// GetVpcResourcesBlockingEncryptionEnforcementPaginator
func NewGetVpcResourcesBlockingEncryptionEnforcementPaginator(client GetVpcResourcesBlockingEncryptionEnforcementAPIClient, params *GetVpcResourcesBlockingEncryptionEnforcementInput, optFns ...func(*GetVpcResourcesBlockingEncryptionEnforcementPaginatorOptions)) *GetVpcResourcesBlockingEncryptionEnforcementPaginator {
	if params == nil {
		params = &GetVpcResourcesBlockingEncryptionEnforcementInput{}
	}

	options := GetVpcResourcesBlockingEncryptionEnforcementPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &GetVpcResourcesBlockingEncryptionEnforcementPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *GetVpcResourcesBlockingEncryptionEnforcementPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next GetVpcResourcesBlockingEncryptionEnforcement page.
func (p *GetVpcResourcesBlockingEncryptionEnforcementPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*GetVpcResourcesBlockingEncryptionEnforcementOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	optFns = append([]func(*Options){
		addIsPaginatorUserAgent,
	}, optFns...)
	result, err := p.client.GetVpcResourcesBlockingEncryptionEnforcement(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}
