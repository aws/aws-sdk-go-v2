package customizations

import (
	"context"
	"fmt"
	"github.com/awslabs/smithy-go"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
	"net/url"
)

// AddPredictEndpointMiddleware adds the middleware required to set the endpoint
// based on Predict's PredictEndpoint input member.
func AddPredictEndpointMiddleware(stack *middleware.Stack, endpoint func(interface{}) (*string, error)) {
	stack.Serialize.Insert(&predictEndpointMiddleware{}, "ResolveEndpoint", middleware.After)
}

// predictEndpointMiddleware rewrites the endpoint with whatever is specified in the
// operation input if it is non-nil and non-empty.
type predictEndpointMiddleware struct{
	fetchPredictEndpoint func(interface{}) (*string, error)
}

// ID returns the id for the middleware.
func (*predictEndpointMiddleware) ID() string { return "MachineLearning:PredictEndpoint" }

// HandleSerialize implements the SerializeMiddleware interface.
func (m *predictEndpointMiddleware) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, &smithy.SerializationError{
			Err: fmt.Errorf("unknown request type %T", in.Request),
		}
	}

	endpoint, err := m.fetchPredictEndpoint(in.Parameters)
	if err != nil {
		return out, metadata, &smithy.SerializationError{
			Err: fmt.Errorf("failed to fetch PredictEndpoint value, %v", err),
		}
	}

	if endpoint != nil && len(*endpoint) != 0 {
		uri, err := url.Parse(*endpoint)
		if err != nil {
			return out, metadata, &smithy.SerializationError{
				Err: fmt.Errorf("unable to parse predict endpoint, %v", err),
			}
		}
		req.URL = uri
	}

	return next.HandleSerialize(ctx, in)
}
