package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restjson"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/internal/awsrestjson"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

const protoOpGetApiKey = "GetApiKey"

// ProtoGetApiKeyRequest returns a request value for making API operation for
// Amazon API Gateway.
//
// Gets information about the current ApiKey resource.
//
//    // Example sending a request using ProtoGetApiKeyRequest.
//    req := client.ProtoGetApiKeyRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *Client) ProtoGetApiKeyRequest(input *types.GetApiKeyInput) ProtoGetApiKeyRequest {
	op := &aws.Operation{
		Name:       protoOpGetApiKey,
		HTTPMethod: "GET",
		HTTPPath:   "/apikeys/{api_Key}",
	}

	if input == nil {
		input = &types.GetApiKeyInput{}
	}

	req := c.newRequest(op, input, &types.GetApiKeyOutput{})
	// swap existing build handler on svc, with a new named build handler
	req.Handlers.Build.Swap(restjson.BuildHandler.Name, awsrestjson.ProtoGetApiKeyMarshaler{Input: input}.GetNamedBuildHandler())
	return ProtoGetApiKeyRequest{Request: req, Input: input, Copy: c.ProtoGetApiKeyRequest}
}

// ProtoGetApiKeyRequest is the request type for the
// GetApiKey API operation.
type ProtoGetApiKeyRequest struct {
	*aws.Request
	Input *types.GetApiKeyInput
	Copy  func(*types.GetApiKeyInput) ProtoGetApiKeyRequest
}

// Send marshals and sends the GetApiKey API request.
func (r ProtoGetApiKeyRequest) Send(ctx context.Context) (*ProtoGetApiKeyResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &ProtoGetApiKeyResponse{
		GetApiKeyOutput: r.Request.Data.(*types.GetApiKeyOutput),
		response:        &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// ProtoGetApiKeyResponse is the response type for the
// GetApiKey API operation.
type ProtoGetApiKeyResponse struct {
	*types.GetApiKeyOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// GetApiKey request.
func (r *ProtoGetApiKeyResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
