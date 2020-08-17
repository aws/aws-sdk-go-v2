package middleware

import (
	"context"
	"fmt"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/smithy-go/httpbinding"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// UserAgent is an interface for modifying the User-Agent of a request.
type UserAgent interface {
	AddComponent(name string)
	AddVersionedComponent(name, version string)
}

// RequestUserAgent is a build middleware that set the User-Agent for the request.
type RequestUserAgent struct {
	uab *httpbinding.UserAgentBuilder
}

// NewRequestUserAgent returns a new RequestUserAgent which will set the User-Agent for the request.
//
// Default Example:
//   aws-sdk-go/2.3.4 GOOS/linux GOARCH/amd64 GO/go1.14
func NewRequestUserAgent() *RequestUserAgent {
	uab := httpbinding.NewUserAgentBuilder()
	uab.AddVersionedComponent(aws.SDKName, aws.SDKVersion)
	uab.AddVersionedComponent("GOOS", runtime.GOOS)
	uab.AddVersionedComponent("GOARCH", runtime.GOARCH)
	uab.AddVersionedComponent("GO", runtime.Version())
	return &RequestUserAgent{uab: uab}
}

// AddToRequestUserAgent retrieves the RequestUserAgent from the provided stack, or initializes
func AddToRequestUserAgent(stack *middleware.Stack, callback func(UserAgent) error) error {
	id := (&RequestUserAgent{}).ID()
	bm, ok := stack.Build.Get(id)
	if !ok {
		bm = NewRequestUserAgent()
		err := stack.Build.Add(bm, middleware.After)
		if err != nil {
			return err
		}
	}

	requestUserAgent, ok := bm.(*RequestUserAgent)
	if !ok {
		return fmt.Errorf("%T for %s middleware did not match expected type", bm, id)
	}

	return callback(requestUserAgent)
}

// AddComponent adds the component identified by name to the User-Agent string.
func (u *RequestUserAgent) AddComponent(name string) {
	u.uab.AddComponent(name)
}

// AddVersionedComponent adds the componenet identified by the given name and version to the User-Agent string.
func (u *RequestUserAgent) AddVersionedComponent(name, version string) {
	u.uab.AddVersionedComponent(name, version)
}

// ID the name of the middleware.
func (u *RequestUserAgent) ID() string {
	return "RequestUserAgent"
}

// HandleBuild adds or appends the constructed user agent to the request.
func (u *RequestUserAgent) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", in)
	}

	current := req.Header.Get("User-Agent")
	if v := u.uab.Build(); len(current) > 0 {
		current += " " + v
	} else {
		current = v
	}
	req.Header.Set("User-Agent", current)

	return next.HandleBuild(ctx, in)
}
