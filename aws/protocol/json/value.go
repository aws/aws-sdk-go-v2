package json

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

// Value represents a JSON Value type
// JSON Value types: Object, Array, String, Number, Boolean, and Null
type Value struct {
	w       *bytes.Buffer
	scratch *[]byte
}

// newValue returns a new Value encoder
func newValue(w *bytes.Buffer, scratch *[]byte) Value {
	return Value{w: w, scratch: scratch}
}

// String encodes v as a JSON string
func (jv Value) String(v string) {
	escapeStringBytes(jv.w, []byte(v))
}

// Integer encodes v as a JSON number
func (jv Value) Integer(v int64) {
	*jv.scratch = strconv.AppendInt((*jv.scratch)[:0], v, 10)
	jv.w.Write(*jv.scratch)
}

// Float encodes v as a JSON number
func (jv Value) Float(v float64) {
	*jv.scratch = encodeFloat((*jv.scratch)[:0], v, 64)
	jv.w.Write(*jv.scratch)
}

// Boolean encodes v as a JSON boolean
func (jv Value) Boolean(v bool) {
	*jv.scratch = strconv.AppendBool((*jv.scratch)[:0], v)
	jv.w.Write(*jv.scratch)
}

// ByteSlice encodes v as a base64 value in JSON string
func (jv Value) ByteSlice(v []byte) {
	encodeByteSlice(jv.w, (*jv.scratch)[:0], v)
}

// Time encodes v using the provided format specifier as a JSON string
func (jv Value) Time(v time.Time, format string) error {
	value, err := protocol.FormatTime(format, v)
	if err != nil {
		return err
	}

	escapeStringBytes(jv.w, []byte(value))

	return nil
}

// Array returns a new Array encoder
func (jv Value) Array() *Array {
	return newArray(jv.w, jv.scratch)
}

// Object returns a new Object encoder
func (jv Value) Object() *Object {
	return newObject(jv.w, jv.scratch)
}

// Null encodes a null JSON value
func (jv Value) Null() {
	jv.w.WriteString(null)
}

// Based on encoding/json floatEncoder from the Go Standard Library
// https://golang.org/src/encoding/json/encode.go
func encodeFloat(dst []byte, v float64, bits int) []byte {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		panic(fmt.Sprintf("invalid float value: %s", strconv.FormatFloat(v, 'g', -1, bits)))
	}

	abs := math.Abs(v)
	fmt := byte('f')

	if abs != 0 {
		if bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}

	dst = strconv.AppendFloat(dst, v, fmt, -1, bits)

	if fmt == 'e' {
		// clean up e-09 to e-9
		n := len(dst)
		if n >= 4 && dst[n-4] == 'e' && dst[n-3] == '-' && dst[n-2] == '0' {
			dst[n-2] = dst[n-1]
			dst = dst[:n-1]
		}
	}

	return dst
}

// Based on encoding/json encodeByteSlice from the Go Standard Library
// https://golang.org/src/encoding/json/encode.go
func encodeByteSlice(w *bytes.Buffer, scratch []byte, v []byte) {
	if v == nil {
		w.WriteString(null)
		return
	}

	w.WriteRune(quote)

	encodedLen := base64.StdEncoding.EncodedLen(len(v))
	if encodedLen <= len(scratch) {
		// If the encoded bytes fit in e.scratch, avoid an extra
		// allocation and use the cheaper Encoding.Encode.
		dst := scratch[:encodedLen]
		base64.StdEncoding.Encode(dst, v)
		w.Write(dst)
	} else if encodedLen <= 1024 {
		// The encoded bytes are short enough to allocate for, and
		// Encoding.Encode is still cheaper.
		dst := make([]byte, encodedLen)
		base64.StdEncoding.Encode(dst, v)
		w.Write(dst)
	} else {
		// The encoded bytes are too long to cheaply allocate, and
		// Encoding.Encode is no longer noticeably cheaper.
		enc := base64.NewEncoder(base64.StdEncoding, w)
		enc.Write(v)
		enc.Close()
	}

	w.WriteRune(quote)
}
