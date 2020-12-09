package attributevalue

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
	"github.com/google/go-cmp/cmp"
)

type testEmptyCollectionsNumericalScalars struct {
	String string

	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64

	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64

	Float32 float32
	Float64 float64
}

type testEmptyCollectionsOmittedNumericalScalars struct {
	String string `dynamodbav:",omitempty"`

	Uint8  uint8  `dynamodbav:",omitempty"`
	Uint16 uint16 `dynamodbav:",omitempty"`
	Uint32 uint32 `dynamodbav:",omitempty"`
	Uint64 uint64 `dynamodbav:",omitempty"`

	Int8  int8  `dynamodbav:",omitempty"`
	Int16 int16 `dynamodbav:",omitempty"`
	Int32 int32 `dynamodbav:",omitempty"`
	Int64 int64 `dynamodbav:",omitempty"`

	Float32 float32 `dynamodbav:",omitempty"`
	Float64 float64 `dynamodbav:",omitempty"`
}

type testEmptyCollectionsNulledNumericalScalars struct {
	String string `dynamodbav:",nullempty"`

	Uint8  uint8  `dynamodbav:",nullempty"`
	Uint16 uint16 `dynamodbav:",nullempty"`
	Uint32 uint32 `dynamodbav:",nullempty"`
	Uint64 uint64 `dynamodbav:",nullempty"`

	Int8  int8  `dynamodbav:",nullempty"`
	Int16 int16 `dynamodbav:",nullempty"`
	Int32 int32 `dynamodbav:",nullempty"`
	Int64 int64 `dynamodbav:",nullempty"`

	Float32 float32 `dynamodbav:",nullempty"`
	Float64 float64 `dynamodbav:",nullempty"`
}

type testEmptyCollectionsPtrScalars struct {
	PtrString *string

	PtrUint8  *uint8
	PtrUint16 *uint16
	PtrUint32 *uint32
	PtrUint64 *uint64

	PtrInt8  *int8
	PtrInt16 *int16
	PtrInt32 *int32
	PtrInt64 *int64

	PtrFloat32 *float32
	PtrFloat64 *float64
}

type testEmptyCollectionsOmittedPtrNumericalScalars struct {
	PtrString *string `dynamodbav:",omitempty"`

	PtrUint8  *uint8  `dynamodbav:",omitempty"`
	PtrUint16 *uint16 `dynamodbav:",omitempty"`
	PtrUint32 *uint32 `dynamodbav:",omitempty"`
	PtrUint64 *uint64 `dynamodbav:",omitempty"`

	PtrInt8  *int8  `dynamodbav:",omitempty"`
	PtrInt16 *int16 `dynamodbav:",omitempty"`
	PtrInt32 *int32 `dynamodbav:",omitempty"`
	PtrInt64 *int64 `dynamodbav:",omitempty"`

	PtrFloat32 *float32 `dynamodbav:",omitempty"`
	PtrFloat64 *float64 `dynamodbav:",omitempty"`
}

type testEmptyCollectionsNulledPtrNumericalScalars struct {
	PtrString *string `dynamodbav:",nullempty"`

	PtrUint8  *uint8  `dynamodbav:",nullempty"`
	PtrUint16 *uint16 `dynamodbav:",nullempty"`
	PtrUint32 *uint32 `dynamodbav:",nullempty"`
	PtrUint64 *uint64 `dynamodbav:",nullempty"`

	PtrInt8  *int8  `dynamodbav:",nullempty"`
	PtrInt16 *int16 `dynamodbav:",nullempty"`
	PtrInt32 *int32 `dynamodbav:",nullempty"`
	PtrInt64 *int64 `dynamodbav:",nullempty"`

	PtrFloat32 *float32 `dynamodbav:",nullempty"`
	PtrFloat64 *float64 `dynamodbav:",nullempty"`
}

type testEmptyCollectionTypes struct {
	Map       map[string]string
	Slice     []string
	ByteSlice []byte
	ByteArray [4]byte
	ZeroArray [0]byte
	BinarySet [][]byte `dynamodbav:",binaryset"`
	NumberSet []int    `dynamodbav:",numberset"`
	StringSet []string `dynamodbav:",stringset"`
}

type testEmptyCollectionTypesOmitted struct {
	Map       map[string]string `dynamodbav:",omitempty"`
	Slice     []string          `dynamodbav:",omitempty"`
	ByteSlice []byte            `dynamodbav:",omitempty"`
	ByteArray [4]byte           `dynamodbav:",omitempty"`
	ZeroArray [0]byte           `dynamodbav:",omitempty"`
	BinarySet [][]byte          `dynamodbav:",binaryset,omitempty"`
	NumberSet []int             `dynamodbav:",numberset,omitempty"`
	StringSet []string          `dynamodbav:",stringset,omitempty"`
}

type testEmptyCollectionTypesNulled struct {
	Map       map[string]string `dynamodbav:",nullempty"`
	Slice     []string          `dynamodbav:",nullempty"`
	ByteSlice []byte            `dynamodbav:",nullempty"`
	ByteArray [4]byte           `dynamodbav:",nullempty"`
	ZeroArray [0]byte           `dynamodbav:",nullempty"`
	BinarySet [][]byte          `dynamodbav:",binaryset,nullempty"`
	NumberSet []int             `dynamodbav:",numberset,nullempty"`
	StringSet []string          `dynamodbav:",stringset,nullempty"`
}

type testEmptyCollectionStruct struct {
	Int int
}

type testEmptyCollectionStructOmitted struct {
	Slice []string `dynamodbav:",omitempty"`
}

var sharedEmptyCollectionsTestCases = map[string]struct {
	in types.AttributeValue
	// alternative input to compare against for marshal flow
	inMarshal types.AttributeValue

	actual, expected interface{}
	err              error
}{
	"scalars with zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"String":  &types.AttributeValueMemberS{Value: ""},
				"Uint8":   &types.AttributeValueMemberN{Value: "0"},
				"Uint16":  &types.AttributeValueMemberN{Value: "0"},
				"Uint32":  &types.AttributeValueMemberN{Value: "0"},
				"Uint64":  &types.AttributeValueMemberN{Value: "0"},
				"Int8":    &types.AttributeValueMemberN{Value: "0"},
				"Int16":   &types.AttributeValueMemberN{Value: "0"},
				"Int32":   &types.AttributeValueMemberN{Value: "0"},
				"Int64":   &types.AttributeValueMemberN{Value: "0"},
				"Float32": &types.AttributeValueMemberN{Value: "0"},
				"Float64": &types.AttributeValueMemberN{Value: "0"},
			},
		},
		actual:   &testEmptyCollectionsNumericalScalars{},
		expected: testEmptyCollectionsNumericalScalars{},
	},
	"scalars with non-zero values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"String":  &types.AttributeValueMemberS{Value: "test string"},
				"Uint8":   &types.AttributeValueMemberN{Value: "1"},
				"Uint16":  &types.AttributeValueMemberN{Value: "2"},
				"Uint32":  &types.AttributeValueMemberN{Value: "3"},
				"Uint64":  &types.AttributeValueMemberN{Value: "4"},
				"Int8":    &types.AttributeValueMemberN{Value: "-5"},
				"Int16":   &types.AttributeValueMemberN{Value: "-6"},
				"Int32":   &types.AttributeValueMemberN{Value: "-7"},
				"Int64":   &types.AttributeValueMemberN{Value: "-8"},
				"Float32": &types.AttributeValueMemberN{Value: "9.9"},
				"Float64": &types.AttributeValueMemberN{Value: "10.1"},
			},
		},
		actual: &testEmptyCollectionsNumericalScalars{},
		expected: testEmptyCollectionsNumericalScalars{
			String:  "test string",
			Uint8:   1,
			Uint16:  2,
			Uint32:  3,
			Uint64:  4,
			Int8:    -5,
			Int16:   -6,
			Int32:   -7,
			Int64:   -8,
			Float32: 9.9,
			Float64: 10.1,
		},
	},
	"omittable scalars with zero value": {
		in:       &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
		actual:   &testEmptyCollectionsOmittedNumericalScalars{},
		expected: testEmptyCollectionsOmittedNumericalScalars{},
	},
	"omittable scalars with non-zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"String":  &types.AttributeValueMemberS{Value: "test string"},
				"Uint8":   &types.AttributeValueMemberN{Value: "1"},
				"Uint16":  &types.AttributeValueMemberN{Value: "2"},
				"Uint32":  &types.AttributeValueMemberN{Value: "3"},
				"Uint64":  &types.AttributeValueMemberN{Value: "4"},
				"Int8":    &types.AttributeValueMemberN{Value: "-5"},
				"Int16":   &types.AttributeValueMemberN{Value: "-6"},
				"Int32":   &types.AttributeValueMemberN{Value: "-7"},
				"Int64":   &types.AttributeValueMemberN{Value: "-8"},
				"Float32": &types.AttributeValueMemberN{Value: "9.9"},
				"Float64": &types.AttributeValueMemberN{Value: "10.1"},
			},
		},
		actual: &testEmptyCollectionsOmittedNumericalScalars{},
		expected: testEmptyCollectionsOmittedNumericalScalars{
			String:  "test string",
			Uint8:   1,
			Uint16:  2,
			Uint32:  3,
			Uint64:  4,
			Int8:    -5,
			Int16:   -6,
			Int32:   -7,
			Int64:   -8,
			Float32: 9.9,
			Float64: 10.1,
		},
	},
	"null scalars with zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"String":  &types.AttributeValueMemberNULL{Value: true},
				"Uint8":   &types.AttributeValueMemberNULL{Value: true},
				"Uint16":  &types.AttributeValueMemberNULL{Value: true},
				"Uint32":  &types.AttributeValueMemberNULL{Value: true},
				"Uint64":  &types.AttributeValueMemberNULL{Value: true},
				"Int8":    &types.AttributeValueMemberNULL{Value: true},
				"Int16":   &types.AttributeValueMemberNULL{Value: true},
				"Int32":   &types.AttributeValueMemberNULL{Value: true},
				"Int64":   &types.AttributeValueMemberNULL{Value: true},
				"Float32": &types.AttributeValueMemberNULL{Value: true},
				"Float64": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual:   &testEmptyCollectionsNulledNumericalScalars{},
		expected: testEmptyCollectionsNulledNumericalScalars{},
	},
	"null scalars with non-zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"String":  &types.AttributeValueMemberS{Value: "test string"},
				"Uint8":   &types.AttributeValueMemberN{Value: "1"},
				"Uint16":  &types.AttributeValueMemberN{Value: "2"},
				"Uint32":  &types.AttributeValueMemberN{Value: "3"},
				"Uint64":  &types.AttributeValueMemberN{Value: "4"},
				"Int8":    &types.AttributeValueMemberN{Value: "-5"},
				"Int16":   &types.AttributeValueMemberN{Value: "-6"},
				"Int32":   &types.AttributeValueMemberN{Value: "-7"},
				"Int64":   &types.AttributeValueMemberN{Value: "-8"},
				"Float32": &types.AttributeValueMemberN{Value: "9.9"},
				"Float64": &types.AttributeValueMemberN{Value: "10.1"},
			},
		},
		actual: &testEmptyCollectionsNulledNumericalScalars{},
		expected: testEmptyCollectionsNulledNumericalScalars{
			String:  "test string",
			Uint8:   1,
			Uint16:  2,
			Uint32:  3,
			Uint64:  4,
			Int8:    -5,
			Int16:   -6,
			Int32:   -7,
			Int64:   -8,
			Float32: 9.9,
			Float64: 10.1,
		},
	},
	"nil pointer scalars": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"PtrString":  &types.AttributeValueMemberNULL{Value: true},
				"PtrUint8":   &types.AttributeValueMemberNULL{Value: true},
				"PtrUint16":  &types.AttributeValueMemberNULL{Value: true},
				"PtrUint32":  &types.AttributeValueMemberNULL{Value: true},
				"PtrUint64":  &types.AttributeValueMemberNULL{Value: true},
				"PtrInt8":    &types.AttributeValueMemberNULL{Value: true},
				"PtrInt16":   &types.AttributeValueMemberNULL{Value: true},
				"PtrInt32":   &types.AttributeValueMemberNULL{Value: true},
				"PtrInt64":   &types.AttributeValueMemberNULL{Value: true},
				"PtrFloat32": &types.AttributeValueMemberNULL{Value: true},
				"PtrFloat64": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual:   &testEmptyCollectionsPtrScalars{},
		expected: testEmptyCollectionsPtrScalars{},
	},
	"non-nil pointer to scalars with zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"PtrString":  &types.AttributeValueMemberNULL{Value: true},
				"PtrUint8":   &types.AttributeValueMemberN{Value: "0"},
				"PtrUint16":  &types.AttributeValueMemberN{Value: "0"},
				"PtrUint32":  &types.AttributeValueMemberN{Value: "0"},
				"PtrUint64":  &types.AttributeValueMemberN{Value: "0"},
				"PtrInt8":    &types.AttributeValueMemberN{Value: "0"},
				"PtrInt16":   &types.AttributeValueMemberN{Value: "0"},
				"PtrInt32":   &types.AttributeValueMemberN{Value: "0"},
				"PtrInt64":   &types.AttributeValueMemberN{Value: "0"},
				"PtrFloat32": &types.AttributeValueMemberN{Value: "0"},
				"PtrFloat64": &types.AttributeValueMemberN{Value: "0"},
			},
		},
		actual: &testEmptyCollectionsPtrScalars{},
		expected: testEmptyCollectionsPtrScalars{
			PtrUint8:   aws.Uint8(0),
			PtrUint16:  aws.Uint16(0),
			PtrUint32:  aws.Uint32(0),
			PtrUint64:  aws.Uint64(0),
			PtrInt8:    aws.Int8(0),
			PtrInt16:   aws.Int16(0),
			PtrInt32:   aws.Int32(0),
			PtrInt64:   aws.Int64(0),
			PtrFloat32: aws.Float32(0),
			PtrFloat64: aws.Float64(0),
		},
	},
	"pointer scalars non-nil non-zero": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"PtrString":  &types.AttributeValueMemberS{Value: "test string"},
				"PtrUint8":   &types.AttributeValueMemberN{Value: "1"},
				"PtrUint16":  &types.AttributeValueMemberN{Value: "2"},
				"PtrUint32":  &types.AttributeValueMemberN{Value: "3"},
				"PtrUint64":  &types.AttributeValueMemberN{Value: "4"},
				"PtrInt8":    &types.AttributeValueMemberN{Value: "-5"},
				"PtrInt16":   &types.AttributeValueMemberN{Value: "-6"},
				"PtrInt32":   &types.AttributeValueMemberN{Value: "-7"},
				"PtrInt64":   &types.AttributeValueMemberN{Value: "-8"},
				"PtrFloat32": &types.AttributeValueMemberN{Value: "9.9"},
				"PtrFloat64": &types.AttributeValueMemberN{Value: "10.1"},
			},
		},
		actual: &testEmptyCollectionsPtrScalars{},
		expected: testEmptyCollectionsPtrScalars{
			PtrString:  aws.String("test string"),
			PtrUint8:   aws.Uint8(1),
			PtrUint16:  aws.Uint16(2),
			PtrUint32:  aws.Uint32(3),
			PtrUint64:  aws.Uint64(4),
			PtrInt8:    aws.Int8(-5),
			PtrInt16:   aws.Int16(-6),
			PtrInt32:   aws.Int32(-7),
			PtrInt64:   aws.Int64(-8),
			PtrFloat32: aws.Float32(9.9),
			PtrFloat64: aws.Float64(10.1),
		},
	},
	"omittable nil pointer scalars": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{},
		},
		actual:   &testEmptyCollectionsOmittedPtrNumericalScalars{},
		expected: testEmptyCollectionsOmittedPtrNumericalScalars{},
	},
	"omittable non-nil pointer to scalars with zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"PtrUint8":   &types.AttributeValueMemberN{Value: "0"},
				"PtrUint16":  &types.AttributeValueMemberN{Value: "0"},
				"PtrUint32":  &types.AttributeValueMemberN{Value: "0"},
				"PtrUint64":  &types.AttributeValueMemberN{Value: "0"},
				"PtrInt8":    &types.AttributeValueMemberN{Value: "0"},
				"PtrInt16":   &types.AttributeValueMemberN{Value: "0"},
				"PtrInt32":   &types.AttributeValueMemberN{Value: "0"},
				"PtrInt64":   &types.AttributeValueMemberN{Value: "0"},
				"PtrFloat32": &types.AttributeValueMemberN{Value: "0"},
				"PtrFloat64": &types.AttributeValueMemberN{Value: "0"},
			},
		},
		actual: &testEmptyCollectionsOmittedPtrNumericalScalars{},
		expected: testEmptyCollectionsOmittedPtrNumericalScalars{
			PtrUint8:   aws.Uint8(0),
			PtrUint16:  aws.Uint16(0),
			PtrUint32:  aws.Uint32(0),
			PtrUint64:  aws.Uint64(0),
			PtrInt8:    aws.Int8(0),
			PtrInt16:   aws.Int16(0),
			PtrInt32:   aws.Int32(0),
			PtrInt64:   aws.Int64(0),
			PtrFloat32: aws.Float32(0),
			PtrFloat64: aws.Float64(0),
		},
	},
	"omittable non-nil pointer to non-zero scalar": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"PtrUint8":   &types.AttributeValueMemberN{Value: "1"},
				"PtrUint16":  &types.AttributeValueMemberN{Value: "2"},
				"PtrUint32":  &types.AttributeValueMemberN{Value: "3"},
				"PtrUint64":  &types.AttributeValueMemberN{Value: "4"},
				"PtrInt8":    &types.AttributeValueMemberN{Value: "-5"},
				"PtrInt16":   &types.AttributeValueMemberN{Value: "-6"},
				"PtrInt32":   &types.AttributeValueMemberN{Value: "-7"},
				"PtrInt64":   &types.AttributeValueMemberN{Value: "-8"},
				"PtrFloat32": &types.AttributeValueMemberN{Value: "9.9"},
				"PtrFloat64": &types.AttributeValueMemberN{Value: "10.1"},
			},
		},
		actual: &testEmptyCollectionsOmittedPtrNumericalScalars{},
		expected: testEmptyCollectionsOmittedPtrNumericalScalars{
			PtrUint8:   aws.Uint8(1),
			PtrUint16:  aws.Uint16(2),
			PtrUint32:  aws.Uint32(3),
			PtrUint64:  aws.Uint64(4),
			PtrInt8:    aws.Int8(-5),
			PtrInt16:   aws.Int16(-6),
			PtrInt32:   aws.Int32(-7),
			PtrInt64:   aws.Int64(-8),
			PtrFloat32: aws.Float32(9.9),
			PtrFloat64: aws.Float64(10.1),
		},
	},
	"maps slices nil values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map":       &types.AttributeValueMemberNULL{Value: true},
				"Slice":     &types.AttributeValueMemberNULL{Value: true},
				"ByteSlice": &types.AttributeValueMemberNULL{Value: true},
				"ByteArray": &types.AttributeValueMemberB{Value: make([]byte, 4)},
				"ZeroArray": &types.AttributeValueMemberB{Value: make([]byte, 0)},
				"BinarySet": &types.AttributeValueMemberNULL{Value: true},
				"NumberSet": &types.AttributeValueMemberNULL{Value: true},
				"StringSet": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual:   &testEmptyCollectionTypes{},
		expected: testEmptyCollectionTypes{},
	},
	"null nil pointer scalars": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"PtrString":  &types.AttributeValueMemberNULL{Value: true},
				"PtrUint8":   &types.AttributeValueMemberNULL{Value: true},
				"PtrUint16":  &types.AttributeValueMemberNULL{Value: true},
				"PtrUint32":  &types.AttributeValueMemberNULL{Value: true},
				"PtrUint64":  &types.AttributeValueMemberNULL{Value: true},
				"PtrInt8":    &types.AttributeValueMemberNULL{Value: true},
				"PtrInt16":   &types.AttributeValueMemberNULL{Value: true},
				"PtrInt32":   &types.AttributeValueMemberNULL{Value: true},
				"PtrInt64":   &types.AttributeValueMemberNULL{Value: true},
				"PtrFloat32": &types.AttributeValueMemberNULL{Value: true},
				"PtrFloat64": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual:   &testEmptyCollectionsNulledPtrNumericalScalars{},
		expected: testEmptyCollectionsNulledPtrNumericalScalars{},
	},
	"null non-nil pointer to scalars with zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"PtrString":  &types.AttributeValueMemberS{Value: ""},
				"PtrUint8":   &types.AttributeValueMemberN{Value: "0"},
				"PtrUint16":  &types.AttributeValueMemberN{Value: "0"},
				"PtrUint32":  &types.AttributeValueMemberN{Value: "0"},
				"PtrUint64":  &types.AttributeValueMemberN{Value: "0"},
				"PtrInt8":    &types.AttributeValueMemberN{Value: "0"},
				"PtrInt16":   &types.AttributeValueMemberN{Value: "0"},
				"PtrInt32":   &types.AttributeValueMemberN{Value: "0"},
				"PtrInt64":   &types.AttributeValueMemberN{Value: "0"},
				"PtrFloat32": &types.AttributeValueMemberN{Value: "0"},
				"PtrFloat64": &types.AttributeValueMemberN{Value: "0"},
			},
		},
		actual: &testEmptyCollectionsNulledPtrNumericalScalars{},
		expected: testEmptyCollectionsNulledPtrNumericalScalars{
			PtrString:  aws.String(""),
			PtrUint8:   aws.Uint8(0),
			PtrUint16:  aws.Uint16(0),
			PtrUint32:  aws.Uint32(0),
			PtrUint64:  aws.Uint64(0),
			PtrInt8:    aws.Int8(0),
			PtrInt16:   aws.Int16(0),
			PtrInt32:   aws.Int32(0),
			PtrInt64:   aws.Int64(0),
			PtrFloat32: aws.Float32(0),
			PtrFloat64: aws.Float64(0),
		},
	},
	"null non-nil pointer to non-zero scalar": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"PtrString":  &types.AttributeValueMemberS{Value: "abc"},
				"PtrUint8":   &types.AttributeValueMemberN{Value: "1"},
				"PtrUint16":  &types.AttributeValueMemberN{Value: "2"},
				"PtrUint32":  &types.AttributeValueMemberN{Value: "3"},
				"PtrUint64":  &types.AttributeValueMemberN{Value: "4"},
				"PtrInt8":    &types.AttributeValueMemberN{Value: "-5"},
				"PtrInt16":   &types.AttributeValueMemberN{Value: "-6"},
				"PtrInt32":   &types.AttributeValueMemberN{Value: "-7"},
				"PtrInt64":   &types.AttributeValueMemberN{Value: "-8"},
				"PtrFloat32": &types.AttributeValueMemberN{Value: "9.9"},
				"PtrFloat64": &types.AttributeValueMemberN{Value: "10.1"},
			},
		},
		actual: &testEmptyCollectionsNulledPtrNumericalScalars{},
		expected: testEmptyCollectionsNulledPtrNumericalScalars{
			PtrString:  aws.String("abc"),
			PtrUint8:   aws.Uint8(1),
			PtrUint16:  aws.Uint16(2),
			PtrUint32:  aws.Uint32(3),
			PtrUint64:  aws.Uint64(4),
			PtrInt8:    aws.Int8(-5),
			PtrInt16:   aws.Int16(-6),
			PtrInt32:   aws.Int32(-7),
			PtrInt64:   aws.Int64(-8),
			PtrFloat32: aws.Float32(9.9),
			PtrFloat64: aws.Float64(10.1),
		},
	},
	"maps slices zero values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map":       &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				"Slice":     &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
				"ByteSlice": &types.AttributeValueMemberB{Value: []byte{}},
				"ByteArray": &types.AttributeValueMemberB{Value: make([]byte, 4)},
				"ZeroArray": &types.AttributeValueMemberB{Value: make([]byte, 0)},
				// sets are special and not serialized to empty if no elements
				"BinarySet": &types.AttributeValueMemberBS{Value: [][]byte{}},
				"NumberSet": &types.AttributeValueMemberNS{Value: []string{}},
				"StringSet": &types.AttributeValueMemberSS{Value: []string{}},
			},
		},
		inMarshal: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map":       &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				"Slice":     &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
				"ByteSlice": &types.AttributeValueMemberB{Value: []byte{}},
				"ByteArray": &types.AttributeValueMemberB{Value: make([]byte, 4)},
				"ZeroArray": &types.AttributeValueMemberB{Value: make([]byte, 0)},
				// sets are special and not serialized to empty if no elements
				"BinarySet": &types.AttributeValueMemberNULL{Value: true},
				"NumberSet": &types.AttributeValueMemberNULL{Value: true},
				"StringSet": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual: &testEmptyCollectionTypes{},
		expected: testEmptyCollectionTypes{
			Map:       map[string]string{},
			Slice:     []string{},
			ByteSlice: []byte{},
			ByteArray: [4]byte{},
			ZeroArray: [0]byte{},
			BinarySet: [][]byte{},
			NumberSet: []int{},
			StringSet: []string{},
		},
	},
	"maps slices non-zero values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"key": &types.AttributeValueMemberS{Value: "value"},
					},
				},
				"Slice": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "test"},
					&types.AttributeValueMemberS{Value: "slice"},
				}},
				"ByteSlice": &types.AttributeValueMemberB{Value: []byte{0, 1}},
				"ByteArray": &types.AttributeValueMemberB{Value: []byte{0, 1, 2, 3}},
				"ZeroArray": &types.AttributeValueMemberB{Value: make([]byte, 0)},
				"BinarySet": &types.AttributeValueMemberBS{Value: [][]byte{{0, 1}, {2, 3}}},
				"NumberSet": &types.AttributeValueMemberNS{Value: []string{"0", "1"}},
				"StringSet": &types.AttributeValueMemberSS{Value: []string{"test", "slice"}},
			},
		},
		actual: &testEmptyCollectionTypes{},
		expected: testEmptyCollectionTypes{
			Map:       map[string]string{"key": "value"},
			Slice:     []string{"test", "slice"},
			ByteSlice: []byte{0, 1},
			ByteArray: [4]byte{0, 1, 2, 3},
			ZeroArray: [0]byte{},
			BinarySet: [][]byte{{0, 1}, {2, 3}},
			NumberSet: []int{0, 1},
			StringSet: []string{"test", "slice"},
		},
	},
	"omittable maps slices nil values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"ByteArray": &types.AttributeValueMemberB{Value: make([]byte, 4)},
			},
		},
		actual:   &testEmptyCollectionTypesOmitted{},
		expected: testEmptyCollectionTypesOmitted{},
	},
	"omittable maps slices zero values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map":       &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				"Slice":     &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
				"ByteSlice": &types.AttributeValueMemberB{Value: []byte{}},
				"ByteArray": &types.AttributeValueMemberB{Value: make([]byte, 4)},
				"BinarySet": &types.AttributeValueMemberBS{Value: [][]byte{}},
				"NumberSet": &types.AttributeValueMemberNS{Value: []string{}},
				"StringSet": &types.AttributeValueMemberSS{Value: []string{}},
			},
		},
		inMarshal: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map":       &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				"Slice":     &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
				"ByteSlice": &types.AttributeValueMemberB{Value: []byte{}},
				"ByteArray": &types.AttributeValueMemberB{Value: make([]byte, 4)},
				"BinarySet": &types.AttributeValueMemberNULL{Value: true},
				"NumberSet": &types.AttributeValueMemberNULL{Value: true},
				"StringSet": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual: &testEmptyCollectionTypesOmitted{},
		expected: testEmptyCollectionTypesOmitted{
			Map:       map[string]string{},
			Slice:     []string{},
			ByteSlice: []byte{},
			ByteArray: [4]byte{},
			BinarySet: [][]byte{},
			NumberSet: []int{},
			StringSet: []string{},
		},
	},
	"omittable maps slices non-zero values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"key": &types.AttributeValueMemberS{Value: "value"},
					},
				},
				"Slice": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "test"},
					&types.AttributeValueMemberS{Value: "slice"},
				}},
				"ByteSlice": &types.AttributeValueMemberB{Value: []byte{0, 1}},
				"ByteArray": &types.AttributeValueMemberB{Value: []byte{0, 1, 2, 3}},
				"BinarySet": &types.AttributeValueMemberBS{Value: [][]byte{{0, 1}, {2, 3}}},
				"NumberSet": &types.AttributeValueMemberNS{Value: []string{"0", "1"}},
				"StringSet": &types.AttributeValueMemberSS{Value: []string{"test", "slice"}},
			},
		},
		actual: &testEmptyCollectionTypesOmitted{},
		expected: testEmptyCollectionTypesOmitted{
			Map:       map[string]string{"key": "value"},
			Slice:     []string{"test", "slice"},
			ByteSlice: []byte{0, 1},
			ByteArray: [4]byte{0, 1, 2, 3},
			ZeroArray: [0]byte{},
			BinarySet: [][]byte{{0, 1}, {2, 3}},
			NumberSet: []int{0, 1},
			StringSet: []string{"test", "slice"},
		},
	},
	"null maps slices nil values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map":       &types.AttributeValueMemberNULL{Value: true},
				"Slice":     &types.AttributeValueMemberNULL{Value: true},
				"ByteSlice": &types.AttributeValueMemberNULL{Value: true},
				"ByteArray": &types.AttributeValueMemberB{Value: make([]byte, 4)},
				"BinarySet": &types.AttributeValueMemberNULL{Value: true},
				"NumberSet": &types.AttributeValueMemberNULL{Value: true},
				"StringSet": &types.AttributeValueMemberNULL{Value: true},
				"ZeroArray": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual:   &testEmptyCollectionTypesNulled{},
		expected: testEmptyCollectionTypesNulled{},
	},
	"null maps slices zero values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map":       &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				"Slice":     &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
				"ByteSlice": &types.AttributeValueMemberB{Value: []byte{}},
				"ByteArray": &types.AttributeValueMemberB{Value: make([]byte, 4)},
				"BinarySet": &types.AttributeValueMemberBS{Value: [][]byte{}},
				"NumberSet": &types.AttributeValueMemberNS{Value: []string{}},
				"StringSet": &types.AttributeValueMemberSS{Value: []string{}},
			},
		},
		inMarshal: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map":       &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				"Slice":     &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
				"ByteSlice": &types.AttributeValueMemberB{Value: []byte{}},
				"ByteArray": &types.AttributeValueMemberB{Value: make([]byte, 4)},
				"BinarySet": &types.AttributeValueMemberNULL{Value: true},
				"NumberSet": &types.AttributeValueMemberNULL{Value: true},
				"StringSet": &types.AttributeValueMemberNULL{Value: true},
				"ZeroArray": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual: &testEmptyCollectionTypesNulled{},
		expected: testEmptyCollectionTypesNulled{
			Map:       map[string]string{},
			Slice:     []string{},
			ByteSlice: []byte{},
			ByteArray: [4]byte{},
			BinarySet: [][]byte{},
			NumberSet: []int{},
			StringSet: []string{},
		},
	},
	"null maps slices non-zero values": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Map": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"key": &types.AttributeValueMemberS{Value: "value"},
					},
				},
				"Slice": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "test"},
					&types.AttributeValueMemberS{Value: "slice"},
				}},
				"ByteSlice": &types.AttributeValueMemberB{Value: []byte{0, 1}},
				"ByteArray": &types.AttributeValueMemberB{Value: []byte{0, 1, 2, 3}},
				"BinarySet": &types.AttributeValueMemberBS{Value: [][]byte{{0, 1}, {2, 3}}},
				"NumberSet": &types.AttributeValueMemberNS{Value: []string{"0", "1"}},
				"StringSet": &types.AttributeValueMemberSS{Value: []string{"test", "slice"}},
				"ZeroArray": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual: &testEmptyCollectionTypesNulled{},
		expected: testEmptyCollectionTypesNulled{
			Map:       map[string]string{"key": "value"},
			Slice:     []string{"test", "slice"},
			ByteSlice: []byte{0, 1},
			ByteArray: [4]byte{0, 1, 2, 3},
			ZeroArray: [0]byte{},
			BinarySet: [][]byte{{0, 1}, {2, 3}},
			NumberSet: []int{0, 1},
			StringSet: []string{"test", "slice"},
		},
	},
	"structs with members zero": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Struct": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Int": &types.AttributeValueMemberN{Value: "0"},
					},
				},
				"PtrStruct": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual: &struct {
			Struct    testEmptyCollectionStruct
			PtrStruct *testEmptyCollectionStruct
		}{},
		expected: struct {
			Struct    testEmptyCollectionStruct
			PtrStruct *testEmptyCollectionStruct
		}{},
	},
	"structs with members non-zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Struct": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Int": &types.AttributeValueMemberN{Value: "1"},
					},
				},
				"PtrStruct": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Int": &types.AttributeValueMemberN{Value: "1"},
					},
				},
			},
		},
		actual: &struct {
			Struct    testEmptyCollectionStruct
			PtrStruct *testEmptyCollectionStruct
		}{},
		expected: struct {
			Struct    testEmptyCollectionStruct
			PtrStruct *testEmptyCollectionStruct
		}{
			Struct:    testEmptyCollectionStruct{Int: 1},
			PtrStruct: &testEmptyCollectionStruct{Int: 1},
		},
	},
	"struct with omittable members zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Struct":    &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
				"PtrStruct": &types.AttributeValueMemberNULL{Value: true},
			},
		},
		actual: &struct {
			Struct    testEmptyCollectionStructOmitted
			PtrStruct *testEmptyCollectionStructOmitted
		}{},
		expected: struct {
			Struct    testEmptyCollectionStructOmitted
			PtrStruct *testEmptyCollectionStructOmitted
		}{},
	},
	"omittable struct with omittable members zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Struct": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
			},
		},
		actual: &struct {
			Struct    testEmptyCollectionStructOmitted  `dynamodbav:",omitempty"`
			PtrStruct *testEmptyCollectionStructOmitted `dynamodbav:",omitempty"`
		}{},
		expected: struct {
			Struct    testEmptyCollectionStructOmitted  `dynamodbav:",omitempty"`
			PtrStruct *testEmptyCollectionStructOmitted `dynamodbav:",omitempty"`
		}{},
	},
	"omittable struct with omittable members non-zero value": {
		in: &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"Struct": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Slice": &types.AttributeValueMemberL{Value: []types.AttributeValue{
							&types.AttributeValueMemberS{Value: "test"},
						}},
					},
				},
				"InitPtrStruct": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Slice": &types.AttributeValueMemberL{Value: []types.AttributeValue{
							&types.AttributeValueMemberS{Value: "test"},
						}},
					},
				},
			},
		},
		actual: &struct {
			Struct        testEmptyCollectionStructOmitted  `dynamodbav:",omitempty"`
			InitPtrStruct *testEmptyCollectionStructOmitted `dynamodbav:",omitempty"`
		}{},
		expected: struct {
			Struct        testEmptyCollectionStructOmitted  `dynamodbav:",omitempty"`
			InitPtrStruct *testEmptyCollectionStructOmitted `dynamodbav:",omitempty"`
		}{
			Struct:        testEmptyCollectionStructOmitted{Slice: []string{"test"}},
			InitPtrStruct: &testEmptyCollectionStructOmitted{Slice: []string{"test"}},
		},
	},
}

func TestMarshalEmptyCollections(t *testing.T) {
	for name, c := range sharedEmptyCollectionsTestCases {
		t.Run(name, func(t *testing.T) {
			av, err := Marshal(c.expected)
			in := c.in
			if c.inMarshal != nil {
				in = c.inMarshal
			}
			assertConvertTest(t, av, in, err, c.err)
		})
	}
}

func TestEmptyCollectionsSpecialCases(t *testing.T) {
	// ptr string non nil with empty value

	type SpecialCases struct {
		PtrString        *string
		OmittedString    string  `dynamodbav:",omitempty"`
		OmittedPtrString *string `dynamodbav:",omitempty"`
	}

	expectedEncode := &types.AttributeValueMemberM{
		Value: map[string]types.AttributeValue{
			"PtrString": &types.AttributeValueMemberS{Value: ""},
		},
	}
	expectedDecode := SpecialCases{}

	actualEncode, err := Marshal(&SpecialCases{
		PtrString:        aws.String(""),
		OmittedString:    "",
		OmittedPtrString: nil,
	})
	if err != nil {
		t.Fatalf("expected no err got %v", err)
	}
	if diff := cmp.Diff(expectedEncode, actualEncode); len(diff) != 0 {
		t.Errorf("expected encode match\n%s", diff)
	}

	var actualDecode SpecialCases
	var av types.AttributeValue
	err = Unmarshal(av, &actualDecode)
	if err != nil {
		t.Fatalf("expected no err got %v", err)
	}
	if diff := cmp.Diff(expectedDecode, actualDecode); len(diff) != 0 {
		t.Errorf("expected dencode match\n%s", diff)
	}
}

func TestUnmarshalEmptyCollections(t *testing.T) {
	for name, c := range sharedEmptyCollectionsTestCases {
		t.Run(name, func(t *testing.T) {
			err := Unmarshal(c.in, c.actual)
			assertConvertTest(t, c.actual, c.expected, err, c.err)
		})
	}
}
