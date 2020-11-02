package customizations

import (
	"github.com/awslabs/smithy-go/middleware"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/s3shared"

	internalendpoints "github.com/aws/aws-sdk-go-v2/service/s3control/internal/endpoints"
)

// EndpointResolver interface for resolving service endpoints.
type EndpointResolver interface {
	ResolveEndpoint(region string, options EndpointResolverOptions) (aws.Endpoint, error)
}

// EndpointResolverOptions is the service endpoint resolver options
type EndpointResolverOptions = internalendpoints.Options

// UpdateEndpointOptions provides the options for the UpdateEndpoint middleware setup.
type UpdateEndpointOptions struct {

	// GetARNInput points to a function that processes an input and returns ARN as string ptr,
	// and bool indicating if ARN is supported or set.
	GetARNInput func(interface{}) (*string, bool)

	// GetOutpostIDInput points to a function that processes an input and returns a outpostID as string ptr,
	// and bool indicating if outpostID is supported or set.
	GetOutpostIDInput func(interface{}) (*string, bool)

	// BackfillAccountID points to a function that validates the input for accountID. If absent, it populates the
	// accountID and returns a copy. If present, but different than passed in accountID value throws an error
	BackfillAccountID func(interface{}, string) (interface{}, error)

	// UpdateARNField points to a function that takes in a copy of input, updates the ARN field with
	// the provided value and returns the input copy, along with a bool indicating if field supports ARN
	UpdateARNField func(interface{}, string) (interface{}, bool)

	// UseARNRegion indicates if region parsed from an ARN should be used.
	UseARNRegion bool

	// UseDualstack instructs if s3 dualstack endpoint config is enabled
	UseDualstack bool

	// EndpointResolver used to resolve endpoints. This may be a custom endpoint resolver
	EndpointResolver EndpointResolver

	// EndpointResolverOptions used by endpoint resolver
	EndpointResolverOptions EndpointResolverOptions
}

// UpdateEndpoint adds the middleware to the middleware stack based on the UpdateEndpointOptions.
func UpdateEndpoint(stack *middleware.Stack, options UpdateEndpointOptions) {
	// validate and backfill account id from ARN
	stack.Initialize.Insert(&BackfillInputMiddleware{
		BackfillAccountID: options.BackfillAccountID,
	}, "OperationInputValidation", middleware.Before)

	// initial arn look up middleware should be before BackfillInputMiddleware
	stack.Initialize.Insert(&s3shared.InitARNLookupMiddleware{
		GetARNValue: options.GetARNInput,
	}, (&BackfillInputMiddleware{}).ID(), middleware.Before)

	// process arn
	stack.Serialize.Insert(&processARNResourceMiddleware{
		UpdateARNField:          options.UpdateARNField,
		UseARNRegion:            options.UseARNRegion,
		UseDualstack:            options.UseDualstack,
		EndpointResolver:        options.EndpointResolver,
		EndpointResolverOptions: options.EndpointResolverOptions,
	}, "OperationSerializer", middleware.Before)

	// outpostID middleware
	stack.Serialize.Insert(&processOutpostIDMiddleware{
		GetOutpostID: options.GetOutpostIDInput,
		UseDualstack: options.UseDualstack,
	}, (&processARNResourceMiddleware{}).ID(), middleware.Before)

	// enable dual stack support
	stack.Serialize.Insert(&s3shared.EnableDualstackMiddleware{
		UseDualstack:     options.UseDualstack,
		DefaultServiceID: "s3-control",
	}, "OperationSerializer", middleware.After)
}
