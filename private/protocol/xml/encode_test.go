package xml

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

func TestEncodeAttribute(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			AttrValue: aws.String("value"),
		},
	})
	if err != nil {
		t.Fatalf("expect no marshal error, %v", err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("expect no read error, %v", err)
	}

	expect := `<payload attrkey="value"></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}

func TestEncodeNamespace(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			Namespace: &nestedShape{
				Prefixed: aws.String("abc"),
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

	expect := `<payload><namespace xmlns:prefix="https://example.com"><prefixed>abc</prefixed></namespace></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}

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

	expect := `<payload><nested><value>expected value</value></nested></payload>`

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

	expect := `<payload><mapstr><entry><key>abc</key><value>123</value></entry></mapstr></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeMapFlatten(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			MapFlatten: map[string]*string{
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

	expect := `<payload><mapFlatten><key>abc</key><value>123</value></mapFlatten></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeMapNamed(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			MapNamed: map[string]*string{
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

	expect := `<payload><mapNamed><entry><namedKey>abc</namedKey><namedValue>123</namedValue></entry></mapNamed></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeMapShape(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			MapShape: map[string]*nestedShape{
				"abc": {Value: aws.String("1")},
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

	expect := `<payload><mapShape><entry><key>abc</key><value><value>1</value></value></entry></mapShape></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeMapFlattenShape(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			MapFlattenShape: map[string]*nestedShape{
				"abc": {Value: aws.String("1")},
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

	expect := `<payload><mapFlattenShape><key>abc</key><value><value>1</value></value></mapFlattenShape></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeMapNamedShape(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			MapNamedShape: map[string]*nestedShape{
				"abc": {Value: aws.String("1")},
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

	expect := `<payload><mapNamedShape><entry><namedKey>abc</namedKey><namedValue><value>1</value></namedValue></entry></mapNamedShape></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
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

	expect := `<payload><liststr><member>abc</member><member>123</member></liststr></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeListFlatten(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			ListFlatten: []*string{
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

	expect := `<payload><listFlatten>abc</listFlatten><listFlatten>123</listFlatten></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeListFlattened(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			ListFlatten: []*string{
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

	expect := `<payload><listFlatten>abc</listFlatten><listFlatten>123</listFlatten></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeListNamed(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			ListNamed: []*string{
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

	expect := `<payload><listNamed><namedMember>abc</namedMember><namedMember>123</namedMember></listNamed></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeListShape(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			ListShape: []*nestedShape{
				{Value: aws.String("abc")},
				{Value: aws.String("123")},
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

	expect := `<payload><listShape><member><value>abc</value></member><member><value>123</value></member></listShape></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeListFlattenShape(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			ListFlattenShape: []*nestedShape{
				{Value: aws.String("abc")},
				{Value: aws.String("123")},
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

	expect := `<payload><listFlattenShape><value>abc</value></listFlattenShape><listFlattenShape><value>123</value></listFlattenShape></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}
func TestEncodeListNamedShape(t *testing.T) {
	r, err := encode(baseShape{
		Payload: &payloadShape{
			ListNamedShape: []*nestedShape{
				{Value: aws.String("abc")},
				{Value: aws.String("123")},
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

	expect := `<payload><listNamedShape><namedMember><value>abc</value></namedMember><namedMember><value>123</value></namedMember></listNamedShape></payload>`

	if e, a := expect, string(b); e != a {
		t.Errorf("expect bodies to match, did not.\n,\tExpect:\n%s\n\tActual:\n%s\n", e, a)
	}
}

type baseShape struct {
	Payload *payloadShape
}

func (s *baseShape) MarshalFields(e protocol.FieldEncoder) error {
	if s.Payload != nil {
		attrs := make([]protocol.Attribute, 0, 1)
		if s.Payload.AttrValue != nil {
			attrs = append(attrs, protocol.Attribute{
				Name:  "attrkey",
				Value: protocol.StringValue(*s.Payload.AttrValue),
				Meta:  protocol.Metadata{},
			})
		}
		e.SetFields(protocol.PayloadTarget, "payload", s.Payload, protocol.Metadata{Attributes: attrs})
	}
	return nil
}

type payloadShape struct {
	AttrValue        *string
	Namespace        *nestedShape
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
	// Attribute values are skipped
	if s.Namespace != nil {
		e.SetFields(protocol.BodyTarget, "namespace", s.Namespace, protocol.Metadata{
			XMLNamespaceURI: "https://example.com", XMLNamespacePrefix: "prefix",
		})
	}
	if s.Nested != nil {
		e.SetFields(protocol.BodyTarget, "nested", s.Nested, protocol.Metadata{})
	}
	if len(s.MapStr) > 0 {
		me := e.Map(protocol.BodyTarget, "mapstr", protocol.Metadata{})
		me.Start()
		for k, v := range s.MapStr {
			me.MapSetValue(k, protocol.StringValue(*v))
		}
		me.End()
	}
	if len(s.MapFlatten) > 0 {
		me := e.Map(protocol.BodyTarget, "mapFlatten", protocol.Metadata{
			Flatten: true,
		})
		me.Start()
		for k, v := range s.MapFlatten {
			me.MapSetValue(k, protocol.StringValue(*v))
		}
		me.End()
	}
	if len(s.MapNamed) > 0 {
		me := e.Map(protocol.BodyTarget, "mapNamed", protocol.Metadata{
			MapLocationNameKey: "namedKey", MapLocationNameValue: "namedValue",
		})
		me.Start()
		for k, v := range s.MapNamed {
			me.MapSetValue(k, protocol.StringValue(*v))
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
			le.ListAddValue(protocol.StringValue(*v))
		}

		le.End()
	}
	if len(s.ListFlatten) > 0 {
		le := e.List(protocol.BodyTarget, "listFlatten", protocol.Metadata{
			Flatten: true,
		})
		le.Start()
		for _, v := range s.ListFlatten {
			le.ListAddValue(protocol.StringValue(*v))
		}

		le.End()
	}
	if len(s.ListNamed) > 0 {
		le := e.List(protocol.BodyTarget, "listNamed", protocol.Metadata{
			ListLocationName: "namedMember",
		})
		le.Start()
		for _, v := range s.ListNamed {
			le.ListAddValue(protocol.StringValue(*v))
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
	Prefixed *string
}

func (s *nestedShape) MarshalFields(e protocol.FieldEncoder) error {
	if s.Value != nil {
		e.SetValue(protocol.BodyTarget, "value", protocol.StringValue(*s.Value), protocol.Metadata{})
	}
	if s.Prefixed != nil {
		e.SetValue(protocol.BodyTarget, "prefixed", protocol.StringValue(*s.Prefixed), protocol.Metadata{
			XMLNamespacePrefix: "prefix",
		})
	}
	return nil
}
func encodeNestedShapeMap(vs map[string]*nestedShape) func(protocol.MapEncoder) {
	return func(me protocol.MapEncoder) {
		for k, v := range vs {
			me.MapSetFields(k, v)
		}
	}
}
func encodeNestedShapeList(vs []*nestedShape) func(protocol.ListEncoder) {
	return func(le protocol.ListEncoder) {
		for _, v := range vs {
			le.ListAddFields(v)
		}
	}
}

func encode(s baseShape) (io.ReadSeeker, error) {
	e := NewEncoder()
	s.MarshalFields(e)
	return e.Encode()
}
