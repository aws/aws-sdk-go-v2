package aws_test

import (
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

func TestBuildableHTTPClient_WithTimeout(t *testing.T) {
	client := &aws.BuildableHTTPClient{}

	expect := 10 * time.Millisecond
	client2 := client.WithTimeout(expect).(*aws.BuildableHTTPClient)

	if e, a := time.Duration(0), client.GetTimeout(); e != a {
		t.Errorf("expect %v initial timeout, got %v", e, a)
	}

	if e, a := expect, client2.GetTimeout(); e != a {
		t.Errorf("expect %v timeout, got %v", e, a)
	}
}
