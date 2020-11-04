package customizations

import (
	"github.com/awslabs/smithy-go/middleware"

	"github.com/aws/aws-sdk-go-v2/service/internal/s3shared"
)

// UpdateEndpointOptions provides the options for the UpdateEndpoint middleware setup.
type UpdateEndpointOptions struct {
	// use dualstack
	UseDualstack bool
}

// UpdateEndpoint adds the middleware to the middleware stack based on the UpdateEndpointOptions.
func UpdateEndpoint(stack *middleware.Stack, options UpdateEndpointOptions) error {
	// enable dual stack support
	return stack.Serialize.Insert(&s3shared.EnableDualstack{
		UseDualstack: options.UseDualstack,
		ServiceID:    "s3-control",
	}, "OperationSerializer", middleware.After)
}
