package s3shared

import (
	"encoding/xml"
	"fmt"
	"io"
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
func GetErrorResponseComponents(r io.Reader) (ErrorComponents, error) {
	var errComponents ErrorComponents
	if err := xml.NewDecoder(r).Decode(&errComponents); err != nil && err != io.EOF {
		return ErrorComponents{}, fmt.Errorf("error while deserializing xml error response : %w", err)
	}
	return errComponents, nil
}
