package rest

import (
	"encoding/base64"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

// Headers is used to encode header keys using a provided prefix
type Headers struct {
	header http.Header
	prefix string
}

// AddHeader returns a HeaderValue used to append values to prefix+key
func (h Headers) AddHeader(key string) HeaderValue {
	return h.newHeaderValue(key, true)
}

// SetHeader returns a HeaderValue used to set the value of prefix+key
func (h Headers) SetHeader(key string) HeaderValue {
	return h.newHeaderValue(key, false)
}

func (h Headers) newHeaderValue(key string, append bool) HeaderValue {
	return newHeaderValue(h.header, h.prefix+strings.TrimSpace(key), append)
}

// HeaderValue is used to encode values to an HTTP header
type HeaderValue struct {
	header http.Header
	key    string
	append bool
}

func newHeaderValue(header http.Header, key string, append bool) HeaderValue {
	return HeaderValue{header: header, key: strings.TrimSpace(key), append: append}
}

func (h HeaderValue) modifyHeader(value string) {
	if h.append {
		h.header.Add(h.key, value)
	} else {
		h.header.Set(h.key, value)
	}
}

// String encodes the value v as the header string value
func (h HeaderValue) String(v string) {
	h.modifyHeader(v)
}

// Byte encodes the value v as a query string value
func (h HeaderValue) Byte(v int8) {
	h.Long(int64(v))
}

// Short encodes the value v as a query string value
func (h HeaderValue) Short(v int16) {
	h.Long(int64(v))
}

// Integer encodes the value v as the header string value
func (h HeaderValue) Integer(v int32) {
	h.Long(int64(v))
}

// Long encodes the value v as the header string value
func (h HeaderValue) Long(v int64) {
	h.modifyHeader(strconv.FormatInt(v, 10))
}

// Boolean encodes the value v as a query string value
func (h HeaderValue) Boolean(v bool) {
	h.modifyHeader(strconv.FormatBool(v))
}

// Float encodes the value v as a query string value
func (h HeaderValue) Float(v float32) {
	h.float(float64(v), 32)
}

// Double encodes the value v as a query string value
func (h HeaderValue) Double(v float64) {
	h.float(v, 64)
}

func (h HeaderValue) float(v float64, bitSize int) {
	h.modifyHeader(strconv.FormatFloat(v, 'f', -1, bitSize))
}

// BigInteger encodes the value v as a query string value
func (h HeaderValue) BigInteger(v *big.Int) {
	h.modifyHeader(v.String())
}

// BigDecimal encodes the value v as a query string value
func (h HeaderValue) BigDecimal(v *big.Float) {
	h.modifyHeader(v.String())
}

// Timestamp encodes the value v using the format name as a header string value
func (h HeaderValue) Timestamp(t time.Time, format string) error {
	value, err := protocol.FormatTime(format, t)
	if err != nil {
		return err
	}
	h.modifyHeader(value)
	return nil
}

// UnixTime encodes the value v using the format name as a query string value
func (h HeaderValue) UnixTime(v time.Time) {
	h.Long(v.Unix())
}

// Blob encodes the value v as a base64 header string value
func (h HeaderValue) Blob(v []byte) {
	encodeToString := base64.StdEncoding.EncodeToString(v)
	h.modifyHeader(encodeToString)
}
