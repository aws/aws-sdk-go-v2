package apigateway

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/rest"
	"github.com/aws/aws-sdk-go-v2/private/protocol/jsonrpc"
	restlegacy "github.com/aws/aws-sdk-go-v2/private/protocol/rest"
)

// protoGetAPIKeyMarshaler defines a marshaler for GetApiKey operation
type protoGetAPIKeyMarshaler struct {
	input *GetApiKeyInput
}

// marshalOperation is the top level method used within a handler stack to marshal an operation
// This method calls appropriate marshal shape functions as per the input shape and protocol used by the service.
func (m protoGetAPIKeyMarshaler) marshalOperation(r *aws.Request) {
	var err error
	encoder := rest.NewEncoder(r.HTTPRequest)

	// We add Content-Type Header if input shape's is not of type streaming payload
	// The value of Content-Type is decided by following:
	// a. if shape's metadata has JSONVersion and protocol is json.
	//    - application/x-amz-json-%s where %s is ths JSONVersion.
	// b. else if protocol is either json or rest-json
	//     - application/json
	// Here protocol is rest-json, and shape is not of type streaming payload,
	// thus content-type header with value application/json is added.
	encoder.AddHeader("Content-Type").String("application/json")
	err = marshalGetAPIKeyInputShapeAWSREST(m.input, encoder)
	if err != nil {
		r.Error = err
		return
	}

	err = encoder.Encode()
	if err != nil {
		r.Error = err
		return
	}

	err = marshalGetAPIKeyInputShapeAWSJSON(m.input, r)
	if err != nil {
		r.Error = err
	}
}

// marshalGetAPIKeyInputShapeAWSREST is a stand alone function used to marshal the HTTP bindings a input shape.
// This method uses the rest encoder utility
func marshalGetAPIKeyInputShapeAWSREST(input *GetApiKeyInput, encoder *rest.Encoder) error {
	if input.ApiKey != nil {
		if err := encoder.SetURI("api_Key").String(*input.ApiKey); err != nil {
			return awserr.New(aws.ErrCodeSerialization, "failed to marshal API KEY to URI using REST encoder", err)
		}
	}
	if input.IncludeValue != nil {
		encoder.AddQuery("includeValue").Boolean(*input.IncludeValue)
	}
	return nil
}

// marshalGetAPIKeyInputShapeAWSJSON is a stand alone function used to marshal the json body
func marshalGetAPIKeyInputShapeAWSJSON(v *GetApiKeyInput, r *aws.Request) error {
	// delegate to reflection based marshaling
	if t := restlegacy.PayloadType(r.Params); t == "structure" || t == "" {
		jsonrpc.Build(r)
	}
	return nil
}

// NamedHandler returns a Named Build Handler for an operation marshal function
func (m protoGetAPIKeyMarshaler) NamedHandler() aws.NamedHandler {
	const buildHandler = "ProtoGetApiKey.BuildHandler"
	return aws.NamedHandler{
		Name: buildHandler,
		Fn:   m.marshalOperation,
	}
}
