package json

import (
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

func TestEncodeNestedShape(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			Nested: &nestedShape{
				Value: aws.String("expected value"),
			},
		},
	})
	if err != nil {
		t.Fatalf("expect no marshal error, %v", err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("expect no read error, %v", err)
	}

	expect := `{"nested":{"value":"expected value"}}`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeMapString(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			MapStr: map[string]*string{
				"abc": aws.String("123"),
			},
		},
	})
	if err != nil {
		t.Fatalf("expect no marshal error, %v", err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("expect no read error, %v", err)
	}

	expect := `{"mapstr":{"abc":"123"}}`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeMapShape(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			MapShape: map[string]*nestedShape{
				"abc": {Value: aws.String("1")},
				"123": {IntVal: aws.Int64(123)},
			},
		},
	})
	if err != nil {
		t.Fatalf("expect no marshal error, %v", err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("expect no read error, %v", err)
	}

	expect := `{"mapShape":{"abc":{"value":"1"},"123":{"intval":123}}}`

	awstesting.AssertJSON(t, expect, string(b), "expect bodies to match")
}
func TestEncodeListString(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			ListStr: []*string{
				aws.String("abc"),
				aws.String("123"),
			},
		},
	})
	if err != nil {
		t.Fatalf("expect no marshal error, %v", err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("expect no read error, %v", err)
	}

	expect := `{"liststr":["abc","123"]}`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeListFlatten(t *testing.T) {
	// TODO no JSON flatten
}
func TestEncodeListFlattened(t *testing.T) {
	// TODO No json flatten
}
func TestEncodeListNamed(t *testing.T) {
	// TODO no json named
}
func TestEncodeListShape(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			ListShape: []*nestedShape{
				{Value: aws.String("abc")},
				{Value: aws.String("123")},
				{IntVal: aws.Int64(123)},
			},
		},
	})
	if err != nil {
		t.Fatalf("expect no marshal error, %v", err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("expect no read error, %v", err)
	}

	expect := `{"listShape":[{"value":"abc"},{"value":"123"},{"intval":123}]}`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}

type baseShape struct {
	Payload *payloadShape
}

func (s *baseShape) MarshalFields(e protocol.FieldEncoder) error {
	if s.Payload != nil {
		e.SetFields(protocol.PayloadTarget, "payload", s.Payload, protocol.Metadata{})
	}
	return nil
}

type payloadShape struct {
	Value            *string
	IntVal           *int64
	TimeVal          *time.Time
	Nested           *nestedShape
	MapStr           map[string]*string
	MapFlatten       map[string]*string
	MapNamed         map[string]*string
	MapShape         map[string]*nestedShape
	MapFlattenShape  map[string]*nestedShape
	MapNamedShape    map[string]*nestedShape
	ListStr          []*string
	ListFlatten      []*string
	ListNamed        []*string
	ListShape        []*nestedShape
	ListFlattenShape []*nestedShape
	ListNamedShape   []*nestedShape
}

func (s *payloadShape) MarshalFields(e protocol.FieldEncoder) error {
	if s.Value != nil {
		e.SetValue(protocol.BodyTarget, "value", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*s.Value)}, protocol.Metadata{})
	}
	if s.IntVal != nil {
		e.SetValue(protocol.BodyTarget, "intval", protocol.Int64Value(*s.IntVal), protocol.Metadata{})
	}
	if s.TimeVal != nil {
		e.SetValue(protocol.BodyTarget, "timeval", protocol.TimeValue{
			V: *s.TimeVal, Format: protocol.UnixTimeFormat,
		}, protocol.Metadata{})
	}
	if s.Nested != nil {
		e.SetFields(protocol.BodyTarget, "nested", s.Nested, protocol.Metadata{})
	}
	if len(s.MapStr) > 0 {
		me := e.Map(protocol.BodyTarget, "mapstr", protocol.Metadata{})
		me.Start()
		for k, v := range s.MapStr {
			me.MapSetValue(k, protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*v)})
		}
		me.End()
	}
	if len(s.MapFlatten) > 0 {
		me := e.Map(protocol.BodyTarget, "mapFlatten", protocol.Metadata{
			Flatten: true,
		})
		me.Start()
		for k, v := range s.MapFlatten {
			me.MapSetValue(k, protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*v)})
		}
		me.End()
	}
	if len(s.MapNamed) > 0 {
		me := e.Map(protocol.BodyTarget, "mapNamed", protocol.Metadata{
			MapLocationNameKey: "namedKey", MapLocationNameValue: "namedValue",
		})
		me.Start()
		for k, v := range s.MapNamed {
			me.MapSetValue(k, protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*v)})
		}
		me.End()
	}
	if len(s.MapShape) > 0 {
		me := e.Map(protocol.BodyTarget, "mapShape", protocol.Metadata{})
		me.Start()
		for k, v := range s.MapShape {
			me.MapSetFields(k, v)
		}
		me.End()
	}
	if len(s.MapFlattenShape) > 0 {
		me := e.Map(protocol.BodyTarget, "mapFlattenShape", protocol.Metadata{
			Flatten: true,
		})
		me.Start()
		for k, v := range s.MapFlattenShape {
			me.MapSetFields(k, v)
		}
		me.End()
	}
	if len(s.MapNamedShape) > 0 {
		me := e.Map(protocol.BodyTarget, "mapNamedShape", protocol.Metadata{
			MapLocationNameKey: "namedKey", MapLocationNameValue: "namedValue",
		})
		me.Start()
		for k, v := range s.MapNamedShape {
			me.MapSetFields(k, v)
		}
		me.End()
	}
	if len(s.ListStr) > 0 {
		le := e.List(protocol.BodyTarget, "liststr", protocol.Metadata{})
		le.Start()
		for _, v := range s.ListStr {
			le.ListAddValue(protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*v)})
		}
		le.End()
	}
	if len(s.ListFlatten) > 0 {
		le := e.List(protocol.BodyTarget, "listFlatten", protocol.Metadata{
			Flatten: true,
		})
		le.Start()
		for _, v := range s.ListFlatten {
			le.ListAddValue(protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*v)})
		}
		le.End()
	}
	if len(s.ListNamed) > 0 {
		le := e.List(protocol.BodyTarget, "listNamed", protocol.Metadata{
			ListLocationName: "namedMember",
		})
		le.Start()
		for _, v := range s.ListNamed {
			le.ListAddValue(protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*v)})
		}
		le.End()
	}
	if len(s.ListShape) > 0 {
		le := e.List(protocol.BodyTarget, "listShape", protocol.Metadata{})
		le.Start()
		for _, v := range s.ListShape {
			le.ListAddFields(v)
		}
		le.End()
	}
	if len(s.ListFlattenShape) > 0 {
		le := e.List(protocol.BodyTarget, "listFlattenShape", protocol.Metadata{
			Flatten: true,
		})
		le.Start()
		for _, v := range s.ListFlattenShape {
			le.ListAddFields(v)
		}
		le.End()
	}
	if len(s.ListNamedShape) > 0 {
		le := e.List(protocol.BodyTarget, "listNamedShape", protocol.Metadata{
			ListLocationName: "namedMember",
		})
		le.Start()
		for _, v := range s.ListNamedShape {
			le.ListAddFields(v)
		}
		le.End()
	}
	return nil
}

type nestedShape struct {
	Value    *string
	IntVal   *int64
	Prefixed *string
}

func (s *nestedShape) MarshalFields(e protocol.FieldEncoder) error {
	if s.Value != nil {
		e.SetValue(protocol.BodyTarget, "value", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*s.Value)}, protocol.Metadata{})
	}
	if s.IntVal != nil {
		e.SetValue(protocol.BodyTarget, "intval", protocol.Int64Value(*s.IntVal), protocol.Metadata{})
	}
	if s.Prefixed != nil {
		e.SetValue(protocol.BodyTarget, "prefixed", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(*s.Prefixed)}, protocol.Metadata{})
	}
	return nil
}

func encode(s baseShape) (io.ReadSeeker, error) {
	e := NewEncoder()
	s.MarshalFields(e)
	return e.Encode()
}
