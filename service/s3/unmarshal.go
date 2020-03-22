package s3

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/internal/s3err"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
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
// This method calls appropriate unmarshal shape functions as per the output shape and protocol used by the service.
func (u protoGetObjectUnmarshaler) unmarshalOperation(r *aws.Request) {
	// isRequestError checks if operation response returned a error response
	if isRequestError(r) {
		defer r.HTTPResponse.Body.Close()
		defer io.Copy(ioutil.Discard, r.HTTPResponse.Body)

		buff := make([]byte, 1024)
		readBuff := make([]byte, 1024)
		ringBuff := sdkio.NewRingBuffer(buff)
		body := io.TeeReader(r.HTTPResponse.Body, ringBuff)
		decoder := xml.NewDecoder(body)

		// recurse thru the xml body till we get the startElement,
		// we ignore the xml preamble.
		for {
			startToken, err := decoder.Token()
			if err != nil {
				ringBuff.Read(readBuff)
				r.Error = awserr.New(aws.ErrCodeSerialization, fmt.Sprintf("Failed to decode response body with invalid XML. Here's a snapshot : %v", readBuff), err)
				return
			}
			// deligate to error unmarshaler if startElement is retrieved
			if start, ok := startToken.(xml.StartElement); ok {
				r.Error = unmarshalErrorPrototype(r, decoder, start, ringBuff)
				return
			}
		}
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
	HostID    string `xml:"HostId"`
}

// unmarshalErrorPrototype unmarshal's the error response.
// The function takes in a startElement; this is to support 200 Errors
// where we will be reading the startElement.
// Note: Services may have customizations that require custom error unmarshaling.
//  also, in this prototype, we do not prototype modeled error unmarshaling.
func unmarshalErrorPrototype(r *aws.Request, d *xml.Decoder, start xml.StartElement, buffer *sdkio.RingBuffer) error {
	// protoXMLErrorResponse is error response struct for xml errors.
	var respErr = protoXMLErrorResponse{}

	// Delegate to reflection based decoding utils
	err := d.DecodeElement(&respErr, &start)
	if err != nil && err != io.EOF {
		readBuff := make([]byte, 1024)
		buffer.Read(readBuff)
		return awserr.New(aws.ErrCodeSerialization, fmt.Sprintf("Failed to unmarshal error with invalid XML. Here's a snapshot : %v", buffer), err)
	}

	reqID := respErr.RequestID
	if len(reqID) == 0 {
		reqID = r.RequestID
	}

	hostID := respErr.HostID
	if len(hostID) == 0 {
		hostID = r.HTTPResponse.Header.Get("X-Amz-Id-2")
	}

	// return s3 specific error
	return s3err.NewRequestFailure(
		awserr.NewRequestFailure(
			awserr.New(respErr.Code, respErr.Message, nil),
			r.HTTPResponse.StatusCode, reqID),
		hostID)
}

// namedHandler returns a named handler for an operation unmarshal function
func (u protoGetObjectUnmarshaler) namedHandler() aws.NamedHandler {
	return aws.NamedHandler{
		Name: "ProtoGetObject.UnmarshalHandler",
		Fn:   u.unmarshalOperation,
	}
}
