package restxml_test

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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restxml"
	"github.com/aws/aws-sdk-go-v2/private/protocol/xml/xmlutil"
	"github.com/aws/aws-sdk-go-v2/private/util"
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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService1ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation1Request(input *types.OutputService1TestShapeOutputService1TestCaseOperation1Input) OutputService1TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService1TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService1TestShapeOutputService1TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService1TestShapeOutputService1TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService1TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService1TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService1TestCaseOperation1Request}
}

// OutputService1TestCaseOperation1Request is the request type for the
// OutputService1TestCaseOperation1 API operation.
type OutputService1TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService1TestShapeOutputService1TestCaseOperation1Input
	Copy  func(*types.OutputService1TestShapeOutputService1TestCaseOperation1Input) OutputService1TestCaseOperation1Request
}

// Send marshals and sends the OutputService1TestCaseOperation1 API request.
func (r OutputService1TestCaseOperation1Request) Send(ctx context.Context) (*OutputService1TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService1TestCaseOperation1Response{
		OutputService1TestShapeOutputService1TestCaseOperation1Output: r.Request.Data.(*types.OutputService1TestShapeOutputService1TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService1TestCaseOperation1Response is the response type for the
// OutputService1TestCaseOperation1 API operation.
type OutputService1TestCaseOperation1Response struct {
	*types.OutputService1TestShapeOutputService1TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService1TestCaseOperation1 request.
func (r *OutputService1TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

const opOutputService1TestCaseOperation2 = "OperationName"

// OutputService1TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService1TestCaseOperation2Request.
//    req := client.OutputService1TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation2Request(input *types.OutputService1TestShapeOutputService1TestCaseOperation2Input) OutputService1TestCaseOperation2Request {
	op := &aws.Operation{
		Name: opOutputService1TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService1TestShapeOutputService1TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &types.OutputService1TestShapeOutputService1TestCaseOperation2Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService1TestCaseOperation2Marshaler(input).MarshalOperation)
	return OutputService1TestCaseOperation2Request{Request: req, Input: input, Copy: c.OutputService1TestCaseOperation2Request}
}

// OutputService1TestCaseOperation2Request is the request type for the
// OutputService1TestCaseOperation2 API operation.
type OutputService1TestCaseOperation2Request struct {
	*aws.Request
	Input *types.OutputService1TestShapeOutputService1TestCaseOperation2Input
	Copy  func(*types.OutputService1TestShapeOutputService1TestCaseOperation2Input) OutputService1TestCaseOperation2Request
}

// Send marshals and sends the OutputService1TestCaseOperation2 API request.
func (r OutputService1TestCaseOperation2Request) Send(ctx context.Context) (*OutputService1TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService1TestCaseOperation2Response{
		OutputService1TestShapeOutputService1TestCaseOperation2Output: r.Request.Data.(*types.OutputService1TestShapeOutputService1TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService1TestCaseOperation2Response is the response type for the
// OutputService1TestCaseOperation2 API operation.
type OutputService1TestCaseOperation2Response struct {
	*types.OutputService1TestShapeOutputService1TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService1TestCaseOperation2 request.
func (r *OutputService1TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService2ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService2ProtocolTest) OutputService2TestCaseOperation1Request(input *types.OutputService2TestShapeOutputService2TestCaseOperation1Input) OutputService2TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService2TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService2TestShapeOutputService2TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService2TestShapeOutputService2TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService2TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService2TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService2TestCaseOperation1Request}
}

// OutputService2TestCaseOperation1Request is the request type for the
// OutputService2TestCaseOperation1 API operation.
type OutputService2TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService2TestShapeOutputService2TestCaseOperation1Input
	Copy  func(*types.OutputService2TestShapeOutputService2TestCaseOperation1Input) OutputService2TestCaseOperation1Request
}

// Send marshals and sends the OutputService2TestCaseOperation1 API request.
func (r OutputService2TestCaseOperation1Request) Send(ctx context.Context) (*OutputService2TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService2TestCaseOperation1Response{
		OutputService2TestShapeOutputService2TestCaseOperation1Output: r.Request.Data.(*types.OutputService2TestShapeOutputService2TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService2TestCaseOperation1Response is the response type for the
// OutputService2TestCaseOperation1 API operation.
type OutputService2TestCaseOperation1Response struct {
	*types.OutputService2TestShapeOutputService2TestCaseOperation1Output

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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService3ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService3ProtocolTest) OutputService3TestCaseOperation1Request(input *types.OutputService3TestShapeOutputService3TestCaseOperation1Input) OutputService3TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService3TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService3TestShapeOutputService3TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService3TestShapeOutputService3TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService3TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService3TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService3TestCaseOperation1Request}
}

// OutputService3TestCaseOperation1Request is the request type for the
// OutputService3TestCaseOperation1 API operation.
type OutputService3TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService3TestShapeOutputService3TestCaseOperation1Input
	Copy  func(*types.OutputService3TestShapeOutputService3TestCaseOperation1Input) OutputService3TestCaseOperation1Request
}

// Send marshals and sends the OutputService3TestCaseOperation1 API request.
func (r OutputService3TestCaseOperation1Request) Send(ctx context.Context) (*OutputService3TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService3TestCaseOperation1Response{
		OutputService3TestShapeOutputService3TestCaseOperation1Output: r.Request.Data.(*types.OutputService3TestShapeOutputService3TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService3TestCaseOperation1Response is the response type for the
// OutputService3TestCaseOperation1 API operation.
type OutputService3TestCaseOperation1Response struct {
	*types.OutputService3TestShapeOutputService3TestCaseOperation1Output

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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService4ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService4ProtocolTest) OutputService4TestCaseOperation1Request(input *types.OutputService4TestShapeOutputService4TestCaseOperation1Input) OutputService4TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService4TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService4TestShapeOutputService4TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService4TestShapeOutputService4TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService4TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService4TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService4TestCaseOperation1Request}
}

// OutputService4TestCaseOperation1Request is the request type for the
// OutputService4TestCaseOperation1 API operation.
type OutputService4TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService4TestShapeOutputService4TestCaseOperation1Input
	Copy  func(*types.OutputService4TestShapeOutputService4TestCaseOperation1Input) OutputService4TestCaseOperation1Request
}

// Send marshals and sends the OutputService4TestCaseOperation1 API request.
func (r OutputService4TestCaseOperation1Request) Send(ctx context.Context) (*OutputService4TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService4TestCaseOperation1Response{
		OutputService4TestShapeOutputService4TestCaseOperation1Output: r.Request.Data.(*types.OutputService4TestShapeOutputService4TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService4TestCaseOperation1Response is the response type for the
// OutputService4TestCaseOperation1 API operation.
type OutputService4TestCaseOperation1Response struct {
	*types.OutputService4TestShapeOutputService4TestCaseOperation1Output

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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService5ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService5ProtocolTest) OutputService5TestCaseOperation1Request(input *types.OutputService5TestShapeOutputService5TestCaseOperation1Input) OutputService5TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService5TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService5TestShapeOutputService5TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService5TestShapeOutputService5TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService5TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService5TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService5TestCaseOperation1Request}
}

// OutputService5TestCaseOperation1Request is the request type for the
// OutputService5TestCaseOperation1 API operation.
type OutputService5TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService5TestShapeOutputService5TestCaseOperation1Input
	Copy  func(*types.OutputService5TestShapeOutputService5TestCaseOperation1Input) OutputService5TestCaseOperation1Request
}

// Send marshals and sends the OutputService5TestCaseOperation1 API request.
func (r OutputService5TestCaseOperation1Request) Send(ctx context.Context) (*OutputService5TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService5TestCaseOperation1Response{
		OutputService5TestShapeOutputService5TestCaseOperation1Output: r.Request.Data.(*types.OutputService5TestShapeOutputService5TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService5TestCaseOperation1Response is the response type for the
// OutputService5TestCaseOperation1 API operation.
type OutputService5TestCaseOperation1Response struct {
	*types.OutputService5TestShapeOutputService5TestCaseOperation1Output

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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService6ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService6ProtocolTest) OutputService6TestCaseOperation1Request(input *types.OutputService6TestShapeOutputService6TestCaseOperation1Input) OutputService6TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService6TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService6TestShapeOutputService6TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService6TestShapeOutputService6TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService6TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService6TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService6TestCaseOperation1Request}
}

// OutputService6TestCaseOperation1Request is the request type for the
// OutputService6TestCaseOperation1 API operation.
type OutputService6TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService6TestShapeOutputService6TestCaseOperation1Input
	Copy  func(*types.OutputService6TestShapeOutputService6TestCaseOperation1Input) OutputService6TestCaseOperation1Request
}

// Send marshals and sends the OutputService6TestCaseOperation1 API request.
func (r OutputService6TestCaseOperation1Request) Send(ctx context.Context) (*OutputService6TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService6TestCaseOperation1Response{
		OutputService6TestShapeOutputService6TestCaseOperation1Output: r.Request.Data.(*types.OutputService6TestShapeOutputService6TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService6TestCaseOperation1Response is the response type for the
// OutputService6TestCaseOperation1 API operation.
type OutputService6TestCaseOperation1Response struct {
	*types.OutputService6TestShapeOutputService6TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService6TestCaseOperation1 request.
func (r *OutputService6TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService6TestShapeSingleStructure struct {
	_ struct{} `type:"structure"`

	Foo *string `locationName:"foo" type:"string"`
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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService7ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService7ProtocolTest) OutputService7TestCaseOperation1Request(input *types.OutputService7TestShapeOutputService7TestCaseOperation1Input) OutputService7TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService7TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService7TestShapeOutputService7TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService7TestShapeOutputService7TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService7TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService7TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService7TestCaseOperation1Request}
}

// OutputService7TestCaseOperation1Request is the request type for the
// OutputService7TestCaseOperation1 API operation.
type OutputService7TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService7TestShapeOutputService7TestCaseOperation1Input
	Copy  func(*types.OutputService7TestShapeOutputService7TestCaseOperation1Input) OutputService7TestCaseOperation1Request
}

// Send marshals and sends the OutputService7TestCaseOperation1 API request.
func (r OutputService7TestCaseOperation1Request) Send(ctx context.Context) (*OutputService7TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService7TestCaseOperation1Response{
		OutputService7TestShapeOutputService7TestCaseOperation1Output: r.Request.Data.(*types.OutputService7TestShapeOutputService7TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService7TestCaseOperation1Response is the response type for the
// OutputService7TestCaseOperation1 API operation.
type OutputService7TestCaseOperation1Response struct {
	*types.OutputService7TestShapeOutputService7TestCaseOperation1Output

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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService8ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService8ProtocolTest) OutputService8TestCaseOperation1Request(input *types.OutputService8TestShapeOutputService8TestCaseOperation1Input) OutputService8TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService8TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService8TestShapeOutputService8TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService8TestShapeOutputService8TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService8TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService8TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService8TestCaseOperation1Request}
}

// OutputService8TestCaseOperation1Request is the request type for the
// OutputService8TestCaseOperation1 API operation.
type OutputService8TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService8TestShapeOutputService8TestCaseOperation1Input
	Copy  func(*types.OutputService8TestShapeOutputService8TestCaseOperation1Input) OutputService8TestCaseOperation1Request
}

// Send marshals and sends the OutputService8TestCaseOperation1 API request.
func (r OutputService8TestCaseOperation1Request) Send(ctx context.Context) (*OutputService8TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService8TestCaseOperation1Response{
		OutputService8TestShapeOutputService8TestCaseOperation1Output: r.Request.Data.(*types.OutputService8TestShapeOutputService8TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService8TestCaseOperation1Response is the response type for the
// OutputService8TestCaseOperation1 API operation.
type OutputService8TestCaseOperation1Response struct {
	*types.OutputService8TestShapeOutputService8TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService8TestCaseOperation1 request.
func (r *OutputService8TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService9ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService9ProtocolTest) OutputService9TestCaseOperation1Request(input *types.OutputService9TestShapeOutputService9TestCaseOperation1Input) OutputService9TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService9TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService9TestShapeOutputService9TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService9TestShapeOutputService9TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService9TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService9TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService9TestCaseOperation1Request}
}

// OutputService9TestCaseOperation1Request is the request type for the
// OutputService9TestCaseOperation1 API operation.
type OutputService9TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService9TestShapeOutputService9TestCaseOperation1Input
	Copy  func(*types.OutputService9TestShapeOutputService9TestCaseOperation1Input) OutputService9TestCaseOperation1Request
}

// Send marshals and sends the OutputService9TestCaseOperation1 API request.
func (r OutputService9TestCaseOperation1Request) Send(ctx context.Context) (*OutputService9TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService9TestCaseOperation1Response{
		OutputService9TestShapeOutputService9TestCaseOperation1Output: r.Request.Data.(*types.OutputService9TestShapeOutputService9TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService9TestCaseOperation1Response is the response type for the
// OutputService9TestCaseOperation1 API operation.
type OutputService9TestCaseOperation1Response struct {
	*types.OutputService9TestShapeOutputService9TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService9TestCaseOperation1 request.
func (r *OutputService9TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService9TestShapeSingleStructure struct {
	_ struct{} `type:"structure"`

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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService10ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService10ProtocolTest) OutputService10TestCaseOperation1Request(input *types.OutputService10TestShapeOutputService10TestCaseOperation1Input) OutputService10TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService10TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService10TestShapeOutputService10TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService10TestShapeOutputService10TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService10TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService10TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService10TestCaseOperation1Request}
}

// OutputService10TestCaseOperation1Request is the request type for the
// OutputService10TestCaseOperation1 API operation.
type OutputService10TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService10TestShapeOutputService10TestCaseOperation1Input
	Copy  func(*types.OutputService10TestShapeOutputService10TestCaseOperation1Input) OutputService10TestCaseOperation1Request
}

// Send marshals and sends the OutputService10TestCaseOperation1 API request.
func (r OutputService10TestCaseOperation1Request) Send(ctx context.Context) (*OutputService10TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService10TestCaseOperation1Response{
		OutputService10TestShapeOutputService10TestCaseOperation1Output: r.Request.Data.(*types.OutputService10TestShapeOutputService10TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService10TestCaseOperation1Response is the response type for the
// OutputService10TestCaseOperation1 API operation.
type OutputService10TestCaseOperation1Response struct {
	*types.OutputService10TestShapeOutputService10TestCaseOperation1Output

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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService11ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService11ProtocolTest) OutputService11TestCaseOperation1Request(input *types.OutputService11TestShapeOutputService11TestCaseOperation1Input) OutputService11TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService11TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService11TestShapeOutputService11TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService11TestShapeOutputService11TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService11TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService11TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService11TestCaseOperation1Request}
}

// OutputService11TestCaseOperation1Request is the request type for the
// OutputService11TestCaseOperation1 API operation.
type OutputService11TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService11TestShapeOutputService11TestCaseOperation1Input
	Copy  func(*types.OutputService11TestShapeOutputService11TestCaseOperation1Input) OutputService11TestCaseOperation1Request
}

// Send marshals and sends the OutputService11TestCaseOperation1 API request.
func (r OutputService11TestCaseOperation1Request) Send(ctx context.Context) (*OutputService11TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService11TestCaseOperation1Response{
		OutputService11TestShapeOutputService11TestCaseOperation1Output: r.Request.Data.(*types.OutputService11TestShapeOutputService11TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService11TestCaseOperation1Response is the response type for the
// OutputService11TestCaseOperation1 API operation.
type OutputService11TestCaseOperation1Response struct {
	*types.OutputService11TestShapeOutputService11TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService11TestCaseOperation1 request.
func (r *OutputService11TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService12ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService12ProtocolTest) OutputService12TestCaseOperation1Request(input *types.OutputService12TestShapeOutputService12TestCaseOperation1Input) OutputService12TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService12TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService12TestShapeOutputService12TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService12TestShapeOutputService12TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService12TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService12TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService12TestCaseOperation1Request}
}

// OutputService12TestCaseOperation1Request is the request type for the
// OutputService12TestCaseOperation1 API operation.
type OutputService12TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService12TestShapeOutputService12TestCaseOperation1Input
	Copy  func(*types.OutputService12TestShapeOutputService12TestCaseOperation1Input) OutputService12TestCaseOperation1Request
}

// Send marshals and sends the OutputService12TestCaseOperation1 API request.
func (r OutputService12TestCaseOperation1Request) Send(ctx context.Context) (*OutputService12TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService12TestCaseOperation1Response{
		OutputService12TestShapeOutputService12TestCaseOperation1Output: r.Request.Data.(*types.OutputService12TestShapeOutputService12TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService12TestCaseOperation1Response is the response type for the
// OutputService12TestCaseOperation1 API operation.
type OutputService12TestCaseOperation1Response struct {
	*types.OutputService12TestShapeOutputService12TestCaseOperation1Output

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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService13ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService13ProtocolTest) OutputService13TestCaseOperation1Request(input *types.OutputService13TestShapeOutputService13TestCaseOperation1Input) OutputService13TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opOutputService13TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &types.OutputService13TestShapeOutputService13TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService13TestShapeOutputService13TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService13TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService13TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService13TestCaseOperation1Request}
}

// OutputService13TestCaseOperation1Request is the request type for the
// OutputService13TestCaseOperation1 API operation.
type OutputService13TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService13TestShapeOutputService13TestCaseOperation1Input
	Copy  func(*types.OutputService13TestShapeOutputService13TestCaseOperation1Input) OutputService13TestCaseOperation1Request
}

// Send marshals and sends the OutputService13TestCaseOperation1 API request.
func (r OutputService13TestCaseOperation1Request) Send(ctx context.Context) (*OutputService13TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService13TestCaseOperation1Response{
		OutputService13TestShapeOutputService13TestCaseOperation1Output: r.Request.Data.(*types.OutputService13TestShapeOutputService13TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService13TestCaseOperation1Response is the response type for the
// OutputService13TestCaseOperation1 API operation.
type OutputService13TestCaseOperation1Response struct {
	*types.OutputService13TestShapeOutputService13TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService13TestCaseOperation1 request.
func (r *OutputService13TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService13TestShapeTimeContainer struct {
	_ struct{} `type:"structure"`

	Bar *time.Time `locationName:"bar" type:"timestamp" timestampFormat:"unixTimestamp"`

	Foo *time.Time `locationName:"foo" type:"timestamp"`
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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService14ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService14ProtocolTest) OutputService14TestCaseOperation1Request(input *types.OutputService14TestShapeOutputService14TestCaseOperation1Input) OutputService14TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opOutputService14TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &types.OutputService14TestShapeOutputService14TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService14TestShapeOutputService14TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService14TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService14TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService14TestCaseOperation1Request}
}

// OutputService14TestCaseOperation1Request is the request type for the
// OutputService14TestCaseOperation1 API operation.
type OutputService14TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService14TestShapeOutputService14TestCaseOperation1Input
	Copy  func(*types.OutputService14TestShapeOutputService14TestCaseOperation1Input) OutputService14TestCaseOperation1Request
}

// Send marshals and sends the OutputService14TestCaseOperation1 API request.
func (r OutputService14TestCaseOperation1Request) Send(ctx context.Context) (*OutputService14TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService14TestCaseOperation1Response{
		OutputService14TestShapeOutputService14TestCaseOperation1Output: r.Request.Data.(*types.OutputService14TestShapeOutputService14TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService14TestCaseOperation1Response is the response type for the
// OutputService14TestCaseOperation1 API operation.
type OutputService14TestCaseOperation1Response struct {
	*types.OutputService14TestShapeOutputService14TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService14TestCaseOperation1 request.
func (r *OutputService14TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

const opOutputService14TestCaseOperation2 = "OperationName"

// OutputService14TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using OutputService14TestCaseOperation2Request.
//    req := client.OutputService14TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *OutputService14ProtocolTest) OutputService14TestCaseOperation2Request(input *types.OutputService14TestShapeOutputService14TestCaseOperation2Input) OutputService14TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opOutputService14TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &types.OutputService14TestShapeOutputService14TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &types.OutputService14TestShapeOutputService14TestCaseOperation2Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService14TestCaseOperation2Marshaler(input).MarshalOperation)
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return OutputService14TestCaseOperation2Request{Request: req, Input: input, Copy: c.OutputService14TestCaseOperation2Request}
}

// OutputService14TestCaseOperation2Request is the request type for the
// OutputService14TestCaseOperation2 API operation.
type OutputService14TestCaseOperation2Request struct {
	*aws.Request
	Input *types.OutputService14TestShapeOutputService14TestCaseOperation2Input
	Copy  func(*types.OutputService14TestShapeOutputService14TestCaseOperation2Input) OutputService14TestCaseOperation2Request
}

// Send marshals and sends the OutputService14TestCaseOperation2 API request.
func (r OutputService14TestCaseOperation2Request) Send(ctx context.Context) (*OutputService14TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService14TestCaseOperation2Response{
		OutputService14TestShapeOutputService14TestCaseOperation2Output: r.Request.Data.(*types.OutputService14TestShapeOutputService14TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService14TestCaseOperation2Response is the response type for the
// OutputService14TestCaseOperation2 API operation.
type OutputService14TestCaseOperation2Response struct {
	*types.OutputService14TestShapeOutputService14TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService14TestCaseOperation2 request.
func (r *OutputService14TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService14TestShapeRESTJSONEnumType string

// Enum values for OutputService14TestShapeRESTJSONEnumType
const (
	RESTJSONEnumTypeFoo OutputService14TestShapeRESTJSONEnumType = "foo"
	RESTJSONEnumTypeBar OutputService14TestShapeRESTJSONEnumType = "bar"
	RESTJSONEnumType0   OutputService14TestShapeRESTJSONEnumType = "0"
	RESTJSONEnumType1   OutputService14TestShapeRESTJSONEnumType = "1"
)

func (enum OutputService14TestShapeRESTJSONEnumType) MarshalValue() (string, error) {
	return string(enum), nil
}

func (enum OutputService14TestShapeRESTJSONEnumType) MarshalValueBuf(b []byte) ([]byte, error) {
	b = b[0:0]
	return append(b, enum...), nil
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

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *OutputService15ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
func (c *OutputService15ProtocolTest) OutputService15TestCaseOperation1Request(input *types.OutputService15TestShapeOutputService15TestCaseOperation1Input) OutputService15TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opOutputService15TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &types.OutputService15TestShapeOutputService15TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &types.OutputService15TestShapeOutputService15TestCaseOperation1Output{})
	req.Handlers.Build.Remove(restxml.BuildHandler)
	req.Handlers.Build.PushBack(aws_restxml.NewOutputService15TestCaseOperation1Marshaler(input).MarshalOperation)
	return OutputService15TestCaseOperation1Request{Request: req, Input: input, Copy: c.OutputService15TestCaseOperation1Request}
}

// OutputService15TestCaseOperation1Request is the request type for the
// OutputService15TestCaseOperation1 API operation.
type OutputService15TestCaseOperation1Request struct {
	*aws.Request
	Input *types.OutputService15TestShapeOutputService15TestCaseOperation1Input
	Copy  func(*types.OutputService15TestShapeOutputService15TestCaseOperation1Input) OutputService15TestCaseOperation1Request
}

// Send marshals and sends the OutputService15TestCaseOperation1 API request.
func (r OutputService15TestCaseOperation1Request) Send(ctx context.Context) (*OutputService15TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &OutputService15TestCaseOperation1Response{
		OutputService15TestShapeOutputService15TestCaseOperation1Output: r.Request.Data.(*types.OutputService15TestShapeOutputService15TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// OutputService15TestCaseOperation1Response is the response type for the
// OutputService15TestCaseOperation1 API operation.
type OutputService15TestCaseOperation1Response struct {
	*types.OutputService15TestShapeOutputService15TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// OutputService15TestCaseOperation1 request.
func (r *OutputService15TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type OutputService15TestShapeItemDetailShape struct {
	_ struct{} `type:"structure" xmlPrefix:"xsi" xmlURI:"http://www.w3.org/2001/XMLSchema-instance"`

	ID *string `type:"string"`

	// Type is a required field
	Type enums.OutputService15TestShapeItemType `locationName:"xsi:type" type:"string" xmlAttribute:"true" required:"true" enum:"true"`
}

type OutputService15TestShapeItemShape struct {
	_ struct{} `type:"structure"`

	ItemDetail *OutputService15TestShapeItemDetailShape `type:"structure" xmlPrefix:"xsi" xmlURI:"http://www.w3.org/2001/XMLSchema-instance"`
}

type OutputService15TestShapeItemType string

// Enum values for OutputService15TestShapeItemType
const (
	ItemTypeType1 OutputService15TestShapeItemType = "Type1"
	ItemTypeType2 OutputService15TestShapeItemType = "Type2"
	ItemTypeType3 OutputService15TestShapeItemType = "Type3"
)

func (enum OutputService15TestShapeItemType) MarshalValue() (string, error) {
	return string(enum), nil
}

func (enum OutputService15TestShapeItemType) MarshalValueBuf(b []byte) ([]byte, error) {
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

	buf := bytes.NewReader([]byte("<OperationNameResponse><Str>myname</Str><FooNum>123</FooNum><FalseBool>false</FalseBool><TrueBool>true</TrueBool><Float>1.2</Float><Double>1.3</Double><Long>200</Long><Char>a</Char><Timestamp>2015-01-25T08:00:00Z</Timestamp></OperationNameResponse>"))
	req := svc.OutputService1TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("ImaHeader", "test")
	req.HTTPResponse.Header.Set("X-Foo", "abc")

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
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
	if e, a := string("test"), *out.ImaHeader; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("abc"), *out.ImaHeaderLocation; e != a {
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

func TestOutputService1ProtocolTestScalarMembersCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService1ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Str></Str><FooNum>123</FooNum><FalseBool>false</FalseBool><TrueBool>true</TrueBool><Float>1.2</Float><Double>1.3</Double><Long>200</Long><Char>a</Char><Timestamp>2015-01-25T08:00:00Z</Timestamp></OperationNameResponse>"))
	req := svc.OutputService1TestCaseOperation2Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("ImaHeader", "test")
	req.HTTPResponse.Header.Set("X-Foo", "abc")

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService1TestShapeOutputService1TestCaseOperation2Output)
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
	if e, a := string("test"), *out.ImaHeader; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("abc"), *out.ImaHeaderLocation; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := int64(200), *out.Long; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := int64(123), *out.Num; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string(""), *out.Str; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.4221728e+09, 0).UTC().String(), out.Timestamp.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := true, *out.TrueBool; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService2ProtocolTestBlobCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService2ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResult><Blob>dmFsdWU=</Blob></OperationNameResult>"))
	req := svc.OutputService2TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService2TestShapeOutputService2TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "value", string(out.Blob); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService3ProtocolTestListsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService3ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResult><ListMember><member>abc</member><member>123</member></ListMember></OperationNameResult>"))
	req := svc.OutputService3TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService3TestShapeOutputService3TestCaseOperation1Output)
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

func TestOutputService4ProtocolTestListWithCustomMemberNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService4ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResult><ListMember><item>abc</item><item>123</item></ListMember></OperationNameResult>"))
	req := svc.OutputService4TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
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

func TestOutputService5ProtocolTestFlattenedListCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService5ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResult><ListMember>abc</ListMember><ListMember>123</ListMember></OperationNameResult>"))
	req := svc.OutputService5TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
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

func TestOutputService6ProtocolTestNormalMapCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService6ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResult><Map><entry><key>qux</key><value><foo>bar</foo></value></entry><entry><key>baz</key><value><foo>bam</foo></value></entry></Map></OperationNameResult>"))
	req := svc.OutputService6TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService6TestShapeOutputService6TestCaseOperation1Output)
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

func TestOutputService7ProtocolTestFlattenedMapCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService7ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResult><Map><key>qux</key><value>bar</value></Map><Map><key>baz</key><value>bam</value></Map></OperationNameResult>"))
	req := svc.OutputService7TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService7TestShapeOutputService7TestCaseOperation1Output)
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

func TestOutputService8ProtocolTestNamedMapCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService8ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResult><Map><entry><foo>qux</foo><bar>bar</bar></entry><entry><foo>baz</foo><bar>bam</bar></entry></Map></OperationNameResult>"))
	req := svc.OutputService8TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService8TestShapeOutputService8TestCaseOperation1Output)
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

func TestOutputService9ProtocolTestXMLPayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService9ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Foo>abc</Foo></OperationNameResponse>"))
	req := svc.OutputService9TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("X-Foo", "baz")

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService9TestShapeOutputService9TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("abc"), *out.Data.Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("baz"), *out.Header; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService10ProtocolTestStreamingPayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService10ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("abc"))
	req := svc.OutputService10TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService10TestShapeOutputService10TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "abc", string(out.Stream); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService11ProtocolTestScalarMembersInHeadersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService11ProtocolTest(cfg)

	buf := bytes.NewReader([]byte(""))
	req := svc.OutputService11TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("x-char", "a")
	req.HTTPResponse.Header.Set("x-double", "1.5")
	req.HTTPResponse.Header.Set("x-false-bool", "false")
	req.HTTPResponse.Header.Set("x-float", "1.5")
	req.HTTPResponse.Header.Set("x-int", "1")
	req.HTTPResponse.Header.Set("x-long", "100")
	req.HTTPResponse.Header.Set("x-str", "string")
	req.HTTPResponse.Header.Set("x-timestamp", "Sun, 25 Jan 2015 08:00:00 GMT")
	req.HTTPResponse.Header.Set("x-true-bool", "true")

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService11TestShapeOutputService11TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "a", *out.Char; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 1.5, *out.Double; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := false, *out.FalseBool; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 1.5, *out.Float; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := int64(1), *out.Integer; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := int64(100), *out.Long; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("string"), *out.Str; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.4221728e+09, 0).UTC().String(), out.Timestamp.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := true, *out.TrueBool; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService12ProtocolTestEmptyStringCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService12ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Foo/><RequestId>requestid</RequestId></OperationNameResponse>"))
	req := svc.OutputService12TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService12TestShapeOutputService12TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string(""), *out.Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService13ProtocolTestTimestampMembersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService13ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><StructMember><foo>2014-04-29T18:30:38Z</foo><bar>1398796238</bar></StructMember><TimeArg>2014-04-29T18:30:38Z</TimeArg><TimeCustom>Tue, 29 Apr 2014 18:30:38 GMT</TimeCustom><TimeFormat>1398796238</TimeFormat><RequestId>requestid</RequestId></OperationNameResponse>"))
	req := svc.OutputService13TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("x-amz-timearg", "Tue, 29 Apr 2014 18:30:38 GMT")
	req.HTTPResponse.Header.Set("x-amz-timecustom", "1398796238")
	req.HTTPResponse.Header.Set("x-amz-timeformat", "1398796238")

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService13TestShapeOutputService13TestCaseOperation1Output)
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
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.TimeArgInHeader.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.TimeCustom.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.TimeCustomInHeader.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.TimeFormat.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := time.Unix(1.398796238e+09, 0).UTC().String(), out.TimeFormatInHeader.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService14ProtocolTestEnumCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService14ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><FooEnum>foo</FooEnum><ListEnums><member>0</member><member>1</member></ListEnums></OperationNameResponse>"))
	req := svc.OutputService14TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers
	req.HTTPResponse.Header.Set("x-amz-enum", "baz")

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService14TestShapeOutputService14TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := OutputService14TestShapeRESTJSONEnumType("foo"), out.FooEnum; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := OutputService14TestShapeRESTJSONEnumType("baz"), out.HeaderEnum; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := OutputService14TestShapeRESTJSONEnumType("0"), out.ListEnums[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := OutputService14TestShapeRESTJSONEnumType("1"), out.ListEnums[1]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService14ProtocolTestEnumCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService14ProtocolTest(cfg)

	buf := bytes.NewReader([]byte(""))
	req := svc.OutputService14TestCaseOperation2Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService14TestShapeOutputService14TestCaseOperation2Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}

}

func TestOutputService15ProtocolTestXMLAttributesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService15ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<SomeOutputDoc xmlns=\"http://s3.amazonaws.com/doc/2006-03-01/\"><ItemsList><Item><ItemDetail xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:type=\"Type1\"><ID>id1</ID></ItemDetail></Item><Item><ItemDetail xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:type=\"Type2\"><ID>id2</ID></ItemDetail></Item><Item><ItemDetail xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:type=\"Type3\"><ID>id3</ID></ItemDetail></Item></ItemsList></SomeOutputDoc>"))
	req := svc.OutputService15TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	restxml.UnmarshalMeta(req.Request)
	restxml.Unmarshal(req.Request)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	out := req.Data.(*OutputService15TestShapeOutputService15TestCaseOperation1Output)
	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := string("id1"), *out.ListItems[0].ItemDetail.ID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := OutputService15TestShapeItemType("Type1"), out.ListItems[0].ItemDetail.Type; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("id2"), *out.ListItems[1].ItemDetail.ID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := OutputService15TestShapeItemType("Type2"), out.ListItems[1].ItemDetail.Type; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := string("id3"), *out.ListItems[2].ItemDetail.ID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := OutputService15TestShapeItemType("Type3"), out.ListItems[2].ItemDetail.Type; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}
