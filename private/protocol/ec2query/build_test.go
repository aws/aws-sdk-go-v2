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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService1ProtocolTest client from just a config.
//     svc := inputservice1protocoltest.New(myConfig)
//
//     // Create a InputService1ProtocolTest client with additional configuration
//     svc := inputservice1protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService1TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService1TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService1TestCaseOperation1 for more information on using the InputService1TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService1TestCaseOperation1Request method.
//    req, resp := client.InputService1TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService1ProtocolTest) InputService1TestCaseOperation1Request(input *InputService1TestShapeInputService1TestCaseOperation1Input) (req *aws.Request, output *InputService1TestShapeInputService1TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService1TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService1TestShapeInputService1TestCaseOperation1Input{}
	}

	output = &InputService1TestShapeInputService1TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService1TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService1TestCaseOperation1 for usage and error information.
func (c *InputService1ProtocolTest) InputService1TestCaseOperation1(input *InputService1TestShapeInputService1TestCaseOperation1Input) (*InputService1TestShapeInputService1TestCaseOperation1Output, error) {
	req, out := c.InputService1TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService1TestCaseOperation1WithContext is the same as InputService1TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService1TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService1ProtocolTest) InputService1TestCaseOperation1WithContext(ctx aws.Context, input *InputService1TestShapeInputService1TestCaseOperation1Input, opts ...aws.Option) (*InputService1TestShapeInputService1TestCaseOperation1Output, error) {
	req, out := c.InputService1TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService1TestShapeInputService1TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Bar *string `type:"string"`

	Foo *string `type:"string"`
}

// SetBar sets the Bar field's value.
func (s *InputService1TestShapeInputService1TestCaseOperation1Input) SetBar(v string) *InputService1TestShapeInputService1TestCaseOperation1Input {
	s.Bar = &v
	return s
}

// SetFoo sets the Foo field's value.
func (s *InputService1TestShapeInputService1TestCaseOperation1Input) SetFoo(v string) *InputService1TestShapeInputService1TestCaseOperation1Input {
	s.Foo = &v
	return s
}

type InputService1TestShapeInputService1TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService2ProtocolTest client from just a config.
//     svc := inputservice2protocoltest.New(myConfig)
//
//     // Create a InputService2ProtocolTest client with additional configuration
//     svc := inputservice2protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService2TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService2TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService2TestCaseOperation1 for more information on using the InputService2TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService2TestCaseOperation1Request method.
//    req, resp := client.InputService2TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService2ProtocolTest) InputService2TestCaseOperation1Request(input *InputService2TestShapeInputService2TestCaseOperation1Input) (req *aws.Request, output *InputService2TestShapeInputService2TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService2TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService2TestShapeInputService2TestCaseOperation1Input{}
	}

	output = &InputService2TestShapeInputService2TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService2TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService2TestCaseOperation1 for usage and error information.
func (c *InputService2ProtocolTest) InputService2TestCaseOperation1(input *InputService2TestShapeInputService2TestCaseOperation1Input) (*InputService2TestShapeInputService2TestCaseOperation1Output, error) {
	req, out := c.InputService2TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService2TestCaseOperation1WithContext is the same as InputService2TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService2TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService2ProtocolTest) InputService2TestCaseOperation1WithContext(ctx aws.Context, input *InputService2TestShapeInputService2TestCaseOperation1Input, opts ...aws.Option) (*InputService2TestShapeInputService2TestCaseOperation1Output, error) {
	req, out := c.InputService2TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService2TestShapeInputService2TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	Bar *string `locationName:"barLocationName" type:"string"`

	Foo *string `type:"string"`

	Yuck *string `locationName:"yuckLocationName" queryName:"yuckQueryName" type:"string"`
}

// SetBar sets the Bar field's value.
func (s *InputService2TestShapeInputService2TestCaseOperation1Input) SetBar(v string) *InputService2TestShapeInputService2TestCaseOperation1Input {
	s.Bar = &v
	return s
}

// SetFoo sets the Foo field's value.
func (s *InputService2TestShapeInputService2TestCaseOperation1Input) SetFoo(v string) *InputService2TestShapeInputService2TestCaseOperation1Input {
	s.Foo = &v
	return s
}

// SetYuck sets the Yuck field's value.
func (s *InputService2TestShapeInputService2TestCaseOperation1Input) SetYuck(v string) *InputService2TestShapeInputService2TestCaseOperation1Input {
	s.Yuck = &v
	return s
}

type InputService2TestShapeInputService2TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService3ProtocolTest client from just a config.
//     svc := inputservice3protocoltest.New(myConfig)
//
//     // Create a InputService3ProtocolTest client with additional configuration
//     svc := inputservice3protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService3TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService3TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService3TestCaseOperation1 for more information on using the InputService3TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService3TestCaseOperation1Request method.
//    req, resp := client.InputService3TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService3ProtocolTest) InputService3TestCaseOperation1Request(input *InputService3TestShapeInputService3TestCaseOperation1Input) (req *aws.Request, output *InputService3TestShapeInputService3TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService3TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService3TestShapeInputService3TestCaseOperation1Input{}
	}

	output = &InputService3TestShapeInputService3TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService3TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService3TestCaseOperation1 for usage and error information.
func (c *InputService3ProtocolTest) InputService3TestCaseOperation1(input *InputService3TestShapeInputService3TestCaseOperation1Input) (*InputService3TestShapeInputService3TestCaseOperation1Output, error) {
	req, out := c.InputService3TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService3TestCaseOperation1WithContext is the same as InputService3TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService3TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService3ProtocolTest) InputService3TestCaseOperation1WithContext(ctx aws.Context, input *InputService3TestShapeInputService3TestCaseOperation1Input, opts ...aws.Option) (*InputService3TestShapeInputService3TestCaseOperation1Output, error) {
	req, out := c.InputService3TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService3TestShapeInputService3TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	StructArg *InputService3TestShapeStructType `locationName:"Struct" type:"structure"`
}

// SetStructArg sets the StructArg field's value.
func (s *InputService3TestShapeInputService3TestCaseOperation1Input) SetStructArg(v *InputService3TestShapeStructType) *InputService3TestShapeInputService3TestCaseOperation1Input {
	s.StructArg = v
	return s
}

type InputService3TestShapeInputService3TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

type InputService3TestShapeStructType struct {
	_ struct{} `type:"structure"`

	ScalarArg *string `locationName:"Scalar" type:"string"`
}

// SetScalarArg sets the ScalarArg field's value.
func (s *InputService3TestShapeStructType) SetScalarArg(v string) *InputService3TestShapeStructType {
	s.ScalarArg = &v
	return s
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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService4ProtocolTest client from just a config.
//     svc := inputservice4protocoltest.New(myConfig)
//
//     // Create a InputService4ProtocolTest client with additional configuration
//     svc := inputservice4protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService4TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService4TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService4TestCaseOperation1 for more information on using the InputService4TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService4TestCaseOperation1Request method.
//    req, resp := client.InputService4TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService4ProtocolTest) InputService4TestCaseOperation1Request(input *InputService4TestShapeInputService4TestCaseOperation1Input) (req *aws.Request, output *InputService4TestShapeInputService4TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService4TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService4TestShapeInputService4TestCaseOperation1Input{}
	}

	output = &InputService4TestShapeInputService4TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService4TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService4TestCaseOperation1 for usage and error information.
func (c *InputService4ProtocolTest) InputService4TestCaseOperation1(input *InputService4TestShapeInputService4TestCaseOperation1Input) (*InputService4TestShapeInputService4TestCaseOperation1Output, error) {
	req, out := c.InputService4TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService4TestCaseOperation1WithContext is the same as InputService4TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService4TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService4ProtocolTest) InputService4TestCaseOperation1WithContext(ctx aws.Context, input *InputService4TestShapeInputService4TestCaseOperation1Input, opts ...aws.Option) (*InputService4TestShapeInputService4TestCaseOperation1Output, error) {
	req, out := c.InputService4TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService4TestShapeInputService4TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListArg []*string `type:"list"`
}

// SetListArg sets the ListArg field's value.
func (s *InputService4TestShapeInputService4TestCaseOperation1Input) SetListArg(v []*string) *InputService4TestShapeInputService4TestCaseOperation1Input {
	s.ListArg = v
	return s
}

type InputService4TestShapeInputService4TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService5ProtocolTest client from just a config.
//     svc := inputservice5protocoltest.New(myConfig)
//
//     // Create a InputService5ProtocolTest client with additional configuration
//     svc := inputservice5protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService5TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService5TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService5TestCaseOperation1 for more information on using the InputService5TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService5TestCaseOperation1Request method.
//    req, resp := client.InputService5TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService5ProtocolTest) InputService5TestCaseOperation1Request(input *InputService5TestShapeInputService5TestCaseOperation1Input) (req *aws.Request, output *InputService5TestShapeInputService5TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService5TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService5TestShapeInputService5TestCaseOperation1Input{}
	}

	output = &InputService5TestShapeInputService5TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService5TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService5TestCaseOperation1 for usage and error information.
func (c *InputService5ProtocolTest) InputService5TestCaseOperation1(input *InputService5TestShapeInputService5TestCaseOperation1Input) (*InputService5TestShapeInputService5TestCaseOperation1Output, error) {
	req, out := c.InputService5TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService5TestCaseOperation1WithContext is the same as InputService5TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService5TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService5ProtocolTest) InputService5TestCaseOperation1WithContext(ctx aws.Context, input *InputService5TestShapeInputService5TestCaseOperation1Input, opts ...aws.Option) (*InputService5TestShapeInputService5TestCaseOperation1Output, error) {
	req, out := c.InputService5TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService5TestShapeInputService5TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListArg []*string `locationName:"ListMemberName" locationNameList:"item" type:"list"`
}

// SetListArg sets the ListArg field's value.
func (s *InputService5TestShapeInputService5TestCaseOperation1Input) SetListArg(v []*string) *InputService5TestShapeInputService5TestCaseOperation1Input {
	s.ListArg = v
	return s
}

type InputService5TestShapeInputService5TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService6ProtocolTest client from just a config.
//     svc := inputservice6protocoltest.New(myConfig)
//
//     // Create a InputService6ProtocolTest client with additional configuration
//     svc := inputservice6protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService6TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService6TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService6TestCaseOperation1 for more information on using the InputService6TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService6TestCaseOperation1Request method.
//    req, resp := client.InputService6TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService6ProtocolTest) InputService6TestCaseOperation1Request(input *InputService6TestShapeInputService6TestCaseOperation1Input) (req *aws.Request, output *InputService6TestShapeInputService6TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService6TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService6TestShapeInputService6TestCaseOperation1Input{}
	}

	output = &InputService6TestShapeInputService6TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService6TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService6TestCaseOperation1 for usage and error information.
func (c *InputService6ProtocolTest) InputService6TestCaseOperation1(input *InputService6TestShapeInputService6TestCaseOperation1Input) (*InputService6TestShapeInputService6TestCaseOperation1Output, error) {
	req, out := c.InputService6TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService6TestCaseOperation1WithContext is the same as InputService6TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService6TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService6ProtocolTest) InputService6TestCaseOperation1WithContext(ctx aws.Context, input *InputService6TestShapeInputService6TestCaseOperation1Input, opts ...aws.Option) (*InputService6TestShapeInputService6TestCaseOperation1Output, error) {
	req, out := c.InputService6TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService6TestShapeInputService6TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	ListArg []*string `locationName:"ListMemberName" queryName:"ListQueryName" locationNameList:"item" type:"list"`
}

// SetListArg sets the ListArg field's value.
func (s *InputService6TestShapeInputService6TestCaseOperation1Input) SetListArg(v []*string) *InputService6TestShapeInputService6TestCaseOperation1Input {
	s.ListArg = v
	return s
}

type InputService6TestShapeInputService6TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService7ProtocolTest client from just a config.
//     svc := inputservice7protocoltest.New(myConfig)
//
//     // Create a InputService7ProtocolTest client with additional configuration
//     svc := inputservice7protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService7TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService7TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService7TestCaseOperation1 for more information on using the InputService7TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService7TestCaseOperation1Request method.
//    req, resp := client.InputService7TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService7ProtocolTest) InputService7TestCaseOperation1Request(input *InputService7TestShapeInputService7TestCaseOperation1Input) (req *aws.Request, output *InputService7TestShapeInputService7TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService7TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService7TestShapeInputService7TestCaseOperation1Input{}
	}

	output = &InputService7TestShapeInputService7TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService7TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService7TestCaseOperation1 for usage and error information.
func (c *InputService7ProtocolTest) InputService7TestCaseOperation1(input *InputService7TestShapeInputService7TestCaseOperation1Input) (*InputService7TestShapeInputService7TestCaseOperation1Output, error) {
	req, out := c.InputService7TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService7TestCaseOperation1WithContext is the same as InputService7TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService7TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService7ProtocolTest) InputService7TestCaseOperation1WithContext(ctx aws.Context, input *InputService7TestShapeInputService7TestCaseOperation1Input, opts ...aws.Option) (*InputService7TestShapeInputService7TestCaseOperation1Output, error) {
	req, out := c.InputService7TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService7TestShapeInputService7TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	// BlobArg is automatically base64 encoded/decoded by the SDK.
	BlobArg []byte `type:"blob"`
}

// SetBlobArg sets the BlobArg field's value.
func (s *InputService7TestShapeInputService7TestCaseOperation1Input) SetBlobArg(v []byte) *InputService7TestShapeInputService7TestCaseOperation1Input {
	s.BlobArg = v
	return s
}

type InputService7TestShapeInputService7TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService8ProtocolTest client from just a config.
//     svc := inputservice8protocoltest.New(myConfig)
//
//     // Create a InputService8ProtocolTest client with additional configuration
//     svc := inputservice8protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService8TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService8TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService8TestCaseOperation1 for more information on using the InputService8TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService8TestCaseOperation1Request method.
//    req, resp := client.InputService8TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService8ProtocolTest) InputService8TestCaseOperation1Request(input *InputService8TestShapeInputService8TestCaseOperation1Input) (req *aws.Request, output *InputService8TestShapeInputService8TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService8TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService8TestShapeInputService8TestCaseOperation1Input{}
	}

	output = &InputService8TestShapeInputService8TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService8TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService8TestCaseOperation1 for usage and error information.
func (c *InputService8ProtocolTest) InputService8TestCaseOperation1(input *InputService8TestShapeInputService8TestCaseOperation1Input) (*InputService8TestShapeInputService8TestCaseOperation1Output, error) {
	req, out := c.InputService8TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService8TestCaseOperation1WithContext is the same as InputService8TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService8TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService8ProtocolTest) InputService8TestCaseOperation1WithContext(ctx aws.Context, input *InputService8TestShapeInputService8TestCaseOperation1Input, opts ...aws.Option) (*InputService8TestShapeInputService8TestCaseOperation1Output, error) {
	req, out := c.InputService8TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService8TestShapeInputService8TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`

	TimeArg *time.Time `type:"timestamp" timestampFormat:"iso8601"`
}

// SetTimeArg sets the TimeArg field's value.
func (s *InputService8TestShapeInputService8TestCaseOperation1Input) SetTimeArg(v time.Time) *InputService8TestShapeInputService8TestCaseOperation1Input {
	s.TimeArg = &v
	return s
}

type InputService8TestShapeInputService8TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService9ProtocolTest client from just a config.
//     svc := inputservice9protocoltest.New(myConfig)
//
//     // Create a InputService9ProtocolTest client with additional configuration
//     svc := inputservice9protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService9TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService9TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService9TestCaseOperation1 for more information on using the InputService9TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService9TestCaseOperation1Request method.
//    req, resp := client.InputService9TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService9ProtocolTest) InputService9TestCaseOperation1Request(input *InputService9TestShapeInputService9TestCaseOperation2Input) (req *aws.Request, output *InputService9TestShapeInputService9TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService9TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation2Input{}
	}

	output = &InputService9TestShapeInputService9TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService9TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService9TestCaseOperation1 for usage and error information.
func (c *InputService9ProtocolTest) InputService9TestCaseOperation1(input *InputService9TestShapeInputService9TestCaseOperation2Input) (*InputService9TestShapeInputService9TestCaseOperation1Output, error) {
	req, out := c.InputService9TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService9TestCaseOperation1WithContext is the same as InputService9TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService9TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService9ProtocolTest) InputService9TestCaseOperation1WithContext(ctx aws.Context, input *InputService9TestShapeInputService9TestCaseOperation2Input, opts ...aws.Option) (*InputService9TestShapeInputService9TestCaseOperation1Output, error) {
	req, out := c.InputService9TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

const opInputService9TestCaseOperation2 = "OperationName"

// InputService9TestCaseOperation2Request generates a "aws.Request" representing the
// client's request for the InputService9TestCaseOperation2 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService9TestCaseOperation2 for more information on using the InputService9TestCaseOperation2
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService9TestCaseOperation2Request method.
//    req, resp := client.InputService9TestCaseOperation2Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService9ProtocolTest) InputService9TestCaseOperation2Request(input *InputService9TestShapeInputService9TestCaseOperation2Input) (req *aws.Request, output *InputService9TestShapeInputService9TestCaseOperation2Output) {
	op := &aws.Operation{
		Name:     opInputService9TestCaseOperation2,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService9TestShapeInputService9TestCaseOperation2Input{}
	}

	output = &InputService9TestShapeInputService9TestCaseOperation2Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService9TestCaseOperation2 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService9TestCaseOperation2 for usage and error information.
func (c *InputService9ProtocolTest) InputService9TestCaseOperation2(input *InputService9TestShapeInputService9TestCaseOperation2Input) (*InputService9TestShapeInputService9TestCaseOperation2Output, error) {
	req, out := c.InputService9TestCaseOperation2Request(input)
	return out, req.Send()
}

// InputService9TestCaseOperation2WithContext is the same as InputService9TestCaseOperation2 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService9TestCaseOperation2 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService9ProtocolTest) InputService9TestCaseOperation2WithContext(ctx aws.Context, input *InputService9TestShapeInputService9TestCaseOperation2Input, opts ...aws.Option) (*InputService9TestShapeInputService9TestCaseOperation2Output, error) {
	req, out := c.InputService9TestCaseOperation2Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService9TestShapeInputService9TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

type InputService9TestShapeInputService9TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	Token *string `type:"string" idempotencyToken:"true"`
}

// SetToken sets the Token field's value.
func (s *InputService9TestShapeInputService9TestCaseOperation2Input) SetToken(v string) *InputService9TestShapeInputService9TestCaseOperation2Input {
	s.Token = &v
	return s
}

type InputService9TestShapeInputService9TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
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
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a InputService10ProtocolTest client from just a config.
//     svc := inputservice10protocoltest.New(myConfig)
//
//     // Create a InputService10ProtocolTest client with additional configuration
//     svc := inputservice10protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
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

// InputService10TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the InputService10TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService10TestCaseOperation1 for more information on using the InputService10TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService10TestCaseOperation1Request method.
//    req, resp := client.InputService10TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService10ProtocolTest) InputService10TestCaseOperation1Request(input *InputService10TestShapeInputService10TestCaseOperation2Input) (req *aws.Request, output *InputService10TestShapeInputService10TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opInputService10TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService10TestShapeInputService10TestCaseOperation2Input{}
	}

	output = &InputService10TestShapeInputService10TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService10TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService10TestCaseOperation1 for usage and error information.
func (c *InputService10ProtocolTest) InputService10TestCaseOperation1(input *InputService10TestShapeInputService10TestCaseOperation2Input) (*InputService10TestShapeInputService10TestCaseOperation1Output, error) {
	req, out := c.InputService10TestCaseOperation1Request(input)
	return out, req.Send()
}

// InputService10TestCaseOperation1WithContext is the same as InputService10TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService10TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService10ProtocolTest) InputService10TestCaseOperation1WithContext(ctx aws.Context, input *InputService10TestShapeInputService10TestCaseOperation2Input, opts ...aws.Option) (*InputService10TestShapeInputService10TestCaseOperation1Output, error) {
	req, out := c.InputService10TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

const opInputService10TestCaseOperation2 = "OperationName"

// InputService10TestCaseOperation2Request generates a "aws.Request" representing the
// client's request for the InputService10TestCaseOperation2 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See InputService10TestCaseOperation2 for more information on using the InputService10TestCaseOperation2
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the InputService10TestCaseOperation2Request method.
//    req, resp := client.InputService10TestCaseOperation2Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *InputService10ProtocolTest) InputService10TestCaseOperation2Request(input *InputService10TestShapeInputService10TestCaseOperation2Input) (req *aws.Request, output *InputService10TestShapeInputService10TestCaseOperation2Output) {
	op := &aws.Operation{
		Name:     opInputService10TestCaseOperation2,
		HTTPPath: "/",
	}

	if input == nil {
		input = &InputService10TestShapeInputService10TestCaseOperation2Input{}
	}

	output = &InputService10TestShapeInputService10TestCaseOperation2Output{}
	req = c.newRequest(op, input, output)
	req.Handlers.Unmarshal.Remove(ec2query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return
}

// InputService10TestCaseOperation2 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation InputService10TestCaseOperation2 for usage and error information.
func (c *InputService10ProtocolTest) InputService10TestCaseOperation2(input *InputService10TestShapeInputService10TestCaseOperation2Input) (*InputService10TestShapeInputService10TestCaseOperation2Output, error) {
	req, out := c.InputService10TestCaseOperation2Request(input)
	return out, req.Send()
}

// InputService10TestCaseOperation2WithContext is the same as InputService10TestCaseOperation2 with the addition of
// the ability to pass a context and additional request options.
//
// See InputService10TestCaseOperation2 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *InputService10ProtocolTest) InputService10TestCaseOperation2WithContext(ctx aws.Context, input *InputService10TestShapeInputService10TestCaseOperation2Input, opts ...aws.Option) (*InputService10TestShapeInputService10TestCaseOperation2Output, error) {
	req, out := c.InputService10TestCaseOperation2Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type InputService10TestShapeInputService10TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`
}

type InputService10TestShapeInputService10TestCaseOperation2Input struct {
	_ struct{} `type:"structure"`

	FooEnum InputService10TestShapeEnumType `type:"string"`

	ListEnums []InputService10TestShapeEnumType `type:"list"`
}

// SetFooEnum sets the FooEnum field's value.
func (s *InputService10TestShapeInputService10TestCaseOperation2Input) SetFooEnum(v InputService10TestShapeEnumType) *InputService10TestShapeInputService10TestCaseOperation2Input {
	s.FooEnum = v
	return s
}

// SetListEnums sets the ListEnums field's value.
func (s *InputService10TestShapeInputService10TestCaseOperation2Input) SetListEnums(v []InputService10TestShapeEnumType) *InputService10TestShapeInputService10TestCaseOperation2Input {
	s.ListEnums = v
	return s
}

type InputService10TestShapeInputService10TestCaseOperation2Output struct {
	_ struct{} `type:"structure"`
}

type InputService10TestShapeEnumType string

// Enum values for InputService10TestShapeEnumType
const (
	EnumTypeFoo InputService10TestShapeEnumType = "foo"
	EnumTypeBar InputService10TestShapeEnumType = "bar"
)

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
	req, _ := svc.InputService1TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
	req, _ := svc.InputService2TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
	req, _ := svc.InputService3TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
		ListArg: []*string{
			aws.String("foo"),
			aws.String("bar"),
			aws.String("baz"),
		},
	}
	req, _ := svc.InputService4TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}

	// assert body
	if r.Body == nil {
		t.Errorf("expect body not to be nil")
	}
	body, _ := ioutil.ReadAll(r.Body)
	awstesting.AssertQuery(t, `Action=OperationName&ListArg.1=foo&ListArg.2=bar&ListArg.3=baz&Version=2014-01-01`, util.Trim(string(body)))

	// assert URL
	awstesting.AssertURL(t, "https://test/", r.URL.String())

	// assert headers

}

func TestInputService5ProtocolTestListWithLocationNameAppliedToMemberCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewInputService5ProtocolTest(cfg)
	input := &InputService5TestShapeInputService5TestCaseOperation1Input{
		ListArg: []*string{
			aws.String("a"),
			aws.String("b"),
			aws.String("c"),
		},
	}
	req, _ := svc.InputService5TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
		ListArg: []*string{
			aws.String("a"),
			aws.String("b"),
			aws.String("c"),
		},
	}
	req, _ := svc.InputService6TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
	req, _ := svc.InputService7TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
	req, _ := svc.InputService8TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
	req, _ := svc.InputService9TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
	req, _ := svc.InputService9TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
	req, _ := svc.InputService10TestCaseOperation1Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
	req, _ := svc.InputService10TestCaseOperation2Request(input)
	r := req.HTTPRequest

	// build request
	ec2query.Build(req)
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
