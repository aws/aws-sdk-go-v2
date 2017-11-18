package s3manager

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func testSetupGetBucketRegionServer(region string, statusCode int, incHeader bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if incHeader {
			w.Header().Set(bucketRegionHeader, region)
		}
		w.WriteHeader(statusCode)
	}))
}

var testGetBucketRegionCases = []struct {
	RespRegion string
	StatusCode int
}{
	{"bucket-region", 301},
	{"bucket-region", 403},
	{"bucket-region", 200},
}

func TestGetBucketRegion_Exists(t *testing.T) {
	for i, c := range testGetBucketRegionCases {
		server := testSetupGetBucketRegionServer(c.RespRegion, c.StatusCode, true)

		cfg := unit.Config()
		cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

		ctx := aws.BackgroundContext()
		region, err := GetBucketRegion(ctx, cfg, "bucket", "region")
		if err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}
		if e, a := c.RespRegion, region; e != a {
			t.Errorf("%d, expect %q region, got %q", i, e, a)
		}
	}
}

func TestGetBucketRegion_NotExists(t *testing.T) {
	server := testSetupGetBucketRegionServer("ignore-region", 404, false)

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	ctx := aws.BackgroundContext()
	region, err := GetBucketRegion(ctx, cfg, "bucket", "region")
	if err == nil {
		t.Fatalf("expect error, but did not get one")
	}
	aerr := err.(awserr.Error)
	if e, a := "NotFound", aerr.Code(); e != a {
		t.Errorf("expect %s error code, got %s", e, a)
	}
	if len(region) != 0 {
		t.Errorf("expect region not to be set, got %q", region)
	}
}

func TestGetBucketRegionWithClient(t *testing.T) {
	for i, c := range testGetBucketRegionCases {
		server := testSetupGetBucketRegionServer(c.RespRegion, c.StatusCode, true)

		cfg := unit.Config()
		cfg.Region = "region"
		cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

		svc := s3.New(cfg)
		svc.ForcePathStyle = true

		ctx := aws.BackgroundContext()

		region, err := GetBucketRegionWithClient(ctx, svc, "bucket")
		if err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}
		if e, a := c.RespRegion, region; e != a {
			t.Errorf("%d, expect %q region, got %q", i, e, a)
		}
	}
}
