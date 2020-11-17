package s3shared

import (
	"context"
)

// clonedInputKey used to denote if request input was cloned.
type clonedInputKey struct{}

// SetClonedInputKey sets a key on context to denote input was cloned previously
func SetClonedInputKey(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, clonedInputKey{}, value)
}

// IsClonedInput retrieves if context key for cloned input was set.
// If set, we can infer that the reuqest input was cloned previously.
func IsClonedInput(ctx context.Context) bool {
	v, _ := ctx.Value(clonedInputKey{}).(bool)
	return v
}
