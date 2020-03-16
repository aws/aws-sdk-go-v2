package apigateway

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/sdkio"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	jsonlegacy "github.com/aws/aws-sdk-go-v2/private/protocol/json/jsonutil"
	restlegacy "github.com/aws/aws-sdk-go-v2/private/protocol/rest"
)

// protoCreateAPIKeyUnmarshaler defines unmarshaler forProtoCreateAPIKey Operation
type protoCreateAPIKeyUnmarshaler struct {
	output *CreateApiKeyOutput
}

// unmarshalOperation is the top level method used with a handler stack to unmarshal an operation response
// This method calls appropriate unmarshal shape functions as per the output shape and protocol used by the service.
func (u protoCreateAPIKeyUnmarshaler) unmarshalOperation(r *aws.Request) {
	if isRequestError(r) {
		unmarshalError(r)
		return
	}
	restlegacy.UnmarshalMeta(r)
	jsonlegacy.UnmarshalJSON(u.output, r.Body)
}

// isRequestError would check if a request response was an error
func isRequestError(r *aws.Request) bool {
	if r.HTTPResponse.StatusCode == 0 || r.HTTPResponse.StatusCode >= 300 {
		return true
	}
	return false
}

// namedHandler returns a named handler for an operation unmarshal function
func (u protoCreateAPIKeyUnmarshaler) namedHandler() aws.NamedHandler {
	return aws.NamedHandler{
		Name: "ProtoCreateAPIKey.UnmarshalHandler",
		Fn:   u.unmarshalOperation,
	}
}

// unmarshalError unmarshal's an error response.
// some service may have custom error handling
// here we do not handle modelled exceptions.
func unmarshalError(req *aws.Request) {
	defer req.HTTPResponse.Body.Close()
	defer io.Copy(ioutil.Discard, req.HTTPResponse.Body)
	buff := make([]byte, 1024)
	readBuff := make([]byte, 1024)
	ringBuff := sdkio.NewRingBuffer(buff)
	body := io.TeeReader(req.HTTPResponse.Body, ringBuff)
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		ringBuff.Read(readBuff)
		req.Error = awserr.New("SerializationError",
			fmt.Sprintf("failed reading JSON error response, Here's a snapshot %s", readBuff), err)
		return
	}
	if len(bodyBytes) == 0 {
		req.Error = awserr.NewRequestFailure(
			awserr.New("SerializationError", req.HTTPResponse.Status, nil),
			req.HTTPResponse.StatusCode,
			"",
		)
		return
	}
	var jsonErr jsonErrorResponse
	if err := json.Unmarshal(bodyBytes, &jsonErr); err != nil {
		req.Error = awserr.New("SerializationError", "failed decoding JSON RPC error response", err)
		return
	}

	codes := strings.SplitN(jsonErr.Code, "#", 2)
	req.Error = awserr.NewRequestFailure(
		awserr.New(codes[len(codes)-1], jsonErr.Message, nil),
		req.HTTPResponse.StatusCode,
		req.RequestID,
	)
}

type jsonErrorResponse struct {
	Code    string `json:"__type"`
	Message string `json:"message"`
}
