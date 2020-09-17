package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/awslabs/smithy-go"
	smithymiddleware "github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// HTTPClient is a client for sending HTTP requests
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Options is the endpoint client configurable options
type Options struct {
	// The endpoint to retrieve credentials from
	Endpoint string

	// The authorization token to pass in the Authorization header with the request
	AuthorizationToken string

	// Set of options to modify how an operation is invoked. These apply to all
	// operations invoked for this client. Use functional options on operation call to
	// modify this list for per operation behavior.
	APIOptions []func(*smithymiddleware.Stack) error

	// The HTTP client to invoke API calls with. Defaults to client's default HTTP
	// implementation if nil.
	HTTPClient HTTPClient

	// Retryer guides how HTTP requests should be retried in case of recoverable
	// failures. When nil the API client will use a default retryer.
	Retryer retry.Retryer
}

// Copy creates a copy of the API options.
func (o Options) Copy() Options {
	to := o
	to.APIOptions = append([]func(*smithymiddleware.Stack) error{}, o.APIOptions...)
	return to
}

// Client is an client for retrieving AWS credentials from an endpoint
type Client struct {
	options Options
}

// New constructs a new Client from the given options
func New(options Options, optFns ...func(*Options)) *Client {
	options = options.Copy()

	if options.HTTPClient == nil {
		options.HTTPClient = aws.NewBuildableHTTPClient().WithTransportOptions(func(t *http.Transport) {
			t.Proxy = nil
		})
	}

	if options.Retryer == nil {
		options.Retryer = retry.NewStandard()
	}

	for _, fn := range optFns {
		fn(&options)
	}

	client := &Client{
		options: options,
	}

	return client
}

// NewFromConfig constructs a new client using the provided SDK config.
func NewFromConfig(cfg aws.Config, optsFns ...func(options *Options)) *Client {
	opts := Options{
		HTTPClient: cfg.HTTPClient,
		Retryer:    cfg.Retryer,
	}

	return New(opts, optsFns...)
}

// GetCredentials retrieves credentials from credential endpoint
func (c *Client) GetCredentials(ctx context.Context, optFns ...func(*Options)) (*GetCredentialsOutput, error) {
	stack := smithymiddleware.NewStack("GetCredentials", smithyhttp.NewStackRequest)
	options := c.options.Copy()
	for _, fn := range optFns {
		fn(&options)
	}

	stack.Serialize.Add(&serializeOpGetCredential{AuthorizationToken: options.AuthorizationToken}, smithymiddleware.After)
	stack.Build.Add(&buildEndpoint{Endpoint: options.Endpoint}, smithymiddleware.After)
	stack.Deserialize.Add(&deserializeOpGetCredential{}, smithymiddleware.After)

	for _, fn := range options.APIOptions {
		fn(stack)
	}

	handler := smithymiddleware.DecorateHandler(smithyhttp.NewClientHandler(options.HTTPClient), stack)
	result, _, err := handler.Handle(ctx, nil)
	if err != nil {
		return nil, err
	}

	return result.(*GetCredentialsOutput), err
}

// GetCredentialsOutput is the response from the credential endpoint
type GetCredentialsOutput struct {
	Expiration      *time.Time
	AccessKeyID     string
	SecretAccessKey string
	Token           string
}

// EndpointError is an error returned from the endpoint service
type EndpointError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error is the error mesage string
func (e *EndpointError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ErrorCode is the error code returned by the endpoint
func (e *EndpointError) ErrorCode() string {
	return e.Code
}

// ErrorMessage is the error message returned by the endpoint
func (e *EndpointError) ErrorMessage() string {
	return e.Message
}

// ErrorFault indicates error fault classification
func (e *EndpointError) ErrorFault() smithy.ErrorFault {
	return smithy.FaultUnknown
}
