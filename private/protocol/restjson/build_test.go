package restjson_test

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
	"github.com/jviney/aws-sdk-go-v2/private/protocol/restjson"
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService1TestShapeInputService1TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	return nil
}

type InputService1TestShapeInputService1TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService1TestShapeInputService1TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobs",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService1TestShapeInputService1TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	PipelineId *string `location:"uri" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService2TestShapeInputService2TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService2TestShapeInputService2TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService2TestShapeInputService2TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService2TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService2TestShapeInputService2TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService2TestShapeInputService2TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	Foo *string `location:"uri" locationName:"PipelineId" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService3TestShapeInputService3TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService3TestShapeInputService3TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService3TestShapeInputService3TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService3TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService3TestShapeInputService3TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService3TestShapeInputService3TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	Items []string `location:"querystring" locationName:"item" type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService4TestShapeInputService4TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Items != nil {
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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService4TestShapeInputService4TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService4TestShapeInputService4TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService4TestShapeInputService4TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
				APIVersion:    "2014-01-01",
			},
		),
	}

	if config.Retryer == nil {
		svc.Retryer = retry.NewStandard()
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	PipelineId *string `location:"uri" type:"string"`

	QueryDoc map[string]string `location:"querystring" type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService5TestShapeInputService5TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.QueryDoc != nil {
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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService5TestShapeInputService5TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService5TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService5TestShapeInputService5TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	PipelineId *string `location:"uri" type:"string"`

	QueryDoc map[string][]string `location:"querystring" type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService6TestShapeInputService6TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.QueryDoc != nil {
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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService6TestShapeInputService6TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService6TestShapeInputService6TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService6TestShapeInputService6TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	BoolQuery *bool `location:"querystring" locationName:"bool-query" type:"boolean"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService7TestShapeInputService7TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.BoolQuery != nil {
		v := *s.BoolQuery

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "bool-query", protocol.BoolValue(v), metadata)
	}
	return nil
}

type InputService7TestShapeInputService7TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService7TestShapeInputService7TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService7TestShapeInputService7TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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

	BoolQuery *bool `location:"querystring" locationName:"bool-query" type:"boolean"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService7TestShapeInputService7TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.BoolQuery != nil {
		v := *s.BoolQuery

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "bool-query", protocol.BoolValue(v), metadata)
	}
	return nil
}

type InputService7TestShapeInputService7TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService7TestShapeInputService7TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService7TestShapeInputService7TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	Ascending *string `location:"querystring" locationName:"Ascending" type:"string"`

	PageToken *string `location:"querystring" locationName:"PageToken" type:"string"`

	PipelineId *string `location:"uri" locationName:"PipelineId" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService8TestShapeInputService8TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService8TestShapeInputService8TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService8TestShapeInputService8TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService8TestShapeInputService8TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	Ascending *string `location:"querystring" locationName:"Ascending" type:"string"`

	Config *InputService9TestShapeStructType `type:"structure"`

	PageToken *string `location:"querystring" locationName:"PageToken" type:"string"`

	PipelineId *string `location:"uri" locationName:"PipelineId" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService9TestShapeInputService9TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService9TestShapeInputService9TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService9TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService9TestShapeInputService9TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	Ascending *string `location:"querystring" locationName:"Ascending" type:"string"`

	Checksum *string `location:"header" locationName:"x-amz-checksum" type:"string"`

	Config *InputService10TestShapeStructType `type:"structure"`

	PageToken *string `location:"querystring" locationName:"PageToken" type:"string"`

	PipelineId *string `location:"uri" locationName:"PipelineId" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService10TestShapeInputService10TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService10TestShapeInputService10TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService10TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService10TestShapeInputService10TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService10TestShapeInputService10TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService11ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService11TestShapeInputService11TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService11TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/vaults/{vaultName}/archives",
	}

	if input == nil {
		input = &InputService11TestShapeInputService11TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService11TestShapeInputService11TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService12TestShapeInputService12TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService12TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/{Foo}",
	}

	if input == nil {
		input = &InputService12TestShapeInputService12TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService12TestShapeInputService12TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService13ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService13TestShapeInputService13TestCaseOperation1Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo []byte `locationName:"foo" type:"blob"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService13TestShapeInputService13TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetStream(protocol.PayloadTarget, "foo", protocol.BytesStream(v), metadata)
	}
	return nil
}

type InputService13TestShapeInputService13TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService13TestShapeInputService13TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService13TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService13TestShapeInputService13TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService13TestShapeInputService13TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService13TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService13TestShapeInputService13TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService14ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService14TestShapeInputService14TestCaseOperation1Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo *InputService14TestShapeFooShape `locationName:"foo" type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService14TestShapeInputService14TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetFields(protocol.PayloadTarget, "foo", v, metadata)
	}
	return nil
}

type InputService14TestShapeInputService14TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService14TestShapeInputService14TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	_ struct{} `type:"structure" payload:"Foo"`

	Foo *InputService14TestShapeFooShape `locationName:"foo" type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService14TestShapeInputService14TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetFields(protocol.PayloadTarget, "foo", v, metadata)
	}
	return nil
}

type InputService14TestShapeInputService14TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService14TestShapeInputService14TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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

type InputService14TestShapeFooShape struct {
	_ struct{} `locationName:"foo" type:"structure"`

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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	Foo *string `location:"querystring" locationName:"param-name" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService15TestShapeInputService15TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "param-name", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService15TestShapeInputService15TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService15TestShapeInputService15TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService15TestShapeInputService15TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService15TestShapeInputService15TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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

	Foo *string `location:"querystring" locationName:"param-name" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService15TestShapeInputService15TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "param-name", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService15TestShapeInputService15TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService15TestShapeInputService15TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		HTTPPath:   "/path?abc=mno",
	}

	if input == nil {
		input = &InputService15TestShapeInputService15TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService15TestShapeInputService15TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

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

	RecursiveStruct *InputService16TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.RecursiveStruct != nil {
		v := s.RecursiveStruct

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
	}
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService16TestCaseOperation1 = "OperationName"

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
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

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

	RecursiveStruct *InputService16TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.RecursiveStruct != nil {
		v := s.RecursiveStruct

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
	}
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService16TestCaseOperation2 = "OperationName"

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
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

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

type InputService16TestShapeInputService16TestCaseOperation3Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService16TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation3Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.RecursiveStruct != nil {
		v := s.RecursiveStruct

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
	}
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation3Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService16TestCaseOperation3 = "OperationName"

// InputService16TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService16TestCaseOperation3Request.
//    req := client.InputService16TestCaseOperation3Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation3Request(input *InputService16TestShapeInputService16TestCaseOperation3Input) InputService16TestCaseOperation3Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation3,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation3Input{}
	}

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation3Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService16TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation3Request}
}

// InputService16TestCaseOperation3Request is the request type for the
// InputService16TestCaseOperation3 API operation.
type InputService16TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation3Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation3Input) InputService16TestCaseOperation3Request
}

// Send marshals and sends the InputService16TestCaseOperation3 API request.
func (r InputService16TestCaseOperation3Request) Send(ctx context.Context) (*InputService16TestCaseOperation3Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService16TestCaseOperation3Response{
		InputService16TestShapeInputService16TestCaseOperation3Output: r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation3Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService16TestCaseOperation3Response is the response type for the
// InputService16TestCaseOperation3 API operation.
type InputService16TestCaseOperation3Response struct {
	*InputService16TestShapeInputService16TestCaseOperation3Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService16TestCaseOperation3 request.
func (r *InputService16TestCaseOperation3Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService16TestShapeInputService16TestCaseOperation4Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService16TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation4Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.RecursiveStruct != nil {
		v := s.RecursiveStruct

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
	}
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation4Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation4Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService16TestCaseOperation4 = "OperationName"

// InputService16TestCaseOperation4Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService16TestCaseOperation4Request.
//    req := client.InputService16TestCaseOperation4Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation4Request(input *InputService16TestShapeInputService16TestCaseOperation4Input) InputService16TestCaseOperation4Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation4,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation4Input{}
	}

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation4Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService16TestCaseOperation4Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation4Request}
}

// InputService16TestCaseOperation4Request is the request type for the
// InputService16TestCaseOperation4 API operation.
type InputService16TestCaseOperation4Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation4Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation4Input) InputService16TestCaseOperation4Request
}

// Send marshals and sends the InputService16TestCaseOperation4 API request.
func (r InputService16TestCaseOperation4Request) Send(ctx context.Context) (*InputService16TestCaseOperation4Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService16TestCaseOperation4Response{
		InputService16TestShapeInputService16TestCaseOperation4Output: r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation4Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService16TestCaseOperation4Response is the response type for the
// InputService16TestCaseOperation4 API operation.
type InputService16TestCaseOperation4Response struct {
	*InputService16TestShapeInputService16TestCaseOperation4Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService16TestCaseOperation4 request.
func (r *InputService16TestCaseOperation4Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService16TestShapeInputService16TestCaseOperation5Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService16TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation5Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.RecursiveStruct != nil {
		v := s.RecursiveStruct

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
	}
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation5Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation5Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService16TestCaseOperation5 = "OperationName"

// InputService16TestCaseOperation5Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService16TestCaseOperation5Request.
//    req := client.InputService16TestCaseOperation5Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService16ProtocolTest) InputService16TestCaseOperation5Request(input *InputService16TestShapeInputService16TestCaseOperation5Input) InputService16TestCaseOperation5Request {
	op := &aws.Operation{
		Name:       opInputService16TestCaseOperation5,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation5Input{}
	}

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation5Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService16TestCaseOperation5Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation5Request}
}

// InputService16TestCaseOperation5Request is the request type for the
// InputService16TestCaseOperation5 API operation.
type InputService16TestCaseOperation5Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation5Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation5Input) InputService16TestCaseOperation5Request
}

// Send marshals and sends the InputService16TestCaseOperation5 API request.
func (r InputService16TestCaseOperation5Request) Send(ctx context.Context) (*InputService16TestCaseOperation5Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService16TestCaseOperation5Response{
		InputService16TestShapeInputService16TestCaseOperation5Output: r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation5Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService16TestCaseOperation5Response is the response type for the
// InputService16TestCaseOperation5 API operation.
type InputService16TestCaseOperation5Response struct {
	*InputService16TestShapeInputService16TestCaseOperation5Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService16TestCaseOperation5 request.
func (r *InputService16TestCaseOperation5Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService16TestShapeInputService16TestCaseOperation6Input struct {
	_ struct{} `type:"structure"`

	RecursiveStruct *InputService16TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation6Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.RecursiveStruct != nil {
		v := s.RecursiveStruct

		metadata := protocol.Metadata{}
		e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
	}
	return nil
}

type InputService16TestShapeInputService16TestCaseOperation6Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation6Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService16TestCaseOperation6 = "OperationName"

// InputService16TestCaseOperation6Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService16TestCaseOperation6Request.
//    req := client.InputService16TestCaseOperation6Request(params)
//    resp, err := req.Send(context.TODO())
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

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation6Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService16TestCaseOperation6Request{Request: req, Input: input, Copy: c.InputService16TestCaseOperation6Request}
}

// InputService16TestCaseOperation6Request is the request type for the
// InputService16TestCaseOperation6 API operation.
type InputService16TestCaseOperation6Request struct {
	*aws.Request
	Input *InputService16TestShapeInputService16TestCaseOperation6Input
	Copy  func(*InputService16TestShapeInputService16TestCaseOperation6Input) InputService16TestCaseOperation6Request
}

// Send marshals and sends the InputService16TestCaseOperation6 API request.
func (r InputService16TestCaseOperation6Request) Send(ctx context.Context) (*InputService16TestCaseOperation6Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService16TestCaseOperation6Response{
		InputService16TestShapeInputService16TestCaseOperation6Output: r.Request.Data.(*InputService16TestShapeInputService16TestCaseOperation6Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService16TestCaseOperation6Response is the response type for the
// InputService16TestCaseOperation6 API operation.
type InputService16TestCaseOperation6Response struct {
	*InputService16TestShapeInputService16TestCaseOperation6Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService16TestCaseOperation6 request.
func (r *InputService16TestCaseOperation6Response) SDKResponseMetdata() *aws.Response {
	return r.response
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
	if s.RecursiveList != nil {
		v := s.RecursiveList

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.BodyTarget, "RecursiveList", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddFields(v1)
		}
		ls0.End()

	}
	if s.RecursiveMap != nil {
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
// InputService17ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService17ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice17protocoltest.New(myConfig)
func NewInputService17ProtocolTest(config aws.Config) *InputService17ProtocolTest {
	svc := &InputService17ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService17ProtocolTest",
				ServiceID:     "InputService17ProtocolTest",
				EndpointsID:   "inputservice17protocoltest",
				SigningName:   "inputservice17protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService17ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService17TestShapeInputService17TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	TimeArg *time.Time `type:"timestamp"`

	TimeArgInHeader *time.Time `location:"header" locationName:"x-amz-timearg" type:"timestamp"`

	TimeArgInQuery *time.Time `location:"querystring" locationName:"TimeQuery" type:"timestamp"`

	TimeCustom *time.Time `type:"timestamp" timestampFormat:"iso8601"`

	TimeCustomInHeader *time.Time `location:"header" locationName:"x-amz-timecustom-header" type:"timestamp" timestampFormat:"unixTimestamp"`

	TimeCustomInQuery *time.Time `location:"querystring" locationName:"TimeCustomQuery" type:"timestamp" timestampFormat:"unixTimestamp"`

	TimeFormat *time.Time `type:"timestamp" timestampFormat:"rfc822"`

	TimeFormatInHeader *time.Time `location:"header" locationName:"x-amz-timeformat-header" type:"timestamp" timestampFormat:"unixTimestamp"`

	TimeFormatInQuery *time.Time `location:"querystring" locationName:"TimeFormatQuery" type:"timestamp" timestampFormat:"unixTimestamp"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService17TestShapeInputService17TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.TimeArg != nil {
		v := *s.TimeArg

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "TimeArg",
			protocol.TimeValue{V: v, Format: protocol.UnixTimeFormatName, QuotedFormatTime: true}, metadata)
	}
	if s.TimeCustom != nil {
		v := *s.TimeCustom

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "TimeCustom",
			protocol.TimeValue{V: v, Format: "iso8601", QuotedFormatTime: true}, metadata)
	}
	if s.TimeFormat != nil {
		v := *s.TimeFormat

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "TimeFormat",
			protocol.TimeValue{V: v, Format: "rfc822", QuotedFormatTime: true}, metadata)
	}
	if s.TimeArgInHeader != nil {
		v := *s.TimeArgInHeader

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "x-amz-timearg",
			protocol.TimeValue{V: v, Format: protocol.RFC822TimeFormatName, QuotedFormatTime: false}, metadata)
	}
	if s.TimeCustomInHeader != nil {
		v := *s.TimeCustomInHeader

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "x-amz-timecustom-header",
			protocol.TimeValue{V: v, Format: "unixTimestamp", QuotedFormatTime: false}, metadata)
	}
	if s.TimeFormatInHeader != nil {
		v := *s.TimeFormatInHeader

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "x-amz-timeformat-header",
			protocol.TimeValue{V: v, Format: "unixTimestamp", QuotedFormatTime: false}, metadata)
	}
	if s.TimeArgInQuery != nil {
		v := *s.TimeArgInQuery

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "TimeQuery",
			protocol.TimeValue{V: v, Format: protocol.ISO8601TimeFormatName, QuotedFormatTime: false}, metadata)
	}
	if s.TimeCustomInQuery != nil {
		v := *s.TimeCustomInQuery

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "TimeCustomQuery",
			protocol.TimeValue{V: v, Format: "unixTimestamp", QuotedFormatTime: false}, metadata)
	}
	if s.TimeFormatInQuery != nil {
		v := *s.TimeFormatInQuery

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "TimeFormatQuery",
			protocol.TimeValue{V: v, Format: "unixTimestamp", QuotedFormatTime: false}, metadata)
	}
	return nil
}

type InputService17TestShapeInputService17TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService17TestShapeInputService17TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService17TestCaseOperation1 = "OperationName"

// InputService17TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService17TestCaseOperation1Request.
//    req := client.InputService17TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService17ProtocolTest) InputService17TestCaseOperation1Request(input *InputService17TestShapeInputService17TestCaseOperation1Input) InputService17TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService17TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService17TestShapeInputService17TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService17TestShapeInputService17TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService17TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService17TestCaseOperation1Request}
}

// InputService17TestCaseOperation1Request is the request type for the
// InputService17TestCaseOperation1 API operation.
type InputService17TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService17TestShapeInputService17TestCaseOperation1Input
	Copy  func(*InputService17TestShapeInputService17TestCaseOperation1Input) InputService17TestCaseOperation1Request
}

// Send marshals and sends the InputService17TestCaseOperation1 API request.
func (r InputService17TestCaseOperation1Request) Send(ctx context.Context) (*InputService17TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService17TestCaseOperation1Response{
		InputService17TestShapeInputService17TestCaseOperation1Output: r.Request.Data.(*InputService17TestShapeInputService17TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService17TestCaseOperation1Response is the response type for the
// InputService17TestCaseOperation1 API operation.
type InputService17TestCaseOperation1Response struct {
	*InputService17TestShapeInputService17TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService17TestCaseOperation1 request.
func (r *InputService17TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService18ProtocolTest provides the API operation methods for making requests to
// InputService18ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService18ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice18protocoltest.New(myConfig)
func NewInputService18ProtocolTest(config aws.Config) *InputService18ProtocolTest {
	svc := &InputService18ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService18ProtocolTest",
				ServiceID:     "InputService18ProtocolTest",
				EndpointsID:   "inputservice18protocoltest",
				SigningName:   "inputservice18protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService18ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService18TestShapeInputService18TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	TimeArg *time.Time `locationName:"timestamp_location" type:"timestamp"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService18TestShapeInputService18TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.TimeArg != nil {
		v := *s.TimeArg

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "timestamp_location",
			protocol.TimeValue{V: v, Format: protocol.UnixTimeFormatName, QuotedFormatTime: true}, metadata)
	}
	return nil
}

type InputService18TestShapeInputService18TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService18TestShapeInputService18TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService18TestCaseOperation1 = "OperationName"

// InputService18TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService18TestCaseOperation1Request.
//    req := client.InputService18TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
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

	req := c.newRequest(op, input, &InputService18TestShapeInputService18TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService18TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService18TestCaseOperation1Request}
}

// InputService18TestCaseOperation1Request is the request type for the
// InputService18TestCaseOperation1 API operation.
type InputService18TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService18TestShapeInputService18TestCaseOperation1Input
	Copy  func(*InputService18TestShapeInputService18TestCaseOperation1Input) InputService18TestCaseOperation1Request
}

// Send marshals and sends the InputService18TestCaseOperation1 API request.
func (r InputService18TestCaseOperation1Request) Send(ctx context.Context) (*InputService18TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService18TestCaseOperation1Response{
		InputService18TestShapeInputService18TestCaseOperation1Output: r.Request.Data.(*InputService18TestShapeInputService18TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService18TestCaseOperation1Response is the response type for the
// InputService18TestCaseOperation1 API operation.
type InputService18TestCaseOperation1Response struct {
	*InputService18TestShapeInputService18TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService18TestCaseOperation1 request.
func (r *InputService18TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService19ProtocolTest provides the API operation methods for making requests to
// InputService19ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService19ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice19protocoltest.New(myConfig)
func NewInputService19ProtocolTest(config aws.Config) *InputService19ProtocolTest {
	svc := &InputService19ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService19ProtocolTest",
				ServiceID:     "InputService19ProtocolTest",
				EndpointsID:   "inputservice19protocoltest",
				SigningName:   "inputservice19protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService19ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService19TestShapeInputService19TestCaseOperation1Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo *string `locationName:"foo" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetStream(protocol.PayloadTarget, "foo", protocol.StringStream(v), metadata)
	}
	return nil
}

type InputService19TestShapeInputService19TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService19TestCaseOperation1 = "OperationName"

// InputService19TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService19TestCaseOperation1Request.
//    req := client.InputService19TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
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

	req := c.newRequest(op, input, &InputService19TestShapeInputService19TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService19TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService19TestCaseOperation1Request}
}

// InputService19TestCaseOperation1Request is the request type for the
// InputService19TestCaseOperation1 API operation.
type InputService19TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService19TestShapeInputService19TestCaseOperation1Input
	Copy  func(*InputService19TestShapeInputService19TestCaseOperation1Input) InputService19TestCaseOperation1Request
}

// Send marshals and sends the InputService19TestCaseOperation1 API request.
func (r InputService19TestCaseOperation1Request) Send(ctx context.Context) (*InputService19TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService19TestCaseOperation1Response{
		InputService19TestShapeInputService19TestCaseOperation1Output: r.Request.Data.(*InputService19TestShapeInputService19TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService19TestCaseOperation1Response is the response type for the
// InputService19TestCaseOperation1 API operation.
type InputService19TestCaseOperation1Response struct {
	*InputService19TestShapeInputService19TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService19TestCaseOperation1 request.
func (r *InputService19TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService20ProtocolTest provides the API operation methods for making requests to
// InputService20ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService20ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice20protocoltest.New(myConfig)
func NewInputService20ProtocolTest(config aws.Config) *InputService20ProtocolTest {
	svc := &InputService20ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService20ProtocolTest",
				ServiceID:     "InputService20ProtocolTest",
				EndpointsID:   "inputservice20protocoltest",
				SigningName:   "inputservice20protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService20ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService20TestShapeInputService20TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Token *string `type:"string" idempotencyToken:"true"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeInputService20TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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

type InputService20TestShapeInputService20TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeInputService20TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService20TestCaseOperation1 = "OperationName"

// InputService20TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService20TestCaseOperation1Request.
//    req := client.InputService20TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService20ProtocolTest) InputService20TestCaseOperation1Request(input *InputService20TestShapeInputService20TestCaseOperation1Input) InputService20TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService20TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService20TestShapeInputService20TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService20TestShapeInputService20TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService20TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService20TestCaseOperation1Request}
}

// InputService20TestCaseOperation1Request is the request type for the
// InputService20TestCaseOperation1 API operation.
type InputService20TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService20TestShapeInputService20TestCaseOperation1Input
	Copy  func(*InputService20TestShapeInputService20TestCaseOperation1Input) InputService20TestCaseOperation1Request
}

// Send marshals and sends the InputService20TestCaseOperation1 API request.
func (r InputService20TestCaseOperation1Request) Send(ctx context.Context) (*InputService20TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService20TestCaseOperation1Response{
		InputService20TestShapeInputService20TestCaseOperation1Output: r.Request.Data.(*InputService20TestShapeInputService20TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService20TestCaseOperation1Response is the response type for the
// InputService20TestCaseOperation1 API operation.
type InputService20TestCaseOperation1Response struct {
	*InputService20TestShapeInputService20TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService20TestCaseOperation1 request.
func (r *InputService20TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService20TestShapeInputService20TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	Token *string `type:"string" idempotencyToken:"true"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeInputService20TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeInputService20TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService20TestCaseOperation2 = "OperationName"

// InputService20TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService20TestCaseOperation2Request.
//    req := client.InputService20TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
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

	req := c.newRequest(op, input, &InputService20TestShapeInputService20TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService20TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService20TestCaseOperation2Request}
}

// InputService20TestCaseOperation2Request is the request type for the
// InputService20TestCaseOperation2 API operation.
type InputService20TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService20TestShapeInputService20TestCaseOperation2Input
	Copy  func(*InputService20TestShapeInputService20TestCaseOperation2Input) InputService20TestCaseOperation2Request
}

// Send marshals and sends the InputService20TestCaseOperation2 API request.
func (r InputService20TestCaseOperation2Request) Send(ctx context.Context) (*InputService20TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService20TestCaseOperation2Response{
		InputService20TestShapeInputService20TestCaseOperation2Output: r.Request.Data.(*InputService20TestShapeInputService20TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService20TestCaseOperation2Response is the response type for the
// InputService20TestCaseOperation2 API operation.
type InputService20TestCaseOperation2Response struct {
	*InputService20TestShapeInputService20TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService20TestCaseOperation2 request.
func (r *InputService20TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService21ProtocolTest provides the API operation methods for making requests to
// InputService21ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService21ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice21protocoltest.New(myConfig)
func NewInputService21ProtocolTest(config aws.Config) *InputService21ProtocolTest {
	svc := &InputService21ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService21ProtocolTest",
				ServiceID:     "InputService21ProtocolTest",
				EndpointsID:   "inputservice21protocoltest",
				SigningName:   "inputservice21protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService21ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService21TestShapeInputService21TestCaseOperation1Input struct {
	_ struct{} `type:"structure" payload:"Body"`

	Body *InputService21TestShapeBodyStructure `type:"structure"`

	HeaderField aws.JSONValue `location:"header" locationName:"X-Amz-Foo" type:"jsonvalue"`

	QueryField aws.JSONValue `location:"querystring" locationName:"Bar" type:"jsonvalue"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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

type InputService21TestShapeInputService21TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService21TestCaseOperation1 = "OperationName"

// InputService21TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService21TestCaseOperation1Request.
//    req := client.InputService21TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService21ProtocolTest) InputService21TestCaseOperation1Request(input *InputService21TestShapeInputService21TestCaseOperation1Input) InputService21TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService21TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService21TestShapeInputService21TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService21TestShapeInputService21TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService21TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService21TestCaseOperation1Request}
}

// InputService21TestCaseOperation1Request is the request type for the
// InputService21TestCaseOperation1 API operation.
type InputService21TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService21TestShapeInputService21TestCaseOperation1Input
	Copy  func(*InputService21TestShapeInputService21TestCaseOperation1Input) InputService21TestCaseOperation1Request
}

// Send marshals and sends the InputService21TestCaseOperation1 API request.
func (r InputService21TestCaseOperation1Request) Send(ctx context.Context) (*InputService21TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService21TestCaseOperation1Response{
		InputService21TestShapeInputService21TestCaseOperation1Output: r.Request.Data.(*InputService21TestShapeInputService21TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService21TestCaseOperation1Response is the response type for the
// InputService21TestCaseOperation1 API operation.
type InputService21TestCaseOperation1Response struct {
	*InputService21TestShapeInputService21TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService21TestCaseOperation1 request.
func (r *InputService21TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService21TestShapeInputService21TestCaseOperation2Input struct {
	_ struct{} `type:"structure" payload:"Body"`

	Body *InputService21TestShapeBodyStructure `type:"structure"`

	HeaderField aws.JSONValue `location:"header" locationName:"X-Amz-Foo" type:"jsonvalue"`

	QueryField aws.JSONValue `location:"querystring" locationName:"Bar" type:"jsonvalue"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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

type InputService21TestShapeInputService21TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService21TestCaseOperation2 = "OperationName"

// InputService21TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService21TestCaseOperation2Request.
//    req := client.InputService21TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService21ProtocolTest) InputService21TestCaseOperation2Request(input *InputService21TestShapeInputService21TestCaseOperation2Input) InputService21TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService21TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService21TestShapeInputService21TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService21TestShapeInputService21TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService21TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService21TestCaseOperation2Request}
}

// InputService21TestCaseOperation2Request is the request type for the
// InputService21TestCaseOperation2 API operation.
type InputService21TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService21TestShapeInputService21TestCaseOperation2Input
	Copy  func(*InputService21TestShapeInputService21TestCaseOperation2Input) InputService21TestCaseOperation2Request
}

// Send marshals and sends the InputService21TestCaseOperation2 API request.
func (r InputService21TestCaseOperation2Request) Send(ctx context.Context) (*InputService21TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService21TestCaseOperation2Response{
		InputService21TestShapeInputService21TestCaseOperation2Output: r.Request.Data.(*InputService21TestShapeInputService21TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService21TestCaseOperation2Response is the response type for the
// InputService21TestCaseOperation2 API operation.
type InputService21TestCaseOperation2Response struct {
	*InputService21TestShapeInputService21TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService21TestCaseOperation2 request.
func (r *InputService21TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService21TestShapeInputService21TestCaseOperation3Input struct {
	_ struct{} `type:"structure" payload:"Body"`

	Body *InputService21TestShapeBodyStructure `type:"structure"`

	HeaderField aws.JSONValue `location:"header" locationName:"X-Amz-Foo" type:"jsonvalue"`

	QueryField aws.JSONValue `location:"querystring" locationName:"Bar" type:"jsonvalue"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation3Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation3Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService21TestCaseOperation3 = "OperationName"

// InputService21TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService21TestCaseOperation3Request.
//    req := client.InputService21TestCaseOperation3Request(params)
//    resp, err := req.Send(context.TODO())
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

	req := c.newRequest(op, input, &InputService21TestShapeInputService21TestCaseOperation3Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService21TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService21TestCaseOperation3Request}
}

// InputService21TestCaseOperation3Request is the request type for the
// InputService21TestCaseOperation3 API operation.
type InputService21TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService21TestShapeInputService21TestCaseOperation3Input
	Copy  func(*InputService21TestShapeInputService21TestCaseOperation3Input) InputService21TestCaseOperation3Request
}

// Send marshals and sends the InputService21TestCaseOperation3 API request.
func (r InputService21TestCaseOperation3Request) Send(ctx context.Context) (*InputService21TestCaseOperation3Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService21TestCaseOperation3Response{
		InputService21TestShapeInputService21TestCaseOperation3Output: r.Request.Data.(*InputService21TestShapeInputService21TestCaseOperation3Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService21TestCaseOperation3Response is the response type for the
// InputService21TestCaseOperation3 API operation.
type InputService21TestCaseOperation3Response struct {
	*InputService21TestShapeInputService21TestCaseOperation3Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService21TestCaseOperation3 request.
func (r *InputService21TestCaseOperation3Response) SDKResponseMetdata() *aws.Response {
	return r.response
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
	if s.BodyListField != nil {
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

// InputService22ProtocolTest provides the API operation methods for making requests to
// InputService22ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService22ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice22protocoltest.New(myConfig)
func NewInputService22ProtocolTest(config aws.Config) *InputService22ProtocolTest {
	svc := &InputService22ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService22ProtocolTest",
				ServiceID:     "InputService22ProtocolTest",
				EndpointsID:   "inputservice22protocoltest",
				SigningName:   "inputservice22protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService22ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService22TestShapeInputService22TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService22TestShapeEnumType `type:"string" enum:"true"`

	HeaderEnum InputService22TestShapeEnumType `location:"header" locationName:"x-amz-enum" type:"string" enum:"true"`

	ListEnums []InputService22TestShapeEnumType `type:"list"`

	QueryFooEnum InputService22TestShapeEnumType `location:"querystring" locationName:"Enum" type:"string" enum:"true"`

	QueryListEnums []InputService22TestShapeEnumType `location:"querystring" locationName:"List" type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService22TestShapeInputService22TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if len(s.FooEnum) > 0 {
		v := s.FooEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "FooEnum", protocol.QuotedValue{ValueMarshaler: v}, metadata)
	}
	if s.ListEnums != nil {
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
	if s.QueryListEnums != nil {
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

type InputService22TestShapeInputService22TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService22TestShapeInputService22TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService22TestCaseOperation1 = "OperationName"

// InputService22TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService22TestCaseOperation1Request.
//    req := client.InputService22TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService22ProtocolTest) InputService22TestCaseOperation1Request(input *InputService22TestShapeInputService22TestCaseOperation1Input) InputService22TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService22TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService22TestShapeInputService22TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService22TestShapeInputService22TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService22TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService22TestCaseOperation1Request}
}

// InputService22TestCaseOperation1Request is the request type for the
// InputService22TestCaseOperation1 API operation.
type InputService22TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService22TestShapeInputService22TestCaseOperation1Input
	Copy  func(*InputService22TestShapeInputService22TestCaseOperation1Input) InputService22TestCaseOperation1Request
}

// Send marshals and sends the InputService22TestCaseOperation1 API request.
func (r InputService22TestCaseOperation1Request) Send(ctx context.Context) (*InputService22TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService22TestCaseOperation1Response{
		InputService22TestShapeInputService22TestCaseOperation1Output: r.Request.Data.(*InputService22TestShapeInputService22TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService22TestCaseOperation1Response is the response type for the
// InputService22TestCaseOperation1 API operation.
type InputService22TestCaseOperation1Response struct {
	*InputService22TestShapeInputService22TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService22TestCaseOperation1 request.
func (r *InputService22TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
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
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if len(s.FooEnum) > 0 {
		v := s.FooEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "FooEnum", protocol.QuotedValue{ValueMarshaler: v}, metadata)
	}
	if s.ListEnums != nil {
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
	if s.QueryListEnums != nil {
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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService22TestShapeInputService22TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService22TestCaseOperation2 = "OperationName"

// InputService22TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService22TestCaseOperation2Request.
//    req := client.InputService22TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
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

	req := c.newRequest(op, input, &InputService22TestShapeInputService22TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService22TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService22TestCaseOperation2Request}
}

// InputService22TestCaseOperation2Request is the request type for the
// InputService22TestCaseOperation2 API operation.
type InputService22TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService22TestShapeInputService22TestCaseOperation2Input
	Copy  func(*InputService22TestShapeInputService22TestCaseOperation2Input) InputService22TestCaseOperation2Request
}

// Send marshals and sends the InputService22TestCaseOperation2 API request.
func (r InputService22TestCaseOperation2Request) Send(ctx context.Context) (*InputService22TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService22TestCaseOperation2Response{
		InputService22TestShapeInputService22TestCaseOperation2Output: r.Request.Data.(*InputService22TestShapeInputService22TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService22TestCaseOperation2Response is the response type for the
// InputService22TestCaseOperation2 API operation.
type InputService22TestCaseOperation2Response struct {
	*InputService22TestShapeInputService22TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService22TestCaseOperation2 request.
func (r *InputService22TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
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
// InputService23ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService23ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice23protocoltest.New(myConfig)
func NewInputService23ProtocolTest(config aws.Config) *InputService23ProtocolTest {
	svc := &InputService23ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService23ProtocolTest",
				ServiceID:     "InputService23ProtocolTest",
				EndpointsID:   "inputservice23protocoltest",
				SigningName:   "inputservice23protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService23ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService23TestShapeInputService23TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Name *string `type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Name != nil {
		v := *s.Name

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Name", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

type InputService23TestShapeInputService23TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService23TestCaseOperation1 = "StaticOp"

// InputService23TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService23TestCaseOperation1Request.
//    req := client.InputService23TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
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

	req := c.newRequest(op, input, &InputService23TestShapeInputService23TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("data-", nil))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)

	return InputService23TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService23TestCaseOperation1Request}
}

// InputService23TestCaseOperation1Request is the request type for the
// InputService23TestCaseOperation1 API operation.
type InputService23TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService23TestShapeInputService23TestCaseOperation1Input
	Copy  func(*InputService23TestShapeInputService23TestCaseOperation1Input) InputService23TestCaseOperation1Request
}

// Send marshals and sends the InputService23TestCaseOperation1 API request.
func (r InputService23TestCaseOperation1Request) Send(ctx context.Context) (*InputService23TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService23TestCaseOperation1Response{
		InputService23TestShapeInputService23TestCaseOperation1Output: r.Request.Data.(*InputService23TestShapeInputService23TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService23TestCaseOperation1Response is the response type for the
// InputService23TestCaseOperation1 API operation.
type InputService23TestCaseOperation1Response struct {
	*InputService23TestShapeInputService23TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService23TestCaseOperation1 request.
func (r *InputService23TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
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
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService23TestCaseOperation2 = "MemberRefOp"

// InputService23TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService23TestCaseOperation2Request.
//    req := client.InputService23TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
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

	req := c.newRequest(op, input, &InputService23TestShapeInputService23TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("foo-{Name}.", input.hostLabels))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)

	return InputService23TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService23TestCaseOperation2Request}
}

// InputService23TestCaseOperation2Request is the request type for the
// InputService23TestCaseOperation2 API operation.
type InputService23TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService23TestShapeInputService23TestCaseOperation2Input
	Copy  func(*InputService23TestShapeInputService23TestCaseOperation2Input) InputService23TestCaseOperation2Request
}

// Send marshals and sends the InputService23TestCaseOperation2 API request.
func (r InputService23TestCaseOperation2Request) Send(ctx context.Context) (*InputService23TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService23TestCaseOperation2Response{
		InputService23TestShapeInputService23TestCaseOperation2Output: r.Request.Data.(*InputService23TestShapeInputService23TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService23TestCaseOperation2Response is the response type for the
// InputService23TestCaseOperation2 API operation.
type InputService23TestCaseOperation2Response struct {
	*InputService23TestShapeInputService23TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService23TestCaseOperation2 request.
func (r *InputService23TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService24ProtocolTest provides the API operation methods for making requests to
// InputService24ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService24ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice24protocoltest.New(myConfig)
func NewInputService24ProtocolTest(config aws.Config) *InputService24ProtocolTest {
	svc := &InputService24ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService24ProtocolTest",
				ServiceID:     "InputService24ProtocolTest",
				EndpointsID:   "inputservice24protocoltest",
				SigningName:   "inputservice24protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService24ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService24TestShapeInputService24TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Header1 *string `location:"header" type:"string"`

	HeaderMap map[string]string `location:"headers" locationName:"header-map-" type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService24TestShapeInputService24TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Header1 != nil {
		v := *s.Header1

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "Header1", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	if s.HeaderMap != nil {
		v := s.HeaderMap

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.HeadersTarget, "header-map-", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ms0.MapSetValue(k1, protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v1)})
		}
		ms0.End()

	}
	return nil
}

type InputService24TestShapeInputService24TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService24TestShapeInputService24TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService24TestCaseOperation1 = "OperationName"

// InputService24TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService24TestCaseOperation1Request.
//    req := client.InputService24TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService24ProtocolTest) InputService24TestCaseOperation1Request(input *InputService24TestShapeInputService24TestCaseOperation1Input) InputService24TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService24TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService24TestShapeInputService24TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService24TestShapeInputService24TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService24TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService24TestCaseOperation1Request}
}

// InputService24TestCaseOperation1Request is the request type for the
// InputService24TestCaseOperation1 API operation.
type InputService24TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService24TestShapeInputService24TestCaseOperation1Input
	Copy  func(*InputService24TestShapeInputService24TestCaseOperation1Input) InputService24TestCaseOperation1Request
}

// Send marshals and sends the InputService24TestCaseOperation1 API request.
func (r InputService24TestCaseOperation1Request) Send(ctx context.Context) (*InputService24TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService24TestCaseOperation1Response{
		InputService24TestShapeInputService24TestCaseOperation1Output: r.Request.Data.(*InputService24TestShapeInputService24TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService24TestCaseOperation1Response is the response type for the
// InputService24TestCaseOperation1 API operation.
type InputService24TestCaseOperation1Response struct {
	*InputService24TestShapeInputService24TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService24TestCaseOperation1 request.
func (r *InputService24TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService25ProtocolTest provides the API operation methods for making requests to
// InputService25ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService25ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice25protocoltest.New(myConfig)
func NewInputService25ProtocolTest(config aws.Config) *InputService25ProtocolTest {
	svc := &InputService25ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService25ProtocolTest",
				ServiceID:     "InputService25ProtocolTest",
				EndpointsID:   "inputservice25protocoltest",
				SigningName:   "inputservice25protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restjson.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restjson.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restjson.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restjson.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService25ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService25TestShapeInputService25TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	MapDefaults map[string]string `type:"map"`

	MapEmpty map[string]string `type:"map"`

	MapNotSet map[string]string `type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService25TestShapeInputService25TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.MapDefaults != nil {
		v := s.MapDefaults

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.BodyTarget, "MapDefaults", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ms0.MapSetValue(k1, protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v1)})
		}
		ms0.End()

	}
	if s.MapEmpty != nil {
		v := s.MapEmpty

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.BodyTarget, "MapEmpty", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ms0.MapSetValue(k1, protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v1)})
		}
		ms0.End()

	}
	if s.MapNotSet != nil {
		v := s.MapNotSet

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.BodyTarget, "MapNotSet", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ms0.MapSetValue(k1, protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v1)})
		}
		ms0.End()

	}
	return nil
}

type InputService25TestShapeInputService25TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService25TestShapeInputService25TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService25TestCaseOperation1 = "OperationName"

// InputService25TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService25TestCaseOperation1Request.
//    req := client.InputService25TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService25ProtocolTest) InputService25TestCaseOperation1Request(input *InputService25TestShapeInputService25TestCaseOperation1Input) InputService25TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService25TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService25TestShapeInputService25TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService25TestShapeInputService25TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restjson.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService25TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService25TestCaseOperation1Request}
}

// InputService25TestCaseOperation1Request is the request type for the
// InputService25TestCaseOperation1 API operation.
type InputService25TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService25TestShapeInputService25TestCaseOperation1Input
	Copy  func(*InputService25TestShapeInputService25TestCaseOperation1Input) InputService25TestCaseOperation1Request
}

// Send marshals and sends the InputService25TestCaseOperation1 API request.
func (r InputService25TestCaseOperation1Request) Send(ctx context.Context) (*InputService25TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService25TestCaseOperation1Response{
		InputService25TestShapeInputService25TestCaseOperation1Output: r.Request.Data.(*InputService25TestShapeInputService25TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService25TestCaseOperation1Response is the response type for the
// InputService25TestCaseOperation1 API operation.
type InputService25TestCaseOperation1Response struct {
	*InputService25TestShapeInputService25TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService25TestCaseOperation1 request.
func (r *InputService25TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
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
	input := &InputService7TestShapeInputService7TestCaseOperation1Input{
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
	awstesting.AssertJSON(t, `{"Config": {"A": "one", "B": "two"}}`, util.Trim(string(body)))

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
	awstesting.AssertJSON(t, `{"Config": {"A": "one", "B": "two"}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/jobsByPipeline/foo?Ascending=true&PageToken=bar", r.URL.String())

	// assert headers
	if e, a := "12345", r.Header.Get("x-amz-checksum"); e != a {
		t.Errorf("expect %v, got %v", e, a)
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
		t.Errorf("expect %v, got %v", e, a)
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
	awstesting.AssertJSON(t, `{"Bar": "QmxvYiBwYXJhbQ=="}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/foo_name", r.URL.String())

	// assert headers

}

func TestInputService13ProtocolTestBlobPayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation1Input{
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
	input := &InputService14TestShapeInputService14TestCaseOperation1Input{
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
	awstesting.AssertJSON(t, `{"baz": "bar"}`, util.Trim(string(body)))

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
	input := &InputService15TestShapeInputService15TestCaseOperation1Input{}

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
	input := &InputService16TestShapeInputService16TestCaseOperation1Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"NoRecurse": "foo"}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation2Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveStruct": {"NoRecurse": "foo"}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation3Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveStruct": {"RecursiveStruct": {"RecursiveStruct": {"NoRecurse": "foo"}}}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase4(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation4Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveList": [{"NoRecurse": "foo"}, {"NoRecurse": "bar"}]}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService16ProtocolTestRecursiveShapesCase5(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation5Input{
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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveList": [{"NoRecurse": "foo"}, {"RecursiveStruct": {"NoRecurse": "bar"}}]}}`, util.Trim(string(body)))

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
	awstesting.AssertJSON(t, `{"RecursiveStruct": {"RecursiveMap": {"foo": {"NoRecurse": "foo"}, "bar": {"NoRecurse": "bar"}}}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService17ProtocolTestTimestampValuesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService17ProtocolTest(cfg)
	input := &InputService17TestShapeInputService17TestCaseOperation1Input{
		TimeArg:            aws.Time(time.Unix(1422172800, 0)),
		TimeArgInHeader:    aws.Time(time.Unix(1422172800, 0)),
		TimeArgInQuery:     aws.Time(time.Unix(1422172800, 0)),
		TimeCustom:         aws.Time(time.Unix(1422172800, 0)),
		TimeCustomInHeader: aws.Time(time.Unix(1422172800, 0)),
		TimeCustomInQuery:  aws.Time(time.Unix(1422172800, 0)),
		TimeFormat:         aws.Time(time.Unix(1422172800, 0)),
		TimeFormatInHeader: aws.Time(time.Unix(1422172800, 0)),
		TimeFormatInQuery:  aws.Time(time.Unix(1422172800, 0)),
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
	awstesting.AssertJSON(t, `{"TimeArg": 1422172800, "TimeCustom": "2015-01-25T08:00:00Z", "TimeFormat": "Sun, 25 Jan 2015 08:00:00 GMT"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path?TimeQuery=2015-01-25T08%3A00%3A00Z&TimeCustomQuery=1422172800&TimeFormatQuery=1422172800", r.URL.String())

	// assert headers
	if e, a := "application/json", r.Header.Get("Content-Type"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "Sun, 25 Jan 2015 08:00:00 GMT", r.Header.Get("x-amz-timearg"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "1422172800", r.Header.Get("x-amz-timecustom-header"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "1422172800", r.Header.Get("x-amz-timeformat-header"); e != a {
		t.Errorf("expect %v, got %v", e, a)
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
	awstesting.AssertJSON(t, `{"timestamp_location": 1422172800}`, util.Trim(string(body)))

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
	input := &InputService20TestShapeInputService20TestCaseOperation1Input{
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
	awstesting.AssertJSON(t, `{"Token": "abc123"}`, util.Trim(string(body)))

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
	awstesting.AssertJSON(t, `{"Token": "00000000-0000-4000-8000-000000000000"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService21ProtocolTestJSONValueTraitCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService21ProtocolTest(cfg)
	input := &InputService21TestShapeInputService21TestCaseOperation1Input{
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
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService21ProtocolTestJSONValueTraitCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService21ProtocolTest(cfg)
	input := &InputService21TestShapeInputService21TestCaseOperation2Input{
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
	input := &InputService22TestShapeInputService22TestCaseOperation1Input{
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
	awstesting.AssertJSON(t, `{"FooEnum": "foo", "ListEnums": ["foo", "", "bar"]}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/path?Enum=bar&List=0&List=1&List=", r.URL.String())

	// assert headers
	if e, a := "baz", r.Header.Get("x-amz-enum"); e != a {
		t.Errorf("expect %v, got %v", e, a)
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
	awstesting.AssertJSON(t, `{"Name": "myname"}`, util.Trim(string(body)))

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
	awstesting.AssertJSON(t, `{"Name": "myname"}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://foo-myname.service.region.amazonaws.com/path", r.URL.String())

	// assert headers

}

func TestInputService24ProtocolTestHeaderWhitespaceCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService24ProtocolTest(cfg)
	input := &InputService24TestShapeInputService24TestCaseOperation1Input{
		Header1: aws.String("   headerValue"),
		HeaderMap: map[string]string{
			"   key-leading-space": "value",
			"   key-with-space   ": "value",
			"leading-space":        "   value",
			"leading-tab":          "    value",
			"with-space":           "   value   ",
		},
	}

	req := svc.InputService24TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "value", r.Header.Get("header-map-key-leading-space"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "value", r.Header.Get("header-map-key-with-space"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "value", r.Header.Get("header-map-leading-space"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "value", r.Header.Get("header-map-leading-tab"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "value", r.Header.Get("header-map-with-space"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "headerValue", r.Header.Get("header1"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService25ProtocolTestEmptyOrUnsetMapsInPayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService25ProtocolTest(cfg)
	input := &InputService25TestShapeInputService25TestCaseOperation1Input{
		MapDefaults: map[string]string{
			"test_key_1": "test_value_1",
			"test_key_2": "test_value_2",
			"test_key_3": "test_value_3",
		},
		MapEmpty: map[string]string{},
	}

	req := svc.InputService25TestCaseOperation1Request(input)
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
	awstesting.AssertJSON(t, `{"MapDefaults": {"test_key_1":"test_value_1", "test_key_2":"test_value_2", "test_key_3":"test_value_3"}, "MapEmpty":{}}`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}
