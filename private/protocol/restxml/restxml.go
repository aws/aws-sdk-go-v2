// Package restxml provides RESTful XML serialization of AWS
// requests and responses.
package restxml

//go:generate go run -tags codegen ../../../models/protocol_tests/generate.go ../../../models/protocol_tests/input/rest-xml.json build_test.go
//go:generate go run -tags codegen ../../../models/protocol_tests/generate.go ../../../models/protocol_tests/output/rest-xml.json unmarshal_test.go

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	request "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
	"github.com/aws/aws-sdk-go-v2/private/protocol/query"
	"github.com/aws/aws-sdk-go-v2/private/protocol/rest"
	"github.com/aws/aws-sdk-go-v2/private/protocol/xml/xmlutil"
)

// BuildHandler is a named request handler for building restxml protocol requests
var BuildHandler = request.NamedHandler{Name: "awssdk.restxml.Build", Fn: Build}

// UnmarshalHandler is a named request handler for unmarshaling restxml protocol requests
var UnmarshalHandler = request.NamedHandler{Name: "awssdk.restxml.Unmarshal", Fn: Unmarshal}

// UnmarshalMetaHandler is a named request handler for unmarshaling restxml protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler{Name: "awssdk.restxml.UnmarshalMeta", Fn: UnmarshalMeta}

// UnmarshalErrorHandler is a named request handler for unmarshaling restxml protocol request errors
var UnmarshalErrorHandler = request.NamedHandler{Name: "awssdk.restxml.UnmarshalError", Fn: UnmarshalError}

// xmlUnmarshaler is an interface that a shape can implement
type xmlUnmarshaler interface {
	UnmarshalAWSXML(*xml.Decoder) error
}

// restUnmarshaler is an interface that a shape can implement
type restUnmarshaler interface {
	UnmarshalAWSREST(*http.Response) error
}

// payloadUnmarshaler is an interface that a shape can implement
type payloadUnmarshaler interface {
	UnmarshalAWSPayload(io.ReadCloser) error
}

// Build builds a request payload for the REST XML protocol.
func Build(r *request.Request) {
	if m, ok := r.Params.(protocol.FieldMarshaler); ok {
		e := NewEncoder(r.HTTPRequest)

		m.MarshalFields(e)

		var body io.ReadSeeker
		var err error
		r.HTTPRequest, body, err = e.Encode()
		if err != nil {
			r.Error = awserr.New(request.ErrCodeSerialization, "failed to encode rest XML request", err)
			return
		}
		if body != nil {
			r.SetReaderBody(body)
		}
		return
	}

	// Fall back to old reflection based marshaler
	rest.Build(r)

	if t := rest.PayloadType(r.Params); t == "structure" || t == "" {
		var buf bytes.Buffer
		err := xmlutil.BuildXML(r.Params, xml.NewEncoder(&buf))
		if err != nil {
			r.Error = awserr.New("SerializationError", "failed to encode rest XML request", err)
			return
		}
		r.SetBufferBody(buf.Bytes())
	}
}

// Unmarshal unmarshals a payload response for the REST XML protocol.
func Unmarshal(r *request.Request) {
	hasGeneratedUmarshaler := false
	if resp, ok := r.Data.(restUnmarshaler); ok {
		hasGeneratedUmarshaler = true
		err := resp.UnmarshalAWSREST(r.HTTPResponse)
		if err != nil {
			r.Error = awserr.New("SerializationError", "failed to decode REST XML response", err)
			return
		}
	}
	if resp, ok := r.Data.(payloadUnmarshaler); ok {
		hasGeneratedUmarshaler = true
		err := resp.UnmarshalAWSPayload(r.HTTPResponse.Body)
		if err != nil {
			r.Error = awserr.New("SerializationError", "failed to decode REST XML response", err)
			return
		}
	}
	if resp, ok := r.Data.(xmlUnmarshaler); ok {
		hasGeneratedUmarshaler = true
		defer r.HTTPResponse.Body.Close()
		decoder := xml.NewDecoder(r.HTTPResponse.Body)
		err := resp.UnmarshalAWSXML(decoder)
		if err != nil {
			r.Error = awserr.New("SerializationError", "failed to decode REST XML response", err)
			return
		}
	}
	if hasGeneratedUmarshaler {
		return
	}

	// Fall back to old reflection based unmarshaler
	if t := rest.PayloadType(r.Data); t == "structure" || t == "" {
		defer r.HTTPResponse.Body.Close()
		decoder := xml.NewDecoder(r.HTTPResponse.Body)
		err := xmlutil.UnmarshalXML(r.Data, decoder, "")
		if err != nil {
			r.Error = awserr.New("SerializationError", "failed to decode REST XML response", err)
			return
		}
	} else {
		rest.Unmarshal(r)
	}
}

// UnmarshalMeta unmarshals response headers for the REST XML protocol.
func UnmarshalMeta(r *request.Request) {
	/*
	   If r.Data has implemented the restUmarshaler interface, the header and status code unmarshaling can be handled by the function
	   shape.UnmarshalAWSREST(*http.Response). Then the UnmarshalMeta(*request.Request) function only needs to handle requestID unmarshaling.
	*/
	if _, ok := r.Data.(restUnmarshaler); ok {
		r.RequestID = r.HTTPResponse.Header.Get("X-Amzn-Requestid")
		if r.RequestID == "" {
			// Alternative version of request id in the header
			r.RequestID = r.HTTPResponse.Header.Get("X-Amz-Request-Id")
		}
		return
	}

	// Fall back to old reflection based unmarshaler
	rest.UnmarshalMeta(r)
}

// UnmarshalError unmarshals a response error for the REST XML protocol.
func UnmarshalError(r *request.Request) {
	query.UnmarshalError(r)
}
