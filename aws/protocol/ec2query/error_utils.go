package ec2query

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

// errorResponse denotes the error response structure for ec2Query error response
type errorResponse struct {
	Code string `xml:"Errors>Error>Code"`
}

// GetEc2QueryResponseErrorCode returns the error code from an ec2query error response body
func GetResponseErrorCode(r io.Reader) (string, error) {
	rb, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	var er errorResponse
	if err := xml.Unmarshal(rb, &er); err != nil {
		return "", fmt.Errorf("error while fetching xml error response code: %w", err)
	}
	return er.Code, nil
}
