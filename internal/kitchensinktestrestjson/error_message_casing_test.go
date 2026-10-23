package kitchensinktestrestjson

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/kitchensinktestrestjson/types"
)

// The service models the error message member as "message", but real services
// are inconsistent about the casing they send on the wire.
func TestModeledErrorMessageCasing(t *testing.T) {
	for _, key := range []string{"message", "Message"} {
		t.Run(key, func(t *testing.T) {
			body := newCloseTrackingBody(`{"__type":"ResourceNotFound","` + key + `":"not here"}`)
			svc := closeTestClient(404, http.Header{}, body)

			_, err := svc.GetResource(context.Background(), &GetResourceInput{Id: strPtr("some-id")})

			var nf *types.ResourceNotFound
			if !errors.As(err, &nf) {
				t.Fatalf("expect *types.ResourceNotFound, got %T: %v", err, err)
			}
			if actual := nf.ErrorMessage(); actual != "not here" {
				t.Errorf("expect message %q, got %q", "not here", actual)
			}
		})
	}
}
