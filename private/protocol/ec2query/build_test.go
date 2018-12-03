package ec2query_test

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
	"github.com/aws/aws-sdk-go-v2/private/protocol/ec2query"
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

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
		Name: opInputService1TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation1Input{}
	}

	output := &InputService1TestShapeInputService1TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService1TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService1TestCaseOperation1Request}
}

type InputService1TestShapeInputService1TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Bar *string `type:"string"`

	Foo *string `type:"string"`
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

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
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService2TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService2TestCaseOperation1Request}
}

type InputService2TestShapeInputService2TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Bar *string `locationName:"barLocationName" type:"string"`

	Foo *string `type:"string"`

	Yuck *string `locationName:"yuckLocationName" queryName:"yuckQueryName" type:"string"`
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

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
	Input *InputService3TestShapeInputService3TestCaseOperation1Input
	Copy  func(*InputService3TestShapeInputService3TestCaseOperation1Input) InputService3TestCaseOperation1Request
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
func (c *InputService3ProtocolTest) InputService3TestCaseOperation1Request(input *InputService3TestShapeInputService3TestCaseOperation1Input) InputService3TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService3TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService3TestShapeInputService3TestCaseOperation1Input{}
	}

	output := &InputService3TestShapeInputService3TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService3TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService3TestCaseOperation1Request}
}

type InputService3TestShapeInputService3TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	StructArg *InputService3TestShapeStructType `locationName:"Struct" type:"structure"`
}

type InputService3TestShapeInputService3TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService3TestShapeInputService3TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService3TestShapeStructType struct {
	_ struct{} `type:"structure"`

	ScalarArg *string `locationName:"Scalar" type:"string"`
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

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
		Name: opInputService4TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService4TestShapeInputService4TestCaseOperation1Input{}
	}

	output := &InputService4TestShapeInputService4TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService4TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService4TestCaseOperation1Request}
}

type InputService4TestShapeInputService4TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListBools []bool `type:"list"`

	ListFloats []float64 `type:"list"`

	ListIntegers []int64 `type:"list"`

	ListStrings []string `type:"list"`
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

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
	Input *InputService5TestShapeInputService5TestCaseOperation1Input
	Copy  func(*InputService5TestShapeInputService5TestCaseOperation1Input) InputService5TestCaseOperation1Request
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
func (c *InputService5ProtocolTest) InputService5TestCaseOperation1Request(input *InputService5TestShapeInputService5TestCaseOperation1Input) InputService5TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService5TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation1Input{}
	}

	output := &InputService5TestShapeInputService5TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService5TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation1Request}
}

type InputService5TestShapeInputService5TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListArg []string `locationName:"ListMemberName" locationNameList:"item" type:"list"`
}

type InputService5TestShapeInputService5TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService5TestShapeInputService5TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

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
		Name: opInputService6TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService6TestShapeInputService6TestCaseOperation1Input{}
	}

	output := &InputService6TestShapeInputService6TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService6TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService6TestCaseOperation1Request}
}

type InputService6TestShapeInputService6TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListArg []string `locationName:"ListMemberName" queryName:"ListQueryName" locationNameList:"item" type:"list"`
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
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

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
	Input *InputService7TestShapeInputService7TestCaseOperation1Input
	Copy  func(*InputService7TestShapeInputService7TestCaseOperation1Input) InputService7TestCaseOperation1Request
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
func (c *InputService7ProtocolTest) InputService7TestCaseOperation1Request(input *InputService7TestShapeInputService7TestCaseOperation1Input) InputService7TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService7TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation1Input{}
	}

	output := &InputService7TestShapeInputService7TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService7TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService7TestCaseOperation1Request}
}

type InputService7TestShapeInputService7TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	// BlobArg is automatically base64 encoded/decoded by the SDK.
	BlobArg []byte `type:"blob"`
}

type InputService7TestShapeInputService7TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService7TestShapeInputService7TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
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
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

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
	Input *InputService8TestShapeInputService8TestCaseOperation1Input
	Copy  func(*InputService8TestShapeInputService8TestCaseOperation1Input) InputService8TestCaseOperation1Request
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
func (c *InputService8ProtocolTest) InputService8TestCaseOperation1Request(input *InputService8TestShapeInputService8TestCaseOperation1Input) InputService8TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService8TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService8TestShapeInputService8TestCaseOperation1Input{}
	}

	output := &InputService8TestShapeInputService8TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService8TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService8TestCaseOperation1Request}
}

type InputService8TestShapeInputService8TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	TimeArg *time.Time `type:"timestamp" timestampFormat:"iso8601"`
}

type InputService8TestShapeInputService8TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService8TestShapeInputService8TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService9ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService9ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService9TestCaseOperation1 = "OperationName"

// InputService9TestCaseOperation1Request is a API request type for the InputService9TestCaseOperation1 API operation.
type InputService9TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService9TestShapeInputService9TestCaseOperation2Input
	Copy  func(*InputService9TestShapeInputService9TestCaseOperation2Input) InputService9TestCaseOperation1Request
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
func (c *InputService9ProtocolTest) InputService9TestCaseOperation1Request(input *InputService9TestShapeInputService9TestCaseOperation2Input) InputService9TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService9TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation2Input{}
	}

	output := &InputService9TestShapeInputService9TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService9TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService9TestCaseOperation1Request}
}

const opInputService9TestCaseOperation2 = "OperationName"

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
		Name: opInputService9TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation2Input{}
	}

	output := &InputService9TestShapeInputService9TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService9TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService9TestCaseOperation2Request}
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

	Token *string `type:"string" idempotencyToken:"true"`
}

type InputService9TestShapeInputService9TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService9TestShapeInputService9TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// InputService10ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService10ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService10ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService10ProtocolTest client with a config.
//
// Example:
//     // Create a InputService10ProtocolTest client from just a config.
//     svc := inputservice10protocoltest.New(myConfig)
func NewInputService10ProtocolTest(config aws.Config) *InputService10ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService10ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice10protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService10ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService10ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService10TestCaseOperation1 = "OperationName"

// InputService10TestCaseOperation1Request is a API request type for the InputService10TestCaseOperation1 API operation.
type InputService10TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService10TestShapeInputService10TestCaseOperation2Input
	Copy  func(*InputService10TestShapeInputService10TestCaseOperation2Input) InputService10TestCaseOperation1Request
}

// Send marshals and sends the InputService10TestCaseOperation1 API request.
func (r InputService10TestCaseOperation1Request) Send() (*InputService10TestShapeInputService10TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService10TestShapeInputService10TestCaseOperation1Output), nil
}

// InputService10TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService10TestCaseOperation1Request method.
//    req := client.InputService10TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService10ProtocolTest) InputService10TestCaseOperation1Request(input *InputService10TestShapeInputService10TestCaseOperation2Input) InputService10TestCaseOperation1Request {
	op := &aws.Operation{
		Name: opInputService10TestCaseOperation1,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService10TestShapeInputService10TestCaseOperation2Input{}
	}

	output := &InputService10TestShapeInputService10TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService10TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService10TestCaseOperation1Request}
}

const opInputService10TestCaseOperation2 = "OperationName"

// InputService10TestCaseOperation2Request is a API request type for the InputService10TestCaseOperation2 API operation.
type InputService10TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService10TestShapeInputService10TestCaseOperation2Input
	Copy  func(*InputService10TestShapeInputService10TestCaseOperation2Input) InputService10TestCaseOperation2Request
}

// Send marshals and sends the InputService10TestCaseOperation2 API request.
func (r InputService10TestCaseOperation2Request) Send() (*InputService10TestShapeInputService10TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService10TestShapeInputService10TestCaseOperation2Output), nil
}

// InputService10TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService10TestCaseOperation2Request method.
//    req := client.InputService10TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService10ProtocolTest) InputService10TestCaseOperation2Request(input *InputService10TestShapeInputService10TestCaseOperation2Input) InputService10TestCaseOperation2Request {
	op := &aws.Operation{
		Name: opInputService10TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService10TestShapeInputService10TestCaseOperation2Input{}
	}

	output := &InputService10TestShapeInputService10TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService10TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService10TestCaseOperation2Request}
}

type InputService10TestShapeInputService10TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService10TestShapeInputService10TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService10TestShapeInputService10TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService10TestShapeEnumType `type:"string" enum:"true"`

	ListEnums []InputService10TestShapeEnumType `type:"list"`
}

type InputService10TestShapeInputService10TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService10TestShapeInputService10TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService10TestShapeEnumType string

// Enum values for InputService10TestShapeEnumType
const (
	EnumTypeFoo InputService10TestShapeEnumType = "foo"
	EnumTypeBar InputService10TestShapeEnumType = "bar"
)

func (enum InputService10TestShapeEnumType) MarshalValue() (string, error) {
	return string(enum), nil
}

func (enum InputService10TestShapeEnumType) MarshalValueBuf(b []byte) ([]byte, error) {
	b = b[0:0]
	return append(b, enum...), nil
}

// InputService11ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService11ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService11ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService11ProtocolTest client with a config.
//
// Example:
//     // Create a InputService11ProtocolTest client from just a config.
//     svc := inputservice11protocoltest.New(myConfig)
func NewInputService11ProtocolTest(config aws.Config) *InputService11ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService11ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice11protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(ec2query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(ec2query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(ec2query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(ec2query.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService11ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService11ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService11TestCaseOperation1 = "StaticOp"

// InputService11TestCaseOperation1Request is a API request type for the InputService11TestCaseOperation1 API operation.
type InputService11TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService11TestShapeInputService11TestCaseOperation1Input
	Copy  func(*InputService11TestShapeInputService11TestCaseOperation1Input) InputService11TestCaseOperation1Request
}

// Send marshals and sends the InputService11TestCaseOperation1 API request.
func (r InputService11TestCaseOperation1Request) Send() (*InputService11TestShapeInputService11TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService11TestShapeInputService11TestCaseOperation1Output), nil
}

// InputService11TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService11TestCaseOperation1Request method.
//    req := client.InputService11TestCaseOperation1Request(params)
//    resp, err := req.Send()
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

	output := &InputService11TestShapeInputService11TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("data-", nil))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)
	return InputService11TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService11TestCaseOperation1Request}
}

const opInputService11TestCaseOperation2 = "MemberRefOp"

// InputService11TestCaseOperation2Request is a API request type for the InputService11TestCaseOperation2 API operation.
type InputService11TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService11TestShapeInputService11TestCaseOperation2Input
	Copy  func(*InputService11TestShapeInputService11TestCaseOperation2Input) InputService11TestCaseOperation2Request
}

// Send marshals and sends the InputService11TestCaseOperation2 API request.
func (r InputService11TestCaseOperation2Request) Send() (*InputService11TestShapeInputService11TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService11TestShapeInputService11TestCaseOperation2Output), nil
}

// InputService11TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService11TestCaseOperation2Request method.
//    req := client.InputService11TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService11ProtocolTest) InputService11TestCaseOperation2Request(input *InputService11TestShapeInputService11TestCaseOperation2Input) InputService11TestCaseOperation2Request {
	op := &aws.Operation{
		Name: opInputService11TestCaseOperation2,

		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService11TestShapeInputService11TestCaseOperation2Input{}
	}

	output := &InputService11TestShapeInputService11TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("foo-{Name}.", input.hostLabels))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)
	return InputService11TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService11TestCaseOperation2Request}
}

type InputService11TestShapeInputService11TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Name *string `type:"string"`
}

type InputService11TestShapeInputService11TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService11TestShapeInputService11TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

type InputService11TestShapeInputService11TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	// Name is a required field
	Name *string `type:"string" required:"true"`
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *InputService11TestShapeInputService11TestCaseOperation2Input) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "InputService11TestShapeInputService11TestCaseOperation2Input"}

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

func (s *InputService11TestShapeInputService11TestCaseOperation2Input) hostLabels() map[string]string {
	return map[string]string{
		"Name": aws.StringValue(s.Name),
	}
}

type InputService11TestShapeInputService11TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService11TestShapeInputService11TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
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

func TestInputService2ProtocolTestStructureWithLocationNameAndQueryNameAppliedToMembersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService2ProtocolTest(cfg)
	input := &InputService2TestShapeInputService2TestCaseOperation1Input{
		Bar:  aws.String("val2"),
		Foo:  aws.String("val1"),
		Yuck: aws.String("val3"),
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
	awstesting.AssertQuery(t, `Action=OperationName&BarLocationName=val2&Foo=val1&Version=2014-01-01&yuckQueryName=val3`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService3ProtocolTestNestedStructureMembersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService3ProtocolTest(cfg)
	input := &InputService3TestShapeInputService3TestCaseOperation1Input{
		StructArg: &InputService3TestShapeStructType{
			ScalarArg: aws.String("foo"),
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
	awstesting.AssertQuery(t, `Action=OperationName&Struct.Scalar=foo&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService4ProtocolTestListTypesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService4ProtocolTest(cfg)
	input := &InputService4TestShapeInputService4TestCaseOperation1Input{
		ListBools: []bool{
			true,
			false,
			false,
		},
		ListFloats: []float64{
			1.1,
			2.718,
			3.14,
		},
		ListIntegers: []int64{
			0,
			1,
			2,
		},
		ListStrings: []string{
			"foo",
			"bar",
			"baz",
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
	awstesting.AssertQuery(t, `Action=OperationName&ListBools.1=true&ListBools.2=false&ListBools.3=false&ListFloats.1=1.1&ListFloats.2=2.718&ListFloats.3=3.14&ListIntegers.1=0&ListIntegers.2=1&ListIntegers.3=2&ListStrings.1=foo&ListStrings.2=bar&ListStrings.3=baz&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService5ProtocolTestListWithLocationNameAppliedToMemberCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation1Input{
		ListArg: []string{
			"a",
			"b",
			"c",
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
	awstesting.AssertQuery(t, `Action=OperationName&ListMemberName.1=a&ListMemberName.2=b&ListMemberName.3=c&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService6ProtocolTestListWithLocationNameAndQueryNameCase1(t *testing.T) {
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
	awstesting.AssertQuery(t, `Action=OperationName&ListQueryName.1=a&ListQueryName.2=b&ListQueryName.3=c&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService7ProtocolTestBase64EncodedBlobsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService7ProtocolTest(cfg)
	input := &InputService7TestShapeInputService7TestCaseOperation1Input{
		BlobArg: []byte("foo"),
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
	awstesting.AssertQuery(t, `Action=OperationName&BlobArg=Zm9v&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService8ProtocolTestTimestampValuesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService8ProtocolTest(cfg)
	input := &InputService8TestShapeInputService8TestCaseOperation1Input{
		TimeArg: aws.Time(time.Unix(1422172800, 0)),
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
	awstesting.AssertQuery(t, `Action=OperationName&TimeArg=2015-01-25T08%3A00%3A00Z&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService9ProtocolTestIdempotencyTokenAutoFillCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService9ProtocolTest(cfg)
	input := &InputService9TestShapeInputService9TestCaseOperation2Input{
		Token: aws.String("abc123"),
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
	awstesting.AssertQuery(t, `Action=OperationName&Token=abc123&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService9ProtocolTestIdempotencyTokenAutoFillCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService9ProtocolTest(cfg)
	input := &InputService9TestShapeInputService9TestCaseOperation2Input{}

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
	awstesting.AssertQuery(t, `Action=OperationName&Token=00000000-0000-4000-8000-000000000000&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService10ProtocolTestEnumCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService10ProtocolTest(cfg)
	input := &InputService10TestShapeInputService10TestCaseOperation2Input{
		ListEnums: []InputService10TestShapeEnumType{
			InputService10TestShapeEnumType("foo"),
			InputService10TestShapeEnumType(""),
			InputService10TestShapeEnumType("bar"),
		},
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
	awstesting.AssertQuery(t, `Action=OperationName&ListEnums.1=foo&ListEnums.2=&ListEnums.3=bar&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService10ProtocolTestEnumCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService10ProtocolTest(cfg)
	input := &InputService10TestShapeInputService10TestCaseOperation2Input{}

	req := svc.InputService10TestCaseOperation2Request(input)
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

func TestInputService11ProtocolTestEndpointHostTraitCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService11ProtocolTest(cfg)
	input := &InputService11TestShapeInputService11TestCaseOperation1Input{
		Name: aws.String("myname"),
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
	awstesting.AssertQuery(t, `Action=StaticOp&Name=myname&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://data-service.region.amazonaws.com/", r.URL.String())

	// assert headers

}

func TestInputService11ProtocolTestEndpointHostTraitCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService11ProtocolTest(cfg)
	input := &InputService11TestShapeInputService11TestCaseOperation2Input{
		Name: aws.String("myname"),
	}

	req := svc.InputService11TestCaseOperation2Request(input)
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
