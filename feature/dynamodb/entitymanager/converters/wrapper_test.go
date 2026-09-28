package converters

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ AttributeConverter[any] = (*mockConverter)(nil)

type mockConverter struct {
	retFrom    any
	retTo      types.AttributeValue
	fromCalled bool
	retFromErr error
	retToErr   error
	toCalled   bool
}

func (d *mockConverter) FromAttributeValue(_ types.AttributeValue, _ []string) (any, error) {
	d.fromCalled = true

	return d.retFrom, d.retFromErr
}

func (d *mockConverter) ToAttributeValue(_ any, _ []string) (types.AttributeValue, error) {
	d.toCalled = true

	return d.retTo, d.retToErr
}

func TestWrapper_FromAttributeValue(t *testing.T) {
	d := &mockConverter{
		retFrom:    0,
		retFromErr: ErrNilValue,
	}
	dw := &Wrapper[any]{Impl: d}

	actualOutput, actualError := dw.FromAttributeValue(nil, nil)

	comparisons := [][]any{
		{d.retFrom, actualOutput},
		{d.retFromErr, actualError},
		{d.fromCalled, true},
	}

	for _, cmp := range comparisons {
		if !reflect.DeepEqual(cmp[0], cmp[1]) {
			t.Fatalf("%#+v != %#+v", cmp[0], cmp[1])
		}
	}
}

func TestWrapper_ToAttributeValue(t *testing.T) {
	d := &mockConverter{
		retTo: &types.AttributeValueMemberS{
			Value: "test",
		},
		retToErr: ErrNilValue,
	}
	dw := &Wrapper[any]{Impl: d}

	actualOutput, actualError := dw.ToAttributeValue(nil, nil)

	comparisons := [][]any{
		{d.retTo, actualOutput},
		{d.retToErr, actualError},
		{d.toCalled, true},
	}

	for _, cmp := range comparisons {
		if !reflect.DeepEqual(cmp[0], cmp[1]) {
			t.Logf("%#+v != %#+v", cmp[0], cmp[1])
		}
	}
}
