package customizations

import (
	"context"
	"testing"

	"github.com/aws/smithy-go/middleware"
)

type accountIDBearer struct {
	AccountID *string
}

func TestDefaultAccountIDMiddleware(t *testing.T) {
	m := &DefaultAccountID{
		setDefaultAccountID: func(input interface{}, defaultAccountID string) interface{} {
			if bearer, ok := input.(accountIDBearer); ok && bearer.AccountID == nil {
				bearer.AccountID = &defaultAccountID
				return bearer
			}
			return input
		},
	}

	_, _, err := m.HandleInitialize(context.Background(),
		middleware.InitializeInput{
			Parameters: accountIDBearer{},
		},
		middleware.InitializeHandlerFunc(
			func(ctx context.Context, input middleware.InitializeInput) (
				out middleware.InitializeOutput, metadata middleware.Metadata, err error,
			) {
				params, ok := input.Parameters.(accountIDBearer)
				if !ok {
					t.Fatalf("expect struct input, got %T", input.Parameters)
				}
				if params.AccountID == nil || *params.AccountID != "-" {
					t.Errorf("expect `-` AccountID, got %v", params.AccountID)
				}
				return out, metadata, err
			}),
	)

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
}
