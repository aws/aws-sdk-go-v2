package query

import "net/url"

// Encoder is a Query encoder that supports construction of Query body
// values using methods.
type Encoder struct {
	// The query values that will be built up to manage encoding.
	values url.Values
	Value
}

// NewEncoder returns a new Query body encoder
func NewEncoder() *Encoder {
	values := url.Values{}
	return &Encoder{values, newBaseValue(values)}
}

// String returns the string output of the Query encoder
func (e Encoder) String() string {
	return e.values.Encode()
}

// Bytes returns the []byte slice of the Query encoder
func (e Encoder) Bytes() []byte {
	return []byte(e.values.Encode())
}