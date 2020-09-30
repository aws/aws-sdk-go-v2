package customizations

import (
	"context"
	"fmt"
	"github.com/awslabs/smithy-go"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
	"net/url"
	"regexp"
)

var reSanitizeURL = regexp.MustCompile(`\/%2F\w+%2F`)

// AddSanitizeURLMiddleware add the middleware necessary to clean up Route53 paths.
func AddSanitizeURLMiddleware(stack *middleware.Stack) {
	stack.Serialize.Insert(&sanitizeURLMiddleware{}, "OperationSerializer", middleware.After)
}

// sanitizeURLMiddleware cleans up potential formatting issues in the Route53 path.
//
// Notably it will strip out an excess `/hostedzone/` prefix that can be present in
// the hosted zone id. That excess prefix is there because some route53 apis return
// the id in that format, so this middleware enables round-tripping those values.
type sanitizeURLMiddleware struct{}

// ID returns the id for the middleware.
func (*sanitizeURLMiddleware) ID() string { return "Route53:SanitizeURL" }

// HandleSerialize implements the SerializeMiddleware interface.
func (*sanitizeURLMiddleware) HandleSerialize(
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

	if len(req.URL.RawPath) != 0 {
		req.URL.RawPath = reSanitizeURL.ReplaceAllString(req.URL.RawPath, "/")

		// Update Path so that it reflects the cleaned RawPath
		updated, err := url.Parse(req.URL.RawPath)
		if err != nil {
			return out, metadata, &smithy.SerializationError{
				Err: fmt.Errorf("failed to clean Route53 URL, %v", err),
			}
		}
		req.URL.Path = updated.Path
	}

	return next.HandleSerialize(ctx, in)
}
