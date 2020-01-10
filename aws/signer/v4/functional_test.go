package v4_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var standaloneSignCases = []struct {
	OrigURI                    string
	OrigQuery                  string
	Region, Service, SubDomain string
	ExpSig                     string
	EscapedURI                 string
}{
	{
		OrigURI:   `/logs-*/_search`,
		OrigQuery: `pretty=true`,
		Region:    "us-west-2", Service: "es", SubDomain: "hostname-clusterkey",
		EscapedURI: `/logs-%2A/_search`,
		ExpSig:     `AWS4-HMAC-SHA256 Credential=AKID/19700101/us-west-2/es/aws4_request, SignedHeaders=host;x-amz-date;x-amz-security-token, Signature=79d0760751907af16f64a537c1242416dacf51204a7dd5284492d15577973b91`,
	},
}

func TestPresignHandler(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	svc := s3.New(cfg)
	req := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:             aws.String("bucket"),
		Key:                aws.String("key"),
		ContentDisposition: aws.String("a+b c$d"),
		ACL:                s3.ObjectCannedACLPublicRead,
	})
	req.Time = time.Unix(0, 0)
	urlstr, err := req.Presign(5 * time.Minute)

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	expectedHost := "bucket.s3.mock-region.amazonaws.com"
	expectedDate := "19700101T000000Z"
	expectedHeaders := "content-disposition;host;x-amz-acl"
	expectedSig := "2d76a414208c0eac2a23ef9c834db9635ecd5a0fbb447a00ad191f82d854f55b"
	expectedCred := "AKID/19700101/mock-region/s3/aws4_request"

	u, _ := url.Parse(urlstr)
	urlQ := u.Query()
	if e, a := expectedHost, u.Host; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedSig, urlQ.Get("X-Amz-Signature"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedCred, urlQ.Get("X-Amz-Credential"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedHeaders, urlQ.Get("X-Amz-SignedHeaders"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedDate, urlQ.Get("X-Amz-Date"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "300", urlQ.Get("X-Amz-Expires"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if a := urlQ.Get("X-Amz-Content-Sha256"); len(a) != 0 {
		t.Errorf("expect no content sha256 got %v", a)
	}

	if e, a := "+", urlstr; strings.Contains(a, e) { // + encoded as %20
		t.Errorf("expect %v not to be in %v", e, a)
	}
}

func TestPresignRequest(t *testing.T) {
	cfg := unit.Config()
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	svc := s3.New(cfg)
	req := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:             aws.String("bucket"),
		Key:                aws.String("key"),
		ContentDisposition: aws.String("a+b c$d"),
		ACL:                s3.ObjectCannedACLPublicRead,
	})
	req.Time = time.Unix(0, 0)
	urlstr, headers, err := req.PresignRequest(5 * time.Minute)

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	expectedHost := "bucket.s3.mock-region.amazonaws.com"
	expectedDate := "19700101T000000Z"
	expectedHeaders := "content-disposition;host;x-amz-acl"
	expectedSig := "2d76a414208c0eac2a23ef9c834db9635ecd5a0fbb447a00ad191f82d854f55b"
	expectedCred := "AKID/19700101/mock-region/s3/aws4_request"
	expectedHeaderMap := http.Header{
		"x-amz-acl":           []string{"public-read"},
		"content-disposition": []string{"a+b c$d"},
	}

	u, _ := url.Parse(urlstr)
	urlQ := u.Query()
	if e, a := expectedHost, u.Host; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedSig, urlQ.Get("X-Amz-Signature"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedCred, urlQ.Get("X-Amz-Credential"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedHeaders, urlQ.Get("X-Amz-SignedHeaders"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedDate, urlQ.Get("X-Amz-Date"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedHeaderMap, headers; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "300", urlQ.Get("X-Amz-Expires"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if a := urlQ.Get("X-Amz-Content-Sha256"); len(a) != 0 {
		t.Errorf("expect no content sha256 got %v", a)
	}

	if e, a := "+", urlstr; strings.Contains(a, e) { // + encoded as %20
		t.Errorf("expect %v not to be in %v", e, a)
	}
}

func TestStandaloneSign_CustomURIEscape(t *testing.T) {
	var expectSig = `AWS4-HMAC-SHA256 Credential=AKID/19700101/us-east-1/es/aws4_request, SignedHeaders=host;x-amz-date;x-amz-security-token, Signature=6601e883cc6d23871fd6c2a394c5677ea2b8c82b04a6446786d64cd74f520967`

	creds := unit.Config().Credentials
	signer := v4.NewSigner(creds, func(s *v4.Signer) {
		s.DisableURIPathEscaping = true
	})

	host := "https://subdomain.us-east-1.es.amazonaws.com"
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	req.URL.Path = `/log-*/_search`
	req.URL.Opaque = "//subdomain.us-east-1.es.amazonaws.com/log-%2A/_search"

	_, err = signer.Sign(context.Background(), req, nil, "es", "us-east-1", time.Unix(0, 0))
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	actual := req.Header.Get("Authorization")
	if e, a := expectSig, actual; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestStandaloneSign(t *testing.T) {
	creds := unit.Config().Credentials
	signer := v4.NewSigner(creds)

	for _, c := range standaloneSignCases {
		host := fmt.Sprintf("https://%s.%s.%s.amazonaws.com",
			c.SubDomain, c.Region, c.Service)

		req, err := http.NewRequest("GET", host, nil)
		if err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		// URL.EscapedPath() will be used by the signer to get the
		// escaped form of the request's URI path.
		req.URL.Path = c.OrigURI
		req.URL.RawQuery = c.OrigQuery

		_, err = signer.Sign(context.Background(), req, nil, c.Service, c.Region, time.Unix(0, 0))
		if err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		actual := req.Header.Get("Authorization")
		if e, a := c.ExpSig, actual; e != a {
			t.Errorf("expected %v, but recieved %v", e, a)
		}
		if e, a := c.OrigURI, req.URL.Path; e != a {
			t.Errorf("expected %v, but recieved %v", e, a)
		}
		if e, a := c.EscapedURI, req.URL.EscapedPath(); e != a {
			t.Errorf("expected %v, but recieved %v", e, a)
		}
	}
}

func TestStandaloneSign_RawPath(t *testing.T) {
	creds := unit.Config().Credentials
	signer := v4.NewSigner(creds)

	for _, c := range standaloneSignCases {
		host := fmt.Sprintf("https://%s.%s.%s.amazonaws.com",
			c.SubDomain, c.Region, c.Service)

		req, err := http.NewRequest("GET", host, nil)
		if err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		// URL.EscapedPath() will be used by the signer to get the
		// escaped form of the request's URI path.
		req.URL.Path = c.OrigURI
		req.URL.RawPath = c.EscapedURI
		req.URL.RawQuery = c.OrigQuery

		_, err = signer.Sign(context.Background(), req, nil, c.Service, c.Region, time.Unix(0, 0))
		if err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		actual := req.Header.Get("Authorization")
		if e, a := c.ExpSig, actual; e != a {
			t.Errorf("expected %v, but recieved %v", e, a)
		}
		if e, a := c.OrigURI, req.URL.Path; e != a {
			t.Errorf("expected %v, but recieved %v", e, a)
		}
		if e, a := c.EscapedURI, req.URL.EscapedPath(); e != a {
			t.Errorf("expected %v, but recieved %v", e, a)
		}
	}
}
