package rest

import (
	"encoding/base64"
	"math/big"
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

// Byte encodes the value v as a query string value
func (qv QueryValue) Byte(v int8) {
	qv.Long(int64(v))
}

// Short encodes the value v as a query string value
func (qv QueryValue) Short(v int16) {
	qv.Long(int64(v))
}

// Integer encodes the value v as a query string value
func (qv QueryValue) Integer(v int32) {
	qv.Long(int64(v))
}

// Long encodes the value v as a query string value
func (qv QueryValue) Long(v int64) {
	qv.updateKey(strconv.FormatInt(v, 10))
}

// Boolean encodes the value v as a query string value
func (qv QueryValue) Boolean(v bool) {
	qv.updateKey(strconv.FormatBool(v))
}

// Float encodes the value v as a query string value
func (qv QueryValue) Float(v float32) {
	qv.float(float64(v), 32)
}

// Double encodes the value v as a query string value
func (qv QueryValue) Double(v float64) {
	qv.float(v, 64)
}

func (qv QueryValue) float(v float64, bitSize int) {
	qv.updateKey(strconv.FormatFloat(v, 'f', -1, bitSize))
}

// BigInteger encodes the value v as a query string value
func (qv QueryValue) BigInteger(v *big.Int) {
	qv.updateKey(v.String())
}

// BigDecimal encodes the value v as a query string value
func (qv QueryValue) BigDecimal(v *big.Float) {
	qv.updateKey(v.String())
}

// Timestamp encodes the value v using the format name as a query string value
func (qv QueryValue) Timestamp(v time.Time, format string) error {
	value, err := protocol.FormatTime(format, v)
	if err != nil {
		return err
	}
	qv.updateKey(value)
	return nil
}

// UnixTime encodes the value v using the format name as a query string value
func (qv QueryValue) UnixTime(v time.Time) {
	qv.Long(v.Unix())
}

// Blob encodes the value v as a base64 query string value
func (qv QueryValue) Blob(v []byte) {
	encodeToString := base64.StdEncoding.EncodeToString(v)
	qv.updateKey(encodeToString)
}

// JSONValue encodes the value v using the format name as a query string value
// deprecated: this will be removed at a later point
func (qv QueryValue) JSONValue(v aws.JSONValue) error {
	encodeJSONValue, err := protocol.EncodeJSONValue(v, protocol.NoEscape)
	if err != nil {
		return err
	}
	qv.updateKey(encodeJSONValue)
	return nil
}
