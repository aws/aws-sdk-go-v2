package manager

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/manager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var mockErrResponse = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
    <Code>MockCode</Code>
    <Message>The error message</Message>
    <RequestId>4442587FB7D0A2F9</RequestId>
</Error>`)

func testSetupGetBucketRegionServer(region string, statusCode int, incHeader bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(ioutil.Discard, r.Body)
		if incHeader {
			w.Header().Set(bucketRegionHeader, region)
		}
		if statusCode >= 300 {
			w.Header().Set("Content-Length", strconv.Itoa(len(mockErrResponse)))
			w.WriteHeader(statusCode)
			w.Write(mockErrResponse)
		} else {
			w.WriteHeader(statusCode)
		}
	}))
}

var testGetBucketRegionCases = []struct {
	RespRegion      string
	StatusCode      int
	ExpectReqRegion string
}{
	{
		RespRegion: "bucket-region",
		StatusCode: 301,
	},
	{
		RespRegion: "bucket-region",
		StatusCode: 403,
	},
	{
		RespRegion: "bucket-region",
		StatusCode: 200,
	},
	{
		RespRegion:      "bucket-region",
		StatusCode:      200,
		ExpectReqRegion: "default-region",
	},
}

func TestGetBucketRegion_Exists(t *testing.T) {
	for i, c := range testGetBucketRegionCases {
		server := testSetupGetBucketRegionServer(c.RespRegion, c.StatusCode, true)

		client := s3.New(s3.Options{
			EndpointResolver: s3testing.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: server.URL,
				}, nil
			}),
		})

		region, err := GetBucketRegion(context.Background(), client, "bucket", func(o *s3.Options) {
			o.UsePathStyle = true
		})
		if err != nil {
			t.Errorf("%d, expect no error, got %v", i, err)
			goto closeServer
		}
		if e, a := c.RespRegion, region; e != a {
			t.Errorf("%d, expect %q region, got %q", i, e, a)
		}

	closeServer:
		server.Close()
	}
}

func TestGetBucketRegion_NotExists(t *testing.T) {
	server := testSetupGetBucketRegionServer("ignore-region", 404, false)
	defer server.Close()

	client := s3.New(s3.Options{
		EndpointResolver: s3testing.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: server.URL,
			}, nil
		}),
	})

	region, err := GetBucketRegion(context.Background(), client, "bucket", func(o *s3.Options) {
		o.UsePathStyle = true
	})
	if err == nil {
		t.Fatalf("expect error, but did not get one")
	}

	var bnf BucketNotFound
	if !errors.As(err, &bnf) {
		t.Errorf("expect %T error, got %v", bnf, err)
	}

	if len(region) != 0 {
		t.Errorf("expect region not to be set, got %q", region)
	}
}
