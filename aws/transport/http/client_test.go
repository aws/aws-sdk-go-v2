package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestBuildableClient_NoFollowRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Moved Permanently", http.StatusMovedPermanently)
		}))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)

	client := NewBuildableClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := http.StatusMovedPermanently, resp.StatusCode; e != a {
		t.Errorf("expect %v code, got %v", e, a)
	}
}

func TestBuildableClient_WithTimeout(t *testing.T) {
	client := &BuildableClient{}

	expect := 10 * time.Millisecond
	client2 := client.WithTimeout(expect)

	if e, a := time.Duration(0), client.GetTimeout(); e != a {
		t.Errorf("expect %v initial timeout, got %v", e, a)
	}

	if e, a := expect, client2.GetTimeout(); e != a {
		t.Errorf("expect %v timeout, got %v", e, a)
	}
}

func TestBuildableClient_concurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
	defer server.Close()

	var client aws.HTTPClient = NewBuildableClient()

	atOnce := 100
	var wg sync.WaitGroup
	wg.Add(atOnce)
	for i := 0; i < atOnce; i++ {
		go func(i int, client aws.HTTPClient) {
			defer wg.Done()

			if v, ok := client.(interface{ GetTimeout() time.Duration }); ok {
				v.GetTimeout()
			}

			if i%3 == 0 {
				if v, ok := client.(interface {
					WithTransportOptions(opts ...func(*http.Transport)) aws.HTTPClient
				}); ok {
					client = v.WithTransportOptions()
				}
			}

			req, _ := http.NewRequest("GET", server.URL, nil)
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("expect no error, got %v", err)
			}
			resp.Body.Close()
		}(i, client)
	}

	wg.Wait()
}

func TestBuildableClient_RemovesSecurityTokenOnHostChange(t *testing.T) {
	// Create a test server that returns a 307 redirect to a different host
	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request has the security token header
		if r.Header.Get("X-Amz-Security-Token") != "" {
			t.Errorf("expected X-Amz-Security-Token to be removed on redirect to different host, but it was present")
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer redirectServer.Close()

	// Create the initial server that returns a redirect
	initialServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectServer.URL, http.StatusTemporaryRedirect)
	}))
	defer initialServer.Close()

	client := NewBuildableClient()

	// Create a request with Authorization and X-Amz-Security-Token headers
	req, _ := http.NewRequest("GET", initialServer.URL, nil)
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=AKID/20210101/us-east-1/s3/aws4_request")
	req.Header.Set("X-Amz-Security-Token", "token123")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expect 200 code, got %d", resp.StatusCode)
	}
}

func TestBuildableClient_KeepsSecurityTokenOnSameHost(t *testing.T) {
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const redirectedPath = "/redirected"
		if !strings.HasSuffix(r.URL.Path, redirectedPath) {
			// First request - redirect to different path on same host
			redirectURL := serverURL + redirectedPath
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}
		if r.Header.Get("X-Amz-Security-Token") == "" {
			t.Errorf("expected X-Amz-Security-Token to be set, but it was present")
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	serverURL = server.URL
	defer server.Close()

	client := NewBuildableClient()
	// Create a request with X-Amz-Security-Token header
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("X-Amz-Security-Token", "token123")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expect 200 code, got %d", resp.StatusCode)
	}
}
