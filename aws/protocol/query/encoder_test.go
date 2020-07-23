package query

import (
	smithytesting "github.com/awslabs/smithy-go/testing"
	"testing"
)

func TestEncodeObject(t *testing.T) {
	encoder := NewEncoder()
	encoder.Object().Key("foo").String("bar")
	result, err := encoder.Encode()
	if err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`foo=bar`), result)
}

func TestEncodeNestedObject(t *testing.T) {
	encoder := NewEncoder()
	encoder.Object().Key("foo").Object().Key("bar").String("baz")
	result, err := encoder.Encode()
	if err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`foo.bar=baz`), result)
}

func TestEncodeList(t *testing.T) {
	encoder := NewEncoder()
	list := encoder.Object().Key("list").Array("spam")
	list.Value().String("spam")
	list.Value().String("eggs")
	result, err := encoder.Encode()
	if err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`list.spam.1=spam&list.spam.2=eggs`), result)
}

func TestEncodeFlatList(t *testing.T) {
	encoder := NewEncoder()
	list := encoder.Object().FlatKey("list").Array("spam")
	list.Value().String("spam")
	list.Value().String("eggs")
	result, err := encoder.Encode()
	if err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`list.1=spam&list.2=eggs`), result)
}

func TestEncodeMap(t *testing.T) {
	encoder := NewEncoder()
	mapValue := encoder.Object().Key("map").Map("key", "value")
	mapValue.Key("bar").String("baz")
	mapValue.Key("foo").String("bin")
	result, err := encoder.Encode()
	if err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`map.entry.1.key=bar&map.entry.1.value=baz&map.entry.2.key=foo&map.entry.2.value=bin`), result)
}

func TestEncodeFlatMap(t *testing.T) {
	encoder := NewEncoder()
	mapValue := encoder.Object().FlatKey("map").Map("key", "value")
	mapValue.Key("bar").String("baz")
	mapValue.Key("foo").String("bin")
	result, err := encoder.Encode()
	if err != nil {
		t.Fatal(err)
	}
	smithytesting.AssertURLFormEqual(t, []byte(`map.1.key=bar&map.1.value=baz&map.2.key=foo&map.2.value=bin`), result)
}
