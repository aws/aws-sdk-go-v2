package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restjson"
)

const protoOpGetAPIKey = "GetApiKey"

// ProtoGetAPIKeyRequest returns a request value for making API operation for
// Amazon API Gateway.
//
// Gets information about the current ApiKey resource.
//
//    // Example sending a request using ProtoGetAPIKeyRequest.
//    req := client.ProtoGetAPIKeyRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
func (c *Client) ProtoGetAPIKeyRequest(input *GetApiKeyInput) ProtoGetAPIKeyRequest {
	op := &aws.Operation{
		Name:       protoOpGetAPIKey,
		HTTPMethod: "GET",
		HTTPPath:   "/apikeys/{api_Key}",
	}

	if input == nil {
		input = &GetApiKeyInput{}
	}

	req := c.newRequest(op, input, &GetApiKeyOutput{})
	// swap existing build handler on svc, with a new named build handler
	req.Handlers.Build.Swap(restjson.BuildHandler.Name, protoGetAPIKeyMarshaler{input: input}.NamedHandler())
	return ProtoGetAPIKeyRequest{Request: req, Input: input, Copy: c.ProtoGetAPIKeyRequest}
}

// ProtoGetAPIKeyRequest is the request type for the
// GetApiKey API operation.
type ProtoGetAPIKeyRequest struct {
	*aws.Request
	Input *GetApiKeyInput
	Copy  func(*GetApiKeyInput) ProtoGetAPIKeyRequest
}

// Send marshals and sends the GetApiKey API request.
func (r ProtoGetAPIKeyRequest) Send(ctx context.Context) (*ProtoGetAPIKeyResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &ProtoGetAPIKeyResponse{
		GetApiKeyOutput: r.Request.Data.(*GetApiKeyOutput),
		response:        &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// ProtoGetAPIKeyResponse is the response type for the
// GetApiKey API operation.
type ProtoGetAPIKeyResponse struct {
	*GetApiKeyOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// GetApiKey request.
func (r *ProtoGetAPIKeyResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
