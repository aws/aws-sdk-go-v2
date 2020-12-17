package retry_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
)

func TestAddWithErrorCodes(t *testing.T) {
	cases := map[string]struct {
		Err    error
		Expect bool
	}{
		"retryable": {
			Err:    &mockErrorCodeError{code: "Error1"},
			Expect: true,
		},
		"not retryable": {
			Err:    &mockErrorCodeError{code: "Error3"},
			Expect: false,
		},
	}

	r := retry.AddWithErrorCodes(aws.NopRetryer{}, "Error1", "Error2")

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.Expect, r.IsErrorRetryable(c.Err); e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
