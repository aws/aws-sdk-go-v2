package http

import (
	"encoding/xml"
	"io"
	"io/ioutil"

	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// XMLRequestIDRetriever used by retriever to retrieve xml response request id
type XMLRequestIDRetriever struct {
	// header representing the request id member
	header string
	nowrap bool
}

// GetRequestID returns request id
func (r XMLRequestIDRetriever) GetRequestID(resp smithyhttp.Response, reader io.Reader) string {
	if v := resp.Header.Get(r.header); len(v) != 0 {
		return v
	}

	type errorResponse struct {
		requestID string `xml:"ErrorResponse>RequestId"`
	}

	type errorResponseNoWrap struct {
		requestID string `xml:"Error>RequestId"`
	}

	rb, _ := ioutil.ReadAll(reader)
	if r.nowrap {
		var errResponse errorResponseNoWrap
		xml.Unmarshal(rb, &errResponse)
		return errResponse.requestID
	}

	var errResponse errorResponse
	xml.Unmarshal(rb, &errResponse)
	return errResponse.requestID
}

// QueryRequestIDRetriever used by retriever to retrieve query response request id
type QueryRequestIDRetriever struct {
	// header representing the request id member
	header string
}

// GetRequestID returns request id
func (r QueryRequestIDRetriever) GetRequestID(resp smithyhttp.Response, reader io.Reader) string {
	if v := resp.Header.Get(r.header); len(v) != 0 {
		return v
	}

	type errorResponse struct {
		requestID string `xml:"ErrorResponse>RequestId"`
	}

	rb, _ := ioutil.ReadAll(reader)
	var errResponse errorResponse
	xml.Unmarshal(rb, &errResponse)
	return errResponse.requestID
}

// EC2QueryRequestIDRetriever used by retriever to retrieve ec2query response request id
type EC2QueryRequestIDRetriever struct {
	// header representing the request id member
	header string
}

// GetRequestID returns request id
func (r EC2QueryRequestIDRetriever) GetRequestID(resp smithyhttp.Response, reader io.Reader) string {
	if v := resp.Header.Get(r.header); len(v) != 0 {
		return v
	}

	type errorResponse struct {
		requestID string `xml:"Response>RequestId"`
	}

	rb, _ := ioutil.ReadAll(reader)
	var errResponse errorResponse
	xml.Unmarshal(rb, &errResponse)
	return errResponse.requestID
}

// JSONRequestIDRetriever used by retriever to retrieve json response request id
type JSONRequestIDRetriever struct {
	// header representing the request id member
	header string
}

// GetRequestID returns request id
func (r JSONRequestIDRetriever) GetRequestID(resp smithyhttp.Response, reader io.Reader) string {
	if v := resp.Header.Get(r.header); len(v) != 0 {
		return v
	}
	return ""
}
