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

// OutputService1ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// OutputService1ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService1ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the OutputService1ProtocolTest client with a config.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a OutputService1ProtocolTest client from just a config.
//     svc := outputservice1protocoltest.New(myConfig)
//
//     // Create a OutputService1ProtocolTest client with additional configuration
//     svc := outputservice1protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
func NewOutputService1ProtocolTest(config aws.Config) *OutputService1ProtocolTest {
	var signingName string
	signingRegion := aws.StringValue(config.Region)

	svc := &OutputService1ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "outputservice1protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
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

// newRequest creates a new request for a OutputService1ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService1ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opOutputService1TestCaseOperation1 = "OperationName"

// OutputService1TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the OutputService1TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See OutputService1TestCaseOperation1 for more information on using the OutputService1TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the OutputService1TestCaseOperation1Request method.
//    req, resp := client.OutputService1TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation1Request(input *OutputService1TestShapeOutputService1TestCaseOperation1Input) (req *aws.Request, output *OutputService1TestShapeOutputService1TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opOutputService1TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService1TestShapeOutputService1TestCaseOperation1Input{}
	}

	output = &OutputService1TestShapeOutputService1TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	return
}

// OutputService1TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation OutputService1TestCaseOperation1 for usage and error information.
func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation1(input *OutputService1TestShapeOutputService1TestCaseOperation1Input) (*OutputService1TestShapeOutputService1TestCaseOperation1Output, error) {
	req, out := c.OutputService1TestCaseOperation1Request(input)
	return out, req.Send()
}

// OutputService1TestCaseOperation1WithContext is the same as OutputService1TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See OutputService1TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *OutputService1ProtocolTest) OutputService1TestCaseOperation1WithContext(ctx aws.Context, input *OutputService1TestShapeOutputService1TestCaseOperation1Input, opts ...aws.Option) (*OutputService1TestShapeOutputService1TestCaseOperation1Output, error) {
	req, out := c.OutputService1TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
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

	TrueBool *bool `type:"boolean"`
}

// SetChar sets the Char field's value.
func (s *OutputService1TestShapeOutputService1TestCaseOperation1Output) SetChar(v string) *OutputService1TestShapeOutputService1TestCaseOperation1Output {
	s.Char = &v
	return s
}

// SetDouble sets the Double field's value.
func (s *OutputService1TestShapeOutputService1TestCaseOperation1Output) SetDouble(v float64) *OutputService1TestShapeOutputService1TestCaseOperation1Output {
	s.Double = &v
	return s
}

// SetFalseBool sets the FalseBool field's value.
func (s *OutputService1TestShapeOutputService1TestCaseOperation1Output) SetFalseBool(v bool) *OutputService1TestShapeOutputService1TestCaseOperation1Output {
	s.FalseBool = &v
	return s
}

// SetFloat sets the Float field's value.
func (s *OutputService1TestShapeOutputService1TestCaseOperation1Output) SetFloat(v float64) *OutputService1TestShapeOutputService1TestCaseOperation1Output {
	s.Float = &v
	return s
}

// SetLong sets the Long field's value.
func (s *OutputService1TestShapeOutputService1TestCaseOperation1Output) SetLong(v int64) *OutputService1TestShapeOutputService1TestCaseOperation1Output {
	s.Long = &v
	return s
}

// SetNum sets the Num field's value.
func (s *OutputService1TestShapeOutputService1TestCaseOperation1Output) SetNum(v int64) *OutputService1TestShapeOutputService1TestCaseOperation1Output {
	s.Num = &v
	return s
}

// SetStr sets the Str field's value.
func (s *OutputService1TestShapeOutputService1TestCaseOperation1Output) SetStr(v string) *OutputService1TestShapeOutputService1TestCaseOperation1Output {
	s.Str = &v
	return s
}

// SetTrueBool sets the TrueBool field's value.
func (s *OutputService1TestShapeOutputService1TestCaseOperation1Output) SetTrueBool(v bool) *OutputService1TestShapeOutputService1TestCaseOperation1Output {
	s.TrueBool = &v
	return s
}

// OutputService2ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// OutputService2ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService2ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the OutputService2ProtocolTest client with a config.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a OutputService2ProtocolTest client from just a config.
//     svc := outputservice2protocoltest.New(myConfig)
//
//     // Create a OutputService2ProtocolTest client with additional configuration
//     svc := outputservice2protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
func NewOutputService2ProtocolTest(config aws.Config) *OutputService2ProtocolTest {
	var signingName string
	signingRegion := aws.StringValue(config.Region)

	svc := &OutputService2ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "outputservice2protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
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

// newRequest creates a new request for a OutputService2ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService2ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opOutputService2TestCaseOperation1 = "OperationName"

// OutputService2TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the OutputService2TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See OutputService2TestCaseOperation1 for more information on using the OutputService2TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the OutputService2TestCaseOperation1Request method.
//    req, resp := client.OutputService2TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *OutputService2ProtocolTest) OutputService2TestCaseOperation1Request(input *OutputService2TestShapeOutputService2TestCaseOperation1Input) (req *aws.Request, output *OutputService2TestShapeOutputService2TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opOutputService2TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService2TestShapeOutputService2TestCaseOperation1Input{}
	}

	output = &OutputService2TestShapeOutputService2TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	return
}

// OutputService2TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation OutputService2TestCaseOperation1 for usage and error information.
func (c *OutputService2ProtocolTest) OutputService2TestCaseOperation1(input *OutputService2TestShapeOutputService2TestCaseOperation1Input) (*OutputService2TestShapeOutputService2TestCaseOperation1Output, error) {
	req, out := c.OutputService2TestCaseOperation1Request(input)
	return out, req.Send()
}

// OutputService2TestCaseOperation1WithContext is the same as OutputService2TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See OutputService2TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *OutputService2ProtocolTest) OutputService2TestCaseOperation1WithContext(ctx aws.Context, input *OutputService2TestShapeOutputService2TestCaseOperation1Input, opts ...aws.Option) (*OutputService2TestShapeOutputService2TestCaseOperation1Output, error) {
	req, out := c.OutputService2TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type OutputService2TestShapeOutputService2TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService2TestShapeOutputService2TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	// Blob is automatically base64 encoded/decoded by the SDK.
	Blob []byte `type:"blob"`
}

// SetBlob sets the Blob field's value.
func (s *OutputService2TestShapeOutputService2TestCaseOperation1Output) SetBlob(v []byte) *OutputService2TestShapeOutputService2TestCaseOperation1Output {
	s.Blob = v
	return s
}

// OutputService3ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// OutputService3ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService3ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the OutputService3ProtocolTest client with a config.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a OutputService3ProtocolTest client from just a config.
//     svc := outputservice3protocoltest.New(myConfig)
//
//     // Create a OutputService3ProtocolTest client with additional configuration
//     svc := outputservice3protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
func NewOutputService3ProtocolTest(config aws.Config) *OutputService3ProtocolTest {
	var signingName string
	signingRegion := aws.StringValue(config.Region)

	svc := &OutputService3ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "outputservice3protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
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

// newRequest creates a new request for a OutputService3ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService3ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opOutputService3TestCaseOperation1 = "OperationName"

// OutputService3TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the OutputService3TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See OutputService3TestCaseOperation1 for more information on using the OutputService3TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the OutputService3TestCaseOperation1Request method.
//    req, resp := client.OutputService3TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *OutputService3ProtocolTest) OutputService3TestCaseOperation1Request(input *OutputService3TestShapeOutputService3TestCaseOperation1Input) (req *aws.Request, output *OutputService3TestShapeOutputService3TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opOutputService3TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService3TestShapeOutputService3TestCaseOperation1Input{}
	}

	output = &OutputService3TestShapeOutputService3TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	return
}

// OutputService3TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation OutputService3TestCaseOperation1 for usage and error information.
func (c *OutputService3ProtocolTest) OutputService3TestCaseOperation1(input *OutputService3TestShapeOutputService3TestCaseOperation1Input) (*OutputService3TestShapeOutputService3TestCaseOperation1Output, error) {
	req, out := c.OutputService3TestCaseOperation1Request(input)
	return out, req.Send()
}

// OutputService3TestCaseOperation1WithContext is the same as OutputService3TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See OutputService3TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *OutputService3ProtocolTest) OutputService3TestCaseOperation1WithContext(ctx aws.Context, input *OutputService3TestShapeOutputService3TestCaseOperation1Input, opts ...aws.Option) (*OutputService3TestShapeOutputService3TestCaseOperation1Output, error) {
	req, out := c.OutputService3TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type OutputService3TestShapeOutputService3TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService3TestShapeOutputService3TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	ListMember []*string `type:"list"`
}

// SetListMember sets the ListMember field's value.
func (s *OutputService3TestShapeOutputService3TestCaseOperation1Output) SetListMember(v []*string) *OutputService3TestShapeOutputService3TestCaseOperation1Output {
	s.ListMember = v
	return s
}

// OutputService4ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// OutputService4ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService4ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the OutputService4ProtocolTest client with a config.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a OutputService4ProtocolTest client from just a config.
//     svc := outputservice4protocoltest.New(myConfig)
//
//     // Create a OutputService4ProtocolTest client with additional configuration
//     svc := outputservice4protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
func NewOutputService4ProtocolTest(config aws.Config) *OutputService4ProtocolTest {
	var signingName string
	signingRegion := aws.StringValue(config.Region)

	svc := &OutputService4ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "outputservice4protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
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

// newRequest creates a new request for a OutputService4ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService4ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opOutputService4TestCaseOperation1 = "OperationName"

// OutputService4TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the OutputService4TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See OutputService4TestCaseOperation1 for more information on using the OutputService4TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the OutputService4TestCaseOperation1Request method.
//    req, resp := client.OutputService4TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *OutputService4ProtocolTest) OutputService4TestCaseOperation1Request(input *OutputService4TestShapeOutputService4TestCaseOperation1Input) (req *aws.Request, output *OutputService4TestShapeOutputService4TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opOutputService4TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService4TestShapeOutputService4TestCaseOperation1Input{}
	}

	output = &OutputService4TestShapeOutputService4TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	return
}

// OutputService4TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation OutputService4TestCaseOperation1 for usage and error information.
func (c *OutputService4ProtocolTest) OutputService4TestCaseOperation1(input *OutputService4TestShapeOutputService4TestCaseOperation1Input) (*OutputService4TestShapeOutputService4TestCaseOperation1Output, error) {
	req, out := c.OutputService4TestCaseOperation1Request(input)
	return out, req.Send()
}

// OutputService4TestCaseOperation1WithContext is the same as OutputService4TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See OutputService4TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *OutputService4ProtocolTest) OutputService4TestCaseOperation1WithContext(ctx aws.Context, input *OutputService4TestShapeOutputService4TestCaseOperation1Input, opts ...aws.Option) (*OutputService4TestShapeOutputService4TestCaseOperation1Output, error) {
	req, out := c.OutputService4TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type OutputService4TestShapeOutputService4TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService4TestShapeOutputService4TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	ListMember []*string `locationNameList:"item" type:"list"`
}

// SetListMember sets the ListMember field's value.
func (s *OutputService4TestShapeOutputService4TestCaseOperation1Output) SetListMember(v []*string) *OutputService4TestShapeOutputService4TestCaseOperation1Output {
	s.ListMember = v
	return s
}

// OutputService5ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// OutputService5ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService5ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the OutputService5ProtocolTest client with a config.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a OutputService5ProtocolTest client from just a config.
//     svc := outputservice5protocoltest.New(myConfig)
//
//     // Create a OutputService5ProtocolTest client with additional configuration
//     svc := outputservice5protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
func NewOutputService5ProtocolTest(config aws.Config) *OutputService5ProtocolTest {
	var signingName string
	signingRegion := aws.StringValue(config.Region)

	svc := &OutputService5ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "outputservice5protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
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

// newRequest creates a new request for a OutputService5ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService5ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opOutputService5TestCaseOperation1 = "OperationName"

// OutputService5TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the OutputService5TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See OutputService5TestCaseOperation1 for more information on using the OutputService5TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the OutputService5TestCaseOperation1Request method.
//    req, resp := client.OutputService5TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *OutputService5ProtocolTest) OutputService5TestCaseOperation1Request(input *OutputService5TestShapeOutputService5TestCaseOperation1Input) (req *aws.Request, output *OutputService5TestShapeOutputService5TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opOutputService5TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService5TestShapeOutputService5TestCaseOperation1Input{}
	}

	output = &OutputService5TestShapeOutputService5TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	return
}

// OutputService5TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation OutputService5TestCaseOperation1 for usage and error information.
func (c *OutputService5ProtocolTest) OutputService5TestCaseOperation1(input *OutputService5TestShapeOutputService5TestCaseOperation1Input) (*OutputService5TestShapeOutputService5TestCaseOperation1Output, error) {
	req, out := c.OutputService5TestCaseOperation1Request(input)
	return out, req.Send()
}

// OutputService5TestCaseOperation1WithContext is the same as OutputService5TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See OutputService5TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *OutputService5ProtocolTest) OutputService5TestCaseOperation1WithContext(ctx aws.Context, input *OutputService5TestShapeOutputService5TestCaseOperation1Input, opts ...aws.Option) (*OutputService5TestShapeOutputService5TestCaseOperation1Output, error) {
	req, out := c.OutputService5TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type OutputService5TestShapeOutputService5TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService5TestShapeOutputService5TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	ListMember []*string `type:"list" flattened:"true"`
}

// SetListMember sets the ListMember field's value.
func (s *OutputService5TestShapeOutputService5TestCaseOperation1Output) SetListMember(v []*string) *OutputService5TestShapeOutputService5TestCaseOperation1Output {
	s.ListMember = v
	return s
}

// OutputService6ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// OutputService6ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService6ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the OutputService6ProtocolTest client with a config.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a OutputService6ProtocolTest client from just a config.
//     svc := outputservice6protocoltest.New(myConfig)
//
//     // Create a OutputService6ProtocolTest client with additional configuration
//     svc := outputservice6protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
func NewOutputService6ProtocolTest(config aws.Config) *OutputService6ProtocolTest {
	var signingName string
	signingRegion := aws.StringValue(config.Region)

	svc := &OutputService6ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "outputservice6protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
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

// newRequest creates a new request for a OutputService6ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService6ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opOutputService6TestCaseOperation1 = "OperationName"

// OutputService6TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the OutputService6TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See OutputService6TestCaseOperation1 for more information on using the OutputService6TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the OutputService6TestCaseOperation1Request method.
//    req, resp := client.OutputService6TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *OutputService6ProtocolTest) OutputService6TestCaseOperation1Request(input *OutputService6TestShapeOutputService6TestCaseOperation1Input) (req *aws.Request, output *OutputService6TestShapeOutputService6TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opOutputService6TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService6TestShapeOutputService6TestCaseOperation1Input{}
	}

	output = &OutputService6TestShapeOutputService6TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	return
}

// OutputService6TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation OutputService6TestCaseOperation1 for usage and error information.
func (c *OutputService6ProtocolTest) OutputService6TestCaseOperation1(input *OutputService6TestShapeOutputService6TestCaseOperation1Input) (*OutputService6TestShapeOutputService6TestCaseOperation1Output, error) {
	req, out := c.OutputService6TestCaseOperation1Request(input)
	return out, req.Send()
}

// OutputService6TestCaseOperation1WithContext is the same as OutputService6TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See OutputService6TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *OutputService6ProtocolTest) OutputService6TestCaseOperation1WithContext(ctx aws.Context, input *OutputService6TestShapeOutputService6TestCaseOperation1Input, opts ...aws.Option) (*OutputService6TestShapeOutputService6TestCaseOperation1Output, error) {
	req, out := c.OutputService6TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type OutputService6TestShapeOutputService6TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService6TestShapeOutputService6TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Map map[string]*OutputService6TestShapeStructureType `type:"map"`
}

// SetMap sets the Map field's value.
func (s *OutputService6TestShapeOutputService6TestCaseOperation1Output) SetMap(v map[string]*OutputService6TestShapeStructureType) *OutputService6TestShapeOutputService6TestCaseOperation1Output {
	s.Map = v
	return s
}

type OutputService6TestShapeStructureType struct {
	_ struct{} `type:"structure"`

	Foo *string `locationName:"foo" type:"string"`
}

// SetFoo sets the Foo field's value.
func (s *OutputService6TestShapeStructureType) SetFoo(v string) *OutputService6TestShapeStructureType {
	s.Foo = &v
	return s
}

// OutputService7ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// OutputService7ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService7ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the OutputService7ProtocolTest client with a config.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a OutputService7ProtocolTest client from just a config.
//     svc := outputservice7protocoltest.New(myConfig)
//
//     // Create a OutputService7ProtocolTest client with additional configuration
//     svc := outputservice7protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
func NewOutputService7ProtocolTest(config aws.Config) *OutputService7ProtocolTest {
	var signingName string
	signingRegion := aws.StringValue(config.Region)

	svc := &OutputService7ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "outputservice7protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
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

// newRequest creates a new request for a OutputService7ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService7ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opOutputService7TestCaseOperation1 = "OperationName"

// OutputService7TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the OutputService7TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See OutputService7TestCaseOperation1 for more information on using the OutputService7TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the OutputService7TestCaseOperation1Request method.
//    req, resp := client.OutputService7TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *OutputService7ProtocolTest) OutputService7TestCaseOperation1Request(input *OutputService7TestShapeOutputService7TestCaseOperation1Input) (req *aws.Request, output *OutputService7TestShapeOutputService7TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opOutputService7TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService7TestShapeOutputService7TestCaseOperation1Input{}
	}

	output = &OutputService7TestShapeOutputService7TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	return
}

// OutputService7TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation OutputService7TestCaseOperation1 for usage and error information.
func (c *OutputService7ProtocolTest) OutputService7TestCaseOperation1(input *OutputService7TestShapeOutputService7TestCaseOperation1Input) (*OutputService7TestShapeOutputService7TestCaseOperation1Output, error) {
	req, out := c.OutputService7TestCaseOperation1Request(input)
	return out, req.Send()
}

// OutputService7TestCaseOperation1WithContext is the same as OutputService7TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See OutputService7TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *OutputService7ProtocolTest) OutputService7TestCaseOperation1WithContext(ctx aws.Context, input *OutputService7TestShapeOutputService7TestCaseOperation1Input, opts ...aws.Option) (*OutputService7TestShapeOutputService7TestCaseOperation1Output, error) {
	req, out := c.OutputService7TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type OutputService7TestShapeOutputService7TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService7TestShapeOutputService7TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Map map[string]*string `type:"map" flattened:"true"`
}

// SetMap sets the Map field's value.
func (s *OutputService7TestShapeOutputService7TestCaseOperation1Output) SetMap(v map[string]*string) *OutputService7TestShapeOutputService7TestCaseOperation1Output {
	s.Map = v
	return s
}

// OutputService8ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// OutputService8ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService8ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the OutputService8ProtocolTest client with a config.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a OutputService8ProtocolTest client from just a config.
//     svc := outputservice8protocoltest.New(myConfig)
//
//     // Create a OutputService8ProtocolTest client with additional configuration
//     svc := outputservice8protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
func NewOutputService8ProtocolTest(config aws.Config) *OutputService8ProtocolTest {
	var signingName string
	signingRegion := aws.StringValue(config.Region)

	svc := &OutputService8ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "outputservice8protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
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

// newRequest creates a new request for a OutputService8ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService8ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opOutputService8TestCaseOperation1 = "OperationName"

// OutputService8TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the OutputService8TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See OutputService8TestCaseOperation1 for more information on using the OutputService8TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the OutputService8TestCaseOperation1Request method.
//    req, resp := client.OutputService8TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *OutputService8ProtocolTest) OutputService8TestCaseOperation1Request(input *OutputService8TestShapeOutputService8TestCaseOperation1Input) (req *aws.Request, output *OutputService8TestShapeOutputService8TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opOutputService8TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService8TestShapeOutputService8TestCaseOperation1Input{}
	}

	output = &OutputService8TestShapeOutputService8TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	return
}

// OutputService8TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation OutputService8TestCaseOperation1 for usage and error information.
func (c *OutputService8ProtocolTest) OutputService8TestCaseOperation1(input *OutputService8TestShapeOutputService8TestCaseOperation1Input) (*OutputService8TestShapeOutputService8TestCaseOperation1Output, error) {
	req, out := c.OutputService8TestCaseOperation1Request(input)
	return out, req.Send()
}

// OutputService8TestCaseOperation1WithContext is the same as OutputService8TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See OutputService8TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *OutputService8ProtocolTest) OutputService8TestCaseOperation1WithContext(ctx aws.Context, input *OutputService8TestShapeOutputService8TestCaseOperation1Input, opts ...aws.Option) (*OutputService8TestShapeOutputService8TestCaseOperation1Output, error) {
	req, out := c.OutputService8TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type OutputService8TestShapeOutputService8TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService8TestShapeOutputService8TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Map map[string]*string `locationNameKey:"foo" locationNameValue:"bar" type:"map" flattened:"true"`
}

// SetMap sets the Map field's value.
func (s *OutputService8TestShapeOutputService8TestCaseOperation1Output) SetMap(v map[string]*string) *OutputService8TestShapeOutputService8TestCaseOperation1Output {
	s.Map = v
	return s
}

// OutputService9ProtocolTest provides the API operation methods for making requests to
// . See this package's package overview docs
// for details on the service.
//
// OutputService9ProtocolTest methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type OutputService9ProtocolTest struct {
	*aws.Client
}

// New creates a new instance of the OutputService9ProtocolTest client with a config.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a OutputService9ProtocolTest client from just a config.
//     svc := outputservice9protocoltest.New(myConfig)
//
//     // Create a OutputService9ProtocolTest client with additional configuration
//     svc := outputservice9protocoltest.New(myConfig, aws.NewConfig().WithRegion("us-west-2"))
func NewOutputService9ProtocolTest(config aws.Config) *OutputService9ProtocolTest {
	var signingName string
	signingRegion := aws.StringValue(config.Region)

	svc := &OutputService9ProtocolTest{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   "outputservice9protocoltest",
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "",
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

// newRequest creates a new request for a OutputService9ProtocolTest operation and runs any
// custom request initialization.
func (c *OutputService9ProtocolTest) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	return req
}

const opOutputService9TestCaseOperation1 = "OperationName"

// OutputService9TestCaseOperation1Request generates a "aws.Request" representing the
// client's request for the OutputService9TestCaseOperation1 operation. The "output" return
// value will be populated with the request's response once the request complets
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See OutputService9TestCaseOperation1 for more information on using the OutputService9TestCaseOperation1
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the OutputService9TestCaseOperation1Request method.
//    req, resp := client.OutputService9TestCaseOperation1Request(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
func (c *OutputService9ProtocolTest) OutputService9TestCaseOperation1Request(input *OutputService9TestShapeOutputService9TestCaseOperation1Input) (req *aws.Request, output *OutputService9TestShapeOutputService9TestCaseOperation1Output) {
	op := &aws.Operation{
		Name:     opOutputService9TestCaseOperation1,
		HTTPPath: "/",
	}

	if input == nil {
		input = &OutputService9TestShapeOutputService9TestCaseOperation1Input{}
	}

	output = &OutputService9TestShapeOutputService9TestCaseOperation1Output{}
	req = c.newRequest(op, input, output)
	return
}

// OutputService9TestCaseOperation1 API operation for .
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for 's
// API operation OutputService9TestCaseOperation1 for usage and error information.
func (c *OutputService9ProtocolTest) OutputService9TestCaseOperation1(input *OutputService9TestShapeOutputService9TestCaseOperation1Input) (*OutputService9TestShapeOutputService9TestCaseOperation1Output, error) {
	req, out := c.OutputService9TestCaseOperation1Request(input)
	return out, req.Send()
}

// OutputService9TestCaseOperation1WithContext is the same as OutputService9TestCaseOperation1 with the addition of
// the ability to pass a context and additional request options.
//
// See OutputService9TestCaseOperation1 for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *OutputService9ProtocolTest) OutputService9TestCaseOperation1WithContext(ctx aws.Context, input *OutputService9TestShapeOutputService9TestCaseOperation1Input, opts ...aws.Option) (*OutputService9TestShapeOutputService9TestCaseOperation1Output, error) {
	req, out := c.OutputService9TestCaseOperation1Request(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

type OutputService9TestShapeOutputService9TestCaseOperation1Input struct {
	_ struct{} `type:"structure"`
}

type OutputService9TestShapeOutputService9TestCaseOperation1Output struct {
	_ struct{} `type:"structure"`

	Foo *string `type:"string"`
}

// SetFoo sets the Foo field's value.
func (s *OutputService9TestShapeOutputService9TestCaseOperation1Output) SetFoo(v string) *OutputService9TestShapeOutputService9TestCaseOperation1Output {
	s.Foo = &v
	return s
}

//
// Tests begin here
//

func TestOutputService1ProtocolTestScalarMembersCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService1ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Str>myname</Str><FooNum>123</FooNum><FalseBool>false</FalseBool><TrueBool>true</TrueBool><Float>1.2</Float><Double>1.3</Double><Long>200</Long><Char>a</Char><RequestId>request-id</RequestId></OperationNameResponse>"))
	req, out := svc.OutputService1TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	ec2query.UnmarshalMeta(req)
	ec2query.Unmarshal(req)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

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
	if e, a := "myname", *out.Str; e != a {
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

	buf := bytes.NewReader([]byte("<OperationNameResponse><Blob>dmFsdWU=</Blob><RequestId>requestid</RequestId></OperationNameResponse>"))
	req, out := svc.OutputService2TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	ec2query.UnmarshalMeta(req)
	ec2query.Unmarshal(req)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

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

	buf := bytes.NewReader([]byte("<OperationNameResponse><ListMember><member>abc</member><member>123</member></ListMember><RequestId>requestid</RequestId></OperationNameResponse>"))
	req, out := svc.OutputService3TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	ec2query.UnmarshalMeta(req)
	ec2query.Unmarshal(req)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "abc", *out.ListMember[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "123", *out.ListMember[1]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService4ProtocolTestListWithCustomMemberNameCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService4ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><ListMember><item>abc</item><item>123</item></ListMember><RequestId>requestid</RequestId></OperationNameResponse>"))
	req, out := svc.OutputService4TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	ec2query.UnmarshalMeta(req)
	ec2query.Unmarshal(req)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "abc", *out.ListMember[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "123", *out.ListMember[1]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService5ProtocolTestFlattenedListCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService5ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><ListMember>abc</ListMember><ListMember>123</ListMember><RequestId>requestid</RequestId></OperationNameResponse>"))
	req, out := svc.OutputService5TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	ec2query.UnmarshalMeta(req)
	ec2query.Unmarshal(req)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "abc", *out.ListMember[0]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "123", *out.ListMember[1]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService6ProtocolTestNormalMapCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService6ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Map><entry><key>qux</key><value><foo>bar</foo></value></entry><entry><key>baz</key><value><foo>bam</foo></value></entry></Map><RequestId>requestid</RequestId></OperationNameResponse>"))
	req, out := svc.OutputService6TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	ec2query.UnmarshalMeta(req)
	ec2query.Unmarshal(req)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "bam", *out.Map["baz"].Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "bar", *out.Map["qux"].Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService7ProtocolTestFlattenedMapCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService7ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Map><key>qux</key><value>bar</value></Map><Map><key>baz</key><value>bam</value></Map><RequestId>requestid</RequestId></OperationNameResponse>"))
	req, out := svc.OutputService7TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	ec2query.UnmarshalMeta(req)
	ec2query.Unmarshal(req)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "bam", *out.Map["baz"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "bar", *out.Map["qux"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService8ProtocolTestNamedMapCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService8ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Map><foo>qux</foo><bar>bar</bar></Map><Map><foo>baz</foo><bar>bam</bar></Map><RequestId>requestid</RequestId></OperationNameResponse>"))
	req, out := svc.OutputService8TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	ec2query.UnmarshalMeta(req)
	ec2query.Unmarshal(req)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "bam", *out.Map["baz"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "bar", *out.Map["qux"]; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestOutputService9ProtocolTestEmptyStringCase1(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://test")

	svc := NewOutputService9ProtocolTest(cfg)

	buf := bytes.NewReader([]byte("<OperationNameResponse><Foo/><RequestId>requestid</RequestId></OperationNameResponse>"))
	req, out := svc.OutputService9TestCaseOperation1Request(nil)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}

	// set headers

	// unmarshal response
	ec2query.UnmarshalMeta(req)
	ec2query.Unmarshal(req)
	if req.Error != nil {
		t.Errorf("expect not error, got %v", req.Error)
	}

	// assert response
	if out == nil {
		t.Errorf("expect not to be nil")
	}
	if e, a := "", *out.Foo; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}
