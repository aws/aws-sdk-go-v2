package s3shared

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
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
	rb, err := ioutil.ReadAll(r)
	if err != nil {
		return ErrorComponents{}, err
	}

	var errComponents ErrorComponents
	if err := xml.Unmarshal(rb, &errComponents); err != nil {
		return ErrorComponents{}, fmt.Errorf("error while deserializingg xml error response : %w", err)
	}
	return errComponents, err
}
