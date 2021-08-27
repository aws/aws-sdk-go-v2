package attributevalue

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// An Unmarshaler is an interface to provide custom unmarshaling of
// AttributeValues. Use this to provide custom logic determining
// how AttributeValues should be unmarshaled.
// 		type ExampleUnmarshaler struct {
// 			Value int
// 		}
//
// 		func (u *ExampleUnmarshaler) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
//          avN, ok := av.(*types.AttributeValueMemberN)
// 			if !ok {
// 				return nil
// 			}
//
// 			n, err := strconv.ParseInt(avN.Value, 10, 0)
// 			if err != nil {
// 				return err
// 			}
//
// 			u.Value = int(n)
// 			return nil
// 		}
type Unmarshaler interface {
	UnmarshalDynamoDBAttributeValue(types.AttributeValue) error
}

// Unmarshal will unmarshal AttributeValues to Go value types.
// Both generic interface{} and concrete types are valid unmarshal
// destination types.
//
// Unmarshal will allocate maps, slices, and pointers as needed to
// unmarshal the AttributeValue into the provided type value.
//
// When unmarshaling AttributeValues into structs Unmarshal matches
// the field names of the struct to the AttributeValue Map keys.
// Initially it will look for exact field name matching, but will
// fall back to case insensitive if not exact match is found.
//
// With the exception of omitempty, omitemptyelem, binaryset, numberset
// and stringset all struct tags used by Marshal are also used by
// Unmarshal.
//
// When decoding AttributeValues to interfaces Unmarshal will use the
// following types.
//
//		[]byte,                 AV Binary (B)
//		[][]byte,               AV Binary Set (BS)
//		bool,                   AV Boolean (BOOL)
//		[]interface{},          AV List (L)
//		map[string]interface{}, AV Map (M)
//		float64,                AV Number (N)
//		Number,                 AV Number (N) with UseNumber set
//		[]float64,              AV Number Set (NS)
//		[]Number,               AV Number Set (NS) with UseNumber set
//		string,                 AV String (S)
//		[]string,               AV String Set (SS)
//
// If the Decoder option, UseNumber is set numbers will be unmarshaled
// as Number values instead of float64. Use this to maintain the original
// string formating of the number as it was represented in the AttributeValue.
// In addition provides additional opportunities to parse the number
// string based on individual use cases.
//
// When unmarshaling any error that occurs will halt the unmarshal
// and return the error.
//
// The output value provided must be a non-nil pointer
func Unmarshal(av types.AttributeValue, out interface{}) error {
	return NewDecoder().Decode(av, out)
}

// UnmarshalMap is an alias for Unmarshal which unmarshals from
// a map of AttributeValues.
//
// The output value provided must be a non-nil pointer
func UnmarshalMap(m map[string]types.AttributeValue, out interface{}) error {
	return NewDecoder().Decode(&types.AttributeValueMemberM{Value: m}, out)
}

// UnmarshalList is an alias for Unmarshal func which unmarshals
// a slice of AttributeValues.
//
// The output value provided must be a non-nil pointer
func UnmarshalList(l []types.AttributeValue, out interface{}) error {
	return NewDecoder().Decode(&types.AttributeValueMemberL{Value: l}, out)
}

// UnmarshalListOfMaps is an alias for Unmarshal func which unmarshals a
// slice of maps of attribute values.
//
// This is useful for when you need to unmarshal the Items from a Query API
// call.
//
// The output value provided must be a non-nil pointer
func UnmarshalListOfMaps(l []map[string]types.AttributeValue, out interface{}) error {
	items := make([]types.AttributeValue, len(l))
	for i, m := range l {
		items[i] = &types.AttributeValueMemberM{Value: m}
	}

	return UnmarshalList(items, out)
}

// DecoderOptions is a collection of options to configure how the decoder
// unmarshalls the value.
type DecoderOptions struct {
	// Support other custom struct tag keys, such as `yaml`, `json`, or `toml`.
	// Note that values provided with a custom TagKey must also be supported
	// by the (un)marshalers in this package.
	//
	// Tag key `dynamodbav` will always be read, but if custom tag key
	// conflicts with `dynamodbav` the custom tag key value will be used.
	TagKey string

	// Instructs the decoder to decode AttributeValue Numbers as
	// Number type instead of float64 when the destination type
	// is interface{}. Similar to encoding/json.Number
	UseNumber bool
}

// A Decoder provides unmarshaling AttributeValues to Go value types.
type Decoder struct {
	options DecoderOptions
}

// NewDecoder creates a new Decoder with default configuration. Use
// the `opts` functional options to override the default configuration.
func NewDecoder(optFns ...func(*DecoderOptions)) *Decoder {
	var options DecoderOptions
	for _, fn := range optFns {
		fn(&options)
	}

	return &Decoder{
		options: options,
	}
}

// Decode will unmarshal an AttributeValue into a Go value type. An error
// will be return if the decoder is unable to unmarshal the AttributeValue
// to the provide Go value type.
//
// The output value provided must be a non-nil pointer
func (d *Decoder) Decode(av types.AttributeValue, out interface{}, opts ...func(*Decoder)) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() || !v.IsValid() {
		return &InvalidUnmarshalError{Type: reflect.TypeOf(out)}
	}

	return d.decode(av, v, tag{})
}

var stringInterfaceMapType = reflect.TypeOf(map[string]interface{}(nil))
var byteSliceType = reflect.TypeOf([]byte(nil))
var byteSliceSliceType = reflect.TypeOf([][]byte(nil))
var timeType = reflect.TypeOf(time.Time{})

func (d *Decoder) decode(av types.AttributeValue, v reflect.Value, fieldTag tag) error {
	var u Unmarshaler
	_, isNull := av.(*types.AttributeValueMemberNULL)
	if av == nil || isNull {
		u, v = indirect(v, true)
		if u != nil {
			return u.UnmarshalDynamoDBAttributeValue(av)
		}
		return d.decodeNull(v)
	}

	u, v = indirect(v, false)
	if u != nil {
		return u.UnmarshalDynamoDBAttributeValue(av)
	}

	switch tv := av.(type) {
	case *types.AttributeValueMemberB:
		return d.decodeBinary(tv.Value, v)

	case *types.AttributeValueMemberBOOL:
		return d.decodeBool(tv.Value, v)

	case *types.AttributeValueMemberBS:
		return d.decodeBinarySet(tv.Value, v)

	case *types.AttributeValueMemberL:
		return d.decodeList(tv.Value, v)

	case *types.AttributeValueMemberM:
		return d.decodeMap(tv.Value, v)

	case *types.AttributeValueMemberN:
		return d.decodeNumber(tv.Value, v, fieldTag)

	case *types.AttributeValueMemberNS:
		return d.decodeNumberSet(tv.Value, v)

	case *types.AttributeValueMemberS:
		return d.decodeString(tv.Value, v, fieldTag)

	case *types.AttributeValueMemberSS:
		return d.decodeStringSet(tv.Value, v)

	default:
		return fmt.Errorf("unsupported AttributeValue type, %T", av)
	}
}

func (d *Decoder) decodeBinary(b []byte, v reflect.Value) error {
	if v.Kind() == reflect.Interface {
		buf := make([]byte, len(b))
		copy(buf, b)
		v.Set(reflect.ValueOf(buf))
		return nil
	}

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return &UnmarshalTypeError{Value: "binary", Type: v.Type()}
	}

	if v.Type() == byteSliceType {
		// Optimization for []byte types
		if v.IsNil() || v.Cap() < len(b) {
			v.Set(reflect.MakeSlice(byteSliceType, len(b), len(b)))
		} else if v.Len() != len(b) {
			v.SetLen(len(b))
		}
		copy(v.Interface().([]byte), b)
		return nil
	}

	switch v.Type().Elem().Kind() {
	case reflect.Uint8:
		// Fallback to reflection copy for type aliased of []byte type
		if v.Kind() != reflect.Array && (v.IsNil() || v.Cap() < len(b)) {
			v.Set(reflect.MakeSlice(v.Type(), len(b), len(b)))
		} else if v.Len() != len(b) {
			v.SetLen(len(b))
		}
		for i := 0; i < len(b); i++ {
			v.Index(i).SetUint(uint64(b[i]))
		}
	default:
		if v.Kind() == reflect.Array && v.Type().Elem().Kind() == reflect.Uint8 {
			reflect.Copy(v, reflect.ValueOf(b))
			break
		}
		return &UnmarshalTypeError{Value: "binary", Type: v.Type()}
	}

	return nil
}

func (d *Decoder) decodeBool(b bool, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool, reflect.Interface:
		v.Set(reflect.ValueOf(b).Convert(v.Type()))

	default:
		return &UnmarshalTypeError{Value: "bool", Type: v.Type()}
	}

	return nil
}

func (d *Decoder) decodeBinarySet(bs [][]byte, v reflect.Value) error {
	var isArray bool

	switch v.Kind() {
	case reflect.Slice:
		// Make room for the slice elements if needed
		if v.IsNil() || v.Cap() < len(bs) {
			// What about if ignoring nil/empty values?
			v.Set(reflect.MakeSlice(v.Type(), 0, len(bs)))
		}
	case reflect.Array:
		// Limited to capacity of existing array.
		isArray = true
	case reflect.Interface:
		set := make([][]byte, len(bs))
		for i, b := range bs {
			if err := d.decodeBinary(b, reflect.ValueOf(&set[i]).Elem()); err != nil {
				return err
			}
		}
		v.Set(reflect.ValueOf(set))
		return nil
	default:
		return &UnmarshalTypeError{Value: "binary set", Type: v.Type()}
	}

	for i := 0; i < v.Cap() && i < len(bs); i++ {
		if !isArray {
			v.SetLen(i + 1)
		}
		u, elem := indirect(v.Index(i), false)
		if u != nil {
			return u.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberBS{Value: bs})
		}
		if err := d.decodeBinary(bs[i], elem); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) decodeNumber(n string, v reflect.Value, fieldTag tag) error {
	switch v.Kind() {
	case reflect.Interface:
		i, err := d.decodeNumberToInterface(n)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(i))
		return nil
	case reflect.String:
		if isNumberValueType(v) {
			v.SetString(n)
			return nil
		}
		v.SetString(n)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowInt(i) {
			return &UnmarshalTypeError{
				Value: fmt.Sprintf("number overflow, %s", n),
				Type:  v.Type(),
			}
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(n, 10, 64)
		if err != nil {
			return err
		}
		if v.OverflowUint(i) {
			return &UnmarshalTypeError{
				Value: fmt.Sprintf("number overflow, %s", n),
				Type:  v.Type(),
			}
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return err
		}
		if v.OverflowFloat(i) {
			return &UnmarshalTypeError{
				Value: fmt.Sprintf("number overflow, %s", n),
				Type:  v.Type(),
			}
		}
		v.SetFloat(i)
	default:
		if v.Type().ConvertibleTo(timeType) && fieldTag.AsUnixTime {
			t, err := decodeUnixTime(n)
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(t).Convert(v.Type()))
			return nil
		}
		return &UnmarshalTypeError{Value: "number", Type: v.Type()}
	}

	return nil
}

func (d *Decoder) decodeNumberToInterface(n string) (interface{}, error) {
	if d.options.UseNumber {
		return Number(n), nil
	}

	// Default to float64 for all numbers
	return strconv.ParseFloat(n, 64)
}

func (d *Decoder) decodeNumberSet(ns []string, v reflect.Value) error {
	var isArray bool

	switch v.Kind() {
	case reflect.Slice:
		// Make room for the slice elements if needed
		if v.IsNil() || v.Cap() < len(ns) {
			// What about if ignoring nil/empty values?
			v.Set(reflect.MakeSlice(v.Type(), 0, len(ns)))
		}
	case reflect.Array:
		// Limited to capacity of existing array.
		isArray = true
	case reflect.Interface:
		if d.options.UseNumber {
			set := make([]Number, len(ns))
			for i, n := range ns {
				if err := d.decodeNumber(n, reflect.ValueOf(&set[i]).Elem(), tag{}); err != nil {
					return err
				}
			}
			v.Set(reflect.ValueOf(set))
		} else {
			set := make([]float64, len(ns))
			for i, n := range ns {
				if err := d.decodeNumber(n, reflect.ValueOf(&set[i]).Elem(), tag{}); err != nil {
					return err
				}
			}
			v.Set(reflect.ValueOf(set))
		}
		return nil
	default:
		return &UnmarshalTypeError{Value: "number set", Type: v.Type()}
	}

	for i := 0; i < v.Cap() && i < len(ns); i++ {
		if !isArray {
			v.SetLen(i + 1)
		}
		u, elem := indirect(v.Index(i), false)
		if u != nil {
			return u.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberNS{Value: ns})
		}
		if err := d.decodeNumber(ns[i], elem, tag{}); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) decodeList(avList []types.AttributeValue, v reflect.Value) error {
	var isArray bool

	switch v.Kind() {
	case reflect.Slice:
		// Make room for the slice elements if needed
		if v.IsNil() || v.Cap() < len(avList) {
			// What about if ignoring nil/empty values?
			v.Set(reflect.MakeSlice(v.Type(), 0, len(avList)))
		}
	case reflect.Array:
		// Limited to capacity of existing array.
		isArray = true
	case reflect.Interface:
		s := make([]interface{}, len(avList))
		for i, av := range avList {
			if err := d.decode(av, reflect.ValueOf(&s[i]).Elem(), tag{}); err != nil {
				return err
			}
		}
		v.Set(reflect.ValueOf(s))
		return nil
	default:
		return &UnmarshalTypeError{Value: "list", Type: v.Type()}
	}

	// If v is not a slice, array
	for i := 0; i < v.Cap() && i < len(avList); i++ {
		if !isArray {
			v.SetLen(i + 1)
		}
		if err := d.decode(avList[i], v.Index(i), tag{}); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) decodeMap(avMap map[string]types.AttributeValue, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Map:
		t := v.Type()
		if t.Key().Kind() != reflect.String {
			return &UnmarshalTypeError{Value: "map string key", Type: t.Key()}
		}
		if v.IsNil() {
			v.Set(reflect.MakeMap(t))
		}
	case reflect.Struct:
	case reflect.Interface:
		v.Set(reflect.MakeMap(stringInterfaceMapType))
		v = v.Elem()
	default:
		return &UnmarshalTypeError{Value: "map", Type: v.Type()}
	}

	if v.Kind() == reflect.Map {
		for k, av := range avMap {
			key := reflect.New(v.Type().Key()).Elem()
			key.SetString(k)
			elem := reflect.New(v.Type().Elem()).Elem()
			if err := d.decode(av, elem, tag{}); err != nil {
				return err
			}
			v.SetMapIndex(key, elem)
		}
	} else if v.Kind() == reflect.Struct {
		fields := unionStructFields(v.Type(), structFieldOptions{
			TagKey: d.options.TagKey,
		})
		for k, av := range avMap {
			if f, ok := fields.FieldByName(k); ok {
				fv := decoderFieldByIndex(v, f.Index)
				if err := d.decode(av, fv, f.tag); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (d *Decoder) decodeNull(v reflect.Value) error {
	if v.IsValid() && v.CanSet() {
		v.Set(reflect.Zero(v.Type()))
	}

	return nil
}

func (d *Decoder) decodeString(s string, v reflect.Value, fieldTag tag) error {
	if fieldTag.AsString {
		return d.decodeNumber(s, v, fieldTag)
	}

	// To maintain backwards compatibility with ConvertFrom family of methods which
	// converted strings to time.Time structs
	if v.Type().ConvertibleTo(timeType) {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(t).Convert(v.Type()))
		return nil
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Interface:
		// Ensure type aliasing is handled properly
		v.Set(reflect.ValueOf(s).Convert(v.Type()))
	default:
		return &UnmarshalTypeError{Value: "string", Type: v.Type()}
	}

	return nil
}

func (d *Decoder) decodeStringSet(ss []string, v reflect.Value) error {
	var isArray bool

	switch v.Kind() {
	case reflect.Slice:
		// Make room for the slice elements if needed
		if v.IsNil() || v.Cap() < len(ss) {
			v.Set(reflect.MakeSlice(v.Type(), 0, len(ss)))
		}
	case reflect.Array:
		// Limited to capacity of existing array.
		isArray = true
	case reflect.Interface:
		set := make([]string, len(ss))
		for i, s := range ss {
			if err := d.decodeString(s, reflect.ValueOf(&set[i]).Elem(), tag{}); err != nil {
				return err
			}
		}
		v.Set(reflect.ValueOf(set))
		return nil
	default:
		return &UnmarshalTypeError{Value: "string set", Type: v.Type()}
	}

	for i := 0; i < v.Cap() && i < len(ss); i++ {
		if !isArray {
			v.SetLen(i + 1)
		}
		u, elem := indirect(v.Index(i), false)
		if u != nil {
			return u.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberSS{Value: ss})
		}
		if err := d.decodeString(ss[i], elem, tag{}); err != nil {
			return err
		}
	}

	return nil
}

func decodeUnixTime(n string) (time.Time, error) {
	v, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return time.Time{}, &UnmarshalError{
			Err: err, Value: n, Type: timeType,
		}
	}

	return time.Unix(v, 0), nil
}

// decoderFieldByIndex finds the field with the provided nested index, allocating
// embedded parent structs if needed
func decoderFieldByIndex(v reflect.Value, index []int) reflect.Value {
	for i, x := range index {
		if i > 0 && v.Kind() == reflect.Ptr && v.Type().Elem().Kind() == reflect.Struct {
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
		v = v.Field(x)
	}
	return v
}

// indirect will walk a value's interface or pointer value types. Returning
// the final value or the value a unmarshaler is defined on.
//
// Based on the enoding/json type reflect value type indirection in Go Stdlib
// https://golang.org/src/encoding/json/decode.go indirect func.
func indirect(v reflect.Value, decodingNull bool) (Unmarshaler, reflect.Value) {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}
		if v.Kind() != reflect.Ptr {
			break
		}
		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(Unmarshaler); ok {
				return u, reflect.Value{}
			}
		}
		v = v.Elem()
	}

	return nil, v
}

// A Number represents a Attributevalue number literal.
type Number string

// Float64 attempts to cast the number ot a float64, returning
// the result of the case or error if the case failed.
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 attempts to cast the number ot a int64, returning
// the result of the case or error if the case failed.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

// Uint64 attempts to cast the number ot a uint64, returning
// the result of the case or error if the case failed.
func (n Number) Uint64() (uint64, error) {
	return strconv.ParseUint(string(n), 10, 64)
}

// String returns the raw number represented as a string
func (n Number) String() string {
	return string(n)
}

// An UnmarshalTypeError is an error type representing a error
// unmarshaling the AttributeValue's element to a Go value type.
// Includes details about the AttributeValue type and Go value type.
type UnmarshalTypeError struct {
	Value string
	Type  reflect.Type
}

// Error returns the string representation of the error.
// satisfying the error interface
func (e *UnmarshalTypeError) Error() string {
	return fmt.Sprintf("unmarshal failed, cannot unmarshal %s into Go value type %s",
		e.Value, e.Type.String())
}

// An InvalidUnmarshalError is an error type representing an invalid type
// encountered while unmarshaling a AttributeValue to a Go value type.
type InvalidUnmarshalError struct {
	Type reflect.Type
}

// Error returns the string representation of the error.
// satisfying the error interface
func (e *InvalidUnmarshalError) Error() string {
	var msg string
	if e.Type == nil {
		msg = "cannot unmarshal to nil value"
	} else if e.Type.Kind() != reflect.Ptr {
		msg = fmt.Sprintf("cannot unmarshal to non-pointer value, got %s", e.Type.String())
	} else {
		msg = fmt.Sprintf("cannot unmarshal to nil value, %s", e.Type.String())
	}

	return fmt.Sprintf("unmarshal failed, %s", msg)
}

// An UnmarshalError wraps an error that occurred while unmarshaling a
// AttributeValue element into a Go type. This is different from
// UnmarshalTypeError in that it wraps the underlying error that occurred.
type UnmarshalError struct {
	Err   error
	Value string
	Type  reflect.Type
}

func (e *UnmarshalError) Unwrap() error {
	return e.Err
}

// Error returns the string representation of the error satisfying the error
// interface.
func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("unmarshal failed, cannot unmarshal %q into %s, %v",
		e.Value, e.Type.String(), e.Err)
}
