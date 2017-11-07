// +build go1.8

package aws_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestResetBody_WithEmptyBody(t *testing.T) {
	r := aws.Request{
		HTTPRequest: &http.Request{},
	}

	reader := strings.NewReader("")
	r.Body = reader

	r.ResetBody()

	if a, e := r.HTTPRequest.Body, http.NoBody; a != e {
		t.Errorf("expected request body to be set to reader, got %#v",
			r.HTTPRequest.Body)
	}
}

func TestRequest_FollowPUTRedirects(t *testing.T) {
	const bodySize = 1024

	redirectHit := 0
	endpointHit := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/redirect-me":
			u := *r.URL
			u.Path = "/endpoint"
			w.Header().Set("Location", u.String())
			w.WriteHeader(307)
			redirectHit++
		case "/endpoint":
			b := bytes.Buffer{}
			io.Copy(&b, r.Body)
			r.Body.Close()
			if e, a := bodySize, b.Len(); e != a {
				t.Fatalf("expect %d body size, got %d", e, a)
			}
			endpointHit++
		default:
			t.Fatalf("unexpected endpoint used, %q", r.URL.String())
		}
	}))

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	svc := awstesting.NewClient(cfg)

	req := svc.NewRequest(&aws.Operation{
		Name:       "Operation",
		HTTPMethod: "PUT",
		HTTPPath:   "/redirect-me",
	}, &struct{}{}, &struct{}{})
	req.SetReaderBody(bytes.NewReader(make([]byte, bodySize)))

	err := req.Send()
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := 1, redirectHit; e != a {
		t.Errorf("expect %d redirect hits, got %d", e, a)
	}
	if e, a := 1, endpointHit; e != a {
		t.Errorf("expect %d endpoint hits, got %d", e, a)
	}
}
