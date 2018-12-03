package restjson_test

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
	"github.com/aws/aws-sdk-go-v2/private/protocol/restjson"
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobs",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation1Input{}
	}

	output := &InputService1TestShapeInputService1TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService1TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService1TestCaseOperation1Request}
}

type InputService1TestShapeInputService1TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService1TestShapeInputService1TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	return nil
}

type InputService1TestShapeInputService1TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService1TestShapeInputService1TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService1TestShapeInputService1TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
		Name:       opInputService2TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService2TestShapeInputService2TestCaseOperation1Input{}
	}

	output := &InputService2TestShapeInputService2TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService2TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService2TestCaseOperation1Request}
}

type InputService2TestShapeInputService2TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	PipelineId *string `location:"uri" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService2TestShapeInputService2TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService2TestShapeInputService2TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService2TestShapeInputService2TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService2TestShapeInputService2TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
		Name:       opInputService3TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService3TestShapeInputService3TestCaseOperation1Input{}
	}

	output := &InputService3TestShapeInputService3TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService3TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService3TestCaseOperation1Request}
}

type InputService3TestShapeInputService3TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Foo *string `location:"uri" locationName:"PipelineId" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService3TestShapeInputService3TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService3TestShapeInputService3TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService3TestShapeInputService3TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService3TestShapeInputService3TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService4TestShapeInputService4TestCaseOperation1Input{}
	}

	output := &InputService4TestShapeInputService4TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService4TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService4TestCaseOperation1Request}
}

type InputService4TestShapeInputService4TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Items []string `location:"querystring" locationName:"item" type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService4TestShapeInputService4TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if len(s.Items) > 0 {
		v := s.Items

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.QueryTarget, "item", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddValue(protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v1)})
		}
		ls0.End()

	}
	return nil
}

type InputService4TestShapeInputService4TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService4TestShapeInputService4TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService4TestShapeInputService4TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
		Name:       opInputService5TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation1Input{}
	}

	output := &InputService5TestShapeInputService5TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService5TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService5TestCaseOperation1Request}
}

type InputService5TestShapeInputService5TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	PipelineId *string `location:"uri" type:"string"`

	QueryDoc map[string]string `location:"querystring" type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService5TestShapeInputService5TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if len(s.QueryDoc) > 0 {
		v := s.QueryDoc

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.QueryTarget, "QueryDoc", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ms0.MapSetValue(k1, protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v1)})
		}
		ms0.End()

	}
	return nil
}

type InputService5TestShapeInputService5TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService5TestShapeInputService5TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService5TestShapeInputService5TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService6TestShapeInputService6TestCaseOperation1Input{}
	}

	output := &InputService6TestShapeInputService6TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService6TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService6TestCaseOperation1Request}
}

type InputService6TestShapeInputService6TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	PipelineId *string `location:"uri" type:"string"`

	QueryDoc map[string][]string `location:"querystring" type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService6TestShapeInputService6TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if len(s.QueryDoc) > 0 {
		v := s.QueryDoc

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.QueryTarget, "QueryDoc", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ls1 := ms0.List(k1)
			ls1.Start()
			for _, v2 := range v1 {
				ls1.ListAddValue(protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v2)})
			}
			ls1.End()
		}
		ms0.End()

	}
	return nil
}

type InputService6TestShapeInputService6TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService6TestShapeInputService6TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService6TestShapeInputService6TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation2Input{}
	}

	output := &InputService7TestShapeInputService7TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation2Input{}
	}

	output := &InputService7TestShapeInputService7TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService7TestShapeInputService7TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService7TestShapeInputService7TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	BoolQuery *bool `location:"querystring" locationName:"bool-query" type:"boolean"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService7TestShapeInputService7TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.BoolQuery != nil {
		v := *s.BoolQuery

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "bool-query", protocol.BoolValue(v), metadata)
	}
	return nil
}

type InputService7TestShapeInputService7TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService7TestShapeInputService7TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService7TestShapeInputService7TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
		Name:       opInputService8TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService8TestShapeInputService8TestCaseOperation1Input{}
	}

	output := &InputService8TestShapeInputService8TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService8TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService8TestCaseOperation1Request}
}

type InputService8TestShapeInputService8TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Ascending *string `location:"querystring" locationName:"Ascending" type:"string"`

	PageToken *string `location:"querystring" locationName:"PageToken" type:"string"`

	PipelineId *string `location:"uri" locationName:"PipelineId" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService8TestShapeInputService8TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.Ascending != nil {
		v := *s.Ascending

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "Ascending", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.PageToken != nil {
		v := *s.PageToken

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "PageToken", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService8TestShapeInputService8TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService8TestShapeInputService8TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService8TestShapeInputService8TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation1Input{}
	}

	output := &InputService9TestShapeInputService9TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService9TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService9TestCaseOperation1Request}
}

type InputService9TestShapeInputService9TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Ascending *string `location:"querystring" locationName:"Ascending" type:"string"`

	Config *InputService9TestShapeStructType `type:"structure"`

	PageToken *string `location:"querystring" locationName:"PageToken" type:"string"`

	PipelineId *string `location:"uri" locationName:"PipelineId" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService9TestShapeInputService9TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Config != nil {
		v := s.Config

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "Config", v, metadata)
	}
	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.Ascending != nil {
		v := *s.Ascending

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "Ascending", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.PageToken != nil {
		v := *s.PageToken

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "PageToken", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService9TestShapeInputService9TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService9TestShapeInputService9TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService9TestShapeInputService9TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService9TestShapeStructType struct {
	_ struct{} `type:"structure"`

	A *string `type:"string"`

	B *string `type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService9TestShapeStructType) MarshalFields(e protocol.FieldEncoder) error {
	if s.A != nil {
		v := *s.A

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "A", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.B != nil {
		v := *s.B

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "B", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
	Input *InputService10TestShapeInputService10TestCaseOperation1Input
	Copy  func(*InputService10TestShapeInputService10TestCaseOperation1Input) InputService10TestCaseOperation1Request
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
func (c *InputService10ProtocolTest) InputService10TestCaseOperation1Request(input *InputService10TestShapeInputService10TestCaseOperation1Input) InputService10TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService10TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService10TestShapeInputService10TestCaseOperation1Input{}
	}

	output := &InputService10TestShapeInputService10TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService10TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService10TestCaseOperation1Request}
}

type InputService10TestShapeInputService10TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Ascending *string `location:"querystring" locationName:"Ascending" type:"string"`

	Checksum *string `location:"header" locationName:"x-amz-checksum" type:"string"`

	Config *InputService10TestShapeStructType `type:"structure"`

	PageToken *string `location:"querystring" locationName:"PageToken" type:"string"`

	PipelineId *string `location:"uri" locationName:"PipelineId" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService10TestShapeInputService10TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Config != nil {
		v := s.Config

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "Config", v, metadata)
	}
	if s.Checksum != nil {
		v := *s.Checksum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "x-amz-checksum", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.Ascending != nil {
		v := *s.Ascending

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "Ascending", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.PageToken != nil {
		v := *s.PageToken

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "PageToken", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService10TestShapeInputService10TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService10TestShapeInputService10TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService10TestShapeInputService10TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService10TestShapeStructType struct {
	_ struct{} `type:"structure"`

	A *string `type:"string"`

	B *string `type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService10TestShapeStructType) MarshalFields(e protocol.FieldEncoder) error {
	if s.A != nil {
		v := *s.A

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "A", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.B != nil {
		v := *s.B

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "B", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService11ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService11ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService11TestCaseOperation1 = "OperationName"

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
		Name:       opInputService11TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/vaults/{vaultName}/archives",
	}

	if input == nil {
		input = &InputService11TestShapeInputService11TestCaseOperation1Input{}
	}

	output := &InputService11TestShapeInputService11TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService11TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService11TestCaseOperation1Request}
}

type InputService11TestShapeInputService11TestCaseOperation1Input struct {
	_ struct{} `type:"structure" payload:"Body"`

	Body io.ReadSeeker `locationName:"body" type:"blob"`

	Checksum *string `location:"header" locationName:"x-amz-sha256-tree-hash" type:"string"`

	// VaultName is a required field
	VaultName *string `location:"uri" locationName:"vaultName" type:"string" required:"true"`
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *InputService11TestShapeInputService11TestCaseOperation1Input) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "InputService11TestShapeInputService11TestCaseOperation1Input"}

	if s.VaultName == nil {
		invalidParams.Add(aws.NewErrParamRequired("VaultName"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService11TestShapeInputService11TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Checksum != nil {
		v := *s.Checksum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "x-amz-sha256-tree-hash", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.VaultName != nil {
		v := *s.VaultName

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "vaultName", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.Body != nil {
		v := s.Body

		metadata := protocol.Metadata{}
		e.SetStream(protocol.PayloadTarget, "body", protocol.ReadSeekerStream{V: v}, metadata)
	}
	return nil
}

type InputService11TestShapeInputService11TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService11TestShapeInputService11TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService11TestShapeInputService11TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService12ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService12ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService12ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService12ProtocolTest client with a config.
//
// Example:
//     // Create a InputService12ProtocolTest client from just a config.
//     svc := inputservice12protocoltest.New(myConfig)
func NewInputService12ProtocolTest(config aws.Config) *InputService12ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService12ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice12protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService12ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService12ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService12TestCaseOperation1 = "OperationName"

// InputService12TestCaseOperation1Request is a API request type for the InputService12TestCaseOperation1 API operation.
type InputService12TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService12TestShapeInputService12TestCaseOperation1Input
	Copy  func(*InputService12TestShapeInputService12TestCaseOperation1Input) InputService12TestCaseOperation1Request
}

// Send marshals and sends the InputService12TestCaseOperation1 API request.
func (r InputService12TestCaseOperation1Request) Send() (*InputService12TestShapeInputService12TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService12TestShapeInputService12TestCaseOperation1Output), nil
}

// InputService12TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService12TestCaseOperation1Request method.
//    req := client.InputService12TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService12ProtocolTest) InputService12TestCaseOperation1Request(input *InputService12TestShapeInputService12TestCaseOperation1Input) InputService12TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService12TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/{Foo}",
	}

	if input == nil {
		input = &InputService12TestShapeInputService12TestCaseOperation1Input{}
	}

	output := &InputService12TestShapeInputService12TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService12TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService12TestCaseOperation1Request}
}

type InputService12TestShapeInputService12TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	// Bar is automatically base64 encoded/decoded by the SDK.
	Bar []byte `type:"blob"`

	// Foo is a required field
	Foo *string `location:"uri" locationName:"Foo" type:"string" required:"true"`
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *InputService12TestShapeInputService12TestCaseOperation1Input) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "InputService12TestShapeInputService12TestCaseOperation1Input"}

	if s.Foo == nil {
		invalidParams.Add(aws.NewErrParamRequired("Foo"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService12TestShapeInputService12TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Bar != nil {
		v := s.Bar

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Bar", protocol.QuotedValue{ValueMarshaler: protocol.BytesValue(v)}, metadata)
	}
	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "Foo", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService12TestShapeInputService12TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService12TestShapeInputService12TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService12TestShapeInputService12TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService13ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService13ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService13ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService13ProtocolTest client with a config.
//
// Example:
//     // Create a InputService13ProtocolTest client from just a config.
//     svc := inputservice13protocoltest.New(myConfig)
func NewInputService13ProtocolTest(config aws.Config) *InputService13ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService13ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice13protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService13ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService13ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService13TestCaseOperation1 = "OperationName"

// InputService13TestCaseOperation1Request is a API request type for the InputService13TestCaseOperation1 API operation.
type InputService13TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService13TestShapeInputService13TestCaseOperation2Input
	Copy  func(*InputService13TestShapeInputService13TestCaseOperation2Input) InputService13TestCaseOperation1Request
}

// Send marshals and sends the InputService13TestCaseOperation1 API request.
func (r InputService13TestCaseOperation1Request) Send() (*InputService13TestShapeInputService13TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService13TestShapeInputService13TestCaseOperation1Output), nil
}

// InputService13TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService13TestCaseOperation1Request method.
//    req := client.InputService13TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService13ProtocolTest) InputService13TestCaseOperation1Request(input *InputService13TestShapeInputService13TestCaseOperation2Input) InputService13TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService13TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation2Input{}
	}

	output := &InputService13TestShapeInputService13TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService13TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService13TestCaseOperation1Request}
}

const opInputService13TestCaseOperation2 = "OperationName"

// InputService13TestCaseOperation2Request is a API request type for the InputService13TestCaseOperation2 API operation.
type InputService13TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService13TestShapeInputService13TestCaseOperation2Input
	Copy  func(*InputService13TestShapeInputService13TestCaseOperation2Input) InputService13TestCaseOperation2Request
}

// Send marshals and sends the InputService13TestCaseOperation2 API request.
func (r InputService13TestCaseOperation2Request) Send() (*InputService13TestShapeInputService13TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService13TestShapeInputService13TestCaseOperation2Output), nil
}

// InputService13TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService13TestCaseOperation2Request method.
//    req := client.InputService13TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService13ProtocolTest) InputService13TestCaseOperation2Request(input *InputService13TestShapeInputService13TestCaseOperation2Input) InputService13TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService13TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation2Input{}
	}

	output := &InputService13TestShapeInputService13TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService13TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService13TestCaseOperation2Request}
}

type InputService13TestShapeInputService13TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService13TestShapeInputService13TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService13TestShapeInputService13TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService13TestShapeInputService13TestCaseOperation2Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo []byte `locationName:"foo" type:"blob"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService13TestShapeInputService13TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetStream(protocol.PayloadTarget, "foo", protocol.BytesStream(v), metadata)
	}
	return nil
}

type InputService13TestShapeInputService13TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService13TestShapeInputService13TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService13TestShapeInputService13TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService14ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService14ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService14ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService14ProtocolTest client with a config.
//
// Example:
//     // Create a InputService14ProtocolTest client from just a config.
//     svc := inputservice14protocoltest.New(myConfig)
func NewInputService14ProtocolTest(config aws.Config) *InputService14ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService14ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice14protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService14ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService14ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService14TestCaseOperation1 = "OperationName"

// InputService14TestCaseOperation1Request is a API request type for the InputService14TestCaseOperation1 API operation.
type InputService14TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService14TestShapeInputService14TestCaseOperation2Input
	Copy  func(*InputService14TestShapeInputService14TestCaseOperation2Input) InputService14TestCaseOperation1Request
}

// Send marshals and sends the InputService14TestCaseOperation1 API request.
func (r InputService14TestCaseOperation1Request) Send() (*InputService14TestShapeInputService14TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService14TestShapeInputService14TestCaseOperation1Output), nil
}

// InputService14TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService14TestCaseOperation1Request method.
//    req := client.InputService14TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService14ProtocolTest) InputService14TestCaseOperation1Request(input *InputService14TestShapeInputService14TestCaseOperation2Input) InputService14TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService14TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService14TestShapeInputService14TestCaseOperation2Input{}
	}

	output := &InputService14TestShapeInputService14TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService14TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService14TestCaseOperation1Request}
}

const opInputService14TestCaseOperation2 = "OperationName"

// InputService14TestCaseOperation2Request is a API request type for the InputService14TestCaseOperation2 API operation.
type InputService14TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService14TestShapeInputService14TestCaseOperation2Input
	Copy  func(*InputService14TestShapeInputService14TestCaseOperation2Input) InputService14TestCaseOperation2Request
}

// Send marshals and sends the InputService14TestCaseOperation2 API request.
func (r InputService14TestCaseOperation2Request) Send() (*InputService14TestShapeInputService14TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService14TestShapeInputService14TestCaseOperation2Output), nil
}

// InputService14TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService14TestCaseOperation2Request method.
//    req := client.InputService14TestCaseOperation2Request(params)
//    resp, err := req.Send()
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

	output := &InputService14TestShapeInputService14TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService14TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService14TestCaseOperation2Request}
}

type InputService14TestShapeFooShape struct {
	_ struct{} `type:"structure"`

	Baz *string `locationName:"baz" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService14TestShapeFooShape) MarshalFields(e protocol.FieldEncoder) error {
	if s.Baz != nil {
		v := *s.Baz

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "baz", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService14TestShapeInputService14TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService14TestShapeInputService14TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService14TestShapeInputService14TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService14TestShapeInputService14TestCaseOperation2Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo *InputService14TestShapeFooShape `locationName:"foo" type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService14TestShapeInputService14TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetFields(protocol.PayloadTarget, "foo", v, metadata)
	}
	return nil
}

type InputService14TestShapeInputService14TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService14TestShapeInputService14TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService14TestShapeInputService14TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService15ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService15ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService15ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService15ProtocolTest client with a config.
//
// Example:
//     // Create a InputService15ProtocolTest client from just a config.
//     svc := inputservice15protocoltest.New(myConfig)
func NewInputService15ProtocolTest(config aws.Config) *InputService15ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService15ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice15protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService15ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService15ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService15TestCaseOperation1 = "OperationName"

// InputService15TestCaseOperation1Request is a API request type for the InputService15TestCaseOperation1 API operation.
type InputService15TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService15TestShapeInputService15TestCaseOperation2Input
	Copy  func(*InputService15TestShapeInputService15TestCaseOperation2Input) InputService15TestCaseOperation1Request
}

// Send marshals and sends the InputService15TestCaseOperation1 API request.
func (r InputService15TestCaseOperation1Request) Send() (*InputService15TestShapeInputService15TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService15TestShapeInputService15TestCaseOperation1Output), nil
}

// InputService15TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService15TestCaseOperation1Request method.
//    req := client.InputService15TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService15ProtocolTest) InputService15TestCaseOperation1Request(input *InputService15TestShapeInputService15TestCaseOperation2Input) InputService15TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService15TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService15TestShapeInputService15TestCaseOperation2Input{}
	}

	output := &InputService15TestShapeInputService15TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService15TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService15TestCaseOperation1Request}
}

const opInputService15TestCaseOperation2 = "OperationName"

// InputService15TestCaseOperation2Request is a API request type for the InputService15TestCaseOperation2 API operation.
type InputService15TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService15TestShapeInputService15TestCaseOperation2Input
	Copy  func(*InputService15TestShapeInputService15TestCaseOperation2Input) InputService15TestCaseOperation2Request
}

// Send marshals and sends the InputService15TestCaseOperation2 API request.
func (r InputService15TestCaseOperation2Request) Send() (*InputService15TestShapeInputService15TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService15TestShapeInputService15TestCaseOperation2Output), nil
}

// InputService15TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService15TestCaseOperation2Request method.
//    req := client.InputService15TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService15ProtocolTest) InputService15TestCaseOperation2Request(input *InputService15TestShapeInputService15TestCaseOperation2Input) InputService15TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService15TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path?abc=mno",
	}

	if input == nil {
		input = &InputService15TestShapeInputService15TestCaseOperation2Input{}
	}

	output := &InputService15TestShapeInputService15TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService15TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService15TestCaseOperation2Request}
}

type InputService15TestShapeInputService15TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService15TestShapeInputService15TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService15TestShapeInputService15TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService15TestShapeInputService15TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	Foo *string `location:"querystring" locationName:"param-name" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService15TestShapeInputService15TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "param-name", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService15TestShapeInputService15TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService15TestShapeInputService15TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService15TestShapeInputService15TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService16ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService16ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService16ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService16ProtocolTest client with a config.
//
// Example:
//     // Create a InputService16ProtocolTest client from just a config.
//     svc := inputservice16protocoltest.New(myConfig)
func NewInputService16ProtocolTest(config aws.Config) *InputService16ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService16ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice16protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService16ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService16ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService16TestCaseOperation1 = "OperationName"

// InputService16TestCaseOperation1Request is a API request type for the InputService16TestCaseOperation1 API operation.
type InputService16TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation6Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation1Request
}

// Send marshals and sends the InputService16TestCaseOperation1 API request.
func (r InputService16TestCaseOperation1Request) Send() (*InputService16TestShapeInputService16TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation1Output), nil
}

// InputService16TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService16TestCaseOperation1Request method.
//    req := client.InputService16TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation1Request(input *InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation6Input{}
	}

	output := &InputService16TestShapeInputService16TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService16TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation1Request}
}

const opInputService16TestCaseOperation2 = "OperationName"

// InputService16TestCaseOperation2Request is a API request type for the InputService16TestCaseOperation2 API operation.
type InputService16TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation6Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation2Request
}

// Send marshals and sends the InputService16TestCaseOperation2 API request.
func (r InputService16TestCaseOperation2Request) Send() (*InputService16TestShapeInputService16TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation2Output), nil
}

// InputService16TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService16TestCaseOperation2Request method.
//    req := client.InputService16TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation2Request(input *InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation6Input{}
	}

	output := &InputService16TestShapeInputService16TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService16TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation2Request}
}

const opInputService16TestCaseOperation3 = "OperationName"

// InputService16TestCaseOperation3Request is a API request type for the InputService16TestCaseOperation3 API operation.
type InputService16TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation6Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation3Request
}

// Send marshals and sends the InputService16TestCaseOperation3 API request.
func (r InputService16TestCaseOperation3Request) Send() (*InputService16TestShapeInputService16TestCaseOperation3Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation3Output), nil
}

// InputService16TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService16TestCaseOperation3Request method.
//    req := client.InputService16TestCaseOperation3Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation3Request(input *InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation3Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation3,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation6Input{}
	}

	output := &InputService16TestShapeInputService16TestCaseOperation3Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService16TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation3Request}
}

const opInputService16TestCaseOperation4 = "OperationName"

// InputService16TestCaseOperation4Request is a API request type for the InputService16TestCaseOperation4 API operation.
type InputService16TestCaseOperation4Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation6Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation4Request
}

// Send marshals and sends the InputService16TestCaseOperation4 API request.
func (r InputService16TestCaseOperation4Request) Send() (*InputService16TestShapeInputService16TestCaseOperation4Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation4Output), nil
}

// InputService16TestCaseOperation4Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService16TestCaseOperation4Request method.
//    req := client.InputService16TestCaseOperation4Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation4Request(input *InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation4Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation4,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation6Input{}
	}

	output := &InputService16TestShapeInputService16TestCaseOperation4Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService16TestCaseOperation4Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation4Request}
}

const opInputService16TestCaseOperation5 = "OperationName"

// InputService16TestCaseOperation5Request is a API request type for the InputService16TestCaseOperation5 API operation.
type InputService16TestCaseOperation5Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation6Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation5Request
}

// Send marshals and sends the InputService16TestCaseOperation5 API request.
func (r InputService16TestCaseOperation5Request) Send() (*InputService16TestShapeInputService16TestCaseOperation5Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation5Output), nil
}

// InputService16TestCaseOperation5Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService16TestCaseOperation5Request method.
//    req := client.InputService16TestCaseOperation5Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation5Request(input *InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation5Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation5,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation6Input{}
	}

	output := &InputService16TestShapeInputService16TestCaseOperation5Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService16TestCaseOperation5Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation5Request}
}

const opInputService16TestCaseOperation6 = "OperationName"

// InputService16TestCaseOperation6Request is a API request type for the InputService16TestCaseOperation6 API operation.
type InputService16TestCaseOperation6Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation6Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation6Request
}

// Send marshals and sends the InputService16TestCaseOperation6 API request.
func (r InputService16TestCaseOperation6Request) Send() (*InputService16TestShapeInputService16TestCaseOperation6Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation6Output), nil
}

// InputService16TestCaseOperation6Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService16TestCaseOperation6Request method.
//    req := client.InputService16TestCaseOperation6Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation6Request(input *InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation6Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation6,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation6Input{}
	}

	output := &InputService16TestShapeInputService16TestCaseOperation6Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService16TestCaseOperation6Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation6Request}
}

type InputService16TestShapeInputService16TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService16TestShapeInputService16TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService16TestShapeInputService16TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService16TestShapeInputService16TestCaseOperation3Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation3Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation4Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService16TestShapeInputService16TestCaseOperation4Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation4Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation5Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService16TestShapeInputService16TestCaseOperation5Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation5Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation6Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService16TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation6Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.RecursiveStruct != nil {
		v := s.RecursiveStruct

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
	}
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation6Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService16TestShapeInputService16TestCaseOperation6Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation6Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService16TestShapeRecursiveStructType struct {
	_ struct{} `type:"structure"`

	NoRecurse *string `type:"string"`

	RecursiveList []InputService16TestShapeRecursiveStructType `type:"list"`

	RecursiveMap map[string]InputService16TestShapeRecursiveStructType `type:"map"`

	RecursiveStruct *InputService16TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeRecursiveStructType) MarshalFields(e protocol.FieldEncoder) error {
	if s.NoRecurse != nil {
		v := *s.NoRecurse

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "NoRecurse", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if len(s.RecursiveList) > 0 {
		v := s.RecursiveList

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.BodyTarget, "RecursiveList", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddFields(v1)
		}
		ls0.End()

	}
	if len(s.RecursiveMap) > 0 {
		v := s.RecursiveMap

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.BodyTarget, "RecursiveMap", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ms0.MapSetFields(k1, v1)
		}
		ms0.End()

	}
	if s.RecursiveStruct != nil {
		v := s.RecursiveStruct

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
	}
	return nil
}

// InputService17ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService17ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService17ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService17ProtocolTest client with a config.
//
// Example:
//     // Create a InputService17ProtocolTest client from just a config.
//     svc := inputservice17protocoltest.New(myConfig)
func NewInputService17ProtocolTest(config aws.Config) *InputService17ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService17ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice17protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService17ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService17ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService17TestCaseOperation1 = "OperationName"

// InputService17TestCaseOperation1Request is a API request type for the InputService17TestCaseOperation1 API operation.
type InputService17TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService17TestShapeInputService17TestCaseOperation2Input
	Copy  func(*InputService17TestShapeInputService17TestCaseOperation2Input) InputService17TestCaseOperation1Request
}

// Send marshals and sends the InputService17TestCaseOperation1 API request.
func (r InputService17TestCaseOperation1Request) Send() (*InputService17TestShapeInputService17TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService17TestShapeInputService17TestCaseOperation1Output), nil
}

// InputService17TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService17TestCaseOperation1Request method.
//    req := client.InputService17TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService17ProtocolTest) InputService17TestCaseOperation1Request(input *InputService17TestShapeInputService17TestCaseOperation2Input) InputService17TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService17TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService17TestShapeInputService17TestCaseOperation2Input{}
	}

	output := &InputService17TestShapeInputService17TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService17TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService17TestCaseOperation1Request}
}

const opInputService17TestCaseOperation2 = "OperationName"

// InputService17TestCaseOperation2Request is a API request type for the InputService17TestCaseOperation2 API operation.
type InputService17TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService17TestShapeInputService17TestCaseOperation2Input
	Copy  func(*InputService17TestShapeInputService17TestCaseOperation2Input) InputService17TestCaseOperation2Request
}

// Send marshals and sends the InputService17TestCaseOperation2 API request.
func (r InputService17TestCaseOperation2Request) Send() (*InputService17TestShapeInputService17TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService17TestShapeInputService17TestCaseOperation2Output), nil
}

// InputService17TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService17TestCaseOperation2Request method.
//    req := client.InputService17TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService17ProtocolTest) InputService17TestCaseOperation2Request(input *InputService17TestShapeInputService17TestCaseOperation2Input) InputService17TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService17TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService17TestShapeInputService17TestCaseOperation2Input{}
	}

	output := &InputService17TestShapeInputService17TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService17TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService17TestCaseOperation2Request}
}

type InputService17TestShapeInputService17TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService17TestShapeInputService17TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService17TestShapeInputService17TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService17TestShapeInputService17TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	TimeArg *time.Time `type:"timestamp" timestampFormat:"unix"`

	TimeArgInHeader *time.Time `location:"header" locationName:"x-amz-timearg" type:"timestamp" timestampFormat:"rfc822"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService17TestShapeInputService17TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.TimeArg != nil {
		v := *s.TimeArg

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "TimeArg", protocol.TimeValue{V: v, Format: protocol.UnixTimeFormat}, metadata)
	}
	if s.TimeArgInHeader != nil {
		v := *s.TimeArgInHeader

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "x-amz-timearg", protocol.TimeValue{V: v, Format: protocol.RFC822TimeFromat}, metadata)
	}
	return nil
}

type InputService17TestShapeInputService17TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService17TestShapeInputService17TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService17TestShapeInputService17TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService18ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService18ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService18ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService18ProtocolTest client with a config.
//
// Example:
//     // Create a InputService18ProtocolTest client from just a config.
//     svc := inputservice18protocoltest.New(myConfig)
func NewInputService18ProtocolTest(config aws.Config) *InputService18ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService18ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice18protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService18ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService18ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService18TestCaseOperation1 = "OperationName"

// InputService18TestCaseOperation1Request is a API request type for the InputService18TestCaseOperation1 API operation.
type InputService18TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService18TestShapeInputService18TestCaseOperation1Input
	Copy  func(*InputService18TestShapeInputService18TestCaseOperation1Input) InputService18TestCaseOperation1Request
}

// Send marshals and sends the InputService18TestCaseOperation1 API request.
func (r InputService18TestCaseOperation1Request) Send() (*InputService18TestShapeInputService18TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService18TestShapeInputService18TestCaseOperation1Output), nil
}

// InputService18TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService18TestCaseOperation1Request method.
//    req := client.InputService18TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService18ProtocolTest) InputService18TestCaseOperation1Request(input *InputService18TestShapeInputService18TestCaseOperation1Input) InputService18TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService18TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService18TestShapeInputService18TestCaseOperation1Input{}
	}

	output := &InputService18TestShapeInputService18TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService18TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService18TestCaseOperation1Request}
}

type InputService18TestShapeInputService18TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	TimeArg *time.Time `locationName:"timestamp_location" type:"timestamp" timestampFormat:"unix"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService18TestShapeInputService18TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.TimeArg != nil {
		v := *s.TimeArg

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "timestamp_location", protocol.TimeValue{V: v, Format: protocol.UnixTimeFormat}, metadata)
	}
	return nil
}

type InputService18TestShapeInputService18TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService18TestShapeInputService18TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService18TestShapeInputService18TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService19ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService19ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService19ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService19ProtocolTest client with a config.
//
// Example:
//     // Create a InputService19ProtocolTest client from just a config.
//     svc := inputservice19protocoltest.New(myConfig)
func NewInputService19ProtocolTest(config aws.Config) *InputService19ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService19ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice19protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService19ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService19ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService19TestCaseOperation1 = "OperationName"

// InputService19TestCaseOperation1Request is a API request type for the InputService19TestCaseOperation1 API operation.
type InputService19TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService19TestShapeInputService19TestCaseOperation1Input
	Copy  func(*InputService19TestShapeInputService19TestCaseOperation1Input) InputService19TestCaseOperation1Request
}

// Send marshals and sends the InputService19TestCaseOperation1 API request.
func (r InputService19TestCaseOperation1Request) Send() (*InputService19TestShapeInputService19TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService19TestShapeInputService19TestCaseOperation1Output), nil
}

// InputService19TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService19TestCaseOperation1Request method.
//    req := client.InputService19TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService19ProtocolTest) InputService19TestCaseOperation1Request(input *InputService19TestShapeInputService19TestCaseOperation1Input) InputService19TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService19TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService19TestShapeInputService19TestCaseOperation1Input{}
	}

	output := &InputService19TestShapeInputService19TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService19TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService19TestCaseOperation1Request}
}

type InputService19TestShapeInputService19TestCaseOperation1Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo *string `locationName:"foo" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetStream(protocol.PayloadTarget, "foo", protocol.StringStream(v), metadata)
	}
	return nil
}

type InputService19TestShapeInputService19TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService19TestShapeInputService19TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService20ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService20ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService20ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService20ProtocolTest client with a config.
//
// Example:
//     // Create a InputService20ProtocolTest client from just a config.
//     svc := inputservice20protocoltest.New(myConfig)
func NewInputService20ProtocolTest(config aws.Config) *InputService20ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService20ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice20protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService20ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService20ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService20TestCaseOperation1 = "OperationName"

// InputService20TestCaseOperation1Request is a API request type for the InputService20TestCaseOperation1 API operation.
type InputService20TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService20TestShapeInputService20TestCaseOperation2Input
	Copy  func(*InputService20TestShapeInputService20TestCaseOperation2Input) InputService20TestCaseOperation1Request
}

// Send marshals and sends the InputService20TestCaseOperation1 API request.
func (r InputService20TestCaseOperation1Request) Send() (*InputService20TestShapeInputService20TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService20TestShapeInputService20TestCaseOperation1Output), nil
}

// InputService20TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService20TestCaseOperation1Request method.
//    req := client.InputService20TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService20ProtocolTest) InputService20TestCaseOperation1Request(input *InputService20TestShapeInputService20TestCaseOperation2Input) InputService20TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService20TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService20TestShapeInputService20TestCaseOperation2Input{}
	}

	output := &InputService20TestShapeInputService20TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService20TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService20TestCaseOperation1Request}
}

const opInputService20TestCaseOperation2 = "OperationName"

// InputService20TestCaseOperation2Request is a API request type for the InputService20TestCaseOperation2 API operation.
type InputService20TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService20TestShapeInputService20TestCaseOperation2Input
	Copy  func(*InputService20TestShapeInputService20TestCaseOperation2Input) InputService20TestCaseOperation2Request
}

// Send marshals and sends the InputService20TestCaseOperation2 API request.
func (r InputService20TestCaseOperation2Request) Send() (*InputService20TestShapeInputService20TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService20TestShapeInputService20TestCaseOperation2Output), nil
}

// InputService20TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService20TestCaseOperation2Request method.
//    req := client.InputService20TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService20ProtocolTest) InputService20TestCaseOperation2Request(input *InputService20TestShapeInputService20TestCaseOperation2Input) InputService20TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService20TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService20TestShapeInputService20TestCaseOperation2Input{}
	}

	output := &InputService20TestShapeInputService20TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService20TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService20TestCaseOperation2Request}
}

type InputService20TestShapeInputService20TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService20TestShapeInputService20TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeInputService20TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService20TestShapeInputService20TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	Token *string `type:"string" idempotencyToken:"true"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeInputService20TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	var Token string
	if s.Token != nil {
		Token = *s.Token
	} else {
		Token = protocol.GetIdempotencyToken()
	}
	{
		v := Token

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Token", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService20TestShapeInputService20TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService20TestShapeInputService20TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeInputService20TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService21ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService21ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService21ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService21ProtocolTest client with a config.
//
// Example:
//     // Create a InputService21ProtocolTest client from just a config.
//     svc := inputservice21protocoltest.New(myConfig)
func NewInputService21ProtocolTest(config aws.Config) *InputService21ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService21ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice21protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService21ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService21ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService21TestCaseOperation1 = "OperationName"

// InputService21TestCaseOperation1Request is a API request type for the InputService21TestCaseOperation1 API operation.
type InputService21TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService21TestShapeInputService21TestCaseOperation3Input
	Copy  func(*InputService21TestShapeInputService21TestCaseOperation3Input) InputService21TestCaseOperation1Request
}

// Send marshals and sends the InputService21TestCaseOperation1 API request.
func (r InputService21TestCaseOperation1Request) Send() (*InputService21TestShapeInputService21TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService21TestShapeInputService21TestCaseOperation1Output), nil
}

// InputService21TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService21TestCaseOperation1Request method.
//    req := client.InputService21TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService21ProtocolTest) InputService21TestCaseOperation1Request(input *InputService21TestShapeInputService21TestCaseOperation3Input) InputService21TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService21TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService21TestShapeInputService21TestCaseOperation3Input{}
	}

	output := &InputService21TestShapeInputService21TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService21TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService21TestCaseOperation1Request}
}

const opInputService21TestCaseOperation2 = "OperationName"

// InputService21TestCaseOperation2Request is a API request type for the InputService21TestCaseOperation2 API operation.
type InputService21TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService21TestShapeInputService21TestCaseOperation3Input
	Copy  func(*InputService21TestShapeInputService21TestCaseOperation3Input) InputService21TestCaseOperation2Request
}

// Send marshals and sends the InputService21TestCaseOperation2 API request.
func (r InputService21TestCaseOperation2Request) Send() (*InputService21TestShapeInputService21TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService21TestShapeInputService21TestCaseOperation2Output), nil
}

// InputService21TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService21TestCaseOperation2Request method.
//    req := client.InputService21TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService21ProtocolTest) InputService21TestCaseOperation2Request(input *InputService21TestShapeInputService21TestCaseOperation3Input) InputService21TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService21TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService21TestShapeInputService21TestCaseOperation3Input{}
	}

	output := &InputService21TestShapeInputService21TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService21TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService21TestCaseOperation2Request}
}

const opInputService21TestCaseOperation3 = "OperationName"

// InputService21TestCaseOperation3Request is a API request type for the InputService21TestCaseOperation3 API operation.
type InputService21TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService21TestShapeInputService21TestCaseOperation3Input
	Copy  func(*InputService21TestShapeInputService21TestCaseOperation3Input) InputService21TestCaseOperation3Request
}

// Send marshals and sends the InputService21TestCaseOperation3 API request.
func (r InputService21TestCaseOperation3Request) Send() (*InputService21TestShapeInputService21TestCaseOperation3Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService21TestShapeInputService21TestCaseOperation3Output), nil
}

// InputService21TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService21TestCaseOperation3Request method.
//    req := client.InputService21TestCaseOperation3Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService21ProtocolTest) InputService21TestCaseOperation3Request(input *InputService21TestShapeInputService21TestCaseOperation3Input) InputService21TestCaseOperation3Request {
	op := &aws.Operation{
		Name:       opInputService21TestCaseOperation3,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService21TestShapeInputService21TestCaseOperation3Input{}
	}

	output := &InputService21TestShapeInputService21TestCaseOperation3Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService21TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService21TestCaseOperation3Request}
}

type InputService21TestShapeBodyStructure struct {
	_ struct{} `type:"structure"`

	BodyField aws.JSONValue `type:"jsonvalue"`

	BodyListField []aws.JSONValue `type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeBodyStructure) MarshalFields(e protocol.FieldEncoder) error {
	if s.BodyField != nil {
		v := s.BodyField

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "BodyField", protocol.JSONValue{V: v, EscapeMode: protocol.QuotedEscape}, metadata)
	}
	if len(s.BodyListField) > 0 {
		v := s.BodyListField

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.BodyTarget, "BodyListField", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddValue(protocol.JSONValue{V: v1, EscapeMode: protocol.QuotedEscape})
		}
		ls0.End()

	}
	return nil
}

type InputService21TestShapeInputService21TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService21TestShapeInputService21TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService21TestShapeInputService21TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService21TestShapeInputService21TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService21TestShapeInputService21TestCaseOperation3Input struct {
	_ struct{} `type:"structure" payload:"Body"`

	Body *InputService21TestShapeBodyStructure `type:"structure"`

	HeaderField aws.JSONValue `location:"header" locationName:"X-Amz-Foo" type:"jsonvalue"`

	QueryField aws.JSONValue `location:"querystring" locationName:"Bar" type:"jsonvalue"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation3Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.HeaderField != nil {
		v := s.HeaderField

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "X-Amz-Foo", protocol.JSONValue{V: v, EscapeMode: protocol.Base64Escape}, metadata)
	}
	if s.Body != nil {
		v := s.Body

		metadata := protocol.Metadata{}
		e.SetFields(protocol.PayloadTarget, "Body", v, metadata)
	}
	if s.QueryField != nil {
		v := s.QueryField

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "Bar", protocol.JSONValue{V: v}, metadata)
	}
	return nil
}

type InputService21TestShapeInputService21TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService21TestShapeInputService21TestCaseOperation3Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation3Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

// InputService22ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService22ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService22ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService22ProtocolTest client with a config.
//
// Example:
//     // Create a InputService22ProtocolTest client from just a config.
//     svc := inputservice22protocoltest.New(myConfig)
func NewInputService22ProtocolTest(config aws.Config) *InputService22ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService22ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice22protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService22ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService22ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService22TestCaseOperation1 = "OperationName"

// InputService22TestCaseOperation1Request is a API request type for the InputService22TestCaseOperation1 API operation.
type InputService22TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService22TestShapeInputService22TestCaseOperation2Input
	Copy  func(*InputService22TestShapeInputService22TestCaseOperation2Input) InputService22TestCaseOperation1Request
}

// Send marshals and sends the InputService22TestCaseOperation1 API request.
func (r InputService22TestCaseOperation1Request) Send() (*InputService22TestShapeInputService22TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService22TestShapeInputService22TestCaseOperation1Output), nil
}

// InputService22TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService22TestCaseOperation1Request method.
//    req := client.InputService22TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService22ProtocolTest) InputService22TestCaseOperation1Request(input *InputService22TestShapeInputService22TestCaseOperation2Input) InputService22TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService22TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService22TestShapeInputService22TestCaseOperation2Input{}
	}

	output := &InputService22TestShapeInputService22TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService22TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService22TestCaseOperation1Request}
}

const opInputService22TestCaseOperation2 = "OperationName"

// InputService22TestCaseOperation2Request is a API request type for the InputService22TestCaseOperation2 API operation.
type InputService22TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService22TestShapeInputService22TestCaseOperation2Input
	Copy  func(*InputService22TestShapeInputService22TestCaseOperation2Input) InputService22TestCaseOperation2Request
}

// Send marshals and sends the InputService22TestCaseOperation2 API request.
func (r InputService22TestCaseOperation2Request) Send() (*InputService22TestShapeInputService22TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService22TestShapeInputService22TestCaseOperation2Output), nil
}

// InputService22TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService22TestCaseOperation2Request method.
//    req := client.InputService22TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService22ProtocolTest) InputService22TestCaseOperation2Request(input *InputService22TestShapeInputService22TestCaseOperation2Input) InputService22TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService22TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService22TestShapeInputService22TestCaseOperation2Input{}
	}

	output := &InputService22TestShapeInputService22TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	return InputService22TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService22TestCaseOperation2Request}
}

type InputService22TestShapeInputService22TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService22TestShapeInputService22TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService22TestShapeInputService22TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService22TestShapeInputService22TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService22TestShapeEnumType `type:"string" enum:"true"`

	HeaderEnum InputService22TestShapeEnumType `location:"header" locationName:"x-amz-enum" type:"string" enum:"true"`

	ListEnums []InputService22TestShapeEnumType `type:"list"`

	QueryFooEnum InputService22TestShapeEnumType `location:"querystring" locationName:"Enum" type:"string" enum:"true"`

	QueryListEnums []InputService22TestShapeEnumType `location:"querystring" locationName:"List" type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService22TestShapeInputService22TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if len(s.FooEnum) > 0 {
		v := s.FooEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "FooEnum", protocol.QuotedValue{ValueMarshaler: v}, metadata)
	}
	if len(s.ListEnums) > 0 {
		v := s.ListEnums

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.BodyTarget, "ListEnums", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddValue(protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v1)})
		}
		ls0.End()

	}
	if len(s.HeaderEnum) > 0 {
		v := s.HeaderEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "x-amz-enum", protocol.QuotedValue{ValueMarshaler: v}, metadata)
	}
	if len(s.QueryFooEnum) > 0 {
		v := s.QueryFooEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "Enum", protocol.QuotedValue{ValueMarshaler: v}, metadata)
	}
	if len(s.QueryListEnums) > 0 {
		v := s.QueryListEnums

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.QueryTarget, "List", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddValue(protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v1)})
		}
		ls0.End()

	}
	return nil
}

type InputService22TestShapeInputService22TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService22TestShapeInputService22TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService22TestShapeInputService22TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService22TestShapeEnumType string

// Enum values for InputService22TestShapeEnumType
const (
	EnumTypeFoo InputService22TestShapeEnumType = "foo"
	EnumTypeBar InputService22TestShapeEnumType = "bar"
	EnumType0   InputService22TestShapeEnumType = "0"
	EnumType1   InputService22TestShapeEnumType = "1"
)

func (enum InputService22TestShapeEnumType) MarshalValue() (string, error) {
	return string(enum), nil
}

func (enum InputService22TestShapeEnumType) MarshalValueBuf(b []byte) ([]byte, error) {
	b = b[0:0]
	return append(b, enum...), nil
}

// InputService23ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// InputService23ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService23ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the InputService23ProtocolTest client with a config.
//
// Example:
//     // Create a InputService23ProtocolTest client from just a config.
//     svc := inputservice23protocoltest.New(myConfig)
func NewInputService23ProtocolTest(config aws.Config) *InputService23ProtocolTest {
	var signingName string
	signingRegion := config.Region

	svc := &InputService23ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "inputservice23protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2014-01-01",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a InputService23ProtocolTest operation and runs any
// custom request initialization.
func (c *InputService23ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opInputService23TestCaseOperation1 = "StaticOp"

// InputService23TestCaseOperation1Request is a API request type for the InputService23TestCaseOperation1 API operation.
type InputService23TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService23TestShapeInputService23TestCaseOperation1Input
	Copy  func(*InputService23TestShapeInputService23TestCaseOperation1Input) InputService23TestCaseOperation1Request
}

// Send marshals and sends the InputService23TestCaseOperation1 API request.
func (r InputService23TestCaseOperation1Request) Send() (*InputService23TestShapeInputService23TestCaseOperation1Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService23TestShapeInputService23TestCaseOperation1Output), nil
}

// InputService23TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService23TestCaseOperation1Request method.
//    req := client.InputService23TestCaseOperation1Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService23ProtocolTest) InputService23TestCaseOperation1Request(input *InputService23TestShapeInputService23TestCaseOperation1Input) InputService23TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService23TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService23TestShapeInputService23TestCaseOperation1Input{}
	}

	output := &InputService23TestShapeInputService23TestCaseOperation1Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("data-", nil))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)
	return InputService23TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService23TestCaseOperation1Request}
}

const opInputService23TestCaseOperation2 = "MemberRefOp"

// InputService23TestCaseOperation2Request is a API request type for the InputService23TestCaseOperation2 API operation.
type InputService23TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService23TestShapeInputService23TestCaseOperation2Input
	Copy  func(*InputService23TestShapeInputService23TestCaseOperation2Input) InputService23TestCaseOperation2Request
}

// Send marshals and sends the InputService23TestCaseOperation2 API request.
func (r InputService23TestCaseOperation2Request) Send() (*InputService23TestShapeInputService23TestCaseOperation2Output, error) {
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	return r.Request.Data.(*InputService23TestShapeInputService23TestCaseOperation2Output), nil
}

// InputService23TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using the InputService23TestCaseOperation2Request method.
//    req := client.InputService23TestCaseOperation2Request(params)
//    resp, err := req.Send()
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService23ProtocolTest) InputService23TestCaseOperation2Request(input *InputService23TestShapeInputService23TestCaseOperation2Input) InputService23TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService23TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService23TestShapeInputService23TestCaseOperation2Input{}
	}

	output := &InputService23TestShapeInputService23TestCaseOperation2Output{}
	req := c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	output.responseMetadata = aws.Response{Request: req}

	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("foo-{Name}.", input.hostLabels))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)
	return InputService23TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService23TestCaseOperation2Request}
}

type InputService23TestShapeInputService23TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Name *string `type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Name != nil {
		v := *s.Name

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Name", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService23TestShapeInputService23TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService23TestShapeInputService23TestCaseOperation1Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

type InputService23TestShapeInputService23TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	// Name is a required field
	Name *string `type:"string" required:"true"`
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *InputService23TestShapeInputService23TestCaseOperation2Input) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "InputService23TestShapeInputService23TestCaseOperation2Input"}

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

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Name != nil {
		v := *s.Name

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Name", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

func (s *InputService23TestShapeInputService23TestCaseOperation2Input) hostLabels() map[string]string {
	return map[string]string{
		"Name": aws.StringValue(s.Name),
	}
}

type InputService23TestShapeInputService23TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`

	responseMetadata aws.Response
}

// SDKResponseMetdata return sthe response metadata for the API.
func (s InputService23TestShapeInputService23TestCaseOperation2Output) SDKResponseMetadata() aws.Response {
	return s.responseMetadata
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

//
// Tests begin here
//

func TestInputService1ProtocolTestNoParametersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService1ProtocolTest(cfg)

	req := svc.InputService1TestCaseOperation1Request(nil)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/jobs", r.URL.String())

	// assert headers

}

func TestInputService2ProtocolTestURIParameterOnlyWithNoLocationNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService2ProtocolTest(cfg)
	input := &InputService2TestShapeInputService2TestCaseOperation1Input{
		PipelineId: aws.String("foo"),
	}

	req := svc.InputService2TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/jobsByPipeline/foo", r.URL.String())

	// assert headers

}

func TestInputService3ProtocolTestURIParameterOnlyWithLocationNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService3ProtocolTest(cfg)
	input := &InputService3TestShapeInputService3TestCaseOperation1Input{
		Foo: aws.String("bar"),
	}

	req := svc.InputService3TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/jobsByPipeline/bar", r.URL.String())

	// assert headers

}

func TestInputService4ProtocolTestQuerystringListOfStringsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService4ProtocolTest(cfg)
	input := &InputService4TestShapeInputService4TestCaseOperation1Input{
		Items: []string{
			"value1",
			"value2",
		},
	}

	req := svc.InputService4TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/path?item=value1&item=value2", r.URL.String())

	// assert headers

}

func TestInputService5ProtocolTestStringToStringMapsInQuerystringCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation1Input{
		PipelineId: aws.String("foo"),
		QueryDoc: map[string]string{
			"bar":  "baz",
			"fizz": "buzz",
		},
	}

	req := svc.InputService5TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/jobsByPipeline/foo?bar=baz&fizz=buzz", r.URL.String())

	// assert headers

}

func TestInputService6ProtocolTestStringToStringListMapsInQuerystringCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService6ProtocolTest(cfg)
	input := &InputService6TestShapeInputService6TestCaseOperation1Input{
		PipelineId: aws.String("id"),
		QueryDoc: map[string][]string{
			"fizz": {
				"buzz",
				"pop",
			},
			"foo": {
				"bar",
				"baz",
			},
		},
	}

	req := svc.InputService6TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/jobsByPipeline/id?foo=bar&foo=baz&fizz=buzz&fizz=pop", r.URL.String())

	// assert headers

}

func TestInputService7ProtocolTestBooleanInQuerystringCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService7ProtocolTest(cfg)
	input := &InputService7TestShapeInputService7TestCaseOperation2Input{
		BoolQuery: aws.Bool(true),
	}

	req := svc.InputService7TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/path?bool-query=true", r.URL.String())

	// assert headers

}

func TestInputService7ProtocolTestBooleanInQuerystringCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService7ProtocolTest(cfg)
	input := &InputService7TestShapeInputService7TestCaseOperation2Input{
		BoolQuery: aws.Bool(false),
	}

	req := svc.InputService7TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/path?bool-query=false", r.URL.String())

	// assert headers

}

func TestInputService8ProtocolTestURIParameterAndQuerystringParamsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService8ProtocolTest(cfg)
	input := &InputService8TestShapeInputService8TestCaseOperation1Input{
		Ascending:  aws.String("true"),
		PageToken:  aws.String("bar"),
		PipelineId: aws.String("foo"),
	}

	req := svc.InputService8TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/jobsByPipeline/foo?Ascending=true&PageToken=bar", r.URL.String())

	// assert headers

}

func TestInputService9ProtocolTestURIParameterQuerystringParamsAndJSONBodyCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService9ProtocolTest(cfg)
	input := &InputService9TestShapeInputService9TestCaseOperation1Input{
		Ascending: aws.String("true"),
		Config: &InputService9TestShapeStructType{
			A: aws.String("one"),
			B: aws.String("two"),
		},
		PageToken:  aws.String("bar"),
		PipelineId: aws.String("foo"),
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
	awstesting.AssertJSON(t, `{"Config":{"A":"one","B":"two"}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/jobsByPipeline/foo?Ascending=true&PageToken=bar", r.URL.String())

	// assert headers

}

func TestInputService10ProtocolTestURIParameterQuerystringParamsHeadersAndJSONBodyCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService10ProtocolTest(cfg)
	input := &InputService10TestShapeInputService10TestCaseOperation1Input{
		Ascending: aws.String("true"),
		Checksum:  aws.String("12345"),
		Config: &InputService10TestShapeStructType{
			A: aws.String("one"),
			B: aws.String("two"),
		},
		PageToken:  aws.String("bar"),
		PipelineId: aws.String("foo"),
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
	awstesting.AssertJSON(t, `{"Config":{"A":"one","B":"two"}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/jobsByPipeline/foo?Ascending=true&PageToken=bar", r.URL.String())

	// assert headers
	if e, a := "12345", r.Header.Get("x-amz-checksum"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService11ProtocolTestStreamingPayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService11ProtocolTest(cfg)
	input := &InputService11TestShapeInputService11TestCaseOperation1Input{
		Body:      bytes.NewReader([]byte("contents")),
		Checksum:  aws.String("foo"),
		VaultName: aws.String("name"),
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
	if e, a := "contents", util.Trim(string(body)); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/vaults/name/archives", r.URL.String())

	// assert headers
	if e, a := "foo", r.Header.Get("x-amz-sha256-tree-hash"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService12ProtocolTestSerializeBlobsInBodyCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService12ProtocolTest(cfg)
	input := &InputService12TestShapeInputService12TestCaseOperation1Input{
		Bar: []byte("Blob param"),
		Foo: aws.String("foo_name"),
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
	awstesting.AssertJSON(t, `{"Bar":"QmxvYiBwYXJhbQ=="}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/foo_name", r.URL.String())

	// assert headers

}

func TestInputService13ProtocolTestBlobPayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation2Input{
		Foo: []byte("bar"),
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
	if e, a := "bar", util.Trim(string(body)); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService13ProtocolTestBlobPayloadCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation2Input{}

	req := svc.InputService13TestCaseOperation2Request(input)
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

func TestInputService14ProtocolTestStructurePayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService14ProtocolTest(cfg)
	input := &InputService14TestShapeInputService14TestCaseOperation2Input{
		Foo: &InputService14TestShapeFooShape{
			Baz: aws.String("bar"),
		},
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
	awstesting.AssertJSON(t, `{"baz":"bar"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService14ProtocolTestStructurePayloadCase2(t *testing.T) {
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

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService15ProtocolTestOmitsNullQueryParamsButSerializesEmptyStringsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService15ProtocolTest(cfg)
	input := &InputService15TestShapeInputService15TestCaseOperation2Input{}

	req := svc.InputService15TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService15ProtocolTestOmitsNullQueryParamsButSerializesEmptyStringsCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService15ProtocolTest(cfg)
	input := &InputService15TestShapeInputService15TestCaseOperation2Input{
		Foo: aws.String(""),
	}

	req := svc.InputService15TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/path?abc=mno&param-name=", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation6Input{
		RecursiveStruct: &InputService16TestShapeRecursiveStructType{
			NoRecurse: aws.String("foo"),
		},
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
	awstesting.AssertJSON(t, `{"RecursiveStruct":{"NoRecurse":"foo"}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation6Input{
		RecursiveStruct: &InputService16TestShapeRecursiveStructType{
			RecursiveStruct: &InputService16TestShapeRecursiveStructType{
				NoRecurse: aws.String("foo"),
			},
		},
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
	awstesting.AssertJSON(t, `{"RecursiveStruct":{"RecursiveStruct":{"NoRecurse":"foo"}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation6Input{
		RecursiveStruct: &InputService16TestShapeRecursiveStructType{
			RecursiveStruct: &InputService16TestShapeRecursiveStructType{
				RecursiveStruct: &InputService16TestShapeRecursiveStructType{
					RecursiveStruct: &InputService16TestShapeRecursiveStructType{
						NoRecurse: aws.String("foo"),
					},
				},
			},
		},
	}

	req := svc.InputService16TestCaseOperation3Request(input)
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
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase4(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation6Input{
		RecursiveStruct: &InputService16TestShapeRecursiveStructType{
			RecursiveList: []InputService16TestShapeRecursiveStructType{
				{
					NoRecurse: aws.String("foo"),
				},
				{
					NoRecurse: aws.String("bar"),
				},
			},
		},
	}

	req := svc.InputService16TestCaseOperation4Request(input)
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
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase5(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation6Input{
		RecursiveStruct: &InputService16TestShapeRecursiveStructType{
			RecursiveList: []InputService16TestShapeRecursiveStructType{
				{
					NoRecurse: aws.String("foo"),
				},
				{
					RecursiveStruct: &InputService16TestShapeRecursiveStructType{
						NoRecurse: aws.String("bar"),
					},
				},
			},
		},
	}

	req := svc.InputService16TestCaseOperation5Request(input)
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
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase6(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation6Input{
		RecursiveStruct: &InputService16TestShapeRecursiveStructType{
			RecursiveMap: map[string]InputService16TestShapeRecursiveStructType{
				"bar": {
					NoRecurse: aws.String("bar"),
				},
				"foo": {
					NoRecurse: aws.String("foo"),
				},
			},
		},
	}

	req := svc.InputService16TestCaseOperation6Request(input)
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
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService17ProtocolTestTimestampValuesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService17ProtocolTest(cfg)
	input := &InputService17TestShapeInputService17TestCaseOperation2Input{
		TimeArg: aws.Time(time.Unix(1422172800, 0)),
	}

	req := svc.InputService17TestCaseOperation1Request(input)
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
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService17ProtocolTestTimestampValuesCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService17ProtocolTest(cfg)
	input := &InputService17TestShapeInputService17TestCaseOperation2Input{
		TimeArgInHeader: aws.Time(time.Unix(1422172800, 0)),
	}

	req := svc.InputService17TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers
	if e, a := "Sun, 25 Jan 2015 08:00:00 GMT", r.Header.Get("x-amz-timearg"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService18ProtocolTestNamedLocationsInJSONBodyCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService18ProtocolTest(cfg)
	input := &InputService18TestShapeInputService18TestCaseOperation1Input{
		TimeArg: aws.Time(time.Unix(1422172800, 0)),
	}

	req := svc.InputService18TestCaseOperation1Request(input)
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
	awstesting.AssertJSON(t, `{"timestamp_location":1422172800}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService19ProtocolTestStringPayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService19ProtocolTest(cfg)
	input := &InputService19TestShapeInputService19TestCaseOperation1Input{
		Foo: aws.String("bar"),
	}

	req := svc.InputService19TestCaseOperation1Request(input)
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
	if e, a := "bar", util.Trim(string(body)); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService20ProtocolTestIdempotencyTokenAutoFillCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService20ProtocolTest(cfg)
	input := &InputService20TestShapeInputService20TestCaseOperation2Input{
		Token: aws.String("abc123"),
	}

	req := svc.InputService20TestCaseOperation1Request(input)
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
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService20ProtocolTestIdempotencyTokenAutoFillCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService20ProtocolTest(cfg)
	input := &InputService20TestShapeInputService20TestCaseOperation2Input{}

	req := svc.InputService20TestCaseOperation2Request(input)
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
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService21ProtocolTestJSONValueTraitCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService21ProtocolTest(cfg)
	input := &InputService21TestShapeInputService21TestCaseOperation3Input{
		Body: &InputService21TestShapeBodyStructure{
			BodyField: func() aws.JSONValue {
				var m aws.JSONValue
				if err := json.Unmarshal([]byte("{\"Foo\":\"Bar\"}"), &m); err != nil {
					panic("failed to unmarshal JSONValue, " + err.Error())
				}
				return m
			}(),
		},
	}
	input.HeaderField = aws.JSONValue{"Foo": "Bar"}
	input.QueryField = aws.JSONValue{"Foo": "Bar"}

	req := svc.InputService21TestCaseOperation1Request(input)
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
	awstesting.AssertJSON(t, `{"BodyField":"{\"Foo\":\"Bar\"}"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/?Bar=%7B%22Foo%22%3A%22Bar%22%7D", r.URL.String())

	// assert headers
	if e, a := "eyJGb28iOiJCYXIifQ==", r.Header.Get("X-Amz-Foo"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService21ProtocolTestJSONValueTraitCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService21ProtocolTest(cfg)
	input := &InputService21TestShapeInputService21TestCaseOperation3Input{
		Body: &InputService21TestShapeBodyStructure{
			BodyListField: []aws.JSONValue{
				func() aws.JSONValue {
					var m aws.JSONValue
					if err := json.Unmarshal([]byte("{\"Foo\":\"Bar\"}"), &m); err != nil {
						panic("failed to unmarshal JSONValue, " + err.Error())
					}
					return m
				}(),
			},
		},
	}

	req := svc.InputService21TestCaseOperation2Request(input)
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
	awstesting.AssertJSON(t, `{"BodyListField":["{\"Foo\":\"Bar\"}"]}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService21ProtocolTestJSONValueTraitCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService21ProtocolTest(cfg)
	input := &InputService21TestShapeInputService21TestCaseOperation3Input{}

	req := svc.InputService21TestCaseOperation3Request(input)
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

func TestInputService22ProtocolTestEnumCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService22ProtocolTest(cfg)
	input := &InputService22TestShapeInputService22TestCaseOperation2Input{
		FooEnum:    InputService22TestShapeEnumType("foo"),
		HeaderEnum: InputService22TestShapeEnumType("baz"),
		ListEnums: []InputService22TestShapeEnumType{
			InputService22TestShapeEnumType("foo"),
			InputService22TestShapeEnumType(""),
			InputService22TestShapeEnumType("bar"),
		},
		QueryFooEnum: InputService22TestShapeEnumType("bar"),
		QueryListEnums: []InputService22TestShapeEnumType{
			InputService22TestShapeEnumType("0"),
			InputService22TestShapeEnumType(""),
			InputService22TestShapeEnumType("1"),
		},
	}

	req := svc.InputService22TestCaseOperation1Request(input)
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
	awstesting.AssertURL(t, "https://test/path?Enum=bar&List=0&List=1&List=", r.URL.String())

	// assert headers
	if e, a := "baz", r.Header.Get("x-amz-enum"); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}

}

func TestInputService22ProtocolTestEnumCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService22ProtocolTest(cfg)
	input := &InputService22TestShapeInputService22TestCaseOperation2Input{}

	req := svc.InputService22TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService23ProtocolTestEndpointHostTraitCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService23ProtocolTest(cfg)
	input := &InputService23TestShapeInputService23TestCaseOperation1Input{
		Name: aws.String("myname"),
	}

	req := svc.InputService23TestCaseOperation1Request(input)
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
	awstesting.AssertURL(t, "https://data-service.region.amazonaws.com/path", r.URL.String())

	// assert headers

}

func TestInputService23ProtocolTestEndpointHostTraitCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService23ProtocolTest(cfg)
	input := &InputService23TestShapeInputService23TestCaseOperation2Input{
		Name: aws.String("myname"),
	}

	req := svc.InputService23TestCaseOperation2Request(input)
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
	awstesting.AssertURL(t, "https://foo-myname.service.region.amazonaws.com/path", r.URL.String())

	// assert headers

}
