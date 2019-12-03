package rest

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
	"net/url"
	"strconv"
	"time"
)

type QueryValue struct {
	query  url.Values
	key    string
	append bool
}

func newQueryValue(query url.Values, key string, append bool) *QueryValue {
	return &QueryValue{
		query:  query,
		key:    key,
		append: append,
	}
}

func (qv *QueryValue) updateKey(value string) {
	if qv.append {
		qv.query.Add(qv.key, value)
	} else {
		qv.query.Set(qv.key, value)
	}
}

func (qv *QueryValue) String(v string) {
	qv.updateKey(v)
}

func (qv *QueryValue) Integer(v int64) {
	qv.updateKey(strconv.FormatInt(v, 10))
}

func (qv *QueryValue) Boolean(v bool) {
	qv.updateKey(strconv.FormatBool(v))
}

func (qv *QueryValue) Float(v float64) {
	qv.updateKey(strconv.FormatFloat(v, 'f', -1, 64))
}

func (qv *QueryValue) Time(v time.Time, format string) error {
	value, err := protocol.FormatTime(format, v)
	if err != nil {
		return err
	}
	qv.updateKey(value)
	return nil
}

func (qv *QueryValue) ByteSlice(v []byte) {
	encodeToString := base64.StdEncoding.EncodeToString(v)
	qv.updateKey(encodeToString)
}

func (qv *QueryValue) JSONValue(v aws.JSONValue) error {
	encodeJSONValue, err := protocol.EncodeJSONValue(v, protocol.NoEscape)
	if err != nil {
		return err
	}
	qv.updateKey(encodeJSONValue)
	return nil
}
