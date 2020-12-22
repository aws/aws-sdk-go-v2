package customizations

import (
	"context"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"testing"
)

func TestGlacierAPIVersionMiddleware(t *testing.T) {
	apiVersion := "2012-06-01"
	m := &GlacierAPIVersion{apiVersion: apiVersion}

	_, _, err := m.HandleSerialize(context.Background(),
		middleware.SerializeInput{
			Request: smithyhttp.NewStackRequest(),
		},
		middleware.SerializeHandlerFunc(
			func(ctx context.Context, input middleware.SerializeInput) (
				output middleware.SerializeOutput, metadata middleware.Metadata, err error,
			) {
				req, ok := input.Request.(*smithyhttp.Request)
				if !ok || req == nil {
					t.Fatalf("expect smithy request, got %T", input.Request)
				}

				actual := req.Header.Get(glacierAPIVersionHeaderKey)
				if actual != apiVersion {
					t.Errorf("expect %s glacier version, got %s", apiVersion, actual)
				}
				return output, metadata, err
			}),
	)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
}
