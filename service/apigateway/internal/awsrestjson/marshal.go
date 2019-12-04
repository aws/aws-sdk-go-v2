package awsrestjson

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/jsonrpc"
	"github.com/aws/aws-sdk-go-v2/private/protocol/rest"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

// ProtoCreateApiKeyMarshaler defines marshaler for ProtoCreateApiKey operation
type ProtoCreateApiKeyMarshaler struct {
	Input *types.CreateApiKeyInput
}

func (m ProtoCreateApiKeyMarshaler) MarshalOperation(r *aws.Request) {
	var err error

	err = MarshalCreateApiKeyInputShapeAWSREST(m.Input, r)
	if err != nil {
		r.Error = err
	}

	err = MarshalCreateApiKeyInputShapeAWSJSON(m.Input, r)
	if err != nil {
		r.Error = err
	}

}

func MarshalCreateApiKeyInputShapeAWSREST(v *types.CreateApiKeyInput, r *aws.Request) error {
	// delegate to reflection based marshaling
	rest.Build(r)
	return nil
}
func MarshalCreateApiKeyInputShapeAWSJSON(v *types.CreateApiKeyInput, r *aws.Request) error {
	// delegate to reflection based marshaling
	if t := rest.PayloadType(r.Params); t == "structure" || t == "" {
		jsonrpc.Build(r)
	}
	return nil
}

// GetNamedBuildHandler returns a Named Build Handler for an operation marshal function
func (m ProtoCreateApiKeyMarshaler) GetNamedBuildHandler() aws.NamedHandler {
	const BuildHandler = "CreateApiKey.BuildHandler"
	return aws.NamedHandler{
		Name: BuildHandler,
		Fn:   m.MarshalOperation,
	}
}
