package s3

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// This file replicates the tests in https://github.com/smithy-lang/smithy/blob/main/smithy-aws-protocol-tests/model/restXml/services/s3.smithy,
// which we cannot generate through normal protocoltest codegen due to
// requirement on handwritten source in S3.

type capturedRequest struct {
	r *http.Request
}

func (cr *capturedRequest) Do(r *http.Request) (*http.Response, error) {
	cr.r = r
	return &http.Response{ // returns are moot, for request tests only
		StatusCode: 400,
		Body:       http.NoBody,
	}, nil
}

func TestS3Protocol_ListObjectsV2_Request(t *testing.T) {
	for name, tt := range map[string]struct {
		Options          func(*Options)
		OperationOptions func(*Options)
		Input            *ListObjectsV2Input

		ExpectMethod string
		ExpectHost   string
		ExpectPath   string
		ExpectQuery  []string
	}{
		"S3DefaultAddressing": {
			Options: func(o *Options) {
				o.Region = "us-west-2"
			},
			OperationOptions: func(o *Options) {},
			Input: &ListObjectsV2Input{
				Bucket: aws.String("mybucket"),
			},

			ExpectMethod: "GET",
			ExpectHost:   "mybucket.s3.us-west-2.amazonaws.com",
			ExpectPath:   "/",
			ExpectQuery:  []string{"list-type=2"},
		},
		"S3VirtualHostAddressing": {
			Options: func(o *Options) {
				o.Region = "us-west-2"
				o.UsePathStyle = false
			},
			OperationOptions: func(o *Options) {},
			Input: &ListObjectsV2Input{
				Bucket: aws.String("mybucket"),
			},

			ExpectMethod: "GET",
			ExpectHost:   "mybucket.s3.us-west-2.amazonaws.com",
			ExpectPath:   "/",
			ExpectQuery:  []string{"list-type=2"},
		},
		"S3PathAddressing": {
			Options: func(o *Options) {
				o.Region = "us-west-2"
				o.UsePathStyle = true
			},
			OperationOptions: func(o *Options) {},
			Input: &ListObjectsV2Input{
				Bucket: aws.String("mybucket"),
			},

			ExpectMethod: "GET",
			ExpectHost:   "s3.us-west-2.amazonaws.com",
			ExpectPath:   "/mybucket",
			ExpectQuery:  []string{"list-type=2"},
		},
		"S3VirtualHostDualstackAddressing": {
			Options: func(o *Options) {
				o.Region = "us-west-2"
				o.UsePathStyle = false
				o.EndpointOptions.UseDualStackEndpoint = aws.DualStackEndpointStateEnabled
			},
			OperationOptions: func(o *Options) {},
			Input: &ListObjectsV2Input{
				Bucket: aws.String("mybucket"),
			},

			ExpectMethod: "GET",
			ExpectHost:   "mybucket.s3.dualstack.us-west-2.amazonaws.com",
			ExpectPath:   "/",
			ExpectQuery:  []string{"list-type=2"},
		},
		"S3VirtualHostAccelerateAddressing": {
			Options: func(o *Options) {
				o.Region = "us-west-2"
				o.UsePathStyle = false
				o.UseAccelerate = true
			},
			OperationOptions: func(o *Options) {},
			Input: &ListObjectsV2Input{
				Bucket: aws.String("mybucket"),
			},

			ExpectMethod: "GET",
			ExpectHost:   "mybucket.s3-accelerate.amazonaws.com",
			ExpectPath:   "/",
			ExpectQuery:  []string{"list-type=2"},
		},
		"S3VirtualHostDualstackAccelerateAddressing": {
			Options: func(o *Options) {
				o.Region = "us-west-2"
				o.UsePathStyle = false
				o.EndpointOptions.UseDualStackEndpoint = aws.DualStackEndpointStateEnabled
				o.UseAccelerate = true
			},
			OperationOptions: func(o *Options) {},
			Input: &ListObjectsV2Input{
				Bucket: aws.String("mybucket"),
			},

			ExpectMethod: "GET",
			ExpectHost:   "mybucket.s3-accelerate.dualstack.amazonaws.com",
			ExpectPath:   "/",
			ExpectQuery:  []string{"list-type=2"},
		},
		"S3OperationAddressingPreferred": {
			Options: func(o *Options) {
				o.Region = "us-west-2"
				o.UsePathStyle = true
			},
			OperationOptions: func(o *Options) {
				o.UsePathStyle = false
			},
			Input: &ListObjectsV2Input{
				Bucket: aws.String("mybucket"),
			},

			ExpectMethod: "GET",
			ExpectHost:   "mybucket.s3.us-west-2.amazonaws.com",
			ExpectPath:   "/",
			ExpectQuery:  []string{"list-type=2"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			var r capturedRequest
			svc := New(Options{HTTPClient: &r}, tt.Options)

			svc.ListObjectsV2(context.Background(), tt.Input, tt.OperationOptions)
			if r.r == nil {
				t.Fatal("captured request is nil")
			}

			req := r.r
			if tt.ExpectMethod != req.Method {
				t.Errorf("expect method: %v != %v", tt.ExpectMethod, req.Method)
			}
			if tt.ExpectHost != req.URL.Host {
				t.Errorf("expect host: %v != %v", tt.ExpectHost, req.URL.Host)
			}
			if tt.ExpectPath != req.URL.RawPath {
				t.Errorf("expect path: %v != %v", tt.ExpectPath, req.URL.RawPath)
			}
			for _, q := range tt.ExpectQuery {
				if !strings.Contains(req.URL.RawQuery, q) {
					t.Errorf("query %v is missing %v", req.URL.RawQuery, q)
				}
			}
		})
	}
}

func TestS3Protocol_DeleteObjectTagging_Request(t *testing.T) {
	for name, tt := range map[string]struct {
		ClientOptions    func(*Options)
		OperationOptions func(*Options)
		Input            *DeleteObjectTaggingInput

		ExpectMethod string
		ExpectHost   string
		ExpectPath   string
		ExpectQuery  []string
	}{
		"S3EscapeObjectKeyInUriLabel": {
			ClientOptions: func(o *Options) {
				o.Region = "us-west-2"
			},
			OperationOptions: func(o *Options) {},
			Input: &DeleteObjectTaggingInput{
				Bucket: aws.String("mybucket"),
				Key:    aws.String("my key.txt"),
			},

			ExpectMethod: "DELETE",
			ExpectHost:   "mybucket.s3.us-west-2.amazonaws.com",
			ExpectPath:   "/my%20key.txt",
			ExpectQuery:  []string{"tagging"},
		},
		"S3EscapePathObjectKeyInUriLabel": {
			ClientOptions: func(o *Options) {
				o.Region = "us-west-2"
			},
			OperationOptions: func(o *Options) {},
			Input: &DeleteObjectTaggingInput{
				Bucket: aws.String("mybucket"),
				Key:    aws.String("foo/bar/my key.txt"),
			},

			ExpectMethod: "DELETE",
			ExpectHost:   "mybucket.s3.us-west-2.amazonaws.com",
			ExpectPath:   "/foo/bar/my%20key.txt",
			ExpectQuery:  []string{"tagging"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			var r capturedRequest
			svc := New(Options{HTTPClient: &r}, tt.ClientOptions)

			svc.DeleteObjectTagging(context.Background(), tt.Input, tt.OperationOptions)
			if r.r == nil {
				t.Fatal("captured request is nil")
			}

			req := r.r
			if tt.ExpectMethod != req.Method {
				t.Errorf("expect method: %v != %v", tt.ExpectMethod, req.Method)
			}
			if tt.ExpectHost != req.URL.Host {
				t.Errorf("expect host: %v != %v", tt.ExpectHost, req.URL.Host)
			}
			if tt.ExpectPath != req.URL.RawPath {
				t.Errorf("expect path: %v != %v", tt.ExpectPath, req.URL.RawPath)
			}
			for _, q := range tt.ExpectQuery {
				if !strings.Contains(req.URL.RawQuery, q) {
					t.Errorf("query %v is missing %v", req.URL.RawQuery, q)
				}
			}
		})
	}

}

func TestS3Protocol_GetObject_Request(t *testing.T) {
	for name, tt := range map[string]struct {
		ClientOptions    func(*Options)
		OperationOptions func(*Options)
		Input            *GetObjectInput

		ExpectMethod string
		ExpectHost   string
		ExpectPath   string
		ExpectQuery  []string
	}{
		"S3PreservesLeadingDotSegmentInUriLabel": {
			ClientOptions: func(o *Options) {
				o.Region = "us-west-2"
				o.UsePathStyle = false
			},
			OperationOptions: func(o *Options) {},
			Input: &GetObjectInput{
				Bucket: aws.String("mybucket"),
				Key:    aws.String("../key.txt"),
			},

			ExpectMethod: "GET",
			ExpectHost:   "mybucket.s3.us-west-2.amazonaws.com",
			ExpectPath:   "/../key.txt",
		},
		"S3PreservesEmbeddedDotSegmentInUriLabel": {
			ClientOptions: func(o *Options) {
				o.Region = "us-west-2"
				o.UsePathStyle = false
			},
			OperationOptions: func(o *Options) {},
			Input: &GetObjectInput{
				Bucket: aws.String("mybucket"),
				Key:    aws.String("foo/../key.txt"),
			},

			ExpectMethod: "GET",
			ExpectHost:   "mybucket.s3.us-west-2.amazonaws.com",
			ExpectPath:   "/foo/../key.txt",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var r capturedRequest
			svc := New(Options{HTTPClient: &r}, tt.ClientOptions)

			svc.GetObject(context.Background(), tt.Input, tt.OperationOptions)
			if r.r == nil {
				t.Fatal("captured request is nil")
			}

			req := r.r
			if tt.ExpectMethod != req.Method {
				t.Errorf("expect method: %v != %v", tt.ExpectMethod, req.Method)
			}
			if tt.ExpectHost != req.URL.Host {
				t.Errorf("expect host: %v != %v", tt.ExpectHost, req.URL.Host)
			}
			if tt.ExpectPath != req.URL.RawPath {
				t.Errorf("expect path: %v != %v", tt.ExpectPath, req.URL.RawPath)
			}
			for _, q := range tt.ExpectQuery {
				if !strings.Contains(req.URL.RawQuery, q) {
					t.Errorf("query %v is missing %v", req.URL.RawQuery, q)
				}
			}
		})
	}

}

type mockHTTPResponse struct {
	resp *http.Response
}

func (m *mockHTTPResponse) Do(r *http.Request) (*http.Response, error) {
	return m.resp, nil
}

func TestS3Protocol_GetBucketLocation_Response(t *testing.T) {
	for name, tt := range map[string]struct {
		Response *http.Response
		Expect   *GetBucketLocationOutput
	}{
		"GetBucketLocationUnwrappedOutput": {
			Response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<LocationConstraint xmlns=\"http://s3.amazonaws.com/doc/2006-03-01/\">us-west-2</LocationConstraint>")),
			},
			Expect: &GetBucketLocationOutput{
				LocationConstraint: types.BucketLocationConstraintUsWest2,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			svc := New(Options{
				Region:     "us-west-2",
				HTTPClient: &mockHTTPResponse{tt.Response},
			})

			out, err := svc.GetBucketLocation(context.Background(), &GetBucketLocationInput{
				Bucket: aws.String("bucket"),
			})
			if err != nil {
				t.Fatalf("get bucket location: %v", err)
			}

			if tt.Expect.LocationConstraint != out.LocationConstraint {
				t.Errorf("LocationConstraint %v != %v", tt.Expect.LocationConstraint, out.LocationConstraint)
			}
		})
	}
}

func TestS3Protocol_Error_NoSuchBucket(t *testing.T) {
	for name, tt := range map[string]struct {
		Response *http.Response
	}{
		"GetBucketLocationUnwrappedOutput": {
			Response: &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(strings.NewReader("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Error>\n\t<Type>Sender</Type>\n\t<Code>NoSuchBucket</Code>\n</Error>")),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			svc := New(Options{
				Region:     "us-west-2",
				HTTPClient: &mockHTTPResponse{tt.Response},
			})

			_, err := svc.GetObject(context.Background(), &GetObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			})
			if err == nil {
				t.Fatal("call operation: expected error, got none")
			}

			// of note: we don't actually return a *types.NoSuchBucket in this
			// case, but we DO capture the right error code
			var terr interface {
				ErrorCode() string
			}
			if !errors.As(err, &terr) {
				t.Errorf("error does not implement ErrorCode(), was %v", err)
			}
			if actual := terr.ErrorCode(); actual != "NoSuchBucket" {
				t.Errorf("error code, expected NoSuchBucket, was %v", actual)
			}
		})
	}

}
