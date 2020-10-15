package s3shared

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ErrorComponents represents the error response fields
// that will be deserialized from an xml error response body
type ErrorComponents struct {
	Code      string `xml:"Code"`
	Message   string `xml:"Message"`
	RequestID string `xml:"RequestId"`
	HostID    string `xml:"HostId"`
}

// GetErrorResponseComponents returns the error fields from an xml error response body
func GetErrorResponseComponents(r io.Reader, statusCode int, isS3service bool) (ErrorComponents, error) {
	var errComponents ErrorComponents
	if err := xml.NewDecoder(r).Decode(&errComponents); err != nil && err != io.EOF {
		return ErrorComponents{}, fmt.Errorf("error while deserializing xml error response : %w", err)
	}

	// for S3 service, we derive err code and message, if none is found
	if isS3service && len(errComponents.Code) == 0 && len(errComponents.Message) == 0 {
		// derive code and message from status code
		statusText := http.StatusText(statusCode)
		errComponents.Code = strings.Replace(statusText, " ", "", -1)
		errComponents.Message = statusText
	}

	return errComponents, nil
}
