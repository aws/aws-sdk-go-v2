package awsrestjson

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/jsonrpc"
	"github.com/aws/aws-sdk-go-v2/private/protocol/rest"
	v2 "github.com/aws/aws-sdk-go-v2/private/protocol/rest/v2"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

// ProtoGetApiKeyMarshaler defines marshaler for ProtoGetApiKey operation
type ProtoGetApiKeyMarshaler struct {
	Input *types.GetApiKeyInput
}

func (m ProtoGetApiKeyMarshaler) MarshalOperation(r *aws.Request) {
	var err error
	encoder := v2.NewEncoder(r.HTTPRequest)
	// adds content-type header
	encoder.AddHeader("Content-Type").String("application/json")

	err = MarshalGetApiKeyInputShapeAWSREST(m.Input, encoder)
	if err != nil {
		r.Error = err
	}

	err = MarshalGetApiKeyInputShapeAWSJSON(m.Input, r)
	if err != nil {
		r.Error = err
	}
}

func MarshalGetApiKeyInputShapeAWSREST(input *types.GetApiKeyInput, encoder *v2.Encoder) error {
	// encode using rest encoder utility
	if err := encoder.SetURI("api_Key").String(*input.ApiKey); err != nil {
		log.Fatalf("Failed to marshal API KEY to URI using REST encoder:\n \t %v", err.Error())
	}
	encoder.AddQuery("includeValue").Boolean(*input.IncludeValue)
	encoder.Encode()
	return nil
}
func MarshalGetApiKeyInputShapeAWSJSON(v *types.GetApiKeyInput, r *aws.Request) error {
	// delegate to reflection based marshaling
	if t := rest.PayloadType(r.Params); t == "structure" || t == "" {
		jsonrpc.Build(r)
	}
	return nil
}

// GetNamedBuildHandler returns a Named Build Handler for an operation marshal function
func (m ProtoGetApiKeyMarshaler) GetNamedBuildHandler() aws.NamedHandler {
	const BuildHandler = "ProtoGetApiKey.BuildHandler"
	return aws.NamedHandler{
		Name: BuildHandler,
		Fn:   m.MarshalOperation,
	}
}
