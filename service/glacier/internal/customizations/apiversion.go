package customizations

import (
	"context"
	"fmt"
	"github.com/awslabs/smithy-go"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

const glacierAPIVersionHeaderKey = "X-Amz-Glacier-Version"

// AddGlacierAPIVersionMiddleware explicitly add handling for the Glacier api version
// middleware to the operation stack.
func AddGlacierAPIVersionMiddleware(stack *middleware.Stack, apiVersion string) {
	stack.Serialize.Add(&GlacierAPIVersionMiddleware{apiVersion: apiVersion}, middleware.Before)
}

// GlacierAPIVersionMiddleware handles automatically setting Glacier's API version header.
type GlacierAPIVersionMiddleware struct{
	apiVersion string
}

// ID returns the id for the middleware.
func (*GlacierAPIVersionMiddleware) ID() string { return "Glacier:APIVersion" }

// HandleSerialize implements the SerializeMiddleware interface
func (m *GlacierAPIVersionMiddleware) HandleSerialize(
	ctx context.Context, input middleware.SerializeInput, next middleware.SerializeHandler,
) (
	output middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := input.Request.(*smithyhttp.Request)
	if !ok {
		return output, metadata, &smithy.SerializationError{
			Err: fmt.Errorf("unknown request type %T", input.Request),
		}
	}

	req.Header.Set(glacierAPIVersionHeaderKey, m.apiVersion)

	return next.HandleSerialize(ctx, input)
}
