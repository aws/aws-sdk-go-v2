package converters

import (
	"strconv"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()

	if r == nil || r.converters == nil {
		t.Fail()
	}
}

func TestRegistry_Clone(t *testing.T) {
	r := DefaultRegistry.Clone()

	if r == nil || r.converters == nil {
		t.Fail()
	}

	if len(r.converters) != len(DefaultRegistry.converters) {
		t.Fail()
	}
}

func TestRegistry_Add(t *testing.T) {
	r := &Registry{}

	l := len(r.converters)

	r.Add("mock", &mockConverter{})

	if l+1 != len(r.converters) {
		t.Fail()
	}
}

func TestRegistry_Remove(t *testing.T) {
	l := len(DefaultRegistry.converters)

	DefaultRegistry.Remove("json")

	if l-1 != len(DefaultRegistry.converters) {
		t.Fail()
	}
}

func TestRegistry_Converter(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{
			name: "404",
			ok:   false,
		},
		{
			name: "json",
			ok:   true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_ = c
		})
	}
}

func TestRegistry_ConverterFor(t *testing.T) {}

//type x struct {
//	// raw
//	Uint      uint      `dynamodbav:",converter|uint"`
//	Uint8     uint8     `dynamodbav:",converter|uint8"`
//	Uint16    uint16    `dynamodbav:",converter|uint16"`
//	Uint32    uint32    `dynamodbav:",converter|uint32"`
//	Uint64    uint64    `dynamodbav:",converter|uint64"`
//	Int       int       `dynamodbav:",converter|int"`
//	Int8      int8      `dynamodbav:",converter|int8"`
//	Int16     int16     `dynamodbav:",converter|int16"`
//	Int32     int32     `dynamodbav:",converter|int32"`
//	Int64     int64     `dynamodbav:",converter|int64"`
//	Float32   float32   `dynamodbav:",converter|float32"`
//	Float64   float64   `dynamodbav:",converter|float64"`
//	Bool      bool      `dynamodbav:",converter|bool"`
//	ByteArray []byte    `dynamodbav:",converter|[]byte"`
//	String    string    `dynamodbav:",converter|string"`
//	TimeTime  time.Time `dynamodbav:",converter|time.Time"`
//	// ptr
//	PtrUint      *uint      `dynamodbav:",converter|uint"`
//	PtrUint8     *uint8     `dynamodbav:",converter|uint8"`
//	PtrUint16    *uint16    `dynamodbav:",converter|uint16"`
//	PtrUint32    *uint32    `dynamodbav:",converter|uint32"`
//	PtrUint64    *uint64    `dynamodbav:",converter|uint64"`
//	PtrInt       *int       `dynamodbav:",converter|int"`
//	PtrInt8      *int8      `dynamodbav:",converter|int8"`
//	PtrInt16     *int16     `dynamodbav:",converter|int16"`
//	PtrInt32     *int32     `dynamodbav:",converter|int32"`
//	PtrInt64     *int64     `dynamodbav:",converter|int64"`
//	PtrFloat32   *float32   `dynamodbav:",converter|float32"`
//	PtrFloat64   *float64   `dynamodbav:",converter|float64"`
//	PtrBool      *bool      `dynamodbav:",converter|bool"`
//	PtrByteArray *[]byte    `dynamodbav:",converter|[]byte"`
//	PtrString    *string    `dynamodbav:",converter|string"`
//	PtrTimeTime  *time.Time `dynamodbav:",converter|time.Time"`
//}
//
//func TestConverters(t *testing.T) {
//	jsn := `{
//	"Uint": 1,
//	"Uint8": 8,
//	"Uint16": 16,
//	"Uint32": 32,
//	"Uint64": 64,
//	"Int": -1,
//	"Int8": -8,
//	"Int16": -16,
//	"Int32": -32,
//	"Int64": -64,
//	"Float32": 32.32,
//	"Float64": 64.64,
//	"Bool": true,
//	"ByteArray": "dGVzdA==",
//	"String": "test",
//	"TimeTime": "2025-06-16T12:34:56Z",
//	"PtrUint": 1,
//	"PtrUint8": 8,
//	"PtrUint16": 16,
//	"PtrUint32": 32,
//	"PtrUint64": 64,
//	"PtrInt": -1,
//	"PtrInt8": -8,
//	"PtrInt16": -16,
//	"PtrInt32": -32,
//	"PtrInt64": -64,
//	"PtrFloat32": 32.32,
//	"PtrFloat64": 64.64,
//	"PtrBool": false,
//	"PtrByteArray": "dGVzdA==",
//	"PtrString": "test",
//	"PtrTimeTime": "2025-06-16T12:34:56Z"
//}`
//	_ = jsn
//
//	p := x{}
//	//err := json.Unmarshal([]byte(jsn), &p)
//	//if err != nil {
//	//	t.Fatalf("json.Unmarshal(): %v", err)
//	//}
//
//	e := json.NewEncoder(os.Stdout)
//	e.SetIndent("", "  ")
//	_ = e.Encode(p)
//
//	ref := reflect.ValueOf(p)
//	typ := ref.Type()
//
//	m := map[string]types.AttributeValue{}
//
//	for c := range typ.NumField() {
//		f := ref.Field(c)
//		t := typ.Field(c)
//		v := f.Interface()
//		var av types.AttributeValue
//		var err error
//
//		st := t.Tag.Get("dynamodbav")
//		stParts := strings.Split(st, ",")
//		cvtName := ""
//		for _, stPart := range stParts {
//			if !strings.HasPrefix(stPart, "converter|") {
//				continue
//			}
//
//			cvtName = stPart[10:]
//			break
//		}
//
//		fmt.Println("field:", t.Name)
//		fmt.Println("\tcvt:", cvtName)
//		fmt.Println("\tst:", st)
//		cvt := GetConverter(DefaultRegistry, cvtName)
//		if cvt == nil {
//			fmt.Printf("\tnil converter: %s\n", cvtName)
//			continue
//		}
//
//		av, err = cvt.ToAttributeValue(v, nil)
//		if err == nil {
//			m[t.Name] = av
//		} else {
//			fmt.Printf("\terr: %v\n", err)
//		}
//		fmt.Printf("\tav: %#+v\n", av)
//	}
//
//	//_ = e.Encode(m)
//	//spew.Dump(m)
//	//_ = e.Encode(p)
//}
