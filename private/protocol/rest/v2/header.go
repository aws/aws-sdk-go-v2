package rest

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

type Headers struct {
	header http.Header
	prefix string
}

func (h *Headers) AddHeader(key string) *HeaderValue {
	return h.newHeaderValue(key, true)
}

func (h *Headers) SetHeader(key string) *HeaderValue {
	return h.newHeaderValue(key, false)
}

func (h *Headers) newHeaderValue(key string, append bool) *HeaderValue {
	return newHeaderValue(h.header, h.prefix+strings.TrimSpace(key), append)
}

type HeaderValue struct {
	header http.Header
	key    string
	append bool
}

func newHeaderValue(header http.Header, key string, append bool) *HeaderValue {
	return &HeaderValue{header: header, key: strings.TrimSpace(key), append: append}
}

func (h *HeaderValue) modifyHeader(value string) {
	if h.append {
		h.header.Add(h.key, value)
	} else {
		h.header.Set(h.key, value)
	}
}

func (h *HeaderValue) String(v string) {
	h.modifyHeader(v)
}

func (h *HeaderValue) Integer(v int64) {
	h.modifyHeader(strconv.FormatInt(v, 10))
}

func (h *HeaderValue) Boolean(v bool) {
	h.modifyHeader(strconv.FormatBool(v))
}

func (h *HeaderValue) Float(v float64) {
	h.modifyHeader(strconv.FormatFloat(v, 'f', -1, 64))
}

func (h *HeaderValue) Time(t time.Time, format string) (err error) {
	value, err := protocol.FormatTime(format, t)
	if err != nil {
		return err
	}
	h.modifyHeader(value)
	return nil
}

func (h *HeaderValue) ByteSlice(v []byte) {
	encodeToString := base64.StdEncoding.EncodeToString(v)
	h.modifyHeader(encodeToString)
}

func (h *HeaderValue) JSONValue(v aws.JSONValue) error {
	encodedValue, err := protocol.EncodeJSONValue(v, protocol.Base64Escape)
	if err != nil {
		return err
	}
	h.modifyHeader(encodedValue)
	return nil
}
