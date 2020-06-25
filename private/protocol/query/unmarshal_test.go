package query_test

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/retry"
	"github.com/jviney/aws-sdk-go-v2/aws/signer/v4"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/private/protocol"
	"github.com/jviney/aws-sdk-go-v2/private/protocol/query"
	"github.com/jviney/aws-sdk-go-v2/private/protocol/xml/xmlutil"
	"github.com/jviney/aws-sdk-go-v2/private/util"
)

var _ bytes.Buffer // always import bytes
var _ http.Request
var _ json.Marshaler
var _ time.Time
var _ xmlutil.XMLNode
var _ xml.Attr
var _ = ioutil.Discard
var _ = util.Trim("")
var _ = url.Values{}
var _ = io.EOF
var _ = aws.String
var _ = fmt.Println
var _ = reflect.Value{}

func init() {
	protocol.RandReader = &awstesting.ZeroReader{}
}

// OutputService1ProtocolTest provides the API operation methods for making requests to
// OutputService1ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService1ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice1protocoltest.New(myConfig)
func NewOutputService1ProtocolTest(config aws.Config) *OutputService1ProtocolTest {
	svc := &OutputService1ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService1ProtocolTest",
				ServiceID:     "OutputService1ProtocolTest",
				EndpointsID:   "outputservice1protocoltest",
				SigningName:   "outputservice1protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService1ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService1TestShapeOutputService1TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService1TestShapeOutputService1TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Char *string `type:"character"`

	Double *float64 `type:"double"`

	FalseBool *bool `type:"boolean"`

	Float *float64 `type:"float"`

	Long *int64 `type:"long"`

	Num *int64 `locationName:"FooNum" type:"integer"`

	Str *string `type:"string"`

	Timestamp *time.Time `type:"timestamp"`

	TrueBool *bool `type:"boolean"`
}

const opOutputService1TestCaseOperation1 = "OperationName"

// OutputService1TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService1TestCaseOperation1Request.
//    req := client.OutputService1TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation1Request(input *OutputService1TestShapeOutputService1TestCaseOperation1Input) OutputService1TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService1TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService1TestShapeOutputService1TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService1TestShapeOutputService1TestCaseOperation1Output{})

	return OutputService1TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService1TestCaseOperation1Request}
}

// OutputService1TestCaseOperation1Request is the request type for the
// OutputService1TestCaseOperation1 API operation.
type OutputService1TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService1TestShapeOutputService1TestCaseOperation1Input
	Copy  func(*OutputService1TestShapeOutputService1TestCaseOperation1Input) OutputService1TestCaseOperation1Request
}

// Send marshals and sends the OutputService1TestCaseOperation1 API request.
func (r OutputService1TestCaseOperation1Request) Send(ctx context.Context) (*OutputService1TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService1TestCaseOperation1Response{
		OutputService1TestShapeOutputService1TestCaseOperation1Output: r.Request.Data.(*OutputService1TestShapeOutputService1TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService1TestCaseOperation1Response is the response type for the
// OutputService1TestCaseOperation1 API operation.
type OutputService1TestCaseOperation1Response struct {
	*OutputService1TestShapeOutputService1TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService1TestCaseOperation1 request.
func (r *OutputService1TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService2ProtocolTest provides the API operation methods for making requests to
// OutputService2ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService2ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice2protocoltest.New(myConfig)
func NewOutputService2ProtocolTest(config aws.Config) *OutputService2ProtocolTest {
	svc := &OutputService2ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService2ProtocolTest",
				ServiceID:     "OutputService2ProtocolTest",
				EndpointsID:   "outputservice2protocoltest",
				SigningName:   "outputservice2protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService2ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService2TestShapeOutputService2TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService2TestShapeOutputService2TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Num *int64 `type:"integer"`

	Str *string `type:"string"`
}

const opOutputService2TestCaseOperation1 = "OperationName"

// OutputService2TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService2TestCaseOperation1Request.
//    req := client.OutputService2TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService2ProtocolTest) OutputService2TestCaseOperation1Request(input *OutputService2TestShapeOutputService2TestCaseOperation1Input) OutputService2TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService2TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService2TestShapeOutputService2TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService2TestShapeOutputService2TestCaseOperation1Output{})

	return OutputService2TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService2TestCaseOperation1Request}
}

// OutputService2TestCaseOperation1Request is the request type for the
// OutputService2TestCaseOperation1 API operation.
type OutputService2TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService2TestShapeOutputService2TestCaseOperation1Input
	Copy  func(*OutputService2TestShapeOutputService2TestCaseOperation1Input) OutputService2TestCaseOperation1Request
}

// Send marshals and sends the OutputService2TestCaseOperation1 API request.
func (r OutputService2TestCaseOperation1Request) Send(ctx context.Context) (*OutputService2TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService2TestCaseOperation1Response{
		OutputService2TestShapeOutputService2TestCaseOperation1Output: r.Request.Data.(*OutputService2TestShapeOutputService2TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService2TestCaseOperation1Response is the response type for the
// OutputService2TestCaseOperation1 API operation.
type OutputService2TestCaseOperation1Response struct {
	*OutputService2TestShapeOutputService2TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService2TestCaseOperation1 request.
func (r *OutputService2TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService3ProtocolTest provides the API operation methods for making requests to
// OutputService3ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService3ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice3protocoltest.New(myConfig)
func NewOutputService3ProtocolTest(config aws.Config) *OutputService3ProtocolTest {
	svc := &OutputService3ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService3ProtocolTest",
				ServiceID:     "OutputService3ProtocolTest",
				EndpointsID:   "outputservice3protocoltest",
				SigningName:   "outputservice3protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService3ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService3TestShapeOutputService3TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService3TestShapeOutputService3TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	// Blob is automatically base64 encoded/decoded by the SDK.
	Blob []byte `type:"blob"`
}

const opOutputService3TestCaseOperation1 = "OperationName"

// OutputService3TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService3TestCaseOperation1Request.
//    req := client.OutputService3TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService3ProtocolTest) OutputService3TestCaseOperation1Request(input *OutputService3TestShapeOutputService3TestCaseOperation1Input) OutputService3TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService3TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService3TestShapeOutputService3TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService3TestShapeOutputService3TestCaseOperation1Output{})

	return OutputService3TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService3TestCaseOperation1Request}
}

// OutputService3TestCaseOperation1Request is the request type for the
// OutputService3TestCaseOperation1 API operation.
type OutputService3TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService3TestShapeOutputService3TestCaseOperation1Input
	Copy  func(*OutputService3TestShapeOutputService3TestCaseOperation1Input) OutputService3TestCaseOperation1Request
}

// Send marshals and sends the OutputService3TestCaseOperation1 API request.
func (r OutputService3TestCaseOperation1Request) Send(ctx context.Context) (*OutputService3TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService3TestCaseOperation1Response{
		OutputService3TestShapeOutputService3TestCaseOperation1Output: r.Request.Data.(*OutputService3TestShapeOutputService3TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService3TestCaseOperation1Response is the response type for the
// OutputService3TestCaseOperation1 API operation.
type OutputService3TestCaseOperation1Response struct {
	*OutputService3TestShapeOutputService3TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService3TestCaseOperation1 request.
func (r *OutputService3TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService4ProtocolTest provides the API operation methods for making requests to
// OutputService4ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService4ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice4protocoltest.New(myConfig)
func NewOutputService4ProtocolTest(config aws.Config) *OutputService4ProtocolTest {
	svc := &OutputService4ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService4ProtocolTest",
				ServiceID:     "OutputService4ProtocolTest",
				EndpointsID:   "outputservice4protocoltest",
				SigningName:   "outputservice4protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService4ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService4TestShapeOutputService4TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService4TestShapeOutputService4TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	ListMember []string `type:"list"`
}

const opOutputService4TestCaseOperation1 = "OperationName"

// OutputService4TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService4TestCaseOperation1Request.
//    req := client.OutputService4TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService4ProtocolTest) OutputService4TestCaseOperation1Request(input *OutputService4TestShapeOutputService4TestCaseOperation1Input) OutputService4TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService4TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService4TestShapeOutputService4TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService4TestShapeOutputService4TestCaseOperation1Output{})

	return OutputService4TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService4TestCaseOperation1Request}
}

// OutputService4TestCaseOperation1Request is the request type for the
// OutputService4TestCaseOperation1 API operation.
type OutputService4TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService4TestShapeOutputService4TestCaseOperation1Input
	Copy  func(*OutputService4TestShapeOutputService4TestCaseOperation1Input) OutputService4TestCaseOperation1Request
}

// Send marshals and sends the OutputService4TestCaseOperation1 API request.
func (r OutputService4TestCaseOperation1Request) Send(ctx context.Context) (*OutputService4TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService4TestCaseOperation1Response{
		OutputService4TestShapeOutputService4TestCaseOperation1Output: r.Request.Data.(*OutputService4TestShapeOutputService4TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService4TestCaseOperation1Response is the response type for the
// OutputService4TestCaseOperation1 API operation.
type OutputService4TestCaseOperation1Response struct {
	*OutputService4TestShapeOutputService4TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService4TestCaseOperation1 request.
func (r *OutputService4TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService5ProtocolTest provides the API operation methods for making requests to
// OutputService5ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService5ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice5protocoltest.New(myConfig)
func NewOutputService5ProtocolTest(config aws.Config) *OutputService5ProtocolTest {
	svc := &OutputService5ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService5ProtocolTest",
				ServiceID:     "OutputService5ProtocolTest",
				EndpointsID:   "outputservice5protocoltest",
				SigningName:   "outputservice5protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService5ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService5TestShapeOutputService5TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService5TestShapeOutputService5TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	ListMember []string `locationNameList:"item" type:"list"`
}

const opOutputService5TestCaseOperation1 = "OperationName"

// OutputService5TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService5TestCaseOperation1Request.
//    req := client.OutputService5TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService5ProtocolTest) OutputService5TestCaseOperation1Request(input *OutputService5TestShapeOutputService5TestCaseOperation1Input) OutputService5TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService5TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService5TestShapeOutputService5TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService5TestShapeOutputService5TestCaseOperation1Output{})

	return OutputService5TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService5TestCaseOperation1Request}
}

// OutputService5TestCaseOperation1Request is the request type for the
// OutputService5TestCaseOperation1 API operation.
type OutputService5TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService5TestShapeOutputService5TestCaseOperation1Input
	Copy  func(*OutputService5TestShapeOutputService5TestCaseOperation1Input) OutputService5TestCaseOperation1Request
}

// Send marshals and sends the OutputService5TestCaseOperation1 API request.
func (r OutputService5TestCaseOperation1Request) Send(ctx context.Context) (*OutputService5TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService5TestCaseOperation1Response{
		OutputService5TestShapeOutputService5TestCaseOperation1Output: r.Request.Data.(*OutputService5TestShapeOutputService5TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService5TestCaseOperation1Response is the response type for the
// OutputService5TestCaseOperation1 API operation.
type OutputService5TestCaseOperation1Response struct {
	*OutputService5TestShapeOutputService5TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService5TestCaseOperation1 request.
func (r *OutputService5TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService6ProtocolTest provides the API operation methods for making requests to
// OutputService6ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService6ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice6protocoltest.New(myConfig)
func NewOutputService6ProtocolTest(config aws.Config) *OutputService6ProtocolTest {
	svc := &OutputService6ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService6ProtocolTest",
				ServiceID:     "OutputService6ProtocolTest",
				EndpointsID:   "outputservice6protocoltest",
				SigningName:   "outputservice6protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService6ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService6TestShapeOutputService6TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService6TestShapeOutputService6TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	ListMember []string `type:"list" flattened:"true"`
}

const opOutputService6TestCaseOperation1 = "OperationName"

// OutputService6TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService6TestCaseOperation1Request.
//    req := client.OutputService6TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService6ProtocolTest) OutputService6TestCaseOperation1Request(input *OutputService6TestShapeOutputService6TestCaseOperation1Input) OutputService6TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService6TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService6TestShapeOutputService6TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService6TestShapeOutputService6TestCaseOperation1Output{})

	return OutputService6TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService6TestCaseOperation1Request}
}

// OutputService6TestCaseOperation1Request is the request type for the
// OutputService6TestCaseOperation1 API operation.
type OutputService6TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService6TestShapeOutputService6TestCaseOperation1Input
	Copy  func(*OutputService6TestShapeOutputService6TestCaseOperation1Input) OutputService6TestCaseOperation1Request
}

// Send marshals and sends the OutputService6TestCaseOperation1 API request.
func (r OutputService6TestCaseOperation1Request) Send(ctx context.Context) (*OutputService6TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService6TestCaseOperation1Response{
		OutputService6TestShapeOutputService6TestCaseOperation1Output: r.Request.Data.(*OutputService6TestShapeOutputService6TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService6TestCaseOperation1Response is the response type for the
// OutputService6TestCaseOperation1 API operation.
type OutputService6TestCaseOperation1Response struct {
	*OutputService6TestShapeOutputService6TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService6TestCaseOperation1 request.
func (r *OutputService6TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService7ProtocolTest provides the API operation methods for making requests to
// OutputService7ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService7ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice7protocoltest.New(myConfig)
func NewOutputService7ProtocolTest(config aws.Config) *OutputService7ProtocolTest {
	svc := &OutputService7ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService7ProtocolTest",
				ServiceID:     "OutputService7ProtocolTest",
				EndpointsID:   "outputservice7protocoltest",
				SigningName:   "outputservice7protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService7ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService7TestShapeOutputService7TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService7TestShapeOutputService7TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	ListMember []string `type:"list" flattened:"true"`
}

const opOutputService7TestCaseOperation1 = "OperationName"

// OutputService7TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService7TestCaseOperation1Request.
//    req := client.OutputService7TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService7ProtocolTest) OutputService7TestCaseOperation1Request(input *OutputService7TestShapeOutputService7TestCaseOperation1Input) OutputService7TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService7TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService7TestShapeOutputService7TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService7TestShapeOutputService7TestCaseOperation1Output{})

	return OutputService7TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService7TestCaseOperation1Request}
}

// OutputService7TestCaseOperation1Request is the request type for the
// OutputService7TestCaseOperation1 API operation.
type OutputService7TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService7TestShapeOutputService7TestCaseOperation1Input
	Copy  func(*OutputService7TestShapeOutputService7TestCaseOperation1Input) OutputService7TestCaseOperation1Request
}

// Send marshals and sends the OutputService7TestCaseOperation1 API request.
func (r OutputService7TestCaseOperation1Request) Send(ctx context.Context) (*OutputService7TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService7TestCaseOperation1Response{
		OutputService7TestShapeOutputService7TestCaseOperation1Output: r.Request.Data.(*OutputService7TestShapeOutputService7TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService7TestCaseOperation1Response is the response type for the
// OutputService7TestCaseOperation1 API operation.
type OutputService7TestCaseOperation1Response struct {
	*OutputService7TestShapeOutputService7TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService7TestCaseOperation1 request.
func (r *OutputService7TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService8ProtocolTest provides the API operation methods for making requests to
// OutputService8ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService8ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice8protocoltest.New(myConfig)
func NewOutputService8ProtocolTest(config aws.Config) *OutputService8ProtocolTest {
	svc := &OutputService8ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService8ProtocolTest",
				ServiceID:     "OutputService8ProtocolTest",
				EndpointsID:   "outputservice8protocoltest",
				SigningName:   "outputservice8protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService8ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService8TestShapeOutputService8TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService8TestShapeOutputService8TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	List []OutputService8TestShapeStructureShape `type:"list"`
}

const opOutputService8TestCaseOperation1 = "OperationName"

// OutputService8TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService8TestCaseOperation1Request.
//    req := client.OutputService8TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService8ProtocolTest) OutputService8TestCaseOperation1Request(input *OutputService8TestShapeOutputService8TestCaseOperation1Input) OutputService8TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService8TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService8TestShapeOutputService8TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService8TestShapeOutputService8TestCaseOperation1Output{})

	return OutputService8TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService8TestCaseOperation1Request}
}

// OutputService8TestCaseOperation1Request is the request type for the
// OutputService8TestCaseOperation1 API operation.
type OutputService8TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService8TestShapeOutputService8TestCaseOperation1Input
	Copy  func(*OutputService8TestShapeOutputService8TestCaseOperation1Input) OutputService8TestCaseOperation1Request
}

// Send marshals and sends the OutputService8TestCaseOperation1 API request.
func (r OutputService8TestCaseOperation1Request) Send(ctx context.Context) (*OutputService8TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService8TestCaseOperation1Response{
		OutputService8TestShapeOutputService8TestCaseOperation1Output: r.Request.Data.(*OutputService8TestShapeOutputService8TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService8TestCaseOperation1Response is the response type for the
// OutputService8TestCaseOperation1 API operation.
type OutputService8TestCaseOperation1Response struct {
	*OutputService8TestShapeOutputService8TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService8TestCaseOperation1 request.
func (r *OutputService8TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService8TestShapeStructureShape struct {
	_ struct{} `type:"structure"`

	Bar *string `type:"string"`

	Baz *string `type:"string"`

	Foo *string `type:"string"`
}

// OutputService9ProtocolTest provides the API operation methods for making requests to
// OutputService9ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService9ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice9protocoltest.New(myConfig)
func NewOutputService9ProtocolTest(config aws.Config) *OutputService9ProtocolTest {
	svc := &OutputService9ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService9ProtocolTest",
				ServiceID:     "OutputService9ProtocolTest",
				EndpointsID:   "outputservice9protocoltest",
				SigningName:   "outputservice9protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService9ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService9TestShapeOutputService9TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService9TestShapeOutputService9TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	List []OutputService9TestShapeStructureShape `type:"list" flattened:"true"`
}

const opOutputService9TestCaseOperation1 = "OperationName"

// OutputService9TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService9TestCaseOperation1Request.
//    req := client.OutputService9TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService9ProtocolTest) OutputService9TestCaseOperation1Request(input *OutputService9TestShapeOutputService9TestCaseOperation1Input) OutputService9TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService9TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService9TestShapeOutputService9TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService9TestShapeOutputService9TestCaseOperation1Output{})

	return OutputService9TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService9TestCaseOperation1Request}
}

// OutputService9TestCaseOperation1Request is the request type for the
// OutputService9TestCaseOperation1 API operation.
type OutputService9TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService9TestShapeOutputService9TestCaseOperation1Input
	Copy  func(*OutputService9TestShapeOutputService9TestCaseOperation1Input) OutputService9TestCaseOperation1Request
}

// Send marshals and sends the OutputService9TestCaseOperation1 API request.
func (r OutputService9TestCaseOperation1Request) Send(ctx context.Context) (*OutputService9TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService9TestCaseOperation1Response{
		OutputService9TestShapeOutputService9TestCaseOperation1Output: r.Request.Data.(*OutputService9TestShapeOutputService9TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService9TestCaseOperation1Response is the response type for the
// OutputService9TestCaseOperation1 API operation.
type OutputService9TestCaseOperation1Response struct {
	*OutputService9TestShapeOutputService9TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService9TestCaseOperation1 request.
func (r *OutputService9TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService9TestShapeStructureShape struct {
	_ struct{} `type:"structure"`

	Bar *string `type:"string"`

	Baz *string `type:"string"`

	Foo *string `type:"string"`
}

// OutputService10ProtocolTest provides the API operation methods for making requests to
// OutputService10ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService10ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice10protocoltest.New(myConfig)
func NewOutputService10ProtocolTest(config aws.Config) *OutputService10ProtocolTest {
	svc := &OutputService10ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService10ProtocolTest",
				ServiceID:     "OutputService10ProtocolTest",
				EndpointsID:   "outputservice10protocoltest",
				SigningName:   "outputservice10protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService10ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService10TestShapeOutputService10TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService10TestShapeOutputService10TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	List []string `locationNameList:"NamedList" type:"list" flattened:"true"`
}

const opOutputService10TestCaseOperation1 = "OperationName"

// OutputService10TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService10TestCaseOperation1Request.
//    req := client.OutputService10TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService10ProtocolTest) OutputService10TestCaseOperation1Request(input *OutputService10TestShapeOutputService10TestCaseOperation1Input) OutputService10TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService10TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService10TestShapeOutputService10TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService10TestShapeOutputService10TestCaseOperation1Output{})

	return OutputService10TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService10TestCaseOperation1Request}
}

// OutputService10TestCaseOperation1Request is the request type for the
// OutputService10TestCaseOperation1 API operation.
type OutputService10TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService10TestShapeOutputService10TestCaseOperation1Input
	Copy  func(*OutputService10TestShapeOutputService10TestCaseOperation1Input) OutputService10TestCaseOperation1Request
}

// Send marshals and sends the OutputService10TestCaseOperation1 API request.
func (r OutputService10TestCaseOperation1Request) Send(ctx context.Context) (*OutputService10TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService10TestCaseOperation1Response{
		OutputService10TestShapeOutputService10TestCaseOperation1Output: r.Request.Data.(*OutputService10TestShapeOutputService10TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService10TestCaseOperation1Response is the response type for the
// OutputService10TestCaseOperation1 API operation.
type OutputService10TestCaseOperation1Response struct {
	*OutputService10TestShapeOutputService10TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService10TestCaseOperation1 request.
func (r *OutputService10TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService11ProtocolTest provides the API operation methods for making requests to
// OutputService11ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService11ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice11protocoltest.New(myConfig)
func NewOutputService11ProtocolTest(config aws.Config) *OutputService11ProtocolTest {
	svc := &OutputService11ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService11ProtocolTest",
				ServiceID:     "OutputService11ProtocolTest",
				EndpointsID:   "outputservice11protocoltest",
				SigningName:   "outputservice11protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService11ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService11TestShapeOutputService11TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService11TestShapeOutputService11TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Map map[string]OutputService11TestShapeStructType `type:"map"`
}

const opOutputService11TestCaseOperation1 = "OperationName"

// OutputService11TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService11TestCaseOperation1Request.
//    req := client.OutputService11TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService11ProtocolTest) OutputService11TestCaseOperation1Request(input *OutputService11TestShapeOutputService11TestCaseOperation1Input) OutputService11TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService11TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService11TestShapeOutputService11TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService11TestShapeOutputService11TestCaseOperation1Output{})

	return OutputService11TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService11TestCaseOperation1Request}
}

// OutputService11TestCaseOperation1Request is the request type for the
// OutputService11TestCaseOperation1 API operation.
type OutputService11TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService11TestShapeOutputService11TestCaseOperation1Input
	Copy  func(*OutputService11TestShapeOutputService11TestCaseOperation1Input) OutputService11TestCaseOperation1Request
}

// Send marshals and sends the OutputService11TestCaseOperation1 API request.
func (r OutputService11TestCaseOperation1Request) Send(ctx context.Context) (*OutputService11TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService11TestCaseOperation1Response{
		OutputService11TestShapeOutputService11TestCaseOperation1Output: r.Request.Data.(*OutputService11TestShapeOutputService11TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService11TestCaseOperation1Response is the response type for the
// OutputService11TestCaseOperation1 API operation.
type OutputService11TestCaseOperation1Response struct {
	*OutputService11TestShapeOutputService11TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService11TestCaseOperation1 request.
func (r *OutputService11TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService11TestShapeStructType struct {
	_ struct{} `type:"structure"`

	Foo *string `locationName:"foo" type:"string"`
}

// OutputService12ProtocolTest provides the API operation methods for making requests to
// OutputService12ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService12ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice12protocoltest.New(myConfig)
func NewOutputService12ProtocolTest(config aws.Config) *OutputService12ProtocolTest {
	svc := &OutputService12ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService12ProtocolTest",
				ServiceID:     "OutputService12ProtocolTest",
				EndpointsID:   "outputservice12protocoltest",
				SigningName:   "outputservice12protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService12ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService12TestShapeOutputService12TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService12TestShapeOutputService12TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Map map[string]string `type:"map" flattened:"true"`
}

const opOutputService12TestCaseOperation1 = "OperationName"

// OutputService12TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService12TestCaseOperation1Request.
//    req := client.OutputService12TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService12ProtocolTest) OutputService12TestCaseOperation1Request(input *OutputService12TestShapeOutputService12TestCaseOperation1Input) OutputService12TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService12TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService12TestShapeOutputService12TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService12TestShapeOutputService12TestCaseOperation1Output{})

	return OutputService12TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService12TestCaseOperation1Request}
}

// OutputService12TestCaseOperation1Request is the request type for the
// OutputService12TestCaseOperation1 API operation.
type OutputService12TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService12TestShapeOutputService12TestCaseOperation1Input
	Copy  func(*OutputService12TestShapeOutputService12TestCaseOperation1Input) OutputService12TestCaseOperation1Request
}

// Send marshals and sends the OutputService12TestCaseOperation1 API request.
func (r OutputService12TestCaseOperation1Request) Send(ctx context.Context) (*OutputService12TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService12TestCaseOperation1Response{
		OutputService12TestShapeOutputService12TestCaseOperation1Output: r.Request.Data.(*OutputService12TestShapeOutputService12TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService12TestCaseOperation1Response is the response type for the
// OutputService12TestCaseOperation1 API operation.
type OutputService12TestCaseOperation1Response struct {
	*OutputService12TestShapeOutputService12TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService12TestCaseOperation1 request.
func (r *OutputService12TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService13ProtocolTest provides the API operation methods for making requests to
// OutputService13ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService13ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice13protocoltest.New(myConfig)
func NewOutputService13ProtocolTest(config aws.Config) *OutputService13ProtocolTest {
	svc := &OutputService13ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService13ProtocolTest",
				ServiceID:     "OutputService13ProtocolTest",
				EndpointsID:   "outputservice13protocoltest",
				SigningName:   "outputservice13protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService13ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService13TestShapeOutputService13TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService13TestShapeOutputService13TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Map map[string]string `locationName:"Attribute" locationNameKey:"Name" locationNameValue:"Value" type:"map" flattened:"true"`
}

const opOutputService13TestCaseOperation1 = "OperationName"

// OutputService13TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService13TestCaseOperation1Request.
//    req := client.OutputService13TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService13ProtocolTest) OutputService13TestCaseOperation1Request(input *OutputService13TestShapeOutputService13TestCaseOperation1Input) OutputService13TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService13TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService13TestShapeOutputService13TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService13TestShapeOutputService13TestCaseOperation1Output{})

	return OutputService13TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService13TestCaseOperation1Request}
}

// OutputService13TestCaseOperation1Request is the request type for the
// OutputService13TestCaseOperation1 API operation.
type OutputService13TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService13TestShapeOutputService13TestCaseOperation1Input
	Copy  func(*OutputService13TestShapeOutputService13TestCaseOperation1Input) OutputService13TestCaseOperation1Request
}

// Send marshals and sends the OutputService13TestCaseOperation1 API request.
func (r OutputService13TestCaseOperation1Request) Send(ctx context.Context) (*OutputService13TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService13TestCaseOperation1Response{
		OutputService13TestShapeOutputService13TestCaseOperation1Output: r.Request.Data.(*OutputService13TestShapeOutputService13TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService13TestCaseOperation1Response is the response type for the
// OutputService13TestCaseOperation1 API operation.
type OutputService13TestCaseOperation1Response struct {
	*OutputService13TestShapeOutputService13TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService13TestCaseOperation1 request.
func (r *OutputService13TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService14ProtocolTest provides the API operation methods for making requests to
// OutputService14ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService14ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice14protocoltest.New(myConfig)
func NewOutputService14ProtocolTest(config aws.Config) *OutputService14ProtocolTest {
	svc := &OutputService14ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService14ProtocolTest",
				ServiceID:     "OutputService14ProtocolTest",
				EndpointsID:   "outputservice14protocoltest",
				SigningName:   "outputservice14protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService14ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService14TestShapeOutputService14TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService14TestShapeOutputService14TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Map map[string]string `locationNameKey:"foo" locationNameValue:"bar" type:"map" flattened:"true"`
}

const opOutputService14TestCaseOperation1 = "OperationName"

// OutputService14TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService14TestCaseOperation1Request.
//    req := client.OutputService14TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService14ProtocolTest) OutputService14TestCaseOperation1Request(input *OutputService14TestShapeOutputService14TestCaseOperation1Input) OutputService14TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService14TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService14TestShapeOutputService14TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService14TestShapeOutputService14TestCaseOperation1Output{})

	return OutputService14TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService14TestCaseOperation1Request}
}

// OutputService14TestCaseOperation1Request is the request type for the
// OutputService14TestCaseOperation1 API operation.
type OutputService14TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService14TestShapeOutputService14TestCaseOperation1Input
	Copy  func(*OutputService14TestShapeOutputService14TestCaseOperation1Input) OutputService14TestCaseOperation1Request
}

// Send marshals and sends the OutputService14TestCaseOperation1 API request.
func (r OutputService14TestCaseOperation1Request) Send(ctx context.Context) (*OutputService14TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService14TestCaseOperation1Response{
		OutputService14TestShapeOutputService14TestCaseOperation1Output: r.Request.Data.(*OutputService14TestShapeOutputService14TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService14TestCaseOperation1Response is the response type for the
// OutputService14TestCaseOperation1 API operation.
type OutputService14TestCaseOperation1Response struct {
	*OutputService14TestShapeOutputService14TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService14TestCaseOperation1 request.
func (r *OutputService14TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService15ProtocolTest provides the API operation methods for making requests to
// OutputService15ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService15ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice15protocoltest.New(myConfig)
func NewOutputService15ProtocolTest(config aws.Config) *OutputService15ProtocolTest {
	svc := &OutputService15ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService15ProtocolTest",
				ServiceID:     "OutputService15ProtocolTest",
				EndpointsID:   "outputservice15protocoltest",
				SigningName:   "outputservice15protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService15ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService15TestShapeOutputService15TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService15TestShapeOutputService15TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Foo *string `type:"string"`
}

const opOutputService15TestCaseOperation1 = "OperationName"

// OutputService15TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService15TestCaseOperation1Request.
//    req := client.OutputService15TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService15ProtocolTest) OutputService15TestCaseOperation1Request(input *OutputService15TestShapeOutputService15TestCaseOperation1Input) OutputService15TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService15TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService15TestShapeOutputService15TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService15TestShapeOutputService15TestCaseOperation1Output{})

	return OutputService15TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService15TestCaseOperation1Request}
}

// OutputService15TestCaseOperation1Request is the request type for the
// OutputService15TestCaseOperation1 API operation.
type OutputService15TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService15TestShapeOutputService15TestCaseOperation1Input
	Copy  func(*OutputService15TestShapeOutputService15TestCaseOperation1Input) OutputService15TestCaseOperation1Request
}

// Send marshals and sends the OutputService15TestCaseOperation1 API request.
func (r OutputService15TestCaseOperation1Request) Send(ctx context.Context) (*OutputService15TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService15TestCaseOperation1Response{
		OutputService15TestShapeOutputService15TestCaseOperation1Output: r.Request.Data.(*OutputService15TestShapeOutputService15TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService15TestCaseOperation1Response is the response type for the
// OutputService15TestCaseOperation1 API operation.
type OutputService15TestCaseOperation1Response struct {
	*OutputService15TestShapeOutputService15TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService15TestCaseOperation1 request.
func (r *OutputService15TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// OutputService16ProtocolTest provides the API operation methods for making requests to
// OutputService16ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService16ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice16protocoltest.New(myConfig)
func NewOutputService16ProtocolTest(config aws.Config) *OutputService16ProtocolTest {
	svc := &OutputService16ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService16ProtocolTest",
				ServiceID:     "OutputService16ProtocolTest",
				EndpointsID:   "outputservice16protocoltest",
				SigningName:   "outputservice16protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService16ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService16TestShapeOutputService16TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService16TestShapeOutputService16TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	StructMember *OutputService16TestShapeTimeContainer `type:"structure"`

	TimeArg *time.Time `type:"timestamp"`

	TimeCustom *time.Time `type:"timestamp" timestampFormat:"rfc822"`

	TimeFormat *time.Time `type:"timestamp" timestampFormat:"unixTimestamp"`
}

const opOutputService16TestCaseOperation1 = "OperationName"

// OutputService16TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService16TestCaseOperation1Request.
//    req := client.OutputService16TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService16ProtocolTest) OutputService16TestCaseOperation1Request(input *OutputService16TestShapeOutputService16TestCaseOperation1Input) OutputService16TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService16TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService16TestShapeOutputService16TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService16TestShapeOutputService16TestCaseOperation1Output{})

	return OutputService16TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService16TestCaseOperation1Request}
}

// OutputService16TestCaseOperation1Request is the request type for the
// OutputService16TestCaseOperation1 API operation.
type OutputService16TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService16TestShapeOutputService16TestCaseOperation1Input
	Copy  func(*OutputService16TestShapeOutputService16TestCaseOperation1Input) OutputService16TestCaseOperation1Request
}

// Send marshals and sends the OutputService16TestCaseOperation1 API request.
func (r OutputService16TestCaseOperation1Request) Send(ctx context.Context) (*OutputService16TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService16TestCaseOperation1Response{
		OutputService16TestShapeOutputService16TestCaseOperation1Output: r.Request.Data.(*OutputService16TestShapeOutputService16TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService16TestCaseOperation1Response is the response type for the
// OutputService16TestCaseOperation1 API operation.
type OutputService16TestCaseOperation1Response struct {
	*OutputService16TestShapeOutputService16TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService16TestCaseOperation1 request.
func (r *OutputService16TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService16TestShapeTimeContainer struct {
	_ struct{} `type:"structure"`

	Bar *time.Time `locationName:"bar" type:"timestamp" timestampFormat:"unixTimestamp"`

	Foo *time.Time `locationName:"foo" type:"timestamp"`
}

// OutputService17ProtocolTest provides the API operation methods for making requests to
// OutputService17ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService17ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := outputservice17protocoltest.New(myConfig)
func NewOutputService17ProtocolTest(config aws.Config) *OutputService17ProtocolTest {
	svc := &OutputService17ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "OutputService17ProtocolTest",
				ServiceID:     "OutputService17ProtocolTest",
				EndpointsID:   "outputservice17protocoltest",
				SigningName:   "outputservice17protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService17ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type OutputService17TestShapeOutputService17TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService17TestShapeOutputService17TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	FooEnum OutputService17TestShapeEC2EnumType `type:"string" enum:"true"`

	ListEnums []OutputService17TestShapeEC2EnumType `type:"list"`
}

const opOutputService17TestCaseOperation1 = "OperationName"

// OutputService17TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService17TestCaseOperation1Request.
//    req := client.OutputService17TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService17ProtocolTest) OutputService17TestCaseOperation1Request(input *OutputService17TestShapeOutputService17TestCaseOperation1Input) OutputService17TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService17TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService17TestShapeOutputService17TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &OutputService17TestShapeOutputService17TestCaseOperation1Output{})

	return OutputService17TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService17TestCaseOperation1Request}
}

// OutputService17TestCaseOperation1Request is the request type for the
// OutputService17TestCaseOperation1 API operation.
type OutputService17TestCaseOperation1Request struct {
	*aws.Request
	Input *OutputService17TestShapeOutputService17TestCaseOperation1Input
	Copy  func(*OutputService17TestShapeOutputService17TestCaseOperation1Input) OutputService17TestCaseOperation1Request
}

// Send marshals and sends the OutputService17TestCaseOperation1 API request.
func (r OutputService17TestCaseOperation1Request) Send(ctx context.Context) (*OutputService17TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService17TestCaseOperation1Response{
		OutputService17TestShapeOutputService17TestCaseOperation1Output: r.Request.Data.(*OutputService17TestShapeOutputService17TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService17TestCaseOperation1Response is the response type for the
// OutputService17TestCaseOperation1 API operation.
type OutputService17TestCaseOperation1Response struct {
	*OutputService17TestShapeOutputService17TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService17TestCaseOperation1 request.
func (r *OutputService17TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService17TestShapeEC2EnumType string

// Enum values for OutputService17TestShapeEC2EnumType
const (
	EC2EnumTypeFoo OutputService17TestShapeEC2EnumType = "foo"
	EC2EnumTypeBar OutputService17TestShapeEC2EnumType = "bar"
)

func (enum OutputService17TestShapeEC2EnumType) MarshalValue() (string, error) {
	return string(enum), nil
}

func (enum OutputService17TestShapeEC2EnumType) MarshalValueBuf(b []byte) ([]byte, error) {
	b = b[0:0]
	return append(b, enum...), nil
}

//
// Tests begin here
//

func TestOutputService1ProtocolTestScalarMembersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService1ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><Str>myname</Str><FooNum>123</FooNum><FalseBool>false</FalseBool><TrueBool>true</TrueBool><Float>1.2</Float><Double>1.3</Double><Long>200</Long><Char>a</Char><Timestamp>2015-01-25T08:00:00Z</Timestamp></OperationNameResult><ResponseMetadata><RequestId>request-id</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService1TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService1TestShapeOutputService1TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "a", *out.Char; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 1.3, *out.Double; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := false, *out.FalseBool; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 1.2, *out.Float; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := int64(200), *out.Long; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := int64(123), *out.Num; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("myname"), *out.Str; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.4221728e+09, 0).UTC().String(), out.Timestamp.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := true, *out.TrueBool; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService2ProtocolTestNotAllMembersInResponseCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService2ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><Str>myname</Str></OperationNameResult><ResponseMetadata><RequestId>request-id</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService2TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService2TestShapeOutputService2TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("myname"), *out.Str; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService3ProtocolTestBlobCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService3ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><Blob>dmFsdWU=</Blob></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService3TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService3TestShapeOutputService3TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "value", string(out.Blob); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService4ProtocolTestListsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService4ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><ListMember><member>abc</member><member>123</member></ListMember></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService4TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService4TestShapeOutputService4TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("abc"), out.ListMember[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("123"), out.ListMember[1]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService5ProtocolTestListWithCustomMemberNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService5ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><ListMember><item>abc</item><item>123</item></ListMember></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService5TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService5TestShapeOutputService5TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("abc"), out.ListMember[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("123"), out.ListMember[1]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService6ProtocolTestFlattenedListCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService6ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><ListMember>abc</ListMember><ListMember>123</ListMember></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService6TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService6TestShapeOutputService6TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("abc"), out.ListMember[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("123"), out.ListMember[1]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService7ProtocolTestFlattenedSingleElementListCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService7ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><ListMember>abc</ListMember></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService7TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService7TestShapeOutputService7TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("abc"), out.ListMember[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService8ProtocolTestListOfStructuresCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService8ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse xmlns=\"https://service.amazonaws.com/doc/2010-05-08/\"><OperationNameResult><List><member><Foo>firstfoo</Foo><Bar>firstbar</Bar><Baz>firstbaz</Baz></member><member><Foo>secondfoo</Foo><Bar>secondbar</Bar><Baz>secondbaz</Baz></member></List></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService8TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService8TestShapeOutputService8TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("firstbar"), *out.List[0].Bar; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("firstbaz"), *out.List[0].Baz; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("firstfoo"), *out.List[0].Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("secondbar"), *out.List[1].Bar; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("secondbaz"), *out.List[1].Baz; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("secondfoo"), *out.List[1].Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService9ProtocolTestFlattenedListOfStructuresCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService9ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse xmlns=\"https://service.amazonaws.com/doc/2010-05-08/\"><OperationNameResult><List><Foo>firstfoo</Foo><Bar>firstbar</Bar><Baz>firstbaz</Baz></List><List><Foo>secondfoo</Foo><Bar>secondbar</Bar><Baz>secondbaz</Baz></List></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService9TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService9TestShapeOutputService9TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("firstbar"), *out.List[0].Bar; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("firstbaz"), *out.List[0].Baz; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("firstfoo"), *out.List[0].Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("secondbar"), *out.List[1].Bar; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("secondbaz"), *out.List[1].Baz; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("secondfoo"), *out.List[1].Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService10ProtocolTestFlattenedListWithLocationNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService10ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse xmlns=\"https://service.amazonaws.com/doc/2010-05-08/\"><OperationNameResult><NamedList>a</NamedList><NamedList>b</NamedList></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService10TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService10TestShapeOutputService10TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("a"), out.List[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("b"), out.List[1]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService11ProtocolTestNormalMapCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService11ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse xmlns=\"https://service.amazonaws.com/doc/2010-05-08\"><OperationNameResult><Map><entry><key>qux</key><value><foo>bar</foo></value></entry><entry><key>baz</key><value><foo>bam</foo></value></entry></Map></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService11TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService11TestShapeOutputService11TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("bam"), *out.Map["baz"].Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("bar"), *out.Map["qux"].Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService12ProtocolTestFlattenedMapCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService12ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><Map><key>qux</key><value>bar</value></Map><Map><key>baz</key><value>bam</value></Map></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService12TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService12TestShapeOutputService12TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("bam"), out.Map["baz"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("bar"), out.Map["qux"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService13ProtocolTestFlattenedMapInShapeDefinitionCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService13ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><Attribute><Name>qux</Name><Value>bar</Value></Attribute></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService13TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService13TestShapeOutputService13TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("bar"), out.Map["qux"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService14ProtocolTestNamedMapCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService14ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><Map><foo>qux</foo><bar>bar</bar></Map><Map><foo>baz</foo><bar>bam</bar></Map></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService14TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService14TestShapeOutputService14TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("bam"), out.Map["baz"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("bar"), out.Map["qux"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService15ProtocolTestEmptyStringCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService15ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><OperationNameResult><Foo/></OperationNameResult><ResponseMetadata><RequestId>requestid</RequestId></ResponseMetadata></OperationNameResponse>"))
	req := svc.OutputService15TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService15TestShapeOutputService15TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string(""), *out.Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService16ProtocolTestTimestampMembersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService16ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><StructMember><foo>2014-04-29T18:30:38Z</foo><bar>1398796238</bar></StructMember><TimeArg>2014-04-29T18:30:38Z</TimeArg><TimeCustom>Tue, 29 Apr 2014 18:30:38 GMT</TimeCustom><TimeFormat>1398796238</TimeFormat><RequestId>requestid</RequestId></OperationNameResponse>"))
	req := svc.OutputService16TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService16TestShapeOutputService16TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.StructMember.Bar.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.StructMember.Foo.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.TimeArg.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.TimeCustom.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.TimeFormat.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService17ProtocolTestEnumOutputCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService17ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><FooEnum>foo</FooEnum><ListEnums><member>foo</member><member>bar</member></ListEnums></OperationNameResponse>"))
	req := svc.OutputService17TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	query.UnmarshalMeta(req.Request)
	query.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService17TestShapeOutputService17TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := OutputService17TestShapeEC2EnumType("foo"), out.FooEnum; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := OutputService17TestShapeEC2EnumType("foo"), out.ListEnums[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := OutputService17TestShapeEC2EnumType("bar"), out.ListEnums[1]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}
