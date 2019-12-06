package rest

import (
	"net/http"
	"net/url"
	"strings"
)

// An Encoder provides encoding of REST URI path, query, and header components
// of an HTTP request. Can also encode a stream as the payload.
//
// Does not support SetFields.
type Encoder struct {
	req *http.Request

	path, rawPath, pathBuffer []byte

	query  url.Values
	header http.Header
}

// NewEncoder creates a new encoder from the passed in request. All query and
// header values will be added on top of the request's existing values. Overwriting
// duplicate values.
func NewEncoder(req *http.Request) *Encoder {
	e := &Encoder{
		req: req,

		path:    []byte(req.URL.Path),
		rawPath: []byte(req.URL.Path),
		query:   req.URL.Query(),
		header:  req.Header,
	}

	return e
}

// Encode will return the request and body if one was set. If the body
// payload was not set the io.ReadSeeker will be nil.
//
// returns any error if one occured while encoding the API's parameters.
func (e *Encoder) Encode() *http.Request {
	e.req.URL.Path, e.req.URL.RawPath = string(e.path), string(e.rawPath)
	e.req.URL.RawQuery = e.query.Encode()
	e.req.Header = e.header

	return e.req
}

func (e *Encoder) AddHeader(key string) *HeaderValue {
	return newHeaderValue(e.header, key, true)
}

func (e *Encoder) SetHeader(key string) *HeaderValue {
	return newHeaderValue(e.header, key, false)
}

func (e *Encoder) Headers(prefix string) *Headers {
	return &Headers{
		header: e.header,
		prefix: strings.TrimSpace(prefix),
	}
}

func (e *Encoder) SetURI(key string) *URIValue {
	return newURIValue(&e.path, &e.rawPath, &e.pathBuffer, key)
}

func (e *Encoder) SetQuery(key string) *QueryValue {
	return newQueryValue(e.query, key, false)
}

func (e *Encoder) AddQuery(key string) *QueryValue {
	return newQueryValue(e.query, key, true)
}
