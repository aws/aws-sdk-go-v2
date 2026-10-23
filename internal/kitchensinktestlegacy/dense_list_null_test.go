package kitchensinktestlegacy

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/kitchensinktestlegacy/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type denseListNullEndpointResolver struct{}

func (*denseListNullEndpointResolver) ResolveEndpoint(ctx context.Context, params EndpointParameters) (smithyendpoints.Endpoint, error) {
	return smithyendpoints.Endpoint{URI: url.URL{Scheme: "https", Host: "test.example.com"}}, nil
}

// An explicit null in a dense list is not valid per the spec, but we tolerate
// it by keeping the element with the zero value.
func TestDenseListNullElement(t *testing.T) {
	body := `{
		"strings": ["x", null],
		"integers": [1, null],
		"structs": [{"name": "x"}, null],
		"sparseStrings": ["x", null]
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
		EndpointResolverV2: &denseListNullEndpointResolver{},
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				s.Finalize.Clear()
				return nil
			},
		},
	})

	out, err := svc.GetLists(context.Background(), &GetListsInput{})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if expect := []string{"x", ""}; !reflect.DeepEqual(expect, out.Strings) {
		t.Errorf("strings: expect %v, got %v", expect, out.Strings)
	}
	if expect := []int32{1, 0}; !reflect.DeepEqual(expect, out.Integers) {
		t.Errorf("integers: expect %v, got %v", expect, out.Integers)
	}
	if expect := []types.ListValueStruct{{Name: aws.String("x")}, {}}; !reflect.DeepEqual(expect, out.Structs) {
		t.Errorf("structs: expect %v, got %v", expect, out.Structs)
	}
	if expect := []*string{aws.String("x"), nil}; !reflect.DeepEqual(expect, out.SparseStrings) {
		t.Errorf("sparseStrings: expect %v, got %v", expect, out.SparseStrings)
	}
}
