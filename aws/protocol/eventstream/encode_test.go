package eventstream

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestEncoder_Encode(t *testing.T) {
	cases, err := readPositiveTests("testdata")
	if err != nil {
		t.Fatalf("failed to load positive tests, %v", err)
	}

	for _, c := range cases {
		var w bytes.Buffer
		encoder := NewEncoder()

		err = encoder.Encode(&w, c.Decoded.Message())
		if err != nil {
			t.Fatalf("%s, failed to encode message, %v", c.Name, err)
		}

		if e, a := c.Encoded, w.Bytes(); !reflect.DeepEqual(e, a) {
			t.Errorf("%s, expect:\n%v\nactual:\n%v\n", c.Name,
				hex.Dump(e), hex.Dump(a))
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	var w bytes.Buffer
	encoder := NewEncoder()
	msg := Message{
		Headers: Headers{
			{Name: "event-id", Value: Int16Value(123)},
		},
		Payload: []byte(`{"abc":123}`),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := encoder.Encode(&w, msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncoder_Limits(t *testing.T) {
	l := 25 * 1024 * 1024 // Previously we failed if message was set to >16 MB
	payload := make([]byte, l)
	encoder := NewEncoder()
	err := encoder.Encode(io.Discard, Message{Payload: payload})
	if err != nil {
		t.Fatalf("Expected encoder being able to encode %d size, failed with %v", l, err)
	}

	h := Header{
		Name: "event-id", Value: Int16Value(123),
	}

	headers := make(Headers, 0, 10_000) // Previously we failed if headers size was above a certain size
	for range 10_000 {
		headers = append(headers, h)
	}

	err = encoder.Encode(io.Discard, Message{Headers: headers})
	if err != nil {
		t.Fatalf("Expected encoder being able to encode %d size, failed with %v", l, err)
	}
}

// TestEncodeHeaders_NameLengthBoundary covers the 255-byte (valid) and
// 256-byte (invalid) encoded header-name boundary from the protocol spec:
// "Header names can contain no more than 255 encoded bytes." A 256-byte name
// must be rejected with a LengthError and leave nothing on the wire.
func TestEncodeHeaders_NameLengthBoundary(t *testing.T) {
	cases := map[string]struct {
		nameLen   int
		expectErr bool
	}{
		"255 bytes is valid": {
			nameLen:   maxHeaderNameLen, // 255
			expectErr: false,
		},
		"256 bytes is rejected": {
			nameLen:   maxHeaderNameLen + 1, // 256
			expectErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			hdrName := strings.Repeat("n", c.nameLen)
			headers := Headers{
				{Name: hdrName, Value: StringValue("v")},
			}

			var w bytes.Buffer
			err := EncodeHeaders(&w, headers)

			if c.expectErr {
				if err == nil {
					t.Fatalf("expected error for %d-byte name, got nil (encoded %d bytes)",
						c.nameLen, w.Len())
				}
				var lenErr LengthError
				if !errors.As(err, &lenErr) {
					t.Fatalf("expected LengthError, got %T: %s", err, safeErr(err))
				}
				if lenErr.Want != maxHeaderNameLen {
					t.Errorf("expected Want=%d, got %d", maxHeaderNameLen, lenErr.Want)
				}
				if w.Len() != 0 {
					t.Errorf("expected no bytes written on rejection, got %d", w.Len())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error for %d-byte name, got %s", c.nameLen, safeErr(err))
			}
			// 1 length byte + name bytes + value (type tag + 2 len + 1 "v").
			if gotLen := int(w.Bytes()[0]); gotLen != c.nameLen {
				t.Errorf("expected encoded name length %d, got %d", c.nameLen, gotLen)
			}
		})
	}
}

// TestEncode_ForgedHeaderCannotProduceTransmissibleFrame proves a forged header
// input (an oversized string value whose length would wrap to a small/zero
// 16-bit value) cannot produce a transmissible event-stream frame. The
// full-message Encoder.Encode path must surface the per-field length error and
// emit no frame; a decoder is then shown to have nothing valid to read.
func TestEncode_ForgedHeaderCannotProduceTransmissibleFrame(t *testing.T) {
	// 65,536 bytes: low 16 bits of the length are zero.
	forged := StringValue(strings.Repeat("a", 1<<16))

	msg := Message{
		Headers: Headers{
			{Name: "x", Value: forged},
		},
		Payload: []byte("payload"),
	}

	var w bytes.Buffer
	err := NewEncoder().Encode(&w, msg)
	if err == nil {
		t.Fatalf("expected Encode to reject forged header, got nil; wrote %d bytes", w.Len())
	}
	var lenErr LengthError
	if !errors.As(err, &lenErr) {
		t.Fatalf("expected LengthError, got %T: %s", err, safeErr(err))
	}

	// Nothing transmissible must have been written.
	if w.Len() != 0 {
		t.Fatalf("expected no frame bytes on the wire, got %d bytes", w.Len())
	}

	// And a decoder has nothing to decode.
	if _, derr := NewDecoder().Decode(bytes.NewReader(w.Bytes()), nil); derr == nil {
		t.Fatal("expected decode of empty output to fail, got nil")
	}
}

// TestEncode_FailureAtomicity proves that when an input is rejected, the
// encoder leaves no partial frame. Asserts atomicity at the full Encoder.Encode
// level (not just the per-value encoder): a message whose header value is
// oversized must result in zero bytes written.
func TestEncode_FailureAtomicity(t *testing.T) {
	msg := Message{
		Headers: Headers{
			{Name: "valid-header", Value: StringValue("ok")},
			{Name: "oversized", Value: StringValue(strings.Repeat("a", maxHeaderValueLen+1))},
		},
		Payload: []byte("payload-bytes"),
	}

	var w bytes.Buffer
	err := NewEncoder().Encode(&w, msg)
	if err == nil {
		t.Fatalf("expected Encode to fail on oversized header value, got nil")
	}
	if w.Len() != 0 {
		t.Errorf("expected zero bytes written (atomic failure), got %d bytes", w.Len())
	}
}
