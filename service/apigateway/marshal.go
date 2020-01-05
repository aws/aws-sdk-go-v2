package apigateway

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/rest"
	"github.com/aws/aws-sdk-go-v2/private/protocol/jsonrpc"
	restlegacy "github.com/aws/aws-sdk-go-v2/private/protocol/rest"
)

// protoGetApiKeyMarshaler defines a marshaler for GetApiKey operation
type protoGetApiKeyMarshaler struct {
	input *GetApiKeyInput
}

// marshalOperation is the top level method used within a handler stack to marshal an operation
// This method calls appropriate marshal shape functions as per the input shape and protocol used by the service.
func (m protoGetApiKeyMarshaler) marshalOperation(r *aws.Request) {
	var err error
	encoder := rest.NewEncoder(r.HTTPRequest)
	// adds content-type header
	encoder.AddHeader("Content-Type").String("application/json")
	err = marshalGetApiKeyInputShapeAWSREST(m.input, encoder)
	if err != nil {
		r.Error = err
		return
	}
	if err := encoder.Encode(); err != nil {
		r.Error = err
		return
	}

	err = marshalGetApiKeyInputShapeAWSJSON(m.input, r)
	if err != nil {
		r.Error = err
	}
}

// marshalGetApiKeyInputShapeAWSREST is a stand alone function used to marshal the HTTP bindings a input shape.
// This method uses the rest encoder utility
func marshalGetApiKeyInputShapeAWSREST(input *GetApiKeyInput, encoder *rest.Encoder) error {
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

// marshalGetApiKeyInputShapeAWSJSON is a stand alone function used to marshal the json body
func marshalGetApiKeyInputShapeAWSJSON(v *GetApiKeyInput, r *aws.Request) error {
	// delegate to reflection based marshaling
	if t := restlegacy.PayloadType(r.Params); t == "structure" || t == "" {
		jsonrpc.Build(r)
	}
	return nil
}

// getNamedBuildHandler returns a Named Build Handler for an operation marshal function
func (m protoGetApiKeyMarshaler) getNamedBuildHandler() aws.NamedHandler {
	const BuildHandler = "ProtoGetApiKey.BuildHandler"
	return aws.NamedHandler{
		Name: BuildHandler,
		Fn:   m.marshalOperation,
	}
}
