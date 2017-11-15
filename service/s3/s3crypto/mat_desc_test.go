package s3crypto

import (
	"reflect"
	"testing"
)

func TestEncodeMaterialDescription(t *testing.T) {
	md := MaterialDescription{}
	md["foo"] = "bar"
	b, err := md.encodeDescription()
	expected := `{"foo":"bar"}`
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if expected != string(b) {
		t.Errorf("expected %s, but received %s", expected, string(b))
	}
}
func TestDecodeMaterialDescription(t *testing.T) {
	md := MaterialDescription{}
	json := `{"foo":"bar"}`
	err := md.decodeDescription([]byte(json))
	expected := MaterialDescription{
		"foo": "bar",
	}
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if !reflect.DeepEqual(expected, md) {
		t.Error("expected material description to be equivalent, but received otherwise")
	}
}
