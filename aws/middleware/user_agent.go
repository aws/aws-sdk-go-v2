package middleware

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/smithy-go/httpbinding"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

const execEnvVar = `AWS_EXECUTION_ENV`
const execEnvUAKey = `exec-env`

// UserAgent is an interface for modifying the User-Agent of a request.
type UserAgent interface {
	AddKey(name string)
	AddKeyValue(name, version string)
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
	uab.AddKeyValue(aws.SDKName, aws.SDKVersion)
	uab.AddKeyValue("GOOS", runtime.GOOS)
	uab.AddKeyValue("GOARCH", runtime.GOARCH)
	uab.AddKeyValue("GO", runtime.Version())
	if ev := os.Getenv(execEnvVar); len(ev) > 0 {
		uab.AddKeyValue(execEnvUAKey, ev)
	}
	return &RequestUserAgent{uab: uab}
}

// AddUserAgentKey retrieves a RequestUserAgent from the provided stack, or initializes one.
func AddUserAgentKey(callback func(interface {
	AddKey(string)
}) error) func(*middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		requestUserAgent, err := getOrAddRequestUserAgent(stack)
		if err != nil {
			return err
		}
		return callback(requestUserAgent)
	}
}

// AddUserAgentKeyValue retrieves a RequestUserAgent from the provided stack, or initializes one.
func AddUserAgentKeyValue(callback func(interface {
	AddKeyValue(string, string)
}) error) func(*middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		requestUserAgent, err := getOrAddRequestUserAgent(stack)
		if err != nil {
			return err
		}
		return callback(requestUserAgent)
	}
}

func getOrAddRequestUserAgent(stack *middleware.Stack) (*RequestUserAgent, error) {
	id := (&RequestUserAgent{}).ID()
	bm, ok := stack.Build.Get(id)
	if !ok {
		bm = NewRequestUserAgent()
		err := stack.Build.Add(bm, middleware.After)
		if err != nil {
			return nil, err
		}
	}

	requestUserAgent, ok := bm.(*RequestUserAgent)
	if !ok {
		return nil, fmt.Errorf("%T for %s middleware did not match expected type", bm, id)
	}

	return requestUserAgent, nil
}

// AddKey adds the component identified by name to the User-Agent string.
func (u *RequestUserAgent) AddKey(key string) {
	u.uab.AddKey(key)
}

// AddKeyValue adds the key identified by the given name and value to the User-Agent string.
func (u *RequestUserAgent) AddKeyValue(key, value string) {
	u.uab.AddKeyValue(key, value)
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
