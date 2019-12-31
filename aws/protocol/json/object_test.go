package json

import (
	"bytes"
	"testing"
)

func TestObject(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	scatch := make([]byte, 64)

	object := newObject(buffer, &scatch)
	object.Key("foo").String("bar")
	object.Key("faz").String("baz")
	object.Close()

	e := []byte(`{"foo":"bar","faz":"baz"}`)
	if a := buffer.Bytes(); bytes.Compare(e, a) != 0 {
		t.Errorf("expected %+q, but got %+q", e, a)
	}
}
