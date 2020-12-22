package customizations

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/ptr"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestPredictEndpointMiddleware(t *testing.T) {
	cases := map[string]struct {
		PredictEndpoint  *string
		ExpectedEndpoint string
		ExpectedErr      string
	}{
		"nil endpoint": {},
		"empty endpoint": {
			PredictEndpoint: ptr.String(""),
		},
		"invalid endpoint": {
			PredictEndpoint: ptr.String("::::::::"),
			ExpectedErr:     "unable to parse",
		},
		"valid endpoint": {
			PredictEndpoint:  ptr.String("https://example.amazonaws.com/"),
			ExpectedEndpoint: "https://example.amazonaws.com/",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			m := &predictEndpoint{
				fetchPredictEndpoint: func(i interface{}) (*string, error) {
					return c.PredictEndpoint, nil
				},
			}
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

						if c.ExpectedEndpoint != req.URL.String() {
							t.Errorf("expected url to be `%v`, but was `%v`", c.ExpectedEndpoint, req.URL.String())
						}

						return output, metadata, err
					}),
			)
			if len(c.ExpectedErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectedErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect error to contain %v, got %v", e, a)
				}
			} else {
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
			}
		})
	}

}
