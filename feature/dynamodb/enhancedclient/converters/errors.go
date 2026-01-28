package converters

import (
	"errors"
	"fmt"
	"strings"
)

// ErrNilValue is returned when a nil value is encountered where a non-nil value is required.
var ErrNilValue = errors.New("nil value error")

// unsupportedType returns a formatted error indicating the provided type is not supported.
// Optionally lists the supported types for better diagnostics.
func unsupportedType(unsupported any, supported ...any) error {
	err := fmt.Errorf("unsupported type: %T", unsupported)

	if len(supported) > 0 {
		var sup []string

		for i := range supported {
			sup = append(sup, fmt.Sprintf("%T", supported[i]))
		}

		err = fmt.Errorf("expected %s, got %s", strings.Join(sup, " or "), err.Error())
	}

	return err
}
