package rest

import (
	"encoding/base64"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

// QueryValue is used to encode query key values
type QueryValue struct {
	query  url.Values
	key    string
	append bool
}

func newQueryValue(query url.Values, key string, append bool) QueryValue {
	return QueryValue{
		query:  query,
		key:    key,
		append: append,
	}
}

func (qv QueryValue) updateKey(value string) {
	if qv.append {
		qv.query.Add(qv.key, value)
	} else {
		qv.query.Set(qv.key, value)
	}
}

// String encodes the value v as a query string value
func (qv QueryValue) String(v string) {
	qv.updateKey(v)
}

// Integer encodes the value v as a query string value
func (qv QueryValue) Integer(v int64) {
	qv.updateKey(strconv.FormatInt(v, 10))
}

// Boolean encodes the value v as a query string value
func (qv QueryValue) Boolean(v bool) {
	qv.updateKey(strconv.FormatBool(v))
}

// Float encodes the value v as a query string value
func (qv QueryValue) Float(v float64) {
	qv.updateKey(strconv.FormatFloat(v, 'f', -1, 64))
}

// Time encodes the value v using the format name as a query string value
func (qv QueryValue) Time(v time.Time, format string) error {
	value, err := protocol.FormatTime(format, v)
	if err != nil {
		return err
	}
	qv.updateKey(value)
	return nil
}

// ByteSlice encodes the value v as a base64 query string value
func (qv QueryValue) ByteSlice(v []byte) {
	encodeToString := base64.StdEncoding.EncodeToString(v)
	qv.updateKey(encodeToString)
}

// JSONValue encodes the value v using the format name as a query string value
func (qv QueryValue) JSONValue(v aws.JSONValue) error {
	encodeJSONValue, err := protocol.EncodeJSONValue(v, protocol.NoEscape)
	if err != nil {
		return err
	}
	qv.updateKey(encodeJSONValue)
	return nil
}
