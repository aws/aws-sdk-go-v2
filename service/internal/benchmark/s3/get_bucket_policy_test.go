package s3

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func init() {
	var s strings.Builder
	io.Copy(&s, io.LimitReader(byteReader('a'), 1024*1024*12))
	largeStringPayload = s.String()
}

type byteReader byte

func (b byteReader) Read(p []byte) (int, error) {
	for i := 0; i < len(p); i++ {
		p[i] = byte(b)
	}
	return len(p), nil
}

var largeStringPayload string

func BenchmarkGetBucketPolicy(b *testing.B) {
	client := s3.New(s3.Options{
		Region: "us-west-2",
		HTTPClient: smithyhttp.ClientDoFunc(
			func(r *http.Request) (*http.Response, error) {
				return newGetBucketPolicyHTTPResponse(largeStringPayload), nil
			}),
	})

	ctx := context.Background()
	params := s3.GetBucketPolicyInput{
		Bucket: aws.String("fooBucket"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetBucketPolicy(ctx, &params)
		if err != nil {
			b.Fatalf("failed to send: %v", err)
		}
	}

}

func newGetBucketPolicyHTTPResponse(payload string) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		ContentLength: int64(len(payload)),
		Body:          ioutil.NopCloser(strings.NewReader(payload)),
	}
}
