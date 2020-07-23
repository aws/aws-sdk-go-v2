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
	return &Encoder{
		values: values,
		Value:  newBaseValue(values),
	}
}

// Encode returns the []byte slice representing the current
// state of the Query encoder.
func (e Encoder) Encode() ([]byte, error) {
	return []byte(e.values.Encode()), nil
}
