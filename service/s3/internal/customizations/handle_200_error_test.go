package customizations_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type EndpointResolverFunc func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error)

func (fn EndpointResolverFunc) ResolveEndpoint(region string, options s3.EndpointResolverOptions) (endpoint aws.Endpoint, err error) {
	return fn(region, options)
}

func TestErrorResponseWith200StatusCode(t *testing.T) {
	cases := map[string]struct {
		response       []byte
		statusCode     int
		expectedError  string
		expectedBucket string
	}{
		"200ErrorBody": {
			response: []byte(`<Error><Type>Sender</Type>
    <Code>InvalidGreeting</Code>
    <Message>Hi</Message>
    <AnotherSetting>setting</AnotherSetting>
    <RequestId>foo-id</RequestId></Error>`),
			statusCode:    200,
			expectedError: "InvalidGreeting",
		},
		"200NoResponse": {
			response:      []byte{},
			statusCode:    200,
			expectedError: "received empty response payload",
		},
		"200InvalidResponse": {
			response: []byte(`<Error><Type>Sender</Type>
    <Code>InvalidGreeting</Code>
    <Message>Hi</Message>
    <AnotherSetting>setting</AnotherSetting>
    <RequestId>foo-id`),
			statusCode:    200,
			expectedError: "unexpected EOF",
		},
		"200SuccessResponse": {
			response: []byte(`<CompleteMultipartUploadResult>
   			<Bucket>bucket</Bucket>
			</CompleteMultipartUploadResult>`),
			statusCode:     200,
			expectedError:  "",
			expectedBucket: "bucket",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(c.statusCode)
					w.Write(c.response)
				}))
			defer server.Close()

			options := s3.Options{
				Credentials: unit.StubCredentialsProvider{},
				Retryer:     aws.NopRetryer{},
				Region:      "mock-region",
				EndpointResolver: EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (e aws.Endpoint, err error) {
					e.URL = server.URL
					e.SigningRegion = "us-west-2"
					return e, err
				}),
				UsePathStyle: true,
			}

			svc := s3.New(options)
			resp, err := svc.CompleteMultipartUpload(context.Background(), &s3.CompleteMultipartUploadInput{
				UploadId:     aws.String("mockID"),
				RequestPayer: "requester",
				Bucket:       aws.String("bucket"),
				Key:          aws.String("mockKey"),
			})

			if len(c.expectedError) != 0 {
				if err == nil {
					t.Fatalf("expected error, got none")
				}

				if e, a := c.expectedError, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expected %v, got %v", e, a)
				}
			}

			if len(c.expectedBucket) != 0 {
				if e, a := c.expectedBucket, *resp.Bucket; !strings.EqualFold(e, a) {
					t.Fatalf("expected bucket name to be %v, got %v", e, a)
				}
			}
		})
	}
}
