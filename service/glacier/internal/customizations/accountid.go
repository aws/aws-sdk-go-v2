package customizations

import (
	"context"
	"github.com/awslabs/smithy-go/middleware"
)

type setDefaultAccountID func(input interface{}, accountID string) interface{}

// AddDefaultAccountIDMiddleware adds the DefaultAccountIDMiddleware to the stack using
// the given options.
func AddDefaultAccountIDMiddleware(stack *middleware.Stack, setDefaultAccountID setDefaultAccountID) {
	stack.Initialize.Add(&DefaultAccountIDMiddleware{
		setDefaultAccountID: setDefaultAccountID,
	}, middleware.Before)
}

// DefaultAccountIDMiddleware sets the account ID to "-" if it isn't already set
type DefaultAccountIDMiddleware struct {
	setDefaultAccountID setDefaultAccountID
}

// ID returns the id of the middleware
func (*DefaultAccountIDMiddleware) ID() string { return "Glacier:DefaultAccountID" }

// HandleInitialize implements the InitializeMiddleware interface
func (m *DefaultAccountIDMiddleware) HandleInitialize(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	in.Parameters = m.setDefaultAccountID(in.Parameters, "-")
	return next.HandleInitialize(ctx, in)
}
