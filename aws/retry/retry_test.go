package retry_test

import (
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/awserr"
	"github.com/jviney/aws-sdk-go-v2/aws/retry"
)

func TestAddWithErrorCodes(t *testing.T) {
	cases := map[string]struct {
		Err    error
		Expect bool
	}{
		"retryable": {
			Err:    awserr.New("Error1", "err", nil),
			Expect: true,
		},
		"not retryable": {
			Err:    awserr.New("Error3", "err", nil),
			Expect: false,
		},
	}

	r := retry.AddWithErrorCodes(aws.NoOpRetryer{}, "Error1", "Error2")

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.Expect, r.IsErrorRetryable(c.Err); e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
