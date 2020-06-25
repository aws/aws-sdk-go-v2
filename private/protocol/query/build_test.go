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

// InputService1ProtocolTest provides the API operation methods for making requests to
// InputService1ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService1ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice1protocoltest.New(myConfig)
func NewInputService1ProtocolTest(config aws.Config) *InputService1ProtocolTest {
	svc := &InputService1ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService1ProtocolTest",
				ServiceID:     "InputService1ProtocolTest",
				EndpointsID:   "inputservice1protocoltest",
				SigningName:   "inputservice1protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService1ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService1TestShapeInputService1TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Bar *string `type:"string"`

	Baz *bool `type:"boolean"`

	Foo *string `type:"string"`
}

type InputService1TestShapeInputService1TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService1TestCaseOperation1 = "OperationName"

// InputService1TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService1TestCaseOperation1Request.
//    req := client.InputService1TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService1ProtocolTest) InputService1TestCaseOperation1Request(input *InputService1TestShapeInputService1TestCaseOperation1Input) InputService1TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService1TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService1TestShapeInputService1TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService1TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService1TestCaseOperation1Request}
}

// InputService1TestCaseOperation1Request is the request type for the
// InputService1TestCaseOperation1 API operation.
type InputService1TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService1TestShapeInputService1TestCaseOperation1Input
	Copy  func(*InputService1TestShapeInputService1TestCaseOperation1Input) InputService1TestCaseOperation1Request
}

// Send marshals and sends the InputService1TestCaseOperation1 API request.
func (r InputService1TestCaseOperation1Request) Send(ctx context.Context) (*InputService1TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService1TestCaseOperation1Response{
		InputService1TestShapeInputService1TestCaseOperation1Output: r.Request.Data.(*InputService1TestShapeInputService1TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService1TestCaseOperation1Response is the response type for the
// InputService1TestCaseOperation1 API operation.
type InputService1TestCaseOperation1Response struct {
	*InputService1TestShapeInputService1TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService1TestCaseOperation1 request.
func (r *InputService1TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService1TestShapeInputService1TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	Bar *string `type:"string"`

	Baz *bool `type:"boolean"`

	Foo *string `type:"string"`
}

type InputService1TestShapeInputService1TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService1TestCaseOperation2 = "OperationName"

// InputService1TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService1TestCaseOperation2Request.
//    req := client.InputService1TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService1ProtocolTest) InputService1TestCaseOperation2Request(input *InputService1TestShapeInputService1TestCaseOperation2Input) InputService1TestCaseOperation2Request {
	op := &aws.Operation{
		Name: opInputService1TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService1TestShapeInputService1TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService1TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService1TestCaseOperation2Request}
}

// InputService1TestCaseOperation2Request is the request type for the
// InputService1TestCaseOperation2 API operation.
type InputService1TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService1TestShapeInputService1TestCaseOperation2Input
	Copy  func(*InputService1TestShapeInputService1TestCaseOperation2Input) InputService1TestCaseOperation2Request
}

// Send marshals and sends the InputService1TestCaseOperation2 API request.
func (r InputService1TestCaseOperation2Request) Send(ctx context.Context) (*InputService1TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService1TestCaseOperation2Response{
		InputService1TestShapeInputService1TestCaseOperation2Output: r.Request.Data.(*InputService1TestShapeInputService1TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService1TestCaseOperation2Response is the response type for the
// InputService1TestCaseOperation2 API operation.
type InputService1TestCaseOperation2Response struct {
	*InputService1TestShapeInputService1TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService1TestCaseOperation2 request.
func (r *InputService1TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService1TestShapeInputService1TestCaseOperation3Input struct {
	_ struct{} `type:"structure"`

	Bar *string `type:"string"`

	Baz *bool `type:"boolean"`

	Foo *string `type:"string"`
}

type InputService1TestShapeInputService1TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`
}

const opInputService1TestCaseOperation3 = "OperationName"

// InputService1TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService1TestCaseOperation3Request.
//    req := client.InputService1TestCaseOperation3Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService1ProtocolTest) InputService1TestCaseOperation3Request(input *InputService1TestShapeInputService1TestCaseOperation3Input) InputService1TestCaseOperation3Request {
	op := &aws.Operation{
		Name: opInputService1TestCaseOperation3,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation3Input{}
	}

	req := c.newRequest(op, input, &InputService1TestShapeInputService1TestCaseOperation3Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService1TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService1TestCaseOperation3Request}
}

// InputService1TestCaseOperation3Request is the request type for the
// InputService1TestCaseOperation3 API operation.
type InputService1TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService1TestShapeInputService1TestCaseOperation3Input
	Copy  func(*InputService1TestShapeInputService1TestCaseOperation3Input) InputService1TestCaseOperation3Request
}

// Send marshals and sends the InputService1TestCaseOperation3 API request.
func (r InputService1TestCaseOperation3Request) Send(ctx context.Context) (*InputService1TestCaseOperation3Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService1TestCaseOperation3Response{
		InputService1TestShapeInputService1TestCaseOperation3Output: r.Request.Data.(*InputService1TestShapeInputService1TestCaseOperation3Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService1TestCaseOperation3Response is the response type for the
// InputService1TestCaseOperation3 API operation.
type InputService1TestCaseOperation3Response struct {
	*InputService1TestShapeInputService1TestCaseOperation3Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService1TestCaseOperation3 request.
func (r *InputService1TestCaseOperation3Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService2ProtocolTest provides the API operation methods for making requests to
// InputService2ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService2ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice2protocoltest.New(myConfig)
func NewInputService2ProtocolTest(config aws.Config) *InputService2ProtocolTest {
	svc := &InputService2ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService2ProtocolTest",
				ServiceID:     "InputService2ProtocolTest",
				EndpointsID:   "inputservice2protocoltest",
				SigningName:   "inputservice2protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService2ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService2TestShapeInputService2TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	StructArg *InputService2TestShapeStructType `type:"structure"`
}

type InputService2TestShapeInputService2TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService2TestCaseOperation1 = "OperationName"

// InputService2TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService2TestCaseOperation1Request.
//    req := client.InputService2TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService2ProtocolTest) InputService2TestCaseOperation1Request(input *InputService2TestShapeInputService2TestCaseOperation1Input) InputService2TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService2TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService2TestShapeInputService2TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService2TestShapeInputService2TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService2TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService2TestCaseOperation1Request}
}

// InputService2TestCaseOperation1Request is the request type for the
// InputService2TestCaseOperation1 API operation.
type InputService2TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService2TestShapeInputService2TestCaseOperation1Input
	Copy  func(*InputService2TestShapeInputService2TestCaseOperation1Input) InputService2TestCaseOperation1Request
}

// Send marshals and sends the InputService2TestCaseOperation1 API request.
func (r InputService2TestCaseOperation1Request) Send(ctx context.Context) (*InputService2TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService2TestCaseOperation1Response{
		InputService2TestShapeInputService2TestCaseOperation1Output: r.Request.Data.(*InputService2TestShapeInputService2TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService2TestCaseOperation1Response is the response type for the
// InputService2TestCaseOperation1 API operation.
type InputService2TestCaseOperation1Response struct {
	*InputService2TestShapeInputService2TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService2TestCaseOperation1 request.
func (r *InputService2TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService2TestShapeStructType struct {
	_ struct{} `type:"structure"`

	ScalarArg *string `type:"string"`
}

// InputService3ProtocolTest provides the API operation methods for making requests to
// InputService3ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService3ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice3protocoltest.New(myConfig)
func NewInputService3ProtocolTest(config aws.Config) *InputService3ProtocolTest {
	svc := &InputService3ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService3ProtocolTest",
				ServiceID:     "InputService3ProtocolTest",
				EndpointsID:   "inputservice3protocoltest",
				SigningName:   "inputservice3protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService3ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService3TestShapeInputService3TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListArg []string `type:"list"`
}

type InputService3TestShapeInputService3TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService3TestCaseOperation1 = "OperationName"

// InputService3TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService3TestCaseOperation1Request.
//    req := client.InputService3TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService3ProtocolTest) InputService3TestCaseOperation1Request(input *InputService3TestShapeInputService3TestCaseOperation1Input) InputService3TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService3TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService3TestShapeInputService3TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService3TestShapeInputService3TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService3TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService3TestCaseOperation1Request}
}

// InputService3TestCaseOperation1Request is the request type for the
// InputService3TestCaseOperation1 API operation.
type InputService3TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService3TestShapeInputService3TestCaseOperation1Input
	Copy  func(*InputService3TestShapeInputService3TestCaseOperation1Input) InputService3TestCaseOperation1Request
}

// Send marshals and sends the InputService3TestCaseOperation1 API request.
func (r InputService3TestCaseOperation1Request) Send(ctx context.Context) (*InputService3TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService3TestCaseOperation1Response{
		InputService3TestShapeInputService3TestCaseOperation1Output: r.Request.Data.(*InputService3TestShapeInputService3TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService3TestCaseOperation1Response is the response type for the
// InputService3TestCaseOperation1 API operation.
type InputService3TestCaseOperation1Response struct {
	*InputService3TestShapeInputService3TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService3TestCaseOperation1 request.
func (r *InputService3TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService3TestShapeInputService3TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	ListArg []string `type:"list"`
}

type InputService3TestShapeInputService3TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService3TestCaseOperation2 = "OperationName"

// InputService3TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService3TestCaseOperation2Request.
//    req := client.InputService3TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService3ProtocolTest) InputService3TestCaseOperation2Request(input *InputService3TestShapeInputService3TestCaseOperation2Input) InputService3TestCaseOperation2Request {
	op := &aws.Operation{
		Name: opInputService3TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService3TestShapeInputService3TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService3TestShapeInputService3TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService3TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService3TestCaseOperation2Request}
}

// InputService3TestCaseOperation2Request is the request type for the
// InputService3TestCaseOperation2 API operation.
type InputService3TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService3TestShapeInputService3TestCaseOperation2Input
	Copy  func(*InputService3TestShapeInputService3TestCaseOperation2Input) InputService3TestCaseOperation2Request
}

// Send marshals and sends the InputService3TestCaseOperation2 API request.
func (r InputService3TestCaseOperation2Request) Send(ctx context.Context) (*InputService3TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService3TestCaseOperation2Response{
		InputService3TestShapeInputService3TestCaseOperation2Output: r.Request.Data.(*InputService3TestShapeInputService3TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService3TestCaseOperation2Response is the response type for the
// InputService3TestCaseOperation2 API operation.
type InputService3TestCaseOperation2Response struct {
	*InputService3TestShapeInputService3TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService3TestCaseOperation2 request.
func (r *InputService3TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService4ProtocolTest provides the API operation methods for making requests to
// InputService4ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService4ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice4protocoltest.New(myConfig)
func NewInputService4ProtocolTest(config aws.Config) *InputService4ProtocolTest {
	svc := &InputService4ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService4ProtocolTest",
				ServiceID:     "InputService4ProtocolTest",
				EndpointsID:   "inputservice4protocoltest",
				SigningName:   "inputservice4protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService4ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService4TestShapeInputService4TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListArg []string `type:"list" flattened:"true"`

	NamedListArg []string `locationNameList:"Foo" type:"list" flattened:"true"`

	ScalarArg *string `type:"string"`
}

type InputService4TestShapeInputService4TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService4TestCaseOperation1 = "OperationName"

// InputService4TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService4TestCaseOperation1Request.
//    req := client.InputService4TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService4ProtocolTest) InputService4TestCaseOperation1Request(input *InputService4TestShapeInputService4TestCaseOperation1Input) InputService4TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService4TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService4TestShapeInputService4TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService4TestShapeInputService4TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService4TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService4TestCaseOperation1Request}
}

// InputService4TestCaseOperation1Request is the request type for the
// InputService4TestCaseOperation1 API operation.
type InputService4TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService4TestShapeInputService4TestCaseOperation1Input
	Copy  func(*InputService4TestShapeInputService4TestCaseOperation1Input) InputService4TestCaseOperation1Request
}

// Send marshals and sends the InputService4TestCaseOperation1 API request.
func (r InputService4TestCaseOperation1Request) Send(ctx context.Context) (*InputService4TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService4TestCaseOperation1Response{
		InputService4TestShapeInputService4TestCaseOperation1Output: r.Request.Data.(*InputService4TestShapeInputService4TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService4TestCaseOperation1Response is the response type for the
// InputService4TestCaseOperation1 API operation.
type InputService4TestCaseOperation1Response struct {
	*InputService4TestShapeInputService4TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService4TestCaseOperation1 request.
func (r *InputService4TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService4TestShapeInputService4TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	ListArg []string `type:"list" flattened:"true"`

	NamedListArg []string `locationNameList:"Foo" type:"list" flattened:"true"`

	ScalarArg *string `type:"string"`
}

type InputService4TestShapeInputService4TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService4TestCaseOperation2 = "OperationName"

// InputService4TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService4TestCaseOperation2Request.
//    req := client.InputService4TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService4ProtocolTest) InputService4TestCaseOperation2Request(input *InputService4TestShapeInputService4TestCaseOperation2Input) InputService4TestCaseOperation2Request {
	op := &aws.Operation{
		Name: opInputService4TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService4TestShapeInputService4TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService4TestShapeInputService4TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService4TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService4TestCaseOperation2Request}
}

// InputService4TestCaseOperation2Request is the request type for the
// InputService4TestCaseOperation2 API operation.
type InputService4TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService4TestShapeInputService4TestCaseOperation2Input
	Copy  func(*InputService4TestShapeInputService4TestCaseOperation2Input) InputService4TestCaseOperation2Request
}

// Send marshals and sends the InputService4TestCaseOperation2 API request.
func (r InputService4TestCaseOperation2Request) Send(ctx context.Context) (*InputService4TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService4TestCaseOperation2Response{
		InputService4TestShapeInputService4TestCaseOperation2Output: r.Request.Data.(*InputService4TestShapeInputService4TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService4TestCaseOperation2Response is the response type for the
// InputService4TestCaseOperation2 API operation.
type InputService4TestCaseOperation2Response struct {
	*InputService4TestShapeInputService4TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService4TestCaseOperation2 request.
func (r *InputService4TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService5ProtocolTest provides the API operation methods for making requests to
// InputService5ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService5ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice5protocoltest.New(myConfig)
func NewInputService5ProtocolTest(config aws.Config) *InputService5ProtocolTest {
	svc := &InputService5ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService5ProtocolTest",
				ServiceID:     "InputService5ProtocolTest",
				EndpointsID:   "inputservice5protocoltest",
				SigningName:   "inputservice5protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService5ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService5TestShapeInputService5TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	MapArg map[string]string `type:"map" flattened:"true"`
}

type InputService5TestShapeInputService5TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService5TestCaseOperation1 = "OperationName"

// InputService5TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService5TestCaseOperation1Request.
//    req := client.InputService5TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation1Request(input *InputService5TestShapeInputService5TestCaseOperation1Input) InputService5TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService5TestShapeInputService5TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService5TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation1Request}
}

// InputService5TestCaseOperation1Request is the request type for the
// InputService5TestCaseOperation1 API operation.
type InputService5TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation1Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation1Input) InputService5TestCaseOperation1Request
}

// Send marshals and sends the InputService5TestCaseOperation1 API request.
func (r InputService5TestCaseOperation1Request) Send(ctx context.Context) (*InputService5TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService5TestCaseOperation1Response{
		InputService5TestShapeInputService5TestCaseOperation1Output: r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService5TestCaseOperation1Response is the response type for the
// InputService5TestCaseOperation1 API operation.
type InputService5TestCaseOperation1Response struct {
	*InputService5TestShapeInputService5TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService5TestCaseOperation1 request.
func (r *InputService5TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService6ProtocolTest provides the API operation methods for making requests to
// InputService6ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService6ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice6protocoltest.New(myConfig)
func NewInputService6ProtocolTest(config aws.Config) *InputService6ProtocolTest {
	svc := &InputService6ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService6ProtocolTest",
				ServiceID:     "InputService6ProtocolTest",
				EndpointsID:   "inputservice6protocoltest",
				SigningName:   "inputservice6protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService6ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService6TestShapeInputService6TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListArg []string `locationNameList:"item" type:"list"`
}

type InputService6TestShapeInputService6TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService6TestCaseOperation1 = "OperationName"

// InputService6TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService6TestCaseOperation1Request.
//    req := client.InputService6TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService6ProtocolTest) InputService6TestCaseOperation1Request(input *InputService6TestShapeInputService6TestCaseOperation1Input) InputService6TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService6TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService6TestShapeInputService6TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService6TestShapeInputService6TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService6TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService6TestCaseOperation1Request}
}

// InputService6TestCaseOperation1Request is the request type for the
// InputService6TestCaseOperation1 API operation.
type InputService6TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService6TestShapeInputService6TestCaseOperation1Input
	Copy  func(*InputService6TestShapeInputService6TestCaseOperation1Input) InputService6TestCaseOperation1Request
}

// Send marshals and sends the InputService6TestCaseOperation1 API request.
func (r InputService6TestCaseOperation1Request) Send(ctx context.Context) (*InputService6TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService6TestCaseOperation1Response{
		InputService6TestShapeInputService6TestCaseOperation1Output: r.Request.Data.(*InputService6TestShapeInputService6TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService6TestCaseOperation1Response is the response type for the
// InputService6TestCaseOperation1 API operation.
type InputService6TestCaseOperation1Response struct {
	*InputService6TestShapeInputService6TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService6TestCaseOperation1 request.
func (r *InputService6TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService7ProtocolTest provides the API operation methods for making requests to
// InputService7ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService7ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice7protocoltest.New(myConfig)
func NewInputService7ProtocolTest(config aws.Config) *InputService7ProtocolTest {
	svc := &InputService7ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService7ProtocolTest",
				ServiceID:     "InputService7ProtocolTest",
				EndpointsID:   "inputservice7protocoltest",
				SigningName:   "inputservice7protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService7ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService7TestShapeInputService7TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListArg []string `locationNameList:"ListArgLocation" type:"list" flattened:"true"`

	ScalarArg *string `type:"string"`
}

type InputService7TestShapeInputService7TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService7TestCaseOperation1 = "OperationName"

// InputService7TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService7TestCaseOperation1Request.
//    req := client.InputService7TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService7ProtocolTest) InputService7TestCaseOperation1Request(input *InputService7TestShapeInputService7TestCaseOperation1Input) InputService7TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService7TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService7TestShapeInputService7TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService7TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService7TestCaseOperation1Request}
}

// InputService7TestCaseOperation1Request is the request type for the
// InputService7TestCaseOperation1 API operation.
type InputService7TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService7TestShapeInputService7TestCaseOperation1Input
	Copy  func(*InputService7TestShapeInputService7TestCaseOperation1Input) InputService7TestCaseOperation1Request
}

// Send marshals and sends the InputService7TestCaseOperation1 API request.
func (r InputService7TestCaseOperation1Request) Send(ctx context.Context) (*InputService7TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService7TestCaseOperation1Response{
		InputService7TestShapeInputService7TestCaseOperation1Output: r.Request.Data.(*InputService7TestShapeInputService7TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService7TestCaseOperation1Response is the response type for the
// InputService7TestCaseOperation1 API operation.
type InputService7TestCaseOperation1Response struct {
	*InputService7TestShapeInputService7TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService7TestCaseOperation1 request.
func (r *InputService7TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService8ProtocolTest provides the API operation methods for making requests to
// InputService8ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService8ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice8protocoltest.New(myConfig)
func NewInputService8ProtocolTest(config aws.Config) *InputService8ProtocolTest {
	svc := &InputService8ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService8ProtocolTest",
				ServiceID:     "InputService8ProtocolTest",
				EndpointsID:   "inputservice8protocoltest",
				SigningName:   "inputservice8protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService8ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService8TestShapeInputService8TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	MapArg map[string]string `type:"map"`
}

type InputService8TestShapeInputService8TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService8TestCaseOperation1 = "OperationName"

// InputService8TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService8TestCaseOperation1Request.
//    req := client.InputService8TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService8ProtocolTest) InputService8TestCaseOperation1Request(input *InputService8TestShapeInputService8TestCaseOperation1Input) InputService8TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService8TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService8TestShapeInputService8TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService8TestShapeInputService8TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService8TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService8TestCaseOperation1Request}
}

// InputService8TestCaseOperation1Request is the request type for the
// InputService8TestCaseOperation1 API operation.
type InputService8TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService8TestShapeInputService8TestCaseOperation1Input
	Copy  func(*InputService8TestShapeInputService8TestCaseOperation1Input) InputService8TestCaseOperation1Request
}

// Send marshals and sends the InputService8TestCaseOperation1 API request.
func (r InputService8TestCaseOperation1Request) Send(ctx context.Context) (*InputService8TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService8TestCaseOperation1Response{
		InputService8TestShapeInputService8TestCaseOperation1Output: r.Request.Data.(*InputService8TestShapeInputService8TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService8TestCaseOperation1Response is the response type for the
// InputService8TestCaseOperation1 API operation.
type InputService8TestCaseOperation1Response struct {
	*InputService8TestShapeInputService8TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService8TestCaseOperation1 request.
func (r *InputService8TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService9ProtocolTest provides the API operation methods for making requests to
// InputService9ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService9ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice9protocoltest.New(myConfig)
func NewInputService9ProtocolTest(config aws.Config) *InputService9ProtocolTest {
	svc := &InputService9ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService9ProtocolTest",
				ServiceID:     "InputService9ProtocolTest",
				EndpointsID:   "inputservice9protocoltest",
				SigningName:   "inputservice9protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService9ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService9TestShapeInputService9TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	MapArg map[string]string `locationNameKey:"TheKey" locationNameValue:"TheValue" type:"map"`
}

type InputService9TestShapeInputService9TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService9TestCaseOperation1 = "OperationName"

// InputService9TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService9TestCaseOperation1Request.
//    req := client.InputService9TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService9ProtocolTest) InputService9TestCaseOperation1Request(input *InputService9TestShapeInputService9TestCaseOperation1Input) InputService9TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService9TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService9TestShapeInputService9TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService9TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService9TestCaseOperation1Request}
}

// InputService9TestCaseOperation1Request is the request type for the
// InputService9TestCaseOperation1 API operation.
type InputService9TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService9TestShapeInputService9TestCaseOperation1Input
	Copy  func(*InputService9TestShapeInputService9TestCaseOperation1Input) InputService9TestCaseOperation1Request
}

// Send marshals and sends the InputService9TestCaseOperation1 API request.
func (r InputService9TestCaseOperation1Request) Send(ctx context.Context) (*InputService9TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService9TestCaseOperation1Response{
		InputService9TestShapeInputService9TestCaseOperation1Output: r.Request.Data.(*InputService9TestShapeInputService9TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService9TestCaseOperation1Response is the response type for the
// InputService9TestCaseOperation1 API operation.
type InputService9TestCaseOperation1Response struct {
	*InputService9TestShapeInputService9TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService9TestCaseOperation1 request.
func (r *InputService9TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService10ProtocolTest provides the API operation methods for making requests to
// InputService10ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService10ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice10protocoltest.New(myConfig)
func NewInputService10ProtocolTest(config aws.Config) *InputService10ProtocolTest {
	svc := &InputService10ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService10ProtocolTest",
				ServiceID:     "InputService10ProtocolTest",
				EndpointsID:   "inputservice10protocoltest",
				SigningName:   "inputservice10protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService10ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService10TestShapeInputService10TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	// BlobArg is automatically base64 encoded/decoded by the SDK.
	BlobArg []byte `type:"blob"`
}

type InputService10TestShapeInputService10TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService10TestCaseOperation1 = "OperationName"

// InputService10TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService10TestCaseOperation1Request.
//    req := client.InputService10TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService10ProtocolTest) InputService10TestCaseOperation1Request(input *InputService10TestShapeInputService10TestCaseOperation1Input) InputService10TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService10TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService10TestShapeInputService10TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService10TestShapeInputService10TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService10TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService10TestCaseOperation1Request}
}

// InputService10TestCaseOperation1Request is the request type for the
// InputService10TestCaseOperation1 API operation.
type InputService10TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService10TestShapeInputService10TestCaseOperation1Input
	Copy  func(*InputService10TestShapeInputService10TestCaseOperation1Input) InputService10TestCaseOperation1Request
}

// Send marshals and sends the InputService10TestCaseOperation1 API request.
func (r InputService10TestCaseOperation1Request) Send(ctx context.Context) (*InputService10TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService10TestCaseOperation1Response{
		InputService10TestShapeInputService10TestCaseOperation1Output: r.Request.Data.(*InputService10TestShapeInputService10TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService10TestCaseOperation1Response is the response type for the
// InputService10TestCaseOperation1 API operation.
type InputService10TestCaseOperation1Response struct {
	*InputService10TestShapeInputService10TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService10TestCaseOperation1 request.
func (r *InputService10TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService11ProtocolTest provides the API operation methods for making requests to
// InputService11ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService11ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice11protocoltest.New(myConfig)
func NewInputService11ProtocolTest(config aws.Config) *InputService11ProtocolTest {
	svc := &InputService11ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService11ProtocolTest",
				ServiceID:     "InputService11ProtocolTest",
				EndpointsID:   "inputservice11protocoltest",
				SigningName:   "inputservice11protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService11ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService11TestShapeInputService11TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	BlobArgs [][]byte `type:"list" flattened:"true"`
}

type InputService11TestShapeInputService11TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService11TestCaseOperation1 = "OperationName"

// InputService11TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService11TestCaseOperation1Request.
//    req := client.InputService11TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService11ProtocolTest) InputService11TestCaseOperation1Request(input *InputService11TestShapeInputService11TestCaseOperation1Input) InputService11TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService11TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService11TestShapeInputService11TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService11TestShapeInputService11TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService11TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService11TestCaseOperation1Request}
}

// InputService11TestCaseOperation1Request is the request type for the
// InputService11TestCaseOperation1 API operation.
type InputService11TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService11TestShapeInputService11TestCaseOperation1Input
	Copy  func(*InputService11TestShapeInputService11TestCaseOperation1Input) InputService11TestCaseOperation1Request
}

// Send marshals and sends the InputService11TestCaseOperation1 API request.
func (r InputService11TestCaseOperation1Request) Send(ctx context.Context) (*InputService11TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService11TestCaseOperation1Response{
		InputService11TestShapeInputService11TestCaseOperation1Output: r.Request.Data.(*InputService11TestShapeInputService11TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService11TestCaseOperation1Response is the response type for the
// InputService11TestCaseOperation1 API operation.
type InputService11TestCaseOperation1Response struct {
	*InputService11TestShapeInputService11TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService11TestCaseOperation1 request.
func (r *InputService11TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService12ProtocolTest provides the API operation methods for making requests to
// InputService12ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService12ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice12protocoltest.New(myConfig)
func NewInputService12ProtocolTest(config aws.Config) *InputService12ProtocolTest {
	svc := &InputService12ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService12ProtocolTest",
				ServiceID:     "InputService12ProtocolTest",
				EndpointsID:   "inputservice12protocoltest",
				SigningName:   "inputservice12protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService12ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService12TestShapeInputService12TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	TimeArg *time.Time `type:"timestamp"`

	TimeCustom *time.Time `type:"timestamp" timestampFormat:"unixTimestamp"`

	TimeFormat *time.Time `type:"timestamp" timestampFormat:"unixTimestamp"`
}

type InputService12TestShapeInputService12TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService12TestCaseOperation1 = "OperationName"

// InputService12TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService12TestCaseOperation1Request.
//    req := client.InputService12TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService12ProtocolTest) InputService12TestCaseOperation1Request(input *InputService12TestShapeInputService12TestCaseOperation1Input) InputService12TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService12TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService12TestShapeInputService12TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService12TestShapeInputService12TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService12TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService12TestCaseOperation1Request}
}

// InputService12TestCaseOperation1Request is the request type for the
// InputService12TestCaseOperation1 API operation.
type InputService12TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService12TestShapeInputService12TestCaseOperation1Input
	Copy  func(*InputService12TestShapeInputService12TestCaseOperation1Input) InputService12TestCaseOperation1Request
}

// Send marshals and sends the InputService12TestCaseOperation1 API request.
func (r InputService12TestCaseOperation1Request) Send(ctx context.Context) (*InputService12TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService12TestCaseOperation1Response{
		InputService12TestShapeInputService12TestCaseOperation1Output: r.Request.Data.(*InputService12TestShapeInputService12TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService12TestCaseOperation1Response is the response type for the
// InputService12TestCaseOperation1 API operation.
type InputService12TestCaseOperation1Response struct {
	*InputService12TestShapeInputService12TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService12TestCaseOperation1 request.
func (r *InputService12TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService13ProtocolTest provides the API operation methods for making requests to
// InputService13ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService13ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice13protocoltest.New(myConfig)
func NewInputService13ProtocolTest(config aws.Config) *InputService13ProtocolTest {
	svc := &InputService13ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService13ProtocolTest",
				ServiceID:     "InputService13ProtocolTest",
				EndpointsID:   "inputservice13protocoltest",
				SigningName:   "inputservice13protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService13ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService13TestShapeInputService13TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService13TestShapeRecursiveStructType `type:"structure"`
}

type InputService13TestShapeInputService13TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService13TestCaseOperation1 = "OperationName"

// InputService13TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService13TestCaseOperation1Request.
//    req := client.InputService13TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService13ProtocolTest) InputService13TestCaseOperation1Request(input *InputService13TestShapeInputService13TestCaseOperation1Input) InputService13TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService13TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService13TestShapeInputService13TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService13TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService13TestCaseOperation1Request}
}

// InputService13TestCaseOperation1Request is the request type for the
// InputService13TestCaseOperation1 API operation.
type InputService13TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService13TestShapeInputService13TestCaseOperation1Input
	Copy  func(*InputService13TestShapeInputService13TestCaseOperation1Input) InputService13TestCaseOperation1Request
}

// Send marshals and sends the InputService13TestCaseOperation1 API request.
func (r InputService13TestCaseOperation1Request) Send(ctx context.Context) (*InputService13TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService13TestCaseOperation1Response{
		InputService13TestShapeInputService13TestCaseOperation1Output: r.Request.Data.(*InputService13TestShapeInputService13TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService13TestCaseOperation1Response is the response type for the
// InputService13TestCaseOperation1 API operation.
type InputService13TestCaseOperation1Response struct {
	*InputService13TestShapeInputService13TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService13TestCaseOperation1 request.
func (r *InputService13TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService13TestShapeInputService13TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService13TestShapeRecursiveStructType `type:"structure"`
}

type InputService13TestShapeInputService13TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService13TestCaseOperation2 = "OperationName"

// InputService13TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService13TestCaseOperation2Request.
//    req := client.InputService13TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService13ProtocolTest) InputService13TestCaseOperation2Request(input *InputService13TestShapeInputService13TestCaseOperation2Input) InputService13TestCaseOperation2Request {
	op := &aws.Operation{
		Name: opInputService13TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService13TestShapeInputService13TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService13TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService13TestCaseOperation2Request}
}

// InputService13TestCaseOperation2Request is the request type for the
// InputService13TestCaseOperation2 API operation.
type InputService13TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService13TestShapeInputService13TestCaseOperation2Input
	Copy  func(*InputService13TestShapeInputService13TestCaseOperation2Input) InputService13TestCaseOperation2Request
}

// Send marshals and sends the InputService13TestCaseOperation2 API request.
func (r InputService13TestCaseOperation2Request) Send(ctx context.Context) (*InputService13TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService13TestCaseOperation2Response{
		InputService13TestShapeInputService13TestCaseOperation2Output: r.Request.Data.(*InputService13TestShapeInputService13TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService13TestCaseOperation2Response is the response type for the
// InputService13TestCaseOperation2 API operation.
type InputService13TestCaseOperation2Response struct {
	*InputService13TestShapeInputService13TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService13TestCaseOperation2 request.
func (r *InputService13TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService13TestShapeInputService13TestCaseOperation3Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService13TestShapeRecursiveStructType `type:"structure"`
}

type InputService13TestShapeInputService13TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`
}

const opInputService13TestCaseOperation3 = "OperationName"

// InputService13TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService13TestCaseOperation3Request.
//    req := client.InputService13TestCaseOperation3Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService13ProtocolTest) InputService13TestCaseOperation3Request(input *InputService13TestShapeInputService13TestCaseOperation3Input) InputService13TestCaseOperation3Request {
	op := &aws.Operation{
		Name: opInputService13TestCaseOperation3,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation3Input{}
	}

	req := c.newRequest(op, input, &InputService13TestShapeInputService13TestCaseOperation3Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService13TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService13TestCaseOperation3Request}
}

// InputService13TestCaseOperation3Request is the request type for the
// InputService13TestCaseOperation3 API operation.
type InputService13TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService13TestShapeInputService13TestCaseOperation3Input
	Copy  func(*InputService13TestShapeInputService13TestCaseOperation3Input) InputService13TestCaseOperation3Request
}

// Send marshals and sends the InputService13TestCaseOperation3 API request.
func (r InputService13TestCaseOperation3Request) Send(ctx context.Context) (*InputService13TestCaseOperation3Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService13TestCaseOperation3Response{
		InputService13TestShapeInputService13TestCaseOperation3Output: r.Request.Data.(*InputService13TestShapeInputService13TestCaseOperation3Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService13TestCaseOperation3Response is the response type for the
// InputService13TestCaseOperation3 API operation.
type InputService13TestCaseOperation3Response struct {
	*InputService13TestShapeInputService13TestCaseOperation3Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService13TestCaseOperation3 request.
func (r *InputService13TestCaseOperation3Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService13TestShapeInputService13TestCaseOperation4Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService13TestShapeRecursiveStructType `type:"structure"`
}

type InputService13TestShapeInputService13TestCaseOperation4Output struct {
	_ struct{} `type:"structure"`
}

const opInputService13TestCaseOperation4 = "OperationName"

// InputService13TestCaseOperation4Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService13TestCaseOperation4Request.
//    req := client.InputService13TestCaseOperation4Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService13ProtocolTest) InputService13TestCaseOperation4Request(input *InputService13TestShapeInputService13TestCaseOperation4Input) InputService13TestCaseOperation4Request {
	op := &aws.Operation{
		Name: opInputService13TestCaseOperation4,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation4Input{}
	}

	req := c.newRequest(op, input, &InputService13TestShapeInputService13TestCaseOperation4Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService13TestCaseOperation4Request{Request: req, Input: input, Copy: c.InputService13TestCaseOperation4Request}
}

// InputService13TestCaseOperation4Request is the request type for the
// InputService13TestCaseOperation4 API operation.
type InputService13TestCaseOperation4Request struct {
	*aws.Request
	Input *InputService13TestShapeInputService13TestCaseOperation4Input
	Copy  func(*InputService13TestShapeInputService13TestCaseOperation4Input) InputService13TestCaseOperation4Request
}

// Send marshals and sends the InputService13TestCaseOperation4 API request.
func (r InputService13TestCaseOperation4Request) Send(ctx context.Context) (*InputService13TestCaseOperation4Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService13TestCaseOperation4Response{
		InputService13TestShapeInputService13TestCaseOperation4Output: r.Request.Data.(*InputService13TestShapeInputService13TestCaseOperation4Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService13TestCaseOperation4Response is the response type for the
// InputService13TestCaseOperation4 API operation.
type InputService13TestCaseOperation4Response struct {
	*InputService13TestShapeInputService13TestCaseOperation4Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService13TestCaseOperation4 request.
func (r *InputService13TestCaseOperation4Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService13TestShapeInputService13TestCaseOperation5Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService13TestShapeRecursiveStructType `type:"structure"`
}

type InputService13TestShapeInputService13TestCaseOperation5Output struct {
	_ struct{} `type:"structure"`
}

const opInputService13TestCaseOperation5 = "OperationName"

// InputService13TestCaseOperation5Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService13TestCaseOperation5Request.
//    req := client.InputService13TestCaseOperation5Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService13ProtocolTest) InputService13TestCaseOperation5Request(input *InputService13TestShapeInputService13TestCaseOperation5Input) InputService13TestCaseOperation5Request {
	op := &aws.Operation{
		Name: opInputService13TestCaseOperation5,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation5Input{}
	}

	req := c.newRequest(op, input, &InputService13TestShapeInputService13TestCaseOperation5Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService13TestCaseOperation5Request{Request: req, Input: input, Copy: c.InputService13TestCaseOperation5Request}
}

// InputService13TestCaseOperation5Request is the request type for the
// InputService13TestCaseOperation5 API operation.
type InputService13TestCaseOperation5Request struct {
	*aws.Request
	Input *InputService13TestShapeInputService13TestCaseOperation5Input
	Copy  func(*InputService13TestShapeInputService13TestCaseOperation5Input) InputService13TestCaseOperation5Request
}

// Send marshals and sends the InputService13TestCaseOperation5 API request.
func (r InputService13TestCaseOperation5Request) Send(ctx context.Context) (*InputService13TestCaseOperation5Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService13TestCaseOperation5Response{
		InputService13TestShapeInputService13TestCaseOperation5Output: r.Request.Data.(*InputService13TestShapeInputService13TestCaseOperation5Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService13TestCaseOperation5Response is the response type for the
// InputService13TestCaseOperation5 API operation.
type InputService13TestCaseOperation5Response struct {
	*InputService13TestShapeInputService13TestCaseOperation5Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService13TestCaseOperation5 request.
func (r *InputService13TestCaseOperation5Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService13TestShapeInputService13TestCaseOperation6Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService13TestShapeRecursiveStructType `type:"structure"`
}

type InputService13TestShapeInputService13TestCaseOperation6Output struct {
	_ struct{} `type:"structure"`
}

const opInputService13TestCaseOperation6 = "OperationName"

// InputService13TestCaseOperation6Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService13TestCaseOperation6Request.
//    req := client.InputService13TestCaseOperation6Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService13ProtocolTest) InputService13TestCaseOperation6Request(input *InputService13TestShapeInputService13TestCaseOperation6Input) InputService13TestCaseOperation6Request {
	op := &aws.Operation{
		Name: opInputService13TestCaseOperation6,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation6Input{}
	}

	req := c.newRequest(op, input, &InputService13TestShapeInputService13TestCaseOperation6Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService13TestCaseOperation6Request{Request: req, Input: input, Copy: c.InputService13TestCaseOperation6Request}
}

// InputService13TestCaseOperation6Request is the request type for the
// InputService13TestCaseOperation6 API operation.
type InputService13TestCaseOperation6Request struct {
	*aws.Request
	Input *InputService13TestShapeInputService13TestCaseOperation6Input
	Copy  func(*InputService13TestShapeInputService13TestCaseOperation6Input) InputService13TestCaseOperation6Request
}

// Send marshals and sends the InputService13TestCaseOperation6 API request.
func (r InputService13TestCaseOperation6Request) Send(ctx context.Context) (*InputService13TestCaseOperation6Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService13TestCaseOperation6Response{
		InputService13TestShapeInputService13TestCaseOperation6Output: r.Request.Data.(*InputService13TestShapeInputService13TestCaseOperation6Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService13TestCaseOperation6Response is the response type for the
// InputService13TestCaseOperation6 API operation.
type InputService13TestCaseOperation6Response struct {
	*InputService13TestShapeInputService13TestCaseOperation6Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService13TestCaseOperation6 request.
func (r *InputService13TestCaseOperation6Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService13TestShapeRecursiveStructType struct {
	_ struct{} `type:"structure"`

	NoRecurse *string `type:"string"`

	RecursiveList []InputService13TestShapeRecursiveStructType `type:"list"`

	RecursiveMap map[string]InputService13TestShapeRecursiveStructType `type:"map"`

	RecursiveStruct *InputService13TestShapeRecursiveStructType `type:"structure"`
}

// InputService14ProtocolTest provides the API operation methods for making requests to
// InputService14ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService14ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice14protocoltest.New(myConfig)
func NewInputService14ProtocolTest(config aws.Config) *InputService14ProtocolTest {
	svc := &InputService14ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService14ProtocolTest",
				ServiceID:     "InputService14ProtocolTest",
				EndpointsID:   "inputservice14protocoltest",
				SigningName:   "inputservice14protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService14ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService14TestShapeInputService14TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Token *string `type:"string" idempotencyToken:"true"`
}

type InputService14TestShapeInputService14TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService14TestCaseOperation1 = "OperationName"

// InputService14TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService14TestCaseOperation1Request.
//    req := client.InputService14TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService14ProtocolTest) InputService14TestCaseOperation1Request(input *InputService14TestShapeInputService14TestCaseOperation1Input) InputService14TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService14TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService14TestShapeInputService14TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService14TestShapeInputService14TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService14TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService14TestCaseOperation1Request}
}

// InputService14TestCaseOperation1Request is the request type for the
// InputService14TestCaseOperation1 API operation.
type InputService14TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService14TestShapeInputService14TestCaseOperation1Input
	Copy  func(*InputService14TestShapeInputService14TestCaseOperation1Input) InputService14TestCaseOperation1Request
}

// Send marshals and sends the InputService14TestCaseOperation1 API request.
func (r InputService14TestCaseOperation1Request) Send(ctx context.Context) (*InputService14TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService14TestCaseOperation1Response{
		InputService14TestShapeInputService14TestCaseOperation1Output: r.Request.Data.(*InputService14TestShapeInputService14TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService14TestCaseOperation1Response is the response type for the
// InputService14TestCaseOperation1 API operation.
type InputService14TestCaseOperation1Response struct {
	*InputService14TestShapeInputService14TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService14TestCaseOperation1 request.
func (r *InputService14TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService14TestShapeInputService14TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	Token *string `type:"string" idempotencyToken:"true"`
}

type InputService14TestShapeInputService14TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService14TestCaseOperation2 = "OperationName"

// InputService14TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService14TestCaseOperation2Request.
//    req := client.InputService14TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService14ProtocolTest) InputService14TestCaseOperation2Request(input *InputService14TestShapeInputService14TestCaseOperation2Input) InputService14TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService14TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService14TestShapeInputService14TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService14TestShapeInputService14TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService14TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService14TestCaseOperation2Request}
}

// InputService14TestCaseOperation2Request is the request type for the
// InputService14TestCaseOperation2 API operation.
type InputService14TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService14TestShapeInputService14TestCaseOperation2Input
	Copy  func(*InputService14TestShapeInputService14TestCaseOperation2Input) InputService14TestCaseOperation2Request
}

// Send marshals and sends the InputService14TestCaseOperation2 API request.
func (r InputService14TestCaseOperation2Request) Send(ctx context.Context) (*InputService14TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService14TestCaseOperation2Response{
		InputService14TestShapeInputService14TestCaseOperation2Output: r.Request.Data.(*InputService14TestShapeInputService14TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService14TestCaseOperation2Response is the response type for the
// InputService14TestCaseOperation2 API operation.
type InputService14TestCaseOperation2Response struct {
	*InputService14TestShapeInputService14TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService14TestCaseOperation2 request.
func (r *InputService14TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService15ProtocolTest provides the API operation methods for making requests to
// InputService15ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService15ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice15protocoltest.New(myConfig)
func NewInputService15ProtocolTest(config aws.Config) *InputService15ProtocolTest {
	svc := &InputService15ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService15ProtocolTest",
				ServiceID:     "InputService15ProtocolTest",
				EndpointsID:   "inputservice15protocoltest",
				SigningName:   "inputservice15protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService15ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService15TestShapeInputService15TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService15TestShapeEnumType `type:"string" enum:"true"`

	ListEnums []InputService15TestShapeEnumType `type:"list"`
}

type InputService15TestShapeInputService15TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService15TestCaseOperation1 = "OperationName"

// InputService15TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService15TestCaseOperation1Request.
//    req := client.InputService15TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService15ProtocolTest) InputService15TestCaseOperation1Request(input *InputService15TestShapeInputService15TestCaseOperation1Input) InputService15TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService15TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService15TestShapeInputService15TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService15TestShapeInputService15TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService15TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService15TestCaseOperation1Request}
}

// InputService15TestCaseOperation1Request is the request type for the
// InputService15TestCaseOperation1 API operation.
type InputService15TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService15TestShapeInputService15TestCaseOperation1Input
	Copy  func(*InputService15TestShapeInputService15TestCaseOperation1Input) InputService15TestCaseOperation1Request
}

// Send marshals and sends the InputService15TestCaseOperation1 API request.
func (r InputService15TestCaseOperation1Request) Send(ctx context.Context) (*InputService15TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService15TestCaseOperation1Response{
		InputService15TestShapeInputService15TestCaseOperation1Output: r.Request.Data.(*InputService15TestShapeInputService15TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService15TestCaseOperation1Response is the response type for the
// InputService15TestCaseOperation1 API operation.
type InputService15TestCaseOperation1Response struct {
	*InputService15TestShapeInputService15TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService15TestCaseOperation1 request.
func (r *InputService15TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService15TestShapeInputService15TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService15TestShapeEnumType `type:"string" enum:"true"`

	ListEnums []InputService15TestShapeEnumType `type:"list"`
}

type InputService15TestShapeInputService15TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService15TestCaseOperation2 = "OperationName"

// InputService15TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService15TestCaseOperation2Request.
//    req := client.InputService15TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService15ProtocolTest) InputService15TestCaseOperation2Request(input *InputService15TestShapeInputService15TestCaseOperation2Input) InputService15TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService15TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService15TestShapeInputService15TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService15TestShapeInputService15TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService15TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService15TestCaseOperation2Request}
}

// InputService15TestCaseOperation2Request is the request type for the
// InputService15TestCaseOperation2 API operation.
type InputService15TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService15TestShapeInputService15TestCaseOperation2Input
	Copy  func(*InputService15TestShapeInputService15TestCaseOperation2Input) InputService15TestCaseOperation2Request
}

// Send marshals and sends the InputService15TestCaseOperation2 API request.
func (r InputService15TestCaseOperation2Request) Send(ctx context.Context) (*InputService15TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService15TestCaseOperation2Response{
		InputService15TestShapeInputService15TestCaseOperation2Output: r.Request.Data.(*InputService15TestShapeInputService15TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService15TestCaseOperation2Response is the response type for the
// InputService15TestCaseOperation2 API operation.
type InputService15TestCaseOperation2Response struct {
	*InputService15TestShapeInputService15TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService15TestCaseOperation2 request.
func (r *InputService15TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService15TestShapeInputService15TestCaseOperation3Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService15TestShapeEnumType `type:"string" enum:"true"`

	ListEnums []InputService15TestShapeEnumType `type:"list"`
}

type InputService15TestShapeInputService15TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`
}

const opInputService15TestCaseOperation3 = "OperationName"

// InputService15TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService15TestCaseOperation3Request.
//    req := client.InputService15TestCaseOperation3Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService15ProtocolTest) InputService15TestCaseOperation3Request(input *InputService15TestShapeInputService15TestCaseOperation3Input) InputService15TestCaseOperation3Request {
	op := &aws.Operation{
		Name:       opInputService15TestCaseOperation3,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService15TestShapeInputService15TestCaseOperation3Input{}
	}

	req := c.newRequest(op, input, &InputService15TestShapeInputService15TestCaseOperation3Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService15TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService15TestCaseOperation3Request}
}

// InputService15TestCaseOperation3Request is the request type for the
// InputService15TestCaseOperation3 API operation.
type InputService15TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService15TestShapeInputService15TestCaseOperation3Input
	Copy  func(*InputService15TestShapeInputService15TestCaseOperation3Input) InputService15TestCaseOperation3Request
}

// Send marshals and sends the InputService15TestCaseOperation3 API request.
func (r InputService15TestCaseOperation3Request) Send(ctx context.Context) (*InputService15TestCaseOperation3Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService15TestCaseOperation3Response{
		InputService15TestShapeInputService15TestCaseOperation3Output: r.Request.Data.(*InputService15TestShapeInputService15TestCaseOperation3Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService15TestCaseOperation3Response is the response type for the
// InputService15TestCaseOperation3 API operation.
type InputService15TestCaseOperation3Response struct {
	*InputService15TestShapeInputService15TestCaseOperation3Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService15TestCaseOperation3 request.
func (r *InputService15TestCaseOperation3Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService15TestShapeEnumType string

// Enum values for InputService15TestShapeEnumType
const (
	EnumTypeFoo InputService15TestShapeEnumType = "foo"
	EnumTypeBar InputService15TestShapeEnumType = "bar"
)

func (enum InputService15TestShapeEnumType) MarshalValue() (string, error) {
	return string(enum), nil
}

func (enum InputService15TestShapeEnumType) MarshalValueBuf(b []byte) ([]byte, error) {
	b = b[0:0]
	return append(b, enum...), nil
}

// InputService16ProtocolTest provides the API operation methods for making requests to
// InputService16ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService16ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice16protocoltest.New(myConfig)
func NewInputService16ProtocolTest(config aws.Config) *InputService16ProtocolTest {
	svc := &InputService16ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService16ProtocolTest",
				ServiceID:     "InputService16ProtocolTest",
				EndpointsID:   "inputservice16protocoltest",
				SigningName:   "inputservice16protocoltest",
				SigningRegion: config.Region,
				APIVersion:    "2014-01-01",
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
func (c *InputService16ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService16TestShapeInputService16TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Name *string `type:"string"`
}

type InputService16TestShapeInputService16TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService16TestCaseOperation1 = "StaticOp"

// InputService16TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService16TestCaseOperation1Request.
//    req := client.InputService16TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation1Request(input *InputService16TestShapeInputService16TestCaseOperation1Input) InputService16TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("data-", nil))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)

	return InputService16TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation1Request}
}

// InputService16TestCaseOperation1Request is the request type for the
// InputService16TestCaseOperation1 API operation.
type InputService16TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation1Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation1Input) InputService16TestCaseOperation1Request
}

// Send marshals and sends the InputService16TestCaseOperation1 API request.
func (r InputService16TestCaseOperation1Request) Send(ctx context.Context) (*InputService16TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService16TestCaseOperation1Response{
		InputService16TestShapeInputService16TestCaseOperation1Output: r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService16TestCaseOperation1Response is the response type for the
// InputService16TestCaseOperation1 API operation.
type InputService16TestCaseOperation1Response struct {
	*InputService16TestShapeInputService16TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService16TestCaseOperation1 request.
func (r *InputService16TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService16TestShapeInputService16TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	// Name is a required field
	Name *string `type:"string" required:"true"`
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *InputService16TestShapeInputService16TestCaseOperation2Input) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "InputService16TestShapeInputService16TestCaseOperation2Input"}

	if s.Name == nil {
		invalidParams.Add(aws.NewErrParamRequired("Name"))
	}
	if s.Name != nil && len(*s.Name) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("Name", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

func (s *InputService16TestShapeInputService16TestCaseOperation2Input) hostLabels() map[string]string {
	return map[string]string{
		"Name": aws.StringValue(s.Name),
	}
}

type InputService16TestShapeInputService16TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService16TestCaseOperation2 = "MemberRefOp"

// InputService16TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService16TestCaseOperation2Request.
//    req := client.InputService16TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation2Request(input *InputService16TestShapeInputService16TestCaseOperation2Input) InputService16TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("foo-{Name}.", input.hostLabels))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)

	return InputService16TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation2Request}
}

// InputService16TestCaseOperation2Request is the request type for the
// InputService16TestCaseOperation2 API operation.
type InputService16TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation2Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation2Input) InputService16TestCaseOperation2Request
}

// Send marshals and sends the InputService16TestCaseOperation2 API request.
func (r InputService16TestCaseOperation2Request) Send(ctx context.Context) (*InputService16TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService16TestCaseOperation2Response{
		InputService16TestShapeInputService16TestCaseOperation2Output: r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService16TestCaseOperation2Response is the response type for the
// InputService16TestCaseOperation2 API operation.
type InputService16TestCaseOperation2Response struct {
	*InputService16TestShapeInputService16TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService16TestCaseOperation2 request.
func (r *InputService16TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

//
// Tests begin here
//

func TestInputService1ProtocolTestScalarMembersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService1ProtocolTest(cfg)
	input := &InputService1TestShapeInputService1TestCaseOperation1Input{
		Bar: aws.String("val2"),
		Foo: aws.String("val1"),
	}

	req := svc.InputService1TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&Bar=val2&Foo=val1&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService1ProtocolTestScalarMembersCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService1ProtocolTest(cfg)
	input := &InputService1TestShapeInputService1TestCaseOperation2Input{
		Baz: aws.Bool(true),
	}

	req := svc.InputService1TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&Baz=true&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService1ProtocolTestScalarMembersCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService1ProtocolTest(cfg)
	input := &InputService1TestShapeInputService1TestCaseOperation3Input{
		Baz: aws.Bool(false),
	}

	req := svc.InputService1TestCaseOperation3Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&Baz=false&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService2ProtocolTestNestedStructureMembersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService2ProtocolTest(cfg)
	input := &InputService2TestShapeInputService2TestCaseOperation1Input{
		StructArg: &InputService2TestShapeStructType{
			ScalarArg: aws.String("foo"),
		},
	}

	req := svc.InputService2TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&StructArg.ScalarArg=foo&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService3ProtocolTestListTypesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService3ProtocolTest(cfg)
	input := &InputService3TestShapeInputService3TestCaseOperation1Input{
		ListArg: []string{
			"foo",
			"bar",
			"baz",
		},
	}

	req := svc.InputService3TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&ListArg.member.1=foo&ListArg.member.2=bar&ListArg.member.3=baz&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService3ProtocolTestListTypesCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService3ProtocolTest(cfg)
	input := &InputService3TestShapeInputService3TestCaseOperation2Input{
		ListArg: []string{},
	}

	req := svc.InputService3TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&ListArg=&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService4ProtocolTestFlattenedListCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService4ProtocolTest(cfg)
	input := &InputService4TestShapeInputService4TestCaseOperation1Input{
		ListArg: []string{
			"a",
			"b",
			"c",
		},
		ScalarArg: aws.String("foo"),
	}

	req := svc.InputService4TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&ListArg.1=a&ListArg.2=b&ListArg.3=c&ScalarArg=foo&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService4ProtocolTestFlattenedListCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService4ProtocolTest(cfg)
	input := &InputService4TestShapeInputService4TestCaseOperation2Input{
		NamedListArg: []string{
			"a",
		},
	}

	req := svc.InputService4TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&Foo.1=a&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService5ProtocolTestSerializeFlattenedMapTypeCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation1Input{
		MapArg: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
	}

	req := svc.InputService5TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&MapArg.1.key=key1&MapArg.1.value=val1&MapArg.2.key=key2&MapArg.2.value=val2&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService6ProtocolTestNonFlattenedListWithLocationNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService6ProtocolTest(cfg)
	input := &InputService6TestShapeInputService6TestCaseOperation1Input{
		ListArg: []string{
			"a",
			"b",
			"c",
		},
	}

	req := svc.InputService6TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&ListArg.item.1=a&ListArg.item.2=b&ListArg.item.3=c&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService7ProtocolTestFlattenedListWithLocationNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService7ProtocolTest(cfg)
	input := &InputService7TestShapeInputService7TestCaseOperation1Input{
		ListArg: []string{
			"a",
			"b",
			"c",
		},
		ScalarArg: aws.String("foo"),
	}

	req := svc.InputService7TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&ListArgLocation.1=a&ListArgLocation.2=b&ListArgLocation.3=c&ScalarArg=foo&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService8ProtocolTestSerializeMapTypeCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService8ProtocolTest(cfg)
	input := &InputService8TestShapeInputService8TestCaseOperation1Input{
		MapArg: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
	}

	req := svc.InputService8TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&MapArg.entry.1.key=key1&MapArg.entry.1.value=val1&MapArg.entry.2.key=key2&MapArg.entry.2.value=val2&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService9ProtocolTestSerializeMapTypeWithLocationNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService9ProtocolTest(cfg)
	input := &InputService9TestShapeInputService9TestCaseOperation1Input{
		MapArg: map[string]string{
			"key1": "val1",
			"key2": "val2",
		},
	}

	req := svc.InputService9TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&MapArg.entry.1.TheKey=key1&MapArg.entry.1.TheValue=val1&MapArg.entry.2.TheKey=key2&MapArg.entry.2.TheValue=val2&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService10ProtocolTestBase64EncodedBlobsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService10ProtocolTest(cfg)
	input := &InputService10TestShapeInputService10TestCaseOperation1Input{
		BlobArg: []byte("foo"),
	}

	req := svc.InputService10TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&BlobArg=Zm9v&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService11ProtocolTestBase64EncodedBlobsNestedCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService11ProtocolTest(cfg)
	input := &InputService11TestShapeInputService11TestCaseOperation1Input{
		BlobArgs: [][]byte{
			[]byte("foo"),
		},
	}

	req := svc.InputService11TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&BlobArgs.1=Zm9v&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService12ProtocolTestTimestampValuesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService12ProtocolTest(cfg)
	input := &InputService12TestShapeInputService12TestCaseOperation1Input{
		TimeArg:    aws.Time(time.Unix(1422172800, 0)),
		TimeCustom: aws.Time(time.Unix(1422172800, 0)),
		TimeFormat: aws.Time(time.Unix(1422172800, 0)),
	}

	req := svc.InputService12TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&TimeArg=2015-01-25T08%3A00%3A00Z&TimeCustom=1422172800&TimeFormat=1422172800&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService13ProtocolTestRecursiveShapesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation1Input{
		RecursiveStruct: &InputService13TestShapeRecursiveStructType{
			NoRecurse: aws.String("foo"),
		},
	}

	req := svc.InputService13TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&RecursiveStruct.NoRecurse=foo&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService13ProtocolTestRecursiveShapesCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation2Input{
		RecursiveStruct: &InputService13TestShapeRecursiveStructType{
			RecursiveStruct: &InputService13TestShapeRecursiveStructType{
				NoRecurse: aws.String("foo"),
			},
		},
	}

	req := svc.InputService13TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&RecursiveStruct.RecursiveStruct.NoRecurse=foo&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService13ProtocolTestRecursiveShapesCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation3Input{
		RecursiveStruct: &InputService13TestShapeRecursiveStructType{
			RecursiveStruct: &InputService13TestShapeRecursiveStructType{
				RecursiveStruct: &InputService13TestShapeRecursiveStructType{
					RecursiveStruct: &InputService13TestShapeRecursiveStructType{
						NoRecurse: aws.String("foo"),
					},
				},
			},
		},
	}

	req := svc.InputService13TestCaseOperation3Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&RecursiveStruct.RecursiveStruct.RecursiveStruct.RecursiveStruct.NoRecurse=foo&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService13ProtocolTestRecursiveShapesCase4(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation4Input{
		RecursiveStruct: &InputService13TestShapeRecursiveStructType{
			RecursiveList: []InputService13TestShapeRecursiveStructType{
				{
					NoRecurse: aws.String("foo"),
				},
				{
					NoRecurse: aws.String("bar"),
				},
			},
		},
	}

	req := svc.InputService13TestCaseOperation4Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&RecursiveStruct.RecursiveList.member.1.NoRecurse=foo&RecursiveStruct.RecursiveList.member.2.NoRecurse=bar&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService13ProtocolTestRecursiveShapesCase5(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation5Input{
		RecursiveStruct: &InputService13TestShapeRecursiveStructType{
			RecursiveList: []InputService13TestShapeRecursiveStructType{
				{
					NoRecurse: aws.String("foo"),
				},
				{
					RecursiveStruct: &InputService13TestShapeRecursiveStructType{
						NoRecurse: aws.String("bar"),
					},
				},
			},
		},
	}

	req := svc.InputService13TestCaseOperation5Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&RecursiveStruct.RecursiveList.member.1.NoRecurse=foo&RecursiveStruct.RecursiveList.member.2.RecursiveStruct.NoRecurse=bar&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService13ProtocolTestRecursiveShapesCase6(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation6Input{
		RecursiveStruct: &InputService13TestShapeRecursiveStructType{
			RecursiveMap: map[string]InputService13TestShapeRecursiveStructType{
				"bar": {
					NoRecurse: aws.String("bar"),
				},
				"foo": {
					NoRecurse: aws.String("foo"),
				},
			},
		},
	}

	req := svc.InputService13TestCaseOperation6Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&RecursiveStruct.RecursiveMap.entry.1.key=foo&RecursiveStruct.RecursiveMap.entry.1.value.NoRecurse=foo&RecursiveStruct.RecursiveMap.entry.2.key=bar&RecursiveStruct.RecursiveMap.entry.2.value.NoRecurse=bar&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService14ProtocolTestIdempotencyTokenAutoFillCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService14ProtocolTest(cfg)
	input := &InputService14TestShapeInputService14TestCaseOperation1Input{
		Token: aws.String("abc123"),
	}

	req := svc.InputService14TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&Token=abc123&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService14ProtocolTestIdempotencyTokenAutoFillCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService14ProtocolTest(cfg)
	input := &InputService14TestShapeInputService14TestCaseOperation2Input{}

	req := svc.InputService14TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&Token=00000000-0000-4000-8000-000000000000&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService15ProtocolTestEnumCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService15ProtocolTest(cfg)
	input := &InputService15TestShapeInputService15TestCaseOperation1Input{
		FooEnum: InputService15TestShapeEnumType("foo"),
		ListEnums: []InputService15TestShapeEnumType{
			InputService15TestShapeEnumType("foo"),
			InputService15TestShapeEnumType(""),
			InputService15TestShapeEnumType("bar"),
		},
	}

	req := svc.InputService15TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&FooEnum=foo&ListEnums.member.1=foo&ListEnums.member.2=&ListEnums.member.3=bar&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService15ProtocolTestEnumCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService15ProtocolTest(cfg)
	input := &InputService15TestShapeInputService15TestCaseOperation2Input{
		FooEnum: InputService15TestShapeEnumType("foo"),
	}

	req := svc.InputService15TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&FooEnum=foo&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService15ProtocolTestEnumCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService15ProtocolTest(cfg)
	input := &InputService15TestShapeInputService15TestCaseOperation3Input{}

	req := svc.InputService15TestCaseOperation3Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestEndpointHostTraitCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation1Input{
		Name: aws.String("myname"),
	}

	req := svc.InputService16TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=StaticOp&Name=myname&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://data-service.region.amazonaws.com/", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestEndpointHostTraitCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation2Input{
		Name: aws.String("myname"),
	}

	req := svc.InputService16TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=MemberRefOp&Name=myname&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://foo-myname.service.region.amazonaws.com/", r.URL.String())

	// assert headers

}
