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

// Encode returns a REST protocol encoder for encoding HTTP bindings
// Returns any error if one occurred during encoding.
func (e *Encoder) Encode() error {
	e.req.URL.Path, e.req.URL.RawPath = string(e.path), string(e.rawPath)
	e.req.URL.RawQuery = e.query.Encode()
	e.req.Header = e.header

	return nil
}

// AddHeader returns a HeaderValue for appending to the given header name
func (e *Encoder) AddHeader(key string) HeaderValue {
	return newHeaderValue(e.header, key, true)
}

// SetHeader returns a HeaderValue for setting the given header name
func (e *Encoder) SetHeader(key string) HeaderValue {
	return newHeaderValue(e.header, key, false)
}

// Headers returns a Header used encoding headers with the given prefix
func (e *Encoder) Headers(prefix string) Headers {
	return Headers{
		header: e.header,
		prefix: strings.TrimSpace(prefix),
	}
}

// SetURI returns a URIValue used for setting the given path key
func (e *Encoder) SetURI(key string) URIValue {
	return newURIValue(&e.path, &e.rawPath, &e.pathBuffer, key)
}

// SetQuery returns a QueryValue used for setting the given query key
func (e *Encoder) SetQuery(key string) QueryValue {
	return newQueryValue(e.query, key, false)
}

// AddQuery returns a QueryValue used for appending the given query key
func (e *Encoder) AddQuery(key string) QueryValue {
	return newQueryValue(e.query, key, true)
}
