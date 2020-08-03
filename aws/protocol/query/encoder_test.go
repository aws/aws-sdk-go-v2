package query

import (
	"bytes"
	"fmt"
	smithytesting "github.com/awslabs/smithy-go/testing"
	"testing"
)

func TestEncodeObject(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	encoder := NewEncoder(buff)
	encoder.Object().Key("foo").String("bar")
	if err := encoder.Encode(); err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`foo=bar`), buff.Bytes())
}

func TestEncodeNestedObject(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	encoder := NewEncoder(buff)
	encoder.Object().Key("foo").Object().Key("bar").String("baz")
	if err := encoder.Encode(); err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`foo.bar=baz`), buff.Bytes())
}

func TestEncodeList(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	encoder := NewEncoder(buff)
	list := encoder.Object().Key("list").Array("spam")
	list.Value().String("spam")
	list.Value().String("eggs")
	if err := encoder.Encode(); err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`list.spam.1=spam&list.spam.2=eggs`), buff.Bytes())
}

func TestEncodeFlatList(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	encoder := NewEncoder(buff)
	list := encoder.Object().FlatKey("list").Array("spam")
	list.Value().String("spam")
	list.Value().String("eggs")
	if err := encoder.Encode(); err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`list.1=spam&list.2=eggs`), buff.Bytes())
}

func TestEncodeMap(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	encoder := NewEncoder(buff)
	mapValue := encoder.Object().Key("map").Map("key", "value")
	mapValue.Key("bar").String("baz")
	mapValue.Key("foo").String("bin")
	if err := encoder.Encode(); err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`map.entry.1.key=bar&map.entry.1.value=baz&map.entry.2.key=foo&map.entry.2.value=bin`), buff.Bytes())
}

func TestEncodeFlatMap(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	encoder := NewEncoder(buff)
	mapValue := encoder.Object().FlatKey("map").Map("key", "value")
	mapValue.Key("bar").String("baz")
	mapValue.Key("foo").String("bin")
	if err := encoder.Encode(); err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`map.1.key=bar&map.1.value=baz&map.2.key=foo&map.2.value=bin`), buff.Bytes())
}

// limitedWriter exists to isolate WriteString to ensure that any writer
// can actually be used
type limitedWriter struct {
	writer *bytes.Buffer
}

func (lw limitedWriter) Write(v []byte) (int, error) {
	return lw.writer.Write(v)
}

func TestEncodeHandlesBareIoWriter(t *testing.T) {
	buff := limitedWriter{writer: bytes.NewBuffer(nil)}
	encoder := NewEncoder(buff)
	encoder.Object().Key("foo").String("bar")
	if err := encoder.Encode(); err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`foo=bar`), buff.writer.Bytes())
}

// stringWriter exists to ensure that WriteString is called when
// available.
type stringWriter struct {
	writer *bytes.Buffer
}

func (w stringWriter) Write(v []byte) (int, error)  {
	return 0, fmt.Errorf("the WriteString method should be used when available")
}

func (w stringWriter) WriteString(v string) (int, error) {
	return w.writer.WriteString(v)
}

func TestEncodeUsesWriteString(t *testing.T) {
	buff := stringWriter{writer: bytes.NewBuffer(nil)}
	encoder := NewEncoder(buff)
	encoder.Object().Key("foo").String("bar")
	if err := encoder.Encode(); err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`foo=bar`), buff.writer.Bytes())
}
