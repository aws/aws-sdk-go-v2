package customizations

import (
	"context"
	"fmt"
	"github.com/awslabs/smithy-go"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
	"net/http"
	"net/url"
	"testing"
)

func TestSanitizeURLMiddleware(t *testing.T) {
	cases := map[string]struct {
		Given    string
		Expected string
	}{
		"includes hostedzone": {
			Given:    "https://example.amazonaws.com/2013-04-01/hostedzone/%2Fhostedzone%2FABCDEFG?abc=123",
			Expected: "https://example.amazonaws.com/2013-04-01/hostedzone/ABCDEFG?abc=123",
		},
		"excludes hostedzone": {
			Given:    "https://example.amazonaws.com/2013-04-01/hostedzone/ABCDEFG?abc=123",
			Expected: "https://example.amazonaws.com/2013-04-01/hostedzone/ABCDEFG?abc=123",
		},
	}

	m := &sanitizeURLMiddleware{}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, _, err := m.HandleSerialize(context.Background(),
				middleware.SerializeInput{
					Request: func() interface{} {
						uri, err := url.Parse(c.Given)
						if err != nil {
							t.Fatalf("failed to parse given uri, %v", c.Given)
						}
						return &smithyhttp.Request{
							Request: &http.Request{
								URL:    uri,
								Header: http.Header{},
							},
						}
					}(),
				},
				middleware.SerializeHandlerFunc(
					func(ctx context.Context, input middleware.SerializeInput) (
						output middleware.SerializeOutput, metadata middleware.Metadata, err error,
					) {
						req, ok := input.Request.(*smithyhttp.Request)
						if !ok {
							return output, metadata, &smithy.SerializationError{
								Err: fmt.Errorf("unknown request type %T", input.Request),
							}
						}

						if req.URL.String() != c.Expected {
							t.Errorf("expected url to be `%s`, but was `%s`", c.Expected, req.URL.String())
						}
						return output, metadata, err
					}),
			)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
