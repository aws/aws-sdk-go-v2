package ec2query

import (
	"encoding/xml"
	"fmt"
	"io"
)

// errorResponse denotes the error response structure for ec2Query error response
type errorResponse struct {
	Code    string `xml:"Errors>Error>Code"`
	Message string `xml:"Errors>Error>Message"`
}

// GetResponseErrorCode returns the error code, error message from an ec2query error response body
func GetResponseErrorCode(r io.Reader) (code string, message string, err error) {
	var er errorResponse
	if err := xml.NewDecoder(r).Decode(&er); err != nil {
		return code, message, fmt.Errorf("error while fetching xml error response code: %w", err)
	}
	return er.Code, er.Message, nil
}
