package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restjson"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/internal/awsrestjson"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

const protoOpCreateApiKey = "CreateApiKey"

// ProtoCreateApiKeyRequest returns a request value for making API operation for
// Amazon API Gateway.
//
// Create an ApiKey resource.
//
// AWS CLI (https://docs.aws.amazon.com/cli/latest/reference/apigateway/create-api-key.html)
//
//    // Example sending a request using ProtoCreateApiKeyRequest.
//    req := client.ProtoCreateApiKeyRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *Client) ProtoCreateApiKeyRequest(input *types.CreateApiKeyInput) ProtoCreateApiKeyRequest {
	op := &aws.Operation{
		Name:       protoOpCreateApiKey,
		HTTPMethod: "POST",
		HTTPPath:   "/apikeys",
	}

	if input == nil {
		input = &types.CreateApiKeyInput{}
	}

	req := c.newRequest(op, input, &types.CreateApiKeyOutput{})
	req.Handlers.Build.Swap(restjson.BuildHandler.Name, awsrestjson.ProtoCreateApiKeyMarshaler{Input: input}.GetNamedBuildHandler())

	return ProtoCreateApiKeyRequest{Request: req, Input: input, Copy: c.ProtoCreateApiKeyRequest}
}

// ProtoCreateApiKeyRequest is the request type for the
// CreateApiKey API operation.
type ProtoCreateApiKeyRequest struct {
	*aws.Request
	Input *types.CreateApiKeyInput
	Copy  func(*types.CreateApiKeyInput) ProtoCreateApiKeyRequest
}

// Send marshals and sends the CreateApiKey API request.
func (r ProtoCreateApiKeyRequest) Send(ctx context.Context) (*ProtoCreateApiKeyResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &ProtoCreateApiKeyResponse{
		CreateApiKeyOutput: r.Request.Data.(*types.CreateApiKeyOutput),
		response:           &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// ProtoCreateApiKeyResponse is the response type for the
// CreateApiKey API operation.
type ProtoCreateApiKeyResponse struct {
	*types.CreateApiKeyOutput
	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// CreateApiKey request.
func (r *ProtoCreateApiKeyResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
