package converters

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// AnyAttributeConverter is a non-generic abstraction for attribute converters.
// It provides a type-erased interface so heterogeneous converters can be stored
// together (e.g., in a map[string]AnyAttributeConverter) and invoked dynamically.
// Implementations must accept a slice of tag-derived options that may influence
// conversion behavior.
type AnyAttributeConverter interface {
	FromAttributeValue(types.AttributeValue, []string) (any, error)
	ToAttributeValue(any, []string) (types.AttributeValue, error)
}

// AttributeConverter defines conversion logic between a concrete Go type T and DynamoDB
// AttributeValues. Implementations encode/decode a value of T, optionally using
// tag-derived options supplied via the []string parameter.
type AttributeConverter[T any] interface {
	// FromAttributeValue converts a DynamoDB AttributeValue to the Go type T.
	// The second argument provides tag options for the converter.
	FromAttributeValue(types.AttributeValue, []string) (T, error)
	// ToAttributeValue converts a value of type T to a DynamoDB AttributeValue.
	ToAttributeValue(T, []string) (types.AttributeValue, error)
}

// Wrapper adapts an AttributeConverter[T] to the non-generic AnyAttributeConverter
// interface so converters for different concrete types can coexist in the same
// registry structure. Without this indirection, generic constraints would prevent
// a uniform collection (e.g. map[string]AttributeConverter[T]) spanning multiple T.
//
// Wrapper assumes the wrapped converter's methods are safe for concurrent use.
// It performs a runtime type assertion when converting values back to AttributeValue.
// If the provided value does not match T (and is not nil) an unsupportedType error
// is returned.
type Wrapper[T any] struct {
	Impl AttributeConverter[T]
}

// FromAttributeValue delegates to the underlying AttributeConverter[T] and returns
// the resulting value boxed as any. Tag-derived options are forwarded unchanged.
// Errors from the underlying converter are propagated as-is.
func (w *Wrapper[T]) FromAttributeValue(attr types.AttributeValue, opts []string) (any, error) {
	return w.Impl.FromAttributeValue(attr, opts)
}

// ToAttributeValue attempts to cast the supplied value to T (allowing nil) and
// delegates to the underlying converter. If the dynamic type does not match T,
// unsupportedType is returned. A nil value is passed through as the zero value
// of T; underlying converters must define how they treat it (often yielding a
// ErrNilValue for pointer types or encoding a zero-value for value types).
func (w *Wrapper[T]) ToAttributeValue(value any, opts []string) (types.AttributeValue, error) {
	if v, ok := value.(T); ok || value == nil {
		return w.Impl.ToAttributeValue(v, opts)
	}

	return nil, unsupportedType(value, *new(T))
}
