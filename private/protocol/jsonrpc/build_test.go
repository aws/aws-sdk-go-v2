package jsonrpc_test

import (
	"bytes"
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
	"github.com/aws/aws-sdk-go-v2/private/protocol/jsonrpc"
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

// InputService1ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService1ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService1ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService1ProtocolTest client with a config.
//
// Example:
//     // Create a InputService1ProtocolTest client from just a config.
//     svc := inputservice1protocoltest.New(myConfig)
func NewInputService1ProtocolTest(config aws.Config) *InputService1ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService1ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice1protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService1ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService1ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService1TestCaseOperation1 = "OperationName"

// InputService1TestCaseOperation1Request is a API request type for the InputService1TestCaseOperation1 API operation.
type InputService1TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService1TestShapeInputService1TestCaseOperation1Input
	Copy  func(*InputService1TestShapeInputService1TestCaseOperation1Input) InputService1TestCaseOperation1Request
}

// Send marshals and sends the InputService1TestCaseOperation1 API request.
func (r InputService1TestCaseOperation1Request) Send() (*InputService1TestShapeInputService1TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService1TestShapeInputService1TestCaseOperation1Output), nil
}

// InputService1TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService1TestCaseOperation1Request method.
//    req := client.InputService1TestCaseOperation1Request(params)
//    resp, err := req.Send()
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

	output := &InputService1TestShapeInputService1TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService1TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService1TestCaseOperation1Request}
}

type InputService1TestShapeInputService1TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Name *string `type:"string"`
}

type InputService1TestShapeInputService1TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService1TestShapeInputService1TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// InputService2ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService2ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService2ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService2ProtocolTest client with a config.
//
// Example:
//     // Create a InputService2ProtocolTest client from just a config.
//     svc := inputservice2protocoltest.New(myConfig)
func NewInputService2ProtocolTest(config aws.Config) *InputService2ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService2ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice2protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService2ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService2ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService2TestCaseOperation1 = "OperationName"

// InputService2TestCaseOperation1Request is a API request type for the InputService2TestCaseOperation1 API operation.
type InputService2TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService2TestShapeInputService2TestCaseOperation1Input
	Copy  func(*InputService2TestShapeInputService2TestCaseOperation1Input) InputService2TestCaseOperation1Request
}

// Send marshals and sends the InputService2TestCaseOperation1 API request.
func (r InputService2TestCaseOperation1Request) Send() (*InputService2TestShapeInputService2TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService2TestShapeInputService2TestCaseOperation1Output), nil
}

// InputService2TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService2TestCaseOperation1Request method.
//    req := client.InputService2TestCaseOperation1Request(params)
//    resp, err := req.Send()
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

	output := &InputService2TestShapeInputService2TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService2TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService2TestCaseOperation1Request}
}

type InputService2TestShapeInputService2TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	TimeArg *time.Time `type:"timestamp" timestampFormat:"unix"`
}

type InputService2TestShapeInputService2TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService2TestShapeInputService2TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// InputService3ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService3ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService3ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService3ProtocolTest client with a config.
//
// Example:
//     // Create a InputService3ProtocolTest client from just a config.
//     svc := inputservice3protocoltest.New(myConfig)
func NewInputService3ProtocolTest(config aws.Config) *InputService3ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService3ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice3protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService3ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService3ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService3TestCaseOperation1 = "OperationName"

// InputService3TestCaseOperation1Request is a API request type for the InputService3TestCaseOperation1 API operation.
type InputService3TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService3TestShapeInputService3TestCaseOperation2Input
	Copy  func(*InputService3TestShapeInputService3TestCaseOperation2Input) InputService3TestCaseOperation1Request
}

// Send marshals and sends the InputService3TestCaseOperation1 API request.
func (r InputService3TestCaseOperation1Request) Send() (*InputService3TestShapeInputService3TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService3TestShapeInputService3TestCaseOperation1Output), nil
}

// InputService3TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService3TestCaseOperation1Request method.
//    req := client.InputService3TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService3ProtocolTest) InputService3TestCaseOperation1Request(input *InputService3TestShapeInputService3TestCaseOperation2Input) InputService3TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService3TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService3TestShapeInputService3TestCaseOperation2Input{}
	}

	output := &InputService3TestShapeInputService3TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService3TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService3TestCaseOperation1Request}
}

const opInputService3TestCaseOperation2 = "OperationName"

// InputService3TestCaseOperation2Request is a API request type for the InputService3TestCaseOperation2 API operation.
type InputService3TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService3TestShapeInputService3TestCaseOperation2Input
	Copy  func(*InputService3TestShapeInputService3TestCaseOperation2Input) InputService3TestCaseOperation2Request
}

// Send marshals and sends the InputService3TestCaseOperation2 API request.
func (r InputService3TestCaseOperation2Request) Send() (*InputService3TestShapeInputService3TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService3TestShapeInputService3TestCaseOperation2Output), nil
}

// InputService3TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService3TestCaseOperation2Request method.
//    req := client.InputService3TestCaseOperation2Request(params)
//    resp, err := req.Send()
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

	output := &InputService3TestShapeInputService3TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService3TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService3TestCaseOperation2Request}
}

type InputService3TestShapeInputService3TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService3TestShapeInputService3TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService3TestShapeInputService3TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	// BlobArg is automatically base64 encoded/decoded by the SDK.
	BlobArg []byte `type:"blob"`

	BlobMap map[string][]byte `type:"map"`
}

type InputService3TestShapeInputService3TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService3TestShapeInputService3TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// InputService4ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService4ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService4ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService4ProtocolTest client with a config.
//
// Example:
//     // Create a InputService4ProtocolTest client from just a config.
//     svc := inputservice4protocoltest.New(myConfig)
func NewInputService4ProtocolTest(config aws.Config) *InputService4ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService4ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice4protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService4ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService4ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService4TestCaseOperation1 = "OperationName"

// InputService4TestCaseOperation1Request is a API request type for the InputService4TestCaseOperation1 API operation.
type InputService4TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService4TestShapeInputService4TestCaseOperation1Input
	Copy  func(*InputService4TestShapeInputService4TestCaseOperation1Input) InputService4TestCaseOperation1Request
}

// Send marshals and sends the InputService4TestCaseOperation1 API request.
func (r InputService4TestCaseOperation1Request) Send() (*InputService4TestShapeInputService4TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService4TestShapeInputService4TestCaseOperation1Output), nil
}

// InputService4TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService4TestCaseOperation1Request method.
//    req := client.InputService4TestCaseOperation1Request(params)
//    resp, err := req.Send()
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

	output := &InputService4TestShapeInputService4TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService4TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService4TestCaseOperation1Request}
}

type InputService4TestShapeInputService4TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListParam [][]byte `type:"list"`
}

type InputService4TestShapeInputService4TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService4TestShapeInputService4TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// InputService5ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService5ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService5ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService5ProtocolTest client with a config.
//
// Example:
//     // Create a InputService5ProtocolTest client from just a config.
//     svc := inputservice5protocoltest.New(myConfig)
func NewInputService5ProtocolTest(config aws.Config) *InputService5ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService5ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice5protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService5ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService5ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService5TestCaseOperation1 = "OperationName"

// InputService5TestCaseOperation1Request is a API request type for the InputService5TestCaseOperation1 API operation.
type InputService5TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation6Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation1Request
}

// Send marshals and sends the InputService5TestCaseOperation1 API request.
func (r InputService5TestCaseOperation1Request) Send() (*InputService5TestShapeInputService5TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation1Output), nil
}

// InputService5TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService5TestCaseOperation1Request method.
//    req := client.InputService5TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation1Request(input *InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation6Input{}
	}

	output := &InputService5TestShapeInputService5TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService5TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation1Request}
}

const opInputService5TestCaseOperation2 = "OperationName"

// InputService5TestCaseOperation2Request is a API request type for the InputService5TestCaseOperation2 API operation.
type InputService5TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation6Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation2Request
}

// Send marshals and sends the InputService5TestCaseOperation2 API request.
func (r InputService5TestCaseOperation2Request) Send() (*InputService5TestShapeInputService5TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation2Output), nil
}

// InputService5TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService5TestCaseOperation2Request method.
//    req := client.InputService5TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation2Request(input *InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation2Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation6Input{}
	}

	output := &InputService5TestShapeInputService5TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService5TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation2Request}
}

const opInputService5TestCaseOperation3 = "OperationName"

// InputService5TestCaseOperation3Request is a API request type for the InputService5TestCaseOperation3 API operation.
type InputService5TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation6Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation3Request
}

// Send marshals and sends the InputService5TestCaseOperation3 API request.
func (r InputService5TestCaseOperation3Request) Send() (*InputService5TestShapeInputService5TestCaseOperation3Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation3Output), nil
}

// InputService5TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService5TestCaseOperation3Request method.
//    req := client.InputService5TestCaseOperation3Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation3Request(input *InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation3Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation3,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation6Input{}
	}

	output := &InputService5TestShapeInputService5TestCaseOperation3Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService5TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation3Request}
}

const opInputService5TestCaseOperation4 = "OperationName"

// InputService5TestCaseOperation4Request is a API request type for the InputService5TestCaseOperation4 API operation.
type InputService5TestCaseOperation4Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation6Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation4Request
}

// Send marshals and sends the InputService5TestCaseOperation4 API request.
func (r InputService5TestCaseOperation4Request) Send() (*InputService5TestShapeInputService5TestCaseOperation4Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation4Output), nil
}

// InputService5TestCaseOperation4Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService5TestCaseOperation4Request method.
//    req := client.InputService5TestCaseOperation4Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation4Request(input *InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation4Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation4,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation6Input{}
	}

	output := &InputService5TestShapeInputService5TestCaseOperation4Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService5TestCaseOperation4Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation4Request}
}

const opInputService5TestCaseOperation5 = "OperationName"

// InputService5TestCaseOperation5Request is a API request type for the InputService5TestCaseOperation5 API operation.
type InputService5TestCaseOperation5Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation6Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation5Request
}

// Send marshals and sends the InputService5TestCaseOperation5 API request.
func (r InputService5TestCaseOperation5Request) Send() (*InputService5TestShapeInputService5TestCaseOperation5Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation5Output), nil
}

// InputService5TestCaseOperation5Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService5TestCaseOperation5Request method.
//    req := client.InputService5TestCaseOperation5Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation5Request(input *InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation5Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation5,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation6Input{}
	}

	output := &InputService5TestShapeInputService5TestCaseOperation5Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService5TestCaseOperation5Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation5Request}
}

const opInputService5TestCaseOperation6 = "OperationName"

// InputService5TestCaseOperation6Request is a API request type for the InputService5TestCaseOperation6 API operation.
type InputService5TestCaseOperation6Request struct {
	*aws.Request
	Input *InputService5TestShapeInputService5TestCaseOperation6Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation6Input) InputService5TestCaseOperation6Request
}

// Send marshals and sends the InputService5TestCaseOperation6 API request.
func (r InputService5TestCaseOperation6Request) Send() (*InputService5TestShapeInputService5TestCaseOperation6Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService5TestShapeInputService5TestCaseOperation6Output), nil
}

// InputService5TestCaseOperation6Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService5TestCaseOperation6Request method.
//    req := client.InputService5TestCaseOperation6Request(params)
//    resp, err := req.Send()
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

	output := &InputService5TestShapeInputService5TestCaseOperation6Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService5TestCaseOperation6Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation6Request}
}

type InputService5TestShapeInputService5TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService5TestShapeInputService5TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService5TestShapeInputService5TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService5TestShapeInputService5TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService5TestShapeInputService5TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService5TestShapeInputService5TestCaseOperation3Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService5TestShapeInputService5TestCaseOperation4Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService5TestShapeInputService5TestCaseOperation4Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService5TestShapeInputService5TestCaseOperation5Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService5TestShapeInputService5TestCaseOperation5Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService5TestShapeInputService5TestCaseOperation6Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService5TestShapeRecursiveStructType `type:"structure"`
}

type InputService5TestShapeInputService5TestCaseOperation6Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService5TestShapeInputService5TestCaseOperation6Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService5TestShapeRecursiveStructType struct {
	_ struct{} `type:"structure"`

	NoRecurse *string `type:"string"`

	RecursiveList []InputService5TestShapeRecursiveStructType `type:"list"`

	RecursiveMap map[string]InputService5TestShapeRecursiveStructType `type:"map"`

	RecursiveStruct *InputService5TestShapeRecursiveStructType `type:"structure"`
}

// InputService6ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService6ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService6ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService6ProtocolTest client with a config.
//
// Example:
//     // Create a InputService6ProtocolTest client from just a config.
//     svc := inputservice6protocoltest.New(myConfig)
func NewInputService6ProtocolTest(config aws.Config) *InputService6ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService6ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice6protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService6ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService6ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService6TestCaseOperation1 = "OperationName"

// InputService6TestCaseOperation1Request is a API request type for the InputService6TestCaseOperation1 API operation.
type InputService6TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService6TestShapeInputService6TestCaseOperation1Input
	Copy  func(*InputService6TestShapeInputService6TestCaseOperation1Input) InputService6TestCaseOperation1Request
}

// Send marshals and sends the InputService6TestCaseOperation1 API request.
func (r InputService6TestCaseOperation1Request) Send() (*InputService6TestShapeInputService6TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService6TestShapeInputService6TestCaseOperation1Output), nil
}

// InputService6TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService6TestCaseOperation1Request method.
//    req := client.InputService6TestCaseOperation1Request(params)
//    resp, err := req.Send()
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

	output := &InputService6TestShapeInputService6TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService6TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService6TestCaseOperation1Request}
}

type InputService6TestShapeInputService6TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Map map[string]string `type:"map"`
}

type InputService6TestShapeInputService6TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService6TestShapeInputService6TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// InputService7ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService7ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService7ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService7ProtocolTest client with a config.
//
// Example:
//     // Create a InputService7ProtocolTest client from just a config.
//     svc := inputservice7protocoltest.New(myConfig)
func NewInputService7ProtocolTest(config aws.Config) *InputService7ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService7ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice7protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService7ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService7ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService7TestCaseOperation1 = "OperationName"

// InputService7TestCaseOperation1Request is a API request type for the InputService7TestCaseOperation1 API operation.
type InputService7TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService7TestShapeInputService7TestCaseOperation2Input
	Copy  func(*InputService7TestShapeInputService7TestCaseOperation2Input) InputService7TestCaseOperation1Request
}

// Send marshals and sends the InputService7TestCaseOperation1 API request.
func (r InputService7TestCaseOperation1Request) Send() (*InputService7TestShapeInputService7TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService7TestShapeInputService7TestCaseOperation1Output), nil
}

// InputService7TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService7TestCaseOperation1Request method.
//    req := client.InputService7TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService7ProtocolTest) InputService7TestCaseOperation1Request(input *InputService7TestShapeInputService7TestCaseOperation2Input) InputService7TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService7TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation2Input{}
	}

	output := &InputService7TestShapeInputService7TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService7TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService7TestCaseOperation1Request}
}

const opInputService7TestCaseOperation2 = "OperationName"

// InputService7TestCaseOperation2Request is a API request type for the InputService7TestCaseOperation2 API operation.
type InputService7TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService7TestShapeInputService7TestCaseOperation2Input
	Copy  func(*InputService7TestShapeInputService7TestCaseOperation2Input) InputService7TestCaseOperation2Request
}

// Send marshals and sends the InputService7TestCaseOperation2 API request.
func (r InputService7TestCaseOperation2Request) Send() (*InputService7TestShapeInputService7TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService7TestShapeInputService7TestCaseOperation2Output), nil
}

// InputService7TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService7TestCaseOperation2Request method.
//    req := client.InputService7TestCaseOperation2Request(params)
//    resp, err := req.Send()
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

	output := &InputService7TestShapeInputService7TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService7TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService7TestCaseOperation2Request}
}

type InputService7TestShapeInputService7TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService7TestShapeInputService7TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService7TestShapeInputService7TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	Token *string `type:"string" idempotencyToken:"true"`
}

type InputService7TestShapeInputService7TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService7TestShapeInputService7TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// InputService8ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService8ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService8ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService8ProtocolTest client with a config.
//
// Example:
//     // Create a InputService8ProtocolTest client from just a config.
//     svc := inputservice8protocoltest.New(myConfig)
func NewInputService8ProtocolTest(config aws.Config) *InputService8ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService8ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice8protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService8ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService8ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService8TestCaseOperation1 = "OperationName"

// InputService8TestCaseOperation1Request is a API request type for the InputService8TestCaseOperation1 API operation.
type InputService8TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService8TestShapeInputService8TestCaseOperation2Input
	Copy  func(*InputService8TestShapeInputService8TestCaseOperation2Input) InputService8TestCaseOperation1Request
}

// Send marshals and sends the InputService8TestCaseOperation1 API request.
func (r InputService8TestCaseOperation1Request) Send() (*InputService8TestShapeInputService8TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService8TestShapeInputService8TestCaseOperation1Output), nil
}

// InputService8TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService8TestCaseOperation1Request method.
//    req := client.InputService8TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService8ProtocolTest) InputService8TestCaseOperation1Request(input *InputService8TestShapeInputService8TestCaseOperation2Input) InputService8TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService8TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService8TestShapeInputService8TestCaseOperation2Input{}
	}

	output := &InputService8TestShapeInputService8TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService8TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService8TestCaseOperation1Request}
}

const opInputService8TestCaseOperation2 = "OperationName"

// InputService8TestCaseOperation2Request is a API request type for the InputService8TestCaseOperation2 API operation.
type InputService8TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService8TestShapeInputService8TestCaseOperation2Input
	Copy  func(*InputService8TestShapeInputService8TestCaseOperation2Input) InputService8TestCaseOperation2Request
}

// Send marshals and sends the InputService8TestCaseOperation2 API request.
func (r InputService8TestCaseOperation2Request) Send() (*InputService8TestShapeInputService8TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService8TestShapeInputService8TestCaseOperation2Output), nil
}

// InputService8TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService8TestCaseOperation2Request method.
//    req := client.InputService8TestCaseOperation2Request(params)
//    resp, err := req.Send()
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

	output := &InputService8TestShapeInputService8TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService8TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService8TestCaseOperation2Request}
}

type InputService8TestShapeInputService8TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService8TestShapeInputService8TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService8TestShapeInputService8TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService8TestShapeEnumType `type:"string" enum:"true"`

	ListEnums []InputService8TestShapeEnumType `type:"list"`
}

type InputService8TestShapeInputService8TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService8TestShapeInputService8TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
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
// . See this package's package overview docs
// for details on the service.
//
// InputService9ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService9ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService9ProtocolTest client with a config.
//
// Example:
//     // Create a InputService9ProtocolTest client from just a config.
//     svc := inputservice9protocoltest.New(myConfig)
func NewInputService9ProtocolTest(config aws.Config) *InputService9ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService9ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice9protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.foo",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService9ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService9ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService9TestCaseOperation1 = "StaticOp"

// InputService9TestCaseOperation1Request is a API request type for the InputService9TestCaseOperation1 API operation.
type InputService9TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService9TestShapeInputService9TestCaseOperation1Input
	Copy  func(*InputService9TestShapeInputService9TestCaseOperation1Input) InputService9TestCaseOperation1Request
}

// Send marshals and sends the InputService9TestCaseOperation1 API request.
func (r InputService9TestCaseOperation1Request) Send() (*InputService9TestShapeInputService9TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService9TestShapeInputService9TestCaseOperation1Output), nil
}

// InputService9TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService9TestCaseOperation1Request method.
//    req := client.InputService9TestCaseOperation1Request(params)
//    resp, err := req.Send()
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

	output := &InputService9TestShapeInputService9TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("data-", nil))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)
	return InputService9TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService9TestCaseOperation1Request}
}

const opInputService9TestCaseOperation2 = "MemberRefOp"

// InputService9TestCaseOperation2Request is a API request type for the InputService9TestCaseOperation2 API operation.
type InputService9TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService9TestShapeInputService9TestCaseOperation2Input
	Copy  func(*InputService9TestShapeInputService9TestCaseOperation2Input) InputService9TestCaseOperation2Request
}

// Send marshals and sends the InputService9TestCaseOperation2 API request.
func (r InputService9TestCaseOperation2Request) Send() (*InputService9TestShapeInputService9TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService9TestShapeInputService9TestCaseOperation2Output), nil
}

// InputService9TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService9TestCaseOperation2Request method.
//    req := client.InputService9TestCaseOperation2Request(params)
//    resp, err := req.Send()
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

	output := &InputService9TestShapeInputService9TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(jsonrpc.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("foo-{Name}.", input.hostLabels))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)
	return InputService9TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService9TestCaseOperation2Request}
}

type InputService9TestShapeInputService9TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Name *string `type:"string"`
}

type InputService9TestShapeInputService9TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService9TestShapeInputService9TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
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

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService9TestShapeInputService9TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
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
	awstesting.AssertJSON(t, `{"Name":"myname"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService2ProtocolTestTimestampValuesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService2ProtocolTest(cfg)
	input := &InputService2TestShapeInputService2TestCaseOperation1Input{
		TimeArg: aws.Time(time.Unix(1422172800, 0)),
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
	awstesting.AssertJSON(t, `{"TimeArg":1422172800}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService3ProtocolTestBase64EncodedBlobsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService3ProtocolTest(cfg)
	input := &InputService3TestShapeInputService3TestCaseOperation2Input{
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
	awstesting.AssertJSON(t, `{"BlobArg":"Zm9v"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
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
	awstesting.AssertJSON(t, `{"BlobMap":{"key1":"Zm9v","key2":"YmFy"}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
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
	awstesting.AssertJSON(t, `{"ListParam":["Zm9v","YmFy"]}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation6Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct":{"NoRecurse":"foo"}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation6Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct":{"RecursiveStruct":{"NoRecurse":"foo"}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation6Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct":{"RecursiveStruct":{"RecursiveStruct":{"RecursiveStruct":{"NoRecurse":"foo"}}}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase4(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation6Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct":{"RecursiveList":[{"NoRecurse":"foo"},{"NoRecurse":"bar"}]}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService5ProtocolTestRecursiveShapesCase5(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation6Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct":{"RecursiveList":[{"NoRecurse":"foo"},{"RecursiveStruct":{"NoRecurse":"bar"}}]}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
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
	awstesting.AssertJSON(t, `{"RecursiveStruct":{"RecursiveMap":{"foo":{"NoRecurse":"foo"},"bar":{"NoRecurse":"bar"}}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
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
	awstesting.AssertJSON(t, `{"Map":{}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.OperationName", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService7ProtocolTestIdempotencyTokenAutoFillCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService7ProtocolTest(cfg)
	input := &InputService7TestShapeInputService7TestCaseOperation2Input{
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
	awstesting.AssertJSON(t, `{"Token":"abc123"}`, util.Trim(string(body)))

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
	awstesting.AssertJSON(t, `{"Token":"00000000-0000-4000-8000-000000000000"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService8ProtocolTestEnumCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService8ProtocolTest(cfg)
	input := &InputService8TestShapeInputService8TestCaseOperation2Input{
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
	awstesting.AssertJSON(t, `{"FooEnum":"foo","ListEnums":["foo","","bar"]}`, util.Trim(string(body)))

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
	awstesting.AssertJSON(t, `{"Name":"myname"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://data-service.region.amazonaws.com/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.StaticOp", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
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
	awstesting.AssertJSON(t, `{"Name":"myname"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://foo-myname.service.region.amazonaws.com/", r.URL.String())

	// assert headers
	if e, a := "application/x-amz-json-1.1", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
	if e, a := "com.amazonaws.foo.MemberRefOp", r.Header.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}
