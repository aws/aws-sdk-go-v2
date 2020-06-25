package jsonrpc_test

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
	"github.com/jviney/aws-sdk-go-v2/private/protocol/jsonrpc"
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
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

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

	Name *string `type:"string"`
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
		Name:       opInputService1TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService1TestShapeInputService1TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
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
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

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

	TimeArg *time.Time `type:"timestamp"`

	TimeCustom *time.Time `type:"timestamp" timestampFormat:"rfc822"`

	TimeFormat *time.Time `type:"timestamp" timestampFormat:"rfc822"`
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
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
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
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

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

	// BlobArg is automatically base64 encoded/decoded by the SDK.
	BlobArg []byte `type:"blob"`

	BlobMap map[string][]byte `type:"map"`
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
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
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

	// BlobArg is automatically base64 encoded/decoded by the SDK.
	BlobArg []byte `type:"blob"`

	BlobMap map[string][]byte `type:"map"`
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
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
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
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

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

	ListParam [][]byte `type:"list"`
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
		Name:       opInputService4TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService4TestShapeInputService4TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService4TestShapeInputService4TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
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
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

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

	RecursiveStruct *InputService5TestShapeRecursiveStructType `type:"structure"`
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
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
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

type InputService5TestShapeInputService5TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService5TestShapeRecursiveStructType `type:"structure"`
}

type InputService5TestShapeInputService5TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService5TestCaseOperation2 = "OperationName"

// InputService5TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService5TestCaseOperation2Request.
//    req := client.InputService5TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation2Request(input *InputService5TestShapeInputService5TestCaseOperation2Input) InputService5TestCaseOperation2Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService5TestShapeInputService5TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService5TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation2Request}
}

// InputService5TestCaseOperation2Request is the request type for the
// InputService5TestCaseOperation2 API operation.
type InputService5TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation2Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation2Input) InputService5TestCaseOperation2Request
}

// Send marshals and sends the InputService5TestCaseOperation2 API request.
func (r InputService5TestCaseOperation2Request) Send(ctx context.Context) (*InputService5TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService5TestCaseOperation2Response{
		InputService5TestShapeInputService5TestCaseOperation2Output: r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService5TestCaseOperation2Response is the response type for the
// InputService5TestCaseOperation2 API operation.
type InputService5TestCaseOperation2Response struct {
	*InputService5TestShapeInputService5TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService5TestCaseOperation2 request.
func (r *InputService5TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService5TestShapeInputService5TestCaseOperation3Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService5TestShapeRecursiveStructType `type:"structure"`
}

type InputService5TestShapeInputService5TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`
}

const opInputService5TestCaseOperation3 = "OperationName"

// InputService5TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService5TestCaseOperation3Request.
//    req := client.InputService5TestCaseOperation3Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation3Request(input *InputService5TestShapeInputService5TestCaseOperation3Input) InputService5TestCaseOperation3Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation3,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation3Input{}
	}

	req := c.newRequest(op, input, &InputService5TestShapeInputService5TestCaseOperation3Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService5TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation3Request}
}

// InputService5TestCaseOperation3Request is the request type for the
// InputService5TestCaseOperation3 API operation.
type InputService5TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation3Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation3Input) InputService5TestCaseOperation3Request
}

// Send marshals and sends the InputService5TestCaseOperation3 API request.
func (r InputService5TestCaseOperation3Request) Send(ctx context.Context) (*InputService5TestCaseOperation3Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService5TestCaseOperation3Response{
		InputService5TestShapeInputService5TestCaseOperation3Output: r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation3Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService5TestCaseOperation3Response is the response type for the
// InputService5TestCaseOperation3 API operation.
type InputService5TestCaseOperation3Response struct {
	*InputService5TestShapeInputService5TestCaseOperation3Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService5TestCaseOperation3 request.
func (r *InputService5TestCaseOperation3Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService5TestShapeInputService5TestCaseOperation4Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService5TestShapeRecursiveStructType `type:"structure"`
}

type InputService5TestShapeInputService5TestCaseOperation4Output struct {
	_ struct{} `type:"structure"`
}

const opInputService5TestCaseOperation4 = "OperationName"

// InputService5TestCaseOperation4Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService5TestCaseOperation4Request.
//    req := client.InputService5TestCaseOperation4Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation4Request(input *InputService5TestShapeInputService5TestCaseOperation4Input) InputService5TestCaseOperation4Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation4,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation4Input{}
	}

	req := c.newRequest(op, input, &InputService5TestShapeInputService5TestCaseOperation4Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService5TestCaseOperation4Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation4Request}
}

// InputService5TestCaseOperation4Request is the request type for the
// InputService5TestCaseOperation4 API operation.
type InputService5TestCaseOperation4Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation4Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation4Input) InputService5TestCaseOperation4Request
}

// Send marshals and sends the InputService5TestCaseOperation4 API request.
func (r InputService5TestCaseOperation4Request) Send(ctx context.Context) (*InputService5TestCaseOperation4Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService5TestCaseOperation4Response{
		InputService5TestShapeInputService5TestCaseOperation4Output: r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation4Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService5TestCaseOperation4Response is the response type for the
// InputService5TestCaseOperation4 API operation.
type InputService5TestCaseOperation4Response struct {
	*InputService5TestShapeInputService5TestCaseOperation4Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService5TestCaseOperation4 request.
func (r *InputService5TestCaseOperation4Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService5TestShapeInputService5TestCaseOperation5Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService5TestShapeRecursiveStructType `type:"structure"`
}

type InputService5TestShapeInputService5TestCaseOperation5Output struct {
	_ struct{} `type:"structure"`
}

const opInputService5TestCaseOperation5 = "OperationName"

// InputService5TestCaseOperation5Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService5TestCaseOperation5Request.
//    req := client.InputService5TestCaseOperation5Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation5Request(input *InputService5TestShapeInputService5TestCaseOperation5Input) InputService5TestCaseOperation5Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation5,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation5Input{}
	}

	req := c.newRequest(op, input, &InputService5TestShapeInputService5TestCaseOperation5Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService5TestCaseOperation5Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation5Request}
}

// InputService5TestCaseOperation5Request is the request type for the
// InputService5TestCaseOperation5 API operation.
type InputService5TestCaseOperation5Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation5Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation5Input) InputService5TestCaseOperation5Request
}

// Send marshals and sends the InputService5TestCaseOperation5 API request.
func (r InputService5TestCaseOperation5Request) Send(ctx context.Context) (*InputService5TestCaseOperation5Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService5TestCaseOperation5Response{
		InputService5TestShapeInputService5TestCaseOperation5Output: r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation5Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService5TestCaseOperation5Response is the response type for the
// InputService5TestCaseOperation5 API operation.
type InputService5TestCaseOperation5Response struct {
	*InputService5TestShapeInputService5TestCaseOperation5Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService5TestCaseOperation5 request.
func (r *InputService5TestCaseOperation5Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService5TestShapeInputService5TestCaseOperation6Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService5TestShapeRecursiveStructType `type:"structure"`
}

type InputService5TestShapeInputService5TestCaseOperation6Output struct {
	_ struct{} `type:"structure"`
}

const opInputService5TestCaseOperation6 = "OperationName"

// InputService5TestCaseOperation6Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService5TestCaseOperation6Request.
//    req := client.InputService5TestCaseOperation6Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation6Request(input *InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation6Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation6,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation6Input{}
	}

	req := c.newRequest(op, input, &InputService5TestShapeInputService5TestCaseOperation6Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService5TestCaseOperation6Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation6Request}
}

// InputService5TestCaseOperation6Request is the request type for the
// InputService5TestCaseOperation6 API operation.
type InputService5TestCaseOperation6Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation6Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation6Request
}

// Send marshals and sends the InputService5TestCaseOperation6 API request.
func (r InputService5TestCaseOperation6Request) Send(ctx context.Context) (*InputService5TestCaseOperation6Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService5TestCaseOperation6Response{
		InputService5TestShapeInputService5TestCaseOperation6Output: r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation6Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService5TestCaseOperation6Response is the response type for the
// InputService5TestCaseOperation6 API operation.
type InputService5TestCaseOperation6Response struct {
	*InputService5TestShapeInputService5TestCaseOperation6Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService5TestCaseOperation6 request.
func (r *InputService5TestCaseOperation6Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService5TestShapeRecursiveStructType struct {
	_ struct{} `type:"structure"`

	NoRecurse *string `type:"string"`

	RecursiveList []InputService5TestShapeRecursiveStructType `type:"list"`

	RecursiveMap map[string]InputService5TestShapeRecursiveStructType `type:"map"`

	RecursiveStruct *InputService5TestShapeRecursiveStructType `type:"structure"`
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
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

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

	Map map[string]string `type:"map"`
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
		Name:       opInputService6TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService6TestShapeInputService6TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService6TestShapeInputService6TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

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

	Token *string `type:"string" idempotencyToken:"true"`
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
		Name:       opInputService7TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService7TestShapeInputService7TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
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

type InputService7TestShapeInputService7TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	Token *string `type:"string" idempotencyToken:"true"`
}

type InputService7TestShapeInputService7TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService7TestCaseOperation2 = "OperationName"

// InputService7TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService7TestCaseOperation2Request.
//    req := client.InputService7TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService7ProtocolTest) InputService7TestCaseOperation2Request(input *InputService7TestShapeInputService7TestCaseOperation2Input) InputService7TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService7TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService7TestShapeInputService7TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService7TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService7TestCaseOperation2Request}
}

// InputService7TestCaseOperation2Request is the request type for the
// InputService7TestCaseOperation2 API operation.
type InputService7TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService7TestShapeInputService7TestCaseOperation2Input
	Copy  func(*InputService7TestShapeInputService7TestCaseOperation2Input) InputService7TestCaseOperation2Request
}

// Send marshals and sends the InputService7TestCaseOperation2 API request.
func (r InputService7TestCaseOperation2Request) Send(ctx context.Context) (*InputService7TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService7TestCaseOperation2Response{
		InputService7TestShapeInputService7TestCaseOperation2Output: r.Request.Data.(*InputService7TestShapeInputService7TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService7TestCaseOperation2Response is the response type for the
// InputService7TestCaseOperation2 API operation.
type InputService7TestCaseOperation2Response struct {
	*InputService7TestShapeInputService7TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService7TestCaseOperation2 request.
func (r *InputService7TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
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
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

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

	FooEnum InputService8TestShapeEnumType `type:"string" enum:"true"`

	ListEnums []InputService8TestShapeEnumType `type:"list"`
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
		Name:       opInputService8TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService8TestShapeInputService8TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService8TestShapeInputService8TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
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

type InputService8TestShapeInputService8TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService8TestShapeEnumType `type:"string" enum:"true"`

	ListEnums []InputService8TestShapeEnumType `type:"list"`
}

type InputService8TestShapeInputService8TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService8TestCaseOperation2 = "OperationName"

// InputService8TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService8TestCaseOperation2Request.
//    req := client.InputService8TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService8ProtocolTest) InputService8TestCaseOperation2Request(input *InputService8TestShapeInputService8TestCaseOperation2Input) InputService8TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService8TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService8TestShapeInputService8TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService8TestShapeInputService8TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService8TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService8TestCaseOperation2Request}
}

// InputService8TestCaseOperation2Request is the request type for the
// InputService8TestCaseOperation2 API operation.
type InputService8TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService8TestShapeInputService8TestCaseOperation2Input
	Copy  func(*InputService8TestShapeInputService8TestCaseOperation2Input) InputService8TestCaseOperation2Request
}

// Send marshals and sends the InputService8TestCaseOperation2 API request.
func (r InputService8TestCaseOperation2Request) Send(ctx context.Context) (*InputService8TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService8TestCaseOperation2Response{
		InputService8TestShapeInputService8TestCaseOperation2Output: r.Request.Data.(*InputService8TestShapeInputService8TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService8TestCaseOperation2Response is the response type for the
// InputService8TestCaseOperation2 API operation.
type InputService8TestCaseOperation2Response struct {
	*InputService8TestShapeInputService8TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService8TestCaseOperation2 request.
func (r *InputService8TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService8TestShapeEnumType string

// Enum values for InputService8TestShapeEnumType
const (
	EnumTypeFoo InputService8TestShapeEnumType = "foo"
	EnumTypeBar InputService8TestShapeEnumType = "bar"
)

func (enum InputService8TestShapeEnumType) MarshalValue() (string, error) {
	return string(enum), nil
}

func (enum InputService8TestShapeEnumType) MarshalValueBuf(b []byte) ([]byte, error) {
	b = b[0:0]
	return append(b, enum...), nil
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
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

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

	Name *string `type:"string"`
}

type InputService9TestShapeInputService9TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

const opInputService9TestCaseOperation1 = "StaticOp"

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
		Name:       opInputService9TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService9TestShapeInputService9TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("data-", nil))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)

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

type InputService9TestShapeInputService9TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	// Name is a required field
	Name *string `type:"string" required:"true"`
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *InputService9TestShapeInputService9TestCaseOperation2Input) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "InputService9TestShapeInputService9TestCaseOperation2Input"}

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

func (s *InputService9TestShapeInputService9TestCaseOperation2Input) hostLabels() map[string]string {
	return map[string]string{
		"Name": aws.StringValue(s.Name),
	}
}

type InputService9TestShapeInputService9TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

const opInputService9TestCaseOperation2 = "MemberRefOp"

// InputService9TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService9TestCaseOperation2Request.
//    req := client.InputService9TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService9ProtocolTest) InputService9TestCaseOperation2Request(input *InputService9TestShapeInputService9TestCaseOperation2Input) InputService9TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService9TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService9TestShapeInputService9TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("foo-{Name}.", input.hostLabels))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)

	return InputService9TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService9TestCaseOperation2Request}
}

// InputService9TestCaseOperation2Request is the request type for the
// InputService9TestCaseOperation2 API operation.
type InputService9TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService9TestShapeInputService9TestCaseOperation2Input
	Copy  func(*InputService9TestShapeInputService9TestCaseOperation2Input) InputService9TestCaseOperation2Request
}

// Send marshals and sends the InputService9TestCaseOperation2 API request.
func (r InputService9TestCaseOperation2Request) Send(ctx context.Context) (*InputService9TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService9TestCaseOperation2Response{
		InputService9TestShapeInputService9TestCaseOperation2Output: r.Request.Data.(*InputService9TestShapeInputService9TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService9TestCaseOperation2Response is the response type for the
// InputService9TestCaseOperation2 API operation.
type InputService9TestCaseOperation2Response struct {
	*InputService9TestShapeInputService9TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService9TestCaseOperation2 request.
func (r *InputService9TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
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
		Name: aws.String("myname"),
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
	awstesting.AssertJSON(t, `{"Name": "myname"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService2ProtocolTestTimestampValuesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService2ProtocolTest(cfg)
	input := &InputService2TestShapeInputService2TestCaseOperation1Input{
		TimeArg:    aws.Time(time.Unix(1422172800, 0)),
		TimeCustom: aws.Time(time.Unix(1422172800, 0)),
		TimeFormat: aws.Time(time.Unix(1422172800, 0)),
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
	awstesting.AssertJSON(t, `{"TimeArg": 1422172800, "TimeCustom": "Sun, 25 Jan 2015 08:00:00 GMT", "TimeFormat": "Sun, 25 Jan 2015 08:00:00 GMT"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService3ProtocolTestBase64EncodedBlobsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService3ProtocolTest(cfg)
	input := &InputService3TestShapeInputService3TestCaseOperation1Input{
		BlobArg: []byte("foo"),
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
	awstesting.AssertJSON(t, `{"BlobArg": "Zm9v"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService3ProtocolTestBase64EncodedBlobsCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService3ProtocolTest(cfg)
	input := &InputService3TestShapeInputService3TestCaseOperation2Input{
		BlobMap: map[string][]byte{
			"key1": []byte("foo"),
			"key2": []byte("bar"),
		},
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
	awstesting.AssertJSON(t, `{"BlobMap": {"key1": "Zm9v", "key2": "YmFy"}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService4ProtocolTestNestedBlobsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService4ProtocolTest(cfg)
	input := &InputService4TestShapeInputService4TestCaseOperation1Input{
		ListParam: [][]byte{
			[]byte("foo"),
			[]byte("bar"),
		},
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
	awstesting.AssertJSON(t, `{"ListParam": ["Zm9v", "YmFy"]}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation1Input{
		RecursiveStruct: &InputService5TestShapeRecursiveStructType{
			NoRecurse: aws.String("foo"),
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"NoRecurse": "foo"}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation2Input{
		RecursiveStruct: &InputService5TestShapeRecursiveStructType{
			RecursiveStruct: &InputService5TestShapeRecursiveStructType{
				NoRecurse: aws.String("foo"),
			},
		},
	}

	req := svc.InputService5TestCaseOperation2Request(input)
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveStruct": {"NoRecurse": "foo"}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation3Input{
		RecursiveStruct: &InputService5TestShapeRecursiveStructType{
			RecursiveStruct: &InputService5TestShapeRecursiveStructType{
				RecursiveStruct: &InputService5TestShapeRecursiveStructType{
					RecursiveStruct: &InputService5TestShapeRecursiveStructType{
						NoRecurse: aws.String("foo"),
					},
				},
			},
		},
	}

	req := svc.InputService5TestCaseOperation3Request(input)
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveStruct": {"RecursiveStruct": {"RecursiveStruct": {"NoRecurse": "foo"}}}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase4(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation4Input{
		RecursiveStruct: &InputService5TestShapeRecursiveStructType{
			RecursiveList: []InputService5TestShapeRecursiveStructType{
				{
					NoRecurse: aws.String("foo"),
				},
				{
					NoRecurse: aws.String("bar"),
				},
			},
		},
	}

	req := svc.InputService5TestCaseOperation4Request(input)
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveList": [{"NoRecurse": "foo"}, {"NoRecurse": "bar"}]}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase5(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation5Input{
		RecursiveStruct: &InputService5TestShapeRecursiveStructType{
			RecursiveList: []InputService5TestShapeRecursiveStructType{
				{
					NoRecurse: aws.String("foo"),
				},
				{
					RecursiveStruct: &InputService5TestShapeRecursiveStructType{
						NoRecurse: aws.String("bar"),
					},
				},
			},
		},
	}

	req := svc.InputService5TestCaseOperation5Request(input)
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveList": [{"NoRecurse": "foo"}, {"RecursiveStruct": {"NoRecurse": "bar"}}]}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase6(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation6Input{
		RecursiveStruct: &InputService5TestShapeRecursiveStructType{
			RecursiveMap: map[string]InputService5TestShapeRecursiveStructType{
				"bar": {
					NoRecurse: aws.String("bar"),
				},
				"foo": {
					NoRecurse: aws.String("foo"),
				},
			},
		},
	}

	req := svc.InputService5TestCaseOperation6Request(input)
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveMap": {"foo": {"NoRecurse": "foo"}, "bar": {"NoRecurse": "bar"}}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService6ProtocolTestEmptyMapsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService6ProtocolTest(cfg)
	input := &InputService6TestShapeInputService6TestCaseOperation1Input{
		Map: map[string]string{},
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
	awstesting.AssertJSON(t, `{"Map": {}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService7ProtocolTestIdempotencyTokenAutoFillCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService7ProtocolTest(cfg)
	input := &InputService7TestShapeInputService7TestCaseOperation1Input{
		Token: aws.String("abc123"),
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
	awstesting.AssertJSON(t, `{"Token": "abc123"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService7ProtocolTestIdempotencyTokenAutoFillCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService7ProtocolTest(cfg)
	input := &InputService7TestShapeInputService7TestCaseOperation2Input{}

	req := svc.InputService7TestCaseOperation2Request(input)
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
	awstesting.AssertJSON(t, `{"Token": "00000000-0000-4000-8000-000000000000"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService8ProtocolTestEnumCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService8ProtocolTest(cfg)
	input := &InputService8TestShapeInputService8TestCaseOperation1Input{
		FooEnum: InputService8TestShapeEnumType("foo"),
		ListEnums: []InputService8TestShapeEnumType{
			InputService8TestShapeEnumType("foo"),
			InputService8TestShapeEnumType(""),
			InputService8TestShapeEnumType("bar"),
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
	awstesting.AssertJSON(t, `{"FooEnum": "foo", "ListEnums": ["foo", "", "bar"]}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService8ProtocolTestEnumCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService8ProtocolTest(cfg)
	input := &InputService8TestShapeInputService8TestCaseOperation2Input{}

	req := svc.InputService8TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService9ProtocolTestEndpointHostTraitStaticPrefixCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService9ProtocolTest(cfg)
	input := &InputService9TestShapeInputService9TestCaseOperation1Input{
		Name: aws.String("myname"),
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
	awstesting.AssertJSON(t, `{"Name": "myname"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://data-service.region.amazonaws.com/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.StaticOp", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService9ProtocolTestEndpointHostTraitStaticPrefixCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService9ProtocolTest(cfg)
	input := &InputService9TestShapeInputService9TestCaseOperation2Input{
		Name: aws.String("myname"),
	}

	req := svc.InputService9TestCaseOperation2Request(input)
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
	awstesting.AssertJSON(t, `{"Name": "myname"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://foo-myname.service.region.amazonaws.com/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "com.amazonaws.foo.MemberRefOp", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}
