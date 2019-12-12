package awsrestjson

import (
	"github.com/aws/aws-sdk-go-v2/aws/awserr"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/jsonrpc"
	restlegacy "github.com/aws/aws-sdk-go-v2/private/protocol/rest"
	restV2 "github.com/aws/aws-sdk-go-v2/private/protocol/rest/v2"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

// ProtoGetAPIKeyMarshaler defines marshaler for ProtoGetApiKey operation
type ProtoGetAPIKeyMarshaler struct {
	Input *types.GetApiKeyInput
}

// MarshalOperation is the top level method used within a handler stack to marshal an operation
// This method calls appropriate marshal shape functions as per the input shape and protocol used by the service.
func (m ProtoGetAPIKeyMarshaler) MarshalOperation(r *aws.Request) {
	var err error
	encoder := restV2.NewEncoder(r.HTTPRequest)
	// adds content-type header
	encoder.AddHeader("Content-Type").String("application/json")

	err = MarshalGetAPIKeyInputShapeAWSREST(m.Input, encoder)
	if err != nil {
		r.Error = err
		return
	}
	encoder.Encode()

	// Todo Instead of passing aws.Request directly to MarshalGetAPIKeyInputShapeAWSJSON;
	//  we should pass the payload as an argument
	err = MarshalGetAPIKeyInputShapeAWSJSON(m.Input, r)
	if err != nil {
		r.Error = err
		return
	}
}

// MarshalGetAPIKeyInputShapeAWSREST is a stand alone function used to marshal the HTTP bindings a input shape.
// This method uses the rest encoder utility
func MarshalGetAPIKeyInputShapeAWSREST(input *types.GetApiKeyInput, encoder *restV2.Encoder) error {
	// encode using rest encoder utility
	if input.ApiKey != nil {
		if err := encoder.SetURI("api_Key").String(*input.ApiKey); err != nil {
			return awserr.New(aws.ErrCodeSerialization, "failed to marshal API KEY to URI using REST encoder:\n \t %v", err)
		}
	}

	if input.IncludeValue != nil {
		encoder.AddQuery("includeValue").Boolean(*input.IncludeValue)
	}

	return nil
}

// MarshalGetAPIKeyInputShapeAWSJSON is a stand alone function used to marshal the json body
func MarshalGetAPIKeyInputShapeAWSJSON(v *types.GetApiKeyInput, r *aws.Request) error {
	// delegate to reflection based marshaling
	if t := restlegacy.PayloadType(r.Params); t == "structure" || t == "" {
		jsonrpc.Build(r)
	}
	return r.Error
}

// GetNamedBuildHandler returns a Named Build Handler for an operation marshal function
func (m ProtoGetAPIKeyMarshaler) GetNamedBuildHandler() aws.NamedHandler {
	const BuildHandler = "ProtoGetApiKey.BuildHandler"
	return aws.NamedHandler{
		Name: BuildHandler,
		Fn:   m.MarshalOperation,
	}
}
