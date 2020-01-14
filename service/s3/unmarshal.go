package s3

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	restlegacy "github.com/aws/aws-sdk-go-v2/private/protocol/rest"
)

/*
   GetObject operation doesn't support 200 status code with errors.

*/

// protoGetObjectUnmarshaler defines unmarshaler for ProtoGetObject operation
type protoGetObjectUnmarshaler struct {
	output *GetObjectOutput
}

// unmarshalOperation is the top level method used with a handler stack to unmarshal an operation response
// This method calls appropriate unmarshal shape functions as per the output shape and protocol used by the
// service.
func (u protoGetObjectUnmarshaler) unmarshalOperation(r *aws.Request) {

	var (
		buff       bytes.Buffer
		body       io.Reader
		decoder    *xml.Decoder
		startToken xml.Token
		err        error
	)

	if r.HTTPResponse.Body != nil {
		body = io.TeeReader(r.HTTPResponse.Body, &buff)
		decoder = xml.NewDecoder(body)
		startToken, err = decoder.Token()
		// we will encounter io.EOF if body does not have an xml
		if err != nil && err != io.EOF {
			r.Error = awserr.New(aws.ErrCodeSerialization, "Failed to decode response body with invalid XML", err)
			return
		}
	}

	// isRequestError checks if operation response returned a error response
	if isRequestError(r) {
		// if startToken is nil, it would mean there is no xml response body
		if startToken == nil {
			// fall back to status code error message
			statusText := http.StatusText(r.HTTPResponse.StatusCode)
			errCode := strings.Replace(statusText, " ", "", -1)
			errMsg := statusText
			r.Error = awserr.NewRequestFailure(
				awserr.New(errCode, errMsg, err),
				r.HTTPResponse.StatusCode,
				r.RequestID,
			)
			return
		}

		r.Error = unmarshalErrorPrototype(r, decoder, startToken)
		return
	}

	// delegate to reflection based rest unmarshaler
	restlegacy.UnmarshalMeta(r)

	// payload unmarshal would directly assign the response payload to unmarshal output.
	u.output.Body = r.HTTPResponse.Body
}

// isRequestError would check if a request response was an error
// This method should also take in startTag to check for error if operation supports 200 errors.
func isRequestError(r *aws.Request) bool {
	if r.HTTPResponse.StatusCode == 0 || r.HTTPResponse.StatusCode >= 300 {
		return true
	}
	return false
}

type protoXMLErrorResponse struct {
	Code      string `xml:"Code"`
	Message   string `xml:"Message"`
	RequestID string `xml:"RequestId"`
}

// unmarshalError unmarshal's the error response
func unmarshalErrorPrototype(r *aws.Request, d *xml.Decoder, startToken xml.Token) error {
	// protoXMLErrorResponse is error response struct for xml errors.
	var respErr = protoXMLErrorResponse{}

	// Delegate to reflection based decoding utils
	if start, ok := startToken.(xml.StartElement); ok {
		err := d.DecodeElement(&respErr, &start)
		if err != nil && err != io.EOF {
			return awserr.New(aws.ErrCodeSerialization, "Serialization error: Failed to unmarshal error", err)
		}
	} else {
		return awserr.New(aws.ErrCodeSerialization, "Serialization error: Failed to unmarshal invalid xml", nil)
	}

	reqID := respErr.RequestID
	if len(reqID) == 0 {
		reqID = r.RequestID
	}
	return awserr.NewRequestFailure(awserr.New(respErr.Code, respErr.Message, nil), r.HTTPResponse.StatusCode, reqID)
}

// NamedHandler returns a Named Handler for an operation unmarshal function
func (u protoGetObjectUnmarshaler) NamedHandler() aws.NamedHandler {
	const unmarshalHandler = "ProtoGetObject.UnmarshalHandler"
	return aws.NamedHandler{
		Name: unmarshalHandler,
		Fn:   u.unmarshalOperation,
	}
}
