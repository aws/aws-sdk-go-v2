package rest

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

// URIValue is used to encode named URI parameters
type URIValue struct {
	path, rawPath, buffer *[]byte

	key string
}

func newURIValue(path *[]byte, rawPath *[]byte, buffer *[]byte, key string) URIValue {
	return URIValue{path: path, rawPath: rawPath, buffer: buffer, key: key}
}

func (u URIValue) modifyURI(value string) (err error) {
	*u.path, *u.buffer, err = protocol.ReplacePathElement(*u.path, *u.buffer, u.key, value, false)
	*u.rawPath, *u.buffer, err = protocol.ReplacePathElement(*u.rawPath, *u.buffer, u.key, value, true)
	return err
}

// String encodes the value v as a URI string value
func (u URIValue) String(v string) error {
	return u.modifyURI(v)
}

// Integer encodes the value v as a URI string value
func (u URIValue) Integer(v int32) error {
	return u.Long(int64(v))
}

// Long encodes the value v as a URI string value
func (u URIValue) Long(v int64) error {
	return u.modifyURI(strconv.FormatInt(v, 10))
}

// Boolean encodes the value v as a URI string value
func (u URIValue) Boolean(v bool) error {
	return u.modifyURI(strconv.FormatBool(v))
}

// Float encodes the value v as a query string value
func (u URIValue) Float(v float32) error {
	return u.float(float64(v), 32)
}

// Double encodes the value v as a query string value
func (u URIValue) Double(v float64) error {
	return u.float(v, 64)
}

func (u URIValue) float(v float64, bitSize int) error {
	return u.modifyURI(strconv.FormatFloat(v, 'f', -1, bitSize))
}

// Timestamp encodes the value v using the format name as a URI string value
func (u URIValue) Timestamp(v time.Time, format string) error {
	value, err := protocol.FormatTime(format, v)
	if err != nil {
		return err
	}

	return u.modifyURI(value)
}

// UnixTime encodes the value v using the format name as a query string value
func (u URIValue) UnixTime(v time.Time) error {
	return u.Long(v.Unix())
}

// Blob encodes the value v as a base64 URI string value
func (u URIValue) Blob(v []byte) error {
	return u.modifyURI(string(v))
}

// JSONValue encodes the value v as a URI string value
// deprecated: this will be removed at a later point
func (u URIValue) JSONValue(v aws.JSONValue) error {
	encodeJSONValue, err := protocol.EncodeJSONValue(v, protocol.NoEscape)
	if err != nil {
		return err
	}

	return u.modifyURI(encodeJSONValue)
}
