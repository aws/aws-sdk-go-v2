package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restjson"
)

const protoOpCreateAPIKey = "CreateApiKey"

// ProtoCreateAPIKeyRequest returns a request value for making API operation for
// Amazon API Gateway.
//
// Create an ApiKey resource.
//
// AWS CLI (https://docs.aws.amazon.com/cli/latest/reference/apigateway/create-api-key.html)
//
//    // Example sending a request using CreateApiKeyRequest.
//    req := client.CreateApiKeyRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *Client) ProtoCreateAPIKeyRequest(input *CreateApiKeyInput) ProtoCreateAPIKeyRequest {
	op := &aws.Operation{
		Name:       protoOpCreateAPIKey,
		HTTPMethod: "POST",
		HTTPPath:   "/apikeys",
	}

	if input == nil {
		input = &CreateApiKeyInput{}
	}

	req := c.newRequest(op, input, &CreateApiKeyOutput{})
	// swap existing build handler on svc, with a new named build handler
	req.Handlers.Build.Swap(restjson.BuildHandler.Name, protoCreateAPIKeyMarshaler{input: input}.namedHandler())
	// swap existing build handler on svc, with a new named build handler
	req.Handlers.Unmarshal.Swap(restjson.UnmarshalHandler.Name, protoCreateAPIKeyUnmarshaler{output: output}.namedHandler())
	return ProtoCreateAPIKeyRequest{Request: req, Input: input, Copy: c.ProtoCreateAPIKeyRequest}
}

// ProtoCreateAPIKeyRequest is the request type for the
// ProtoCreateApiKey API operation.
type ProtoCreateAPIKeyRequest struct {
	*aws.Request
	Input *CreateApiKeyInput
	Copy  func(*CreateApiKeyInput) ProtoCreateAPIKeyRequest
}

// Send marshals and sends the CreateApiKey API request.
func (r ProtoCreateAPIKeyRequest) Send(ctx context.Context) (*ProtoCreateAPIKeyResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &ProtoCreateAPIKeyResponse{
		CreateApiKeyOutput: r.Request.Data.(*CreateApiKeyOutput),
		response:           &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// ProtoCreateAPIKeyResponse is the response type for the
// ProtoCreateApiKey API operation.
type ProtoCreateAPIKeyResponse struct {
	*CreateApiKeyOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// CreateApiKey request.
func (r *ProtoCreateAPIKeyResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
