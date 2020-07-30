package query

import (
	"io"
	"net/url"
	"sort"
)

// Encoder is a Query encoder that supports construction of Query body
// values using methods.
type Encoder struct {
	// The query values that will be built up to manage encoding.
	values url.Values
	// The writer that the encoded body will be written to.
	writer io.Writer
	Value
}

// NewEncoder returns a new Query body encoder
func NewEncoder(writer io.Writer) *Encoder {
	values := url.Values{}
	return &Encoder{
		values: values,
		writer: writer,
		Value:  newBaseValue(values),
	}
}

// Encode returns the []byte slice representing the current
// state of the Query encoder.
func (e Encoder) Encode() error {
	// Get the keys and sort them to have a stable output
	keys := make([]string, 0, len(e.values))
	for k := range e.values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	isFirstEntry := true
	for _, key := range keys {
		queryValues := e.values[key]
		escapedKey := url.QueryEscape(key)
		for _, value := range queryValues {
			if !isFirstEntry {
				if _, err := e.writer.Write([]byte(`&`)); err != nil {
					return err
				}
			} else {
				isFirstEntry = false
			}
			if _, err := e.writer.Write([]byte(escapedKey)); err != nil {
				return err
			}
			if _, err := e.writer.Write([]byte(`=`)); err != nil {
				return err
			}
			if _, err := e.writer.Write([]byte(url.QueryEscape(value))); err != nil {
				return err
			}
		}
	}
	return nil
}
