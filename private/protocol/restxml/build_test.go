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

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/retry"
	"github.com/jviney/aws-sdk-go-v2/aws/signer/v4"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/private/protocol"
	"github.com/jviney/aws-sdk-go-v2/private/protocol/restxml"
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService1ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService1TestShapeInputService1TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	Description *string `type:"string"`

	Name *string `type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService1TestShapeInputService1TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.Description != nil {
			v := *s.Description

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Description", protocol.StringValue(v), metadata)
		}
		if s.Name != nil {
			v := *s.Name

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Name", protocol.StringValue(v), metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService1TestShapeInputService1TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	Description *string `type:"string"`

	Name *string `type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService1TestShapeInputService1TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.Description != nil {
			v := *s.Description

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Description", protocol.StringValue(v), metadata)
		}
		if s.Name != nil {
			v := *s.Name

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Name", protocol.StringValue(v), metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
	return nil
}

type InputService1TestShapeInputService1TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService1TestShapeInputService1TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService1TestCaseOperation2,
		HTTPMethod: "PUT",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService1TestShapeInputService1TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService1TestShapeInputService1TestCaseOperation3Input) MarshalFields(e protocol.FieldEncoder) error {

	return nil
}

type InputService1TestShapeInputService1TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService1TestShapeInputService1TestCaseOperation3Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService1TestCaseOperation3,
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation3Input{}
	}

	req := c.newRequest(op, input, &InputService1TestShapeInputService1TestCaseOperation3Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService2ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService2TestShapeInputService2TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	First *bool `type:"boolean"`

	Fourth *int64 `type:"integer"`

	Second *bool `type:"boolean"`

	Third *float64 `type:"float"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService2TestShapeInputService2TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.First != nil {
			v := *s.First

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "First", protocol.BoolValue(v), metadata)
		}
		if s.Fourth != nil {
			v := *s.Fourth

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Fourth", protocol.Int64Value(v), metadata)
		}
		if s.Second != nil {
			v := *s.Second

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Second", protocol.BoolValue(v), metadata)
		}
		if s.Third != nil {
			v := *s.Third

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Third", protocol.Float64Value(v), metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService2TestShapeInputService2TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService2TestShapeInputService2TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService3ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService3TestShapeInputService3TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	Description *string `type:"string"`

	SubStructure *InputService3TestShapeSubStructure `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService3TestShapeInputService3TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.Description != nil {
			v := *s.Description

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Description", protocol.StringValue(v), metadata)
		}
		if s.SubStructure != nil {
			v := s.SubStructure

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "SubStructure", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService3TestShapeInputService3TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService3TestShapeInputService3TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	Description *string `type:"string"`

	SubStructure *InputService3TestShapeSubStructure `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService3TestShapeInputService3TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.Description != nil {
			v := *s.Description

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Description", protocol.StringValue(v), metadata)
		}
		if s.SubStructure != nil {
			v := s.SubStructure

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "SubStructure", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
	return nil
}

type InputService3TestShapeInputService3TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService3TestShapeInputService3TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
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
		Name:       opInputService3TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService3TestShapeInputService3TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService3TestShapeInputService3TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

type InputService3TestShapeSubStructure struct {
	_ struct{} `type:"structure"`

	Bar *string `type:"string"`

	Foo *string `type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService3TestShapeSubStructure) MarshalFields(e protocol.FieldEncoder) error {
	if s.Bar != nil {
		v := *s.Bar

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Bar", protocol.StringValue(v), metadata)
	}
	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Foo", protocol.StringValue(v), metadata)
	}
	return nil
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService4ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService4TestShapeInputService4TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	Description *string `type:"string"`

	SubStructure *InputService4TestShapeSubStructure `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService4TestShapeInputService4TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.Description != nil {
			v := *s.Description

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Description", protocol.StringValue(v), metadata)
		}
		if s.SubStructure != nil {
			v := s.SubStructure

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "SubStructure", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService4TestShapeInputService4TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService4TestShapeInputService4TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

type InputService4TestShapeSubStructure struct {
	_ struct{} `type:"structure"`

	Bar *string `type:"string"`

	Foo *string `type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService4TestShapeSubStructure) MarshalFields(e protocol.FieldEncoder) error {
	if s.Bar != nil {
		v := *s.Bar

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Bar", protocol.StringValue(v), metadata)
	}
	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Foo", protocol.StringValue(v), metadata)
	}
	return nil
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService5ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService5TestShapeInputService5TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	ListParam []string `type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService5TestShapeInputService5TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.ListParam != nil {
			v := s.ListParam

			metadata := protocol.Metadata{}
			ls0 := e.List(protocol.BodyTarget, "ListParam", metadata)
			ls0.Start()
			for _, v1 := range v {
				ls0.ListAddValue(protocol.StringValue(v1))
			}
			ls0.End()

		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService5TestShapeInputService5TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService6ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService6TestShapeInputService6TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	ListParam []string `locationName:"AlternateName" locationNameList:"NotMember" type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService6TestShapeInputService6TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.ListParam != nil {
			v := s.ListParam

			metadata := protocol.Metadata{ListLocationName: "NotMember"}
			ls0 := e.List(protocol.BodyTarget, "AlternateName", metadata)
			ls0.Start()
			for _, v1 := range v {
				ls0.ListAddValue(protocol.StringValue(v1))
			}
			ls0.End()

		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService6TestShapeInputService6TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService6TestShapeInputService6TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService7ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService7TestShapeInputService7TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	ListParam []string `type:"list" flattened:"true"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService7TestShapeInputService7TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.ListParam != nil {
			v := s.ListParam

			metadata := protocol.Metadata{Flatten: true}
			ls0 := e.List(protocol.BodyTarget, "ListParam", metadata)
			ls0.Start()
			for _, v1 := range v {
				ls0.ListAddValue(protocol.StringValue(v1))
			}
			ls0.End()

		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService7TestShapeInputService7TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService8ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService8TestShapeInputService8TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	ListParam []string `locationName:"item" type:"list" flattened:"true"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService8TestShapeInputService8TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.ListParam != nil {
			v := s.ListParam

			metadata := protocol.Metadata{Flatten: true}
			ls0 := e.List(protocol.BodyTarget, "item", metadata)
			ls0.Start()
			for _, v1 := range v {
				ls0.ListAddValue(protocol.StringValue(v1))
			}
			ls0.End()

		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPMethod: "POST",
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService8TestShapeInputService8TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService8TestShapeInputService8TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService9ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService9TestShapeInputService9TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	ListParam []InputService9TestShapeSingleFieldStruct `locationName:"item" type:"list" flattened:"true"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService9TestShapeInputService9TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.ListParam != nil {
			v := s.ListParam

			metadata := protocol.Metadata{Flatten: true}
			ls0 := e.List(protocol.BodyTarget, "item", metadata)
			ls0.Start()
			for _, v1 := range v {
				ls0.ListAddFields(v1)
			}
			ls0.End()

		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService9TestShapeInputService9TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

type InputService9TestShapeSingleFieldStruct struct {
	_ struct{} `type:"structure"`

	Element *string `locationName:"value" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService9TestShapeSingleFieldStruct) MarshalFields(e protocol.FieldEncoder) error {
	if s.Element != nil {
		v := *s.Element

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "value", protocol.StringValue(v), metadata)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService10ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService10TestShapeInputService10TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	StructureParam *InputService10TestShapeStructureShape `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService10TestShapeInputService10TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.StructureParam != nil {
			v := s.StructureParam

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "StructureParam", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService10TestShapeInputService10TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService10TestShapeInputService10TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

type InputService10TestShapeStructureShape struct {
	_ struct{} `type:"structure"`

	// B is automatically base64 encoded/decoded by the SDK.
	B []byte `locationName:"b" type:"blob"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService10TestShapeStructureShape) MarshalFields(e protocol.FieldEncoder) error {
	if s.B != nil {
		v := s.B

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "b", protocol.BytesValue(v), metadata)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService11ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService11TestShapeInputService11TestCaseOperation1Input struct {
	_ struct{} `locationName:"TimestampStructure" type:"structure" xmlURI:"https://foo/"`

	TimeArg *time.Time `type:"timestamp"`

	TimeArgInHeader *time.Time `location:"header" locationName:"x-amz-timearg" type:"timestamp"`

	TimeArgInQuery *time.Time `location:"querystring" locationName:"TimeQuery" type:"timestamp"`

	TimeCustom *time.Time `type:"timestamp" timestampFormat:"rfc822"`

	TimeCustomInHeader *time.Time `location:"header" locationName:"x-amz-timecustom-header" type:"timestamp" timestampFormat:"unixTimestamp"`

	TimeCustomInQuery *time.Time `location:"querystring" locationName:"TimeCustomQuery" type:"timestamp" timestampFormat:"unixTimestamp"`

	TimeFormat *time.Time `type:"timestamp" timestampFormat:"rfc822"`

	TimeFormatInHeader *time.Time `location:"header" locationName:"x-amz-timeformat-header" type:"timestamp" timestampFormat:"unixTimestamp"`

	TimeFormatInQuery *time.Time `location:"querystring" locationName:"TimeFormatQuery" type:"timestamp" timestampFormat:"unixTimestamp"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService11TestShapeInputService11TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "TimestampStructure", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.TimeArg != nil {
			v := *s.TimeArg

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "TimeArg",
				protocol.TimeValue{V: v, Format: protocol.ISO8601TimeFormatName, QuotedFormatTime: false}, metadata)
		}
		if s.TimeCustom != nil {
			v := *s.TimeCustom

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "TimeCustom",
				protocol.TimeValue{V: v, Format: "rfc822", QuotedFormatTime: false}, metadata)
		}
		if s.TimeFormat != nil {
			v := *s.TimeFormat

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "TimeFormat",
				protocol.TimeValue{V: v, Format: "rfc822", QuotedFormatTime: false}, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
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
		HTTPPath:   "/2014-01-01/hostedzone",
	}

	if input == nil {
		input = &InputService11TestShapeInputService11TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService11TestShapeInputService11TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService12ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService12TestShapeInputService12TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	Foo map[string]string `location:"headers" locationName:"x-foo-" type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService12TestShapeInputService12TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.HeadersTarget, "x-foo-", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ms0.MapSetValue(k1, protocol.StringValue(v1))
		}
		ms0.End()

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
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService12TestShapeInputService12TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService12TestShapeInputService12TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

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

	Items []string `location:"querystring" locationName:"item" type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService13TestShapeInputService13TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Items != nil {
		v := s.Items

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.QueryTarget, "item", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddValue(protocol.StringValue(v1))
		}
		ls0.End()

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
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService13TestShapeInputService13TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService13TestShapeInputService13TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

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

	PipelineId *string `location:"uri" type:"string"`

	QueryDoc map[string]string `location:"querystring" type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService14TestShapeInputService14TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.StringValue(v), metadata)
	}
	if s.QueryDoc != nil {
		v := s.QueryDoc

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.QueryTarget, "QueryDoc", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ms0.MapSetValue(k1, protocol.StringValue(v1))
		}
		ms0.End()

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
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService14TestShapeInputService14TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService14TestShapeInputService14TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

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

	PipelineId *string `location:"uri" type:"string"`

	QueryDoc map[string][]string `location:"querystring" type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService15TestShapeInputService15TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.PipelineId != nil {
		v := *s.PipelineId

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "PipelineId", protocol.StringValue(v), metadata)
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
				ls1.ListAddValue(protocol.StringValue(v2))
			}
			ls1.End()
		}
		ms0.End()

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
		HTTPMethod: "GET",
		HTTPPath:   "/2014-01-01/jobsByPipeline/{PipelineId}",
	}

	if input == nil {
		input = &InputService15TestShapeInputService15TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService15TestShapeInputService15TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

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

	BoolQuery *bool `location:"querystring" locationName:"bool-query" type:"boolean"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.BoolQuery != nil {
		v := *s.BoolQuery

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "bool-query", protocol.BoolValue(v), metadata)
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
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

	BoolQuery *bool `location:"querystring" locationName:"bool-query" type:"boolean"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService16TestShapeInputService16TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.BoolQuery != nil {
		v := *s.BoolQuery

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "bool-query", protocol.BoolValue(v), metadata)
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
		HTTPMethod: "GET",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService16TestShapeInputService16TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService16TestShapeInputService16TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService17ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService17TestShapeInputService17TestCaseOperation1Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo *string `locationName:"foo" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService17TestShapeInputService17TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetStream(protocol.PayloadTarget, "foo", protocol.StringStream(v), metadata)
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
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService17TestShapeInputService17TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService17TestShapeInputService17TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService18ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService18TestShapeInputService18TestCaseOperation1Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo []byte `locationName:"foo" type:"blob"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService18TestShapeInputService18TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetStream(protocol.PayloadTarget, "foo", protocol.BytesStream(v), metadata)
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
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService18TestShapeInputService18TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService18TestShapeInputService18TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

type InputService18TestShapeInputService18TestCaseOperation2Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo []byte `locationName:"foo" type:"blob"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService18TestShapeInputService18TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetStream(protocol.PayloadTarget, "foo", protocol.BytesStream(v), metadata)
	}
	return nil
}

type InputService18TestShapeInputService18TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService18TestShapeInputService18TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService18TestCaseOperation2 = "OperationName"

// InputService18TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService18TestCaseOperation2Request.
//    req := client.InputService18TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService18ProtocolTest) InputService18TestCaseOperation2Request(input *InputService18TestShapeInputService18TestCaseOperation2Input) InputService18TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService18TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService18TestShapeInputService18TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService18TestShapeInputService18TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService18TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService18TestCaseOperation2Request}
}

// InputService18TestCaseOperation2Request is the request type for the
// InputService18TestCaseOperation2 API operation.
type InputService18TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService18TestShapeInputService18TestCaseOperation2Input
	Copy  func(*InputService18TestShapeInputService18TestCaseOperation2Input) InputService18TestCaseOperation2Request
}

// Send marshals and sends the InputService18TestCaseOperation2 API request.
func (r InputService18TestCaseOperation2Request) Send(ctx context.Context) (*InputService18TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService18TestCaseOperation2Response{
		InputService18TestShapeInputService18TestCaseOperation2Output: r.Request.Data.(*InputService18TestShapeInputService18TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService18TestCaseOperation2Response is the response type for the
// InputService18TestCaseOperation2 API operation.
type InputService18TestCaseOperation2Response struct {
	*InputService18TestShapeInputService18TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService18TestCaseOperation2 request.
func (r *InputService18TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

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

	Foo *InputService19TestShapeFooShape `locationName:"foo" type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetFields(protocol.PayloadTarget, "foo", v, metadata)
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
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

type InputService19TestShapeInputService19TestCaseOperation2Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo *InputService19TestShapeFooShape `locationName:"foo" type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetFields(protocol.PayloadTarget, "foo", v, metadata)
	}
	return nil
}

type InputService19TestShapeInputService19TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService19TestCaseOperation2 = "OperationName"

// InputService19TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService19TestCaseOperation2Request.
//    req := client.InputService19TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService19ProtocolTest) InputService19TestCaseOperation2Request(input *InputService19TestShapeInputService19TestCaseOperation2Input) InputService19TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService19TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService19TestShapeInputService19TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService19TestShapeInputService19TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService19TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService19TestCaseOperation2Request}
}

// InputService19TestCaseOperation2Request is the request type for the
// InputService19TestCaseOperation2 API operation.
type InputService19TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService19TestShapeInputService19TestCaseOperation2Input
	Copy  func(*InputService19TestShapeInputService19TestCaseOperation2Input) InputService19TestCaseOperation2Request
}

// Send marshals and sends the InputService19TestCaseOperation2 API request.
func (r InputService19TestCaseOperation2Request) Send(ctx context.Context) (*InputService19TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService19TestCaseOperation2Response{
		InputService19TestShapeInputService19TestCaseOperation2Output: r.Request.Data.(*InputService19TestShapeInputService19TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService19TestCaseOperation2Response is the response type for the
// InputService19TestCaseOperation2 API operation.
type InputService19TestCaseOperation2Response struct {
	*InputService19TestShapeInputService19TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService19TestCaseOperation2 request.
func (r *InputService19TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService19TestShapeInputService19TestCaseOperation3Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo *InputService19TestShapeFooShape `locationName:"foo" type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation3Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetFields(protocol.PayloadTarget, "foo", v, metadata)
	}
	return nil
}

type InputService19TestShapeInputService19TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation3Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService19TestCaseOperation3 = "OperationName"

// InputService19TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService19TestCaseOperation3Request.
//    req := client.InputService19TestCaseOperation3Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService19ProtocolTest) InputService19TestCaseOperation3Request(input *InputService19TestShapeInputService19TestCaseOperation3Input) InputService19TestCaseOperation3Request {
	op := &aws.Operation{
		Name:       opInputService19TestCaseOperation3,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService19TestShapeInputService19TestCaseOperation3Input{}
	}

	req := c.newRequest(op, input, &InputService19TestShapeInputService19TestCaseOperation3Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService19TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService19TestCaseOperation3Request}
}

// InputService19TestCaseOperation3Request is the request type for the
// InputService19TestCaseOperation3 API operation.
type InputService19TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService19TestShapeInputService19TestCaseOperation3Input
	Copy  func(*InputService19TestShapeInputService19TestCaseOperation3Input) InputService19TestCaseOperation3Request
}

// Send marshals and sends the InputService19TestCaseOperation3 API request.
func (r InputService19TestCaseOperation3Request) Send(ctx context.Context) (*InputService19TestCaseOperation3Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService19TestCaseOperation3Response{
		InputService19TestShapeInputService19TestCaseOperation3Output: r.Request.Data.(*InputService19TestShapeInputService19TestCaseOperation3Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService19TestCaseOperation3Response is the response type for the
// InputService19TestCaseOperation3 API operation.
type InputService19TestCaseOperation3Response struct {
	*InputService19TestShapeInputService19TestCaseOperation3Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService19TestCaseOperation3 request.
func (r *InputService19TestCaseOperation3Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService19TestShapeInputService19TestCaseOperation4Input struct {
	_ struct{} `type:"structure" payload:"Foo"`

	Foo *InputService19TestShapeFooShape `locationName:"foo" type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation4Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := s.Foo

		metadata := protocol.Metadata{}
		e.SetFields(protocol.PayloadTarget, "foo", v, metadata)
	}
	return nil
}

type InputService19TestShapeInputService19TestCaseOperation4Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeInputService19TestCaseOperation4Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService19TestCaseOperation4 = "OperationName"

// InputService19TestCaseOperation4Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService19TestCaseOperation4Request.
//    req := client.InputService19TestCaseOperation4Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService19ProtocolTest) InputService19TestCaseOperation4Request(input *InputService19TestShapeInputService19TestCaseOperation4Input) InputService19TestCaseOperation4Request {
	op := &aws.Operation{
		Name:       opInputService19TestCaseOperation4,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService19TestShapeInputService19TestCaseOperation4Input{}
	}

	req := c.newRequest(op, input, &InputService19TestShapeInputService19TestCaseOperation4Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService19TestCaseOperation4Request{Request: req, Input: input, Copy: c.InputService19TestCaseOperation4Request}
}

// InputService19TestCaseOperation4Request is the request type for the
// InputService19TestCaseOperation4 API operation.
type InputService19TestCaseOperation4Request struct {
	*aws.Request
	Input *InputService19TestShapeInputService19TestCaseOperation4Input
	Copy  func(*InputService19TestShapeInputService19TestCaseOperation4Input) InputService19TestCaseOperation4Request
}

// Send marshals and sends the InputService19TestCaseOperation4 API request.
func (r InputService19TestCaseOperation4Request) Send(ctx context.Context) (*InputService19TestCaseOperation4Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService19TestCaseOperation4Response{
		InputService19TestShapeInputService19TestCaseOperation4Output: r.Request.Data.(*InputService19TestShapeInputService19TestCaseOperation4Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService19TestCaseOperation4Response is the response type for the
// InputService19TestCaseOperation4 API operation.
type InputService19TestCaseOperation4Response struct {
	*InputService19TestShapeInputService19TestCaseOperation4Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService19TestCaseOperation4 request.
func (r *InputService19TestCaseOperation4Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService19TestShapeFooShape struct {
	_ struct{} `locationName:"foo" type:"structure"`

	Baz *string `locationName:"baz" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService19TestShapeFooShape) MarshalFields(e protocol.FieldEncoder) error {
	e.SetFields(protocol.BodyTarget, "foo", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.Baz != nil {
			v := *s.Baz

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "baz", protocol.StringValue(v), metadata)
		}
		return nil
	}), protocol.Metadata{})
	return nil
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService20ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService20TestShapeInputService20TestCaseOperation1Input struct {
	_ struct{} `type:"structure" payload:"Grant"`

	Grant *InputService20TestShapeGrant `locationName:"Grant" type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeInputService20TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Grant != nil {
		v := s.Grant

		metadata := protocol.Metadata{}
		e.SetFields(protocol.PayloadTarget, "Grant", v, metadata)
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
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService20TestShapeInputService20TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService20TestShapeInputService20TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

type InputService20TestShapeGrant struct {
	_ struct{} `locationName:"Grant" type:"structure"`

	Grantee *InputService20TestShapeGrantee `type:"structure" xmlPrefix:"xsi" xmlURI:"http://www.w3.org/2001/XMLSchema-instance"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeGrant) MarshalFields(e protocol.FieldEncoder) error {
	e.SetFields(protocol.BodyTarget, "Grant", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.Grantee != nil {
			v := s.Grantee
			attrs := make([]protocol.Attribute, 0, 1)

			if s.Grantee.Type != nil {
				v := *s.Grantee.Type
				attrs = append(attrs, protocol.Attribute{Name: "xsi:type", Value: protocol.StringValue(v), Meta: protocol.Metadata{}})
			}
			metadata := protocol.Metadata{Attributes: attrs, XMLNamespacePrefix: "xsi", XMLNamespaceURI: "http://www.w3.org/2001/XMLSchema-instance"}
			e.SetFields(protocol.BodyTarget, "Grantee", v, metadata)
		}
		return nil
	}), protocol.Metadata{})
	return nil
}

type InputService20TestShapeGrantee struct {
	_ struct{} `type:"structure" xmlPrefix:"xsi" xmlURI:"http://www.w3.org/2001/XMLSchema-instance"`

	EmailAddress *string `type:"string"`

	Type *string `locationName:"xsi:type" type:"string" xmlAttribute:"true"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService20TestShapeGrantee) MarshalFields(e protocol.FieldEncoder) error {
	if s.EmailAddress != nil {
		v := *s.EmailAddress

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "EmailAddress", protocol.StringValue(v), metadata)
	}
	// Skipping Type XML Attribute.
	return nil
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService21ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService21TestShapeInputService21TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Bucket *string `location:"uri" type:"string"`

	Key *string `location:"uri" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService21TestShapeInputService21TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Bucket != nil {
		v := *s.Bucket

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "Bucket", protocol.StringValue(v), metadata)
	}
	if s.Key != nil {
		v := *s.Key

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "Key", protocol.StringValue(v), metadata)
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
		HTTPMethod: "GET",
		HTTPPath:   "/{Bucket}/{Key+}",
	}

	if input == nil {
		input = &InputService21TestShapeInputService21TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService21TestShapeInputService21TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

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

	Foo *string `location:"querystring" locationName:"param-name" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService22TestShapeInputService22TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "param-name", protocol.StringValue(v), metadata)
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
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

	Foo *string `location:"querystring" locationName:"param-name" type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService22TestShapeInputService22TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Foo != nil {
		v := *s.Foo

		metadata := protocol.Metadata{}
		e.SetValue(protocol.QueryTarget, "param-name", protocol.StringValue(v), metadata)
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
		HTTPPath:   "/path?abc=mno",
	}

	if input == nil {
		input = &InputService22TestShapeInputService22TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService22TestShapeInputService22TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService23ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService23TestShapeInputService23TestCaseOperation1Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	RecursiveStruct *InputService23TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.RecursiveStruct != nil {
			v := s.RecursiveStruct

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
	return nil
}

type InputService23TestShapeInputService23TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService23TestCaseOperation1 = "OperationName"

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
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

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
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	RecursiveStruct *InputService23TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.RecursiveStruct != nil {
			v := s.RecursiveStruct

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
	return nil
}

type InputService23TestShapeInputService23TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService23TestCaseOperation2 = "OperationName"

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
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

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

type InputService23TestShapeInputService23TestCaseOperation3Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	RecursiveStruct *InputService23TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation3Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.RecursiveStruct != nil {
			v := s.RecursiveStruct

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
	return nil
}

type InputService23TestShapeInputService23TestCaseOperation3Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation3Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService23TestCaseOperation3 = "OperationName"

// InputService23TestCaseOperation3Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService23TestCaseOperation3Request.
//    req := client.InputService23TestCaseOperation3Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService23ProtocolTest) InputService23TestCaseOperation3Request(input *InputService23TestShapeInputService23TestCaseOperation3Input) InputService23TestCaseOperation3Request {
	op := &aws.Operation{
		Name:       opInputService23TestCaseOperation3,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService23TestShapeInputService23TestCaseOperation3Input{}
	}

	req := c.newRequest(op, input, &InputService23TestShapeInputService23TestCaseOperation3Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService23TestCaseOperation3Request{Request: req, Input: input, Copy: c.InputService23TestCaseOperation3Request}
}

// InputService23TestCaseOperation3Request is the request type for the
// InputService23TestCaseOperation3 API operation.
type InputService23TestCaseOperation3Request struct {
	*aws.Request
	Input *InputService23TestShapeInputService23TestCaseOperation3Input
	Copy  func(*InputService23TestShapeInputService23TestCaseOperation3Input) InputService23TestCaseOperation3Request
}

// Send marshals and sends the InputService23TestCaseOperation3 API request.
func (r InputService23TestCaseOperation3Request) Send(ctx context.Context) (*InputService23TestCaseOperation3Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService23TestCaseOperation3Response{
		InputService23TestShapeInputService23TestCaseOperation3Output: r.Request.Data.(*InputService23TestShapeInputService23TestCaseOperation3Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService23TestCaseOperation3Response is the response type for the
// InputService23TestCaseOperation3 API operation.
type InputService23TestCaseOperation3Response struct {
	*InputService23TestShapeInputService23TestCaseOperation3Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService23TestCaseOperation3 request.
func (r *InputService23TestCaseOperation3Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService23TestShapeInputService23TestCaseOperation4Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	RecursiveStruct *InputService23TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation4Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.RecursiveStruct != nil {
			v := s.RecursiveStruct

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
	return nil
}

type InputService23TestShapeInputService23TestCaseOperation4Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation4Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService23TestCaseOperation4 = "OperationName"

// InputService23TestCaseOperation4Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService23TestCaseOperation4Request.
//    req := client.InputService23TestCaseOperation4Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService23ProtocolTest) InputService23TestCaseOperation4Request(input *InputService23TestShapeInputService23TestCaseOperation4Input) InputService23TestCaseOperation4Request {
	op := &aws.Operation{
		Name:       opInputService23TestCaseOperation4,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService23TestShapeInputService23TestCaseOperation4Input{}
	}

	req := c.newRequest(op, input, &InputService23TestShapeInputService23TestCaseOperation4Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService23TestCaseOperation4Request{Request: req, Input: input, Copy: c.InputService23TestCaseOperation4Request}
}

// InputService23TestCaseOperation4Request is the request type for the
// InputService23TestCaseOperation4 API operation.
type InputService23TestCaseOperation4Request struct {
	*aws.Request
	Input *InputService23TestShapeInputService23TestCaseOperation4Input
	Copy  func(*InputService23TestShapeInputService23TestCaseOperation4Input) InputService23TestCaseOperation4Request
}

// Send marshals and sends the InputService23TestCaseOperation4 API request.
func (r InputService23TestCaseOperation4Request) Send(ctx context.Context) (*InputService23TestCaseOperation4Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService23TestCaseOperation4Response{
		InputService23TestShapeInputService23TestCaseOperation4Output: r.Request.Data.(*InputService23TestShapeInputService23TestCaseOperation4Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService23TestCaseOperation4Response is the response type for the
// InputService23TestCaseOperation4 API operation.
type InputService23TestCaseOperation4Response struct {
	*InputService23TestShapeInputService23TestCaseOperation4Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService23TestCaseOperation4 request.
func (r *InputService23TestCaseOperation4Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService23TestShapeInputService23TestCaseOperation5Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	RecursiveStruct *InputService23TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation5Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.RecursiveStruct != nil {
			v := s.RecursiveStruct

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
	return nil
}

type InputService23TestShapeInputService23TestCaseOperation5Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation5Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService23TestCaseOperation5 = "OperationName"

// InputService23TestCaseOperation5Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService23TestCaseOperation5Request.
//    req := client.InputService23TestCaseOperation5Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService23ProtocolTest) InputService23TestCaseOperation5Request(input *InputService23TestShapeInputService23TestCaseOperation5Input) InputService23TestCaseOperation5Request {
	op := &aws.Operation{
		Name:       opInputService23TestCaseOperation5,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService23TestShapeInputService23TestCaseOperation5Input{}
	}

	req := c.newRequest(op, input, &InputService23TestShapeInputService23TestCaseOperation5Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService23TestCaseOperation5Request{Request: req, Input: input, Copy: c.InputService23TestCaseOperation5Request}
}

// InputService23TestCaseOperation5Request is the request type for the
// InputService23TestCaseOperation5 API operation.
type InputService23TestCaseOperation5Request struct {
	*aws.Request
	Input *InputService23TestShapeInputService23TestCaseOperation5Input
	Copy  func(*InputService23TestShapeInputService23TestCaseOperation5Input) InputService23TestCaseOperation5Request
}

// Send marshals and sends the InputService23TestCaseOperation5 API request.
func (r InputService23TestCaseOperation5Request) Send(ctx context.Context) (*InputService23TestCaseOperation5Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService23TestCaseOperation5Response{
		InputService23TestShapeInputService23TestCaseOperation5Output: r.Request.Data.(*InputService23TestShapeInputService23TestCaseOperation5Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService23TestCaseOperation5Response is the response type for the
// InputService23TestCaseOperation5 API operation.
type InputService23TestCaseOperation5Response struct {
	*InputService23TestShapeInputService23TestCaseOperation5Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService23TestCaseOperation5 request.
func (r *InputService23TestCaseOperation5Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService23TestShapeInputService23TestCaseOperation6Input struct {
	_ struct{} `locationName:"OperationRequest" type:"structure" xmlURI:"https://foo/"`

	RecursiveStruct *InputService23TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation6Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "OperationRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.RecursiveStruct != nil {
			v := s.RecursiveStruct

			metadata := protocol.Metadata{}
			e.SetFields(protocol.BodyTarget, "RecursiveStruct", v, metadata)
		}
		return nil
	}), protocol.Metadata{XMLNamespaceURI: "https://foo/"})
	return nil
}

type InputService23TestShapeInputService23TestCaseOperation6Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeInputService23TestCaseOperation6Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService23TestCaseOperation6 = "OperationName"

// InputService23TestCaseOperation6Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService23TestCaseOperation6Request.
//    req := client.InputService23TestCaseOperation6Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService23ProtocolTest) InputService23TestCaseOperation6Request(input *InputService23TestShapeInputService23TestCaseOperation6Input) InputService23TestCaseOperation6Request {
	op := &aws.Operation{
		Name:       opInputService23TestCaseOperation6,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService23TestShapeInputService23TestCaseOperation6Input{}
	}

	req := c.newRequest(op, input, &InputService23TestShapeInputService23TestCaseOperation6Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService23TestCaseOperation6Request{Request: req, Input: input, Copy: c.InputService23TestCaseOperation6Request}
}

// InputService23TestCaseOperation6Request is the request type for the
// InputService23TestCaseOperation6 API operation.
type InputService23TestCaseOperation6Request struct {
	*aws.Request
	Input *InputService23TestShapeInputService23TestCaseOperation6Input
	Copy  func(*InputService23TestShapeInputService23TestCaseOperation6Input) InputService23TestCaseOperation6Request
}

// Send marshals and sends the InputService23TestCaseOperation6 API request.
func (r InputService23TestCaseOperation6Request) Send(ctx context.Context) (*InputService23TestCaseOperation6Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService23TestCaseOperation6Response{
		InputService23TestShapeInputService23TestCaseOperation6Output: r.Request.Data.(*InputService23TestShapeInputService23TestCaseOperation6Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService23TestCaseOperation6Response is the response type for the
// InputService23TestCaseOperation6 API operation.
type InputService23TestCaseOperation6Response struct {
	*InputService23TestShapeInputService23TestCaseOperation6Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService23TestCaseOperation6 request.
func (r *InputService23TestCaseOperation6Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService23TestShapeRecursiveStructType struct {
	_ struct{} `type:"structure"`

	NoRecurse *string `type:"string"`

	RecursiveList []InputService23TestShapeRecursiveStructType `type:"list"`

	RecursiveMap map[string]InputService23TestShapeRecursiveStructType `type:"map"`

	RecursiveStruct *InputService23TestShapeRecursiveStructType `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService23TestShapeRecursiveStructType) MarshalFields(e protocol.FieldEncoder) error {
	if s.NoRecurse != nil {
		v := *s.NoRecurse

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "NoRecurse", protocol.StringValue(v), metadata)
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

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

	Token *string `type:"string" idempotencyToken:"true"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService24TestShapeInputService24TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	var Token string
	if s.Token != nil {
		Token = *s.Token
	} else {
		Token = protocol.GetIdempotencyToken()
	}
	{
		v := Token

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Token", protocol.StringValue(v), metadata)
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
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService24TestShapeInputService24TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService24TestShapeInputService24TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

type InputService24TestShapeInputService24TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	Token *string `type:"string" idempotencyToken:"true"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService24TestShapeInputService24TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	var Token string
	if s.Token != nil {
		Token = *s.Token
	} else {
		Token = protocol.GetIdempotencyToken()
	}
	{
		v := Token

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "Token", protocol.StringValue(v), metadata)
	}
	return nil
}

type InputService24TestShapeInputService24TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService24TestShapeInputService24TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService24TestCaseOperation2 = "OperationName"

// InputService24TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService24TestCaseOperation2Request.
//    req := client.InputService24TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService24ProtocolTest) InputService24TestCaseOperation2Request(input *InputService24TestShapeInputService24TestCaseOperation2Input) InputService24TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService24TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService24TestShapeInputService24TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService24TestShapeInputService24TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService24TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService24TestCaseOperation2Request}
}

// InputService24TestCaseOperation2Request is the request type for the
// InputService24TestCaseOperation2 API operation.
type InputService24TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService24TestShapeInputService24TestCaseOperation2Input
	Copy  func(*InputService24TestShapeInputService24TestCaseOperation2Input) InputService24TestCaseOperation2Request
}

// Send marshals and sends the InputService24TestCaseOperation2 API request.
func (r InputService24TestCaseOperation2Request) Send(ctx context.Context) (*InputService24TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService24TestCaseOperation2Response{
		InputService24TestShapeInputService24TestCaseOperation2Output: r.Request.Data.(*InputService24TestShapeInputService24TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService24TestCaseOperation2Response is the response type for the
// InputService24TestCaseOperation2 API operation.
type InputService24TestCaseOperation2Response struct {
	*InputService24TestShapeInputService24TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService24TestCaseOperation2 request.
func (r *InputService24TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

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

	FooEnum InputService25TestShapeEnumType `type:"string" enum:"true"`

	HeaderEnum InputService25TestShapeEnumType `location:"header" locationName:"x-amz-enum" type:"string" enum:"true"`

	ListEnums []InputService25TestShapeEnumType `type:"list"`

	URIFooEnum InputService25TestShapeEnumType `location:"uri" locationName:"URIEnum" type:"string" enum:"true"`

	URIListEnums []InputService25TestShapeEnumType `location:"querystring" locationName:"ListEnums" type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService25TestShapeInputService25TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if len(s.FooEnum) > 0 {
		v := s.FooEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "FooEnum", v, metadata)
	}
	if s.ListEnums != nil {
		v := s.ListEnums

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.BodyTarget, "ListEnums", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddValue(protocol.StringValue(v1))
		}
		ls0.End()

	}
	if len(s.HeaderEnum) > 0 {
		v := s.HeaderEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "x-amz-enum", v, metadata)
	}
	if len(s.URIFooEnum) > 0 {
		v := s.URIFooEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "URIEnum", v, metadata)
	}
	if s.URIListEnums != nil {
		v := s.URIListEnums

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.QueryTarget, "ListEnums", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddValue(protocol.StringValue(v1))
		}
		ls0.End()

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
		HTTPPath:   "/Enum/{URIEnum}",
	}

	if input == nil {
		input = &InputService25TestShapeInputService25TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService25TestShapeInputService25TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
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

type InputService25TestShapeInputService25TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService25TestShapeEnumType `type:"string" enum:"true"`

	HeaderEnum InputService25TestShapeEnumType `location:"header" locationName:"x-amz-enum" type:"string" enum:"true"`

	ListEnums []InputService25TestShapeEnumType `type:"list"`

	URIFooEnum InputService25TestShapeEnumType `location:"uri" locationName:"URIEnum" type:"string" enum:"true"`

	URIListEnums []InputService25TestShapeEnumType `location:"querystring" locationName:"ListEnums" type:"list"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService25TestShapeInputService25TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	if len(s.FooEnum) > 0 {
		v := s.FooEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.BodyTarget, "FooEnum", v, metadata)
	}
	if s.ListEnums != nil {
		v := s.ListEnums

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.BodyTarget, "ListEnums", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddValue(protocol.StringValue(v1))
		}
		ls0.End()

	}
	if len(s.HeaderEnum) > 0 {
		v := s.HeaderEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "x-amz-enum", v, metadata)
	}
	if len(s.URIFooEnum) > 0 {
		v := s.URIFooEnum

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "URIEnum", v, metadata)
	}
	if s.URIListEnums != nil {
		v := s.URIListEnums

		metadata := protocol.Metadata{}
		ls0 := e.List(protocol.QueryTarget, "ListEnums", metadata)
		ls0.Start()
		for _, v1 := range v {
			ls0.ListAddValue(protocol.StringValue(v1))
		}
		ls0.End()

	}
	return nil
}

type InputService25TestShapeInputService25TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService25TestShapeInputService25TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService25TestCaseOperation2 = "OperationName"

// InputService25TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService25TestCaseOperation2Request.
//    req := client.InputService25TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService25ProtocolTest) InputService25TestCaseOperation2Request(input *InputService25TestShapeInputService25TestCaseOperation2Input) InputService25TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService25TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService25TestShapeInputService25TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService25TestShapeInputService25TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService25TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService25TestCaseOperation2Request}
}

// InputService25TestCaseOperation2Request is the request type for the
// InputService25TestCaseOperation2 API operation.
type InputService25TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService25TestShapeInputService25TestCaseOperation2Input
	Copy  func(*InputService25TestShapeInputService25TestCaseOperation2Input) InputService25TestCaseOperation2Request
}

// Send marshals and sends the InputService25TestCaseOperation2 API request.
func (r InputService25TestCaseOperation2Request) Send(ctx context.Context) (*InputService25TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService25TestCaseOperation2Response{
		InputService25TestShapeInputService25TestCaseOperation2Output: r.Request.Data.(*InputService25TestShapeInputService25TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService25TestCaseOperation2Response is the response type for the
// InputService25TestCaseOperation2 API operation.
type InputService25TestCaseOperation2Response struct {
	*InputService25TestShapeInputService25TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService25TestCaseOperation2 request.
func (r *InputService25TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService25TestShapeEnumType string

// Enum values for InputService25TestShapeEnumType
const (
	EnumTypeFoo InputService25TestShapeEnumType = "foo"
	EnumTypeBar InputService25TestShapeEnumType = "bar"
	EnumType0   InputService25TestShapeEnumType = "0"
	EnumType1   InputService25TestShapeEnumType = "1"
)

func (enum InputService25TestShapeEnumType) MarshalValue() (string, error) {
	return string(enum), nil
}

func (enum InputService25TestShapeEnumType) MarshalValueBuf(b []byte) ([]byte, error) {
	b = b[0:0]
	return append(b, enum...), nil
}

// InputService26ProtocolTest provides the API operation methods for making requests to
// InputService26ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService26ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice26protocoltest.New(myConfig)
func NewInputService26ProtocolTest(config aws.Config) *InputService26ProtocolTest {
	svc := &InputService26ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService26ProtocolTest",
				ServiceID:     "InputService26ProtocolTest",
				EndpointsID:   "inputservice26protocoltest",
				SigningName:   "inputservice26protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService26ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService26TestShapeInputService26TestCaseOperation1Input struct {
	_ struct{} `locationName:"StaticOpRequest" type:"structure"`

	Name *string `type:"string"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService26TestShapeInputService26TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "StaticOpRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.Name != nil {
			v := *s.Name

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Name", protocol.StringValue(v), metadata)
		}
		return nil
	}), protocol.Metadata{})
	return nil
}

type InputService26TestShapeInputService26TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService26TestShapeInputService26TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService26TestCaseOperation1 = "StaticOp"

// InputService26TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService26TestCaseOperation1Request.
//    req := client.InputService26TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService26ProtocolTest) InputService26TestCaseOperation1Request(input *InputService26TestShapeInputService26TestCaseOperation1Input) InputService26TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService26TestCaseOperation1,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService26TestShapeInputService26TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService26TestShapeInputService26TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("data-", nil))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)

	return InputService26TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService26TestCaseOperation1Request}
}

// InputService26TestCaseOperation1Request is the request type for the
// InputService26TestCaseOperation1 API operation.
type InputService26TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService26TestShapeInputService26TestCaseOperation1Input
	Copy  func(*InputService26TestShapeInputService26TestCaseOperation1Input) InputService26TestCaseOperation1Request
}

// Send marshals and sends the InputService26TestCaseOperation1 API request.
func (r InputService26TestCaseOperation1Request) Send(ctx context.Context) (*InputService26TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService26TestCaseOperation1Response{
		InputService26TestShapeInputService26TestCaseOperation1Output: r.Request.Data.(*InputService26TestShapeInputService26TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService26TestCaseOperation1Response is the response type for the
// InputService26TestCaseOperation1 API operation.
type InputService26TestCaseOperation1Response struct {
	*InputService26TestShapeInputService26TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService26TestCaseOperation1 request.
func (r *InputService26TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

type InputService26TestShapeInputService26TestCaseOperation2Input struct {
	_ struct{} `locationName:"MemberRefOpRequest" type:"structure"`

	// Name is a required field
	Name *string `type:"string" required:"true"`
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *InputService26TestShapeInputService26TestCaseOperation2Input) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "InputService26TestShapeInputService26TestCaseOperation2Input"}

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
func (s InputService26TestShapeInputService26TestCaseOperation2Input) MarshalFields(e protocol.FieldEncoder) error {

	e.SetFields(protocol.BodyTarget, "MemberRefOpRequest", protocol.FieldMarshalerFunc(func(e protocol.FieldEncoder) error {
		if s.Name != nil {
			v := *s.Name

			metadata := protocol.Metadata{}
			e.SetValue(protocol.BodyTarget, "Name", protocol.StringValue(v), metadata)
		}
		return nil
	}), protocol.Metadata{})
	return nil
}

func (s *InputService26TestShapeInputService26TestCaseOperation2Input) hostLabels() map[string]string {
	return map[string]string{
		"Name": aws.StringValue(s.Name),
	}
}

type InputService26TestShapeInputService26TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService26TestShapeInputService26TestCaseOperation2Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService26TestCaseOperation2 = "MemberRefOp"

// InputService26TestCaseOperation2Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService26TestCaseOperation2Request.
//    req := client.InputService26TestCaseOperation2Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService26ProtocolTest) InputService26TestCaseOperation2Request(input *InputService26TestShapeInputService26TestCaseOperation2Input) InputService26TestCaseOperation2Request {
	op := &aws.Operation{
		Name:       opInputService26TestCaseOperation2,
		HTTPMethod: "POST",
		HTTPPath:   "/path",
	}

	if input == nil {
		input = &InputService26TestShapeInputService26TestCaseOperation2Input{}
	}

	req := c.newRequest(op, input, &InputService26TestShapeInputService26TestCaseOperation2Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	req.Handlers.Build.PushBackNamed(protocol.NewHostPrefixHandler("foo-{Name}.", input.hostLabels))
	req.Handlers.Build.PushBackNamed(protocol.ValidateEndpointHostHandler)

	return InputService26TestCaseOperation2Request{Request: req, Input: input, Copy: c.InputService26TestCaseOperation2Request}
}

// InputService26TestCaseOperation2Request is the request type for the
// InputService26TestCaseOperation2 API operation.
type InputService26TestCaseOperation2Request struct {
	*aws.Request
	Input *InputService26TestShapeInputService26TestCaseOperation2Input
	Copy  func(*InputService26TestShapeInputService26TestCaseOperation2Input) InputService26TestCaseOperation2Request
}

// Send marshals and sends the InputService26TestCaseOperation2 API request.
func (r InputService26TestCaseOperation2Request) Send(ctx context.Context) (*InputService26TestCaseOperation2Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService26TestCaseOperation2Response{
		InputService26TestShapeInputService26TestCaseOperation2Output: r.Request.Data.(*InputService26TestShapeInputService26TestCaseOperation2Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService26TestCaseOperation2Response is the response type for the
// InputService26TestCaseOperation2 API operation.
type InputService26TestCaseOperation2Response struct {
	*InputService26TestShapeInputService26TestCaseOperation2Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService26TestCaseOperation2 request.
func (r *InputService26TestCaseOperation2Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

// InputService27ProtocolTest provides the API operation methods for making requests to
// InputService27ProtocolTest. See this package's package overview docs
// for details on the service.
//
// The client's methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type InputService27ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the client from the provided Config.
//
// Example:
//     // Create a client from just a config.
//     svc := inputservice27protocoltest.New(myConfig)
func NewInputService27ProtocolTest(config aws.Config) *InputService27ProtocolTest {
	svc := &InputService27ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "InputService27ProtocolTest",
				ServiceID:     "InputService27ProtocolTest",
				EndpointsID:   "inputservice27protocoltest",
				SigningName:   "inputservice27protocoltest",
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
	svc.Handlers.Build.PushBackNamed(restxml.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(restxml.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(restxml.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(restxml.UnmarshalErrorHandler)

	return svc
}

// newRequest creates a new request for a client operation and runs any
// custom request initialization.
func (c *InputService27ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

type InputService27TestShapeInputService27TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Header1 *string `location:"header" type:"string"`

	HeaderMap map[string]string `location:"headers" locationName:"header-map-" type:"map"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService27TestShapeInputService27TestCaseOperation1Input) MarshalFields(e protocol.FieldEncoder) error {

	if s.Header1 != nil {
		v := *s.Header1

		metadata := protocol.Metadata{}
		e.SetValue(protocol.HeaderTarget, "Header1", protocol.StringValue(v), metadata)
	}
	if s.HeaderMap != nil {
		v := s.HeaderMap

		metadata := protocol.Metadata{}
		ms0 := e.Map(protocol.HeadersTarget, "header-map-", metadata)
		ms0.Start()
		for k1, v1 := range v {
			ms0.MapSetValue(k1, protocol.StringValue(v1))
		}
		ms0.End()

	}
	return nil
}

type InputService27TestShapeInputService27TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s InputService27TestShapeInputService27TestCaseOperation1Output) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opInputService27TestCaseOperation1 = "OperationName"

// InputService27TestCaseOperation1Request returns a request value for making API operation for
// .
//
//    // Example sending a request using InputService27TestCaseOperation1Request.
//    req := client.InputService27TestCaseOperation1Request(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *InputService27ProtocolTest) InputService27TestCaseOperation1Request(input *InputService27TestShapeInputService27TestCaseOperation1Input) InputService27TestCaseOperation1Request {
	op := &aws.Operation{
		Name:       opInputService27TestCaseOperation1,
		HTTPMethod: "GET",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &InputService27TestShapeInputService27TestCaseOperation1Input{}
	}

	req := c.newRequest(op, input, &InputService27TestShapeInputService27TestCaseOperation1Output{})
	req.Handlers.Unmarshal.Remove(restxml.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)

	return InputService27TestCaseOperation1Request{Request: req, Input: input, Copy: c.InputService27TestCaseOperation1Request}
}

// InputService27TestCaseOperation1Request is the request type for the
// InputService27TestCaseOperation1 API operation.
type InputService27TestCaseOperation1Request struct {
	*aws.Request
	Input *InputService27TestShapeInputService27TestCaseOperation1Input
	Copy  func(*InputService27TestShapeInputService27TestCaseOperation1Input) InputService27TestCaseOperation1Request
}

// Send marshals and sends the InputService27TestCaseOperation1 API request.
func (r InputService27TestCaseOperation1Request) Send(ctx context.Context) (*InputService27TestCaseOperation1Response, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &InputService27TestCaseOperation1Response{
		InputService27TestShapeInputService27TestCaseOperation1Output: r.Request.Data.(*InputService27TestShapeInputService27TestCaseOperation1Output),
		response: &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// InputService27TestCaseOperation1Response is the response type for the
// InputService27TestCaseOperation1 API operation.
type InputService27TestCaseOperation1Response struct {
	*InputService27TestShapeInputService27TestCaseOperation1Output

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// InputService27TestCaseOperation1 request.
func (r *InputService27TestCaseOperation1Response) SDKResponseMetdata() *aws.Response {
	return r.response
}

//
// Tests begin here
//

func TestInputService1ProtocolTestBasicXMLSerializationCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService1ProtocolTest(cfg)
	input := &InputService1TestShapeInputService1TestCaseOperation1Input{
		Description: aws.String("bar"),
		Name:        aws.String("foo"),
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><Description xmlns="https://foo/">bar</Description><Name xmlns="https://foo/">foo</Name></OperationRequest>`, util.Trim(string(body)), InputService1TestShapeInputService1TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService1ProtocolTestBasicXMLSerializationCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService1ProtocolTest(cfg)
	input := &InputService1TestShapeInputService1TestCaseOperation2Input{
		Description: aws.String("bar"),
		Name:        aws.String("foo"),
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><Description xmlns="https://foo/">bar</Description><Name xmlns="https://foo/">foo</Name></OperationRequest>`, util.Trim(string(body)), InputService1TestShapeInputService1TestCaseOperation2Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService1ProtocolTestBasicXMLSerializationCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService1ProtocolTest(cfg)
	input := &InputService1TestShapeInputService1TestCaseOperation3Input{}

	req := svc.InputService1TestCaseOperation3Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService2ProtocolTestSerializeOtherScalarTypesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService2ProtocolTest(cfg)
	input := &InputService2TestShapeInputService2TestCaseOperation1Input{
		First:  aws.Bool(true),
		Fourth: aws.Int64(3),
		Second: aws.Bool(false),
		Third:  aws.Float64(1.2),
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><First xmlns="https://foo/">true</First><Fourth xmlns="https://foo/">3</Fourth><Second xmlns="https://foo/">false</Second><Third xmlns="https://foo/">1.2</Third></OperationRequest>`, util.Trim(string(body)), InputService2TestShapeInputService2TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService3ProtocolTestNestedStructuresCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService3ProtocolTest(cfg)
	input := &InputService3TestShapeInputService3TestCaseOperation1Input{
		Description: aws.String("baz"),
		SubStructure: &InputService3TestShapeSubStructure{
			Bar: aws.String("b"),
			Foo: aws.String("a"),
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><Description xmlns="https://foo/">baz</Description><SubStructure xmlns="https://foo/"><Bar xmlns="https://foo/">b</Bar><Foo xmlns="https://foo/">a</Foo></SubStructure></OperationRequest>`, util.Trim(string(body)), InputService3TestShapeInputService3TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService3ProtocolTestNestedStructuresCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService3ProtocolTest(cfg)
	input := &InputService3TestShapeInputService3TestCaseOperation2Input{
		Description: aws.String("baz"),
		SubStructure: &InputService3TestShapeSubStructure{
			Foo: aws.String("a"),
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><Description xmlns="https://foo/">baz</Description><SubStructure xmlns="https://foo/"><Foo xmlns="https://foo/">a</Foo></SubStructure></OperationRequest>`, util.Trim(string(body)), InputService3TestShapeInputService3TestCaseOperation2Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService4ProtocolTestNestedStructuresCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService4ProtocolTest(cfg)
	input := &InputService4TestShapeInputService4TestCaseOperation1Input{
		Description:  aws.String("baz"),
		SubStructure: &InputService4TestShapeSubStructure{},
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><Description xmlns="https://foo/">baz</Description><SubStructure xmlns="https://foo/"></SubStructure></OperationRequest>`, util.Trim(string(body)), InputService4TestShapeInputService4TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService5ProtocolTestNonFlattenedListsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation1Input{
		ListParam: []string{
			"one",
			"two",
			"three",
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><ListParam xmlns="https://foo/"><member xmlns="https://foo/">one</member><member xmlns="https://foo/">two</member><member xmlns="https://foo/">three</member></ListParam></OperationRequest>`, util.Trim(string(body)), InputService5TestShapeInputService5TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService6ProtocolTestNonFlattenedListsWithLocationNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService6ProtocolTest(cfg)
	input := &InputService6TestShapeInputService6TestCaseOperation1Input{
		ListParam: []string{
			"one",
			"two",
			"three",
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><AlternateName xmlns="https://foo/"><NotMember xmlns="https://foo/">one</NotMember><NotMember xmlns="https://foo/">two</NotMember><NotMember xmlns="https://foo/">three</NotMember></AlternateName></OperationRequest>`, util.Trim(string(body)), InputService6TestShapeInputService6TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService7ProtocolTestFlattenedListsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService7ProtocolTest(cfg)
	input := &InputService7TestShapeInputService7TestCaseOperation1Input{
		ListParam: []string{
			"one",
			"two",
			"three",
		},
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><ListParam xmlns="https://foo/">one</ListParam><ListParam xmlns="https://foo/">two</ListParam><ListParam xmlns="https://foo/">three</ListParam></OperationRequest>`, util.Trim(string(body)), InputService7TestShapeInputService7TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService8ProtocolTestFlattenedListsWithLocationNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService8ProtocolTest(cfg)
	input := &InputService8TestShapeInputService8TestCaseOperation1Input{
		ListParam: []string{
			"one",
			"two",
			"three",
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><item xmlns="https://foo/">one</item><item xmlns="https://foo/">two</item><item xmlns="https://foo/">three</item></OperationRequest>`, util.Trim(string(body)), InputService8TestShapeInputService8TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService9ProtocolTestListOfStructuresCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService9ProtocolTest(cfg)
	input := &InputService9TestShapeInputService9TestCaseOperation1Input{
		ListParam: []InputService9TestShapeSingleFieldStruct{
			{
				Element: aws.String("one"),
			},
			{
				Element: aws.String("two"),
			},
			{
				Element: aws.String("three"),
			},
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><item xmlns="https://foo/"><value xmlns="https://foo/">one</value></item><item xmlns="https://foo/"><value xmlns="https://foo/">two</value></item><item xmlns="https://foo/"><value xmlns="https://foo/">three</value></item></OperationRequest>`, util.Trim(string(body)), InputService9TestShapeInputService9TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService10ProtocolTestBlobShapesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService10ProtocolTest(cfg)
	input := &InputService10TestShapeInputService10TestCaseOperation1Input{
		StructureParam: &InputService10TestShapeStructureShape{
			B: []byte("foo"),
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><StructureParam xmlns="https://foo/"><b xmlns="https://foo/">Zm9v</b></StructureParam></OperationRequest>`, util.Trim(string(body)), InputService10TestShapeInputService10TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone", r.URL.String())

	// assert headers

}

func TestInputService11ProtocolTestTimestampShapesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService11ProtocolTest(cfg)
	input := &InputService11TestShapeInputService11TestCaseOperation1Input{
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<TimestampStructure xmlns="https://foo/"><TimeArg xmlns="https://foo/">2015-01-25T08:00:00Z</TimeArg><TimeCustom xmlns="https://foo/">Sun, 25 Jan 2015 08:00:00 GMT</TimeCustom><TimeFormat xmlns="https://foo/">Sun, 25 Jan 2015 08:00:00 GMT</TimeFormat></TimestampStructure>`, util.Trim(string(body)), InputService11TestShapeInputService11TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/2014-01-01/hostedzone?TimeQuery=2015-01-25T08%3A00%3A00Z&TimeCustomQuery=1422172800&TimeFormatQuery=1422172800", r.URL.String())

	// assert headers
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

func TestInputService12ProtocolTestHeaderMapsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService12ProtocolTest(cfg)
	input := &InputService12TestShapeInputService12TestCaseOperation1Input{
		Foo: map[string]string{
			"a": "b",
			"c": "d",
		},
	}

	req := svc.InputService12TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers
	if e, a := "b", r.Header.Get("x-foo-a"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "d", r.Header.Get("x-foo-c"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService13ProtocolTestQuerystringListOfStringsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService13ProtocolTest(cfg)
	input := &InputService13TestShapeInputService13TestCaseOperation1Input{
		Items: []string{
			"value1",
			"value2",
		},
	}

	req := svc.InputService13TestCaseOperation1Request(input)
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

func TestInputService14ProtocolTestStringToStringMapsInQuerystringCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService14ProtocolTest(cfg)
	input := &InputService14TestShapeInputService14TestCaseOperation1Input{
		PipelineId: aws.String("foo"),
		QueryDoc: map[string]string{
			"bar":  "baz",
			"fizz": "buzz",
		},
	}

	req := svc.InputService14TestCaseOperation1Request(input)
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

func TestInputService15ProtocolTestStringToStringListMapsInQuerystringCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService15ProtocolTest(cfg)
	input := &InputService15TestShapeInputService15TestCaseOperation1Input{
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

	req := svc.InputService15TestCaseOperation1Request(input)
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

func TestInputService16ProtocolTestBooleanInQuerystringCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation1Input{
		BoolQuery: aws.Bool(true),
	}

	req := svc.InputService16TestCaseOperation1Request(input)
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

func TestInputService16ProtocolTestBooleanInQuerystringCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService16ProtocolTest(cfg)
	input := &InputService16TestShapeInputService16TestCaseOperation2Input{
		BoolQuery: aws.Bool(false),
	}

	req := svc.InputService16TestCaseOperation2Request(input)
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

func TestInputService17ProtocolTestStringPayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService17ProtocolTest(cfg)
	input := &InputService17TestShapeInputService17TestCaseOperation1Input{
		Foo: aws.String("bar"),
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
	body := util.SortXML(r.Body)
	if e, a := "bar", util.Trim(string(body)); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService18ProtocolTestBlobPayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService18ProtocolTest(cfg)
	input := &InputService18TestShapeInputService18TestCaseOperation1Input{
		Foo: []byte("bar"),
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
	body := util.SortXML(r.Body)
	if e, a := "bar", util.Trim(string(body)); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService18ProtocolTestBlobPayloadCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService18ProtocolTest(cfg)
	input := &InputService18TestShapeInputService18TestCaseOperation2Input{}

	req := svc.InputService18TestCaseOperation2Request(input)
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

func TestInputService19ProtocolTestStructurePayloadCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService19ProtocolTest(cfg)
	input := &InputService19TestShapeInputService19TestCaseOperation1Input{
		Foo: &InputService19TestShapeFooShape{
			Baz: aws.String("bar"),
		},
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<foo><baz>bar</baz></foo>`, util.Trim(string(body)), InputService19TestShapeInputService19TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService19ProtocolTestStructurePayloadCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService19ProtocolTest(cfg)
	input := &InputService19TestShapeInputService19TestCaseOperation2Input{}

	req := svc.InputService19TestCaseOperation2Request(input)
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

func TestInputService19ProtocolTestStructurePayloadCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService19ProtocolTest(cfg)
	input := &InputService19TestShapeInputService19TestCaseOperation3Input{
		Foo: &InputService19TestShapeFooShape{},
	}

	req := svc.InputService19TestCaseOperation3Request(input)
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<foo></foo>`, util.Trim(string(body)), InputService19TestShapeInputService19TestCaseOperation3Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService19ProtocolTestStructurePayloadCase4(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService19ProtocolTest(cfg)
	input := &InputService19TestShapeInputService19TestCaseOperation4Input{}

	req := svc.InputService19TestCaseOperation4Request(input)
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

func TestInputService20ProtocolTestXMLAttributeCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService20ProtocolTest(cfg)
	input := &InputService20TestShapeInputService20TestCaseOperation1Input{
		Grant: &InputService20TestShapeGrant{
			Grantee: &InputService20TestShapeGrantee{
				EmailAddress: aws.String("foo@example.com"),
				Type:         aws.String("CanonicalUser"),
			},
		},
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<Grant xmlns:_xmlns="xmlns" _xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:XMLSchema-instance="http://www.w3.org/2001/XMLSchema-instance" XMLSchema-instance:type="CanonicalUser"><Grantee><EmailAddress>foo@example.com</EmailAddress></Grantee></Grant>`, util.Trim(string(body)), InputService20TestShapeInputService20TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService21ProtocolTestGreedyKeysCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService21ProtocolTest(cfg)
	input := &InputService21TestShapeInputService21TestCaseOperation1Input{
		Bucket: aws.String("my/bucket"),
		Key:    aws.String("testing /123"),
	}

	req := svc.InputService21TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	req.Build()
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert URL
	awstesting.AssertURL(t, "https://test/my%2Fbucket/testing%20/123", r.URL.String())

	// assert headers

}

func TestInputService22ProtocolTestOmitsNullQueryParamsButSerializesEmptyStringsCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService22ProtocolTest(cfg)
	input := &InputService22TestShapeInputService22TestCaseOperation1Input{}

	req := svc.InputService22TestCaseOperation1Request(input)
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

func TestInputService22ProtocolTestOmitsNullQueryParamsButSerializesEmptyStringsCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService22ProtocolTest(cfg)
	input := &InputService22TestShapeInputService22TestCaseOperation2Input{
		Foo: aws.String(""),
	}

	req := svc.InputService22TestCaseOperation2Request(input)
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

func TestInputService23ProtocolTestRecursiveShapesCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService23ProtocolTest(cfg)
	input := &InputService23TestShapeInputService23TestCaseOperation1Input{
		RecursiveStruct: &InputService23TestShapeRecursiveStructType{
			NoRecurse: aws.String("foo"),
		},
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><NoRecurse xmlns="https://foo/">foo</NoRecurse></RecursiveStruct></OperationRequest>`, util.Trim(string(body)), InputService23TestShapeInputService23TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService23ProtocolTestRecursiveShapesCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService23ProtocolTest(cfg)
	input := &InputService23TestShapeInputService23TestCaseOperation2Input{
		RecursiveStruct: &InputService23TestShapeRecursiveStructType{
			RecursiveStruct: &InputService23TestShapeRecursiveStructType{
				NoRecurse: aws.String("foo"),
			},
		},
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><NoRecurse xmlns="https://foo/">foo</NoRecurse></RecursiveStruct></RecursiveStruct></OperationRequest>`, util.Trim(string(body)), InputService23TestShapeInputService23TestCaseOperation2Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService23ProtocolTestRecursiveShapesCase3(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService23ProtocolTest(cfg)
	input := &InputService23TestShapeInputService23TestCaseOperation3Input{
		RecursiveStruct: &InputService23TestShapeRecursiveStructType{
			RecursiveStruct: &InputService23TestShapeRecursiveStructType{
				RecursiveStruct: &InputService23TestShapeRecursiveStructType{
					RecursiveStruct: &InputService23TestShapeRecursiveStructType{
						NoRecurse: aws.String("foo"),
					},
				},
			},
		},
	}

	req := svc.InputService23TestCaseOperation3Request(input)
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><NoRecurse xmlns="https://foo/">foo</NoRecurse></RecursiveStruct></RecursiveStruct></RecursiveStruct></RecursiveStruct></OperationRequest>`, util.Trim(string(body)), InputService23TestShapeInputService23TestCaseOperation3Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService23ProtocolTestRecursiveShapesCase4(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService23ProtocolTest(cfg)
	input := &InputService23TestShapeInputService23TestCaseOperation4Input{
		RecursiveStruct: &InputService23TestShapeRecursiveStructType{
			RecursiveList: []InputService23TestShapeRecursiveStructType{
				{
					NoRecurse: aws.String("foo"),
				},
				{
					NoRecurse: aws.String("bar"),
				},
			},
		},
	}

	req := svc.InputService23TestCaseOperation4Request(input)
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><RecursiveList xmlns="https://foo/"><member xmlns="https://foo/"><NoRecurse xmlns="https://foo/">foo</NoRecurse></member><member xmlns="https://foo/"><NoRecurse xmlns="https://foo/">bar</NoRecurse></member></RecursiveList></RecursiveStruct></OperationRequest>`, util.Trim(string(body)), InputService23TestShapeInputService23TestCaseOperation4Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService23ProtocolTestRecursiveShapesCase5(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService23ProtocolTest(cfg)
	input := &InputService23TestShapeInputService23TestCaseOperation5Input{
		RecursiveStruct: &InputService23TestShapeRecursiveStructType{
			RecursiveList: []InputService23TestShapeRecursiveStructType{
				{
					NoRecurse: aws.String("foo"),
				},
				{
					RecursiveStruct: &InputService23TestShapeRecursiveStructType{
						NoRecurse: aws.String("bar"),
					},
				},
			},
		},
	}

	req := svc.InputService23TestCaseOperation5Request(input)
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><RecursiveList xmlns="https://foo/"><member xmlns="https://foo/"><NoRecurse xmlns="https://foo/">foo</NoRecurse></member><member xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><NoRecurse xmlns="https://foo/">bar</NoRecurse></RecursiveStruct></member></RecursiveList></RecursiveStruct></OperationRequest>`, util.Trim(string(body)), InputService23TestShapeInputService23TestCaseOperation5Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService23ProtocolTestRecursiveShapesCase6(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService23ProtocolTest(cfg)
	input := &InputService23TestShapeInputService23TestCaseOperation6Input{
		RecursiveStruct: &InputService23TestShapeRecursiveStructType{
			RecursiveMap: map[string]InputService23TestShapeRecursiveStructType{
				"bar": {
					NoRecurse: aws.String("bar"),
				},
				"foo": {
					NoRecurse: aws.String("foo"),
				},
			},
		},
	}

	req := svc.InputService23TestCaseOperation6Request(input)
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<OperationRequest xmlns="https://foo/"><RecursiveStruct xmlns="https://foo/"><RecursiveMap xmlns="https://foo/"><entry xmlns="https://foo/"><key xmlns="https://foo/">foo</key><value xmlns="https://foo/"><NoRecurse xmlns="https://foo/">foo</NoRecurse></value></entry><entry xmlns="https://foo/"><key xmlns="https://foo/">bar</key><value xmlns="https://foo/"><NoRecurse xmlns="https://foo/">bar</NoRecurse></value></entry></RecursiveMap></RecursiveStruct></OperationRequest>`, util.Trim(string(body)), InputService23TestShapeInputService23TestCaseOperation6Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService24ProtocolTestIdempotencyTokenAutoFillCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService24ProtocolTest(cfg)
	input := &InputService24TestShapeInputService24TestCaseOperation1Input{
		Token: aws.String("abc123"),
	}

	req := svc.InputService24TestCaseOperation1Request(input)
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<InputShape><Token>abc123</Token></InputShape>`, util.Trim(string(body)), InputService24TestShapeInputService24TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService24ProtocolTestIdempotencyTokenAutoFillCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService24ProtocolTest(cfg)
	input := &InputService24TestShapeInputService24TestCaseOperation2Input{}

	req := svc.InputService24TestCaseOperation2Request(input)
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<InputShape><Token>00000000-0000-4000-8000-000000000000</Token></InputShape>`, util.Trim(string(body)), InputService24TestShapeInputService24TestCaseOperation2Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/path", r.URL.String())

	// assert headers

}

func TestInputService25ProtocolTestEnumCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService25ProtocolTest(cfg)
	input := &InputService25TestShapeInputService25TestCaseOperation1Input{
		FooEnum:    InputService25TestShapeEnumType("foo"),
		HeaderEnum: InputService25TestShapeEnumType("baz"),
		ListEnums: []InputService25TestShapeEnumType{
			InputService25TestShapeEnumType("foo"),
			InputService25TestShapeEnumType(""),
			InputService25TestShapeEnumType("bar"),
		},
		URIFooEnum: InputService25TestShapeEnumType("bar"),
		URIListEnums: []InputService25TestShapeEnumType{
			InputService25TestShapeEnumType("0"),
			InputService25TestShapeEnumType(""),
			InputService25TestShapeEnumType("1"),
		},
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<InputShape><FooEnum>foo</FooEnum><ListEnums><member>foo</member><member></member><member>bar</member></ListEnums></InputShape>`, util.Trim(string(body)), InputService25TestShapeInputService25TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://test/Enum/bar?ListEnums=0&ListEnums=&ListEnums=1", r.URL.String())

	// assert headers
	if e, a := "baz", r.Header.Get("x-amz-enum"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestInputService25ProtocolTestEnumCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService25ProtocolTest(cfg)
	input := &InputService25TestShapeInputService25TestCaseOperation2Input{}

	req := svc.InputService25TestCaseOperation2Request(input)
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

func TestInputService26ProtocolTestEndpointHostTraitCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService26ProtocolTest(cfg)
	input := &InputService26TestShapeInputService26TestCaseOperation1Input{
		Name: aws.String("myname"),
	}

	req := svc.InputService26TestCaseOperation1Request(input)
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<StaticOpRequest><Name>myname</Name></StaticOpRequest>`, util.Trim(string(body)), InputService26TestShapeInputService26TestCaseOperation1Input{})

	// assert URL
	awstesting.AssertURL(t, "https://data-service.region.amazonaws.com/path", r.URL.String())

	// assert headers

}

func TestInputService26ProtocolTestEndpointHostTraitCase2(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://service.region.amazonaws.com")

	svc := NewInputService26ProtocolTest(cfg)
	input := &InputService26TestShapeInputService26TestCaseOperation2Input{
		Name: aws.String("myname"),
	}

	req := svc.InputService26TestCaseOperation2Request(input)
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
	body := util.SortXML(r.Body)
	awstesting.AssertXML(t, `<MemberRefOpRequest><Name>myname</Name></MemberRefOpRequest>`, util.Trim(string(body)), InputService26TestShapeInputService26TestCaseOperation2Input{})

	// assert URL
	awstesting.AssertURL(t, "https://foo-myname.service.region.amazonaws.com/path", r.URL.String())

	// assert headers

}

func TestInputService27ProtocolTestHeaderWhitespaceCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService27ProtocolTest(cfg)
	input := &InputService27TestShapeInputService27TestCaseOperation1Input{
		Header1: aws.String("   headerValue"),
		HeaderMap: map[string]string{
			"   key-leading-space": "value",
			"   key-with-space   ": "value",
			"leading-space":        "   value",
			"leading-tab":          "    value",
			"with-space":           "   value   ",
		},
	}

	req := svc.InputService27TestCaseOperation1Request(input)
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
