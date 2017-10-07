// +build !go1.6

package aws_test

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
)

var errTimeout = awserr.New("foo", "bar", errors.New("net/http: request canceled Timeout"))
