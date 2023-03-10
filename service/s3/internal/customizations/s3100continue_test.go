package customizations_test

import (
	"context"
	"crypto/tls"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

// unit test for service/internal/s3shared/s3100continue.go
func TestAdd100ContinueHttpHeader(t *testing.T) {
	const HeaderKey = "Expect"
	HeaderValue := "100-continue"
	testBucket := "testBucket"
	testKey := "testKey"

	cases := map[string]struct {
		Handler          func(*testing.T) http.Handler
		Input            interface{}
		ExpectValueFound string
	}{
		"http put request smaller than 2MB": {
			Handler: func(t *testing.T) http.Handler {
				return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					if diff := cmp.Diff(request.Header.Get(HeaderKey), ""); len(diff) > 0 {
						t.Error(diff)
					}

					writer.WriteHeader(200)
				})
			},
			Input: &s3.PutObjectInput{
				Bucket:        &testBucket,
				Key:           &testKey,
				ContentLength: 1,
				Body:          &awstesting.ReadCloser{Size: 1},
			},
			ExpectValueFound: "",
		},
		"http put request larger than 2MB": {
			Handler: func(t *testing.T) http.Handler {
				return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					if diff := cmp.Diff(request.Header.Get(HeaderKey), HeaderValue); len(diff) > 0 {
						t.Error(diff)
					}

					writer.WriteHeader(200)
				})
			},
			Input: &s3.PutObjectInput{
				Bucket:        &testBucket,
				Key:           &testKey,
				ContentLength: 1024 * 1024 * 3,
				Body:          &awstesting.ReadCloser{Size: 1024 * 1024 * 3},
			},
			ExpectValueFound: HeaderValue,
		},
		"http put request with unknown -1 ContentLength": {
			Handler: func(t *testing.T) http.Handler {
				return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					if diff := cmp.Diff(request.Header.Get(HeaderKey), HeaderValue); len(diff) > 0 {
						t.Error(diff)
					}

					writer.WriteHeader(200)
				})
			},
			Input: &s3.PutObjectInput{
				Bucket:        &testBucket,
				Key:           &testKey,
				ContentLength: -1,
				Body:          &awstesting.ReadCloser{Size: 1024 * 1024 * 10},
			},
			ExpectValueFound: HeaderValue,
		},
		"http put request with 0 ContentLength but unknown non-nil body": {
			Handler: func(t *testing.T) http.Handler {
				return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					if diff := cmp.Diff(request.Header.Get(HeaderKey), HeaderValue); len(diff) > 0 {
						t.Error(diff)
					}

					writer.WriteHeader(200)
				})
			},
			Input: &s3.PutObjectInput{
				Bucket:        &testBucket,
				Key:           &testKey,
				ContentLength: 0,
				Body:          &awstesting.ReadCloser{Size: 1024 * 1024 * 3},
			},
			ExpectValueFound: HeaderValue,
		},
		"http get request": {
			Handler: func(t *testing.T) http.Handler {
				return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					if diff := cmp.Diff(request.Header.Get(HeaderKey), ""); len(diff) > 0 {
						t.Error(diff)
					}

					writer.WriteHeader(200)
				})
			},
			Input: &s3.GetObjectInput{
				Bucket: &testBucket,
				Key:    &testKey,
			},
			ExpectValueFound: "",
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewTLSServer(tt.Handler(t))
			defer server.Close()
			client := s3.New(s3.Options{
				Region: "us-west-2",
				HTTPClient: &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				},
				EndpointResolver: s3.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               server.URL,
						SigningName:       "s3-object-lambda",
						SigningRegion:     region,
						Source:            aws.EndpointSourceCustom,
						HostnameImmutable: true,
					}, nil
				}),
			})

			switch tt.Input.(type) {
			case *s3.PutObjectInput:
				_, err := client.PutObject(context.Background(), tt.Input.(*s3.PutObjectInput))
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
			default:
				_, err := client.GetObject(context.Background(), tt.Input.(*s3.GetObjectInput))
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
			}
		})
	}
}
