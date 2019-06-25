package aws_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestBuildableHTTPClient_NoFollowRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Moved Permanently", http.StatusMovedPermanently)
		}))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)

	client := aws.NewBuildableHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := http.StatusMovedPermanently, resp.StatusCode; e != a {
		t.Errorf("expect %v code, got %v", e, a)
	}
}

func ExampleBuildableHTTPClient_WithTransportOptions() {
	// Create a new client by calling the constructor.
	// Add custom *http.Transport configuration to the client. The
	// modifications will be on the new client returned.
	client := aws.NewBuildableHTTPClient().
		WithTransportOptions(func(tr *http.Transport) {
			// Only wait 10 seconds for the full response headers to be read.
			tr.ResponseHeaderTimeout = 10 * time.Second
		})

	// Set the configured HTTP client to the SDK's Config, and use the
	// Config to create API clients with.
	cfg.HTTPClient = client
}

func ExampleBuildableHTTPClient_WithDialerOptions() {
	// Create a new client by calling the constructor.
	client := aws.NewBuildableHTTPClient().
		WithDialerOptions(func(d *net.Dialer) {
			// Set the network (e.g. TCP) dial timeout to 10 seconds.
			d.Timeout = 10 * time.Second
		})

	// Set the configured HTTP client to the SDK's Config, and use the
	// Config to create API clients with.
	cfg.HTTPClient = client
}

var cfg aws.Config
