package converters

import (
	"reflect"
)

// getType returns the canonical string form of the dynamic type of x,
// as produced by reflect.TypeOf(x).String(). This value is used as the
// key in converter registries (e.g., "int", "*time.Time", "[]byte").
// It is faster than fmt.Sprintf("%T", x) (avoids formatting machinery)
// and slightly faster than obtaining a generic type via reflect.TypeFor[T].
//
// Benchmark (indicative only; values vary by Go version, hardware, flags):
//
//	fmt.Sprintf(\"%T\", x)    ~ slower
//	reflect.TypeFor[T]()     ~ medium
//	reflect.TypeOf(x).String() ~ fastest in our measurements
//
// Caveat: Passing a nil interface (x == nil) will cause a panic because
// reflect.TypeOf(nil) == nil and the subsequent call to .String() dereferences nil.
// Callers must ensure x is non-nil.
func getType(x any) string {
	return reflect.TypeOf(x).String()
}
