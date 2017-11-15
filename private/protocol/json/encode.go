package json

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

// An Encoder provides encoding of the AWS JSON protocol. This encoder will will
// write all content to JSON. Only supports body and payload targets.
type Encoder struct {
	encoder
	root bool
}

// NewEncoder creates a new encoder for encoding AWS JSON protocol. Only encodes
// fields into the JSON body, and error is returned if target is anything other
// than Body or Payload.
func NewEncoder() *Encoder {
	buf := bytes.NewBuffer([]byte{'{'})
	s := newScope(buf, nil)
	e := &Encoder{
		encoder: encoder{
			buf:      buf,
			scope:    s,
			fieldBuf: &protocol.FieldBuffer{},
		},
		root: true,
	}

	return e
}

// Encode returns the encoded XMl reader. An error will be returned if one was
// encountered while building the JSON body.
func (e *Encoder) Encode() (io.ReadSeeker, error) {
	b, err := e.encode()
	if err != nil {
		return nil, err
	}

	if len(b) == 2 {
		// Account for first starting object in buffer
		return nil, nil
	}

	return bytes.NewReader(b), nil
}

// SetValue sets an individual value to the JSON body.
func (e *Encoder) SetValue(t protocol.Target, k string, v protocol.ValueMarshaler, meta protocol.Metadata) {
	e.encoder.scope.writeSep()
	e.writeKey(k)
	e.writeValue(v)
}

// SetStream is not supported for JSON protocol marshaling.
func (e *Encoder) SetStream(t protocol.Target, k string, v protocol.StreamMarshaler, meta protocol.Metadata) {
	if e.err != nil {
		return
	}
	e.err = fmt.Errorf("json encoder SetStream not supported, %s, %s", t, k)
}

// Map will return a new mapEncoder and create a new scope for the map encoding.
func (e *Encoder) Map(t protocol.Target, k string, meta protocol.Metadata) protocol.MapEncoder {
	temp := e.encoder
	temp.scope = newScope(e.encoder.buf, temp.scope)
	return &mapEncoder{temp, k}
}

// List will return a new listEncoder and create a new scope for the list encoding.
func (e *Encoder) List(t protocol.Target, k string, meta protocol.Metadata) protocol.ListEncoder {
	temp := e.encoder
	temp.scope = newScope(e.encoder.buf, temp.scope)
	return &listEncoder{temp, k}
}

// SetFields sets the nested fields to the JSON body.
func (e *Encoder) SetFields(t protocol.Target, k string, m protocol.FieldMarshaler, meta protocol.Metadata) {
	if t == protocol.PayloadTarget {
		// Ignore payload key and only marshal body without wrapping in object first.
		nested := Encoder{
			encoder: encoder{
				buf:      e.encoder.buf,
				fieldBuf: e.encoder.fieldBuf,
				scope:    newScope(e.encoder.buf, e.encoder.scope),
			},
		}
		m.MarshalFields(&nested)
		e.err = nested.err
		return
	}

	e.scope.writeSep()
	e.writeKey(k)
	e.writeObject(func(enc encoder) error {
		enc.scope = newScope(enc.buf, enc.scope)
		nested := Encoder{encoder: enc}
		m.MarshalFields(&nested)
		return nested.err
	})
}

// A listEncoder encodes elements within a list for the JSON encoder.
type listEncoder struct {
	encoder
	k string
}

// Map return a new mapEncoder while creating a new scope for the encoder.
func (e *listEncoder) Map() protocol.MapEncoder {
	temp := e.encoder
	temp.scope = newScope(e.buf, temp.scope)
	return &mapEncoder{temp, ""}
}

// List return a new listEncoder while creating a new scope for the encoder.
func (e *listEncoder) List() protocol.ListEncoder {
	temp := e.encoder
	temp.scope = newScope(e.buf, temp.scope)
	return &listEncoder{temp, ""}
}

// Start will open a new scope for a list and write the given key.
func (e *listEncoder) Start() {
	e.encoder.scope.parent.writeSep()
	e.writeKey(e.k)
	e.WriteListStart()
}

// End will close the list.
func (e *listEncoder) End() {
	e.WriteListEnd()
}

// ListAddValue will add the value to the list.
func (e *listEncoder) ListAddValue(v protocol.ValueMarshaler) {
	e.encoder.scope.writeSep()
	e.writeValue(v)
}

// ListAddFields will set the nested type's fields to the list.
func (e *listEncoder) ListAddFields(m protocol.FieldMarshaler) {
	e.encoder.scope.writeSep()
	e.writeObject(func(enc encoder) error {
		enc.scope = newScope(enc.buf, enc.scope)
		nested := Encoder{encoder: enc}
		m.MarshalFields(&nested)
		return nested.err
	})
}

// A mapEncoder encodes key values pair map values for the JSON encoder.
type mapEncoder struct {
	encoder encoder
	k       string
}

// Start will open a new scope for a list and write the given key.
func (e *mapEncoder) Start() {
	e.encoder.scope.parent.writeSep()
	e.encoder.writeKey(e.k)
	e.encoder.WriteMapStart()
}

// End will close the list.
func (e *mapEncoder) End() {
	e.encoder.WriteMapEnd()
}

// Map will create a new scope and return a mapEncoder.
func (e *mapEncoder) Map(k string) protocol.MapEncoder {
	temp := e.encoder
	temp.scope = newScope(e.encoder.buf, temp.scope)
	return &mapEncoder{temp, k}
}

// List will create a new scope and return a listEncoder
func (e *mapEncoder) List(k string) protocol.ListEncoder {
	temp := e.encoder
	temp.scope = newScope(e.encoder.buf, temp.scope)
	return &listEncoder{temp, k}
}

// MapSetValue sets a map value.
func (e *mapEncoder) MapSetValue(k string, v protocol.ValueMarshaler) {
	e.encoder.scope.writeSep()
	e.encoder.writeKey(k)
	e.encoder.writeValue(v)
}

// MapSetFields will set the nested type's fields under the map.
func (e *mapEncoder) MapSetFields(k string, m protocol.FieldMarshaler) {
	e.encoder.scope.writeSep()
	e.encoder.writeKey(k)
	e.encoder.writeObject(func(enc encoder) error {
		enc.scope = newScope(enc.buf, enc.scope)
		nested := Encoder{encoder: enc}
		m.MarshalFields(&nested)
		return nested.err
	})
}

type encoder struct {
	buf      *bytes.Buffer
	fieldBuf *protocol.FieldBuffer
	scope    *scope
	err      error
}

func (e encoder) encode() ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}

	// Close the root object
	e.buf.WriteByte('}')

	return e.buf.Bytes(), nil
}

func (e *encoder) writeKey(k string) {
	e.buf.WriteByte('"')
	e.buf.WriteString(k) // TODO escape?
	e.buf.WriteByte('"')
	e.buf.WriteByte(':')
}

func (e *encoder) writeValue(v protocol.ValueMarshaler) {
	if e.err != nil {
		return
	}

	b, err := e.fieldBuf.GetValue(v)
	if err != nil {
		e.err = err
		return
	}

	var asStr bool
	switch v.(type) {
	case protocol.QuotedValue:
		asStr = true
	}

	if asStr {
		escapeStringBytes(e.buf, b)
	} else {
		e.buf.Write(b)
	}
}

func (e *encoder) writeObject(fn func(encoder) error) {
	if e.err != nil {
		return
	}

	e.buf.WriteByte('{')
	e.err = fn(*e)
	e.buf.WriteByte('}')
}

func (e *encoder) WriteListStart() {
	e.buf.WriteByte('[')
}

func (e *encoder) WriteListEnd() {
	e.buf.WriteByte(']')
}

func (e *encoder) WriteMapStart() {
	e.buf.WriteByte('{')
}

func (e *encoder) WriteMapEnd() {
	e.buf.WriteByte('}')
}

type scope struct {
	started bool
	buf     *bytes.Buffer
	parent  *scope
}

func newScope(buf *bytes.Buffer, parent *scope) *scope {
	return &scope{false, buf, parent}
}

func (s *scope) writeSep() {
	if s.started {
		s.buf.WriteByte(',')
	} else {
		s.started = true
	}

}
