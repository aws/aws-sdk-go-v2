package retry

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	smithytraits "github.com/aws/smithy-go/traits"
)

// retryableTraitError is implemented by every modeled error type. The bool
// reports whether the error is modeled with the smithy.api#retryable trait.
type retryableTraitError interface {
	RetryableTrait() (smithytraits.Retryable, bool)
}

// modeledRetryableTrait returns the smithy.api#retryable trait of the first
// modeled error in err's chain, and whether that error is modeled with it.
func modeledRetryableTrait(err error) (smithytraits.Retryable, bool) {
	var v retryableTraitError
	if !errors.As(err, &v) {
		return smithytraits.Retryable{}, false
	}
	return v.RetryableTrait()
}

// RetryableTrait is a retryable check for modeled errors targeted by the
// smithy.api#retryable trait. An error without the trait is not decided by
// this check, so other checks such as status codes and error codes still
// apply to it.
type RetryableTrait struct{}

// IsErrorRetryable returns true if the error is modeled with the
// smithy.api#retryable trait. Returns unknown otherwise.
func (RetryableTrait) IsErrorRetryable(err error) aws.Ternary {
	if _, ok := modeledRetryableTrait(err); ok {
		return aws.TrueTernary
	}
	return aws.UnknownTernary
}

// ThrottleTrait is a throttle check for modeled errors targeted by the
// smithy.api#retryable trait with throttling set to true. An error without
// the trait, or with throttling unset, is not decided by this check, so other
// checks such as error codes still apply to it.
type ThrottleTrait struct{}

// IsErrorThrottle returns true if the error is modeled with the
// smithy.api#retryable trait and throttling set to true. Returns unknown
// otherwise.
func (ThrottleTrait) IsErrorThrottle(err error) aws.Ternary {
	if trait, ok := modeledRetryableTrait(err); ok && trait.Throttling {
		return aws.TrueTernary
	}
	return aws.UnknownTernary
}
