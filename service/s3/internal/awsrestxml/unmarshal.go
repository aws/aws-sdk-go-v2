package awsrestxml

import (
	"bytes"
	"encoding/xml"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/rest"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ProtoGetObjectUnmarshaler defines unmarshaler for ProtoGetObject operation
type ProtoGetObjectUnmarshaler struct {
	Output *types.GetObjectOutput
}

// UnmarshalOperation is the top level method used with a handler stack to unmarshal an operation response
// This method calls appropriate unmarshal shape functions as per the output shape and protocol used by the
// service.
func (u ProtoGetObjectUnmarshaler) UnmarshalOperation(r *aws.Request) {

	var buff bytes.Buffer
	body := io.TeeReader(r.HTTPResponse.Body, &buff)
	// an xml decoder should be used to decode an xml document
	d := xml.NewDecoder(body)
	startToken, err := d.Token()
	if err != nil && err != io.EOF {
		r.Error = err
		return
	}

	// isRequestError checks if operation response returned a error response
	if isRequestError(r) {
		r.Error = UnmarshalError(r, d, startToken)
		return
	}

	// UnmarshalGetObjectOutputShapeBodyAWSREST unmarshal's Rest Components of response
	if err := UnmarshalGetObjectOutputShapeBodyAWSREST(u.Output, r); err != nil {
		r.Error = err
		return
	}

	// UnmarshalGetObjectOutputShapeAWSPayload unmarshal's payload for GetObject output shape
	if err := UnmarshalGetObjectOutputShapeAWSPayload(u.Output, r); err != nil {
		r.Error = err
		return
	}
}

// isRequestError would check if a request response was an error
// It takes in *aws.Request, & a startTag to check for error.
func isRequestError(r *aws.Request) bool {
	if r.HTTPResponse.StatusCode != http.StatusOK {
		return true
	}
	return false
}

// UnmarshalGetObjectOutputShapeBodyAWSREST  is a stand alone function used to a REST component of payload response
// for the REST XML protocol.
// Currently this delegates to reflection based unmarshal meta function
func UnmarshalGetObjectOutputShapeBodyAWSREST(output *types.GetObjectOutput, r *aws.Request) error {
	rest.UnmarshalMeta(r)
	return r.Error
}

// UnmarshalGetObjectOutputShapeAWSPayload is a stand alone function used to unmarshal response body
func UnmarshalGetObjectOutputShapeAWSPayload(output *types.GetObjectOutput, r *aws.Request) error {
	output.Body = r.HTTPResponse.Body
	return r.Error
}

// UnmarshalError unmarshal's the error response
func UnmarshalError(r *aws.Request, d *xml.Decoder, startToken xml.Token) error {
	// xmlErrorResponse is error reponse struct for xml errors.
	// TODO: These types will be eliminated when we do code generated unmarshaling for errors,
	//  we will be walking the document and no type assertion would be required.
	type xmlErrorResponse struct {
		Code      string `xml:"Error>Code"`
		Message   string `xml:"Error>Message"`
		RequestID string `xml:"RequestId"`
	}

	// xmlResponseError wraps xmlErrorResponse struct
	type xmlResponseError struct {
		xmlErrorResponse
	}

	var respErr xmlResponseError
	// delegate to reflection based error decoder
	err := d.Decode(respErr)
	if err != nil && err != io.EOF {
		return awserr.New(aws.ErrCodeSerialization, "Serialization error: Failed to unmarshal error", err)
	}
	reqID := respErr.RequestID
	if len(reqID) == 0 {
		reqID = r.RequestID
	}
	return awserr.NewRequestFailure(awserr.New(respErr.Code, respErr.Message, nil), r.HTTPResponse.StatusCode, reqID)
}

// GetNamedUnmarshalHandler returns a Named Handler for an operation unmarshal function
func (u ProtoGetObjectUnmarshaler) GetNamedUnmarshalHandler() aws.NamedHandler {
	const UnmarshalHandler = "ProtoGetObject.UnmarshalHandler"
	return aws.NamedHandler{
		Name: UnmarshalHandler,
		Fn:   u.UnmarshalOperation,
	}
}
