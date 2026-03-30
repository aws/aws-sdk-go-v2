package schemas

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/schemas"
	"github.com/aws/smithy-go/ptr"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func init() {
	largeBytePayload = make([]byte, 1024*1024*12)
}

var largeBytePayload []byte

func BenchmarkGetCodeBindingSource(b *testing.B) {
	client := schemas.New(schemas.Options{
		Region: "us-west-2",
		HTTPClient: smithyhttp.ClientDoFunc(
			func(r *http.Request) (*http.Response, error) {
				return newGetCodeBindingSourceHTTPResponse(largeBytePayload), nil
			}),
	})

	ctx := context.Background()
	params := schemas.GetCodeBindingSourceInput{
		Language:     ptr.String("fooLanguage"),
		RegistryName: ptr.String("fooRegistry"),
		SchemaName:   ptr.String("fooSchema"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetCodeBindingSource(ctx, &params)
		if err != nil {
			b.Fatalf("failed to send: %v", err)
		}
	}

}

func newGetCodeBindingSourceHTTPResponse(payload []byte) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Header: map[string][]string{
			"Content-Type": {"application/octet-stream"},
		},
		ContentLength: int64(len(payload)),
		Body:          ioutil.NopCloser(bytes.NewReader(payload)),
	}
}
