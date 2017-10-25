// +build appengine plan9

package aws_test

import (
	"errors"
)

var stubConnectionResetError = errors.New("connection reset")
