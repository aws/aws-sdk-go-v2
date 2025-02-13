package eventstream

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
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
	for i := 0; i < 10_000; i++ {
		headers = append(headers, h)
	}

	err = encoder.Encode(io.Discard, Message{Headers: headers})
	if err != nil {
		t.Fatalf("Expected encoder being able to encode %d size, failed with %v", l, err)
	}
}
