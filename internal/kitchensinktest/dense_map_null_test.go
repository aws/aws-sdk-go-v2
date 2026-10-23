package kitchensinktest

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/kitchensinktest/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type denseMapNullEndpointResolver struct{}

func (*denseMapNullEndpointResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (smithyendpoints.Endpoint, error) {
	return smithyendpoints.Endpoint{URI: url.URL{Scheme: "https", Host: "test.example.com"}}, nil
}

// An explicit null in a dense map is not valid per the spec, but we tolerate
// it by keeping the key with the zero value.
func TestDenseMapNullValue(t *testing.T) {
	body := `{
		"strings": {"a": "x", "b": null},
		"integers": {"a": 1, "b": null},
		"structs": {"a": {"name": "x"}, "b": null},
		"sparseStrings": {"a": "x", "b": null}
	}`

	svc := New(Options{
		Region: "us-east-1",
		HTTPClient: smithyhttp.ClientDoFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode:    200,
				Header:        http.Header{"Content-Type": []string{"application/x-amz-json-1.0"}},
				Body:          io.NopCloser(strings.NewReader(body)),
				ContentLength: -1,
				Request:       req,
			}, nil
		}),
		EndpointResolverV2: &denseMapNullEndpointResolver{},
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				s.Finalize.Clear()
				return nil
			},
		},
	})

	out, err := svc.GetMaps(context.Background(), &GetMapsInput{})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if expect := map[string]string{"a": "x", "b": ""}; !reflect.DeepEqual(expect, out.Strings) {
		t.Errorf("strings: expect %v, got %v", expect, out.Strings)
	}
	if expect := map[string]int32{"a": 1, "b": 0}; !reflect.DeepEqual(expect, out.Integers) {
		t.Errorf("integers: expect %v, got %v", expect, out.Integers)
	}
	if expect := map[string]types.MapValueStruct{"a": {Name: aws.String("x")}, "b": {}}; !reflect.DeepEqual(expect, out.Structs) {
		t.Errorf("structs: expect %v, got %v", expect, out.Structs)
	}
	if expect := map[string]*string{"a": aws.String("x"), "b": nil}; !reflect.DeepEqual(expect, out.SparseStrings) {
		t.Errorf("sparseStrings: expect %v, got %v", expect, out.SparseStrings)
	}
}
