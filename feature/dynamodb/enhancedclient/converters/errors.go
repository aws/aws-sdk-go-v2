package converters

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNilValue = errors.New("nil value error")

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
