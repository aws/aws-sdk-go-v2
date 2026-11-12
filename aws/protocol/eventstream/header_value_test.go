package eventstream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func binWrite(v any) []byte {
	var w bytes.Buffer
	binary.Write(&w, binary.BigEndian, v)
	return w.Bytes()
}

var testValueEncodingCases = []struct {
	Val    Value
	Expect []byte
	Decode func(io.Reader) (Value, error)
}{
	{
		BoolValue(true),
		[]byte{byte(trueValueType)},
		nil,
	},
	{
		BoolValue(false),
		[]byte{byte(falseValueType)},
		nil,
	},
	{
		Int8Value(0x0f),
		[]byte{byte(int8ValueType), 0x0f},
		func(r io.Reader) (Value, error) {
			var v Int8Value
			err := v.decode(r)
			return v, err
		},
	},
	{
		Int16Value(0x0f),
		append([]byte{byte(int16ValueType)}, binWrite(int16(0x0f))...),
		func(r io.Reader) (Value, error) {
			var v Int16Value
			err := v.decode(r)
			return v, err
		},
	},
	{
		Int32Value(0x0f),
		append([]byte{byte(int32ValueType)}, binWrite(int32(0x0f))...),
		func(r io.Reader) (Value, error) {
			var v Int32Value
			err := v.decode(r)
			return v, err
		},
	},
	{
		Int64Value(0x0f),
		append([]byte{byte(int64ValueType)}, binWrite(int64(0x0f))...),
		func(r io.Reader) (Value, error) {
			var v Int64Value
			err := v.decode(r)
			return v, err
		},
	},
	{
		BytesValue([]byte{0, 1, 2, 3}),
		[]byte{byte(bytesValueType), 0x00, 0x04, 0, 1, 2, 3},
		func(r io.Reader) (Value, error) {
			var v BytesValue
			err := v.decode(r)
			return v, err
		},
	},
	{
		StringValue("abc123"),
		append([]byte{byte(stringValueType), 0, 6}, []byte("abc123")...),
		func(r io.Reader) (Value, error) {
			var v StringValue
			err := v.decode(r)
			return v, err
		},
	},
	{
		TimestampValue(
			time.Date(2014, 04, 04, 0, 1, 0, 0, time.FixedZone("PDT", -7)),
		),
		append([]byte{byte(timestampValueType)}, binWrite(int64(1396569667000))...),
		func(r io.Reader) (Value, error) {
			var v TimestampValue
			err := v.decode(r)
			return v, err
		},
	},
	{
		UUIDValue(
			[16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf},
		),
		[]byte{byte(uuidValueType), 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf},
		func(r io.Reader) (Value, error) {
			var v UUIDValue
			err := v.decode(r)
			return v, err
		},
	},
}

func TestValue_MarshalValue(t *testing.T) {
	for i, c := range testValueEncodingCases {
		var w bytes.Buffer

		if err := c.Val.encode(&w); err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}

		if e, a := c.Expect, w.Bytes(); !reflect.DeepEqual(e, a) {
			t.Errorf("%d, expect %v, got %v", i, e, a)
		}
	}
}

func TestHeader_DecodeValues(t *testing.T) {
	for i, c := range testValueEncodingCases {
		r := bytes.NewBuffer(c.Expect)
		v, err := decodeHeaderValue(r)
		if err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}

		switch tv := v.(type) {
		case TimestampValue:
			exp := time.Time(c.Val.(TimestampValue))
			if e, a := exp, time.Time(tv); !e.Equal(a) {
				t.Errorf("%d, expect %v, got %v", i, e, a)
			}
		default:
			if e, a := c.Val, v; !reflect.DeepEqual(e, a) {
				t.Errorf("%d, expect %v, got %v", i, e, a)
			}
		}
	}
}

func TestValue_Decode(t *testing.T) {
	for i, c := range testValueEncodingCases {
		if c.Decode == nil {
			continue
		}

		r := bytes.NewBuffer(c.Expect)
		r.ReadByte() // strip off Type field

		v, err := c.Decode(r)
		if err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}

		switch tv := v.(type) {
		case TimestampValue:
			exp := time.Time(c.Val.(TimestampValue))
			if e, a := exp, time.Time(tv); !e.Equal(a) {
				t.Errorf("%d, expect %v, got %v", i, e, a)
			}
		default:
			if e, a := c.Val, v; !reflect.DeepEqual(e, a) {
				t.Errorf("%d, expect %v, got %v", i, e, a)
			}
		}
	}
}

func TestValue_String(t *testing.T) {
	cases := []struct {
		Val    Value
		Expect string
	}{
		{BoolValue(true), "true"},
		{BoolValue(false), "false"},
		{Int8Value(0x0f), "0x0f"},
		{Int16Value(0x0f), "0x000f"},
		{Int32Value(0x0f), "0x0000000f"},
		{Int64Value(0x0f), "0x000000000000000f"},
		{BytesValue([]byte{0, 1, 2, 3}), "AAECAw=="},
		{StringValue("abc123"), "abc123"},
		{TimestampValue(
			time.Date(2014, 04, 04, 0, 1, 0, 0, time.FixedZone("PDT", -7)),
		),
			"1396569667000",
		},
		{UUIDValue([16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}),
			"00010203-0405-0607-0809-0a0b0c0d0e0f",
		},
	}

	for i, c := range cases {
		if e, a := c.Expect, c.Val.String(); e != a {
			t.Errorf("%d, expect %v, got %v", i, e, a)
		}
	}
}

// safeErr formats an error for test output without dumping an oversized
// LengthError.Value payload (which can be tens of kilobytes). For a
// LengthError it prints only Part/Want/Have; for anything else it falls back
// to the normal error string.
func safeErr(err error) string {
	if err == nil {
		return "<nil>"
	}
	var le LengthError
	if errors.As(err, &le) {
		return fmt.Sprintf("LengthError{Part:%q Want:%d Have:%d}", le.Part, le.Want, le.Have)
	}
	return err.Error()
}

// TestStringValue_EncodeRejectsOversizedLength verifies that the string header
// encoder validates the encoded byte length before writing the 16-bit length
// field. It must return an error on overflow rather than truncating the length
// to its low 16 bits while still appending the full value (the header-smuggling
// primitive). This is the core protocol invariant for outbound event-stream
// headers.
func TestStringValue_EncodeRejectsOversizedLength(t *testing.T) {
	cases := map[string]struct {
		len       int
		expectErr bool
	}{
		"at Go limit": {
			len:       maxHeaderValueLen,
			expectErr: false,
		},
		"one over Go limit": {
			len:       maxHeaderValueLen + 1,
			expectErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			v := StringValue(strings.Repeat("a", c.len))

			var w bytes.Buffer
			err := v.encode(&w)

			if c.expectErr {
				if err == nil {
					t.Fatalf("expected error for value of length %d, got nil", c.len)
				}

				var lenErr LengthError
				if !errors.As(err, &lenErr) {
					t.Fatalf("expected LengthError, got %T: %s", err, safeErr(err))
				}
				if lenErr.Have != c.len {
					t.Errorf("expected Have=%d, got %d", c.len, lenErr.Have)
				}
				if lenErr.Want != maxHeaderValueLen {
					t.Errorf("expected Want=%d, got %d", maxHeaderValueLen, lenErr.Want)
				}

				// No partially serialized frame may be left on the wire
				// after validation fails.
				if w.Len() != 0 {
					t.Errorf("expected no bytes written on rejection, got %d bytes", w.Len())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error for value of length %d, got %s", c.len, safeErr(err))
			}
			// 1 byte type tag + 2 byte length + value bytes.
			if e, a := 1+2+c.len, w.Len(); e != a {
				t.Errorf("expected %d bytes written, got %d", e, a)
			}
		})
	}
}

// TestBytesValue_EncodeRejectsOversizedLength verifies the same bounds check
// for byte-array header values, the other variable-length value type that
// carries a 16-bit length prefix.
func TestBytesValue_EncodeRejectsOversizedLength(t *testing.T) {
	cases := map[string]struct {
		len       int
		expectErr bool
	}{
		"at Go limit": {
			len:       maxHeaderValueLen,
			expectErr: false,
		},
		"one over Go limit": {
			len:       maxHeaderValueLen + 1,
			expectErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			v := BytesValue(bytes.Repeat([]byte{0x61}, c.len))

			var w bytes.Buffer
			err := v.encode(&w)

			if c.expectErr {
				var lenErr LengthError
				if !errors.As(err, &lenErr) {
					t.Fatalf("expected LengthError, got %T: %s", err, safeErr(err))
				}
				if w.Len() != 0 {
					t.Errorf("expected no bytes written on rejection, got %d bytes", w.Len())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error for value of length %d, got %s", c.len, safeErr(err))
			}
			// 1 byte type tag + 2 byte length + value bytes.
			if e, a := 1+2+c.len, w.Len(); e != a {
				t.Errorf("expected %d bytes written, got %d", e, a)
			}
		})
	}
}

// TestStringValue_MultibyteByteLengthEnforced verifies that the length check is
// on ENCODED BYTE LENGTH, not rune/character count. A string whose rune count
// is within the limit but whose UTF-8 encoding exceeds it must be rejected.
//
// Uses the 3-byte rune U+20AC (euro sign, "€"). The rune count stays at or
// below the limit while the byte length is ~3x, so a correct byte-length check
// rejects it. This case PASSES against current behavior because the Go encoder
// correctly uses len(v) (bytes), not utf8.RuneCountInString.
func TestStringValue_MultibyteByteLengthEnforced(t *testing.T) {
	const euro = "\u20ac" // 3 UTF-8 bytes, 1 rune

	// Rune count = maxHeaderValueLen (within the per-character limit), but
	// byte length = 3 * maxHeaderValueLen (well over the limit).
	runeCount := maxHeaderValueLen
	v := StringValue(strings.Repeat(euro, runeCount))

	byteLen := len(string(v))
	if byteLen <= maxHeaderValueLen {
		t.Fatalf("test setup: expected byte length > %d, got %d", maxHeaderValueLen, byteLen)
	}

	var w bytes.Buffer
	err := v.encode(&w)

	var lenErr LengthError
	if !errors.As(err, &lenErr) {
		t.Fatalf("expected LengthError for %d-byte (%d-rune) value, got %T: %s",
			byteLen, runeCount, err, safeErr(err))
	}
	if lenErr.Have != byteLen {
		t.Errorf("expected LengthError.Have to be byte length %d, got %d", byteLen, lenErr.Have)
	}
	if w.Len() != 0 {
		t.Errorf("expected no bytes written on rejection, got %d bytes", w.Len())
	}
}

// TestFixedWidthValue_UUIDDecodeRejectsTruncated covers the invalid-width half
// of fixed-width value handling for UUID values. UUIDValue is structurally
// fixed on the ENCODE side in Go ([16]byte), so an encoder can never emit a
// wrong width, and the valid-width round-trip is already covered by
// TestValue_MarshalValue / TestHeader_DecodeValues. The gap is the DECODE
// side: an independent, spec-compliant decoder must reject a truncated
// fixed-width value rather than silently accepting it.
func TestFixedWidthValue_UUIDDecodeRejectsTruncated(t *testing.T) {
	// Type tag for UUID followed by only 8 value bytes (should be 16).
	truncated := append([]byte{byte(uuidValueType)}, make([]byte, 8)...)

	_, err := decodeHeaderValue(bytes.NewReader(truncated))
	if err == nil {
		t.Fatal("expected error decoding a 8-byte UUID value, got nil")
	}
	if !errors.Is(err, io.ErrUnexpectedEOF) && err != io.EOF {
		// decodeFixedBytesValue uses io.ReadFull, which returns
		// ErrUnexpectedEOF for a short read.
		t.Logf("decode error (acceptable as long as non-nil): %v", err)
	}
}
