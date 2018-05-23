package s3_test

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3BucketTest struct {
	bucket  string
	url     string
	errCode string
}

var (
	sslTests = []s3BucketTest{
		{"abc", "https://abc.s3.mock-region.amazonaws.com/", ""},
		{"a$b$c", "https://s3.mock-region.amazonaws.com/a%24b%24c", ""},
		{"a.b.c", "https://s3.mock-region.amazonaws.com/a.b.c", ""},
		{"a..bc", "https://s3.mock-region.amazonaws.com/a..bc", ""},
	}

	nosslTests = []s3BucketTest{
		{"a.b.c", "http://a.b.c.s3.mock-region.amazonaws.com/", ""},
		{"a..bc", "http://s3.mock-region.amazonaws.com/a..bc", ""},
	}

	forcepathTests = []s3BucketTest{
		{"abc", "https://s3.mock-region.amazonaws.com/abc", ""},
		{"a$b$c", "https://s3.mock-region.amazonaws.com/a%24b%24c", ""},
		{"a.b.c", "https://s3.mock-region.amazonaws.com/a.b.c", ""},
		{"a..bc", "https://s3.mock-region.amazonaws.com/a..bc", ""},
	}

	accelerateTests = []s3BucketTest{
		{"abc", "https://abc.s3-accelerate.amazonaws.com/", ""},
		{"a.b.c", "https://s3.mock-region.amazonaws.com/%7BBucket%7D", "InvalidParameterException"},
		{"a$b$c", "https://s3.mock-region.amazonaws.com/%7BBucket%7D", "InvalidParameterException"},
	}

	accelerateNoSSLTests = []s3BucketTest{
		{"abc", "http://abc.s3-accelerate.amazonaws.com/", ""},
		{"a.b.c", "http://a.b.c.s3-accelerate.amazonaws.com/", ""},
		{"a$b$c", "http://s3.mock-region.amazonaws.com/%7BBucket%7D", "InvalidParameterException"},
	}

	accelerateDualstack = []s3BucketTest{
		{"abc", "https://abc.s3-accelerate.dualstack.amazonaws.com/", ""},
		{"a.b.c", "https://s3.dualstack.mock-region.amazonaws.com/%7BBucket%7D", "InvalidParameterException"},
		{"a$b$c", "https://s3.dualstack.mock-region.amazonaws.com/%7BBucket%7D", "InvalidParameterException"},
	}
)

func runTests(t *testing.T, svc *s3.S3, tests []s3BucketTest) {
	t.Helper()

	for i, test := range tests {
		req := svc.ListObjectsRequest(&s3.ListObjectsInput{Bucket: &test.bucket})
		req.Build()

		if test.errCode != "" {
			if err := req.Error; err == nil {
				t.Fatalf("%d, expect error, got none", i)
			}
			if a, e := req.Error.(awserr.Error).Code(), test.errCode; !strings.Contains(a, e) {
				t.Errorf("%d, expect error code to contain %q, got %q", i, e, a)
			}
		} else {
			if err := req.Error; err != nil {
				t.Fatalf("%d, expect no error, got %v", i, err)
			}
			if e, a := test.url, req.HTTPRequest.URL.String(); e != a {
				t.Errorf("%d, expect url %s, got %s", i, e, a)
			}
		}
	}
}

func TestAccelerateBucketBuild(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	s := s3.New(cfg)
	s.UseAccelerate = true
	runTests(t, s, accelerateTests)
}

func TestAccelerateNoSSLBucketBuild(t *testing.T) {
	cfg := unit.Config()
	resolver := endpoints.NewDefaultResolver()
	resolver.DisableSSL = true
	cfg.EndpointResolver = resolver

	s := s3.New(cfg)
	s.UseAccelerate = true
	runTests(t, s, accelerateNoSSLTests)
}

func TestAccelerateDualstackBucketBuild(t *testing.T) {
	cfg := unit.Config()
	resolver := endpoints.NewDefaultResolver()
	resolver.UseDualStack = true
	cfg.EndpointResolver = resolver

	s := s3.New(cfg)
	s.UseAccelerate = true
	runTests(t, s, accelerateDualstack)
}

func TestHostStyleBucketBuild(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	s := s3.New(cfg)
	runTests(t, s, sslTests)
}

func TestHostStyleBucketBuildNoSSL(t *testing.T) {
	cfg := unit.Config()
	resolver := endpoints.NewDefaultResolver()
	resolver.DisableSSL = true
	cfg.EndpointResolver = resolver

	s := s3.New(cfg)
	runTests(t, s, nosslTests)
}

func TestPathStyleBucketBuild(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	s := s3.New(cfg)
	s.ForcePathStyle = true
	runTests(t, s, forcepathTests)
}

func TestHostStyleBucketGetBucketLocation(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	s := s3.New(cfg)
	req := s.GetBucketLocationRequest(&s3.GetBucketLocationInput{
		Bucket: aws.String("bucket"),
	})

	req.Build()
	if req.Error != nil {
		t.Fatalf("expect no error, got %v", req.Error)
	}
	u, _ := url.Parse(req.HTTPRequest.URL.String())
	if e, a := "bucket", u.Host; strings.Contains(a, e) {
		t.Errorf("expect %s to not be in %s", e, a)
	}
	if e, a := "bucket", u.Path; !strings.Contains(a, e) {
		t.Errorf("expect %s to be in %s", e, a)
	}
}
func TestVirtualHostStyleSuite(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "virtual_host.json"))
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}

	cases := []struct {
		Bucket                    string
		Region                    string
		UseDualStack              bool
		UseS3Accelerate           bool
		ConfiguredAddressingStyle string

		ExpectedURI string
	}{}

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&cases); err != nil {
		t.Fatalf("expect no error, %v", err)
	}

	const testPathStyle = "path"
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cfg := unit.Config()
			resolver := endpoints.NewDefaultResolver()
			resolver.UseDualStack = c.UseDualStack
			cfg.EndpointResolver = resolver
			cfg.Region = c.Region

			svc := s3.New(cfg)
			svc.ForcePathStyle = c.ConfiguredAddressingStyle == testPathStyle
			svc.UseAccelerate = c.UseS3Accelerate

			req := svc.HeadBucketRequest(&s3.HeadBucketInput{
				Bucket: &c.Bucket,
			})
			req.Build()
			if req.Error != nil {
				t.Fatalf("expect no error, got %v", req.Error)
			}

			// Trim trailing '/' that are added by the SDK but not in the tests.
			actualURI := strings.TrimRightFunc(
				req.HTTPRequest.URL.String(),
				func(r rune) bool { return r == '/' },
			)
			if e, a := c.ExpectedURI, actualURI; e != a {
				t.Errorf("expect\n%s\nurl to be\n%s", e, a)
			}
		})
	}
}
