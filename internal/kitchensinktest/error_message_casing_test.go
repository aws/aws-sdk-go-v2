package kitchensinktest

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/kitchensinktest/types"
)

// The service models the error message member as "message", but real services
// are inconsistent about the casing they send on the wire.
func TestModeledErrorMessageCasing(t *testing.T) {
	for _, key := range []string{"message", "Message"} {
		t.Run(key, func(t *testing.T) {
			body := newTrackedBody([]byte(`{"__type":"ItemNotFound","` + key + `":"not here"}`))
			svc := closeTestClient(400, body)

			_, err := svc.GetItem(context.Background(), &GetItemInput{})

			var nf *types.ItemNotFound
			if !errors.As(err, &nf) {
				t.Fatalf("expect *types.ItemNotFound, got %T: %v", err, err)
			}
			if actual := nf.ErrorMessage(); actual != "not here" {
				t.Errorf("expect message %q, got %q", "not here", actual)
			}
		})
	}
}
