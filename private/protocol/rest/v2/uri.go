package rest

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

type URIValue struct {
	path, rawPath, buffer *[]byte

	key string
}

func newURIValue(path *[]byte, rawPath *[]byte, buffer *[]byte, key string) *URIValue {
	return &URIValue{path: path, rawPath: rawPath, buffer: buffer, key: key}
}

func (u *URIValue) modifyURI(value string) (err error) {
	*u.path, *u.buffer, err = protocol.ReplacePathElement(*u.path, *u.buffer, u.key, value, false)
	*u.rawPath, *u.buffer, err = protocol.ReplacePathElement(*u.rawPath, *u.buffer, u.key, value, true)
	return err
}

func (u *URIValue) String(v string) error {
	return u.modifyURI(v)
}

func (u *URIValue) Integer(v int64) error {
	return u.modifyURI(strconv.FormatInt(v, 10))
}

func (u *URIValue) Boolean(v bool) error {
	return u.modifyURI(strconv.FormatBool(v))
}

func (u *URIValue) Float(v float64) error {
	return u.modifyURI(strconv.FormatFloat(v, 'f', -1, 64))
}

func (u *URIValue) Time(v time.Time, format string) error {
	value, err := protocol.FormatTime(format, v)
	if err != nil {
		return err
	}

	return u.modifyURI(value)
}

func (u *URIValue) ByteSlice(v []byte) error {
	return u.modifyURI(string(v))
}

func (u *URIValue) JSONValue(v aws.JSONValue) error {
	encodeJSONValue, err := protocol.EncodeJSONValue(v, protocol.NoEscape)
	if err != nil {
		return err
	}

	return u.modifyURI(encodeJSONValue)
}
